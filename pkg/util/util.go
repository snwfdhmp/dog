package util

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/sirupsen/logrus"
	"github.com/snwfdhmp/dog/pkg/config"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

func init() {
	log.SetLevel(logrus.ErrorLevel)
}

var (
	fs  = afero.NewOsFs()
	log = logrus.New()

	buildContentTemplate = func() *template.Template {
		return template.New("content").Option("missingkey=error").Funcs(template.FuncMap{
			"unexported": func(packageName string) string {
				return strings.ToLower(packageName)
			},
		})
	}

	buildPathTemplate = func() *template.Template {
		return template.New("path")
	}
)

func TemplateLocation(templateName string) (path string, err error) {
	path = filepath.Join(config.TemplatesLocation, templateName)

	exists, err := afero.Exists(fs, path)
	if err != nil {
		return
	} else if !exists {
		err = fmt.Errorf("%s: no such directory", path)
		return
	}

	return
}

func TemplateFile(templateName string) (content []byte, err error) {
	content, err = afero.ReadFile(fs, filepath.Join(config.TemplatesLocation, templateName+".yaml"))
	return
}

func TemplateData(cmd *cobra.Command, templateName string, args []string) (data map[string]*string, err error) {
	data = make(map[string]*string)
	dataYaml := make(map[string]interface{})

	templateFile, err := TemplateFile(templateName)
	if err != nil {
		return
	}

	if err = yaml.Unmarshal(templateFile, &dataYaml); err != nil {
		return
	}

	for name, value := range dataYaml["vars"].(map[interface{}]interface{}) {
		data[name.(string)] = cmd.Flags().StringP(name.(string), "", value.(string), "usage")
	}

	if err = cmd.Flags().Parse(args); err != nil {
		return
	}

	return
}

type File struct {
	Path    string
	Content []byte
	Perm    os.FileMode
}

func (f *File) WriteToFs() error {
	return afero.WriteFile(fs, f.Path, f.Content, f.Perm)
}

func TranslateFile(templatePath, src, dst string, data interface{}) (*File, error) {
	info, err := fs.Stat(templatePath) // get file informations
	if err != nil {
		return nil, err
	}

	log.Infof("Processing %s", templatePath)

	newPath, err := TranslatePath(templatePath, src, dst, data)
	if err != nil {
		return nil, err
	}

	//declare buffers
	contentBuf := bytes.NewBuffer([]byte{})

	//read template content
	b, err := afero.ReadFile(fs, templatePath)
	if err != nil {
		return nil, err
	}

	//initialize content template
	contentTmpl, err := buildContentTemplate().Parse(string(b))
	if err != nil {
		return nil, err
	}

	//execute content template
	if err := contentTmpl.Execute(contentBuf, data); err != nil {
		if strings.Contains(err.Error(), "empty template") {
			log.Warnf("%s is an empty template", templatePath)
			return &File{
				Path:    newPath,
				Content: b,
				Perm:    info.Mode().Perm(),
			}, nil
		}
		return nil, err
	}

	file := &File{
		Path:    newPath,
		Content: contentBuf.Bytes(),
		Perm:    info.Mode().Perm(),
	}

	return file, file.WriteToFs()
}

func TranslatePath(path, src, dst string, data interface{}) (string, error) {
	path = strings.Replace(path, src, dst, -1)       //switch root directory
	pathTmpl, err := buildPathTemplate().Parse(path) //parse path template
	if err != nil {
		return "", err
	}

	buf := bytes.NewBuffer([]byte{}) //prepare buffer

	if err := pathTmpl.Execute(buf, data); err != nil { //execute path template
		return "", err
	}

	return buf.String(), nil //return result
}

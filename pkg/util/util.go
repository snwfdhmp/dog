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
	log.SetLevel(logrus.DebugLevel)
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

type Builder struct {
	cmd    *cobra.Command
	args   []string
	config map[string]interface{}
	path   string
	target string
	name   string
	vars   map[string]*string
}

func NewBuilder(name string) (builder *Builder, err error) {
	builder = &Builder{
		args:   make([]string, 0),
		config: make(map[string]interface{}),
		path:   "",
		target: "",
		name:   "",
		vars:   make(map[string]*string),
	}
	builder.name = name

	builder.path, err = TemplateLocation(name)
	if err != nil {
		return
	}

	builder.config, err = builder.GetConfig()
	if err != nil {
		return
	}

	return
}

func (b *Builder) Configure(cmd *cobra.Command, args []string) error {
	b.cmd = cmd
	b.args = args

	_, err := b.GetVars()
	return err
}

func (b *Builder) Run(action, target string) error {
	b.target = target
	config, err := b.GetConfig()
	if err != nil {
		return err
	}

	for _, file := range config[action].(map[interface{}]interface{})["files"].([]interface{}) {
		if err := b.Translate(file.(string)); err != nil {
			return err
		}
	}
	return nil
}

func (b *Builder) Translate(path string) error {
	log.Infof("Translating %s...", path)
	path = filepath.Join(b.GetPath(), path)
	info, err := fs.Stat(path)
	if err != nil {
		return err
	}

	if info.IsDir() {
		translatedPath, err := b.TranslatePath(path)
		if err != nil {
			return err
		}
		if err := fs.Mkdir(translatedPath, info.Mode().Perm()); err != nil {
			return err
		}
	} else {
		_, err := b.TranslateFile(path)
		if err != nil {
			return err
		}
	}
	return nil
}

func (b *Builder) GetName() string {
	return b.name
}

func (b *Builder) GetPath() string {
	return b.path
}

func (b *Builder) GetTarget() string {
	return b.target
}

func (b *Builder) GetConfig() (map[string]interface{}, error) {
	if len(b.config) > 0 {
		return b.config, nil
	}

	templateFile, err := TemplateFile(b.GetName())
	if err != nil {
		return nil, err
	}

	if err = yaml.Unmarshal(templateFile, &b.config); err != nil {
		return nil, err
	}

	return b.config, nil
}

func (b *Builder) GetArgs() []string {
	return b.args
}

func (b *Builder) GetVars() (map[string]*string, error) {
	if len(b.vars) > 0 {
		return b.vars, nil
	}

	config, err := b.GetConfig()
	if err != nil {
		return nil, err
	}

	for name, value := range config["vars"].(map[interface{}]interface{}) {
		b.vars[name.(string)] = b.cmd.Flags().StringP(name.(string), "", value.(string), "usage")
	}

	if err = b.cmd.Flags().Parse(b.GetArgs()); err != nil {
		return nil, err
	}

	return b.vars, nil
}

type File struct {
	Path    string
	Content []byte
	Perm    os.FileMode
}

func (f *File) WriteToFs() error {
	return afero.WriteFile(fs, f.Path, f.Content, f.Perm)
}

func (b *Builder) TranslateFile(filePath string) (*File, error) {
	info, err := fs.Stat(filePath) // get file informations
	if err != nil {
		return nil, err
	}

	newPath, err := b.TranslatePath(filePath)
	if err != nil {
		return nil, err
	}

	log.Debugf("Translating %s to %s", filePath, newPath)

	targetDir := filepath.Dir(newPath)
	exists, err := afero.Exists(fs, targetDir)
	if err != nil {
		return nil, err
	}
	if !exists {
		dirInfo, err := fs.Stat(filepath.Dir(filePath))
		if err != nil {
			return nil, err
		}
		if err := fs.MkdirAll(targetDir, dirInfo.Mode().Perm()); err != nil {
			return nil, err
		}
	}

	//declare buffers
	contentBuf := bytes.NewBuffer([]byte{})

	//read template content
	content, err := afero.ReadFile(fs, filePath)
	if err != nil {
		return nil, err
	}

	//initialize content template
	contentTmpl, err := buildContentTemplate().Parse(string(content))
	if err != nil {
		return nil, err
	}

	data, err := b.GetVars()
	if err != nil {
		return nil, err
	}
	//execute content template
	if err := contentTmpl.Execute(contentBuf, data); err != nil {
		if strings.Contains(err.Error(), "empty template") {
			log.Warnf("%s is an empty template", filePath)
			return &File{
				Path:    newPath,
				Content: content,
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

func (b *Builder) TranslatePath(path string) (string, error) {
	path = strings.Replace(path, b.GetPath(), b.GetTarget(), -1) //switch root directory
	pathTmpl, err := buildPathTemplate().Parse(path)             //parse path template
	if err != nil {
		return "", err
	}

	buf := bytes.NewBuffer([]byte{}) //prepare buffer

	data, err := b.GetVars()
	if err != nil {
		return "", err
	}
	if err := pathTmpl.Execute(buf, data); err != nil { //execute path template
		return "", err
	}

	return buf.String(), nil //return result
}

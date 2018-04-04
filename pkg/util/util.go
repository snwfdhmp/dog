package util

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
)

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

func TemplateLocation(templateName string) (string, error) {
	path := filepath.Join("/Users/snwfdhmp/go/src/github.com/snwfdhmp/dog/templates", templateName)

	exists, err := afero.Exists(fs, path)
	if err != nil {
		return "", err
	} else if !exists {
		return "", fmt.Errorf("%s: no such directory", path)
	}

	return path, nil
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

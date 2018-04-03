package util

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/afero"
)

var (
	fs = afero.NewOsFs()
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

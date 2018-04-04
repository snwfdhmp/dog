// Copyright © 2018 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"strings"
	"text/template"

	"github.com/snwfdhmp/dog/pkg/util"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

var (
	fs = afero.NewOsFs()

	buildUsage = `Usage: build <templateName> <action> [--name=value]`
)

// buildCmd represents the build command
var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "A brief description of your command",
	Long:  buildUsage,
	Run: func(cmd *cobra.Command, args []string) {
		if err := buildFunc(cmd, args); err != nil {
			fmt.Println("fatal:", err)
			return
		}
	},
}

func buildFunc(cmd *cobra.Command, args []string) error {
	if len(args) < 2 {
		return errors.New(buildUsage)
	}

	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	fmt.Println("WD:", wd)

	templateName := args[0]
	// action := args[1]

	location, err := util.TemplateLocation(templateName)
	if err != nil {
		return err
	}

	data := struct {
		ResourceName string
		PackageName  string
		PackagePath  string
	}{
		"patient",
		"testPackage",
		"github.com/snwfdhmp/testPackage",
	}

	if err := afero.Walk(fs, location, func(path string, info os.FileInfo, err error) error {
		realPath := path
		if path == location {
			return nil
		}
		path = strings.Replace(path, location, wd, -1)
		tmpl, err := template.New("name").Parse(path)
		if err != nil {
			return err
		}

		buf := bytes.NewBuffer([]byte{})
		if err := tmpl.Execute(buf, data); err != nil {
			return err
		}

		path = buf.String()
		// fmt.Println(path)
		if info.IsDir() {
			if err := fs.Mkdir(path, info.Mode().Perm()); err != nil {
				return err
			}
		} else {
			b, err := afero.ReadFile(fs, realPath)
			if err != nil {
				return err
			}

			contentTmpl, err := template.New("content").Funcs(template.FuncMap{
				"unexported": func(packageName string) string {
					return strings.ToLower(packageName)
				},
			}).Parse(string(b))
			if err != nil {
				return err
			}

			content := bytes.NewBuffer([]byte{})

			contentTmpl.Execute(content, data)

			if err := afero.WriteFile(fs, path, content.Bytes(), info.Mode().Perm()); err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		return err
	}

	return nil
}

func init() {
	RootCmd.AddCommand(buildCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// buildCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// buildCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

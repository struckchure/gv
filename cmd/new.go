package main

import (
	"fmt"
	"io/fs"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/fatih/color"
	"github.com/go-resty/resty/v2"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
	"github.com/struckchure/gv"
)

const (
	TEMPLATE_REMOTE_BASE = "https://api.github.com/repos/struckchure/gv/contents/"
	TEMPLATE_ROOT        = "./_templates"
	TEMPLATE_IS_REMOTE   = true
)

type newAnswers struct {
	ProjectName string
	PackageName string
	Template    string
}

var newCommand = &cobra.Command{
	Use:   "new",
	Short: "Setup a new project from a list of templates",
	Run: func(cmd *cobra.Command, args []string) {
		answers := newAnswers{}
		templateOptions := []string{
			"Blank",
			"HTML File Router",
			"React",
			"Vue",
			"Svelte",
		}

		var qs = []*survey.Question{
			{
				Name:      "ProjectName",
				Prompt:    &survey.Input{Message: "Project Name:"},
				Validate:  survey.Required,
				Transform: survey.ToLower,
			},
			{
				Name:      "PackageName",
				Prompt:    &survey.Input{Message: "Package Name:"},
				Validate:  survey.Required,
				Transform: survey.ToLower,
			},
			{
				Name:     "Template",
				Prompt:   &survey.Select{Message: "Choose a template:", Options: templateOptions},
				Validate: survey.Required,
			},
		}
		err := survey.Ask(qs, &answers)
		if err != nil {
			color.Red(err.Error())
			return
		}

		newService(answers)
	},
}

type gitHubFileInfo struct {
	Name        string      `json:"name"`
	Path        string      `json:"path"`
	Sha         string      `json:"sha"`
	Size        int         `json:"size"`
	URL         string      `json:"url"`
	HTMLURL     string      `json:"html_url"`
	GitURL      string      `json:"git_url"`
	DownloadURL string      `json:"download_url"`
	Type        string      `json:"type"`
	Links       gitHubLinks `json:"_links"`
}

type gitHubLinks struct {
	Self string `json:"self"`
	Git  string `json:"git"`
	HTML string `json:"html"`
}

// For parsing the array wrapper
type gitHubFileResponse []gitHubFileInfo

func downloadFromDir(client *resty.Client, dir, outputDir string) error {
	files := &gitHubFileResponse{}
	res, err := client.R().SetResult(files).Get(dir)
	if err != nil {
		return err
	}

	if res.StatusCode() != 200 {
		return fmt.Errorf("%d %s", res.StatusCode(), res.String())
	}

	for _, file := range *files {
		path, err := filepath.Rel(TEMPLATE_ROOT, file.Path)
		if err != nil {
			return err
		}

		if file.Type == "dir" {
			err = downloadFromDir(client, path, outputDir)
			if err != nil {
				return err
			}
		}

		content, err := client.R().Get(file.DownloadURL)
		if err != nil {
			return err
		}

		path = strings.Join(strings.Split(path, "/")[1:], "/")

		dirPath := filepath.Dir(filepath.Join(outputDir, path))
		err = os.MkdirAll(dirPath, os.ModePerm)
		if err != nil {
			return err
		}

		err = os.WriteFile(filepath.Join(outputDir, path), content.Body(), os.ModePerm)
		if err != nil {
			info, err := os.Stat(filepath.Join(outputDir, path))
			if err == nil && !info.IsDir() {
				return err
			}
		}
	}

	return nil
}

func remoteTemplate(projectRoot, template string) error {
	client := resty.New().
		SetBaseURL(lo.Must(url.JoinPath(TEMPLATE_REMOTE_BASE, TEMPLATE_ROOT))).
		SetQueryParam("ref", "gh/issue-11")

	return downloadFromDir(client, template, projectRoot)
}

func localTemplate(projectRoot, template string) error {
	_, err := os.Stat(lo.Must(url.JoinPath(
		lo.Must(os.Getwd()),
		TEMPLATE_ROOT,
		template,
	)))
	if err != nil {
		return err
	}

	templateDir := filepath.Join(lo.Must(os.Getwd()), TEMPLATE_ROOT, template)
	_, err = os.Stat(templateDir)
	if err != nil {
		return err
	}

	return filepath.Walk(
		templateDir,
		func(path string, info fs.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() {
				return nil
			}

			outputPath, err := filepath.Rel(templateDir, path)
			if err != nil {
				return err
			}

			err = gv.CopyFile(path, filepath.Join(projectRoot, outputPath))
			if err != nil {
				return err
			}

			return nil
		},
	)
}

func newService(answers newAnswers) {
	answers.ProjectName = strings.ToLower(answers.ProjectName)
	answers.ProjectName = strings.ReplaceAll(answers.ProjectName, " ", "-")

	answers.Template = strings.ToLower(answers.Template)
	answers.Template = strings.ReplaceAll(answers.Template, " ", "-")

	err := os.Mkdir(answers.ProjectName, os.ModePerm)
	if err != nil {
		color.Red(err.Error())
		return
	}

	projectRoot := filepath.Join(lo.Must(os.Getwd()), answers.ProjectName)

	if TEMPLATE_IS_REMOTE {
		err = remoteTemplate(projectRoot, answers.Template)
		if err != nil {
			color.Red(err.Error())
			return
		}

		return
	}

	err = localTemplate(projectRoot, answers.Template)
	if err != nil {
		color.Red(err.Error())
		return
	}

	err = os.Chdir(projectRoot)
	if err != nil {
		color.Red(err.Error())
		return
	}

	color.Blue("$ go mod init %s", answers.PackageName)
	err = gv.ExecCommandWithCallback("go", []string{"mod", "init", answers.PackageName}, func(output string) {
		color.Green(output)
	})
	if err != nil {
		color.Red(err.Error())
		return
	}

	color.Blue("$ go mod tidy")
	err = gv.ExecCommandWithCallback("go", []string{"mod", "tidy"}, func(output string) {
		color.Green(output)
	})
	if err != nil {
		color.Red(err.Error())
		return
	}

	color.Cyan(`
cd %s
go run .
	`, answers.ProjectName)
}

func init() {
	rootCmd.AddCommand(newCommand)
}

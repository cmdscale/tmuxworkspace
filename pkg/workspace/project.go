package workspace

import (
	"errors"
	"fmt"
	"os"
	"text/template"

	"github.com/spf13/cobra"
)

const (
	TmuxProjectPathConfig = "tmuxspace"
)

var (
	errOpenProjectYamlFile  = errors.New("Could not open yaml project file")
	errUnMarshalProjectFile = errors.New("Could unmarshal yaml project file")
)

/*
  A project provides automation logic and session managed for software development
  projects. It provides a project template functionality corresponding to tmux
  layouts, but with persistence and automations such as executing recurring commands
  on session / project up and down. E.g. docker-compose up, docker-compose down,
  docker system prune, open your prefered editor (neovim) and so on.
  A startDirectory will be set to all windows and panes and passed down by the project
  if not explicitely set for windows and panes. The same is applied to panes if not set,
  but windows have set a startDirectory.
*/

// Project as a one to one association with tmux session.
type Project struct {
	Name           string   `yaml:"name"`
	StartDirectory string   `yaml:"startDirectory"`
	Windows        []Window `yaml:"windows"`
}

// Window as a one to one association with tmux window.
type Window struct {
	Name           string `yaml:"name"`
	StartDirectory string `yaml:"startDirectory"`
	Panes          []Pane `yaml:"panes"`
}

// Pane as a one to one association with tmux pane.
type Pane struct {
	Name           string `yaml:"name"`
	StartDirectory string `yaml:"startDirectory"`
	Ratio          string `yaml:"ratio"`
}

func NewCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "new-project",
		Short: "create a new tmux project",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 2 {
				return errors.New("requires two positional arguments 0: project name 1: start directory")
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			NewProject(args[0], args[1])
		},
	}
}

// NewProject will create a defaultProjectTemplate.
func NewProject(name string, startDirectory string) (Project, error) {
	p := Project{
		Name:           name,
		StartDirectory: startDirectory,
		Windows: []Window{
			{
				Name:           fmt.Sprintf("%s:%d", name, 0),
				StartDirectory: startDirectory,
				Panes: []Pane{
					{
						Name:           fmt.Sprintf("%d:%d", 0, 0),
						StartDirectory: startDirectory,
					},
				},
			},
		},
	}

	tmpl, err := template.New("project").Parse(defaultProjectTemplate)
	if err != nil {
		return p, err
	}

	file, err := openProjectFile(name)
	if err != nil {
		return p, err
	}
	defer file.Close()

	err = tmpl.Execute(file, p)
	if err != nil {
		return p, err
	}
	return p, nil
}

func openProjectFile(name string) (*os.File, error) {
	var (
		path string
		file *os.File
	)

	path = os.Getenv("XDG_CONFIG_HOME")
	if path == "" {
		path = os.Getenv("HOME")
	}
	path = fmt.Sprintf("%s/%s", path, TmuxProjectPathConfig)
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(path, os.ModePerm)
		if err != nil {
			return nil, errors.New("Could not create tmux-project directory")
		}
	}
	filePath := fmt.Sprintf("%s/%s.yaml", path, name)

	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0660)
	if err != nil {
		return nil, err
	}

	fileInfo, err := os.Stat(filePath)
	mode := fileInfo.Mode()
	if !mode.IsRegular() {
		return nil, errors.New("not a regular file")
	}

	return file, err
}

// defaultProjectTemplate will create a default project template which can either be adapted
// by passing additional cli args or directly editing the generated yam config.
var defaultProjectTemplate = `
name: {{ .Name }}
startDirectory: {{ .StartDirectory }}
windows: {{ range .Windows }}
- name: {{ .Name }}
  startDirectory: {{ .StartDirectory}}
  panes: {{ range .Panes }}
  - {{ .StartDirectory }}
  {{ end }}
{{ end }}
`

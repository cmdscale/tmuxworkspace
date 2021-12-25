package tmux

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"path/filepath"

	"github.com/cmdscale/tmux-project/pkg/workspace"
	"github.com/magefile/mage/sh"
	"github.com/spf13/cobra"
)

func NewCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "new-session",
		Short: "create a new tmux session",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("requires at least session name")
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			var path string
			exists, err := checkIfProjectExists(args[0])
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
			if exists {
				fmt.Fprintln(os.Stderr, errors.New(fmt.Sprintf("Project %s already exists", args[0])))
				os.Exit(1)
			}

			// no path was given
			if len(args) == 1 {
				path, err = os.Getwd()
				if err != nil {
					fmt.Fprintln(os.Stderr, err)
					os.Exit(1)
				}
			}

			NewSession(args[0], path)
		},
	}
}

func NewSession(name string, path string) error {
	p, err := workspace.NewProject(name, path)
	if err != nil {
		return err
	}
	err = applyProjectLayoutToSession(p)
	if err != nil {
		return err
	}
	return nil
}

func checkIfProjectExists(name string) (bool, error) {
	var (
		projectExists bool
	)
	path := os.Getenv("XDG_CONFIG_HOME")
	if path == "" {
		path = os.Getenv("HOME")
	}
	path = fmt.Sprintf("%s/%s", path, workspace.TmuxProjectPathConfig)
	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			fileName := strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
			if fileName == name {
				projectExists = true
			}
		}
		return nil
	})
	if err != nil {
		return projectExists, errors.New("Could not walk through tmux-projects")
	}
	return projectExists, nil
}

func applyProjectLayoutToSession(project workspace.Project) error {
	// var paneRatio string
	// -d: detached, -s: session name, -n: window name, -c: start directory
	tmuxNewSession := sh.RunCmd("tmux", "new-session", "-Ad", "-s", project.Name, "-n", project.Windows[0].Name, "-c", project.StartDirectory)
	err := tmuxNewSession()
	if sh.ExitStatus(err) != 0 {
		return errors.New(fmt.Sprintf("Could not create session, %s", project.Name))
	}
	for ii, window := range project.Windows {
		// only create a new window when it's not first as already created in new session command via -n option.
		if ii > 0 {
			// -t: target window, -a: inserted after target window (index), -n: window name, -c: start directory
			tmuxNewWindow := sh.RunCmd("tmux", "new-window", "-a", "-t", project.Name, "-n", window.Name, "-c", window.StartDirectory)
			err := tmuxNewWindow()
			if sh.ExitStatus(err) != 0 {
				return errors.New(fmt.Sprintf("Could not create window, %s", project.Name))
			}
		}
		// for _, pane := range window.Panes {
		// 	// select the window before creating panes.
		// 	tmuxSelectTargetWindow := sh.RunCmd("tmux", "select-window", "-t", window.Name)
		// 	err := tmuxSelectTargetWindow()
		// 	if sh.ExitStatus(err) != 0 {
		// 		return errors.New(fmt.Sprintf("Could not select window, %s", project.Name))
		// 	}
		// 	// -v: vertical split, -h: horizontal split, -t: target pane, -c: start directory
		// 	if pane.Ratio == "" {
		// 		paneRatio = "30"
		// 	}
		// 	log.Println(window.Name)
		// 	tmuxNewPane := sh.RunCmd("tmux", "split-window", "-v", "-t", window.Name, "-c", pane.StartDirectory, "-p", paneRatio)
		// 	err = tmuxNewPane()
		// 	if sh.ExitStatus(err) != 0 {
		// 		return errors.New(fmt.Sprintf("Could not create pane, %s", project.Name))
		// 	}
		// }
	}
	return nil
}

func hasSession(name string) (bool, error) {
	tmuxHasSession := sh.OutCmd("tmux", "has-session", "-t")
	_, err := tmuxHasSession(name)
	if sh.ExitStatus(err) != 0 {
		return false, errors.New(fmt.Sprintf("Session exists, %s", name))
	}
	return true, nil
}

func attachToSession(name string) (bool, error) {
	tmuxAttachSession := sh.OutCmd("tmux", "attach", "-t")
	_, err := tmuxAttachSession(name)
	if sh.ExitStatus(err) != 0 {
		return false, errors.New(fmt.Sprintf("Session could not be attached, %s", name))
	}
	return true, nil
}

func switchSession(name string) (bool, error) {
	tmuxSwitchSession := sh.OutCmd("tmux", "switch-client", "-t")
	_, err := tmuxSwitchSession(name)
	if sh.ExitStatus(err) != 0 {
		return false, errors.New(fmt.Sprintf("Session could not be switched, %s", name))
	}
	return true, nil
}

package main

import (
	"log"

	"github.com/cmdscale/tmux-project/cmd"
	"github.com/cmdscale/tmux-project/pkg/tmux"
	"github.com/cmdscale/tmux-project/pkg/workspace"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	rootCmd := cmd.NewTmuxProjectCLI()
	workspaceCmd := workspace.NewCmd()
	tmuxCmd := tmux.NewCmd()
	rootCmd.AddCommand(workspaceCmd)
	rootCmd.AddCommand(tmuxCmd)
	rootCmd.Execute()
}

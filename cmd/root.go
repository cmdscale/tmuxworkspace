package cmd

import (
	"fmt"
	"log"
	"os"

	// "github.com/cmdscale/tmux-project/pkg/project"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile string
)

type TmuxProjectCLI struct {
	*cobra.Command
}

func NewTmuxProjectCLI() *TmuxProjectCLI {
	cobra.OnInitialize(initConfig)
	t := &TmuxProjectCLI{
		&cobra.Command{
			Use:   "tmuxctl",
			Short: "tmuxctl helps you manage your tmux projects",
			Run: func(cmd *cobra.Command, args []string) {
			},
		},
	}
	t.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.config/tmux-project.yaml)")
	return t
}

func (t *TmuxProjectCLI) Execute() {
	if err := t.Command.Execute(); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}

func (t *TmuxProjectCLI) AddCommand(cmd *cobra.Command) {
	t.Command.AddCommand(cmd)
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".cobra")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

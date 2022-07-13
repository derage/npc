/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"os"

	"github.com/derage/npc/pkg/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var logger = utils.GetLogger()
var version, commit, date, builtBy string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "npc",
	Short: "Templating tool to help with your every day needs",
	Long: `NPC: who you go to when you want to continue your journey
in your quest for devops.

This tool will use template files from the given repo
that is configured in your yaml file to bootstrap the 
directory you give it. With the help of various template
functions, this tool can be used to generate files that
would otherwise need to be copied from various locations`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		if viper.GetBool("version") {
			fmt.Printf("npc %s, commit %s, built at %s by %s", version, commit, date, builtBy)
		} else {
			cmd.Help()
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute(bVersion string, bCommit string, bDate string, bBuiltBy string) {
	version = bVersion
	commit = bCommit
	date = bDate
	builtBy = bBuiltBy
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	home, err := os.UserHomeDir()
	cobra.CheckErr(err)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.npc.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.PersistentFlags().StringP("template", "t", "", "name of the template")
	rootCmd.PersistentFlags().StringP("template-path", "p", home+"/.npc", "Where do you want your templates to be cached?")
	rootCmd.PersistentFlags().StringP("directory", "d", "./", "Where do you want to bootstrap the files to?")
	rootCmd.PersistentFlags().BoolP("version", "v", false, "Print build version")
	viper.BindPFlags(rootCmd.PersistentFlags())
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".npc" (without extension).
		viper.AddConfigPath(home)
		viper.AddConfigPath("./")
		viper.SetConfigType("yaml")
		viper.SetConfigName(".npc")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}

}

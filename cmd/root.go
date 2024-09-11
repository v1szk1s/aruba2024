/*
Copyright © 2024 Attila Ambrus, Kvak Barnabás, Litavecz Máté, Sárközi Gergő, Ördög Csaba  attila.ambrus0022@gmail.com

*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var(
    cfgFile string

    key string
    secret string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
    Use:   "ruba",
    Short: "Cli tool for communicating with uba Cloud api",
    Long: `ruba is a ... longer desc`,
    Version: "0.1",
    // Uncomment the following line if your bare application
    // has an action associated with it:
    // Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

    rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file (default: $HOME/.config/ruba/ruba.yaml)")

    rootCmd.PersistentFlags().StringVar(&key, "key", "", "Aruba api key (required)")
    rootCmd.PersistentFlags().StringVar(&secret, "secret", "", "Aruba api secret (required)")

    viper.BindPFlag("key", rootCmd.PersistentFlags().Lookup("key"))
    viper.BindPFlag("secret", rootCmd.PersistentFlags().Lookup("secret"))

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
        


		// Search config in .config directory with name "ruba" (without extension).
		viper.AddConfigPath(home + "/.config/ruba")
		viper.SetConfigType("yaml")
		viper.SetConfigName("ruba")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
    key = viper.GetString("key")
    secret = viper.GetString("secret")

	}
}

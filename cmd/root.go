/*
Copyright © 2024 Attila Ambrus, Litavecz Máté, Sárközi Gergő, Ördög Csaba  attila.ambrus0022@gmail.com
*/
package cmd

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile string

	key        string
	secret     string
	configPath string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "ruba",
	Short:   "Cli tool for communicating with uba Cloud api",
	Long:    `ruba is a ... longer desc`,
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

	//rootCmd

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
		configPath = filepath.Join(home, ".config", "ruba")

		err = os.MkdirAll(filepath.Dir(configPath), 0700)
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

func authenticate() {
	urlStr := "https://login.aruba.it/auth/realms/cmp-new-apikey/protocol/openid-connect/token"

	// Data to be sent in the request body
	data := url.Values{}
	data.Set("grant_type", "client_credentials")
	data.Set("client_id", key)        // Replace with your client_id
	data.Set("client_secret", secret) // Replace with your client_secret

	fmt.Fprintln(os.Stderr, "helo", key, secret)
	// Create a new POST request with the form data
	req, err := http.NewRequest("POST", urlStr, bytes.NewBufferString(data.Encode()))
	cobra.CheckErr(err)

	// Set the Content-Type header to application/x-www-form-urlencoded
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Send the request using http.DefaultClient
	client := &http.Client{}
	resp, err := client.Do(req)

	cobra.CheckErr(err)
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	cobra.CheckErr(err)

	// Print the response
	fmt.Println("Response Status:", resp.Status)
	fmt.Println("Response Body:", string(body))
}

func saveTokenToFile(token string) error {

	err := os.WriteFile(configPath, []byte(token), 0600)
	if err != nil {
		return err
	}

	return nil
}

func readTokenFromFile() string {

	token, err := os.ReadFile(configPath)
	cobra.CheckErr(err)

	return string(token)
}

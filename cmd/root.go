package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

var configFile string
var verbose bool

// TODO Change me
const cmdName = "storyline-api"

var (
	Version  string = ""
	Revision string = ""
	Time     string = ""
)

var rootCmd = &cobra.Command{
	Use:   cmdName,
	Short: "",
	Long:  ``,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "", fmt.Sprintf("config file (default \"./.%s.yml\")", cmdName))
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")

	// Set Default Viper configs.
	viper.SetDefault("build.version", Version)
	viper.SetDefault("build.revision", Revision)
	viper.SetDefault("build.time", Time)
	viper.SetDefault("http.handler", ":8080")
	viper.SetDefault("https.handler", ":8443")
	viper.SetDefault("https.keystore", "./store.key")
	viper.SetDefault("https.cert", "./store.crt")
	viper.SetDefault("mongo.uri", "localhost:27017")
	viper.SetDefault("mongo.database", "storyline_cqrs")
	viper.SetDefault("nats.uri", "nats://localhost:4222")
}

func initConfig() {
	if configFile != "" {
		viper.SetConfigFile(configFile)
	} else {
		filename := filepath.Join(".", fmt.Sprintf(".%s.yml", cmdName))
		if _, err := os.Stat(filename); os.IsNotExist(err) {
			// Default config file.
			configYml := `
http:
  handler: ":8080"
https:
  handler: ":8443"
  keystore: "./store.key"
  cert: "./store.crt"
mongo:
  uri: "localhost:27017"
  database: "storyline_cqrs"
nats:
  uri: "nats://localhost:4222"
`
			err = ioutil.WriteFile(filename, []byte(configYml), 0644)
			if err != nil {
				log.Fatal(err)
			}
		}

		viper.SetConfigName(fmt.Sprintf(".%s", cmdName))
		viper.AddConfigPath(".")
	}
	viper.SetEnvPrefix(cmdName)
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		log.Fatal(err)
	}
}

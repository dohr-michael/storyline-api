package cmd

import (
	"fmt"
	"github.com/dohr-michael/storyline-api/config"
	"github.com/dohr-michael/storyline-api/pkg"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
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
	RunE: func(cmd *cobra.Command, args []string) error {
		log.Println("storyline-api")

		// Configure base router
		router := chi.NewMux()
		// Middleware
		corsMiddleware := cors.New(cors.Options{
			AllowedOrigins:   []string{"*"},
			AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
			ExposedHeaders:   []string{"Link"},
			AllowCredentials: true,
			MaxAge:           300,
		})

		router.Use(
			middleware.DefaultCompress,
			middleware.DefaultLogger,
			corsMiddleware.Handler,
		)

		router.Get("/@/health", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"status": "OK"}`))
		})
		now := time.Now()
		router.Get("/@/info", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(
				fmt.Sprintf(
					`{
    "version": "%s",
    "revision": "%s",
    "build_at": "%s",
    "started_at": "%s"
}`,
					config.Config.BuildVersion(),
					config.Config.BuildRevision(),
					config.Config.BuildTime(),
					now.Format(time.RFC3339),
				),
			))
		})

		// Configure application.
		err := pkg.Start(router)
		if err != nil {
			return err
		}

		// Configure http / https handlers
		httpHandler := config.Config.HttpHandler()
		httpsHandler := config.Config.HttpsHandler()
		keystore := config.Config.HttpsKeystoreFile()
		cert := config.Config.HttpsCertFile()

		_, keystoreErr := os.Stat(keystore)
		_, certErr := os.Stat(cert)

		errors := make(chan error)

		go func() {
			log.Printf(" - Start http server on port %s\n", httpHandler)
			errors <- http.ListenAndServe(httpHandler, router)
		}()
		if keystoreErr == nil && certErr == nil {
			go func() {
				log.Printf(" - Start https server on port %s\n", httpsHandler)
				errors <- http.ListenAndServeTLS(httpsHandler, cert, keystore, router)
			}()
		}

		return <-errors
	},
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
	viper.SetDefault("nats.uri", "nats://localhost:4222")
	//viper.SetDefault("arango.endpoints", []string{"http://michael:azerty@localhost:8529"})
	//viper.SetDefault("arango.database", "storyline")
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
arango:
  endpoints:
  - "http://localhost:8529"
  database: "storyline"
  username: "michael"
  password: "azerty"
nats:
  uri: "nats://localhost:4222"
auth:
  iss: ""
  jwks: ""
  clientId: ""
  secret: ""
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

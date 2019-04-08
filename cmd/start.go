package cmd

import (
	"github.com/dohr-michael/storyline-api/config"
	"github.com/dohr-michael/storyline-api/pkg"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/spf13/cobra"
	"log"
	"net/http"
	"os"
)

var startCmd = &cobra.Command{
	Use: "start",
	RunE: func(cmd *cobra.Command, args []string) error {
		log.Println("storyline-api")

		// Configure base router
		router := chi.NewMux()
		router.Use(
			middleware.DefaultCompress,
			middleware.DefaultLogger,
		)

		router.Get("/@/health", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"status": "OK"}`))
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

func init() {
	rootCmd.AddCommand(startCmd)
}

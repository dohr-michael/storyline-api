package config

import (
	"github.com/arangodb/go-driver"
	arangohttp "github.com/arangodb/go-driver/http"
	"github.com/spf13/viper"
)

type arangoConfig struct {
	ClientConfig *driver.ClientConfig
	Database     string
}

var Config *config

func init() {
	Config = &config{}
}

type config struct{}

func (c *config) BuildVersion() string {
	return viper.GetString("build.version")
}

func (c *config) BuildRevision() string {
	return viper.GetString("build.revision")
}

func (c *config) BuildTime() string {
	return viper.GetString("build.time")
}

func (c *config) HttpHandler() string {
	return viper.GetString("http.handler")
}

func (c *config) HttpsHandler() string {
	return viper.GetString("https.handler")
}

func (c *config) HttpsKeystoreFile() string {
	return viper.GetString("https.keystore")
}

func (c *config) HttpsCertFile() string {
	return viper.GetString("https.cert")
}

func (c *config) Arango() (*arangoConfig, error) {
	conn, err := arangohttp.NewConnection(arangohttp.ConnectionConfig{
		Endpoints: viper.GetStringSlice("arango.endpoints"),
	})
	if err != nil {
		return nil, err
	}
	return &arangoConfig{
		ClientConfig: &driver.ClientConfig{
			Connection:     conn,
			Authentication: driver.BasicAuthentication(viper.GetString("arango.username"), viper.GetString("arango.password")),
		},
		Database: viper.GetString("arango.database"),
	}, nil
}

func (c *config) NatsUri() string {
	return viper.GetString("nats.uri")
}

func (c *config) AuthJwks() string {
	return viper.GetString("auth.jwks")
}

func (c *config) AuthClientId() string {
	return viper.GetString("auth.clientId")
}

func (c *config) AuthIss() string {
	return viper.GetString("auth.iss")
}

func (c *config) AuthSecret() string {
	return viper.GetString("auth.secret")
}

package config

import "github.com/spf13/viper"

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

func (c *config) MongoUri() string {
	return viper.GetString("mongo.uri")
}

func (c *config) MongoDatabase() string {
	return viper.GetString("mongo.database")
}

func (c *config) NatsUri() string {
	return viper.GetString("nats.uri")
}

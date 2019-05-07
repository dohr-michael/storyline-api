package universe

type Create struct {
	Name        string   `json:"name" mapstructure:"name"`
	Description string   `json:"description" mapstructure:"description"`
	Picture     string   `json:"picture,omitempty" mapstructure:"picture,omitempty"`
	Tags        []string `json:"tags" mapstructure:"tags"`
}

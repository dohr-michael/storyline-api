package user

import (
	"time"
)

type User struct {
	Email          string     `json:"email" mapstructure:"_key"`
	Name           string     `json:"name" mapstructure:"name"`
	Picture        string     `json:"picture" mapstructure:"picture"`
	Locale         string     `json:"locale" mapstructure:"locale"`
	Gender         string     `json:"gender" mapstructure:"gender"`
	CreatedAt      *time.Time `json:"createdAt" mapstructure:"createdAt"`
	LastConnection time.Time  `json:"lastConnection" mapstructure:"lastConnection"`
}
type Users []*User

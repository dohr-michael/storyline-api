package core

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/auth0/go-jwt-middleware"
	"github.com/dgrijalva/jwt-go"
	"github.com/dohr-michael/storyline-api/config"
	"net/http"
	"sync"
)

const UserContext = "UserContext"

func GetUserContext(ctx context.Context) *AuthenticatedUser {
	result, ok := ctx.Value(UserContext).(*AuthenticatedUser)
	if !ok {
		return nil
	}
	return result
}

type AuthenticatedUser struct {
	Token      *jwt.Token
	Email      string
	Name       string
	GivenName  string
	FamilyName string
	Locale     string
	Picture    string
	Gender     string
}

type Jwks struct {
	Keys []JSONWebKeys `json:"keys"`
}

type JSONWebKeys struct {
	Kty string   `json:"kty"`
	Kid string   `json:"kid"`
	Use string   `json:"use"`
	N   string   `json:"n"`
	E   string   `json:"e"`
	X5c []string `json:"x5c"`
}

type Auth struct {
	mutex sync.Mutex
	jwks  *Jwks
}

func NewAuth() *Auth {
	return &Auth{}
}

func (auth *Auth) loadJwks() error {
	auth.mutex.Lock()
	defer auth.mutex.Unlock()
	if auth.jwks == nil {
		resp, err := http.Get(config.Config.AuthJwks())
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		var jwks = Jwks{}
		err = json.NewDecoder(resp.Body).Decode(&jwks)
		if err != nil {
			return err
		}
		auth.jwks = &jwks
	}
	return nil
}

func (auth *Auth) getPem(token *jwt.Token) (string, error) {
	if auth.jwks == nil {
		err := auth.loadJwks()
		if err != nil {
			return "", err
		}
	}
	cert := ""

	for _, key := range auth.jwks.Keys {
		if token.Header["kid"] == key.Kid {
			cert = "-----BEGIN CERTIFICATE-----\n" + key.X5c[0] + "\n-----END CERTIFICATE-----"
		}
	}

	if cert == "" {
		err := errors.New("unable to find appropriate key")
		return cert, err
	}

	return cert, nil
}

func (auth *Auth) validateAndGetPk(token *jwt.Token) (interface{}, error) {
	// Check audience.
	if ok := token.Claims.(jwt.MapClaims).VerifyAudience(config.Config.AuthClientId(), true); !ok {
		return nil, fmt.Errorf("invalid audience")
	}
	// Check issuer
	if ok := token.Claims.(jwt.MapClaims).VerifyIssuer(config.Config.AuthIss(), true); !ok {
		return nil, fmt.Errorf("invalid iss")
	}
	// Check validity
	err := token.Claims.(jwt.MapClaims).Valid()
	if err != nil {
		return nil, err
	}
	// Check signature
	cert, err := auth.getPem(token)
	if err != nil {
		return nil, err
	}
	pk, err := jwt.ParseRSAPublicKeyFromPEM([]byte(cert))
	if err != nil {
		return nil, err
	}

	return pk, nil
}

func (auth *Auth) Middleware() func(http.Handler) http.Handler {
	jwtMiddleware := jwtmiddleware.New(jwtmiddleware.Options{
		SigningMethod:       jwt.SigningMethodRS256,
		ValidationKeyGetter: auth.validateAndGetPk,
		CredentialsOptional: true,
	})
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			token, ok := ctx.Value("user").(*jwt.Token)
			if !ok {
				next.ServeHTTP(w, r)
				return
			}
			claims := token.Claims.(jwt.MapClaims)
			readClaim := func(name string) string {
				v, ok := claims[name].(string)
				if !ok {
					return ""
				}
				return v
			}
			user := &AuthenticatedUser{
				Token:      token,
				Email:      readClaim("email"),
				Name:       readClaim("name"),
				GivenName:  readClaim("given_name"),
				FamilyName: readClaim("family_name"),
				Locale:     readClaim("locale"),
				Picture:    readClaim("picture"),
				Gender:     readClaim("gender"),
			}
			next.ServeHTTP(w, r.WithContext(context.WithValue(ctx, UserContext, user)))
		}
		return jwtMiddleware.Handler(http.HandlerFunc(fn))
	}
}

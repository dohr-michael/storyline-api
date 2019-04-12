package pkg

import (
	"context"
	"encoding/json"
	"github.com/dohr-michael/storyline-api/config"
	"github.com/dohr-michael/storyline-api/pkg/graphql"
	"github.com/dohr-michael/storyline-api/pkg/repo"
	"github.com/go-chi/chi"
	gographql "github.com/graphql-go/graphql"
	"github.com/nats-io/go-nats"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type requestOptions struct {
	Query         string                 `json:"query" url:"query" schema:"query"`
	Variables     map[string]interface{} `json:"variables" url:"variables" schema:"variables"`
	OperationName string                 `json:"operationName" url:"operationName" schema:"operationName"`
}

// a workaround for getting`variables` as a JSON string
type requestOptionsCompatibility struct {
	Query         string `json:"query" url:"query" schema:"query"`
	Variables     string `json:"variables" url:"variables" schema:"variables"`
	OperationName string `json:"operationName" url:"operationName" schema:"operationName"`
}

func getFromForm(values url.Values) *requestOptions {
	query := values.Get("query")
	if query != "" {
		// get variables map
		variables := make(map[string]interface{}, len(values))
		variablesStr := values.Get("variables")
		_ = json.Unmarshal([]byte(variablesStr), &variables)

		return &requestOptions{
			Query:         query,
			Variables:     variables,
			OperationName: values.Get("operationName"),
		}
	}

	return nil
}

func newRequestOptions(r *http.Request) *requestOptions {
	contentTypeStr := r.Header.Get("Content-Type")
	contentTypeTokens := strings.Split(contentTypeStr, ";")
	contentType := contentTypeTokens[0]
	switch contentType {
	case "application/graphql":
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			return &requestOptions{}
		}
		return &requestOptions{
			Query: string(body),
		}
	case "application/x-www-form-urlencoded":
		if err := r.ParseForm(); err != nil {
			return &requestOptions{}
		}

		if reqOpt := getFromForm(r.PostForm); reqOpt != nil {
			return reqOpt
		}

		return &requestOptions{}
	case "application/json":
		fallthrough
	default:
		var opts requestOptions
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			return &opts
		}
		err = json.Unmarshal(body, &opts)
		if err != nil {
			// Probably `variables` was sent as a string instead of an object.
			// So, we try to be polite and try to parse that as a JSON string
			var optsCompatible requestOptionsCompatibility
			_ = json.Unmarshal(body, &optsCompatible)
			_ = json.Unmarshal([]byte(optsCompatible.Variables), &opts.Variables)
		}
		return &opts
	}
}

func Start(
	router chi.Router,
) error {
	// Initialize tools
	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "text/html")
		_, _ = w.Write([]byte(graphqlPlayground))
	})
	router.Get("/@/voyager", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "text/html")
		_, _ = w.Write([]byte(graphqlVoyager))
	})

	// Initialize NATS connection.
	nc, err := nats.Connect(config.Config.NatsUri())
	if err != nil {
		return err
	}
	defer nc.Drain()
	c, err := nats.NewEncodedConn(nc, nats.JSON_ENCODER)
	if err != nil {
		return err
	}
	defer c.Drain()

	// Initialize repositories
	userRepo, err := repo.NewUserRepo()
	if err != nil {
		return err
	}

	// Initialize graphql route.
	schema, err := graphql.NewSchema()
	if err != nil {
		return err
	}
	router.Post("/", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		// Register all repository here.
		ctx = context.WithValue(ctx, graphql.UserRepoKey, userRepo)

		w.Header().Add("Content-Type", "application/json; charset=utf-8")
		opts := newRequestOptions(r)
		result := gographql.Do(gographql.Params{
			Schema:         schema,
			RequestString:  opts.Query,
			VariableValues: opts.Variables,
			OperationName:  opts.OperationName,
			Context:        ctx,
		})
		e := json.NewEncoder(w)
		_ = e.Encode(result)
	})

	return nil
}

package pkg

import (
	"github.com/dohr-michael/storyline-api/config"
	"github.com/dohr-michael/storyline-api/pkg/graphql"
	"github.com/go-chi/chi"
	ggraphql "github.com/graph-gophers/graphql-go"
	"github.com/graph-gophers/graphql-go/relay"
	"github.com/nats-io/go-nats"
	"net/http"
)

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

	// Initialize graphql route.
	resolver := &graphql.Resolver{}
	graph, err := ggraphql.ParseSchema(graphql.GetRootSchema(), resolver)
	if err != nil {
		return err
	}
	router.Post("/", (&relay.Handler{Schema: graph}).ServeHTTP)

	return nil
}

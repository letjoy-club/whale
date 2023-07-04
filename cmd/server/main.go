package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"whale/graph"
	"whale/pkg/loader"
	"whale/pkg/restfulserver"
	"whale/pkg/whaleconf"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/letjoy-club/mida-tool/authenticator"
	"github.com/letjoy-club/mida-tool/dbutil"
	"github.com/letjoy-club/mida-tool/graphqlutil"
	"github.com/letjoy-club/mida-tool/midacode"
	"github.com/letjoy-club/mida-tool/midacontext"
	"github.com/letjoy-club/mida-tool/proxy"
	"github.com/letjoy-club/mida-tool/pulsarutil"
	"github.com/letjoy-club/mida-tool/qcloudutil/clsutil"
	"github.com/letjoy-club/mida-tool/redisutil"
	"github.com/letjoy-club/mida-tool/tracerutil"
	"go.opentelemetry.io/otel/trace"
)

func main() {
	port := flag.Int("port", 11013, "port to run server on")
	confPath := flag.String("conf", "conf.yaml", "path to config file")

	upstream := flag.String("upstream", "", "upstream url, example: ws://localhost:11013/whale/v1/ws?userIds=1000,u_abc")

	flag.Parse()

	whaleConf := whaleconf.ReadConf(*confPath)
	db := whaleConf.DB()
	redis := whaleConf.Redis()
	loader := loader.NewLoader(db)
	cls := whaleConf.CLS()
	publisher := whaleConf.MatchingPublisher()

	tp := whaleConf.Trace()
	defer func() { tp.Shutdown(context.Background()) }()
	tr := tp.Tracer("whale-graph")

	gqlConf := graph.Config{Resolvers: &graph.Resolver{}}
	gqlConf.Directives.AdminOnly = func(ctx context.Context, obj interface{}, next graphql.Resolver) (interface{}, error) {
		token := midacontext.GetClientToken(ctx)
		if token.IsAdmin() {
			return next(ctx)
		}
		return nil, midacode.ErrNotPermitted
	}

	srv := handler.NewDefaultServer(graph.NewExecutableSchema(gqlConf))
	srv.AroundOperations(func(ctx context.Context, next graphql.OperationHandler) graphql.ResponseHandler {
		oc := graphql.GetOperationContext(ctx)
		opName := "unknown"
		if oc.OperationName != "" {
			opName = oc.OperationName
		}
		ctx, span := tr.Start(ctx, "whale."+opName, trace.WithSpanKind(trace.SpanKindServer))
		defer span.End()
		fmt.Printf("RawQuery: %s\n", oc.RawQuery)
		return next(ctx)
	})
	srv.SetErrorPresenter(graphqlutil.ErrorPresenter)

	r := chi.NewRouter()

	secret := []byte(whaleConf.Secret)
	auth := authenticator.Authenticator{Key: secret}

	services := midacontext.Services{
		Hoopoe: midacontext.NewServices(whaleConf.ServiceConf.Hoopoe, "1000"),
		Smew:   midacontext.NewServices(whaleConf.ServiceConf.Smew, "1000"),
		Scream: midacontext.NewServices(whaleConf.ServiceConf.Scream, "1000"),
	}

	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			token := midacontext.ParseToken(r, auth)

			if token.IsInvalid() {
				w.Header().Set("Content-Type", "application/json")
				encoder := json.NewEncoder(w)
				encoder.Encode(midacontext.GraphQLResp{
					Data:   struct{}{},
					Errors: []midacontext.GraphQLErr{{Message: string(token), Path: []string{}}},
				})
				return
			}

			if proxy.RedirectProxy(token.String(), w, r) {
				return
			}

			ctx = dbutil.WithDB(ctx, db)
			ctx = clsutil.WithGraphLogger(ctx, cls, whaleConf.QCloud.CLSConf.TopicID, "whale")
			ctx = redisutil.WithRedis(ctx, redis)
			ctx = pulsarutil.WithMQ[*whaleconf.Publisher](ctx, publisher)
			ctx = midacontext.WithLoader(ctx, loader)
			ctx = midacontext.WithServices(ctx, services)
			ctx = midacontext.WithClientToken(ctx, token)
			ctx = midacontext.WithAuthenticator(ctx, auth)
			ctx = tracerutil.WithSpanContext(ctx, r.Header)

			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
		})
	})

	r.Use(middleware.Logger)

	r.Route("/whale/", func(r chi.Router) {
		r.Route("/v2", func(r chi.Router) {
			r.Handle("/query", srv)
			r.Handle("/", playground.Handler("GraphQL playground", "/whale/v2/query"))
		})
		r.Route("/v1", func(r chi.Router) {
			restfulserver.Mount(r)
		})
	})

	if *upstream != "" {
		go proxy.Gateway(*upstream, r)
	}

	log.Printf("connect to http://localhost:%d/whale/v2 for GraphQL playground", *port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), r))
}

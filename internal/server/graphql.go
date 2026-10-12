package server

import (
	"context"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/OleksUMD/ecommerce_api/graph"
	"github.com/OleksUMD/ecommerce_api/graph/resolver"
	"github.com/gin-gonic/gin"
	"github.com/vektah/gqlparser/v2/ast"
)

func (s *Server) createGraphQLHandler() *handler.Server {
	rsvr := resolver.NewResolver(
		s.authService,
		s.userService,
		s.productService,
		s.cartService,
		s.orderService,
	)

	schema := graph.NewExecutableSchema(graph.Config{Resolvers: rsvr})
	srv := handler.New(schema)

	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})
	srv.AddTransport(transport.MultipartForm{})

	srv.SetQueryCache(lru.New[*ast.QueryDocument](1000))

	srv.Use(extension.Introspection{})
	srv.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New[string](100),
	})
	return srv
}

func (s *Server) graphQLHandler() gin.HandlerFunc {
	h := s.createGraphQLHandler()
	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

func (s *Server) graphQLPlaygroundHandler() gin.HandlerFunc {
	h := playground.Handler("GraphQL Playground", "/graphql/")
	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

func (s *Server) graphQLMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, _ := c.Get("user_id")
		userRole, _ := c.Get("user_role")
		userEmail, _ := c.Get("user_email")

		ctx := context.WithValue(c.Request.Context(), "user_id", userID)
		ctx = context.WithValue(ctx, "user_role", userRole)
		ctx = context.WithValue(ctx, "user_email", userEmail)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

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
	"github.com/OleksUMD/ecommerce_api/internal/utils"
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

func (s *Server) graphQLPlaygroundPublicHandler() gin.HandlerFunc {
	h := playground.Handler("GraphQL Playground(Public)", "/graphql/public")
	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

func (s *Server) graphQLMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, _ := c.Get(utils.UserIDKey)
		userRole, _ := c.Get(utils.UserRoleKey)
		userEmail, _ := c.Get(utils.UserEmailKey)

		ctx := context.WithValue(c.Request.Context(), utils.UserIDKey, userID)
		ctx = context.WithValue(ctx, utils.UserRoleKey, userRole)
		ctx = context.WithValue(ctx, utils.UserEmailKey, userEmail)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

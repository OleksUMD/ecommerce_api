package server_test

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/OleksUMD/ecommerce_api/internal/cache"
	"github.com/OleksUMD/ecommerce_api/internal/config"
	"github.com/OleksUMD/ecommerce_api/internal/database"
	"github.com/OleksUMD/ecommerce_api/internal/dto"
	"github.com/OleksUMD/ecommerce_api/internal/logger"
	"github.com/OleksUMD/ecommerce_api/internal/providers"
	"github.com/OleksUMD/ecommerce_api/internal/repositories"
	"github.com/OleksUMD/ecommerce_api/internal/server"
	"github.com/OleksUMD/ecommerce_api/internal/services"
	"github.com/OleksUMD/ecommerce_api/internal/utils"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"gorm.io/gorm"
)

var testDB *sql.DB
var testServer *http.Server
var testProductService *services.ProductService

func TestMain(m *testing.M) {
	os.Exit(runTests(m))
}

func runTests(m *testing.M) int {
	fmt.Println("Starting to test the server")
	var err error
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}
	postgresContainer, err := utils.StartTestDatabase(&cfg.Database)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := testcontainers.TerminateContainer(postgresContainer); err != nil {
			log.Printf("failed to terminate postgres container: %s", err)
		}
	}()
	redisContainer, err := utils.StartTestRedis(&cfg.Cache)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := testcontainers.TerminateContainer(redisContainer); err != nil {
			log.Printf("failed to terminate redis container: %s", err)
		}
	}()

	gormDB, err := database.New(&cfg.Database)
	if err != nil {
		panic(err)
	}
	testDB, err = gormDB.DB()
	if err != nil {
		panic(err)
	}
	migRunner, err := utils.NewMigrationsRunner(testDB)
	if err != nil {
		panic(err)
	}
	fmt.Println("Run migrations up")
	if err := migRunner.Up(); err != nil {
		panic(err)
	}
	cacheService := cache.NewCache(&cfg.Cache)
	defer func() {
		if err := cacheService.Close(); err != nil {
			fmt.Println("failed to close redis connection")
		}
	}()

	testServer = createTestServer(gormDB, cacheService, cfg)
	if err := seedDB(); err != nil {
		println("before panic seed db")
		panic(err)
	} else {
		println("database seeded")
	}
	code := m.Run()

	fmt.Println("Run migrations down")
	if err := migRunner.Down(); err != nil {
		panic(err)
	}
	if err := testDB.Close(); err != nil {
		panic(err)
	}
	return code
}

func cleanupTestData(t *testing.T) {
	var err error
	deleteTmpl := "DELETE FROM %s;"
	tables := []string{
		"users",
		"carts",
	}
	for _, tName := range tables {
		_, err = testDB.Exec(fmt.Sprintf(deleteTmpl, tName))
		assert.NoError(t, err)
	}
}

func createTestServer(db *gorm.DB, cache cache.ICacheProvider, cfg *config.Config) *http.Server {
	log := logger.New()
	eventPublisher := &utils.TestPublisher{}
	userRepo := repositories.NewUserRepository(db, cache)
	cartRepo := repositories.NewCartRepository(db)
	authService := services.NewAuthService(
		userRepo,
		cartRepo,
		cfg,
		eventPublisher,
	)
	userService := services.NewUserService(db, cfg)
	testProductService = services.NewProductService(db, cfg)

	uploadProvider := providers.NewLocalUploadProvider(cfg.Upload.Path)

	uploadService := services.NewUploadService(uploadProvider)
	cartService := services.NewCartService(db, cfg)
	orderService := services.NewOrderService(db, cfg)
	srv := server.New(
		cfg,
		log,
		authService,
		testProductService,
		userService,
		uploadService,
		cartService,
		orderService,
		cache,
	)
	router := srv.SetupRoutes()
	return server.NewHttpServer(
		fmt.Sprintf(":%s", cfg.Server.Port),
		router,
		10*time.Second,
		10*time.Second,
	)
}

func seedDB() error {
	if testProductService == nil {
		return errors.New("ProductService is not initialized")
	}
	categories := []dto.CreateCategoryRequest{
		dto.CreateCategoryRequest{
			Name:        "Category 1",
			Description: "Description Category 1",
		},
		dto.CreateCategoryRequest{
			Name:        "Category 2",
			Description: "Description Category 2",
		},
	}
	products := []dto.CreateProductRequest{
		dto.CreateProductRequest{
			Name:        "Product 1",
			Description: "Description Product 1",
			Price:       19.99,
			Stock:       10,
			SKU:         "SKU product 1",
		},
		dto.CreateProductRequest{
			Name:        "Product 2",
			Description: "Description Product 2",
			Price:       13.99,
			Stock:       13,
			SKU:         "SKU product 2",
		},
	}
	for ix, catReq := range categories {
		res, err := testProductService.CreateCategory(&catReq)
		if err != nil {
			return err
		}
		prodReq := products[ix]
		prodReq.CategoryID = res.ID
		_, err = testProductService.CreateProduct(&prodReq)
		if err != nil {
			return err
		}
	}
	return nil
}

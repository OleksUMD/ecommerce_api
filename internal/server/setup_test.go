package server_test

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"testing"

	"github.com/OleksUMD/ecommerce_api/internal/config"
	"github.com/OleksUMD/ecommerce_api/internal/database"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

var testDB *sql.DB

type Migration struct {
	OrderNo int
	Path    string
}
type MigrationRunner struct {
	up      []Migration
	down    []Migration
	db      *sql.DB
	dirPath string
}

func (m *MigrationRunner) Up() error {
	return m.run("up")
}

func (m *MigrationRunner) Down() error {
	return m.run("down")
}

func (m *MigrationRunner) run(direction string) error {
	var migrations []Migration
	switch direction {
	case "up":
		migrations = m.up
	case "down":
		migrations = m.down
	default:
		return errors.New("unknown migration command")
	}

	// All or nothing principle
	ctx := context.Background()
	tx, err := m.db.BeginTx(ctx, nil)
	if err != nil {
		println("failed start transaction")
		return err
	}
	defer func() {
		if err := tx.Rollback(); err != nil && !errors.Is(err, sql.ErrTxDone) {
			// The error is not due to a committed transaction
			panic(err)
		}
	}()
	for i := range migrations {
		fmt.Printf("Applying: %s", migrations[i].Path)
		rawSQL, err := os.ReadFile(migrations[i].Path)
		if err != nil {
			return err
		}
		_, err = tx.Exec(string(rawSQL))
		if err != nil {
			fmt.Println(" Failed!")
			return err
		}
		fmt.Println(" Success!")
	}
	return tx.Commit()
}

func (m *MigrationRunner) sort() error {
	sort.Slice(m.up, func(i, j int) bool {
		return m.up[i].OrderNo < m.up[j].OrderNo
	})
	// reverse order for downgrade
	sort.Slice(m.down, func(i, j int) bool {
		return m.down[i].OrderNo > m.down[j].OrderNo
	})
	return nil
}

func (m *MigrationRunner) setup() error {
	var migFileRE = `^([0-9]+)_([a-zA-Z0-9_]+).(up|down).sql$`
	migRE, err := regexp.Compile(migFileRE)
	if err != nil {
		return err
	}
	err = filepath.Walk(m.dirPath, func(path string, d fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			fileName := d.Name()
			submatch := migRE.FindStringSubmatch(fileName)
			if len(submatch) != 4 {
				return nil
			}
			orderNoStr := submatch[1]
			ordNoInt, err := strconv.Atoi(orderNoStr)
			if err != nil {
				return err
			}
			direction := submatch[3]
			migration := Migration{
				OrderNo: ordNoInt,
				Path:    path,
			}
			if direction == "up" {
				m.up = append(m.up, migration)
			} else {
				m.down = append(m.down, migration)
			}
		}
		return nil
	})
	if err != nil {
		return err
	}
	if err := m.sort(); err != nil {
		return err
	}
	return nil

}

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
	postgresContainer, err := startTestDatabase(&cfg.Database)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := testcontainers.TerminateContainer(postgresContainer); err != nil {
			log.Printf("failed to terminate container: %s", err)
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
	migrationsDirPath, err := getMigrationsDir()
	if err != nil {
		panic(err)
	}
	migRunner, err := NewMigrationsRunner(testDB, migrationsDirPath)
	if err != nil {
		panic(err)
	}
	fmt.Println("Run migrations up")
	if err := migRunner.Up(); err != nil {
		panic(err)
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

func NewMigrationsRunner(db *sql.DB, migrationsDirPath string) (*MigrationRunner, error) {
	migrationRunner := &MigrationRunner{
		up:      make([]Migration, 0),
		down:    make([]Migration, 0),
		db:      db,
		dirPath: migrationsDirPath,
	}
	if err := migrationRunner.setup(); err != nil {
		return nil, err
	}
	return migrationRunner, nil
}

func cleanupTestData(t *testing.T) {
	var err error
	deleteTmpl := "DELETE FROM %s;"
	tables := []string{
		"users",
		"carts",
	}
	println("check for nil dereference")
	println(testDB)
	for _, tName := range tables {
		_, err = testDB.Exec(fmt.Sprintf(deleteTmpl, tName))
		assert.NoError(t, err)
	}
}

func getMigrationsDir() (string, error) {
	migrationsDir := "db/migrations"
	out, err := exec.Command("go", "env", "GOMOD").Output()
	if err != nil {
		return "", err
	}
	return filepath.Join(filepath.Dir(string(out)), migrationsDir), nil
}

func startTestDatabase(cfg *config.DatabaseConfig) (*postgres.PostgresContainer, error) {
	ctx := context.Background()
	postgresContainer, err := postgres.Run(ctx,
		"postgres:17-alpine",
		postgres.WithDatabase(cfg.Name),
		postgres.WithUsername(cfg.User),
		postgres.WithPassword(cfg.Password),
		postgres.BasicWaitStrategies(),
	)
	if err != nil {
		log.Printf("failed to start container: %s", err)
		return nil, err
	}
	mappedPort, err := postgresContainer.MappedPort(ctx, "5432/tcp")
	if err != nil {
		return nil, err
	}
	cfg.Port = mappedPort.Port()

	return postgresContainer, nil
}

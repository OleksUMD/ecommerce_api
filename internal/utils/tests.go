package utils

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

	"github.com/OleksUMD/ecommerce_api/internal/config"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/modules/redis"
)

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
				return fmt.Errorf("unprocessible migration name: %s", path)
			}
			orderNoStr := submatch[1]
			ordNoInt, err := strconv.Atoi(orderNoStr)
			if err != nil {
				return err
			}
			migration := Migration{
				OrderNo: ordNoInt,
				Path:    path,
			}
			switch submatch[3] {
			case "up":
				m.up = append(m.up, migration)
			case "down":
				m.down = append(m.down, migration)
			default:
				return fmt.Errorf("unprocessible migration direction: %s", path)
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

func NewMigrationsRunner(db *sql.DB) (*MigrationRunner, error) {
	migrationsDirPath, err := getMigrationsDir()
	if err != nil {
		return nil, err
	}
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

type TestPublisher struct{}

func (tp *TestPublisher) Publish(eventType string, payload any, metadata map[string]string) error {
	return nil
}
func (tp *TestPublisher) Close() error {
	return nil
}

func StartTestDatabase(cfg *config.DatabaseConfig) (*postgres.PostgresContainer, error) {
	ctx := context.Background()
	postgresContainer, err := postgres.Run(ctx,
		"postgres:17-alpine",
		postgres.WithDatabase(cfg.Name),
		postgres.WithUsername(cfg.User),
		postgres.WithPassword(cfg.Password),
		postgres.BasicWaitStrategies(),
	)
	if err != nil {
		log.Printf("failed to start postgres container: %s", err)
		return nil, err
	}
	mappedPort, err := postgresContainer.MappedPort(ctx, "5432/tcp")
	if err != nil {
		return nil, err
	}
	cfg.Port = mappedPort.Port()

	return postgresContainer, nil
}

func StartTestRedis(cfg *config.CacheConfig) (*redis.RedisContainer, error) {
	ctx := context.Background()
	redisContainer, err := redis.Run(ctx,
		"redis:7",
		testcontainers.WithCmd("redis-server", "--requirepass", cfg.Password),
	)
	if err != nil {
		log.Printf("failed to start redis container: %s", err)
		return nil, err
	}
	mappedPort, err := redisContainer.MappedPort(ctx, "6379/tcp")
	if err != nil {
		return nil, err
	}
	cfg.Port = mappedPort.Port()

	return redisContainer, nil
}

func getMigrationsDir() (string, error) {
	migrationsDir := "db/migrations"
	out, err := exec.Command("go", "env", "GOMOD").Output()
	if err != nil {
		return "", err
	}
	return filepath.Join(filepath.Dir(string(out)), migrationsDir), nil
}

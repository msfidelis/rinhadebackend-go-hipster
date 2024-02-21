package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"runtime"
	"sync"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
)

var onceDB sync.Once
var onceDBPGX sync.Once
var onceBun sync.Once
var pgxInstance *sql.DB
var dbInstance *sql.DB
var BunInstance *bun.DB

func GetDBConn() *sql.DB {
	onceDB.Do(func() {
		var err error
		connectionString := getDBUrl()
		dbInstance, err = sql.Open("postgres", connectionString)
		if err != nil {
			log.Fatalf("Erro ao conectar com o banco de dados: %v", err)
		}

		// Verifica a conexão
		err = dbInstance.Ping()
		if err != nil {
			log.Fatalf("Erro ao estabelecer uma conexão com o banco de dados: %v", err)
		}
	})
	return dbInstance
}

// Retorna a conexão com o database em utilizando uma estratégia de Singleton
func GetPGX() *sql.DB {
	onceDBPGX.Do(func() {
		var err error
		config, err := pgx.ParseConfig(getDBUrl())
		if err != nil {
			panic(err)
		}
		pgxInstance = stdlib.OpenDB(*config)

		maxOpenConns := 4 * runtime.GOMAXPROCS(0)
		pgxInstance.SetMaxOpenConns(maxOpenConns)
		pgxInstance.SetMaxIdleConns(maxOpenConns)
	})
	return pgxInstance
}

func getDBUrl() string {
	user := os.Getenv("DATABASE_USER")
	pass := os.Getenv("DATABASE_PASSWORD")
	host := os.Getenv("DATABASE_HOST")
	port := os.Getenv("DATABASE_PORT")
	schema := os.Getenv("DATABASE_DB")

	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", user, pass, host, port, schema)
}

func GetDB() *bun.DB {
	onceBun.Do(func() {
		conn := GetPGX()
		// conn := GetDBConn()
		BunInstance = bun.NewDB(conn, pgdialect.New(), bun.WithDiscardUnknownColumns())
	})
	return BunInstance
}

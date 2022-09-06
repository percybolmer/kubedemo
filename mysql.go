package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type DatabaseConfig struct {
	username string
	password string
	hostname string
	port     string
	dbName   string
}

var (
	databaseConn *sql.DB
)

func connectDatabase() error {
	log.Println("Trying to connect to DB")
	db, err := sql.Open("mysql", createDSN(true))
	if err != nil {
		return fmt.Errorf("failed to open mysql connection: %w", err)
	}

	databaseConn = db

	if err := createDatabase(os.Getenv("DATABASE_NAME")); err != nil {
		return err
	}

	if err := db.Ping(); err != nil {
		return fmt.Errorf("failed to ping db: %w", err)
	}

	db, err = sql.Open("mysql", createDSN(false))
	if err != nil {
		return fmt.Errorf("failed to open mysql connection using databasename: %w", err)
	}

	log.Println("connected to database")
	databaseConn = db

	return nil

}

func createDatabase(dbname string) error {
	log.Println("Creating database")

	ctx, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFunc()

	tx, err := databaseConn.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	res, err := tx.ExecContext(ctx, fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s", dbname))
	if err != nil {
		return err
	}

	no, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if no == 0 {
		return errors.New("failed to create database, no rows affected")
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commmit tx: %w", err)
	}
	return nil

}

func createDSN(skipDB bool) string {
	dbCfg := getDatabaseConfig()
	if skipDB {
		return fmt.Sprintf("%s:%s@tcp(%s)/%s", dbCfg.username, dbCfg.password, dbCfg.hostname, "")
	}
	return fmt.Sprintf("%s:%s@tcp(%s)/%s", dbCfg.username, dbCfg.password, dbCfg.hostname, dbCfg.dbName)
}

func getDatabaseConfig() DatabaseConfig {
	return DatabaseConfig{
		username: os.Getenv("DATABASE_USERNAME"),
		password: os.Getenv("DATABASE_PASSWORD"),
		dbName:   os.Getenv("DATABASE_NAME"),
		hostname: os.Getenv("MYSQL_SERVICE_HOST"),
		port:     os.Getenv("MYSQL_SERVICE_port"),
	}
}

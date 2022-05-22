package database

import (
	"Noctobot/utils"
	"Noctobot/utils/colorize"
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/pollen5/minori"
	"github.com/spf13/cast"
	"os"
	"time"
)

const SCHEMA = `
CREATE TABLE IF NOT EXISTS guilds (
  id TEXT PRIMARY KEY NOT NULL UNIQUE,
  prefix TEXT DEFAULT '=',
);
CREATE TABLE IF NOT EXISTS members (
  id TEXT PRIMARY KEY NOT NULL UNIQUE,
  guild TEXT NOT NULL,
  points BIGINT DEFAULT 0,
  level INTEGER DEFAULT 0
);
CREATE TABLE IF NOT EXISTS modlogs (
  id TEXT PRIMARY KEY NOT NULL UNIQUE,
  "gid" TEXT,
  "uid" TEXT,
  "mid" TEXT,
  "aud" TEXT,
  "reason" TEXT,
  "type" TEXT
);
CREATE TABLE IF NOT EXISTS config (
  id TEXT PRIMARY KEY NOT NULL UNIQUE,
  qrcode BOOLEAN,
  "case" TEXT
);
CREATE TABLE IF NOT EXISTS plonked (
  id TEXT PRIMARY KEY NOT NULL UNIQUE,
  plonked BOOLEAN
);
`

var DB *sqlx.DB
var logger = minori.GetLogger("PostgreSQL")

func printParams(args ...interface{}) string {
	str := ""
	for i, arg := range args {
		str += fmt.Sprint(arg)
		if i != len(args)-1 {
			str += ", "
		}
	}
	return colorize.Bright(colorize.Yellow(str))
}

func Query(query string, args ...interface{}) (*sql.Rows, error) {
	logger.Debugf("SQL Query: %s (%s)", colorize.Bright(colorize.Cyan(query)), printParams(args...))
	return DB.Query(query, args...)
}

func QueryRow(query string, args ...interface{}) *sql.Row {
	logger.Debugf("SQL Query: %s (%s)", colorize.Bright(colorize.Cyan(query)), printParams(args...))
	return DB.QueryRow(query, args...)
}

func Exec(query string, args ...interface{}) (sql.Result, error) {
	logger.Debugf("SQL Query: %s (%s)", colorize.Bright(colorize.Cyan(query)), printParams(args...))
	return DB.Exec(query, args...)
}

func MustExec(query string, args ...interface{}) sql.Result {
	logger.Debugf("SQL Query: %s (%s)", colorize.Bright(colorize.Cyan(query)), printParams(args...))
	return DB.MustExec(query, args...)
}

func Begin() (*sqlx.Tx, error) {
	return DB.Beginx()
}

func MustBegin() *sqlx.Tx {
	return DB.MustBegin()
}

func Close() error {
	logger.Info("Shutting down Database connection...")
	return DB.Close()
}

func Get(dest interface{}, query string, args ...interface{}) error {
	logger.Debugf("SQL Query: %s (%s)", colorize.Bright(colorize.Cyan(query)), printParams(args...))
	return DB.Get(dest, query, args...)
}

func Select(dest interface{}, query string, args ...interface{}) error {
	logger.Debugf("SQL SEL: %s (%s)", colorize.Bright(colorize.Cyan(query)), printParams(args...))
	return DB.Select(dest, query, args...)
}

func init() {
	logger.Info("Connecting to PostgreSQL...")
	before := time.Now()
	if cast.ToBool(utils.GetConfig("dev")) {
		DB = sqlx.MustConnect("postgres", utils.GetConfig("postgresql"))
	} else {
		DB = sqlx.MustConnect("postgres", os.Getenv("DATABASE_URL"))
	}
	after := time.Now()
	logger.Infof("Connected to PostgreSQL! (took: %d ms)", after.Sub(before).Milliseconds())
}

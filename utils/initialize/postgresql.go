package initialize

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"

	"go.uber.org/fx"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

const (
	envPgDB       = "postgresql_database"
	envPgUser     = "postgresql_user"
	envPgPassword = "postgresql_password"
	envPgHost     = "postgresql_host"
	envPgPort     = "postgresql_port"
)

type UserDBOut struct {
	fx.Out

	DB *sql.DB `name:"userDB"`
}

// NewUserDB creates a connection to User
func NewUserDB(conf *viper.Viper, log *zap.Logger) (out UserDBOut, err error) {
	pgDB := conf.GetString(envPgDB)
	pgUser := conf.GetString(envPgUser)
	pgPassword := conf.GetString(envPgPassword)
	pgHost := conf.GetString(envPgHost)
	pgPort := conf.GetString(envPgPort)

	db, err := postgresqlInit(pgDB, pgUser, pgPassword, pgHost, pgPort, log)
	out = UserDBOut{
		DB: db,
	}
	return
}

func postgresqlInit(dbName, dbUser, dbPassword, dbHost, dbPort string, log *zap.Logger) (
	DB *sql.DB, err error) {
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)
	DB, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer DB.Close()

	err = DB.Ping()
	if err != nil {
		log.Fatal("postgresql connection failed", zap.Error(err), zap.String("uri", fmt.Sprint(dbName, " ", dbHost)))
		return
	}

	fmt.Println("Successfully connected!")
	log.Info("postgresql connected", zap.String("db", dbName))
	return
}

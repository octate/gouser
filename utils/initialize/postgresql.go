package initialize

import (
	"context"
	"fmt"
	"gouser/pkg/user"

	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"go.uber.org/fx"
)

const (
	envPgDB       = "postgresql_db"
	envPgUser     = "postgresql_user"
	envPgPassword = "postgresql_password"
	envPgHost     = "postgresql_host"
	envPgPort     = "postgresql_port"
)

type GoUserDBOut struct {
	fx.Out

	DB *pg.DB `name:"gouserDB"`
}
type dbLogger struct{}

func (d dbLogger) BeforeQuery(c context.Context, q *pg.QueryEvent) (context.Context, error) {
	return c, nil
}

func (d dbLogger) AfterQuery(c context.Context, q *pg.QueryEvent) error {
	qry, e := q.FormattedQuery()
	log.Info("query:", string(qry), "error: ", e)
	return nil
}

// NewGoUserDB creates a connection to User
func NewGoUserDB(conf *viper.Viper, log *logrus.Logger) (out GoUserDBOut, err error) {
	pgDB := conf.GetString(envPgDB)
	pgUser := conf.GetString(envPgUser)
	pgPassword := conf.GetString(envPgPassword)
	pgHost := conf.GetString(envPgHost)
	pgPort := conf.GetString(envPgPort)

	db, err := postgresqlInit(pgDB, pgUser, pgPassword, pgHost, pgPort, log)
	if err != nil {
		log.Error(err)
		return
	}
	out = GoUserDBOut{
		DB: db,
	}
	return
}

func postgresqlInit(dbName, dbUser, dbPassword, dbHost, dbPort string, log *logrus.Logger) (
	DB *pg.DB, err error) {

	//the DB variable below is a connection pool.
	DB = pg.Connect(
		&pg.Options{
			Addr:     fmt.Sprintf(dbHost + ":" + dbPort),
			User:     dbUser,
			Password: dbPassword,
			Database: dbName,
			OnConnect: func(ctx context.Context, db *pg.Conn) error {
				_, err := db.Exec("SET timezone = 'Asia/Calcutta'")
				return err
			},
		},
	)

	DB.AddQueryHook(dbLogger{})
	err = DB.Ping(DB.Context())

	if err != nil {
		log.WithFields(logrus.Fields{
			"error": err.Error(),
			"uri":   fmt.Sprint(dbName, " ", dbHost),
		}).Fatal("postgresql connection failed")
		return
	}

	createSchema(DB)
	log.Info("Successfully connected!")
	log.WithFields(logrus.Fields{
		"database": dbName,
	}).Info("postgresql connected")
	return
}

func createSchema(db *pg.DB) error {
	models := []interface{}{
		(*user.User)(nil),
	}

	for _, model := range models {
		err := db.Model(model).CreateTable(&orm.CreateTableOptions{
			Temp:        false,
			IfNotExists: true,
		})
		if err != nil {
			return err
		}
	}
	return nil
}

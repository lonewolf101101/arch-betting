package app

import (
	"log"
	"os"
	"time"

	"git.bolor.net/bolorsoft/micro"
	"github.com/lonewolf101101/Architect-betting/backend/common/apputils"
	"github.com/lonewolf101101/Architect-betting/backend/pkg/actionlogman"
	"github.com/lonewolf101101/Architect-betting/backend/pkg/customerman"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	ErrorLog  *log.Logger
	InfoLog   *log.Logger
	DB        *gorm.DB
	Mode      string
	Customers *customerman.Service
	Actions   *actionlogman.Service

	Config      = conf{}
	Location    *time.Location
	KafkaServer *micro.KafkaServer
)

func Init(path, mode string) {
	InfoLog = log.New(os.Stdout, "[superset]\tINFO\t", log.Ldate|log.Ltime)
	ErrorLog = log.New(os.Stderr, "[superset]\tERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	Mode = mode

	loc, err := time.LoadLocation("Asia/Ulaanbaatar")
	if err != nil {
		ErrorLog.Fatal(err)
	}
	Location = loc

	apputils.LoadConfig(&Config, path, mode)
	DB = OpenDB(Config.DSN)

	Customers = customerman.NewService(DB, InfoLog, ErrorLog)
	Actions = actionlogman.NewService(DB, InfoLog, ErrorLog)

	m, err := micro.NewKafkaServer(&micro.Options{
		URL:      Config.KafkaHost,
		InfoLog:  InfoLog,
		ErrorLog: ErrorLog,
	})
	if err != nil {
		ErrorLog.Fatal(err)
	}
	KafkaServer = m
}

func Close() {
	KafkaServer.Close()
}

func OpenDB(dsn string) *gorm.DB {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		// Logger:                                   logger.Default.LogMode(logger.Info),
		SkipDefaultTransaction:                   true,
		PrepareStmt:                              true,
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		log.Fatal(err)
	}
	return db
}

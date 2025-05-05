package app

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"git.bolor.net/bolorsoft/micro"
	"github.com/golangcollege/sessions"
	"github.com/lonewolf101101/Architect-betting/backend/common/apputils"
	"github.com/lonewolf101101/Architect-betting/backend/common/easyOAuth2"
	"github.com/lonewolf101101/Architect-betting/backend/pkg/actionlogman"
	"github.com/lonewolf101101/Architect-betting/backend/pkg/customerman"
	"github.com/lonewolf101101/Architect-betting/backend/pkg/mailerman"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"gorm.io/gorm"
)

var (
	ErrorLog                *log.Logger
	InfoLog                 *log.Logger
	Session                 *sessions.Session
	Config                  = conf{}
	Mode                    string
	Location                *time.Location
	CustomerConnectionMutex = new(sync.RWMutex)
	Google                  *easyOAuth2.EasyOAuthClient
	DB                      *gorm.DB
	KafkaServer             *micro.KafkaServer

	// #region Services
	Customers  *customerman.Service
	ActionLogs *actionlogman.Service
	Mailer     *mailerman.Service
)

const (
	GB = 1 << 30
	MB = 1 << 20
	KB = 1 << 10
)

//#region Init

func Init(path, mode string) {
	InfoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	ErrorLog = log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	Mode = mode
	fmt.Println(mode)
	loc, err := time.LoadLocation("Asia/Ulaanbaatar")
	if err != nil {
		panic(err)
	}
	Location = loc

	apputils.LoadConfig(&Config, path, mode)

	DB = apputils.OpenDB(Config.DSN)

	Customers = customerman.NewService(DB, InfoLog, ErrorLog)
	ActionLogs = actionlogman.NewService(DB, InfoLog, ErrorLog)
	// Mailer = mailerman.NewService(DB, InfoLog, ErrorLog)
	// FrontendWS = websocket.New()

	Session = sessions.New([]byte(Config.SessionSecret))
	Session.Lifetime = 72 * time.Hour

	Google = &easyOAuth2.EasyOAuthClient{
		Name: "google",
		Config: &oauth2.Config{
			RedirectURL:  Config.OAuth2.RedirectURL,
			ClientID:     Config.OAuth2.ClientID,
			ClientSecret: Config.OAuth2.ClientSecret,
			Scopes:       Config.OAuth2.Scopes,
			Endpoint:     google.Endpoint,
		},
		UserInfoEndpoint: Config.OAuth2.UserInfoEndpoint,
	}

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
}

// #region use with caution
func PanicOnError(err error) {
	if err != nil {
		panic(err)
	}
}

func PrintOnError(err error) {
	if err != nil {
		ErrorLog.Println(err)
	}
}

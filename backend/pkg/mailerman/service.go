package mailerman

import (
	"log"

	"gorm.io/gorm"
)

type Service struct {
	DB       *gorm.DB
	infoLog  *log.Logger
	errorLog *log.Logger
}

func NewService(db *gorm.DB, infoLog, errorLog *log.Logger) *Service {
	return &Service{
		DB:       db,
		infoLog:  infoLog,
		errorLog: errorLog,
	}
}

// func (s *Service) GenerateEmail(email string, mailer Mailer, resultInfos *kudos.CompanyResultInfos) bool {

// 	s.infoLog.Println("daily task: generate email")

// 	to := email
// 	m := gomail.NewMessage()
// 	m.SetHeader("From", mailer.From)
// 	m.SetHeader("To", to)
// 	m.SetHeader("Subject", mailer.Subject)

// 	tpl, err := os.ReadFile(mailer.Path)
// 	if err != nil {
// 		panic(err)
// 	}

// 	ctx := map[string]any{
// 		"CompanyName": resultInfos.CompanyName,
// 		"Year":        resultInfos.Year,
// 		"Month":       resultInfos.Month,
// 		"Day":         resultInfos.Day,
// 		"Infos":       resultInfos.ResultInfos,
// 		"Date":        mailer.Date,
// 	}

// 	result, err := raymond.Render(string(tpl), ctx)
// 	if err != nil {
// 		panic(err)
// 	}

// 	s.infoLog.Println("daily task: parsed template")

// 	m.SetBody("text/html", result)

// 	d := gomail.NewDialer(mailer.Host, mailer.Port, mailer.Username, mailer.Password)
// 	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
// 	// Send emails using d.

// 	if err := d.DialAndSend(m); err != nil {
// 		s.errorLog.Printf("Error at sending message to '%s'. Error:%s.", to, err)
// 		return false
// 	}
// 	s.infoLog.Printf("Message to %s.", to)

// 	return true
// }

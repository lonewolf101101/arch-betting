package mailerman

type Mailer struct {
	Host     string
	Port     int
	Username string
	Password string
	From     string
	Subject  string
	Date     string
	Path     string
}
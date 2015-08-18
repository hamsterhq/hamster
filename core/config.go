package core

/*.toml reader struct*/
type config struct {
	Title   string
	Author  author
	Servers map[string]server
	Logins  map[string]login
	DB      map[string]database `toml:"databases"`
	Clients map[string]client
}

type author struct {
	Name         string
	Organization string
	Email        string
}

type server struct {
	Port         int
	Host         string
	CookieSecret string
}

type login struct {
	ID     string
	Secret string
}

type database struct {
	Username string
	Password string
	Name     string
	Host     string
	Port     int
}

type client struct {
	IP     string
	Secret string
	Token  string
}

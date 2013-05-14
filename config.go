package hamster

/*.toml reader struct*/
type Config struct {
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
	Port int
	Host string
}

type login struct {
	Id     string
	Secret string
}

type database struct {
	Username string
	Password string
	Name     string
	Host     string
}

type client struct {
	Secret string
}

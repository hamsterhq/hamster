package hamster

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"testing"
)

func TestTomlParsing(t *testing.T) {
	var config Config
	if _, err := toml.DecodeFile("hamster_sample.toml", &config); err != nil {
		t.Fatalf("toml failed to parse: %v", err)
		return
	}

	fmt.Printf("Title: %s\n", config.Title)
	fmt.Printf("Author: %s (%s, %s)\n",
		config.Author.Name, config.Author.Organization, config.Author.Email)

	for serverName, server := range config.Servers {
		fmt.Printf("Server: %s (%d, %s)\n", serverName, server.Port, server.Host)
	}

	for loginName, login := range config.Logins {
		fmt.Printf("Login: %s (%s, %s)\n", loginName, login.Id, login.Secret)
	}

	for dbName, db := range config.DB {
		fmt.Printf("Database: %s (%s, %s, %s, %s)\n", dbName, db.Username, db.Password, db.Name, db.Host)
	}

	for clientName, client := range config.Clients {
		fmt.Printf("Client: %s , %s\n", clientName, client.Secret)
	}

}

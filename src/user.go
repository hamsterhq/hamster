package main

import (
	"database/sql"
	"fmt"
	"time"
)

type timestamp time.Time

func (t timestamp) MarshalJSON() ([]byte, error) {
	return []byte(time.Time(t).Format(`"2006-01-02T15:04:05Z"`)), nil
}

type User struct {
	Id        string    `json:"id"`
	Email     string    `json:"email,omitempty"`
	CreatedAt timestamp `json:"created_at"`
	UpdatedAt timestamp `json:"updated_at"`
	Name      string    `json:"name,omitempty"`
	Admin     bool      `json:"admin"`
	Active    bool      `json:"-"`
}

func GetUserByApiKey(db *sql.DB, key Credentials) (*User, error) {

	rows, err := db.Query("select id, email, name, admin, active, created_at, updated_at from users where apikey = $1", string(key))
	fmt.Println("executed query")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	if rows.Next() {
		fmt.Println("found user!")
		user := new(User)
		var id []byte
		var email []byte
		var name []byte
		err = rows.Scan(&id, &email, &name, &user.Admin, &user.Active, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			fmt.Println("can't scan!")
			return nil, err
		}
		user.Id = string(id)
		user.Email = string(email)
		user.Name = string(name)
		return user, nil
	}
	return nil, nil
}

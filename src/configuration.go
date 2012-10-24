package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

type Config struct {
	Db struct {
		Name     string
		User     string
		Password string
	}
}

func ReadConfig() Config {
	if len(os.Args) != 2 {
		log.Fatal("You must supply a configuration filename")
	}
	filename := os.Args[1]
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	var c Config
	err = json.Unmarshal(bytes, &c)
	if err != nil {
		panic(err)
	}
	return c
}

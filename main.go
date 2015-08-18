package main

import "github.com/hamsterhq/hamster/core"

func main() {
	done := make(chan bool)
	server := core.NewServer("hamster.toml")
	server.ListenAndServe()
	<-done
}

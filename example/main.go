package main

import (
	"net/http"

	"github.com/ThatsMrTalbot/scaffold"
)

func main() {
	app := NewApp()

	dispatcher := scaffold.DefaultDispacher()
	scaffold.Scaffold(dispatcher, app)

	http.ListenAndServe(":8080", dispatcher)
}

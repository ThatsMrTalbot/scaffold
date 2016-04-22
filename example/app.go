package main

import "github.com/ThatsMrTalbot/scaffold"

// App is an example platform
type App struct {
	hello *Hello `inject:""`
}

// NewApp creates an example app
func NewApp() *App {
	return &App{
		hello: &Hello{},
	}
}

// Routes implements scaffold.Platform.Routes
func (a *App) Routes(router *scaffold.Router) {
	router.Platform("hello", a.hello)
}

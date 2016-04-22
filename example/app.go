package main

import "github.com/ThatsMrTalbot/scaffold"

type App struct {
	hello *Hello `inject:""`
}

func NewApp() *App {
	return &App{
		hello: &Hello{},
	}
}

func (a *App) Routes(router *scaffold.Router) {
	router.Platform("hello", a.hello)
}

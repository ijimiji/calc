package app

import "calc/internal/view"

func New(view *view.View) *App {
	return &App{
		view: view,
	}
}

type App struct {
	view *view.View
}

func (a *App) Run() {
	a.view.Show()
}

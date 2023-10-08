package main

import (
	"calc/internal/app"
	"calc/internal/calculator"
	"calc/internal/view"
)

func main() {
	app.New(view.New(calculator.New())).Run()
}

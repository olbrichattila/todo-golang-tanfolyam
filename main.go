package main

import (
	"todo/internal/app"
)

func main() {
	todoApp := app.New()
	todoApp.Serve()
}

package main

import (
	"context"
	"fmt"

	"github.com/weswest/msds431wk9/backend"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

func (a *App) CheckTermOriginal(term string) string {
	answer, exists, err := backend.TermExistsOriginal(term)
	if err != nil {
		return fmt.Sprintf("Error: %s", err.Error())
	}
	if exists {
		return fmt.Sprintf("For term '%s', found answer: %s", term, answer)
	}
	return fmt.Sprintf("No answer found for term '%s'", term)
}

func (a *App) CheckTermEmbed(term string) string {
	answer, exists, err := backend.TermExistsEmbed(term)
	if err != nil {
		return fmt.Sprintf("Error: %s", err.Error())
	}
	if exists {
		return fmt.Sprintf("For term '%s', found answer: %s", term, answer)
	}
	return fmt.Sprintf("No answer found for term '%s'", term)
}

func (a *App) CheckDir(term string) string {
	answer, _, _ := backend.CheckDirectory(term)
	return fmt.Sprintf("For term '%s', found answer: %s", term, answer)
}

func (a *App) RunDebug() string {
	answer, err := backend.ListEmbeddedFiles()
	if err != nil {
		return fmt.Sprintf("Debug result is: %s; error: %s", answer, err.Error())
	}
	return fmt.Sprintf("Debug result is: %s; error: nil", answer)
}

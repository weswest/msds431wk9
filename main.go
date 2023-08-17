package main

import (
	"embed"
	"fmt"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"

	"github.com/weswest/msds431wk9/backend"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {

	create := true
	if create {
		// Create an instance of the app structure
		app := NewApp()

		// Create application with options
		err := wails.Run(&options.App{
			Title:  "Wk9",
			Width:  1024,
			Height: 768,
			AssetServer: &assetserver.Options{
				Assets: assets,
			},
			BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
			OnStartup:        app.startup,
			Bind: []interface{}{
				app,
			},
		})

		if err != nil {
			println("Error:", err.Error())
		}
	}

	// This creates the db if it doesn't exist
	forcedRebuild := true
	err := backend.CreateDatabaseIfNeeded(forcedRebuild)
	if err != nil {
		fmt.Println(err)
		return
	}

	dbDebug := false
	if dbDebug {
		results, err := backend.ListDatabase()
		if err != nil {
			fmt.Println(err)
			return
		}

		for _, result := range results {
			fmt.Println(result)
		}

		// Testing TermExists function
		term := "defer"
		answer, exists, _ := backend.TermExistsOriginal(term)
		if exists {
			fmt.Printf("For term '%s', found answer: %s\n", term, answer)
		} else {
			fmt.Printf("No answer found for term '%s'\n", term)
		}
	}
}

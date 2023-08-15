package main

import (
	"fmt"

	"github.com/weswest/msds431wk9/backend"
)

func testMain() {
	// This creates the db if it doesn't exist
	err := backend.CreateDatabase()
	if err != nil {
		fmt.Println(err)
		return
	}

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
	answer, exists := backend.TermExists(term)
	if exists {
		fmt.Printf("For term '%s', found answer: %s\n", term, answer)
	} else {
		fmt.Printf("No answer found for term '%s'\n", term)
	}
}

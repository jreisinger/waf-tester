package main

import (
	"log"
)

func main() {
	// Get inputs via CLI flags and/or parameters.
	flags, err := ParseFlags()
	if err != nil {
		log.Fatalf("cannot parse flags: %v", err)
	}
	log.Printf("flags: %+v", flags)

	// Get tests to execute.
	/* 	alltests, err := httptest.GetTests(*testspath, *only, *scheme)
	   	if err != nil {
	   		log.Fatalf("can't get tests: %v", err)
	   	}
	   	log.Printf("flags: %+v", alltests) */

	/*
		Execute the tests.
		Print the results.
	*/
}

package main

import (
	"fmt"
	"mock-api/databases"
)

func main() {
	/*
	*
	*This is the main file for running the migrations
	*It uses the local "databases" pkg to import DB related functions
	 */
	fmt.Println("Running Migrations...")
	databases.CreateTables()
}

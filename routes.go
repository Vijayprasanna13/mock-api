package main

import (
	"fmt"
	"github.com/gorilla/handlers"
	"github.com/julienschmidt/httprouter"
	"log"
	"mock-api/Users"
	"net/http"
	"os"
)

func main() {

	port := "8000"
	fmt.Println("Serving on port : " + port)
	router := httprouter.New()

	/*
	*API routes
	 */
	router.POST("/mock-api/user/auth", Users.Authenticate())
    router.POST("/mock-api/user/add", Users.AddUser())

	/*
	*
	*Serve the app via port <8000>
	 */
	log.Fatal(http.ListenAndServe(":"+port, handlers.LoggingHandler(os.Stdout, router)))
}

package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"cloud.google.com/go/firestore"
	api "github.com/eelcoh/go-api"
)

var db_user string
var db_password string
var db_name string

var signingKey string
var passphrase string

var ctx context.Context
var client *firestore.Client

func init() {

	fmt.Println("initialising firestore")

	var err error

	ctx := context.Background()
	client, err := firestore.NewClient(ctx, "websites-394411")
	if err != nil {
		log.Fatalf("Error initializing Cloud Firestore client: %v", err)
	}
	defer client.Close()

	fmt.Println("initialising environment variables")

	signingKey = os.Getenv("SIGNINGKEY")

	passphrase = os.Getenv("PASSPHRASE")

}

func main() {
	fmt.Println("starting")

	// set port
	var port string

	if len(os.Args) == 2 {
		if _, err := strconv.Atoi(os.Args[1]); err == nil {
			port = fmt.Sprintf(":%s", os.Args[1])
		} else {
			port = ":8080"
		}
	} else {
		port = ":8080"
	}
	log.Printf("Serving at port %s", port)

	// defining routes
	router := api.NewRouter(routes)

	fmt.Println("routes defined")

	// starting the engine
	log.Fatal(http.ListenAndServe(port, router))
	fmt.Println("running")

}

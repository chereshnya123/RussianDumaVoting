package main

import (
	"dumaVote/server"
	"fmt"
	"log"
	"net/http"
	"os"
)

func isExist(env_name string) bool {
	_, isExists := os.LookupEnv("APP_API_KEY")
	return isExists
}

func checkEnv() {
	if !isExist("APP_API_KEY") {
		log.Fatal("App api key is not set. APP_API_KEY env is empty")
	}
	if !isExist("PERSONAL_API_KEY") {
		log.Fatal("Personal api key is not set. PERSONAL_API_KEY env is empty")
	}
}

func main() {
	checkEnv()
	port := "8080"
	dumaVotesServer := server.NewDumaVotesServer(os.Getenv("APP_API_KEY"), os.Getenv("PERSONAL_API_KEY"))
	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: http.HandlerFunc(dumaVotesServer.MainHandler),
	}
	log.Printf("Start server on http://localhost%s", port)
	err := server.ListenAndServe()
	if err != nil {
		fmt.Fprintf(os.Stderr, "[INFO ]: Http serve shut down. Error = %s\n", err.Error())
	}
}

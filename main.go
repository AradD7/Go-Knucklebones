package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/AradD7/Go-Knuclebones/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	db 			*database.Queries
	tokenSecret	string
}

func main() {
	godotenv.Load()

	secret 				:= os.Getenv("TOKEN_SECERT")
	const port 			= "8080"
	const filepathRoot 	= "."
	dbURL 				:= os.Getenv("DB_URL")

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}

	apiCfg := apiConfig{
		db: 		 database.New(db),
		tokenSecret: secret,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/players", apiCfg.handlerNewPlayer)
	mux.HandleFunc("POST /api/games", apiCfg.handlerNewGame)



	srv := &http.Server{
		Handler: mux,
		Addr: 	 ":" + port,
	}

	fmt.Printf("Seving files from %s on port:%s\n", filepathRoot, port)
	srv.ListenAndServe()
}

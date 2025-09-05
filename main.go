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
	platform 	string
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
		platform: 	 os.Getenv("PLATFORM"),
	}

	mux := http.NewServeMux()
	mux.HandleFunc("POST /admin/reset", apiCfg.handlerReset)

	mux.HandleFunc("POST /api/players", apiCfg.handlerNewPlayer)
	mux.HandleFunc("POST /api/players/login", apiCfg.handlerPlayerLogin)

	//mux.HandleFunc("POST /api/rolls", apiCfg.handlerRoll)

	//mux.HandleFunc("GET /api/games", apiCfg.handlerNewGame) returns all games for player
	//mux.HandleFunc("POST /api/games/{game_id}", apiCfg.handlerNewGame) returns specific game
	mux.HandleFunc("POST /api/games/new", apiCfg.handlerNewGame)

	//mux.HandleFunc("POST /api/moves", apiCfg.handlerNewGame)

	srv := &http.Server{
		Handler: mux,
		Addr: 	 ":" + port,
	}

	fmt.Printf("Seving files from %s on port:%s\n", filepathRoot, port)
	srv.ListenAndServe()
}

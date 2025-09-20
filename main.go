package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/AradD7/Go-Knuclebones/internal/database"
	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173") // React dev server
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Handle preflight OPTIONS request
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Continue to the next handler
		next.ServeHTTP(w, r)
	})
}

type gameServer struct {
	connections map[string][]*websocket.Conn
	rwMux 		*sync.RWMutex
}

type apiConfig struct {
	db 			*database.Queries
	tokenSecret	string
	platform 	string
	gs  *gameServer
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
		gs: 		 &gameServer{
			connections: make(map[string][]*websocket.Conn),
			rwMux: 		 &sync.RWMutex{},
		},
	}

	mux := http.NewServeMux()

	mux.Handle("/app/", http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot))))

	mux.HandleFunc("POST /admin/reset", apiCfg.handlerReset)

	mux.HandleFunc("POST /api/players/new", apiCfg.handlerNewPlayer)
	mux.HandleFunc("POST /api/players/login", apiCfg.handlerPlayerLogin)

	mux.HandleFunc("GET /api/rolls", handlerRoll)

	mux.HandleFunc("GET /api/games", apiCfg.handlerGetGames)
	mux.HandleFunc("GET /api/games/{game_id}", apiCfg.handlerGetGame)
	mux.HandleFunc("GET /api/games/new", apiCfg.handlerNewGame)
	mux.HandleFunc("GET /api/games/{game_id}/join", apiCfg.handlerJoinGame)
	mux.HandleFunc("POST /api/games/move/{game_id}", apiCfg.handlerMakeMove)

	mux.HandleFunc("/ws/games/{game_id}", apiCfg.handlerWebSocket)

	mux.HandleFunc("POST /api/games/localgame", apiCfg.handlerLocalGame)


	srv := &http.Server{
		Handler: corsMiddleware(mux),
		Addr: 	 ":" + port,
	}

	fmt.Printf("Api available on port:%s\n", port)
	srv.ListenAndServe()
}

package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	httphandler "github.com/example/game-designer-server/internal/http"
	"github.com/example/game-designer-server/internal/gamestate"
	"github.com/example/game-designer-server/internal/leaderboard"
	"github.com/example/game-designer-server/internal/profile"
	"github.com/example/game-designer-server/internal/session"
	"github.com/example/game-designer-server/internal/store"
)

func main() {
	s := store.New()

	sessSvc := session.NewService(s, 24*time.Hour)
	profSvc := profile.NewService(s)
	gsSvc := gamestate.NewService(s)
	lbSvc := leaderboard.NewService(s)

	handler := httphandler.NewHandler(sessSvc, profSvc, gsSvc, lbSvc)

	mux := http.NewServeMux()
	handler.RegisterRoutes(mux)

	addr := ":8080"
	fmt.Printf("Game Designer Server starting on %s\n", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}

package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/example/game-designer-server/internal/balance"
	httphandler "github.com/example/game-designer-server/internal/http"
	"github.com/example/game-designer-server/internal/profile"
	"github.com/example/game-designer-server/internal/session"
	"github.com/example/game-designer-server/internal/slot"
	"github.com/example/game-designer-server/internal/store"
)

func main() {
	s := store.New()

	balSvc := balance.NewService(s)
	sessSvc := session.NewService(s, 24*time.Hour, balSvc)
	profSvc := profile.NewService(s)
	slotSvc := slot.NewService(s)

	handler := httphandler.NewHandler(sessSvc, profSvc, slotSvc, balSvc)

	mux := http.NewServeMux()
	handler.RegisterRoutes(mux)

	addr := ":8080"
	fmt.Printf("Game Designer Slot Server starting on %s\n", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}

package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/HumanDignityGuardian/HDGP-protocol/internal/engine"
	"github.com/HumanDignityGuardian/HDGP-protocol/internal/gateway"
)

func main() {
	addr := os.Getenv("HDGP_ENGINE_ADDR")
	if addr == "" {
		addr = ":8080"
	}

	mux := http.NewServeMux()
	mux.Handle("/hdgp/v1/evaluate", engine.NewEvaluateHandler())
	mux.Handle("/hdgp/v1/chat", gateway.NewChatHandler(engine.Evaluate))
	mux.Handle("/hdgp/v1/audit", engine.NewAuditHandler())
	mux.Handle("/hdgp/v1/appeal", engine.NewAppealHandler())
	mux.Handle("/hdgp/v1/status", engine.NewStatusHandler())

	server := &http.Server{
		Addr:              addr,
		Handler:           loggingMiddleware(mux),
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Printf("HDGP Engine listening on %s", addr)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("server error: %v", err)
	}
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%s %s %s", r.Method, r.URL.Path, time.Since(start))
	})
}


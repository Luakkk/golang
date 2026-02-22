package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"

	"practice3-go/internal/app"
	"practice3-go/internal/handler"
	"practice3-go/internal/middleware"
)

func main() {
	loadDotEnv()

	cfg := app.LoadConfig()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	db, userUC, err := app.Build(ctx, cfg)
	if err != nil {
		log.Fatalf("build error: %v", err)
	}
	defer db.DB.Close()

	r := chi.NewRouter()
	r.Use(middleware.LoggingMiddleware)
	r.Use(middleware.AuthMiddleware(cfg.APIKey))

	h := handler.NewUserHandler(userUC)

	r.Get("/health", h.Health)
	r.Route("/users", func(r chi.Router) {
		r.Get("/", h.GetUsers)
		r.Post("/", h.CreateUser)
		r.Get("/{id}", h.GetUserByID)
		r.Put("/{id}", h.UpdateUser)
		r.Delete("/{id}", h.DeleteUserByID)
	})

	srv := &http.Server{
		Addr:         ":8080",
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	go func() {
		log.Printf("server started on %s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctxShutdown, cancelShutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelShutdown()

	_ = srv.Shutdown(ctxShutdown)
	log.Printf("server stopped")
}

// optional .env loader (в PDF optional EASY) :contentReference[oaicite:19]{index=19}
func loadDotEnv() {
	b, err := os.ReadFile(".env")
	if err != nil {
		return
	}

	lines := splitLines(string(b))
	for _, ln := range lines {
		ln = trimSpace(ln)
		if ln == "" || ln[0] == '#' {
			continue
		}
		eq := indexByte(ln, '=')
		if eq <= 0 {
			continue
		}
		k := trimSpace(ln[:eq])
		v := trimSpace(ln[eq+1:])
		if k != "" && os.Getenv(k) == "" {
			_ = os.Setenv(k, v)
		}
	}
}

func splitLines(s string) []string {
	var out []string
	start := 0
	for i := 0; i < len(s); i++ {
		if s[i] == '\n' {
			out = append(out, s[start:i])
			start = i + 1
		}
	}
	if start <= len(s)-1 {
		out = append(out, s[start:])
	}
	return out
}

func trimSpace(s string) string {
	i := 0
	j := len(s) - 1
	for i <= j && (s[i] == ' ' || s[i] == '\r' || s[i] == '\t') {
		i++
	}
	for j >= i && (s[j] == ' ' || s[j] == '\r' || s[j] == '\t') {
		j--
	}
	if i > j {
		return ""
	}
	return s[i : j+1]
}

func indexByte(s string, b byte) int {
	for i := 0; i < len(s); i++ {
		if s[i] == b {
			return i
		}
	}
	return -1
}
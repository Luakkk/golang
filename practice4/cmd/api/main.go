package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/lib/pq"
)

type Movie struct {
	ID      int64  `json:"id"`
	Title   string `json:"title"`
	Genre   string `json:"genre"`
	Budget  int64  `json:"budget"`
	Hero    string `json:"hero"`
	Heroine string `json:"heroine"`
}

type Config struct {
	AppPort   string
	DBHost    string
	DBPort    string
	DBUser    string
	DBPass    string
	DBName    string
	DBSSLMode string
}

type App struct {
	db *sql.DB
}

func main() {
	cfg := loadConfig()

	// Wait for DB (helps match the required demo scenario where app waits for db healthcheck)
	db := mustConnectWithRetry(cfg, 60*time.Second)
	defer db.Close()

	app := &App{db: db}

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(30 * time.Second))
	r.Use(middleware.Logger)

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	r.Route("/movies", func(r chi.Router) {
		r.Get("/", app.getMovies)
		r.Post("/", app.createMovie)

		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", app.getMovie)
			r.Put("/", app.updateMovie)
			r.Delete("/", app.deleteMovie)
		})
	})

	srv := &http.Server{
		Addr:              ":" + cfg.AppPort,
		Handler:           r,
		ReadHeaderTimeout: 10 * time.Second,
	}

	go func() {
		log.Printf("Starting the Server on :%s...", cfg.AppPort)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("server error: %v", err)
		}
	}()

	// BONUS (ADVANCED): graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down gracefully...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_ = srv.Shutdown(ctx)
	_ = db.Close()

	log.Println("Bye.")
}

func loadConfig() Config {
	// NO hardcoded creds: everything from env / compose (required in PDF)
	return Config{
		AppPort:   getEnv("APP_PORT", "8080"),
		DBHost:    getEnv("DB_HOST", "db"), // IMPORTANT: connect by service name "db", not localhost
		DBPort:    getEnv("DB_PORT", "5432"),
		DBUser:    getEnv("DB_USER", "postgres"),
		DBPass:    getEnv("DB_PASSWORD", "postgres"),
		DBName:    getEnv("DB_NAME", "appdb"),
		DBSSLMode: getEnv("DB_SSLMODE", "disable"),
	}
}

func getEnv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}

func (c Config) dsn() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.DBHost, c.DBPort, c.DBUser, c.DBPass, c.DBName, c.DBSSLMode,
	)
}

func mustConnectWithRetry(cfg Config, timeout time.Duration) *sql.DB {
	deadline := time.Now().Add(timeout)
	var lastErr error

	for time.Now().Before(deadline) {
		db, err := sql.Open("postgres", cfg.dsn())
		if err != nil {
			lastErr = err
			log.Printf("Waiting for database (open error): %v", err)
			time.Sleep(1 * time.Second)
			continue
		}

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		err = db.PingContext(ctx)
		cancel()
		if err == nil {
			db.SetMaxOpenConns(10)
			db.SetMaxIdleConns(10)
			db.SetConnMaxLifetime(30 * time.Minute)
			return db
		}

		lastErr = err
		log.Printf("Waiting for database (ping): %v", err)
		time.Sleep(1 * time.Second)
	}

	log.Fatalf("DB not ready after %s: %v", timeout, lastErr)
	return nil
}

// ---- handlers ----

func (a *App) getMovies(w http.ResponseWriter, r *http.Request) {
	rows, err := a.db.QueryContext(r.Context(), `
		SELECT id, title, genre, budget, hero, heroine
		FROM movies
		ORDER BY id ASC
	`)
	if err != nil {
		httpError(w, err, http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var out []Movie
	for rows.Next() {
		var m Movie
		if err := rows.Scan(&m.ID, &m.Title, &m.Genre, &m.Budget, &m.Hero, &m.Heroine); err != nil {
			httpError(w, err, http.StatusInternalServerError)
			return
		}
		out = append(out, m)
	}
	writeJSON(w, http.StatusOK, out)
}

func (a *App) getMovie(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(chi.URLParam(r, "id"))
	if err != nil {
		httpError(w, err, http.StatusBadRequest)
		return
	}

	var m Movie
	err = a.db.QueryRowContext(r.Context(), `
		SELECT id, title, genre, budget, hero, heroine
		FROM movies
		WHERE id = $1
	`, id).Scan(&m.ID, &m.Title, &m.Genre, &m.Budget, &m.Hero, &m.Heroine)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			httpError(w, errors.New("not found"), http.StatusNotFound)
			return
		}
		httpError(w, err, http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, m)
}

func (a *App) createMovie(w http.ResponseWriter, r *http.Request) {
	var in Movie
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		httpError(w, err, http.StatusBadRequest)
		return
	}
	if in.Title == "" || in.Genre == "" {
		httpError(w, errors.New("title and genre are required"), http.StatusBadRequest)
		return
	}

	err := a.db.QueryRowContext(r.Context(), `
		INSERT INTO movies (title, genre, budget, hero, heroine)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`, in.Title, in.Genre, in.Budget, in.Hero, in.Heroine).Scan(&in.ID)
	if err != nil {
		httpError(w, err, http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusCreated, in)
}

func (a *App) updateMovie(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(chi.URLParam(r, "id"))
	if err != nil {
		httpError(w, err, http.StatusBadRequest)
		return
	}

	var in Movie
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		httpError(w, err, http.StatusBadRequest)
		return
	}

	res, err := a.db.ExecContext(r.Context(), `
		UPDATE movies
		SET title=$1, genre=$2, budget=$3, hero=$4, heroine=$5
		WHERE id=$6
	`, in.Title, in.Genre, in.Budget, in.Hero, in.Heroine, id)
	if err != nil {
		httpError(w, err, http.StatusInternalServerError)
		return
	}

	n, _ := res.RowsAffected()
	if n == 0 {
		httpError(w, errors.New("not found"), http.StatusNotFound)
		return
	}

	in.ID = id
	writeJSON(w, http.StatusOK, in)
}

func (a *App) deleteMovie(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(chi.URLParam(r, "id"))
	if err != nil {
		httpError(w, err, http.StatusBadRequest)
		return
	}

	res, err := a.db.ExecContext(r.Context(), `DELETE FROM movies WHERE id=$1`, id)
	if err != nil {
		httpError(w, err, http.StatusInternalServerError)
		return
	}

	n, _ := res.RowsAffected()
	if n == 0 {
		httpError(w, errors.New("not found"), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func parseID(s string) (int64, error) {
	id, err := strconv.ParseInt(s, 10, 64)
	if err != nil || id <= 0 {
		return 0, errors.New("invalid id")
	}
	return id, nil
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func httpError(w http.ResponseWriter, err error, status int) {
	writeJSON(w, status, map[string]any{"error": err.Error()})
}

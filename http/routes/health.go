package routes

import (
	"context"
	"encoding/json"

	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/ong-gtp/go-chat/config"
	"github.com/ong-gtp/go-chat/services/rabbitmq"
)

type healthResponse struct {
	Status string            `json:"status"`
	Checks map[string]string `json:"checks,omitempty"`
}

var RegisterHealthRoutes = func(router *mux.Router) {
	router.HandleFunc("/health/live", livenessHandler).Methods(http.MethodGet)
	router.HandleFunc("/health/ready", readinessHandler).Methods(http.MethodGet)
}

func livenessHandler(w http.ResponseWriter, r *http.Request) {
	writeHealthResponse(w, http.StatusOK, healthResponse{
		Status: "ok",
	})
}

func readinessHandler(w http.ResponseWriter, r *http.Request) {
	checks := make(map[string]string)

	dbReady := checkDatabase(r.Context())
	checks["database"] = statusString(dbReady)

	rabbitMQReady := rabbitmq.GetRabbitMQBroker().Ready()
	checks["rabbitmq"] = statusString(rabbitMQReady)

	if !dbReady || !rabbitMQReady {
		writeHealthResponse(w, http.StatusServiceUnavailable, healthResponse{
			Status: "not_ready",
			Checks: checks,
		})
		return
	}

	writeHealthResponse(w, http.StatusOK, healthResponse{
		Status: "ready",
		Checks: checks,
	})
}

func checkDatabase(parent context.Context) bool {
	db := config.GetDB()

	if db == nil {
		log.Println("database readiness check failed: database is nil")
		return false
	}

	sqlDB, err := db.DB()

	if err != nil {
		log.Printf("database readiness check failed getting sql DB: %v", err)
		return false
	}

	ctx, cancel := context.WithTimeout(parent, 6*time.Second)
	defer cancel()

	if err := sqlDB.PingContext(ctx); err != nil {
		log.Printf("database readiness check failed: %v", err)
		return false
	}

	return true
}

func statusString(ready bool) string {
	if ready {
		return "ok"
	}
	return "failed"
}

func writeHealthResponse(w http.ResponseWriter, status int, response healthResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		return
	}
}

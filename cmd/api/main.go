package main

import (
	"golang-stady/internal/db"
	"golang-stady/internal/handlers"
	"log"
	"net/http"
	"os"
)

func main() {
	/*cfg, err := config.BuildConfig()
	if err != nil {
		log.Fatal(err)
	}*/
	dbUrl := os.Getenv("APP_DB_URL")

	serverPort := os.Getenv("APP_HTTP_PORT")

	connect, err := db.Connect(dbUrl)
	if err != nil {
		log.Fatal(err)
	}

	defer connect.Close()

	taskStore := db.NewTaskStore(connect)
	handler := handlers.NewHandlers(taskStore)

	mux := http.NewServeMux()

	mux.HandleFunc("/tasks", methodHandler(handler.GetAllTasks, http.MethodGet))
	mux.HandleFunc("/tasks/create", methodHandler(handler.CreateTask, http.MethodPost))
	mux.HandleFunc("/tasks/", taskIdHandler(handler))
	mux.HandleFunc("/tasks/", taskIdHandler(handler))

	loggedMux := loggingMiddleware(mux)

	serverAddr := ":" + serverPort
	err = http.ListenAndServe(serverAddr, loggedMux)
	if err != nil {
		log.Fatal(err)
	}
}

func methodHandler(handlerFunc http.HandlerFunc, allowedMethod string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != allowedMethod {
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		}
		handlerFunc(w, r)
	}
}

func taskIdHandler(handler *handlers.Handlers) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handler.GetById(w, r)
		case http.MethodPut:
			handler.UpdateCompleted(w, r)
		case http.MethodDelete:
			handler.Delete(w, r)
		default:
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		}
	}
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s", r.Method, r.URL.Path, r.RemoteAddr)
		next.ServeHTTP(w, r)
	})
}

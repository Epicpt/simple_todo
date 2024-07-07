package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	_ "modernc.org/sqlite"

	"main.go/database"
	"main.go/handlers"
	"main.go/tests"
)

func main() {
	// Проверка базы данных
	// Создание базы данных и таблицы, если не существует
	db, err := database.СheckDB()
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	repository := database.NewRepository(db)

	webDir := "./web"

	r := chi.NewRouter()
	fs := http.FileServer(http.Dir(webDir))
	r.Handle("/*", fs)

	// Регистрация обработчиков
	r.Get("/api/nextdate", handlers.HandleNextDate)
	r.Post("/api/task", handlers.HandleTaskPost(repository))
	r.Get("/api/tasks", handlers.HandleTaskGet(repository))
	r.Put("/api/task", handlers.HandleTaskPut(repository))
	r.Delete("/api/task", handlers.HandleTaskDelete(repository))
	r.Get("/api/task", handlers.HandleTaskID(repository))
	r.Post("/api/task/done", handlers.HandleTaskDone(repository))

	port := fmt.Sprintf(":%d", tests.Port)

	log.Printf("Сервер запущен. Порт %s\n", port)

	err = http.ListenAndServe(port, r)
	if err != nil {
		log.Fatal(err)
	}

}

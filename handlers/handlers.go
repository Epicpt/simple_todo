package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"main.go/api"
	"main.go/database"
	"main.go/model"
)

const layoutDate string = "20060102" // формат даты

func HandleNextDate(w http.ResponseWriter, r *http.Request) {

	nowRaw := r.URL.Query().Get("now")
	now, err := time.Parse(layoutDate, nowRaw)
	if err != nil {
		responseWithError(w, err)
	}
	date := r.URL.Query().Get("date")
	repeat := r.URL.Query().Get("repeat")

	result, err := api.NextDate(now, date, repeat, layoutDate)
	out := result
	if err != nil {
		out = err.Error()
		log.Printf("Error next date - %s", out)
	}

	_, err = w.Write([]byte(out))
	if err != nil {
		log.Println(err)
	}
}

func responseWithError(w http.ResponseWriter, err error) {
	log.Printf("%v\n", err)
	if errEncode := json.NewEncoder(w).Encode(model.ResponseError{Error: err.Error()}); errEncode != nil {
		log.Println(errEncode)
	}

}

func HandleTaskPost(repo database.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		var task model.Task
		err := json.NewDecoder(r.Body).Decode(&task)
		if err != nil {
			responseWithError(w, errors.New("Ошибка десериализации JSON"))
			return
		}

		if task.Title == "" {
			responseWithError(w, errors.New("Не указан заголовок задачи"))
			return
		}
		task.Date, err = api.GetNextDate(task, layoutDate)
		if err != nil {
			responseWithError(w, errors.New("Ошибка получения следующей даты"))
			return
		}

		id, err := repo.AddTask(task)
		if err != nil {
			responseWithError(w, errors.New("Ошибка добавления задачи"))
			return
		}
		response := map[string]interface{}{
			"id": id,
		}

		if errEncode := json.NewEncoder(w).Encode(response); errEncode != nil {
			log.Println(errEncode)
		}
	}

}

func HandleTaskGet(repo database.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		tasks, err := repo.GetTasks()
		if err != nil {
			responseWithError(w, errors.New("Ошибка получения задач"))
		}

		result := model.Tasks{Tasks: tasks}

		if tasks == nil {
			result = model.Tasks{Tasks: []model.Task{}}

		}

		if errEncode := json.NewEncoder(w).Encode(result); errEncode != nil {
			log.Println(errEncode)
		}
	}

}

func HandleTaskPut(repo database.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		var task model.Task

		err := json.NewDecoder(r.Body).Decode(&task)
		if err != nil {
			responseWithError(w, err)
			return
		}

		if task.ID == "" {
			responseWithError(w, errors.New("Не указан ID задачи"))
			return
		}

		if task.Title == "" {
			responseWithError(w, errors.New("Не указан заголовок задачи"))
			return
		}

		task.Date, err = api.GetNextDate(task, layoutDate)
		if err != nil {
			responseWithError(w, errors.New("Ошибка получения следующей даты"))
			return
		}

		if err = repo.UpdateTask(task); err != nil {
			responseWithError(w, err)
			return

		}

		if errEncode := json.NewEncoder(w).Encode(map[string]interface{}{}); errEncode != nil {
			log.Println(errEncode)
		}
	}

}

func HandleTaskDelete(repo database.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		id := r.URL.Query().Get("id")
		if id == "" {
			responseWithError(w, errors.New("Задача не найдена"))
			return
		}

		if err := repo.DeleteTask(id); err != nil {
			responseWithError(w, err)
			return

		}

		if errEncode := json.NewEncoder(w).Encode(map[string]interface{}{}); errEncode != nil {
			log.Println(errEncode)
		}
	}

}

func HandleTaskID(repo database.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")

		id := r.URL.Query().Get("id")
		if id == "" {
			responseWithError(w, errors.New("Задача не найдена"))
			return
		}
		task, err := repo.GetTask(id)
		if err != nil {
			responseWithError(w, errors.New("Ошибка получения задачи"))
			return
		}
		if errEncode := json.NewEncoder(w).Encode(task); errEncode != nil {
			log.Println(errEncode)
		}
	}
}

func HandleTaskDone(repo database.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")

		var task model.Task

		id := r.URL.Query().Get("id")

		if id == "" {
			responseWithError(w, errors.New("Задача не найдена"))
			return
		}
		task, err := repo.GetTask(id)
		if err != nil {
			responseWithError(w, errors.New("Ошибка получения задачи"))
			return
		}

		if task.Repeat != "" {

			nextDate, err := api.NextDate(time.Now(), task.Date, task.Repeat, layoutDate)
			if err != nil {
				responseWithError(w, err)
				return
			}
			task.Date = nextDate

			if err := repo.UpdateTaskDate(task); err != nil {
				responseWithError(w, err)
				return

			}

		} else {
			if err := repo.DeleteTask(id); err != nil {
				responseWithError(w, err)
				return

			}

		}
		if errEncode := json.NewEncoder(w).Encode(map[string]interface{}{}); errEncode != nil {
			log.Println(errEncode)
		}
	}
}

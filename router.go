// Запуск сервера и RESTful API.
// В случае неудачной обработки запроса, вернуть ошибку в формате RFC 7807.

package main

import (
	"bytes"
    "context"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"
    "os"
    "os/signal"
    "syscall"
	"strconv"

	"github.com/gorilla/mux"
)

const (
	BAD_REQUEST      = "Неправильный формат запроса"
	SERVER_ERROR     = "Ошибка сервера"
	NO_CONTENT_ERROR = "Записей не найдено"
	UNKNOWN_ERROR    = "Неизвестный статус ошибки"

    ADDRESS = "localhost"
    API_PATH = "/api/v1/employees"
	PORT      = ":33890"
	TYPE_JSON = "application/json"
    SERVER_URL = "http://" + ADDRESS + PORT

	INFO = `{"name": "employees", "version":"1.0.0"}`
)

// Структура сервер
type Router struct {
	M Model
	R *mux.Router
    Srv http.Server
}

// Установить функции для обработки запросов.
func (r *Router) MakeRoutes() {
	r.R.HandleFunc(API_PATH, r.AddEmployee).Methods("POST")
	r.R.HandleFunc(API_PATH, r.GetAllEmployees).Methods("GET")
	r.R.HandleFunc(API_PATH + "/{id}", r.GetEmployeeById).Methods("GET")
	r.R.HandleFunc(API_PATH + "/{id}", r.RemoveEmployee).Methods("DELETE")
	r.R.HandleFunc(API_PATH, r.UpdateEmployee).Methods("PUT")
	r.R.HandleFunc("/tech/info", r.GetTechInfo).Methods("GET")
}

// Запустить сервер
func (r *Router) Start() {
	r.MakeRoutes()
    go func() {
        shutdownCh := make(chan os.Signal, 1)
        signal.Notify(shutdownCh, syscall.SIGINT, syscall.SIGTERM)
        <-shutdownCh
        close(shutdownCh)
        if err := r.Shutdown(context.Background()); err != nil {
            log.Fatal(err)
        }
    }()
    log.Printf("Запуск сервера на %s%s\n", ADDRESS, PORT)
	log.Fatal(r.Srv.ListenAndServe())
}

// Завершить работу сервера. Закрыть переменную модель.
func (r *Router) Shutdown(ctx context.Context) error {
    r.M.Close()
    return r.Srv.Shutdown(ctx)
}

// Обработать запрос добавления данных в модель.
// В случае успеха вернуть код 201.
func (r *Router) AddEmployee(w http.ResponseWriter, req *http.Request) {
	// 1. Проверить запрос
	if !isJson(req) {
		writeErr(w, http.StatusBadRequest, "тип запроса не JSON")
		return
	}
	// 2. Выполнить запрос
	var buf bytes.Buffer
	n, err := io.Copy(&buf, req.Body)
    if err != nil {
		writeErr(w, http.StatusInternalServerError, err.Error())
		return
	} else if n == 0 {
		writeErr(w, http.StatusBadRequest, ("тело не задано"))
		return
    }

	if err := r.M.Add(buf.String(), req.Context()); err != nil {
		writeErr(w, http.StatusInternalServerError, err.Error())
		return
	}
	// 3. Вернуть код 201
	w.WriteHeader(http.StatusCreated)
}

// Обработать запрос чтения всех данных из модели.
// Логировать ошибку записи ответа.
// В случае успеха вернуть данные в формате JSON.
func (r *Router) GetAllEmployees(w http.ResponseWriter, req *http.Request) {
	// 1. Сделать запрос
	data, err := r.M.GetAll(req.Context())
	if err != nil {
		writeErr(w, http.StatusInternalServerError, err.Error())
		return
	} else if data == nil {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	// 2. Перекодировать результат запроса в JSON или XML и вернуть клиенту
    var convData []byte
    switch req.Header.Get("Accept") {
    case "application/xml": // кодировать в XML
        convData, err = xml.Marshal(data)
        if err != nil {
            writeErr(w, http.StatusInternalServerError, err.Error())
            return
        }
        w.Header().Set("Content-Type", "application/xml")
    default: // кодировать в JSON
        convData, err = json.Marshal(data)
        if err != nil {
            writeErr(w, http.StatusInternalServerError, err.Error())
            return
        }
        w.Header().Set("Content-Type", "application/json")
    }
	_, err = w.Write(convData)
	if err != nil {
		log.Printf("Write: %v\n", err)
		return
	}
}

// Запросить данные из модели с переданным аргументом id.
// Логировать ошибку записи ответа.
// В случае успеха вернуть данные в формате JSON.
func (r *Router) GetEmployeeById(w http.ResponseWriter, req *http.Request) {
	// 1. Проверить переданные данные
	vars := mux.Vars(req)
	idstr, ok := vars["id"]
	if !ok {
		writeErr(w, http.StatusBadRequest, "отсутствует параметр id")
		return
	}
	id, err := strconv.Atoi(idstr)
	if err != nil {
		writeErr(w, http.StatusBadRequest, "плохой параметр id")
		return
	}
	// 2. Сделать запрос
	data, err := r.M.GetId(id, req.Context())
	if err != nil {
		writeErr(w, http.StatusInternalServerError, err.Error())
		return
	} else if data == nil {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	// 3. Перекодировать запрос в JSON и вернуть клиенту
	datajson, err := json.Marshal(data)
	if err != nil {
		writeErr(w, http.StatusInternalServerError, err.Error())
		return
	}
	_, err = w.Write(datajson)
	if err != nil {
		log.Printf("Write: %v\n", err)
		return
	}
}

// Вернуть данные о программе
func (r *Router) GetTechInfo(w http.ResponseWriter, req *http.Request) {
	_, err := w.Write([]byte(INFO))
	if err != nil {
		log.Printf("Write: %v\n", err)
		return
	}
}

// Удалить данные из модели.
// В случае успеха вернуть код 200.
func (r *Router) RemoveEmployee(w http.ResponseWriter, req *http.Request) {
	// 1. Проверить переданные данные
	vars := mux.Vars(req)
	idstr, ok := vars["id"]
	if !ok {
		writeErr(w, http.StatusBadRequest, "отсутствует параметр id")
		return
	}
	id, err := strconv.Atoi(idstr)
	if err != nil {
		writeErr(w, http.StatusBadRequest, "плохой параметр id")
		return
	}
	// 2. Сделать запрос
	if err := r.M.Remove(id, req.Context()); err != nil {
		writeErr(w, http.StatusInternalServerError, err.Error())
		return
	}
}

// Заменить данные модели.
// В случае успеха вернуть код 200.
func (r *Router) UpdateEmployee(w http.ResponseWriter, req *http.Request) {
	// 1. Проверить запрос
	if !isJson(req) {
		writeErr(w, http.StatusBadRequest, "тип запроса не JSON")
		return
	}
	// 2. Выполнить запрос
	var buf bytes.Buffer
	n, err := io.Copy(&buf, req.Body)
    if err != nil {
		writeErr(w, http.StatusInternalServerError, err.Error())
		return
	} else if n == 0 {
		writeErr(w, http.StatusBadRequest, ("тело не задано"))
		return
    }

	if err := r.M.Update(buf.String(), req.Context()); err != nil {
		writeErr(w, http.StatusInternalServerError, err.Error())
		return
	}
}


// Создать новую переменную Router.
func NewRouter() (*Router, error) {
    var r Router
	m, err := NewModel()
	if err != nil {
		return nil, fmt.Errorf("NewModel: %v\n", err)
	}
    r.M = m
    r.R = mux.NewRouter()
    r.Srv.Addr = ADDRESS + PORT
    r.Srv.Handler = r.R

	return &r, nil
}

// Проверить content-type запроса.
// Вернуть true, если это - JSON. Иначе вернуть false.
func isJson(r *http.Request) bool {
	contentType := r.Header.Get("Content-Type")
	return contentType == TYPE_JSON
}

// Структура ошибки в формате RFC 7807.
type ResponseError struct {
	Status int    `json:"status"`
	Detail string `json:"detail"`
	Title  string `json:"title"`
}

// Записать ошибку в формате RFC 7807 в w.
func writeErr(w http.ResponseWriter, status int, detail string) {
	var respErr ResponseError

	// 1. Установить поля в соответствии с типом ошибки.
	switch status {
	case http.StatusBadRequest:
		respErr.Title = BAD_REQUEST
		w.WriteHeader(http.StatusBadRequest)
	case http.StatusInternalServerError:
		respErr.Title = SERVER_ERROR
		w.WriteHeader(http.StatusInternalServerError)
	default:
		log.Printf(UNKNOWN_ERROR)
		respErr.Title = UNKNOWN_ERROR
		w.WriteHeader(http.StatusInternalServerError)
	}
	respErr.Status = status
	respErr.Detail = detail

	// 2. Вернуть перекодированную в JSON ошибку клиенту.
	json, err := json.Marshal(respErr)
	if err != nil {
		log.Printf("json.Marshal: %v\n", err)
		return
	}
	_, err = w.Write(json)
	if err != nil {
		fmt.Printf("jjjjjjj: %d\n", status)
		log.Printf("Write: %v\n", err)
		return
	}
}

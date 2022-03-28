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
	"strconv"
	"syscall"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	BAD_REQUEST      = "Неправильный формат запроса"
	SERVER_ERROR     = "Ошибка сервера"
	NO_CONTENT_ERROR = "Записей не найдено"
	UNKNOWN_ERROR    = "Неизвестный статус ошибки"

	API_EMPL_PATH = "/api/v1/employees"
	API_ERR_PATH  = "/api/v1/error"
	API_INFO_PATH = "/tech/info"
	API_METRICS   = "/metrics"

	TYPE_JSON = "application/json"
)

type Router struct {
	M   Model
	R   *mux.Router
	Srv http.Server
}

var httpDuration = promauto.NewHistogram(prometheus.HistogramOpts{
	Name: "http_response_time_seconds",
	Help: "Продолжительность HTTP запросов.",
})

// Установить функции для обработки запросов.
func (r *Router) MakeRoutes() {
	r.R.HandleFunc(API_ERR_PATH+"/{id}", r.GetErr).Methods("GET")
	r.R.HandleFunc(API_EMPL_PATH+"/{id}", r.FireEmployee).Methods("DELETE")
	r.R.HandleFunc(API_EMPL_PATH, accept(r.GetAllEmployees)).Methods("GET")
	r.R.HandleFunc(API_EMPL_PATH+"/{id}", r.GetEmployee).Methods("GET")
	r.R.HandleFunc(API_EMPL_PATH, r.HireEmployee).Methods("POST")
	r.R.HandleFunc(API_EMPL_PATH, r.UpdateEmployee).Methods("PUT")
	r.R.HandleFunc(API_INFO_PATH, r.GetTechInfo).Methods("GET")
	r.R.Handle(API_METRICS, promhttp.Handler()).Methods("GET")
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
	log.Printf("Запуск сервера на %s:%s\n", cfg.Prog.Addr, cfg.Prog.Port)
	log.Fatal(r.Srv.ListenAndServe())
}

// Завершить работу сервера. Закрыть переменную модель.
func (r *Router) Shutdown(ctx context.Context) error {
	r.M.Close()
	return r.Srv.Shutdown(ctx)
}

// Удалить данные из модели.
// В случае успеха вернуть код 200.
func (r *Router) FireEmployee(w http.ResponseWriter, req *http.Request) {
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
	if err := r.M.FireEmployee(id, req.Context()); err != nil {
		writeDbErr(w, err)
		return
	}
}

// Вернуть middleware функцию.
// Убедиться, что Accept header запроса - XML или JSON.
// Иначе вернуть ошибку.
func accept(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Header.Get("Accept") {
		case "application/json": // OK
		case "application/xml": // OK
		default: // не OK
			writeErr(w, http.StatusBadRequest, "не XML/JSON Accept header")
			return
		}
		next(w, r)
	}
}

// Обработать запрос чтения всех данных из модели.
// Логировать ошибку записи ответа.
// В случае успеха вернуть данные в формате JSON или XML.
func (r *Router) GetAllEmployees(w http.ResponseWriter, req *http.Request) {
	// 1. Сделать запрос.
	data, err := r.M.GetAllEmployees(req.Context())
	if err != nil {
		writeDbErr(w, err)
		return
	}
	// 2. Перекодировать запрос в формат, соответствующий header Accept.
	var buf bytes.Buffer // конвертированные данные
	switch req.Header.Get("Accept") {
	case "application/json":
		if err := json.NewEncoder(&buf).Encode(data); err != nil {
			writeErr(w, http.StatusInternalServerError, err.Error())
			return
		}
		w.Header().Set("Content-Type", "application/json")
	case "application/xml":
		if err := xml.NewEncoder(&buf).Encode(data); err != nil {
			writeErr(w, http.StatusInternalServerError, err.Error())
			return
		}
		w.Header().Set("Content-Type", "application/xml")
	default:
		writeErr(w, http.StatusBadRequest, "не XML/JSON accept header")
		return
	}
	// 3. Записать данные клиенту
	_, err = w.Write(buf.Bytes())
	if err != nil {
		log.Printf("Write: %v\n", err)
		return
	}
}

// Используется только для тестов.
// Вернуть пользовательскую ошибку на id = 1.
// Вернуть ошибку базы данных на id = 2.
// Иначе, вернуть nil.
func (r *Router) GetErr(w http.ResponseWriter, req *http.Request) {
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
	err = r.M.GetErr(id)
	if err != nil {
		writeDbErr(w, err)
		return
	}
}

// Запросить данные из модели с переданным аргументом id.
// Логировать ошибку записи ответа.
// В случае успеха вернуть данные в формате JSON.
func (r *Router) GetEmployee(w http.ResponseWriter, req *http.Request) {
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
	data, err := r.M.GetEmployee(id, req.Context())
	if err != nil {
		writeDbErr(w, err)
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

// Вернуть данные о программе.
func (r *Router) GetTechInfo(w http.ResponseWriter, req *http.Request) {
	_, err := w.Write([]byte(cfg.Prog.Info))
	if err != nil {
		log.Printf("Write: %v\n", err)
		return
	}
}

// Обработать запрос добавления данных в модель.
// В случае успеха вернуть код 201.
func (r *Router) HireEmployee(w http.ResponseWriter, req *http.Request) {
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

	if err := r.M.HireEmployee(buf.String(), req.Context()); err != nil {
		writeDbErr(w, err)
		return
	}
	// 3. Вернуть код 201
	w.WriteHeader(http.StatusCreated)
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

	if err := r.M.UpdateEmployee(buf.String(), req.Context()); err != nil {
		writeDbErr(w, err)
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
	r.Srv.Addr = cfg.Prog.Addr + ":" + cfg.Prog.Port
	r.Srv.Handler = r.R

	r.R.Use(durationMiddleware)

	return &r, nil
}

// Измерить длительность удовлетворения HTTP запроса
func durationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		timer := prometheus.NewTimer(httpDuration)
		next.ServeHTTP(w, r)
		timer.ObserveDuration()
	})
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

// Вызвать writeErr с аргументами зависящими от типа ошибки базы данных.
func writeDbErr(w http.ResponseWriter, err error) {
	if err == dbNA {
		writeErr(w, http.StatusInternalServerError, dbNA.Error())
	} else {
		writeErr(w, http.StatusBadRequest, err.Error())
	}
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

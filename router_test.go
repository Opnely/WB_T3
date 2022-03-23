package main

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const (
	ADD       = "/api/v1/add"
	GET_ALL   = "/api/v1/get_all"
	GET_ID    = "/api/v1/get_id/"
	REMOVE    = "/api/v1/remove/"
	UPDATE    = "/api/v1/update"
	TECH_INFO = "/tech/info"

    EMPLOYEES_URL = SERVER_URL + API_PATH
)

// Запустить сервер
func TestMain(m *testing.M) {
	router, err := NewRouter()
	if err != nil {
		log.Fatalf("NewRouter: %v\n", err)
	}
	go router.Start()
	log.Println("Ожидание запуска сервера")
	time.Sleep(2 * time.Second)

	os.Exit(m.Run())
}

// Тесты запроса добавления записи.
func TestRouterAdd(t *testing.T) {
	t1 := `{"name":"A", "last_name":"Al", "patrnonymic":"Ap", "phone":"X",` +
		`"position":"a", "good_job_count":1}`
	var tests = []struct {
		json, content string
		status        int
	}{
		{"", "none", http.StatusBadRequest},                      // не json
		{"", "application/json", http.StatusBadRequest}, // "" JSON
		{t1, "application/json", http.StatusCreated},             // успех
	}
	assert := assert.New(t)
	for i, test := range tests {
		// 1. Сделать запрос
		buf := bytes.NewBuffer([]byte(test.json))
        req := &Request{Method: "POST", URL: EMPLOYEES_URL, Body: buf }
        req.ContentType = test.content
		resp, err := makeRequest(req)
		if err != nil {
			t.Logf("makeRequest: %v\n", err)
			continue
		}
		resp.Body.Close()

		// 2. Проверить результат
		assert.Equal(test.status, resp.StatusCode, "Тест %d\n", i)
	}
}

// Тесты запроса считывания всех записей.
// Ожидается, что число элементов в базе данных больше нуля.
func TestRouterGetAll(t *testing.T) {
	var tests = []struct{acceptHeader string; status int}{
		{"", http.StatusBadRequest}, // отсутствует accept Header
		{"application/json", http.StatusOK}, // успех
		{"application/xml", http.StatusOK}, // успех
	}
	assert := assert.New(t)
	for i, test := range tests {
        req := &Request{Method: "GET", URL: EMPLOYEES_URL}
        req.AcceptHeader = test.acceptHeader
		resp, err := makeRequest(req)
		if err != nil {
			t.Logf("makeRequest: %v\n", err)
			continue
		}
		// 1. Проверить код возврата в результате
		assert.Equal(test.status, resp.StatusCode, "Тест %d\n", i)
		if resp.StatusCode != http.StatusOK {
			continue
		}
		// 2. Проверить результат
		var buf bytes.Buffer
		if _, err := io.Copy(&buf, resp.Body); err != nil {
			t.Logf("io.Copy: %v\n", err)
			return
		}
		resp.Body.Close()
		var d []Data
        switch v := resp.Header.Get("Content-Type"); v {
        case "application/json":
            if err := json.Unmarshal(buf.Bytes(), &d); err != nil {
                t.Logf("json.Unmarshal: %v\n", err)
                continue
            }
        case "application/xml":
            if err := xml.Unmarshal(buf.Bytes(), &d); err != nil {
                t.Logf("xml.Unmarshal: %v\n", err)
                continue
            }
        default:
            t.Logf("Неожиданный http.Response Content-Type: %v\n", v)
            continue
        }
		assert.Less(0, len(d), "Тест %d\n", i)
	}
}

// Тесты запроса считывания записи по id.
// Проверить результат только по id.
func TestRouterGetId(t *testing.T) {
	var tests = []struct {
		id     string
		status int
	}{
		{"", http.StatusNotFound},    // в запросе отсутствует id
		{"x", http.StatusBadRequest}, // плохой id
		{"0", http.StatusNoContent},  // запись с id не найдена
		{"3", http.StatusOK},         // успех
	}
	assert := assert.New(t)
	for i, test := range tests {
        req := &Request{Method: "GET", URL: EMPLOYEES_URL + "/" + test.id}
		resp, err := makeRequest(req)
		if err != nil {
			t.Logf("makeRequest: %v\n", err)
			continue
		}
		// 1. Проверить statusCode в результате
		assert.Equal(test.status, resp.StatusCode, "Тест %d\n", i)
		if resp.StatusCode != http.StatusOK {
			continue
		}
		// 2. Проверить id в результате
		var buf bytes.Buffer
		io.Copy(&buf, resp.Body)
		resp.Body.Close()
		var d Data
		if err := json.Unmarshal(buf.Bytes(), &d); err != nil {
			t.Logf("json.Unmarshal: %v\n", err)
			continue
		}
		assert.Equal(test.id, strconv.Itoa(d.Id), "Тест %d\n", i)
	}
}

// Тесты запроса удаления записи по id.
// Перед ожидаемым успешным удалением, найти запись с самым большим id.
func TestRouterRemove(t *testing.T) {
	var tests = []struct {
		id     string
		status int
	}{
		{"a", http.StatusBadRequest},             // плохой id
		{"9999", http.StatusInternalServerError}, // id не найден
		{"last", http.StatusOK},                  // успех
	}
	assert := assert.New(t)
	for i, test := range tests {
		if test.id == "last" { // найти последний id и удалить
			id, err := getLastId()
			if err != nil {
				t.Logf("getLastId: %v\n", err)
				continue
			}
			test.id = strconv.Itoa(id)
		}
        req := &Request{Method: "DELETE", URL: EMPLOYEES_URL + "/" + test.id}
		resp, err := makeRequest(req)
		if err != nil {
			t.Logf("makeRequest: %v\n", err)
			continue
		}
		resp.Body.Close()
		assert.Equal(test.status, resp.StatusCode, "Тест %d\n", i)
	}
}

// Тесты запроса обновления записи.
func TestRouterUpdate(t *testing.T) {
	t1 := `{"name":"A", "last_name":"Al", "patrnonymic":"Ap", "phone":"X",` +
		`"position":"a", "good_job_count":1}`
	t2 := `{"name":"B", "last_name":"Bl", "patrnonymic":"Bp", "phone":"Y",` +
		`"position":"b", "good_job_count":2, "id":4 }`
	var tests = []struct {
		json, content string
		status        int
	}{
		{"", "none", http.StatusBadRequest},                      // не json
		{"", "application/json", http.StatusBadRequest}, // "" JSON
		{t1, "application/json", http.StatusInternalServerError}, // нет id
		{t2, "application/json", http.StatusOK},                  // успех
	}
	assert := assert.New(t)
	for i, test := range tests {
		buf := bytes.NewBuffer([]byte(test.json))
        req := &Request{Method: "PUT", URL: EMPLOYEES_URL, Body: buf}
        req.ContentType = test.content
		resp, err := makeRequest(req)
		if err != nil {
			t.Logf("makeRequest: %v\n", err)
			continue
		}
		resp.Body.Close()
		assert.Equal(test.status, resp.StatusCode, "Тест %d\n", i)
	}
}

// Тест пяти запросов подряд: добавление записи, изменение записи, проверка
// записи и удаление записи. Считывание всех записей происходит после
// добавления записи с целью установления её id.
func TestRouterAll(t *testing.T) {
	t1 := `{"name":"A", "last_name":"Al", "patrnonymic":"Ap", "phone":"X",` +
		`"position":"a", "good_job_count":1}`
	t2 := `{"name":"B", "last_name":"Bl", "patrnonymic":"Bp", "phone":"Y",` +
		`"position":"b", "good_job_count":2, ` // "id": }`

	// 1. Создать новую запись.
	buf := bytes.NewBuffer([]byte(t1))
    req := &Request{Method: "POST", URL: EMPLOYEES_URL, Body: buf}
    req.ContentType = "application/json"
	resp, err := makeRequest(req)
	if err != nil {
		t.Logf("makeRequest: %v\n", err)
		return
	}
	resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		t.Logf("StatusCode: %v\n", resp.StatusCode)
		return
	}
	// 2. Найти id последней записи
	id, err := getLastId()
	if err != nil {
		t.Logf("getLastId: %v\n", err)
		return
	}
	// 3. Изменить последнюю запись.
	idstr := strconv.Itoa(id)
	if idstr == "" {
		t.Logf("strconv.Itoa(%d): нулевой результат\n", id)
		return
	}
	t2 += `"id":` + idstr + "}"
	buf3 := bytes.NewBuffer([]byte(t2))
    req = &Request{Method: "PUT", URL: EMPLOYEES_URL, Body: buf3}
    req.ContentType = "application/json"
	resp, err = makeRequest(req)
	if err != nil {
		t.Logf("makeRequest: %v\n", err)
		return
	}
	resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Logf("StatusCode: %v\n", resp.StatusCode)
		return
	}

	// 4. Проверить изменение записи.
    req = &Request{Method: "GET", URL: EMPLOYEES_URL+"/"+idstr, Body: nil}
    req.ContentType = ""
	resp, err = makeRequest(req)
	if err != nil {
		t.Logf("makeRequest: %v\n", err)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Logf("StatusCode: %v\n", resp.StatusCode)
		return
	}
	decoder := json.NewDecoder(resp.Body)
	var orig, fresh Data
	if err := decoder.Decode(&fresh); err != nil {
		t.Logf("Decode: %v\n", err)
		return
	}
	if err := json.Unmarshal([]byte(t2), &orig); err != nil {
		t.Logf("json.Unmarshal: %v\n", err)
		return
	}
	if orig != fresh {
		t.Logf("%v != %v\n", orig, fresh)
	}

	// 5. Удалить запись
    req = &Request{Method: "DELETE", URL: EMPLOYEES_URL+"/"+idstr, Body: nil}
    req.ContentType = ""
	resp2, err := makeRequest(req)
	if err != nil {
		t.Logf("makeRequest: %v\n", err)
		return
	}
	resp2.Body.Close()
	if resp2.StatusCode != http.StatusOK {
		t.Logf("StatusCode: %v\n", resp.StatusCode)
		return
	}
}

type Request struct {
    Body io.Reader
    AcceptHeader string
    ContentType string
    Method string
    URL string
}

// Сделать запрос.
// Вернуть адрес перемнной ответа и ошибку.
func makeRequest(r *Request) (*http.Response, error) {
	req, err := http.NewRequest(r.Method, r.URL, r.Body)
	if err != nil {
		return nil, fmt.Errorf("http.NewRequest: %v", err)
	}
	req.Header.Set("Content-Type", r.ContentType)
	req.Header.Set("Accept-Encoding", r.AcceptHeader)
	client := &http.Client{}

	return client.Do(req)
}

// Вернуть последний добавленный id.
// Считать все записи, выбрать последний добавленный id методом сравнения.
func getLastId() (int, error) {
	// 1. Получить все записи.
    req := &Request{Method: "GET", URL: EMPLOYEES_URL}
    req.AcceptHeader = "application/json"
	resp, err := makeRequest(req)
	if err != nil {
		return 0, fmt.Errorf("makeRequest: %v\n", err)
	}
	//req, err := http.NewRequest("GET", EMPLOYEES_URL, nil)
	//if err != nil {
	//	return 0, fmt.Errorf("http.NewRequest: %v", err)
	//}
	//client := &http.Client{}
	//resp, err := client.Do(req)
	//if err != nil {
	//	return 0, fmt.Errorf("Do: %v", err)
	//}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("StatusCode: %v", resp.StatusCode)
	}

	var buf2 bytes.Buffer
	if _, err := io.Copy(&buf2, resp.Body); err != nil {
		return 0, fmt.Errorf("io.Copy: %v", err)
	}

	var data []Data
	if err := json.Unmarshal(buf2.Bytes(), &data); err != nil {
		return 0, fmt.Errorf("json.Unmarshal: %v", err)
	}

	// 2. Найти id последней записи.
	var id int
	for _, d := range data {
		if d.Id > id {
			id = d.Id
		}
	}
	return id, nil
}

func TestHandleTechInfo(t *testing.T) {
	assert := assert.New(t)
    req := &Request{Method: "GET", URL: SERVER_URL + TECH_INFO}
	resp, err := makeRequest(req)
	if err != nil {
		t.Logf("makeRequest: %v\n", err)
		return
	}
	// 1. Проверить statusCode в результате
	if resp.StatusCode != http.StatusOK {
		t.Logf("StatusCode: %v\n", resp.StatusCode)
		return
	}
	// 2. Проверить результат
	var buf bytes.Buffer
	_, err = io.Copy(&buf, resp.Body)
	if err != nil {
		t.Logf("io.Copy: %v\n", err)
		return
	}
	resp.Body.Close()
	assert.Equal(INFO, buf.String(), "%q != %q\n", INFO, buf.String())
}

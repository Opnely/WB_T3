package service

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

// Тип запроса для функции makeRequest
type Request struct {
	Body         io.Reader
	AcceptHeader string
	ContentType  string
	Method       string
	URL          string
}

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
func TestRouterAddEmployee(t *testing.T) {
	t1 := `{"name":"A", "last_name":"Al", "patrnonymic":"Ap", "phone":"X",` +
		`"position":"a", "good_job_count":1}`
	var tests = []struct {
		json, content string
		status        int
	}{
		{"", "none", http.StatusBadRequest},             // не json
		{"", "application/json", http.StatusBadRequest}, // "" JSON
		{t1, "application/json", http.StatusCreated},    // успех
	}
	assert := assert.New(t)
	for i, test := range tests {
		// 1. Сделать запрос
		buf := bytes.NewBuffer([]byte(test.json))
		req := &Request{Method: "POST", URL: cfg.Prog.EmplUrl, Body: buf}
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
func TestRouterGetAllEmployees(t *testing.T) {
	var tests = []struct {
		acceptHeader string
		status       int
	}{
		{"", http.StatusBadRequest},         // отсутствует accept Header
		{"application/json", http.StatusOK}, // успех
		{"application/xml", http.StatusOK},  // успех
	}
	assert := assert.New(t)
	for i, test := range tests {
		req := &Request{Method: "GET", URL: cfg.Prog.EmplUrl}
		req.AcceptHeader = test.acceptHeader
		resp, err := makeRequest(req)
		if err != nil {
			t.Logf("makeRequest: %v\n", err)
			continue
		}
		// 1. Проверить код возврата в результате
		assert.Equal(test.status, resp.StatusCode, "Тест %d\n", i)
		if resp.StatusCode != http.StatusOK {
			resp.Body.Close()
			continue
		}
		// 2. Проверить результат
		var buf bytes.Buffer
		if _, err := io.Copy(&buf, resp.Body); err != nil {
			t.Logf("io.Copy: %v\n", err)
			resp.Body.Close()
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
func TestRouterGetEmployee(t *testing.T) {
	var tests = []struct {
		id     string
		status int
	}{
		{"", http.StatusNotFound},    // в запросе отсутствует id
		{"x", http.StatusBadRequest}, // плохой id
		{"0", http.StatusBadRequest}, // запись с id не найдена
		{"3", http.StatusOK},         // успех
	}
	assert := assert.New(t)
	for i, test := range tests {
		req := &Request{Method: "GET", URL: cfg.Prog.EmplUrl + "/" + test.id}
		resp, err := makeRequest(req)
		if err != nil {
			t.Logf("makeRequest: %v\n", err)
			continue
		}
		// 1. Проверить statusCode в результате
		assert.Equal(test.status, resp.StatusCode, "Тест %d\n", i)
		if resp.StatusCode != http.StatusOK {
			resp.Body.Close()
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

// Тесты ошибок базы данных.
func TestRouterErrors(t *testing.T) {
	var tests = []struct {
		id     string
		status int
	}{
		{"", http.StatusNotFound},             // неизвестная страница
		{"1", http.StatusBadRequest},          // плохой запрос
		{"2", http.StatusInternalServerError}, // внутренняя ошибка
		{"3", http.StatusOK},                  // успех
	}
	assert := assert.New(t)
	for i, test := range tests {
		req := &Request{Method: "GET", URL: cfg.Prog.ErrUrl + "/" + test.id}
		resp, err := makeRequest(req)
		if err != nil {
			t.Logf("makeRequest: %v\n", err)
			continue
		}
		//io.Copy(os.Stdout, resp.Body)
		resp.Body.Close()
		assert.Equal(test.status, resp.StatusCode, "Тест %d\n", i)
	}
}

// Тесты запроса удаления записи по id.
// Перед ожидаемым успешным удалением, найти запись с самым большим id.
func TestRouterRemoveEmployee(t *testing.T) {
	var tests = []struct {
		id     string
		status int
	}{
		{"a", http.StatusBadRequest},    // плохой id
		{"9999", http.StatusBadRequest}, // id не найден
		{"last", http.StatusOK},         // успех
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
		req := &Request{Method: "DELETE", URL: cfg.Prog.EmplUrl + "/" + test.id}
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
func TestRouterUpdateEmployee(t *testing.T) {
	t1 := `{"name":"A", "last_name":"Al", "patrnonymic":"Ap", "phone":"X",` +
		`"position":"a", "good_job_count":1}`
	t2 := `{"name":"B", "last_name":"Bl", "patrnonymic":"Bp", "phone":"Y",` +
		`"position":"b", "good_job_count":2, "id":4 }`
	var tests = []struct {
		json, content string
		status        int
	}{
		{"", "none", http.StatusBadRequest},             // не json
		{"", "application/json", http.StatusBadRequest}, // "" JSON
		{t1, "application/json", http.StatusBadRequest}, // нет id
		{t2, "application/json", http.StatusOK},         // успех
	}
	assert := assert.New(t)
	for i, test := range tests {
		buf := bytes.NewBuffer([]byte(test.json))
		req := &Request{Method: "PUT", URL: cfg.Prog.EmplUrl, Body: buf}
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

// Тест пяти запросов подряд: добавление записи t1, изменение записи t1 на t2,
// проверка записи и удаление записи. Считывание всех записей происходит после
// добавления записи с целью установления её id.
func TestRouterAllMethods(t *testing.T) {
	t1 := `{"name":"A", "last_name":"Al", "patrnonymic":"Ap", "phone":"X",` +
		`"position":"a", "good_job_count":1}`
	t2 := `{"name":"B", "last_name":"Bl", "patrnonymic":"Bp", "phone":"Y",` +
		`"position":"b", "good_job_count":2, ` // "id": }`

	// 1. Создать новую запись.
	buf := bytes.NewBuffer([]byte(t1))
	req := &Request{Method: "POST", URL: cfg.Prog.EmplUrl, Body: buf}
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
	req = &Request{Method: "PUT", URL: cfg.Prog.EmplUrl, Body: buf3}
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
	req = &Request{Method: "GET", URL: cfg.Prog.EmplUrl + "/" + idstr, Body: nil}
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
		t.Logf("Orig: %d, %s, %s, %s, %s, %d, %s\n", orig.Id, orig.FirstName,
			orig.LastName, orig.MidName, orig.PhoneNum, orig.DoneJobs,
			orig.Position)
		t.Logf("Orig: %d, %s, %s, %s, %s, %d, %s\n", fresh.Id,
			fresh.FirstName, fresh.LastName, fresh.MidName, fresh.PhoneNum,
			fresh.DoneJobs, fresh.Position)
	}

	// 5. Удалить запись
	req = &Request{Method: "DELETE", URL: cfg.Prog.EmplUrl + "/" + idstr}
	req.Body = nil
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

// Сделать запрос.
// Вернуть адрес перемнной ответа и ошибку.
func makeRequest(r *Request) (*http.Response, error) {
	req, err := http.NewRequest(r.Method, r.URL, r.Body)
	if err != nil {
		return nil, fmt.Errorf("http.NewRequest: %v", err)
	}
	req.Header.Set("Content-Type", r.ContentType)
	req.Header.Set("Accept", r.AcceptHeader)
	client := &http.Client{}

	return client.Do(req)
}

// Вернуть последний добавленный id.
// Считать все записи, выбрать последний добавленный id методом сравнения.
func getLastId() (int, error) {
	// 1. Получить все записи.
	req := &Request{Method: "GET", URL: cfg.Prog.EmplUrl}
	req.AcceptHeader = "application/json"
	resp, err := makeRequest(req)
	if err != nil {
		return 0, fmt.Errorf("makeRequest: %v\n", err)
	}
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

func TestRouterTechInfo(t *testing.T) {
	assert := assert.New(t)
	req := &Request{Method: "GET", URL: cfg.Prog.ServerUrl + API_INFO_PATH}
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
	assert.Equal(cfg.Prog.Info, buf.String(), "%q != %q\n",
		cfg.Prog.Info, buf.String())
}

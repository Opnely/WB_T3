// Функции и тестовые структуры для тестирования model.go.
// Большая часть Mock-методов возвращает результат в зависимости от
// полученных аргументов. Иными слова, методы не представляют собой работающую
// эмуляцию методов базы данных.

package main

import (
	"bytes"
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

var _ = bytes.NewBuffer

// 1. Тестовые структуры
// 1.1 Тестовая бд
type MockPgdb struct {
	Rows []string // содержимое таблиц базы данных
}

// Добавить строку в базу данных. Функция не взаимодействует в m.Rows.
func (m *MockPgdb) HireEmployee(d Data, _ context.Context) error {
	if d.FirstName == "" {
		return fmt.Errorf("пустое имя")
	} else if d.FirstName == "John" {
		return fmt.Errorf("запись существует")
	}
	return nil
}

func (m *MockPgdb) Close() {
}

// Считать строку из базы данных. Вернуть поле в интерфейсе Rows и ошибку.
func (m *MockPgdb) GetEmployee(id int, _ context.Context) (Rows, error) {
	if id == 0 {
		return nil, fmt.Errorf("id < 1")
	} else if id == 1 {
		return NewMockRows([]string{m.Rows[0]}), nil
	} else if id == 9999 {
		return nil, nil // id не существует
	}
	return NewMockRows(nil), nil
}

// Считать все строки из базы данных. Вернуть поля в интерфейсе Rows и ошибку.
func (m *MockPgdb) GetAllEmployees(_ context.Context) (Rows, error) {
	return NewMockRows(m.Rows), nil
}

func (m *MockPgdb) GetAllEmployeeNames(_ context.Context) (Rows, error) {
	return NewMockRows(m.Rows), nil
}

func (m *MockPgdb) GetAllEmployeeNonNames(_ context.Context) (Rows, error) {
	return NewMockRows(m.Rows), nil
}

func (m *MockPgdb) GetHighestId(_ context.Context) (int, error) {
	return 0, nil
}

// Удалить строку из базы данных. Функция не взаимодействует в m.Rows.
func (m *MockPgdb) FireEmployee(id int, _ context.Context) error {
	if id == 0 {
		return fmt.Errorf("id < 1")
	} else if id == 1 {
		return fmt.Errorf("id не существует")
	}
	return nil
}

// Удалить строку в базе данных. Функция не взаимодействует в m.Rows.
func (m *MockPgdb) UpdateEmployee(d Data, _ context.Context) error {
	if d.FirstName == "" {
		return fmt.Errorf("пустое имя")
	} else if d.FirstName == "John" {
		return fmt.Errorf("записи не существует")
	}
	return nil
}

// Создать переменную типа MockPgdb. Вернуть её адрес.
func NewMockPgdb() *MockPgdb {
	return &MockPgdb{Rows: []string{"1", "2"}}
}

// 1.2 Тестовый результат запроса бд
type MockRows struct {
	Rows []string // строки в результате
	Idx  int      // следующая непрочитанная строка
	Len  int      // количество строк
}

func (r *MockRows) Close() {
}

func (r *MockRows) Err() error {
	return nil
}

func (r *MockRows) Next() bool {
	return r.Idx < r.Len
}

// Используется в тесте для извлечения id.
// Тест учитывает порядок аргументов, установленный в функциях
// GetAllEmployees и GetEmployee.
func (r *MockRows) Scan(dest ...interface{}) error {
	_, err := fmt.Sscanf(r.Rows[r.Idx], "%d", dest[2])
	if err != nil {
		return err
	}
	r.Idx++
	return nil
}

// Создать новую переменную структуры MockRows. Вернуть как интерфейс Rows.
func NewMockRows(rows []string) Rows {
	return &MockRows{Rows: rows, Idx: 0, Len: len(rows)}
}

// 2. Тесты
// Тесты добавления записи в базу данных.
func TestModelHireEmployee(t *testing.T) {
	var tests = []struct {
		json string
		err  bool
	}{
		{"", true},                 // ошибка json.Unmarshal
		{`{"sane":"no"}`, true},    // отсутствие необходимого поля
		{`{"name":"John"}`, true},  // запись существует
		{`{"name":"jane"}`, false}, // успех
	}
	assert := assert.New(t)
	service, _ := NewModel()
	service.ChangeDatabase(NewMockPgdb())
	for i, test := range tests {
		err := service.HireEmployee(test.json, context.Background())
		assert.Equal(test.err, err != nil, "Тест %d: %v\n", i, err)
	}
}

// Тесты считывания всех строк из базы данных.
func TestModelGetAllEmployees(t *testing.T) {
	var tests = []struct{ err bool }{
		{false}, // успех
	}
	assert := assert.New(t)
	service, _ := NewModel()
	service.ChangeDatabase(NewMockPgdb())
	for i, test := range tests {
		_, err := service.GetAllEmployees(context.Background())
		assert.Equal(test.err, err != nil, "Тест %d: %v\n", i, err)
	}
}

// Тесты получения записи по id.
// Параметр out = false, если результат функции GetEmployee = nil.
// Параметр err = true, если функция GetEmployee возвращает ошибку.
func TestModelGetEmployee(t *testing.T) {
	var tests = []struct {
		id       int
		out, err bool
	}{
		{0, false, true},     // неверный id
		{9999, false, false}, // записи не существует
		{1, true, false},     // успех
	}
	assert := assert.New(t)
	service, _ := NewModel()
	service.ChangeDatabase(NewMockPgdb())
	for i, test := range tests {
		d, err := service.GetEmployee(test.id, context.Background())
		if err != nil {
			assert.Equal(test.err, true, "Тест %d: %v\n", i, err)
			continue
		} else if d == nil {
			assert.Equal(test.out, false, "Тест %d\n", i)
			continue
		}
		assert.Equal(test.id, d.Id, "Тест %d\n", i)
	}
}

// Тесты функции удаления строки.
func TestModelFireEmployee(t *testing.T) {
	var tests = []struct {
		id  int
		err bool
	}{
		{0, true},  // неверный id
		{1, true},  // записи не существует
		{2, false}, // успех
	}
	assert := assert.New(t)
	service, _ := NewModel()
	service.ChangeDatabase(NewMockPgdb())
	for i, test := range tests {
		err := service.FireEmployee(test.id, context.Background())
		assert.Equal(test.err, err != nil, "Тест %d: %v\n", i, err)
	}
}

// Тесты обновления записи в базе данных.
func TestModelUpdateEmployee(t *testing.T) {
	var tests = []struct {
		json string
		err  bool
	}{
		{"", true},                 // ошибка json.Unmarshal
		{`{"sane":"no"}`, true},    // отсутствие необходимого поля
		{`{"name":"John"}`, true},  // записи не существует
		{`{"name":"jane"}`, false}, // успех
	}
	assert := assert.New(t)
	service, _ := NewModel()
	service.ChangeDatabase(NewMockPgdb())
	for i, test := range tests {
		err := service.UpdateEmployee(test.json, context.Background())
		assert.Equal(test.err, err != nil, "Тест %d: %v\n", i, err)
	}
}

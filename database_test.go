package main

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Тесты ошибок базы данных.
func TestDbErrors(t *testing.T) {
	var tests = []struct {
		id           int
		err, dbNAerr bool
	}{
		{1, true, false},  // плохой запрос
		{2, true, true},   // внутренняя ошибка
		{3, false, false}, // успех
	}
	assert := assert.New(t)
	db, err := NewPostgresdb()
	defer db.Close()
	if err != nil {
		t.Logf("NewPostgresdb: %v\n", err)
		return
	}
	for i, test := range tests {
		err := db.GetErr(test.id)
		if err != nil {
			assert.Equal(test.err, true, "Тест %d\n", i)
			if test.id == 2 {
				assert.Equal(err, dbNA, "Тест %d\n", i)
			}
			continue
		}
		assert.Equal(test.err, false, "Тест %d\n", i)
	}
}

// Тесты фукнции удаления строки. Все тесты, кроме первого запрашивают
// самый высокий id и удаляют его.
func TestDbFireEmployee(t *testing.T) {
	var tests = []struct{ err bool }{
		{true},  // id не найден
		{false}, // успех
	}
	assert := assert.New(t)
	db, err := NewPostgresdb()
	defer db.Close()
	if err != nil {
		t.Logf("NewPostgresdb: %v\n", err)
		return
	}
	for i, test := range tests {
		id := 0
		if i > 0 {
			var err error
			id, err = db.GetHighestId(context.Background())
			if err != nil {
				t.Logf("GetHighestId: %v\n", err)
				continue
			}
		}
		err := db.FireEmployee(id, context.Background())
		assert.Equal(test.err, err != nil, "Тест %d: %v\n", i, err)
	}
}

// Тесты считывания всех строк из базы данных. Проверить результат по числу
// рядов. Тест ожидает, что число рядов в базе данных больше нуля.
func TestDbGetEmployees(t *testing.T) {
	var tests = []struct{ err bool }{
		{false}, // успех
	}
	assert := assert.New(t)
	db, err := NewPostgresdb()
	if err != nil {
		t.Logf("NewPostgresdb: %v\n", err)
		return
	}
	defer db.Close()

	for i, test := range tests {
		rows, err := db.GetAllEmployees(context.Background())
		if err != nil {
			t.Logf("Test #%d: %v\n", i, err)
			continue
		} else if rows == nil { // записей не обнаружено
			assert.Equal(test.err, true, "Тест %d: %v\n", i, err)
			continue
		}
		var count int
		var d Data
		for rows.Next() {
			err := rows.Scan(&d.LastName, &d.FirstName, &d.Id, &d.MidName,
				&d.PhoneNum, &d.Position, &d.DoneJobs)
			if err != nil {
				t.Logf("Test #%d: %v\n", i, err)
				break
			}
			count++
		}
		err = rows.Err()
		rows.Close()
		if err != nil {
			assert.Equal(test.err, err != nil, "Тест %d: %v\n", i, err)
			continue
		}
		assert.Less(0, count, "Тест %d: нет результатов", i)
	}
}

// Тесты считывания имени и фамилии всех строк из базы данных.
// Проверить результат по числу рядов.
// Тест ожидает, что число рядов в базе данных больше нуля.
func TestDbGetAllEmployeeNames(t *testing.T) {
	var tests = []struct{ err bool }{
		{false}, // успех
	}
	assert := assert.New(t)
	db, err := NewPostgresdb()
	if err != nil {
		t.Logf("NewPostgresdb: %v\n", err)
		return
	}
	defer db.Close()

	for i, test := range tests {
		rows, err := db.GetAllEmployeeNames(context.Background())
		if err != nil {
			t.Logf("Test #%d: %v\n", i, err)
			continue
		} else if rows == nil { // записей не обнаружено
			assert.Equal(test.err, true, "Тест %d: %v\n", i, err)
			continue
		}
		var count int
		for rows.Next() {
			count++
		}
		err = rows.Err()
		rows.Close()
		if err != nil {
			assert.Equal(test.err, err != nil, "Тест %d: %v\n", i, err)
			continue
		}
		assert.Less(0, count, "Тест %d: нет результатов", i)
	}
}

// Тесты считывания всех полей, кроме имени и фамилии всех строк из базы
// данных. Проверить результат по числу рядов.
// Тест ожидает, что число рядов в базе данных больше нуля.
func TestDbGetAllEmployeeNonNames(t *testing.T) {
	var tests = []struct{ err bool }{
		{false}, // успех
	}
	assert := assert.New(t)
	db, err := NewPostgresdb()
	if err != nil {
		t.Logf("NewPostgresdb: %v\n", err)
		return
	}
	defer db.Close()

	for i, test := range tests {
		rows, err := db.GetAllEmployeeNonNames(context.Background())
		if err != nil {
			t.Logf("Test #%d: %v\n", i, err)
			continue
		} else if rows == nil { // записей не обнаружено
			assert.Equal(test.err, true, "Тест %d: %v\n", i, err)
			continue
		}
		var count int
		for rows.Next() {
			count++
		}
		err = rows.Err()
		rows.Close()
		if err != nil {
			assert.Equal(test.err, err != nil, "Тест %d: %v\n", i, err)
			continue
		}
		assert.Less(0, count, "Тест %d: нет результатов", i)
	}
}

// Тесты считывания строки по id из базы данных.
func TestDbGetEmployee(t *testing.T) {
	var tests = []struct {
		id  int
		err bool
	}{
		{9999, true}, // id не найден
		{2, false},   // успех
	}
	assert := assert.New(t)
	db, err := NewPostgresdb()
	defer db.Close()
	if err != nil {
		t.Logf("NewPostgresdb: %v\n", err)
		return
	}
	for i, test := range tests {
		id, err := getId(db, test.id)
		if err != nil {
			assert.Equal(test.err, true, "Тест %d: %v\n", i, err)
			continue
		}
		assert.Equal(test.id, id, "Тест %d\n", i)
	}
}

// Считать строку по id из базы данных.
// Вернуть id и ошибку.
func getId(db Postgresdb, testId int) (int, error) {
	var count int
	rows, err := db.GetEmployee(testId, context.Background())
	if err != nil {
		return 0, fmt.Errorf("GetRowById: %v", err)
	}
	if rows == nil { // нет результатов
		return 0, fmt.Errorf("нет результатов")
	}
	defer rows.Close()
	var d Data
	for rows.Next() {
		err := rows.Scan(&d.LastName, &d.FirstName, &d.Id, &d.MidName,
			&d.PhoneNum, &d.Position, &d.DoneJobs)
		if err != nil {
			return 0, fmt.Errorf("Scan: %v", err)
		}
		count++
	}
	if err = rows.Err(); err != nil {
		return 0, fmt.Errorf("Err: %v", err)
	}
	if count == 0 {
		return 0, fmt.Errorf("нет результатов")
	}
	return d.Id, nil
}

// Тесты добавления записи в базу данных.
func TestDbHireEmployee(t *testing.T) {
	d1 := Data{FirstName: "Test", LastName: "Testov", MidName: "Testovich",
		PhoneNum: "84951112233", Position: "Tester", DoneJobs: 9}
	var tests = []struct {
		data Data
		err  bool
	}{
		{d1, false}, // успех
	}
	assert := assert.New(t)
	db, err := NewPostgresdb()
	if err != nil {
		t.Logf("NewPostgresdb: %v\n", err)
		return
	}
	defer db.Close()

	for i, test := range tests {
		err := db.HireEmployee(test.data, context.Background())
		assert.Equal(test.err, err != nil, "Тест %d: %v\n", i, err)
	}
}

// Тесты обновления записи в базе данных.
func TestDbUpdateEmployee(t *testing.T) {
	d1 := Data{FirstName: "Test", LastName: "Testov", MidName: "Testovich",
		PhoneNum: "84951112233", Position: "Tester", DoneJobs: 10, Id: 10}
	var tests = []struct {
		data Data
		err  bool
	}{
		{d1, false}, // успех
	}
	assert := assert.New(t)
	db, err := NewPostgresdb()
	if err != nil {
		t.Logf("NewPostgresdb: %v\n", err)
		return
	}
	defer db.Close()

	for i, test := range tests {
		err := db.UpdateEmployee(test.data, context.Background())
		assert.Equal(test.err, err != nil, "Тест %d: %v\n", i, err)
	}
}

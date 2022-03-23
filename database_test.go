package main

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Тесты добавления записи в базу данных.
func TestDbAdd(t *testing.T) {
	d1 := Data{FirstName: "Test", LastName: "Testov", MidName: "Testovich",
		PhoneNum: "84951112233", Position: "Tester", DoneJobs: 9}
	var tests = []struct {
		data Data
		err  bool
	}{
		{d1, false}, // успех
	}
	assert := assert.New(t)
	db, err := NewDatabase()
	if err != nil {
		t.Logf("NewDatabase: %v\n", err)
		return
	}
	defer db.Close()

	for i, test := range tests {
		err := db.Add(test.data, context.Background())
		assert.Equal(test.err, err != nil, "Тест %d: %v\n", i, err)
	}
}

// Тесты считывания всех строк из базы данных. Проверить результат по числу
// рядов. Тест ожидает, что число рядов в базе данных больше нуля.
func TestDbGetAll(t *testing.T) {
	var tests = []struct{ err bool }{
		{false}, // успех
	}
	assert := assert.New(t)
	db, err := NewDatabase()
	if err != nil {
		t.Logf("NewDatabase: %v\n", err)
		return
	}
	defer db.Close()

	for i, test := range tests {
		rows, err := db.GetAll(context.Background())
		if err != nil {
			t.Logf("Test #%d: %v\n", i, err)
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
func TestDbGetId(t *testing.T) {
	var tests = []struct {
		id  int32
		err bool
	}{
		{9999, true}, // id не найден
		{2, false},   // успех
	}
	assert := assert.New(t)
	db, err := NewDatabase()
	defer db.Close()
	if err != nil {
		t.Logf("NewDatabase: %v\n", err)
		return
	}
	for i, test := range tests {
		id, err := getId(db, test.id)
		if err != nil {
			assert.Equal(test.err, true, "Тест %d\n", i)
			continue
		}
		assert.Equal(test.id, id, "Тест %d\n", i)
	}
}

// Считать строку по id из базы данных.
// Вернуть id и ошибку.
func getId(db Database, testId int32) (int32, error) {
	var count, id int32
	rows, err := db.GetId(int(testId), context.Background())
	if err != nil {
		return 0, fmt.Errorf("GetId: %v", err)
	}
	defer rows.Close()
	for rows.Next() {
		count++
		v, err := rows.Values()
		if err != nil {
			return 0, fmt.Errorf("Values: %v", err)
		}
		if len(v) == 0 {
			return 0, fmt.Errorf("количество строк равно 0")
		}
		id = v[0].(int32)
	}
	if err = rows.Err(); err != nil {
		return 0, fmt.Errorf("Err: %v", err)
	}
	if count == 0 {
		return 0, fmt.Errorf("нет результатов")
	}
	return id, nil
}

// Тесты фукнции удаления строки. Все тесты, кроме первого запрашивают
// самый высокий id и удаляют его.
func TestDbRemove(t *testing.T) {
	var tests = []struct{ err bool }{
		{true},  // id не найден
		{false}, // успех
	}
	assert := assert.New(t)
	db, err := NewDatabase()
	defer db.Close()
	if err != nil {
		t.Logf("NewDatabase: %v\n", err)
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
		err := db.Remove(id, context.Background())
		assert.Equal(test.err, err != nil, "Тест %d: %v\n", i, err)
	}
}

// Тесты обновления записи в базе данных.
func TestDbUpdate(t *testing.T) {
	d1 := Data{FirstName: "Test", LastName: "Testov", MidName: "Testovich",
		PhoneNum: "84951112233", Position: "Tester", DoneJobs: 10, Id: 10}
	var tests = []struct {
		data Data
		err  bool
	}{
		{d1, false}, // успех
	}
	assert := assert.New(t)
	db, err := NewDatabase()
	if err != nil {
		t.Logf("NewDatabase: %v\n", err)
		return
	}
	defer db.Close()

	for i, test := range tests {
		err := db.Update(test.data, context.Background())
		assert.Equal(test.err, err != nil, "Тест %d: %v\n", i, err)
	}
}

// Функции и тестовые структуры для тестирования model.go и database.go.

package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Тесты считывания всех строк из базы данных.
func TestStorageDbGetAllEmployees(t *testing.T) {
	var tests = []struct{ err bool }{
		{false}, // успех
	}
	assert := assert.New(t)
	service, _ := NewStorage()
	for i, test := range tests {
		_, err := service.GetAllEmployees(context.Background())
		assert.Equal(test.err, err != nil, "Тест %d: %v\n", i, err)
	}
}

// Тесты конкурентного считывания всех строк из базы данных.
func TestStorageDbGetAllEmployeesConcur(t *testing.T) {
	var tests = []struct{ err bool }{
		{false}, // успех
	}
	assert := assert.New(t)
	service, _ := NewStorage()
	for i, test := range tests {
		_, err := service.GetAllEmployeesConcur(context.Background())
		assert.Equal(test.err, err != nil, "Тест %d: %v\n", i, err)
	}
}

// Тест сравнения результатов обычного считывания всех строк из базы данных
// и конкурентного.
func TestStorageDbGetAllEmployeesComparison(t *testing.T) {
	assert := assert.New(t)
	service, _ := NewStorage()

	// 1. Считать всех работников двумя способами
	d1, err := service.GetAllEmployees(context.Background())
	if err != nil {
		t.Logf("GetAllEmployees: %v\n", err)
		return
	}
	d2, err := service.GetAllEmployeesConcur(context.Background())
	if err != nil {
		t.Logf("GetAllEmployeesConcur: %v\n", err)
		return
	}

	// 2. Сравнить результаты между собой
	if !assert.Equal(len(d1), len(d2), "Количество записей разное\n") {
		return
	}

	entries := make(map[int]Data)
	for _, entry := range d1 {
		entries[entry.Id] = entry
	}
	for _, entry := range d2 {
		if entries[entry.Id] != entry {
			t.Logf("Записи отличаются %v != %v\n", entries[entry.Id], entry)
			break
		}
	}
}

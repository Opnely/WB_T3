// Интерфейс для манипуляции и передачи данных базы данных.

package main

import (
	"context"
	"encoding/json"
	"fmt"
    "log"
)

const (
	MIN_ENTRIES = 64 // начальный размер среза для считывания записей из бд
)

// Поля одной записи таблицы базы данных
type Data struct {
	Id        int    `json:"id"`
	DoneJobs  int    `json:"good_job_count"`
	FirstName string `json:"name"`
	MidName   string `json:"patronynic"`
	LastName  string `json:"last_name"`
	PhoneNum  string `json:"phone"`
	Position  string `json:"position"`
}

// Абстракция структуры Service
type Model interface {
	Add(string, context.Context) error
	ChangeDatabase(Database)
	Close()
	GetAll(context.Context) ([]Data, error)
	GetId(int, context.Context) (*Data, error)
	Remove(int, context.Context) error
	Update(string, context.Context) error
}

type Service struct {
	Db Database // интерфейс для работы с базой данных
}

// Добавить перекодированную в Data строку JSON в базу данных.
func (s *Service) Add(req string, ctx context.Context) error {
	var d Data
	if err := json.Unmarshal([]byte(req), &d); err != nil {
		return err
	}
	return s.Db.Add(d, ctx)
}

// Закрыть соединение с базой данных
func (s *Service) Close() {
    log.Println("Закрытие соединения с базой данных")
	s.Db.Close()
}

// Вернуть адрес структуры Data с данными по id из базы данных.
// Если функция возвращает (nil, nil), запрос выполнен успешно, но данных не
// найдено.
func (s *Service) GetId(id int, ctx context.Context) (*Data, error) {
	rows, err := s.Db.GetId(id, ctx)
	if err != nil {
		return nil, err
	}
	var d Data
	var count int
	for rows.Next() {
		count++
		err := rows.Scan(&d.Id, &d.LastName, &d.FirstName, &d.MidName,
			&d.PhoneNum, &d.Position, &d.DoneJobs)
		if err != nil {
			return nil, fmt.Errorf("Scan: %v", err)
		}
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	if count == 0 { // данных не найдено
		return nil, nil
	}
	return &d, nil
}

// Вернуть срез структур Data со всеми данными из базы данных.
func (s *Service) GetAll(ctx context.Context) ([]Data, error) {
	rows, err := s.Db.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	slice := make([]Data, 0, MIN_ENTRIES)
	for rows.Next() {
		var d Data
		err := rows.Scan(&d.Id, &d.LastName, &d.FirstName, &d.MidName,
			&d.PhoneNum, &d.Position, &d.DoneJobs)
		if err != nil {
			return nil, fmt.Errorf("Scan: %v", err)
		}
		slice = append(slice, d)
	}
	return slice, nil
}

// Удалить запись id из базы данных.
func (s *Service) Remove(id int, ctx context.Context) error {
	return s.Db.Remove(id, ctx)
}

// Изменить строку JSON в базе данных. Перекодировать строку в Data.
func (s *Service) Update(req string, ctx context.Context) error {
	var d Data
	if err := json.Unmarshal([]byte(req), &d); err != nil {
		return err
	}
	return s.Db.Update(d, ctx)
}

// Изменить базу данных. Используется в тестах.
func (s *Service) ChangeDatabase(db Database) {
	s.Db = db
}

// Создать новую переменную Service.
func NewModel() (Model, error) {
	db, err := NewDatabase()
	if err != nil {
		return nil, err
	}
	return &Service{Db: db}, nil
}

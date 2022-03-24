// Интерфейс для манипуляции и передачи данных базы данных.

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
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
	ChangeDatabase(Postgresdb)
	Close()
	FireEmployee(int, context.Context) error
	GetAllEmployees(context.Context) ([]Data, error)
	GetAllEmployeesConcur(context.Context) ([]Data, error)
	GetEmployee(int, context.Context) (*Data, error)
	HireEmployee(string, context.Context) error
	UpdateEmployee(string, context.Context) error
}

type Service struct {
	Pgdb Postgresdb // интерфейс для работы с базой данных postgresdb
}

// Добавить перекодированную в Data строку JSON в базу данных.
func (s *Service) HireEmployee(req string, ctx context.Context) error {
	var d Data
	if err := json.Unmarshal([]byte(req), &d); err != nil {
		return err
	}
	return s.Pgdb.HireEmployee(d, ctx)
}

// Закрыть соединение с базой данных
func (s *Service) Close() {
	log.Println("Закрытие соединения с базой данных")
	s.Pgdb.Close()
}

// Вернуть адрес структуры Data с данными по id из базы данных.
// Если функция возвращает (nil, nil), запрос выполнен успешно, но данных не
// найдено.
func (s *Service) GetEmployee(id int, ctx context.Context) (*Data, error) {
	rows, err := s.Pgdb.GetEmployee(id, ctx)
	if err != nil {
		return nil, err
	}
	if rows == nil { // данных не найдено
		return nil, nil
	}
	var d Data
	for rows.Next() {
		err := rows.Scan(&d.FirstName, &d.LastName, &d.Id, &d.MidName,
			&d.PhoneNum, &d.Position, &d.DoneJobs)
		if err != nil {
			return nil, fmt.Errorf("Scan: %v", err)
		}
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return &d, nil
}

// Вернуть срез структур Data со всеми данными из базы данных.
func (s *Service) GetAllEmployees(ctx context.Context) ([]Data, error) {
	rows, err := s.Pgdb.GetAllEmployees(ctx)
	if err != nil {
		return nil, err
	} else if rows == nil { // данных не найдено
		return nil, nil
	}
	slice := make([]Data, 0, MIN_ENTRIES)
	for rows.Next() {
		var d Data
		err := rows.Scan(&d.FirstName, &d.LastName, &d.Id, &d.MidName,
			&d.PhoneNum, &d.Position, &d.DoneJobs)
		if err != nil {
			return nil, fmt.Errorf("Scan: %v", err)
		}
		slice = append(slice, d)
	}
	return slice, nil
}

// Вернуть срез структур Data со всеми данными из базы данных.
// Считать данные конкурентно используя две функции.
func (s *Service) GetAllEmployeesConcur(ctx context.Context) ([]Data, error) {
	var wg sync.WaitGroup
	var goerr error // ошибка в горутине
	names := make(map[int][2]string)

	// 1. Считать имя и фамилию в карту
	wg.Add(1)
	go func() {
		defer wg.Done()
		var rows Rows
		rows, goerr = s.Pgdb.GetAllEmployeeNames(ctx)
		if goerr != nil {
			return
		} else if rows == nil { // данных не найдено
			return // вне горутины ожидается такая же ошибка
		}
		for rows.Next() {
			var id int
			var fname, lname string
			goerr = rows.Scan(&fname, &lname, &id)
			if goerr != nil {
				return
			}
			names[id] = [2]string{fname, lname}
		}
	}()
	// 2. Считать данные в срез структур Data
	rows, err := s.Pgdb.GetAllEmployeeNonNames(ctx)
	if err != nil {
		return nil, err
	} else if rows == nil { // данных не найдено
		return nil, nil
	}
	slice := make([]Data, 0, MIN_ENTRIES)
	for rows.Next() {
		var d Data
		err := rows.Scan(&d.Id, &d.MidName, &d.PhoneNum,
			&d.Position, &d.DoneJobs)
		if err != nil {
			return nil, fmt.Errorf("Scan: %v", err)
		}
		slice = append(slice, d)
	}
	wg.Wait()
	if goerr != nil {
		return nil, fmt.Errorf("%v", goerr)
	}

	// 3. Добавить данные из карты в срез
	for i, d := range slice {
		n, ok := names[d.Id]
		if !ok {
			continue
		}
		slice[i].FirstName = n[0]
		slice[i].LastName = n[1]
	}
	return slice, nil
}

// Удалить запись id из базы данных.
func (s *Service) FireEmployee(id int, ctx context.Context) error {
	return s.Pgdb.FireEmployee(id, ctx)
}

// Изменить строку JSON в базе данных. Перекодировать строку в Data.
func (s *Service) UpdateEmployee(req string, ctx context.Context) error {
	var d Data
	if err := json.Unmarshal([]byte(req), &d); err != nil {
		return err
	}
	return s.Pgdb.UpdateEmployee(d, ctx)
}

// Изменить базу данных. Используется в тестах.
func (s *Service) ChangeDatabase(db Postgresdb) {
	s.Pgdb = db
}

// Создать новую переменную Service.
func NewModel() (Model, error) {
	db, err := NewPostgresdb()
	if err != nil {
		return nil, err
	}
	return &Service{Pgdb: db}, nil
}

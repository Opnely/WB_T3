// Интерфейс для манипулирования базой данных postgreSQL.
// Соединение осуществляется с помощью database/sql и lib/pq.

package main

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/lib/pq"
)

const (
	DB_URL = "postgres://postgresadmin:admin123@localhost:5432/postgresdb" +
		"?sslmode=disable"
	NO_RESULTS = "50003" // pq.Error.Code
)

// Интерфейс базы данных сотрудников
type Postgresdb interface {
	FireEmployee(int, context.Context) error
	Close()
	GetAllEmployees(context.Context) (Rows, error)
	GetAllEmployeeNames(context.Context) (Rows, error)
	GetAllEmployeeNonNames(context.Context) (Rows, error)
	GetHighestId(context.Context) (int, error)
	GetEmployee(int, context.Context) (Rows, error)
	HireEmployee(Data, context.Context) error
	UpdateEmployee(Data, context.Context) error
}

// Структура с методами для манипулирования информацией в базе данных.
type PostgreSQL struct {
	Conn *sql.DB // установленное соединение с базой данных
}

// Добавить запись d в таблицу employees.employees базы данных.
func (p *PostgreSQL) HireEmployee(d Data, ctx context.Context) error {
	res, err := p.Conn.ExecContext(ctx,
		"SELECT employees.employee_add($1, $2, $3, $4, $5, $6);", d.FirstName,
		d.LastName, d.MidName, d.PhoneNum, d.Position, d.DoneJobs)
	if err != nil {
		return fmt.Errorf("ExecContext: %v\n", err)
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("RowsAffected: %v\n", err)
	}
	if rows != 1 {
		return fmt.Errorf("ожидалось изменение одной строки, получилось: %d",
			rows)
	}
	return nil
}

// Закрыть соединение с базой данных.
func (p *PostgreSQL) Close() {
	p.Conn.Close()
}

// Удалить запись id таблицы employees.employees базы данных.
func (p *PostgreSQL) FireEmployee(id int, ctx context.Context) error {
	res, err := p.Conn.ExecContext(ctx,
		"SELECT employees.employee_remove($1);", id)
	if err != nil {
		return fmt.Errorf("ExecContext: %v\n", err)
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("RowsAffected: %v\n", err)
	}
	if rows != 1 {
		return fmt.Errorf("ожидалось изменение одной строки, получилось: %d",
			rows)
	}
	return nil
}

// Запросить все поля таблицы employees.employees базы данных.
// Вернуть поля в интерфейсе Rows и ошибку.
func (p *PostgreSQL) GetAllEmployees(ctx context.Context) (Rows, error) {
	rows, err := p.Conn.QueryContext(ctx,
		"SELECT * FROM employees.get_all();")
	if err != nil {
		err, ok := err.(*pq.Error)
		if !ok { // невозможная ошибка
			return nil, fmt.Errorf("assertion failed QueryContext: %v\n", err)
		}
		if err.Code == NO_RESULTS {
			return nil, nil
		}
		return nil, fmt.Errorf("QueryContext: %v\n", err)
	}
	return NewRows(rows), nil
}

// Запросить поля name, last_name и id таблицы employees.employees бд.
// Вернуть поля в интерфейсе Rows и ошибку.
func (p *PostgreSQL) GetAllEmployeeNames(ctx context.Context) (Rows, error) {
	rows, err := p.Conn.QueryContext(ctx,
		"SELECT * FROM employees.employees_get_all_part1();")
	if err != nil {
		err, ok := err.(*pq.Error)
		if !ok { // невозможная ошибка
			return nil, fmt.Errorf("assertion failed QueryContext: %v\n", err)
		}
		if err.Code == NO_RESULTS {
			return nil, nil
		}
		return nil, fmt.Errorf("QueryContext: %v\n", err)
	}
	return NewRows(rows), nil
}

// Запросить все поля, кроме name и last_name таблицы employees.employees бд.
// Вернуть поля в интерфейсе Rows и ошибку.
func (p *PostgreSQL) GetAllEmployeeNonNames(ctx context.Context) (Rows, error) {
	rows, err := p.Conn.QueryContext(ctx,
		"SELECT * FROM employees.employees_get_all_part2();")
	if err != nil {
		err, ok := err.(*pq.Error)
		if !ok { // невозможная ошибка
			return nil, fmt.Errorf("assertion failed QueryContext: %v\n", err)
		}
		if err.Code == NO_RESULTS {
			return nil, nil
		}
		return nil, fmt.Errorf("QueryContext: %v\n", err)
	}
	return NewRows(rows), nil
}

// Запросить соответствующие id поля таблицы employees.emoloyess базы данных.
// Если данных не обнаружено, функция возвращает (nil, nil).
// Вернуть поле в интерфейсе Rows и ошибку.
func (p *PostgreSQL) GetEmployee(id int, ctx context.Context) (Rows, error) {
	rows, err := p.Conn.QueryContext(ctx,
		"SELECT * FROM employees.employee_get($1);", id)
	if err != nil {
		err, ok := err.(*pq.Error)
		if !ok { // невозможная ошибка
			return nil, fmt.Errorf("assertion failed QueryContext: %v\n", err)
		}
		if err.Code == NO_RESULTS {
			return nil, nil
		}
		return nil, fmt.Errorf("QueryContext: %v\n", err)
	}
	return NewRows(rows), nil
}

// Обновить запись d таблицы employees.employees базы данных.
func (p *PostgreSQL) UpdateEmployee(d Data, ctx context.Context) error {
	res, err := p.Conn.ExecContext(ctx,
		"SELECT employees.employee_upd($1, $2, $3, $4, $5, $6, $7);", d.Id,
		d.FirstName, d.LastName, d.MidName, d.PhoneNum, d.Position, d.DoneJobs)
	if err != nil {
		return fmt.Errorf("ExecContext: %v\n", err)
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("RowsAffected: %v\n", err)
	}
	if rows != 1 {
		return fmt.Errorf("ожидалось изменение одной строки, получилось: %d",
			rows)
	}
	return nil
}

// Вернуть самый высокий id таблицы employees.employees базы данных.
// Вернуть id и ошибку.
func (p *PostgreSQL) GetHighestId(ctx context.Context) (int, error) {
	rows, err := p.Conn.QueryContext(ctx,
		"SELECT id FROM employees.employees ORDER BY id DESC LIMIT 1;")
	if err != nil {
		err, ok := err.(*pq.Error)
		if !ok { // невозможная ошибка
			return 0, fmt.Errorf("assertion failed QueryContext: %v\n", err)
		}
		if err.Code == NO_RESULTS {
			return 0, nil
		}
		return 0, fmt.Errorf("QueryContext: %v\n", err)
	}
	defer rows.Close()

	var id int
	rows.Next()
	rows.Scan(&id)
	if err := rows.Err(); err != nil {
		return 0, err
	}
	return id, nil
}

// Создать новую переменную интерфейса Postgresdb.
func NewPostgresdb() (Postgresdb, error) {
	cn, err := sql.Open("postgres", DB_URL)
	if err != nil {
		return nil, fmt.Errorf("sql.Open: %v", err)
	}
	return &PostgreSQL{Conn: cn}, nil
}

// 2. Интерфейс для абстракции результатов Query.
// Методы идентичны методам pgx.Rows.
type Rows interface {
	Close()
	Err() error
	Next() bool
	Scan(...interface{}) error
}

type pgRows struct {
	Rows *sql.Rows
}

func (r *pgRows) Close() {
	r.Rows.Close()
}

func (r *pgRows) Err() error {
	return r.Rows.Err()
}

func (r *pgRows) Next() bool {
	return r.Rows.Next()
}

func (r *pgRows) Scan(dest ...interface{}) error {
	return r.Rows.Scan(dest...)
}

func NewRows(rows *sql.Rows) Rows {
	return &pgRows{Rows: rows}
}

// Интерфейс для манипулирования базой данных postgreSQL.
// Соединение осуществляется с помощью pgxpool.

package main

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

const (
	DB_URL = "postgres://postgresadmin:admin123@localhost:5432/postgresdb"
)

// Интерфейс базы данных
type Database interface {
	AddRow(Data, context.Context) error
	Close()
	GetAllRows(context.Context) (Rows, error)
	GetAllRowsNamesId(context.Context) (Rows, error)
	GetAllRowsNonNames(context.Context) (Rows, error)
	GetHighestId(context.Context) (int, error)
	GetRowById(int, context.Context) (Rows, error)
	RemoveRow(int, context.Context) error
	UpdateRow(Data, context.Context) error
}

// Структура с методами для манипулирования информацией в базе данных.
type PostgreSQL struct {
	Conn *pgxpool.Pool // установленное соединение с базой данных
}

// Добавить запись d в таблицу employees.employees базы данных.
func (p *PostgreSQL) AddRow(d Data, ctx context.Context) error {
	rows, err := p.Conn.Query(ctx,
		"SELECT employees.employee_add($1, $2, $3, $4, $5, $6);", d.FirstName,
		d.LastName, d.MidName, d.PhoneNum, d.Position, d.DoneJobs)
	if err != nil {
		return fmt.Errorf("Query: %v\n", err)
	}
	defer rows.Close()
	rows.Next() // прочитать ошибку, если есть
	return rows.Err()
}

// Закрыть соединение с базой данных.
func (p *PostgreSQL) Close() {
	p.Conn.Close()
}

// Запросить все поля таблицы employees.employees базы данных.
// Вернуть поля в интерфейсе Rows и ошибку.
func (p *PostgreSQL) GetAllRows(ctx context.Context) (Rows, error) {
	//rows, err := p.Conn.Query(ctx, "SELECT employees.get_all();")
	rows, err := p.Conn.Query(ctx, "SELECT * FROM employees.employees;")
	if err != nil {
		return nil, fmt.Errorf("Query: %v\n", err)
	}
	return NewRows(rows), nil
}

// Запросить поля name, last_name и id таблицы employees.employees бд.
// Вернуть поля в интерфейсе Rows и ошибку.
func (p *PostgreSQL) GetAllRowsNamesId(ctx context.Context) (Rows, error) {
	rows, err := p.Conn.Query(ctx,
		"SELECT employees.employees_get_all_part1();")
	if err != nil {
		return nil, fmt.Errorf("Query: %v\n", err)
	}
	return NewRows(rows), nil
}

// Запросить все поля, кроме name и last_name таблицы employees.employees бд.
// Вернуть поля в интерфейсе Rows и ошибку.
func (p *PostgreSQL) GetAllRowsNonNames(ctx context.Context) (Rows, error) {
	rows, err := p.Conn.Query(ctx,
		"SELECT employees.employees_get_all_part2();")
	if err != nil {
		return nil, fmt.Errorf("Query: %v\n", err)
	}
	return NewRows(rows), nil
}

// Запросить соответствующие id поля таблицы employees.emoloyess базы данных.
// Вернуть поле в интерфейсе Rows и ошибку.
func (p *PostgreSQL) GetRowById(id int, ctx context.Context) (Rows, error) {
	//rows, err := p.Conn.Query(ctx, "SELECT employees.employee_get($1);", id)
	rows, err := p.Conn.Query(ctx,
		"SELECT * FROM employees.employees WHERE id = $1;", id)
	if err != nil {
		return nil, fmt.Errorf("Query: %v\n", err)
	}
	return NewRows(rows), nil
}

// Удалить запись id таблицы employees.employees базы данных.
func (p *PostgreSQL) RemoveRow(id int, ctx context.Context) error {
	rs, err := p.Conn.Query(ctx, "SELECT employees.employee_remove($1);", id)
	if err != nil {
		return fmt.Errorf("Query: %v\n", err)
	}
	defer rs.Close()
	rs.Next() // прочитать ошибку, если есть
	return rs.Err()
}

// Обновить запись d таблицы employees.employees базы данных.
func (p *PostgreSQL) UpdateRow(d Data, ctx context.Context) error {
	rows, err := p.Conn.Query(ctx,
		"SELECT employees.employee_upd($1, $2, $3, $4, $5, $6, $7);", d.Id,
		d.FirstName, d.LastName, d.MidName, d.PhoneNum, d.Position, d.DoneJobs)
	if err != nil {
		return fmt.Errorf("Query: %v\n", err)
	}
	defer rows.Close()
	rows.Next() // прочитать ошибку, если есть
	return rows.Err()
}

// Вернуть самый высокий id таблицы employees.employees базы данных.
// Вернуть id и ошибку.
func (p *PostgreSQL) GetHighestId(ctx context.Context) (int, error) {
	rows, err := p.Conn.Query(ctx,
		"SELECT id FROM employees.employees ORDER BY id DESC LIMIT 1;")
	if err != nil {
		return 0, fmt.Errorf("Query: %v\n", err)
	}
	defer rows.Close()

	var id int32
	rows.Next()
	rows.Scan(&id)
	if err := rows.Err(); err != nil {
		return 0, err
	}
	return int(id), nil
}

// Создать новую переменную интерфейса Database.
func NewDatabase() (Database, error) {
	cn, err := pgxpool.Connect(context.Background(), DB_URL)
	if err != nil {
		return nil, fmt.Errorf("Connect: %v", err)
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
	Values() ([]interface{}, error)
}

type pgRows struct {
	Rows pgx.Rows
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

func (r *pgRows) Values() ([]interface{}, error) {
	return r.Rows.Values()
}

func NewRows(rows pgx.Rows) Rows {
	return &pgRows{Rows: rows}
}

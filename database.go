// Интерфейс для манипулирования базой данных postgreSQL.
// Соединение осуществляется с помощью database/sql и lib/pq.

package main

import (
	"context"
	"database/sql"
	"fmt"
    "log"
    "os"

	"github.com/lib/pq"
)

const (
    PGDB_URL_FMT = "postgres://%s:%s@%s:%s/%s%s"
    ENV_VAR_PGDB_USER = "PGDB_USER"
    ENV_VAR_PGDB_PWD = "PGDB_PWD"
    //PGDB_USER = ENVIRONMENTALS
    //PGDB_PWD = ENVIRONMENTALS
    // @
    // ADDRESS
    // PORT
    // NAME
    PGDB_SETTINGS = "?sslmode=disable"

    NO_ENV_VAR_MSG_FMT = "переменная окружения %v не установлена" 

	PGDB_URL = "postgres://postgresadmin:admin123@localhost:5432/postgresdb" +
		"?sslmode=disable"
	USER_ERR = "50" // pq.Error.Code.Class
)

// Интерфейс базы данных сотрудников
type Postgresdb interface {
	FireEmployee(int, context.Context) error
	Close()
	GetAllEmployees(context.Context) (Rows, error)
	GetAllEmployeeNames(context.Context) (Rows, error)
	GetAllEmployeeNonNames(context.Context) (Rows, error)
    GetErr(int) error
	GetHighestId(context.Context) (int, error)
	GetEmployee(int, context.Context) (Rows, error)
	HireEmployee(Data, context.Context) error
	UpdateEmployee(Data, context.Context) error
}

// Структура с методами для манипулирования информацией в базе данных.
type PostgreSQL struct {
	Conn *sql.DB // установленное соединение с базой данных
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
        return p.handleDbErr(err)
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("RowsAffected: %v", err)
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
        return nil, p.handleDbErr(err)
	}
	return NewRows(rows), nil
}

// Запросить поля name, last_name и id таблицы employees.employees бд.
// Вернуть поля в интерфейсе Rows и ошибку.
func (p *PostgreSQL) GetAllEmployeeNames(ctx context.Context) (Rows, error) {
	rows, err := p.Conn.QueryContext(ctx,
		"SELECT * FROM employees.employees_get_all_part1();")
	if err != nil {
        return nil, p.handleDbErr(err)
	}
	return NewRows(rows), nil
}

// Запросить все поля, кроме name и last_name таблицы employees.employees бд.
// Вернуть поля в интерфейсе Rows и ошибку.
func (p *PostgreSQL) GetAllEmployeeNonNames(ctx context.Context) (Rows, error) {
	rows, err := p.Conn.QueryContext(ctx,
		"SELECT * FROM employees.employees_get_all_part2();")
	if err != nil {
        return nil, p.handleDbErr(err)
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
        return nil, p.handleDbErr(err)
	}
	return NewRows(rows), nil
}

// Вернуть самый высокий id таблицы employees.employees базы данных и ошибку.
func (p *PostgreSQL) GetHighestId(ctx context.Context) (int, error) {
	rows, err := p.Conn.QueryContext(ctx,
		"SELECT id FROM employees.employees ORDER BY id DESC LIMIT 1;")
	if err != nil {
        return 0, p.handleDbErr(err)
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

// Вернуть пользовательскую ошибку на id = 1.
// Вернуть ошибку базы данных на id = 2.
// Иначе, вернуть nil.
func (p *PostgreSQL) GetErr(id int) error {
	_, err := p.Conn.Exec("SELECT test.get_db_error($1);", id)
    if err != nil {
        return p.handleDbErr(err)
    }
    return nil
}

// Добавить запись d в таблицу employees.employees базы данных.
// Вернуть пользовательскую ошибку как есть.
// Вернуть ошибку базы данных как dbNA, логировать изначальную ошибку.
func (p *PostgreSQL) handleDbErr(origErr error) error {
	err, ok := origErr.(*pq.Error)
	if !ok { // невозможная ошибка
		log.Printf("формат не pq.Error Query: %v", origErr)
        return dbNA
	}
	if err.Code.Class() == USER_ERR {
		return origErr
	}
    log.Printf("Query: %v", origErr)
	return dbNA
}

// Добавить запись d в таблицу employees.employees базы данных.
func (p *PostgreSQL) HireEmployee(d Data, ctx context.Context) error {
	res, err := p.Conn.ExecContext(ctx,
		"SELECT employees.employee_add($1, $2, $3, $4, $5, $6);", d.FirstName,
		d.LastName, d.MidName, d.PhoneNum, d.Position, d.DoneJobs)
	if err != nil {
        return p.handleDbErr(err)
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("RowsAffected: %v", err)
	}
	if rows != 1 {
		return fmt.Errorf("ожидалось изменение одной строки, получилось: %d",
			rows)
	}
	return nil
}


// Обновить запись d таблицы employees.employees базы данных.
func (p *PostgreSQL) UpdateEmployee(d Data, ctx context.Context) error {
	res, err := p.Conn.ExecContext(ctx,
		"SELECT employees.employee_upd($1, $2, $3, $4, $5, $6, $7);", d.Id,
		d.FirstName, d.LastName, d.MidName, d.PhoneNum, d.Position, d.DoneJobs)
	if err != nil {
		return fmt.Errorf("ExecContext: %v", err)
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("RowsAffected: %v", err)
	}
	if rows != 1 {
		return fmt.Errorf("ожидалось изменение одной строки, получилось: %d",
			rows)
	}
	return nil
}


// Создать новую переменную интерфейса Postgresdb.
func NewPostgresdb() (Postgresdb, error) {
    user, ok := os.LookupEnv(ENV_VAR_PGDB_USER)
    if !ok {
        return nil, fmt.Errorf(NO_ENV_VAR_MSG_FMT, ENV_VAR_PGDB_USER)
    }
    pwd, ok := os.LookupEnv(ENV_VAR_PGDB_PWD)
    if !ok {
        return nil, fmt.Errorf(NO_ENV_VAR_MSG_FMT, ENV_VAR_PGDB_PWD)
    }
	cn, err := sql.Open("postgres", fmt.Sprintf(PGDB_URL_FMT, user, pwd,
        cfg.Pgdb.Addr, cfg.Pgdb.Port, cfg.Pgdb.Name, PGDB_SETTINGS))
	if err != nil {
		return nil, fmt.Errorf("sql.Open: %v", err)
	}
    if err := cn.Ping(); err != nil {
		return nil, fmt.Errorf("Ping: %v", err)
	}
	return &PostgreSQL{Conn: cn}, nil
}


// 2. Интерфейс для абстракции результатов Query.
// Методы идентичны методам sql.Rows и pgx.Rows.
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

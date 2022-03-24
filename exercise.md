# Написание первого RESTful API сервиса

Цель: научиться писать RESTful API микросервисы

```
База (postgres)                          # change to whatever
Host: postgres.finance.svc.k8s.dldevel   # change to local
Port: 5432
Name: postgresdb
User: postgresadmin
Password: admin123
```

[A function call example](https://medium.com/geekculture/work-with-go-postgresql-using-pgx-caee4573672)




## Задание 1.

написать сервис, у которого будут методы:

1) нанять сотрудника (201 сreated при успешном выполнении)

```sql
employees.employee_add(
	_id integer,
	_name varchar,
	_last_name varchar,
	_patronymic varchar,
	_phone varchar,
	_position varchar,
	_good_job_count integer
)
```
Ex: SELECT employees.employee_add(
    'test', 'testov', 'testovich', '89291111111', 'cool dude', 1);
    # note the absense of the id field.



2) уволить сотрудника

```sql
employees.employee_remove(_id integer)
```
SELECT employees.employee_remove(6);



3) изменить личные данные сотрудника

```sql
employees.employee_upd(
	_id integer,
	_name varchar,
	_last_name varchar,
	_patronymic varchar,
	_phone varchar,
	_position varchar,
	_good_job_count integer
)
```
SELECT employees.employee_upd(
    6, 'test', 'testov', 'patron', '89291111111', 'cool girl', 2);
    # note the presense of the id field.


4) получить всех сотрудников  

```sql
employees. get_all(_id integer)
```
Ex: SELECT employees.get_all();



5) получить сотрудника по его ID  

```sql
employees.employee_get(_id integer)
```
Ex: SELECT employees.employee_get(1);


Условия:

* структура проекта – mvc https://ru.wikipedia.org/wiki/Model-View-Controller
* ошибки необходимо возвращать по стандарту [rfc7807](https://tools.ietf.org/html/rfc7807)
  {
     "type": "string"               # link to a URI that 
                                    # describes the problem
                                    # in detail. Optional.
     "title": "string"              # general error type.
                                    # Should be short and
                                    # human-readable.
                                    # Shouldn't ever change
     "status": "HTTP return code"   # must be an integer.
     "detail": "string"             # detailed human-
                                    # -readable explanation
     "instance": "string"           # optional unique URI.
  }

* у каждого метода путь должен начинаться с префикса /api/v1
* методы, принимающие в запросе JSON должны иметь миддлвар проверяющий Content-Type. передавать должны только json
* метод 4) должен иметь accept миддлвар, отдающий xml или json
* метод 5) должен принимать employeeId через path-параметр (http://localhost:8000/api/v1/{employeeId}). для этого можно использовать библиотеку https://github.com/gorilla/mux

* при запросах в БД необходимо передавать контекст запроса (r.Context()). Если клиент сервиса перестал ожидать ответ, запрос в базу должен прекращаться
https://medium.com/@nitronick600/extending-gorilla-mux-tutorial-pt-1-90d0ef3affec


## Задание 2.

Добавить технический метод GET /tech/info, который вернет JSON с информацией о приложении:

```js
{ 
	"name": "employees",
	"version": "1.0.0"
}
```

## Задание 3.

Написать API сервис, возвращающий данные обо всех сотрудниках за счёт 
конкурентного вызова employees_get_all_part1() и employees_get_all_part2().



Вспомогательная литература:

* Building RESTful Web services with Go.pdf
* The_Ultimate_Guide_To_Building_Database-Driven_Apps_with_Go.pdf
* Clean_Code.pdf

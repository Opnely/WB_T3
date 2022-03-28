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



## Задание 4.
Цель - научиться обрабатывать коды ошибок при запросах БД. В зависимости от типа ошибки возвращать соответствующий статус код.
Процедура test.get_db_error(_id integer), где
Id = 1 вернет пользовательскую ошибку
Id = 2 вернет ошибку бд,
Любой другой id вернет bool (то есть отработает успешно, без ошибок)_
Ошибка с кодом больше 50 000 (либо уровень ошибки 50) – пользовательская, иначе – внутренняя, в бд в зависимости от уровня ошибки необходимо отдавать сервисом 400 или 500 код. 
При пользовательской ошибке необходимо в ответ также отдавать текст ошибки из БД, при внутренней ошибке нужно модифицировать текст ответа от базы и отдавать общий текст - "хранилище временно не доступно", при этом писать в лог полную ошибку из базы # Конфиг. GitLab CI.

Цель: Научиться конфигурировать сервис. 



## Задание 5.

Сконфигурировать сервис следующим образом:

Для хранения настроек приложения мы используем формат [TOML](https://en.wikipedia.org/wiki/TOML).
Для чтения данных из toml-файлов используем библиотеку "github.com/BurntSushi/toml" для конфига библиотека.

Пример данных, подходящих для хранения в config.toml:

* название приложения
* версия приложения
* название хранимых процедур
* claim для авторизации
* и т.п.

Пример данных, подходящих для хранения в секретах:

* клиент авторизации и его секрет
* имя пользователя и пароль БД
* и другие чувствительные данные
Задание 3.
Изучить переменные окружения и основные инструменты для работы с ними в Go. 
Переписать программу таким образом, чтобы секретные данные, такие как User и Password считывались из переменных окружения.
Вспомогательная литература:

https://gobyexample.com/environment-variables# Сбор метрик приложения



## Задание 6.

Цель: научиться работать с метриками приложения 


Ознакомиться со статьями  
https://prometheus.io/docs/concepts/data_model/#metric-names-and-labels  
https://medium.com/southbridge/prometheus-monitoring-ba8fbda6e83   ([original](https://kjanshair.github.io/2018/02/20/prometheus-monitoring/) )
[recommended blog](https://blog.pvincent.io/2017/12/prometheus-blog-series-part-1-metrics-and-labels/)
[example](https://dev.to/metonymicsmokey/custom-prometheus-metrics-with-go-520n)
[metric types](https://chronosphere.io/learn/an-introduction-to-the-four-primary-types-of-prometheus-metrics/)
[http response time example](https://github.com/brancz/prometheus-example-app/blob/master/main.go)
[basic tutorial and http response time full example](https://www.jajaldoang.com/post/monitor-golang-app-with-prometheus-grafana/)
[another good tutorial](https://blog.pvincent.io/2017/12/prometheus-blog-series-part-4-instrumenting-code-in-go-and-java/)
[tutorial with db](https://percona.community/blog/2021/07/21/create-your-own-exporter-in-go/)
[example with http request and timer](https://gabrieltanner.org/blog/collecting-prometheus-metrics-in-golang)
[official Go guide](https://prometheus.io/docs/guides/go-application/)

Затем добавить в сервис отдельный роут /metrics (без префиксов /api/v1 или /tech) в котором будут отражены:

* метрики на успешные (200+) и неуспешные запросы (400+ и 500+)
default: promhttp_metric_handler_requests_total for 200, 500, and 503
* метрики на время выполнения запроса и время ответа от бд
* метрики на используемую память и количество потраченного процессорного времени
heap go_gc_heap_allocs_bytes_total 
cpu process_cpu_seconds_total counter
VM process_virtual_memory_bytes gauge

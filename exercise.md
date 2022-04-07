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


## Задание 7.
# Docker контейнеры. Minikube/Kubernetes. Health-checks.
Цель: научиться запускать приложение в kubernetes.

* Установить Docker 
* Установить [minikube](https://kubernetes.io/ru/docs/tutorials/hello-minikube/)
[original](https://kubernetes.io/docs/tutorials/hello-minikube/)
[tutorial](https://minikube.sigs.k8s.io/docs/start/)
[common ops](https://minikube.sigs.k8s.io/docs/handbook/)
[gitlab pulling](https://juju.is/tutorials/using-gitlab-as-a-container-registry#7-pull-your-container)
[gitlab secret](https://blog.cloudhelix.io/using-a-private-docker-registry-with-kubernetes-f8d5f6b8f646)
[secret docs](https://kubernetes.io/docs/concepts/configuration/secret/)
* Написать kube manifest (deployment, service).
* Запустить Docker-контейнер с приложением в кубе.
* Научиться работать с секретами. Поместить в секреты db-user, db-password. Убрать их из конфига, если они там были.
* Написать healthcheck пробы для deployment, сделанного ранее (liveness и readiness probes)

Чтобы запустить приложение в кубе, необходимо дать возможность кубу выкачивать Docker-образ из registry GitLab'а. Для этого в манифест нужно указать секрет gitlab-registry-secret.

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: gitlab-registry-secret
data:
  .dockerconfigjson: ew0KICAiYXV0aHMiOiB7DQogICAgImdpdC53aWxkYmVycmllcy5ydTo0NTY3Ijogew0KICAgICAgImF1dGgiOiAiY21WbmFYTjBjbmxmYzNaak9uTTBXbE5JVDJkSWQwYzNZdz09Ig0KICAgIH0NCiAgfQ0KfQ==
type: kubernetes.io/dockerconfigjson
```

Вспомогательная литература:

* Using docker by Adrian Mouat.pdf
* Kubernetes in Action.pdf
* (by-Bilgin-Ibryam,-Roland-Hu)-Kubernetes-Patterns-5233872-(z-lib.org).pdf
# Nats



## Задание 8.
Цель: научиться работать с nats

Часть 1:

1) Написать таблицу на базе(очередь сотрудников), куда будут попадать созданные сотрудники. Нужна новая таблица, которая должна называться например employees_queue. Таким образом в эту таблицу будут падать только новые сотрудники, с которыми можно производить манипуляции не трогая основную большую таблицу.
2) Написать 2 хранимые процедуры - получить данные из этой таблицы, конфирм - удаление данных из этой таблицы. 
3) Реализовать на новый бэкграунд сервис - в нем реализовать процесс забора данных из очереди на базе, публикации их в натс, и конфирм на базе.

Часть 2:

1) Реализовать бэкграунд сервис который читает данные из натса и сохраняет их в файлик

Темы для изучения:
nats, nats streaming

библиотеки для работе с натсом в го:
"github.com/nats-io/stan.go"
"github.com/nats-io/nats.go"

Прим 1.
он должен запускаться сам и работать в отдельной горутине
технические пути в сервисе должны быть доступны в сервисе (tech/info,  metrics)
иными словами нужен новый сервис, в котором не будет тех путей которые были в первом сервисе (с префиксом api/v1) иначе это будет лишний код который будет мешаться

Прим 2.
первые несколько заданий были посвящены рест апи сервису(rest api), так обычно называют сервисы, которые имеют конечные точки (endpoints), на эти конечные точки происходят запросы - пользователем сервиса может выступать другой сервис, пользователь или интерфейс.
вторая часто используемая категория сервисов это background сервис, такой сервис не имеет endpoints кроме технических путей(метрики, health check)
бэкграунд сервисы выполняют разнообразные функции без вмешательства пользователя, при старте сервиса запускается какой-то процесс и он выполняется на протяжении всей жизни этого сервиса, пока он не будет выключен.

по сути есть два главных вида таких сервисов:
1) выполнение задания, которое выполняется в сервисе в бесконечном цикле каждые N минут (практически все сервисы синхронизаторы разных систем работают таким образом, раз в N минут/секунд делают запрос на забор данных из одной системы, и затем данные сохраняются в другую систему); 
2) процесс прослушивания натса/кафки и других подобных систем (там происходит бесконечное ожидание новых событий, при поступлении событий выполняются необходимые действия)



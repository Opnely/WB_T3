Условия:

4) получить всех сотрудников  

```sql
employees. get_all(_id integer)
```

* методы, принимающие в запросе JSON должны иметь миддлвар проверяющий Content-Type. передавать должны только json
[general](https://stackoverflow.com/questions/24455147/how-do-i-send-a-json-string-in-a-post-request-in-go)
[r.Header.Get("Content-type"](https://gist.github.com/rjz/fe283b02cbaa50c5991e1ba921adf7c9)

* при запросах в БД необходимо передавать контекст запроса (r.Context()). Если клиент сервиса перестал ожидать ответ, запрос в базу должен прекращаться
* ошибки необходимо возвращать по стандарту [rfc7807](https://tools.ietf.org/html/rfc7807)
* метод 4) должен иметь accept миддлвар, отдающий xml или json
* метод 5) должен принимать employeeId через path-параметр (http://localhost:8000/api/v1/{employeeId}). для этого можно использовать библиотеку https://github.com/gorilla/mux
Example: handler("/{code}") == http://localhost/123
[var ex1](https://stackoverflow.com/questions/46045756/retrieve-optional-query-variables-with-gorilla-mux)
[full tutorial](https://semaphoreci.com/community/tutorials/building-and-testing-a-rest-api-in-go-with-gorilla-mux-and-postgresql)
[var ex2](https://stackoverflow.com/questions/31371111/mux-vars-not-working)



TESTS
1. Model.go
    create:   bad json, important field is missing, entry already exists, created successfully
    read:     bad id, no id, yes id.
    update:   bad json, important field is missing, no entry to update, updated successfully
    delete:   bad id, no id, success
    read_all: error, yes

2. Router.go
    create:   non-json, json
    read:     no id, bad id, yes id no res, yes id and res
    update:   non-json, json, don't exist, exists
    delete:   no id, bad id, yes id
    read_all: -

3. postgresql.go
    create:   bad query, query exists, fine query
    read:     no entry, no id, yes entry
    update:   bad query, query doesn't exist, fine query
    delete:   no entry, no id, yes entry
    read_all: no results, one result, some results


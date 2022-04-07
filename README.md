# UML
postgresql -> Datbase -> Model <- Router <- Client


# Тесты
```
PGDB_USER=postgresadmin PGDB_PWD=admin123 go test -v *.go
```


# Запуск сервера
```
PGDB_USER=postgresadmin PGDB_PWD=admin123 go run main.go database.go model.go router.go
```

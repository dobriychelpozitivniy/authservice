# запустить сервер
run:
	go run ./cmd/main.go

# посмотреть покрытие проекта тестами
test_coverage:
	go test ./... -coverprofile cover.out
	go tool cover -func cover.out
	del cover.out

# создать swagger документацию\перезаписать на новую
swag:
	swag init -g ./cmd/main.go


.PHONY: update linter linter-fix goose-up goose-down

# Обновить зависимости
update:
	go mod tidy

# Запустить линтер
linter:
	golangci-lint run

# Автофикс линтера
linter-fix:
	golangci-lint run --fix


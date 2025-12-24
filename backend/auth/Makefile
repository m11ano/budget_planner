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


# goose

DB_DSN := "postgres://postgres:12345678@localhost:5433/budget-planner-auth"

.PHONY: goose-up goose-down test test-cov

goose-up:
	@echo "== Applying migrations using DSN: $(DB_DSN) =="
	goose -dir $(CURDIR)/migrations/sql postgres "$(DB_DSN)" up

goose-down:
	@echo "== Rolling back last migration using DSN: $(DB_DSN) =="
	goose -dir $(CURDIR)/migrations/sql postgres "$(DB_DSN)" down

test:
	go test -covermode=atomic -coverpkg=./internal/domain/auth/usecase/account,./internal/domain/auth/usecase/auth,./internal/domain/auth/usecase/session -coverprofile=coverage.out ./... && go tool cover -html=coverage.out -o coverage.html

test-cov:
	$(MAKE) test && go tool cover -func=coverage.out
CURDIR    := $(abspath .)
LOCAL_BIN := $(CURDIR)/.bin

ifeq ($(OS),Windows_NT)
  SHELL := cmd
  .SHELLFLAGS := /C
  EXE := .exe
  WIN_LOCAL_BIN := $(subst /,\\,$(LOCAL_BIN))
  MKDIR_P := if not exist "$(WIN_LOCAL_BIN)" mkdir "$(WIN_LOCAL_BIN)"
  BUF := $(WIN_LOCAL_BIN)\\buf$(EXE)
else
  EXE :=
  MKDIR_P := mkdir -p
  BUF := $(LOCAL_BIN)/buf
endif

.PHONY: bin-deps generate-auth generate-ledger generate-gateway

# Установка buf + локальных плагинов
bin-deps:
ifeq ($(OS),Windows_NT)
	@$(MKDIR_P)
	@IF NOT EXIST "$(BUF)" ( \
	  ECHO == Installing buf into $(WIN_LOCAL_BIN) & \
	  set "GOBIN=$(WIN_LOCAL_BIN)" && go install github.com/bufbuild/buf/cmd/buf@latest || (ECHO !! buf install failed & exit /b 1) \
	) ELSE ( ECHO == buf already installed at $(BUF) )
	@ECHO == Installing local protoc plugins into $(WIN_LOCAL_BIN)
	@set "GOBIN=$(WIN_LOCAL_BIN)" && go install google.golang.org/protobuf/cmd/protoc-gen-go@latest || (ECHO !! protoc-gen-go failed & exit /b 1)
	@set "GOBIN=$(WIN_LOCAL_BIN)" && go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest || (ECHO !! protoc-gen-go-grpc failed & exit /b 1)
	@REM ВАЖНО: PGV ставим из envoyproxy, это корректный module path
	@set "GOBIN=$(WIN_LOCAL_BIN)" && go install github.com/envoyproxy/protoc-gen-validate@latest || (ECHO !! protoc-gen-validate failed & exit /b 1)
	@set "GOBIN=$(WIN_LOCAL_BIN)" && go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest || (ECHO !! grpc-gateway failed & exit /b 1)
	@set "GOBIN=$(WIN_LOCAL_BIN)" && go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest || (ECHO !! openapiv2 failed & exit /b 1)
else
	@$(MKDIR_P) $(LOCAL_BIN)
	@if [ ! -x "$(BUF)" ]; then \
	  echo "== Installing buf into $(LOCAL_BIN)"; \
	  GOBIN=$(LOCAL_BIN) go install github.com/bufbuild/buf/cmd/buf@latest || exit 1; \
	else \
	  echo "== buf already installed at $(BUF)"; \
	fi
	@echo "== Installing local protoc plugins into $(LOCAL_BIN)"
	@GOBIN=$(LOCAL_BIN) go install google.golang.org/protobuf/cmd/protoc-gen-go@latest || exit 1
	@GOBIN=$(LOCAL_BIN) go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest || exit 1
	@# ВАЖНО: PGV из envoyproxy
	@GOBIN=$(LOCAL_BIN) go install github.com/envoyproxy/protoc-gen-validate@latest || exit 1
	@GOBIN=$(LOCAL_BIN) go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest || exit 1
	@GOBIN=$(LOCAL_BIN) go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest || exit 1
endif


generate-auth: bin-deps
	"$(BUF)" generate --template auth/buf.gen.yaml 

generate-ledger: bin-deps
	"$(BUF)" generate --template ledger/buf.gen.yaml 

generate-gateway: bin-deps
	"$(BUF)" generate --template gateway/buf.gen.yaml 


generate: generate-auth generate-ledger generate-gateway
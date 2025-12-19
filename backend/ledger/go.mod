module github.com/m11ano/budget_planner/backend/ledger

go 1.24.3

require github.com/m11ano/budget_planner/backend/auth v0.0.0

replace github.com/m11ano/budget_planner/backend/auth => ../auth

require (
	cloud.google.com/go v0.123.0
	github.com/Masterminds/squirrel v1.5.4
	github.com/georgysavva/scany/v2 v2.1.4
	github.com/go-playground/validator/v10 v10.26.0
	github.com/google/uuid v1.6.0
	github.com/govalues/decimal v0.1.36
	github.com/grpc-ecosystem/go-grpc-middleware/v2 v2.3.3
	github.com/ilyakaznacheev/cleanenv v1.5.0
	github.com/jackc/pgx/v5 v5.7.5
	github.com/pashagolub/pgxmock/v4 v4.8.0
	github.com/pressly/goose/v3 v3.24.3
	github.com/samber/lo v1.51.0
	github.com/stretchr/testify v1.10.0
	go.uber.org/fx v1.24.0
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250818200422-3122310a409c
	google.golang.org/grpc v1.74.2
	google.golang.org/protobuf v1.36.7
)

require (
	github.com/BurntSushi/toml v1.5.0 // indirect
	github.com/LastPossum/kamino v0.0.2 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/gabriel-vasile/mimetype v1.4.9 // indirect
	github.com/go-playground/locales v0.14.1 // indirect
	github.com/go-playground/universal-translator v0.18.1 // indirect
	github.com/golang-jwt/jwt/v5 v5.2.2 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	github.com/jackc/puddle/v2 v2.2.2 // indirect
	github.com/joho/godotenv v1.5.1 // indirect
	github.com/lann/builder v0.0.0-20180802200727-47ae307949d0 // indirect
	github.com/lann/ps v0.0.0-20150810152359-62de8c46ede0 // indirect
	github.com/leodido/go-urn v1.4.0 // indirect
	github.com/mfridman/interpolate v0.0.2 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/sethvargo/go-retry v0.3.0 // indirect
	go.uber.org/dig v1.19.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	go.uber.org/zap v1.27.0 // indirect
	golang.org/x/crypto v0.41.0 // indirect
	golang.org/x/net v0.43.0 // indirect
	golang.org/x/sync v0.16.0 // indirect
	golang.org/x/sys v0.35.0 // indirect
	golang.org/x/text v0.28.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	olympos.io/encoding/edn v0.0.0-20201019073823-d3554ca0b0a3 // indirect
)

module github.com/m11ano/budget_planner/backend/gateway

go 1.24.3

require github.com/m11ano/budget_planner/backend/auth v0.0.0

replace github.com/m11ano/budget_planner/backend/auth => ../auth

require (
	cloud.google.com/go v0.123.0
	github.com/go-playground/validator/v10 v10.26.0
	github.com/gofiber/fiber/v2 v2.52.10
	github.com/gofiber/swagger v1.1.1
	github.com/google/uuid v1.6.0
	github.com/grpc-ecosystem/go-grpc-middleware/v2 v2.3.3
	github.com/ilyakaznacheev/cleanenv v1.5.0
	github.com/jackc/pgx/v5 v5.7.5
	github.com/samber/lo v1.51.0
	github.com/swaggo/swag v1.16.4
	go.uber.org/fx v1.24.0
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250818200422-3122310a409c
	google.golang.org/grpc v1.74.2
	google.golang.org/protobuf v1.36.7
)

require (
	github.com/BurntSushi/toml v1.5.0 // indirect
	github.com/KyleBanks/depth v1.2.1 // indirect
	github.com/LastPossum/kamino v0.0.2 // indirect
	github.com/PuerkitoBio/purell v1.1.1 // indirect
	github.com/PuerkitoBio/urlesc v0.0.0-20170810143723-de5bf2ad4578 // indirect
	github.com/andybalholm/brotli v1.1.1 // indirect
	github.com/gabriel-vasile/mimetype v1.4.9 // indirect
	github.com/go-openapi/jsonpointer v0.19.5 // indirect
	github.com/go-openapi/jsonreference v0.19.6 // indirect
	github.com/go-openapi/spec v0.20.4 // indirect
	github.com/go-openapi/swag v0.19.15 // indirect
	github.com/go-playground/locales v0.14.1 // indirect
	github.com/go-playground/universal-translator v0.18.1 // indirect
	github.com/golang-jwt/jwt/v5 v5.2.2 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	github.com/joho/godotenv v1.5.1 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/klauspost/compress v1.18.0 // indirect
	github.com/leodido/go-urn v1.4.0 // indirect
	github.com/mailru/easyjson v0.7.6 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/mattn/go-runewidth v0.0.16 // indirect
	github.com/rivo/uniseg v0.2.0 // indirect
	github.com/rogpeppe/go-internal v1.14.1 // indirect
	github.com/swaggo/files/v2 v2.0.2 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/valyala/fasthttp v1.51.0 // indirect
	github.com/valyala/tcplisten v1.0.0 // indirect
	go.uber.org/dig v1.19.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	go.uber.org/zap v1.27.0 // indirect
	golang.org/x/crypto v0.41.0 // indirect
	golang.org/x/net v0.43.0 // indirect
	golang.org/x/sys v0.35.0 // indirect
	golang.org/x/text v0.28.0 // indirect
	golang.org/x/tools v0.35.0 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	olympos.io/encoding/edn v0.0.0-20201019073823-d3554ca0b0a3 // indirect
)

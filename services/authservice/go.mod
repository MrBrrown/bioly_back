module bioly/auth

go 1.24.3

replace bioly/storage => ../../common/storage

replace bioly/yamlconf => ../../common/yamlconf

replace bioly/asynclogger => ../../common/asynclogger

require (
	bioly/asynclogger v0.0.0-00010101000000-000000000000 // indirect
	bioly/storage v0.0.0-00010101000000-000000000000 // indirect
	bioly/yamlconf v0.0.0-00010101000000-000000000000 // indirect
	github.com/DATA-DOG/go-sqlmock v1.5.2 // indirect
	github.com/ajg/form v1.5.1 // indirect
	github.com/alexedwards/argon2id v1.0.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/go-chi/chi/v5 v5.2.3 // indirect
	github.com/go-chi/render v1.0.3 // indirect
	github.com/golang-jwt/jwt/v5 v5.3.0 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/jmoiron/sqlx v1.4.0 // indirect
	github.com/lib/pq v1.10.9 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/stretchr/testify v1.11.1 // indirect
	golang.org/x/crypto v0.14.0 // indirect
	golang.org/x/sys v0.13.0 // indirect
	gopkg.in/natefinch/lumberjack.v2 v2.2.1 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

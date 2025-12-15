module bioly/profileservice

go 1.24.3

replace bioly/common/yamlconf => ../../common/yamlconf

replace bioly/common/asynclogger => ../../common/asynclogger

replace bioly/common/storage => ../../common/storage

require (
	bioly/common/asynclogger v0.0.0-00010101000000-000000000000
	bioly/common/storage v0.0.0-00010101000000-000000000000
	bioly/common/yamlconf v0.0.0-00010101000000-000000000000
	github.com/DATA-DOG/go-sqlmock v1.5.2
	github.com/jmoiron/sqlx v1.4.0
	github.com/stretchr/testify v1.11.1
)

require (
	github.com/ajg/form v1.5.1 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/go-chi/chi/v5 v5.2.3 // indirect
	github.com/go-chi/render v1.0.3 // indirect
	github.com/hashicorp/golang-lru/v2 v2.0.7 // indirect
	github.com/lib/pq v1.10.9 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	gopkg.in/natefinch/lumberjack.v2 v2.2.1 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

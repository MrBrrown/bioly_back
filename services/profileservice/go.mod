module bioly/profileservice

go 1.24.3

replace bioly/common/yamlconf => ../../common/yamlconf

replace bioly/common/asynclogger => ../../common/asynclogger

replace bioly/common/storage => ../../common/storage

require (
	bioly/common/asynclogger v0.0.0-00010101000000-000000000000 // indirect
	bioly/common/storage v0.0.0-00010101000000-000000000000 // indirect
	bioly/common/yamlconf v0.0.0-00010101000000-000000000000 // indirect
	github.com/jmoiron/sqlx v1.4.0 // indirect
	github.com/lib/pq v1.10.9 // indirect
	gopkg.in/natefinch/lumberjack.v2 v2.2.1 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
)

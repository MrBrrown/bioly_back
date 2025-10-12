module bioly/profileservice

go 1.24.3

replace bioly/common/yamlconf => ../../common/yamlconf

replace bioly/common/asynclogger => ../../common/asynclogger

require (
	bioly/common/asynclogger v0.0.0-00010101000000-000000000000 // indirect
	bioly/common/yamlconf v0.0.0-00010101000000-000000000000 // indirect
	gopkg.in/natefinch/lumberjack.v2 v2.2.1 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
)

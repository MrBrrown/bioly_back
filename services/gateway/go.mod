module bioly/gateway

go 1.24.3

require bioly/common/yamlconf v0.0.0

require bioly/common/asynclogger v0.0.0

require (
	gopkg.in/natefinch/lumberjack.v2 v2.2.1 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
)

replace bioly/common/yamlconf => ../../common/yamlconf

replace bioly/common/asynclogger => ../../common/asynclogger

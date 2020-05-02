module github.com/MagalixTechnologies/core/middleware

go 1.14

replace github.com/MagalixTechnologies/core/logger => ../logger
require (
	github.com/MagalixTechnologies/core/logger v0.0.0-20200429222314-736083e276b6
	go.uber.org/zap v1.15.0
	goa.design/goa/v3 v3.1.2
)

module github.com/tutorials/go/crud

go 1.16

require (
	github.com/gorilla/mux v1.8.0
	go.opentelemetry.io/otel v1.19.0
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc v1.19.0 // indirect
	google.golang.org/grpc v1.59.0 // indirect
	gorm.io/driver/postgres v1.2.2
	gorm.io/gorm v1.25.1
	gorm.io/plugin/opentelemetry v0.1.4
)

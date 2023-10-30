module github.com/tutorials/go/crud

go 1.16

require (
	github.com/gorilla/mux v1.8.0
	github.com/jackc/pgx/v4 v4.13.0 // indirect
	go.opentelemetry.io/otel v1.19.0
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc v1.19.0
	go.opentelemetry.io/otel/sdk v1.19.0
	google.golang.org/grpc v1.59.0
	gorm.io/driver/postgres v1.5.4
	gorm.io/gorm v1.25.5
	gorm.io/plugin/opentelemetry v0.1.4
)

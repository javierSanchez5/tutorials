package db

import (
	"log"
	"context"
	"fmt"
	"time"
	"os"
	"os/signal"

	"github.com/tutorials/go/crud/pkg/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/plugin/opentelemetry/tracing"
	"go.opentelemetry.io/otel"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
)

func initProvider() (func(context.Context) error, error) {
	ctx := context.Background()

	res, err := resource.New(ctx,
		resource.WithAttributes(
			// the service name used to display traces in backends
			semconv.ServiceName("ServiceName"),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}


	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	conn, err := grpc.DialContext(ctx, "localhost:4317",
		// Note the use of insecure transport here. TLS is recommended in production.
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create gRPC connection to collector: %w", err)
	}

	// Set up a trace exporter
	traceExporter, err := otlptracegrpc.New(ctx, otlptracegrpc.WithGRPCConn(conn))
	if err != nil {
		return nil, fmt.Errorf("failed to create trace exporter: %w", err)
	}

	// Register the trace exporter with a TracerProvider, using a batch
	// span processor to aggregate spans before export.
	bsp := sdktrace.NewBatchSpanProcessor(traceExporter)
	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithResource(res),
		sdktrace.WithSpanProcessor(bsp),
	)
	otel.SetTracerProvider(tracerProvider)

	// set global propagator to tracecontext (the default is no-op).
	otel.SetTextMapPropagator(propagation.TraceContext{})

	// Shutdown will flush any remaining spans and shut down the exporter.
	return tracerProvider.Shutdown, nil
}

func Init() *gorm.DB {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	Shutdown, err := initProvider()
	if err != nil {
		log.Fatal(err)
	}

	defer func ()  {
		if err := Shutdown(ctx); err != nil {
			log.Fatal("failed to shutdown TracerProvider: %w", err)
		}
	}()

	dbURL := "postgres://pg:pass@localhost:5432/crud"

	db, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{})

	if err != nil {
		log.Fatalln(err)
	}

	if err := db.Use(tracing.NewPlugin()); err != nil {
		panic(err)
	}

	tracer := otel.Tracer("gorm.io/plugin/opentelemetry")

	ctx, span := tracer.Start(ctx, "root")
	defer span.End()

	db.AutoMigrate(&models.Book{})

	return db.WithContext(ctx)
}

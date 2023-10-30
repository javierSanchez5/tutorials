package db

import (
	"log"
	"context"

	"github.com/tutorials/go/crud/pkg/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/plugin/opentelemetry/tracing"
	"go.opentelemetry.io/otel"
)

func Init() *gorm.DB {
	ctx := context.Background()

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

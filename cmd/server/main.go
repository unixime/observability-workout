package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"condominium-manager/internal/handler"
	"condominium-manager/internal/model"
	"condominium-manager/internal/repository"
)

func main() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "host=localhost user=postgres password=postgres dbname=condominium port=5432 sslmode=disable"
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{TranslateError: true})
	if err != nil {
		log.Fatalf("connect database: %v", err)
	}
	if err := db.AutoMigrate(&model.Tenant{}); err != nil {
		log.Fatalf("migrate: %v", err)
	}

	r := gin.Default()
	handler.Register(r, repository.NewTenantRepository(db))

	specPath := os.Getenv("SPEC_PATH")
	if specPath == "" {
		specPath = "api/spec.yml"
	}
	handler.RegisterSpec(r, specPath)

	addr := os.Getenv("ADDR")
	if addr == "" {
		addr = ":8080"
	}
	if err := r.Run(addr); err != nil {
		log.Fatal(err)
	}
}

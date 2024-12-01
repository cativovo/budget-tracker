package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/cativovo/budget-tracker/internal/config"
	"github.com/cativovo/budget-tracker/internal/models"
	"github.com/cativovo/budget-tracker/internal/repository"
	"go.uber.org/zap"
)

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}

	defer logger.Sync()
	suggaredLogger := logger.Sugar()

	cfg, err := config.LoadConfig(suggaredLogger)
	if err != nil {
		panic(err)
	}

	r, err := repository.NewRepository(cfg.DBPath, suggaredLogger)
	if err != nil {
		log.Fatal(err)
	}
	defer r.Close()

	if err := r.Migrate(); err != nil {
		suggaredLogger.Fatal(err)
	}

	ey, err := r.ListEntriesByDate(context.Background(), repository.ListEntriesByDateParams{
		StartDate: time.Now().AddDate(0, 0, -5).Format("2006-01-02"),
		EndDate:   time.Now().Format("2006-01-02"),
		AccountID: "237B59CB8AFFF758",
		EntryType: []models.EntryType{models.EntryTypeExpense, models.EntryTypeIncome},
		Limit:     10,
		Offset:    0,
		OrderBy:   repository.Desc,
	})
	if err != nil {
		panic(err)
	}

	b, err := json.MarshalIndent(&ey, "  ", "  ")
	if err != nil {
		panic(err)
	}

	fmt.Println(string(b))
}

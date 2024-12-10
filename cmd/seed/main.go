package main

import (
	"context"
	"flag"
	"math/rand/v2"
	"strings"
	"sync"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/cativovo/budget-tracker/internal/config"
	"github.com/cativovo/budget-tracker/internal/constants"
	"github.com/cativovo/budget-tracker/internal/repository"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

type flags struct {
	Clean bool `json:"clean"`
}

func getFlags() flags {
	cleanPtr := flag.Bool("c", false, "clean db")
	flag.Parse()

	return flags{
		Clean: *cleanPtr,
	}
}

func main() {
	logger := zap.Must(zap.NewProduction()).Sugar()

	flags := getFlags()
	logger.Infow("Flags", "flags", flags)

	cfg, err := config.LoadConfig(logger)
	if err != nil {
		logger.Fatal(err)
	}

	r, err := repository.NewRepository(cfg.DBPath)
	if err != nil {
		logger.Fatal(err)
	}

	if err := r.Migrate(logger); err != nil {
		logger.Fatal(err)
	}

	defer r.Close()

	if flags.Clean {
		cleanDB(r.NonConcurrentDB(), logger)
		return
	}

	logger.Info("seeding...")

	account, err := r.CreateAccount(context.Background(), logger, repository.CreateAccountParams{
		Name: gofakeit.Name(),
	})

	var wg sync.WaitGroup
	categoryIDChan := make(chan *string)

	minCategory := 2
	maxCategory := 10
	categoryCount := rand.IntN(maxCategory) + minCategory + 1
	for i := 0; i < categoryCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			if i%2 == 0 {
				categoryIDChan <- nil
				return
			}

			category, err := r.CreateCategory(context.Background(), logger, repository.CreateCategoryParams{
				Name:      gofakeit.Noun(),
				Icon:      strings.ToLower(gofakeit.Noun()),
				ColorHex:  gofakeit.HexColor(),
				AccountID: account.ID,
			})
			if err != nil {
				logger.Fatal(err)
			}
			categoryIDChan <- &category.ID
		}()
	}

	expenseCountChan := make(chan int)
	incomeCountChan := make(chan int)

	numTask := categoryCount

	for i := 0; i < numTask; i++ {
		go func() {
			for categoryID := range categoryIDChan {
				createExpenses := func() {
					minTransaction := 0
					maxTransaction := 100
					transactionCount := rand.IntN(maxTransaction) + minTransaction + 1

					for i := 0; i < transactionCount; i++ {
						startDate := time.Date(2024, time.September, 1, 0, 0, 0, 0, time.UTC)
						endDate := time.Date(2024, time.December, 31, 0, 0, 0, 0, time.UTC)
						fakeDate := gofakeit.DateRange(startDate, endDate)

						entryTypes := []constants.EntryType{constants.EntryTypeExpense, constants.EntryTypeIncome}
						entryType := entryTypes[rand.IntN(len(entryTypes))]

						if entryType == constants.EntryTypeExpense {
							expenseCountChan <- 1
						} else {
							incomeCountChan <- 1
						}

						description := gofakeit.SentenceSimple()

						result, err := r.CreateEntry(context.Background(), logger, repository.CreateEntryParams{
							EntryType:   entryType,
							Name:        gofakeit.Noun(),
							Amount:      gofakeit.IntRange(1000, 10000),
							Description: &description,
							Date:        fakeDate.Format("2006-01-02"),
							CategoryID:  categoryID,
							AccountID:   account.ID,
						})
						if err != nil {
							logger.Fatal(err)
						}
						_ = result
					}
				}

				wg.Add(1)
				go func() {
					defer wg.Done()
					createExpenses()
				}()
			}
		}()
	}

	type done struct{}
	doneChan := make(chan done)

	go func() {
		wg.Wait()
		doneChan <- done{}
	}()

	var expenseCount int
	var incomeCount int

LOOP:
	for {
		select {
		case ec := <-expenseCountChan:
			expenseCount += ec
		case ic := <-incomeCountChan:
			incomeCount += ic
		case <-doneChan:
			break LOOP
		}
	}

	accountID := account.ID
	if err != nil {
		logger.Fatal(err)
	}
	logger.Infow("Done seeding ID", "account_id", accountID)
	logger.Infow("Results", "category_count", categoryCount, "expense_count", expenseCount, "income_count", incomeCount)
}

func cleanDB(db *sqlx.DB, logger *zap.SugaredLogger) {
	logger.Info("cleaning db...")
	_, err := db.ExecContext(context.Background(), `
		DELETE FROM entry;
		DELETE FROM category;
		DELETE FROM account;
		`)
	if err != nil {
		logger.Fatal(err)
	}

	logger.Info("done")
}

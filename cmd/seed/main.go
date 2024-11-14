package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"math/rand/v2"
	"strings"
	"sync"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/cativovo/budget-tracker/internal/config"
	"github.com/cativovo/budget-tracker/internal/repository"
	"github.com/jackc/pgx/v5/pgxpool"
)

type flags struct {
	clean bool
}

func getFlags() flags {
	cleanPtr := flag.Bool("c", false, "clean db")
	flag.Parse()

	return flags{
		clean: *cleanPtr,
	}
}

func main() {
	flags := getFlags()
	fmt.Println(flags)

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	r, err := repository.NewRepository(context.Background(), cfg.DB)
	if err != nil {
		panic(err)
	}
	defer r.Close()

	if flags.clean {
		cleanDB(r.DBPool())
		return
	}

	log.Println("seeding...")

	account, err := r.CreateAccount(context.Background(), repository.CreateAccountParams{
		Name: gofakeit.Name(),
	})

	var wg sync.WaitGroup
	categoryChan := make(chan repository.CreateCategoryRow)

	minCategory := 2
	maxCategory := 10
	categoryCount := rand.IntN(maxCategory) + minCategory + 1
	for i := 0; i < categoryCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			category, err := r.CreateCategory(context.Background(), repository.CreateCategoryParams{
				Name:      gofakeit.Noun(),
				Icon:      strings.ToLower(gofakeit.Noun()),
				ColorHex:  gofakeit.HexColor()[1:],
				AccountID: account.ID,
			})
			if err != nil {
				log.Fatal("encountered an error in category seed:", err)
			}
			categoryChan <- category
		}()
	}

	expenseCountChan := make(chan int)
	incomeCountChan := make(chan int)

	numTask := categoryCount

	for i := 0; i < numTask; i++ {
		go func() {
			for category := range categoryChan {
				createExpenses := func() {
					minTransaction := 0
					maxTransaction := 100
					transactionCount := rand.IntN(maxTransaction) + minTransaction + 1

					for i := 0; i < transactionCount; i++ {
						startDate := time.Date(2024, time.September, 1, 0, 0, 0, 0, time.UTC)
						endDate := time.Date(2024, time.December, 31, 0, 0, 0, 0, time.UTC)
						fakeDate := gofakeit.DateRange(startDate, endDate)

						trasactionTypes := []repository.TransactionType{repository.TransactionTypeExpense, repository.TransactionTypeIncome}
						transactionType := trasactionTypes[rand.IntN(len(trasactionTypes))]

						if transactionType == repository.TransactionTypeExpense {
							expenseCountChan <- 1
						} else {
							incomeCountChan <- 1
						}

						description := gofakeit.SentenceSimple()

						result, err := r.CreateTransaction(context.Background(), repository.CreateTransactionParams{
							TransactionType: transactionType,
							Name:            gofakeit.Noun(),
							Amount:          gofakeit.IntRange(1000, 10000),
							Description:     &description,
							Date:            &fakeDate,
							CategoryID:      &category.ID,
							AccountID:       account.ID,
						})
						if err != nil {
							log.Fatal("encountered an error in expense seed:", err)
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

	log.Println("done seeding")
	accountID := account.ID.String()
	if err != nil {
		log.Fatal("unmarshal account id:", err)
	}
	log.Println("account id:", accountID)
	log.Printf("results: category: %d, expense: %d, income: %d", categoryCount, expenseCount, incomeCount)
}

func cleanDB(dbpool *pgxpool.Pool) {
	log.Println("cleaning db...")
	c, err := dbpool.Exec(context.Background(), `
		DELETE FROM transaction;
		DELETE FROM category;
		DELETE FROM account;
		`)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(c.RowsAffected())

	log.Println("done")
}

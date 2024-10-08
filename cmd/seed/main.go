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
	"github.com/cativovo/budget-tracker/internal/store"
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
	dbpool, err := store.InitDB(context.Background(), cfg.DB)
	if err != nil {
		log.Fatal(err)
	}
	defer dbpool.Close()

	if flags.clean {
		cleanDB(dbpool)
		return
	}

	log.Println("seeding...")

	queries := store.New(dbpool)
	account, err := queries.CreateAccount(context.Background(), gofakeit.Name())
	if err != nil {
		log.Fatal(err)
	}

	var wg sync.WaitGroup
	categoryChan := make(chan store.CreateCategoryRow)

	minCategory := 5
	maxCategory := 15
	categoryCount := rand.IntN(maxCategory) + minCategory + 1
	for i := 0; i < categoryCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			category, err := queries.CreateCategory(context.Background(), store.CreateCategoryParams{
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
					minTransaction := 10
					maxTransaction := 300
					transactionCount := rand.IntN(maxTransaction) + minTransaction + 1

					for i := 0; i < transactionCount; i++ {
						amount, err := store.NewNumeric(fmt.Sprintf("%.2f", gofakeit.Price(10.0, 10000.0)))
						if err != nil {
							log.Fatal("encountered an error in expense seed:", err)
						}

						description, err := store.NewText(gofakeit.SentenceSimple())
						if err != nil {
							log.Fatal("encountered an error in expense seed:", err)
						}

						startDate := time.Date(2024, time.September, 1, 0, 0, 0, 0, time.UTC)
						endDate := time.Date(2024, time.October, 31, 0, 0, 0, 0, time.UTC)
						fakeDate := gofakeit.DateRange(startDate, endDate).Format("2006-01-02")
						date, err := store.NewDate(fakeDate)
						if err != nil {
							log.Fatal("encountered an error in expense seed:", err)
						}

						trasactionTypes := []int16{0, 1}
						transactionType := trasactionTypes[rand.IntN(len(trasactionTypes))]

						if transactionType == 0 {
							expenseCountChan <- 1
						} else {
							incomeCountChan <- 1
						}

						result, err := queries.CreateTransaction(context.Background(), store.CreateTransactionParams{
							TransactionType: transactionType,
							Name:            gofakeit.Noun(),
							Amount:          amount,
							Description:     description,
							Date:            date,
							CategoryID:      category.ID,
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
	accountID, err := account.ID.MarshalJSON()
	if err != nil {
		log.Fatal("unmarshal account id:", err)
	}
	log.Println("account id:", string(accountID))
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

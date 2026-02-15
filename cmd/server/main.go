package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/iwansofian0512/billing_service/config/db"
	"github.com/iwansofian0512/billing_service/internal/constant"
	delivery "github.com/iwansofian0512/billing_service/internal/handler"
	"github.com/iwansofian0512/billing_service/internal/handler/borrower_handler"
	"github.com/iwansofian0512/billing_service/internal/handler/loan_handler"
	"github.com/iwansofian0512/billing_service/internal/handler/payment_handler"
	"github.com/iwansofian0512/billing_service/internal/repository/borrower_repository"
	"github.com/iwansofian0512/billing_service/internal/repository/loan_repository"
	"github.com/iwansofian0512/billing_service/internal/repository/payment_repository"
	"github.com/iwansofian0512/billing_service/internal/service/borrower_service"
	"github.com/iwansofian0512/billing_service/internal/service/loan_service"
	"github.com/iwansofian0512/billing_service/internal/service/payment_service"
	"github.com/joho/godotenv"
)

type operation func(ctx context.Context) error

func main() {
	_ = godotenv.Load()

	database, err := db.NewPostgresDB()
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	LoanRepo := loan_repository.NewPostgresLoanRepository(database)
	paymentRepo := payment_repository.NewPostgresPaymentRepository(database)
	borrowerRepo := borrower_repository.NewPostgresBorrowerRepository(database)

	loanService := loan_service.NewLoanService(LoanRepo)
	borrowerService := borrower_service.NewBorrowerService(borrowerRepo, LoanRepo, loanService)
	paymentService := payment_service.NewPaymentService(LoanRepo, paymentRepo)

	handler := loan_handler.NewLoanHandler(loanService)
	borrowerHandler := borrower_handler.NewBorrowerHandler(borrowerService)
	paymentHandler := payment_handler.NewPaymentHandler(paymentService)

	router := delivery.NewRouter(handler, borrowerHandler, paymentHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = constant.DefaultPort
	}

	srv := &http.Server{
		Addr:    os.Getenv("HOST") + ":" + port,
		Handler: router,
	}

	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)
	defer stop()

	go func() {
		log.Printf("server starting on port %s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("server error: %v", err)
		}
		stop()
	}()

	<-ctx.Done()

	wait := gracefulShutdown(context.Background(), constant.ShutdownTimeout, map[string]operation{
		"http-server": func(ctx context.Context) error {
			return srv.Shutdown(ctx)
		},
		"postgres": func(ctx context.Context) error {
			return database.Close()
		},
	})

	<-wait
}

func gracefulShutdown(ctx context.Context, timeout time.Duration, ops map[string]operation) <-chan struct{} {
	wait := make(chan struct{})

	go func() {
		timeoutFunc := time.AfterFunc(timeout, func() {
			log.Printf("shutdown timeout %d ms elapsed, force exit", timeout.Milliseconds())
			os.Exit(0)
		})
		defer timeoutFunc.Stop()

		var wg sync.WaitGroup

		for key, op := range ops {
			wg.Add(1)
			innerKey := key
			innerOp := op

			go func() {
				defer wg.Done()

				if err := innerOp(ctx); err != nil {
					log.Printf("%s cleanup failed: %s", innerKey, err.Error())
					return
				}
			}()
		}

		wg.Wait()
		close(wait)
	}()

	return wait
}

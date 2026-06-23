package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"

	"racha-historico/handler"
	"racha-historico/repository"
	"racha-historico/service"
)

func main() {
	godotenv.Load()

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("DATABASE_URL não configurada")
	}

	db, err := sql.Open("mysql", databaseURL)
	if err != nil {
		log.Fatalf("erro ao conectar ao banco: %v", err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatalf("banco de dados inacessível: %v", err)
	}

	userRepo := repository.NewUserRepository(db)
	groupRepo := repository.NewGroupRepository(db)
	expenseRepo := repository.NewExpenseRepository(db)

	authService := service.NewAuthService(userRepo)
	groupService := service.NewGroupService(groupRepo)
	notificationService := service.NewNotificationService()
	expenseService := service.NewExpenseService(expenseRepo, groupRepo, notificationService)
	ocrService := service.NewOCRService()

	authHandler := handler.NewAuthHandler(authService)
	groupHandler := handler.NewGroupHandler(groupService)
	expenseHandler := handler.NewExpenseHandler(expenseService, ocrService)

	r := chi.NewRouter()
	r.Use(chimiddleware.Logger)
	r.Use(chimiddleware.Recoverer)

	r.Post("/auth/register", authHandler.Register)
	r.Post("/auth/login", authHandler.Login)

	r.Group(func(r chi.Router) {
		r.Use(handler.AuthRequired)

		r.Get("/groups", groupHandler.ListGroups)
		r.Post("/groups", groupHandler.CreateGroup)
		r.Get("/groups/{id}", groupHandler.GetGroup)
		r.Post("/groups/{id}/members", groupHandler.AddMember)
		r.Get("/groups/{id}/balances", expenseHandler.GetGroupBalances)
		r.Get("/groups/{id}/expenses", expenseHandler.ListGroupExpenses)
		r.Post("/groups/{id}/expenses", expenseHandler.CreateExpense)

		r.Delete("/expenses/{id}", expenseHandler.DeleteExpense)
		r.Post("/expenses/ocr", expenseHandler.ProcessOCR)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("servidor rodando na porta %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}

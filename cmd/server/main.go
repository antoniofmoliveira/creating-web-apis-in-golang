package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/antoniofmoliveira/apis/configs"
	_ "github.com/antoniofmoliveira/apis/docs"
	"github.com/antoniofmoliveira/apis/internal/entity"
	"github.com/antoniofmoliveira/apis/internal/infra/database"
	"github.com/antoniofmoliveira/apis/internal/infra/webserver/handlers"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/jwtauth"
	httpSwagger "github.com/swaggo/http-swagger"
	"golang.org/x/exp/slog"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// @title         API
// @version       1.0.0
// @description   API for demonstration purposes
// @termsOfService https://github.com/antoniofmoliveira/apis
// @contact.name  Antonio Francisco M Oliveira
// @contact.url   https://github.com/antoniofmoliveira
// @contact.email antoniofmoliveira@gmail.com
// @license.name  MIT
// @license.url   https://opensource.org/licenses/MIT
// @host      	  localhost:8080
// @BasePath      /
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	cfg, err := configs.LoadConfig(".")
	if err != nil {
		panic(err)
	}
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&entity.Product{}, &entity.User{})

	productDB := database.NewProductRepository(db)
	productHandler := handlers.NewProductHandler(productDB)

	userDB := database.NewUserRepository(db)
	userHandler := handlers.NewUserHandler(userDB)

	// public middlewares
	public := func(next http.Handler) http.Handler {
		return middleware.Logger(
			middleware.Recoverer(
				middleware.WithValue("jwt", cfg.TokenAuth)(
					middleware.WithValue("jwtExpiresIn", cfg.JWTExpiresIn)(
						next))))
	}
	// public middlewares plus verification
	private := func(next http.Handler) http.Handler {
		return public(
			jwtauth.Verifier(cfg.TokenAuth)(
				jwtauth.Authenticator(
					next)))
	}

	r := http.NewServeMux()

	r.Handle("GET /products", private(http.HandlerFunc(productHandler.FindAllProducts)))
	r.Handle("POST /products", private(http.HandlerFunc(productHandler.CreateProduct)))
	r.Handle("GET /products/{id}", private(http.HandlerFunc(productHandler.GetProduct)))
	r.Handle("PUT /products/{id}", private(http.HandlerFunc(productHandler.UpdateProduct)))
	r.Handle("DELETE /products/{id}", private(http.HandlerFunc(productHandler.DeleteProduct)))

	r.Handle("POST /users", private(http.HandlerFunc(userHandler.CreateUser)))
	r.Handle("GET /users", private(http.HandlerFunc(userHandler.FindByEmail)))

	r.Handle("POST /users/generate_token", public(http.HandlerFunc(userHandler.GetJwt)))

	r.Handle("GET /docs/", public(httpSwagger.Handler(httpSwagger.URL("http://localhost:8080/docs/doc.json"))))

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.WebServerPort),
		Handler: r,
	}

	go func() {
		url := fmt.Sprintf("http://%s:%s", cfg.WebServerHost, cfg.WebServerPort)
		slog.Info("Server is running at ", url)
		if err := server.ListenAndServe(); err != nil && http.ErrServerClosed != err {
			slog.Error("Could not listen on %s: %v\n", server.Addr, err)
			os.Exit(1)
		}
	}()

	termChan := make(chan os.Signal, 1)
	signal.Notify(termChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	<-termChan
	slog.Info("server: shutting down")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		slog.Error("Could not shutdown the server: %v\n", err)
		os.Exit(1)
	}
	slog.Info("Server stopped")
	os.Exit(0)
}

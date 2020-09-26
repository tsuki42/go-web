package cmd

import (
	"context"
	"crypto/tls"
	"github.com/aka-achu/go-web/controller"
	"github.com/aka-achu/go-web/repo"
	"github.com/aka-achu/go-web/service"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

// Execute, initializes the web application
// It creates a router,
//    configures the middlewares,
//    initializes the controllers,
//    initialized the repo layers,
//    initialized the service layers,
//    registers the handle functions,
//    starts the web server.
func Execute() {

	// Creating a new router
	router := mux.NewRouter()
	// Handling cors access
	router.Use(
		cors.AllowAll().Handler,
	)
	// Getting a new database connection
	db, err := repo.GetConnection()
	if err != nil {
		panic(err)
	}

	// Registering handle functions
	controller.InitUserRoute(
		router,
		controller.NewUserController(),
		repo.NewUserRepo(db),
		service.NewUserService(),
	)
	controller.InitAuthenticationRoute(
		router,
		controller.NewUserController(),
		repo.NewUserRepo(db),
		service.NewAuthenticationService(),
	)

	server := &http.Server{
		Addr:         os.Getenv("SERVER_ADDRESS"),
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
		TLSConfig: &tls.Config{
			PreferServerCipherSuites: true,
			CurvePreferences: []tls.CurveID{
				tls.CurveP256,
				tls.X25519, // Go 1.8 only
			},
			MinVersion: tls.VersionTLS12,
			CipherSuites: []uint16{
				tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305, // Go 1.8 only
				tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,   // Go 1.8 only
				tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			},
		},
	}

	done := make(chan bool)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	go func() {
		<-quit
		log.Println("Server is shutting down...")
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		server.SetKeepAlivesEnabled(false)
		if err := server.Shutdown(ctx); err != nil {
			log.Fatalf("Could not gracefully shutdown the server: %v\n", err)
		}
		close(done)
	}()

	log.Println("Server has started at", os.Getenv("SERVER_ADDRESS"))
	if err := server.ListenAndServeTLS(
		os.Getenv("SSL_CERT"),
		os.Getenv("SSL_KEY"),
	); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Could not listen on %s: %v\n", os.Getenv("SERVER_ADDRESS"), err)
	}
	<-done
	log.Println("Server has stopped")
}

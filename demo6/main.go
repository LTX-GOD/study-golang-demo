package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/zsm/CurrencyExchangeApp/config"
	"github.com/zsm/CurrencyExchangeApp/routes"
)

func main() {
	config.InitConfig()
	r := routes.SetupRouter()

	port := os.Getenv("PORT")
	if port == "" {
		port = config.AppConfig.App.Port
	}
	if port == "" {
		port = ":8080"
	}

	srv := &http.Server{
		Addr:    port,
		Handler: r,
	}

	fmt.Println("server is running on", port)

	//优雅关闭
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	log.Println("Server exiting")

}

package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/prajwal-huggi/context_chatbot/internal/config"
	"github.com/prajwal-huggi/context_chatbot/internal/http/handlers/rag"
	"github.com/prajwal-huggi/context_chatbot/internal/utils/response"
)

func main(){
	fmt.Println("Running the go")

	err := godotenv.Load(".env") // or backend/.env depending on your structure
	if err != nil {
		log.Println("No .env file found, falling back to system env")
	}

	// 1) Load config
	cfg:= config.MustLoad()

	// 2) Setup the database

	// 3) Setup the router
	router:= http.NewServeMux()

	router.HandleFunc("GET /api/", func(w http.ResponseWriter, r *http.Request){
		response.WriteJson(w, http.StatusOK, map[string]string {"message":"Hello GoLang"})
	})
	router.HandleFunc("POST /api/reset", rag.ResetRAG())
	// router.HandleFunc("POST /api/document", rag.UploadDocument())
	router.HandleFunc("POST /api/answer", rag.GetAnswer())

	// 4) Setup the server
	server:= http.Server{
		Addr: cfg.Addr,
		Handler: router,
	}

	slog.Info("Server started successfully: ", slog.String("address", cfg.Addr))

	// Implementing the Graceful Shutdown
	done:= make(chan os.Signal, 1)

	signal.Notify(done , os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func(){
		err:= server.ListenAndServe()
		if err!= nil && err != http.ErrServerClosed{
			log.Fatal("Failed to start the server ")
		}
	}()

	<- done

	slog.Info("Shutting down the server")

	ctx, cancel:=context.WithTimeout(context.Background(), 5* time.Second)
	defer cancel()

	err= server.Shutdown(ctx)
	if err!= nil{
		slog.Error("Failed to shutdown server", slog.String("error", err.Error()))
	}

	slog.Info("Server Shutdown Successfully ")
}

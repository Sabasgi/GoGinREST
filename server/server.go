package main

import (
	"context"
	"goGinRest/modules/api"
	"goGinRest/modules/database"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func init() {
	logfile, err := os.OpenFile("./logs/app_"+time.Now().Format("2006-01-02")+".log", os.O_APPEND|os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		log.Fatal(err)
	}
	mw := io.MultiWriter(os.Stdout, logfile)
	log.SetFlags(log.LstdFlags | log.Lshortfile | log.LUTC)
	log.SetOutput(mw)
	// log.SetFlags(log.Lshortfile | log.LstdFlags)
	// log.Println("Starting Logging")
	err = godotenv.Load(".env.dev")
	if err != nil {
		log.Fatalf("Error loading Configs : %v", err)
	}
}

func main() {
	log.Println("Hey Welcome!")
	log.Println("initializing the Exam server")
	log.Println("Starting server...")
	router := gin.Default()
	api.Init(router)
	port := os.Getenv("port")
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to initialize server: %v\n", err)
		}
	}()
	log.Println("Server sarted runnning on :", port)
	quit := make(chan os.Signal)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// This blocks until a signal is passed into the quit channel
	<-quit
	//TODO KILL MONGO CONNECTIONS HERE
	// shutdown data sources
	if err := killDBConnections(); err != nil {
		log.Fatalf("A problem occurred gracefully shutting down data sources: %v\n", err)
	}

	if err := srv.Shutdown(nil); err != nil {
		log.Fatalf("Server shutdown failed: %v", err)
	}
	log.Println("Server shutdown complete")
}
func killDBConnections() error {
	sqlConn := database.Sqlinstance
	mongoConn := database.MongoInstance
	var err error
	for k, v := range sqlConn {
		if err = v.Close(); err != nil {
			log.Printf("Error closing SQL conection %s", k)
			return err
		}
	}
	// keys := mongoConn
	for _, v := range mongoConn {
		// conn:=mongoConn[Key]
		// conn := mongoConn.Get(v)
		if err = v.Disconnect(context.TODO()); err != nil {
			log.Printf("Error closing Mongo conection %s", v)
			return err
		}
	}
	return err
}

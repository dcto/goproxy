package main

import (
	"context"
	"goproxy/app/apis"
	"goproxy/app/jobs"
	"goproxy/pkg/config"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
	"github.com/lestrrat-go/file-rotatelogs"
)


var (
	g errgroup.Group
)


func main() {
	router := gin.Default()
	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Welcome GoProxy Server")
	})

	router.GET("/get", apis.Get)
	router.GET("/pop", apis.Pop)
	router.GET("/all", apis.All)
	router.GET("/raw", apis.Raw)
	


	server := &http.Server{
		Addr:   config.GetString("server.addr"),
		Handler: router,
	}

	g.Go(func() error {
		return server.ListenAndServe();
	})

	g.Go(func() error {		
		jobs.Cron(config.GetString("check.time"), jobs.Ping)		
		jobs.Crontab.Start()
		select {} 
	})

	logger, _ := rotatelogs.New(
		config.GetString("logger.dir")+"/"+config.GetString("logger.file"),
		rotatelogs.WithRotationTime(time.Hour),
		)	
	log.SetOutput(logger)
	
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)

	}

	if err := jobs.Crontab.Stop().Err(); err != nil {
		log.Fatal("Crontab Shutdown:", err)
	}
	log.Println("Server halt done")
	log.Println("Crontab halt done")

}
package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	port := flag.Int("port", 5000, "HTTP Listen Port(default 5000)")
	flag.Parse()
	r := gin.Default()
	r.GET("/api/services", ListServicesApi)
	r.PATCH("/api/services/:unit", CommandServiceApi)

	srv := &http.Server{
		Addr:    ":" + strconv.Itoa(*port),
		Handler: r,
	}

	go func() {
		// 服务连接
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	cleanup := func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			log.Fatal("Server Shutdown:", err)
		}
	}

	Run(cleanup)
}

func Run(cleanup func()) error {
	state := 1
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

EXIT:
	for sig := range sc {
		log.Printf("Receive signal[%s]", sig.String())
		switch sig {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			state = 0
			break EXIT
		default:
			break EXIT
		}
	}
	cleanup()
	log.Println("Server exit.")
	time.Sleep(time.Second)
	os.Exit(state)
	return nil
}

func ListServicesApi(c *gin.Context) {
	items, err := ListServices()
	if err != nil {
		c.JSON(500, gin.H{
			"code":    500,
			"message": err.Error(),
		})
	} else {
		c.JSON(200, gin.H{
			"code":    0,
			"message": "ok",
			"data":    items,
		})
	}
}

func CommandServiceApi(c *gin.Context) {
	unit := c.Param("unit")
	command := c.DefaultQuery("command", "")
	err := CommandService(command, unit)
	if err != nil {
		c.JSON(500, gin.H{
			"code":    500,
			"message": err.Error(),
		})
	} else {
		c.JSON(200, gin.H{
			"code":    0,
			"message": "ok",
		})
	}
}

package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/my1562/api/config"
	"github.com/my1562/api/models"
	"github.com/my1562/api/router"
	"github.com/my1562/api/routes"
	"go.uber.org/dig"
)

func main() {

	c := dig.New()
	c.Provide(config.NewConfig)
	c.Provide(routes.NewSubscriptionService)
	c.Provide(routes.NewAddressService)
	c.Provide(router.NewRouter)
	c.Provide(models.NewDatabase)

	err := c.Invoke(func(r *mux.Router, config *config.Config) {
		fmt.Printf("Listening at: %s", config.Port)
		srv := &http.Server{
			Addr:           ":" + config.Port,
			Handler:        r,
			ReadTimeout:    10 * time.Second,
			WriteTimeout:   10 * time.Second,
			MaxHeaderBytes: 1 << 20,
		}

		err := srv.ListenAndServe()

		if err != nil {
			fmt.Print(err)
		}
	})
	if err != nil {
		panic(err)
	}

}

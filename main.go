package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/handlers"
	"github.com/my1562/api/config"
	"github.com/my1562/api/models"
	"github.com/my1562/api/notifier"
	"github.com/my1562/api/router"
	"github.com/my1562/api/routes"
	"github.com/my1562/geocoder"
	"go.uber.org/dig"
)

func main() {

	c := dig.New()
	c.Provide(config.NewConfig)
	c.Provide(routes.NewSubscriptionService)
	c.Provide(routes.NewAddressService)
	c.Provide(router.NewGinRouter)
	c.Provide(models.NewDatabase)
	c.Provide(notifier.NewNotifier)
	c.Provide(
		func(conf *config.Config) (*geocoder.Geocoder, error) {
			geo := geocoder.NewGeocoder()
			geo.BuildSpatialIndex(100)
			return geo, nil
		})

	err := c.Invoke(func(config *config.Config, grouter *gin.Engine) error {
		fmt.Printf("Listening at: %s", config.Port)
		srv := &http.Server{
			Addr:           ":" + config.Port,
			Handler:        handlers.CORS(handlers.ExposedHeaders([]string{"Content-Range"}))(grouter),
			ReadTimeout:    10 * time.Second,
			WriteTimeout:   10 * time.Second,
			MaxHeaderBytes: 1 << 20,
		}

		err := srv.ListenAndServe()

		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
}

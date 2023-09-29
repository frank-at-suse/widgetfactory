package main

import (
	_ "embed"
	"flag"
	"github.com/ebauman/widgetfactory/database"
	"github.com/ebauman/widgetfactory/web"
	"github.com/gofiber/fiber/v2"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
)

var (
	dsn               string
	staticContentPath string
)

func init() {
	logrus.SetOutput(os.Stdout)

	flag.StringVar(&dsn, "dsn", "", "MySQL DSN")
	flag.StringVar(&staticContentPath, "static-content-path", "./app",
		"Path from where static content is served")

	flag.Parse()
}

var (
	requestCount = promauto.NewCounter(prometheus.CounterOpts{
		Name: "widgetfactory_requests_total",
		Help: "The total number of received http requests (non-metrics)",
	})
)

func main() {
	logrus.Infof("Starting with mysql dsn %s", dsn)
	logrus.Infof("Starting with static content path %s", staticContentPath)
	prom := http.NewServeMux()
	prom.Handle("/metrics", promhttp.Handler())

	stopCh := make(chan error)

	go func() {
		logrus.Infof("Metrics on 0.0.0.0:9090")

		err := http.ListenAndServe("0.0.0.0:9090", prom)

		stopCh <- err
	}()

	app := fiber.New(fiber.Config{
		AppName: "Widget Factory",
	})

	app.Use(func(c *fiber.Ctx) error {
		logrus.Infof("Request for %s", c.Request().RequestURI())
		requestCount.Inc()

		return c.Next()
	})

	db := database.New(dsn)

	svr := web.New(app, db, staticContentPath)

	go func() {
		logrus.Infof("Data on 0.0.0.0:8080")

		err := svr.Listen("0.0.0.0:8080")

		stopCh <- err
	}()

	err := <-stopCh

	logrus.Fatal(err)
}

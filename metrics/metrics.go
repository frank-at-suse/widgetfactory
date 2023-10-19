package metrics

import (
	"fmt"
	"github.com/ebauman/widgetfactory/database"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/sirupsen/logrus"
	"time"
)

var (
	totalOrdersMetric = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "total_orders",
		Help: "Total number orders stored in the database",
	})

	totalWidgetsMetric = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "total_widgets",
		Help: "Total number of widgets stored in the database",
	})

	orderTotalWidget = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "total_widgets_on_order",
		Help: "Total number of a specific type of widget currently on order",
	}, []string{"widget_id"})
)

func Start(db *database.DB, stopCh chan error) {
	for {
		select {
		case <-stopCh:
			return
		default:
			orders, err := db.ListOrders()
			if err != nil {
				logrus.Errorf("error listing orders: %s", err.Error())
			}

			totalOrdersMetric.Set(float64(len(orders)))

			totals := map[int]int{}

			for _, o := range orders {
				totals[o.Widget] += o.Quantity
			}

			widgets, err := db.ListWidgets()
			if err != nil {
				logrus.Errorf("error listing widgets: %s", err.Error())
			}

			totalWidgetsMetric.Set(float64(len(widgets)))

			for k, v := range totals {
				orderTotalWidget.With(prometheus.Labels{"widget_id": fmt.Sprintf("%d", k)}).Set(float64(v))
			}

			time.Sleep(30 * time.Second)
		}
	}
}

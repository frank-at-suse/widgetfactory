package web

import (
	_ "embed"
	"encoding/json"
	"github.com/ebauman/widgetfactory/database"
	"github.com/ebauman/widgetfactory/pubsub"
	"github.com/ebauman/widgetfactory/types"
	"github.com/gofiber/contrib/websocket"
	_ "github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

var (
	stream = pubsub.New()
)

type Server struct {
	app *fiber.App
	db  *database.DB
}

func New(app *fiber.App, db *database.DB, staticContentPath string) *Server {
	svr := &Server{
		app: app,
		db:  db,
	}

	svr.register(staticContentPath)

	return svr
}

func (s *Server) Listen(addr string) error {
	return s.app.Listen(addr)
}

func (s *Server) register(staticContentPath string) {
	s.app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "*",
	}))

	s.app.Get("/widget", s.listWidget)
	s.app.Get("/order", s.listOrders)
	s.app.Post("/widget", s.createWidget)
	s.app.Post("/order", s.createOrder)
	s.app.Delete("/widget", s.deleteWidget)
	s.app.Delete("/order", s.deleteOrder)
	s.app.Post("/sql", s.sql)

	s.app.Use("/ws", func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	s.app.Get("/ws/orders", websocket.New(s.streamOrders))
	s.app.Get("/ws/widgets", websocket.New(s.streamWidgets))

	s.app.Static("/", staticContentPath)
}

func (s *Server) sql(c *fiber.Ctx) error {
	stmt := string(c.Body())

	res, err := s.db.Query(stmt)
	if err != nil {
		return err
	}

	return c.JSON(res)
}

func (s *Server) listWidget(c *fiber.Ctx) error {
	widgets, err := s.db.ListWidgets()
	if err != nil {
		return err
	}

	return c.JSON(widgets)
}

func (s *Server) listOrders(c *fiber.Ctx) error {
	orders, err := s.db.ListOrders()
	if err != nil {
		return err
	}

	return c.JSON(orders)
}

func (s *Server) createWidget(c *fiber.Ctx) error {
	widget, err := parseWidget(c)

	if err != nil {
		return err
	}

	widget, err = s.db.CreateWidget(widget)
	if err != nil {
		return err
	}

	go func() {
		logrus.Infof("sending stream message create with widget id %d", widget.ID)
		stream.Publish("widget", pubsub.StreamMessage{pubsub.StreamMessageKindCreate, widget})
	}()

	return c.JSON(widget)
}

func (s *Server) createOrder(c *fiber.Ctx) error {
	order, err := parseOrder(c)

	if err != nil {
		return err
	}

	order, err = s.db.CreateOrder(order)
	if err != nil {
		return err
	}

	go func() {
		logrus.Infof("sending stream message create with order id %s", order.ID)
		stream.Publish("order", pubsub.StreamMessage{pubsub.StreamMessageKindCreate, order})
	}()

	return c.JSON(order)
}

func (s *Server) deleteWidget(c *fiber.Ctx) error {
	widget, err := parseWidget(c)

	if err != nil {
		return err
	}

	if err = s.db.DeleteWidget(widget); err != nil {
		return err
	}

	go func() {
		logrus.Infof("sending stream message delete with widget id %d", widget.ID)
		stream.Publish("widget", pubsub.StreamMessage{pubsub.StreamMessageKindDelete, widget})
	}()

	return nil
}

func (s *Server) deleteOrder(c *fiber.Ctx) error {
	order, err := parseOrder(c)

	if err != nil {
		return err
	}

	if err := s.db.DeleteOrder(order); err != nil {
		return err
	}

	go func() {
		logrus.Infof("sending stream message delete with order id %d", order.ID)
		stream.Publish("order", pubsub.StreamMessage{pubsub.StreamMessageKindDelete, order})
	}()

	return nil
}

func (s *Server) streamWidgets(c *websocket.Conn) {
	// start by getting all widgets and sending a load message
	widgets, err := s.db.ListWidgets()
	if err != nil {
		c.WriteJSON(pubsub.StreamMessage{pubsub.StreamMessageKindError, err})
		return
	}

	c.WriteJSON(pubsub.StreamMessage{pubsub.StreamMessageKindLoad, widgets})

	subChan := stream.Subscribe("widget")
	closeChan := make(chan struct{})

	c.SetCloseHandler(func(code int, text string) error {
		close(closeChan)
		return nil
	})
	subId := uuid.NewString()
	for {
		select {
		case <-closeChan:
			logrus.Infof("closing subscriber with id %s", subId)
			return
		case widgetMessage := <-subChan:
			logrus.Infof("received message from widget topic on subscriber id %s: %v", subId, widgetMessage)
			c.WriteJSON(widgetMessage)
		}
	}

}

func (s *Server) streamOrders(c *websocket.Conn) {
	// start by getting all orders and sending a load message
	orders, err := s.db.ListOrders()
	if err != nil {
		c.WriteJSON(pubsub.StreamMessage{pubsub.StreamMessageKindError, err})
		return
	}

	c.WriteJSON(pubsub.StreamMessage{pubsub.StreamMessageKindLoad, orders})

	subChan := stream.Subscribe("order")

	for {
		select {
		case orderMessage := <-subChan:
			logrus.Infof("received message from order topic: %v", orderMessage)
			c.WriteJSON(orderMessage)
		}
	}

}

func parseWidget(c *fiber.Ctx) (*types.Widget, error) {
	data := c.Request().Body()

	var widget = types.Widget{}
	err := json.Unmarshal(data, &widget)
	if err != nil {
		return nil, err
	}

	return &widget, nil
}

func parseOrder(c *fiber.Ctx) (*types.Order, error) {
	data := c.Request().Body()

	var order = types.Order{}
	err := json.Unmarshal(data, &order)
	if err != nil {
		return nil, err
	}

	return &order, nil
}

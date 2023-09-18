package router

import (
	"context"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"html/template"
	"l0/internal/domain/service"
	"net/http"
)

type EchoRouter struct {
	service *service.OrderProcessing
	r       *echo.Echo
}

func NewEchoRouter(service *service.OrderProcessing) *EchoRouter {
	return &EchoRouter{
		service: service,
		r:       echo.New(),
	}
}

func (s *EchoRouter) Start(host, port string) error {
	s.r.Use(middleware.Logger())
	s.r.GET("/order/:id", s.GetOrderByID)
	s.r.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	return s.r.Start(":" + port)
}

func (s *EchoRouter) GetOrderByID(c echo.Context) error {
	id := c.Param("id")
	order, err := s.service.GetOrder(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, err)
	}
	tmpl, err := template.ParseFiles("templates/order.html")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "ERROR")
	}
	return tmpl.Execute(c.Response().Writer, order)
}

func (s *EchoRouter) Stop() error {
	return s.r.Shutdown(context.TODO())
}

package controller

import (
	"controller/internal/service"
	"log"
)

type Controller struct {
	service *service.Service
}

func NewController(service *service.Service) *Controller {
	return &Controller{
		service: service,
	}
}

func (c *Controller) InitListener() {
	for {
		m, err := c.service.ReaderWriterKafka.ReadMessage()
		if err != nil {
			log.Println(err)
		}
		log.Println(m)

		err = c.service.Processing(m)
		if err != nil {
			log.Println(err)
		}
	}
}

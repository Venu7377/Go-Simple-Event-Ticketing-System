package model

import (
	"gorm.io/gorm"
)

type Event struct {
	gorm.Model
	Name            string `json:"name"`
	Date            string `json:"date"`
	Place           string `json:"place"`
	NumberOfTickets int32  `json:"noOfTickets"`
}

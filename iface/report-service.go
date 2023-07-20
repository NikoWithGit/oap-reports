package iface

import (
	"oap-reposts/model"
	"time"
)

type ReportService interface {
	AddOrder(order *model.Order) error
	GetAll(from time.Time, to time.Time) ([]model.Order, error)
}

package iface

import (
	"oap-reposts/model"
	"time"
)

type ReportRepo interface {
	AddOrder(order *model.Order) error
	GetAll(from time.Time, to time.Time) ([]model.Order, error)
}

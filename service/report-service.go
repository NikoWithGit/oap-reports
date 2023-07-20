package service

import (
	"oap-reposts/iface"
	"oap-reposts/model"
	"time"
)

type ReportService struct {
	rr iface.ReportRepo
}

func NewReportService(rr iface.ReportRepo) *ReportService {
	return &ReportService{rr}
}

func (rs *ReportService) GetAll(from time.Time, to time.Time) ([]model.Order, error) {
	return rs.rr.GetAll(from, to)
}

func (rs *ReportService) AddOrder(order *model.Order) error {
	return rs.rr.AddOrder(order)
}

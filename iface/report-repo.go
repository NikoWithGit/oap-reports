package iface

type ReportRepo interface {
	AddReport(report string) error
}

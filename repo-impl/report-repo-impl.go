package repoimpl

import "database/sql"

type ReportRepoImpl struct {
	db *sql.DB
}

func NewReportRepoImpl(db *sql.DB) *ReportRepoImpl {
	return &ReportRepoImpl{db}
}

func (rri *ReportRepoImpl) AddReport(report string) error {
	_, err := rri.db.Query("INSERT INTO completed_orders(info) VALUES($1)", report)
	return err
}

package controller

import (
	"net/http"
	"oap-reposts/iface"
	"time"

	"github.com/gin-gonic/gin"
)

type ReportController struct {
	rs     iface.ReportService
	logger iface.Ilogger
}

func NewReportController(rs iface.ReportService, l iface.Ilogger) *ReportController {
	return &ReportController{rs, l}
}

func (rc *ReportController) GetAll(ctx *gin.Context) {
	fromStr := ctx.Query("from")
	toStr := ctx.Query("to")
	if fromStr == "" || toStr == "" {
		ctx.String(http.StatusBadRequest, "'from' and 'to' parameters are required")
		return
	}
	from, err := time.Parse("2006-01-02", fromStr)
	if err != nil {
		ctx.String(http.StatusBadRequest, "Wrong 'from' value")
		return
	}
	to, err := time.Parse("2006-01-02", toStr)
	if err != nil {
		ctx.String(http.StatusBadRequest, "Wrong 'to' value")
		return
	}
	orders, err := rc.rs.GetAll(from, to)
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, &orders)
}

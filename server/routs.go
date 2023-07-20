package server

import "oap-reposts/controller"

func (s *server) RegisterRoutes(rc *controller.ReportController) {
	s.Gin.GET("/reports", rc.GetAll)
}

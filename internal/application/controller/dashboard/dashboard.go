package dashboard

import (
	"gochat/internal/pkg/response"
)

type DashboardController struct {
	response *response.SvcRequest
}

func NewDashboardController() *DashboardController {
	return &DashboardController{
		response: response.NewSvcRequest(),
	}
}

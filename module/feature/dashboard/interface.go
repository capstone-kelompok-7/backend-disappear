package dashboard

import "github.com/labstack/echo/v4"

type RepositoryDashboardInterface interface {
	CountProducts() (int64, error)
	CountUsers() (int64, error)
	CountOrder() (int64, error)
	CountIncome() (float64, error)
}

type ServiceDashboardInterface interface {
	GetCardDashboard() (int64, int64, int64, float64, error)
}

type HandlerDashboardInterface interface {
	GetCardDashboard() echo.HandlerFunc
}

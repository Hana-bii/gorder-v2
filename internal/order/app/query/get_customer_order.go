package query

import (
	"context"
	"github.com/Hana-bii/gorder-v2/common/tracing"

	"github.com/Hana-bii/gorder-v2/common/decorator"
	domain "github.com/Hana-bii/gorder-v2/order/domain/order"
	"github.com/sirupsen/logrus"
)

type GetCustomerOrder struct {
	CustomerID string
	OrderID    string
}

type GetCustomerOrderHandler decorator.QueryHandler[GetCustomerOrder, *domain.Order]

type getCustomerOrderHandler struct {
	orderRepo domain.Repository
}

func NewGetCustomerOrderHandler(
	orderRepo domain.Repository,
	logger *logrus.Entry,
	metricsClient decorator.MetricsClient,
) GetCustomerOrderHandler {
	if orderRepo == nil {
		panic("nil orderRepo")
	}
	return decorator.ApplyQueryDecorators[GetCustomerOrder, *domain.Order](
		getCustomerOrderHandler{orderRepo: orderRepo},
		logger,
		metricsClient,
	)
}

// 具体查询方法
func (g getCustomerOrderHandler) Handle(ctx context.Context, query GetCustomerOrder) (*domain.Order, error) {
	_, span := tracing.Start(ctx, "handle...")

	o, err := g.orderRepo.Get(ctx, query.OrderID, query.CustomerID)
	if err != nil {
		return nil, err
	}
	span.AddEvent("get_success")
	span.End()
	return o, nil
}

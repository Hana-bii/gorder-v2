package ports

import (
	"context"
	"github.com/Hana-bii/gorder-v2/order/convertor"

	"github.com/Hana-bii/gorder-v2/common/genproto/orderpb"
	"github.com/Hana-bii/gorder-v2/order/app"
	"github.com/Hana-bii/gorder-v2/order/app/command"
	"github.com/Hana-bii/gorder-v2/order/app/query"
	domain "github.com/Hana-bii/gorder-v2/order/domain/order"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GRPCServer struct {
	app app.Application
}

func NewGRPCServer(app app.Application) *GRPCServer {
	return &GRPCServer{app: app}
}

func (G GRPCServer) CreateOrder(ctx context.Context, request *orderpb.CreateOrderRequest) (*empty.Empty, error) {
	_, err := G.app.Commands.CreateOrder.Handle(ctx, command.CreateOrder{
		CustomerID: request.CustomerID,
		Items:      convertor.NewItemWithQuantityConvertor().ProtosToEntities(request.Items),
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &empty.Empty{}, nil
}

func (G GRPCServer) GetOrder(ctx context.Context, request *orderpb.GetOrderRequest) (*orderpb.Order, error) {
	o, err := G.app.Queries.GetCustomerOrder.Handle(ctx, query.GetCustomerOrder{
		CustomerID: request.CustomerID,
		OrderID:    request.OrderID,
	})
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}
	return convertor.NewOrderConvertor().EntityToProto(o), nil

}

func (G GRPCServer) UpdateOrder(ctx context.Context, request *orderpb.Order) (_ *empty.Empty, err error) {
	logrus.Infof("order_grpc || request_in || request=#{request}")
	order, err := domain.NewOrder(
		request.ID,
		request.CustomerID,
		request.Status,
		request.PaymentLink,
		convertor.NewItemConvertor().ProtosToEntities(request.Items),
	)
	if err != nil {
		err = status.Error(codes.Internal, err.Error())
		return nil, err
	}
	logrus.Infof("order_grpc || order=#{order}")
	_, err = G.app.Commands.UpdateOrder.Handle(ctx, command.UpdateOrder{
		Order: order,
		UpdateFn: func(ctx context.Context, order *domain.Order) (*domain.Order, error) {
			return order, nil
		},
	})
	return nil, err
}

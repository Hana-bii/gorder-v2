package service

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"time"

	"github.com/Hana-bii/gorder-v2/common/broker"
	grpcClient "github.com/Hana-bii/gorder-v2/common/client"
	"github.com/Hana-bii/gorder-v2/common/metrics"
	"github.com/Hana-bii/gorder-v2/order/adapters"
	"github.com/Hana-bii/gorder-v2/order/adapters/grpc"
	"github.com/Hana-bii/gorder-v2/order/app"
	"github.com/Hana-bii/gorder-v2/order/app/command"
	"github.com/Hana-bii/gorder-v2/order/app/query"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func NewApplication(ctx context.Context) (app.Application, func()) {
	// 实际的存储基础设施在这里更换，不影响业务逻辑

	stockClient, closeStockClient, err := grpcClient.NewStockGRPCClient(ctx)
	if err != nil {
		panic(err)
	}

	// 注册消息队列
	ch, closeCh := broker.Connect(
		viper.GetString("rabbitmq.user"),
		viper.GetString("rabbitmq.passwd"),
		viper.GetString("rabbitmq.host"),
		viper.GetString("rabbitmq.port"),
	)
	stockGRPC := grpc.NewStockGRPC(stockClient)
	return newApplication(ctx, stockGRPC, ch), func() {
		_ = closeStockClient()
		_ = closeCh()
		_ = ch.Close()
	}

}

func newApplication(_ context.Context, stockGRPC query.StockService, ch *amqp.Channel) app.Application {
	//orderRepo := adapters.NewMemoryOrderRepository()
	mongoClient := newMongoClient()
	orderRepo := adapters.NewOrderRepositoryMongo(mongoClient)
	logger := logrus.NewEntry(logrus.StandardLogger())
	metricsClient := metrics.TodoMetrics{}

	// 注入服务
	return app.Application{
		Commands: app.Commands{
			CreateOrder: command.NewCreateOrderHandler(orderRepo, stockGRPC, ch, logger, metricsClient),
			UpdateOrder: command.NewUpdateOrderHandler(orderRepo, logger, metricsClient),
		},
		Queries: app.Queries{
			GetCustomerOrder: query.NewGetCustomerOrderHandler(orderRepo, logger, metricsClient),
		},
	}
}

func newMongoClient() *mongo.Client {
	// 获取客户端信息
	uri := fmt.Sprintf(
		"mongodb://%s:%s@%s:%s",
		viper.GetString("mongo.user"),
		viper.GetString("mongo.password"),
		viper.GetString("mongo.host"),
		viper.GetString("mongo.port"))
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	c, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}
	if err = c.Ping(ctx, readpref.Primary()); err != nil {
		panic(err)
	}
	return c
}

package order

import (
	"errors"
	"github.com/Hana-bii/gorder-v2/common/genproto/orderpb"
)

type Order struct {
	ID          string
	CustomerID  string
	Status      string
	PaymentLink string
	Items       []*orderpb.Item
}

func NewOrder(id, customerID, status, paymentLink string, items []*orderpb.Item) (*Order, error) {
	if id == "" {
		return nil, errors.New("empty id")
	}
	if customerID == "" {
		return nil, errors.New("empty customerID")
	}
	if status != "" {
		return nil, errors.New("empty status")
	}
	if items == nil {
		return nil, errors.New("empty items")
	}
	return &Order{
		ID:          id,
		CustomerID:  customerID,
		Status:      status,
		PaymentLink: paymentLink,
		Items:       items,
	}, nil
}

func (o *Order) ToProto() *orderpb.Order {
	return &orderpb.Order{
		ID:          o.ID,
		CustomerID:  o.CustomerID,
		Status:      o.Status,
		PaymentLink: o.PaymentLink,
		Items:       o.Items,
	}
}

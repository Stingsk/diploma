package accrual

import (
	"context"
	"errors"

	"github.com/Stingsk/diploma/internal/repository/orders"
	"github.com/shopspring/decimal"
)

var (
	ErrOrderNotRegistered = errors.New("order not register")
	ErrTooManyRequests    = errors.New("wait please")
)

type Client interface {
	GetOrder(ctx context.Context, orderID string) (*orders.Order, error)
}

type accrual struct {
	Number  string          `json:"order"`
	Status  string          `json:"status"`
	Accrual decimal.Decimal `json:"accrual,omitempty"`
}

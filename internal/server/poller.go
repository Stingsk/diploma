package server

import (
	"context"
	"time"

	"github.com/Stingsk/diploma/internal/accrual"
	"github.com/Stingsk/diploma/internal/repository/orders"
	"github.com/sirupsen/logrus"
)

const (
	pollTimeout = 1 * time.Second
)

type PollerConfig struct {
	PollInterval   time.Duration
	AccrualAddress string
}

type PollerWorker struct {
	Cfg PollerConfig
}

func (pw *PollerWorker) Run(ctx context.Context, ordersStore orders.Store) {
	pollTicker := time.NewTicker(pw.Cfg.PollInterval)
	defer pollTicker.Stop()

	accrualClient := accrual.NewAccrualClient(pw.Cfg.AccrualAddress)

	for {
		select {
		case <-ctx.Done():
			return
		case <-pollTicker.C:
			UpdateOrders(ctx, accrualClient, ordersStore)
		}
	}
}

func UpdateOrders(ctx context.Context, accrualClient accrual.Client, ordersStore orders.Store) {
	skipStatuses := map[string]struct{}{
		"REGISTERED": {},
	}

	getContext, getCancel := context.WithTimeout(ctx, pollTimeout)
	defer getCancel()

	ordersSlice, err := ordersStore.GetUnprocessedOrders(getContext)
	if err != nil {
		logrus.Error("poller couldn't get orders from store")
	}

	for _, order := range ordersSlice {
		accrualOrder, err := accrualClient.GetOrder(getContext, order.Number)
		if err != nil {
			logrus.Error("fail to get order from accrual")

			continue
		}
		if _, ok := skipStatuses[accrualOrder.Status]; ok {
			continue
		}

		if err := ordersStore.UpdateOrder(ctx, accrualOrder); err != nil {
			logrus.Error("fail to update " + accrualOrder.Number + "order")
		}
	}
}

package server

import (
	"context"
	"testing"

	mockaccrual "github.com/Stingsk/diploma/internal/accrual/moks"
	"github.com/Stingsk/diploma/internal/repository/orders"
	mockorders "github.com/Stingsk/diploma/internal/repository/orders/mocks"
	"github.com/golang/mock/gomock"
	"github.com/shopspring/decimal"
)

type testOrders struct {
	order        orders.Order
	accrualOrder *orders.Order
}

type testPoller struct {
	name       string
	buildStubs func(client *mockaccrual.MockClient, store *mockorders.MockStore)
}

func TestUpdateOrders(t *testing.T) {
	ordersForTests := []testOrders{
		{
			order: orders.Order{
				Number: "8954741254",
				Status: "NEW",
			},
			accrualOrder: &orders.Order{
				Number:  "8954741254",
				Status:  "PROCESSED",
				Accrual: decimal.NewFromInt(500),
			},
		},
		{
			order: orders.Order{
				Number: "4565323212",
				Status: "NEW",
			},
			accrualOrder: &orders.Order{
				Number: "4565323212",
				Status: "REGISTERED",
			},
		},
		{
			order: orders.Order{
				Number: "874587458",
				Status: "NEW",
			},
			accrualOrder: &orders.Order{
				Number: "874587458",
				Status: "INVALID",
			},
		},
	}
	tests := []testPoller{
		{
			name: "positive test #1",
			buildStubs: func(client *mockaccrual.MockClient, store *mockorders.MockStore) {
				store.EXPECT().GetUnprocessedOrders(gomock.Any()).Return([]orders.Order{ordersForTests[0].order}, nil)
				client.EXPECT().GetOrder(gomock.Any(),
					ordersForTests[0].order.Number).Return(ordersForTests[0].accrualOrder, nil).Times(1)
				store.EXPECT().UpdateOrder(gomock.Any(), ordersForTests[0].accrualOrder).Return(nil).Times(1)
			},
		},
		{
			name: "positive test #2",
			buildStubs: func(client *mockaccrual.MockClient, store *mockorders.MockStore) {
				store.EXPECT().GetUnprocessedOrders(gomock.Any()).Return([]orders.Order{ordersForTests[1].order}, nil)
				client.EXPECT().GetOrder(gomock.Any(),
					ordersForTests[1].order.Number).Return(ordersForTests[1].accrualOrder, nil).Times(1)
				store.EXPECT().UpdateOrder(gomock.Any(), ordersForTests[1].accrualOrder).Return(nil).Times(0)
			},
		},
		{
			name: "negative test #1",
			buildStubs: func(client *mockaccrual.MockClient, store *mockorders.MockStore) {
				store.EXPECT().GetUnprocessedOrders(gomock.Any()).Return([]orders.Order{ordersForTests[2].order}, nil)
				client.EXPECT().GetOrder(gomock.Any(),
					ordersForTests[2].order.Number).Return(ordersForTests[2].accrualOrder, nil).Times(1)
				store.EXPECT().UpdateOrder(gomock.Any(), ordersForTests[2].accrualOrder).Return(nil).Times(1)
			},
		},
	}

	client, store := getMocks(t)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.buildStubs(client, store)
			UpdateOrders(context.Background(), client, store)
		})
	}
}

func getMocks(t *testing.T) (*mockaccrual.MockClient, *mockorders.MockStore) {
	t.Helper()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	a := mockaccrual.NewMockClient(ctrl)
	s := mockorders.NewMockStore(ctrl)

	return a, s
}

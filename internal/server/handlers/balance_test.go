package handlers

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Stingsk/diploma/internal/repository/orders"
	"github.com/Stingsk/diploma/internal/repository/orders/mocks"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	"github.com/golang/mock/gomock"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type wantBalance struct {
	code int
	data string
}

type testBalance struct {
	name       string
	method     string
	url        string
	authHeader string
	buildStubs func(store *mock_orders.MockStore)
	want       wantBalance
}

func TestBalanceHandlers(t *testing.T) {
	processedOrders := []orders.Order{
		{
			Accrual: decimal.NewFromFloat(200.2),
		},
		{
			Accrual: decimal.NewFromFloat(251.0),
		},
	}
	withdrawals := []orders.Withdraw{
		{
			Sum: decimal.NewFromFloat(25.1),
		},
		{
			Sum: decimal.NewFromFloat(25.3),
		},
	}
	testOrders := []testBalance{
		{
			name:       "get balance",
			method:     http.MethodGet,
			url:        "/api/user/balance",
			authHeader: authHeader,
			want: wantBalance{
				code: http.StatusOK,
				data: "{\"current\":400.8,\"withdrawn\":50.4}\n",
			},
			buildStubs: func(store *mock_orders.MockStore) {
				store.EXPECT().GetProcessedOrders(gomock.Any(), "test").Return(processedOrders, nil).Times(1)
				store.EXPECT().GetWithdrawals(gomock.Any(), "test").Return(withdrawals, nil).Times(1)
			},
		},
		{
			name:       "get withdrawals",
			method:     http.MethodGet,
			url:        "/api/user/balance/withdrawals",
			authHeader: authHeader,
			want: wantBalance{
				code: http.StatusOK,
				data: "[{\"order\":\"5896541234\",\"sum\":50, \"processed_at\":\"2022-06-06T10:15:10.987Z\"}]\n",
			},
			buildStubs: func(store *mock_orders.MockStore) {
				store.EXPECT().GetWithdrawals(gomock.Any(), "test").Return([]orders.Withdraw{
					{
						Order:       "5896541234",
						Sum:         decimal.NewFromFloat(50),
						ProcessedAt: getDate(),
					},
				}, nil)
			},
		},
	}

	jwtToken := jwtauth.New("HS256", []byte("test"), []byte("test"))

	mux := chi.NewRouter()
	store := getBalanceStore(t)
	RegisterPrivateHandlers(mux, store, jwtToken)

	ts := httptest.NewServer(mux)
	defer ts.Close()

	for _, tt := range testOrders {
		t.Run(tt.name, func(t *testing.T) {
			tt.buildStubs(store)
			testBalanceRequest(t, ts, tt)
		})
	}
}

func testBalanceRequest(t *testing.T, ts *httptest.Server, testData testBalance) {
	t.Helper()

	req, err := http.NewRequest(testData.method, ts.URL+testData.url, nil)
	require.NoError(t, err)

	req.Header.Set("Authorization", testData.authHeader)

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)

	assert.Equal(t, testData.want.code, resp.StatusCode)

	respBody, err := ioutil.ReadAll(resp.Body)
	assert.JSONEq(t, testData.want.data, string(respBody))

	require.NoError(t, err)

	err = resp.Body.Close()
	if err != nil {
		return
	}
}

func getBalanceStore(t *testing.T) *mock_orders.MockStore {
	t.Helper()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	s := mock_orders.NewMockStore(ctrl)

	return s
}

func getDate() time.Time {
	layout := "2006-01-02T15:04:05.000Z"
	str := "2022-06-06T10:15:10.987Z"
	t, _ := time.Parse(layout, str)

	return t
}

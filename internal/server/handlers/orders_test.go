package handlers

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Stingsk/diploma/internal/repository/orders"
	"github.com/Stingsk/diploma/internal/repository/orders/mocks"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	authHeader = "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9." +
		"eyJsb2dpbiI6InRlc3QifQ." +
		"pV-CiPiB0QjqUno0SPHOlO4NWx_Gd0vHMrtRIjUqdhQ"
)

type wantOrders struct {
	code int
	data string
}

type testOrder struct {
	name       string
	method     string
	url        string
	order      string
	authHeader string
	buildStubs func(store *mock_orders.MockStore)
	want       wantOrders
}

func TestOrdersHandlers(t *testing.T) {
	testOrders := []testOrder{
		{
			name:       "positive test #1",
			method:     http.MethodPost,
			url:        "/api/user/orders",
			order:      "5874589614521",
			authHeader: authHeader,
			want: wantOrders{
				code: http.StatusAccepted,
				data: "",
			},
			buildStubs: func(store *mock_orders.MockStore) {
				store.EXPECT().CreateOrder(gomock.Any(), "test", "5874589614521")
			},
		},
		{
			name:       "positive test #2",
			method:     http.MethodPost,
			url:        "/api/user/orders",
			order:      "5874589614521",
			authHeader: authHeader,
			want: wantOrders{
				code: http.StatusOK,
				data: "",
			},
			buildStubs: func(store *mock_orders.MockStore) {
				store.EXPECT().CreateOrder(gomock.Any(), "test", "5874589614521").Return(orders.ErrOrderExists)
			},
		},
		{
			name:       "negative test #1",
			method:     http.MethodPost,
			url:        "/api/user/orders",
			order:      "1111",
			authHeader: authHeader,
			want: wantOrders{
				code: http.StatusUnprocessableEntity,
				data: "",
			},
			buildStubs: func(store *mock_orders.MockStore) {
				store.EXPECT().CreateOrder(gomock.Any(), "test", "1111").Times(0)
			},
		},
		{
			name:       "negative test #2",
			method:     http.MethodPost,
			url:        "/api/user/orders",
			order:      "5874589614521",
			authHeader: authHeader,
			want: wantOrders{
				code: http.StatusConflict,
				data: "",
			},
			buildStubs: func(store *mock_orders.MockStore) {
				store.EXPECT().CreateOrder(gomock.Any(), "test", "5874589614521").Return(orders.ErrOtherOrderExists)
			},
		},
		{
			name:       "negative test #3",
			method:     http.MethodPost,
			url:        "/api/user/orders",
			order:      "5874589614521",
			authHeader: "",
			want: wantOrders{
				code: http.StatusUnauthorized,
				data: "",
			},
			buildStubs: func(store *mock_orders.MockStore) {
				store.EXPECT().CreateOrder(gomock.Any(), "test", "5874589614521").Times(0)
			},
		},
	}

	jwtToken := jwtauth.New("HS256", []byte("test"), []byte("test"))

	mux := chi.NewRouter()
	store := getOrdersStore(t)
	RegisterPrivateHandlers(mux, store, jwtToken)

	ts := httptest.NewServer(mux)
	defer ts.Close()

	for _, tt := range testOrders {
		t.Run(tt.name, func(t *testing.T) {
			tt.buildStubs(store)
			testOrdersRequest(t, ts, tt)
		})
	}
}

func testOrdersRequest(t *testing.T, ts *httptest.Server, testData testOrder) {
	t.Helper()

	var body bytes.Buffer
	body.WriteString(testData.order)

	req, err := http.NewRequest(testData.method, ts.URL+testData.url, &body)
	require.NoError(t, err)

	req.Header.Set("Content-Type", "text/plain")
	req.Header.Set("Authorization", testData.authHeader)

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, testData.want.code, resp.StatusCode)
}

func getOrdersStore(t *testing.T) *mock_orders.MockStore {
	t.Helper()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	s := mock_orders.NewMockStore(ctrl)

	return s
}

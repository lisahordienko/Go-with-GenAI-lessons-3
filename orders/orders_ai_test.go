package orders

import (
	"errors"
	"testing"
)

type aiOrderStoreMock struct {
	execFn func(query string, args ...any) error
}

func (m *aiOrderStoreMock) Exec(query string, args ...any) error {
	if m.execFn != nil {
		return m.execFn(query, args...)
	}
	return nil
}

func TestPlaceOrder_AIGeneratedMock(t *testing.T) {
	const (
		query   = "INSERT INTO orders (id, amount) VALUES (?, ?)"
		orderID = "order-ai-42"
		amount  = 99.99
	)

	var (
		called   bool
		gotQuery string
		gotArgs  []any
	)

	store := &aiOrderStoreMock{
		execFn: func(query string, args ...any) error {
			called = true
			gotQuery = query
			gotArgs = append([]any(nil), args...)
			return nil
		},
	}

	svc := NewOrderService(store)
	if err := svc.PlaceOrder(orderID, amount); err != nil {
		t.Fatalf("PlaceOrder повернув помилку: %v", err)
	}

	if !called {
		t.Fatal("Exec не був викликаний")
	}
	if gotQuery != query {
		t.Fatalf("Exec query = %q, want %q", gotQuery, query)
	}
	if len(gotArgs) != 2 {
		t.Fatalf("len(args) = %d, want 2", len(gotArgs))
	}
	if gotArgs[0] != orderID {
		t.Fatalf("args[0] = %v, want %q", gotArgs[0], orderID)
	}
	if gotArgs[1] != amount {
		t.Fatalf("args[1] = %v, want %v", gotArgs[1], amount)
	}
}

func TestPlaceOrder_AIGeneratedMock_PropagatesError(t *testing.T) {
	expectedErr := errors.New("simulated db failure")
	store := &aiOrderStoreMock{
		execFn: func(query string, args ...any) error {
			return expectedErr
		},
	}

	svc := NewOrderService(store)
	err := svc.PlaceOrder("order-ai-error", 12.5)
	if !errors.Is(err, expectedErr) {
		t.Fatalf("PlaceOrder повернув %v, want %v", err, expectedErr)
	}
}

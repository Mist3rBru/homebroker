package entity

import (
	"time"

	"github.com/google/uuid"
)

type Transaction struct {
	Id           string
	SellingOrder *Order
	BuyingOrder  *Order
	Shares       int
	Price        float64
	Total        float64
	DateTime     time.Time
}

func NewTransaction(sellingOrder *Order, buyingOrder *Order, shares int, price float64) *Transaction {
	return &Transaction{
		Id:           uuid.New().String(),
		SellingOrder: sellingOrder,
		BuyingOrder:  buyingOrder,
		Shares:       shares,
		Price:        price,
		Total:        float64(shares) * price,
		DateTime:     time.Now(),
	}
}

func (t *Transaction) CalculateTotal(shares int, price float64) float64 {
	return float64(shares) * price
}

func (t *Transaction) CloseSellOrder() {
	if t.SellingOrder.PendingShares == 0 {
		t.SellingOrder.Status = "CLOSED"
	}
}

func (t *Transaction) CloseBuyOrder() {
	if t.BuyingOrder.PendingShares == 0 {
		t.BuyingOrder.Status = "CLOSED"
	}
}

func (t *Transaction) CloseOrders() {
	t.CloseSellOrder()
	t.CloseBuyOrder()
}

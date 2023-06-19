package entity

type Order struct {
	Id            string
	Investor      *Investor
	Asset         *Asset
	Shares        int
	PendingShares int
	Price         float64
	Type          string
	Status        string
	Transactions  []*Transaction
}

func NewOrder(orderId string, investor *Investor, asset *Asset, shares int, price float64, orderType string) *Order {
	return &Order{
		Id:            orderId,
		Investor:      investor,
		Asset:         asset,
		Shares:        shares,
		PendingShares: shares,
		Price:         price,
		Type:          orderType,
		Status:        "OPEN",
		Transactions:  []*Transaction{},
	}
}

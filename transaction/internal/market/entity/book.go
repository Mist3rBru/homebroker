package entity

import (
	"container/heap"
	"sync"
)

type Book struct {
	Orders        []*Order
	Transactions  []*Transaction
	OrdersChan    chan *Order
	OrdersChanOut chan *Order
	Wg            *sync.WaitGroup
}

func NewBook(orderChan chan *Order, orderChanOut chan *Order, wg *sync.WaitGroup) *Book {
	return &Book{
		Orders:        []*Order{},
		Transactions:  []*Transaction{},
		OrdersChan:    orderChan,
		OrdersChanOut: orderChanOut,
		Wg:            wg,
	}
}

func (b *Book) Trade() {
	buyOrders := make(map[string]*OrderQueue)
	sellOrders := make(map[string]*OrderQueue)

	for order := range b.OrdersChan {
		asset := order.Asset.Id

		if buyOrders[asset] == nil {
			buyOrders[asset] = NewOrderQueue()
			heap.Init(buyOrders[asset])
		}

		if sellOrders[asset] == nil {
			sellOrders[asset] = NewOrderQueue()
			heap.Init(sellOrders[asset])
		}

		if order.Type == "BUY" {
			buyOrder := order
			buyOrders[asset].Push(buyOrder)
			if sellOrders[asset].Len() > 0 && sellOrders[asset].Orders[0].Price <= buyOrder.Price {
				sellOrder := sellOrders[asset].Pop().(*Order)
				if sellOrder.PendingShares > 0 {
					transaction := NewTransaction(sellOrder, buyOrder, buyOrder.Shares, sellOrder.Price)
					b.AddTransation(transaction, b.Wg)
					sellOrder.Transactions = append(sellOrder.Transactions, transaction)
					buyOrder.Transactions = append(buyOrder.Transactions, transaction)
					b.OrdersChanOut <- sellOrder
					b.OrdersChanOut <- buyOrder
					if sellOrder.PendingShares > 0 {
						heap.Push(sellOrders[asset], sellOrder)
					}
				}
			}
		} else if order.Type == "SELL" {
			sellOrder := order
			sellOrders[asset].Push(sellOrder)
			if buyOrders[asset].Len() > 0 && buyOrders[asset].Orders[0].Price >= sellOrder.Price {
				buyOrder := buyOrders[asset].Pop().(*Order)
				if buyOrder.PendingShares > 0 {
					transaction := NewTransaction(sellOrder, buyOrder, sellOrder.Shares, sellOrder.Price)
					b.AddTransation(transaction, b.Wg)
					buyOrder.Transactions = append(buyOrder.Transactions, transaction)
					sellOrder.Transactions = append(sellOrder.Transactions, transaction)
					b.OrdersChanOut <- buyOrder
					b.OrdersChanOut <- sellOrder
					if buyOrder.PendingShares > 0 {
						heap.Push(buyOrders[asset], buyOrder)
					}
				}
			}
		}
	}
}

func (b *Book) AddTransation(transaction *Transaction, wg *sync.WaitGroup) {
	defer wg.Done()

	sellingShares := transaction.SellingOrder.PendingShares
	buyingShares := transaction.BuyingOrder.PendingShares

	minShares := sellingShares
	if buyingShares < minShares {
		minShares = buyingShares
	}

	transaction.SellingOrder.Investor.UpdateAssetPosition(transaction.SellingOrder.Asset.Id, -minShares)
	transaction.SellingOrder.PendingShares -= minShares

	transaction.BuyingOrder.Investor.UpdateAssetPosition(transaction.BuyingOrder.Asset.Id, minShares)
	transaction.BuyingOrder.PendingShares -= minShares

	transaction.Total = transaction.CalculateTotal(transaction.Shares, transaction.SellingOrder.Price)

	transaction.CloseOrders()
	b.Transactions = append(b.Transactions, transaction)
}

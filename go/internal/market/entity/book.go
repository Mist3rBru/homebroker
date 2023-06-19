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
	buyOrders := NewOrderQueue()
	sellOrders := NewOrderQueue()

	heap.Init(buyOrders)
	heap.Init(sellOrders)

	for order := range b.OrdersChan {
		if order.Type == "BUY" {
			buyOrder := order
			buyOrders.Push(buyOrder)
			if sellOrders.Len() > 0 && sellOrders.Orders[0].Price <= buyOrder.Price {
				sellOrder := sellOrders.Pop().(*Order)
				if sellOrder.PendingShares > 0 {
					transaction := NewTransaction(sellOrder, buyOrder, buyOrder.Shares, sellOrder.Price)
					b.AddTransation(transaction, b.Wg)
					sellOrder.Transactions = append(sellOrder.Transactions, transaction)
					buyOrder.Transactions = append(buyOrder.Transactions, transaction)
					b.OrdersChanOut <- sellOrder
					b.OrdersChanOut <- buyOrder
					if sellOrder.PendingShares > 0 {
						heap.Push(sellOrders, sellOrder)
					}
				}
			}
		} else if order.Type == "SELL" {
			sellOrder := order
			sellOrders.Push(sellOrder)
			if buyOrders.Len() > 0 && buyOrders.Orders[0].Price >= sellOrder.Price {
				buyOrder := buyOrders.Pop().(*Order)
				if buyOrder.PendingShares > 0 {
					transaction := NewTransaction(sellOrder, buyOrder, sellOrder.Shares, sellOrder.Price)
					b.AddTransation(transaction, b.Wg)
					buyOrder.Transactions = append(buyOrder.Transactions, transaction)
					sellOrder.Transactions = append(sellOrder.Transactions, transaction)
					b.OrdersChanOut <- buyOrder
					b.OrdersChanOut <- sellOrder
					if buyOrder.PendingShares > 0 {
						heap.Push(buyOrders, buyOrder)
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

package entity

type OrderQueue struct {
	Orders []*Order
	Length int
}

func (oq *OrderQueue) Less(i int, j int) bool {
	return oq.Orders[i].Price < oq.Orders[j].Price
}

func (oq *OrderQueue) Swap(i int, j int) {
	oq.Orders[i], oq.Orders[j] = oq.Orders[j], oq.Orders[i]
}

func (oq *OrderQueue) Len() int {
	return oq.Length
}

func (oq *OrderQueue) Push(x interface{}) {
	oq.Orders = append(oq.Orders, x.(*Order))
	oq.Length++
}

func (oq *OrderQueue) Pop() interface{} {
	lastItem := oq.Last()
	oq.Orders = oq.Orders[0 : oq.Length-1]
	oq.Length--
	return lastItem
}

func (oq *OrderQueue) Last() *Order {
	return oq.Orders[oq.Length-1]
}

func (oq *OrderQueue) Lowest() *Order {
	lowest := oq.Orders[0]

	for _, order := range oq.Orders {
		if order.Price < lowest.Price {
			lowest = order
		}
	}

	return lowest
}

func NewOrderQueue() *OrderQueue {
	return &OrderQueue{
		Orders: []*Order{},
	}
}

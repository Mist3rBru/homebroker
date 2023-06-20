package transformer

import (
	"github.com/Mist3rBru/homebroker/internal/market/dto"
	"github.com/Mist3rBru/homebroker/internal/market/entity"
)

func TransformInput(input dto.TradeInput) *entity.Order {
	asset := entity.NewAsset(input.AssetId, input.AssetId, 1000)
	investor := entity.NewInvestor(input.InvestorId)
	order := entity.NewOrder(input.OrderId, investor, asset, input.Shares, input.Price, input.OrderType)

	if input.CurrentShares > 0 {
		assetPosition := entity.NewInvestorAssetPosition(input.AssetId, input.Shares)
		investor.AddAssetPosition(assetPosition)
	}

	return order
}

func TransformOutput(order *entity.Order) *dto.OrderOutput {
	output := &dto.OrderOutput{
		OrderId:    order.Id,
		InvestorId: order.Investor.Id,
		AssetId:    order.Asset.Id,
		Status:     order.Status,
		Partial:    order.PendingShares,
		Shares:     order.Shares,
	}

	var transactionsOutput []*dto.TransactionOutput
	for _, t := range order.Transactions {
		transactionOutput := &dto.TransactionOutput{
			TransactionId: t.Id,
			BuyerId:       t.BuyingOrder.Investor.Id,
			SellerId:      t.SellingOrder.Investor.Id,
			Price:         t.Price,
			AssetId:       order.Asset.Id,
			Shares:        t.Shares,
		}
		transactionsOutput = append(transactionsOutput, transactionOutput)
	}

	output.TransactionsOutput = transactionsOutput
	return output
}

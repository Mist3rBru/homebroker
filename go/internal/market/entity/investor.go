package entity

type Investor struct {
	Id            string
	Name          string
	AssetPosition []*InvestorAssetPosition
}

func NewInvestor(id string) *Investor {
	return &Investor{
		Id:            id,
		AssetPosition: []*InvestorAssetPosition{},
	}
}

func (i *Investor) AddAssetPosition(assetPosition *InvestorAssetPosition) {
	i.AssetPosition = append(i.AssetPosition, assetPosition)
}

func (i *Investor) UpdateAssetPosition(assetId string, shares int) {
	assetPosition := i.GetAssetPosition(assetId)
	if assetPosition == nil {
		i.AssetPosition = append(i.AssetPosition, NewInvestorAssetPosition(assetId, shares))
	} else {
		assetPosition.Shares += shares
	}
}

func (i *Investor) GetAssetPosition(assetId string) *InvestorAssetPosition {
	for _, assetPosition := range i.AssetPosition {
		if assetPosition.AssetId == assetId {
			return assetPosition
		}
	}
	return nil
}

type InvestorAssetPosition struct {
	AssetId string
	Shares  int
}

func NewInvestorAssetPosition(assetId string, shares int) *InvestorAssetPosition {
	return &InvestorAssetPosition{
		AssetId: assetId,
		Shares:  shares,
	}
}

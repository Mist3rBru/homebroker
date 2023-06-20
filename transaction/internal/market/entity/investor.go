package entity

type Investor struct {
	Id            string
	Name          string
	AssetPositions []*InvestorAssetPosition
}

func NewInvestor(id string) *Investor {
	return &Investor{
		Id:            id,
		AssetPositions: []*InvestorAssetPosition{},
	}
}

func (i *Investor) AddAssetPosition(assetPosition *InvestorAssetPosition) {
	i.AssetPositions = append(i.AssetPositions, assetPosition)
}

func (i *Investor) UpdateAssetPosition(assetId string, shares int) {
	assetPosition := i.GetAssetPosition(assetId)
	if assetPosition == nil {
		i.AssetPositions= append(i.AssetPositions, NewInvestorAssetPosition(assetId, shares))
	} else {
		assetPosition.Shares += shares
	}
}

func (i *Investor) GetAssetPosition(assetId string) *InvestorAssetPosition {
	for _, assetPosition := range i.AssetPositions {
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

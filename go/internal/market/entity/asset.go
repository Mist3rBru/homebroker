package entity

type Asset struct {
	Id           string
	Name         string
	MarketVolume int
}

func NewAsset(id string, name string, marketVolume int) *Asset {
	return &Asset{
		Id:           id,
		Name:         name,
		MarketVolume: marketVolume,
	}
}

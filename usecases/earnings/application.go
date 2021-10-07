package earnings

import (
	"stockfyApi/entity"
	"time"
)

type Application struct {
	repo Repository
}

//NewApplication create new use case
func NewApplication(r Repository) *Application {
	return &Application{
		repo: r,
	}
}

func (a *Application) CreateEarning(earningType string, earnings float64,
	currency string, date string, country string, assetId string,
	userUid string) (*entity.Earnings, error) {

	layout := "2006-01-02"
	dateFormatted, _ := time.Parse(layout, date)
	eargningFormatted, err := entity.NewEarnings(earningType, earnings, currency,
		dateFormatted, country, assetId, userUid)
	if err != nil {
		return nil, err
	}

	earningCreated, err := a.repo.Create(*eargningFormatted)
	if err != nil {
		return nil, err
	}

	return &earningCreated[0], nil
}

func (a *Application) SearchEarningsFromAssetUser(assetId string, userUid string) (
	[]entity.Earnings, error) {
	earnings, err := a.repo.SearchFromAssetUser(assetId, userUid)
	if err != nil {
		return nil, err
	}

	return earnings, nil
}

func (a *Application) DeleteEarningsFromUser(earningId string,
	userUid string) (*string, error) {
	orderId, err := a.repo.DeleteFromUser(earningId, userUid)
	if err != nil {
		return nil, err
	}

	return &orderId, nil
}

func (a *Application) DeleteEarningsFromAsset(assetId string) ([]entity.Earnings,
	error) {

	deletedEarnings, err := a.repo.DeleteFromAsset(assetId)
	if err != nil {
		return nil, err
	}

	return deletedEarnings, nil
}

func (a *Application) DeleteEarningsFromAssetUser(assetId, userUid string) (
	*[]entity.Earnings, error) {
	deletedEarnings, err := a.repo.DeleteFromAssetUser(assetId, userUid)
	if err != nil {
		return nil, err
	}

	return &deletedEarnings, nil
}

func (a *Application) EarningsVerification(symbol string, currency string,
	earningType string, date string, earning float64) error {

	if symbol == "" || currency == "" || earningType == "" || date == "" {
		return entity.ErrInvalidApiMissedKeysBody
	}

	if earning <= 0 {
		return entity.ErrInvalidApiEarningsAmount
	}

	if !entity.ValidEarningTypes[earningType] {
		return entity.ErrInvalidApiEarningType
	}

	return nil
}

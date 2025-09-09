package application_test

import (
	usecase "api/internal/application/useCase"
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockLocationRepo struct {
	GetCityIDFunc func(ctx context.Context, countryName, cityName string) (int, error)
}

func (m *mockLocationRepo) GetCityID(ctx context.Context, countryName, cityName string) (int, error) {
	if m.GetCityIDFunc != nil {
		return m.GetCityIDFunc(ctx, countryName, cityName)
	}
	return 0, nil
}

func TestLocationService_GetCityID_Success(t *testing.T) {
	repo := &mockLocationRepo{GetCityIDFunc: func(ctx context.Context, country, city string) (int, error) {
		assert.Equal(t, "peru", country)
		assert.Equal(t, "lima", city)
		return 77, nil
	}}
	svc := usecase.NewLocationService(repo)

	id, err := svc.GetCityID(context.Background(), "peru", "lima")
	assert.NoError(t, err)
	assert.Equal(t, 77, id)
}

func TestLocationService_GetCityID_Error(t *testing.T) {
	repo := &mockLocationRepo{GetCityIDFunc: func(ctx context.Context, country, city string) (int, error) {
		return 0, errors.New("db error")
	}}
	svc := usecase.NewLocationService(repo)

	_, err := svc.GetCityID(context.Background(), "x", "y")
	assert.Error(t, err)
}

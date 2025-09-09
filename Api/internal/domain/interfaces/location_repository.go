package interfaces

import "context"

// LocationRepository expone operaciones de lectura para país/ciudad.
// Se usa para resolver city_id a partir de nombres (country, city).
type LocationRepository interface {
	GetCityID(ctx context.Context, countryName, cityName string) (int, error)
}

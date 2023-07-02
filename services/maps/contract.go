package maps

import (
	"context"
	"github.com/maxheckel/parks/models"
	"googlemaps.github.io/maps"
	"os"
)

type Service interface {
	GetLocationLatAndLng(ctx *context.Context, park *models.Park) (string, float64, float64, error)
	GetDrivingDirections(ctx *context.Context, start *models.Location, parks []*models.Park) (string, error)
}

func New() (Service, error) {
	APIKey := os.Getenv("GOOGLE_MAPS_API_KEY")
	c, err := maps.NewClient(maps.WithAPIKey(APIKey))
	if err != nil {
		return nil, err
	}
	return &google{
		client: c,
	}, nil
}

package maps

import (
	"context"
	"errors"
	"fmt"
	"github.com/maxheckel/parks/models"
	"googlemaps.github.io/maps"
	"net/url"
)

type google struct {
	client *maps.Client
}

var (
	ErrNoCity   = errors.New("no city provided for park")
	ErrNotFound = errors.New("no results found for park")
)

func (g google) GetDrivingDirections(ctx *context.Context, start *models.Location, parks []*models.Park) (string, error) {
	baseStr := fmt.Sprintf("https://www.google.com/maps/dir/?api=1&origin=%s&destination=%s&travelmode=driving&waypoints=", url.QueryEscape(start.ToString()), url.QueryEscape(parks[len(parks)-1].LatLngString()))
	for i, park := range parks {
		if i != len(parks)-1 {
			baseStr = fmt.Sprintf("%s%s%s", baseStr, url.QueryEscape(park.LatLngString()), "%7C")
		}
	}
	return baseStr, nil
}

func (g google) GetLocationLatAndLng(ctx *context.Context, park *models.Park) (string, float64, float64, error) {
	if park.City.Name == "" {
		return "", 0, 0, ErrNoCity
	}
	resp, err := g.client.FindPlaceFromText(*ctx, &maps.FindPlaceFromTextRequest{
		Input:     fmt.Sprintf("%s %s", park.Name, park.City.ToLocationName()),
		InputType: maps.FindPlaceFromTextInputTypeTextQuery,
		Fields: []maps.PlaceSearchFieldMask{
			maps.PlaceSearchFieldMaskGeometry,
		},
	})
	if err != nil {
		return "", 0, 0, err
	}
	for _, candidate := range resp.Candidates {
		return candidate.PlaceID, candidate.Geometry.Location.Lat, candidate.Geometry.Location.Lng, nil
	}
	return "", 0, 0, ErrNotFound
}

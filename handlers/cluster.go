package handlers

import (
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/maxheckel/parks/database"
	"github.com/maxheckel/parks/models"
	"github.com/maxheckel/parks/services/maps"
	"github.com/muesli/clusters"
	"github.com/muesli/kmeans"
	"math"
	"strconv"
)

func GetClusters(c *fiber.Ctx) error {
	cityName := c.Query("city", "Columbus")
	city := &models.City{}
	database.DB.Db.Model(&models.City{}).Preload("Parks").Where("name = ?", cityName).First(&city)
	if city.Name == "" {
		return c.Status(400).JSON(map[string]string{
			"type":    "BAD_REQUEST",
			"message": "could not get city",
			"error":   "no results for query",
		})
	}
	daysStr := c.Query("tourDays", "10")
	tourDays, err := strconv.Atoi(daysStr)
	if err != nil {
		return c.Status(400).JSON(map[string]string{
			"type":    "BAD_REQUEST",
			"message": "could not get tourDays count",
			"error":   err.Error(),
		})
	}
	startingLatStr := c.Query("start_lat", "40.02713809932883")
	startingLngStr := c.Query("start_lng", "-83.02027792396306")
	startingLat, err := strconv.ParseFloat(startingLatStr, 64)
	if err != nil {
		return c.Status(400).JSON(map[string]string{
			"type":    "BAD_REQUEST",
			"message": "bad starting lng or lat",
			"error":   err.Error(),
		})
	}
	startingLng, err := strconv.ParseFloat(startingLngStr, 64)
	if err != nil {
		return c.Status(400).JSON(map[string]string{
			"type":    "BAD_REQUEST",
			"message": "bad starting lng or lat",
			"error":   err.Error(),
		})
	}
	tourModel := &models.Tour{}
	database.DB.Db.
		Where("city_id = ?", city.ID).
		Where("start_lat = ?", startingLat).
		Where("start_lng = ?", startingLng).
		Where("days_count = ?", tourDays).
		Preload("Days.Parks.DayPark").
		First(&tourModel)
	if tourModel.ID != 0 {
		return c.Status(200).JSON(tourModel)
	}
	days, err := getTourDays(city, tourDays)
	if err != nil {
		return c.Status(500).JSON(map[string]string{
			"type":    "INTERNAL_ERROR",
			"message": "could not cluster tourDays",
			"error":   err.Error(),
		})
	}
	sortParks(days, startingLat, startingLng)

	mapsClient, err := maps.New()
	if err != nil {
		return c.Status(500).JSON(map[string]string{
			"type":    "INTERNAL_ERROR",
			"message": "could not initiate maps client",
			"error":   err.Error(),
		})
	}

	tour := &models.Tour{
		CityID:    city.ID,
		DaysCount: tourDays,
		StartLat:  startingLat,
		StartLng:  startingLng,
	}
	database.DB.Db.Save(&tour)

	ctx := context.Background()
	for _, day := range days {
		day.DirectionsURL, err = mapsClient.GetDrivingDirections(&ctx, &models.Location{startingLat, startingLng}, day.Parks)
		day.TourID = tour.ID
		if err != nil {
			return c.Status(500).JSON(map[string]string{
				"type":    "INTERNAL_ERROR",
				"message": "could not get driving directions",
				"error":   err.Error(),
			})
		}
	}
	database.DB.Db.Save(&days)
	for _, day := range days {
		for _, park := range day.Parks {
			database.DB.Db.Save(&models.DayPark{
				DayID:  day.ID,
				ParkID: park.ID,
				Order:  park.Sort,
			})
		}
	}

	return c.Status(200).JSON(days)
}

func sortParks(tourDaysArr []*models.Day, startingLat float64, startingLng float64) {
	for _, day := range tourDaysArr {
		firstPark := nextClosest(startingLat, startingLng, day.Parks)

		i := indexOf(firstPark, day.Parks)
		if i == -1 {
			panic("i is -1")
		}
		firstPark.Sort = 1
		day.Parks = append(day.Parks[:i], day.Parks[i+1:]...)
		parksSorted := []*models.Park{firstPark}
		lastPark := firstPark
		index := 2
		for len(day.Parks) > 0 {
			lastPark = nextClosest(lastPark.Latitude, lastPark.Longitude, day.Parks)
			lastPark.Sort = index
			index++
			parksSorted = append(parksSorted, lastPark)
			i := indexOf(lastPark, day.Parks)
			if i == -1 {
				panic("i is -1")
			}
			day.Parks = append(day.Parks[:i], day.Parks[i+1:]...)
		}
		day.Parks = parksSorted
	}
}

func indexOf(park *models.Park, parks []*models.Park) int {
	for i, p := range parks {
		if p.ID == park.ID {
			return i
		}
	}
	return -1
}

func nextClosest(startingLat, startingLng float64, available []*models.Park) *models.Park {
	minDistance := float64(10000000)
	var minPark *models.Park
	for _, park := range available {
		newDistance := distance(startingLat, startingLng, park.Latitude, park.Longitude)
		if newDistance < minDistance {
			minDistance = newDistance
			minPark = park
		}
	}
	return minPark
}

func distance(lat1 float64, lng1 float64, lat2 float64, lng2 float64, unit ...string) float64 {
	radlat1 := float64(math.Pi * lat1 / 180)
	radlat2 := float64(math.Pi * lat2 / 180)
	theta := float64(lng1 - lng2)
	radtheta := float64(math.Pi * theta / 180)

	dist := math.Sin(radlat1)*math.Sin(radlat2) + math.Cos(radlat1)*math.Cos(radlat2)*math.Cos(radtheta)
	if dist > 1 {
		dist = 1
	}

	dist = math.Acos(dist)
	dist = dist * 180 / math.Pi
	dist = dist * 60 * 1.1515

	if len(unit) > 0 {
		if unit[0] == "K" {
			dist = dist * 1.609344
		} else if unit[0] == "N" {
			dist = dist * 0.8684
		}
	}

	return dist
}

func getTourDays(city *models.City, days int) ([]*models.Day, error) {
	var d clusters.Observations
	for _, park := range city.Parks {
		d = append(d, clusters.Coordinates{
			park.Latitude,
			park.Longitude,
		})
	}

	// Partition the data points into 16 clusters
	km := kmeans.New()
	clusters, err := km.Partition(d, days)
	if err != nil {
		return nil, err
	}
	tourDaysArr := []*models.Day{}
	for i, c := range clusters {
		day := &models.Day{
			Name: fmt.Sprintf("Day %d", i+1),
		}
		for _, o := range c.Observations {
			day.Parks = append(day.Parks, parkForLatLng(city.Parks, o.Coordinates()))
		}
		tourDaysArr = append(tourDaysArr, day)
	}
	return tourDaysArr, nil
}

func parkForLatLng(parks []models.Park, coords clusters.Observation) *models.Park {
	for _, park := range parks {
		if park.Latitude != coords.Coordinates()[0] {
			continue
		}
		if park.Longitude != coords.Coordinates()[1] {
			continue
		}
		return &park
	}
	return nil
}

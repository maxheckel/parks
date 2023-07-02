package models

import (
	"fmt"
	"gorm.io/gorm"
)

type Location struct {
	Latitude  float64
	Longitude float64
}

func (l Location) ToString() string {
	return fmt.Sprintf("%0f %0f", l.Latitude, l.Longitude)
}

type Tour struct {
	gorm.Model
	CityID uint
	Days   int
}

type Park struct {
	gorm.Model
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Name      string
	CityID    uint
	PlaceID   string
	Sort      int `sql:"-" gorm:"-"`
	City      *City
}

func (p Park) LatLngString() string {
	return fmt.Sprintf("%0f %0f", p.Latitude, p.Longitude)
}

type City struct {
	gorm.Model
	Name    string `json:"name" gorm:"index:idx_city,unique"`
	Country string `json:"country" gorm:"index:idx_city,unique"`
	State   string `json:"state" gorm:"index:idx_city,unique"`
	Parks   []Park
}

func (c City) ToLocationName() string {
	return fmt.Sprintf("%s, %s %s", c.Name, c.State, c.Country)
}

type Day struct {
	gorm.Model
	Name          string
	DirectionsURL string
	Parks         []*Park `gorm:"many2many:day_parks;"`
}

type DayPark struct {
	DayID  int `gorm:"primaryKey"`
	ParkID int `gorm:"primaryKey"`
	Order  int
}

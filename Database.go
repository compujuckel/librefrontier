package main

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type Device struct {
	gorm.Model
	Mac              string    `gorm:"UNIQUE; NOT NULL"`
	FavoriteStations []Station `gorm:"many2many:favorite"`
}

type Station struct {
	gorm.Model
	StationId uint   `gorm:"NOT NULL"`
	Name      string `gorm:"NOT NULL"`
}

type Database struct {
	dbHandle *gorm.DB
}

func NewDatabase(config *Config) (*Database, error) {
	database := Database{}

	db, err := gorm.Open("postgres", config.dbConnString)
	if err != nil {
		return nil, errors.Wrap(err, "Cannot connect to dbHandle")
	}

	db.LogMode(true)
	database.dbHandle = db

	db.AutoMigrate(Device{}, Station{})

	return &database, nil
}

func (d *Database) AddFavorite(mac string, stationId uint, stationName string) {
	var device Device
	d.dbHandle.FirstOrCreate(&device, Device{Mac: mac})
	log.Info("Device", device)

	d.dbHandle.Model(&device).Association("FavoriteStations").Append(Station{
		StationId: stationId,
		Name:      stationName,
	})
}

func (d *Database) GetFavoriteStations(mac string) []Station {
	var device Device
	d.dbHandle.First(&device, Device{Mac: mac})

	var favorites []Station
	d.dbHandle.Model(&device).Association("FavoriteStations").Find(&favorites)

	return favorites
}

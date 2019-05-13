package common

import (
	"database/sql"
	"github.com/compujuckel/librefrontier/common/radioprovider"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type Database struct {
	db *sql.DB
}

func NewDatabase(config *Config) (*Database, error) {
	database := Database{}

	db, err := sql.Open("postgres", config.dbConnString)
	if err != nil {
		return nil, errors.Wrap(err, "Cannot connect to db")
	}

	database.db = db

	return &database, nil
}

func (d *Database) CreateDevice(mac string) {
	s := "INSERT INTO device (mac) VALUES ($1) ON CONFLICT DO NOTHING;"

	_, err := d.db.Exec(s, mac)
	if err != nil {
		log.Error("Error creating device: ", err)
	}
}

func (d *Database) createRadioBrowserStation(stationId int64, stationName string) {
	s := "INSERT INTO station (radiobrowser_id, name) VALUES ($1, $2) ON CONFLICT DO NOTHING;"

	_, err := d.db.Exec(s, stationId, stationName)
	if err != nil {
		log.Error("Error creating station: ", err)
	}
}

func (d *Database) AddFavorite(mac string, stationId int64, stationName string) {
	d.createRadioBrowserStation(stationId, stationName)

	s := `INSERT INTO favorite (device_id, station_id) SELECT (SELECT d.device_id FROM device d WHERE d.mac = $1), (SELECT s.station_id FROM station s WHERE s.radiobrowser_id = $2)`

	_, err := d.db.Exec(s, mac, stationId)
	if err != nil {
		log.Error("Error creating station: ", err)
	}
}

func (d *Database) RemoveFavorite(mac string, stationId uint64) {
	s := `DELETE FROM favorite f
                USING device d, station s
           	    WHERE d.device_id = f.device_id
                  AND s.station_id = f.station_id
           	      AND s.radiobrowser_id = $1`

	_, err := d.db.Exec(s, stationId)
	if err != nil {
		log.Error("Error removing favorite: ", err)
	}
}

func (d *Database) IsFavorite(mac string, stationId uint64) bool {
	s := "SELECT EXISTS(SELECT * FROM favorite f JOIN device d ON d.device_id = f.device_id JOIN station s on s.station_id = f.station_id WHERE d.mac = $1 AND s.radiobrowser_id = $2)"

	row := d.db.QueryRow(s, mac, stationId)

	var exists bool
	err := row.Scan(&exists)
	if err != nil {
		log.Error("cannot parse row", err)
		return false
	}

	return exists
}

func (d *Database) GetFavoriteStations(mac string) []radioprovider.Station {
	s := `SELECT s.radiobrowser_id,
                 s.name
            FROM favorite f
            JOIN device d ON d.device_id = f.device_id
            JOIN station s ON s.station_id = f.station_id
           WHERE d.mac = $1;`

	rows, err := d.db.Query(s, mac)
	if err != nil {
		log.Error("error getting favorite stations", err)
		return []radioprovider.Station{}
	}

	var stations []radioprovider.Station
	for rows.Next() {
		var s radioprovider.Station

		err := rows.Scan(&s.Id, &s.Name)
		if err != nil {
			log.Error("error scanning row", err)
			return []radioprovider.Station{}
		}

		stations = append(stations, s)
	}

	return stations
}

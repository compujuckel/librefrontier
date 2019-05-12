CREATE TABLE device (
    device_id SERIAL PRIMARY KEY,
    mac TEXT UNIQUE NOT NULL
);

CREATE TABLE station (
    station_id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    radiobrowser_id INTEGER UNIQUE,
    stream_url TEXT
);

CREATE TABLE favorite (
    favorite_id SERIAL PRIMARY KEY,
    device_id INTEGER NOT NULL REFERENCES device (device_id),
    station_id INTEGER NOT NULL REFERENCES station (station_id),
    UNIQUE(device_id, station_id)
)
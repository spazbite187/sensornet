package sensornet

// SensorData ...
type SensorData struct {
	Pk         int     `json:"key" storm:"id,increment"` // Primary key
	ID         string  `json:"id" storm:"index"`
	Location   string  `json:"location" storm:"index"`
	Uptime     string  `json:"uptime" storm:"index"`
	IP         string  `json:"ip" storm:"index"`
	SSID       string  `json:"ssid" storm:"index"`
	Signal     int     `json:"signal" storm:"index"`
	TempC      float64 `json:"temp" storm:"index"`
	TempF      float64 `json:"tempf" storm:"index"`
	LastUpdate string  `json:"lastupdate" storm:"index"`
	TimeSince  string  `json:"timesince,omitempty" storm:"index"`
}

// Sensor ...
type Sensor struct {
	ID          string `json:"id,omitempty" storm:"index"`
	Location    string `json:"location,omitempty" storm:"index"`
	NumReadings int    `json:"numreads,omitempty"`
}

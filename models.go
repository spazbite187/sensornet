package sensornet

// SensorData ...
type SensorData struct {
	Pk         int     `json:"key" storm:"id,increment"` // Primary key
	ID         string  `json:"id" storm:"index"`
	Location   string  `json:"location" storm:"index"` // location at time of reading
	Uptime     string  `json:"uptime" storm:"index"`
	IP         string  `json:"ip" storm:"index"`
	SSID       string  `json:"ssid" storm:"index"`
	Signal     int     `json:"signal" storm:"index"`
	TempC      float64 `json:"temp" storm:"index"`
	TempF      float64 `json:"tempf" storm:"index"`
	LastUpdate string  `json:"lastupdate" storm:"index"`
	TimeSince  string  `json:"timesince,omitempty" storm:"index"`
}

// Sensor contains information about the sensor and is designed to be updated asynchronously. This will
// enable faster returns to the requestor.
type Sensor struct {
	ID          string     `json:"id" storm:"index"`
	Location    string     `json:"location"`
	Uptime      string     `json:"uptime,omitempty"`
	IP          string     `json:"ip,omitempty"`
	SSID        string     `json:"ssid,omitempty"`
	Signal      int        `json:"signal,omitempty"`
	NumReadings int        `json:"numreads,omitempty"`
	TempC       float64    `json:"temp,omitempty"`
	TempF       float64    `json:"tempf,omitempty"`
	HighTemp    SensorData `json:"high_temp_data,omitempty"`
	LowTemp     SensorData `json:"low_temp_data,omitempty"`
	AvgTemp     float64    `json:"avg_temp,omitempty"`
	TempGraph   []byte     `json:"temp_graph,omitempty"`
	HighSignal  SensorData `json:"high_signal_data,omitempty"`
	LowSignal   SensorData `json:"low_signal_data,omitempty"`
	AvgSignal   float64    `json:"avg_signal,omitempty"`
	SignalGraph []byte     `json:"signal_graph,omitempty"`
	LastUpdate  string     `json:"lastupdate,omitempty"`
}

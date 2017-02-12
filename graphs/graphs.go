package graphs

import (
	"bufio"
	"bytes"
	"image/color"
	"time"

	"github.com/gonum/plot"
	"github.com/gonum/plot/plotter"
	"github.com/gonum/plot/vg"
	"github.com/gonum/plot/vg/draw"
	"github.com/gonum/plot/vg/vgsvg"
	"github.com/spazbite187/sensornet"
)

const (
	xSize = 7
	ySize = 3
)

// GetTempGraph ...
func GetTempGraph(input []*sensornet.SensorData) ([]byte, error) {
	// xticks defines how we convert and display time.Time values.
	xticks := plot.TimeTicks{Format: "15:04\n2006-01-02"}
	currentLocation := time.Now().Location()
	xticks.Time = plot.UnixTimeIn(currentLocation)
	p, err := plot.New()
	if err != nil {
		return []byte{}, nil
	}

	p.X.Label.Text = "Time"
	p.Y.Label.Text = "Temp (F)"

	p.X.Tick.Marker = xticks

	data, err := tempPoints(input)
	if err != nil {
		return []byte{}, nil
	}

	line, points, err := plotter.NewLinePoints(data)
	if err != nil {
		return []byte{}, nil
	}
	points.Shape = draw.CircleGlyph{}
	points.GlyphStyle.Radius = 0
	points.Color = color.RGBA{R: 255, G: 255, B: 255, A: 255}

	line.Color = color.RGBA{R: 255, G: 255, B: 255, A: 255}

	setColors(p)

	p.Add(plotter.NewGrid())
	p.Add(line, points)

	// Save to []byte section
	c := vgsvg.New(xSize*vg.Inch, ySize*vg.Inch) // Create a Canvas for writing SVG images.
	p.Draw(draw.New(c))                          // Draw to the Canvas.
	image, err := getSVGBytes(c)
	if err != nil {
		return []byte{}, nil
	}

	return image, nil
}

// GetSignalGraph ...
func GetSignalGraph(input []*sensornet.SensorData) ([]byte, error) {
	// xticks defines how we convert and display time.Time values.
	xticks := plot.TimeTicks{Format: "15:04\n2006-01-02"}
	currentLocation := time.Now().Location()
	xticks.Time = plot.UnixTimeIn(currentLocation)
	p, err := plot.New()
	if err != nil {
		return []byte{}, nil
	}

	p.X.Label.Text = "Time"
	p.Y.Label.Text = "Signal (dBm)"
	p.X.Tick.Marker = xticks

	data, err := signalPoints(input)
	if err != nil {
		return []byte{}, nil
	}

	line, points, err := plotter.NewLinePoints(data)
	if err != nil {
		return []byte{}, nil
	}
	points.Shape = draw.CircleGlyph{}
	points.GlyphStyle.Radius = 0
	points.Color = color.RGBA{R: 255, G: 255, B: 255, A: 255}

	line.Color = color.RGBA{R: 255, G: 255, B: 255, A: 255}

	setColors(p)

	p.Add(plotter.NewGrid())
	p.Add(line, points)

	// Save to []byte section
	c := vgsvg.New(xSize*vg.Inch, ySize*vg.Inch) // Create a Canvas for writing SVG images.
	p.Draw(draw.New(c))                          // Draw to the Canvas.
	image, err := getSVGBytes(c)
	if err != nil {
		return []byte{}, nil
	}

	return image, nil
}

// tempPoints returns x, y points based on SensorData.
func tempPoints(data []*sensornet.SensorData) (plotter.XYs, error) {
	pts := make(plotter.XYs, len(data))
	for k, v := range data {
		lastUpdate, err := time.Parse(time.ANSIC, v.LastUpdate)
		if err != nil {
			return pts, err
		}

		timeFloat := float64(lastUpdate.Unix())
		pts[k].X = timeFloat
		pts[k].Y = v.TempF
	}

	return pts, nil
}

// signalPoints returns x, y points based on SensorData.
func signalPoints(data []*sensornet.SensorData) (plotter.XYs, error) {
	pts := make(plotter.XYs, len(data))
	for k, v := range data {
		lastUpdate, err := time.Parse(time.ANSIC, v.LastUpdate)
		if err != nil {
			return pts, err
		}

		timeFloat := float64(lastUpdate.Unix())
		pts[k].X = timeFloat
		pts[k].Y = float64(v.Signal)
	}

	return pts, nil
}

func getSVGBytes(canvas *vgsvg.Canvas) ([]byte, error) {
	var b bytes.Buffer
	buffer := bufio.NewWriter(&b)
	if _, err := canvas.WriteTo(buffer); err != nil {
		return []byte{}, err
	}

	return b.Bytes(), nil
}

func setColors(p *plot.Plot) {
	p.X.Label.Color = color.RGBA{R: 255, G: 255, B: 255, A: 255}
	p.Y.Label.Color = color.RGBA{R: 255, G: 255, B: 255, A: 255}
	p.X.Tick.Color = color.RGBA{R: 255, G: 255, B: 255, A: 255}
	p.Y.Tick.Color = color.RGBA{R: 255, G: 255, B: 255, A: 255}
	p.X.Tick.LineStyle.Color = color.RGBA{R: 255, G: 255, B: 255, A: 255}
	p.Y.Tick.LineStyle.Color = color.RGBA{R: 255, G: 255, B: 255, A: 255}
	p.X.Tick.Label.Color = color.RGBA{R: 200, G: 200, B: 200, A: 255}
	p.Y.Tick.Label.Color = color.RGBA{R: 200, G: 200, B: 200, A: 255}
	p.BackgroundColor = color.RGBA{R: 55, G: 58, B: 60, A: 255}
	p.Title.Color = color.RGBA{R: 255, G: 255, B: 255, A: 255}
	p.Title.TextStyle.Color = color.RGBA{R: 255, G: 255, B: 255, A: 255}
	p.X.LineStyle.Color = color.RGBA{R: 255, G: 255, B: 255, A: 255}
	p.Y.LineStyle.Color = color.RGBA{R: 255, G: 255, B: 255, A: 255}
	p.Legend.Color = color.RGBA{R: 255, G: 255, B: 255, A: 255}
}

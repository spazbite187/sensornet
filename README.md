# SensorNet 1.1.0
SensorNet is a web service used receive sensor data readings in JSON format. This data is stored in a file DB
and used to produce SVG graphs, and html front end. JSON API end points are also available.

# Get
```
$ go get -u github.com/spazbite187/sensornet
```

# Installation
```
$ go install github.com/spazbite187/sensornet/cmd/sensornet
```

# Run
```
$ cd src/github.com/spazbite187/sensornet/
$ PORT=8080 MODE=debug sensornet
```

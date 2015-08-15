package main

import (
	"os"
	"os/signal"

	//"github.com/paulcull/go-webbrick" // The magic part that lets us control devices

	//"github.com/davecgh/go-spew/spew" // For neatly outputting stuff
	//"strconv" // For String construction
	//"time" // Used as part of "setInterval" and for pausing code to allow for data to come back
	//"fmt" //Output stuff to the screen

	"github.com/alecthomas/kingpin"
	"github.com/juju/loggo"
	//"github.com/rcrowley/go-metrics"
)

//const statsTopic = "$device/stats"
const wbTopic = "$device/wb"

var (
	debug    = kingpin.Flag("debug", "Enable debug mode.").OverrideDefaultFromEnvar("DEBUG").Bool()
	daemon   = kingpin.Flag("daemon", "Run in daemon mode.").Short('d').Bool()
	mqttURL  = kingpin.Flag("mqttUrl", "The MQTT url to publish too.").Short('u').Default("tcp://localhost:1883").String()
	logName  = kingpin.Flag("logName", "The Log Name.").Short('u').Default("webbrick-mqtt").String()
	//port     = kingpin.Flag("port", "HTTP Port.").Short('i').OverrideDefaultFromEnvar("PORT").Default("9980").Int()
	//path     = kingpin.Flag("path", "Path to static content.").Short('p').OverrideDefaultFromEnvar("CONTENT_PATH").Default("./public").String()
	interval = kingpin.Flag("interval", "Publish interval.").Short('i').Default("30").Int()

	//log = loggo.GetLogger("mqtt_webbrick")
)

func main() {
	kingpin.Version(Version)
	kingpin.Parse()

	setupLoggo(*debug)

	//localRegistry := metrics.NewRegistry()
	//publisher := newPublisher(localRegistry)
	//publisher := newLitePublisher()

	

	mqttEngine, mqerr := newMqttEngine(*mqttURL)
	wbEngine, wberr := NewWebBrickDriver(mqttEngine)

	if mqerr != nil && wberr != nil {
		panic(err)
	}

	wbEngine.Start()

	//ws := newWsServer(localRegistry)

	//go ws.listenAndServ(*port, *path)

	// Set up channel on which to send signal notifications.
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, os.Interrupt, os.Kill)

	// Wait for receiving a signal.
	<-sigc

	// Disconnect the Network Connection.
	if err := engine.disconnect(); err != nil {
		panic(err)
	}
}

func setupLoggo(debug bool) {
	// apply flags
	if debug {
		loggo.GetLogger(*logName).SetLogLevel(loggo.DEBUG)
	} else {
		loggo.GetLogger(*logName).SetLogLevel(loggo.INFO)
	}
}

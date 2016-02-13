package main

import (
	"encoding/json"
	"flag"
	"github.com/labstack/echo"
	mw "github.com/labstack/echo/middleware"
	"golang.org/x/net/websocket"
	"strconv"
)

type controlMessage struct {
	Message string
	Command string
}

var (
	port         int
	imagesPath   string
	soundsPath   string
	publicPath   string
	daemonize    bool
	enableLogger bool
)

var (
	controlChan = make(chan controlMessage, 10)
	widgetChans = make([]chan controlMessage, 0)
)

func main() {
	flag.IntVar(&port, "p", 8998, "server port")
	flag.StringVar(&imagesPath, "images", "images", "images folder path")
	flag.StringVar(&soundsPath, "sounds", "sounds", "sounds folder path")
	flag.StringVar(&soundsPath, "public", "public", "public folder path")
	flag.BoolVar(&daemonize, "d", false, "run as daemon")
	flag.BoolVar(&enableLogger, "l", false, "enable logging")

	flag.Parse()

	runServer()
}

func runServer() {
	go fanOutToWidgets()

	e := echo.New()

	if enableLogger {
		e.Use(mw.Logger())
	}
	e.Use(mw.Recover())

	e.Static("/", publicPath)
	e.Static("/sounds", soundsPath)
	e.Static("/images", imagesPath)

	e.WebSocket("/ws/widget", handleWidgetWS)
	e.WebSocket("/ws/command", handleCommandWS)

	e.Run(":" + strconv.Itoa(port))
}

func handleWidgetWS(c *echo.Context) error {
	ws := c.Socket()
	defer ws.Close()

	in := make(chan controlMessage, 100)
	registerClient(in)
	defer deregisterClient(in)

	for {
		message := <-in
		if err := websocket.Message.Send(ws, message.ToJSON()); err != nil {
			return err
		}
	}
}

func handleCommandWS(c *echo.Context) error {
	ws := c.Socket()
	defer ws.Close()

	websocket.JSON.Send(ws, GetFileList())

	var message controlMessage
	for {
		if err := websocket.JSON.Receive(ws, &message); err != nil {
			return err
		}
		controlChan <- message
	}
}

func fanOutToWidgets() {
	for {
		select {
		case message := <-controlChan:
			for _, widget := range widgetChans {
				widget <- message
			}
		}
	}
}

func (msg *controlMessage) ToJSON() string {
	result, _ := json.Marshal(msg)
	return string(result)
}

func registerClient(ch chan controlMessage) {
	widgetChans = append(widgetChans, ch)
}

func deregisterClient(ch chan controlMessage) {
	for index, widget := range widgetChans {
		if ch == widget {
			widgetChans = append(widgetChans[:index], widgetChans[index+1:]...)
			break
		}
	}
}

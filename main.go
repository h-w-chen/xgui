package main

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"time"

	"github.com/BurntSushi/xgb/xproto"
	"github.com/BurntSushi/xgbutil/keybind"
	"github.com/BurntSushi/xgbutil/mousebind"
	"github.com/BurntSushi/xgbutil/xevent"
	"github.com/BurntSushi/xgbutil/xgraphics"

	"github.com/BurntSushi/xgbutil"
)

var (
	width  = 400
	height = 600
	co     = color.RGBA{0, 200, 0, 0}

	ticker     = time.NewTicker(50 * time.Millisecond)
	value      = 0
	rightmostX = 0
)

// type of function to get next Y value to draw
type getYToDraw func() int

func getYOnTicker() int {
	<-ticker.C
	y := value
	value++
	return y
}

func drawY(y int, canvas *xgraphics.Image, winID xproto.Window) {
	if y > height {
		y = height
	}

	canvas.Set(rightmostX, y, co)
	canvas.XDraw()
	canvas.XPaint(winID)

	rightmostX++
}

func main() {
	X, err := xgbutil.NewConn()
	if err != nil {
		log.Fatalf("[x11] init error: %s\n", err)
	}

	keybind.Initialize(X)
	mousebind.Initialize(X)

	canvas := xgraphics.New(X, image.Rect(0, 0, width, height))
	win := canvas.XShowExtra("X11-GUI", true)
	win.Listen(xproto.EventMaskButtonPress | xproto.EventMaskButtonRelease | xproto.EventMaskKeyPress)

	go func(op getYToDraw) {
		for {
			drawY(op(), canvas, win.Id)
		}
	}(getYOnTicker)

	xevent.Main(X)

	fmt.Println("exited")
}

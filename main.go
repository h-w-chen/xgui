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
	width  = 500
	height = 500
	co     = color.RGBA{0, 200, 0, 0}

	ticker = time.NewTicker(50 * time.Millisecond)
	x, y   = 0, 0
)

// type of function to get next point to draw
type getPointToDraw func() image.Point

func nextPointOnTicker() image.Point {
	<-ticker.C
	xx, yy := x, y
	x++
	y++
	return image.Pt(xx, yy)
}

func drawPoint(pt image.Point, canvas *xgraphics.Image, winID xproto.Window) {
	canvas.Set(pt.X, pt.Y, co)
	canvas.XDraw()
	canvas.XPaint(winID)
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

	go func(op getPointToDraw) {
		for {
			drawPoint(op(), canvas, win.Id)
		}
	}(nextPointOnTicker)

	xevent.Main(X)

	fmt.Println("exited")
}

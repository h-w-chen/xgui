package main

import (
	"fmt"
	"image"
	"image/color"
	"log"

	"github.com/BurntSushi/xgb/xproto"
	"github.com/BurntSushi/xgbutil/keybind"
	"github.com/BurntSushi/xgbutil/mousebind"
	"github.com/BurntSushi/xgbutil/xevent"
	"github.com/BurntSushi/xgbutil/xgraphics"

	"github.com/BurntSushi/xgbutil"

	"./rs232"
)

var (
	width  = 600
	height = 400
	co     = color.RGBA{0, 200, 0, 0}
	bg     = xgraphics.BGRA{R: 0x0, G: 0x0, B: 0x0, A: 0xff}

	x      = 0
	boundX = width - 1
	boundY = height - 1
)

// type of function to get next Y value to draw
type getYToDraw func() int

func drawY(y int, canvas *xgraphics.Image, winID xproto.Window) {
	if y > boundY {
		y = boundY
	}

	// if the current x is beyond the (right) drawing boundary (should be by 1 only),
	// left shift the canvas by 1, in order to make room for the current point
	if x > boundX { // MUST x == boundX + 1
		if x-boundX != 1 {
			log.Fatalf("internal error: x=%d, boundX=%d, distance is not 1\n", x, boundX)
		}
		leftShiftCanvasByX(canvas, 1)
		x = boundX
	}

	canvas.Set(x, y, co)
	canvas.XDraw()
	canvas.XPaint(winID)

	x++
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

	//	go func() {
	//		web.StartServer(8777)
	//	}()

	go rs232.StartSioReceiver("/dev/ttyUSB0", 115200)

	go func(op getYToDraw) {
		for {
			drawY(op(), canvas, win.Id)
		}
	}(rs232.GetY)

	xevent.Main(X)

	fmt.Println("exited")
}

func leftShiftCanvasByX(canvas *xgraphics.Image, distance int) {
	canvas.For(func(x, y int) xgraphics.BGRA {
		//if x >= canvas.Rect.Dx()-distance {
		//	return bg
		//}
		return canvas.At(x+distance, y).(xgraphics.BGRA)
	})
}

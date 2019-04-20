package rs232

import (
	"fmt"
	"log"
	"strconv"
	"syscall"

	"github.com/schleibinger/sio"
)

var chY = make(chan int)

// GetY produces Y value based on feed received at RS232 port
func GetY() int {
	return <-chY
}

// StartSioReceiver starts receiving at rs232 port
func StartSioReceiver(dev string, baud int) {
	// initialize rs232 setting
	port, err := sio.Open(dev, bauds[baud])
	if err != nil {
		log.Fatalf("open dev: %s\n", err)
	}

	for {
		// todo: robust read - each feed has 3 chars
		const length = 3
		fullFeed := []byte{}

		for {
			rxbuf := []byte{0}
			_, err := port.Read(rxbuf)
			if err != nil {
				log.Fatalf("read: %s", err)
			}

			fullFeed = append(fullFeed, rxbuf...)
			if len(fullFeed) >= length {
				break
			}
		}

		// decode whatever received (string to int)
		fmt.Println(string(fullFeed))
		y, err := strconv.Atoi(string(fullFeed))
		if err != nil {
			continue // OK to ignore invalid input
		}

		chY <- y
	}
}

var bauds = map[int]uint32{
	50:      syscall.B50,
	75:      syscall.B75,
	110:     syscall.B110,
	134:     syscall.B134,
	150:     syscall.B150,
	200:     syscall.B200,
	300:     syscall.B300,
	600:     syscall.B600,
	1200:    syscall.B1200,
	1800:    syscall.B1800,
	2400:    syscall.B2400,
	4800:    syscall.B4800,
	9600:    syscall.B9600,
	19200:   syscall.B19200,
	38400:   syscall.B38400,
	57600:   syscall.B57600,
	115200:  syscall.B115200,
	230400:  syscall.B230400,
	460800:  syscall.B460800,
	500000:  syscall.B500000,
	576000:  syscall.B576000,
	921600:  syscall.B921600,
	1000000: syscall.B1000000,
	1152000: syscall.B1152000,
	1500000: syscall.B1500000,
	2000000: syscall.B2000000,
	2500000: syscall.B2500000,
	3000000: syscall.B3000000,
	3500000: syscall.B3500000,
	4000000: syscall.B4000000,
}

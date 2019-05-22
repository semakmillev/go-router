package main

import (
	"time"
	"fmt"
)

func main() {
	a := make(chan int, 1)
	finish := make(chan bool, 1)
	go func(p chan int) {
		time.Sleep(5 * time.Second)
		fmt.Println("!!!")
		select {
		case p <- 10:
			fmt.Println("sent")
		default:
			fmt.Println("not sent")

			//fmt.Println("wait for recieve...")
			//time.Sleep(1 * time.Second)

		}
		select {
		case p <- 20:
			fmt.Println("sent 20")
		default:
			fmt.Println("not sent 20")

			//fmt.Println("wait for recieve...")
			//time.Sleep(1 * time.Second)

		}
		fmt.Println("???")

	}(a)
	go func(p chan int) {
		x := 0
		waiting := true
		for waiting {
			select {
			case x = <-p:
				fmt.Println("Got it!!!", x)
				waiting = false
				time.Sleep(1 * time.Second)
			default:
				time.Sleep(2 * time.Second)
				fmt.Println("waiting...")
			}
		}

		finish <- true
	}(a)
	<-finish

}

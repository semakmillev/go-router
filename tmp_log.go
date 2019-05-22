package main

import (
	"github.com/chasex/glog"
	"math/rand"
	"time"
)

func testChannel(param string) {

}

const logPath = "/some/location/example.log"

func main() {

	options := glog.LogOptions{
		File:  "./abc.log",
		Flag:  glog.LstdFlags | glog.Lmicroseconds | glog.Lshortfile,
		Level: glog.Ldebug,
		Mode:  glog.R_Day,
		//Maxsize: 1024 * 1024 * 16,
	}
	logger, err := glog.New(options)
	if err != nil {
		panic(err)
	}
	/*
	con := 100

	done := make(chan int, con)
	for k := 0; k < con; k++ {
		go func() {
			for i := 0; i < 100; i++ {
				currentTime := time.Now()
				s := currentTime.Format("2006.01.02 15:04:05")
				//logger.Info(s)
				logger.Error(s)
				a := rand.Intn(5)
				time.Sleep(time.Duration(a) * time.Second)
			}
			done <- 1
		}()
	}

	for i := 0; i < con; i++ {
		<-done
	}
	*/
	for i := 0; i < 100; i++ {
		currentTime := time.Now()
		s := currentTime.Format("2006.01.02 15:04:05")
		//logger.Info(s)
		logger.Error(s)
		a := rand.Intn(5)
		time.Sleep(time.Duration(a) * time.Second)
	}
	defer logger.Flush()

	/*
	b := []byte("<test/>")
	request := bytes.NewReader(b)

	resp, err := http.Post("http://0.0.0.0:5002/err", "", request)
	if err != nil {
		fmt.Println("err", err, resp.StatusCode)
	}

	respBody, readErr := ioutil.ReadAll(resp.Body)
	if readErr != nil {
		fmt.Println("readErr", readErr)
	}
	for k, v := range resp.Header {
		fmt.Println(k, v, reflect.TypeOf(v), strings.Join(v, ";"))
	}
	fmt.Println(string(respBody[:len(respBody)]), resp.StatusCode)*/
}

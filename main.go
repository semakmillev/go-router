package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
	"./logging"
	"./routerStruct"
	"github.com/gorilla/mux"
	"bytes"
	"net"
	"net/url"
	"github.com/chasex/glog"
	"strconv"
)

func getNotFoundResponse() *http.Response {
	return &http.Response{
		Status:        "404 Not Found",
		StatusCode:    404,
		Proto:         "HTTP/1.1",
		ProtoMajor:    1,
		ProtoMinor:    1,
		Body:          ioutil.NopCloser(bytes.NewReader([]byte(""))),
		ContentLength: int64(0),
		Header:        make(http.Header, 0),
	}

}

func sleep(n int) {
	time.Sleep(time.Duration(n) * time.Second)
	fmt.Println("123")
}

func ping(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, r.URL.Path)
	fmt.Fprintf(w, "Pong")
}

func call(url_ string, requestMethod string, requestBody []byte, requestHeader http.Header) (http.Response, error) {
	s := string(requestBody[:len(requestBody)])
	fmt.Println("send", url_, s)
	proxyReq, err := http.NewRequest(requestMethod, url_, bytes.NewReader(requestBody))
	proxyReq.Header = make(http.Header)
	for k, v := range requestHeader {
		proxyReq.Header.Set(k, strings.Join(v, ""))
	}
	httpClient := http.Client{}
	resp, err := httpClient.Do(proxyReq) // http.Post(url, , request)
	if resp == nil{
		resp = getNotFoundResponse()
	}
	return *resp, err
}

func redirect(w http.ResponseWriter, r *http.Request, route *routerStruct.Route, rCh chan routerStruct.Route, logger *glog.Logger) {
	//fmt.Println("got it")
	requestBody, errRead := ioutil.ReadAll(r.Body)
	last := false
	for !last {
		select {
		case *route = <-rCh:
		default:
			last = true
		}
	}
	currentRoute := *route
	if errRead != nil {
		fmt.Println(errRead)
	}
	var gotAnswer bool = false
	for !gotAnswer {
		gotAnswer = true
	}
	var resp http.Response
	var found = false
	var err error
	allDead := routerStruct.CheckAllDead(currentRoute)
	//fmt.Println(allDead)
	for _, direction := range currentRoute.Directions {
		if direction.Alive == 1 || allDead {
			found = true
			url_ := direction.URL + r.URL.Path
			resp, err = call(url_, r.Method, requestBody, r.Header)
			if err != nil{
				fmt.Println(err)
			}
			//fmt.Println("received answer", url_, resp.StatusCode)
			logging.Info(logger, "received answer: "+url_+" "+strconv.Itoa(resp.StatusCode))
			//logging.Info(logger, "rec")
			if resp.StatusCode == 200 {
				break
			} else {
				logging.Error(logger, "received error: "+url_+" "+strconv.Itoa(resp.StatusCode))
			}
		}

	}
	if !found {
		resp = *getNotFoundResponse()
		logging.Error(logger,"Router "+route.Name+" is unreachable for url " + r.URL.Path)
	}
	for k, v := range resp.Header {
		w.Header().Set(k, strings.Join(v, ";"))
	}
	respBody, errRead := ioutil.ReadAll(resp.Body)
	if errRead != nil {
		fmt.Println(errRead)
	}
	//fmt.Println("recieved body", string(respBody[:len(respBody)]))
	w.WriteHeader(resp.StatusCode)
	w.Write(respBody)

}

func runRoute(router routerStruct.Route, rCh chan routerStruct.Route, logger *glog.Logger) {
	// fmt.Println(router.Name + ":" + router.Port)
	logger.Info("run Route - ", router.Name+":"+router.Port)
	r := mux.NewRouter()
	r.HandleFunc("/{.*}", func(w http.ResponseWriter, r *http.Request) {
		currentRouter := router
		redirect(w, r, &currentRouter, rCh, logger)
	}).Methods("POST")
	http.ListenAndServe("0.0.0.0"+":"+router.Port, r)
}

func clearChannel(rCh chan routerStruct.Route) {
	select {
	case <-rCh:
	default:
	}
}

func checkAvalable(router *routerStruct.Route, rCh chan routerStruct.Route, logger *glog.Logger)  {
	// fmt.Println("checking ...", router.Name)
	logging.Info(logger, "checking ..."+router.Name)
	for true {
		timeout := time.Duration(3 * time.Second)
		for i, direction := range router.Directions {
			checkingUrl, _ := url.Parse(direction.URL)
			_, err := net.DialTimeout("tcp", checkingUrl.Host, timeout)
			if err == nil {
				if router.Directions[i].Alive == 0 {
					logging.Warn(logger, direction.URL+" became alive")
					clearChannel(rCh)
					select {
					case rCh <- *router:
					default:
					}

					//logger.Warn(direction.URL, " became alive")

				}
				router.Directions[i].Alive = 1

			} else {
				if router.Directions[i].Alive == 1 {
					clearChannel(rCh)
					select {
					case rCh <- *router:
					default:
					}
					logging.Warn(logger, direction.URL+" became dead")
					//logger.Warn(direction.URL, " became dead")
				}
				router.Directions[i].Alive = 0
				//fmt.Println(direction.URL, " is dead")
			}
		}
		time.Sleep(5 * time.Second)
	}
}

func main() {
	operationDone := make(chan bool)

	jsonFile, err := os.Open("local_config.json")
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()
	val, _ := ioutil.ReadAll(jsonFile)
	var res routerStruct.JsonRouter
	json.Unmarshal([]byte(val), &res)
	fmt.Println("Succesfully opened local_config.json")
	routerChannels := make(map[string]chan routerStruct.Route, 1)
	for _, element := range res.Routes {
		routerChannels[element.Name] = make(chan routerStruct.Route)
		routerStruct.ResetDirections(&element)
		//fmt.Println("ruote-" + element.Name + ".log")
		logger := logging.GetLogger("ruote-" + element.Name + ".log")
		defer logger.Flush()
		//routerChannels[element.Name] <- element

		go func(r routerStruct.Route, rCh chan routerStruct.Route) {
			runRoute(r, rCh, logger)
		}(element, routerChannels[element.Name])

		go func(r routerStruct.Route, rCh chan routerStruct.Route) {
			go checkAvalable(&r, routerChannels[element.Name], logger)
		}(element, routerChannels[element.Name])

	}

	<-operationDone
	fmt.Println("Stop")
	//r := mux.NewRouter()
	//r.HandleFunc("/sleep", sl).Methods("GET")
	//r.HandleFunc("/ping", ping).Methods("GET")
	//r.HandleFunc("/jobresult/{jobfolder}/{jobname}", routeHandler).Methods("GET")

}

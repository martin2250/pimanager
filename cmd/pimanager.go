package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/martin2250/pimanager/dev"
	"periph.io/x/periph/host"
)

var pwm = gammapwm.GammaPWM{Bus: `/dev/i2c-1`, Address: 0x33}

func handleLEDset(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	chans := params["channel"]
	vals := params["value"]

	channel, err := strconv.ParseInt(chans, 10, 32)
	if err != nil {
		fmt.Fprintln(w, "could not parse channel")
		return
	}
	if channel < 0 || channel > 7 {
		fmt.Fprintln(w, "invalid channel")
		return
	}

	value, err := strconv.ParseInt(vals, 10, 32)
	if err != nil {
		fmt.Fprintln(w, "could not parse value")
		return
	}

	if strings.HasPrefix(vals, `+`) || strings.HasPrefix(vals, `-`) {
		value = int64(pwm.Value[channel]) + value
	}

	if value > 100 {
		value = 100
	}
	if value < 0 {
		value = 0
	}

	pwm.Value[channel] = byte(value)

	err = pwm.Update()

	if err != nil {
		fmt.Fprintln(w, "error: ", err)
	} else {
		handleLED(w, r)
	}
}

func handleLED(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	chans, ok := params["channel"]

	if !ok || chans == `all` {
		fmt.Fprintln(w, pwm.Value)
	} else {
		channel, err := strconv.ParseInt(chans, 10, 32)
		if err != nil {
			fmt.Fprintln(w, "could not parse channel")
			return
		}
		if channel < 0 || channel > 7 {
			fmt.Fprintln(w, "invalid channel")
			return
		}
		fmt.Fprintln(w, pwm.Value[channel])
	}
}

func main() {
	if _, err := host.Init(); err != nil {
		log.Fatal(err)
	}

	pwm.Init()

	rtr := mux.NewRouter()
	rtr.HandleFunc("/led/set/{channel:[0-7]}/{value:(?:\\+|-)?[0-9]+}", handleLEDset).Methods("GET")
	rtr.HandleFunc("/led/get", handleLED).Methods("GET")
	rtr.HandleFunc("/led/get/{channel:(?:[0-7]|all)}", handleLED).Methods("GET")
	http.Handle("/", rtr)

	log.Println("Listening...")
	http.ListenAndServe(":8000", nil)
}

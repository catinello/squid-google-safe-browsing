package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log/syslog"
	"net"
	"net/http"
	"net/url"
	"os"
	"time"
)

//constants
const version = "0.1"
const name = "squid-gsb"
const envvar = "GSB_APIKEY"
const keylen = 58 //length google safe browsing api key
const timeout = 4 //timeout for http.Get in seconds

func usage() {
	fmt.Println("Usage: " + name + " [" + envvar + "]")
	fmt.Println("Version: " + version + " | License: MIT")
	os.Exit(1)
}

func main() {
	//check env variable
	key := os.Getenv(envvar)

	if len(key) == 0 {
		//check arguments
		if len(os.Args) != 2 {
			//wrong usage
			usage()
		} else {
			key = os.Args[1]
		}
	} else if len(key) != keylen {
	        fmt.Println("Is your APIKEY correct?")
                usage()
	}

	//logging to syslog
	logger, err := syslog.New(syslog.LOG_INFO, name)
	defer logger.Close()

	//syslog error check
	if err != nil {
		panic(err)
	}

	//loop forever
	for {
		//read stdin line
		bio := bufio.NewReader(os.Stdin)
		line, _, err := bio.ReadLine()

		processQuery(line, logger, key)

		if err != nil {
			logger.Crit(err.Error())
			os.Exit(2)
		}
	}
}

func processQuery(line []byte, logger *syslog.Writer, key string) {
	//split line
	sl := bytes.Split(line, []byte(" "))
	//get first blob of line
	raw, err := url.Parse(string(bytes.ToLower(sl[0])))

	if err != nil {
		logger.Warning(err.Error())
	}

	//use only scheme and host
	addr := []byte(raw.Scheme + "://" + raw.Host)
	//addr := []byte(raw.Host)

	var retval []byte
	result := askGoogle(url.QueryEscape(string(addr)), key)

	if result == 200 {
		retval = []byte("OK url=https://www.google.com/safebrowsing/diagnostic?site=" + raw.Host + " 302:https://www.google.com/safebrowsing/diagnostic?site=" + raw.Host)
		logger.Alert("Blocked Site: " + string(addr))
	} else if result > 500 {
		retval = []byte("BH message=Service Unavailable")
		logger.Crit("Service Unavailable")
	} else if result > 400 {
		retval = []byte("BH message=Not Authorized")
		logger.Crit("Not Authorized")
	} else {
		retval = []byte("ERR")
	}

	fmt.Printf("%s\n", retval)

}

// Wrapper to pass timeout constant for http.Get -> net.DialTimeout
func customDialer(timeout time.Duration) func(network, address string) (net.Conn, error) {
	return func(network, address string) (net.Conn, error) {
		c, err := net.DialTimeout(network, address, timeout)

		return c, err
	}
}

// Ask Google safe browsing with API call
func askGoogle(url, apikey string) int {
	call := fmt.Sprintf("https://sb-ssl.google.com/safebrowsing/api/lookup?client=api&apikey=%s&appver=%s&pver=3.0&url=%s", apikey, version, url)

	//http client for custom timeout
	client := http.Client{
		Transport: &http.Transport{
			Dial: customDialer(time.Duration(timeout * time.Second)),
		},
	}

	resp, err := client.Get(call)

	if err != nil {
		//resp.StatusCode = 599 //service not available
		resp = &http.Response{StatusCode: 599}
	}

	return resp.StatusCode
}

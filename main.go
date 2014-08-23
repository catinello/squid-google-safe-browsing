package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log/syslog"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
)

//constants
const version = "0.2"
const name = "squid-gsb"
const envvar = "GSB_APIKEY"
const keylen = 58 //length google safe browsing api key
const debugfile = "/tmp/squid-gsb.debug"

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

	//to debug touch /tmp/squid-gsb.debug
	var debug bool
	if _, err := os.Stat(debugfile); err == nil {
		debug = true
	}

	//loop forever
	for {
		//read stdin line
		bio := bufio.NewReader(os.Stdin)
		line, _, err := bio.ReadLine()

		processQuery(line, logger, key, debug)

		if err != nil {
			logger.Crit("Loop kill: " + err.Error())
			os.Exit(2)
		}
	}
}

func processQuery(line []byte, logger *syslog.Writer, key string, debug bool) {
	//split line
	sl := bytes.Split(line, []byte(" "))

	//get first blob of line (channel-id)
	id := string(sl[0])

	var addr []byte

	//use scheme and host
	if strings.HasPrefix(string(sl[1]), "http") {
		//get second blob of line (url)
		raw, err := url.Parse(string(bytes.ToLower(sl[1])))

		if err != nil {
			logger.Warning(err.Error())
		}

		addr = []byte(raw.Scheme + "://" + raw.Host)
	} else {
		addr = []byte("http://" + strings.Split(string(sl[1]), "/")[0])
	}

	var retval []byte
	result := askGoogle(url.QueryEscape(string(addr)), key)

	if result == 200 {
		retval = []byte(id + " OK url=https://www.google.com/safebrowsing/diagnostic?site=" + string(addr) + " 302:https://www.google.com/safebrowsing/diagnostic?site=" + string(addr))
		logger.Alert("Blocked Site: " + string(addr))
	} else if result > 500 {
		retval = []byte(id + " BH")
		logger.Crit("Service Unavailable")
	} else if result > 400 {
		retval = []byte(id + " BH")
		logger.Crit("Not Authorized")
	} else {
		retval = []byte(id + " ERR")
	}

	if debug {
		logger.Info(string(addr) + " -> " + strconv.Itoa(result) + ": " + string(retval))
	}

	fmt.Printf("%s\n", retval)
}

// Google safe browsing API call
func askGoogle(url, apikey string) int {
	call := fmt.Sprintf("https://sb-ssl.google.com/safebrowsing/api/lookup?client=api&apikey=%s&appver=%s&pver=3.0&url=%s", apikey, version, url)
	resp, err := http.Get(call)

	if err != nil {
		resp.StatusCode = 599 //service not available
	}

	return resp.StatusCode
}

package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

var targetServer string

func ProxyHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Requested URL: %s", r.URL)
	log.Printf("Will forward to %s", targetServer)

	target, err := url.Parse(targetServer)
	log.Printf("Connecting to: %s", r.URL)
	if err != nil {
		http.Error(w, "Could not parse the target URL", http.StatusInternalServerError)
		return
	}

	// Preserve the original path and query parameters using RawQuery
	target.Path = r.URL.Path
	target.RawQuery = r.URL.RawQuery // Copy query parameters as is

	log.Printf("%s %s %s\n", r.RemoteAddr, r.Method, target)

	// Create a new request to the target server
	req, err := http.NewRequest(r.Method, target.String(), r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Copy headers from original request
	for k, v := range r.Header {
		req.Header.Set(k, v[0])
	}

	// As needed, add custom headers from the proxy intended for the target server
	// req.Header.Set("Via", "1.0 go-proxligithub.com/osuritz/go-proxli")
	// req.Header.Set("X-Forwarded-For", r.RemoteAddr)

	// req.Header.Set("x-foo-header", "my-bar-value")

	// Make the request to the target server
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	// Copy response headers and body to the client
	for k, v := range resp.Header {
		w.Header().Set(k, v[0])
	}
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}

func main() {
	var (
		listenPort = flag.Int("port", 8080, "Incoming port")
		tsFlag = flag.String("target", "", "Target server")
	)

	flag.Usage = func() {
		fmt.Println("Usage: proxy [-port <port>] [-target <server>]")
		flag.PrintDefaults()
	}

	flag.Parse()

	if *tsFlag == "" {
		log.Println("Target server must be specified")
		flag.Usage()
		os.Exit(1)
	}

	if !strings.HasPrefix(*tsFlag, "http://") && !strings.HasPrefix(*tsFlag, "https://") {
		updatedTS := "http://" + *tsFlag
		tsFlag = &updatedTS
	}

	target, err := url.Parse(*tsFlag)
	if err != nil {
		log.Fatalf("Error: could not parse target URL: %v\n", err)
	}



	// if target.Scheme == "" {
	// 	target, err = url.Parse(fmt.Sprintf("http://%s", *tsFlag))
	// 	log.Printf("Target is now %s", target.String())
	// 	if err != nil {
	// 		log.Fatalf("Error: could not parse target URL: %v\n", err)
	// 	}
	// }
	targetServer = target.String()

	log.Printf("Starting proxy on port: %d, targeting %s\n", *listenPort, targetServer)

	http.HandleFunc("/", ProxyHandler)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *listenPort), nil))
}
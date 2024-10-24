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
	"time"
)

var (
	targetServer   string
	headers        map[string]string
	shouldColorize bool
)

func proxyHandler(w http.ResponseWriter, r *http.Request) {
	rqTime := time.Now()

	target, err := url.Parse(targetServer)
	if err != nil {
		http.Error(w, "Could not parse the target URL", http.StatusInternalServerError)
		return
	}

	// Preserve the original path and query parameters using RawQuery
	target.Path = r.URL.Path
	target.RawQuery = r.URL.RawQuery // Copy query parameters as is

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

	req.Header.Set("Via", "go-proxli/0.1")
	ip, err := getRemoteIP(r)
	if err == nil {
		req.Header.Set("X-Forwarded-For", ip)
	} else {
		log.Printf("Error obtaining remote IP: %v\n", err)
	}

	if len(headers) > 0 {
		for k, v := range headers {
			req.Header.Set(k, v)
		}
	}

	// As needed, add custom headers from the proxy intended for the target
	// server. E.g.
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
	bytesWritten, err := io.Copy(w, resp.Body)
	if err != nil {
		bytesWritten = -1
	}

	fmt.Println(getCLFEntry(rqTime, r, resp, bytesWritten,
		/* client-id */ "-",
		/* user-id */ "-",
		shouldColorize))
}

func main() {
	var (
		listenPort                     = flag.Int("port", 8080, "Incoming port")
		tsFlag                         = flag.String("target", "", "Target server")
		noColor                        = flag.Bool("noColor", false, "Disable colorized output")
		extraHeaders map[string]string = make(map[string]string)
	)

	// Register the `--header` flag(s) with a custom func to collect multiple values
	flag.Func("header", "Additional header to be sent to the target server", func(h string) error {
		tokens := strings.Split(h, ":")
		if len(tokens) != 1 && len(tokens) != 2 {
			return fmt.Errorf("invalid error format: %s", h)
		}
		if len(tokens) == 1 {
			extraHeaders[strings.TrimSpace(tokens[0])] = ""
		} else {
			extraHeaders[strings.TrimSpace(tokens[0])] = strings.TrimSpace(tokens[1])
		}
		return nil
	})

	flag.Usage = func() {
		fmt.Println("Usage: proxy [-port <port>] [-target <server>] [--noColor] [--header '<name>: <value>']")
		flag.PrintDefaults()
	}

	flag.Parse()
	localAddr := fmt.Sprintf("localhost:%d", *listenPort)
	log.Printf("Starting go-proxli on %s...", localAddr)

	shouldColorize = !*noColor

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

	targetServer = target.String()

	// Copy the extra headers to the global headers map
	headers = extraHeaders

	format := "[go-proxli] Running on http://%s → %s\n"
	if !*noColor {
		format = "[go-proxli] Running on " + ColorCyanNormal("http://%s") + " → " + ColorCyanBright("http://%s") + "\n"
	}
	log.Printf(format, localAddr, targetServer)

	http.HandleFunc("/", proxyHandler)
	log.Fatal(http.ListenAndServe(localAddr, nil))
}

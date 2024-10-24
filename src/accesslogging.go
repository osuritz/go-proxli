package main

import (
	"fmt"
	"net"
	"net/http"
	"time"
)

// A plain text version of the Common Log Format (CLF).
const CLFNoColor = "%s %s %s [%s] \"%s %s %s\" %d %s"

// An ANSI-color version of the Common Log Format (CLF).
var CLFColorized =
/* IP addr. */ ColorWhiteNormal("%s") + " " +
	/* client ID */ ColorWhiteNormal("%s") + " " +
	/* user ID */ ColorWhiteNormal("%s") + " " +
	/* timestamp */ ColorWhiteNormal("[%s]") + " " +
	/* HTTP method */ "\"" + ColorMagentaNormal("%s") + " " +
	/* path+querystring */ ColorMagentaNormal("%s") + " " +
	/* HTTP protocol */ "%s\" " +
	/* status code */ ColorCyanNormal("%d") + " " +
	/* payload size (bytes) */ ColorWhiteBright("%s")

func getCLFEntry(t time.Time, rq *http.Request, resp *http.Response, bSize int64, clientId string, userId string, colorize bool) string {
	format := CLFNoColor
	if colorize {
		format = CLFColorized
	}

	cLen := "-"
	if (*resp).ContentLength < 0 {
		cLen = fmt.Sprintf("%d", bSize)
	} else {
		cLen = fmt.Sprintf("%d", resp.ContentLength)
	}

	return fmt.Sprintf(format,
		getIPOrDefault(rq, "-"),
		clientId,
		userId,
		t.Format("02/Jan/2006:15:04:05 -0700"),
		rq.Method,
		rq.RequestURI,
		rq.Proto,
		resp.StatusCode,
		cLen)
}

func getRemoteIP(r *http.Request) (string, error) {
	forwarded := r.Header.Get("X-Forwarded-For")
	if forwarded != "" {
		return forwarded, nil
	}

	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	return ip, err
}

func getIPOrDefault(r *http.Request, defaultVal string) string {
	ip, err := getRemoteIP(r)
	if err != nil {
		return defaultVal
	}
	return ip
}

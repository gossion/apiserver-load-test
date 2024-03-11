package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"time"
)

var (
	concurrentConnections int
	varbose               bool
	apiServer             string
	podName               string
	podNamespace          string
)

type ChanType string

const (
	Connected    ChanType = "connected"
	Disconnected ChanType = "disconnected"
)

type ChanConf struct {
	Type ChanType // connected, disconnected
}

const (
	tokenFile  = "/var/run/secrets/kubernetes.io/serviceaccount/token"
	rootCAFile = "/var/run/secrets/kubernetes.io/serviceaccount/ca.crt"
)

func main() {
	// parse args with flag package
	// flag.Parse()
	flag.StringVar(&apiServer, "api-server", "https://<APISERVER_IP>:<PORT>", "API server URL")
	flag.IntVar(&concurrentConnections, "concurrent-connections", 10, "Number of concurrent connections to establish")
	flag.BoolVar(&varbose, "verbose", false, "Print verbose logs")
	flag.StringVar(&podName, "pod-name", "example", "pod name")
	flag.StringVar(&podNamespace, "pod-namespace", "default", "pod namespace")

	flag.Parse()

	// read token
	token, err := ioutil.ReadFile(tokenFile)
	if err != nil {
		panic(err.Error())
	}

	// Channel to count successful connections
	chanconf := make(chan ChanConf)

	// Function to establish a connection to the API server
	connectToAPIServer := func() {
		sleepRandom(0, 120) // don't push api server too hard

		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // For testing purposes only; not recommended for production
			// Disable HTTP/2.
			TLSNextProto: make(map[string]func(authority string, c *tls.Conn) http.RoundTripper),
		}
		client := &http.Client{Transport: tr}

		connected := false
		for {
			// Create a new request
			req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/v1/namespaces/%s/pods/%s", apiServer, podNamespace, podName), nil)
			if err != nil {
				fmt.Printf("Error creating request: %s\n", err)
				if connected {
					chanconf <- ChanConf{Type: Disconnected}
					connected = false
				}

				sleepRandom(60, 120)
				continue
			}

			// Set the authorization header
			req.Header.Set("Authorization", "Bearer "+string(token))

			// Perform the request
			resp, err := client.Do(req)
			if err != nil {
				fmt.Printf("Error performing request: %s\n", err)
				if connected {
					chanconf <- ChanConf{Type: Disconnected}
					connected = false
				}
				sleepRandom(60, 120)
				continue
			}

			// Read the response body
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				fmt.Printf("Error reading response body: %s\n", err)
				if connected {
					chanconf <- ChanConf{Type: Disconnected}
					connected = false
				}

				resp.Body.Close()
				sleepRandom(60, 120)
				continue
			}

			// Count successful connections
			if resp.StatusCode == http.StatusOK {
				if varbose {
					fmt.Printf("established a connection and received a response, body: %v\n", string(body))
				}
				if !connected {
					chanconf <- ChanConf{Type: Connected}
					connected = true
				}
			} else {
				fmt.Printf("received non-OK response: %s\n", string(body))
				if connected {
					chanconf <- ChanConf{Type: Disconnected}
					connected = false
				}
			}

			resp.Body.Close()
			sleepRandom(60, 120)
		}
	}

	// Initialize the connections
	for i := 0; i < concurrentConnections; i++ {
		go connectToAPIServer()
	}

	// Count and output the number of successful connections
	count := 0
	for c := range chanconf {
		if c.Type == Disconnected {
			count--
		} else if c.Type == Connected {
			count++
		}
		fmt.Printf("Number of successful connections: %d\n", count)
	}
}

func sleepRandom(min, max int) {
	rand.Seed(time.Now().UnixNano())
	time.Sleep(time.Duration(rand.Intn(max-min)+min) * time.Second)
}

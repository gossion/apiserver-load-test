package main

import (
	"context"
	"flag"
	"fmt"
	"math/rand"
	"time"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

var (
	concurrentConnections int
	varbose               bool
	reconnectOnFailure    bool
)

type ChanType string

const (
	Connected    ChanType = "connected"
	Disconnected ChanType = "disconnected"
)

type ChanConf struct {
	Type ChanType // connected, disconnected
}

func main() {
	// parse args with flag package
	// flag.Parse()
	flag.IntVar(&concurrentConnections, "concurrent-connections", 1000, "Number of concurrent connections to establish")
	flag.BoolVar(&varbose, "verbose", false, "Print verbose logs")
	flag.BoolVar(&reconnectOnFailure, "reconnect-on-failure", false, "reconnection on failure")
	flag.Parse()

	// Channel to count successful connections
	chanconf := make(chan ChanConf)

	// Function to establish a connection to the API server
	connectToAPIServer := func() {
		sleepRandom() // don't push api server too hard

		config, err := rest.InClusterConfig()
		if err != nil {
			panic(err.Error())
		}

		// Create the clientset using the in-cluster config
		clientset, err := kubernetes.NewForConfig(config)
		if err != nil {
			panic(err.Error())
		}

		connected := false
		for !connected {
			// Here we will just list pods as an example of a long-lived connection,
			watcher, err := clientset.CoreV1().Pods("kube-system").Watch(context.Background(), metav1.ListOptions{})
			if err != nil {
				fmt.Printf("Error watching pods: %v\n", err)

				time.Sleep(2 * time.Second) // sleep so it won't put too much pressure on the API server
				continue
			}
			chanconf <- ChanConf{Type: Connected}
			connected = true

			for event := range watcher.ResultChan() {
				if varbose {
					pod := event.Object.(*v1.Pod)
					fmt.Printf("Event: %v %v\n", event.Type, pod.Name)
				}
			}

			// if the watcher is closed, we will try to reconnect if the reconnectOnFailure flag is set
			fmt.Println("Watcher closed.")
			if reconnectOnFailure {
				connected = false // reset state so it will be reconnected
			}
			chanconf <- ChanConf{Type: Disconnected}
			time.Sleep(2 * time.Second)
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

func sleepRandom() {
	// sleep for a random interval between 0 and 60 seconds
	rand.Seed(time.Now().UnixNano())
	time.Sleep(time.Duration(rand.Intn(60)) * time.Second)
}

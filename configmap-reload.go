package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"encoding/json"

	fsnotify "gopkg.in/fsnotify.v1"
)

var volumeDirs volumeDirsFlag
var webhookStatusCode = flag.Int("webhook-status-code", 200, "the HTTP status code indicating successful triggering of reload")
var jenkins = flag.String("jenkins-url", "", "the jenkins url")
var username = flag.String("username", "", "the jenkins username")
var password = flag.String("password", "", "the jenkins password")
var crumbIssuerPath = "crumbIssuer/api/json"
var cascReload = "configuration-as-code/reload"

func main() {
	flag.Var(&volumeDirs, "volume-dir", "the config map volume directory to watch for updates; may be used multiple times")
	flag.Parse()

	if len(volumeDirs) < 1 {
		log.Println("Missing volume-dir")
		log.Println()
		flag.Usage()
		os.Exit(1)
	}
	if *jenkins == "" {
		log.Println("Missing jenkins")
		log.Println()
		flag.Usage()
		os.Exit(1)
	}
	if *username == "" {
		log.Println("Missing username")
		log.Println()
		flag.Usage()
		os.Exit(1)
	}
	if *password == "" {
		log.Println("Missing password")
		log.Println()
		flag.Usage()
		os.Exit(1)
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	done := make(chan bool)
	go func() {
		for {
			select {
			case event := <-watcher.Events:
				if event.Op&fsnotify.Create == fsnotify.Create {
					if filepath.Base(event.Name) == "..data" {
						log.Println("config map updated")
						req, err := http.NewRequest("GET", *jenkins + crumbIssuerPath, nil)
						req.SetBasicAuth(*username, *password)
						if err != nil {
							log.Println("error:", err)
							continue
						}
						resp, err := http.DefaultClient.Do(req)
						if err != nil {
							log.Println("error:", err)
							continue
						}
						crumbData := map[string]string{}
						json.NewDecoder(resp.Body).Decode(&crumbData)
						resp.Body.Close()

						req, err = http.NewRequest("POST", *jenkins + cascReload, nil)
						req.SetBasicAuth(*username, *password)
						req.Header.Add(crumbData["crumbRequestField"], crumbData["crumb"])
						if err != nil {
							log.Println("error:", err)
							continue
						}
						resp, err = http.DefaultClient.Do(req)
						if err != nil {
							log.Println("error:", err)
							continue
						}

						if resp.StatusCode != *webhookStatusCode {
							log.Println("error:", "Received response code", resp.StatusCode, ", expected", *webhookStatusCode)
							continue
						}
						log.Println("successfully triggered reload")
					}
				}
			case err := <-watcher.Errors:
				log.Println("error:", err)
			}
		}
	}()

	for _, d := range volumeDirs {
		log.Printf("Watching directory: %q", d)
		err = watcher.Add(d)
		if err != nil {
			log.Fatal(err)
		}
	}
	<-done
}

type volumeDirsFlag []string

func (v *volumeDirsFlag) Set(value string) error {
	*v = append(*v, value)
	return nil
}

func (v *volumeDirsFlag) String() string {
	return fmt.Sprint(*v)
}

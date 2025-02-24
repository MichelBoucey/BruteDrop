package main

import (
	"flag"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

var commitHash = ""

var config Config

func main() {

	version := "3.0.0"

	var httpLines []string

	var sshLines []string

	var invalidUser = regexp.MustCompile(`^(.*?\d{2}:\d{2}:\d{2}).*?invalid\suser\s(\w+)\s(\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3})\sport\s\d{1,5}`)

	var http404Error = regexp.MustCompile(`^(\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}).*?"\s404\s\d*\s"`)

	var http404ErrorsCount = make(map[string]int)

	versionFlag := flag.Bool("version", false, "Show version")

	configFilepathFlag := flag.String("configuration-filepath", "/etc/brutedrop.conf", "Full configuration file path")

	flag.Parse()

	if *versionFlag == true {
		fmt.Println("brutedrop v" + version + " (" + commitHash + ")\nCopyright @2025 Michel Boucey\nReleased under 3-Clause BSD License")
		os.Exit(0)
	}

	// Get and check BruteDrop configuration
	data, err := ioutil.ReadFile(*configFilepathFlag)
	if err != nil {
		log.Fatal(err)
	}
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		log.Fatal(err)
	}
	if _, err := os.Stat(config.Iptables); os.IsNotExist(err) {
		log.Fatal("Can't find iptables at path " + config.Iptables)
	}
	if _, err := os.Stat(config.Journalctl); os.IsNotExist(err) {
		log.Fatal("Can't find journalctl at path " + config.Journalctl)
	}
	// AuthorizedUsers and AuthorizedAddresses can't be both empty
	if len(config.AuthorizedUsers) == 0 && len(config.AuthorizedAddresses) == 0 {
		log.Fatal("To run brutedrop you have to add authorized users and/or IP addresses in /etc/brutedrop.conf")
	}

	// Set how to log
	if config.LoggingTo != "stdout" {
		lf, err := os.OpenFile(config.LoggingTo, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			log.Fatalf("error opening file: %v", err)
		}
		defer lf.Close()
		log.SetOutput(lf)
	} else {
		log.SetFlags(0)
	}

	//
	// Ban IP addresses for SSH login attempt
	//

	// Okay, now get some of the latest log lines of failed SSH login attempts from journalctl
	sshOut, err := exec.Command("sh", "-c", config.Journalctl+" --since \""+strconv.Itoa(config.LogEntriesSince)+" minutes ago\" -u sshd --no-pager | grep invalid").Output()

	if len(sshOut) >= 0 {

		sshLines = strings.Split(string(sshOut), "\n")

		// Iterating over log lines searching invalid users
		// who fails to login to ban their IP addresses
		for i := 0; i < len(sshLines); i++ {

			if sshLines[i] != "" {

				matches := invalidUser.FindStringSubmatch(sshLines[i])

				if len(matches) == 4 {

					if isElement(matches[2], config.AuthorizedUsers) {

						log.Println("Authorized user " + matches[2] + " failed to login from " + matches[3] + " at " + matches[1])

					} else if !isElement(matches[3], config.AuthorizedAddresses) {

						if !isAlreadyBanned(matches[3]) {
							dropCommand := config.Iptables + " -w -A INPUT -s " + matches[3] + " -j DROP"
							if config.DryRun == false {
								err := exec.Command("sh", "-c", dropCommand).Run()
								if err != nil {
									log.Fatal("Can't execute \"" + dropCommand + "\"")
								}
								// log.Println("Ban " + matches[2] + "@" + matches[3] + " at " + matches[1])
								log.Println("Ban " + matches[3] + " for SSH login attempt as " + matches[2] + " at " + matches[1])
							} else {
								log.Println("BruteDrop is currently in dry run mode (" + dropCommand + ")")
							}
						}

					} else {

						log.Println("Invalid user " + matches[2] + " from authorized IP address " + matches[3] + " at " + matches[1])

					}
				}
			}
		}
	}

	//
	// Ban IP addresses for too many HTTP 404 Errors
	//

	httpOut, err := exec.Command("sh", "-c", "tail -n 300 /opt/WebSites/volumes/nginx/logs/access.log").Output()

	if len(httpOut) >= 0 {

		httpLines = strings.Split(string(httpOut), "\n")

		for i := 0; i < len(httpLines); i++ {

			a404Error := http404Error.FindStringSubmatch(httpLines[i])

			if len(a404Error) == 2 {

				http404ErrorsCount[a404Error[1]]++

				if http404ErrorsCount[a404Error[1]] == config.MaxHTTP404Errors {

					if !isAlreadyBanned(a404Error[1]) {
						dropCommand := config.Iptables + " -w -A INPUT -s " + a404Error[1] + " -j DROP"
						if config.DryRun == false {
							err := exec.Command("sh", "-c", dropCommand).Run()
							if err != nil {
								log.Fatal("Can't execute \"" + dropCommand + "\"")
							}
							log.Println("Ban " + a404Error[1] + " for too many HTTP 404 errors")
						} else {
							log.Println("BruteDrop is currently in dry run mode (" + dropCommand + ")")
						}
					}
				}
			}
		}
	}
}

func isAlreadyBanned (ipAddr string) bool {

	_, err := exec.Command("sh", "-c", config.Iptables+" -w -C INPUT -s "+ ipAddr +" -j DROP").Output()

	if err == nil { return true } else { return false }

}


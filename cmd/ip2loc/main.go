package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/meinside/ipstack-go"
)

const (
	configFilepath = ".config/ip2loc.json" // $HOME/.config/ip2loc.json
)

type config struct {
	AccessKey string `json:"access_key"`
	IsPremium bool   `json:"is_premium"`
}

// loggers
var _stdout = log.New(os.Stdout, "", 0)
var _stderr = log.New(os.Stderr, "", 0)

// load config file
func loadConf() (conf config, err error) {
	var usr *user.User
	if usr, err = user.Current(); err == nil {
		fpath := filepath.Join(usr.HomeDir, configFilepath)

		var bytes []byte
		if bytes, err = ioutil.ReadFile(fpath); err == nil {
			if err = json.Unmarshal(bytes, &conf); err == nil {
				return conf, nil
			}
		}
	}

	return config{}, err
}

// print error and exit
func printErrorAndExit(err error) {
	_stderr.Fatalf(err.Error())
}

func printRes(res ipstack.Response) {
	strs := []string{}

	// ip/hostname
	var host string
	if res.IP != res.Hostname && res.Hostname != "" {
		host = fmt.Sprintf("%s / %s", res.IP, res.Hostname)
	} else {
		host = res.IP
	}

	// other information
	if res.CountryName != "" {
		strs = append(strs, res.CountryName)
	}
	if res.RegionName != "" {
		strs = append(strs, res.RegionName)
	}
	if res.City != "" {
		strs = append(strs, res.City)
	}

	// print them
	if len(strs) > 0 {
		_stdout.Printf("%s (%s)\n", host, strings.Join(strs, ", "))
	} else {
		_stdout.Printf("%s\n", host)
	}
}

// print result and exit
func printResponseAndExit(res ipstack.Response) {
	printRes(res)

	os.Exit(0)
}

// print results and exit
func printResponsesAndExit(res []ipstack.Response) {
	// print each response
	for _, r := range res {
		printRes(r)
	}

	os.Exit(0)
}

func main() {
	var conf config
	var err error

	if conf, err = loadConf(); err == nil {
		client := ipstack.NewClient(conf.AccessKey, !conf.IsPremium)
		switch len(os.Args) {
		case 1: // no params
			var res ipstack.Response
			if res, err = client.LookupRequester(); err == nil {
				printResponseAndExit(res)
			}
		case 2: // one param
			param := os.Args[1]

			var res ipstack.Response
			if res, err = client.LookupStandard(param); err == nil {
				printResponseAndExit(res)
			}
		default: // more than one
			params := os.Args[1:]

			var res []ipstack.Response
			if res, err = client.LookupBulk(params); err == nil {
				printResponsesAndExit(res)
			}
		}
	}

	if err != nil {
		printErrorAndExit(err)
	}
}

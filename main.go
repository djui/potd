package main

import (
	"path/filepath"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"time"
	"os"
)

type AppCredentials struct {
	consumerKey string
	consumerSecret string
	accessToken string
	accessTokenSecret string
}

type ParsedTweet struct {
	time time.Time
	text string
	url string
}

const configFilename = ".potdrc"
const template = "\033[1;30m@%s\033[0m: %s \033[4m%s\033[0m"
const shouldDebug = true
const lineIndentation = 2
const maxLineLength = 50 - lineIndentation

var accounts = []string{
	"adriancolyer", // "themorningpaper",
	"onepaperperday",
	"onecspaperaday",
}

func init() {
	log.SetFlags(0)
}

func debug(errs ...error) {
	for _, err := range errs {
		if err != nil && shouldDebug {
			log.Printf("Error: %v\n", err)
		}
	}
}

func checkArgument(cond bool, expl string) {
	if !cond {
		log.Fatal(expl)
	}
}

func loadCredentials(filename string) (credentials AppCredentials, err error) {
	c, err := ioutil.ReadFile(filename)
	if err != nil {
		return
	}
	lines := strings.Split(string(c), "\n")
	return AppCredentials{
		consumerKey: lines[0],
		consumerSecret: lines[1],
		accessToken: lines[2],
		accessTokenSecret: lines[3],
	}, nil
}

func main() {
	var configFilepath string
	if dir, err := filepath.Abs(filepath.Dir(os.Args[0])); err != nil {
		configFilepath = filepath.Join(".", configFilename)
	} else {
		configFilepath = filepath.Join(dir, configFilename)
	}

	credentials, err := loadCredentials(configFilepath)
	checkArgument(err == nil, fmt.Sprintf("Could not parse %s file: %v", configFilename, err))

	client := makeClient(credentials)

	for i, account := range accounts {
		tweets, err := getLatestTweets(client, account, 1)
		checkArgument(len(tweets) > 0, "No tweets found")

		tweet, err := parseTweet(tweets[0])
		checkArgument(err == nil, "Could not parse tweet")

		output := fmt.Sprintf(template, account, tweet.text, tweet.url)
		fmt.Println(indentMultilines(breakLongLine(output, maxLineLength), lineIndentation))
		if i < len(accounts) - 1 {
			fmt.Println()
		}
	}
}

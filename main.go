package main

import (
	"fmt"
	"github.com/ChimeraCoder/anaconda"
	"io/ioutil"
	"log"
	"net/url"
	"strings"
	"time"
	"os"
)

type appCredentials struct {
	consumerKey string
	consumerSecret string
	accessToken string
	accessTokenSecret string
}

type parsedTweet struct {
	time time.Time
	text string
	url string
}

const credentialsFilename = "CREDENTIALS"
const shouldDebug = true


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

func checkArgument(cond bool) {
	if cond {
		os.Exit(1)
	}
}

func loadCredentials(filename string) (credentials appCredentials, err error) {
	c, err := ioutil.ReadFile(filename)
	if err != nil {
		return
	}
	lines := strings.Split(string(c), "\n")
	return appCredentials{
		consumerKey: lines[0],
		consumerSecret: lines[1],
		accessToken: lines[2],
		accessTokenSecret: lines[3],
	}, nil
}

func getLatestTweets(client *anaconda.TwitterApi, screenName string, length int) (tweets []anaconda.Tweet, err error) {
	query := url.Values{}
	query.Set("screen_name", screenName)
	query.Set("count", fmt.Sprintf("%v", length))
	query.Set("trim_user", "true")
	query.Set("exclude_replies", "true")
	query.Set("include_rts", "true")
	return client.GetUserTimeline(query)
}

func parseTweet(tweet anaconda.Tweet) (t parsedTweet, err error) {
	time, err:= time.Parse(time.RubyDate, tweet.CreatedAt)
	if err != nil {
		return
	}

	if retweet := tweet.RetweetedStatus; retweet != nil {
		tweet = *retweet
	}

	text := strings.Replace(tweet.Text, "\n", " ", -1)

	var url string
	if entities := tweet.Entities; &entities != nil {
	    if urls := entities.Urls; &urls != nil && len(urls) > 0 {
			url = urls[0].Expanded_url
		}
	}

	return parsedTweet{
		time: time,
		text: text,
		url: url,
	}, nil
}

func main() {
	credentials, err := loadCredentials(credentialsFilename)
	if err != nil {
		log.Fatalf("Could not parse %s file: %v\n", credentialsFilename, err)
	}

	anaconda.SetConsumerKey(credentials.consumerKey)
	anaconda.SetConsumerSecret(credentials.consumerSecret)
	client := anaconda.NewTwitterApi(credentials.accessToken, credentials.accessTokenSecret)

	onepaperperdayTweets, err1 := getLatestTweets(client, "onepaperperday", 1)
	onecspaperadayTweets, err2 := getLatestTweets(client, "onecspaperaday", 1)
	debug(err1, err2)
	checkArgument(len(onepaperperdayTweets) == 0 || len(onecspaperadayTweets) == 0)
	
	onepaperperdayTweet, err1 := parseTweet(onepaperperdayTweets[0])
	onecspaperadayTweet, err2 := parseTweet(onecspaperadayTweets[0])
	debug(err1, err2)
	checkArgument(err1 != nil || err2 != nil)

	var tweet parsedTweet
	if onepaperperdayTweet.time.After(onecspaperadayTweet.time) {
		tweet = onepaperperdayTweet
	} else {
		tweet = onecspaperadayTweet
	}

	fmt.Printf("%s %s\n", tweet.text, tweet.url)
}

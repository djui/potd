package main

import (
	"fmt"
	"github.com/ChimeraCoder/anaconda"
	"io/ioutil"
	"log"
	"net/url"
	"strings"
)

type appCredentials struct {
	consumerKey string
	consumerSecret string
	accessToken string
	accessTokenSecret string
}

const credentialsFilename = "CREDENTIALS"


func init() {
	log.SetFlags(0)
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

func getTimeline(client *anaconda.TwitterApi, screenName string) (tweets []anaconda.Tweet, err error) {
	query := url.Values{}
	query.Set("screen_name", screenName)
	query.Set("count", "5")
	query.Set("trim_user", "true")
	query.Set("exclude_replies", "true")
	query.Set("include_rts", "true")
	return client.GetUserTimeline(query)
}

func main() {
	var credentials appCredentials
	var err error
	
	if credentials, err = loadCredentials(credentialsFilename); err != nil {
		log.Fatalf("Could not parse %s file: %v\n", credentialsFilename, err)
	}
	anaconda.SetConsumerKey(credentials.consumerKey)
	anaconda.SetConsumerSecret(credentials.consumerSecret)
	//client := anaconda.NewTwitterApi(credentials.accessToken, credentials.accessTokenSecret)
	client := anaconda.NewTwitterApi("", "")
	
	var tweets []anaconda.Tweet
	if tweets, err = getTimeline(client, "onepaperperday"); err != nil {
		log.Fatal(err)
	}
	
	for _, tweet := range tweets {
		var text string
		if retweet := tweet.RetweetedStatus; retweet != nil {
			text = retweet.Text
		} else {
			text = tweet.Text
		}
		text = strings.Replace(text, "\n", " ", -1)
		fmt.Printf("Text: %v\n", text)
	}
}

package main

import (
	"fmt"
	"github.com/ChimeraCoder/anaconda"
	"net/url"
	"strings"
	"time"
)

type parsedTweet struct {
	time time.Time
	text string
	url string
}

func makeClient(credentials AppCredentials) *anaconda.TwitterApi {
	anaconda.SetConsumerKey(credentials.consumerKey)
	anaconda.SetConsumerSecret(credentials.consumerSecret)
	return anaconda.NewTwitterApi(credentials.accessToken, credentials.accessTokenSecret)
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

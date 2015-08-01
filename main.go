package main

import (
	"fmt"
	"github.com/kurrik/oauth1a"
	"github.com/kurrik/twittergo"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

const credentialsFilename = "CREDENTIALS"

type credentials struct {
	consumerKey       string
	consumerSecret    string
	accessToken       string
	accessTokenSecret string
}


func init() {
	log.SetFlags(0)
}	

func loadCredentials(credentialsFilename string) (cred credentials, err error) {
	c, err := ioutil.ReadFile(credentialsFilename)
	if err != nil {
		return
	}
	lines := strings.Split(string(c), "\n")
	cred = credentials{
		consumerKey:       lines[0],
		consumerSecret:    lines[1],
		accessToken:       lines[2],
		accessTokenSecret: lines[3],
	}
	return
}

func newTwitterClient(cred credentials) *twittergo.Client {
	config := &oauth1a.ClientConfig{
		ConsumerKey:    cred.consumerKey,
		ConsumerSecret: cred.consumerSecret,
	}
	auth := oauth1a.NewAuthorizedConfig(cred.accessToken, cred.accessTokenSecret)
	return twittergo.NewClient(config, auth)	
}

func main() {
	cred, err := loadCredentials(credentialsFilename)
	if err != nil {
		log.Fatal("Could not parse %s file: %v\n", credentialsFilename, err)
	}
	client := newTwitterClient(cred)
	
	req, _ := http.NewRequest("GET", "/1.1/account/verify_credentials.json", nil)
	resp, err := client.SendRequest(req)
	if err != nil {
		log.Fatal("Could not send request: %v\n", err)
	}
	
	user := &twittergo.User{}
	err = resp.Parse(user)
	if err != nil {
		log.Fatal("Problem parsing response: %v\n", err)
	}
	
	fmt.Printf("ID:                   %v\n", user.Id())
	fmt.Printf("Name:                 %v\n", user.Name())
	if resp.HasRateLimit() {
		fmt.Printf("Rate limit:           %v\n", resp.RateLimit())
		fmt.Printf("Rate limit remaining: %v\n", resp.RateLimitRemaining())
		fmt.Printf("Rate limit reset:     %v\n", resp.RateLimitReset())
	} else {
		fmt.Printf("Could not parse rate limit from response.\n")
	}
}

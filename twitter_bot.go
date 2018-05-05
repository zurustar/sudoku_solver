package main

import (
	"encoding/json"
	"fmt"
	"github.com/ChimeraCoder/anaconda"
	"io/ioutil"
	"log"
	"net/url"
	"os/exec"
	"regexp"
	"strings"
)

type Config struct {
	Username          string `json:"username"`
	ConsumerKey       string `json:"consumer_key"`
	ConsumerSecret    string `json:"consumer_secret"`
	AccessToken       string `json:"access_token"`
	AccessTokenSecret string `json:"access_token_secret"`
}

func load(filename string) *Config {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}
	config := new(Config)
	err = json.Unmarshal(bytes, &config)
	if err != nil {
		log.Fatal(err)
	}
	if strings.HasPrefix(config.Username, "@") == false {
		config.Username = "@" + config.Username
	}
	return config
}

func run() {
	re := regexp.MustCompile(`[^0-9]`)
	config := load("../twitter_conf.json")
	anaconda.SetConsumerKey(config.ConsumerKey)
	anaconda.SetConsumerSecret(config.ConsumerSecret)
	api := anaconda.NewTwitterApi(config.AccessToken, config.AccessTokenSecret)
	v := url.Values{}
	v.Set("track", config.Username)
	stream := api.PublicStreamFilter(v)
	for {
		select {
		case stream := <-stream.C:
			switch tweet := stream.(type) {
			case anaconda.Tweet:
				s := strings.Replace(tweet.Text, config.Username, "", -1)
				s = re.ReplaceAllString(s, "")
				result := ""
				if len(s) != 81 {
					result = "問題がなにかおかしい気がします。"
				} else {
					out, err := exec.Command("./sudoku_solver", s).Output()
					if err != nil {
						result = "私には解けませんでした。ごめんなさい。"
					} else {
						result = string(out)
					}
				}
				result = "@" + tweet.User.ScreenName + "\n" + result
				posted, err := api.PostTweet(result, nil)
				if err != nil {
					fmt.Println(err)
				} else {
					fmt.Println(posted.Text)
				}
			default:
			}
		}
	}
}

func main() {
	run()
}

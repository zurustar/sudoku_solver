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

const ConfigFilePath = "/etc/sudoku_solver.json"
const SudokuSolverFilePath = "/usr/local/bin/sudoku_solver/sudoku_solver"

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
	config := load(ConfigFilePath)
	anaconda.SetConsumerKey(config.ConsumerKey)
	anaconda.SetConsumerSecret(config.ConsumerSecret)
	api := anaconda.NewTwitterApi(config.AccessToken, config.AccessTokenSecret)
	v := url.Values{}
	v.Set("track", config.Username)
	stream := api.PublicStreamFilter(v)
	fmt.Println("ok")
	for {
		select {
		case stream := <-stream.C:
			switch tweet := stream.(type) {
			case anaconda.Tweet:
				s := strings.Replace(tweet.Text, config.Username, "", -1)
				s = re.ReplaceAllString(s, "")
				result := ""
				if len(s) != 81 {
					result = "問題がおかしい気がします。"
				} else {
					out, err := exec.Command(SudokuSolverFilePath, s).Output()
					if err != nil {
						result = "私には解けませんでした。ごめんなさい。"
					} else {
						result = string(out)
					}
				}
				result = "@" + tweet.User.ScreenName + "\n" + result
				v := url.Values{}
				v.Set("in_reply_to_status_id", tweet.IdStr)
				posted, err := api.PostTweet(result, v)
				if err != nil {
					fmt.Println("ERROR ->", err)
				} else {
					fmt.Println("tweeted ->", posted.Text)
				}
			default:
			}
		}
	}
}

func main() {
	run()
}

package main

import (
	"os"
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

const DefaultConfigFilePath = "./config/sudoku_solver.json"

type Config struct {
	Username          string `json:"username"`
	ConsumerKey       string `json:"consumer_key"`
	ConsumerSecret    string `json:"consumer_secret"`
	AccessToken       string `json:"access_token"`
	AccessTokenSecret string `json:"access_token_secret"`
	SudokuSolverCommand string `json:"sudoku_solver_command"`
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

func run(config_file_path string) {
	re := regexp.MustCompile(`[^0-9]`)
	config := load(config_file_path)
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
					out, err := exec.Command(config.SudokuSolverCommand, s).Output()
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
	if len(os.Args) == 1 {
		run(DefaultConfigFilePath)
	} else {
		run(os.Args[1])
	}
}

package main

import (
	"encoding/json"
	"fmt"
	"github.com/ChimeraCoder/anaconda"
	"io/ioutil"
	"log"
	"log/syslog"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type BoardState int

const (
	SOLVED BoardState = iota
	NOT_SOLVED
	INVALID
)

var dic = [9 * 9][3]int{
	{0, 0, 0}, {1, 0, 0}, {2, 0, 0},
	{3, 0, 1}, {4, 0, 1}, {5, 0, 1},
	{6, 0, 2}, {7, 0, 2}, {8, 0, 2},
	{0, 1, 0}, {1, 1, 0}, {2, 1, 0},
	{3, 1, 1}, {4, 1, 1}, {5, 1, 1},
	{6, 1, 2}, {7, 1, 2}, {8, 1, 2},
	{0, 2, 0}, {1, 2, 0}, {2, 2, 0},
	{3, 2, 1}, {4, 2, 1}, {5, 2, 1},
	{6, 2, 2}, {7, 2, 2}, {8, 2, 2},
	{0, 3, 3}, {1, 3, 3}, {2, 3, 3},
	{3, 3, 4}, {4, 3, 4}, {5, 3, 4},
	{6, 3, 5}, {7, 3, 5}, {8, 3, 5},
	{0, 4, 3}, {1, 4, 3}, {2, 4, 3},
	{3, 4, 4}, {4, 4, 4}, {5, 4, 4},
	{6, 4, 5}, {7, 4, 5}, {8, 4, 5},
	{0, 5, 3}, {1, 5, 3}, {2, 5, 3},
	{3, 5, 4}, {4, 5, 4}, {5, 5, 4},
	{6, 5, 5}, {7, 5, 5}, {8, 5, 5},
	{0, 6, 6}, {1, 6, 6}, {2, 6, 6},
	{3, 6, 7}, {4, 6, 7}, {5, 6, 7},
	{6, 6, 8}, {7, 6, 8}, {8, 6, 8},
	{0, 7, 6}, {1, 7, 6}, {2, 7, 6},
	{3, 7, 7}, {4, 7, 7}, {5, 7, 7},
	{6, 7, 8}, {7, 7, 8}, {8, 7, 8},
	{0, 8, 6}, {1, 8, 6}, {2, 8, 6},
	{3, 8, 7}, {4, 8, 7}, {5, 8, 7},
	{6, 8, 8}, {7, 8, 8}, {8, 8, 8},
}

type Board struct {
	cells [9 * 9][]int
}

func NewBoard() *Board {
	b := new(Board)
	for row := 0; row < 9; row++ {
		for column := 0; column < 9; column++ {
			b.cells[row*9+column] = []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
		}
	}
	return b
}

func (board *Board) ToString() string {
	s := ""
	for row := 0; row < 9; row++ {
		if row%3 == 0 {
			s += "\n"
		}
		for column := 0; column < 9; column++ {
			i := row*9 + column
			if len(board.cells[i]) == 1 {
				s += strconv.Itoa(board.cells[i][0])
			} else {
				s += "0"
			}
			if column%3 == 2 {
				s += " "
			}
		}
		s += "\n"
	}
	return s
}

func (board *Board) ShowDetail() {
	s := ""
	for row := 0; row < 9; row++ {
		if row%3 == 0 {
			s += "+---+---+---+\n"
		}
		for column := 0; column < 9; column++ {
			if column%3 == 0 {
				s += "|"
			}
			i := row*9 + column
			if len(board.cells[i]) == 1 {
				s += strconv.Itoa(board.cells[i][0])
			} else {
				s += "0"
			}
		}
		s += "|"
		for column := 0; column < 9; column++ {
			if column%3 == 0 {
				s += " "
			}
			s += strconv.Itoa(len(board.cells[row*9+column]))
		}
		s += "\n"
	}
	s += "+---+---+---+"
	fmt.Println(s)
}

func (board *Board) CopyFrom(src *Board) {
	for pos := 0; pos < 9*9; pos++ {
		board.cells[pos] = src.cells[pos]
	}
}

func (board *Board) Find(pos, value int) int {
	for i, v := range board.cells[pos] {
		if v == value {
			return i
		}
	}
	return -1
}

func (board *Board) Remove(pos, value int) bool {
	result := false
	cands := []int{}
	for _, v := range board.cells[pos] {
		if v == value {
			result = true
		} else {
			cands = append(cands, v)
		}
	}
	if result {
		board.cells[pos] = cands
	}
	return result
}

func (board *Board) _Update1() bool {
	updated := false
	for src_pos := 0; src_pos < 9*9; src_pos++ {
		if len(board.cells[src_pos]) != 1 {
			continue
		}
		for dst_pos := 0; dst_pos < 9*9; dst_pos++ {
			if src_pos == dst_pos {
				continue
			}
			for target := 0; target < 3; target++ {
				if dic[src_pos][target] == dic[dst_pos][target] {
					if board.Remove(dst_pos, board.cells[src_pos][0]) {
						updated = true
					}
				}
			}
		}
	}
	return updated
}

func (board *Board) _Update2() bool {
	updated := false
	for src_pos := 0; src_pos < 9*9; src_pos++ {
		if len(board.cells[src_pos]) == 1 {
			continue
		}
		for target := 0; target < 3; target++ {
			for _, cand := range board.cells[src_pos] {
				found := false
				for dst_pos := 0; dst_pos < 9*9; dst_pos++ {
					if src_pos == dst_pos {
						continue
					}
					if dic[src_pos][target] != dic[dst_pos][target] {
						continue
					}
					if board.Find(dst_pos, cand) >= 0 {
						found = true
					}
				}
				if found == false {
					board.cells[src_pos] = []int{cand}
					updated = true
				}
			}
		}
	}
	return updated
}

func (board *Board) Update() {
	updated := true
	for updated {
		updated = board._Update1()
		if !updated {
			updated = board._Update2()
		}
		if !updated {
			break
		}
	}
}

func Solve(board *Board) (BoardState, *Board) {
	board.Update()
	len_list := [10][]int{{}, {}, {}, {}, {}, {}, {}, {}, {}, {}}
	for pos := 0; pos < 9*9; pos++ {
		l := len(board.cells[pos])
		len_list[l] = append(len_list[l], pos)
	}
	if len(len_list[0]) > 0 {
		return INVALID, nil
	}
	if len(len_list[1]) == 9*9 {
		return SOLVED, board
	}
	//board.ShowDetail()
	for i := 2; i < 10; i++ {
		for _, pos := range len_list[i] {
			if len(board.cells[pos]) > 1 {
				for _, cand := range board.cells[pos] {
					new_board := NewBoard()
					new_board.CopyFrom(board)
					new_board.cells[pos] = []int{cand}
					result, new_board := Solve(new_board)
					if result == SOLVED {
						return SOLVED, new_board
					}
				}
			}
		}
	}
	return NOT_SOLVED, board
}

func RunSolver(src string) string {
	re := regexp.MustCompile(`[^0-9]`)
	q := re.ReplaceAllString(src, "")
	if len(q) != 9*9 {
		log.Println("invalid question")
		return ""
	}
	log.Println("solve", q)
	board := NewBoard()
	for i, c := range q {
		v, err := strconv.Atoi(string(c))
		if err != nil {
			// bug?
			return ""
		}
		if v != 0 {
			board.cells[i] = []int{v}
		}
	}
	result, board := Solve(board)
	if result == SOLVED {
		return board.ToString()
	}
	return ""
}

func Load(filename string) string {
	log.Println("load", filename)
	fp, err := os.Open(filename)
	if err != nil {
		return ""
	}
	defer fp.Close()
	result := ""
	buf := make([]byte, 1024)
	for {
		n, err := fp.Read(buf)
		if n == 0 {
			break
		}
		if err != nil {
			return ""
		}
		result += string(buf[:n])
	}
	return result
}

type Config struct {
	Username            string `json:"username"`
	ConsumerKey         string `json:"consumer_key"`
	ConsumerSecret      string `json:"consumer_secret"`
	AccessToken         string `json:"access_token"`
	AccessTokenSecret   string `json:"access_token_secret"`
	SudokuSolverCommand string `json:"sudoku_solver_command"`
}

func LoadBotConfiguration(filename string) *Config {
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

func RunTwitterBot(config_filename string) {
	logger, _ := syslog.New(syslog.LOG_NOTICE|syslog.LOG_USER, "twitter_bot")
	log.SetOutput(logger)
	config := LoadBotConfiguration(config_filename)
	re := regexp.MustCompile(`[^0-9]`)
	anaconda.SetConsumerKey(config.ConsumerKey)
	anaconda.SetConsumerSecret(config.ConsumerSecret)
	api := anaconda.NewTwitterApi(config.AccessToken, config.AccessTokenSecret)
	v := url.Values{}
	v.Set("track", config.Username)
	stream := api.PublicStreamFilter(v)
	log.Println("ok", config.Username)
	for {
		select {
		case stream := <-stream.C:
			switch tweet := stream.(type) {
			case anaconda.Tweet:
				s := strings.Replace(tweet.Text, config.Username, "", -1)
				s = re.ReplaceAllString(s, "")
				log.Println("received", s, "from", tweet.User.ScreenName)
				result := ""
				if len(s) != 81 {
					result = "問題がおかしい気がします。"
				} else {
					result = RunSolver(s)
				}
				result = "@" + tweet.User.ScreenName + "\n" + result
				v := url.Values{}
				v.Set("in_reply_to_status_id", tweet.IdStr)
				posted, err := api.PostTweet(result, v)
				if err != nil {
					log.Println("ERROR ->", err)
				} else {
					fmt.Println("tweeted ->", posted.Text)
				}
			default:
			}
		}
	}

}

func main() {
	if len(os.Args) == 3 {
		switch os.Args[1] {
		case "-d":
			RunTwitterBot(os.Args[2])
		case "-f":
			fmt.Println(RunSolver(Load(os.Args[2])))
		case "-q":
			fmt.Println(RunSolver(os.Args[2]))
		default:
			log.Fatal("invalid parameter ", os.Args[1])
		}
	}
}

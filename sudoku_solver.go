package main

import (
	"encoding/json"
	"fmt"
	"github.com/ChimeraCoder/anaconda"
	"io/ioutil"
	"log"
	"log/syslog"
	"math"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type BoardState int

const (
	SOLVED BoardState = iota
	NOT_SOLVED
	INVALID
)

type UpdateResult int

const (
	UPDATED UpdateResult = iota
	NOT_UPDATED
	ERROR
)

var peers = [][]int{}

func initPeer() {
	for pos := 0; pos < 81; pos++ {
		peers = append(peers, []int{})
	}
	for sy := 0; sy < 9; sy++ {
		for sx := 0; sx < 9; sx++ {
			sp := sy*9 + sx
			for dy := 0; dy < 9; dy++ {
				for dx := 0; dx < 9; dx++ {
					dp := dy*9 + dx
					if sp == dp {
						continue
					}
					if sx == dx {
						peers[sp] = append(peers[sp], dp)
					} else if sy == dy {
						peers[sp] = append(peers[sp], dp)
					} else {
						sb := math.Trunc(float64(sy)/3.0)*3 + math.Trunc(float64(sx)/3.0)
						db := math.Trunc(float64(dy)/3.0)*3 + math.Trunc(float64(dx)/3.0)
						if sb == db {
							peers[sp] = append(peers[sp], dp)
						}
					}
				}
			}
		}
	}
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
	for i := 0; i < 9*9; i++ {
		if len(board.cells[i]) == 1 {
			s += strconv.Itoa(board.cells[i][0])
		} else {
			s += "0"
		}
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

//
// 指定された場所のセルの候補リストから指定された値を削除する
//
func (board *Board) Remove(pos, value int) UpdateResult {
	found := false
	cands := []int{}
	for _, v := range board.cells[pos] {
		if v == value {
			found = true
		} else {
			cands = append(cands, v)
		}
	}
	if found {
		if len(cands) == 0 {
			return ERROR // 消したらデータがなくなってしまった
		}
		board.cells[pos] = cands
		return UPDATED // 消した
	}
	return NOT_UPDATED // 該当のデータがなかったので消さなかった
}

//
// 各セルに着目し、候補数があとひとつに減っていたら、そのセルはその値で
// 確定ということになるので、横方向、縦方向、同じグループ内の他のセルの
// 候補からその値を削除する
//
func (board *Board) _Update1() UpdateResult {
	updated := false
	for sp := 0; sp < 81; sp++ {
		if len(board.cells[sp]) != 1 {
			continue
		}
		for _, dp := range peers[sp] {
			switch board.Remove(dp, board.cells[sp][0]) {
			case UPDATED:
				updated = true
			case NOT_UPDATED:
				// do nothing
			case ERROR:
				return ERROR
			}
		}
	}
	if updated {
		return UPDATED
	}
	return NOT_UPDATED
}

//
// 各セルの各候補に着目し、その値が横方向、縦方向、同じグループ内の他の
// セルの候補にその値が含まれていなかったら、その値に確定
//
func (board *Board) _Update2() bool {
	updated := false
	for sp := 0; sp < 81; sp++ {
		if len(board.cells[sp]) == 1 {
			continue
		}
		for _, v := range board.cells[sp] {
			found := false
			for _, dp := range peers[sp] {
				if board.Find(dp, v) >= 0 {
					found = true
					break
				}
			}
			if found == false {
				board.cells[sp] = []int{v}
				updated = true
			}
		}
	}
	return updated
}

//
// 各セルの候補を絞り込む
//
func (board *Board) Update() UpdateResult {
	updated := true
	for updated {

		switch board._Update1() {
		case UPDATED:
			updated = true
		case NOT_UPDATED:
			updated = board._Update2()
		case ERROR:
			return ERROR
		}
	}
	return UPDATED
}

//
// 解く。ルールどおりに試して答えがでなかったら仮置きして再帰呼び出し
//
func Solve(board *Board) (BoardState, *Board) {
	board.Update()
	// 各セルそれぞれにのこりいくつの候補が残っているかを調査
	// 残りの個数ごとにそのセルの番号を配列に格納する
	len_list := [10][]int{{}, {}, {}, {}, {}, {}, {}, {}, {}, {}}
	for pos := 0; pos < 9*9; pos++ {
		l := len(board.cells[pos])
		len_list[l] = append(len_list[l], pos)
	}
	// 候補数がゼロのセルが存在したら異常な状態になっている
	if len(len_list[0]) > 0 {
		return INVALID, nil
	}
	// すべてのセルの候補数が1になっていたら解けている
	if len(len_list[1]) == 9*9 {
		return SOLVED, board
	}
	// まだ候補数が１になっていないセルについて、
	// 適当にいっこの値に絞り込んで、ためしに解いてみる
	// 再帰呼び出しになるので、結果的にかたっぱしから試すことになる
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

//
// 81個の数値からなる文字列を人間が見やすい形に変換
//
func ToString(src string) string {
	s := ""
	for row := 0; row < 9; row++ {
		if row%3 == 0 {
			s += "\n"
		}
		for column := 0; column < 9; column++ {
			s += string(src[row*9+column])
			if column%3 == 2 {
				s += " "
			}
		}
		s += "\n"
	}
	return s
}

//
// 問題文字列から余計な文字をカットしてから
// 問題情報を保持するBoardインスタンスに突っ込んで解く
//
func RunSolver(src string) string {
	re := regexp.MustCompile(`[^0-9]`)
	q := re.ReplaceAllString(src, "")
	if len(q) != 9*9 {
		log.Println("invalid question")
		return ""
	}
	board := NewBoard()
	for i, c := range q {
		v, err := strconv.Atoi(string(c))
		if err != nil {
			return "" // bug?
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

// 問題が書いてあるファイルを読み込む
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

type TwitterConfig struct {
	Username          string `json:"username"`
	ConsumerKey       string `json:"consumer_key"`
	ConsumerSecret    string `json:"consumer_secret"`
	AccessToken       string `json:"access_token"`
	AccessTokenSecret string `json:"access_token_secret"`
}

// TwitterBot用設定ファイルのロード処理
func LoadBotConfiguration(filename string) *TwitterConfig {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}
	config := new(TwitterConfig)
	err = json.Unmarshal(bytes, &config)
	if err != nil {
		log.Fatal(err)
	}
	// 設定ファイルのユーザ名のあたまに@がなかったらつける
	if strings.HasPrefix(config.Username, "@") == false {
		config.Username = "@" + config.Username
	}
	return config
}

//
// TwitterBotとして動く
//
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
					start := time.Now()
					result = RunSolver(s)
					end := time.Now()
					sec := float64(end.Sub(start).Nanoseconds()) / 1000000000.0
					result = fmt.Sprintf("こたえは\n%s\nだと思います。%f秒で解けました。",
						ToString(result), sec)
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
			}
		}
	}
}

func main() {
	if len(os.Args) == 3 {
		initPeer()
		switch os.Args[1] {
		case "-d":
			RunTwitterBot(os.Args[2])
		case "-f":
			fmt.Println(ToString(RunSolver(Load(os.Args[2]))))
		case "-q":
			fmt.Println(ToString(RunSolver(os.Args[2])))
		default:
			log.Fatal("invalid parameter ", os.Args[1])
		}
	}
}

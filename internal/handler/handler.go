package handler

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"shuangpin-tutorial-go/internal/data"
)

const strTipsSeparator = "———"

// Options holds the configuration for the tutorial.
type Options struct {
	WordsDisplayed int    // number of words to display each round (default 6)
	RemindOff      bool   // turn off input reminder
	ZhuyinOff      bool   // turn off zhuyin display
	CwordOff       bool   // turn off Chinese word display
}

// DefaultOptions returns the default options.
func DefaultOptions() Options {
	return Options{
		WordsDisplayed: 6,
		RemindOff:      true, // reverse of Ruby's OPTIONS_DEFAULT[:r] = true (meaning reminder is ON by default)
		ZhuyinOff:      true, // reverse of Ruby's OPTIONS_DEFAULT[:z] = true (meaning zhuyin is ON by default)
		CwordOff:       true, // reverse of Ruby's OPTIONS_DEFAULT[:c] = true (meaning cword is ON by default)
	}
}

// Handler manages the tutorial session.
type Handler struct {
	opts           Options
	spData         *data.Data
	wordsDisplayed int
	isEnableRemind bool
	isDisplayZhuyin bool
	isDisplayCword  bool
	wordsCounts    [3]int // total, correct, incorrect
	startTime      time.Time
	endTime        time.Time
	totalTime      float64 // seconds
	typeSpeed      float64 // seconds per correct word
}

// New creates a new Handler with the given options.
func New(opts Options) *Handler {
	return &Handler{
		opts:            opts,
		wordsDisplayed:  opts.WordsDisplayed,
		isEnableRemind:  !opts.RemindOff,
		isDisplayZhuyin: !opts.ZhuyinOff,
		isDisplayCword:  !opts.CwordOff,
	}
}

// Start begins the tutorial session.
func (h *Handler) Start() {
	h.spData = data.NewData()
	h.startTime = time.Now()

	fmt.Println("歡迎使用 快速學習雙拼 -- 微軟雙拼（繁體中文版）！要好好練習哦！（1.3.0 版）")
	fmt.Println("要輸入下面中文字對應雙拼字碼哦！使用空格隔開，要休息請輸入 'exit' 或按 ctrl+c 吧！")

	scanner := bufio.NewScanner(os.Stdin)
	for {
		resultWords, resultDisplay := h.getQuestionLines(h.wordsDisplayed)
		fmt.Printf("\n%s\n", strings.Join(resultDisplay, "\n"))

		if !scanner.Scan() {
			break // EOF or error
		}
		inputAnswers := scanner.Text()

		lowerInput := strings.ToLower(strings.TrimSpace(inputAnswers))
		if lowerInput == "exit" {
			fmt.Println("累了就該休息一下吧！")
			return
		}

		h.evalResult(inputAnswers, resultWords)
	}
}

// Stop ends the tutorial session and calculates statistics.
func (h *Handler) Stop() {
	h.endTime = time.Now()
	h.totalTime = h.endTime.Sub(h.startTime).Seconds()
}

// GetQuestionLines generates a question with display lines.
func (h *Handler) getQuestionLines(num int) ([]string, []string) {
	resultWords, resultDisplay := h.getWordsShowlines(num)

	if h.isEnableRemind {
		tips := make([]string, 0, len(resultWords))
		for _, res := range resultWords {
			tips = append(tips, data.GetCorrectKeys(res))
		}
		resultDisplay[0] = fmt.Sprintf("%s  %s  【%s】", resultDisplay[0], strTipsSeparator, strings.Join(tips, " "))
	}

	return resultWords, resultDisplay
}

// GetWordsShowlines generates display lines for the given number of words.
func (h *Handler) getWordsShowlines(num int) ([]string, []string) {
	result := make([]string, 0, num)
	displayPy := make([]string, 0, num)
	displayZy := make([]string, 0, num)
	displayCw := make([]string, 0, num)

	for i := 0; i < num; i++ {
		sheng, yun := h.spData.GetRandShengYunSym()
		shengyun := data.GetSy(sheng, yun)
		zy := data.GetZy(sheng, yun)
		cw := data.GetCw(sheng, yun)

		result = append(result, shengyun)
		displayPy = append(displayPy, shengyun+"\t")
		displayZy = append(displayZy, zy+"\t")
		displayCw = append(displayCw, cw+"\t")
	}

	display := make([]string, 0, 3)
	display = append(display, strings.Join(displayPy, "|"))
	if h.isDisplayZhuyin {
		display = append(display, strings.Join(displayZy, "|"))
	}
	if h.isDisplayCword {
		display = append(display, strings.Join(displayCw, "|"))
	}

	return result, display
}

// EvalResult evaluates the user's input against the correct answers.
func (h *Handler) evalResult(inputAnswer string, resultWords []string) {
	incorrect := make(map[int]string)
	correct := 0
	isAllFinished := len(strings.Fields(inputAnswer)) >= len(resultWords)

	inputFields := strings.Fields(inputAnswer)
	for index, inputKeys := range inputFields {
		if index >= len(resultWords) {
			break // ignore extra input
		}

		if len(inputKeys) != 2 {
			incorrect[index+1] = fmt.Sprintf("%s： 雙拼只能由兩個字母組合而成！", inputKeys)
			continue
		}

		res := resultWords[index]
		rightAnswer := data.GetCorrectKeys(res)
		if !strings.EqualFold(inputKeys, rightAnswer) {
			incorrect[index+1] = fmt.Sprintf("%s： 不能得到 %s 哦！或許你可以試試 %s", inputKeys, res, rightAnswer)
		} else {
			correct++
		}
	}

	h.wordsCounts[0] += len(resultWords)
	h.wordsCounts[1] += correct
	h.wordsCounts[2] += len(resultWords) - correct

	if !isAllFinished {
		fmt.Printf("你好像寫不齊哦！共有：%d，已輸入：%d\n", len(resultWords), len(inputFields))
	} else if len(incorrect) == 0 {
		fmt.Println("恭喜你全對耶！")
	}

	for key, value := range incorrect {
		fmt.Printf("%d) 錯誤 => %s\n", key, value)
	}
}

// PrintResults prints the session statistics.
func (h *Handler) PrintResults() {
	totalTimeStr := h.getTotalTimeString()
	if totalTimeStr != "" {
		fmt.Println(totalTimeStr)
	}

	totalWordsResult := h.getTotalWordsResult()
	if totalWordsResult != "" {
		fmt.Println(totalWordsResult)
	}

	percentCorrectness := h.getPercentCorrectness()
	if percentCorrectness != "" {
		fmt.Println(percentCorrectness)
	}

	typeSpeedStr := h.getTypeSpeedString()
	if typeSpeedStr != "" {
		fmt.Println(typeSpeedStr)
	}
}

func (h *Handler) getTotalTimeString() string {
	if h.wordsCounts[0] == 0 {
		return ""
	}

	totalSeconds := int(h.totalTime)
	hours := totalSeconds / 3600
	minutes := (totalSeconds % 3600) / 60
	seconds := totalSeconds % 60

	var strTotal string
	switch {
	case h.totalTime < 60:
		strTotal = fmt.Sprintf("%d 秒", seconds)
	case h.totalTime < 3600:
		strTotal = fmt.Sprintf("%d 分 %d 秒", minutes, seconds)
	default:
		strTotal = fmt.Sprintf("%d 小時 %d 分 %d 秒", hours, minutes, seconds)
	}

	return fmt.Sprintf("訓練時間： %s", strTotal)
}

func (h *Handler) getTotalWordsResult() string {
	if h.wordsCounts[0] == 0 {
		return ""
	}
	return fmt.Sprintf("共有字數：%d 字（%d－%d）", h.wordsCounts[0], h.wordsCounts[1], h.wordsCounts[2])
}

func (h *Handler) getPercentCorrectness() string {
	if h.wordsCounts[0] == 0 {
		return ""
	}
	percent := float64(h.wordsCounts[1]) / float64(h.wordsCounts[0]) * 100
	return fmt.Sprintf("總正確率： %.2f ％", percent)
}

func (h *Handler) getTypeSpeedString() string {
	if h.wordsCounts[0] == 0 || h.wordsCounts[1] == 0 {
		return ""
	}
	h.typeSpeed = h.totalTime / float64(h.wordsCounts[1]) // only correct words
	return fmt.Sprintf("平均速度： %.2f 秒／字", h.typeSpeed)
}

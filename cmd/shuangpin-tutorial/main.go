package main

import (
	"flag"
	"fmt"
	"os"

	"shuangpin-tutorial-go/internal/handler"
)

const version = "1.3.0"

func main() {
	w := flag.Int("w", 6, "Words display in each section")
	r := flag.Bool("r", false, "Turn off the input reminder (defaults on)")
	z := flag.Bool("z", false, "Turn off the Taiwan zhuyin display (defaults on)")
	c := flag.Bool("c", false, "Turn off the Chinese word display (defaults on)")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "快速學習雙拼 -- 微軟雙拼（繁體中文版） v%s\n\n", version)
		fmt.Fprintf(os.Stderr, "Usage: shuangpin-tutorial [OPTIONS]\n")
		fmt.Fprintf(os.Stderr, "\nOptions:\n")
		flag.PrintDefaults()
	}

	flag.Parse()

	opts := handler.DefaultOptions()
	opts.WordsDisplayed = *w
	opts.RemindOff = *r
	opts.ZhuyinOff = *z
	opts.CwordOff = *c

	h := handler.New(opts)

	defer func() {
		h.Stop()
		fmt.Println("要努力學習哦，感謝使用！")
	}()

	h.Start()
	h.PrintResults()
}

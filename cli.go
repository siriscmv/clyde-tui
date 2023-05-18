package main

import (
	"fmt"
	"time"

	"github.com/briandowns/spinner"
)

var CLIChan chan string = make(chan string)

func RunCLI(prompt string) {
	s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
	s.Suffix = " Processing"
	s.Color("magenta")
	s.Start()

	<-Ready
	go AskClyde(prompt, "Answer the question while being as specific and short as possible. DO NOT add extra details. Use Markdown when needed")

	resp := <-CLIChan
	_, formatted := FormatClydeReponse(resp)

	s.Stop()
	fmt.Print(formatted)
}

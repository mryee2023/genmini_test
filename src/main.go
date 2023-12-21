package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"genmini/src/vars"
	"github.com/fatih/color"
	"github.com/go-resty/resty/v2"
	"gopkg.in/yaml.v3"

	markdown "github.com/MichaelMure/go-term-markdown"
)

var cli = resty.New().SetTimeout(time.Second * 10)
var configFile = flag.String("f", "etc/config.yaml", "the config file")
var config vars.GenMiniConfig

func sendToGenMini(flag chan struct{}, s string) (string, error) {
	url := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/gemini-pro:generateContent?key=%s", config.AppKey)
	req := vars.GenminiRequest{
		Contents: []vars.Contents{
			{
				Parts: vars.Parts{
					Text: s,
				},
			},
		},
	}
	res, err := cli.
		SetTimeout(time.Second * 10).
		R().
		SetBody(req).
		Post(url)
	if err != nil {
		return "", err
	}
	go func() {
		flag <- struct{}{}
	}()

	return res.String(), nil
}

func main() {

	flag.Parse()
	content, err := os.ReadFile(*configFile)
	if err != nil {
		color.Red("parse config failure, %s", err.Error())
		os.Exit(1)
	}
	err = yaml.Unmarshal(content, &config)
	if err != nil {
		panic(err)
	}
	reader := bufio.NewReader(os.Stdin)
	color.Cyan("Hello GenMini Pro......")
	color.Cyan("---------------------")
	//request.SetTimeout(time.Second * 10)
	flagCh := make(chan struct{})
	for {
		fmt.Print("-> ")
		text, _ := reader.ReadString('\n')
		text = strings.Replace(text, "\n", "", -1)
		if strings.EqualFold(text, "exit") {
			color.HiCyan("Bye")
			break
		}
		if strings.EqualFold(text, "") {
			continue
		}
		thinking(text, flagCh)
	}
}
func thinking(input string, flag chan struct{}) {
	var v string
	var e error
	go func() {
		v, e = sendToGenMini(flag, input)
	}()
	ts := time.NewTicker(time.Second)
	timeout := time.NewTimer(time.Second * 10)
	for {
		select {
		case <-flag:
			if e != nil {
				color.Red(e.Error())
				return
			}
			var res vars.GenminiResponse
			if e = json.Unmarshal([]byte(v), &res); e != nil {
				color.Red("Ops, something wrong, ", e.Error())
				return
			}
			if len(res.Candidates) == 0 {
				color.Red("Ops, something wrong")
				return
			}

			source := res.Candidates[0].Content.Parts[0].Text
			result := markdown.Render(source, 80, 6)
			//clearLastLine()
			fmt.Println()
			fmt.Println(string(result))
			fmt.Println()
			ts.Stop()
			return
		case <-ts.C:
			fmt.Print(".")
		case <-timeout.C:
			color.Red("Ops, something wrong, timeout")
		}
	}
}

func clearLastLine() {
	// 向上移动光标一行
	fmt.Print("\033[1A")
	// 清除光标所在行
	fmt.Print("\033[K")
}

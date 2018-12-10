package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

type App struct {
	APIToken string
}

func NewApp(APIToken string) *App {
	return &App{
		APIToken: APIToken,
	}
}

// func loadTestData() []byte {
// 	data, err := ioutil.ReadFile("slack-emoji-list.json")
// 	if err != nil {
// 		panic(err)
// 	}
// 	return data
// }

func fetchEmojiList() ([]byte, error) {
	url := "https://slack.com/api/emoji.list?token=" + app.APIToken + "&pretty=1"
	req, _ := http.NewRequest("GET", url, nil)

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Slack API - emoji.list <error>", err)
		return nil, err
	}

	resBody, _ := ioutil.ReadAll(res.Body)
	return resBody, nil
}

func downloadEmoji(name string, url string) error {
	if strings.HasPrefix(url, "alias:") {
		return nil
	}

	// Fetch Request
	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	res, err := client.Do(req)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to download emoji image:", err)
		return err
	}

	if res.StatusCode != 200 {
		message := fmt.Sprintf("HTTP Request Error: %v %v", res.StatusCode, res.Status)
		return errors.New(message)
	}

	// Read Response Body
	resBody, _ := ioutil.ReadAll(res.Body)

	filename := "emojis/" + name + ".png"
	err = ioutil.WriteFile(filename, resBody, 0644)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to write a file:", filename, err)
		return err
	}
	return nil
}

// Run is
func (app *App) Run() int {
	err := os.MkdirAll("emojis", os.ModePerm)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to create a \"emojis\" folder:", err)
		return ExitCodeError
	}

	bytes, err := fetchEmojiList()
	if err != nil {
		return ExitCodeError
	}

	var data interface{}
	if err := json.Unmarshal(bytes, &data); err != nil {
		fmt.Fprintln(os.Stderr, "Slack API - <error> json unmarshal:", err)
		return ExitCodeError
	}

	count := 1
	emoji := data.(map[string]interface{})["emoji"]
	for name, url := range emoji.(map[string]interface{}) {
		fmt.Println(count, name)
		err := downloadEmoji(name, url.(string))
		if err != nil {
			return ExitCodeError
		}
		count++
	}

	return ExitCodeOK
}

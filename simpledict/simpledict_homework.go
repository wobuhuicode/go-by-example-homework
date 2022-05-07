package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

type CaiyunDictRequest struct {
	TransType string `json:"trans_type"`
	Source    string `json:"source"`
	UserID    string `json:"user_id"`
}

type CaiyunDictResponse struct {
	Rc   int `json:"rc"`
	Wiki struct {
		KnownInLaguages int `json:"known_in_laguages"`
		Description     struct {
			Source string      `json:"source"`
			Target interface{} `json:"target"`
		} `json:"description"`
		ID   string `json:"id"`
		Item struct {
			Source string `json:"source"`
			Target string `json:"target"`
		} `json:"item"`
		ImageURL  string `json:"image_url"`
		IsSubject string `json:"is_subject"`
		Sitelink  string `json:"sitelink"`
	} `json:"wiki"`
	Dictionary struct {
		Prons struct {
			EnUs string `json:"en-us"`
			En   string `json:"en"`
		} `json:"prons"`
		Explanations []string      `json:"explanations"`
		Synonym      []string      `json:"synonym"`
		Antonym      []string      `json:"antonym"`
		WqxExample   [][]string    `json:"wqx_example"`
		Entry        string        `json:"entry"`
		Type         string        `json:"type"`
		Related      []interface{} `json:"related"`
		Source       string        `json:"source"`
	} `json:"dictionary"`
}

func CaiyunQuery(word string, cancel context.CancelFunc) {
	client := &http.Client{}
	request := CaiyunDictRequest{TransType: "en2zh", Source: word}
	buf, err := json.Marshal(request)
	if err != nil {
		log.Fatal(err)
	}
	var data = bytes.NewReader(buf)
	req, err := http.NewRequest("POST", "https://api.interpreter.caiyunai.com/v1/dict", data)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("DNT", "1")
	req.Header.Set("os-version", "")
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/99.0.4844.51 Safari/537.36")
	req.Header.Set("app-name", "xy")
	req.Header.Set("Content-Type", "application/json;charset=UTF-8")
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("device-id", "")
	req.Header.Set("os-type", "web")
	req.Header.Set("X-Authorization", "token:qgemv4jr1y38jyq6vhvi")
	req.Header.Set("Origin", "https://fanyi.caiyunapp.com")
	req.Header.Set("Sec-Fetch-Site", "cross-site")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Referer", "https://fanyi.caiyunapp.com/")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Cookie", "_ym_uid=16456948721020430059; _ym_d=1645694872")
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	bodyText, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	if resp.StatusCode != 200 {
		log.Fatal("bad StatusCode:", resp.StatusCode, "body", string(bodyText))
	}
	var dictResponse CaiyunDictResponse
	err = json.Unmarshal(bodyText, &dictResponse)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(word, "UK:", dictResponse.Dictionary.Prons.En, "US:", dictResponse.Dictionary.Prons.EnUs)
	for _, item := range dictResponse.Dictionary.Explanations {
		fmt.Println(item)
	}

	cancel()
}

type YoudaoResp struct {
	Basic struct {
		ExamType    []string `json:"exam_type"`
		Explains    []string `json:"explains"`
		Phonetic    string   `json:"phonetic"`
		Uk_phonetic string   `json:"uk-phonetic"`
		Uk_speech   string   `json:"uk-speech"`
		Us_phonetic string   `json:"us-phonetic"`
		Us_speech   string   `json:"us-speech"`
		Wfs         []struct {
			Wf struct {
				Name  string `json:"name"`
				Value string `json:"value"`
			} `json:"wf"`
		} `json:"wfs"`
	} `json:"basic"`
	Dict struct {
		URL string `json:"url"`
	} `json:"dict"`
	ErrorCode    string   `json:"errorCode"`
	IsWord       bool     `json:"isWord"`
	L            string   `json:"l"`
	Query        string   `json:"query"`
	RequestID    string   `json:"requestId"`
	ReturnPhrase []string `json:"returnPhrase"`
	SpeakURL     string   `json:"speakUrl"`
	TSpeakURL    string   `json:"tSpeakUrl"`
	Translation  []string `json:"translation"`
	Web          []struct {
		Key   string   `json:"key"`
		Value []string `json:"value"`
	} `json:"web"`
	Webdict struct {
		URL string `json:"url"`
	} `json:"webdict"`
}

func YoudaoQuery(word string, cancel context.CancelFunc) {
	client := &http.Client{}
	youdaoReqBody := fmt.Sprintf("q=%s&from=en&to=zh-CHS", word)

	// var data = strings.NewReader(`q=sexy&from=en&to=zh-CHS`)
	data := strings.NewReader(youdaoReqBody)
	req, err := http.NewRequest("POST", "https://aidemo.youdao.com/trans", data)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("authority", "aidemo.youdao.com")
	req.Header.Set("accept", "application/json, text/javascript, */*; q=0.01")
	req.Header.Set("accept-language", "zh-CN,zh;q=0.9,en;q=0.8,en-GB;q=0.7,en-US;q=0.6")
	req.Header.Set("content-type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Set("origin", "https://ai.youdao.com")
	req.Header.Set("referer", "https://ai.youdao.com/")
	req.Header.Set("sec-ch-ua", `" Not A;Brand";v="99", "Chromium";v="101", "Microsoft Edge";v="101"`)
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("sec-ch-ua-platform", `"Windows"`)
	req.Header.Set("sec-fetch-dest", "empty")
	req.Header.Set("sec-fetch-mode", "cors")
	req.Header.Set("sec-fetch-site", "same-site")
	req.Header.Set("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/101.0.4951.41 Safari/537.36 Edg/101.0.1210.32")
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	bodyText, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	var youdaoResp = YoudaoResp{}
	err = json.Unmarshal(bodyText, &youdaoResp)
	if err != nil {
		log.Fatal(err)
	}
	// fmt.Printf("%s\n", bodyText)
	for _, item := range youdaoResp.Basic.Explains {
		fmt.Println(item)
	}

	cancel()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, `usage: simpleDict WORD
example: simpleDict hello
		`)
		os.Exit(1)
	}
	word := os.Args[1]

	ctx, cancel := context.WithCancel(context.Background())

	go CaiyunQuery(word, cancel)
	go YoudaoQuery(word, cancel)

	<-ctx.Done()

}

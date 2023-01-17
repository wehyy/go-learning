package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

type DictRequest struct {
	TransType string `json:"trans_type"`
	Source    string `json:"source"`
}

type DictResponse struct {
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

type VolcRequest struct {
	Source         string   `json:"source"`
	Words          []string `json:"words"`
	SourceLanguage string   `json:"source_language"`
	TargetLanguage string   `json:"target_language"`
}

type VolcResponse struct {
	Details []struct {
		Detail string `json:"detail"`
		Extra  string `json:"extra"`
	} `json:"details"`
	BaseResp struct {
		StatusCode    int    `json:"status_code"`
		StatusMessage string `json:"status_message"`
	} `json:"base_resp"`
}

type VolcResponseDetail struct {
	ErrorCode string `json:"errorCode"`
	RequestID string `json:"requestId"`
	Msg       string `json:"msg"`
	Result    []struct {
		Ec struct {
			ReturnPhrase []string `json:"returnPhrase"`
			Synonyms     []struct {
				Pos   string   `json:"pos"`
				Words []string `json:"words"`
				Trans string   `json:"trans"`
			} `json:"synonyms"`
			Etymology struct {
				ZhCHS []struct {
					Description string `json:"description"`
					Detail      string `json:"detail"`
					Source      string `json:"source"`
				} `json:"zh-CHS"`
			} `json:"etymology"`
			SentenceSample []struct {
				Sentence     string `json:"sentence"`
				SentenceBold string `json:"sentenceBold"`
				Translation  string `json:"translation"`
				Source       string `json:"source"`
			} `json:"sentenceSample"`
			WebDict string `json:"webDict"`
			Web     []struct {
				Phrase   string   `json:"phrase"`
				Meanings []string `json:"meanings"`
			} `json:"web"`
			MTerminalDict string `json:"mTerminalDict"`
			RelWord       struct {
				Word string `json:"word"`
				Stem string `json:"stem"`
				Rels []struct {
					Rel struct {
						Pos   string `json:"pos"`
						Words []struct {
							Word string `json:"word"`
							Tran string `json:"tran"`
						} `json:"words"`
					} `json:"rel"`
				} `json:"rels"`
			} `json:"relWord"`
			Dict  string `json:"dict"`
			Basic struct {
				UsPhonetic string   `json:"usPhonetic"`
				UsSpeech   string   `json:"usSpeech"`
				Phonetic   string   `json:"phonetic"`
				UkSpeech   string   `json:"ukSpeech"`
				ExamType   []string `json:"examType"`
				Explains   []struct {
					Pos   string `json:"pos"`
					Trans string `json:"trans"`
				} `json:"explains"`
				UkPhonetic  string `json:"ukPhonetic"`
				WordFormats []struct {
					Name  string `json:"name"`
					Value string `json:"value"`
				} `json:"wordFormats"`
			} `json:"basic"`
			Phrases []struct {
				Phrase   string   `json:"phrase"`
				Meanings []string `json:"meanings"`
			} `json:"phrases"`
			Lang   string `json:"lang"`
			IsWord bool   `json:"isWord"`
		} `json:"ec"`
	} `json:"result"`
}

func timeMeasurement(start time.Time, progname string) {
	elapsed := time.Since(start)
	fmt.Printf(progname+" Execution time: %s\n", elapsed)
}

var wg sync.WaitGroup

func caiyunfanyi(word string) {
	defer timeMeasurement(time.Now(), "caiyunfanyi")
	client := &http.Client{}
	request := DictRequest{TransType: "en2zh", Source: word}
	buf, err := json.Marshal(request)
	if err != nil {
		log.Fatal(err)
	}
	var data = bytes.NewReader(buf)
	req, err := http.NewRequest("POST", "https://api.interpreter.caiyunai.com/v1/dict", data)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("authority", "api.interpreter.caiyunai.com")
	req.Header.Set("accept", "application/json, text/plain, */*")
	req.Header.Set("accept-language", "en-US,en;q=0.9,zh-CN;q=0.8,zh;q=0.7")
	req.Header.Set("app-name", "xy")
	req.Header.Set("content-type", "application/json;charset=UTF-8")
	req.Header.Set("device-id", "")
	req.Header.Set("origin", "https://fanyi.caiyunapp.com")
	req.Header.Set("os-type", "web")
	req.Header.Set("os-version", "")
	req.Header.Set("referer", "https://fanyi.caiyunapp.com/")
	req.Header.Set("sec-ch-ua", `"Not_A Brand";v="99", "Google Chrome";v="109", "Chromium";v="109"`)
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("sec-ch-ua-platform", `"Windows"`)
	req.Header.Set("sec-fetch-dest", "empty")
	req.Header.Set("sec-fetch-mode", "cors")
	req.Header.Set("sec-fetch-site", "cross-site")
	req.Header.Set("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.36")
	req.Header.Set("x-authorization", "token:qgemv4jr1y38jyq6vhvi")
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
	var dictResponse DictResponse
	err = json.Unmarshal(bodyText, &dictResponse)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("彩云翻译:")
	fmt.Println(word, "UK:", dictResponse.Dictionary.Prons.En, "US:", dictResponse.Dictionary.Prons.EnUs)
	for _, item := range dictResponse.Dictionary.Explanations {
		fmt.Println(item)
	}
	wg.Done()
}

func huoshanfanyi(word []string) {
	defer timeMeasurement(time.Now(), "huoshanfanyi")
	client := &http.Client{}
	// var data = strings.NewReader(`{"source":"youdao","words":["excellent"],"source_language":"en","target_language":"zh"}`)
	request := VolcRequest{Source: "youdao", Words: word, SourceLanguage: "en", TargetLanguage: "zh"}
	buf, err := json.Marshal(request)
	if err != nil {
		log.Fatal(err)
	}
	var data = bytes.NewReader(buf)
	req, err := http.NewRequest("POST", "https://translate.volcengine.com/web/dict/detail/v1/?msToken=&X-Bogus=DFSzswVLQDcpnL39SZ424M9WX7j6&_signature=_02B4Z6wo000013xqzTAAAIDD.Gg3carFv4d8asmAALzrhGDhtpz8vEyf4GhaIBm4Ym41aWZRJQhvfyj6E-Wl.M42e.st2iz-TGG0vl0ckZQCFxXdJtao-4TB8wQeXIRaFzW.-1gxVfi6hph12e", data)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("authority", "translate.volcengine.com")
	req.Header.Set("accept", "application/json, text/plain, */*")
	req.Header.Set("accept-language", "en-US,en;q=0.9,zh-CN;q=0.8,zh;q=0.7")
	req.Header.Set("content-type", "application/json")
	req.Header.Set("cookie", "x-jupiter-uuid=16738780947556560; i18next=zh-CN; s_v_web_id=verify_lcyvssf6_mUEV5aFD_1A82_4USq_BMlc_AbGdWeDI9Aov; ttcid=7c781585c9c044b0afce063da48ccfbf72; tt_scid=AIbf9CnHXogUDnV90GYwJWLUZodfvfaws-5BEbsR3U4gN-4m9qFajDED.QTUHh-w95be")
	req.Header.Set("origin", "https://translate.volcengine.com")
	req.Header.Set("referer", "https://translate.volcengine.com/?category=&home_language=zh&source_language=detect&target_language=zh&text=excellent")
	req.Header.Set("sec-ch-ua", `"Not_A Brand";v="99", "Google Chrome";v="109", "Chromium";v="109"`)
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("sec-ch-ua-platform", `"Windows"`)
	req.Header.Set("sec-fetch-dest", "empty")
	req.Header.Set("sec-fetch-mode", "cors")
	req.Header.Set("sec-fetch-site", "same-origin")
	req.Header.Set("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.36")
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
	var dictResponse VolcResponse
	err = json.Unmarshal(bodyText, &dictResponse)
	if err != nil {
		log.Fatal(err)
	}
	var dictResponseDetail VolcResponseDetail
	detail := dictResponse.Details[0].Detail
	err = json.Unmarshal([]byte(detail), &dictResponseDetail)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("火山翻译:")
	fmt.Println("[US: ", dictResponseDetail.Result[0].Ec.Basic.UsPhonetic, "]", "[UK: ", dictResponseDetail.Result[0].Ec.Basic.UkPhonetic, "]")
	for i, item := range dictResponseDetail.Result[0].Ec.Basic.Explains {
		fmt.Println(i+1, item.Pos, item.Trans)
	}
	wg.Done()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprint(os.Stderr, `usage: simpleDict WORD example: simpleDict hello`, os.Args[0])
		os.Exit(1)
	}
	wg.Add(2)
	defer timeMeasurement(time.Now(), "main")
	go caiyunfanyi(os.Args[1])
	go huoshanfanyi(os.Args[1:])
	wg.Wait()
}

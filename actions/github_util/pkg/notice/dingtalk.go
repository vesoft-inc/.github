package notice

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"text/template"
	"time"
)

type DingDing struct {
	secret string
	token  string
}

func NewDingdinng(secret, token string) *DingDing {
	return &DingDing{secret: secret, token: token}
}



type noticeParams struct {
	JobName      string
	Status       string
	SuiteName    string
	ProductName  string
	Commit       string
	CurrentScore string
	RecentScore  []string
	Domain       string
	ID           int
}

const title = "NMeter Result"
const templateString = "## {{ .JobName }} \n\n" +
	"* Suite: {{ .SuiteName }}" + "\n" +
	"* Product: {{ .ProductName }}" + "\n" +
	"* Status: {{ .Status }}" + "\n" +
	"* Commit: {{ .Commit }}" + "\n" +
	"* Current: {{ .CurrentScore }}" + "\n" +
	"* Recent: " + "\n" +
	"{{ range $index, $element := .RecentScore }}  - {{ $element }} \n {{ end }}"

func (d *DingDing) SendDingTalk(title, content string) error {
	ts := time.Now().UnixNano() / 1e6
	h := hmac.New(sha256.New, []byte(d.secret))
	h.Write([]byte(fmt.Sprintf("%d\n%s", ts, d.secret)))
	sign := url.QueryEscape(base64.StdEncoding.EncodeToString(h.Sum(nil)))
	url := fmt.Sprintf("https://oapi.dingtalk.com/robot/send?access_token=%s&timestamp=%d&sign=%s", d.token, ts, sign)
	client := http.Client{
		Timeout: 10 * time.Second,
	}
	body := map[string]interface{}{
		"msgtype": "markdown",
		"markdown": map[string]string{
			"text":  content,
			"title": title,
		},
	}
	bs, err := json.Marshal(body)
	if err != nil {
		return err
	}

	resp, err := client.Post(url, "application/json; charset=utf-8", bytes.NewBuffer(bs))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		return fmt.Errorf("status code: %d", resp.StatusCode)
	}
	return nil
}

func getContent(params *noticeParams) string {
	buf := bytes.NewBufferString("")
	tpl, err := template.New("tpl").Parse(templateString)
	if err != nil {
		return ""
	}
	if err := tpl.Execute(buf, params); err != nil {
		return ""
	}
	return buf.String()
}

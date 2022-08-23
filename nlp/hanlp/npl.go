package hanlp

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"strconv"
	"strings"

	"github.com/samber/lo"
	"github.com/tulanz/base/nlp"
	"github.com/tulanz/base/utils/simple"
)

type HanlpAI struct {
	// client *tnlp.Client
	Token string
}

// "ed2e3bfb7c1d42c6ba773304df7690c21661222719552token"
func NewHanlpAI(token string) (nlp.Summary, error) {
	return &HanlpAI{token}, nil
}

func (t *HanlpAI) Default(title, text string, length uint64) (string, error) {

	data := &bytes.Buffer{}
	w := multipart.NewWriter(data)
	w.WriteField("size", strconv.FormatUint(length, 10))
	w.WriteField("text", text)

	if err := w.Close(); err != nil {
		return "", err
	}

	req, _ := http.NewRequest(http.MethodPost, "http://comdo.hanlp.com/hanlp/v1/summary/extract", data)
	req.Header.Set("Content-Type", w.FormDataContentType())
	req.Header.Set("token", t.Token)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	var rex Response
	_ = json.Unmarshal(body, &rex)
	str := lo.Map(rex.Data, func(t Line, _ int) string { return t.Word })

	textx := strings.Join(str, ",")
	return simple.Substr(textx, 0, 200), nil
}

type Response struct {
	Code int    `json:"code"`
	Data []Line `json:"data"`
}

type Line struct {
	Nature string `json:"nature"`
	Word   string `json:"word"`
}

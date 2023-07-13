package router_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"majiang/router"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"
)

func TestLogin(t *testing.T) {
	go router.InitRouters()
	cases := []struct {
		name     string
		username string
		password string
		want     string
	}{
		{"test1", time.Now().Format("20060102150405"), "password", "操作成功"},
		{"test2", "1", "password", "用户名已存在"},
	}
	for i := range cases {
		//marshal, _ := jsoniter.Marshal(struct {
		//	Username string `json:"username"`
		//	Password string `json:"password"`
		//}{Username: cases[i].username, Password: cases[i].password})
		//request := httptest.NewRequest(http.MethodPost, "/user/register", strings.NewReader(string(marshal)))
		//engine := router.InitRouters("test")
		//w := httptest.NewRecorder()
		//engine.ServeHTTP(w, request)
		//time.Sleep(time.Second)

		b, err := Poster{
			Url: "http://127.0.0.1:8080/user/register",
			Fields: map[string]string{
				"username": cases[i].username,
				"password": cases[i].password,
			},
		}.Post()
		if err != nil {
			t.Errorf("err: %v", err)
		}
		var resp map[string]any

		err = json.Unmarshal(b, &resp)
		if err != nil {
			t.Errorf("err: %v", err)
		}
		v, ok := resp["msg"].(string)
		if !ok {
			t.Errorf("field msg is not a string , v= %v", resp["msg"])
			return
		}
		if v != cases[i].want {
			t.Errorf("got: %v, want: %v", v, cases[i].want)

			return
		}
	}
}

type Poster struct {
	Url    string
	Fields map[string]string
	Files  map[string]string
	Header map[string]string
}

type Getter struct {
	Url    string
	Query  map[string]string
	Header map[string]string
}

func (g Getter) Get() ([]byte, error) {
	var urlSlice []string
	for k, v := range g.Query {
		urlSlice = append(urlSlice, k+"="+v)
	}
	var url = g.Url
	if len(urlSlice) != 0 {
		url = strings.Join(urlSlice, "&")
		url = g.Url + "?" + url
	}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	for k, v := range g.Header {
		req.Header.Add(k, v)
	}
	var client = &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.New("创建请求失败")
	}
	defer resp.Body.Close()
	d, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New("读取数据失败")
	}
	//fmt.Println(string(d))
	dst := map[string]interface{}{}
	if err := json.Unmarshal(d, &dst); err != nil {
		return nil, err
	}
	res, ok := dst["status"].(float64)
	if !ok {
		return nil, err
	}
	if res != 200 {
		return nil, errors.New(fmt.Sprintf("status=%.0f", res))
	}
	return d, nil
}

func (r Poster) Post() ([]byte, error) {
	var buff bytes.Buffer

	writer := multipart.NewWriter(&buff)

	for k, v := range r.Fields {
		writer.WriteField(k, v)
	}
	for k, v := range r.Files {
		w, err := writer.CreateFormFile(k, v)
		file, err := os.Open(v)
		if err != nil {
			fmt.Println(err)
			return nil, errors.New("open file")
		}
		io.Copy(w, file)

	}
	writer.Close()
	req, err := http.NewRequest(http.MethodPost, r.Url, &buff)
	if err != nil {
		return nil, errors.New("create request " + err.Error())
	}
	for k, v := range r.Header {
		req.Header.Set(k, v)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	var client = &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	d, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("读取失败")
		return nil, errors.New("")
	}
	return d, nil
}

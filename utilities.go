package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"syscall"

	"golang.org/x/term"
)

func my_obj_to_json(v interface{}) string {
	A, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	var buf bytes.Buffer
	err2 := json.Indent(&buf, []byte(A), "", "  ")
	if err2 != nil {
		panic(err2)
	}
	return buf.String()
}

func file_exists(fname string) bool {
	_, err := os.Stat(fname)
	return err == nil
}

var executable_path string

func global_init() {
	temp1, err := os.Executable()
	if err != nil {
		panic(err)
	}

	executable_path = filepath.Dir(temp1)
}

func input_text() string {
	var str string
	fmt.Scan(&str)
	return str
}

func input_password() string {
	var str string
	bytePassword, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		panic(err)
	}

	str = string(bytePassword)
	return str
}

func http_call(Origin_URL, URL string, data map[string]string, client *http.Client, call_type string, cookie []*http.Cookie) (string, *http.Response) {

	values := url.Values{}
	for k, e := range data {
		values.Set(k, e)
	}

	req, _ := http.NewRequest(
		call_type,
		URL,
		strings.NewReader(values.Encode()),
	)
	for _, e := range cookie {
		req.AddCookie(e)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Origin", Origin_URL)

	resp, _ := client.Do(req)
	defer resp.Body.Close()

	temp, _ := ioutil.ReadAll(resp.Body)

	return string(temp), resp
}

func http_post(Origin_URL, URL string, data map[string]string, client *http.Client) (string, *http.Response) {
	return http_call(Origin_URL, URL, data, client, "POST", []*http.Cookie{&http.Cookie{}})
}

package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/fatih/color"
)

type Narou_secrets struct {
	Id    string `json:"id"`
	Userl string `json:"userl"`
	Ses   string `json:"ses"`
}

func Project_check_of_narou(episodes []episode_t, summary string, title string) []error {
	A := []error{}
	if title == "" {
		A = append(A, errors.New("タイトルを空にできません"))
	}
	str_temp_summary := strings.Replace(summary, "\n", "", -1)
	str_temp_summary = strings.Replace(str_temp_summary, " ", "", -1)
	str_len_summary := len([]rune(str_temp_summary))
	if str_len_summary < 10 {
		A = append(A, errors.New("あらすじが10文字未満です。あらすじは10文字以上1000文字以下でなければなりません。"))
	}
	if str_len_summary > 1000 {
		A = append(A, errors.New("あらすじが1000文字以上です。あらすじは10文字以上1000文字以下でなければなりません。"))
	}
	for _, e := range episodes {
		str_temp := strings.Replace(e.Body, "\n", "", -1)
		str_temp = strings.Replace(str_temp, " ", "", -1)
		str_len := len([]rune(str_temp))
		if 200 > str_len {
			A = append(A, errors.New(fmt.Sprintf("エピソード「%s」の本文が200文字以下です。エピソードの本文は200文字より大きく70000文字より小さくなければなりません。", e.Meta.Title)))
		}
		if str_len > 70000 {
			A = append(A, errors.New(fmt.Sprintf("エピソード「%s」の本文が70000文字以上です。エピソードの本文は200文字より大きく70000文字より小さくなければなりません。", e.Meta.Title)))
		}
	}

	return A
}

func login_narou(id string, password string) error {
	jar, err := cookiejar.New(nil)
	if err != nil {
		log.Fatalln(err)
	}

	client := &http.Client{Jar: jar}
	client.Get("https://ssl.syosetu.com/login/login/")

	_, req := http_post("https://ssl.syosetu.com/", "https://ssl.syosetu.com/login/login/", map[string]string{
		"narouid": id,
		"pass":    password,
	}, client)

	C := req.Cookies()

	Narou_secrets_data := Narou_secrets{}
	Narou_secrets_data.Id = id

	fmt.Println("\ne.Name")
	success := false
	for _, e := range C {
		if e.Name == "userl" {
			Narou_secrets_data.Userl = e.Value
		}
		if e.Name == "ses" {
			Narou_secrets_data.Ses = e.Value
			success = true
		}
	}
	if !success {
		return errors.New("Inputted Email or password is wrong.")
	}

	S_bytes, err := json.Marshal(Narou_secrets_data)
	if err != nil {
		panic(err)
	}

	f, err := os.Create(filepath.Join(executable_path, "secrets", "narou.json"))
	if err != nil {
		panic(err)
	}
	f.Write(S_bytes)
	f.Close()

	return nil
}

func load_narou_secret() Narou_secrets {
	Narou_secrets_data := Narou_secrets{}

	bytes_data, err := ioutil.ReadFile(filepath.Join(executable_path, "secrets", "narou.json"))
	if err != nil {
		log.Fatalln("You need to login Narou.\n\nnwp narou l")
	}
	json.Unmarshal(bytes_data, &Narou_secrets_data)

	return Narou_secrets_data
}

type Get_list_of_narou_one_result_t struct {
	Ncode      string `json:"ncode"`
	Name       string `json:"name"`
	Writing    bool   `json:"writing"`
	LastUpdate string `json:"last_update"`
	EditLink   string `json:"edit_link"`
}

func Get_list_of_narou(Narou_secrets_data *Narou_secrets) []Get_list_of_narou_one_result_t {
	jar, err := cookiejar.New(nil)
	if err != nil {
		log.Fatalln(err)
	}

	client := &http.Client{Jar: jar}
	//client.Get("https://syosetu.com/usernovel/list/")

	temp, _ := http_call("https://syosetu.com", "https://syosetu.com/usernovel/list/", map[string]string{}, client, "GET",
		[]*http.Cookie{&http.Cookie{
			Name:     "userl",
			Value:    Narou_secrets_data.Userl,
			HttpOnly: true,
		},
		})

	bufferReader := bytes.NewReader([]byte(temp))

	document, _ := goquery.NewDocumentFromReader(bufferReader)

	A := []Get_list_of_narou_one_result_t{}

	novel_list := document.Find("#novellist")
	novel_list.Children().Children().Each(func(i int, s *goquery.Selection) {
		if i == 0 {
			return
		}
		temp := Get_list_of_narou_one_result_t{}
		s.Children().Each(func(j int, s2 *goquery.Selection) {
			if j == 0 {
				if "連載中" == strings.Replace(s2.Text(), "\n", "", -1) {
					temp.Writing = true
				} else {
					temp.Writing = false
				}
			}
			if j == 1 {
				temp.Ncode = s2.Text()
			}
			if j == 2 {
				temp.Name = strings.Replace(s2.Text(), "\n", "", -1)
			}
			if j == 3 {
				temp.LastUpdate = s2.Text()
			}
			if j == 4 {
				path, _ := s2.Children().Attr("href")
				temp.EditLink = "https://syosetu.com" + path
			}
		})
		A = append(A, temp)
	})

	//fmt.Printf("%v\n", A)
	//fmt.Printf("%s\n", temp)

	return A
}

func Get_episode_body_of_narou(URL string, Narou_secrets_data *Narou_secrets) string {
	jar, err := cookiejar.New(nil)
	if err != nil {
		log.Fatalln(err)
	}

	client := &http.Client{Jar: jar}
	temp, _ := http_call("https://syosetu.com", URL, map[string]string{}, client, "GET",
		[]*http.Cookie{&http.Cookie{
			Name:     "userl",
			Value:    Narou_secrets_data.Userl,
			HttpOnly: true,
		},
		})

	bufferReader := bytes.NewReader([]byte(temp))

	document, _ := goquery.NewDocumentFromReader(bufferReader)
	A := document.Find("#novel").Text()

	return A
}

type Get_list_of_episode_of_narou_one_result_t struct {
	Data      episode_t `json:"data"`
	Edit_link string    `json:"edit_link"`
}

func Get_list_of_episode_of_narou(edit_id string, Narou_secrets_data *Narou_secrets) []Get_list_of_episode_of_narou_one_result_t {
	jar, err := cookiejar.New(nil)
	if err != nil {
		log.Fatalln(err)
	}

	client := &http.Client{Jar: jar}

	temp, _ := http_call("https://syosetu.com", fmt.Sprintf("https://syosetu.com/usernovelmanage/top/ncode/%s/", edit_id), map[string]string{}, client, "GET",
		[]*http.Cookie{&http.Cookie{
			Name:     "userl",
			Value:    Narou_secrets_data.Userl,
			HttpOnly: true,
		},
		})

	bufferReader := bytes.NewReader([]byte(temp))

	document, _ := goquery.NewDocumentFromReader(bufferReader)

	A := []Get_list_of_episode_of_narou_one_result_t{}

	episode_list := document.Find("#novelsublist").Find("table")
	episode_list.Children().Children().Each(func(i int, s *goquery.Selection) {
		if i == 0 {
			return
		}
		temp := Get_list_of_episode_of_narou_one_result_t{}
		temp.Data = episode_t{}
		temp.Data.Meta = episode_meta_t{}
		temp.Data.Meta.Index = i + 1
		s.Children().Each(func(j int, s2 *goquery.Selection) {
			if j == 1 {
				temp.Data.Meta.Title = s2.Text()
			}
			if j == 3 {
				href, Is_exist := s2.Children().Attr("href")
				if Is_exist {
				}
				URL := "https://syosetu.com" + href
				body := Get_episode_body_of_narou(URL, Narou_secrets_data)
				temp.Data.Body = body
				temp.Edit_link = URL
			}
		})
		A = append(A, temp)
	})

	//fmt.Printf("%v\n", A)
	//fmt.Printf("%s\n", temp)

	return A
}

func Get_edit_id_from_link_of_narou(URL string) string {
	r := regexp.MustCompile(`https://syosetu.com/usernovelmanage/updateinput/ncode/([0-9]+)/`)
	R := r.FindStringSubmatch(URL)
	return R[1]
}

func Get_edit_episode_id_from_edit_link_of_narou(URL string) string {
	r := regexp.MustCompile(`https://syosetu.com/usernoveldatamanage/updateinput/ncode/([0-9]+)/noveldataid/([0-9]+)/`)
	R := r.FindStringSubmatch(URL)
	return R[2]
}

func make_new_novel_or_narou(episodes []episode_t, Project_Setting_data Project_Setting, Narou_secrets_data *Narou_secrets, summary string) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		log.Fatalln(err)
	}

	client := &http.Client{Jar: jar}
	temp, _ := http_call("https://syosetu.com", "https://syosetu.com/userwrittingnovel/input/", map[string]string{}, client, "GET",
		[]*http.Cookie{&http.Cookie{
			Name:     "userl",
			Value:    Narou_secrets_data.Userl,
			HttpOnly: true,
		},
		})

	//fmt.Println(temp)
	bufferReader := bytes.NewReader([]byte(temp))

	document, _ := goquery.NewDocumentFromReader(bufferReader)

	to_Path, _ := document.Find("form").Attr("action")

	to_URL := "https://syosetu.com" + to_Path

	FORM_data := map[string]string{}

	FORM_data["karititle"] = Project_Setting_data.Title
	FORM_data["novel"] = episodes[0].Body
	//fmt.Printf("%v\n", FORM_data)

	jar, err = cookiejar.New(nil)
	if err != nil {
		log.Fatalln(err)
	}

	client = &http.Client{Jar: jar}
	temp2, _ := http_call("https://syosetu.com", to_URL, FORM_data, client, "POST",
		[]*http.Cookie{&http.Cookie{
			Name:     "userl",
			Value:    Narou_secrets_data.Userl,
			HttpOnly: true,
		},
		})
	//fmt.Println(temp2)

	bufferReader2 := bytes.NewReader([]byte(temp2))

	document2, _ := goquery.NewDocumentFromReader(bufferReader2)

	//fmt.Println(req.Request.URL.String())

	publish_button_URL := ""

	document2.Find(".button_box").Children().Each(func(i int, s *goquery.Selection) {
		//fmt.Println(s.Text())
		if s.Text() == "投稿する" {
			//fmt.Printf("OKKKKKKKKKKKKKKい")
			temp, _ := s.Attr("href")
			publish_button_URL = "https://syosetu.com" + temp
		}
	})

	//fmt.Println(publish_button_URL)
	temp3, _ := http_call("https://syosetu.com", publish_button_URL, map[string]string{}, client, "GET",
		[]*http.Cookie{&http.Cookie{
			Name:     "userl",
			Value:    Narou_secrets_data.Userl,
			HttpOnly: true,
		},
		})
	//fmt.Println(temp3)
	bufferReader3 := bytes.NewReader([]byte(temp3))

	document3, _ := goquery.NewDocumentFromReader(bufferReader3)
	to_Path3, _ := document3.Find("form").Attr("action")

	to_URL3 := "https://syosetu.com" + to_Path3
	//fmt.Println(to_URL3)

	FORM_data3 := map[string]string{}

	FORM_data3["age_limit"] = "1"
	FORM_data3["noveltype"] = "1"
	FORM_data3["genre"] = "9999"
	FORM_data3["ex"] = summary
	FORM_data3["title"] = Project_Setting_data.Title
	FORM_data3["subtitle"] = episodes[0].Meta.Title
	//fmt.Printf("%v\n", FORM_data3)

	temp4, _ := http_call("https://syosetu.com", to_URL3, FORM_data3, client, "POST",
		[]*http.Cookie{&http.Cookie{
			Name:     "userl",
			Value:    Narou_secrets_data.Userl,
			HttpOnly: true,
		},
		})
	//fmt.Println(temp4)

	bufferReader4 := bytes.NewReader([]byte(temp4))

	document4, _ := goquery.NewDocumentFromReader(bufferReader4)
	to_Path4, _ := document4.Find("form").Attr("action")
	to_URL4 := "https://syosetu.com" + to_Path4

	http_call("https://syosetu.com", to_URL4, map[string]string{}, client, "GET",
		[]*http.Cookie{&http.Cookie{
			Name:     "userl",
			Value:    Narou_secrets_data.Userl,
			HttpOnly: true,
		},
		})
	//fmt.Println(temp5)
}

func update_episode_of_narou(to_episode *episode_t, edit_URL string, Project_Setting_data Project_Setting, Narou_secrets_data *Narou_secrets) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		log.Fatalln(err)
	}

	client := &http.Client{Jar: jar}

	FORM_data := map[string]string{
		"subtitle": to_episode.Meta.Title,
		"novel":    to_episode.Body,
	}

	temp, _ := http_call("https://syosetu.com", edit_URL, map[string]string{}, client, "GET",
		[]*http.Cookie{&http.Cookie{
			Name:     "userl",
			Value:    Narou_secrets_data.Userl,
			HttpOnly: true,
		},
		})

	bufferReader := bytes.NewReader([]byte(temp))

	document, _ := goquery.NewDocumentFromReader(bufferReader)
	to_Path, _ := document.Find("form").Attr("action")
	to_URL := "https://syosetu.com" + to_Path

	temp2, _ := http_call("https://syosetu.com", to_URL, FORM_data, client, "POST",
		[]*http.Cookie{&http.Cookie{
			Name:     "userl",
			Value:    Narou_secrets_data.Userl,
			HttpOnly: true,
		},
		})

	bufferReader2 := bytes.NewReader([]byte(temp2))

	document2, _ := goquery.NewDocumentFromReader(bufferReader2)
	to_Path2, _ := document2.Find("form").Attr("action")
	to_URL2 := "https://syosetu.com" + to_Path2

	csrf, _ := document2.Find("form").Find("[name ='csrf_onetimepass']").Attr("value")
	FORM_data["csrf_onetimepass"] = csrf

	http_call("https://syosetu.com", to_URL2, FORM_data, client, "POST",
		[]*http.Cookie{&http.Cookie{
			Name:     "userl",
			Value:    Narou_secrets_data.Userl,
			HttpOnly: true,
		},
		})
}

func add_episode_of_narou(to_episode *episode_t, Edit_id string, Project_Setting_data Project_Setting, Narou_secrets_data *Narou_secrets) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		log.Fatalln(err)
	}

	client := &http.Client{Jar: jar}

	FORM_data := map[string]string{
		"subtitle": to_episode.Meta.Title,
		"novel":    to_episode.Body,
		"end":      "1",
	}

	temp, _ := http_call("https://syosetu.com", fmt.Sprintf("https://syosetu.com/usernovelmanage/ziwainput/ncode/%s/", Edit_id), map[string]string{}, client, "GET",
		[]*http.Cookie{&http.Cookie{
			Name:     "userl",
			Value:    Narou_secrets_data.Userl,
			HttpOnly: true,
		},
		})

	//fmt.Println(temp)

	bufferReader := bytes.NewReader([]byte(temp))

	document, _ := goquery.NewDocumentFromReader(bufferReader)
	to_Path, _ := document.Find("form").Attr("action")
	to_URL := "https://syosetu.com" + to_Path

	temp2, _ := http_call("https://syosetu.com", to_URL, FORM_data, client, "POST",
		[]*http.Cookie{&http.Cookie{
			Name:     "userl",
			Value:    Narou_secrets_data.Userl,
			HttpOnly: true,
		},
		})
	//fmt.Println(temp2)

	bufferReader2 := bytes.NewReader([]byte(temp2))

	document2, _ := goquery.NewDocumentFromReader(bufferReader2)
	to_Path2, _ := document2.Find("form").Attr("action")
	to_URL2 := "https://syosetu.com" + to_Path2

	csrf, _ := document2.Find("form").Find("[name ='csrf_onetimepass']").Attr("value")
	//fmt.Println(csrf)
	FORM_data["csrf_onetimepass"] = csrf

	http_call("https://syosetu.com", to_URL2, FORM_data, client, "POST",
		[]*http.Cookie{&http.Cookie{
			Name:     "userl",
			Value:    Narou_secrets_data.Userl,
			HttpOnly: true,
		},
		})
	//fmt.Println(temp4)

}

func delete_episode_of_narou(Edit_id, edit_episode_id string, Project_Setting_data Project_Setting, Narou_secrets_data *Narou_secrets) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		log.Fatalln(err)
	}

	client := &http.Client{Jar: jar}

	FORM_data := map[string]string{}

	temp, _ := http_call("https://syosetu.com", fmt.Sprintf("https://syosetu.com/usernoveldatamanage/deleteconfirm/ncode/%s/noveldataid/%s/", Edit_id, edit_episode_id), map[string]string{}, client, "GET",
		[]*http.Cookie{&http.Cookie{
			Name:     "userl",
			Value:    Narou_secrets_data.Userl,
			HttpOnly: true,
		},
		})

	//fmt.Println(temp)

	bufferReader := bytes.NewReader([]byte(temp))

	document, _ := goquery.NewDocumentFromReader(bufferReader)
	to_Path, _ := document.Find("form").Attr("action")
	to_URL := "https://syosetu.com" + to_Path

	csrf, _ := document.Find("form").Find("[name ='csrf_onetimepass']").Attr("value")
	FORM_data["csrf_onetimepass"] = csrf

	temp2, _ := http_call("https://syosetu.com", to_URL, FORM_data, client, "POST",
		[]*http.Cookie{&http.Cookie{
			Name:     "userl",
			Value:    Narou_secrets_data.Userl,
			HttpOnly: true,
		},
		})
	//fmt.Println(temp2)

	bufferReader2 := bytes.NewReader([]byte(temp2))

	document2, _ := goquery.NewDocumentFromReader(bufferReader2)
	to_Path2, _ := document2.Find("form").Attr("action")
	to_URL2 := "https://syosetu.com" + to_Path2

	csrf, _ = document.Find("form").Find("[name ='csrf_onetimepass']").Attr("value")
	FORM_data["csrf_onetimepass"] = csrf
	//fmt.Println(csrf)

	http_call("https://syosetu.com", to_URL2, FORM_data, client, "POST",
		[]*http.Cookie{&http.Cookie{
			Name:     "userl",
			Value:    Narou_secrets_data.Userl,
			HttpOnly: true,
		},
		})
	//fmt.Println(temp4)
}

func update_info_of_narou(Edit_id string, Project_Setting_data Project_Setting, Narou_secrets_data *Narou_secrets, summary string) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		log.Fatalln(err)
	}

	client := &http.Client{Jar: jar}

	FORM_data := map[string]string{}

	FORM_data["end"] = "1"
	FORM_data["ex"] = summary
	FORM_data["genre"] = "9999"
	FORM_data["title"] = Project_Setting_data.Title

	temp, _ := http_call("https://syosetu.com", fmt.Sprintf("https://syosetu.com/usernovelmanage/updateinput/ncode/%s/", Edit_id), map[string]string{}, client, "GET",
		[]*http.Cookie{&http.Cookie{
			Name:     "userl",
			Value:    Narou_secrets_data.Userl,
			HttpOnly: true,
		},
		})

	//fmt.Println(temp)

	bufferReader := bytes.NewReader([]byte(temp))

	document, _ := goquery.NewDocumentFromReader(bufferReader)
	to_Path, _ := document.Find("form").Attr("action")
	to_URL := "https://syosetu.com" + to_Path

	csrf, _ := document.Find("form").Find("[name ='csrf_onetimepass']").Attr("value")
	FORM_data["csrf_onetimepass"] = csrf

	temp2, _ := http_call("https://syosetu.com", to_URL, FORM_data, client, "POST",
		[]*http.Cookie{&http.Cookie{
			Name:     "userl",
			Value:    Narou_secrets_data.Userl,
			HttpOnly: true,
		},
		})

	bufferReader2 := bytes.NewReader([]byte(temp2))

	document2, _ := goquery.NewDocumentFromReader(bufferReader2)
	to_Path2, _ := document2.Find("form").Attr("action")
	to_URL2 := "https://syosetu.com" + to_Path2

	csrf, _ = document2.Find("form").Find("[name ='csrf_onetimepass']").Attr("value")
	FORM_data["csrf_onetimepass"] = csrf
	//fmt.Println(csrf)

	http_call("https://syosetu.com", to_URL2, FORM_data, client, "POST",
		[]*http.Cookie{&http.Cookie{
			Name:     "userl",
			Value:    Narou_secrets_data.Userl,
			HttpOnly: true,
		},
		})
}

func deploy_of_Narou(episodes []episode_t, Project_Setting_data Project_Setting, summary string) []error {
	Project_Erros := Project_check_of_narou(episodes, summary, Project_Setting_data.Title)
	if len(Project_Erros) > 0 {
		for _, e := range Project_Erros {
			fmt.Println(color.RedString(fmt.Sprintf("Error:%s", e.Error())))
		}
		return Project_Erros
	}
	Narou_secrets_data := load_narou_secret()
	novel_list := Get_list_of_narou(&Narou_secrets_data)
	Is_exist_novel := false
	for _, e := range novel_list {
		if e.Name == Project_Setting_data.Title {
			Is_exist_novel = true
		}
	}

	if !Is_exist_novel {
		make_new_novel_or_narou(episodes, Project_Setting_data, &Narou_secrets_data, summary)
	}

	novel_list = Get_list_of_narou(&Narou_secrets_data)
	this_novel := Get_list_of_narou_one_result_t{}
	for _, e := range novel_list {
		if e.Name == Project_Setting_data.Title {
			this_novel = e
		}
	}
	edit_id := Get_edit_id_from_link_of_narou(this_novel.EditLink)

	Narou_episodes := Get_list_of_episode_of_narou(edit_id, &Narou_secrets_data)

	update_info_of_narou(edit_id, Project_Setting_data, &Narou_secrets_data, summary)

	for i, e := range episodes {
		if i < len(Narou_episodes) {
			if e.Body == Narou_episodes[i].Data.Body && e.Meta.Title == Narou_episodes[i].Data.Meta.Title {

			} else {
				update_episode_of_narou(&e, Narou_episodes[i].Edit_link, Project_Setting_data, &Narou_secrets_data)
			}
		} else {
			add_episode_of_narou(
				&e, edit_id, Project_Setting_data, &Narou_secrets_data,
			)
		}
	}

	for i, e := range Narou_episodes {
		if i >= len(episodes) {
			edit_episode_id := Get_edit_episode_id_from_edit_link_of_narou(e.Edit_link)
			delete_episode_of_narou(edit_id, edit_episode_id, Project_Setting_data, &Narou_secrets_data)
		} else {
		}
	}

	return []error{}
}

package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"sort"

	"gopkg.in/yaml.v2"
)

type Project_Setting struct {
	Title   string   `yaml:"title"`
	Deploys []string `yaml:"deploys"`
}

func new_project(Title string) {
	Project_Setting_data := Project_Setting{}
	nwp_f, err := os.Create("nwp.yml")
	if err != nil {
		panic(err)
	}

	Project_Setting_data.Title = Title
	Project_Setting_data.Deploys = []string{"narou"}

	Project_Setting_byte_data, err := yaml.Marshal(Project_Setting_data)
	nwp_f.Write(Project_Setting_byte_data)
	if err != nil {
		panic(err)
	}

	_, err2 := os.Create("_summary.txt")
	if err2 != nil {
		panic(err2)
	}
}

func load_project() (Project_Setting, string) {
	Project_Setting_data := Project_Setting{}
	bytes_data, err := ioutil.ReadFile("nwp.yml")
	if err != nil {
		panic(err)
	}

	yaml.Unmarshal(bytes_data, &Project_Setting_data)

	bytes_data2, err := ioutil.ReadFile("_summary.txt")
	if err != nil {
		panic(err)
	}
	return Project_Setting_data, string(bytes_data2)
}

type episode_meta_t struct {
	Title string `yaml:"title"`
	Index int    `yaml:"index"`
}

type episode_t struct {
	meta_yaml string
	Meta      episode_meta_t
	Body      string
}

func add_episode_in_project(episodes []episode_t, sub_title string) {
	fname := sub_title + ".txt"
	if file_exists(fname) {
		fmt.Println("episode is already exist.")
		return
	}
	f, err := os.Create(fname)
	if err != nil {
		panic(err)
	}

	max_index := Get_max_episodes_index(episodes)

	episode_data := episode_t{}
	episode_data.Meta = episode_meta_t{
		Title: sub_title,
		Index: max_index + 1,
	}

	temp, err := yaml.Marshal(episode_data.Meta)
	episode_data.meta_yaml = string(temp)
	if err != nil {
		panic(err)
	}

	f.Write([]byte(fmt.Sprintf("---\n%s\n---\n%s", episode_data.meta_yaml, "")))

	f.Close()
}

func load_episode_in_project(sub_title string) episode_t {
	episode_data := episode_t{}
	episode_data.Meta = episode_meta_t{}
	fname := sub_title

	bytes_data, err := ioutil.ReadFile(fname)
	if err != nil {
		panic(err)
	}
	//str_data := string(bytes_data)
	episode_data.meta_yaml = string(bytes_data)

	r := regexp.MustCompile(`^\-\-\-[\s\S]*\-\-\-`)

	index := r.FindIndex(bytes_data)

	yaml.Unmarshal(bytes_data[index[0]+3:index[1]-3], &episode_data.Meta)

	episode_data.Body = string(bytes_data[index[1]:])

	return episode_data
}

func load_episodes_in_project() []episode_t {
	files, err := ioutil.ReadDir(".")
	if err != nil {
		panic(err)
	}
	r1 := regexp.MustCompile(`\_.*`)
	r2 := regexp.MustCompile(`.*.txt`)

	A := []episode_t{}

	for _, file := range files {
		if file.IsDir() {
			continue
		}
		if r1.MatchString(file.Name()) {
			continue
		}
		if !r2.MatchString(file.Name()) {
			continue
		}
		A = append(A, load_episode_in_project(file.Name()))
	}

	sort.Slice(A[:], func(i, j int) bool {
		return A[i].Meta.Index < A[j].Meta.Index
	})

	return A
}

func Get_max_episodes_index(v []episode_t) int {
	max_index := 0
	for _, e := range v {
		if max_index < e.Meta.Index {
			max_index = e.Meta.Index
		}
	}
	return max_index
}

package main

import (
	"log"
	"time"
	"regexp"
	"strings"
	"io/ioutil"
	"net/http"
	"crypto/tls"
	"encoding/json"
)

var transport = &http.Transport{
	TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
}

var client = &http.Client{
	Transport: transport,
	Timeout: 5*time.Second,
}

func backend_github(owner, repo, pattern string) {
	var data []interface{}
	var item map[string]interface{}
	var version string
	var date string
	var value int64

	// Retrieve list of tags
	url := "https://api.github.com/repos/" + owner + "/" + repo + "/tags"
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Add("Accept", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		log.Print(err)
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	s := []byte(body)
	if resp.StatusCode != 200 {
		return
	}

	err = json.Unmarshal(s, &data)
	if err != nil {
		log.Print(err)
		return
	}

	// Compile pattern
	re := regexp.MustCompile(pattern)

	// Find latest release
	for i := 0; i < len(data); i++ {
		item = data[i].(map[string]interface{})
		cver := item["name"].(string)
		item = item["commit"].(map[string]interface{})

		match := re.FindStringSubmatch(cver)
		if len(match) > 0 {
			cver = match[1]
		} else {
			continue
		}

		if version_compare(cver, version) > 0 {
			version = cver
			url = item["url"].(string)
		}
	}

	// No match found
	if len(version) == 0 {
		return
	}

	// Retrieve commit
	req, err = http.NewRequest("GET", url, nil)
	req.Header.Add("Accept", "application/json")
	resp, err = client.Do(req)
	if err != nil {
		log.Print(err)
		return
	}

	body, err = ioutil.ReadAll(resp.Body)
	s = []byte(body)
	if resp.StatusCode != 200 {
		return
	}

	err = json.Unmarshal(s, &item)
	if err != nil {
		log.Print(err)
		return
	}

	// Parse release date
	item = item["commit"].(map[string]interface{})
	item = item["author"].(map[string]interface{})
	date = item["date"].(string)

	ts, err := time.Parse(time.RFC3339, date)
	if err != nil {
		log.Print(err)
	}

	// Generate metrics
	value = ts.Unix()
	url = "https://github.com/" + owner + "/" + repo
	metricsAppend(repo, version, url, value)
}

func backend_folder(name, info_url, url, pattern string) {
	var lines []string
	var version string
	var value int64

	// Retrive list of files
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Add("Accept", "text/html")
	resp, err := client.Do(req)
	if err != nil {
		log.Print(err)
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		return
	}

	// Compile regexp patterns
	d1 := regexp.MustCompile(`(\d{4}-\d{2}-\d{2} \d{2}:\d{2})`) // 2006-01-02 15:04
	d2 := regexp.MustCompile(`(\d{2}-\w{3}-\d{4} \d{2}:\d{2})`) // 02-Jan-2006 15:04
	re := regexp.MustCompile(pattern)

	value = 0
	lines = strings.Split(string(body), "\n")

	// Loop through all lines
	for _,line := range lines {
		match := re.FindStringSubmatch(line)
		if len(match) == 0 {
			continue
		}

		cver := match[1]
		cval := int64(0)

		match = d1.FindStringSubmatch(line)
		if len(match) > 0 {
			t,_ := time.Parse("2006-01-02 15:04", match[1])
			cval = t.Unix()
		}

		match = d2.FindStringSubmatch(line)
		if len(match) > 0 {
			t,_ := time.Parse("02-Jan-2006 15:04", match[1])
			cval = t.Unix()
		}

		if version_compare(cver, version) > 0 {
			version = cver
			value = cval
		}
	}

	// No match found
	if len(version) == 0 {
		return
	}

	// Generate metrics
	metricsAppend(name, version, info_url, value)
}

func backend_regexp(name, info_url, url, vre, dre, dfmt string) {
	var version string
	var value int64

	// Retrive file
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Add("Accept", "text/html")
	resp, err := client.Do(req)
	if err != nil {
		log.Print(err)
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		return
	}

	// Find version and date
	str := string(body)
	d1 := regexp.MustCompile(vre)
	d2 := regexp.MustCompile(dre)

	match := d1.FindStringSubmatch(str)
	if len(match) == 0 {
		return
	}
	version = match[1]

	match = d2.FindStringSubmatch(str)
	if len(match) == 0 {
		return
	}

	t, err := time.Parse(dfmt, match[1])
	if err != nil {
		log.Print(err)
		return
	}

	// Generate metrics
	value = t.Unix()
	metricsAppend(name, version, info_url, value)
}

func collectMetrics() {
	// Github
	for _,v := range config.Github {
		pattern := v.Regexp
		if len(pattern) == 0 {
			pattern = "(.*)"
		}
		backend_github(v.Owner, v.Repo, pattern)
	}

	// Folder
	for _,v := range config.Folder {
		backend_folder(v.Name, v.Info, v.Path, v.Regexp)
	}

	// Regexp
	for _,v := range config.Regexp {
		backend_regexp(v.Name, v.Info, v.Path, v.Regexp, v.Date, v.Format)
	}
}

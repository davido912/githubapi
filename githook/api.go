package githook

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/fatih/color"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"strconv"
)

func RaiseError(err error) {
	if err != nil {
		log.Printf("Error encountered: %s", err)
		os.Exit(1)
	}
}

type GitHook struct {
	Endpoint string
	Method   string
}

// marshal the json data first
func (gh *GitHook) auth(jsonData []byte) *http.Request {
	req, err := http.NewRequest(gh.Method, gh.Endpoint, bytes.NewReader(jsonData))
	req.Header.Set("Authorization", "token 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")
	RaiseError(err)
	return req
}

func (gh *GitHook) apiRequest(issue *Issue, params map[string]string) int {
	jsonData, _ := json.Marshal(issue)
	status := 0
	req := gh.auth(jsonData)
	client := &http.Client{}
	if params != nil {
		q := req.URL.Query()

		for k, v := range params {
			q.Add(k, v)
		}
		req.URL.RawQuery = q.Encode()
	}
	resp, err := client.Do(req)
	RaiseError(err)
	defer resp.Body.Close()
	bs, err := io.ReadAll(resp.Body)
	RaiseError(err)

	if !bytes.Equal(bs, []byte{91, 93}) { // empty response
		var sIssues []Issue // turned into a slice of issues to load, pretty cool
		x := bytes.TrimLeft(bs, " \t\r\n")
		if string(x[0]) == "[" {
			json.Unmarshal(bs, &sIssues)
		} else if string(x[0]) == "{" {
			json.Unmarshal(bs, issue)
			sIssues = append(sIssues, *issue)
		}
		logResults(sIssues)
	} else {
		status = 1
	}
	return status
}

func (gh *GitHook) PostIssue(issue Issue) {
	gh.Method = http.MethodPost
	gh.apiRequest(&issue, nil)
}

func (gh *GitHook) ModifyIssue(issue *Issue) {
	gh.Endpoint = UrlGit
	gh.Method = http.MethodPatch
	u, _ := url.Parse(gh.Endpoint)
	n := strconv.Itoa(issue.Number)
	u.Path = path.Join(u.Path, n)
	gh.Endpoint = u.String()
	gh.apiRequest(issue, nil)
}

func (gh *GitHook) CloseIssue(issue *Issue) {
	gh.GetIssue(issue)
	issue.State = "closed"
	gh.ModifyIssue(issue)	
}

func (gh *GitHook) GetIssue(issue *Issue) {
	gh.Method = http.MethodGet
	u, _ := url.Parse(gh.Endpoint)
	n := strconv.Itoa(issue.Number)
	u.Path = path.Join(u.Path, n)
	gh.Endpoint = u.String()
	gh.apiRequest(issue, nil)
}

func (gh *GitHook) GetIssues() {
	gh.Method = http.MethodGet
	gh.Endpoint = UrlGit
	p := 1
	for {
		params := map[string]string{"page": fmt.Sprintf("%d", p)}
		status := gh.apiRequest(nil, params)
		if status == 1 {
			return
		}
		p++
	}
}

func logResults(result []Issue) {
	// s := "#####\nIssue nr %d\nTitle: %s\nBody: %s\nUrl: %s\nState: %s\n#####\n"
	red := color.New(color.FgRed).SprintFunc()
	s := "Issue nr %s, Title: %s, Url: %s, State: %s\n"
	for _, r := range result {
		fmt.Printf(s, red(r.Number), red(r.Title), red(r.Url), red(r.State))
	}
}

type Issue struct {
	Title  string `json:"title"`
	Body   string `json:"body"`
	Number int    `json:"number"`
	Url    string `json:"url"`
	State  string `json:"state"`
}


var UrlGit string = fmt.Sprintf("https://api.github.com/repos/%s/issues", os.Getenv("GIT_REPO"))

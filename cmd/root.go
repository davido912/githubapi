package commands

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/davido912/githubapi/githook"
	"io/ioutil"
	"os"
	"os/exec"
)

const (
	issueIdFlag    string = "issue-id"
	issueTitleFlag string = "title"
	issueBodyFlag  string = "body"
	allFlag        string = "all"
)

var gh githook.GitHook = githook.GitHook{Endpoint: githook.UrlGit}

var rootCmd = &cobra.Command{
	Use:   "gitapi",
	Short: "Git API is a CLI built for a project done in Go (rookie learning stuff) ;)",
}

var getIssue = &cobra.Command{
	Use:   "get",
	Short: "Get Github Issue details",
	Run: func(cmd *cobra.Command, args []string) {
		issue := &githook.Issue{}

		if cmd.Flag(issueIdFlag).Changed {
			getIssueIdFlag(cmd, issue)
			gh.GetIssue(issue)
		}
		if cmd.Flag(allFlag).Changed {
			gh.GetIssues()
		}
	},
}

var modifyIssue = &cobra.Command{
	Use:   "modify",
	Short: "Modify Github Issue",
	Run: func(cmd *cobra.Command, args []string) {
		issue := &githook.Issue{}
		getIssueIdFlag(cmd, issue)
		getIssueTitleFlag(cmd, issue)
		getIssueBodyFlag(cmd, issue)
		gh.ModifyIssue(issue)
	},
}

var postIssue = &cobra.Command{
	Use:   "post",
	Short: "Post Github Issue",
	Run: func(cmd *cobra.Command, args []string) {
		issue := githook.Issue{}
		getIssueTitleFlag(cmd, &issue)
		getIssueBodyFlag(cmd, &issue)
		gh.PostIssue(issue)
	},
}

var closeIssue = &cobra.Command{
	Use:   "close",
	Short: "Close Github Issue",
	Run: func(cmd *cobra.Command, args []string) {
		issue := githook.Issue{}
		getIssueIdFlag(cmd, &issue)
		gh.CloseIssue(&issue)
	},
}

func getIssueTitleFlag(cmd *cobra.Command, issue *githook.Issue) {
	issue_title, err := cmd.Flags().GetString(issueTitleFlag)
	if cmd.Flag(issueTitleFlag).Changed {
		githook.RaiseError(err)
		issue.Title = issue_title
	} else {
		gh.GetIssue(issue)
		gh.Endpoint = githook.UrlGit
	}
}

// using vim for entering body
func getIssueBodyFlag(cmd *cobra.Command, issue *githook.Issue) {
	if cmd.Flag(issueBodyFlag).Changed {
		vimPath, err := exec.LookPath("vim")
		githook.RaiseError(err)
		f, err := ioutil.TempFile("", "git")
		githook.RaiseError(err)
		defer os.Remove(f.Name())
		execCmd := exec.Cmd{Path: vimPath, Args: []string{"vi", f.Name()}, Stdin: os.Stdin, Stdout: os.Stdout, Stderr: os.Stderr}
		execCmd.Run()
		fmt.Println(f.Name())
		t, _ := os.Open(f.Name())
		bs, _ := ioutil.ReadAll(t)
		issue.Body = string(bs)

	}
}

func getIssueIdFlag(cmd *cobra.Command, issue *githook.Issue) {
	issue_id, err := cmd.Flags().GetInt(issueIdFlag)
	if cmd.Flag(issueIdFlag).Changed {
		githook.RaiseError(err)
		issue.Number = issue_id
	}
}

func init() {
	rootCmd.AddCommand(getIssue)
	rootCmd.AddCommand(modifyIssue)
	rootCmd.AddCommand(postIssue)
	rootCmd.AddCommand(closeIssue)

	postIssue.Flags().String(issueTitleFlag, "", "Issue title in GITHUB")
	postIssue.Flags().Bool(issueBodyFlag, false, "Issue body in GITHUB")
	postIssue.MarkFlagRequired(issueTitleFlag)

	modifyIssue.Flags().String(issueTitleFlag, "", "Issue title in GITHUB")
	modifyIssue.Flags().Bool(issueBodyFlag, false, "Issue body in GITHUB")
	modifyIssue.Flags().Int(issueIdFlag, 0, "Issue ID in GITHUB")
	modifyIssue.MarkFlagRequired(issueIdFlag)

	getIssue.Flags().Bool(allFlag, false, "Get all issues in designated repo")
	getIssue.Flags().Int(issueIdFlag, 0, "Issue ID in GITHUB")

	closeIssue.Flags().Int(issueIdFlag, 0, "Issue ID in GITHUB")
	closeIssue.MarkFlagRequired(issueIdFlag)

}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

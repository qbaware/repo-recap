package github

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

// Repository represents a GitHub repository
type Repository struct {
	ID              int       `json:"id"`
	Name            string    `json:"name"`
	FullName        string    `json:"full_name"`
	Description     string    `json:"description"`
	HTMLURL         string    `json:"html_url"`
	CloneURL        string    `json:"clone_url"`
	Language        string    `json:"language"`
	StargazersCount int       `json:"stargazers_count"`
	ForksCount      int       `json:"forks_count"`
	OpenIssuesCount int       `json:"open_issues_count"`
	DefaultBranch   string    `json:"default_branch"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	PushedAt        time.Time `json:"pushed_at"`
	Size            int       `json:"size"`
	Private         bool      `json:"private"`
	Fork            bool      `json:"fork"`
	Archived        bool      `json:"archived"`
	Disabled        bool      `json:"disabled"`
	Owner           Owner     `json:"owner"`
	License         License   `json:"license"`
}

// Owner represents a repository owner
type Owner struct {
	Login     string `json:"login"`
	ID        int    `json:"id"`
	AvatarURL string `json:"avatar_url"`
	HTMLURL   string `json:"html_url"`
	Type      string `json:"type"`
}

// License represents repository license
type License struct {
	Key    string `json:"key"`
	Name   string `json:"name"`
	SPDXID string `json:"spdx_id"`
	URL    string `json:"url"`
}

// Commit represents a GitHub commit
type Commit struct {
	SHA         string       `json:"sha"`
	NodeID      string       `json:"node_id"`
	Commit      CommitData   `json:"commit"`
	Author      *User        `json:"author"`
	Committer   *User        `json:"committer"`
	Parents     []CommitRef  `json:"parents"`
	HTMLURL     string       `json:"html_url"`
	CommentsURL string       `json:"comments_url"`
	Stats       *CommitStats `json:"stats,omitempty"`
	Files       []CommitFile `json:"files,omitempty"`
}

// CommitData represents commit data
type CommitData struct {
	Message      string       `json:"message"`
	Author       GitUser      `json:"author"`
	Committer    GitUser      `json:"committer"`
	Tree         CommitRef    `json:"tree"`
	URL          string       `json:"url"`
	CommentCount int          `json:"comment_count"`
	Verification Verification `json:"verification"`
}

// GitUser represents git user info
type GitUser struct {
	Name  string    `json:"name"`
	Email string    `json:"email"`
	Date  time.Time `json:"date"`
}

// CommitRef represents a commit reference
type CommitRef struct {
	SHA string `json:"sha"`
	URL string `json:"url"`
}

// CommitStats represents commit statistics
type CommitStats struct {
	Additions int `json:"additions"`
	Deletions int `json:"deletions"`
	Total     int `json:"total"`
}

// CommitFile represents a file changed in a commit
type CommitFile struct {
	Filename         string `json:"filename"`
	Additions        int    `json:"additions"`
	Deletions        int    `json:"deletions"`
	Changes          int    `json:"changes"`
	Status           string `json:"status"` // "added", "removed", "modified", "renamed"
	RawURL           string `json:"raw_url"`
	BlobURL          string `json:"blob_url"`
	Patch            string `json:"patch,omitempty"`
	PreviousFilename string `json:"previous_filename,omitempty"`
}

// Verification represents commit signature verification
type Verification struct {
	Verified  bool   `json:"verified"`
	Reason    string `json:"reason"`
	Signature string `json:"signature"`
	Payload   string `json:"payload"`
}

// CommitComparison represents a comparison between two commits
type CommitComparison struct {
	BaseCommit      Commit       `json:"base_commit"`
	MergeBaseCommit Commit       `json:"merge_base_commit"`
	Status          string       `json:"status"`
	AheadBy         int          `json:"ahead_by"`
	BehindBy        int          `json:"behind_by"`
	TotalCommits    int          `json:"total_commits"`
	Commits         []Commit     `json:"commits"`
	Files           []CommitFile `json:"files"`
	HTMLURL         string       `json:"html_url"`
	PermalinkURL    string       `json:"permalink_url"`
	DiffURL         string       `json:"diff_url"`
	PatchURL        string       `json:"patch_url"`
}

// User represents a GitHub user
type User struct {
	Login     string `json:"login"`
	ID        int    `json:"id"`
	AvatarURL string `json:"avatar_url"`
	HTMLURL   string `json:"html_url"`
}

// PullRequest represents a GitHub pull request
type PullRequest struct {
	ID      int    `json:"id"`
	Number  int    `json:"number"`
	Title   string `json:"title"`
	Body    string `json:"body"`
	State   string `json:"state"`
	HTMLURL string `json:"html_url"`
	User    User   `json:"user"`
	Head    struct {
		Ref string `json:"ref"`
		SHA string `json:"sha"`
	} `json:"head"`
	Base struct {
		Ref string `json:"ref"`
		SHA string `json:"sha"`
	} `json:"base"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	MergedAt  *time.Time `json:"merged_at"`
}

// Client represents a GitHub API client
type Client struct {
	Token   string
	BaseURL string
}

// NewClient creates a new GitHub client
func NewClient(token string) *Client {
	return &Client{
		Token:   token,
		BaseURL: "https://api.github.com",
	}
}

// GetRepository retrieves repository information
func (c *Client) GetRepository(owner, repo string) (*Repository, error) {
	url := fmt.Sprintf("%s/repos/%s/%s", c.BaseURL, owner, repo)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	if c.Token != "" {
		req.Header.Set("Authorization", "token "+c.Token)
	}
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GitHub API error: %s", resp.Status)
	}

	var repository Repository
	if err := json.NewDecoder(resp.Body).Decode(&repository); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &repository, nil
}

// GetCommits retrieves repository commits with basic information
func (c *Client) GetCommits(owner, repo string, since *time.Time, limit int) ([]Commit, error) {
	url := fmt.Sprintf("%s/repos/%s/%s/commits", c.BaseURL, owner, repo)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Add query parameters
	q := req.URL.Query()
	if since != nil {
		q.Add("since", since.Format(time.RFC3339))
	}
	if limit > 0 {
		q.Add("per_page", fmt.Sprintf("%d", limit))
	}
	req.URL.RawQuery = q.Encode()

	if c.Token != "" {
		req.Header.Set("Authorization", "token "+c.Token)
	}
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GitHub API error: %s", resp.Status)
	}

	var commits []Commit
	if err := json.NewDecoder(resp.Body).Decode(&commits); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return commits, nil
}

// GetCommit retrieves detailed information about a specific commit
func (c *Client) GetCommit(owner, repo, sha string) (*Commit, error) {
	url := fmt.Sprintf("%s/repos/%s/%s/commits/%s", c.BaseURL, owner, repo, sha)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	if c.Token != "" {
		req.Header.Set("Authorization", "token "+c.Token)
	}
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GitHub API error: %s", resp.Status)
	}

	var commit Commit
	if err := json.NewDecoder(resp.Body).Decode(&commit); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &commit, nil
}

// GetCommitsInRange retrieves commits between two references (commits, branches, tags)
func (c *Client) GetCommitsInRange(owner, repo, base, head string, limit int) ([]Commit, error) {
	url := fmt.Sprintf("%s/repos/%s/%s/commits", c.BaseURL, owner, repo)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Add query parameters
	q := req.URL.Query()
	q.Add("sha", head) // Get commits from this branch/commit
	if limit > 0 {
		q.Add("per_page", fmt.Sprintf("%d", limit))
	}
	req.URL.RawQuery = q.Encode()

	if c.Token != "" {
		req.Header.Set("Authorization", "token "+c.Token)
	}
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GitHub API error: %s", resp.Status)
	}

	var allCommits []Commit
	if err := json.NewDecoder(resp.Body).Decode(&allCommits); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	// Filter commits to only include those after the base commit
	if base != "" {
		baseCommit, err := c.GetCommit(owner, repo, base)
		if err != nil {
			return nil, fmt.Errorf("failed to get base commit: %w", err)
		}

		var filteredCommits []Commit
		for _, commit := range allCommits {
			if commit.Commit.Author.Date.After(baseCommit.Commit.Author.Date) {
				filteredCommits = append(filteredCommits, commit)
			}
		}
		return filteredCommits, nil
	}

	return allCommits, nil
}

// GetPullRequests retrieves repository pull requests
func (c *Client) GetPullRequests(owner, repo, state string, limit int) ([]PullRequest, error) {
	url := fmt.Sprintf("%s/repos/%s/%s/pulls", c.BaseURL, owner, repo)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Add query parameters
	q := req.URL.Query()
	if state != "" {
		q.Add("state", state) // "open", "closed", "all"
	}
	if limit > 0 {
		q.Add("per_page", fmt.Sprintf("%d", limit))
	}
	req.URL.RawQuery = q.Encode()

	if c.Token != "" {
		req.Header.Set("Authorization", "token "+c.Token)
	}
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GitHub API error: %s", resp.Status)
	}

	var pullRequests []PullRequest
	if err := json.NewDecoder(resp.Body).Decode(&pullRequests); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return pullRequests, nil
}

func (c *Client) GetDailyCommitSummary(owner, repo string) {
	since := time.Now().AddDate(0, 0, -1) // Last 24 hours
	commits, err := c.GetCommits(owner, repo, &since, 100)
	if err != nil {
		log.Fatal(err)
	}

	totalFiles := 0
	totalAdditions := 0
	totalDeletions := 0
	authors := make(map[string]int)

	for _, commit := range commits {
		// Get detailed commit info
		detailed, err := c.GetCommit(owner, repo, commit.SHA)
		if err != nil {
			continue
		}

		if detailed.Stats != nil {
			totalAdditions += detailed.Stats.Additions
			totalDeletions += detailed.Stats.Deletions
		}
		totalFiles += len(detailed.Files)

		if detailed.Author != nil {
			authors[detailed.Author.Login]++
		}
	}

	fmt.Printf("Daily Summary:\n")
	fmt.Printf("- Commits: %d\n", len(commits))
	fmt.Printf("- Files changed: %d\n", totalFiles)
	fmt.Printf("- Lines added: %d\n", totalAdditions)
	fmt.Printf("- Lines deleted: %d\n", totalDeletions)
	fmt.Printf("- Contributors: %d\n", len(authors))
}

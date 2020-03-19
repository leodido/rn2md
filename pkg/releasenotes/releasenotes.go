package releasenotes

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/google/go-github/v28/github"
)

var (
	releaseNoteRegexp = regexp.MustCompile("(?s)```release-note(.+?)```")
	typologyRegexp    = regexp.MustCompile(`(?m)(.+?)(\((.+)\))?: ?(.*)`)
)

const defaultGitHubBaseURI = "https://github.com"

// ReleaseNote ...
type ReleaseNote struct {
	Typology    string
	Scope       string
	Description string
	URI         string
	Num         int
}

// Client ...
type Client struct {
	c *github.Client
	s *statistics
}

// NewClient ...
func NewClient() *Client {
	return &Client{
		c: github.NewClient(nil),
		s: &statistics{
			total:     0,
			totalNone: 0,
			authors:   map[string]int{},
		},
	}
}

// Get returns the list of release notes found for the given parameters.
func (c *Client) Get(org, repo, branch, milestone string) ([]ReleaseNote, error) {
	ctx := context.Background()
	listingOpts := &github.PullRequestListOptions{
		State:     "closed",
		Base:      branch,
		Sort:      "updated",
		Direction: "desc",
		ListOptions: github.ListOptions{
			PerPage: 1000,
		},
	}
	prs, _, err := c.c.PullRequests.List(ctx, org, repo, listingOpts)
	if err != nil {
		return nil, err
	}

	releaseNotes := []ReleaseNote{}
	for _, p := range prs {
		if p.GetMilestone().GetTitle() != milestone {
			continue
		}
		c.s.total++
		c.s.authors[p.GetUser().GetLogin()] = c.s.authors[p.GetUser().GetLogin()] + 1

		res := releaseNoteRegexp.FindStringSubmatch(p.GetBody())
		if len(res) < 1 {
			continue
		}
		note := strings.TrimSpace(res[1])
		if note == "NONE" || note == "none" {
			c.s.totalNone++
			continue
		}
		notes := strings.Split(note, "\n")
		for _, n := range notes {
			n = strings.Trim(n, "\r")
			matches := typologyRegexp.FindStringSubmatch(n)
			if len(matches) < 5 {
				return nil, fmt.Errorf("error extracting type from release note")
			}

			rn := ReleaseNote{
				Typology:    matches[1],
				Scope:       matches[3],
				Description: n,
				URI:         fmt.Sprintf("%s/%s/%s/pull/%d", defaultGitHubBaseURI, org, repo, p.GetNumber()),
				Num:         p.GetNumber(),
			}
			releaseNotes = append(releaseNotes, rn)
		}
	}

	return releaseNotes, nil
}

// TotalNone ...
func (c *Client) TotalNone() int {
	return c.s.totalNone
}

// TotalWithNotes ...
func (c *Client) TotalWithNotes() int {
	return c.s.total - c.s.totalNone
}

// TotalByAuthors returns the number of PRs by author username.
func (c *Client) TotalByAuthors() map[string]int {
	return c.s.authors
}

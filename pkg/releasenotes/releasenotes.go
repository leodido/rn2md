package releasenotes

import (
	"context"
	"errors"
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/gofri/go-github-ratelimit/github_ratelimit"
	"github.com/google/go-github/v28/github"
	"golang.org/x/oauth2"
)

var (
	releaseNoteRegexp = regexp.MustCompile("(?s)```release-note(.+?)```")
	typologyRegexp    = regexp.MustCompile(`(?m)(.+?)(\((.+)\))?(!)?: ?(.*)`)
)

const defaultGitHubBaseURI = "https://github.com"

// ReleaseNote ...
type ReleaseNote struct {
	Typology    string
	Scope       string
	Description string
	URI         string
	Num         int
	Author      string
	AuthorURL   string
}

type ReleaseNotes []ReleaseNote

// Client ...
type Client struct {
	c *github.Client
}

// NewClient ...
func NewClient(token string) *Client {
	var client *github.Client
	// Eventually create an authenticated client
	if token != "" {
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: token},
		)
		tc := oauth2.NewClient(context.Background(), ts)
		client = github.NewClient(tc)
	} else {
		// Force use a rate limited client. It will be slow but should overcome rate limiting issues.
		rateLimiter, err := github_ratelimit.NewRateLimitWaiterClient(nil)
		if err != nil {
			log.Fatal(err)
		}
		client = github.NewClient(rateLimiter)
	}

	return &Client{
		c: client,
	}
}

// Get returns the list of release notes found for the given parameters.
func (c *Client) Get(org, repo, branch, milestone string) (ReleaseNotes, *Statistics, error) {
	ctx := context.Background()
	listingOpts := &github.PullRequestListOptions{
		State:     "closed",
		Base:      branch,
		Sort:      "updated",
		Direction: "desc",
		ListOptions: github.ListOptions{
			// https://docs.github.com/en/rest/pulls/pulls?apiVersion=2022-11-28#list-pull-requests
			// per_page integer
			//
			// The number of results per page (max 100). For more information, see "Using pagination in the REST API."
			// Default: 30
			PerPage: 100,
		},
	}

	// https://docs.github.com/en/rest/pulls/pulls?apiVersion=2022-11-28#list-pull-requests
	// page integer
	//
	// The page number of the results to fetch. For more information, see "Using pagination in the REST API."
	// Default: 1
	page := 1
	prs := make([]*github.PullRequest, 0)
	for {
		listingOpts.Page = page
		pagedPrs, _, err := c.c.PullRequests.List(ctx, org, repo, listingOpts)
		var rateLimitErr *github.RateLimitError
		if errors.As(err, &rateLimitErr) {
			return nil, nil, fmt.Errorf("hit rate limiting")
		}
		if err != nil {
			return nil, nil, err
		}
		prs = append(prs, pagedPrs...)
		if len(pagedPrs) < listingOpts.PerPage {
			// We collected all prs!
			break
		}
		page++
	}

	var releaseNotes []ReleaseNote
	s := &Statistics{
		total:     0,
		nonFacing: 0,
		authors:   make(map[string]int64),
	}
	for _, p := range prs {
		num := p.GetNumber()
		if p.GetMergedAt().Equal(time.Time{}) {
			continue
		}
		if p.GetMilestone().GetTitle() != milestone {
			continue
		}
		s.total++
		s.authors[p.GetUser().GetLogin()] = s.authors[p.GetUser().GetLogin()] + 1

		res := releaseNoteRegexp.FindStringSubmatch(p.GetBody())
		if len(res) < 1 {
			continue
		}
		note := strings.TrimSpace(res[1])
		if strings.EqualFold(note, "NONE") {
			s.nonFacing++
			rn := ReleaseNote{
				Typology:    "none",
				Scope:       "",
				Description: p.GetTitle(),
				URI:         fmt.Sprintf("%s/%s/%s/pull/%d", defaultGitHubBaseURI, org, repo, num),
				Num:         num,
				Author:      fmt.Sprintf("@%s", p.GetUser().GetLogin()),
				AuthorURL:   p.GetUser().GetHTMLURL(),
			}
			releaseNotes = append(releaseNotes, rn)
			continue
		}

		notes := strings.Split(note, "\n")
		for _, n := range notes {
			n = strings.Trim(n, "\r")
			matches := typologyRegexp.FindStringSubmatch(n)
			if len(matches) < 6 {
				return nil, nil, fmt.Errorf("error extracting type from release note, pr: %d", num)
			}

			rn := ReleaseNote{
				Typology:    matches[1],
				Scope:       matches[3],
				Description: n,
				URI:         fmt.Sprintf("%s/%s/%s/pull/%d", defaultGitHubBaseURI, org, repo, num),
				Num:         num,
				Author:      fmt.Sprintf("@%s", p.GetUser().GetLogin()),
				AuthorURL:   p.GetUser().GetHTMLURL(),
			}
			if matches[4] == "!" {
				rn.Typology = "BREAKING CHANGE"
			}
			releaseNotes = append(releaseNotes, rn)
		}
	}

	return releaseNotes, s, nil
}

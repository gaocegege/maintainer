package user

import (
	"context"
	"fmt"
	"log"
	"sort"
	"time"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

type Interface interface {
	GetDaily(username string, date time.Time) (string, error)
	GetDailyToday(username string) (string, error)
	GetReport(username string, begin, end time.Time) (string, error)
}

type imp struct {
	token   string
	perPage int
}

func New(token string) Interface {
	return &imp{
		token: token,
		// Ref https://developer.github.com/v3/activity/events/
		// Events support pagination, however the per_page option is unsupported.
		// The fixed page size is 30 items. Fetching up to ten pages is supported,
		// for a total of 300 events.
		perPage: 30,
	}
}

func (i *imp) GetReport(username string, begin, end time.Time) (string, error) {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: i.token},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)
	index := 1
	shouldBreak := false
	log.SetPrefix("daily-generator")
	res := []string{}
	for {
		log.Printf("Getting events page %d from Github", index)
		events, response, err := client.Activity.ListEventsPerformedByUser(username, false, &github.ListOptions{
			Page: index,
		})
		if err != nil {
			return "", err
		}
		if response.StatusCode != 200 {
			return "", fmt.Errorf("Failed to get the events for the user: %s",
				response.Status)
		}
		for _, event := range events {
			eventTime := event.CreatedAt
			if eventTime == nil {
				continue
			}
			if eventTime.After(end) {
				continue
			} else if eventTime.Before(begin) {
				shouldBreak = true
				break
			}
			res = append(res, i.ComposeEvent(event))
		}
		if shouldBreak {
			break
		}
		index++
	}

	// Avoid dup.
	res = removeDuplicates(res)
	sort.Strings(res)
	returnStr := ""
	for _, s := range res {
		returnStr += s
	}
	return returnStr, nil
}

func (i *imp) GetDailyToday(username string) (string, error) {
	return i.GetReport(username, beginningTime(), beginningTime().Add(24*time.Hour))
}

func (i *imp) GetDaily(username string, date time.Time) (string, error) {
	return i.GetReport(username, date, date.Add(24*time.Hour))
}

func (i imp) ComposeEvent(event *github.Event) string {
	template := "- "
	if event.Public != nil {
		if *event.Public {
			template = template + "[Public]"
		} else {
			template = template + "[Private]"
		}
	}
	switch *event.Type {
	// case "CommitCommentEvent":
	// 	e := event.Payload().(*github.CommitCommentEvent)
	// 	template += "[CommitComment]"
	case "PullRequestEvent":
		e := event.Payload().(*github.PullRequestEvent)
		template += "[PullRequest]"
		template += fmt.Sprintf(" %s %s\n", *e.PullRequest.HTMLURL, *e.PullRequest.Title)
	case "IssuesEvent":
		e := event.Payload().(*github.IssuesEvent)
		template += "[Issue]"
		template += fmt.Sprintf(" %s %s\n", *e.Issue.HTMLURL, *e.Issue.Title)
	case "PullRequestReviewCommentEvent":
		e := event.Payload().(*github.PullRequestReviewCommentEvent)
		template += "[PullRequestReview]"
		template += fmt.Sprintf(" %s %s\n", *e.PullRequest.HTMLURL, *e.PullRequest.Title)
	case "IssueCommentEvent":
		e := event.Payload().(*github.IssueCommentEvent)
		template += "[IssueComment]"
		template += fmt.Sprintf(" %s %s\n", *e.Issue.HTMLURL, *e.Issue.Title)
	default:
		return ""
	}
	return template
}

func beginningTime() time.Time {
	return time.Now().Local().Truncate(24 * time.Hour)
}

func removeDuplicates(elements []string) []string {
	// Use map to record duplicates as we find them.
	encountered := map[string]bool{}
	result := []string{}

	for v := range elements {
		if encountered[elements[v]] == true {
			// Do not add duplicate.
		} else {
			// Record this element as an encountered element.
			encountered[elements[v]] = true
			// Append to result slice.
			result = append(result, elements[v])
		}
	}
	// Return the new slice.
	return result
}
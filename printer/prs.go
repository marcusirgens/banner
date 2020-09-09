package printer

import (
	"context"
	"fmt"
	gh "github.com/google/go-github/v32/github"
	"github.com/logrusorgru/aurora/v3"
	"github.com/marcusirgens/banner/github"
	"time"
)

func PrintEvents(ctx context.Context, accesstoken string) {
	max := 100
	c := 0
	d := time.Hour * 24
	maxAge := time.Now().Add(d * -1)
	for event := range getEvents(ctx, accesstoken) {
		if event.GetCreatedAt().Before(maxAge) {
			return
		}
		if printEvent(event) {
			c++
		}
		if c >= max {
			return
		}
	}
}

func getEvents(ctx context.Context, accesstoken string) <-chan *gh.Event {
	ch := make(chan *gh.Event)

	go func() {
		defer close(ch)

		cli := github.NewClient(ctx, accesstoken)

		user, _, err := cli.Users.Get(ctx, "")
		if err != nil {
			fmt.Println("Could not get current GitHub user")
			return
		}

		un := user.GetLogin()

		opts := &gh.ListOptions{}

		i := 0
		max := 10

		for {
			es, resp, err := cli.Activity.ListEventsReceivedByUser(ctx, un, false, opts)
			if err != nil {
				fmt.Println("Could not get activity stream")
				return
			}

			for _, e := range es {
				ch <- e
			}

			if resp.NextPage == 0 || i > max {
				return
			}

			opts.Page = resp.NextPage
		}
	}()

	return ch
}

func printEvent(event *gh.Event) (handled bool) {
	pp, err := event.ParsePayload()
	if err != nil {
		fmt.Println(err)
	}
	switch e := pp.(type) {
	case *gh.PushEvent:
		return printPushEvent(event, e)
	case *gh.PullRequestEvent:
		return printPullRequestEvent(event, e)
	}

	return
}

func printPullRequestEvent(event *gh.Event, pr *gh.PullRequestEvent) (handled bool) {
	n := aurora.Cyan(event.Repo.GetName())

	switch pr.GetAction() {
	case "closed":
		fmt.Printf("%s: pr #%d \"%s\" closed\n", n, pr.GetNumber(), pr.PullRequest.GetTitle())
		handled = true
	case "opened":
		fmt.Printf("%s: pr #%d \"%s\" opened\n", n, pr.GetNumber(), pr.PullRequest.GetTitle())
		handled = true
	case "reopened":
		fmt.Printf("%s: pr #%d \"%s\" reopened\n", n, pr.GetNumber(), pr.PullRequest.GetTitle())
		handled = true
	case "edited":
		fmt.Printf("%s: pr #%d \"%s\" edited\n", n, pr.GetNumber(), pr.PullRequest.GetTitle())
		handled = true
	default:
		fmt.Printf("%s: pr #%d \"%s\" updated\n", n, pr.GetNumber(), pr.PullRequest.GetTitle())
		handled = true
	}
	return
}

func printPushEvent(event *gh.Event, pe *gh.PushEvent) (handled bool) {
	// fmt.Printf("%s pushed to %s\n", event.Actor.GetLogin(), event.Repo.GetName())
	return false
}

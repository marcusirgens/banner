package printer

import (
	"context"
	"fmt"
	gh "github.com/google/go-github/v32/github"
	"github.com/logrusorgru/aurora/v3"
	"github.com/marcusirgens/banner/github"
	"os"
	"text/template"
	"time"
)

func PrintEvents(ctx context.Context, accesstoken string) {
	bannerShown := false

	roundedTime := roundTime(oneDayAgo())

	for notification := range getNotifications(ctx, accesstoken, roundedTime) {
		if !bannerShown {
			fmt.Fprintln(os.Stdout, "Github notifications:")
			bannerShown = true
		}
		printNotification(notification)
	}
	if bannerShown {
		fmt.Fprintln(os.Stdout, "")
	}
}

func roundTime(then time.Time) time.Time {
	sec := time.Second * time.Duration(then.Second())
	min := time.Minute * (time.Duration(then.Minute()) % 5)

	diff := sec + min
	return then.Add(-diff)
}

func printNotification(notification *gh.Notification) {
	tmplTxt := fmt.Sprintf(
		`  %s in %s: {{.Type}} "{{.Title}}"`,
		aurora.Cyan("{{.Time}}"),
		aurora.Cyan("{{.Repo}}"),
	)
	tmpl, err := template.New("notification").Parse(tmplTxt)

	updateTime := notification.GetUpdatedAt().In(inferTimeZone())

	if err != nil {
		// nop
		return
	}

	typeAbbr, err := getTypeAbbr(notification.GetSubject().GetType())
	if err != nil {
		typeAbbr = notification.GetSubject().GetType()
	}

	err = tmpl.Execute(os.Stdout, struct {
		Time   string
		Action string
		Repo   string
		Type   string
		Title  string
	}{
		Time:   updateTime.Format("Mon 15:04"),
		Action: notification.GetReason(),
		Repo:   notification.GetRepository().GetFullName(),
		Type:   typeAbbr,
		Title:  notification.GetSubject().GetTitle(),
	})
	if err != nil {
		fmt.Printf("error when parsing notification: %v", err)
	}
	fmt.Fprint(os.Stdout, "\n")
}

func getTypeAbbr(typ string) (string, error) {
	switch typ {
	case "PullRequest":
		return "pr", nil
	case "Issue":
		return "issue", nil
	default:
		return "", fmt.Errorf("unknown type %s", typ)
	}
}

func inferTimeZone() *time.Location {
	t := time.Now()
	return t.Location()

}

func oneDayAgo() time.Time {
	diff := time.Hour * 24 * 7
	now := time.Now()
	return now.Add(-diff)
}

func getNotifications(ctx context.Context, accesstoken string, since time.Time) <-chan *gh.Notification {
	ch := make(chan *gh.Notification)

	go func() {
		defer close(ch)

		cli := github.NewClient(ctx, accesstoken)

		// defined here for pagination
		opts := &gh.NotificationListOptions{
			Participating: true,
			All:           false,
			Since:         since,
		}

		for {
			notifications, resp, err := cli.Activity.ListNotifications(ctx, opts)
			if err != nil {
				// handle this
				return
			}
			for _, n := range notifications {
				ch <- n
			}

			if resp.NextPage == 0 {
				return
			}

			opts.Page = resp.NextPage
		}

	}()

	return ch
}

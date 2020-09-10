package github

import (
	"context"
	"fmt"
	"github.com/google/go-github/v32/github"
	"github.com/gregjones/httpcache"
	"github.com/gregjones/httpcache/diskcache"
	"golang.org/x/oauth2"
	"net/http"
	"os"
)

func NewClient(ctx context.Context, accesstoken string) *github.Client {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: accesstoken},
	)
	_ = oauth2.NewClient(ctx, ts)

	tr := newCachedAuthTransport(accesstoken)
	cl := &http.Client{Transport: tr}
	return github.NewClient(cl)
}

type cachedAuthTransport struct {
	inner http.RoundTripper
	token string
}

// newCachedAuthTransport creates a cached HTTP transport that attempts to store
// http requests in the user's cache dir as returned by os.UserCacheDir.
func newCachedAuthTransport(accesstoken string) *cachedAuthTransport {
	dir, err := os.UserCacheDir()
	if err != nil {
		// Can't cache, use a memory transport instead
		return &cachedAuthTransport{
			inner: httpcache.NewMemoryCacheTransport(),
			token: accesstoken,
		}
	}

	cDir := fmt.Sprintf("%s/gh-banner", dir)
	err = os.MkdirAll(cDir, os.ModePerm)
	if err != nil {
		// Can't create cache directory, use memory transport
		return &cachedAuthTransport{
			inner: httpcache.NewMemoryCacheTransport(),
			token: accesstoken,
		}
	}

	c := diskcache.New(cDir)
	tr := httpcache.NewTransport(c)
	return &cachedAuthTransport{
		inner: tr,
		token: accesstoken,
	}
}

func (t *cachedAuthTransport) RoundTrip(request *http.Request) (*http.Response, error) {
	request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", t.token))
	return t.inner.RoundTrip(request)
}

package providers

import (
	"context"
	"github.com/nbedos/citop/cache"
	"os"
	"testing"
	"time"
)

func TestGitLabGetUserBuilds(t *testing.T) {
	token := os.Getenv("GITLAB_API_TOKEN")
	if token == "" {
		t.Fatal("Environment variable GITLAB_API_TOKEN is not set")
	}
	account := cache.Account{
		ID:       "gitlab",
		URL:      "example.com/api/v3",
		UserID:   "42",
		Username: "oops",
	}
	client := NewGitLabClient(account.ID, token, 100*time.Millisecond)

	c := make(chan []cache.Inserter)
	errc := make(chan error)
	ctx := context.Background()
	go func() {
		repository, err := client.Repository(ctx, "https://gitlab.com/nbedos/citop")
		if err != nil {
			close(c)
			errc <- err
			return
		}
		c <- []cache.Inserter{repository}
		err = client.LastBuilds(ctx, repository, 20, c)
		close(c)
		errc <- err
	}()

	inserters := []cache.Inserter{account}
	for is := range c {
		inserters = append(inserters, is...)
	}

	if err := <-errc; err != nil {
		t.Fatal(err)
	}

	if _, err := cache.TemporaryCache(context.Background(), "gitlab", inserters); err != nil {
		t.Fatal(err)
	}
}

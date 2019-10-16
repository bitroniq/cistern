package providers

import (
	"context"
	"github.com/nbedos/citop/cache"
	"os"
	"testing"
	"time"
)

var citopURL = "https://github.com/nbedos/citop"

func Test_RecentRepoBuilds(t *testing.T) {
	token := os.Getenv("TRAVIS_API_TOKEN")
	if token == "" {
		t.Fatal("Environment variable TRAVIS_API_TOKEN is not set")
	}

	client := NewTravisClient(TravisOrgURL, TravisPusherHost, token, "travis", time.Millisecond*time.Duration(50))

	errc := make(chan error)
	c := make(chan []cache.Inserter)
	go func() {
		err := client.RepositoryBuilds(context.Background(), citopURL, 20, 5, c)
		close(c)
		errc <- err
	}()

	inserters := make([]cache.Inserter, 0)
	for is := range c {
		inserters = append(inserters, is...)
	}

	if err := <-errc; err != nil {
		t.Fatal(err)
	}
}

func TestTravisclient_Repository(t *testing.T) {
	token := os.Getenv("TRAVIS_API_TOKEN")
	if token == "" {
		t.Fatal("Environment variable TRAVIS_API_TOKEN is not set")
	}

	client := NewTravisClient(TravisOrgURL, TravisPusherHost, token, "travis", time.Millisecond*time.Duration(50))

	t.Run("repository found", func(t *testing.T) {
		repo, err := client.Repository(context.Background(), "https://github.com/nbedos/citop")
		if err != nil {
			t.Fatal(err)
		}
		if repo.Name != "citop" {
			t.Fatalf("invalid name: %s", repo.Name)
		}
	})

	t.Run("repository not found", func(t *testing.T) {
		_, err := client.Repository(context.Background(), "https://github.com/nbedos/citop-404")
		switch e := err.(type) {
		case HTTPError:
			if e.Status != 404 {
				t.Fatalf("expected status 404 but got %d", e.Status)
			}
		default:
			t.Fatal("expected HTTPError")
		}
	})
}

package ui

import (
    "testing"
    "time"
    "os"

    "github.com/agnivo988/Repo-lyzer/internal/github"
)

func TestCacheStatusFromExpires_FreshAndStale(t *testing.T) {
    future := time.Now().Add(1 * time.Hour)
    past := time.Now().Add(-1 * time.Hour)

    if s := cacheStatusFromExpires(future); s != "fresh" {
        t.Fatalf("expected fresh for future expires, got %s", s)
    }
    if s := cacheStatusFromExpires(past); s != "stale" {
        t.Fatalf("expected stale for past expires, got %s", s)
    }
}

func TestCachedThenLive_ProducesRefreshed(t *testing.T) {
    // ensure no persisted current analysis interferes with tests
    _ = os.Remove("exports/current_analysis.json")

    m := NewMainModel(nil, nil)

    // simulate loading cached result (stateLoading expected)
    m.state = stateLoading

    repo := &github.Repo{FullName: "owner/repo"}
    cached := CachedAnalysisResult{
        Result:   AnalysisResult{Repo: repo},
        IsCached: true,
        CachedAt: time.Now().Add(-2 * time.Hour),
        ExpiresAt: time.Now().Add(1 * time.Hour), // not expired => fresh
    }

    // send cached result
    mm, _ := m.Update(cached)
    m1 := mm.(MainModel)
    if m1.dashboard.cacheStatus != "fresh" {
        t.Fatalf("expected dashboard status fresh after cached result, got %s", m1.dashboard.cacheStatus)
    }

    // now send live AnalysisResult for same repo
    live := AnalysisResult{Repo: repo}
    mm2, _ := m1.Update(live)
    m2 := mm2.(MainModel)
    if m2.dashboard.cacheStatus != "refreshed" {
        t.Fatalf("expected dashboard status refreshed after live result, got %s", m2.dashboard.cacheStatus)
    }

    // cleanup
    _ = os.Remove("exports/current_analysis.json")
}

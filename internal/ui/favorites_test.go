package ui

import (
	"os"
	"testing"
	"time"
)

func TestLoadFavorites_Empty(t *testing.T) {
	favs, err := LoadFavorites()
	if err != nil && !os.IsNotExist(err) {
		t.Fatalf("LoadFavorites() unexpected error = %v", err)
	}
	if favs == nil {
		t.Fatal("LoadFavorites() returned nil")
	}
}

func TestFavorites_Add(t *testing.T) {
	favs := &Favorites{Items: []Favorite{}}

	// Add first favorite
	favs.Add("owner/repo1")
	
	if len(favs.Items) != 1 {
		t.Errorf("After Add(), len = %d, want 1", len(favs.Items))
	}
	if favs.Items[0].RepoName != "owner/repo1" {
		t.Errorf("RepoName = %s, want owner/repo1", favs.Items[0].RepoName)
	}
	if favs.Items[0].UseCount != 1 {
		t.Errorf("UseCount = %d, want 1", favs.Items[0].UseCount)
	}

	// Add same favorite again (should update, not duplicate)
	favs.Add("owner/repo1")
	
	if len(favs.Items) != 1 {
		t.Errorf("After duplicate Add(), len = %d, want 1", len(favs.Items))
	}
	if favs.Items[0].UseCount != 2 {
		t.Errorf("UseCount after duplicate = %d, want 2", favs.Items[0].UseCount)
	}

	// Add different favorite
	favs.Add("owner/repo2")
	
	if len(favs.Items) != 2 {
		t.Errorf("After second Add(), len = %d, want 2", len(favs.Items))
	}
}

func TestFavorites_Remove(t *testing.T) {
	favs := &Favorites{Items: []Favorite{
		{RepoName: "owner/repo1", UseCount: 1},
		{RepoName: "owner/repo2", UseCount: 2},
		{RepoName: "owner/repo3", UseCount: 3},
	}}

	// Remove middle item
	favs.Remove("owner/repo2")
	
	if len(favs.Items) != 2 {
		t.Errorf("After Remove(), len = %d, want 2", len(favs.Items))
	}
	
	// Verify correct items remain
	for _, fav := range favs.Items {
		if fav.RepoName == "owner/repo2" {
			t.Error("Removed item still exists")
		}
	}

	// Remove non-existent (should not panic)
	favs.Remove("owner/nonexistent")
	if len(favs.Items) != 2 {
		t.Error("Remove of non-existent changed length")
	}
}

func TestFavorites_IsFavorite(t *testing.T) {
	favs := &Favorites{Items: []Favorite{
		{RepoName: "owner/repo1"},
		{RepoName: "owner/repo2"},
	}}

	if !favs.IsFavorite("owner/repo1") {
		t.Error("IsFavorite() should return true for existing favorite")
	}
	if !favs.IsFavorite("owner/repo2") {
		t.Error("IsFavorite() should return true for existing favorite")
	}
	if favs.IsFavorite("owner/repo3") {
		t.Error("IsFavorite() should return false for non-favorite")
	}
}

func TestFavorites_UpdateUsage(t *testing.T) {
	now := time.Now()
	favs := &Favorites{Items: []Favorite{
		{RepoName: "owner/repo1", UseCount: 5, LastUsed: now.Add(-24 * time.Hour)},
	}}

	favs.UpdateUsage("owner/repo1")

	if favs.Items[0].UseCount != 6 {
		t.Errorf("UseCount = %d, want 6", favs.Items[0].UseCount)
	}
	if favs.Items[0].LastUsed.Before(now) {
		t.Error("LastUsed should be updated to recent time")
	}

	// Update non-existent (should not panic)
	favs.UpdateUsage("owner/nonexistent")
}

func TestFavorites_GetTopFavorites(t *testing.T) {
	favs := &Favorites{Items: []Favorite{
		{RepoName: "repo1"},
		{RepoName: "repo2"},
		{RepoName: "repo3"},
		{RepoName: "repo4"},
		{RepoName: "repo5"},
	}}

	// Get top 3
	top := favs.GetTopFavorites(3)
	if len(top) != 3 {
		t.Errorf("GetTopFavorites(3) len = %d, want 3", len(top))
	}

	// Get more than available
	top = favs.GetTopFavorites(10)
	if len(top) != 5 {
		t.Errorf("GetTopFavorites(10) len = %d, want 5", len(top))
	}

	// Get 0
	top = favs.GetTopFavorites(0)
	if len(top) != 0 {
		t.Errorf("GetTopFavorites(0) len = %d, want 0", len(top))
	}
}

func TestFavorites_Clear(t *testing.T) {
	favs := &Favorites{Items: []Favorite{
		{RepoName: "repo1"},
		{RepoName: "repo2"},
	}}

	favs.Clear()

	if len(favs.Items) != 0 {
		t.Errorf("After Clear(), len = %d, want 0", len(favs.Items))
	}
}

func TestFavorite_Fields(t *testing.T) {
	now := time.Now()
	fav := Favorite{
		RepoName: "owner/repo",
		AddedAt:  now,
		LastUsed: now,
		UseCount: 10,
		Notes:    "Test notes",
	}

	if fav.RepoName != "owner/repo" {
		t.Errorf("RepoName = %s, want owner/repo", fav.RepoName)
	}
	if fav.UseCount != 10 {
		t.Errorf("UseCount = %d, want 10", fav.UseCount)
	}
	if fav.Notes != "Test notes" {
		t.Errorf("Notes = %s, want 'Test notes'", fav.Notes)
	}
}

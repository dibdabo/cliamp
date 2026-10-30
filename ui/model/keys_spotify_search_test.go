package model

import (
	"context"
	"errors"
	"slices"
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"

	"github.com/bjarneo/cliamp/playlist"
)

func spotifySearchTestTracks() []playlist.Track {
	return []playlist.Track{
		{Title: "First", Artist: "Artist", Path: "spotify:track:first"},
		{Title: "Second", Artist: "Artist", Path: "spotify:track:second"},
		{Title: "Third", Artist: "Artist", Path: "spotify:track:third"},
	}
}

func TestSpotSearchSpaceTogglesSelectionAndAdvances(t *testing.T) {
	m := Model{
		plVisible: 10,
		spotSearch: spotSearchState{
			screen:  spotSearchResults,
			results: spotifySearchTestTracks(),
			cursor:  1,
		},
	}

	m.handleSpotSearchResultsKey(tea.KeyPressMsg{Code: tea.KeySpace})
	if !m.spotSearch.selected[1] {
		t.Fatal("result 1 was not selected")
	}
	if m.spotSearch.cursor != 2 {
		t.Fatalf("cursor = %d, want 2 after selecting", m.spotSearch.cursor)
	}

	m.spotSearch.cursor = 1
	m.handleSpotSearchResultsKey(tea.KeyPressMsg{Code: tea.KeySpace})
	if m.spotSearch.selected[1] {
		t.Fatal("result 1 remained selected after toggling it off")
	}
}

func TestSpotSearchBatchAppendUsesResultOrder(t *testing.T) {
	p := playlist.New()
	p.Add(playlist.Track{Title: "Existing", Path: "existing"})
	m := Model{
		player:   &fakeEngine{},
		playlist: p,
		spotSearch: spotSearchState{
			visible:  true,
			screen:   spotSearchResults,
			results:  spotifySearchTestTracks(),
			selected: map[int]bool{0: true, 2: true},
		},
	}

	if cmd := m.handleSpotSearchResultsKey(tea.KeyPressMsg{Text: "a"}); cmd != nil {
		t.Fatal("append command was non-nil while playback was active")
	}
	got := m.playlist.Tracks()
	if len(got) != 3 || got[0].Title != "Existing" || got[1].Title != "First" || got[2].Title != "Third" {
		t.Fatalf("playlist titles = %v, want [Existing First Third]", trackTitles(got))
	}
	if m.spotSearch.visible {
		t.Fatal("search remained open after batch append")
	}
}

func TestSpotSearchBatchQueueUsesResultOrder(t *testing.T) {
	p := playlist.New()
	p.Add(playlist.Track{Title: "Playing", Path: "playing"})
	m := Model{
		player:   &fakeEngine{},
		playlist: p,
		spotSearch: spotSearchState{
			visible:  true,
			screen:   spotSearchResults,
			results:  spotifySearchTestTracks(),
			selected: map[int]bool{0: true, 2: true},
		},
	}

	m.handleSpotSearchResultsKey(tea.KeyPressMsg{Text: "q"})
	got := m.playlist.QueueTracks()
	if len(got) != 2 || got[0].Title != "First" || got[1].Title != "Third" {
		t.Fatalf("queued titles = %v, want [First Third]", trackTitles(got))
	}
}

func TestSpotSearchResultsRenderSelectionAndHelp(t *testing.T) {
	m := Model{
		plVisible: 4,
		spotSearch: spotSearchState{
			screen:   spotSearchResults,
			results:  spotifySearchTestTracks(),
			selected: map[int]bool{1: true},
		},
	}

	if body := m.renderSpotSearchBody(); !strings.Contains(body, "[x] Artist - Second") {
		t.Fatalf("search body does not mark selected track: %q", body)
	}
	for _, want := range []string{"Space", "Mark"} {
		if help := m.spotSearchResultsHelpLine(); !strings.Contains(help, want) {
			t.Fatalf("search help %q does not contain %q", help, want)
		}
	}
}

func trackTitles(tracks []playlist.Track) []string {
	titles := make([]string, len(tracks))
	for i, track := range tracks {
		titles[i] = track.Title
	}
	return titles
}

func TestSpotSearchMixedBatch(t *testing.T) {
	for _, action := range []string{"a", "q"} {
		t.Run(action, func(t *testing.T) {
			p := playlist.New()
			p.Add(playlist.Track{Title: "Existing"})
			p.Queue(0)
			m := Model{
				player: &fakeEngine{}, playlist: p,
				spotSearch: spotSearchState{
					visible: true, screen: spotSearchResults,
					prov: &subProv{episodes: map[string][]playlist.Track{
						"Album": {{Title: "Album 1"}, {Title: "Album 2"}},
					}},
					results:  []playlist.Track{{Title: "First"}, albumResult("Album"), {Title: "Last"}},
					selected: map[int]bool{2: true, 0: true, 1: true},
				},
			}
			cmd := m.handleSpotSearchResultsKey(tea.KeyPressMsg{Text: action})
			if cmd == nil || !m.spotSearch.albumLoading || p.Len() != 1 {
				t.Fatal("batch should wait for albums without mutating the playlist")
			}
			updated, _ := m.Update(cmd())
			m = updated.(Model)
			want := []string{"Existing", "First", "Album 1", "Album 2", "Last"}
			if got := trackTitles(m.playlist.Tracks()); !slices.Equal(got, want) {
				t.Fatalf("playlist = %v, want %v", got, want)
			}
			if action == "q" {
				if got := trackTitles(m.playlist.QueueTracks()); !slices.Equal(got, want) {
					t.Fatalf("queue = %v, want %v", got, want)
				}
			}
			if m.spotSearch.visible {
				t.Fatal("batch should close search after success")
			}
		})
	}
}

func TestSpotSearchBatchFailurePreservesSelection(t *testing.T) {
	failure := errors.New("feed unavailable")
	m := Model{playlist: playlist.New(), spotSearch: spotSearchState{
		visible: true, screen: spotSearchResults,
		prov:     &subProv{err: failure},
		results:  []playlist.Track{{Title: "First"}, albumResult("Show")},
		selected: map[int]bool{0: true, 1: true},
	}}
	cmd := m.handleSpotSearchResultsKey(tea.KeyPressMsg{Text: "a"})
	updated, _ := m.Update(cmd())
	m = updated.(Model)
	if m.playlist.Len() != 0 || !m.spotSearch.visible || len(m.spotSearch.selected) != 2 || !strings.Contains(m.spotSearch.err, failure.Error()) {
		t.Fatal("failed expansion should preserve selection and leave playlist unchanged")
	}
}

func TestSpotSearchBatchCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	msg := fetchSpotSelectionCmd(ctx, contextAlbumLoader{}, []playlist.Track{albumResult("Album")}, spotAlbumAppend, 7)().(spotAlbumTracksMsg)
	if !errors.Is(msg.err, context.Canceled) {
		t.Fatalf("error = %v, want canceled", msg.err)
	}
	m := Model{playlist: playlist.New(), spotSearch: spotSearchState{visible: true, screen: spotSearchResults}}
	m.requests.spotAlbum = 7
	m.closeSpotSearch()
	msg.err = nil
	msg.tracks = []playlist.Track{{Title: "Late"}}
	updated, cmd := m.Update(msg)
	if cmd != nil || updated.(Model).playlist.Len() != 0 {
		t.Fatal("a canceled batch must not add late results")
	}
}

func TestSpotSearchResultsHelpKeepsBatchActions(t *testing.T) {
	m := newLayoutTestModel(40, 10)
	for _, key := range []string{"Space", "a", "q"} {
		if help := stripAnsi(m.spotSearchResultsHelpLine()); !strings.Contains(help, key) {
			t.Fatalf("help %q omits %q", help, key)
		}
	}
}

package spotify

import "testing"

func TestSpotifySearchQuery(t *testing.T) {
	for _, tt := range []struct{ query, want string }{
		{"Aja album Steely Dan", "Aja Steely Dan"},
		{"ALBUM Aja", "Aja"},
		{"album:Aja artist:Steely", "album:Aja artist:Steely"},
		{"album", "album"},
		{"Radiohead", "Radiohead"},
	} {
		t.Run(tt.query, func(t *testing.T) {
			if got := spotifySearchQuery(tt.query); got != tt.want {
				t.Fatalf("query = %q, want %q", got, tt.want)
			}
		})
	}
}

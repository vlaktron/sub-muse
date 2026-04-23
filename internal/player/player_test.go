package player

import (
	"testing"
	"time"

	"sub-muse/internal/subsonic"
)

func TestNewPlayer(t *testing.T) {
	p := NewPlayer()
	if p == nil {
		t.Error("Expected player to be non-nil")
	}
}

func TestFindPlayerBinary_Mpv(t *testing.T) {
	p := NewPlayer()
	binary := p.findPlayerBinary()
	if binary == "" {
		t.Skip("No player binary found (mpv or ffplay)")
	}
	if !containsAny(binary, []string{"mpv", "ffplay"}) {
		t.Errorf("Expected mpv or ffplay, got %s", binary)
	}
}

func TestGetStateWhenNotPlaying(t *testing.T) {
	p := NewPlayer()
	state := p.GetState()
	if state.IsPlaying {
		t.Error("Expected IsPlaying to be false")
	}
	if state.Song != nil {
		t.Error("Expected Song to be nil")
	}
	if state.Elapsed != 0 {
		t.Error("Expected Elapsed to be 0")
	}
}

func TestStopWhenNotPlaying(t *testing.T) {
	p := NewPlayer()
	err := p.Stop()
	if err != nil {
		t.Errorf("Stop() should not error when not playing, got %v", err)
	}
}

func TestPlayerStruct(t *testing.T) {
	p := &Player{}
	if p.playing {
		t.Error("Expected playing to be false initially")
	}
	if p.song != nil {
		t.Error("Expected song to be nil initially")
	}
}

func TestStateStruct(t *testing.T) {
	state := State{}
	if state.IsPlaying {
		t.Error("Expected IsPlaying to be false initially")
	}
	if state.Song != nil {
		t.Error("Expected Song to be nil initially")
	}
	if state.Elapsed != 0 {
		t.Error("Expected Elapsed to be 0 initially")
	}
}

func TestStateWithPlayingSong(t *testing.T) {
	song := &subsonic.Song{ID: "test", Title: "Test"}
	state := State{
		IsPlaying: true,
		Song:      song,
		Elapsed:   time.Second,
	}

	if !state.IsPlaying {
		t.Error("Expected IsPlaying to be true")
	}
	if state.Song == nil || state.Song.ID != "test" {
		t.Error("Expected Song to be set")
	}
	if state.Elapsed != time.Second {
		t.Errorf("Expected Elapsed to be 1s, got %v", state.Elapsed)
	}
}

func containsAny(s string, subs []string) bool {
	for _, sub := range subs {
		if len(s) >= len(sub) {
			for i := 0; i <= len(s)-len(sub); i++ {
				if s[i:i+len(sub)] == sub {
					return true
				}
			}
		}
	}
	return false
}

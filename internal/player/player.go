package player

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"sub-muse/internal/subsonic"
)

type State struct {
	IsPlaying bool
	Song      *subsonic.Song
	Elapsed   time.Duration
}

type Player struct {
	mu      sync.Mutex
	cmd     *exec.Cmd
	cancel  context.CancelFunc
	playing bool
	song    *subsonic.Song
	started time.Time
}

func NewPlayer() *Player {
	return &Player{}
}

func (p *Player) findPlayerBinary() string {
	if path, err := exec.LookPath("mpv"); err == nil {
		return path
	}
	if path, err := exec.LookPath("ffplay"); err == nil {
		return path
	}
	return ""
}

func (p *Player) Play(song subsonic.Song, data []byte) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.playing {
		_ = p.Stop()
	}

	binary := p.findPlayerBinary()
	if binary == "" {
		return fmt.Errorf("no player found (mpv or ffplay)")
	}

	tmpDir := os.TempDir()
	tmpFile := filepath.Join(tmpDir, fmt.Sprintf("sub-muse-%s.%s", song.ID, song.Suffix))

	if err := os.WriteFile(tmpFile, data, 0644); err != nil {
		return fmt.Errorf("failed to write temp file: %w", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	p.cancel = cancel

	args := []string{"-nodisp", "-autoexit", tmpFile}
	if strings.HasSuffix(binary, "ffplay") {
		args = []string{"-nodisp", "-autoexit", "-vn", tmpFile}
	}

	cmd := exec.CommandContext(ctx, binary, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		os.Remove(tmpFile)
		return fmt.Errorf("failed to start player: %w", err)
	}

	p.cmd = cmd
	p.playing = true
	p.song = &song
	p.started = time.Now()

	go func() {
		_ = cmd.Wait()
		os.Remove(tmpFile)
		p.mu.Lock()
		p.playing = false
		p.song = nil
		p.cancel = nil
		p.mu.Unlock()
	}()

	return nil
}

func (p *Player) Stop() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if !p.playing {
		return nil
	}

	if p.cancel != nil {
		p.cancel()
	}

	if p.cmd != nil && p.cmd.Process != nil {
		_ = p.cmd.Process.Kill()
		_ = p.cmd.Wait()
	}

	p.playing = false
	p.song = nil
	p.cancel = nil

	return nil
}

func (p *Player) GetState() State {
	p.mu.Lock()
	defer p.mu.Unlock()

	if !p.playing {
		return State{IsPlaying: false}
	}

	return State{
		IsPlaying: true,
		Song:      p.song,
		Elapsed:   time.Since(p.started),
	}
}

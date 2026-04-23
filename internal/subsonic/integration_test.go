package subsonic

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

// Integration tests hit a real Subsonic server.
// They are skipped unless the TEST_INTEGRATION env var is set.
//
// Usage:
//   TEST_INTEGRATION=1 go test -v ./internal/subsonic/... -run TestIntegration

func TestIntegration(t *testing.T) {
	if os.Getenv("TEST_INTEGRATION") == "" {
		t.Skip("set TEST_INTEGRATION=1 to run integration tests")
	}

	clientName := testClientName
	if clientName == "" {
		clientName = "sub-muse-test"
	}

	client := NewClient(testBaseURL, testUsername, testPassword, clientName)

	t.Run("Ping", func(t *testing.T) {
		err := client.Ping()
		require.NoError(t, err)
	})

	t.Run("GetMusicFolders", func(t *testing.T) {
		folders, err := client.GetMusicFolders()
		require.NoError(t, err)
		t.Logf("music folders (%d):", len(folders))
		for _, f := range folders {
			t.Logf("  [%d] %s", f.ID, f.Name)
		}
	})

	t.Run("GetArtists", func(t *testing.T) {
		artists, err := client.GetArtists()
		require.NoError(t, err)
		require.NotEmpty(t, artists, "expected at least one artist")
		t.Logf("artists: %d total, first: %s (id=%s)", len(artists), artists[0].Name, artists[0].ID)
	})

	t.Run("GetArtist", func(t *testing.T) {
		artists, err := client.GetArtists()
		require.NoError(t, err)
		require.NotEmpty(t, artists)

		artist, err := client.GetArtist(artists[0].ID)
		require.NoError(t, err)
		require.Equal(t, artists[0].Name, artist.Name)
		t.Logf("artist %q has %d album(s)", artist.Name, len(artist.Albums))
	})

	t.Run("GetAlbums", func(t *testing.T) {
		albums, err := client.GetAlbums()
		require.NoError(t, err)
		require.NotEmpty(t, albums, "expected at least one album")
		t.Logf("albums: %d total, first: %s by %s", len(albums), albums[0].Name, albums[0].Artist)
	})

	t.Run("GetAlbum", func(t *testing.T) {
		albums, err := client.GetAlbums()
		require.NoError(t, err)
		require.NotEmpty(t, albums)

		album, err := client.GetAlbum(albums[0].ID)
		require.NoError(t, err)
		require.Equal(t, albums[0].Name, album.Name)
		t.Logf("album %q has %d song(s)", album.Name, len(album.Songs))
	})

	t.Run("GetSongs", func(t *testing.T) {
		songs, err := client.GetSongs()
		require.NoError(t, err)
		require.NotEmpty(t, songs, "expected at least one song")
		t.Logf("songs: %d total, first: %q by %s", len(songs), songs[0].Title, songs[0].Artist)
	})

	t.Run("Stream", func(t *testing.T) {
		songs, err := client.GetSongs()
		require.NoError(t, err)
		require.NotEmpty(t, songs, "expected at least one song")

		data, err := client.Stream(WithID(songs[0].ID))
		require.NoError(t, err)
		require.NotEmpty(t, data, "expected non-empty audio data")
		t.Logf("streamed %d bytes for song %q", len(data), songs[0].Title)
	})

	t.Run("TokenAuth", func(t *testing.T) {
		tokenClient := NewClientWithTokenAuth(testBaseURL, testUsername, testPassword, clientName)
		err := tokenClient.Ping()
		require.NoError(t, err, "token auth ping should succeed")
	})
}

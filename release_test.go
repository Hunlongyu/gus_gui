package main

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"crypto/ed25519"
	"crypto/sha256"
	"encoding/json"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestSignedReleaseArchive(t *testing.T) {
	if os.Getenv("GUS_VERIFY_RELEASE") != "1" {
		t.Skip("requires packaged release")
	}
	decode := func(path string, dst any) {
		t.Helper()
		b, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		if err = json.Unmarshal(b, dst); err != nil {
			t.Fatal(err)
		}
	}
	var cfg struct {
		Version, Out string
		Updates      struct {
			PublicKey []byte
			GitHub    string
		}
	}
	decode("mygo.json", &cfg)
	platform := os.Getenv("GUS_RELEASE_PLATFORM")
	if platform == "" {
		platform = runtime.GOOS + "-" + runtime.GOARCH
	}
	dir := filepath.Join(cfg.Out, platform)
	var m struct {
		Version, URL, Notes string
		Size                int
		Signature           []byte
	}
	decode(filepath.Join(dir, "update-"+platform+".json"), &m)
	u, err := url.Parse(m.URL)
	if err != nil || u.Scheme != "https" || u.Host != "github.com" || !strings.HasPrefix(u.Path, "/"+cfg.Updates.GitHub+"/releases/download/v"+cfg.Version+"/") {
		t.Fatalf("unexpected update URL %q", m.URL)
	}
	if m.Version != cfg.Version || strings.TrimSpace(m.Notes) == "" || !strings.Contains(releaseHistory, m.Notes) {
		t.Fatal("release version or notes mismatch")
	}
	archive, err := os.ReadFile(filepath.Join(dir, filepath.Base(u.Path)))
	if err != nil {
		t.Fatal(err)
	}
	digest := sha256.Sum256(archive)
	if len(archive) != m.Size || len(cfg.Updates.PublicKey) != ed25519.PublicKeySize || !ed25519.Verify(cfg.Updates.PublicKey, digest[:], m.Signature) {
		t.Fatal("release archive verification failed")
	}
	// Older clients relaunch their original filename. Verify that both
	// entry points are actually in the signed update, not just the installer.
	stream, err := gzip.NewReader(bytes.NewReader(archive))
	if err != nil {
		t.Fatal(err)
	}
	defer stream.Close()
	files := tar.NewReader(stream)
	required := map[string]bool{"gus.exe": true, "GitUserSwitch.exe": true}
	if strings.HasPrefix(platform, "linux-") {
		required = map[string]bool{"gus": true, "gus.desktop": true}
	} else if strings.HasPrefix(platform, "darwin-") {
		required = map[string]bool{"gus.app/Contents/MacOS/gus": true, "gus.app/Contents/Info.plist": true}
	}
	for {
		entry, err := files.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Fatal(err)
		}
		if entry.Typeflag == tar.TypeReg && entry.Size > 0 {
			delete(required, entry.Name)
		}
	}
	if len(required) > 0 {
		t.Fatalf("missing release entry points: %v", required)
	}
}

package main

import (
	"crypto/ed25519"
	"crypto/sha256"
	"encoding/json"
	"net/url"
	"os"
	"path/filepath"
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
	dir := filepath.Join(cfg.Out, "windows-amd64")
	var m struct {
		Version, URL, Notes string
		Size                int
		Signature           []byte
	}
	decode(filepath.Join(dir, "update-windows-amd64.json"), &m)
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
}

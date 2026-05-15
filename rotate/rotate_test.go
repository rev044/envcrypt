package rotate_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yourorg/envcrypt/crypto"
	"github.com/yourorg/envcrypt/envfile"
	"github.com/yourorg/envcrypt/keystore"
	"github.com/yourorg/envcrypt/rotate"
)

func newTempStore(t *testing.T) *keystore.Store {
	t.Helper()
	dir := t.TempDir()
	store, err := keystore.New(dir)
	if err != nil {
		t.Fatalf("keystore.New: %v", err)
	}
	return store
}

func TestRotateReEncryptsEntries(t *testing.T) {
	store := newTempStore(t)
	dir := t.TempDir()

	// Generate and save the initial key.
	rawKey, err := crypto.GenerateKey()
	if err != nil {
		t.Fatalf("GenerateKey: %v", err)
	}
	encoded := crypto.EncodeKey(rawKey)
	if err := store.Save("test", encoded); err != nil {
		t.Fatalf("store.Save: %v", err)
	}

	// Encrypt two values with the initial key.
	derived, _ := crypto.DeriveKey(rawKey)
	val1, _ := crypto.Encrypt(derived, "secret1")
	val2, _ := crypto.Encrypt(derived, "secret2")

	entries := []envfile.Entry{
		{Key: "FOO", Value: val1},
		{Key: "BAR", Value: val2},
	}
	src := filepath.Join(dir, "src.env")
	if err := envfile.Write(src, entries); err != nil {
		t.Fatalf("envfile.Write: %v", err)
	}

	dst := filepath.Join(dir, "dst.env")
	rotator := rotate.New(store)
	res, err := rotator.Rotate("test", src, dst)
	if err != nil {
		t.Fatalf("Rotate: %v", err)
	}

	if res.EntriesRotated != 2 {
		t.Errorf("expected 2 rotated entries, got %d", res.EntriesRotated)
	}

	// Load the new key and verify decryption of dst.
	newKeyEncoded, err := store.Load("test")
	if err != nil {
		t.Fatalf("store.Load after rotate: %v", err)
	}
	newDerived, _ := crypto.DeriveKey(newKeyEncoded)

	rotated, err := envfile.Parse(dst)
	if err != nil {
		t.Fatalf("envfile.Parse dst: %v", err)
	}

	expected := map[string]string{"FOO": "secret1", "BAR": "secret2"}
	for _, e := range rotated {
		plain, err := crypto.Decrypt(newDerived, e.Value)
		if err != nil {
			t.Errorf("decrypt %q after rotate: %v", e.Key, err)
			continue
		}
		if plain != expected[e.Key] {
			t.Errorf("key %q: want %q, got %q", e.Key, expected[e.Key], plain)
		}
	}

	// Ensure old key no longer decrypts the new ciphertext.
	_, err = crypto.Decrypt(derived, rotated[0].Value)
	if err == nil {
		t.Error("expected old key to fail decryption after rotation")
	}

	// Clean up.
	os.Remove(src)
	os.Remove(dst)
}

func TestRotateMissingKey(t *testing.T) {
	store := newTempStore(t)
	rotator := rotate.New(store)
	_, err := rotator.Rotate("nonexistent", "/dev/null", "/dev/null")
	if err == nil {
		t.Error("expected error for missing key, got nil")
	}
}

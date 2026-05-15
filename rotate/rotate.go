// Package rotate provides key rotation logic for encrypted .env files.
package rotate

import (
	"fmt"

	"github.com/yourorg/envcrypt/crypto"
	"github.com/yourorg/envcrypt/envfile"
	"github.com/yourorg/envcrypt/keystore"
)

// Result holds the outcome of a rotation operation.
type Result struct {
	Environment string
	OldKeyID    string
	NewKeyID    string
	EntriesRotated int
}

// Rotator performs key rotation for a given environment.
type Rotator struct {
	Store *keystore.Store
}

// New creates a new Rotator backed by the provided keystore.
func New(store *keystore.Store) *Rotator {
	return &Rotator{Store: store}
}

// Rotate decrypts all entries in src using oldKey, re-encrypts them with a
// freshly generated key, persists the new key in the store, and writes the
// updated file to dst. It returns a Result describing what happened.
func (r *Rotator) Rotate(env, src, dst string) (*Result, error) {
	oldKey, err := r.Store.Load(env)
	if err != nil {
		return nil, fmt.Errorf("rotate: load old key for %q: %w", env, err)
	}

	entries, err := envfile.Parse(src)
	if err != nil {
		return nil, fmt.Errorf("rotate: parse %q: %w", src, err)
	}

	newRaw, err := crypto.GenerateKey()
	if err != nil {
		return nil, fmt.Errorf("rotate: generate new key: %w", err)
	}
	newKey, err := crypto.DeriveKey(newRaw)
	if err != nil {
		return nil, fmt.Errorf("rotate: derive new key: %w", err)
	}

	oldDerived, err := crypto.DeriveKey(oldKey)
	if err != nil {
		return nil, fmt.Errorf("rotate: derive old key: %w", err)
	}

	count := 0
	for i, e := range entries {
		if e.Comment || e.Value == "" {
			continue
		}
		plain, err := crypto.Decrypt(oldDerived, e.Value)
		if err != nil {
			return nil, fmt.Errorf("rotate: decrypt key %q: %w", e.Key, err)
		}
		cipher, err := crypto.Encrypt(newKey, plain)
		if err != nil {
			return nil, fmt.Errorf("rotate: re-encrypt key %q: %w", e.Key, err)
		}
		entries[i].Value = cipher
		count++
	}

	if err := envfile.Write(dst, entries); err != nil {
		return nil, fmt.Errorf("rotate: write %q: %w", dst, err)
	}

	newEncoded := crypto.EncodeKey(newRaw)
	if err := r.Store.Save(env, newEncoded); err != nil {
		return nil, fmt.Errorf("rotate: save new key: %w", err)
	}

	return &Result{
		Environment:    env,
		OldKeyID:       crypto.EncodeKey(oldKey)[:8],
		NewKeyID:       newEncoded[:8],
		EntriesRotated: count,
	}, nil
}

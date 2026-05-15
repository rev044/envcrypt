// Package rotate implements key rotation for envcrypt-managed .env files.
//
// A rotation operation:
//  1. Loads the current encryption key for the target environment from the
//     keystore.
//  2. Decrypts every value in the source .env file using that key.
//  3. Generates a fresh random key.
//  4. Re-encrypts all values with the new key.
//  5. Writes the updated entries to the destination file.
//  6. Persists the new key in the keystore, replacing the old one.
//
// The old key is never deleted until the new key has been successfully saved,
// minimising the window during which data could become inaccessible.
//
// Usage:
//
//	store, _ := keystore.New(keystoreDir)
//	r := rotate.New(store)
//	result, err := r.Rotate("production", "production.env.enc", "production.env.enc")
package rotate

package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var version = "0.1.0"

// rootCmd is the base command for the envcrypt CLI
var rootCmd = &cobra.Command{
	Use:   "envcrypt",
	Short: "Encrypt and manage .env files per environment",
	Long: `envcrypt is a CLI tool for encrypting, decrypting, and managing
.env files across multiple environments with key rotation support.

Example usage:
  envcrypt encrypt --env production --key ./keys/prod.key
  envcrypt decrypt --env staging --output .env.staging
  envcrypt rotate  --env production --new-key ./keys/prod-new.key`,
	Version: version,
}

// encryptCmd handles encryption of .env files
var encryptCmd = &cobra.Command{
	Use:   "encrypt",
	Short: "Encrypt a .env file for a target environment",
	RunE:  runEncrypt,
}

// decryptCmd handles decryption of .env files
var decryptCmd = &cobra.Command{
	Use:   "decrypt",
	Short: "Decrypt an encrypted .env file",
	RunE:  runDecrypt,
}

// rotateCmd handles key rotation for encrypted .env files
var rotateCmd = &cobra.Command{
	Use:   "rotate",
	Short: "Rotate encryption key for an existing encrypted .env file",
	RunE:  runRotate,
}

// keygen generates a new encryption key
var keygenCmd = &cobra.Command{
	Use:   "keygen",
	Short: "Generate a new AES-256 encryption key",
	RunE:  runKeygen,
}

func init() {
	// Encrypt flags
	encryptCmd.Flags().StringP("env", "e", "development", "Target environment name")
	encryptCmd.Flags().StringP("input", "i", ".env", "Input .env file path")
	encryptCmd.Flags().StringP("key", "k", "", "Path to encryption key file (required)")
	encryptCmd.Flags().StringP("output", "o", "", "Output file path (defaults to .env.<env>.enc)")
	_ = encryptCmd.MarkFlagRequired("key")

	// Decrypt flags
	decryptCmd.Flags().StringP("env", "e", "development", "Target environment name")
	decryptCmd.Flags().StringP("input", "i", "", "Input encrypted file path (defaults to .env.<env>.enc)")
	decryptCmd.Flags().StringP("key", "k", "", "Path to decryption key file (required)")
	decryptCmd.Flags().StringP("output", "o", ".env", "Output .env file path")
	_ = decryptCmd.MarkFlagRequired("key")

	// Rotate flags
	rotateCmd.Flags().StringP("env", "e", "development", "Target environment name")
	rotateCmd.Flags().StringP("input", "i", "", "Input encrypted file path (defaults to .env.<env>.enc)")
	rotateCmd.Flags().StringP("old-key", "", "", "Path to current key file (required)")
	rotateCmd.Flags().StringP("new-key", "", "", "Path to new key file (required)")
	_ = rotateCmd.MarkFlagRequired("old-key")
	_ = rotateCmd.MarkFlagRequired("new-key")

	// Keygen flags
	keygenCmd.Flags().StringP("output", "o", "", "Output file path for the generated key")

	// Register subcommands
	rootCmd.AddCommand(encryptCmd)
	rootCmd.AddCommand(decryptCmd)
	rootCmd.AddCommand(rotateCmd)
	rootCmd.AddCommand(keygenCmd)
}

func runEncrypt(cmd *cobra.Command, args []string) error {
	fmt.Println("encrypt: not yet implemented")
	return nil
}

func runDecrypt(cmd *cobra.Command, args []string) error {
	fmt.Println("decrypt: not yet implemented")
	return nil
}

func runRotate(cmd *cobra.Command, args []string) error {
	fmt.Println("rotate: not yet implemented")
	return nil
}

func runKeygen(cmd *cobra.Command, args []string) error {
	fmt.Println("keygen: not yet implemented")
	return nil
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

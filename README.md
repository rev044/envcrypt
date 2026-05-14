# envcrypt

> Encrypt and manage `.env` files per environment with key rotation support.

---

## Installation

```bash
go install github.com/yourusername/envcrypt@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/envcrypt.git && cd envcrypt && go build ./...
```

---

## Usage

**Encrypt a `.env` file:**

```bash
envcrypt encrypt --env production --file .env.production --key $MASTER_KEY
```

**Decrypt a `.env` file:**

```bash
envcrypt decrypt --env production --out .env --key $MASTER_KEY
```

**Rotate encryption keys:**

```bash
envcrypt rotate --env production --old-key $OLD_KEY --new-key $NEW_KEY
```

Encrypted files are stored as `.env.<environment>.enc` and can be safely committed to version control.

**Example project structure:**

```
.
├── .env.development.enc
├── .env.staging.enc
├── .env.production.enc
└── envcrypt.yaml
```

---

## Configuration

Define environments and key sources in `envcrypt.yaml`:

```yaml
environments:
  production:
    file: .env.production
    key_source: env:MASTER_KEY
  staging:
    file: .env.staging
    key_source: env:STAGING_KEY
```

---

## License

[MIT](LICENSE) © 2024 yourusername
# Security Policy

## Supported Versions

| Version | Supported |
|---------|-----------|
| latest  | Yes       |

## Reporting a Vulnerability

If you discover a security vulnerability, please report it responsibly:

1. **Do not** open a public issue
2. Use [GitHub's private vulnerability reporting](https://github.com/ggfevans/li-cli/security/advisories/new)
3. Or email the maintainer directly

## Response Timeline

- **Acknowledgment:** Within 48 hours
- **Assessment:** Within 1 week
- **Fix/disclosure:** Coordinated with reporter

## Scope

This project handles sensitive data including:

- LinkedIn session cookies and authentication tokens
- Credential storage on disk (`~/.config/li-cli/credentials.json`)
- LinkedIn API interactions over HTTPS

Security reports related to credential handling, session management, or data exposure are especially welcome.

# Security Policy

## Supported Versions

| Version | Supported          |
| ------- | ------------------ |
| Latest  | ✅ Active support   |

We actively maintain the `main` branch. Security patches are applied promptly and deployed to [zaps.ai](https://zaps.ai).

## Reporting a Vulnerability

**Please do NOT open a public GitHub issue for security vulnerabilities.**

If you discover a security vulnerability, please report it responsibly:

1. **Email**: Send details to **security@zaps.ai**
2. **Include**:
   - Description of the vulnerability
   - Steps to reproduce
   - Potential impact
   - Suggested fix (if any)

We will acknowledge your report within **48 hours** and aim to provide a fix within **7 days** for critical issues.

## Security Measures

### Branch Protection
- The `main` branch is protected with required pull request reviews
- All changes require at least 1 approving review before merging
- Stale approvals are automatically dismissed when new commits are pushed
- Force pushes to `main` are blocked
- Branch deletion is blocked

### Data Handling
- **Stateless proxy design** — Zaps does not store your conversation data
- **Temporary token mapping** — Redacted PII is cached in Redis with a strict **2-hour TTL** for response rehydration, then permanently destroyed
- **No logging of PII** — Audit logs record metadata (e.g., "EMAIL was redacted") without storing the actual sensitive values
- **In-memory processing** — All redaction and rehydration happens in-memory; nothing is written to disk

### Infrastructure
- Security headers enforced (HSTS, CSP, X-Frame-Options, etc.)
- JWT-based authentication with bcrypt password hashing
- Account lockout after repeated failed login attempts
- Rate limiting on authentication endpoints
- TOTP-based two-factor authentication available

### PII Detection
Zaps currently detects and redacts **24 categories** of sensitive data:

| Category | Examples |
| --- | --- |
| **Personal** | Email, Phone, SSN |
| **Financial** | Credit Card, Bank Account, Routing Number |
| **API Keys** | OpenAI, GitHub, Stripe, Google, AWS, Twilio |
| **Secrets** | JWT, Private Keys, Client Secrets, Generic API Keys |
| **Infrastructure** | MongoDB URIs, Azure Connection Strings, Docker Auth |

## Disclosure Policy

- We follow [responsible disclosure](https://en.wikipedia.org/wiki/Responsible_disclosure) practices
- Security researchers who report valid vulnerabilities will be credited (with permission) in our release notes
- We will not take legal action against researchers who follow this policy

## Contact

For security concerns: **security@zaps.ai**

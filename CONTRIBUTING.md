# Contributing to Zaps.ai

Thank you for your interest in contributing to Zaps! We welcome contributions of all kinds — bug fixes, new PII patterns, documentation improvements, and feature enhancements.

## Getting Started

### 1. Fork & Clone

```bash
# Fork the repo on GitHub, then:
git clone https://github.com/YOUR_USERNAME/zaps.git
cd zaps
```

### 2. Set Up Your Development Environment

**Prerequisites:**
- Docker & Docker Compose
- Go 1.21+
- Node.js 18+

```bash
# Copy environment template
cp .env.example .env

# Start all services (PostgreSQL, Redis, Backend, Frontend)
docker-compose up -d

# Run database migrations
cd backend && ./migrate
```

- **Frontend**: http://localhost:3001
- **Backend API**: http://localhost:3000

For detailed setup instructions, see [docs/development.md](docs/development.md).

### 3. Create a Feature Branch

```bash
git checkout -b feature/your-feature-name
# or
git checkout -b fix/issue-description
```

## How to Contribute

### Submitting a Pull Request

1. **Fork** the repository and create your branch from `main`
2. **Make your changes** with clear, descriptive commits
3. **Test thoroughly** — ensure existing tests pass and add new tests for new functionality
4. **Update documentation** if your change affects the API, configuration, or user-facing behavior
5. **Open a Pull Request** against `main` with:
   - A clear title describing the change
   - A description of what the PR does and why
   - Reference to any related issues (e.g., `Fixes #42`)

### Branch Protection

The `main` branch requires:
- **At least 1 approving review** before merging
- Stale approvals are dismissed on new commits
- Force pushes are blocked

This means your PR will need approval from a maintainer before it can be merged.

### Adding a New PII Pattern

One of the most impactful contributions is adding new PII detection patterns. Here's how:

1. Open `backend/services/redaction.go`
2. Add your pattern to the `SecretPatterns` map:

```go
var SecretPatterns = map[string]*regexp.Regexp{
    // ... existing patterns ...
    "YOUR_PATTERN_NAME": regexp.MustCompile(`your-regex-here`),
}
```

3. **Test your pattern** to ensure:
   - It catches the intended sensitive data
   - It does NOT cause false positives on non-sensitive data (this is critical!)
   - Context-aware patterns are preferred for ambiguous number formats (see `ROUTING_NUMBER` and `US_BANK_ACCOUNT` for examples)

4. Update the pattern count in `SECURITY.md` if applicable

### Reporting Bugs

- **Security vulnerabilities**: Please email **security@zaps.ai** (see [SECURITY.md](SECURITY.md))
- **Everything else**: Open a [GitHub Issue](https://github.com/argosautomation/zaps/issues) with:
  - Steps to reproduce
  - Expected vs actual behavior
  - Environment details (OS, Go/Node version, Docker version)

### Suggesting Features

Open a [GitHub Issue](https://github.com/argosautomation/zaps/issues) with the `enhancement` label. Include:
- The problem you're trying to solve
- Your proposed solution
- Any alternatives you've considered

## Code Style

### Go (Backend)
- Follow standard Go conventions (`gofmt`, `go vet`)
- Use descriptive variable names
- Add comments for non-obvious logic
- Handle errors explicitly — do not ignore them

### TypeScript/React (Frontend)
- Use functional components with hooks
- Follow the existing component structure in `frontend/components/`
- Use the existing CSS design system — avoid introducing new styling approaches

### Commit Messages

Use clear, descriptive commit messages:

```
feat(redaction): Add IBAN pattern for EU bank accounts
fix(proxy): Handle streaming responses with split tokens
docs: Update README with new PII categories
chore: Update Go dependencies
```

Format: `type(scope): description`

Types: `feat`, `fix`, `docs`, `chore`, `test`, `refactor`, `perf`

## Project Structure

```
zaps/
├── backend/                # Go backend (Fiber framework)
│   ├── api/                # HTTP handlers (proxy, auth, admin)
│   ├── db/                 # Database migrations and connection
│   ├── services/           # Business logic (redaction, auth, billing)
│   └── main.go             # Entry point
├── frontend/               # Next.js frontend
│   ├── app/                # Pages (App Router)
│   └── components/         # Reusable UI components
├── cli/                    # CLI tools (system tray, utilities)
├── docs/                   # Documentation
├── scripts/                # Deployment and setup scripts
├── docker-compose.yml      # Local development stack
└── deploy.sh               # Deployment script
```

## License

By contributing to Zaps, you agree that your contributions will be licensed under the [MIT License](LICENSE).

## Questions?

- Check the [Development Guide](docs/development.md) for detailed setup help
- Open a [Discussion](https://github.com/argosautomation/zaps/discussions) for general questions
- Review existing [Pull Requests](https://github.com/argosautomation/zaps/pulls) to see how others contribute

Thank you for helping make AI privacy accessible to everyone! ⚡

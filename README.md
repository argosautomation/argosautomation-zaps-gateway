<div align="center">
  <br />
  <a href="https://zaps.ai">
    <h1>⚡ Zaps Gateway</h1>
  </a>
  
  <p align="center">
    <strong>The High-Performance PII Redaction Gateway for LLMs</strong>
  </p>

  <p align="center">
    <a href="https://github.com/argosautomation/argosautomation-zaps-gateway/blob/main/LICENSE">
      <img src="https://img.shields.io/github/license/argosautomation/argosautomation-zaps-gateway?style=for-the-badge&color=blue" alt="License">
    </a>
    <a href="https://zaps.ai">
      <img src="https://img.shields.io/badge/Managed_Cloud-zaps.ai-22d3ee?style=for-the-badge" alt="Zaps.ai Cloud">
    </a>
  </p>
</div>

<br />

## 🚀 Overview

**Zaps Gateway** is an open-source, high-performance API gateway that sits between your applications and LLM providers (OpenAI, Anthropic, etc.). It automatically detects and redacts Personally Identifiable Information (PII) in real-time, ensuring your customer data never leaves your infrastructure in plain text.

> "Stop sending customer secrets to AI companies."

**Want a fully managed version?** → [zaps.ai](https://zaps.ai) — Free tier with 1,000 requests/month. No credit card required.

## ✨ Features

| Feature | Description |
| :--- | :--- |
| **⚡ Ultra-Low Latency** | Built in **Go**, adding **< 10ms** overhead to your requests. |
| **🔒 Stateless Design** | PII is redacted in-memory. Token mappings cached in Redis with a **2-hour TTL** for response rehydration, then permanently destroyed. |
| **🐳 Docker Native** | Deploy anywhere with a single container. Kubernetes ready. |
| **🛡️ Smart Redaction** | **24 PII & secret types** detected automatically — emails, phones, SSNs, credit cards, API keys, JWTs, private keys, and more. |
| **📦 Multi-Tenant** | Built-in isolation for distinct teams or customers. |
| **📊 Audit Logs** | Complete visibility into what data was redacted (without storing the data). |

## 🛠️ Quick Start

Get up and running in seconds with Docker.

```bash
docker run -p 3000:3000 \
  -e DATABASE_URL=postgres://user:pass@db:5432/zaps \
  -e REDIS_URL=redis:6379 \
  zapsai/zaps-gateway
```

Or clone and run with Docker Compose:

```bash
git clone https://github.com/argosautomation/argosautomation-zaps-gateway.git
cd argosautomation-zaps-gateway
docker-compose up -d
```

Your gateway is now running at `http://localhost:3000`.

## 📚 Documentation

- **[Deployment Guide](docs/deployment.md)** — Production setup, env vars, and security.
- **[API Reference](docs/api.md)** — Endpoints for chat, completion, and administration.
- **[Development Guide](docs/development.md)** — Building from source and local setup.
- **[Contributing](CONTRIBUTING.md)** — How to submit pull requests and add PII patterns.

## 🏗️ Architecture

```
Your App → Zaps Gateway → LLM Provider
              ↓
         PII Redacted
         in <10ms
```

Zaps acts as a transparent proxy. Swap your OpenAI base URL, and every request is automatically sanitized:

```bash
# Before: Direct to OpenAI
curl https://api.openai.com/v1/chat/completions ...

# After: Through Zaps (same API, privacy added)
curl http://localhost:3000/v1/chat/completions ...
```

## 🤝 Contributing

We welcome contributions! See [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

## 📄 License

This project is licensed under the MIT License — see [LICENSE](LICENSE) for details.

## ☁️ Managed Cloud

Don't want to self-host? **[Zaps.ai](https://zaps.ai)** offers a fully managed cloud gateway with:

- Free tier (1,000 requests/month)
- Dashboard with real-time analytics
- Team management and API key rotation
- Priority support

[**Get started for free →**](https://zaps.ai/signup)

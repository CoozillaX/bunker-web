<h1 align="center">Bunker-Web</h1>

<p align="center">
A backend service built with <b>Go, Gin, and GORM</b>, forming the server-side core of a fully separated front-end/back-end authentication system.
</p>

<p align="center">
<a href="https://liliya233.uk" target="_blank"><b>Live Demo</b></a>
</p>

---

### Backend
- Go (Golang)
- Gin — HTTP routing and middleware
- GORM — ORM with MySQL driver
- WebAuthn / FIDO2 support
- SMTP (Gmail) token-based mail delivery
- Environment-based configuration (`.env` or systemd EnvironmentFile)
- Cloudflare Turnstile Captcha

---

### Infrastructure
- AWS Lightsail
- Nginx reverse proxy
- Cloudflare DNS + Origin Certificates
- Cloudflare IP Access Control (auto-updated)
- GitHub Actions automated build + deploy pipeline

---

### Related Repositories
- [Bunker-Web (Frontend)](https://github.com/Cynic158/PhoenixAuth-Web)

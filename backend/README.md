# Live Polling Backend

Go/Gin service for creator authentication, poll lifecycle management, Redis-backed voting, and SSE live results.

Run locally with MongoDB and Redis available:

```powershell
go mod tidy
go run ./cmd
```

Configuration is documented in `.env.example`. The service exposes `/healthz` and `/readyz` plus the `/api` routes described by the technical specification.

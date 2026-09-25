# Signal Polls frontend

The frontend is a Vite + React + TypeScript single-page app for the Go live-polling API.

## Structure

```text
frontend/
├── src/
│   ├── components/   reusable shell, form, result, and loading UI
│   ├── lib/          API client, domain types, formatting, and live polling hook
│   ├── pages/        route-level screens
│   ├── App.tsx       route map
│   └── styles.css    visual system and responsive layout
├── Dockerfile        production build served by nginx
└── nginx.conf        SPA fallback for deep links
```

## Local development

```powershell
npm install
npm run dev
```

The browser client defaults to `http://localhost:8080/api`. Override it with
`VITE_API_BASE_URL` when the API is hosted elsewhere.

## Production container

From the repository root, `docker compose up --build frontend` serves the app at
`http://localhost:5173`.

# Product

<!-- impeccable:product-schema 1 -->

## Platform

web

## Stack

Vite + React + TypeScript SPA, separate from the Go/Gin backend.

## Users

Poll creators who need to publish and manage a poll, and anonymous audience members who need to vote and watch results update live.

## Product Purpose

A live polling tool that lets a creator make a poll, share a link, accept one anonymous vote per visitor, and show live results without refreshes. Success means the full creator-to-audience loop is understandable and works end to end.

## Positioning

Live results are driven by Redis-backed voting and Server-Sent Events, so the audience sees the poll working in real time rather than periodically refreshing.

## Operating Context

Creators use authenticated signup, login, dashboard, and poll creation flows. Audience members use a public shared poll URL without an account. The backend runs locally or in Docker with MongoDB and Redis.

## Capabilities and Constraints

The frontend must consume the existing Go API, preserve cookie-based authentication, support responsive desktop/mobile layouts, show loading/error/empty states, and reconnect/resync when SSE drops. No frontend deployment or public hosting is included in this pass.

## Brand Commitments

No existing brand system was supplied. The interface should feel trustworthy, direct, and visibly live without generic dashboard styling.

## Evidence on Hand

The internship brief and live polling specification define the flow, API behavior, and required live-update mechanism. No production logo, photography, testimonials, or marketing claims are available; the UI must not invent them.

## Product Principles

- Make the live mechanism visible and understandable.
- Keep the voting path anonymous and frictionless.
- Give creators clear control over poll lifecycle.
- Treat errors and reconnection as first-class states.

## Accessibility & Inclusion

Use semantic HTML, keyboard-visible focus, sufficient contrast, reduced-motion support, labeled controls, and responsive layouts.

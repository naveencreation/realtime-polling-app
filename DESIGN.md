# Live Polling Visual System

## Direction

Live editorial noticeboard: a public poll should feel like a sheet pinned to a busy wall while the orange signal marks the moments the room changes. The interface stays calm and legible for operation, with character in the ink-like type scale, ruled surfaces, and live result movement.

## Mode

Operate. The user is completing a real task: create, share, vote, or manage a poll. Familiar form behavior and clear state are more important than spectacle.

## Tokens

- Background: warm paper `#f4f0e8`; content surface `#fffdf8`; ink `#171717`; muted ink `#6e6b64`.
- Signal: orange `#e85d2a`; success green `#2f7d5a`; warning yellow `#d7a72c`; error red `#b63b34`.
- Typography: Space Grotesk for interface and headlines; IBM Plex Mono for live/status metadata.
- Radius: 14px for panels and controls; 999px only for compact status pills.
- Elevation: one soft shadow system, with rules used for structure rather than stacked cards.
- Motion: 180–240ms state transitions; result bars animate width; reduced-motion disables transforms.

## Layout

- Desktop uses a centered max-width shell with a narrow status rail and a wide task surface.
- Public poll pages lead with the question and voting action; results are the proof of the live mechanism.
- Auth and creator flows use split editorial panels without making the form compete with the task.
- Mobile collapses all rails into a single reading column and keeps actions thumb-reachable.

## Signature interaction

Accepted votes briefly pulse the live signal and animate the changed result bar from its previous width. SSE updates use the same result transition, while closed polls replace the action state with a clear final-state notice.

## Accessibility and resilience

Every form has labels, keyboard focus, inline errors, disabled/loading states, and reduced-motion handling. Realtime status is announced through a live region. Network errors explain recovery instead of exposing raw backend details.

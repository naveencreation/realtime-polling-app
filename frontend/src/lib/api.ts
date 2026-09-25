import type { ApiError, Poll, User } from './types'

const API_BASE = (import.meta.env.VITE_API_BASE_URL || '/api').replace(/\/$/, '')

async function request<T>(path: string, init?: RequestInit): Promise<T> {
  const response = await fetch(`${API_BASE}${path}`, {
    credentials: 'include',
    headers: { 'Content-Type': 'application/json', ...init?.headers },
    ...init,
  })
  const payload = await response.json().catch(() => null)
  if (!response.ok) {
    const error = payload as ApiError | null
    throw new Error(error?.message ?? 'Something went wrong. Please try again.')
  }
  return payload as T
}

export const api = {
  signup: (body: { username: string; email: string; password: string }) => request<User>('/auth/signup', { method: 'POST', body: JSON.stringify(body) }),
  login: (body: { email: string; password: string }) => request<User>('/auth/login', { method: 'POST', body: JSON.stringify(body) }),
  logout: () => request<void>('/auth/logout', { method: 'POST' }),
  createPoll: (body: { question: string; options: string[]; expiresAt?: string }) => request<Poll>('/polls', { method: 'POST', body: JSON.stringify(body) }),
  listMine: () => request<Poll[]>('/polls/mine'),
  getPoll: (id: string) => request<Poll>(`/polls/${id}`),
  vote: (id: string, optionId: string) => request<{ accepted: boolean; counts: Record<string, number> }>(`/polls/${id}/vote`, { method: 'POST', body: JSON.stringify({ optionId }) }),
  closePoll: (id: string) => request<{ id: string; status: 'closed' }>(`/polls/${id}/close`, { method: 'PATCH' }),
}

export function apiBaseForEvents() { return API_BASE }

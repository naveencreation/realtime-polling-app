import type { ApiError, Poll, User } from './types'

const API_BASE = (import.meta.env.VITE_API_BASE_URL || '/api').replace(/\/$/, '')

async function request<T>(path: string, init?: RequestInit): Promise<T> {
  const headers = new Headers(init?.headers)
  if (!headers.has('Content-Type')) {
    headers.set('Content-Type', 'application/json')
  }

  // Attach auth token if available (supports cross-origin / third-party cookie restrictions)
  const token = typeof window !== 'undefined' ? localStorage.getItem('signal_auth_token') : null
  if (token && !headers.has('Authorization')) {
    headers.set('Authorization', `Bearer ${token}`)
  }

  // Attach voter token for client-side deduplication support
  const voterToken = typeof window !== 'undefined' ? localStorage.getItem('signal_voter_token') : null
  if (voterToken && !headers.has('X-Voter-Token')) {
    headers.set('X-Voter-Token', voterToken)
  }

  const response = await fetch(`${API_BASE}${path}`, {
    credentials: 'include',
    headers,
    ...init,
  })

  // Capture voter token if returned by server
  const returnedVoterToken = response.headers.get('X-Voter-Token')
  if (returnedVoterToken && typeof window !== 'undefined') {
    localStorage.setItem('signal_voter_token', returnedVoterToken)
  }

  const payload = await response.json().catch(() => null)
  if (!response.ok) {
    const error = payload as ApiError | null
    throw new Error(error?.message ?? 'Something went wrong. Please try again.')
  }
  return payload as T
}

export const api = {
  signup: async (body: { username: string; email: string; password: string }) => {
    const res = await request<User & { token?: string }>('/auth/signup', { method: 'POST', body: JSON.stringify(body) })
    if (res?.token && typeof window !== 'undefined') localStorage.setItem('signal_auth_token', res.token)
    return res
  },
  login: async (body: { email: string; password: string }) => {
    const res = await request<User & { token?: string }>('/auth/login', { method: 'POST', body: JSON.stringify(body) })
    if (res?.token && typeof window !== 'undefined') localStorage.setItem('signal_auth_token', res.token)
    return res
  },
  logout: async () => {
    if (typeof window !== 'undefined') localStorage.removeItem('signal_auth_token')
    return request<void>('/auth/logout', { method: 'POST' })
  },
  createPoll: (body: { question: string; options: string[]; expiresAt?: string }) => request<Poll>('/polls', { method: 'POST', body: JSON.stringify(body) }),
  listMine: () => request<Poll[]>('/polls/mine'),
  getPoll: (id: string) => request<Poll>(`/polls/${id}`),
  vote: (id: string, optionId: string) => request<{ accepted: boolean; counts: Record<string, number> }>(`/polls/${id}/vote`, { method: 'POST', body: JSON.stringify({ optionId }) }),
  closePoll: (id: string) => request<{ id: string; status: 'closed' }>(`/polls/${id}/close`, { method: 'PATCH' }),
}

export function apiBaseForEvents() { return API_BASE }

export type PollStatus = 'open' | 'closed'

export type PollOption = { id: string; text: string }

export type Poll = {
  id: string
  question: string
  options: PollOption[]
  status: PollStatus
  expiresAt?: string | null
  createdAt: string
  shareUrl: string
  hasVoted?: boolean
  results?: Record<string, number>
}

export type User = { id: string; username: string }

export type ApiError = { error?: string; message?: string }

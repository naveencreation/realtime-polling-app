import { useState } from 'react'
import { Link } from 'react-router-dom'
import { formatDate, formatTimeLeft } from '../lib/format'
import type { Poll } from '../lib/types'
import { SharePoll } from './SharePoll'

interface PollCardProps {
  poll: Poll
  onClose?: (pollId: string) => Promise<void> | void
  onDelete?: (pollId: string) => Promise<void> | void
}

export function PollCard({ poll, onClose, onDelete }: PollCardProps) {
  const [busy, setBusy] = useState(false)

  const handleClose = async () => {
    if (!onClose || busy) return
    if (!window.confirm('Close this poll to new responses? Live results will remain visible.')) return
    setBusy(true)
    try {
      await onClose(poll.id)
    } finally {
      setBusy(false)
    }
  }

  const handleDelete = async () => {
    if (!onDelete || busy) return
    if (!window.confirm('Permanently delete this poll and all its responses? This cannot be undone.')) return
    setBusy(true)
    try {
      await onDelete(poll.id)
    } finally {
      setBusy(false)
    }
  }

  return (
    <article className="poll-card">
      <div className="poll-card-top">
        <div style={{ display: 'flex', alignItems: 'center', gap: '8px' }}>
          <span className={`status-chip status-${poll.status}`}>{poll.status}</span>
          {poll.status === 'open' && onClose && (
            <button
              type="button"
              className="text-button"
              onClick={handleClose}
              disabled={busy}
              style={{ fontSize: '11px', padding: '2px 8px', color: 'var(--muted)', cursor: 'pointer' }}
              title="Close voting on this poll"
            >
              {busy ? 'Closing…' : 'Close room'}
            </button>
          )}
          {poll.status === 'closed' && onDelete && (
            <button
              type="button"
              className="text-button"
              onClick={handleDelete}
              disabled={busy}
              style={{ fontSize: '11px', padding: '2px 8px', color: '#c53030', cursor: 'pointer' }}
              title="Permanently delete this closed poll"
            >
              {busy ? 'Deleting…' : 'Delete'}
            </button>
          )}
        </div>
        <span className="mono">{formatDate(poll.createdAt)}</span>
      </div>
      <h2>{poll.question}</h2>
      <div className="poll-card-bottom">
        <span className="muted">{poll.options.length} options · {formatTimeLeft(poll.expiresAt)}</span>
        <Link className="arrow-link" to={`/poll/${poll.id}`}>Open poll <span aria-hidden="true">↗</span></Link>
      </div>
      <SharePoll poll={poll} compact />
    </article>
  )
}

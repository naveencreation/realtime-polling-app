import { Link } from 'react-router-dom'
import { formatDate, formatTimeLeft } from '../lib/format'
import type { Poll } from '../lib/types'
import { SharePoll } from './SharePoll'

interface PollCardProps {
  poll: Poll
  onRequestClose?: (poll: Poll) => void
  onRequestDelete?: (poll: Poll) => void
}

export function PollCard({ poll, onRequestClose, onRequestDelete }: PollCardProps) {
  return (
    <article className={`poll-card is-${poll.status}`}>
      <div className="poll-card-top">
        <div style={{ display: 'flex', alignItems: 'center', gap: '8px' }}>
          <span className={`status-chip status-${poll.status}`}>
            {poll.status === 'open' ? 'Active' : 'Closed'}
          </span>
          {poll.status === 'open' && onRequestClose && (
            <button
              type="button"
              className="text-button"
              onClick={() => onRequestClose(poll)}
              style={{ fontSize: '11px', padding: '2px 8px', color: 'var(--muted)', cursor: 'pointer' }}
              title="Close voting on this poll"
            >
              Close poll
            </button>
          )}
          {poll.status === 'closed' && onRequestDelete && (
            <button
              type="button"
              className="text-button"
              onClick={() => onRequestDelete(poll)}
              style={{ fontSize: '11px', padding: '2px 8px', color: '#c53030', cursor: 'pointer' }}
              title="Permanently delete this closed poll"
            >
              Delete poll
            </button>
          )}
        </div>
        <span className="mono">{formatDate(poll.createdAt)}</span>
      </div>
      <h2>{poll.question}</h2>
      <div className="poll-card-bottom">
        <span className="muted">
          {poll.options.length} options · {poll.status === 'closed' ? 'Closed' : formatTimeLeft(poll.expiresAt)}
        </span>
        <Link
          className="arrow-link"
          to={`/poll/${poll.id}`}
          target="_blank"
          rel="noopener noreferrer"
        >
          {poll.status === 'closed' ? 'View results ↗' : 'View poll ↗'}
        </Link>
      </div>
      <SharePoll poll={poll} compact />
    </article>
  )
}

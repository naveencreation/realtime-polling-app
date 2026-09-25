import { Link } from 'react-router-dom'
import { formatDate, formatTimeLeft } from '../lib/format'
import type { Poll } from '../lib/types'
import { SharePoll } from './SharePoll'

export function PollCard({ poll }: { poll: Poll }) {
  return <article className="poll-card">
    <div className="poll-card-top"><span className={`status-chip status-${poll.status}`}>{poll.status}</span><span className="mono">{formatDate(poll.createdAt)}</span></div>
    <h2>{poll.question}</h2>
    <div className="poll-card-bottom"><span className="muted">{poll.options.length} options · {formatTimeLeft(poll.expiresAt)}</span><Link className="arrow-link" to={`/poll/${poll.id}`}>Open poll <span aria-hidden="true">↗</span></Link></div>
    <SharePoll poll={poll} compact />
  </article>
}

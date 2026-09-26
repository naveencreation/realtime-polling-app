import { Link } from 'react-router-dom'
import type { Poll } from '../lib/types'
import { SharePoll } from './SharePoll'

export function PollCreated({ poll, onReset }: { poll: Poll; onReset?: () => void }) {
  return (
    <div className="poll-created-layout">
      <div className="poll-created-info">
        <span className="section-label">Poll is active</span>
        <h1>
          Share your poll
          <br />
          <em>with your audience.</em>
        </h1>
        <p className="created-lede">
          Anyone with this link can vote once. Live results will update instantly on your screen as votes come in.
        </p>

        <div className="created-share-card">
          <SharePoll poll={poll} />
        </div>

        <div className="created-actions">
          <Link className="primary-button" to={`/poll/${poll.id}`}>
            Go to poll page <span aria-hidden="true">↗</span>
          </Link>
          <Link className="secondary-button" to="/dashboard">
            Back to dashboard
          </Link>
          {onReset && (
            <button type="button" className="text-button" onClick={onReset}>
              Create another poll
            </button>
          )}
        </div>
      </div>

      <div className="poll-created-preview" aria-label="Poll preview">
        <div className="created-preview-badge">
          <span className="live-mark live-live">
            <i aria-hidden="true" /> Live Audience Preview
          </span>
          <span className="mono">{poll.options.length} options</span>
        </div>
        <h3 className="preview-question">{poll.question}</h3>
        <div className="preview-options-list">
          {poll.options.map((opt, i) => (
            <div key={opt.id} className="preview-option-item">
              <span className="preview-opt-index">{String(i + 1).padStart(2, '0')}</span>
              <span className="preview-opt-text">{opt.text}</span>
            </div>
          ))}
        </div>
        <div className="preview-footnote">
          <span>Ready to accept audience votes</span>
          <span className="mono">0 votes</span>
        </div>
      </div>
    </div>
  )
}

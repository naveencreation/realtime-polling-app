import { useEffect, useState } from 'react'
import { Link, useParams } from 'react-router-dom'
import { api } from '../lib/api'
import { formatDate } from '../lib/format'
import { Loading } from '../components/Loading'
import { Shell } from '../components/Shell'
import { usePoll } from '../lib/usePoll'
import { copyText, publicPollUrl, sharePoll } from '../lib/share'

export function PollPage() {
  const { id = '' } = useParams()
  const { poll, setPoll, loading, error } = usePoll(id)
  const [selected, setSelected] = useState('')
  const [voting, setVoting] = useState(false)
  const [voteError, setVoteError] = useState('')
  const [copyStatus, setCopyStatus] = useState('')

  useEffect(() => {
    if (poll?.question) {
      document.title = `${poll.question} — Signal`
    }
  }, [poll?.question])

  if (loading) {
    return (
      <Shell minimal>
        <Loading label="Loading poll..." />
      </Shell>
    )
  }

  if (error || !poll) {
    return (
      <Shell minimal>
        <div className="not-found">
          <span className="section-label">404 Not Found</span>
          <h1>Poll Not Found</h1>
          <p>{error || 'This poll does not exist, was deleted, or the link is incorrect.'}</p>
          <Link className="primary-button" to="/">
            Back to Home
          </Link>
        </div>
      </Shell>
    )
  }

  const results = poll.results ?? {}
  const totalVotes = poll.options.reduce((sum, option) => sum + (results[option.id] ?? 0), 0)

  async function vote() {
    if (!selected) return
    setVoting(true)
    setVoteError('')
    try {
      const result = await api.vote(id, selected)
      setPoll((current) => (current ? { ...current, hasVoted: true, results: result.counts } : current))
    } catch (caught) {
      setVoteError(caught instanceof Error ? caught.message : 'Your vote could not be recorded.')
    } finally {
      setVoting(false)
    }
  }

  async function handleCopy() {
    try {
      const url = publicPollUrl(poll?.shareUrl, poll?.id || '')
      await copyText(url)
      setCopyStatus('Link copied!')
      setTimeout(() => setCopyStatus(''), 2500)
    } catch {
      setCopyStatus('Failed to copy')
      setTimeout(() => setCopyStatus(''), 2500)
    }
  }

  async function handleShare() {
    const url = publicPollUrl(poll?.shareUrl, poll?.id || '')
    try {
      const shared = await sharePoll(url, poll?.question || '')
      if (!shared) {
        await handleCopy()
      }
    } catch (caught) {
      if (caught instanceof DOMException && caught.name === 'AbortError') return
      await handleCopy()
    }
  }

  return (
    <Shell minimal>
      <div className="poll-page-container">
        {/* Unified Focused Card */}
        <div className="poll-unified-card">
          <div className="poll-card-header">
            <div className="poll-badge-row">
              <span className={`status-chip status-${poll.status}`}>
                {poll.status === 'open' ? 'Active' : 'Closed'}
              </span>
              <span className="mono">{formatDate(poll.createdAt)}</span>
            </div>

            <h1 className="poll-main-question">{poll.question}</h1>

            <p className="poll-helper-text">
              {poll.status === 'closed'
                ? 'Voting has ended. Here are the final results:'
                : poll.hasVoted
                ? 'Thank you! Your vote is counted. Live results update below:'
                : 'Choose an option below to submit your vote:'}
            </p>
          </div>

          {/* Body: Vote Form or Results */}
          {poll.status === 'open' && !poll.hasVoted ? (
            <div className="poll-voting-body">
              <div className="vote-options" role="radiogroup" aria-label="Poll options">
                {poll.options.map((option) => (
                  <label
                    className={`vote-option ${selected === option.id ? 'is-selected' : ''}`}
                    key={option.id}
                  >
                    <input
                      type="radio"
                      name="poll-option"
                      value={option.id}
                      checked={selected === option.id}
                      onChange={() => setSelected(option.id)}
                    />
                    <span className="radio-mark" aria-hidden="true" />
                    <span>{option.text}</span>
                  </label>
                ))}
              </div>
              {voteError && (
                <div className="error-banner" role="alert">
                  {voteError}
                </div>
              )}
              <button
                className="primary-button vote-submit-button"
                disabled={!selected || voting}
                onClick={vote}
              >
                {voting ? 'Submitting your vote…' : 'Submit Vote'}
              </button>
            </div>
          ) : (
            <div className="poll-results-body">
              <div className="poll-results-stats">
                <span className="results-subtitle">
                  {poll.status === 'closed' ? 'Final Results' : 'Live Results'}
                </span>
                <span className="results-total-count">
                  <strong>{totalVotes}</strong> {totalVotes === 1 ? 'vote' : 'votes'}
                </span>
              </div>

              <div className="result-list">
                {poll.options.map((option) => {
                  const count = results[option.id] ?? 0
                  const percentage = totalVotes ? Math.round((count / totalVotes) * 100) : 0
                  return (
                    <div className="result-row" key={option.id}>
                      <div className="result-meta">
                        <span className="result-label">{option.text}</span>
                        <span className="result-percent-mono">
                          {percentage}% ({count})
                        </span>
                      </div>
                      <div className="result-track">
                        <div className="result-bar" style={{ width: `${percentage}%` }} />
                      </div>
                    </div>
                  )
                })}
              </div>

              {totalVotes === 0 && (
                <p className="empty-copy">No votes were recorded for this poll.</p>
              )}
            </div>
          )}

          {/* Footer: Poll ID and Single Share Action */}
          <div className="poll-card-footer">
            <span className="mono poll-id-tag">Poll ID: {poll.id.slice(0, 8)}</span>
            <button
              type="button"
              className="share-compact-btn"
              onClick={handleShare}
              title="Share or copy poll link"
            >
              <svg
                width="13"
                height="13"
                viewBox="0 0 24 24"
                fill="none"
                stroke="currentColor"
                strokeWidth="2"
                strokeLinecap="round"
                strokeLinejoin="round"
                aria-hidden="true"
              >
                <circle cx="18" cy="5" r="3" />
                <circle cx="6" cy="12" r="3" />
                <circle cx="18" cy="19" r="3" />
                <line x1="8.59" y1="13.51" x2="15.42" y2="17.49" />
                <line x1="15.41" y1="6.51" x2="8.59" y2="10.49" />
              </svg>
              <span>{copyStatus || 'Share poll'}</span>
            </button>
          </div>
        </div>
      </div>
    </Shell>
  )
}

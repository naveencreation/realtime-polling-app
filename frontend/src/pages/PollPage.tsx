import { useState } from 'react'
import { Link, useParams } from 'react-router-dom'
import { api } from '../lib/api'
import { formatDate, formatTimeLeft } from '../lib/format'
import { Loading } from '../components/Loading'
import { LiveMark, Shell } from '../components/Shell'
import { Results } from '../components/Results'
import { usePoll } from '../lib/usePoll'
import { SharePoll } from '../components/SharePoll'

export function PollPage() {
  const { id = '' } = useParams()
  const { poll, setPoll, loading, error, connection, notice } = usePoll(id)
  const [selected, setSelected] = useState('')
  const [voting, setVoting] = useState(false)
  const [voteError, setVoteError] = useState('')

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

  return (
    <Shell minimal>
      <div className="poll-page">
        <div className="poll-context">
          <Link className="back-link" to="/">
            ← Home
          </Link>
          <LiveMark status={connection} />
        </div>

        <section className="poll-hero">
          <div className="poll-hero-meta">
            <span className={`status-chip status-${poll.status}`}>
              {poll.status === 'open' ? 'Open for voting' : 'Closed'}
            </span>
            <span className="mono">
              {formatDate(poll.createdAt)} · {formatTimeLeft(poll.expiresAt)}
            </span>
          </div>
          <h1>{poll.question}</h1>
          <p className="poll-helper">
            {poll.status === 'closed'
              ? 'This poll is now closed. Final results are shown below.'
              : poll.hasVoted
              ? 'Thank you! Your vote has been recorded. Results update live.'
              : 'Select an option below and submit your vote. Each person can vote once.'}
          </p>
          <SharePoll poll={poll} />
        </section>

        {notice && (
          <div className="notice-banner" role="status">
            {notice}
          </div>
        )}

        {poll.status === 'open' && !poll.hasVoted ? (
          <section className="vote-panel">
            <div className="section-heading">
              <span className="section-label">Cast Your Vote</span>
              <span className="mono">SINGLE CHOICE</span>
            </div>
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
            <button className="primary-button vote-button" disabled={!selected || voting} onClick={vote}>
              {voting ? 'Submitting your vote…' : 'Submit Vote'}{' '}
              <span aria-hidden="true">↗</span>
            </button>
          </section>
        ) : (
          <Results poll={poll} />
        )}

        <div className="poll-footnote">
          <span className="mono">POLL ID / {poll.id.slice(0, 8).toUpperCase()}</span>
          <span>Results update automatically in real time.</span>
        </div>
      </div>
    </Shell>
  )
}

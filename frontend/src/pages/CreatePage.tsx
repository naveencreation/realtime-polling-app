import { Link } from 'react-router-dom'
import { useEffect, useState } from 'react'
import { PollForm } from '../components/PollForm'
import { PollCreated } from '../components/PollCreated'
import { Shell } from '../components/Shell'
import type { Poll } from '../lib/types'

export function CreatePage() {
  const [createdPoll, setCreatedPoll] = useState<Poll | null>(null)

  useEffect(() => {
    document.title = 'Create Poll — Signal'
  }, [])

  return (
    <Shell>
      {createdPoll ? (
        <PollCreated poll={createdPoll} onReset={() => setCreatedPoll(null)} />
      ) : (
        <div className="create-layout">
          <div className="create-intro">
            <Link className="create-back-nav" to="/dashboard">
              <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" aria-hidden="true">
                <line x1="19" y1="12" x2="5" y2="12" />
                <polyline points="12 19 5 12 12 5" />
              </svg>
              <span>Back to dashboard</span>
            </Link>
            <h1 className="create-headline">
              Create your
              <br />
              live poll.
            </h1>
            <p className="create-lede">
              Type your question, add 2 to 6 options, and share the link with your audience to gather instant votes.
            </p>

            <div className="create-features-list">
              <div className="create-feature-item">
                <span className="create-feature-dot" />
                <span>Real-time instant live streaming results</span>
              </div>
              <div className="create-feature-item">
                <span className="create-feature-dot" />
                <span>100% anonymous — no account needed for voters</span>
              </div>
              <div className="create-feature-item">
                <span className="create-feature-dot" />
                <span>Shareable link and mobile-friendly voting</span>
              </div>
            </div>
          </div>
          <section className="form-panel">
            <PollForm onCreated={setCreatedPoll} />
          </section>
        </div>
      )}
    </Shell>
  )
}

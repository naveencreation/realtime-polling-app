import { Link } from 'react-router-dom'
import { useState } from 'react'
import { PollForm } from '../components/PollForm'
import { PollCreated } from '../components/PollCreated'
import { Shell } from '../components/Shell'
import type { Poll } from '../lib/types'

export function CreatePage() {
  const [createdPoll, setCreatedPoll] = useState<Poll | null>(null)

  return (
    <Shell>
      {createdPoll ? (
        <PollCreated poll={createdPoll} onReset={() => setCreatedPoll(null)} />
      ) : (
        <div className="create-layout">
          <div className="create-intro">
            <Link className="back-link" to="/dashboard">
              ← Back to dashboard
            </Link>
            <span className="section-label">New Poll</span>
            <h1>
              Create your
              <br />
              live poll.
            </h1>
            <p>
              Type your question, add 2 to 6 options, and share the link with your audience to start gathering instant votes.
            </p>
            <div className="create-note">
              <span className="mono">01</span>
              <span>Your audience votes anonymously without signing up. Results update in real time.</span>
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

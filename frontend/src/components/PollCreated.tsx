import { Link } from 'react-router-dom'
import type { Poll } from '../lib/types'
import { SharePoll } from './SharePoll'

export function PollCreated({ poll }: { poll: Poll }) {
  return <section className="created-panel">
    <span className="section-label">Room is live</span>
    <h1>Send the question<br /><em>into the room.</em></h1>
    <p>Anyone with this link can answer once. Results will move live as people join.</p>
    <SharePoll poll={poll} />
    <div className="created-actions"><Link className="primary-button" to={`/poll/${poll.id}`}>Open poll <span aria-hidden="true">↗</span></Link><Link className="quiet-link" to="/dashboard">Back to desk</Link></div>
  </section>
}

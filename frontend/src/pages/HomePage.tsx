import { Link } from 'react-router-dom'
import { Shell } from '../components/Shell'

export function HomePage() {
  return (
    <Shell minimal>
      <div className="home-page">
        <section className="home-hero">
          <div className="hero-copy">
            <div className="signal-line">
              <span className="signal-pulse" /> Real-time audience polling
            </div>
            <h1>
              Ask questions.<br />
              <em>See votes live.</em>
            </h1>
            <p className="hero-lede">
              Create a poll in seconds, share the link with your audience, and watch votes roll in live with zero page refreshes.
            </p>
            <div className="hero-actions">
              <Link className="primary-button" to="/signup">
                Create a poll <span aria-hidden="true">↗</span>
              </Link>
              <Link className="quiet-link" to="/login">
                Sign in to your account
              </Link>
            </div>
          </div>
          <div className="hero-demo" aria-label="Live poll preview">
            <div className="demo-top">
              <span className="status-chip status-open">open</span>
              <span className="mono">LIVE DEMO</span>
            </div>
            <p className="demo-question">Where should the team have lunch?</p>
            <div className="demo-bars">
              <div>
                <span>Somewhere new</span>
                <i style={{ width: '68%' }} />
              </div>
              <div>
                <span>The usual place</span>
                <i style={{ width: '42%' }} />
              </div>
              <div>
                <span>Order takeout</span>
                <i style={{ width: '24%' }} />
              </div>
            </div>
            <div className="demo-foot">
              <span className="live-mark live-live">
                <i />Live results
              </span>
              <span className="mono">24 votes</span>
            </div>
          </div>
        </section>
      </div>
    </Shell>
  )
}

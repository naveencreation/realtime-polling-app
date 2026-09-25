import { Link, useNavigate } from 'react-router-dom'
import { api } from '../lib/api'

export function Shell({ children, minimal = false }: { children: React.ReactNode; minimal?: boolean }) {
  const navigate = useNavigate()
  async function signOut() { await api.logout().catch(() => undefined); navigate('/') }
  return <div className="app-shell">
    <header className="site-header">
      <Link className="wordmark" to="/" aria-label="Signal Polls home"><span className="wordmark-mark">S</span><span>signal<span className="wordmark-dot">.</span></span></Link>
      {!minimal && <nav className="site-nav" aria-label="Primary navigation"><Link to="/dashboard">Dashboard</Link><Link to="/create">New poll</Link><button className="text-button" onClick={signOut}>Sign out</button></nav>}
    </header>
    <main>{children}</main>
  </div>
}

export function FieldError({ children }: { children: React.ReactNode }) { return <p className="field-error" role="alert">{children}</p> }
export function LiveMark({ status = 'live' }: { status?: 'live' | 'offline' | 'connecting' }) { return <span className={`live-mark live-${status}`}><i aria-hidden="true" />{status === 'live' ? 'Live' : status === 'connecting' ? 'Connecting' : 'Reconnecting'}</span> }

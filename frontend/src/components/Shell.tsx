import { Link, NavLink, useNavigate } from 'react-router-dom'
import { api } from '../lib/api'

export function Shell({ children, minimal = false }: { children: React.ReactNode; minimal?: boolean }) {
  const navigate = useNavigate()
  async function signOut() {
    await api.logout().catch(() => undefined)
    navigate('/')
  }
  return (
    <div className="app-shell">
      <header className="site-header">
        <Link className="wordmark" to="/" aria-label="Signal Polls home">
          <span className="wordmark-mark">S</span>
          <span>
            signal<span className="wordmark-dot">.</span>
          </span>
        </Link>
        {!minimal && (
          <nav className="site-nav" aria-label="Primary navigation">
            <NavLink
              to="/dashboard"
              className={({ isActive }) => (isActive ? 'nav-link is-active' : 'nav-link')}
            >
              Dashboard
            </NavLink>
            <NavLink
              to="/create"
              className={({ isActive }) => (isActive ? 'nav-link is-active' : 'nav-link')}
            >
              + New poll
            </NavLink>
            <button className="sign-out-btn" onClick={signOut}>
              Sign out
            </button>
          </nav>
        )}
      </header>
      <main>{children}</main>
    </div>
  )
}

export function FieldError({ children }: { children: React.ReactNode }) {
  return (
    <p className="field-error" role="alert">
      {children}
    </p>
  )
}

export function LiveMark({ status = 'live' }: { status?: 'live' | 'offline' | 'connecting' }) {
  return (
    <span className={`live-mark live-${status}`}>
      <i aria-hidden="true" />
      {status === 'live' ? 'Live' : status === 'connecting' ? 'Connecting' : 'Reconnecting'}
    </span>
  )
}

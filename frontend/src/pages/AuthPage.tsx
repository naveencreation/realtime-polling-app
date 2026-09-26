import { useState } from 'react'
import { Link, useNavigate } from 'react-router-dom'
import { api } from '../lib/api'

function EyeIcon({ visible }: { visible: boolean }) {
  if (visible) {
    return (
      <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" aria-hidden="true">
        <path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z" />
        <circle cx="12" cy="12" r="3" />
      </svg>
    )
  }
  return (
    <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" aria-hidden="true">
      <path d="M17.94 17.94A10.07 10.07 0 0 1 12 20c-7 0-11-8-11-8a18.45 18.45 0 0 1 5.06-5.94M9.9 4.24A9.12 9.12 0 0 1 12 4c7 0 11 8 11 8a18.5 18.5 0 0 1-2.16 3.19m-6.72-1.07a3 3 0 1 1-4.24-4.24" />
      <line x1="1" y1="1" x2="23" y2="23" />
    </svg>
  )
}

export function AuthPage({ mode }: { mode: 'login' | 'signup' }) {
  const isSignup = mode === 'signup'
  const navigate = useNavigate()
  const [email, setEmail] = useState('')
  const [username, setUsername] = useState('')
  const [password, setPassword] = useState('')
  const [confirmPassword, setConfirmPassword] = useState('')
  const [showPassword, setShowPassword] = useState(false)
  const [showConfirmPassword, setShowConfirmPassword] = useState(false)
  const [error, setError] = useState('')
  const [loading, setLoading] = useState(false)

  async function submit(event: React.FormEvent) {
    event.preventDefault()
    setError('')

    if (isSignup && password !== confirmPassword) {
      setError('Passwords do not match.')
      return
    }

    setLoading(true)
    try {
      if (isSignup) await api.signup({ username, email, password })
      await api.login({ email, password })
      navigate('/dashboard')
    } catch (caught) {
      setError(caught instanceof Error ? caught.message : 'Could not continue.')
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="auth-layout">
      <aside className="auth-aside">
        <div className="auth-aside-header">
          <Link className="wordmark wordmark-light" to="/" aria-label="Signal Polls home">
            <span className="wordmark-mark wordmark-mark-light">S</span>
            <span>signal<span className="wordmark-dot">.</span></span>
          </Link>
          <Link className="auth-back-link" to="/">
            ← Back to home
          </Link>
        </div>

        <div className="auth-aside-content">
          <span className="section-label">{isSignup ? 'Get Started' : 'Welcome Back'}</span>
          <h1>{isSignup ? 'Create live polls for your audience.' : 'Sign in to manage your polls.'}</h1>
          <p>
            {isSignup
              ? 'Sign up in seconds to create and share polls. Your audience can vote instantly without creating an account.'
              : 'Sign in to create new polls, check live results, and manage your voting sessions.'}
          </p>
        </div>
      </aside>

      <section className="auth-panel">
        <div className="auth-form-wrap">
          <div className="form-heading">
            <h2>{isSignup ? 'Create your account' : 'Welcome back'}</h2>
            <p className="form-subheading">
              {isSignup
                ? 'Free for creators and audiences. No credit card needed.'
                : 'Sign in to manage your polls and view live results.'}
            </p>
          </div>
          {error && <div className="error-banner" role="alert">{error}</div>}
          <form className="form-stack" onSubmit={submit}>
            {isSignup && (
              <>
                <label className="field-label" htmlFor="username">
                  Full name
                </label>
                <input
                  id="username"
                  value={username}
                  onChange={(e) => setUsername(e.target.value)}
                  placeholder="e.g. Naveen Selvan"
                  autoComplete="name"
                  required
                />
              </>
            )}
            <label className="field-label" htmlFor="email">
              Email address
            </label>
            <input
              id="email"
              type="email"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              placeholder="you@example.com"
              autoComplete="email"
              required
            />

            <label className="field-label" htmlFor="password">
              Password
            </label>
            <div className="password-input-wrap">
              <input
                id="password"
                type={showPassword ? 'text' : 'password'}
                minLength={isSignup ? 8 : undefined}
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                placeholder="At least 8 characters"
                autoComplete={isSignup ? 'new-password' : 'current-password'}
                required
              />
              <button
                type="button"
                className="password-toggle-btn"
                onClick={() => setShowPassword(!showPassword)}
                aria-label={showPassword ? 'Hide password' : 'Show password'}
                title={showPassword ? 'Hide password' : 'Show password'}
              >
                <EyeIcon visible={showPassword} />
              </button>
            </div>
            {isSignup && password.length > 0 && (
              <p className={`field-hint ${password.length >= 8 ? 'is-valid' : ''}`}>
                {password.length >= 8
                  ? '✓ 8+ characters'
                  : `At least 8 characters (${password.length}/8)`}
              </p>
            )}

            {isSignup && (
              <>
                <label className="field-label" htmlFor="confirm-password">
                  Confirm password
                </label>
                <div className="password-input-wrap">
                  <input
                    id="confirm-password"
                    type={showConfirmPassword ? 'text' : 'password'}
                    minLength={8}
                    value={confirmPassword}
                    onChange={(e) => setConfirmPassword(e.target.value)}
                    placeholder="Repeat your password"
                    autoComplete="new-password"
                    required
                  />
                  <button
                    type="button"
                    className="password-toggle-btn"
                    onClick={() => setShowConfirmPassword(!showConfirmPassword)}
                    aria-label={showConfirmPassword ? 'Hide password' : 'Show password'}
                    title={showConfirmPassword ? 'Hide password' : 'Show password'}
                  >
                    <EyeIcon visible={showConfirmPassword} />
                  </button>
                </div>
                {confirmPassword.length > 0 && (
                  <p className={`field-hint ${confirmPassword === password ? 'is-valid' : 'is-invalid'}`}>
                    {confirmPassword === password ? '✓ Passwords match' : 'Passwords do not match yet'}
                  </p>
                )}
              </>
            )}

            <button className="primary-button" disabled={loading} style={{ justifyContent: 'center' }}>
              {loading ? (isSignup ? 'Creating account…' : 'Signing in…') : isSignup ? 'Create account' : 'Sign in'}
            </button>
          </form>
          <p className="form-switch">
            {isSignup ? 'Already have an account? ' : "Don't have an account? "}
            <Link to={isSignup ? '/login' : '/signup'}>{isSignup ? 'Sign in' : 'Create one'}</Link>
          </p>
        </div>
      </section>
    </div>
  )
}

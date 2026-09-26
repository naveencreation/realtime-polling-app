import { useEffect, useState } from 'react'
import { Link, useNavigate } from 'react-router-dom'
import { api } from '../lib/api'
import type { Poll } from '../lib/types'
import { PollCard } from '../components/PollCard'
import { Loading } from '../components/Loading'
import { Shell } from '../components/Shell'

export function DashboardPage() {
  const navigate = useNavigate()
  const [polls, setPolls] = useState<Poll[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')

  useEffect(() => {
    api.listMine()
      .then(setPolls)
      .catch(() => {
        setError('Your creator session has ended.')
        navigate('/login')
      })
      .finally(() => setLoading(false))
  }, [navigate])

  const handleClose = async (pollId: string) => {
    try {
      await api.closePoll(pollId)
      setPolls(prev => prev.map(p => p.id === pollId ? { ...p, status: 'closed' } : p))
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Could not close poll.')
    }
  }

  const handleDelete = async (pollId: string) => {
    try {
      await api.deletePoll(pollId)
      setPolls(prev => prev.filter(p => p.id !== pollId))
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Could not delete poll.')
    }
  }

  return (
    <Shell>
      <div className="page-shell">
        <div className="page-heading">
          <div>
            <span className="section-label">Creator Dashboard</span>
            <h1>Your Polls</h1>
            <p>Manage your active polls and view live audience responses.</p>
          </div>
          <Link className="primary-button" to="/create">
            New poll <span aria-hidden="true">↗</span>
          </Link>
        </div>
        {error && <div className="error-banner">{error}</div>}
        {loading ? (
          <Loading label="Loading your polls…" />
        ) : polls.length === 0 ? (
          <div className="empty-state">
            <span className="empty-number">01</span>
            <h2>No polls created yet</h2>
            <p>You haven't created any polls yet. Ask a question and share the link with your audience to start collecting live votes.</p>
            <Link className="primary-button" to="/create">
              Create your first poll <span aria-hidden="true">↗</span>
            </Link>
          </div>
        ) : (
          <div className="poll-grid">
            {polls.map((poll) => (
              <PollCard
                key={poll.id}
                poll={poll}
                onClose={handleClose}
                onDelete={handleDelete}
              />
            ))}
          </div>
        )}
      </div>
    </Shell>
  )
}

import { useEffect, useState } from 'react'
import { Link, useNavigate } from 'react-router-dom'
import { api } from '../lib/api'
import type { Poll } from '../lib/types'
import { PollCard } from '../components/PollCard'
import { Loading } from '../components/Loading'
import { Shell } from '../components/Shell'
import { ConfirmModal } from '../components/ConfirmModal'

export function DashboardPage() {
  const navigate = useNavigate()
  const [polls, setPolls] = useState<Poll[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')

  const [pollToDelete, setPollToDelete] = useState<Poll | null>(null)
  const [pollToClose, setPollToClose] = useState<Poll | null>(null)
  const [modalLoading, setModalLoading] = useState(false)

  useEffect(() => {
    document.title = 'Dashboard — Signal'
    api.listMine()
      .then(setPolls)
      .catch(() => {
        setError('Your creator session has ended.')
        navigate('/login')
      })
      .finally(() => setLoading(false))
  }, [navigate])

  const handleConfirmClose = async () => {
    if (!pollToClose) return
    setModalLoading(true)
    try {
      await api.closePoll(pollToClose.id)
      setPolls(prev => prev.map(p => p.id === pollToClose.id ? { ...p, status: 'closed' } : p))
      setPollToClose(null)
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Could not close poll.')
    } finally {
      setModalLoading(false)
    }
  }

  const handleConfirmDelete = async () => {
    if (!pollToDelete) return
    setModalLoading(true)
    try {
      await api.deletePoll(pollToDelete.id)
      setPolls(prev => prev.filter(p => p.id !== pollToDelete.id))
      setPollToDelete(null)
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Could not delete poll.')
    } finally {
      setModalLoading(false)
    }
  }

  return (
    <Shell>
      <div className="page-shell">
        <div className="page-heading">
          <div>
            <h1>Your Polls</h1>
            <p>Track votes and manage your polls in real time.</p>
          </div>
          <Link className="primary-button" to="/create">
            + Create poll
          </Link>
        </div>
        {error && <div className="error-banner">{error}</div>}
        {loading ? (
          <Loading label="Loading your polls…" />
        ) : polls.length === 0 ? (
          <div className="empty-state">
            <h2>No polls created yet</h2>
            <p>You haven't created any polls yet. Ask a question and share the link with your audience to start collecting live votes.</p>
            <Link className="primary-button" to="/create">
              + Create your first poll
            </Link>
          </div>
        ) : (
          <div className="poll-grid">
            {polls.map((poll) => (
              <PollCard
                key={poll.id}
                poll={poll}
                onRequestClose={setPollToClose}
                onRequestDelete={setPollToDelete}
              />
            ))}
          </div>
        )}
      </div>

      <ConfirmModal
        isOpen={!!pollToDelete}
        title="Delete poll?"
        message={`Are you sure you want to permanently delete "${pollToDelete?.question}"? All votes and responses will be permanently removed. This action cannot be undone.`}
        confirmText="Delete poll"
        cancelText="Cancel"
        isDestructive={true}
        loading={modalLoading}
        onConfirm={handleConfirmDelete}
        onCancel={() => !modalLoading && setPollToDelete(null)}
      />

      <ConfirmModal
        isOpen={!!pollToClose}
        title="Close voting on this poll?"
        message={`Voting on "${pollToClose?.question}" will stop immediately. Existing responses and final results will remain visible to you and your audience.`}
        confirmText="Close poll"
        cancelText="Cancel"
        isDestructive={false}
        loading={modalLoading}
        onConfirm={handleConfirmClose}
        onCancel={() => !modalLoading && setPollToClose(null)}
      />
    </Shell>
  )
}

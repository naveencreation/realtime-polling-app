import { useCallback, useEffect, useRef, useState } from 'react'
import { api, apiBaseForEvents } from './api'
import type { Poll } from './types'

export function usePoll(id: string) {
  const [poll, setPoll] = useState<Poll | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')
  const [connection, setConnection] = useState<'connecting' | 'live' | 'offline'>('connecting')
  const [notice, setNotice] = useState('')
  const retryTimer = useRef<number | undefined>(undefined)

  const refresh = useCallback(async () => {
    try {
      const next = await api.getPoll(id)
      setPoll(next)
      setError('')
      return next
    } catch (caught) {
      setError(caught instanceof Error ? caught.message : 'This poll could not be loaded.')
      return null
    } finally {
      setLoading(false)
    }
  }, [id])

  useEffect(() => { void refresh() }, [refresh])

  useEffect(() => {
    let disposed = false
    let source: EventSource | undefined
    const connect = () => {
      if (disposed) return
      setConnection('connecting')
      source = new EventSource(`${apiBaseForEvents()}/polls/${id}/stream`, { withCredentials: true })
      source.addEventListener('results', (event) => {
        const data = JSON.parse((event as MessageEvent).data) as { counts: Record<string, number> }
        setPoll((current) => current ? { ...current, results: data.counts } : current)
        setConnection('live')
      })
      source.addEventListener('closed', () => {
        setPoll((current) => current ? { ...current, status: 'closed' } : current)
        setNotice('This poll is now closed. These are the final results.')
        setConnection('live')
        source?.close()
      })
      source.onopen = () => setConnection('live')
      source.onerror = () => {
        setConnection('offline')
        source?.close()
        if (!disposed) {
          retryTimer.current = window.setTimeout(() => { void refresh(); connect() }, 2500)
        }
      }
    }
    connect()
    return () => { disposed = true; source?.close(); if (retryTimer.current) window.clearTimeout(retryTimer.current) }
  }, [id, refresh])

  return { poll, setPoll, loading, error, connection, notice, refresh }
}

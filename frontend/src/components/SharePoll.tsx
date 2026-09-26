import { useState } from 'react'
import type { Poll } from '../lib/types'
import { copyText, publicPollUrl, sharePoll } from '../lib/share'

export function SharePoll({ poll, compact = false }: { poll: Poll; compact?: boolean }) {
  const [status, setStatus] = useState('')
  const url = publicPollUrl(poll.shareUrl, poll.id)

  async function copy() {
    try {
      await copyText(url)
      setStatus('Link copied!')
      setTimeout(() => setStatus(''), 2500)
    } catch {
      setStatus('Failed to copy')
      setTimeout(() => setStatus(''), 2500)
    }
  }

  async function share() {
    try {
      const usedNativeShare = await sharePoll(url, poll.question)
      if (!usedNativeShare) await copy()
    } catch (caught) {
      if (caught instanceof DOMException && caught.name === 'AbortError') return
      await copy()
    }
  }

  if (compact) {
    return (
      <div className="share-compact-wrap">
        <button type="button" className="share-compact-btn" onClick={share} title="Share or copy poll link">
          <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" aria-hidden="true">
            <circle cx="18" cy="5" r="3"/>
            <circle cx="6" cy="12" r="3"/>
            <circle cx="18" cy="19" r="3"/>
            <line x1="8.59" y1="13.51" x2="15.42" y2="17.49"/>
            <line x1="15.41" y1="6.51" x2="8.59" y2="10.49"/>
          </svg>
          <span>{status || 'Share poll'}</span>
        </button>
      </div>
    )
  }

  return (
    <div className="share-block">
      <label className="share-label" htmlFor={`share-${poll.id}`}>Shareable Poll Link</label>
      <div className="share-row">
        <input id={`share-${poll.id}`} className="share-input" value={url} readOnly aria-label="Public poll link" onFocus={(event) => event.currentTarget.select()} />
        <button type="button" className="secondary-button" onClick={copy}>Copy Link</button>
        <button type="button" className="share-button" onClick={share}>Share</button>
      </div>
      {status && <span className="share-status" role="status">{status}</span>}
    </div>
  )
}

import { useState } from 'react'
import type { Poll } from '../lib/types'
import { copyText, publicPollUrl, sharePoll } from '../lib/share'

export function SharePoll({ poll, compact = false }: { poll: Poll; compact?: boolean }) {
  const [status, setStatus] = useState('')
  const url = publicPollUrl(poll.shareUrl, poll.id)

  async function copy() {
    try {
      await copyText(url)
      setStatus('Link copied to clipboard')
    } catch {
      setStatus('Failed to copy link')
    }
  }

  async function share() {
    try {
      const usedNativeShare = await sharePoll(url, poll.question)
      if (!usedNativeShare) await copy()
      else setStatus('Sharing options opened')
    } catch (caught) {
      if (caught instanceof DOMException && caught.name === 'AbortError') return
      setStatus('Failed to copy link')
    }
  }

  return <div className={`share-block ${compact ? 'share-compact' : ''}`}>
    {!compact && <label className="share-label" htmlFor={`share-${poll.id}`}>Shareable Poll Link</label>}
    <div className="share-row">
      <input id={`share-${poll.id}`} className="share-input" value={url} readOnly aria-label="Public poll link" onFocus={(event) => event.currentTarget.select()} />
      <button type="button" className="secondary-button" onClick={copy}>Copy Link</button>
      <button type="button" className="share-button" onClick={share}>Share</button>
    </div>
    {status && <span className="share-status" role="status">{status}</span>}
  </div>
}

import type { Poll } from '../lib/types'

export function Results({ poll }: { poll: Poll }) {
  const results = poll.results ?? {}
  const total = poll.options.reduce((sum, option) => sum + (results[option.id] ?? 0), 0)
  return <section className="results-panel" aria-label="Live poll results">
    <div className="results-heading"><div><span className="section-label">Live Standings</span><h2>Poll Results</h2></div><span className="total-votes"><strong>{total}</strong> {total === 1 ? 'vote' : 'votes'}</span></div>
    <div className="result-list">{poll.options.map((option) => {
      const count = results[option.id] ?? 0
      const percentage = total ? Math.round((count / total) * 100) : 0
      return <div className="result-row" key={option.id}><div className="result-meta"><span>{option.text}</span><span className="mono">{percentage}% · {count}</span></div><div className="result-track"><div className="result-bar" style={{ width: `${percentage}%` }} /></div></div>
    })}</div>
    {total === 0 && <p className="empty-copy">No votes have been cast yet. Be the first to vote!</p>}
  </section>
}

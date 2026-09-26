import { useState } from 'react'
import { api } from '../lib/api'
import { FieldError } from './Shell'
import type { Poll } from '../lib/types'

export function PollForm({ onCreated }: { onCreated: (poll: Poll) => void }) {
  const [question, setQuestion] = useState(''); const [options, setOptions] = useState(['', '']); const [expiresAt, setExpiresAt] = useState(''); const [error, setError] = useState(''); const [saving, setSaving] = useState(false)
  function updateOption(index: number, value: string) { setOptions((current) => current.map((option, i) => i === index ? value : option)) }
  function addOption() { if (options.length < 6) setOptions((current) => [...current, '']) }
  function removeOption(index: number) { if (options.length > 2) setOptions((current) => current.filter((_, i) => i !== index)) }
  async function submit(event: React.FormEvent) { event.preventDefault(); setError(''); setSaving(true); try { const poll = await api.createPoll({ question, options, ...(expiresAt ? { expiresAt: new Date(expiresAt).toISOString() } : {}) }); onCreated(poll) } catch (caught) { setError(caught instanceof Error ? caught.message : 'Could not create this poll.') } finally { setSaving(false) } }
  return (
    <form className="form-stack" onSubmit={submit}>
      {error && <div className="error-banner" role="alert">{error}</div>}
      <label className="field-label" htmlFor="question">
        Poll Question <span>Required</span>
      </label>
      <textarea
        id="question"
        value={question}
        onChange={(event) => setQuestion(event.target.value)}
        placeholder="e.g. Which project should we prioritize for next sprint?"
        maxLength={200}
        rows={3}
        required
      />
      <div className="field-heading">
        <label className="field-label" htmlFor="option-0">
          Voting Options <span>2–6 choices</span>
        </label>
        <span className="mono">{options.length}/6</span>
      </div>
      <div className="option-fields">
        {options.map((option, index) => (
          <div className="option-field" key={index}>
            <span className="option-index">0{index + 1}</span>
            <input
              id={`option-${index}`}
              value={option}
              onChange={(event) => updateOption(index, event.target.value)}
              placeholder={`Option ${index + 1}`}
              required
            />
            <button
              type="button"
              className="remove-button"
              onClick={() => removeOption(index)}
              disabled={options.length <= 2}
              aria-label={`Remove option ${index + 1}`}
            >
              ×
            </button>
          </div>
        ))}
      </div>
      {options.length < 6 && (
        <button type="button" className="add-option" onClick={addOption}>
          + Add another option
        </button>
      )}
      <label className="field-label" htmlFor="expires">
        Auto-close deadline <span>Optional</span>
      </label>
      <input
        id="expires"
        type="datetime-local"
        value={expiresAt}
        onChange={(event) => setExpiresAt(event.target.value)}
      />
      <button className="primary-button" disabled={saving}>
        {saving ? 'Creating poll…' : 'Create & Share Poll'} <span aria-hidden="true">↗</span>
      </button>
      {error && <FieldError>Check the form and try again.</FieldError>}
    </form>
  )
}

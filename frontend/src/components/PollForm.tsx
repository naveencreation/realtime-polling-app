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

      <div className="form-group">
        <div className="field-header">
          <label className="field-label-text" htmlFor="question">
            Poll Question
          </label>
          <span className="field-badge-required">Required</span>
        </div>
        <textarea
          id="question"
          value={question}
          onChange={(event) => setQuestion(event.target.value)}
          placeholder="e.g. Which project should we prioritize for next sprint?"
          maxLength={200}
          rows={3}
          required
        />
        <div className="field-char-count">
          <span>{question.length}/200 characters</span>
        </div>
      </div>

      <div className="form-group">
        <div className="field-header">
          <label className="field-label-text" htmlFor="option-0">
            Voting Options
          </label>
          <span className="mono">{options.length} / 6</span>
        </div>

        <div className="option-fields">
          {options.map((option, index) => (
            <div className="option-field" key={index}>
              <span className="option-index">{String(index + 1).padStart(2, '0')}</span>
              <input
                id={`option-${index}`}
                value={option}
                onChange={(event) => updateOption(index, event.target.value)}
                placeholder={`Option ${index + 1}`}
                required
              />
              {options.length > 2 ? (
                <button
                  type="button"
                  className="remove-button"
                  onClick={() => removeOption(index)}
                  aria-label={`Remove option ${index + 1}`}
                  title="Remove this option"
                >
                  <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2.5" strokeLinecap="round" strokeLinejoin="round" aria-hidden="true">
                    <line x1="18" y1="6" x2="6" y2="18" />
                    <line x1="6" y1="6" x2="18" y2="18" />
                  </svg>
                </button>
              ) : (
                <span className="remove-placeholder" aria-hidden="true" />
              )}
            </div>
          ))}
        </div>

        {options.length < 6 && (
          <div className="add-option-wrap">
            <button type="button" className="add-option-btn" onClick={addOption}>
              <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2.5" strokeLinecap="round" strokeLinejoin="round" aria-hidden="true">
                <line x1="12" y1="5" x2="12" y2="19" />
                <line x1="5" y1="12" x2="19" y2="12" />
              </svg>
              <span>Add another option</span>
            </button>
          </div>
        )}
      </div>

      <div className="form-group">
        <div className="field-header">
          <label className="field-label-text" htmlFor="expires">
            Voting Deadline
          </label>
          <span className="field-badge-info">Optional</span>
        </div>
        <input
          id="expires"
          type="datetime-local"
          value={expiresAt}
          onChange={(event) => setExpiresAt(event.target.value)}
        />
        <span className="field-hint-text">Leave blank to keep voting open until you close it manually.</span>
      </div>

      <button className="primary-button create-submit-btn" disabled={saving}>
        {saving ? 'Creating poll…' : 'Create & share poll'}
      </button>
      {error && <FieldError>Check the form and try again.</FieldError>}
    </form>
  )
}

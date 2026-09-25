export function publicPollUrl(shareUrl: string, pollId: string) {
  if (typeof window !== 'undefined' && window.location?.origin) {
    return `${window.location.origin}/poll/${pollId}`
  }
  return shareUrl || `/poll/${pollId}`
}

export async function copyText(value: string) {
  if (navigator.clipboard?.writeText) {
    await navigator.clipboard.writeText(value)
    return
  }

  const textarea = document.createElement('textarea')
  textarea.value = value
  textarea.setAttribute('readonly', '')
  textarea.style.position = 'fixed'
  textarea.style.opacity = '0'
  document.body.appendChild(textarea)
  textarea.select()
  const copied = document.execCommand('copy')
  textarea.remove()
  if (!copied) throw new Error('Copy is unavailable in this browser.')
}

export async function sharePoll(url: string, question: string) {
  if (!navigator.share) return false
  await navigator.share({ title: 'Signal Poll', text: question, url })
  return true
}

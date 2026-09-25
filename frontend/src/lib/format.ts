export function formatDate(value: string) {
  return new Intl.DateTimeFormat(undefined, { month: 'short', day: 'numeric', year: 'numeric' }).format(new Date(value))
}

export function formatTimeLeft(value?: string | null) {
  if (!value) return 'No closing time'
  const difference = new Date(value).getTime() - Date.now()
  if (difference <= 0) return 'Closing now'
  const minutes = Math.round(difference / 60000)
  if (minutes < 60) return `${minutes} min left`
  const hours = Math.round(minutes / 60)
  if (hours < 24) return `${hours} hr left`
  return `${Math.round(hours / 24)} days left`
}

export function formatTime(value?: string | null): string {
  if (!value) return '-'
  return value.replace('T', ' ').slice(0, 19)
}

export function formatMoney(value?: number | null): string {
  if (value === null || value === undefined) return '¥0.00'
  return `¥${Number(value).toFixed(2)}`
}

export function formatDuration(minutes?: number | null): string {
  if (!minutes) return '0分钟'
  const m = Number(minutes)
  if (m < 60) return `${m}分钟`
  return `${Math.floor(m / 60)}小时${m % 60}分钟`
}

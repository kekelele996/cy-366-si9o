import { get } from '@/utils/request'

export interface AuditLog {
  id: number
  user_id: number
  username: string
  action: string
  module: string
  detail: string
  ip: string
  created_at: string
}

export function listAudits(params: { page: number; page_size: number; module?: string; action?: string }) {
  return get<{ list: AuditLog[]; total: number }>('/audits', params)
}

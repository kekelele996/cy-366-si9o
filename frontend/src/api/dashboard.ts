import { get } from '@/utils/request'

export interface DashboardSummary {
  station_total: number
  station_idle: number
  station_using: number
  station_fault: number
  station_reserved: number
  active_session: number
  member_total: number
  tournament_running: number
}

export function getSummary() {
  return get<DashboardSummary>('/dashboard/summary')
}

import { get, post, put, del } from '@/utils/request'

export interface Tournament {
  id: number
  name: string
  game_type: string
  description: string
  register_start: string | null
  register_end: string | null
  start_time: string | null
  max_teams: number
  status: string
}

export interface Team {
  id: number
  tournament_id: number
  name: string
  leader_id: number
  member_count: number
}

export interface Registration {
  id: number
  tournament_id: number
  user_id: number
  team_id: number
  mode: string
  group_no: number
  status: string
}

export interface Match {
  id: number
  tournament_id: number
  round: number
  group_no: number
  team_a_id: number
  team_b_id: number
  score_a: number
  score_b: number
  winner_id: number
  status: string
}

export function listTournaments(params: { page: number; page_size: number; status?: string }) {
  return get<{ list: Tournament[]; total: number }>('/tournaments', params)
}

export function getTournament(id: number) {
  return get<Tournament>(`/tournaments/${id}`)
}

export function createTournament(data: Partial<Tournament>) {
  return post<Tournament>('/tournaments', data)
}

export function updateTournament(id: number, data: Partial<Tournament>) {
  return put<Tournament>(`/tournaments/${id}`, data)
}

export function deleteTournament(id: number) {
  return del(`/tournaments/${id}`)
}

export function registerTournament(id: number, data: { mode: string; team_id?: number }) {
  return post<Registration>(`/tournaments/${id}/register`, data)
}

export function listRegistrations(id: number) {
  return get<Registration[]>(`/tournaments/${id}/registrations`)
}

export function drawGroups(id: number) {
  return post<{ groups: number; registrations: number }>(`/tournaments/${id}/draw`)
}

export function listMatches(id: number) {
  return get<Match[]>(`/tournaments/${id}/matches`)
}

export function submitMatchResult(id: number, data: { score_a: number; score_b: number; winner_id?: number }) {
  return post<Match>(`/matches/${id}/result`, data)
}

export function listMyTeams() {
  return get<Team[]>('/teams/mine')
}

export function createTeam(data: { name: string; member_count?: number }) {
  return post<Team>('/teams', data)
}

import { get, post } from '@/utils/request'

export interface UserItem {
  id: number
  username: string
  nickname: string
  phone: string
  role: string
  balance: number
  status: string
}

export interface LoginResult {
  token: string
  user: UserItem
}

export function login(data: { username: string; password: string }) {
  return post<LoginResult>('/auth/login', data)
}

export function register(data: { username: string; password: string; nickname: string; phone?: string }) {
  return post<LoginResult>('/auth/register', data)
}

export function fetchProfile() {
  return get<UserItem>('/auth/profile')
}

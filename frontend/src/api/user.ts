import { get, post, put, del } from '@/utils/request'
import type { UserItem } from './auth'

export function listUsers(params: { page: number; page_size: number }) {
  return get<{ list: UserItem[]; total: number; page: number; page_size: number }>('/users', params)
}

export function getUser(id: number) {
  return get<UserItem>(`/users/${id}`)
}

export function createUser(data: { username: string; password: string; nickname: string; phone?: string; role: string }) {
  return post<UserItem>('/users', data)
}

export function updateUser(id: number, data: Partial<UserItem>) {
  return put<UserItem>(`/users/${id}`, data)
}

export function deleteUser(id: number) {
  return del(`/users/${id}`)
}

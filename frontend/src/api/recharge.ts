import { get, post, del } from '@/utils/request'

export interface Recharge {
  id: number
  user_id: number
  amount: number
  payment_method: string
  remark: string
  created_at: string
}

export interface PackageOrder {
  id: number
  order_no: string
  package_id: number
  package_name: string
  amount: number
  hours: number
  payment_method: string
  status: string
  created_at: string
}

export interface TimePackage {
  id: number
  name: string
  hours: number
  price: number
  valid_days: number
  status: string
}

export function recharge(data: { user_id: number; amount: number; payment_method: string; remark?: string }) {
  return post('/recharges', data)
}

export function listMyRecharges(params: { page: number; page_size: number }) {
  return get<{ list: Recharge[]; total: number }>('/recharges/mine', params)
}

export function buyPackage(data: { package_id: number; payment_method: string }) {
  return post<PackageOrder>('/recharges/packages', data)
}

export function listMyOrders(params: { page: number; page_size: number }) {
  return get<{ list: PackageOrder[]; total: number }>('/recharges/packages/orders', params)
}

export function listPackages(params: { page: number; page_size: number }) {
  return get<{ list: TimePackage[]; total: number }>('/packages', params)
}

export function listActivePackages() {
  return get<TimePackage[]>('/packages/active')
}

export function createPackage(data: Partial<TimePackage>) {
  return post<TimePackage>('/packages', data)
}

export function deletePackage(id: number) {
  return del('/packages/' + id)
}

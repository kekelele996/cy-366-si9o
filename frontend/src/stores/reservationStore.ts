import { defineStore } from 'pinia'
import { listReservations, type Reservation } from '@/api/reservation'

export const useReservationStore = defineStore('reservation', () => {
  async function fetchList(params: { page: number; page_size: number; status?: string }) {
    return listReservations(params)
  }
  function reserveByUser(list: Reservation[], userId: number): Reservation[] {
    return list.filter((r) => r.user_id === userId)
  }
  return { fetchList, reserveByUser }
})

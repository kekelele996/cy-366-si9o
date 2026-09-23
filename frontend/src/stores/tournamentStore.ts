import { defineStore } from 'pinia'
import { ref } from 'vue'
import { listTournaments, type Tournament } from '@/api/tournament'

export const useTournamentStore = defineStore('tournament', () => {
  const tournaments = ref<Tournament[]>([])

  async function loadList(params: { page: number; page_size: number; status?: string }) {
    const data = await listTournaments(params)
    tournaments.value = data.list
    return data
  }

  return { tournaments, loadList }
})

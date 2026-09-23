<template>
  <div class="dashboard">
    <van-notice-bar left-icon="volume-o" text="机位状态每 5 秒通过 WebSocket 实时推送，刷新自动同步" />
    <div class="stats">
      <div class="stat-card" v-for="s in statItems" :key="s.label">
        <div class="stat-value" :style="{ color: s.color }">{{ s.value }}</div>
        <div class="stat-label">{{ s.label }}</div>
      </div>
    </div>
    <van-cell-group inset title="机位实时状态">
      <van-cell v-for="st in stationStore.stations" :key="st.id" :title="`${st.name}（${st.area}）`" :label="`¥${st.price_per_hour}/小时`">
        <template #value><StatusBadge kind="station" :status="st.status" /></template>
      </van-cell>
    </van-cell-group>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import StatusBadge from '@/components/StatusBadge.vue'
import { useStationStore } from '@/stores/stationStore'
import { getSummary, type DashboardSummary } from '@/api/dashboard'

const stationStore = useStationStore()
const summary = ref<DashboardSummary>({ station_total: 0, station_idle: 0, station_using: 0, station_fault: 0, station_reserved: 0, active_session: 0, member_total: 0, tournament_running: 0 })

const statItems = computed(() => [
  { label: '机位总数', value: summary.value.station_total, color: '#1989fa' },
  { label: '空闲机位', value: summary.value.station_idle, color: '#07c160' },
  { label: '使用中', value: summary.value.station_using, color: '#ff976a' },
  { label: '进行中上机', value: summary.value.active_session, color: '#7232dd' },
  { label: '会员总数', value: summary.value.member_total, color: '#ee0a24' },
  { label: '进行中赛事', value: summary.value.tournament_running, color: '#ffc300' },
])

onMounted(async () => {
  try {
    summary.value = await getSummary()
  } catch { /* 看板统计失败不影响页面 */ }
  await stationStore.loadAll()
})
</script>

<style scoped>
.stats { display: grid; grid-template-columns: repeat(3, 1fr); gap: 10px; margin: 12px 0; }
.stat-card { background: #fff; border-radius: 10px; padding: 14px 8px; text-align: center; }
.stat-value { font-size: 22px; font-weight: 600; }
.stat-label { font-size: 12px; color: #969799; margin-top: 4px; }
</style>

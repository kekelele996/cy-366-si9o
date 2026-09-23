<template>
  <div class="sessions-page">
    <van-tabs v-model:active="tab">
      <van-tab title="上机记录" name="list">
        <van-cell-group inset title="上机/下机操作">
          <van-field v-model="startForm.station_id" type="number" label="机位ID" placeholder="输入机位ID" />
          <van-field v-model="startForm.game_type" label="游戏类型" placeholder="lol/csgo/kog/other" />
        </van-cell-group>
        <div class="submit-btn"><van-button round block type="primary" @click="start">开机上机</van-button></div>

        <van-cell-group inset>
          <van-cell v-for="s in list" :key="s.id" :title="`上机 #${s.id} · 机位 ${s.station_id} · 会员 ${s.user_id}`" :label="`${formatTime(s.start_time)} ~ ${formatTime(s.end_time)} · ${formatDuration(s.duration_minutes)}`">
            <template #value>
              <StatusBadge kind="session" :status="s.status" />
              <van-button v-if="s.status === 'active'" size="mini" type="warning" plain class="op-btn" @click="renew(s)">续费</van-button>
              <van-button v-if="s.status === 'active'" size="mini" type="danger" plain class="op-btn" @click="end(s)">下机</van-button>
            </template>
          </van-cell>
        </van-cell-group>
        <van-pagination v-model="page" :total-items="total" :items-per-page="pageSize" @change="load" />
      </van-tab>

      <van-tab title="时长排行榜" name="rank">
        <van-dropdown-menu>
          <van-dropdown-item v-model="period" :options="periodOptions" @change="loadRank" />
        </van-dropdown-menu>
        <van-cell-group inset title="累计上机时长 TOP10">
          <van-cell v-for="r in ranks" :key="r.rank" :title="`#${r.rank} ${r.nickname || r.username}`" :label="`会员ID ${r.user_id}`" :value="formatDuration(r.total_minutes)" />
        </van-cell-group>
      </van-tab>
    </van-tabs>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { showSuccessToast, showToast } from 'vant'
import StatusBadge from '@/components/StatusBadge.vue'
import { listSessions, startSession, renewSession, endSession, getRank, type Session, type RankItem } from '@/api/session'
import { formatTime, formatDuration } from '@/utils/format'

const tab = ref('list')
const list = ref<Session[]>([])
const ranks = ref<RankItem[]>([])
const total = ref(0)
const page = ref(1)
const pageSize = 10
const period = ref('week')
const periodOptions = [
  { text: '近7天', value: 'week' },
  { text: '近30天', value: 'month' },
  { text: '近1天', value: 'day' },
]
const startForm = ref({ station_id: '', game_type: 'other' })

async function load() {
  const data = await listSessions({ page: page.value, page_size: pageSize })
  list.value = data.list
  total.value = data.total
}

async function loadRank() {
  ranks.value = await getRank({ period: period.value, limit: 10 })
}

async function start() {
  const stationId = Number(startForm.value.station_id)
  if (!stationId) {
    showToast('请输入机位ID')
    return
  }
  const s = await startSession({ station_id: stationId, game_type: startForm.value.game_type || 'other' })
  showSuccessToast(`开机成功，上机记录 #${s.id}`)
  startForm.value = { station_id: '', game_type: 'other' }
  load()
}

async function renew(s: Session) {
  try {
    await renewSession(s.id, 60)
    showSuccessToast('续费1小时成功')
    load()
  } catch { /* 余额不足提示 */ }
}

async function end(s: Session) {
  await endSession(s.id)
  showSuccessToast('下机结算成功')
  load()
  loadRank()
}

onMounted(() => {
  load()
  loadRank()
})
</script>

<style scoped>
.submit-btn { margin: 12px 16px; }
.op-btn { margin-left: 6px; }
</style>

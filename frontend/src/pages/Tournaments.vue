<template>
  <div class="tournaments-page">
    <van-cell-group inset>
      <van-cell v-for="t in list" :key="t.id" :title="t.name" :label="`${GAME_TYPE_TEXT[t.game_type] || t.game_type} · 最多${t.max_teams}队`" @click="detail(t)">
        <template #value><StatusBadge kind="tournament" :status="t.status" /></template>
      </van-cell>
    </van-cell-group>
    <van-pagination v-model="page" :total-items="total" :items-per-page="pageSize" @change="load" />

    <van-action-sheet v-model:show="showDetail" :title="selected?.name || '赛事详情'" :actions="detailActions" @select="onAction" cancel-text="关闭" />

    <van-dialog v-model:show="showCreate" title="创建赛事" show-cancel-button @confirm="submitCreate">
      <van-form>
        <van-cell-group inset>
          <van-field v-model="createForm.name" label="赛事名称" placeholder="如 峡谷之巅杯" />
          <van-field v-model="createForm.game_type" label="游戏类型" placeholder="lol/csgo/kog/other" />
          <van-field v-model="createForm.max_teams" type="number" label="最大队伍数" />
        </van-cell-group>
      </van-form>
    </van-dialog>

    <van-dialog v-model:show="showRegister" title="赛事报名" show-cancel-button @confirm="submitRegister">
      <van-form>
        <van-cell-group inset>
          <van-field v-model="registerForm.mode" label="报名方式" placeholder="solo 或 team" />
          <van-field v-model="registerForm.team_id" type="number" label="战队ID（战队报名时）" />
        </van-cell-group>
      </van-form>
    </van-dialog>

    <van-dialog v-model:show="showMatches" title="比赛场次" @confirm="showMatches = false">
      <van-cell-group inset>
        <van-cell v-for="m in matches" :key="m.id" :title="`第${m.round}轮 · ${m.group_no}组`" :label="`队${m.team_a_id} vs 队${m.team_b_id} · ${m.score_a}:${m.score_b}`">
          <template #value><StatusBadge kind="session" :status="m.status" /></template>
        </van-cell>
        <van-cell v-if="!matches.length" title="暂无比赛场次" />
      </van-cell-group>
    </van-dialog>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { showSuccessToast, showToast } from 'vant'
import StatusBadge from '@/components/StatusBadge.vue'
import { listTournaments, createTournament, registerTournament, drawGroups, listMatches, type Tournament, type Match } from '@/api/tournament'
import { GAME_TYPE_TEXT } from '@/constants'
import { useAuth } from '@/hooks/useAuth'

const { isStaffOrAdmin } = useAuth()
const list = ref<Tournament[]>([])
const matches = ref<Match[]>([])
const total = ref(0)
const page = ref(1)
const pageSize = 10
const selected = ref<Tournament | null>(null)
const showDetail = ref(false)
const showCreate = ref(false)
const showRegister = ref(false)
const showMatches = ref(false)
const createForm = ref({ name: '', game_type: 'lol', max_teams: 16 })
const registerForm = ref({ mode: 'solo', team_id: '' })

const detailActions = computed(() => {
  if (!selected.value) return []
  const acts: any[] = [{ name: '查看比赛场次', value: 'matches' }]
  if (selected.value.status === 'open') {
    acts.push({ name: '我要报名', value: 'register' })
  }
  if (isStaffOrAdmin.value && ['open', 'ready'].includes(selected.value.status)) {
    acts.push({ name: '自动抽签分组', value: 'draw' })
  }
  if (isStaffOrAdmin.value) {
    acts.push({ name: '创建新赛事', value: 'create' })
  }
  return acts
})

async function load() {
  const data = await listTournaments({ page: page.value, page_size: pageSize })
  list.value = data.list
  total.value = data.total
}

function detail(t: Tournament) {
  selected.value = t
  showDetail.value = true
}

async function onAction(action: any) {
  showDetail.value = false
  const t = selected.value
  if (!t) return
  if (action.value === 'matches') {
    matches.value = await listMatches(t.id)
    showMatches.value = true
  } else if (action.value === 'register') {
    registerForm.value = { mode: 'solo', team_id: '' }
    showRegister.value = true
  } else if (action.value === 'draw') {
    await drawGroups(t.id)
    showSuccessToast('抽签分组完成')
    load()
  } else if (action.value === 'create') {
    createForm.value = { name: '', game_type: 'lol', max_teams: 16 }
    showCreate.value = true
  }
}

async function submitCreate() {
  if (!createForm.value.name) {
    showToast('请输入赛事名称')
    return
  }
  await createTournament({ ...createForm.value, max_teams: Number(createForm.value.max_teams) || 16 })
  showSuccessToast('创建成功')
  load()
}

async function submitRegister() {
  if (!selected.value) return
  const payload: any = { mode: registerForm.value.mode }
  if (registerForm.value.mode === 'team') {
    payload.team_id = Number(registerForm.value.team_id)
  }
  await registerTournament(selected.value.id, payload)
  showSuccessToast('报名成功')
}

onMounted(load)
</script>

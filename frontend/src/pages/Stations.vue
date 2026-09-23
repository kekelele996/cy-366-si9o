<template>
  <div class="stations-page">
    <van-search v-model="area" placeholder="按区域搜索，如 A区" @search="load" />
    <van-dropdown-menu>
      <van-dropdown-item v-model="status" :options="statusOptions" @change="load" />
    </van-dropdown-menu>

    <van-cell-group inset title="机位列表">
      <van-cell v-for="st in list" :key="st.id" :title="`${st.name}`" :label="`${st.area} · ${STATION_TYPE_TEXT[st.station_type]} · ¥${st.price_per_hour}/小时`" @click="select(st)">
        <template #value><StatusBadge kind="station" :status="st.status" /></template>
      </van-cell>
    </van-cell-group>
    <van-pagination v-model="page" :total-items="total" :items-per-page="pageSize" @change="load" />

    <van-action-sheet v-model:show="showDetail" :title="selected?.name || '机位详情'" :actions="detailActions" @select="onDetailAction" cancel-text="关闭" />
    <van-dialog v-model:show="showForm" :title="formTitle" show-cancel-button @confirm="submitForm">
      <van-form>
        <van-cell-group inset>
          <van-field v-model="form.name" label="名称" placeholder="如 A区-01" />
          <van-field v-model="form.area" label="区域" placeholder="A区/B区/包厢区" />
          <van-field v-model="form.station_type" label="类型" placeholder="seat 或 box" />
          <van-field v-model="form.price_per_hour" type="number" label="时价(元)" />
        </van-cell-group>
      </van-form>
    </van-dialog>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref, computed } from 'vue'
import { showSuccessToast, showConfirmDialog } from 'vant'
import StatusBadge from '@/components/StatusBadge.vue'
import { listStations, createStation, updateStation, deleteStation, updateStationStatus, type Station } from '@/api/station'
import { STATION_STATUS_TEXT, STATION_TYPE_TEXT, AREA_OPTIONS } from '@/constants'
import { useAuth } from '@/hooks/useAuth'

const { isStaffOrAdmin } = useAuth()
const list = ref<Station[]>([])
const total = ref(0)
const page = ref(1)
const pageSize = 10
const area = ref('')
const status = ref('')
const statusOptions = [
  { text: '全部状态', value: '' },
  { text: '空闲', value: 'idle' },
  { text: '使用中', value: 'using' },
  { text: '故障', value: 'fault' },
  { text: '已预约', value: 'reserved' },
]

const selected = ref<Station | null>(null)
const showDetail = ref(false)
const showForm = ref(false)
const editingId = ref<number | null>(null)
const form = ref({ name: '', area: '', station_type: 'seat', price_per_hour: 8 })
const formTitle = computed(() => (editingId.value ? '编辑机位' : '新增机位'))

const detailActions = computed(() => {
  if (!selected.value) return []
  const acts: any[] = []
  if (isStaffOrAdmin.value) {
    acts.push({ name: '标记为空闲', value: 'idle' })
    acts.push({ name: '标记为故障', value: 'fault' })
    acts.push({ name: '编辑机位', value: 'edit' })
    acts.push({ name: '删除机位', value: 'delete', color: '#ee0a24' })
  }
  return acts
})

async function load() {
  try {
    const data = await listStations({ page: page.value, page_size: pageSize, area: area.value || undefined, status: status.value || undefined })
    list.value = data.list
    total.value = data.total
  } catch { /* toast 已处理 */ }
}

function select(st: Station) {
  selected.value = st
  showDetail.value = true
}

async function onDetailAction(action: any) {
  showDetail.value = false
  const st = selected.value
  if (!st) return
  if (action.value === 'edit') {
    editingId.value = st.id
    form.value = { name: st.name, area: st.area, station_type: st.station_type, price_per_hour: st.price_per_hour }
    showForm.value = true
  } else if (action.value === 'delete') {
    try {
      await showConfirmDialog({ title: '删除机位', message: `确定删除机位「${st.name}」吗？` })
      await deleteStation(st.id)
      showSuccessToast('删除成功')
      load()
    } catch { /* 取消 */ }
  } else {
    await updateStationStatus(st.id, action.value)
    showSuccessToast('状态已更新')
    load()
  }
}

function openCreate() {
  if (!isStaffOrAdmin.value) return
  editingId.value = null
  form.value = { name: '', area: AREA_OPTIONS[0], station_type: 'seat', price_per_hour: 8 }
  showForm.value = true
}

async function submitForm() {
  const payload = { ...form.value }
  if (editingId.value) {
    await updateStation(editingId.value, payload)
    showSuccessToast('更新成功')
  } else {
    await createStation(payload)
    showSuccessToast('创建成功')
  }
  load()
}

onMounted(load)
</script>

<style scoped>
.extra {
  margin: 8px 16px;
}
</style>

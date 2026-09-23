<template>
  <div class="data-table">
    <template v-if="items.length">
      <van-cell-group inset>
        <van-cell v-for="item in items" :key="item.id" :title="titleOf(item)" :label="labelOf(item)" is-link>
          <template #value>
            <StatusBadge v-if="statusKind" :kind="statusKind" :status="item[statusField]" />
          </template>
        </van-cell>
      </van-cell-group>
      <div v-if="total > pageSize" class="pager">
        <van-pagination v-model="current" :total-items="total" :items-per-page="pageSize" @change="onChange" />
      </div>
    </template>
    <EmptyState v-else description="暂无数据" />
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import StatusBadge from '@/components/StatusBadge.vue'
import EmptyState from '@/components/EmptyState.vue'

const props = withDefaults(defineProps<{
  items: any[]
  total: number
  pageSize?: number
  titleField?: string
  labelField?: string
  statusKind?: 'station' | 'reservation' | 'tournament' | 'session'
  statusField?: string
}>(), {
  pageSize: 10,
  titleField: 'name',
  labelField: 'id',
  statusField: 'status',
})

const emit = defineEmits<{ (e: 'change', page: number): void }>()
const current = ref(1)

function titleOf(item: any): string {
  return item[props.titleField] ?? `#${item.id}`
}

function labelOf(item: any): string {
  const v = item[props.labelField]
  return v !== undefined && v !== null ? String(v) : ''
}

function onChange() {
  emit('change', current.value)
}
</script>

<style scoped>
.data-table { margin-bottom: 12px; }
.pager { display: flex; justify-content: center; margin: 12px 0; }
</style>

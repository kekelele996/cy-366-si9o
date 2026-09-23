<template>
  <div class="audits-page">
    <van-cell-group inset>
      <van-cell v-for="a in list" :key="a.id" :title="a.action" :label="`${a.username || '系统'} · ${a.module} · ${a.ip}`" :value="formatTime(a.created_at)" />
    </van-cell-group>
    <EmptyState v-if="!list.length" description="暂无审计记录" />
    <van-pagination v-model="page" :total-items="total" :items-per-page="pageSize" @change="load" />
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import EmptyState from '@/components/EmptyState.vue'
import { listAudits, type AuditLog } from '@/api/audit'
import { formatTime } from '@/utils/format'

const list = ref<AuditLog[]>([])
const total = ref(0)
const page = ref(1)
const pageSize = 10

async function load() {
  const data = await listAudits({ page: page.value, page_size: pageSize })
  list.value = data.list
  total.value = data.total
}

onMounted(load)
</script>

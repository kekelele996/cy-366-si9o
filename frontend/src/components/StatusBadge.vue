<template>
  <van-tag :type="type" plain>{{ text }}</van-tag>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import {
  STATION_STATUS_TEXT, STATION_STATUS_TYPE,
  RESERVATION_STATUS_TEXT, RESERVATION_STATUS_TYPE,
  TOURNAMENT_STATUS_TEXT, TOURNAMENT_STATUS_TYPE,
  SESSION_STATUS_TEXT,
} from '@/constants'

const props = defineProps<{ kind: 'station' | 'reservation' | 'tournament' | 'session'; status: string }>()

const text = computed(() => {
  switch (props.kind) {
    case 'station': return STATION_STATUS_TEXT[props.status] || props.status
    case 'reservation': return RESERVATION_STATUS_TEXT[props.status] || props.status
    case 'tournament': return TOURNAMENT_STATUS_TEXT[props.status] || props.status
    case 'session': return SESSION_STATUS_TEXT[props.status] || props.status
  }
})

const type = computed(() => {
  switch (props.kind) {
    case 'station': return STATION_STATUS_TYPE[props.status] || 'default'
    case 'reservation': return RESERVATION_STATUS_TYPE[props.status] || 'default'
    case 'tournament': return TOURNAMENT_STATUS_TYPE[props.status] || 'default'
    default: return 'default'
  }
})
</script>

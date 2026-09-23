<template>
  <van-dialog v-model:show="visible" :title="title" :message="message" show-cancel-button @confirm="$emit('confirm')" @cancel="$emit('cancel')" />
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'

const props = defineProps<{ show: boolean; title?: string; message?: string }>()
const emit = defineEmits<{ (e: 'confirm'): void; (e: 'cancel'): void }>()
const visible = ref(false)
watch(() => props.show, (v) => { visible.value = v })
</script>

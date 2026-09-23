<template>
  <div class="empty-state">
    <van-empty :image-size="80" :description="description" />
  </div>
</template>

<script setup lang="ts">
defineProps<{ description?: string }>()
</script>

<style scoped>
.empty-state {
  padding: 24px 0;
}
</style>

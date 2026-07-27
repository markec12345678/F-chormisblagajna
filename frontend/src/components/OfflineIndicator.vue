<template>
  <div v-if="!isOnline" class="offline-bar">
    <span>⚡ {{ t('offline_mode') }}</span>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()
const isOnline = ref(navigator.onLine)

onMounted(() => {
  window.addEventListener('online', () => {
    isOnline.value = true
  })
  window.addEventListener('offline', () => {
    isOnline.value = false
  })
})
</script>

<style scoped>
.offline-bar {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  z-index: 10000;
  background: #dc2626;
  color: white;
  text-align: center;
  padding: 0.5rem;
  font-weight: 600;
  font-size: 0.875rem;
}
</style>

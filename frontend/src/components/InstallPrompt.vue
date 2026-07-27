<template>
  <div v-if="showInstall" class="install-prompt">
    <div class="install-content">
      <span class="install-icon">📱</span>
      <div>
        <strong>{{ t('install_app') }}</strong>
        <p>{{ t('install_description') }}</p>
      </div>
      <Button :label="t('install')" @click="installApp" />
      <Button :label="t('dismiss')" severity="secondary" text @click="dismiss" />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import Button from 'primevue/button'

const { t } = useI18n()
const showInstall = ref(false)
let deferredPrompt: any = null

onMounted(() => {
  if (localStorage.getItem('pwa-install-dismissed') === 'true') return
  window.addEventListener('beforeinstallprompt', (e) => {
    e.preventDefault()
    deferredPrompt = e
    showInstall.value = true
  })
})

async function installApp() {
  if (!deferredPrompt) return
  deferredPrompt.prompt()
  const { outcome } = await deferredPrompt.userChoice
  if (outcome === 'accepted') showInstall.value = false
  deferredPrompt = null
}

function dismiss() {
  showInstall.value = false
  localStorage.setItem('pwa-install-dismissed', 'true')
}
</script>

<style scoped>
.install-prompt {
  position: fixed;
  bottom: 0;
  left: 0;
  right: 0;
  z-index: 9999;
  background: #1e293b;
  border-top: 1px solid #334155;
  padding: 1rem;
  box-shadow: 0 -4px 20px rgba(0, 0, 0, 0.3);
}
.install-content {
  display: flex;
  align-items: center;
  gap: 1rem;
  max-width: 1200px;
  margin: 0 auto;
}
.install-icon {
  font-size: 2rem;
}
.install-content p {
  margin: 0;
  color: #94a3b8;
  font-size: 0.875rem;
}
</style>

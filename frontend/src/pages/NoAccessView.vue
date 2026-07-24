<template>
  <div
    class="flex justify-content-start align-items-center flex-column"
    style="height: 100vh; width: 100vw"
  >
    <i class="pi pi-ban mt-8" style="font-size: 4rem; color: #8c737c"></i>
    <h2>{{ $t('access_denied') }}</h2>
    <Button
      :label="$t('relogin')"
      severity="info"
      class="mt-2"
      @click="handleRelogin"
    />
  </div>
</template>
<script setup lang="ts">
import Button from 'primevue/button'
import auth from '@/services/auth'

function handleRelogin() {
  auth.signOut()
  window.location.href = '/login'
}
</script>

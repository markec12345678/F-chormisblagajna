<template>
  <div v-if="error" class="error-boundary">
    <div class="error-content">
      <i class="pi pi-exclamation-triangle" style="font-size: 3rem; color: #ef4444"></i>
      <h2>{{ $t('error_occurred') }}</h2>
      <p style="color: #64748b; max-width: 500px; text-align: center">
        {{ error.message || $t('error_occurred') }}
      </p>
      <Button :label="$t('retry')" @click="reset" class="mt-3" />
    </div>
  </div>
  <slot v-else />
</template>

<script setup lang="ts">
import { ref, onErrorCaptured } from 'vue'
import Button from 'primevue/button'

const error = ref<Error | null>(null)

onErrorCaptured((err) => {
  error.value = err
  return false
})

function reset() {
  error.value = null
}
</script>

<style scoped>
.error-boundary {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 50vh;
}
.error-content {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 1rem;
  text-align: center;
}
</style>

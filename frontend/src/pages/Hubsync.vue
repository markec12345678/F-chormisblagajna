<template>
  <div class="w-full">
    <div class="grid mx-2">
      <div class="col-12 flex">
        <div class="grid w-full">
          <div class="col-12">
            <h3>{{ $t('hubsync') }}</h3>
          </div>
          <div class="col-12">
            <div class="col-12 lg:col-4 flex flex-column gap-2">
              <label for="binary" class="font-bold">{{ $t('enabled') }}</label>
              <ToggleSwitch v-model="hubSync.settings.enabled" class="mb-2" />
            </div>
            <div class="col-12 lg:col-4 flex flex-column gap-2">
              <label for="binary" class="font-bold">{{ $t('server_host') }}</label>
              <InputText v-model="hubSync.settings.server_host" type="text" class="w-full mb-2" />
            </div>
            <div class="col-12 lg:col-4 flex flex-column gap-2">
              <label for="binary" class="font-bold">{{ $t('token') }}</label>
              <InputText v-model="hubSync.settings.token" type="text" class="w-full mb-2" />
            </div>
            <div class="col-12 lg:col-4 flex flex-column gap-2">
              <label for="binary" class="font-bold">{{ $t('sync_interval') }}</label>
              <InputText
                v-model.number="hubSync.settings.sync_interval"
                type="text"
                class="w-full mb-2"
              />
            </div>
            <div class="col-12 lg:col-4 flex flex-column gap-2">
              <label for="binary" class="font-bold">{{ $t('buffer_size') }}</label>
              <InputNumber
                v-model.number="hubSync.settings.buffer_size"
                mode="decimal"
                class="w-full mb-2"
              />
            </div>
            <div class="col-12 lg:col-4 flex flex-column gap-2">
              <label for="binary" class="font-bold">{{ $t('last_synced') }}</label>
              <span>{{ hubSync.last_synced || '-' }}</span>
            </div>
            <div class="col-12 lg:col-4 flex flex-column gap-2">
              <label for="binary" class="font-bold">{{ $t('sync_progress') }}</label>
              <span>{{ hubSync.sync_progress || '-' }}</span>
            </div>
            <div class="col-12 lg:col-4 flex flex-column gap-2">
              <label for="binary" class="font-bold">{{ $t('sync_sales') }}</label>
              <ToggleSwitch v-model="hubSync.settings.sync_sales" class="mb-2" />
            </div>
            <div class="col-12 lg:col-4 flex flex-column gap-2">
              <label for="binary" class="font-bold">{{ $t('sync_inventory') }}</label>
              <ToggleSwitch v-model="hubSync.settings.sync_inventory" class="mb-2" />
            </div>

            <div class="col-12">
              <Button :label="$t('save')" class="w-5rem" @click="save" />
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { InputText, ToggleSwitch, Button, InputNumber } from 'primevue'
import { useI18n } from 'vue-i18n'
import axios from 'axios'
import { useToast } from 'primevue/usetoast'
import auth from '../services/auth'

const { t } = useI18n()
const toast = useToast()

const hubSync = ref({
  settings: {
    enabled: false,
    server_host: '',
    token: '',
    buffer_size: 0,
    sync_interval: 0,
    sync_sales: true,
    sync_inventory: false,
  },
  last_synced: '',
  sync_progress: 0,
})

const save = () => {
  axios
    .patch(
      `http://${import.meta.env.VITE_APP_BACKEND_HOST}/hubsync/api/settings`,
      {
        data: hubSync.value,
      },
      {
        headers: {
          Authorization: `Bearer ${auth.accessToken.value}`,
        },
      },
    )
    .then((response) => {
      toast.add({
        severity: 'success',
        summary: 'Success',
        detail: t('data_saved_success'),
        life: 3000,
        group: 'br',
      })
    })
    .catch((err) => {
      toast.add({
        severity: 'error',
        summary: 'Error ' + err.response.status,
        detail: err.response.data,
        group: 'br',
      })
    })
}

const getHubsyncInfo = () => {
  axios
    .get(`http://${import.meta.env.VITE_APP_BACKEND_HOST}/hubsync/api/settings`, {
      headers: {
        Authorization: `Bearer ${auth.accessToken.value}`,
      },
    })
    .then((response) => {
      hubSync.value = response.data.data
    })
    .catch((err) => {
      toast.add({
        severity: 'error',
        summary: 'Error ' + err.response.status,
        detail: err.response.data,
        group: 'br',
      })
    })
}

getHubsyncInfo()
</script>

<style scoped>
label {
  font-size: 1rem;
  color: var(--p-button-text-primary-color);
}
</style>

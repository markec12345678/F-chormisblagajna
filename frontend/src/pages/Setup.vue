<template>
  <div class="setup-wrapper">
    <div class="setup-card">
      <div class="flex justify-content-end">
        <a href="https://nutrixpos.com/userguide/installation.html" target="_blank">
            <Button icon="pi pi-info-circle" class="p-button-text" />
        </a>
      </div>
      <!-- Logo / Icon -->
      <div class="setup-icon">
        <i class="pi pi-database"></i>
      </div>

      <h1 class="setup-title">Welcome to NutrixPOS</h1>
      <p class="setup-subtitle">
        No database connection has been configured yet.<br />
        Please provide your MongoDB connection details below.
      </p>

      <form class="setup-form" @submit.prevent="submit">
        <!-- Host -->
        <div class="field">
          <label for="host">Host</label>
          <InputText
            id="host"
            v-model="form.host"
            placeholder="e.g. localhost"
            :class="{ 'p-invalid': errors.host }"
            class="w-full"
          />
          <small v-if="errors.host" class="p-error">{{ errors.host }}</small>
        </div>

        <!-- Port -->
        <div class="field">
          <label for="port">Port</label>
          <InputNumber
            id="port"
            v-model="form.port"
            :useGrouping="false"
            placeholder="e.g. 27017"
            :class="{ 'p-invalid': errors.port }"
            class="w-full"
          />
          <small v-if="errors.port" class="p-error">{{ errors.port }}</small>
        </div>

        <!-- Database -->
        <div class="field">
          <label for="database">Database name</label>
          <InputText
            id="database"
            v-model="form.database"
            placeholder="e.g. nutrix"
            :class="{ 'p-invalid': errors.database }"
            class="w-full"
          />
          <small v-if="errors.database" class="p-error">{{ errors.database }}</small>
        </div>

        <!-- Database -->
        <div class="field">
          <label for="username">Username</label>
          <InputText
            id="username"
            v-model="form.username"
            placeholder="e.g. webadmin"
            :class="{ 'p-invalid': errors.username }"
            class="w-full"
          />
          <small v-if="errors.username" class="p-error">{{ errors.username }}</small>
        </div>

        <!-- Database -->
        <div class="field">
          <label for="password">Password</label>
          <InputText
            id="password"
            v-model="form.password"
            type="password"
            placeholder="Your database user password"
            :class="{ 'p-invalid': errors.password }"
            class="w-full"
          />
          <small v-if="errors.password" class="p-error">{{ errors.password }}</small>
        </div>

        <!-- Database -->
        <div class="field">
          <label for="password">Confirm Password</label>
          <InputText
            id="confirm_password"
            v-model="form.confirm_password"
            type="password"
            placeholder="Enter your password again"
            :class="{ 'p-invalid': errors.confirm_password }"
            class="w-full"
          />
          <small v-if="errors.confirm_password" class="p-error">{{ errors.confirm_password }}</small>
        </div>

        <!-- Success banner -->
        <div v-if="saved" class="success-banner mb-2 mt-4">
          <i class="pi pi-check-circle"></i>
          Configuration saved! The server is setup — please restart the application to apply the new settings.
        </div>

        <!-- Connection Test Success banner -->
        <div v-if="testSuccess && !saved" class="success-banner mb-2 mt-4">
          <i class="pi pi-check-circle"></i>
          Connection successful! You can now save the configuration.
        </div>

        <!-- Error banner -->
        <div v-if="serverError" class="error-banner mt-3 mb-2">
          <i class="pi pi-times-circle"></i>
          {{ serverError }}
        </div>

        <div v-if="!saved" class="flex flex-column gap-1 mt-1">
          <Button
            type="button"
            label="Test Connection"
            icon="pi pi-bolt"
            class="submit-btn p-button-outlined"
            style="background: transparent; border: 1px solid rgba(255, 255, 255, 0.4) !important; color: #1a2a3a;"
            :loading="testing"
            :disabled="loading || testing || saved"
            @click="testConnection"
          />

          <Button
            type="submit"
            label="Save & Continue"
            icon="pi pi-arrow-right"
            iconPos="right"
            class="submit-btn"
            :loading="loading"
            :disabled="loading || testing || saved || !testSuccess"
          />
        </div>

        <Button v-if="saved" class="submit-btn mt-2" label="Let's go 🚀" @click="router.push({ path: '/home' })" />
      </form>
    </div>
  </div>
</template>

<script setup lang="ts">
import '../styles/setup-shared.css'
import { ref, reactive, watch } from 'vue'
import InputText from 'primevue/inputtext'
import InputNumber from 'primevue/inputnumber'
import Button from 'primevue/button'
import axios from 'axios'
import { useRouter } from 'vue-router'
const router = useRouter()

const apiUrl = `http://${import.meta.env.VITE_APP_BACKEND_HOST}`

const form = reactive({
  host: '',
  port: null as number | null,
  database: '',
  username: '',
  password: '',
  confirm_password: '',
})



const errors = reactive({
  host: '',
  port: '',
  database: '',
  username: '',
  password: '',
  confirm_password: '',
})

const loading = ref(false)
const testing = ref(false)
const saved   = ref(false)
const testSuccess = ref(false)
const serverError = ref('')

watch(form, () => {
  testSuccess.value = false
})

function validate(): boolean {
  errors.host = form.host.trim() ? '' : 'Host is required'
  errors.port = form.port && form.port > 0 ? '' : 'A valid port number is required'
  errors.database = form.database.trim() ? '' : 'Database name is required'
  errors.confirm_password = form.password == form.confirm_password ? '' : 'Passwords do not match'

  return !errors.host && !errors.port && !errors.database && !errors.username && !errors.password && !errors.confirm_password
}

async function testConnection() {
  if (!validate()) return

  testing.value = true
  serverError.value = ''
  testSuccess.value = false

  try {
    await axios.post(`${apiUrl}/api/setup/test-connection`, {
      host:     form.host.trim(),
      port:     form.port,
      database: form.database.trim(),
      username: form.username.trim(),
      password: form.password
    })
    testSuccess.value = true
  } catch (err: any) {
    serverError.value =
      err?.response?.data || 'Failed to connect to database. Check the credentials.'
    testSuccess.value = false
  } finally {
    testing.value = false
  }
}

async function submit() {
  if (!validate()) return

  loading.value = true
  serverError.value = ''

  try {
    await axios.post(`${apiUrl}/api/setup/config`, {
      host:     form.host.trim(),
      port:     form.port,
      database: form.database.trim(),
      username: form.username.trim(),
      password: form.password
    })
    saved.value = true
    testSuccess.value = false
  } catch (err: any) {
    serverError.value =
      err?.response?.data || 'Failed to save configuration. Check the backend logs.'
  } finally {
    loading.value = false
  }
}
</script>



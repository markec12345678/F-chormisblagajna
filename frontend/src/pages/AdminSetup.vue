<template>
  <div class="setup-wrapper">
    <div class="setup-card">
      <div class="flex justify-content-end">
        <a href="https://nutrixpos.com/userguide/installation.html" target="_blank">
          <Button icon="pi pi-info-circle" class="p-button-text" />
        </a>
      </div>
      <div class="setup-icon">
        <i class="pi pi-user-plus"></i>
      </div>

      <h1 class="setup-title">{{ $t('create_admin_user') }}</h1>
      <p class="setup-subtitle">
        {{ $t('no_admin_exists') }}<br />
        {{ $t('create_admin_subtitle') }}
      </p>

      <form class="setup-form" @submit.prevent="submit">
        <div class="field">
          <label for="username">{{ $t('username') }}</label>
          <InputText
            id="username"
            v-model="form.username"
            :placeholder="$t('admin_placeholder')"
            :class="{ 'p-invalid': errors.username }"
            class="w-full"
          />
          <small v-if="errors.username" class="p-error">{{ errors.username }}</small>
        </div>

        <div class="field">
          <label for="email">{{ $t('email') }}</label>
          <InputText
            id="email"
            v-model="form.email"
            type="email"
            :placeholder="$t('email_placeholder')"
            :class="{ 'p-invalid': errors.email }"
            class="w-full"
          />
          <small v-if="errors.email" class="p-error">{{ errors.email }}</small>
        </div>

        <div class="field">
          <label for="password">{{ $t('password') }}</label>
          <InputText
            id="password"
            v-model="form.password"
            type="password"
            :placeholder="$t('password_placeholder')"
            :class="{ 'p-invalid': errors.password }"
            class="w-full"
          />
          <small v-if="errors.password" class="p-error">{{ errors.password }}</small>
        </div>

        <div class="field">
          <label for="confirm_password">{{ $t('confirm_password') }}</label>
          <InputText
            id="confirm_password"
            v-model="form.confirm_password"
            type="password"
            :placeholder="$t('confirm_password_placeholder')"
            :class="{ 'p-invalid': errors.confirm_password }"
            class="w-full"
          />
          <small v-if="errors.confirm_password" class="p-error">{{
            errors.confirm_password
          }}</small>
        </div>

        <div v-if="serverError" class="error-banner mt-3 mb-2">
          <i class="pi pi-times-circle"></i>
          {{ serverError }}
        </div>

        <div v-if="success" class="success-banner mb-2 mt-4">
          <i class="pi pi-check-circle"></i>
          {{ $t('admin_created_success') }}
        </div>

        <div v-if="!success" class="flex flex-column gap-1 mt-1">
          <Button
            type="submit"
            :label="$t('create_admin')"
            icon="pi pi-user-plus"
            iconPos="right"
            class="submit-btn"
            :loading="loading"
            :disabled="loading"
          />
        </div>
      </form>
    </div>
  </div>
</template>

<script setup lang="ts">
import '../styles/setup-shared.css'
import { ref, reactive } from 'vue'
import { useI18n } from 'vue-i18n'
import InputText from 'primevue/inputtext'
import Button from 'primevue/button'
import axios from 'axios'
import { useRouter } from 'vue-router'
import auth from '@/services/auth'

const { t } = useI18n()
const router = useRouter()

const apiUrl = `http://${import.meta.env.VITE_APP_BACKEND_HOST}${import.meta.env.VITE_APP_MODULE_CORE_API_PREFIX}`

const form = reactive({
  username: '',
  email: '',
  password: '',
  confirm_password: '',
})

const errors = reactive({
  username: '',
  email: '',
  password: '',
  confirm_password: '',
})

const loading = ref(false)
const serverError = ref('')
const success = ref(false)

function validate(): boolean {
  errors.username = form.username.trim() ? '' : t('username_required')
  errors.email = form.email.trim() ? '' : t('email_required')
  errors.email = /^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(form.email) ? '' : t('email_invalid')
  errors.password = form.password.length >= 6 ? '' : t('password_min_length')
  errors.confirm_password = form.password === form.confirm_password ? '' : t('passwords_no_match')

  return !errors.username && !errors.email && !errors.password && !errors.confirm_password
}

async function submit() {
  if (!validate()) return

  loading.value = true
  serverError.value = ''

  try {
    const response = await axios.post(`${apiUrl}/api/auth/register`, {
      username: form.username.trim(),
      email: form.email.trim(),
      password: form.password,
    })

    const data = response.data

    auth.accessToken.value = data.token
    auth.currentUser.value = data.user
    auth.isAuthenticated.value = true

    localStorage.setItem('nutrix_token', data.token)
    localStorage.setItem('nutrix_user', JSON.stringify(data.user))

    success.value = true

    setTimeout(() => {
      router.push({ path: '/' })
    }, 1500)
  } catch (err) {
    const data = axios.isAxiosError(err) ? err.response?.data : undefined
    serverError.value = data?.error || data?.message || t('admin_create_failed')
  } finally {
    loading.value = false
  }
}
</script>

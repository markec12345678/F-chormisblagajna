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

      <h1 class="setup-title">Create Admin User</h1>
      <p class="setup-subtitle">
        No admin user exists yet.<br />
        Create your administrator account to get started.
      </p>

      <form class="setup-form" @submit.prevent="submit">
        <div class="field">
          <label for="username">Username</label>
          <InputText
            id="username"
            v-model="form.username"
            placeholder="e.g. admin"
            :class="{ 'p-invalid': errors.username }"
            class="w-full"
          />
          <small v-if="errors.username" class="p-error">{{ errors.username }}</small>
        </div>

        <div class="field">
          <label for="email">Email</label>
          <InputText
            id="email"
            v-model="form.email"
            type="email"
            placeholder="e.g. admin@example.com"
            :class="{ 'p-invalid': errors.email }"
            class="w-full"
          />
          <small v-if="errors.email" class="p-error">{{ errors.email }}</small>
        </div>

        <div class="field">
          <label for="password">Password</label>
          <InputText
            id="password"
            v-model="form.password"
            type="password"
            placeholder="Your password"
            :class="{ 'p-invalid': errors.password }"
            class="w-full"
          />
          <small v-if="errors.password" class="p-error">{{ errors.password }}</small>
        </div>

        <div class="field">
          <label for="confirm_password">Confirm Password</label>
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

        <div v-if="serverError" class="error-banner mt-3 mb-2">
          <i class="pi pi-times-circle"></i>
          {{ serverError }}
        </div>

        <div v-if="success" class="success-banner mb-2 mt-4">
          <i class="pi pi-check-circle"></i>
          Admin created successfully! Redirecting to home...
        </div>

        <div v-if="!success" class="flex flex-column gap-1 mt-1">
          <Button
            type="submit"
            label="Create Admin"
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
import InputText from 'primevue/inputtext'
import Button from 'primevue/button'
import axios from 'axios'
import { useRouter } from 'vue-router'
import auth from '@/services/auth'

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
  errors.username = form.username.trim() ? '' : 'Username is required'
  errors.email = form.email.trim() ? '' : 'Email is required'
  errors.email = /^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(form.email) ? '' : 'Invalid email format'
  errors.password = form.password.length >= 6 ? '' : 'Password must be at least 6 characters'
  errors.confirm_password = form.password === form.confirm_password ? '' : 'Passwords do not match'

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
      password: form.password
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
  } catch (err: any) {
    serverError.value = err?.response?.data?.error || err?.response?.data?.message || 'Failed to create admin user'
  } finally {
    loading.value = false
  }
}
</script>


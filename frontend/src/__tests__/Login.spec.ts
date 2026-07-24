import { describe, it, expect, beforeEach, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import { createPinia, setActivePinia } from 'pinia'
import { createI18n } from 'vue-i18n'
import Login from '@/pages/Login.vue'
import ToastService from 'primevue/toastservice'
import ConfirmationService from 'primevue/confirmationservice'

vi.mock('vue-router', () => ({
  useRouter: () => ({
    push: vi.fn(),
    replace: vi.fn(),
    back: vi.fn(),
  }),
  useRoute: () => ({
    params: {},
    query: {},
  }),
}))

vi.mock('@/services/auth', () => ({
  default: {
    accessToken: { value: null },
    currentUser: { value: null },
    signOut: vi.fn(),
    login: vi.fn(),
  },
}))

const mockGet = vi.hoisted(() => vi.fn())
const mockPost = vi.hoisted(() => vi.fn())

vi.mock('axios', () => ({
  default: {
    get: mockGet,
    post: mockPost,
    put: vi.fn(),
    delete: vi.fn(),
    patch: vi.fn(),
    defaults: { headers: { common: {} } },
  },
}))

const i18n = createI18n({
  legacy: false,
  locale: 'en',
  fallbackLocale: 'en',
  messages: {
    en: {
      welcome_back: 'Welcome Back',
      sign_in_subtitle: 'Sign in to continue',
      username: 'Username',
      password: 'Password',
      enter_username: 'Enter username',
      enter_password: 'Enter password',
      sign_in: 'Sign In',
      or: 'or',
      no_account: "Don't have an account?",
      contact_admin: 'Contact admin',
      username_required: 'Username is required',
      password_required: 'Password is required',
      invalid_credentials: 'Invalid credentials',
      login_failed: 'Login failed',
    },
  },
})

const stubs = {
  InputText: {
    template:
      '<input class="input-text-stub" :id="id" :value="modelValue" @input="$emit(\'update:modelValue\', $event.target.value)" />',
    props: ['modelValue', 'placeholder', 'id', 'class'],
  },
  Password: {
    template: '<input class="password-stub" type="password" :id="id" :value="modelValue" />',
    props: ['modelValue', 'placeholder', 'id', 'feedback', 'toggleMask', 'inputClass'],
  },
  Button: {
    template: '<button class="btn-stub" type="submit">{{ label }}</button>',
    props: ['label', 'icon', 'loading', 'disabled', 'class'],
  },
  Divider: { template: '<hr class="divider-stub" />', props: ['align'] },
  RouterLink: { template: '<a class="router-link-stub"><slot/></a>', props: ['to'] },
  Toast: { template: '<div />' },
}

describe('Login', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    mockGet.mockReset()
    mockPost.mockReset()
  })

  it('renders the login form', () => {
    const wrapper = mount(Login, {
      global: { plugins: [i18n, ToastService, ConfirmationService], stubs },
    })

    expect(wrapper.find('.login-form').exists()).toBe(true)
  })

  it('shows welcome title', () => {
    const wrapper = mount(Login, {
      global: { plugins: [i18n, ToastService, ConfirmationService], stubs },
    })

    expect(wrapper.text()).toContain('Welcome Back')
  })

  it('shows username input field', () => {
    const wrapper = mount(Login, {
      global: { plugins: [i18n, ToastService, ConfirmationService], stubs },
    })

    expect(wrapper.find('.input-text-stub').exists()).toBe(true)
  })

  it('shows password input field', () => {
    const wrapper = mount(Login, {
      global: { plugins: [i18n, ToastService, ConfirmationService], stubs },
    })

    expect(wrapper.find('.password-stub').exists()).toBe(true)
  })

  it('shows sign in button', () => {
    const wrapper = mount(Login, {
      global: { plugins: [i18n, ToastService, ConfirmationService], stubs },
    })

    expect(wrapper.findAll('.btn-stub').length).toBeGreaterThan(0)
  })

  it('shows sign in subtitle', () => {
    const wrapper = mount(Login, {
      global: { plugins: [i18n, ToastService, ConfirmationService], stubs },
    })

    expect(wrapper.text()).toContain('Sign in to continue')
  })
})

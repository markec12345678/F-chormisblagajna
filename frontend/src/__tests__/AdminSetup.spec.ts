import { describe, it, expect, beforeEach, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import { createPinia, setActivePinia } from 'pinia'
import { createI18n } from 'vue-i18n'
import AdminSetup from '@/pages/AdminSetup.vue'

vi.mock('vue-router', () => ({
  useRouter: () => ({ push: vi.fn(), replace: vi.fn(), back: vi.fn() }),
  useRoute: () => ({ params: {}, query: {} }),
}))

vi.mock('@/services/auth', () => ({
  default: {
    accessToken: { value: 'test-token' },
    currentUser: { value: null },
    isAuthenticated: { value: false },
    signOut: vi.fn(),
  },
}))

const mockPost = vi.hoisted(() => vi.fn())

vi.mock('axios', () => ({
  default: {
    get: vi.fn(),
    post: mockPost,
    put: vi.fn(),
    delete: vi.fn(),
    patch: vi.fn(),
    isAxiosError: vi.fn(),
    defaults: { headers: { common: {} } },
  },
}))

const i18n = createI18n({
  legacy: false,
  locale: 'en',
  fallbackLocale: 'en',
  messages: {
    en: {
      create_admin_user: 'Create Admin User',
      no_admin_exists: 'No admin user exists.',
      create_admin_subtitle: 'Create the first admin account.',
      username: 'Username',
      email: 'Email',
      password: 'Password',
      confirm_password: 'Confirm Password',
      admin_placeholder: 'Enter username',
      email_placeholder: 'Enter email',
      password_placeholder: 'Enter password',
      confirm_password_placeholder: 'Confirm password',
      create_admin: 'Create Admin',
      admin_created_success: 'Admin created successfully',
    },
  },
})

const stubs = {
  InputText: { template: '<input class="input-text-stub" />', props: ['modelValue', 'placeholder', 'type'] },
  Button: { template: '<button class="btn-stub">{{ label }}</button>', props: ['label', 'icon', 'iconPos', 'loading', 'disabled'] },
}

describe('AdminSetup', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    mockPost.mockReset()
  })

  it('renders setup heading', () => {
    const wrapper = mount(AdminSetup, {
      global: { plugins: [i18n], stubs },
    })
    expect(wrapper.text()).toContain('Create Admin User')
  })

  it('renders form fields for username, email and password', () => {
    const wrapper = mount(AdminSetup, {
      global: { plugins: [i18n], stubs },
    })
    expect(wrapper.text()).toContain('Username')
    expect(wrapper.text()).toContain('Email')
    expect(wrapper.text()).toContain('Password')
    expect(wrapper.text()).toContain('Confirm Password')
  })

  it('renders create admin submit button', () => {
    const wrapper = mount(AdminSetup, {
      global: { plugins: [i18n], stubs },
    })
    expect(wrapper.text()).toContain('Create Admin')
  })

  it('renders setup icon', () => {
    const wrapper = mount(AdminSetup, {
      global: { plugins: [i18n], stubs },
    })
    expect(wrapper.find('.setup-icon').exists()).toBe(true)
  })
})

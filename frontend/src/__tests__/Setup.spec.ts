import { describe, it, expect, beforeEach, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import { createPinia, setActivePinia } from 'pinia'
import { createI18n } from 'vue-i18n'
import Setup from '@/pages/Setup.vue'

vi.mock('vue-router', () => ({
  useRouter: () => ({ push: vi.fn(), replace: vi.fn(), back: vi.fn() }),
  useRoute: () => ({ params: {}, query: {} }),
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
      welcome_to_nutrix: 'Welcome to Nutrix',
      no_db_configured: 'No database configured.',
      provide_db_details: 'Provide your database details.',
      host: 'Host',
      port: 'Port',
      database_name: 'Database Name',
      username: 'Username',
      password: 'Password',
      confirm_password: 'Confirm Password',
      host_placeholder: 'Enter host',
      port_placeholder: 'Enter port',
      database_placeholder: 'Enter database name',
      username_placeholder: 'Enter username',
      password_placeholder: 'Enter password',
      confirm_password_placeholder: 'Confirm password',
      test_connection: 'Test Connection',
      save_and_continue: 'Save and Continue',
      lets_go: "Let's Go",
      connection_successful: 'Connection successful',
      config_saved_restart: 'Config saved, restart required',
    },
  },
})

const stubs = {
  InputText: { template: '<input class="input-text-stub" />', props: ['modelValue', 'placeholder', 'type'] },
  InputNumber: { template: '<input class="input-number-stub" />', props: ['modelValue', 'useGrouping', 'placeholder'] },
  Button: { template: '<button class="btn-stub">{{ label }}</button>', props: ['label', 'icon', 'iconPos', 'loading', 'disabled'] },
}

describe('Setup', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    mockPost.mockReset()
  })

  it('renders welcome heading', () => {
    const wrapper = mount(Setup, {
      global: { plugins: [i18n], stubs },
    })
    expect(wrapper.text()).toContain('Welcome to Nutrix')
  })

  it('renders database connection form fields', () => {
    const wrapper = mount(Setup, {
      global: { plugins: [i18n], stubs },
    })
    expect(wrapper.text()).toContain('Host')
    expect(wrapper.text()).toContain('Port')
    expect(wrapper.text()).toContain('Database Name')
    expect(wrapper.text()).toContain('Username')
    expect(wrapper.text()).toContain('Password')
  })

  it('renders test connection and save buttons', () => {
    const wrapper = mount(Setup, {
      global: { plugins: [i18n], stubs },
    })
    expect(wrapper.text()).toContain('Test Connection')
    expect(wrapper.text()).toContain('Save and Continue')
  })

  it('renders setup icon', () => {
    const wrapper = mount(Setup, {
      global: { plugins: [i18n], stubs },
    })
    expect(wrapper.find('.setup-icon').exists()).toBe(true)
  })
})

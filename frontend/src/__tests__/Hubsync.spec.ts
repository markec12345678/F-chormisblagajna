import { describe, it, expect, beforeEach, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import { createPinia, setActivePinia } from 'pinia'
import { createI18n } from 'vue-i18n'
import Hubsync from '@/pages/Hubsync.vue'
import ToastService from 'primevue/toastservice'

vi.mock('@/services/auth', () => ({
  default: {
    accessToken: { value: 'test-token' },
    currentUser: { value: null },
    signOut: vi.fn(),
  },
}))

const mockGet = vi.hoisted(() => vi.fn())

vi.mock('axios', () => ({
  default: {
    get: mockGet,
    post: vi.fn(),
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
      hubsync: 'Hub Sync',
      enabled: 'Enabled',
      server_host: 'Server Host',
      token: 'Token',
      sync_interval: 'Sync Interval',
      buffer_size: 'Buffer Size',
      last_synced: 'Last Synced',
      sync_progress: 'Sync Progress',
      sync_sales: 'Sync Sales',
      sync_inventory: 'Sync Inventory',
      save: 'Save',
    },
  },
})

const stubs = {
  InputText: { template: '<input class="input-text-stub" />', props: ['modelValue', 'type'] },
  InputNumber: { template: '<input class="input-number-stub" />', props: ['modelValue', 'mode'] },
  Button: { template: '<button class="btn-stub">{{ label }}</button>', props: ['label'] },
  ToggleSwitch: {
    template: '<input type="checkbox" class="toggle-stub" />',
    props: ['modelValue'],
  },
}

describe('Hubsync', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    mockGet.mockReset()
    mockGet.mockResolvedValue({
      data: {
        data: {
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
        },
      },
    })
  })

  it('renders hubsync page heading', async () => {
    const wrapper = mount(Hubsync, {
      global: { plugins: [i18n, ToastService], stubs },
    })
    await vi.waitFor(() => {
      expect(wrapper.text()).toContain('Hub Sync')
    })
  })

  it('renders form fields for hubsync settings', async () => {
    const wrapper = mount(Hubsync, {
      global: { plugins: [i18n, ToastService], stubs },
    })
    await vi.waitFor(() => {
      expect(wrapper.text()).toContain('Server Host')
      expect(wrapper.text()).toContain('Token')
      expect(wrapper.text()).toContain('Sync Interval')
      expect(wrapper.text()).toContain('Buffer Size')
    })
  })

  it('renders save button', async () => {
    const wrapper = mount(Hubsync, {
      global: { plugins: [i18n, ToastService], stubs },
    })
    await vi.waitFor(() => {
      expect(wrapper.text()).toContain('Save')
    })
  })

  it('renders toggle switches for sync options', async () => {
    const wrapper = mount(Hubsync, {
      global: { plugins: [i18n, ToastService], stubs },
    })
    await vi.waitFor(() => {
      expect(wrapper.findAll('.toggle-stub').length).toBeGreaterThanOrEqual(2)
    })
  })
})

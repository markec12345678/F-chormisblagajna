import { describe, it, expect, beforeEach, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import { createPinia, setActivePinia } from 'pinia'
import { createI18n } from 'vue-i18n'
import Kitchen from '@/pages/Kitchen.vue'
import ToastService from 'primevue/toastservice'

vi.mock('@/stores', () => ({
  globalStore: () => ({
    shop_mode: '',
    orientation: 'ltr',
    getColorMode: 'light',
    setOrientation: vi.fn(),
  }),
}))

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
    patch: vi.fn(),
    delete: vi.fn(),
  },
}))

const i18n = createI18n({
  legacy: false,
  locale: 'en',
  fallbackLocale: 'en',
  globalInjection: true,
  messages: {
    en: {
      order: 'Order',
      new_orders_appear: 'New orders will appear here',
      order_finished: 'Order Finished',
      order_ready_to_serve: 'is ready to serve',
    },
  },
})

vi.stubGlobal(
  'WebSocket',
  class {
    send = vi.fn()
    close = vi.fn()
    readyState = 1
    static OPEN = 1
    onopen: ((...args: any[]) => void) | null = null
    onmessage: ((...args: any[]) => void) | null = null
    onerror: ((...args: any[]) => void) | null = null
    onclose: ((...args: any[]) => void) | null = null
    constructor() {
      setTimeout(() => this.onopen?.(), 0)
    }
  },
)

const stubs = {
  QueueOrder: { template: '<div class="queue-order-stub" />', props: ['order', 'number'] },
  ProgressSpinner: { template: '<div class="spinner-stub" />' },
  Toast: { template: '<div />' },
}

function mockSettingsAndLanguage() {
  mockGet.mockImplementation((url: string) => {
    if (url.includes('/settings')) {
      return Promise.resolve({ data: { data: { language: { code: 'en' } } } })
    }
    if (url.includes('/languages/')) {
      return Promise.resolve({ data: { data: { code: 'en', pack: {}, orientation: 'ltr' } } })
    }
    return Promise.resolve({ data: { data: null } })
  })
}

describe('Kitchen', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    mockGet.mockReset()
    Object.defineProperty(window, 'innerWidth', { value: 1024, writable: true, configurable: true })
  })

  it('shows loading spinner initially', () => {
    mockGet.mockReturnValue(new Promise(() => {}))

    const wrapper = mount(Kitchen, {
      global: { plugins: [i18n, ToastService], stubs },
    })

    expect(wrapper.find('.spinner-stub').exists()).toBe(true)
  })

  it('shows empty state when no orders', async () => {
    mockSettingsAndLanguage()

    const wrapper = mount(Kitchen, {
      global: { plugins: [i18n, ToastService], stubs },
    })

    await vi.waitFor(() => {
      expect(wrapper.find('.spinner-stub').exists()).toBe(false)
      expect(wrapper.findAll('p').length).toBeGreaterThan(0)
    })
  })
})

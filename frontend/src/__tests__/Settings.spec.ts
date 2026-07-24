import { describe, it, expect, beforeEach, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import { createPinia, setActivePinia } from 'pinia'
import { createI18n } from 'vue-i18n'
import Settings from '@/pages/Settings.vue'
import ToastService from 'primevue/toastservice'

vi.mock('@/stores', () => ({
  globalStore: () => ({
    shop_mode: '',
    orientation: 'ltr',
    getColorMode: 'light',
    getSettings: { orders: { default_cost_calculation_method: 'exact' } },
    setShopMode: vi.fn(),
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
const mockPatch = vi.hoisted(() => vi.fn())

vi.mock('axios', () => ({
  default: {
    get: mockGet,
    patch: mockPatch,
    post: vi.fn(),
    delete: vi.fn(),
  },
}))

const i18n = createI18n({
  legacy: false,
  locale: 'en',
  fallbackLocale: 'en',
  messages: {
    en: {
      settings: 'Settings',
      queues: 'Queues',
      prefix: 'Prefix',
      next: 'Next',
      add: 'Add',
      default_cost_calculation_method: 'Default Cost Calculation Method',
      exact: 'Exact',
      average: 'Average',
      save_settings: 'Save Settings',
      inventory_item: 'Inventory Item',
      stock_alert_threshold: 'Stock Alert Threshold',
      order: 'Order',
      language: 'Language',
      select_language: 'Select Language',
      shop_mode: 'Shop Mode',
      retail: 'Retail',
      kitchen: 'Kitchen',
    },
  },
})

const stubs = {
  Skeleton: { template: '<div class="skeleton-stub" />' },
  InlineMessage: { template: '<div class="inline-msg-stub"><slot/></div>', props: ['severity'] },
  InputNumber: { template: '<input class="input-number-stub" />', props: ['modelValue', 'min'] },
  InputText: { template: '<input class="input-text-stub" />', props: ['modelValue'] },
  Divider: { template: '<hr class="divider-stub" />' },
  Button: { template: '<button class="btn-stub">{{ label }}</button>', props: ['label', 'icon', 'severity'] },
  RadioButton: { template: '<input type="radio" class="radio-stub" />', props: ['modelValue', 'value', 'inputId', 'name'] },
  Badge: { template: '<span class="badge-stub" />', props: ['value', 'size', 'severity'] },
  Select: { template: '<select class="select-stub" />', props: ['modelValue', 'options', 'optionLabel', 'placeholder'] },
  ToggleSwitch: { template: '<input type="checkbox" class="toggle-stub" />', props: ['modelValue'] },
  Toast: { template: '<div />' },
}

describe('Settings', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    mockGet.mockReset()
    mockPatch.mockReset()
  })

  it('shows skeleton while loading', () => {
    mockGet.mockReturnValue(new Promise(() => {}))

    const wrapper = mount(Settings, {
      global: { plugins: [i18n, ToastService], stubs },
    })

    expect(wrapper.findAll('.skeleton-stub').length).toBeGreaterThan(0)
  })

  it('shows settings heading after loading', async () => {
    mockGet.mockResolvedValue({
      data: {
        data: {
          inventory: { stock_alert_treshold: 0 },
          orders: { queues: [], default_cost_calculation_method: 'average' },
          language: { language: 'English', code: 'en' },
          payment_sources: [{ name: 'Cash' }],
          shop_mode: '',
          auto_open_cash_drawer: false,
        },
      },
    })

    const wrapper = mount(Settings, {
      global: { plugins: [i18n, ToastService], stubs },
    })

    await vi.waitFor(() => {
      expect(wrapper.text()).toContain('Settings')
    })
  })

  it('shows error message when settings fail to load', async () => {
    mockGet.mockRejectedValue(new Error('Network error'))

    const wrapper = mount(Settings, {
      global: { plugins: [i18n, ToastService], stubs },
    })

    await vi.waitFor(() => {
      expect(wrapper.find('.inline-msg-stub').exists()).toBe(true)
    })
  })

  it('renders queue section after loading', async () => {
    mockGet.mockResolvedValue({
      data: {
        data: {
          inventory: { stock_alert_treshold: 5 },
          orders: {
            queues: [{ prefix: 'A', next: 1 }],
            default_cost_calculation_method: 'exact',
          },
          language: { language: 'English', code: 'en' },
          payment_sources: [{ name: 'Cash' }],
          shop_mode: '',
          auto_open_cash_drawer: false,
        },
      },
    })

    const wrapper = mount(Settings, {
      global: { plugins: [i18n, ToastService], stubs },
    })

    await vi.waitFor(() => {
      expect(wrapper.text()).toContain('Queues')
    })
  })
})

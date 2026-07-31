import { describe, it, expect, beforeEach, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import { createPinia, setActivePinia } from 'pinia'
import { createI18n } from 'vue-i18n'
import OnlineOrders from '@/pages/OnlineOrders.vue'
import ToastService from 'primevue/toastservice'

vi.mock('@/services/auth', () => ({
  default: {
    accessToken: { value: 'test-token' },
    currentUser: { value: null },
    signOut: vi.fn(),
  },
}))

const mockGet = vi.hoisted(() => vi.fn())
const mockPut = vi.hoisted(() => vi.fn())

vi.mock('axios', () => ({
  default: {
    get: mockGet,
    post: vi.fn(),
    put: mockPut,
    delete: vi.fn(),
    defaults: { headers: { common: {} } },
  },
}))

const i18n = createI18n({
  legacy: false,
  locale: 'en',
  fallbackLocale: 'en',
  messages: {
    en: {
      online_orders: 'Online Orders',
      all_statuses: 'All Statuses',
      refresh: 'Refresh',
      no_online_orders: 'No online orders',
      customer: 'Customer',
      phone: 'Phone',
      order_type: 'Order Type',
      address: 'Address',
      total: 'Total',
    },
  },
})

const stubs = {
  Button: { template: '<button class="btn-stub">{{ label }}</button>', props: ['label', 'icon', 'loading'] },
  Card: { template: '<div class="card-stub"><slot name="title" /><slot name="content" /></div>' },
  Tag: { template: '<span class="tag-stub"><slot /></span>', props: ['severity', 'value'] },
  Dropdown: { template: '<select class="dropdown-stub" />', props: ['modelValue', 'options', 'placeholder'] },
  Divider: { template: '<hr class="divider-stub" />' },
}

describe('OnlineOrders', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    mockGet.mockReset()
    mockPut.mockReset()
    mockGet.mockResolvedValue({ data: { data: [] } })
  })

  it('renders page heading', async () => {
    const wrapper = mount(OnlineOrders, { global: { plugins: [i18n, ToastService], stubs } })
    await vi.waitFor(() => expect(wrapper.text()).toContain('Online Orders'))
  })

  it('shows empty state when no orders', async () => {
    const wrapper = mount(OnlineOrders, { global: { plugins: [i18n, ToastService], stubs } })
    await vi.waitFor(() => expect(wrapper.text()).toContain('No online orders'))
  })

  it('shows refresh button', async () => {
    const wrapper = mount(OnlineOrders, { global: { plugins: [i18n, ToastService], stubs } })
    await vi.waitFor(() => expect(wrapper.text()).toContain('Refresh'))
  })

  it('loads orders from API', async () => {
    mount(OnlineOrders, { global: { plugins: [i18n, ToastService], stubs } })
    await vi.waitFor(() => expect(mockGet).toHaveBeenCalled())
  })
})

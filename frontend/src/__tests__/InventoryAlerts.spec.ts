import { describe, it, expect, beforeEach, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import { createPinia, setActivePinia } from 'pinia'
import { createI18n } from 'vue-i18n'
import InventoryAlerts from '@/pages/InventoryAlerts.vue'
import ToastService from 'primevue/toastservice'

vi.mock('@/services/auth', () => ({
  default: {
    accessToken: { value: 'test-token' },
    currentUser: { value: null },
    signOut: vi.fn(),
  },
}))

const mockGet = vi.hoisted(() => vi.fn())
const mockPost = vi.hoisted(() => vi.fn())
const mockDelete = vi.hoisted(() => vi.fn())

vi.mock('axios', () => ({
  default: {
    get: mockGet,
    post: mockPost,
    put: vi.fn(),
    delete: mockDelete,
    defaults: { headers: { common: {} } },
  },
}))

const i18n = createI18n({
  legacy: false,
  locale: 'en',
  fallbackLocale: 'en',
  messages: {
    en: {
      inventory_alerts: 'Inventory Alerts',
      add_rule: 'Add Rule',
      active_rules: 'Active Rules',
      unread_alerts: 'Unread Alerts',
      critical: 'Critical',
      low_stock: 'Low Stock',
      alert_rules: 'Alert Rules',
      no_rules: 'No rules',
      material: 'Material',
      status: 'Status',
      active: 'Active',
      inactive: 'Inactive',
      actions: 'Actions',
    },
  },
})

const stubs = {
  DataTable: { template: '<div class="datatable-stub"><slot name="empty" /><slot /></div>', props: ['value', 'loading'] },
  Column: { template: '<div class="column-stub"><slot name="body" :data="{}" /></div>', props: ['field', 'header'] },
  Button: { template: '<button class="btn-stub">{{ label }}</button>', props: ['label', 'icon', 'severity', 'text', 'rounded', 'size'] },
  Card: { template: '<div class="card-stub"><slot name="title" /><slot name="content" /></div>' },
  Tag: { template: '<span class="tag-stub"><slot /></span>', props: ['severity', 'value'] },
  InputText: { template: '<input class="input-text-stub" />', props: ['modelValue'] },
  InputNumber: { template: '<input class="input-number-stub" />', props: ['modelValue'] },
}

describe('InventoryAlerts', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    mockGet.mockReset()
    mockPost.mockReset()
    mockDelete.mockReset()
    mockGet
      .mockResolvedValueOnce({ data: { data: [] } })
      .mockResolvedValueOnce({ data: { data: [] } })
      .mockResolvedValueOnce({ data: { data: { total_active: 0, unread_count: 0, critical_count: 0, low_count: 0 } } })
  })

  it('renders page heading', async () => {
    const wrapper = mount(InventoryAlerts, { global: { plugins: [i18n, ToastService], stubs } })
    await vi.waitFor(() => expect(wrapper.text()).toContain('Inventory Alerts'))
  })

  it('shows add rule button', async () => {
    const wrapper = mount(InventoryAlerts, { global: { plugins: [i18n, ToastService], stubs } })
    await vi.waitFor(() => expect(wrapper.text()).toContain('Add Rule'))
  })

  it('loads data from API', async () => {
    mount(InventoryAlerts, { global: { plugins: [i18n, ToastService], stubs } })
    await vi.waitFor(() => expect(mockGet).toHaveBeenCalledTimes(3))
  })

  it('shows summary cards when loaded', async () => {
    mockGet.mockReset()
    mockGet
      .mockResolvedValueOnce({ data: { data: [] } })
      .mockResolvedValueOnce({ data: { data: [] } })
      .mockResolvedValueOnce({ data: { data: { total_active: 5, unread_count: 3, critical_count: 1, low_count: 4 } } })
    const wrapper = mount(InventoryAlerts, { global: { plugins: [i18n, ToastService], stubs } })
    await vi.waitFor(() => {
      expect(wrapper.text()).toContain('5')
      expect(wrapper.text()).toContain('3')
    })
  })
})

import { describe, it, expect, beforeEach, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import { createPinia, setActivePinia } from 'pinia'
import { createI18n } from 'vue-i18n'
import PurchaseOrders from '@/pages/PurchaseOrders.vue'
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
      purchase_orders: 'Purchase Orders',
      no_purchase_orders: 'No purchase orders',
      supplier: 'Supplier',
      items: 'Items',
      total: 'Total',
      status: 'Status',
      date: 'Date',
      actions: 'Actions',
      new_purchase_order: 'New Purchase Order',
      expected_date: 'Expected Date',
      notes: 'Notes',
      add_item: 'Add Item',
      price: 'Price',
      create_order: 'Create Order',
    },
  },
})

const stubs = {
  DataTable: {
    template: '<div class="datatable-stub"><slot name="empty" /><slot /></div>',
    props: ['value', 'loading'],
  },
  Column: { template: '<div class="column-stub"><slot name="body" :data="{}" /></div>', props: ['field', 'header'] },
  Button: { template: '<button class="btn-stub" />', props: ['icon', 'severity', 'size', 'loading'] },
  InputText: { template: '<input class="input-text-stub" />', props: ['modelValue'] },
  InputNumber: { template: '<input class="input-number-stub" />', props: ['modelValue', 'min'] },
  Calendar: { template: '<input class="calendar-stub" />', props: ['modelValue', 'dateFormat'] },
  Tag: { template: '<span class="tag-stub"><slot /></span>', props: ['severity', 'value'] },
  Card: { template: '<div class="card-stub"><slot name="title" /><slot name="content" /></div>' },
}

describe('PurchaseOrders', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    mockGet.mockReset()
    mockPost.mockReset()
    mockGet.mockResolvedValue({ data: { data: [] } })
  })

  it('renders page heading', async () => {
    const wrapper = mount(PurchaseOrders, {
      global: { plugins: [i18n, ToastService], stubs },
    })
    await vi.waitFor(() => {
      expect(wrapper.text()).toContain('Purchase Orders')
    })
  })

  it('loads purchase orders from API', async () => {
    mount(PurchaseOrders, {
      global: { plugins: [i18n, ToastService], stubs },
    })
    await vi.waitFor(() => {
      expect(mockGet).toHaveBeenCalled()
    })
  })

  it('shows new purchase order form', async () => {
    const wrapper = mount(PurchaseOrders, {
      global: { plugins: [i18n, ToastService], stubs },
    })
    await vi.waitFor(() => {
      expect(wrapper.text()).toContain('New Purchase Order')
    })
  })

  it('shows empty state when no orders', async () => {
    const wrapper = mount(PurchaseOrders, {
      global: { plugins: [i18n, ToastService], stubs },
    })
    await vi.waitFor(() => {
      expect(wrapper.text()).toContain('No purchase orders')
    })
  })

  it('renders form fields', async () => {
    const wrapper = mount(PurchaseOrders, {
      global: { plugins: [i18n, ToastService], stubs },
    })
    await vi.waitFor(() => {
      expect(wrapper.find('.input-text-stub').exists()).toBe(true)
      expect(wrapper.find('.calendar-stub').exists()).toBe(true)
    })
  })
})

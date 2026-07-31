import { describe, it, expect, beforeEach, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import { createPinia, setActivePinia } from 'pinia'
import { createI18n } from 'vue-i18n'
import DeliveryManagement from '@/pages/DeliveryManagement.vue'
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
const mockPut = vi.hoisted(() => vi.fn())
const mockDelete = vi.hoisted(() => vi.fn())

vi.mock('axios', () => ({
  default: {
    get: mockGet,
    post: mockPost,
    put: mockPut,
    delete: mockDelete,
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
      delivery_management: 'Delivery Management',
      zones: 'Zones',
      no_results: 'No results',
      name: 'Name',
      fee: 'Fee',
      min_order: 'Min Order',
      active: 'Active',
      actions: 'Actions',
      orders: 'Orders',
      customer: 'Customer',
      address: 'Address',
      status: 'Status',
    },
  },
})

const stubs = {
  DataTable: {
    template: '<div class="datatable-stub"><slot name="empty" /><slot /></div>',
    props: ['value', 'loading'],
  },
  Column: { template: '<div class="column-stub"><slot name="body" :data="{}" /></div>', props: ['field', 'header'] },
  Button: { template: '<button class="btn-stub" @click="$emit(\'click\')" />', props: ['icon', 'severity', 'size'] },
  InputText: { template: '<input class="input-text-stub" />', props: ['modelValue'] },
  InputNumber: { template: '<input class="input-number-stub" />', props: ['modelValue'] },
  TabView: { template: '<div class="tabview-stub"><slot /></div>' },
  TabPanel: { template: '<div class="tabpanel-stub">{{ header }}<slot /></div>', props: ['header'] },
  Tag: { template: '<span class="tag-stub"><slot /></span>', props: ['severity', 'value'] },
}

describe('DeliveryManagement', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    mockGet.mockReset()
    mockPost.mockReset()
    mockPut.mockReset()
    mockDelete.mockReset()
    mockGet.mockResolvedValue({ data: { data: [] } })
  })

  it('renders page heading', async () => {
    const wrapper = mount(DeliveryManagement, {
      global: { plugins: [i18n, ToastService], stubs },
    })
    await vi.waitFor(() => {
      expect(wrapper.text()).toContain('Delivery Management')
    })
  })

  it('renders zones and orders tabs', async () => {
    const wrapper = mount(DeliveryManagement, {
      global: { plugins: [i18n, ToastService], stubs },
    })
    await vi.waitFor(() => {
      expect(wrapper.text()).toContain('Zones')
      expect(wrapper.text()).toContain('Orders')
    })
  })

  it('loads zones and orders from API', async () => {
    mount(DeliveryManagement, {
      global: { plugins: [i18n, ToastService], stubs },
    })
    await vi.waitFor(() => {
      expect(mockGet).toHaveBeenCalledTimes(2)
    })
  })

  it('shows empty state when no data', async () => {
    const wrapper = mount(DeliveryManagement, {
      global: { plugins: [i18n, ToastService], stubs },
    })
    await vi.waitFor(() => {
      expect(wrapper.text()).toContain('No results')
    })
  })
})

import { describe, it, expect, beforeEach, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import { createPinia, setActivePinia } from 'pinia'
import { createI18n } from 'vue-i18n'
import InventoryTransfers from '@/pages/InventoryTransfers.vue'
import ToastService from 'primevue/toastservice'

const mockGet = vi.hoisted(() => vi.fn())
const mockPost = vi.hoisted(() => vi.fn())
const mockPut = vi.hoisted(() => vi.fn())

vi.mock('axios', () => ({
  default: {
    get: mockGet,
    post: mockPost,
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
      inventory_transfers: 'Inventory Transfers',
      create_transfer: 'Create Transfer',
      material: 'Material',
      quantity: 'Quantity',
      from_branch: 'From Branch',
      to_branch: 'To Branch',
      status: 'Status',
      date: 'Date',
      actions: 'Actions',
      cancel: 'Cancel',
      create: 'Create',
      unit: 'Unit',
      notes: 'Notes',
    },
  },
})

const stubs = {
  DataTable: { template: '<div class="datatable-stub"><slot name="empty" /><slot /></div>', props: ['value', 'loading'] },
  Column: { template: '<div class="column-stub"><slot name="body" :data="{}" /></div>', props: ['field', 'header'] },
  Button: { template: '<button class="btn-stub">{{ label }}</button>', props: ['label', 'icon', 'severity', 'loading', 'text', 'rounded'] },
  Dialog: { template: '<div class="dialog-stub" v-if="visible"><slot /><slot name="footer" /></div>', props: ['visible', 'header'] },
  InputText: { template: '<input class="input-text-stub" />', props: ['modelValue'] },
  InputNumber: { template: '<input class="input-number-stub" />', props: ['modelValue'] },
  Tag: { template: '<span class="tag-stub"><slot /></span>', props: ['severity', 'value'] },
  ProgressSpinner: { template: '<div class="spinner-stub" />' },
  Toast: { template: '<div />' },
}

describe('InventoryTransfers', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    mockGet.mockReset()
    mockPost.mockReset()
    mockPut.mockReset()
    mockGet.mockResolvedValue({ data: { data: [] } })
  })

  it('renders page heading', async () => {
    const wrapper = mount(InventoryTransfers, { global: { plugins: [i18n, ToastService], stubs } })
    await vi.waitFor(() => expect(wrapper.text()).toContain('Inventory Transfers'))
  })

  it('shows create transfer button', async () => {
    const wrapper = mount(InventoryTransfers, { global: { plugins: [i18n, ToastService], stubs } })
    await vi.waitFor(() => expect(wrapper.text()).toContain('Create Transfer'))
  })

  it('loads transfers from API', async () => {
    mount(InventoryTransfers, { global: { plugins: [i18n, ToastService], stubs } })
    await vi.waitFor(() => expect(mockGet).toHaveBeenCalledTimes(1))
  })

  it('hides loading spinner after load', async () => {
    const wrapper = mount(InventoryTransfers, { global: { plugins: [i18n, ToastService], stubs } })
    await vi.waitFor(() => expect(wrapper.find('.spinner-stub').exists()).toBe(false))
  })

  it('shows create transfer dialog on button click', async () => {
    const wrapper = mount(InventoryTransfers, { global: { plugins: [i18n, ToastService], stubs } })
    await vi.waitFor(() => {
      const btn = wrapper.find('.btn-stub')
      if (btn.exists()) btn.trigger('click')
    })
    expect(wrapper.find('.dialog-stub').exists()).toBe(true)
  })
})

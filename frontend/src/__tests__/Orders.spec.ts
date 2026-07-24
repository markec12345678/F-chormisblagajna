import { describe, it, expect, beforeEach, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import { createPinia, setActivePinia } from 'pinia'
import { createI18n } from 'vue-i18n'
import Orders from '@/pages/Orders.vue'
import ToastService from 'primevue/toastservice'
import ConfirmationService from 'primevue/confirmationservice'

vi.mock('@/stores', () => ({
  globalStore: () => ({
    shop_mode: '',
    orientation: 'ltr',
    getColorMode: 'light',
  }),
}))

vi.mock('@/services/auth', () => ({
  default: {
    accessToken: { value: 'test-token' },
    currentUser: { value: null },
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
  messages: {
    en: {
      order: 'Order',
      search: 'Search',
      id: 'ID',
      status: 'Status',
      submitted_at: 'Submitted At',
      actions: 'Actions',
      no_results: 'No results found',
      confirm_cancel_order: 'Are you sure you want to cancel this order?',
      cancel: 'Cancel',
      yes: 'Yes',
      order_cancelled_success: 'Order cancelled successfully',
      order_cancel_failed: 'Failed to cancel order',
    },
  },
})

const stubs = {
  DataTable: {
    template:
      '<div class="datatable-stub"><slot name="header" /><slot name="empty" /><slot /></div>',
    props: [
      'value',
      'loading',
      'totalRecords',
      'rows',
      'paginator',
      'stripedRows',
      'lazy',
      'paginatorPosition',
      'rowsPerPageOptions',
    ],
  },
  Column: { template: '<div class="column-stub" />', props: ['field', 'header', 'sortable'] },
  Dialog: {
    template: '<div class="dialog-stub" v-if="visible"><slot/></div>',
    props: ['visible', 'modal', 'header', 'style', 'breakpoints'],
  },
  Tag: { template: '<span class="tag-stub">{{ value }}</span>', props: ['value', 'severity'] },
  Button: {
    template: '<button class="btn-stub" @click="$emit(\'click\')">{{ label }}</button>',
    props: ['label', 'icon', 'severity'],
  },
  ButtonGroup: { template: '<div class="btn-group-stub"><slot/></div>' },
  IconField: { template: '<div class="iconfield-stub"><slot/></div>' },
  InputIcon: { template: '<span class="inputicon-stub"><slot/></span>' },
  InputText: {
    template: '<input class="input-text-stub" />',
    props: ['modelValue', 'placeholder'],
  },
  ConfirmPopup: { template: '<div />' },
  OrderView: { template: '<div class="orderview-stub" />', props: ['order'] },
}

describe('Orders', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    mockGet.mockReset()
  })

  it('renders page heading', () => {
    mockGet.mockResolvedValue({ data: { data: [], meta: { total_records: 0 } } })

    const wrapper = mount(Orders, {
      global: { plugins: [i18n, ToastService, ConfirmationService], stubs },
    })
    expect(wrapper.text()).toContain('Order')
  })

  it('shows empty state when no orders', async () => {
    mockGet.mockResolvedValue({ data: { data: [], meta: { total_records: 0 } } })

    const wrapper = mount(Orders, {
      global: { plugins: [i18n, ToastService, ConfirmationService], stubs },
    })

    await vi.waitFor(() => {
      expect(wrapper.text()).toContain('No results found')
    })
  })

  it('renders search input', async () => {
    mockGet.mockResolvedValue({ data: { data: [], meta: { total_records: 0 } } })

    const wrapper = mount(Orders, {
      global: { plugins: [i18n, ToastService, ConfirmationService], stubs },
    })

    await vi.waitFor(() => {
      expect(wrapper.find('.input-text-stub').exists()).toBe(true)
    })
  })
})

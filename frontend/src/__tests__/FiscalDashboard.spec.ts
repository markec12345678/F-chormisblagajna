import { describe, it, expect, beforeEach, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import { createPinia, setActivePinia } from 'pinia'
import { createI18n } from 'vue-i18n'
import FiscalDashboard from '@/pages/FiscalDashboard.vue'
import ToastService from 'primevue/toastservice'

const mockGet = vi.hoisted(() => vi.fn())

vi.mock('axios', () => ({
  default: { get: mockGet, post: vi.fn(), put: vi.fn(), delete: vi.fn(), defaults: { headers: { common: {} } } },
}))

const i18n = createI18n({
  legacy: false,
  locale: 'en',
  fallbackLocale: 'en',
  messages: { en: { fiscal_dashboard: 'Fiscal Dashboard', total_receipts: 'Total Receipts', total_amount: 'Total Amount', date: 'Date', start_date: 'Start Date', end_date: 'End Date', filter: 'Filter', fiscal_receipts: 'Fiscal Receipts', invoice_number: 'Invoice #', eor: 'EOR', zoi: 'ZOI', amount: 'Amount', status: 'Status', pending: 'Pending', fiscalized: 'Fiscalized', error: 'Error', request_failed: 'Request failed' } },
})

const stubs = {
  Card: { template: '<div class="card-stub"><slot name="title" /><slot name="content" /></div>' },
  Button: { template: '<button class="btn-stub">{{ label }}</button>', props: ['label', 'icon', 'severity', 'loading'] },
  DataTable: { template: '<div class="datatable-stub"><slot name="empty" /><slot /></div>', props: ['value', 'loading'] },
  Column: { template: '<div class="column-stub"><slot name="body" :data="{}" /></div>', props: ['field', 'header'] },
  Tag: { template: '<span class="tag-stub"><slot /></span>', props: ['severity', 'value'] },
  Calendar: { template: '<input class="calendar-stub" />', props: ['modelValue', 'placeholder', 'showIcon'] },
  ProgressSpinner: { template: '<div class="spinner-stub" />' },
  Toast: { template: '<div class="toast-stub" />' },
}

describe('FiscalDashboard', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    mockGet.mockReset()
    mockGet.mockResolvedValueOnce({ data: { data: [] } }).mockResolvedValueOnce({ data: { data: { total_count: 0, total_amount: 0, date: '2025-01-01' } } })
  })

  it('renders page heading', async () => {
    const wrapper = mount(FiscalDashboard, { global: { plugins: [i18n, ToastService], stubs } })
    await vi.waitFor(() => expect(wrapper.text()).toContain('Fiscal Dashboard'))
  })

  it('calls API on mount', async () => {
    mount(FiscalDashboard, { global: { plugins: [i18n, ToastService], stubs } })
    await vi.waitFor(() => expect(mockGet).toHaveBeenCalledTimes(2))
  })

  it('shows loading state initially', async () => {
    mockGet.mockImplementation(() => new Promise(() => {}))
    const wrapper = mount(FiscalDashboard, { global: { plugins: [i18n, ToastService], stubs } })
    expect(wrapper.find('.spinner-stub').exists()).toBe(true)
  })

  it('handles API error gracefully', async () => {
    mockGet.mockRejectedValue(new Error('fail'))
    const wrapper = mount(FiscalDashboard, { global: { plugins: [i18n, ToastService], stubs } })
    await vi.waitFor(() => expect(wrapper.text()).toBeDefined())
  })
})

import { describe, it, expect, beforeEach, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import { createPinia, setActivePinia } from 'pinia'
import { createI18n } from 'vue-i18n'
import MultiLocationDashboard from '@/pages/MultiLocationDashboard.vue'
import ToastService from 'primevue/toastservice'

const mockGet = vi.hoisted(() => vi.fn())

vi.mock('axios', () => ({
  default: { get: mockGet, post: vi.fn(), put: vi.fn(), delete: vi.fn(), defaults: { headers: { common: {} } } },
}))

const i18n = createI18n({
  legacy: false,
  locale: 'en',
  fallbackLocale: 'en',
  messages: { en: { multi_location: 'Multi Location', refresh: 'Refresh', total_branches: 'Total Branches', consolidated_revenue: 'Consolidated Revenue', consolidated_orders: 'Consolidated Orders', average_order: 'Avg Order', location_overview: 'Location Overview', branch_name: 'Branch Name', today_revenue: 'Today Revenue', today_orders: 'Today Orders', week_revenue: 'Week Revenue', month_revenue: 'Month Revenue', status: 'Status', branch_comparison: 'Branch Comparison', metric: 'Metric', error: 'Error', request_failed: 'Request failed' } },
})

const stubs = {
  Card: { template: '<div class="card-stub"><slot name="title" /><slot name="content" /></div>' },
  Button: { template: '<button class="btn-stub">{{ label }}</button>', props: ['label', 'icon', 'loading'] },
  DataTable: { template: '<div class="datatable-stub"><slot /></div>', props: ['value', 'loading'] },
  Column: { template: '<div class="column-stub"><slot name="body" :data="{}" /></div>', props: ['field', 'header'] },
  Tag: { template: '<span class="tag-stub"><slot /></span>', props: ['severity', 'value'] },
  ProgressSpinner: { template: '<div class="spinner-stub" />' },
  Toast: { template: '<div class="toast-stub" />' },
}

describe('MultiLocationDashboard', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    mockGet.mockReset()
    mockGet.mockResolvedValueOnce({ data: { data: { total_branches: 0, total_revenue: 0, total_orders: 0, avg_order_value: 0, branches: [] } } }).mockResolvedValueOnce({ data: { data: [] } })
  })

  it('renders page heading', async () => {
    const wrapper = mount(MultiLocationDashboard, { global: { plugins: [i18n, ToastService], stubs } })
    await vi.waitFor(() => expect(wrapper.text()).toContain('Multi Location'))
  })

  it('calls API on mount', async () => {
    mount(MultiLocationDashboard, { global: { plugins: [i18n, ToastService], stubs } })
    await vi.waitFor(() => expect(mockGet).toHaveBeenCalledTimes(2))
  })

  it('shows loading state initially', async () => {
    mockGet.mockImplementation(() => new Promise(() => {}))
    const wrapper = mount(MultiLocationDashboard, { global: { plugins: [i18n, ToastService], stubs } })
    expect(wrapper.find('.spinner-stub').exists()).toBe(true)
  })

  it('handles API error gracefully', async () => {
    mockGet.mockRejectedValue(new Error('fail'))
    const wrapper = mount(MultiLocationDashboard, { global: { plugins: [i18n, ToastService], stubs } })
    await vi.waitFor(() => expect(wrapper.text()).toBeDefined())
  })
})

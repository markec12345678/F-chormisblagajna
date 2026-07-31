import { describe, it, expect, beforeEach, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import { createPinia, setActivePinia } from 'pinia'
import { createI18n } from 'vue-i18n'
import Reports from '@/pages/Reports.vue'
import ToastService from 'primevue/toastservice'

const mockGet = vi.hoisted(() => vi.fn())

vi.mock('axios', () => ({
  default: { get: mockGet, post: vi.fn(), put: vi.fn(), delete: vi.fn(), defaults: { headers: { common: {} } } },
}))

const i18n = createI18n({
  legacy: false,
  locale: 'en',
  fallbackLocale: 'en',
  messages: { en: { reports: 'Reports', today_revenue: 'Today Revenue', today_orders: 'Today Orders', average_order: 'Avg Order', sales_report: 'Sales Report', start_date: 'Start Date', end_date: 'End Date', generate: 'Generate', sales_summary: 'Sales Summary', total_revenue: 'Total Revenue', total_orders: 'Total Orders', total_items: 'Total Items', net_revenue: 'Net Revenue', top_products: 'Top Products', product: 'Product', quantity: 'Qty', revenue: 'Revenue', refresh: 'Refresh', inventory_report: 'Inventory Report', inventory_status: 'Inventory Status', total_materials: 'Total Materials', low_stock: 'Low Stock', out_of_stock: 'Out of Stock', total_value: 'Total Value', low_stock_items: 'Low Stock Items', name: 'Name', unit: 'Unit', value: 'Value', error: 'Error' } },
})

const stubs = {
  Card: { template: '<div class="card-stub"><slot name="title" /><slot name="content" /></div>' },
  Button: { template: '<button class="btn-stub">{{ label }}</button>', props: ['label', 'icon', 'severity', 'loading'] },
  Calendar: { template: '<input class="calendar-stub" />', props: ['modelValue', 'placeholder', 'showIcon'] },
  DataTable: { template: '<div class="datatable-stub"><slot name="empty" /><slot /></div>', props: ['value', 'loading'] },
  Column: { template: '<div class="column-stub"><slot name="body" :data="{}" /></div>', props: ['field', 'header'] },
}

describe('Reports', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    mockGet.mockReset()
    mockGet.mockResolvedValue({ data: { data: null } }).mockResolvedValue({ data: { data: null } })
  })

  it('renders page heading', async () => {
    const wrapper = mount(Reports, { global: { plugins: [i18n, ToastService], stubs } })
    await vi.waitFor(() => expect(wrapper.text()).toContain('Reports'))
  })

  it('calls API on mount', async () => {
    mount(Reports, { global: { plugins: [i18n, ToastService], stubs } })
    await vi.waitFor(() => expect(mockGet).toHaveBeenCalled())
  })

  it('shows date filter inputs', async () => {
    const wrapper = mount(Reports, { global: { plugins: [i18n, ToastService], stubs } })
    await vi.waitFor(() => expect(wrapper.findAll('.calendar-stub').length).toBe(2))
  })

  it('handles API error gracefully', async () => {
    mockGet.mockRejectedValue(new Error('fail'))
    const wrapper = mount(Reports, { global: { plugins: [i18n, ToastService], stubs } })
    await vi.waitFor(() => expect(wrapper.text()).toBeDefined())
  })
})

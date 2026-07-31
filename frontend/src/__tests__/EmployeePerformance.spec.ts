import { describe, it, expect, beforeEach, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import { createPinia, setActivePinia } from 'pinia'
import { createI18n } from 'vue-i18n'
import EmployeePerformance from '@/pages/EmployeePerformance.vue'
import ToastService from 'primevue/toastservice'

vi.mock('@/services/auth', () => ({
  default: {
    accessToken: { value: 'test-token' },
    currentUser: { value: null },
    signOut: vi.fn(),
  },
}))

const mockGet = vi.hoisted(() => vi.fn())

vi.mock('axios', () => ({
  default: {
    get: mockGet,
    post: vi.fn(),
    put: vi.fn(),
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
      employee_performance: 'Employee Performance',
      start_date: 'Start Date',
      end_date: 'End Date',
      load: 'Load',
      total_employees: 'Total Employees',
      total_revenue: 'Total Revenue',
      avg_sales_per_hour: 'Avg Sales/Hour',
      leaderboard: 'Leaderboard',
      select_date_range: 'Select date range',
      rank: 'Rank',
      employee: 'Employee',
      total_sales: 'Total Sales',
      orders: 'Orders',
      avg_order: 'Avg Order',
      sales_per_hour: 'Sales/Hour',
      tips: 'Tips',
      product: 'Product',
      qty: 'Qty',
      revenue: 'Revenue',
      no_data: 'No data',
      top_products_per_employee: 'Top Products',
    },
  },
})

const stubs = {
  DataTable: { template: '<div class="datatable-stub"><slot name="empty" /><slot /></div>', props: ['value', 'loading'] },
  Column: { template: '<div class="column-stub"><slot name="body" :data="{}" /></div>', props: ['field', 'header'] },
  Button: { template: '<button class="btn-stub">{{ label }}</button>', props: ['label', 'icon', 'loading'] },
  Calendar: { template: '<input class="calendar-stub" />', props: ['modelValue', 'placeholder', 'dateFormat'] },
  Card: { template: '<div class="card-stub"><slot name="title" /><slot name="content" /></div>' },
  Tag: { template: '<span class="tag-stub"><slot /></span>', props: ['severity', 'value'] },
}

describe('EmployeePerformance', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    mockGet.mockReset()
    mockGet.mockResolvedValue({ data: { data: null } })
  })

  it('renders page heading', async () => {
    const wrapper = mount(EmployeePerformance, { global: { plugins: [i18n, ToastService], stubs } })
    await vi.waitFor(() => expect(wrapper.text()).toContain('Employee Performance'))
  })

  it('shows select date range prompt when no data', async () => {
    const wrapper = mount(EmployeePerformance, { global: { plugins: [i18n, ToastService], stubs } })
    await vi.waitFor(() => expect(wrapper.text()).toContain('Select date range'))
  })

  it('shows load button', async () => {
    const wrapper = mount(EmployeePerformance, { global: { plugins: [i18n, ToastService], stubs } })
    await vi.waitFor(() => expect(wrapper.text()).toContain('Load'))
  })

  it('calls API on load button click', async () => {
    const wrapper = mount(EmployeePerformance, { global: { plugins: [i18n, ToastService], stubs } })
    await wrapper.find('button').trigger('click')
    await vi.waitFor(() => expect(mockGet).toHaveBeenCalled())
  })
})

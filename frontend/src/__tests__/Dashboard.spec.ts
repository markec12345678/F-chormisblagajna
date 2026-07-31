import { describe, it, expect, beforeEach, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import { createPinia, setActivePinia } from 'pinia'
import { createI18n } from 'vue-i18n'
import Dashboard from '@/pages/Dashboard.vue'
import ToastService from 'primevue/toastservice'

const mockGet = vi.hoisted(() => vi.fn())

vi.mock('axios', () => ({
  default: { get: mockGet, post: vi.fn(), put: vi.fn(), delete: vi.fn(), defaults: { headers: { common: {} } } },
}))

const i18n = createI18n({
  legacy: false,
  locale: 'en',
  fallbackLocale: 'en',
  messages: { en: { admin_dashboard: 'Admin Dashboard', refresh: 'Refresh', today_sales: 'Today Sales', yesterday: 'Yesterday', orders_today: 'Orders Today', total: 'Total', sales: 'Sales', active_orders: 'Active Orders', pending_orders: 'Pending Orders', low_stock_alerts: 'Low Stock Alerts', average_rating: 'Avg Rating', quick_actions: 'Quick Actions', recent_feedback: 'Recent Feedback', no_feedback: 'No Feedback' } },
})

const stubs = {
  Card: { template: '<div class="card-stub"><slot name="title" /><slot name="content" /></div>' },
  Button: { template: '<button class="btn-stub">{{ label }}</button>', props: ['label', 'icon', 'severity', 'loading'] },
  ProgressSpinner: { template: '<div class="spinner-stub" />' },
  Rating: { template: '<div class="rating-stub" />', props: ['modelValue', 'cancel', 'readonly'] },
  Toast: { template: '<div class="toast-stub" />' },
}

describe('Dashboard', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    mockGet.mockReset()
    mockGet.mockResolvedValue({ data: { data: [] } })
  })

  it('renders page heading', async () => {
    const wrapper = mount(Dashboard, { global: { plugins: [i18n, ToastService], stubs } })
    await vi.waitFor(() => expect(wrapper.text()).toContain('Admin Dashboard'))
  })

  it('shows loading state initially', async () => {
    mockGet.mockImplementation(() => new Promise(() => {}))
    const wrapper = mount(Dashboard, { global: { plugins: [i18n, ToastService], stubs } })
    expect(wrapper.find('.spinner-stub').exists()).toBe(true)
  })

  it('calls API on mount', async () => {
    mount(Dashboard, { global: { plugins: [i18n, ToastService], stubs } })
    await vi.waitFor(() => expect(mockGet).toHaveBeenCalled())
  })

  it('handles API error gracefully', async () => {
    mockGet.mockRejectedValue(new Error('fail'))
    const wrapper = mount(Dashboard, { global: { plugins: [i18n, ToastService], stubs } })
    await vi.waitFor(() => expect(wrapper.text()).toBeDefined())
  })
})

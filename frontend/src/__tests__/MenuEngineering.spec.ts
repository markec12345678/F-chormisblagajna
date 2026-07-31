import { describe, it, expect, beforeEach, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import { createPinia, setActivePinia } from 'pinia'
import { createI18n } from 'vue-i18n'
import MenuEngineering from '@/pages/MenuEngineering.vue'
import ToastService from 'primevue/toastservice'

vi.mock('@/services/auth', () => ({
  default: { accessToken: { value: 'test-token' }, currentUser: { value: null }, signOut: vi.fn() },
}))

const mockGet = vi.hoisted(() => vi.fn())

vi.mock('axios', () => ({
  default: { get: mockGet, post: vi.fn(), put: vi.fn(), delete: vi.fn(), defaults: { headers: { common: {} } } },
}))

const i18n = createI18n({
  legacy: false,
  locale: 'en',
  fallbackLocale: 'en',
  messages: { en: { menu_engineering: 'Menu Engineering', start_date: 'Start Date', end_date: 'End Date', analyze: 'Analyze', total_items: 'Total Items', total_revenue: 'Total Revenue', avg_profit: 'Avg Profit', profitability_matrix: 'Profitability Matrix', stars: 'Stars', plowhorses: 'Plowhorses', puzzles: 'Puzzles', dogs: 'Dogs', no_data: 'No data', product: 'Product', sold: 'Sold', profit_per_item: 'Profit/Item', margin: 'Margin', stars_description: 'High profit, high popularity', plowhorses_description: 'Low profit, high popularity', puzzles_description: 'High profit, low popularity', dogs_description: 'Low profit, low popularity', select_date_range: 'Select date range', failed: 'Failed', analysis_failed: 'Analysis failed' } },
})

const stubs = {
  Card: { template: '<div class="card-stub"><slot name="title" /><slot name="content" /></div>' },
  Button: { template: '<button class="btn-stub">{{ label }}</button>', props: ['label', 'icon', 'loading'] },
  Calendar: { template: '<input class="calendar-stub" />', props: ['modelValue', 'placeholder', 'dateFormat'] },
  DataTable: { template: '<div class="datatable-stub"><slot name="empty" /><slot /></div>', props: ['value', 'loading'] },
  Column: { template: '<div class="column-stub"><slot name="body" :data="{}" /></div>', props: ['field', 'header'] },
  Tag: { template: '<span class="tag-stub"><slot /></span>', props: ['severity', 'value'] },
}

describe('MenuEngineering', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    mockGet.mockReset()
  })

  it('renders page heading', async () => {
    const wrapper = mount(MenuEngineering, { global: { plugins: [i18n, ToastService], stubs } })
    await vi.waitFor(() => expect(wrapper.text()).toContain('Menu Engineering'))
  })

  it('shows analyze button', async () => {
    const wrapper = mount(MenuEngineering, { global: { plugins: [i18n, ToastService], stubs } })
    await vi.waitFor(() => expect(wrapper.text()).toContain('Analyze'))
  })

  it('shows date pickers', async () => {
    const wrapper = mount(MenuEngineering, { global: { plugins: [i18n, ToastService], stubs } })
    await vi.waitFor(() => expect(wrapper.findAll('.calendar-stub').length).toBe(2))
  })

  it('shows select date range when no data', async () => {
    const wrapper = mount(MenuEngineering, { global: { plugins: [i18n, ToastService], stubs } })
    await vi.waitFor(() => expect(wrapper.text()).toContain('Select date range'))
  })
})

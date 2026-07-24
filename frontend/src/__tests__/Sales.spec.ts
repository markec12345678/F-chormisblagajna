import { describe, it, expect, beforeEach, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import { createPinia, setActivePinia } from 'pinia'
import { createI18n } from 'vue-i18n'
import Sales from '@/pages/Sales.vue'
import ToastService from 'primevue/toastservice'
import ConfirmationService from 'primevue/confirmationservice'

vi.mock('vue-router', () => ({
  useRouter: () => ({
    push: vi.fn(),
    replace: vi.fn(),
    back: vi.fn(),
  }),
  useRoute: () => ({
    params: {},
    query: {},
  }),
}))

vi.mock('@/stores', () => ({
  globalStore: () => ({
    shop_mode: '',
    orientation: 'ltr',
    getColorMode: 'light',
    getSettings: {},
  }),
}))

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

vi.mock('chart.js', () => ({
  Chart: { register: vi.fn() },
  Title: {},
  Tooltip: {},
  Legend: {},
  BarElement: {},
  CategoryScale: {},
  LinearScale: {},
  PointElement: {},
  LineElement: {},
  ArcElement: {},
  Filler: {},
}))

vi.mock('vue-chartjs', () => ({
  Line: { template: '<div class="line-stub" />', props: ['data', 'chartOptions'] },
  Pie: { template: '<div class="pie-stub" />', props: ['data', 'chartOptions'] },
}))

const i18n = createI18n({
  legacy: false,
  locale: 'en',
  fallbackLocale: 'en',
  messages: {
    en: {
      sales: 'Sales',
      date: 'Date',
      cost: 'Cost',
      refunds: 'Refunds',
      profit: 'Profit',
      export: 'Export',
      id: 'ID',
      has_refunds: 'Has refunds',
      tips: 'Tips',
      returns_to_inventory: 'Returns to Inventory',
    },
  },
})

const stubs = {
  DataTable: {
    template: '<div class="datatable-stub"><slot name="header"/><slot/></div>',
    props: ['value', 'stripedRows', 'paginator', 'rows', 'lazy', 'totalRecords', 'loading'],
    emits: ['page'],
  },
  Column: { template: '<div class="column-stub"/>', props: ['field', 'header', 'sortable'] },
  Button: {
    template: '<button class="btn-stub">{{ label }}</button>',
    props: ['label', 'icon', 'severity', 'raised'],
  },
  Card: { template: '<div class="card-stub"><slot name="content"/></div>' },
  Tag: { template: '<span class="tag-stub">{{ value }}</span>', props: ['value', 'severity'] },
  Badge: { template: '<span class="badge-stub">{{ value }}</span>', props: ['value', 'severity'] },
  Toast: { template: '<div />' },
  Tooltip: { template: '<span class="tooltip-stub"><slot/></span>' },
  SalesLogTableItems: {
    template: '<div class="sales-log-table-items-stub" />',
    props: ['items', 'order_refunds'],
  },
}

describe('Sales', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    mockGet.mockReset()
    mockPost.mockReset()

    mockGet.mockResolvedValue({
      data: {
        data: [],
        meta: { total_records: 0 },
      },
    })
  })

  it('renders the sales page heading', async () => {
    const wrapper = mount(Sales, {
      global: { plugins: [i18n, ToastService, ConfirmationService], stubs },
    })

    await wrapper.vm.$nextTick()
    expect(wrapper.text()).toContain('Sales')
  })

  it('shows the data table area', async () => {
    const wrapper = mount(Sales, {
      global: { plugins: [i18n, ToastService, ConfirmationService], stubs },
    })

    await wrapper.vm.$nextTick()
    expect(wrapper.find('.datatable-stub').exists()).toBe(true)
  })

  it('shows chart containers', async () => {
    const wrapper = mount(Sales, {
      global: { plugins: [i18n, ToastService, ConfirmationService], stubs },
    })

    await wrapper.vm.$nextTick()
    expect(wrapper.find('.card-stub').exists()).toBe(true)
  })

  it('calls API to load sales data on mount', async () => {
    mount(Sales, {
      global: { plugins: [i18n, ToastService, ConfirmationService], stubs },
    })

    await vi.waitFor(() => {
      expect(mockGet).toHaveBeenCalled()
    })
  })

  it('shows export button', async () => {
    const wrapper = mount(Sales, {
      global: { plugins: [i18n, ToastService, ConfirmationService], stubs },
    })

    await wrapper.vm.$nextTick()
    expect(wrapper.find('.btn-stub').exists()).toBe(true)
  })
})

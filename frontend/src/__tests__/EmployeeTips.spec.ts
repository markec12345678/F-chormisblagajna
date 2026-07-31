import { describe, it, expect, beforeEach, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import { createPinia, setActivePinia } from 'pinia'
import { createI18n } from 'vue-i18n'
import EmployeeTips from '@/pages/EmployeeTips.vue'
import ToastService from 'primevue/toastservice'

const mockGet = vi.hoisted(() => vi.fn())
const mockPost = vi.hoisted(() => vi.fn())

vi.mock('axios', () => ({
  default: {
    get: mockGet,
    post: mockPost,
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
      employee_tips: 'Employee Tips',
      payout_tips: 'Payout Tips',
      record_tip: 'Record Tip',
      total_tips: 'Total Tips',
      tip_count: 'Tip Count',
      average_tip: 'Average Tip',
      payout_history: 'Payout History',
      tip_summary: 'Tip Summary',
      employee: 'Employee',
      amount: 'Amount',
      payment_method: 'Payment Method',
      payout_method: 'Payout Method',
      date: 'Date',
      cancel: 'Cancel',
      payout: 'Payout',
    },
  },
})

const stubs = {
  DataTable: { template: '<div class="datatable-stub"><slot name="empty" /><slot /></div>', props: ['value', 'loading'] },
  Column: { template: '<div class="column-stub"><slot name="body" :data="{}" /></div>', props: ['field', 'header'] },
  Button: { template: '<button class="btn-stub">{{ label }}</button>', props: ['label', 'icon', 'severity', 'loading'] },
  Card: { template: '<div class="card-stub"><slot name="title" /><slot name="content" /></div>' },
  Dialog: { template: '<div class="dialog-stub" v-if="visible"><slot /><slot name="footer" /></div>', props: ['visible', 'header'] },
  InputText: { template: '<input class="input-text-stub" />', props: ['modelValue'] },
  InputNumber: { template: '<input class="input-number-stub" />', props: ['modelValue'] },
  Dropdown: { template: '<select class="dropdown-stub" />', props: ['modelValue', 'options'] },
  ProgressSpinner: { template: '<div class="spinner-stub" />' },
  Toast: { template: '<div />' },
}

describe('EmployeeTips', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    mockGet.mockReset()
    mockPost.mockReset()
    mockGet.mockResolvedValue({ data: { data: [] } })
  })

  it('renders page heading', async () => {
    const wrapper = mount(EmployeeTips, { global: { plugins: [i18n, ToastService], stubs } })
    await vi.waitFor(() => expect(wrapper.text()).toContain('Employee Tips'))
  })

  it('shows payout and record buttons', async () => {
    const wrapper = mount(EmployeeTips, { global: { plugins: [i18n, ToastService], stubs } })
    await vi.waitFor(() => {
      expect(wrapper.text()).toContain('Payout Tips')
      expect(wrapper.text()).toContain('Record Tip')
    })
  })

  it('loads data from API', async () => {
    mount(EmployeeTips, { global: { plugins: [i18n, ToastService], stubs } })
    await vi.waitFor(() => expect(mockGet).toHaveBeenCalledTimes(2))
  })

  it('shows summary cards when data loaded', async () => {
    mockGet.mockResolvedValue({
      data: {
        data: [
          { employee_name: 'Ana', total_tips: 150, tip_count: 5, average_tip: 30 },
        ],
      },
    })
    const wrapper = mount(EmployeeTips, { global: { plugins: [i18n, ToastService], stubs } })
    await vi.waitFor(() => {
      expect(wrapper.text()).toContain('150')
      expect(wrapper.text()).toContain('30')
    })
  })
})

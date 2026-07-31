import { describe, it, expect, beforeEach, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import { createPinia, setActivePinia } from 'pinia'
import { createI18n } from 'vue-i18n'
import Expenses from '@/pages/Expenses.vue'
import ToastService from 'primevue/toastservice'

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
    defaults: { headers: { common: {} } },
  },
}))

const i18n = createI18n({
  legacy: false,
  locale: 'en',
  fallbackLocale: 'en',
  messages: {
    en: {
      expense_tracking: 'Expense Tracking',
      start_date: 'Start Date',
      end_date: 'End Date',
      add_expense: 'Add Expense',
      total_expenses: 'Total Expenses',
      monthly_total: 'Monthly Total',
      categories: 'Categories',
      no_results: 'No results',
      name: 'Name',
      amount: 'Amount',
      date: 'Date',
      actions: 'Actions',
      save: 'Save',
      cancel: 'Cancel',
    },
  },
})

const stubs = {
  DataTable: { template: '<div class="datatable-stub"><slot name="empty" /><slot /></div>', props: ['value', 'loading'] },
  Column: { template: '<div class="column-stub"><slot name="body" :data="{}" /></div>', props: ['field', 'header'] },
  Button: { template: '<button class="btn-stub">{{ label }}</button>', props: ['label', 'icon', 'severity', 'size'] },
  Calendar: { template: '<input class="calendar-stub" />', props: ['modelValue', 'placeholder', 'dateFormat'] },
  Card: { template: '<div class="card-stub"><slot name="title" /><slot name="content" /></div>' },
  Dialog: { template: '<div class="dialog-stub" v-if="visible"><slot /></div>', props: ['visible', 'header'] },
  InputText: { template: '<input class="input-text-stub" />', props: ['modelValue'] },
  InputNumber: { template: '<input class="input-number-stub" />', props: ['modelValue'] },
  Select: { template: '<select class="select-stub" />', props: ['modelValue', 'options'] },
  ButtonGroup: { template: '<div class="btn-group-stub"><slot /></div>' },
}

describe('Expenses', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    mockGet.mockReset()
    mockPost.mockReset()
    mockGet.mockResolvedValue({ data: { data: [] } })
  })

  it('renders page heading', async () => {
    const wrapper = mount(Expenses, { global: { plugins: [i18n, ToastService], stubs } })
    await vi.waitFor(() => expect(wrapper.text()).toContain('Expense Tracking'))
  })

  it('shows add expense button', async () => {
    const wrapper = mount(Expenses, { global: { plugins: [i18n, ToastService], stubs } })
    await vi.waitFor(() => expect(wrapper.text()).toContain('Add Expense'))
  })

  it('loads data from API', async () => {
    mount(Expenses, { global: { plugins: [i18n, ToastService], stubs } })
    await vi.waitFor(() => expect(mockGet).toHaveBeenCalled())
  })
})

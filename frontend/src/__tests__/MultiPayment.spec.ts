import { describe, it, expect, beforeEach, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import { createPinia, setActivePinia } from 'pinia'
import { createI18n } from 'vue-i18n'
import MultiPayment from '@/pages/MultiPayment.vue'
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
  messages: { en: { multi_payment: 'Multi Payment', start_date: 'Start Date', end_date: 'End Date', payment_summary: 'Payment Summary', grand_total: 'Grand Total', cash: 'Cash', card: 'Card', voucher: 'Voucher', mobile_pay: 'Mobile Pay', gift_cards: 'Gift Cards', daily_breakdown: 'Daily Breakdown', no_payments: 'No payments', date: 'Date', transactions: 'Transactions', total: 'Total', failed: 'Failed', load_failed: 'Load failed' } },
})

const stubs = {
  Card: { template: '<div class="card-stub"><slot name="title" /><slot name="content" /></div>' },
  Button: { template: '<button class="btn-stub">{{ label }}</button>', props: ['label', 'icon', 'severity', 'loading'] },
  Calendar: { template: '<input class="calendar-stub" />', props: ['modelValue', 'placeholder', 'dateFormat'] },
  DataTable: { template: '<div class="datatable-stub"><slot name="empty" /><slot /></div>', props: ['value', 'loading'] },
  Column: { template: '<div class="column-stub"><slot name="body" :data="{}" /></div>', props: ['field', 'header'] },
}

describe('MultiPayment', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    mockGet.mockReset()
    mockGet.mockResolvedValue({ data: { data: [] } })
  })

  it('renders page heading', async () => {
    const wrapper = mount(MultiPayment, { global: { plugins: [i18n, ToastService], stubs } })
    await vi.waitFor(() => expect(wrapper.text()).toContain('Multi Payment'))
  })

  it('loads payment data from API', async () => {
    mount(MultiPayment, { global: { plugins: [i18n, ToastService], stubs } })
    await vi.waitFor(() => expect(mockGet).toHaveBeenCalled())
  })

  it('shows date filter inputs', async () => {
    const wrapper = mount(MultiPayment, { global: { plugins: [i18n, ToastService], stubs } })
    await vi.waitFor(() => expect(wrapper.findAll('.calendar-stub').length).toBe(2))
  })

  it('shows empty state', async () => {
    const wrapper = mount(MultiPayment, { global: { plugins: [i18n, ToastService], stubs } })
    await vi.waitFor(() => expect(wrapper.text()).toContain('No payments'))
  })
})

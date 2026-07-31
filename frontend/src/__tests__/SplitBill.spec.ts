import { describe, it, expect, beforeEach, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import { createPinia, setActivePinia } from 'pinia'
import { createI18n } from 'vue-i18n'
import SplitBill from '@/pages/SplitBill.vue'
import ToastService from 'primevue/toastservice'

const mockGet = vi.hoisted(() => vi.fn())
const mockPost = vi.hoisted(() => vi.fn())

vi.mock('axios', () => ({
  default: { get: mockGet, post: mockPost, put: vi.fn(), delete: vi.fn(), defaults: { headers: { common: {} } } },
}))

const i18n = createI18n({
  legacy: false,
  locale: 'en',
  fallbackLocale: 'en',
  messages: { en: { split_bill: 'Split Bill', create_split_bill: 'Create Split Bill', no_results: 'No results', order_id: 'Order ID', split_type: 'Split Type', status: 'Status', parts: 'Parts', paid: 'Paid', actions: 'Actions', view: 'View', pay: 'Pay', split_bill_details: 'Split Bill Details', order: 'Order', amount: 'Amount', payment_method: 'Payment Method', pending: 'Pending', pay_split_part: 'Pay Split Part', select_part: 'Select Part', cancel: 'Cancel', create: 'Create', save: 'Save', equal: 'Equal', custom: 'Custom', by_item: 'By Item', cash: 'Cash', card: 'Card', voucher: 'Voucher', mobile_pay: 'Mobile Pay', success: 'Success', split_bill_created: 'Split bill created', payment_recorded: 'Payment recorded', error: 'Error', validation_required: 'Required' } },
})

const stubs = {
  DataTable: { template: '<div class="datatable-stub"><slot name="header" /><slot name="empty" /><slot /></div>', props: ['value', 'loading', 'totalRecords', 'rows', 'rowsPerPageOptions', 'paginator', 'lazy'] },
  Column: { template: '<div class="column-stub"><slot name="body" :data="{}" /></div>', props: ['field', 'header'] },
  Button: { template: '<button class="btn-stub">{{ label }}</button>', props: ['label', 'icon', 'rounded', 'raised', 'severity', 'loading', 'disabled', 'aria-label'] },
  ButtonGroup: { template: '<div class="btn-group-stub"><slot /></div>' },
  Tag: { template: '<span class="tag-stub"><slot /></span>', props: ['severity', 'value'] },
  Dialog: { template: '<div class="dialog-stub" v-if="visible"><slot /><slot name="footer" /></div>', props: ['visible', 'header', 'modal'] },
  InputText: { template: '<input class="input-text-stub" />', props: ['modelValue'] },
  InputNumber: { template: '<input class="input-number-stub" />', props: ['modelValue', 'min', 'max'] },
  Dropdown: { template: '<select class="dropdown-stub" />', props: ['modelValue', 'options', 'optionLabel', 'optionValue', 'placeholder'] },
}

describe('SplitBill', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    mockGet.mockReset()
    mockPost.mockReset()
    mockGet.mockResolvedValue({ data: { data: [], meta: { total_records: 0 } } })
  })

  it('renders page heading', async () => {
    const wrapper = mount(SplitBill, { global: { plugins: [i18n, ToastService], stubs } })
    await vi.waitFor(() => expect(wrapper.text()).toContain('Split Bill'))
  })

  it('shows create split bill button', async () => {
    const wrapper = mount(SplitBill, { global: { plugins: [i18n, ToastService], stubs } })
    await vi.waitFor(() => expect(wrapper.text()).toContain('Create Split Bill'))
  })

  it('loads split bills from API', async () => {
    mount(SplitBill, { global: { plugins: [i18n, ToastService], stubs } })
    await vi.waitFor(() => expect(mockGet).toHaveBeenCalled())
  })

  it('shows empty state', async () => {
    const wrapper = mount(SplitBill, { global: { plugins: [i18n, ToastService], stubs } })
    await vi.waitFor(() => expect(wrapper.text()).toContain('No results'))
  })
})

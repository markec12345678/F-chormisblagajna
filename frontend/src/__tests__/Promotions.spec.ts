import { describe, it, expect, beforeEach, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import { createPinia, setActivePinia } from 'pinia'
import { createI18n } from 'vue-i18n'
import Promotions from '@/pages/Promotions.vue'
import ToastService from 'primevue/toastservice'
import ConfirmationService from 'primevue/confirmationservice'

const mockGet = vi.hoisted(() => vi.fn())
const mockPost = vi.hoisted(() => vi.fn())
const mockPatch = vi.hoisted(() => vi.fn())
const mockDelete = vi.hoisted(() => vi.fn())

vi.mock('axios', () => ({
  default: { get: mockGet, post: mockPost, patch: mockPatch, delete: mockDelete, put: vi.fn(), defaults: { headers: { common: {} } } },
}))

const i18n = createI18n({
  legacy: false,
  locale: 'en',
  fallbackLocale: 'en',
  messages: { en: { promotions: 'Promotions', add_promotion: 'Add Promotion', no_results: 'No results', name: 'Name', code: 'Code', type: 'Type', value: 'Value', start_date: 'Start Date', end_date: 'End Date', status: 'Status', active: 'Active', inactive: 'Inactive', actions: 'Actions', edit_promotion: 'Edit Promotion', save: 'Save', cancel: 'Cancel', success: 'Success', error: 'Error', promotion_added: 'Promotion added', promotion_updated: 'Promotion updated', promotion_deleted: 'Promotion deleted', do_you_confirm: 'Confirm?', delete: 'Delete', min_order: 'Min Order' } },
})

const stubs = {
  DataTable: { template: '<div class="datatable-stub"><slot name="header" /><slot name="empty" /><slot /><slot name="body" /></div>', props: ['value', 'loading', 'totalRecords', 'rows', 'rowsPerPageOptions', 'paginator', 'lazy'] },
  Column: { template: '<div class="column-stub"><slot name="body" :data="{}" /></div>', props: ['field', 'header'] },
  Button: { template: '<button class="btn-stub">{{ label }}</button>', props: ['label', 'icon', 'rounded', 'raised', 'severity', 'loading', 'disabled'] },
  ButtonGroup: { template: '<div class="btn-group-stub"><slot /></div>' },
  Tag: { template: '<span class="tag-stub"><slot /></span>', props: ['severity', 'value'] },
  Dialog: { template: '<div class="dialog-stub" v-if="visible"><slot /><slot name="footer" /></div>', props: ['visible', 'header', 'modal'] },
  InputText: { template: '<input class="input-text-stub" />', props: ['modelValue'] },
  InputNumber: { template: '<input class="input-number-stub" />', props: ['modelValue', 'min'] },
  Dropdown: { template: '<select class="dropdown-stub" />', props: ['modelValue', 'options', 'optionLabel', 'optionValue'] },
  Calendar: { template: '<input class="calendar-stub" />', props: ['modelValue', 'dateFormat', 'showIcon'] },
  Checkbox: { template: '<input class="checkbox-stub" type="checkbox" />', props: ['modelValue', 'binary', 'inputId'] },
  ConfirmPopup: { template: '<div class="confirm-popup-stub" />' },
}

describe('Promotions', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    mockGet.mockReset()
    mockPost.mockReset()
    mockPatch.mockReset()
    mockDelete.mockReset()
    mockGet.mockResolvedValue({ data: { data: [], meta: { total_records: 0 } } })
  })

  it('renders page heading', async () => {
    const wrapper = mount(Promotions, { global: { plugins: [i18n, ToastService, ConfirmationService], stubs } })
    await vi.waitFor(() => expect(wrapper.text()).toContain('Promotions'))
  })

  it('loads promotions from API', async () => {
    mount(Promotions, { global: { plugins: [i18n, ToastService, ConfirmationService], stubs } })
    await vi.waitFor(() => expect(mockGet).toHaveBeenCalled())
  })

  it('shows add promotion button', async () => {
    const wrapper = mount(Promotions, { global: { plugins: [i18n, ToastService, ConfirmationService], stubs } })
    await vi.waitFor(() => expect(wrapper.text()).toContain('Add Promotion'))
  })

  it('shows empty state', async () => {
    const wrapper = mount(Promotions, { global: { plugins: [i18n, ToastService, ConfirmationService], stubs } })
    await vi.waitFor(() => expect(wrapper.text()).toContain('No results'))
  })
})

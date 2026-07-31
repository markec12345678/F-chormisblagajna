import { describe, it, expect, beforeEach, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import { createPinia, setActivePinia } from 'pinia'
import { createI18n } from 'vue-i18n'
import ReceiptCustomization from '@/pages/ReceiptCustomization.vue'
import ToastService from 'primevue/toastservice'

vi.mock('@/services/auth', () => ({
  default: { accessToken: { value: 'test-token' }, currentUser: { value: null }, signOut: vi.fn() },
}))

const mockGet = vi.hoisted(() => vi.fn())
const mockPost = vi.hoisted(() => vi.fn())
const mockDelete = vi.hoisted(() => vi.fn())

vi.mock('axios', () => ({
  default: { get: mockGet, post: mockPost, delete: mockDelete, put: vi.fn(), patch: vi.fn(), defaults: { headers: { common: {} } } },
}))

const i18n = createI18n({
  legacy: false,
  locale: 'en',
  fallbackLocale: 'en',
  messages: { en: { receipt_customization: 'Receipt Customization', add_template: 'Add Template', receipt_templates: 'Receipt Templates', no_templates: 'No templates', name: 'Name', default: 'Default', paper_width: 'Paper Width', actions: 'Actions', print_settings: 'Print Settings', printer_name: 'Printer Name', printer_ip: 'Printer IP', print_copies: 'Print Copies', auto_print: 'Auto Print', save_settings: 'Save Settings', template_editor: 'Template Editor', business_name: 'Business Name', business_address: 'Address', business_phone: 'Phone', business_tax_id: 'Tax ID', header: 'Header', footer: 'Footer', show_logo: 'Show Logo', show_tax_id: 'Show Tax ID', show_qr_code: 'Show QR Code', show_server: 'Show Server', show_table: 'Show Table', save_template: 'Save Template', preview: 'Preview', receipt_preview: 'Receipt Preview', close: 'Close', create: 'Create', cancel: 'Cancel', success: 'Success', template_created: 'Template created', template_saved: 'Template saved', template_deleted: 'Template deleted', settings_saved: 'Settings saved', failed: 'Failed', save_failed: 'Save failed', delete_failed: 'Delete failed', load_failed: 'Load failed' } },
})

const stubs = {
  Card: { template: '<div class="card-stub"><slot name="title" /><slot name="content" /></div>' },
  Button: { template: '<button class="btn-stub">{{ label }}</button>', props: ['label', 'icon', 'severity', 'text', 'rounded', 'size', 'loading'] },
  DataTable: { template: '<div class="datatable-stub"><slot name="empty" /><slot /></div>', props: ['value', 'loading', 'selectionMode'] },
  Column: { template: '<div class="column-stub"><slot name="body" :data="{}" /></div>', props: ['field', 'header'] },
  InputText: { template: '<input class="input-text-stub" />', props: ['modelValue', 'placeholder'] },
  InputNumber: { template: '<input class="input-number-stub" />', props: ['modelValue', 'min', 'max'] },
  Textarea: { template: '<textarea class="textarea-stub" />', props: ['modelValue'] },
  Dropdown: { template: '<select class="dropdown-stub" />', props: ['modelValue', 'options', 'optionLabel', 'optionValue'] },
  Dialog: { template: '<div class="dialog-stub" v-if="visible"><slot /><slot name="footer" /></div>', props: ['visible', 'header', 'modal'] },
  Divider: { template: '<hr class="divider-stub" />' },
  Checkbox: { template: '<input class="checkbox-stub" type="checkbox" />', props: ['modelValue', 'binary'] },
  Tag: { template: '<span class="tag-stub"><slot /></span>', props: ['severity', 'value'] },
}

describe('ReceiptCustomization', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    mockGet.mockReset()
    mockPost.mockReset()
    mockDelete.mockReset()
    mockGet.mockResolvedValueOnce({ data: { data: [] } }).mockResolvedValueOnce({ data: { data: { printer_name: '', printer_ip: '', auto_print: true, print_copies: 1 } } })
  })

  it('renders page heading', async () => {
    const wrapper = mount(ReceiptCustomization, { global: { plugins: [i18n, ToastService], stubs } })
    await vi.waitFor(() => expect(wrapper.text()).toContain('Receipt Customization'))
  })

  it('calls API on mount', async () => {
    mount(ReceiptCustomization, { global: { plugins: [i18n, ToastService], stubs } })
    await vi.waitFor(() => expect(mockGet).toHaveBeenCalledTimes(2))
  })

  it('shows add template button', async () => {
    const wrapper = mount(ReceiptCustomization, { global: { plugins: [i18n, ToastService], stubs } })
    await vi.waitFor(() => expect(wrapper.text()).toContain('Add Template'))
  })

  it('shows empty state', async () => {
    const wrapper = mount(ReceiptCustomization, { global: { plugins: [i18n, ToastService], stubs } })
    await vi.waitFor(() => expect(wrapper.text()).toContain('No templates'))
  })
})

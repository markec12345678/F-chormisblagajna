import { describe, it, expect, beforeEach, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import { createPinia, setActivePinia } from 'pinia'
import { createI18n } from 'vue-i18n'
import Branches from '@/pages/Branches.vue'
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
  messages: { en: { branches: 'Branches', add_branch: 'Add Branch', branch_name: 'Branch Name', branch_address: 'Address', branch_phone: 'Phone', branch_email: 'Email', status: 'Status', active: 'Active', inactive: 'Inactive', actions: 'Actions', edit: 'Edit', remove: 'Remove', no_results: 'No results', edit_branch: 'Edit Branch', save: 'Save', cancel: 'Cancel', name: 'Name', success: 'Success', branch_added: 'Branch added', branch_updated: 'Branch updated', branch_deleted: 'Branch deleted', error: 'Error', do_you_confirm: 'Confirm?', delete: 'Delete', validation_required: 'Required' } },
})

const stubs = {
  DataTable: { template: '<div class="datatable-stub"><slot name="header" /><slot name="empty" /><slot /></div>', props: ['value', 'loading', 'totalRecords', 'rows', 'rowsPerPageOptions', 'paginator', 'lazy'] },
  Column: { template: '<div class="column-stub"><slot name="body" :data="{}" /></div>', props: ['field', 'header'] },
  Button: { template: '<button class="btn-stub">{{ label }}</button>', props: ['label', 'icon', 'rounded', 'raised', 'severity', 'loading', 'disabled', 'aria-label'] },
  ButtonGroup: { template: '<div class="btn-group-stub"><slot /></div>' },
  Tag: { template: '<span class="tag-stub"><slot /></span>', props: ['severity', 'value'] },
  Dialog: { template: '<div class="dialog-stub" v-if="visible"><slot /><slot name="footer" /></div>', props: ['visible', 'header', 'modal'] },
  InputText: { template: '<input class="input-text-stub" />', props: ['modelValue'] },
  Checkbox: { template: '<input class="checkbox-stub" type="checkbox" />', props: ['modelValue', 'binary', 'inputId'] },
  ConfirmPopup: { template: '<div class="confirm-popup-stub" />' },
}

describe('Branches', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    mockGet.mockReset()
    mockPost.mockReset()
    mockPatch.mockReset()
    mockDelete.mockReset()
    mockGet.mockResolvedValue({ data: { data: [], meta: { total_records: 0 } } })
  })

  it('renders page heading', async () => {
    const wrapper = mount(Branches, { global: { plugins: [i18n, ToastService, ConfirmationService], stubs } })
    await vi.waitFor(() => expect(wrapper.text()).toContain('Branches'))
  })

  it('loads branches from API', async () => {
    mount(Branches, { global: { plugins: [i18n, ToastService, ConfirmationService], stubs } })
    await vi.waitFor(() => expect(mockGet).toHaveBeenCalled())
  })

  it('shows add branch button', async () => {
    const wrapper = mount(Branches, { global: { plugins: [i18n, ToastService, ConfirmationService], stubs } })
    await vi.waitFor(() => expect(wrapper.text()).toContain('Add Branch'))
  })

  it('shows empty state', async () => {
    const wrapper = mount(Branches, { global: { plugins: [i18n, ToastService, ConfirmationService], stubs } })
    await vi.waitFor(() => expect(wrapper.text()).toContain('No results'))
  })
})

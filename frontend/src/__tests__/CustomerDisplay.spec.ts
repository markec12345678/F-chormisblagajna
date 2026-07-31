import { describe, it, expect, beforeEach, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import { createPinia, setActivePinia } from 'pinia'
import { createI18n } from 'vue-i18n'
import CustomerDisplay from '@/pages/CustomerDisplay.vue'
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
  messages: { en: { customer_display: 'Customer Display', add_display: 'Add Display', display_configs: 'Display Configs', no_displays: 'No displays', name: 'Name', active: 'Active', inactive: 'Inactive', slide_interval: 'Slide Interval', theme: 'Theme', actions: 'Actions', edit_display: 'Edit Display', display_name: 'Display Name', welcome_message: 'Welcome Message', show_promotions: 'Show Promotions', show_menu: 'Show Menu', show_order_status: 'Show Order Status', show_wait_time: 'Show Wait Time', save: 'Save', cancel: 'Cancel', confirm_delete: 'Confirm Delete', delete_display_confirm: 'Delete this display?', delete: 'Delete', saved: 'Saved', failed: 'Failed', save_failed: 'Save failed', deleted: 'Deleted', delete_failed: 'Delete failed', load_failed: 'Load failed' } },
})

const stubs = {
  Card: { template: '<div class="card-stub"><slot name="title" /><slot name="content" /></div>' },
  Button: { template: '<button class="btn-stub">{{ label }}</button>', props: ['label', 'icon', 'severity', 'class', 'loading', 'disabled'] },
  DataTable: { template: '<div class="datatable-stub"><slot name="empty" /><slot /></div>', props: ['value', 'loading'] },
  Column: { template: '<div class="column-stub"><slot name="body" :data="{}" /></div>', props: ['field', 'header'] },
  Dialog: { template: '<div class="dialog-stub" v-if="visible"><slot /><slot name="footer" /></div>', props: ['visible', 'header'] },
  InputText: { template: '<input class="input-text-stub" />', props: ['modelValue'] },
  InputNumber: { template: '<input class="input-number-stub" />', props: ['modelValue', 'min', 'max'] },
  Dropdown: { template: '<select class="dropdown-stub" />', props: ['modelValue', 'options'] },
  Checkbox: { template: '<input class="checkbox-stub" type="checkbox" />', props: ['modelValue', 'binary', 'inputId'] },
  Tag: { template: '<span class="tag-stub"><slot /></span>', props: ['severity'] },
}

describe('CustomerDisplay', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    mockGet.mockReset()
    mockPost.mockReset()
    mockDelete.mockReset()
    mockGet.mockResolvedValue({ data: { data: [] } })
  })

  it('renders page heading', async () => {
    const wrapper = mount(CustomerDisplay, { global: { plugins: [i18n, ToastService], stubs } })
    await vi.waitFor(() => expect(wrapper.text()).toContain('Customer Display'))
  })

  it('loads configs from API', async () => {
    mount(CustomerDisplay, { global: { plugins: [i18n, ToastService], stubs } })
    await vi.waitFor(() => expect(mockGet).toHaveBeenCalled())
  })

  it('shows add display button', async () => {
    const wrapper = mount(CustomerDisplay, { global: { plugins: [i18n, ToastService], stubs } })
    await vi.waitFor(() => expect(wrapper.text()).toContain('Add Display'))
  })

  it('shows empty state in table', async () => {
    const wrapper = mount(CustomerDisplay, { global: { plugins: [i18n, ToastService], stubs } })
    await vi.waitFor(() => expect(wrapper.text()).toContain('No displays'))
  })
})

import { describe, it, expect, beforeEach, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import { createPinia, setActivePinia } from 'pinia'
import { createI18n } from 'vue-i18n'
import TablesideOrdering from '@/pages/TablesideOrdering.vue'
import ToastService from 'primevue/toastservice'

vi.mock('@/services/auth', () => ({
  default: { accessToken: { value: 'test-token' }, currentUser: { value: null }, signOut: vi.fn() },
}))

const mockGet = vi.hoisted(() => vi.fn())
const mockPost = vi.hoisted(() => vi.fn())

vi.mock('axios', () => ({
  default: { get: mockGet, post: mockPost, put: vi.fn(), delete: vi.fn(), defaults: { headers: { common: {} } } },
}))

const i18n = createI18n({
  legacy: false,
  locale: 'en',
  fallbackLocale: 'en',
  messages: { en: { tableside_ordering: 'Tableside Ordering', add_table: 'Add Table', table_sessions: 'Table Sessions', no_sessions: 'No sessions', table: 'Table', zone: 'Zone', guests: 'Guests', status: 'Status', active: 'Active', closed: 'Closed', opened: 'Opened', actions: 'Actions', qr_code: 'QR Code', view_orders: 'View Orders', close: 'Close', table_label: 'Table Label', guest_count: 'Guest Count', save: 'Save', cancel: 'Cancel', orders: 'Orders', no_orders: 'No orders', items: 'Items', total: 'Total', time: 'Time', saved: 'Saved', failed: 'Failed', load_failed: 'Load failed' } },
})

const stubs = {
  Card: { template: '<div class="card-stub"><slot name="title" /><slot name="content" /></div>' },
  Button: { template: '<button class="btn-stub">{{ label }}</button>', props: ['label', 'icon', 'severity', 'loading', 'disabled'] },
  DataTable: { template: '<div class="datatable-stub"><slot name="empty" /><slot /></div>', props: ['value', 'loading'] },
  Column: { template: '<div class="column-stub"><slot name="body" :data="{}" /></div>', props: ['field', 'header'] },
  Dialog: { template: '<div class="dialog-stub" v-if="visible"><slot /><slot name="footer" /></div>', props: ['visible', 'header'] },
  InputText: { template: '<input class="input-text-stub" />', props: ['modelValue'] },
  InputNumber: { template: '<input class="input-number-stub" />', props: ['modelValue', 'min'] },
  Dropdown: { template: '<select class="dropdown-stub" />', props: ['modelValue', 'options'] },
  Tag: { template: '<span class="tag-stub"><slot /></span>', props: ['severity', 'value'] },
}

describe('TablesideOrdering', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    mockGet.mockReset()
    mockPost.mockReset()
    mockGet.mockResolvedValue({ data: { data: [] } })
  })

  it('renders page heading', async () => {
    const wrapper = mount(TablesideOrdering, { global: { plugins: [i18n, ToastService], stubs } })
    await vi.waitFor(() => expect(wrapper.text()).toContain('Tableside Ordering'))
  })

  it('loads sessions from API', async () => {
    mount(TablesideOrdering, { global: { plugins: [i18n, ToastService], stubs } })
    await vi.waitFor(() => expect(mockGet).toHaveBeenCalled())
  })

  it('shows add table button', async () => {
    const wrapper = mount(TablesideOrdering, { global: { plugins: [i18n, ToastService], stubs } })
    await vi.waitFor(() => expect(wrapper.text()).toContain('Add Table'))
  })

  it('shows empty state', async () => {
    const wrapper = mount(TablesideOrdering, { global: { plugins: [i18n, ToastService], stubs } })
    await vi.waitFor(() => expect(wrapper.text()).toContain('No sessions'))
  })
})

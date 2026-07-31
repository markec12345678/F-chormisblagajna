import { describe, it, expect, beforeEach, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import { createPinia, setActivePinia } from 'pinia'
import { createI18n } from 'vue-i18n'
import Loyalty from '@/pages/Loyalty.vue'
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
      loyalty_program: 'Loyalty Program',
      add_account: 'Add Account',
      all_tiers: 'All Tiers',
      no_results: 'No results',
      customer: 'Customer',
      points: 'Points',
      tier: 'Tier',
      total_spent: 'Total Spent',
      actions: 'Actions',
    },
  },
})

const stubs = {
  DataTable: { template: '<div class="datatable-stub"><slot name="header" /><slot name="empty" /><slot /></div>', props: ['value', 'loading', 'totalRecords'] },
  Column: { template: '<div class="column-stub"><slot name="body" :data="{}" /></div>', props: ['field', 'header'] },
  Button: { template: '<button class="btn-stub">{{ label }}</button>', props: ['label', 'icon', 'rounded', 'raised', 'severity', 'aria-label'] },
  ButtonGroup: { template: '<div class="btn-group-stub"><slot /></div>' },
  Tag: { template: '<span class="tag-stub"><slot /></span>', props: ['severity', 'value'] },
  Dropdown: { template: '<select class="dropdown-stub" />', props: ['modelValue', 'options', 'placeholder'] },
  Dialog: { template: '<div class="dialog-stub" v-if="visible"><slot /><slot name="footer" /></div>', props: ['visible', 'header'] },
  InputText: { template: '<input class="input-text-stub" />', props: ['modelValue'] },
  InputNumber: { template: '<input class="input-number-stub" />', props: ['modelValue'] },
}

describe('Loyalty', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    mockGet.mockReset()
    mockPost.mockReset()
    mockGet.mockResolvedValue({ data: { data: [], meta: { total_records: 0 } } })
  })

  it('renders page heading', async () => {
    const wrapper = mount(Loyalty, { global: { plugins: [i18n, ToastService], stubs } })
    await vi.waitFor(() => expect(wrapper.text()).toContain('Loyalty Program'))
  })

  it('shows add account button', async () => {
    const wrapper = mount(Loyalty, { global: { plugins: [i18n, ToastService], stubs } })
    await vi.waitFor(() => expect(wrapper.text()).toContain('Add Account'))
  })

  it('loads accounts from API', async () => {
    mount(Loyalty, { global: { plugins: [i18n, ToastService], stubs } })
    await vi.waitFor(() => expect(mockGet).toHaveBeenCalled())
  })

  it('shows empty state', async () => {
    const wrapper = mount(Loyalty, { global: { plugins: [i18n, ToastService], stubs } })
    await vi.waitFor(() => expect(wrapper.text()).toContain('No results'))
  })
})

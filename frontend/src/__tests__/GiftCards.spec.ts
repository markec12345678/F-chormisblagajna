import { describe, it, expect, beforeEach, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import { createPinia, setActivePinia } from 'pinia'
import { createI18n } from 'vue-i18n'
import GiftCards from '@/pages/GiftCards.vue'
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
      gift_cards: 'Gift Cards',
      code: 'Code',
      balance: 'Balance',
      initial_amount: 'Initial Amount',
      issued_to: 'Issued To',
      issued: 'Issued',
      active: 'Active',
      actions: 'Actions',
      issue_gift_card: 'Issue Gift Card',
      expiry_date: 'Expiry Date',
      issue: 'Issue',
      no_data: 'No data',
    },
  },
})

const stubs = {
  DataTable: { template: '<div class="datatable-stub"><slot name="empty" /><slot /></div>', props: ['value', 'loading'] },
  Column: { template: '<div class="column-stub"><slot name="body" :data="{}" /></div>', props: ['field', 'header'] },
  Button: { template: '<button class="btn-stub">{{ label }}</button>', props: ['label', 'icon', 'severity', 'size', 'loading'] },
  Card: { template: '<div class="card-stub"><slot name="title" /><slot name="content" /></div>' },
  InputText: { template: '<input class="input-text-stub" />', props: ['modelValue'] },
  InputNumber: { template: '<input class="input-number-stub" />', props: ['modelValue'] },
  Calendar: { template: '<input class="calendar-stub" />', props: ['modelValue', 'dateFormat'] },
}

describe('GiftCards', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    mockGet.mockReset()
    mockPost.mockReset()
    mockGet.mockResolvedValue({ data: { data: [] } })
  })

  it('renders page heading', async () => {
    const wrapper = mount(GiftCards, { global: { plugins: [i18n, ToastService], stubs } })
    await vi.waitFor(() => expect(wrapper.text()).toContain('Gift Cards'))
  })

  it('shows issue gift card form', async () => {
    const wrapper = mount(GiftCards, { global: { plugins: [i18n, ToastService], stubs } })
    await vi.waitFor(() => expect(wrapper.text()).toContain('Issue Gift Card'))
  })

  it('shows issue button', async () => {
    const wrapper = mount(GiftCards, { global: { plugins: [i18n, ToastService], stubs } })
    await vi.waitFor(() => expect(wrapper.text()).toContain('Issue'))
  })

  it('loads cards from API', async () => {
    mount(GiftCards, { global: { plugins: [i18n, ToastService], stubs } })
    await vi.waitFor(() => expect(mockGet).toHaveBeenCalled())
  })

  it('shows empty state', async () => {
    const wrapper = mount(GiftCards, { global: { plugins: [i18n, ToastService], stubs } })
    await vi.waitFor(() => expect(wrapper.text()).toContain('No data'))
  })
})

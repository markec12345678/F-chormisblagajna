import { describe, it, expect, beforeEach, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import { createPinia, setActivePinia } from 'pinia'
import { createI18n } from 'vue-i18n'
import MarketingCampaigns from '@/pages/MarketingCampaigns.vue'
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
const mockDelete = vi.hoisted(() => vi.fn())

vi.mock('axios', () => ({
  default: {
    get: mockGet,
    post: mockPost,
    put: vi.fn(),
    delete: mockDelete,
    patch: vi.fn(),
    defaults: { headers: { common: {} } },
  },
}))

const i18n = createI18n({
  legacy: false,
  locale: 'en',
  fallbackLocale: 'en',
  messages: {
    en: {
      marketing_campaigns: 'Marketing Campaigns',
      no_results: 'No results',
      name: 'Name',
      type: 'Type',
      discount: 'Discount',
      start_date: 'Start Date',
      end_date: 'End Date',
      active: 'Active',
      actions: 'Actions',
      new_campaign: 'New Campaign',
      description: 'Description',
      create: 'Create',
    },
  },
})

const stubs = {
  DataTable: {
    template: '<div class="datatable-stub"><slot name="empty" /><slot /></div>',
    props: ['value', 'loading'],
  },
  Column: { template: '<div class="column-stub"><slot name="body" :data="{}" /></div>', props: ['field', 'header'] },
  Button: { template: '<button class="btn-stub" />', props: ['icon', 'severity', 'size', 'loading'] },
  InputText: { template: '<input class="input-text-stub" />', props: ['modelValue'] },
  InputNumber: { template: '<input class="input-number-stub" />', props: ['modelValue', 'min', 'max'] },
  Dropdown: { template: '<select class="dropdown-stub" />', props: ['modelValue', 'options'] },
  Calendar: { template: '<input class="calendar-stub" />', props: ['modelValue', 'dateFormat'] },
  Textarea: { template: '<textarea class="textarea-stub" />', props: ['modelValue'] },
  Card: { template: '<div class="card-stub"><slot name="title" /><slot name="content" /></div>' },
}

describe('MarketingCampaigns', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    mockGet.mockReset()
    mockPost.mockReset()
    mockDelete.mockReset()
    mockGet.mockResolvedValue({ data: { data: [] } })
  })

  it('renders page heading', async () => {
    const wrapper = mount(MarketingCampaigns, {
      global: { plugins: [i18n, ToastService], stubs },
    })
    await vi.waitFor(() => {
      expect(wrapper.text()).toContain('Marketing Campaigns')
    })
  })

  it('loads campaigns from API', async () => {
    mockGet.mockResolvedValue({ data: { data: [{ id: '1', name: 'Summer Sale' }] } })
    mount(MarketingCampaigns, {
      global: { plugins: [i18n, ToastService], stubs },
    })
    await vi.waitFor(() => {
      expect(mockGet).toHaveBeenCalled()
    })
  })

  it('shows new campaign form', async () => {
    const wrapper = mount(MarketingCampaigns, {
      global: { plugins: [i18n, ToastService], stubs },
    })
    await vi.waitFor(() => {
      expect(wrapper.text()).toContain('New Campaign')
    })
  })

  it('renders form fields for campaign creation', async () => {
    const wrapper = mount(MarketingCampaigns, {
      global: { plugins: [i18n, ToastService], stubs },
    })
    await vi.waitFor(() => {
      expect(wrapper.find('.input-text-stub').exists()).toBe(true)
      expect(wrapper.find('.dropdown-stub').exists()).toBe(true)
      expect(wrapper.find('.input-number-stub').exists()).toBe(true)
      expect(wrapper.find('.calendar-stub').exists()).toBe(true)
      expect(wrapper.find('.textarea-stub').exists()).toBe(true)
    })
  })

  it('shows empty state when no campaigns', async () => {
    const wrapper = mount(MarketingCampaigns, {
      global: { plugins: [i18n, ToastService], stubs },
    })
    await vi.waitFor(() => {
      expect(wrapper.text()).toContain('No results')
    })
  })
})

import { describe, it, expect, beforeEach, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import { createPinia, setActivePinia } from 'pinia'
import { createI18n } from 'vue-i18n'
import WasteTracking from '@/pages/WasteTracking.vue'
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
    defaults: { headers: { common: {} } },
  },
}))

const i18n = createI18n({
  legacy: false,
  locale: 'en',
  fallbackLocale: 'en',
  messages: {
    en: {
      waste_tracking: 'Waste Tracking',
      start_date: 'Start Date',
      end_date: 'End Date',
      add_waste_entry: 'Add Waste Entry',
      total_waste_cost: 'Total Waste Cost',
      total_entries: 'Total Entries',
      by_reason: 'By Reason',
      no_results: 'No results',
      save: 'Save',
      cancel: 'Cancel',
    },
  },
})

const stubs = {
  DataTable: { template: '<div class="datatable-stub"><slot name="empty" /><slot /></div>', props: ['value', 'loading'] },
  Column: { template: '<div class="column-stub"><slot name="body" :data="{}" /></div>', props: ['field', 'header'] },
  Button: { template: '<button class="btn-stub">{{ label }}</button>', props: ['label', 'icon', 'severity', 'size'] },
  Tag: { template: '<span class="tag-stub"><slot /></span>', props: ['severity', 'value'] },
  Calendar: { template: '<input class="calendar-stub" />', props: ['modelValue', 'placeholder', 'dateFormat'] },
  Card: { template: '<div class="card-stub"><slot name="title" /><slot name="content" /></div>' },
  Dialog: { template: '<div class="dialog-stub" v-if="visible"><slot /></div>', props: ['visible', 'header'] },
  InputText: { template: '<input class="input-text-stub" />', props: ['modelValue'] },
  InputNumber: { template: '<input class="input-number-stub" />', props: ['modelValue'] },
  Textarea: { template: '<textarea class="textarea-stub" />', props: ['modelValue'] },
  Select: { template: '<select class="select-stub" />', props: ['modelValue', 'options'] },
  ButtonGroup: { template: '<div class="btn-group-stub"><slot /></div>' },
}

describe('WasteTracking', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    mockGet.mockReset()
    mockPost.mockReset()
    mockDelete.mockReset()
    mockGet.mockResolvedValue({ data: { data: [] } })
  })

  it('renders page heading', async () => {
    const wrapper = mount(WasteTracking, { global: { plugins: [i18n, ToastService], stubs } })
    await vi.waitFor(() => expect(wrapper.text()).toContain('Waste Tracking'))
  })

  it('shows add waste entry button', async () => {
    const wrapper = mount(WasteTracking, { global: { plugins: [i18n, ToastService], stubs } })
    await vi.waitFor(() => expect(wrapper.text()).toContain('Add Waste Entry'))
  })

  it('loads data from API', async () => {
    mount(WasteTracking, { global: { plugins: [i18n, ToastService], stubs } })
    await vi.waitFor(() => expect(mockGet).toHaveBeenCalled())
  })
})

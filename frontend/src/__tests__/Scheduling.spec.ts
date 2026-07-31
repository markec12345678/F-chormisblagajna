import { describe, it, expect, beforeEach, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import { createPinia, setActivePinia } from 'pinia'
import { createI18n } from 'vue-i18n'
import Scheduling from '@/pages/Scheduling.vue'
import ToastService from 'primevue/toastservice'
import ConfirmationService from 'primevue/confirmationservice'

vi.mock('@/services/auth', () => ({
  default: {
    accessToken: { value: 'test-token' },
    currentUser: { value: null },
    signOut: vi.fn(),
  },
}))

const mockGet = vi.hoisted(() => vi.fn())

vi.mock('axios', () => ({
  default: {
    get: mockGet,
    post: vi.fn(),
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
      employee_scheduling: 'Employee Scheduling',
      add_shift: 'Add Shift',
      date_range: 'Date Range',
      no_results: 'No results',
      employee: 'Employee',
      date: 'Date',
      start_time: 'Start Time',
      end_time: 'End Time',
      role: 'Role',
      status: 'Status',
    },
  },
})

const stubs = {
  DataTable: { template: '<div class="datatable-stub"><slot name="header" /><slot name="empty" /><slot /></div>', props: ['value', 'loading'] },
  Column: { template: '<div class="column-stub"><slot name="body" :data="{}" /></div>', props: ['field', 'header'] },
  Button: { template: '<button class="btn-stub">{{ label }}</button>', props: ['label', 'icon', 'rounded', 'raised'] },
  Tag: { template: '<span class="tag-stub"><slot /></span>', props: ['severity', 'value'] },
  Calendar: { template: '<input class="calendar-stub" />', props: ['modelValue', 'selectionMode', 'placeholder'] },
  Dialog: { template: '<div class="dialog-stub" v-if="visible"><slot /></div>', props: ['visible', 'header'] },
  Dropdown: { template: '<select class="dropdown-stub" />', props: ['modelValue', 'options'] },
}

describe('Scheduling', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    mockGet.mockReset()
    mockGet.mockResolvedValue({ data: { data: [], meta: { total_records: 0 } } })
  })

  it('renders page heading', async () => {
    const wrapper = mount(Scheduling, { global: { plugins: [i18n, ToastService, ConfirmationService], stubs } })
    await vi.waitFor(() => expect(wrapper.text()).toContain('Employee Scheduling'))
  })

  it('shows add shift button', async () => {
    const wrapper = mount(Scheduling, { global: { plugins: [i18n, ToastService, ConfirmationService], stubs } })
    await vi.waitFor(() => expect(wrapper.text()).toContain('Add Shift'))
  })

  it('loads shifts from API', async () => {
    mount(Scheduling, { global: { plugins: [i18n, ToastService, ConfirmationService], stubs } })
    await vi.waitFor(() => expect(mockGet).toHaveBeenCalled())
  })

  it('shows empty state', async () => {
    const wrapper = mount(Scheduling, { global: { plugins: [i18n, ToastService, ConfirmationService], stubs } })
    await vi.waitFor(() => expect(wrapper.text()).toContain('No results'))
  })
})

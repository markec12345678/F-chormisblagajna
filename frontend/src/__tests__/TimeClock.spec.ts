import { describe, it, expect, beforeEach, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import { createPinia, setActivePinia } from 'pinia'
import { createI18n } from 'vue-i18n'
import TimeClock from '@/pages/TimeClock.vue'
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
      time_clock: 'Time Clock',
      clock_in_out: 'Clock In / Out',
      employee: 'Employee',
      select_employee: 'Select Employee',
      notes: 'Notes',
      clock_in: 'Clock In',
      clock_out: 'Clock Out',
      currently_clocked_in: 'Currently Clocked In',
      no_active_shifts: 'No active shifts',
    },
  },
})

const stubs = {
  DataTable: { template: '<div class="datatable-stub"><slot name="empty" /><slot /></div>', props: ['value', 'loading'] },
  Column: { template: '<div class="column-stub"><slot name="body" :data="{}" /></div>', props: ['field', 'header'] },
  Button: { template: '<button class="btn-stub">{{ label }}</button>', props: ['label', 'icon', 'severity', 'loading'] },
  Card: { template: '<div class="card-stub"><slot name="title" /><slot name="content" /></div>' },
  Dropdown: { template: '<select class="dropdown-stub" />', props: ['modelValue', 'options', 'placeholder'] },
  InputText: { template: '<input class="input-text-stub" />', props: ['modelValue'] },
}

describe('TimeClock', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    mockGet.mockReset()
    mockPost.mockReset()
    mockGet.mockResolvedValue({ data: { data: { currently_clocked_in: [], employees: [] } } })
  })

  it('renders page heading', async () => {
    const wrapper = mount(TimeClock, { global: { plugins: [i18n, ToastService], stubs } })
    await vi.waitFor(() => expect(wrapper.text()).toContain('Time Clock'))
  })

  it('shows clock in button when no active shift', async () => {
    const wrapper = mount(TimeClock, { global: { plugins: [i18n, ToastService], stubs } })
    await vi.waitFor(() => expect(wrapper.text()).toContain('Clock In'))
  })

  it('shows no active shifts empty state', async () => {
    const wrapper = mount(TimeClock, { global: { plugins: [i18n, ToastService], stubs } })
    await vi.waitFor(() => expect(wrapper.text()).toContain('No active shifts'))
  })

  it('loads data from API', async () => {
    mount(TimeClock, { global: { plugins: [i18n, ToastService], stubs } })
    await vi.waitFor(() => expect(mockGet).toHaveBeenCalled())
  })
})

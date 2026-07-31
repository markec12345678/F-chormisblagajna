import { describe, it, expect, beforeEach, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import { createPinia, setActivePinia } from 'pinia'
import { createI18n } from 'vue-i18n'
import QueueWaitlist from '@/pages/QueueWaitlist.vue'
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
const mockPut = vi.hoisted(() => vi.fn())
const mockDelete = vi.hoisted(() => vi.fn())

vi.mock('axios', () => ({
  default: {
    get: mockGet,
    post: mockPost,
    put: mockPut,
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
      queue_waitlist: 'Queue Waitlist',
      queue_empty: 'Queue is empty',
      customer: 'Customer',
      guests: 'Guests',
      estimated_wait: 'Est. Wait',
      status: 'Status',
      actions: 'Actions',
      add_to_queue: 'Add to Queue',
      customer_name: 'Name',
      phone: 'Phone',
      notes: 'Notes',
      add: 'Add',
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
  InputNumber: { template: '<input class="input-number-stub" />', props: ['modelValue', 'min'] },
  Tag: { template: '<span class="tag-stub"><slot /></span>', props: ['severity', 'value'] },
  Card: { template: '<div class="card-stub"><slot name="title" /><slot name="content" /></div>' },
}

describe('QueueWaitlist', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    mockGet.mockReset()
    mockPost.mockReset()
    mockPut.mockReset()
    mockDelete.mockReset()
    mockGet.mockResolvedValue({ data: { data: [] } })
  })

  it('renders page heading', async () => {
    const wrapper = mount(QueueWaitlist, {
      global: { plugins: [i18n, ToastService], stubs },
    })
    await vi.waitFor(() => {
      expect(wrapper.text()).toContain('Queue Waitlist')
    })
  })

  it('loads queue entries from API', async () => {
    mount(QueueWaitlist, {
      global: { plugins: [i18n, ToastService], stubs },
    })
    await vi.waitFor(() => {
      expect(mockGet).toHaveBeenCalled()
    })
  })

  it('shows add to queue form', async () => {
    const wrapper = mount(QueueWaitlist, {
      global: { plugins: [i18n, ToastService], stubs },
    })
    await vi.waitFor(() => {
      expect(wrapper.text()).toContain('Add to Queue')
    })
  })

  it('shows empty state when no entries', async () => {
    const wrapper = mount(QueueWaitlist, {
      global: { plugins: [i18n, ToastService], stubs },
    })
    await vi.waitFor(() => {
      expect(wrapper.text()).toContain('Queue is empty')
    })
  })

  it('renders form fields', async () => {
    const wrapper = mount(QueueWaitlist, {
      global: { plugins: [i18n, ToastService], stubs },
    })
    await vi.waitFor(() => {
      expect(wrapper.find('.input-text-stub').exists()).toBe(true)
      expect(wrapper.find('.input-number-stub').exists()).toBe(true)
    })
  })
})

import { describe, it, expect, beforeEach, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import { createPinia, setActivePinia } from 'pinia'
import { createI18n } from 'vue-i18n'
import KioskSettings from '@/pages/KioskSettings.vue'
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
      self_service_kiosk: 'Self Service Kiosk',
      name: 'Name',
      location: 'Location',
      theme: 'Theme',
      active: 'Active',
      actions: 'Actions',
      edit: 'Edit',
      add_kiosk: 'Add Kiosk',
      show_categories: 'Show Categories',
      save: 'Save',
    },
  },
})

const stubs = {
  DataTable: {
    template: '<div class="datatable-stub"><slot /></div>',
    props: ['value', 'loading'],
  },
  Column: { template: '<div class="column-stub"><slot name="body" :data="{}" /></div>', props: ['field', 'header'] },
  Button: { template: '<button class="btn-stub" @click="$emit(\'click\')" />', props: ['icon', 'size', 'loading'] },
  InputText: { template: '<input class="input-text-stub" />', props: ['modelValue'] },
  Dropdown: { template: '<select class="dropdown-stub" />', props: ['modelValue', 'options'] },
  Checkbox: { template: '<input type="checkbox" class="checkbox-stub" />', props: ['modelValue', 'binary'] },
  Card: { template: '<div class="card-stub"><slot name="title" /><slot name="content" /></div>' },
}

describe('KioskSettings', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    mockGet.mockReset()
    mockPost.mockReset()
    mockGet.mockResolvedValue({ data: { data: [] } })
  })

  it('renders page heading', async () => {
    const wrapper = mount(KioskSettings, {
      global: { plugins: [i18n, ToastService], stubs },
    })
    await vi.waitFor(() => {
      expect(wrapper.text()).toContain('Self Service Kiosk')
    })
  })

  it('loads kiosk configs from API', async () => {
    mockGet.mockResolvedValue({ data: { data: [{ id: '1', name: 'Kiosk 1' }] } })
    mount(KioskSettings, {
      global: { plugins: [i18n, ToastService], stubs },
    })
    await vi.waitFor(() => {
      expect(mockGet).toHaveBeenCalled()
    })
  })

  it('shows add kiosk form', async () => {
    const wrapper = mount(KioskSettings, {
      global: { plugins: [i18n, ToastService], stubs },
    })
    await vi.waitFor(() => {
      expect(wrapper.text()).toContain('Add Kiosk')
    })
  })

  it('renders form fields', async () => {
    const wrapper = mount(KioskSettings, {
      global: { plugins: [i18n, ToastService], stubs },
    })
    await vi.waitFor(() => {
      expect(wrapper.find('.input-text-stub').exists()).toBe(true)
      expect(wrapper.find('.dropdown-stub').exists()).toBe(true)
      expect(wrapper.find('.checkbox-stub').exists()).toBe(true)
    })
  })
})

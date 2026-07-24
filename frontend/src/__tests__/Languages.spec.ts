import { describe, it, expect, beforeEach, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import { createPinia, setActivePinia } from 'pinia'
import { createI18n } from 'vue-i18n'
import Languages from '@/pages/Languages.vue'
import ToastService from 'primevue/toastservice'

vi.mock('@/stores', () => ({
  globalStore: () => ({
    shop_mode: '',
    orientation: 'ltr',
    getColorMode: 'light',
    setOrientation: vi.fn(),
  }),
}))

vi.mock('@/services/auth', () => ({
  default: {
    accessToken: { value: 'test-token' },
    currentUser: { value: null },
  },
}))

const mockGet = vi.hoisted(() => vi.fn())

vi.mock('axios', () => ({
  default: {
    get: mockGet,
    post: vi.fn(),
    patch: vi.fn(),
    delete: vi.fn(),
  },
}))

const i18n = createI18n({
  legacy: false,
  locale: 'en',
  fallbackLocale: 'en',
  messages: {
    en: {
      language: 'Language',
      select_language: 'Select Language',
      apply: 'Apply',
      language_load_error: 'Failed to load language',
    },
  },
})

const stubs = {
  Dropdown: {
    template: '<select class="dropdown-stub" @change="$emit(\'change\')" />',
    props: ['modelValue', 'options', 'optionLabel', 'placeholder'],
  },
  Button: {
    template: '<button class="btn-stub" @click="$emit(\'click\')">{{ label }}</button>',
    props: ['label', 'type', 'severity'],
  },
}

describe('Languages', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    mockGet.mockReset()
  })

  it('renders page heading', async () => {
    mockGet.mockResolvedValue({ data: { data: [{ language: 'English', code: 'en' }] } })

    const wrapper = mount(Languages, {
      global: { plugins: [i18n, ToastService], stubs },
    })

    await vi.waitFor(() => {
      expect(wrapper.text()).toContain('Language')
    })
  })

  it('shows select language placeholder', async () => {
    mockGet.mockResolvedValue({ data: { data: [{ language: 'English', code: 'en' }] } })

    const wrapper = mount(Languages, {
      global: { plugins: [i18n, ToastService], stubs },
    })

    await vi.waitFor(() => {
      expect(wrapper.find('.dropdown-stub').exists()).toBe(true)
    })
  })

  it('hides apply button initially', async () => {
    mockGet.mockResolvedValue({ data: { data: [{ language: 'English', code: 'en' }] } })

    const wrapper = mount(Languages, {
      global: { plugins: [i18n, ToastService], stubs },
    })

    await vi.waitFor(() => {
      expect(wrapper.find('.btn-stub').exists()).toBe(false)
    })
  })

  it('loads available languages from API', async () => {
    mockGet.mockResolvedValue({
      data: {
        data: [
          { language: 'English', code: 'en' },
          { language: 'Slovenščina', code: 'sl' },
        ],
      },
    })

    const wrapper = mount(Languages, {
      global: { plugins: [i18n, ToastService], stubs },
    })

    await vi.waitFor(() => {
      expect(mockGet).toHaveBeenCalled()
    })
  })
})

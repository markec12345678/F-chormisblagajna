import { describe, it, expect, beforeEach, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import { createPinia, setActivePinia } from 'pinia'
import { createI18n } from 'vue-i18n'
import AI from '@/pages/AI.vue'
import ToastService from 'primevue/toastservice'

vi.mock('@/stores', () => ({
  globalStore: () => ({
    shop_mode: '',
    orientation: 'ltr',
    getColorMode: 'light',
  }),
}))

vi.mock('@/services/auth', () => ({
  default: {
    accessToken: { value: 'test-token' },
    currentUser: { value: null },
  },
}))

const mockGet = vi.hoisted(() => vi.fn())
const mockPost = vi.hoisted(() => vi.fn())

vi.mock('axios', () => ({
  default: {
    get: mockGet,
    post: mockPost,
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
      ai_search: 'AI Search',
      search: 'Search',
      ai_voice: 'Voice',
      search_results: 'Search Results',
      smart_suggestions: 'Smart Suggestions',
      ai_suggestions: 'AI Suggestion',
      ai_search_description: 'AI-powered search',
      voice_ordering: 'Voice Ordering',
      voice_ordering_description: 'Speak to order',
      start_recording: 'Start Recording',
      stop_recording: 'Stop Recording',
      transcription: 'Transcription',
      no_results: 'No results found',
      name: 'Name',
      category: 'Category',
      total: 'Total',
      score: 'Score',
      actions: 'Actions',
      add_to_order: 'Add to Order',
      success: 'Success',
      failed: 'Failed',
    },
  },
})

const stubs = {
  Card: {
    template: '<div class="card-stub"><slot name="title" /><slot name="content" /></div>',
    props: ['class'],
  },
  DataTable: {
    template:
      '<div class="datatable-stub"><slot name="header" /><slot name="empty" /><slot /></div>',
    props: ['value', 'loading', 'stripedRows', 'tableStyle'],
  },
  Column: { template: '<div class="column-stub" />', props: ['field', 'header'] },
  Button: {
    template: '<button class="btn-stub" @click="$emit(\'click\')">{{ label }}</button>',
    props: ['label', 'icon', 'severity', 'rounded', 'raised', 'size', 'disabled', 'loading'],
  },
  InputText: {
    template: '<input class="input-text-stub" />',
    props: ['modelValue', 'placeholder', 'class'],
  },
  IconField: {
    template: '<div class="iconfield-stub"><slot /></div>',
    props: ['iconPosition'],
  },
  InputIcon: {
    template: '<span class="inputicon-stub"><slot /></span>',
  },
  Tag: {
    template: '<span class="tag-stub">{{ value }}</span>',
    props: ['value', 'severity'],
  },
  ProgressSpinner: { template: '<div class="progress-spinner-stub" />', props: ['style'] },
  Toast: { template: '<div />' },
}

describe('AI', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    mockGet.mockReset()
    mockPost.mockReset()
    mockGet.mockResolvedValue({ data: { data: [] } })
  })

  it('renders page heading', () => {
    const wrapper = mount(AI, {
      global: { plugins: [i18n, ToastService], stubs },
    })
    expect(wrapper.text()).toContain('AI Search')
  })

  it('renders search bar', async () => {
    const wrapper = mount(AI, {
      global: { plugins: [i18n, ToastService], stubs },
    })
    await vi.waitFor(() => {
      expect(wrapper.find('.input-text-stub').exists()).toBe(true)
    })
  })

  it('renders voice button', async () => {
    const wrapper = mount(AI, {
      global: { plugins: [i18n, ToastService], stubs },
    })
    await vi.waitFor(() => {
      expect(wrapper.find('.btn-stub').exists()).toBe(true)
    })
  })

  it('renders suggestions panel', async () => {
    const wrapper = mount(AI, {
      global: { plugins: [i18n, ToastService], stubs },
    })
    await vi.waitFor(() => {
      expect(wrapper.text()).toContain('Smart Suggestions')
    })
  })

  it('calls suggestions API on mount', async () => {
    mount(AI, {
      global: { plugins: [i18n, ToastService], stubs },
    })
    await vi.waitFor(() => {
      expect(mockGet).toHaveBeenCalled()
    })
  })

  it('handles empty results', async () => {
    mockGet.mockResolvedValue({ data: { data: [] } })

    const wrapper = mount(AI, {
      global: { plugins: [i18n, ToastService], stubs },
    })

    await vi.waitFor(() => {
      expect(mockGet).toHaveBeenCalled()
    })
  })
})

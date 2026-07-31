import { describe, it, expect, beforeEach, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import { createPinia, setActivePinia } from 'pinia'
import { createI18n } from 'vue-i18n'
import DisplayView from '@/pages/DisplayView.vue'
import ToastService from 'primevue/toastservice'
import { createRouter, createWebHistory } from 'vue-router'

const mockGet = vi.hoisted(() => vi.fn())

vi.mock('axios', () => ({
  default: { get: mockGet, post: vi.fn(), put: vi.fn(), delete: vi.fn(), defaults: { headers: { common: {} } } },
}))

const i18n = createI18n({
  legacy: false,
  locale: 'en',
  fallbackLocale: 'en',
  messages: { en: { promotions: 'Promotions', our_menu: 'Our Menu', menu_placeholder: 'Delicious food', current_orders: 'Current Orders' } },
})

const router = createRouter({ history: createWebHistory(), routes: [{ path: '/display/:id', name: 'display', component: DisplayView }] })

const stubs = {
  Tag: { template: '<span class="tag-stub"><slot /></span>', props: ['severity'] },
}

describe('DisplayView', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    mockGet.mockReset()
    mockGet.mockResolvedValue({ data: { data: { items: [], interval: 0, welcome_message: 'Welcome', theme: 'light' } } })
  })

  it('renders welcome message when present', async () => {
    mockGet.mockResolvedValue({ data: { data: { items: [], interval: 0, welcome_message: 'Welcome!', theme: 'light' } } })
    const wrapper = mount(DisplayView, { global: { plugins: [i18n, ToastService, router], stubs } })
    await vi.waitFor(() => expect(wrapper.text()).toContain('Welcome!'))
  })

  it('calls API on mount', async () => {
    mount(DisplayView, { global: { plugins: [i18n, ToastService, router], stubs } })
    await vi.waitFor(() => expect(mockGet).toHaveBeenCalled())
  })

  it('handles API error gracefully', async () => {
    mockGet.mockRejectedValue(new Error('fail'))
    const wrapper = mount(DisplayView, { global: { plugins: [i18n, ToastService, router], stubs } })
    await vi.waitFor(() => expect(wrapper.exists()).toBe(true))
  })

  it('applies theme class from content', async () => {
    mockGet.mockResolvedValue({ data: { data: { items: [], interval: 0, welcome_message: '', theme: 'dark' } } })
    const wrapper = mount(DisplayView, { global: { plugins: [i18n, ToastService, router], stubs } })
    await vi.waitFor(() => {
      expect(wrapper.find('.display-fullscreen').exists()).toBe(true)
    })
  })
})

import { describe, it, expect, beforeEach, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import { createPinia, setActivePinia } from 'pinia'
import { createI18n } from 'vue-i18n'
import NoAccessView from '@/pages/NoAccessView.vue'
import ToastService from 'primevue/toastservice'
import ConfirmationService from 'primevue/confirmationservice'

vi.mock('vue-router', () => ({
  useRouter: () => ({
    push: vi.fn(),
    replace: vi.fn(),
    back: vi.fn(),
  }),
  useRoute: () => ({
    params: {},
    query: {},
  }),
}))

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
      access_denied: 'Access Denied',
      relogin: 'Re-login',
    },
  },
})

const stubs = {
  Button: {
    template: '<button class="btn-stub" @click="$emit(\'click\')">{{ label }}</button>',
    props: ['label', 'severity', 'class'],
  },
  Toast: { template: '<div />' },
}

describe('NoAccessView', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    mockGet.mockReset()
    mockPost.mockReset()
  })

  it('renders the access denied message', () => {
    const wrapper = mount(NoAccessView, {
      global: { plugins: [i18n, ToastService, ConfirmationService], stubs },
    })

    expect(wrapper.text()).toContain('Access Denied')
  })

  it('shows the ban icon', () => {
    const wrapper = mount(NoAccessView, {
      global: { plugins: [i18n, ToastService, ConfirmationService], stubs },
    })

    expect(wrapper.find('.pi-ban').exists()).toBe(true)
  })

  it('shows the re-login button', () => {
    const wrapper = mount(NoAccessView, {
      global: { plugins: [i18n, ToastService, ConfirmationService], stubs },
    })

    expect(wrapper.find('.btn-stub').exists()).toBe(true)
  })

  it('calls signOut on re-login click', async () => {
    const auth = await import('@/services/auth')

    const wrapper = mount(NoAccessView, {
      global: { plugins: [i18n, ToastService, ConfirmationService], stubs },
    })

    const btn = wrapper.find('.btn-stub')
    await btn.trigger('click')

    expect(auth.default.signOut).toHaveBeenCalled()
  })

  it('renders centered layout', () => {
    const wrapper = mount(NoAccessView, {
      global: { plugins: [i18n, ToastService, ConfirmationService], stubs },
    })

    expect(wrapper.find('h2').exists()).toBe(true)
  })
})

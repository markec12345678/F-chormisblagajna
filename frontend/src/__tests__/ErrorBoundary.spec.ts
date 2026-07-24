import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import ErrorBoundary from '@/components/ErrorBoundary.vue'
import PrimeVue from 'primevue/config'
import { createI18n } from 'vue-i18n'

const i18n = createI18n({
  legacy: false,
  locale: 'en',
  fallbackLocale: 'en',
  messages: {
    en: {
      error_occurred: 'An error occurred',
      retry: 'Retry',
    },
  },
})

describe('ErrorBoundary', () => {
  it('renders slot content when no error', () => {
    const wrapper = mount(ErrorBoundary, {
      global: { plugins: [PrimeVue, i18n] },
      slots: { default: '<div class="child">Hello</div>' },
    })
    expect(wrapper.find('.child').exists()).toBe(true)
    expect(wrapper.find('.error-boundary').exists()).toBe(false)
  })

  it('does not render error-boundary div by default', () => {
    const wrapper = mount(ErrorBoundary, {
      global: { plugins: [PrimeVue, i18n] },
      slots: { default: '<div>Content</div>' },
    })
    expect(wrapper.find('.error-boundary').exists()).toBe(false)
    expect(wrapper.find('.error-content').exists()).toBe(false)
  })

  it('exposes error and reset methods', () => {
    const wrapper = mount(ErrorBoundary, {
      global: { plugins: [PrimeVue, i18n] },
      slots: { default: '<div>Content</div>' },
    })
    expect(typeof wrapper.vm.error).toBe('object')
    expect(typeof wrapper.vm.reset).toBe('function')
  })

  it('has null error initially', () => {
    const wrapper = mount(ErrorBoundary, {
      global: { plugins: [PrimeVue, i18n] },
      slots: { default: '<div>Content</div>' },
    })
    expect(wrapper.vm.error).toBeNull()
  })

  it('reset sets error back to null', () => {
    const wrapper = mount(ErrorBoundary, {
      global: { plugins: [PrimeVue, i18n] },
      slots: { default: '<div>Content</div>' },
    })
    wrapper.vm.reset()
    expect(wrapper.vm.error).toBeNull()
  })

  it('renders multiple slot children', () => {
    const wrapper = mount(ErrorBoundary, {
      global: { plugins: [PrimeVue, i18n] },
      slots: {
        default: '<div class="a">A</div><div class="b">B</div>',
      },
    })
    expect(wrapper.find('.a').exists()).toBe(true)
    expect(wrapper.find('.b').exists()).toBe(true)
    expect(wrapper.find('.error-boundary').exists()).toBe(false)
  })
})

import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'
import { mount } from '@vue/test-utils'
import NotificationView from '@/components/NotificationView.vue'
import PrimeVue from 'primevue/config'
import { createI18n } from 'vue-i18n'
import { Notification } from '@/classes/Notification'

vi.mock('date-fns', () => ({
  formatDistanceToNow: vi.fn(() => '5 minutes ago'),
}))

const i18n = createI18n({
  legacy: false,
  locale: 'en',
  fallbackLocale: 'en',
  messages: { en: {} },
})

const stubs = {
  Message: {
    template: '<div class="message-stub" :data-severity="severity"><slot/></div>',
    props: ['severity'],
    emits: ['close'],
  },
}

describe('NotificationView', () => {
  beforeEach(() => {
    vi.useFakeTimers()
  })

  afterEach(() => {
    vi.useRealTimers()
  })

  const makeNotification = (overrides = {}) => {
    const n = new Notification()
    n.topic_name = 'order_finished'
    n.description = 'Order #5 is ready'
    n.severity = 'success'
    n.date = new Date('2026-07-25T10:00:00Z')
    Object.assign(n, overrides)
    return n
  }

  it('renders notification topic_name', () => {
    const wrapper = mount(NotificationView, {
      props: { notification: makeNotification() },
      global: { plugins: [PrimeVue, i18n], stubs },
    })
    expect(wrapper.text()).toContain('order_finished')
  })

  it('renders notification description', () => {
    const wrapper = mount(NotificationView, {
      props: { notification: makeNotification() },
      global: { plugins: [PrimeVue, i18n], stubs },
    })
    expect(wrapper.text()).toContain('Order #5 is ready')
  })

  it('shows elapsed time', () => {
    const wrapper = mount(NotificationView, {
      props: { notification: makeNotification() },
      global: { plugins: [PrimeVue, i18n], stubs },
    })
    expect(wrapper.text()).toContain('5 minutes ago')
  })

  it('passes severity to Message', () => {
    const wrapper = mount(NotificationView, {
      props: { notification: makeNotification({ severity: 'warn' }) },
      global: { plugins: [PrimeVue, i18n], stubs },
    })
    const msg = wrapper.find('.message-stub')
    expect(msg.attributes('data-severity')).toBe('warn')
  })

  it('has closed emit defined', async () => {
    const wrapper = mount(NotificationView, {
      props: { notification: makeNotification() },
      global: { plugins: [PrimeVue, i18n], stubs },
    })
    expect(wrapper.emitted('closed')).toBeUndefined()
  })
})

import { describe, it, expect, beforeEach, vi } from 'vitest'
import { createPinia, setActivePinia } from 'pinia'
import { mount } from '@vue/test-utils'
import MainSearchResultView from '@/components/MainSearchResultView.vue'
import PrimeVue from 'primevue/config'
import { createI18n } from 'vue-i18n'
import Order from '@/classes/Order'
import { OrderItem } from '@/classes/OrderItem'

vi.mock('@/stores', () => ({
  globalStore: () => ({
    shop_mode: '',
    orientation: 'ltr',
    getColorMode: 'light',
  }),
}))

const i18n = createI18n({
  legacy: false,
  locale: 'en',
  fallbackLocale: 'en',
  messages: {
    en: {
      name: 'Name',
      items: 'items',
      info: 'Info',
      paid: 'Paid',
      unpaid: 'Unpaid',
      pending: 'Pending',
      cancelled: 'Cancelled',
      in_progress: 'In Progress',
      finished: 'Finished',
      stashed: 'Stashed',
    },
  },
})

const stubs = {
  Message: { template: '<div class="message-stub"><slot name="container"/></div>' },
  Badge: { template: '<span class="badge-stub">{{ value }}</span>', props: ['value', 'severity'] },
  Button: {
    template: '<button @click="$emit(\'click\')"><slot/></button>',
    props: ['icon', 'severity', 'aria-label'],
  },
}

const makeOrder = (state: string, isPaid: boolean, displayId = 'ORD-1') => {
  const order = new Order()
  order.state = state
  order.is_paid = isPaid
  order.display_id = displayId
  const item = new OrderItem()
  order.items = [item]
  return order
}

describe('MainSearchResultView', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
  })

  it('renders order display_id', () => {
    const wrapper = mount(MainSearchResultView, {
      props: { order: makeOrder('pending', false, 'ORD-42') },
      global: { plugins: [PrimeVue, i18n], stubs },
    })
    expect(wrapper.text()).toContain('ORD-42')
  })

  it('shows item count', () => {
    const wrapper = mount(MainSearchResultView, {
      props: { order: makeOrder('pending', false) },
      global: { plugins: [PrimeVue, i18n], stubs },
    })
    expect(wrapper.text()).toContain('1')
    expect(wrapper.text()).toContain('items')
  })

  it('shows paid badge when order is paid', () => {
    const wrapper = mount(MainSearchResultView, {
      props: { order: makeOrder('finished', true) },
      global: { plugins: [PrimeVue, i18n], stubs },
    })
    expect(wrapper.text()).toContain('Paid')
  })

  it('shows unpaid badge when order is not paid', () => {
    const wrapper = mount(MainSearchResultView, {
      props: { order: makeOrder('pending', false) },
      global: { plugins: [PrimeVue, i18n], stubs },
    })
    expect(wrapper.text()).toContain('Unpaid')
  })

  it('emits view-order-pressed when button clicked', async () => {
    const wrapper = mount(MainSearchResultView, {
      props: { order: makeOrder('pending', false) },
      global: { plugins: [PrimeVue, i18n], stubs },
    })
    await wrapper.find('button').trigger('click')
    expect(wrapper.emitted('view-order-pressed')).toBeTruthy()
  })

  it.each([
    ['', 'Pending'],
    ['pending', 'Pending'],
    ['cancelled', 'Cancelled'],
    ['in_progress', 'In Progress'],
    ['finished', 'Finished'],
    ['stashed', 'Stashed'],
  ])('shows correct status for state=%s', (state, expected) => {
    const wrapper = mount(MainSearchResultView, {
      props: { order: makeOrder(state, false) },
      global: { plugins: [PrimeVue, i18n], stubs },
    })
    expect(wrapper.text()).toContain(expected)
  })
})

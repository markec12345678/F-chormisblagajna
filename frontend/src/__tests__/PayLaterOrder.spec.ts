import { describe, it, expect, vi, beforeEach } from 'vitest'
import { createPinia, setActivePinia } from 'pinia'
import { mount } from '@vue/test-utils'
import PayLaterOrder from '@/components/PayLaterOrder.vue'
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

vi.mock('axios', () => ({
  default: { get: vi.fn(() => Promise.resolve({})) },
}))

vi.mock('../services/auth', () => ({
  default: { accessToken: { value: 'token' } },
}))

const mockAdd = vi.fn()
vi.mock('primevue/usetoast', () => ({
  useToast: () => ({ add: mockAdd }),
}))

const i18n = createI18n({
  legacy: false,
  locale: 'en',
  fallbackLocale: 'en',
  messages: {
    en: {
      items: 'items',
      checkout: 'Checkout',
      success: 'Success',
      pending: 'Pending',
      in_progress: 'In Progress',
      cancelled: 'Cancelled',
      finished: 'Finished',
    },
  },
})

const stubs = {
  Message: { template: '<div class="message-stub"><slot name="container"/></div>' },
  Badge: { template: '<span class="badge-stub">{{ value }}</span>', props: ['value'] },
  Button: {
    template: '<button @click="$emit(\'click\')"><slot/></button>',
    props: ['aria-label'],
  },
}

const makeOrder = (state: string, price = 100) => {
  const order = new Order()
  order.state = state
  order.display_id = 'ORD-1'
  order.sale_price = price
  const item = new OrderItem()
  order.items = [item]
  return order
}

describe('PayLaterOrder', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    vi.clearAllMocks()
  })

  it('renders order display_id', () => {
    const wrapper = mount(PayLaterOrder, {
      props: { order: makeOrder('pending') },
      global: { plugins: [PrimeVue, i18n], stubs },
    })
    expect(wrapper.text()).toContain('ORD-1')
  })

  it('shows unpaid amount', () => {
    const wrapper = mount(PayLaterOrder, {
      props: { order: makeOrder('pending', 250) },
      global: { plugins: [PrimeVue, i18n], stubs },
    })
    expect(wrapper.text()).toContain('250')
  })

  it('shows item count', () => {
    const wrapper = mount(PayLaterOrder, {
      props: { order: makeOrder('pending') },
      global: { plugins: [PrimeVue, i18n], stubs },
    })
    expect(wrapper.text()).toContain('1')
    expect(wrapper.text()).toContain('items')
  })

  it.each([
    ['', 'Pending'],
    ['pending', 'Pending'],
    ['cancelled', 'Cancelled'],
    ['in_progress', 'In Progress'],
    ['finished', 'Finished'],
  ])('shows correct status for state=%s', (state, expected) => {
    const wrapper = mount(PayLaterOrder, {
      props: { order: makeOrder(state) },
      global: { plugins: [PrimeVue, i18n], stubs },
    })
    expect(wrapper.text()).toContain(expected)
  })

  it('calls axios.get on checkout button click', async () => {
    const axios = await import('axios')
    const wrapper = mount(PayLaterOrder, {
      props: { order: makeOrder('pending') },
      global: { plugins: [PrimeVue, i18n], stubs },
    })
    await wrapper.find('button').trigger('click')
    expect(axios.default.get).toHaveBeenCalled()
  })

  it('emits order_paid on successful payment', async () => {
    const wrapper = mount(PayLaterOrder, {
      props: { order: makeOrder('pending') },
      global: { plugins: [PrimeVue, i18n], stubs },
    })
    await wrapper.find('button').trigger('click')
    await vi.waitFor(() => {
      expect(wrapper.emitted('order_paid')).toBeTruthy()
    })
  })
})

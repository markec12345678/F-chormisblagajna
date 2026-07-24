import { describe, it, expect, beforeEach, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import { createPinia, setActivePinia } from 'pinia'
import { createI18n } from 'vue-i18n'
import StashedOrder from '@/components/StashedOrder.vue'

vi.mock('@/stores', () => ({
  globalStore: () => ({
    shop_mode: '',
    orientation: 'ltr',
    getColorMode: 'light',
    getSettings: { orders: { default_cost_calculation_method: 'exact' } },
  }),
}))

vi.mock('@/services/auth', () => ({
  default: {
    accessToken: { value: '' },
    currentUser: { value: null },
  },
}))

vi.mock('axios', () => ({
  default: { get: vi.fn(), post: vi.fn(), patch: vi.fn(), delete: vi.fn() },
}))

const i18n = createI18n({
  legacy: false,
  locale: 'en',
  fallbackLocale: 'en',
  messages: {
    en: {
      pending: 'Pending',
      stashed: 'Stashed',
      item: 'item',
    },
  },
})

function createOrder(overrides = {}) {
  return {
    id: 'order-1',
    display_id: '#001',
    state: 'stashed',
    items: [{ id: '1' }, { id: '2' }, { id: '3' }],
    submitted_at: new Date().toISOString(),
    started_at: null,
    comment: '',
    sale_price: 25,
    cost: 10,
    discount: 0,
    tips: 0,
    is_paid: false,
    is_pay_later: false,
    payment_source: '',
    customer: { id: '', name: '', phone: '', address: '' },
    delivery_info: null,
    custom_data: null,
    ...overrides,
  }
}

describe('StashedOrder', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
  })

  it('renders display_id and item count', () => {
    const order = createOrder({ items: [{ id: '1' }, { id: '2' }] })
    const wrapper = mount(StashedOrder, {
      props: { order },
      global: { plugins: [i18n] },
    })

    expect(wrapper.text()).toContain('#001')
    expect(wrapper.text()).toContain('2')
  })

  it('emits back_to_checkout when button clicked', async () => {
    const order = createOrder()
    const wrapper = mount(StashedOrder, {
      props: { order },
      global: { plugins: [i18n] },
    })

    const buttons = wrapper.findAllComponents({ name: 'Button' })
    const restoreBtn = buttons.find((btn) => {
      const icon = btn.props('icon')
      return icon === 'pi pi-arrow-down-right'
    })
    expect(restoreBtn).toBeDefined()

    await restoreBtn!.vm.$emit('click')
    await wrapper.vm.$nextTick()
    expect(wrapper.emitted('back_to_checkout')).toBeTruthy()
    expect(wrapper.emitted('back_to_checkout')!.length).toBe(1)
  })

  it('shows stashed status tag for stashed state', () => {
    const order = createOrder({ state: 'stashed' })
    const wrapper = mount(StashedOrder, {
      props: { order },
      global: { plugins: [i18n] },
    })

    expect(wrapper.text()).toContain('Stashed')
  })

  it('shows pending status tag for pending state', () => {
    const order = createOrder({ state: 'pending' })
    const wrapper = mount(StashedOrder, {
      props: { order },
      global: { plugins: [i18n] },
    })

    expect(wrapper.text()).toContain('Pending')
  })
})

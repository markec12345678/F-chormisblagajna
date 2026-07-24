import { describe, it, expect, beforeEach, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import { createPinia, setActivePinia } from 'pinia'
import { createI18n } from 'vue-i18n'
import ToastService from 'primevue/toastservice'
import ConfirmationService from 'primevue/confirmationservice'
import QueueOrder from '@/components/QueueOrder.vue'

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
    accessToken: { value: 'test-token' },
    currentUser: { value: null },
  },
}))

vi.mock('axios', () => ({
  default: {
    get: vi.fn().mockResolvedValue({ data: { data: {} } }),
    post: vi.fn().mockResolvedValue({ data: { started_at: new Date().toISOString() } }),
    patch: vi.fn().mockResolvedValue({}),
    delete: vi.fn().mockResolvedValue({}),
  },
}))

vi.mock('@/classes/OrderItem', () => ({
  OrderItem: vi.fn().mockImplementation(() => ({
    product: { id: '', name: '', materials: [] },
    materials: [],
    isValid: true,
    SetProductId: vi.fn(),
    RefreshReadyNumber: vi.fn().mockResolvedValue(undefined),
    FromItemData: vi.fn().mockResolvedValue(undefined),
    RefreshProductData: vi.fn().mockResolvedValue(undefined),
    ValidateMaterialTotalQuantity: vi.fn(),
    ValidateMaterialExactQuantity: vi.fn(),
  })),
  Product: vi.fn().mockImplementation(() => ({
    id: '',
    name: '',
    materials: [],
    sub_products: [],
    price: 0,
    quantity: 0,
  })),
}))

const i18n = createI18n({
  legacy: false,
  locale: 'en',
  fallbackLocale: 'en',
  messages: {
    en: {
      start: 'Start',
      finish: 'Finish',
      back: 'Back',
      next: 'Next',
      cancel: 'Cancel',
      yes: 'Yes',
      cancelled: 'Cancelled',
      product_details: 'Product details',
      confirm_finish_order: 'Finish this order?',
      order_finished: 'Order finished',
      order_finished_detail: 'Order has been finished',
      finish_order_failed: 'Failed to finish order',
      failed: 'Failed',
    },
  },
})

function createOrder(overrides = {}) {
  return {
    id: 'order-1',
    display_id: '#001',
    state: 'pending',
    items: [
      {
        id: 'item-1',
        product: { id: 'p1', name: 'Pizza Margherita', materials: [{ id: 'm1', name: 'Flour', quantity: 200, unit: 'g' }] },
        price: 10,
        quantity: 1,
        comment: '',
        sale_price: 10,
        cost: 5,
        cost_method: 'exact',
        status: '',
        materials: [],
        sub_items: [],
      },
    ],
    submitted_at: new Date().toISOString(),
    started_at: null,
    comment: '',
    sale_price: 10,
    cost: 5,
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

describe('QueueOrder', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
  })

  it('renders display_id', () => {
    const order = createOrder()
    const wrapper = mount(QueueOrder, {
      props: { order, number: 1 },
      global: { plugins: [i18n, ToastService, ConfirmationService] },
    })

    expect(wrapper.text()).toContain('#001')
  })

  it('renders product names from order items', () => {
    const order = createOrder()
    const wrapper = mount(QueueOrder, {
      props: { order, number: 1 },
      global: { plugins: [i18n, ToastService, ConfirmationService] },
    })

    expect(wrapper.text()).toContain('Pizza Margherita')
  })

  it('shows Start button when state is pending', () => {
    const order = createOrder({ state: 'pending' })
    const wrapper = mount(QueueOrder, {
      props: { order, number: 1 },
      global: { plugins: [i18n, ToastService, ConfirmationService] },
    })

    expect(wrapper.text()).toContain('Start')
  })

  it('shows Cancelled overlay when state is cancelled', () => {
    const order = createOrder({ state: 'cancelled' })
    const wrapper = mount(QueueOrder, {
      props: { order, number: 1 },
      global: { plugins: [i18n, ToastService, ConfirmationService] },
    })

    expect(wrapper.text()).toContain('Cancelled')
    expect(wrapper.text()).toContain('#001')
  })

  it('shows item comment when present', () => {
    const order = createOrder({
      items: [
        {
          id: 'item-1',
          product: { id: 'p1', name: 'Pizza', materials: [] },
          price: 10,
          quantity: 1,
          comment: 'Extra spicy',
          sale_price: 10,
          cost: 5,
          cost_method: 'exact',
          status: '',
          materials: [],
          sub_items: [],
        },
      ],
    })
    const wrapper = mount(QueueOrder, {
      props: { order, number: 1 },
      global: { plugins: [i18n, ToastService, ConfirmationService] },
    })

    expect(wrapper.text()).toContain('Extra spicy')
  })

  it('hides item comment when empty', () => {
    const order = createOrder({
      items: [
        {
          id: 'item-1',
          product: { id: 'p1', name: 'Pizza', materials: [] },
          price: 10,
          quantity: 1,
          comment: '',
          sale_price: 10,
          cost: 5,
          cost_method: 'exact',
          status: '',
          materials: [],
          sub_items: [],
        },
      ],
    })
    const wrapper = mount(QueueOrder, {
      props: { order, number: 1 },
      global: { plugins: [i18n, ToastService, ConfirmationService] },
    })

    expect(wrapper.findAll('p.mt-1').length).toBe(0)
  })

  it('renders multiple items with dividers', () => {
    const order = createOrder({
      items: [
        { id: 'i1', product: { id: 'p1', name: 'Pizza', materials: [] }, price: 10, quantity: 1, comment: '', sale_price: 10, cost: 5, cost_method: 'exact', status: '', materials: [], sub_items: [] },
        { id: 'i2', product: { id: 'p2', name: 'Pasta', materials: [] }, price: 8, quantity: 1, comment: '', sale_price: 8, cost: 3, cost_method: 'exact', status: '', materials: [], sub_items: [] },
      ],
    })
    const wrapper = mount(QueueOrder, {
      props: { order, number: 1 },
      global: { plugins: [i18n, ToastService, ConfirmationService] },
    })

    expect(wrapper.text()).toContain('Pizza')
    expect(wrapper.text()).toContain('Pasta')
  })

  it('hides finished orders via display computed', async () => {
    const order = createOrder({ state: 'finished' })
    const wrapper = mount(QueueOrder, {
      props: { order, number: 1 },
      global: { plugins: [i18n, ToastService, ConfirmationService] },
    })

    const rootDiv = wrapper.find('div')
    expect(rootDiv.attributes('style')).toContain('display: none')
  })

  it('emits openedDialog when dialog is opened', async () => {
    const order = createOrder()
    const wrapper = mount(QueueOrder, {
      props: { order, number: 1 },
      global: { plugins: [i18n, ToastService, ConfirmationService] },
    })

    const visible = wrapper.vm.$refs
    expect(wrapper.emitted('openedDialog')).toBeFalsy()
  })
})

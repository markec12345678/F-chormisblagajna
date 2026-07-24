import { describe, it, expect, beforeEach, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import { createPinia, setActivePinia } from 'pinia'
import { createI18n } from 'vue-i18n'
import ToastService from 'primevue/toastservice'
import ConfirmationService from 'primevue/confirmationservice'
import OrderView from '@/components/OrderView.vue'

vi.mock('@/components/OrderItemsInfo.vue', () => ({
  default: { name: 'OrderItemsInfo', template: '<div />' },
}))

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
    get: vi.fn().mockResolvedValue({ data: { data: [] } }),
    post: vi.fn().mockResolvedValue({}),
    patch: vi.fn().mockResolvedValue({}),
    delete: vi.fn().mockResolvedValue({}),
  },
}))

const i18n = createI18n({
  legacy: false,
  locale: 'en',
  fallbackLocale: 'en',
  messages: {
    en: {
      items: 'Items',
      item: 'item',
      display_id: 'Display ID',
      customer: 'Customer',
      status: 'Status',
      submitted_at: 'Submitted at',
      started_at: 'Started at',
      comment: 'Comment',
      cost: 'Cost',
      sale_price: 'Sale price',
      pay_later: 'Pay later',
      paid: 'Paid',
      unpaid: 'Unpaid',
      payment_source: 'Payment source',
      discount: 'Discount',
      tips: 'Tips',
      add_tip: 'Add tip',
      remove_tip: 'Remove tip',
      delivery_info: 'Delivery info',
      custom_data: 'Custom data',
      actions: 'Actions',
      client_receipt: 'Client receipt',
      kitchen_receipt: 'Kitchen receipt',
      collect_money: 'Collect money',
      finish: 'Finish',
      cancel: 'Cancel',
      order: 'Order',
      logs: 'Logs',
      pending: 'Pending',
      in_progress: 'In progress',
      finished: 'Finished',
      cancelled: 'Cancelled',
      tip_added: 'Tip added',
      tip_removed: 'Tip removed',
      order_finished: 'Order finished',
      money_collected: 'Money collected',
      print_signal_sent: 'Print signal sent',
      request_failed: 'Request failed',
      key: 'Key',
      value: 'Value',
      save: 'Save',
      add: 'Add',
      edit_custom_data: 'Edit custom data',
      add_custom_data: 'Add custom data',
      confirm_delete_custom_data: 'Are you sure?',
      confirm_cancel_order: 'Cancel this order?',
      order_cancelled_success: 'Order cancelled',
      order_cancel_failed: 'Failed to cancel order',
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
        id: '1',
        product: { name: 'Pizza', id: 'p1', materials: [], sub_products: [] },
        price: 10,
        quantity: 1,
        comment: '',
        sale_price: 10,
        cost: 5,
        cost_method: 'exact',
        status: '',
      },
    ],
    submitted_at: new Date().toISOString(),
    started_at: null,
    comment: 'Extra cheese',
    sale_price: 25,
    cost: 10,
    discount: 2,
    tips: 3,
    is_paid: false,
    is_pay_later: false,
    payment_source: 'Cash',
    customer: { id: 'c1', name: 'Janez', phone: '040123', address: 'Ljubljana 1' },
    delivery_info: null,
    custom_data: null,
    ...overrides,
  }
}

describe('OrderView', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
  })

  it('renders order display_id', () => {
    const order = createOrder()
    const wrapper = mount(OrderView, {
      props: { order },
      global: {
        plugins: [i18n, ToastService, ConfirmationService],
        stubs: { Dialog: true, Popover: true, ConfirmPopup: true },
      },
    })

    expect(wrapper.text()).toContain('#001')
  })

  it('renders customer info', () => {
    const order = createOrder()
    const wrapper = mount(OrderView, {
      props: { order },
      global: {
        plugins: [i18n, ToastService, ConfirmationService],
        stubs: { Dialog: true, Popover: true, ConfirmPopup: true },
      },
    })

    expect(wrapper.text()).toContain('Janez')
    expect(wrapper.text()).toContain('Ljubljana 1')
    expect(wrapper.text()).toContain('040123')
  })

  it('shows unpaid badge when not paid', () => {
    const order = createOrder({ is_paid: false })
    const wrapper = mount(OrderView, {
      props: { order },
      global: {
        plugins: [i18n, ToastService, ConfirmationService],
        stubs: { Dialog: true, Popover: true, ConfirmPopup: true },
      },
    })

    expect(wrapper.text()).toContain('Unpaid')
  })

  it('shows paid badge when paid', () => {
    const order = createOrder({ is_paid: true })
    const wrapper = mount(OrderView, {
      props: { order },
      global: {
        plugins: [i18n, ToastService, ConfirmationService],
        stubs: { Dialog: true, Popover: true, ConfirmPopup: true },
      },
    })

    expect(wrapper.text()).toContain('Paid')
  })

  it('shows pending status', () => {
    const order = createOrder({ state: 'pending' })
    const wrapper = mount(OrderView, {
      props: { order },
      global: {
        plugins: [i18n, ToastService, ConfirmationService],
        stubs: { Dialog: true, Popover: true, ConfirmPopup: true },
      },
    })

    expect(wrapper.text()).toContain('Pending')
  })

  it('shows in_progress status', () => {
    const order = createOrder({ state: 'in_progress' })
    const wrapper = mount(OrderView, {
      props: { order },
      global: {
        plugins: [i18n, ToastService, ConfirmationService],
        stubs: { Dialog: true, Popover: true, ConfirmPopup: true },
      },
    })

    expect(wrapper.text()).toContain('In progress')
  })

  it('shows finished status', () => {
    const order = createOrder({ state: 'finished' })
    const wrapper = mount(OrderView, {
      props: { order },
      global: {
        plugins: [i18n, ToastService, ConfirmationService],
        stubs: { Dialog: true, Popover: true, ConfirmPopup: true },
      },
    })

    expect(wrapper.text()).toContain('Finished')
  })

  it('shows cancelled status', () => {
    const order = createOrder({ state: 'cancelled' })
    const wrapper = mount(OrderView, {
      props: { order },
      global: {
        plugins: [i18n, ToastService, ConfirmationService],
        stubs: { Dialog: true, Popover: true, ConfirmPopup: true },
      },
    })

    expect(wrapper.text()).toContain('Cancelled')
  })

  it('shows delivery info when present', () => {
    const delivery = { receiver_name: 'Marko', address: 'Maribor 5', phone: '041999' }
    const order = createOrder({ delivery_info: delivery })
    const wrapper = mount(OrderView, {
      props: { order },
      global: {
        plugins: [i18n, ToastService, ConfirmationService],
        stubs: { Dialog: true, Popover: true, ConfirmPopup: true },
      },
    })

    expect(wrapper.text()).toContain('Marko')
    expect(wrapper.text()).toContain('Maribor 5')
  })

  it('hides delivery info when null', () => {
    const order = createOrder({ delivery_info: null })
    const wrapper = mount(OrderView, {
      props: { order },
      global: {
        plugins: [i18n, ToastService, ConfirmationService],
        stubs: { Dialog: true, Popover: true, ConfirmPopup: true },
      },
    })

    expect(wrapper.text()).not.toContain('Delivery info')
  })

  it('shows custom_data when present', () => {
    const order = createOrder({ custom_data: { table: '5', notes: 'window seat' } })
    const wrapper = mount(OrderView, {
      props: { order },
      global: {
        plugins: [i18n, ToastService, ConfirmationService],
        stubs: { Dialog: true, Popover: true, ConfirmPopup: true },
      },
    })

    expect(wrapper.text()).toContain('table')
    expect(wrapper.text()).toContain('5')
    expect(wrapper.text()).toContain('window seat')
  })

  it('hides custom_data when null', () => {
    const order = createOrder({ custom_data: null })
    const wrapper = mount(OrderView, {
      props: { order },
      global: {
        plugins: [i18n, ToastService, ConfirmationService],
        stubs: { Dialog: true, Popover: true, ConfirmPopup: true },
      },
    })

    expect(wrapper.text()).not.toContain('Custom data')
  })

  it('shows collect_money button when unpaid', () => {
    const order = createOrder({ is_paid: false, state: 'in_progress' })
    const wrapper = mount(OrderView, {
      props: { order },
      global: {
        plugins: [i18n, ToastService, ConfirmationService],
        stubs: { Dialog: true, Popover: true, ConfirmPopup: true },
      },
    })

    expect(wrapper.text()).toContain('Collect money')
  })

  it('hides collect_money button when paid', () => {
    const order = createOrder({ is_paid: true })
    const wrapper = mount(OrderView, {
      props: { order },
      global: {
        plugins: [i18n, ToastService, ConfirmationService],
        stubs: { Dialog: true, Popover: true, ConfirmPopup: true },
      },
    })

    expect(wrapper.text()).not.toContain('Collect money')
  })

  it('shows discount and tips values', () => {
    const order = createOrder({ discount: 5, tips: 3.5 })
    const wrapper = mount(OrderView, {
      props: { order },
      global: {
        plugins: [i18n, ToastService, ConfirmationService],
        stubs: { Dialog: true, Popover: true, ConfirmPopup: true },
      },
    })

    expect(wrapper.text()).toContain('5')
    expect(wrapper.text()).toContain('3.5')
  })

  it('shows comment', () => {
    const order = createOrder({ comment: 'No onions please' })
    const wrapper = mount(OrderView, {
      props: { order },
      global: {
        plugins: [i18n, ToastService, ConfirmationService],
        stubs: { Dialog: true, Popover: true, ConfirmPopup: true },
      },
    })

    expect(wrapper.text()).toContain('No onions please')
  })
})

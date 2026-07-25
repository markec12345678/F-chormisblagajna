import { describe, it, expect, vi, beforeEach } from 'vitest'
import { createPinia, setActivePinia } from 'pinia'
import { mount } from '@vue/test-utils'
import OrderItemsInfo from '@/components/OrderItemsInfo.vue'
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

vi.mock('primevue/usetoast', () => ({
  useToast: () => ({ add: vi.fn() }),
}))

const i18n = createI18n({
  legacy: false,
  locale: 'en',
  fallbackLocale: 'en',
  messages: {
    en: {
      name: 'Name',
      price: 'Price',
      quantity: 'Quantity',
      cost_method: 'Cost Method',
      status: 'Status',
      actions: 'Actions',
      return: 'Return',
      refund: 'Refund',
      normal: 'Normal',
      pending: 'Pending',
      refunded: 'Refunded',
      returning: 'Returning',
      order: 'Order',
      material: 'Materials',
    },
  },
})

const stubs = {
  DataTable: {
    template:
      '<div class="datatable-stub"><template v-for="(item, idx) in value" :key="idx"><div class="row-stub"><slot :data="item" :index="idx"/></div></template></div>',
    props: ['value', 'expandedRows', 'stripedRows', 'tableStyle'],
  },
  Column: {
    template: '<div class="col-stub"><slot :data="data"/></div>',
    props: ['header', 'field', 'expander', 'style'],
  },
  Badge: { template: '<span class="badge-stub">{{ value }}</span>', props: ['value', 'severity'] },
  Button: {
    template: '<button @click="$emit(\'click\')">{{ label }}</button>',
    props: ['icon', 'label', 'severity', 'aria-label'],
  },
  ButtonGroup: { template: '<div class="bg-stub"><slot/></div>' },
  Dialog: {
    template: '<div class="dialog-stub"><slot/></div>',
    props: ['visible', 'modal', 'header', 'style', 'breakpoints'],
  },
  OrderItemRefund: { template: '<div class="refund-stub"/>' },
}

const makeOrder = (items: OrderItem[]) => {
  const order = new Order()
  order.display_id = 'ORD-1'
  order.items = items
  return order
}

const makeItem = (name: string, status: string) => {
  const item = new OrderItem()
  item.product.name = name
  item.status = status
  item.materials = []
  return item
}

describe('OrderItemsInfo', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
  })

  it('renders DataTable with order items', () => {
    const wrapper = mount(OrderItemsInfo, {
      props: { order: makeOrder([makeItem('Pizza', '')]) },
      global: { plugins: [PrimeVue, i18n], stubs },
    })
    expect(wrapper.find('.datatable-stub').exists()).toBe(true)
  })

  it('renders row for each item', () => {
    const wrapper = mount(OrderItemsInfo, {
      props: { order: makeOrder([makeItem('Pizza', ''), makeItem('Burger', 'pending')]) },
      global: { plugins: [PrimeVue, i18n], stubs },
    })
    expect(wrapper.findAll('.row-stub').length).toBe(2)
  })

  it('passes item data to row slot', () => {
    const wrapper = mount(OrderItemsInfo, {
      props: { order: makeOrder([makeItem('Pizza', '')]) },
      global: { plugins: [PrimeVue, i18n], stubs },
    })
    expect(wrapper.find('.row-stub').exists()).toBe(true)
  })

  it('has expandedRows ref', () => {
    const wrapper = mount(OrderItemsInfo, {
      props: { order: makeOrder([makeItem('Pizza', '')]) },
      global: { plugins: [PrimeVue, i18n], stubs },
    })
    expect(Array.isArray(wrapper.vm.expandedRows)).toBe(true)
    expect(wrapper.vm.expandedRows).toHaveLength(0)
  })

  it('has dialog hidden by default', () => {
    const wrapper = mount(OrderItemsInfo, {
      props: { order: makeOrder([makeItem('Pizza', '')]) },
      global: { plugins: [PrimeVue, i18n], stubs },
    })
    expect(wrapper.vm.item_refund_dialog).toBe(false)
  })

  it('does not emit updated initially', () => {
    const wrapper = mount(OrderItemsInfo, {
      props: { order: makeOrder([makeItem('Pizza', '')]) },
      global: { plugins: [PrimeVue, i18n], stubs },
    })
    expect(wrapper.emitted('updated')).toBeUndefined()
  })

  it('accepts order with empty items', () => {
    const wrapper = mount(OrderItemsInfo, {
      props: { order: makeOrder([]) },
      global: { plugins: [PrimeVue, i18n], stubs },
    })
    expect(wrapper.find('.datatable-stub').exists()).toBe(true)
    expect(wrapper.findAll('.row-stub')).toHaveLength(0)
  })
})

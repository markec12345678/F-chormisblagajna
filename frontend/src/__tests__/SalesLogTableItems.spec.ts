import { describe, it, expect, vi, beforeEach } from 'vitest'
import { createPinia, setActivePinia } from 'pinia'
import { mount, flushPromises } from '@vue/test-utils'
import SalesLogTableItems from '@/components/SalesLogTableItems.vue'
import PrimeVue from 'primevue/config'
import { createI18n } from 'vue-i18n'

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
      cost: 'Cost',
      quantity: 'Quantity',
      sale_price: 'Sale Price',
      refunds: 'Refunds',
      profit: 'Profit',
      refunded: 'Refunded',
      entry: 'Entry',
    },
  },
})

const stubs = {
  DataTable: {
    template: '<div class="datatable-stub"><slot name="header"/><slot name="empty"/><slot/></div>',
    props: ['value', 'expandedRows', 'stripedRows', 'tableStyle'],
  },
  Column: {
    template: '<div><slot :data="data"/></div>',
    props: ['header', 'field', 'sortable', 'style'],
  },
  Badge: { template: '<span class="badge-stub">{{ value }}</span>', props: ['value', 'severity'] },
}

describe('SalesLogTableItems', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
  })

  const makeItems = () => [
    {
      item_id: 'item-1',
      item_name: 'Pizza',
      cost: 5,
      quantity: 2,
      sale_price: 12,
      components: [
        { component_id: 'comp-1', component_name: 'Flour', cost: 2, quantity: 1, item_id: '' },
      ],
    },
    {
      item_id: 'item-2',
      item_name: 'Burger',
      cost: 3,
      quantity: 1,
      sale_price: 8,
      components: [],
    },
  ]

  const makeOrderRefunds = () => ({ refunds: [] })

  it('renders with items and no refunds', () => {
    const wrapper = mount(SalesLogTableItems, {
      props: { items: makeItems(), order_refunds: makeOrderRefunds() },
      global: { plugins: [PrimeVue, i18n], stubs },
    })
    expect(wrapper.find('.datatable-stub').exists()).toBe(true)
  })

  it('initializes localItems from props', () => {
    const wrapper = mount(SalesLogTableItems, {
      props: { items: makeItems(), order_refunds: makeOrderRefunds() },
      global: { plugins: [PrimeVue, i18n], stubs },
    })
    expect(wrapper.vm.localItems).toHaveLength(2)
  })

  it('has expanded rows ref', () => {
    const wrapper = mount(SalesLogTableItems, {
      props: { items: makeItems(), order_refunds: makeOrderRefunds() },
      global: { plugins: [PrimeVue, i18n], stubs },
    })
    expect(Array.isArray(wrapper.vm.expandedSalesLogOrderItemComponents)).toBe(true)
  })

  it('handles items with refunds', () => {
    const items = makeItems()
    const order_refunds = {
      refunds: [
        {
          order_item_id: 'item-1',
          amount: 6,
          reason: 'wrong order',
          item_cost: 5,
          destination: 'inventory',
          material_refunds: [],
        },
      ],
    }
    const wrapper = mount(SalesLogTableItems, {
      props: { items, order_refunds },
      global: { plugins: [PrimeVue, i18n], stubs },
    })
    expect(wrapper.vm.items_refunds['item-1']).toBeDefined()
    expect(wrapper.vm.items_refunds['item-1'].is_refunded).toBe(true)
    expect(wrapper.vm.items_refunds['item-1'].refund_amount).toBe(6)
  })

  it('handles custom destination refunds', () => {
    const items = makeItems()
    const order_refunds = {
      refunds: [
        {
          order_item_id: 'item-1',
          amount: 6,
          reason: 'wrong',
          item_cost: 5,
          destination: 'custom',
          material_refunds: [{ material_id: 'mat-1', inventory_return_qty: 2, cost_per_unit: 1.5 }],
        },
      ],
    }
    const wrapper = mount(SalesLogTableItems, {
      props: { items, order_refunds },
      global: { plugins: [PrimeVue, i18n], stubs },
    })
    expect(wrapper.vm.items_refunds['item-1'].material_refunds['mat-1']).toBeDefined()
    expect(wrapper.vm.items_refunds['item-1'].inventory_refunds).toBe(3)
  })

  it('handles empty items array', () => {
    const wrapper = mount(SalesLogTableItems, {
      props: { items: [], order_refunds: makeOrderRefunds() },
      global: { plugins: [PrimeVue, i18n], stubs },
    })
    expect(wrapper.vm.localItems).toHaveLength(0)
  })
})

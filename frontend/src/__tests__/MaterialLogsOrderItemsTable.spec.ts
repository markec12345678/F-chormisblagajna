import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import MaterialLogsOrderItemsTable from '@/components/MaterialLogsOrderItemsTable.vue'
import PrimeVue from 'primevue/config'
import { createI18n } from 'vue-i18n'

const i18n = createI18n({
  legacy: false,
  locale: 'en',
  fallbackLocale: 'en',
  messages: { en: { name: 'Name', materials: 'Materials' } },
})

const stubs = {
  DataTable: { template: '<div><slot/></div>', props: ['value'] },
  Column: {
    template: '<div><slot :index="index" :data="data"/></div>',
    props: ['header', 'field'],
  },
  Tag: { template: '<span class="tag-stub"><slot/></span>', props: ['icon', 'class'] },
}

describe('MaterialLogsOrderItemsTable', () => {
  const items = [
    {
      product: { name: 'Pizza' },
      materials: [
        { material: { name: 'Flour', unit: 'kg' }, quantity: 2 },
        { material: { name: 'Cheese', unit: 'g' }, quantity: 200 },
      ],
    },
    {
      product: { name: 'Burger' },
      materials: [{ material: { name: 'Beef', unit: 'g' }, quantity: 300 }],
    },
  ]

  it('renders table with items', () => {
    const wrapper = mount(MaterialLogsOrderItemsTable, {
      props: { items },
      global: { plugins: [PrimeVue, i18n], stubs },
    })
    expect(wrapper.exists()).toBe(true)
  })

  it('accepts empty items array', () => {
    const wrapper = mount(MaterialLogsOrderItemsTable, {
      props: { items: [] },
      global: { plugins: [PrimeVue, i18n], stubs },
    })
    expect(wrapper.exists()).toBe(true)
  })

  it('accepts order_item_index prop', () => {
    const wrapper = mount(MaterialLogsOrderItemsTable, {
      props: { items, order_item_index: 0 },
      global: { plugins: [PrimeVue, i18n], stubs },
    })
    expect(wrapper.props().order_item_index).toBe(0)
  })

  it('accepts undefined order_item_index', () => {
    const wrapper = mount(MaterialLogsOrderItemsTable, {
      props: { items },
      global: { plugins: [PrimeVue, i18n], stubs },
    })
    expect(wrapper.props().order_item_index).toBeUndefined()
  })
})

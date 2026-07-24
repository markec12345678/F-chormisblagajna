import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import InventoryItem from '@/components/InventoryItem.vue'

describe('InventoryItem', () => {
  it('renders props correctly', () => {
    const wrapper = mount(InventoryItem, {
      props: { name: 'Milk', quantity: 5, unit: 'L' },
    })
    expect(wrapper.text()).toContain('Milk')
    expect(wrapper.text()).toContain('5')
    expect(wrapper.text()).toContain('L')
  })

  it('renders with default empty props', () => {
    const wrapper = mount(InventoryItem)
    expect(wrapper.find('.container').exists()).toBe(true)
  })

  it('renders with numeric quantity', () => {
    const wrapper = mount(InventoryItem, {
      props: { name: 'Flour', quantity: 25, unit: 'kg' },
    })
    expect(wrapper.text()).toContain('Flour')
    expect(wrapper.text()).toContain('25')
    expect(wrapper.text()).toContain('kg')
  })

  it('renders quantity-container with background', () => {
    const wrapper = mount(InventoryItem, {
      props: { name: 'Oil', quantity: 10, unit: 'L' },
    })
    const qtyContainer = wrapper.find('.quantity-container')
    expect(qtyContainer.exists()).toBe(true)
  })
})

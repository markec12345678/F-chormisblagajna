import { describe, it, expect, beforeEach, vi } from 'vitest'
import { createPinia, setActivePinia } from 'pinia'
import Order from '@/classes/Order'

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

describe('Order', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
  })

  it('should initialize with default values', () => {
    const order = new Order()

    expect(order.id).toBe('')
    expect(order.display_id).toBe('')
    expect(order.state).toBe('')
    expect(order.comment).toBe('')
    expect(order.payment_source).toBe('')
    expect(order.discount).toBe(0)
    expect(order.sale_price).toBe(0)
    expect(order.tips).toBe(0)
    expect(order.is_paid).toBe(false)
    expect(order.is_auto_start).toBe(false)
    expect(order.items).toEqual([])
    expect(order.customer).toEqual({ id: '', name: '', phone: '', address: '' })
    expect(order.delivery_info).toBeNull()
    expect(order.custom_data).toBeNull()
  })

  it('should create dates on construction', () => {
    const before = new Date()
    const order = new Order()
    const after = new Date()

    expect(order.submitted_at).toBeInstanceOf(Date)
    expect(order.started_at).toBeInstanceOf(Date)
    expect(order.submitted_at.getTime()).toBeGreaterThanOrEqual(before.getTime())
    expect(order.submitted_at.getTime()).toBeLessThanOrEqual(after.getTime())
  })

  it('should allow setting all properties', () => {
    const order = new Order()
    order.id = 'abc123'
    order.display_id = 'D-001'
    order.state = 'pending'
    order.comment = 'extra spicy'
    order.payment_source = 'cash'
    order.discount = 10
    order.sale_price = 150.5
    order.tips = 20
    order.is_paid = true
    order.is_auto_start = true

    expect(order.id).toBe('abc123')
    expect(order.display_id).toBe('D-001')
    expect(order.state).toBe('pending')
    expect(order.comment).toBe('extra spicy')
    expect(order.payment_source).toBe('cash')
    expect(order.discount).toBe(10)
    expect(order.sale_price).toBe(150.5)
    expect(order.tips).toBe(20)
    expect(order.is_paid).toBe(true)
    expect(order.is_auto_start).toBe(true)
  })

  it('should set customer details', () => {
    const order = new Order()
    order.customer = { id: 'c1', name: 'Janez', phone: '+386123', address: 'Ljubljana 1' }

    expect(order.customer.id).toBe('c1')
    expect(order.customer.name).toBe('Janez')
    expect(order.customer.phone).toBe('+386123')
    expect(order.customer.address).toBe('Ljubljana 1')
  })

  it('should set delivery_info', () => {
    const order = new Order()
    expect(order.delivery_info).toBeNull()

    order.delivery_info = { receiver_name: 'Peter', address: 'Maribor 5', phone: '+386456' }
    expect(order.delivery_info.receiver_name).toBe('Peter')
    expect(order.delivery_info.address).toBe('Maribor 5')
  })

  it('should set custom_data', () => {
    const order = new Order()
    expect(order.custom_data).toBeNull()

    order.custom_data = { table: '5', section: 'outdoor' }
    expect(order.custom_data?.table).toBe('5')
    expect(order.custom_data?.section).toBe('outdoor')
  })
})

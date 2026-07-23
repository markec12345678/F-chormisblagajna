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
    expect(order.customer).toEqual({})
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
})

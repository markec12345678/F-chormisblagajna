import { describe, it, expect, beforeEach, vi } from 'vitest'
import { createPinia, setActivePinia } from 'pinia'

vi.mock('@/stores', () => ({
  globalStore: () => ({
    shop_mode: '',
    orientation: 'ltr',
    getColorMode: 'light',
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

import { Material, MaterialEntry, Product } from '@/classes/OrderItem'

describe('Material', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
  })

  it('should initialize with default values', () => {
    const material = new Material()

    expect(material.id).toBe('')
    expect(material.name).toBe('')
    expect(material.quantity).toBe(0)
    expect(material.entries).toEqual([])
    expect(material.unit).toBe('')
    expect(material.label).toBe('')
  })
})

describe('MaterialEntry', () => {
  it('should initialize with default values', () => {
    const entry = new MaterialEntry()

    expect(entry.id).toBe('')
    expect(entry.purchase_quantity).toBe(0)
    expect(entry.purchase_price).toBe(0)
    expect(entry.quantity).toBe(0)
    expect(entry.company).toBe('')
    expect(entry.cost).toBe(0)
    expect(entry.sku).toBe('')
    expect(entry.label).toBe('')
  })
})

describe('Product', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
  })

  it('should initialize with default values', () => {
    const product = new Product()

    expect(product.id).toBe('')
    expect(product.name).toBe('')
    expect(product.materials).toEqual([])
    expect(product.sub_products).toEqual([])
    expect(product.price).toBe(0)
    expect(product.quantity).toBe(0)
  })
})

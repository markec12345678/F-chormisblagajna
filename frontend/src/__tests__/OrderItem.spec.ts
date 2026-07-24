import { describe, it, expect, beforeEach, vi } from 'vitest'
import { createPinia, setActivePinia } from 'pinia'

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

import {
  Material,
  MaterialEntry,
  MaterialSettings,
  Product,
  ProductEntry,
  OrderItem,
  OrderItemMaterial,
} from '@/classes/OrderItem'

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

  it('should allow setting properties', () => {
    const material = new Material()
    material.id = 'm1'
    material.name = 'Flour'
    material.quantity = 10
    material.unit = 'kg'
    material.label = 'Flour Label'

    expect(material.id).toBe('m1')
    expect(material.name).toBe('Flour')
    expect(material.quantity).toBe(10)
    expect(material.unit).toBe('kg')
    expect(material.label).toBe('Flour Label')
  })

  it('should have settings with default treshold', () => {
    const material = new Material()
    expect(material.settings).toBeDefined()
    expect(material.settings.stock_alert_treshold).toBe(0)
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

  it('should allow setting all properties', () => {
    const entry = new MaterialEntry()
    entry.id = 'e1'
    entry.purchase_quantity = 50
    entry.purchase_price = 25.5
    entry.quantity = 100
    entry.company = 'Acme Corp'
    entry.cost = 12.75
    entry.sku = 'FLR-001'
    entry.label = 'Acme Corp - 100 kg'

    expect(entry.id).toBe('e1')
    expect(entry.purchase_quantity).toBe(50)
    expect(entry.purchase_price).toBe(25.5)
    expect(entry.quantity).toBe(100)
    expect(entry.company).toBe('Acme Corp')
    expect(entry.cost).toBe(12.75)
    expect(entry.sku).toBe('FLR-001')
    expect(entry.label).toBe('Acme Corp - 100 kg')
  })
})

describe('MaterialSettings', () => {
  it('should initialize with default treshold 0', () => {
    const settings = new MaterialSettings()
    expect(settings.stock_alert_treshold).toBe(0)
  })

  it('should allow setting treshold', () => {
    const settings = new MaterialSettings()
    settings.stock_alert_treshold = 10
    expect(settings.stock_alert_treshold).toBe(10)
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
    expect(product.entries).toEqual([])
    expect(product.price).toBe(0)
    expect(product.image_url).toBe('')
    expect(product.measuring_unit).toBe('')
    expect(product.quantity).toBe(0)
    expect(product.ready).toBe(0)
  })

  it('should allow setting all properties', () => {
    const product = new Product()
    product.id = 'p1'
    product.name = 'Pizza'
    product.price = 12.5
    product.image_url = 'pizza.jpg'
    product.measuring_unit = 'pcs'
    product.quantity = 1
    product.ready = 5

    expect(product.id).toBe('p1')
    expect(product.name).toBe('Pizza')
    expect(product.price).toBe(12.5)
    expect(product.image_url).toBe('pizza.jpg')
    expect(product.measuring_unit).toBe('pcs')
    expect(product.quantity).toBe(1)
    expect(product.ready).toBe(5)
  })
})

describe('ProductEntry', () => {
  it('should initialize with default values', () => {
    const entry = new ProductEntry()

    expect(entry.id).toBe('')
    expect(entry.purchase_quantity).toBe(0)
    expect(entry.purchase_price).toBe(0)
    expect(entry.quantity).toBe(0)
    expect(entry.company).toBe('')
    expect(entry.unit).toBe('')
    expect(entry.sku).toBe('')
  })
})

describe('OrderItemMaterial', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
  })

  it('should initialize without material', () => {
    const itemMat = new OrderItemMaterial()

    expect(itemMat.material).toBeDefined()
    expect(itemMat.entry).toBeDefined()
    expect(itemMat.entries).toEqual([])
    expect(itemMat.quantity).toBe(0)
    expect(itemMat.avgcost).toBe(0)
    expect(itemMat.isQuantityValid).toBe(true)
  })

  it('should initialize with material', () => {
    const material = new Material()
    material.entries = [
      { ...new MaterialEntry(), id: 'e1', company: 'Acme', quantity: 50, purchase_quantity: 50, purchase_price: 10 },
    ]
    material.unit = 'kg'
    material.quantity = 5

    const itemMat = new OrderItemMaterial(material)

    expect(itemMat.material).toBe(material)
    expect(itemMat.entry).toBe(material.entries[0])
    expect(itemMat.entries).toEqual(material.entries)
    expect(itemMat.quantity).toBe(0)
  })
})

describe('OrderItem', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
  })

  it('should initialize without product', () => {
    const item = new OrderItem()

    expect(item.product).toBeDefined()
    expect(item.materials).toEqual([])
    expect(item.comment).toBe('')
    expect(item.Id).toBe('')
    expect(item.is_consume_from_ready).toBe(false)
    expect(item.ready).toBe(0)
    expect(item.sub_items).toEqual([])
    expect(item.quantity).toBe(1)
    expect(item.can_change_ready_toggle).toBe(false)
    expect(item.price).toBe(0)
    expect(item.isValid).toBe(true)
  })

  it('should allow setting SetProductId', () => {
    const item = new OrderItem()
    item.SetProductId('abc123')

    expect(item.Id).toBe('abc123')
    expect(item.product.id).toBe('abc123')
  })

  it('should validate item as valid with no materials', () => {
    const item = new OrderItem()
    item.materials = []
    item.ValidateItem()
    expect(item.isValid).toBe(true)
  })

  it('should invalidate item when material is invalid', () => {
    const item = new OrderItem()
    item.materials = [{ ...new OrderItemMaterial(), isQuantityValid: false } as any]
    item.ValidateItem()
    expect(item.isValid).toBe(false)
  })

  it('should remove material by index', () => {
    const item = new OrderItem()
    item.materials = [
      { ...new OrderItemMaterial(), isQuantityValid: true } as any,
      { ...new OrderItemMaterial(), isQuantityValid: true } as any,
    ]
    item.RemoveMaterialByIndex(0)
    expect(item.materials.length).toBe(1)
  })
})

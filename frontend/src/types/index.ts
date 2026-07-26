export interface Category {
  id: string
  name: string
  products: CategoryProduct[]
}

export interface CategoryProduct {
  id: string
  name: string
  price: number
  image_url: string
  index?: number
}

export interface Customer {
  id: string
  name: string
  phone: string
  address: string
}

export interface OrderListItem {
  id: string
  display_id: string
  state: string
  sale_price: number
  cost: number
  submitted_at: string
  is_pay_later: boolean
  is_paid: boolean
  payment_source: string
  is_delivery: boolean
  is_take_away: boolean
  is_dine_in: boolean
  comment: string
  tips: number
  customer: Customer
  discount: number
  is_fiscalized: boolean
  fiscal_id: string
}

export interface OrderQueue {
  prefix: string
  next: number
}

export interface Language {
  language: string
  code: string
}

export interface PaymentSource {
  name: string
}

export interface DataTablePageEvent {
  first: number
  rows: number
}

export interface ProductItem {
  id: string
  name: string
  price: number
  image_url: string
  materials: ProductItemMaterial[]
  sub_products: ProductItemSubProduct[]
  unit: string
  quantity: number
  ready: number
  enable_inventory_consumption: boolean
  enable_fixed_cost: boolean
  fixed_cost: number
}

export interface ProductItemMaterial {
  id: string
  name: string
  quantity: number
  unit: string
  index?: number
}

export interface ProductItemSubProduct {
  id: string
  name: string
  quantity: number
  index?: number
}

export interface TableItem {
  id: string
  number: number
  name: string
  capacity: number
  zone: string
  status: string
  qr_code: string
  order_id: string
}

export interface AISearchResult {
  product_id: string
  name: string
  description: string
  price: number
  score: number
  category: string
}

export interface SmartSuggestion {
  product_id: string
  product_name: string
  reason: string
  score: number
}

import { OrderItem } from '@/classes/OrderItem'

export interface Customer {
  id: string
  name: string
  phone: string
  address: string
}

export interface DeliveryInfo {
  receiver_name: string
  address: string
  phone: string
}

export default class Order {
  submitted_at: Date
  id: string
  display_id: string
  items: OrderItem[]
  discount: number
  state: string
  started_at: Date
  comment: string
  sale_price: number
  is_auto_start: boolean
  is_paid: boolean
  tips: number
  payment_source: string
  customer: Customer
  delivery_info: DeliveryInfo | null
  custom_data: Record<string, string> | null

  constructor() {
    this.submitted_at = new Date()
    this.id = ''
    this.display_id = ''
    this.items = []
    this.discount = 0
    this.state = ''
    this.started_at = new Date()
    this.comment = ''
    this.sale_price = 0
    this.is_auto_start = false
    this.is_paid = false
    this.tips = 0
    this.payment_source = ''
    this.customer = { id: '', name: '', phone: '', address: '' }
    this.delivery_info = null
    this.custom_data = null
  }
}

import { describe, it, expect } from 'vitest'
import { Notification, OrderFinishNotification } from '@/classes/Notification'

describe('Notification', () => {
  it('should initialize with default values', () => {
    const notification = new Notification()

    expect(notification.id).toBeTruthy()
    expect(notification.description).toBe('')
    expect(notification.type).toBe('')
    expect(notification.topic_name).toBe('')
    expect(notification.severity).toBe('')
    expect(notification.date).toBeInstanceOf(Date)
  })

  it('should generate unique ids', () => {
    const notification1 = new Notification()
    const notification2 = new Notification()

    expect(notification1.id).not.toBe(notification2.id)
  })

  it('should allow setting all properties', () => {
    const notification = new Notification()
    notification.description = 'Order #5 ready'
    notification.type = 'topic_message'
    notification.topic_name = 'order_finished'
    notification.severity = 'success'

    expect(notification.description).toBe('Order #5 ready')
    expect(notification.type).toBe('topic_message')
    expect(notification.topic_name).toBe('order_finished')
    expect(notification.severity).toBe('success')
  })

  it('should have id as non-empty string', () => {
    const notification = new Notification()
    expect(typeof notification.id).toBe('string')
    expect(notification.id.length).toBeGreaterThan(0)
  })
})

describe('OrderFinishNotification', () => {
  it('should initialize with default values', () => {
    const notification = new OrderFinishNotification()

    expect(notification.id).toBeTruthy()
    expect(notification.order_id).toBe('')
  })

  it('should allow setting order_id', () => {
    const notification = new OrderFinishNotification()
    notification.order_id = 'order-123'

    expect(notification.order_id).toBe('order-123')
  })

  it('should have unique id', () => {
    const n1 = new OrderFinishNotification()
    const n2 = new OrderFinishNotification()
    expect(n1.id).not.toBe(n2.id)
  })
})

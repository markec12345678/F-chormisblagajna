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
})

describe('OrderFinishNotification', () => {
  it('should initialize with default values', () => {
    const notification = new OrderFinishNotification()
    
    expect(notification.id).toBeTruthy()
    expect(notification.order_id).toBe('')
  })
})

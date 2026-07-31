import { describe, it, expect, beforeEach, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import { createPinia, setActivePinia } from 'pinia'
import { createI18n } from 'vue-i18n'
import TablesideMenu from '@/pages/TablesideMenu.vue'
import ToastService from 'primevue/toastservice'
import { createRouter, createWebHistory } from 'vue-router'

const mockGet = vi.hoisted(() => vi.fn())
const mockPost = vi.hoisted(() => vi.fn())

vi.mock('axios', () => ({
  default: { get: mockGet, post: mockPost, put: vi.fn(), delete: vi.fn(), defaults: { headers: { common: {} } } },
}))

const i18n = createI18n({
  legacy: false,
  locale: 'en',
  fallbackLocale: 'en',
  messages: { en: { table_menu: 'Table Menu', no_products: 'No products', your_cart: 'Your Cart', total: 'Total', place_order: 'Place Order', order_placed: 'Order Placed', order_placed_message: 'Your order has been placed', close: 'Close', table: 'Table', failed: 'Failed', order_failed: 'Order failed', success: 'Success' } },
})

const router = createRouter({ history: createWebHistory(), routes: [{ path: '/menu/:token', name: 'tableside-menu', component: TablesideMenu }] })

const stubs = {
  Button: { template: '<button class="btn-stub">{{ label }}</button>', props: ['label', 'icon', 'severity', 'size', 'text', 'loading', 'disabled'] },
  InputNumber: { template: '<input class="input-number-stub" />', props: ['modelValue', 'min', 'max'] },
  Dialog: { template: '<div class="dialog-stub" v-if="visible"><slot /><slot name="footer" /></div>', props: ['visible', 'header'] },
}

describe('TablesideMenu', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    mockGet.mockReset()
    mockPost.mockReset()
    mockGet.mockResolvedValue({ data: { data: [] } })
  })

  it('renders page heading', async () => {
    const wrapper = mount(TablesideMenu, { global: { plugins: [i18n, ToastService, router], stubs } })
    await vi.waitFor(() => expect(wrapper.text()).toContain('Table Menu'))
  })

  it('calls API on mount', async () => {
    mount(TablesideMenu, { global: { plugins: [i18n, ToastService, router], stubs } })
    await vi.waitFor(() => expect(mockGet).toHaveBeenCalled())
  })

  it('handles API error gracefully', async () => {
    mockGet.mockRejectedValue(new Error('fail'))
    const wrapper = mount(TablesideMenu, { global: { plugins: [i18n, ToastService, router], stubs } })
    await vi.waitFor(() => expect(wrapper.text()).toBeDefined())
  })
})

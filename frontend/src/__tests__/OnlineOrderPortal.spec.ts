import { describe, it, expect, beforeEach, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import { createPinia, setActivePinia } from 'pinia'
import { createI18n } from 'vue-i18n'
import OnlineOrderPortal from '@/pages/OnlineOrderPortal.vue'
import ToastService from 'primevue/toastservice'

const mockGet = vi.hoisted(() => vi.fn())
const mockPost = vi.hoisted(() => vi.fn())

vi.mock('axios', () => ({
  default: { get: mockGet, post: mockPost, put: vi.fn(), delete: vi.fn(), defaults: { headers: { common: {} } } },
}))

const i18n = createI18n({
  legacy: false,
  locale: 'en',
  fallbackLocale: 'en',
  messages: { en: { order_online: 'Order Online', cart: 'Cart', all: 'All', no_products_available: 'No products available', your_cart: 'Your Cart', cart_empty: 'Cart is empty', item: 'Item', quantity: 'Qty', subtotal: 'Subtotal', actions: 'Actions', customer_name: 'Name', phone: 'Phone', email: 'Email', order_type: 'Order Type', notes: 'Notes', delivery_address: 'Delivery Address', total: 'Total', cancel: 'Cancel', place_order: 'Place Order', warning: 'Warning', fill_required_fields: 'Fill required fields', success: 'Success', order_placed: 'Order placed', failed: 'Failed', order_failed: 'Order failed', menu_load_failed: 'Menu load failed' } },
})

const stubs = {
  Toolbar: { template: '<div class="toolbar-stub"><slot name="start" /><slot name="end" /></div>' },
  Card: { template: '<div class="card-stub"><slot name="header" /><slot name="title" /><slot name="content" /></div>' },
  Button: { template: '<button class="btn-stub">{{ label }}</button>', props: ['label', 'icon', 'severity', 'text', 'rounded', 'size', 'loading', 'badge', 'badgeClass'] },
  Tag: { template: '<span class="tag-stub"><slot /></span>', props: ['severity', 'value'] },
  Dialog: { template: '<div class="dialog-stub" v-if="visible"><slot /><slot name="footer" /></div>', props: ['visible', 'header'] },
  DataTable: { template: '<div class="datatable-stub"><slot name="empty" /><slot /></div>', props: ['value'] },
  Column: { template: '<div class="column-stub"><slot name="body" :data="{}" /></div>', props: ['field', 'header'] },
  InputText: { template: '<input class="input-text-stub" />', props: ['modelValue'] },
  Textarea: { template: '<textarea class="textarea-stub" />', props: ['modelValue'] },
  Dropdown: { template: '<select class="dropdown-stub" />', props: ['modelValue', 'options', 'optionLabel', 'optionValue'] },
  Divider: { template: '<hr class="divider-stub" />' },
  ProgressSpinner: { template: '<div class="spinner-stub" />' },
}

describe('OnlineOrderPortal', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    mockGet.mockReset()
    mockPost.mockReset()
    mockGet.mockResolvedValue({ data: { data: [] } })
  })

  it('renders page heading', async () => {
    const wrapper = mount(OnlineOrderPortal, { global: { plugins: [i18n, ToastService], stubs } })
    await vi.waitFor(() => expect(wrapper.text()).toContain('Order Online'))
  })

  it('calls API to load menu', async () => {
    mount(OnlineOrderPortal, { global: { plugins: [i18n, ToastService], stubs } })
    await vi.waitFor(() => expect(mockGet).toHaveBeenCalled())
  })

  it('shows cart button', async () => {
    const wrapper = mount(OnlineOrderPortal, { global: { plugins: [i18n, ToastService], stubs } })
    await vi.waitFor(() => expect(wrapper.text()).toContain('Cart'))
  })

  it('shows category filter buttons', async () => {
    const wrapper = mount(OnlineOrderPortal, { global: { plugins: [i18n, ToastService], stubs } })
    await vi.waitFor(() => expect(wrapper.text()).toContain('All'))
  })

  it('shows loading spinner initially', async () => {
    mockGet.mockImplementation(() => new Promise(() => {}))
    const wrapper = mount(OnlineOrderPortal, { global: { plugins: [i18n, ToastService], stubs } })
    expect(wrapper.find('.spinner-stub').exists()).toBe(true)
  })
})

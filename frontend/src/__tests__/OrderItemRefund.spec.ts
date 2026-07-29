import { describe, it, expect, vi, beforeEach } from 'vitest'
import { createPinia, setActivePinia } from 'pinia'
import { mount, flushPromises } from '@vue/test-utils'
import OrderItemRefund from '@/components/OrderItemRefund.vue'
import { createI18n } from 'vue-i18n'

vi.mock('@/stores', () => ({
  globalStore: () => ({
    shop_mode: '',
    orientation: 'ltr',
    getColorMode: 'light',
  }),
}))

vi.mock('@/services/auth', () => ({
  default: { accessToken: { value: 'token' } },
}))

const mockPost = vi.fn().mockResolvedValue({ data: {} })
vi.mock('axios', () => ({
  default: {
    get: vi.fn().mockResolvedValue({ data: { data: {} } }),
    post: (...args: unknown[]) => mockPost(...args),
  },
}))

vi.mock('primevue/usetoast', () => ({
  useToast: () => ({ add: vi.fn() }),
}))

const i18n = createI18n({
  legacy: false,
  locale: 'en',
  fallbackLocale: 'en',
  messages: {
    en: {
      money_to_refund: 'Money to refund',
      destination: 'Destination',
      inventory: 'Inventory',
      dispose: 'Dispose',
      waste: 'Waste',
      custom_materials_products: 'Custom',
      materials: 'Materials',
      reason: 'Reason',
      submit: 'Submit',
      success: 'Success',
      failed: 'Failed',
      item_refunded: 'Item refunded',
      inventory_return: 'Inventory return',
      product: 'Products',
      quantity: 'Quantity',
      unit: 'Unit',
      cost: 'Cost',
      exact: 'Exact',
      average: 'Average',
    },
  },
})

function createItem(overrides = {}) {
  return {
    id: 'item-1',
    product: { id: 'p1', name: 'Pizza', materials: [] },
    price: 25.5,
    quantity: 2,
    materials: [
      {
        material: { id: 'm1', name: 'Flour', unit: 'kg' },
        entry: { id: 'e1', cost: 2.5, quantity: 50 },
        quantity: 1,
      },
    ],
    comment: '',
    sale_price: 25.5,
    cost: 10,
    cost_method: 'exact',
    isValid: true,
    ...overrides,
  }
}

function createOrder(overrides = {}) {
  return {
    id: 'order-1',
    display_id: '#001',
    ...overrides,
  }
}

const stubs = {
  Slider: {
    template: '<div class="slider-stub" />',
    props: ['modelValue', 'max'],
  },
  RadioButton: {
    template:
      '<input type="radio" :value="value" :checked="modelValue === value" @change="$emit(\'update:modelValue\', value)" />',
    props: ['modelValue', 'value', 'inputId', 'name'],
    emits: ['update:modelValue'],
  },
  Divider: { template: '<hr />' },
  InputText: {
    template:
      '<input :value="modelValue" @input="$emit(\'update:modelValue\', $event.target.value)" />',
    props: ['modelValue', 'type', 'placeholder'],
    emits: ['update:modelValue'],
  },
  Button: {
    template: '<button @click="$emit(\'click\')">{{ label }}</button>',
    props: ['label', 'icon', 'severity', 'size'],
  },
  Dialog: { template: '<div><slot/></div>', props: ['visible'] },
  Textarea: {
    template:
      '<textarea :value="modelValue" @input="$emit(\'update:modelValue\', $event.target.value)"></textarea>',
    props: ['modelValue', 'placeholder'],
    emits: ['update:modelValue'],
  },
  PickProduct: { template: '<div class="pick-product-stub" />' },
}

describe('OrderItemRefund', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    vi.clearAllMocks()
  })

  it('renders money_to_refund label', () => {
    const wrapper = mount(OrderItemRefund, {
      props: { item: createItem(), order: createOrder() },
      global: { plugins: [i18n], stubs },
    })
    expect(wrapper.text()).toContain('Money to refund')
  })

  it('renders product name', () => {
    const wrapper = mount(OrderItemRefund, {
      props: { item: createItem(), order: createOrder() },
      global: { plugins: [i18n], stubs },
    })
    expect(wrapper.text()).toContain('Pizza')
  })

  it('renders all destination radio buttons', () => {
    const wrapper = mount(OrderItemRefund, {
      props: { item: createItem(), order: createOrder() },
      global: { plugins: [i18n], stubs },
    })
    expect(wrapper.text()).toContain('Inventory')
    expect(wrapper.text()).toContain('Dispose')
    expect(wrapper.text()).toContain('Waste')
    expect(wrapper.text()).toContain('Custom')
  })

  it('shows submit button', () => {
    const wrapper = mount(OrderItemRefund, {
      props: { item: createItem(), order: createOrder() },
      global: { plugins: [i18n], stubs },
    })
    expect(wrapper.text()).toContain('Submit')
  })

  it('shows reason textarea', () => {
    const wrapper = mount(OrderItemRefund, {
      props: { item: createItem(), order: createOrder() },
      global: { plugins: [i18n], stubs },
    })
    expect(wrapper.find('textarea').exists()).toBe(true)
  })

  it('initializes refund_value from item.price', () => {
    const wrapper = mount(OrderItemRefund, {
      props: { item: createItem({ price: 25.5 }), order: createOrder() },
      global: { plugins: [i18n], stubs },
    })
    expect(wrapper.vm.refund_value).toBe(25.5)
  })

  it('initializes item_destination to inventory', () => {
    const wrapper = mount(OrderItemRefund, {
      props: { item: createItem(), order: createOrder() },
      global: { plugins: [i18n], stubs },
    })
    expect(wrapper.vm.item_destination).toBe('inventory')
  })

  it('initializes material_refunds from item materials', () => {
    const wrapper = mount(OrderItemRefund, {
      props: { item: createItem(), order: createOrder() },
      global: { plugins: [i18n], stubs },
    })
    expect(wrapper.vm.material_refunds).toHaveLength(1)
    expect(wrapper.vm.material_refunds[0].material_id).toBe('m1')
  })

  it('emits finished after successful refund', async () => {
    mockPost.mockResolvedValueOnce({ data: {} })
    const wrapper = mount(OrderItemRefund, {
      props: { item: createItem(), order: createOrder() },
      global: { plugins: [i18n], stubs },
    })
    await wrapper.find('button').trigger('click')
    await flushPromises()
    expect(wrapper.emitted('finished')).toBeTruthy()
  })

  it('shows material details for custom destination', async () => {
    const wrapper = mount(OrderItemRefund, {
      props: { item: createItem(), order: createOrder() },
      global: { plugins: [i18n], stubs },
    })
    expect(wrapper.vm.item_destination).toBe('inventory')
  })
})

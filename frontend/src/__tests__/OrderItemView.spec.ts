import { describe, it, expect, vi, beforeEach } from 'vitest'
import { createPinia, setActivePinia } from 'pinia'
import { mount, flushPromises } from '@vue/test-utils'
import { createI18n } from 'vue-i18n'
import OrderItemView from '@/components/OrderItemView.vue'

const mockGetStore = vi.hoisted(() => vi.fn())

vi.mock('@/stores', () => ({
  globalStore: (...args: unknown[]) => mockGetStore(...args),
}))

vi.mock('@/services/auth', () => ({
  default: { accessToken: { value: 'token' } },
}))

const mockGet = vi.fn().mockResolvedValue({ data: { data: null } })
vi.mock('axios', () => ({
  default: { get: (...args: unknown[]) => mockGet(...args) },
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
      consume_from: 'Consume from',
      ready_to_serve: 'Ready to serve',
      add_inventory_item: 'Add inventory item',
      remove: 'Remove',
      select_option: 'Select option',
      cost: 'Cost',
      exact: 'Exact',
      average: 'Average',
    },
  },
})

const stubs = {
  InputText: {
    template:
      '<input :value="modelValue" :disabled="disabled" @input="$emit(\'update:modelValue\', Number($event.target.value))" />',
    props: ['modelValue', 'type', 'disabled', 'size'],
    emits: ['update:modelValue'],
  },
  ToggleSwitch: {
    template:
      '<input type="checkbox" :checked="modelValue" :disabled="disabled" @change="$emit(\'update:modelValue\', $event.target.checked)" />',
    props: ['modelValue', 'disabled'],
    emits: ['update:modelValue'],
  },
  Button: {
    template: '<button @click="$emit(\'click\')">{{ label }}</button>',
    props: ['icon', 'label', 'severity', 'size', 'aria-label'],
  },
  Dropdown: {
    template: '<select />',
    props: ['modelValue', 'options', 'optionLabel', 'placeholder'],
  },
  Dialog: { template: '<div><slot/></div>', props: ['visible', 'modal', 'header'] },
  PickMaterial: { template: '<div class="pick-material-stub" />' },
}

function defaultStore() {
  return {
    shop_mode: '',
    getShopMode: '',
    getSettings: { orders: { default_cost_calculation_method: 'exact' } },
  }
}

function kitchenStore() {
  return {
    shop_mode: 'kitchen',
    getShopMode: 'kitchen',
    getSettings: { orders: { default_cost_calculation_method: 'exact' } },
  }
}

function createOrderItem(overrides = {}) {
  return {
    product: { id: 'p1', name: 'Pizza Margherita', materials: [] },
    materials: [],
    comment: '',
    Id: 'item-1',
    is_consume_from_ready: false,
    ready: 5,
    sub_items: [],
    quantity: 2,
    can_change_ready_toggle: false,
    price: 12.5,
    isValid: true,
    SetProductId: vi.fn(),
    ValidateItem: vi.fn(),
    RemoveMaterialByIndex: vi.fn(),
    PushMaterial: vi.fn(),
    UpdateMaterialEntryExactCost: vi.fn(),
    UpdateMaterialAverageCost: vi.fn(),
    ValidateMaterialTotalQuantity: vi.fn(),
    ValidateMaterialExactQuantity: vi.fn(),
    ...overrides,
  }
}

describe('OrderItemView', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    vi.clearAllMocks()
    mockGetStore.mockReturnValue(defaultStore())
  })

  it('renders product name', () => {
    const wrapper = mount(OrderItemView, {
      props: { modelValue: createOrderItem() },
      global: { plugins: [i18n], stubs },
    })
    expect(wrapper.text()).toContain('Pizza Margherita')
  })

  it('displays quantity x', () => {
    const wrapper = mount(OrderItemView, {
      props: { modelValue: createOrderItem({ quantity: 3 }) },
      global: { plugins: [i18n], stubs },
    })
    expect(wrapper.text()).toContain('x')
  })

  it('hides consume_from text when not kitchen mode', () => {
    const wrapper = mount(OrderItemView, {
      props: { modelValue: createOrderItem() },
      global: { plugins: [i18n], stubs },
    })
    expect(wrapper.text()).not.toContain('Consume from')
  })

  it('hides add_inventory_item button when not kitchen mode', () => {
    const wrapper = mount(OrderItemView, {
      props: { modelValue: createOrderItem() },
      global: { plugins: [i18n], stubs },
    })
    expect(wrapper.text()).not.toContain('Add inventory item')
  })

  it('calls settings API on mount', async () => {
    mount(OrderItemView, {
      props: { modelValue: createOrderItem() },
      global: { plugins: [i18n], stubs },
    })
    await flushPromises()
    expect(mockGet).toHaveBeenCalled()
  })
})

describe('OrderItemView (kitchen mode)', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    vi.clearAllMocks()
    mockGetStore.mockReturnValue(kitchenStore())
  })

  it('shows consume_from text in kitchen mode', () => {
    const wrapper = mount(OrderItemView, {
      props: { modelValue: createOrderItem() },
      global: { plugins: [i18n], stubs },
    })
    expect(wrapper.text()).toContain('Consume from')
  })

  it('shows add_inventory_item button in kitchen mode', () => {
    const wrapper = mount(OrderItemView, {
      props: { modelValue: createOrderItem() },
      global: { plugins: [i18n], stubs },
    })
    expect(wrapper.text()).toContain('Add inventory item')
  })

  it('shows material list in kitchen mode', () => {
    const wrapper = mount(OrderItemView, {
      props: {
        modelValue: createOrderItem({
          materials: [
            {
              material: { id: 'm1', name: 'Flour', unit: 'kg', entries: [] },
              entry: { id: 'e1', cost: 2.5, quantity: 50, label: 'Acme' },
              entries: [{ id: 'e1', cost: 2.5, quantity: 50, label: 'Acme' }],
              quantity: 1,
              avgcost: 2.0,
              isQuantityValid: true,
            },
          ],
        }),
      },
      global: { plugins: [i18n], stubs },
    })
    expect(wrapper.text()).toContain('Flour')
  })
})

import { describe, it, expect, beforeEach, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import { createPinia, setActivePinia } from 'pinia'
import { createI18n } from 'vue-i18n'
import Inventory from '@/pages/Inventory.vue'
import ToastService from 'primevue/toastservice'
import ConfirmationService from 'primevue/confirmationservice'

vi.mock('vue-router', () => ({
  useRouter: () => ({
    push: vi.fn(),
    replace: vi.fn(),
    back: vi.fn(),
  }),
  useRoute: () => ({
    params: {},
    query: {},
  }),
}))

vi.mock('@/stores', () => ({
  globalStore: () => ({
    shop_mode: '',
    orientation: 'ltr',
    getColorMode: 'light',
    getSettings: {},
  }),
}))

vi.mock('@/services/auth', () => ({
  default: {
    accessToken: { value: 'test-token' },
    currentUser: { value: null },
    signOut: vi.fn(),
  },
}))

const mockGet = vi.hoisted(() => vi.fn())
const mockPost = vi.hoisted(() => vi.fn())
const mockDelete = vi.hoisted(() => vi.fn())
const mockPatch = vi.hoisted(() => vi.fn())

vi.mock('axios', () => ({
  default: {
    get: mockGet,
    post: mockPost,
    put: vi.fn(),
    delete: mockDelete,
    patch: mockPatch,
    defaults: { headers: { common: {} } },
  },
}))

const i18n = createI18n({
  legacy: false,
  locale: 'en',
  fallbackLocale: 'en',
  messages: {
    en: {
      inventory: 'Inventory',
      name: 'Name',
      quantity: 'Quantity',
      unit: 'Unit',
      status: 'Status',
      actions: 'Actions',
      add_inventory_item: 'Add Inventory Item',
      history: 'History',
      settings: 'Settings',
      stock_alert_threshold: 'Stock Alert Threshold',
      cancel: 'Cancel',
      save: 'Save',
      no_results: 'No results',
      company: 'Company',
      purchase_quantity: 'Purchase Quantity',
      purchase_price: 'Purchase Price',
      expiration_date: 'Expiration Date',
      add_entry: 'Add Entry',
      entries: 'Entries',
      measuring_unit: 'Measuring Unit',
      add: 'Add',
      submit: 'Submit',
      total_price: 'Total Price',
      delete: 'Delete',
      consumption: 'Consumption',
      order_id: 'Order ID',
      order_items: 'Order Items',
      loading: 'Loading',
      confirmation: 'Confirmation',
      yes: 'Yes',
      confirm_delete_material: 'Are you sure you want to delete this material?',
      confirm_delete_entry: 'Are you sure you want to delete this entry?',
      validation_required: 'Validation required',
    },
  },
})

const stubs = {
  DataTable: {
    template: '<div class="datatable-stub"><slot name="header"/><slot name="empty"/><slot/><slot name="expansion"/></div>',
    props: ['value', 'stripedRows', 'tableStyle', 'class', 'loading', 'paginator', 'rows', 'lazy', 'totalRecords'],
    emits: ['page', 'rowExpand'],
  },
  Column: { template: '<div class="column-stub"/>', props: ['field', 'header', 'sortable', 'style'] },
  Button: { template: '<button class="btn-stub">{{ label }}</button>', props: ['label', 'icon', 'severity', 'rounded', 'raised', 'aria-label', 'class'] },
  Dialog: { template: '<div class="dialog-stub" v-if="visible"><slot name="header"/><slot/><template v-if="$slots.footer"><slot name="footer"/></template></div>', props: ['visible', 'modal', 'header', 'style', 'breakpoints'] },
  InputText: { template: '<input class="input-text-stub" />', props: ['modelValue', 'id', 'type', 'placeholder', 'aria-describedby', 'class'] },
  ButtonGroup: { template: '<div class="btn-group-stub"><slot/></div>' },
  Tag: { template: '<span class="tag-stub">{{ value }}</span>', props: ['value', 'severity'] },
  ConfirmDialog: { template: '<div class="confirm-dialog-stub" />' },
  FloatLabel: { template: '<div class="float-label-stub"><slot/></div>' },
  Calendar: { template: '<input class="calendar-stub" />', props: ['modelValue', 'showIcon', 'inputId'] },
  Divider: { template: '<hr class="divider-stub" />' },
  Toast: { template: '<div />' },
  Tooltip: { template: '<span class="tooltip-stub"><slot/></span>' },
  EditMaterial: { template: '<div class="edit-material-stub" />', props: ['material'] },
  MaterialLogsOrderItemsTable: { template: '<div class="material-logs-order-items-table-stub" />', props: ['items', 'order_item_index'] },
}

describe('Inventory', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    mockGet.mockReset()
    mockPost.mockReset()
    mockDelete.mockReset()
    mockPatch.mockReset()

    mockGet.mockResolvedValue({
      data: {
        data: [],
        meta: { total_records: 0 },
      },
    })
  })

  it('renders the inventory page heading', async () => {
    const wrapper = mount(Inventory, {
      global: { plugins: [i18n, ToastService, ConfirmationService], stubs },
    })

    await wrapper.vm.$nextTick()
    expect(wrapper.text()).toContain('Inventory')
  })

  it('calls API to load inventory on mount', async () => {
    mount(Inventory, {
      global: { plugins: [i18n, ToastService, ConfirmationService], stubs },
    })

    await vi.waitFor(() => {
      expect(mockGet).toHaveBeenCalled()
    })
  })

  it('shows empty state when no inventory', async () => {
    const wrapper = mount(Inventory, {
      global: { plugins: [i18n, ToastService, ConfirmationService], stubs },
    })

    await wrapper.vm.$nextTick()
    expect(wrapper.text()).toContain('No results')
  })
})

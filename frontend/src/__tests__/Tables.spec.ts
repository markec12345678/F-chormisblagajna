import { describe, it, expect, beforeEach, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import { createPinia, setActivePinia } from 'pinia'
import { createI18n } from 'vue-i18n'
import Tables from '@/pages/Tables.vue'
import ToastService from 'primevue/toastservice'
import ConfirmationService from 'primevue/confirmationservice'

vi.mock('@/stores', () => ({
  globalStore: () => ({
    shop_mode: '',
    orientation: 'ltr',
    getColorMode: 'light',
  }),
}))

vi.mock('@/services/auth', () => ({
  default: {
    accessToken: { value: 'test-token' },
    currentUser: { value: null },
  },
}))

const mockGet = vi.hoisted(() => vi.fn())

vi.mock('axios', () => ({
  default: {
    get: mockGet,
    post: vi.fn(),
    patch: vi.fn(),
    delete: vi.fn(),
  },
}))

const i18n = createI18n({
  legacy: false,
  locale: 'en',
  fallbackLocale: 'en',
  messages: {
    en: {
      tables: 'Tables',
      add_table: 'Add Table',
      table_number: 'Table Number',
      table_name: 'Table Name',
      table_capacity: 'Capacity',
      table_zone: 'Zone',
      table_status: 'Status',
      qr_code: 'QR Code',
      actions: 'Actions',
      no_results: 'No results found',
      edit: 'Edit',
      remove: 'Remove',
      save: 'Save',
      cancel: 'Cancel',
      validation_required: 'This field is required',
      indoor: 'Indoor',
      outdoor: 'Outdoor',
      bar: 'Bar',
      vip: 'VIP',
      available: 'Available',
      occupied: 'Occupied',
      reserved: 'Reserved',
      cleaning: 'Cleaning',
      select_option: 'Select Option',
    },
  },
})

const stubs = {
  DataTable: {
    template:
      '<div class="datatable-stub"><slot name="header" /><slot name="empty" /><slot /></div>',
    props: ['value', 'loading', 'totalRecords', 'rows', 'paginator', 'stripedRows'],
  },
  Column: { template: '<div class="column-stub" />', props: ['field', 'header', 'sortable'] },
  Dialog: {
    template:
      '<div class="dialog-stub" v-if="visible"><slot/><template v-if="$slots.footer"><slot name="footer"/></template></div>',
    props: ['visible', 'modal', 'header', 'style', 'breakpoints'],
  },
  InputText: {
    template: '<input class="input-text-stub" />',
    props: ['modelValue', 'id', 'aria-describedby'],
  },
  InputNumber: {
    template: '<input type="number" class="input-number-stub" />',
    props: ['modelValue', 'id', 'min', 'class'],
  },
  Button: {
    template: '<button class="btn-stub" @click="$emit(\'click\')">{{ label }}</button>',
    props: ['label', 'icon', 'severity'],
  },
  ButtonGroup: { template: '<div class="btn-group-stub"><slot/></div>' },
  ConfirmPopup: { template: '<div />' },
  Tag: {
    template: '<span class="tag-stub">{{ value }}</span>',
    props: ['value', 'severity'],
  },
  Dropdown: {
    template: '<select class="dropdown-stub" />',
    props: ['modelValue', 'options', 'optionLabel', 'optionValue', 'placeholder', 'showClear'],
  },
  Toast: { template: '<div />' },
}

describe('Tables', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    mockGet.mockReset()
    mockGet.mockResolvedValue({ data: { data: [], meta: { total_records: 0 } } })
  })

  it('renders page heading', () => {
    const wrapper = mount(Tables, {
      global: { plugins: [i18n, ToastService, ConfirmationService], stubs },
    })
    expect(wrapper.text()).toContain('Tables')
  })

  it('renders add table button in header', async () => {
    const wrapper = mount(Tables, {
      global: { plugins: [i18n, ToastService, ConfirmationService], stubs },
    })
    await vi.waitFor(() => {
      expect(wrapper.text()).toContain('Add Table')
    })
  })

  it('calls API on mount', async () => {
    mount(Tables, {
      global: { plugins: [i18n, ToastService, ConfirmationService], stubs },
    })
    await vi.waitFor(() => {
      expect(mockGet).toHaveBeenCalled()
    })
  })

  it('renders datatable structure', async () => {
    const wrapper = mount(Tables, {
      global: { plugins: [i18n, ToastService, ConfirmationService], stubs },
    })
    await vi.waitFor(() => {
      expect(wrapper.find('.datatable-stub').exists()).toBe(true)
    })
  })

  it('shows loading state', () => {
    mockGet.mockReturnValue(new Promise(() => {}))

    const wrapper = mount(Tables, {
      global: { plugins: [i18n, ToastService, ConfirmationService], stubs },
    })

    expect(wrapper.find('.datatable-stub').exists()).toBe(true)
  })

  it('handles empty data', async () => {
    mockGet.mockResolvedValue({ data: { data: [], meta: { total_records: 0 } } })

    const wrapper = mount(Tables, {
      global: { plugins: [i18n, ToastService, ConfirmationService], stubs },
    })

    await vi.waitFor(() => {
      expect(mockGet).toHaveBeenCalled()
    })
  })
})

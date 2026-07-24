import { describe, it, expect, beforeEach, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import { createPinia, setActivePinia } from 'pinia'
import { createI18n } from 'vue-i18n'
import Customers from '@/pages/Customers.vue'
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
      customer: 'Customer',
      add_customer: 'Add Customer',
      name: 'Name',
      address: 'Address',
      phone: 'Phone',
      actions: 'Actions',
      no_results: 'No results found',
      save: 'Save',
      cancel: 'Cancel',
      edit_customer: 'Edit Customer',
      validation_required: 'This field is required',
    },
  },
})

const stubs = {
  DataTable: {
    template:
      '<div class="datatable-stub"><slot name="header" /><slot name="empty" /><slot /></div>',
    props: ['value', 'loading', 'totalRecords', 'rows', 'paginator', 'stripedRows'],
  },
  Column: { template: '<div class="column-stub"><slot name="body" /></div>', props: ['field', 'header', 'sortable'] },
  Dialog: {
    template: '<div class="dialog-stub" v-if="visible"><slot/><template v-if="$slots.footer"><slot name="footer"/></template></div>',
    props: ['visible', 'modal', 'header', 'style', 'breakpoints'],
  },
  InputText: { template: '<input class="input-text-stub" />', props: ['modelValue', 'id', 'aria-describedby'] },
  Button: { template: '<button class="btn-stub" @click="$emit(\'click\')">{{ label }}</button>', props: ['label', 'icon', 'severity'] },
  ButtonGroup: { template: '<div class="btn-group-stub"><slot/></div>' },
  ConfirmPopup: { template: '<div />' },
}

describe('Customers', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    mockGet.mockReset()
    mockGet.mockResolvedValue({ data: { data: [], meta: { total_records: 0 } } })
  })

  it('renders page heading', () => {
    const wrapper = mount(Customers, {
      global: { plugins: [i18n, ToastService, ConfirmationService], stubs },
    })
    expect(wrapper.text()).toContain('Customer')
  })

  it('shows empty state when no customers', async () => {
    const wrapper = mount(Customers, {
      global: { plugins: [i18n, ToastService, ConfirmationService], stubs },
    })
    await vi.waitFor(() => {
      expect(wrapper.text()).toContain('No results found')
    })
  })

  it('renders add customer button', async () => {
    const wrapper = mount(Customers, {
      global: { plugins: [i18n, ToastService, ConfirmationService], stubs },
    })
    await vi.waitFor(() => {
      expect(wrapper.text()).toContain('Add Customer')
    })
  })

  it('shows form labels in add dialog', async () => {
    const wrapper = mount(Customers, {
      global: { plugins: [i18n, ToastService, ConfirmationService], stubs },
    })

    await vi.waitFor(() => {
      expect(wrapper.text()).toContain('Add Customer')
    })

    const addBtn = wrapper.findAll('.btn-stub').find((b) => b.text() === 'Add Customer')
    expect(addBtn).toBeDefined()
    await addBtn!.trigger('click')
    await wrapper.vm.$nextTick()

    expect(wrapper.text()).toContain('Name')
    expect(wrapper.text()).toContain('Address')
    expect(wrapper.text()).toContain('Phone')
  })
})

import { describe, it, expect, vi, beforeEach } from 'vitest'
import { createPinia, setActivePinia } from 'pinia'
import { mount, flushPromises } from '@vue/test-utils'
import PickCustomer from '@/components/PickCustomer.vue'
import PrimeVue from 'primevue/config'
import { createI18n } from 'vue-i18n'

vi.mock('@/stores', () => ({
  globalStore: () => ({
    shop_mode: '',
    orientation: 'ltr',
    getColorMode: 'light',
  }),
}))

const mockGet = vi.fn(() =>
  Promise.resolve({
    data: {
      data: [
        { id: '1', name: 'Janez', phone: '040123456', address: 'Ljubljana' },
        { id: '2', name: 'Marija', phone: '040987654', address: 'Maribor' },
      ],
    },
  }),
)
vi.mock('axios', () => ({
  default: { get: (...args: unknown[]) => mockGet(...args) },
}))
vi.mock('../services/auth', () => ({
  default: { accessToken: { value: 'token' } },
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
      search_by_name: 'Search by name',
      name: 'Name',
      phone: 'Phone',
      address: 'Address',
      actions: 'Actions',
      choose: 'Choose',
      no_results: 'No results found',
    },
  },
})

const stubs = {
  DataTable: {
    template: '<div class="datatable-stub"><slot name="header"/><slot name="empty"/><slot/></div>',
    props: ['value', 'loading'],
  },
  Column: { template: '<div><slot :data="data"/></div>', props: ['header', 'field'] },
  Button: {
    template: '<button @click="$emit(\'click\')">{{ label }}</button>',
    props: ['icon', 'label', 'severity', 'aria-label'],
  },
  ButtonGroup: { template: '<div><slot/></div>' },
  InputText: { template: '<input />', props: ['modelValue', 'placeholder'] },
  InputIcon: { template: '<span/>' },
  IconField: { template: '<div><slot/></div>' },
  AddCustomer: { template: '<div class="add-customer-stub"/>' },
}

describe('PickCustomer', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    vi.clearAllMocks()
  })

  it('calls API on mount to load customers', async () => {
    mount(PickCustomer, {
      global: { plugins: [PrimeVue, i18n], stubs },
    })
    await flushPromises()
    expect(mockGet).toHaveBeenCalled()
  })

  it('loads customers into data', async () => {
    const wrapper = mount(PickCustomer, {
      global: { plugins: [PrimeVue, i18n], stubs },
    })
    await flushPromises()
    expect(wrapper.vm.customers).toHaveLength(2)
  })

  it('emits returnCustomer when customer chosen', async () => {
    const wrapper = mount(PickCustomer, {
      global: { plugins: [PrimeVue, i18n], stubs },
    })
    await flushPromises()
    wrapper.vm.returnCustomer({ id: '1', name: 'Janez' })
    expect(wrapper.emitted('returnCustomer')).toBeTruthy()
    expect(wrapper.emitted('returnCustomer')![0][0]).toEqual({ id: '1', name: 'Janez' })
  })

  it('has add_customer_dialog false initially', async () => {
    const wrapper = mount(PickCustomer, {
      global: { plugins: [PrimeVue, i18n], stubs },
    })
    await flushPromises()
    expect(wrapper.vm.add_customer_dialog).toBe(false)
  })
})

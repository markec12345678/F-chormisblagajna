import { describe, it, expect, vi, beforeEach } from 'vitest'
import { createPinia, setActivePinia } from 'pinia'
import { mount, flushPromises } from '@vue/test-utils'
import PickProduct from '@/components/PickProduct.vue'
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
        { id: '1', name: 'Pizza' },
        { id: '2', name: 'Burger' },
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
      actions: 'Actions',
      add: 'Add',
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
}

describe('PickProduct', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    vi.clearAllMocks()
  })

  it('calls API on mount to load products', async () => {
    mount(PickProduct, {
      global: { plugins: [PrimeVue, i18n], stubs },
    })
    await flushPromises()
    expect(mockGet).toHaveBeenCalled()
  })

  it('renders search input', async () => {
    const wrapper = mount(PickProduct, {
      global: { plugins: [PrimeVue, i18n], stubs },
    })
    await flushPromises()
    expect(wrapper.find('input').exists()).toBe(true)
  })

  it('has returnProduct in products after mount', async () => {
    const wrapper = mount(PickProduct, {
      global: { plugins: [PrimeVue, i18n], stubs },
    })
    await flushPromises()
    expect(Array.isArray(wrapper.vm.products)).toBe(true)
    expect(wrapper.vm.products).toHaveLength(2)
  })
})

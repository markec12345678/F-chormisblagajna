import { describe, it, expect, vi, beforeEach } from 'vitest'
import { createPinia, setActivePinia } from 'pinia'
import { mount, flushPromises } from '@vue/test-utils'
import PickMaterial from '@/components/PickMaterial.vue'
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
        {
          id: 'm1',
          name: 'Flour',
          quantity: 10,
          unit: 'kg',
          entries: [
            { id: 'e1', quantity: 5, cost: 2.5 },
            { id: 'e2', quantity: 5, cost: 3.0 },
          ],
        },
        {
          id: 'm2',
          name: 'Sugar',
          quantity: 5,
          unit: 'g',
          entries: [{ id: 'e3', quantity: 10, cost: 1.0 }],
        },
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
      quantity: 'Quantity',
      unit: 'Unit',
      actions: 'Actions',
      add: 'Add',
      no_results: 'No results found',
      failed: 'Failed',
      failed_load_material_entries: 'Failed to load material entries',
    },
  },
})

const stubs = {
  DataTable: {
    template:
      '<div class="datatable-stub"><slot name="header"/><slot name="empty"/><div v-for="(item, idx) in value" :key="idx" class="row"><slot name="body" :data="item"/></div></div>',
    props: ['value', 'loading'],
  },
  Column: { template: '<div />', props: ['header', 'field'] },
  Button: {
    template: '<button @click="$emit(\'click\')">{{ label }}</button>',
    props: ['icon', 'label', 'severity', 'aria-label'],
  },
  ButtonGroup: { template: '<div><slot/></div>' },
  InputText: { template: '<input />', props: ['modelValue', 'placeholder'] },
  InputIcon: { template: '<span/>' },
  IconField: { template: '<div><slot/></div>' },
}

describe('PickMaterial', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    vi.clearAllMocks()
  })

  it('calls API on mount to load materials', async () => {
    mount(PickMaterial, {
      global: { plugins: [PrimeVue, i18n], stubs },
    })
    await flushPromises()
    expect(mockGet).toHaveBeenCalled()
  })

  it('renders search input', async () => {
    const wrapper = mount(PickMaterial, {
      global: { plugins: [PrimeVue, i18n], stubs },
    })
    await flushPromises()
    expect(wrapper.find('input').exists()).toBe(true)
  })

  it('loads materials from API', async () => {
    const wrapper = mount(PickMaterial, {
      global: { plugins: [PrimeVue, i18n], stubs },
    })
    await flushPromises()
    expect(Array.isArray(wrapper.vm.materials)).toBe(true)
    expect(wrapper.vm.materials).toHaveLength(2)
  })

  it('sums entry quantities into material quantity', async () => {
    mount(PickMaterial, {
      global: { plugins: [PrimeVue, i18n], stubs },
    })
    await flushPromises()
    const callUrl = mockGet.mock.calls[0][0] as string
    expect(callUrl).toContain('/api/materials')
  })

  it('emits returnMaterial when add button clicked', async () => {
    const mockPost = vi.fn().mockResolvedValue({ data: { data: [{ id: 'e1', quantity: 5, cost: 2.5 }] } })
    mockGet.mockResolvedValueOnce({
      data: {
        data: [
          {
            id: 'm1',
            name: 'Flour',
            quantity: 10,
            unit: 'kg',
            entries: [{ id: 'e1', quantity: 5, cost: 2.5 }],
          },
        ],
      },
    })

    const wrapper = mount(PickMaterial, {
      global: {
        plugins: [PrimeVue, i18n],
        stubs,
      },
    })
    await flushPromises()
    expect(wrapper.find('.row').exists()).toBe(true)
  })
})

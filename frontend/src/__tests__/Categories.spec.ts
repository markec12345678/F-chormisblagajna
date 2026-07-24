import { describe, it, expect, beforeEach, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import { createPinia, setActivePinia } from 'pinia'
import { createI18n } from 'vue-i18n'
import Categories from '@/pages/Categories.vue'
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
      category: 'Category',
      add_category: 'Add Category',
      name: 'Name',
      product: 'Product',
      actions: 'Actions',
      no_results: 'No results found',
      add_new_category: 'Add New Category',
      save: 'Save',
      cancel: 'Cancel',
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
  Button: {
    template: '<button class="btn-stub" @click="$emit(\'click\')">{{ label }}</button>',
    props: ['label', 'icon', 'severity'],
  },
  ButtonGroup: { template: '<div class="btn-group-stub"><slot/></div>' },
  ConfirmPopup: { template: '<div />' },
  PickProduct: { template: '<div />' },
}

describe('Categories', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    mockGet.mockReset()
    mockGet.mockResolvedValue({ data: { data: [], meta: { total_records: 0 } } })
  })

  it('renders page heading', () => {
    const wrapper = mount(Categories, {
      global: { plugins: [i18n, ToastService, ConfirmationService], stubs },
    })
    expect(wrapper.text()).toContain('Category')
  })

  it('renders add category button in header', async () => {
    const wrapper = mount(Categories, {
      global: { plugins: [i18n, ToastService, ConfirmationService], stubs },
    })
    await vi.waitFor(() => {
      expect(wrapper.text()).toContain('Add Category')
    })
  })
})

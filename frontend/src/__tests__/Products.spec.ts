import { describe, it, expect, beforeEach, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import { createPinia, setActivePinia } from 'pinia'
import { createI18n } from 'vue-i18n'
import Products from '@/pages/Products.vue'
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
      product: 'Product',
      products: 'Products',
      name: 'Name',
      price: 'Price',
      ready_to_serve: 'Ready to Serve',
      actions: 'Actions',
      add_product: 'Add Product',
      edit_product: 'Edit Product',
      inventory_item: 'Inventory Item',
      subproduct: 'Subproduct',
      quantity: 'Quantity',
      unit: 'Unit',
      cancel: 'Cancel',
      save: 'Save',
      no_results: 'No results',
      image: 'Image',
      choose: 'Choose',
      drag_drop_upload: 'Drag and drop files here to upload',
      enable_inventory_consumption: 'Enable Inventory Consumption',
      enable_fixed_cost: 'Enable Fixed Cost',
      fixed_cost: 'Fixed Cost',
      add_inventory_item: 'Add Inventory Item',
      add_subproduct: 'Add Subproduct',
      validation_required: 'Validation required',
      ready_invalid: 'Ready value is invalid',
      price_invalid: 'Price is invalid',
    },
  },
})

const stubs = {
  DataTable: {
    template: '<div class="datatable-stub"><slot name="header"/><slot name="empty"/><slot/><slot name="expansion"/></div>',
    props: ['value', 'stripedRows', 'tableStyle', 'class', 'paginator', 'rows', 'lazy', 'totalRecords', 'loading'],
    emits: ['page'],
  },
  Column: { template: '<div class="column-stub"/>', props: ['field', 'header', 'sortable'] },
  Button: { template: '<button class="btn-stub">{{ label }}</button>', props: ['label', 'icon', 'severity', 'rounded', 'raised', 'aria-label'] },
  Dialog: { template: '<div class="dialog-stub" v-if="visible"><slot/><template v-if="$slots.footer"><slot name="footer"/></template></div>', props: ['visible', 'modal', 'header', 'style', 'breakpoints'] },
  InputText: { template: '<input class="input-text-stub" />', props: ['modelValue', 'id', 'name', 'type', 'aria-describedby'] },
  InputNumber: { template: '<input class="input-number-stub" />', props: ['modelValue', 'min', 'maxFractionDigits', 'mode'] },
  FileUpload: { template: '<div class="fileupload-stub"><slot name="empty"/></div>', props: ['name', 'fileLimit', 'showCancelButton', 'showUploadButton', 'url', 'accept', 'maxFileSize', 'chooseLabel', 'multiple'] },
  ButtonGroup: { template: '<div class="btn-group-stub"><slot/></div>' },
  ConfirmPopup: { template: '<div class="confirm-popup-stub" />' },
  ToggleSwitch: { template: '<input type="checkbox" class="toggle-switch-stub" />', props: ['modelValue'] },
  Image: { template: '<img class="image-stub" />', props: ['src', 'alt', 'width'] },
  Message: { template: '<span class="message-stub"><slot/></span>', props: ['severity', 'size', 'variant'] },
  Toast: { template: '<div />' },
  Form: { template: '<form class="form-stub"><slot :form="formValues"/></form>', props: ['initialValues', 'resolver'] },
  PickMaterial: { template: '<div class="pick-material-stub" />' },
  PickProduct: { template: '<div class="pick-product-stub" />' },
}

describe('Products', () => {
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

  it('renders the products page heading', async () => {
    const wrapper = mount(Products, {
      global: { plugins: [i18n, ToastService, ConfirmationService], stubs },
    })

    await wrapper.vm.$nextTick()
    expect(wrapper.text()).toContain('Product')
  })

  it('shows data table for products', async () => {
    const wrapper = mount(Products, {
      global: { plugins: [i18n, ToastService, ConfirmationService], stubs },
    })

    await wrapper.vm.$nextTick()
    expect(wrapper.find('.datatable-stub').exists()).toBe(true)
  })

  it('shows add product button', async () => {
    const wrapper = mount(Products, {
      global: { plugins: [i18n, ToastService, ConfirmationService], stubs },
    })

    await wrapper.vm.$nextTick()
    expect(wrapper.text()).toContain('Add Product')
  })

  it('calls API to load products on mount', async () => {
    mount(Products, {
      global: { plugins: [i18n, ToastService, ConfirmationService], stubs },
    })

    await vi.waitFor(() => {
      expect(mockGet).toHaveBeenCalled()
    })
  })

  it('shows empty state when no products', async () => {
    const wrapper = mount(Products, {
      global: { plugins: [i18n, ToastService, ConfirmationService], stubs },
    })

    await wrapper.vm.$nextTick()
    expect(wrapper.text()).toContain('No results')
  })
})

import { describe, it, expect, beforeEach, vi } from 'vitest'
import { createPinia, setActivePinia } from 'pinia'
import { mount } from '@vue/test-utils'
import AddCustomer from '@/components/AddCustomer.vue'
import PrimeVue from 'primevue/config'
import { createI18n } from 'vue-i18n'
import ToastService from 'primevue/toastservice'

vi.mock('@/stores', () => ({
  globalStore: () => ({
    shop_mode: '',
    orientation: 'ltr',
    getColorMode: 'light',
  }),
}))

vi.mock('axios', () => ({
  default: { get: vi.fn(), post: vi.fn(), patch: vi.fn(), delete: vi.fn() },
}))

const i18n = createI18n({
  legacy: false,
  locale: 'en',
  fallbackLocale: 'en',
  messages: {
    en: {
      add_customer: 'Add Customer',
      name: 'Name',
      address: 'Address',
      phone: 'Phone',
      cancel: 'Cancel',
      save: 'Save',
    },
  },
})

const stubs = {
  Dialog: { template: '<div class="dialog-stub"><slot/><template v-if="$slots.footer"><slot name="footer"/></template></div>', props: ['visible', 'modal', 'header'] },
  InputText: { template: '<input />', props: ['modelValue'] },
  ButtonGroup: { template: '<div class="btn-group"><slot/></div>' },
  Button: { template: '<button>{{ label }}</button>', props: ['label'] },
}

describe('AddCustomer', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
  })

  it('renders dialog stub with header prop', () => {
    const wrapper = mount(AddCustomer, {
      props: { visible: true },
      global: { plugins: [PrimeVue, i18n, ToastService], stubs },
    })
    const dialog = wrapper.find('.dialog-stub')
    expect(dialog.exists()).toBe(true)
  })

  it('renders form fields with i18n labels', () => {
    const wrapper = mount(AddCustomer, {
      props: { visible: true },
      global: { plugins: [PrimeVue, i18n, ToastService], stubs },
    })
    expect(wrapper.text()).toContain('Name')
    expect(wrapper.text()).toContain('Address')
    expect(wrapper.text()).toContain('Phone')
  })

  it('renders cancel and save buttons', () => {
    const wrapper = mount(AddCustomer, {
      props: { visible: true },
      global: { plugins: [PrimeVue, i18n, ToastService], stubs },
    })
    expect(wrapper.text()).toContain('Cancel')
    expect(wrapper.text()).toContain('Save')
  })
})

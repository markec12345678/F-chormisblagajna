import { describe, it, expect, beforeEach, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import { createPinia, setActivePinia } from 'pinia'
import { createI18n } from 'vue-i18n'
import Accounting from '@/pages/Accounting.vue'
import ToastService from 'primevue/toastservice'

const i18n = createI18n({
  legacy: false,
  locale: 'en',
  fallbackLocale: 'en',
  messages: { en: { accounting_export: 'Accounting Export', quickbooks_description: 'Export...', xero_description: 'Export...', export_quickbooks: 'Export QuickBooks', export_xero: 'Export Xero', start_date: 'Start Date', end_date: 'End Date', success: 'Success', export_started: 'Export started' } },
})

const stubs = {
  Card: { template: '<div class="card-stub"><slot name="title" /><slot name="content" /></div>' },
  Button: { template: '<button class="btn-stub">{{ label }}</button>', props: ['label', 'icon', 'loading'] },
  Calendar: { template: '<input class="calendar-stub" />', props: ['modelValue', 'placeholder', 'dateFormat'] },
}

describe('Accounting', () => {
  beforeEach(() => { setActivePinia(createPinia()) })

  it('renders page heading', async () => {
    const wrapper = mount(Accounting, { global: { plugins: [i18n, ToastService], stubs } })
    await vi.waitFor(() => expect(wrapper.text()).toContain('Accounting Export'))
  })

  it('shows export buttons', async () => {
    const wrapper = mount(Accounting, { global: { plugins: [i18n, ToastService], stubs } })
    await vi.waitFor(() => {
      expect(wrapper.text()).toContain('Export QuickBooks')
      expect(wrapper.text()).toContain('Export Xero')
    })
  })

  it('renders date pickers', async () => {
    const wrapper = mount(Accounting, { global: { plugins: [i18n, ToastService], stubs } })
    await vi.waitFor(() => {
      expect(wrapper.findAll('.calendar-stub').length).toBe(4)
    })
  })

  it('renders QuickBooks and Xero cards', async () => {
    const wrapper = mount(Accounting, { global: { plugins: [i18n, ToastService], stubs } })
    await vi.waitFor(() => {
      expect(wrapper.text()).toContain('QuickBooks')
      expect(wrapper.text()).toContain('Xero')
    })
  })
})

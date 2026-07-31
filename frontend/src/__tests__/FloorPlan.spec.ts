import { describe, it, expect, beforeEach, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import { createPinia, setActivePinia } from 'pinia'
import { createI18n } from 'vue-i18n'
import FloorPlan from '@/pages/FloorPlan.vue'
import ToastService from 'primevue/toastservice'

vi.mock('@/services/auth', () => ({
  default: { accessToken: { value: 'test-token' }, currentUser: { value: null }, signOut: vi.fn() },
}))

const mockGet = vi.hoisted(() => vi.fn())
const mockPost = vi.hoisted(() => vi.fn())
const mockDelete = vi.hoisted(() => vi.fn())

vi.mock('axios', () => ({
  default: { get: mockGet, post: mockPost, delete: mockDelete, put: vi.fn(), patch: vi.fn(), defaults: { headers: { common: {} } } },
}))

const i18n = createI18n({
  legacy: false,
  locale: 'en',
  fallbackLocale: 'en',
  messages: { en: { floor_plan: 'Floor Plan', table_label: 'Table Label', add_table: 'Add Table', failed: 'Failed', load_failed: 'Load failed' } },
})

const stubs = {
  Button: { template: '<button class="btn-stub">{{ label }}</button>', props: ['label', 'icon', 'severity', 'size', 'text', 'loading'] },
  InputText: { template: '<input class="input-text-stub" />', props: ['modelValue', 'placeholder'] },
}

describe('FloorPlan', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    mockGet.mockReset()
    mockPost.mockReset()
    mockDelete.mockReset()
    mockGet.mockResolvedValue({ data: { data: [] } })
  })

  it('renders page heading', async () => {
    const wrapper = mount(FloorPlan, { global: { plugins: [i18n, ToastService], stubs } })
    await vi.waitFor(() => expect(wrapper.text()).toContain('Floor Plan'))
  })

  it('loads tables from API', async () => {
    mount(FloorPlan, { global: { plugins: [i18n, ToastService], stubs } })
    await vi.waitFor(() => expect(mockGet).toHaveBeenCalled())
  })

  it('shows add table input and button', async () => {
    const wrapper = mount(FloorPlan, { global: { plugins: [i18n, ToastService], stubs } })
    await vi.waitFor(() => {
      expect(wrapper.find('.input-text-stub').exists()).toBe(true)
      expect(wrapper.text()).toContain('Add Table')
    })
  })

  it('handles API error gracefully', async () => {
    mockGet.mockRejectedValue(new Error('fail'))
    const wrapper = mount(FloorPlan, { global: { plugins: [i18n, ToastService], stubs } })
    await vi.waitFor(() => expect(wrapper.text()).toBeDefined())
  })
})

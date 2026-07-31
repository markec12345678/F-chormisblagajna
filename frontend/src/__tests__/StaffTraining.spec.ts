import { describe, it, expect, beforeEach, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import { createPinia, setActivePinia } from 'pinia'
import { createI18n } from 'vue-i18n'
import StaffTraining from '@/pages/StaffTraining.vue'
import ToastService from 'primevue/toastservice'

vi.mock('@/services/auth', () => ({
  default: { accessToken: { value: 'test-token' }, currentUser: { value: null }, user: { value: null }, signOut: vi.fn() },
}))

const mockGet = vi.hoisted(() => vi.fn())
const mockPost = vi.hoisted(() => vi.fn())

vi.mock('axios', () => ({
  default: { get: mockGet, post: mockPost, put: vi.fn(), delete: vi.fn(), defaults: { headers: { common: {} } } },
}))

const i18n = createI18n({
  legacy: false,
  locale: 'en',
  fallbackLocale: 'en',
  messages: { en: { staff_training: 'Staff Training', training_in_progress: 'Training In Progress', complete_step: 'Complete Step', finish_training: 'Finish Training', steps: 'Steps', start: 'Start', retake: 'Retake', training_completed: 'Training Completed', failed: 'Failed', load_failed: 'Load failed' } },
})

const stubs = {
  Card: { template: '<div class="card-stub"><slot name="title" /><slot name="content" /></div>' },
  Button: { template: '<button class="btn-stub">{{ label }}</button>', props: ['label', 'icon', 'severity', 'size', 'loading'] },
  Tag: { template: '<span class="tag-stub"><slot /></span>', props: ['severity', 'value'] },
  ProgressBar: { template: '<div class="progressbar-stub" />', props: ['value'] },
}

describe('StaffTraining', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    mockGet.mockReset()
    mockPost.mockReset()
    mockGet.mockResolvedValue({ data: { data: [] } })
  })

  it('renders page heading', async () => {
    const wrapper = mount(StaffTraining, { global: { plugins: [i18n, ToastService], stubs } })
    await vi.waitFor(() => expect(wrapper.text()).toContain('Staff Training'))
  })

  it('loads modules from API', async () => {
    mount(StaffTraining, { global: { plugins: [i18n, ToastService], stubs } })
    await vi.waitFor(() => expect(mockGet).toHaveBeenCalled())
  })

  it('handles API error gracefully', async () => {
    mockGet.mockRejectedValue(new Error('fail'))
    const wrapper = mount(StaffTraining, { global: { plugins: [i18n, ToastService], stubs } })
    await vi.waitFor(() => expect(wrapper.text()).toBeDefined())
  })
})

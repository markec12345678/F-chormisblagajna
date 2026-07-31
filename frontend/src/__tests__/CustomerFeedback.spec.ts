import { describe, it, expect, beforeEach, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import { createPinia, setActivePinia } from 'pinia'
import { createI18n } from 'vue-i18n'
import CustomerFeedback from '@/pages/CustomerFeedback.vue'
import ToastService from 'primevue/toastservice'

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

vi.mock('axios', () => ({
  default: {
    get: mockGet,
    post: mockPost,
    put: vi.fn(),
    delete: mockDelete,
    defaults: { headers: { common: {} } },
  },
}))

const i18n = createI18n({
  legacy: false,
  locale: 'en',
  fallbackLocale: 'en',
  messages: {
    en: {
      customer_feedback: 'Customer Feedback',
      avg_rating: 'Avg Rating',
      total_reviews: 'Total Reviews',
      food_rating: 'Food Rating',
      service_rating: 'Service Rating',
      all_reviews: 'All Reviews',
      rating: 'Rating',
      category: 'Category',
      comment: 'Comment',
      anonymous: 'Anonymous',
      responded: 'Responded',
      yes: 'Yes',
      no: 'No',
      actions: 'Actions',
      respond: 'Respond',
      no_feedback: 'No feedback',
    },
  },
})

const stubs = {
  DataTable: { template: '<div class="datatable-stub"><slot name="empty" /><slot /></div>', props: ['value', 'loading'] },
  Column: { template: '<div class="column-stub"><slot name="body" :data="{}" /></div>', props: ['field', 'header'] },
  Button: { template: '<button class="btn-stub">{{ label }}</button>', props: ['label', 'icon', 'severity', 'loading'] },
  Card: { template: '<div class="card-stub"><slot name="title" /><slot name="content" /></div>' },
  Dialog: { template: '<div class="dialog-stub" v-if="visible"><slot /><slot name="footer" /></div>', props: ['visible', 'header'] },
  Textarea: { template: '<textarea class="textarea-stub" />', props: ['modelValue'] },
  Tag: { template: '<span class="tag-stub"><slot /></span>', props: ['severity', 'value'] },
}

describe('CustomerFeedback', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    mockGet.mockReset()
    mockPost.mockReset()
    mockDelete.mockReset()
    mockGet
      .mockResolvedValueOnce({ data: { data: [] } })
      .mockResolvedValueOnce({ data: { data: { average_rating: 4.5, total_feedbacks: 10, rating_distribution: {}, category_averages: {} } } })
  })

  it('renders page heading', async () => {
    const wrapper = mount(CustomerFeedback, { global: { plugins: [i18n, ToastService], stubs } })
    await vi.waitFor(() => expect(wrapper.text()).toContain('Customer Feedback'))
  })

  it('loads data from API', async () => {
    mount(CustomerFeedback, { global: { plugins: [i18n, ToastService], stubs } })
    await vi.waitFor(() => expect(mockGet).toHaveBeenCalledTimes(2))
  })

  it('shows summary cards when data loaded', async () => {
    mockGet
      .mockResolvedValueOnce({ data: { data: [{ comment: 'Great food!', rating: 5, category: 'food', anonymous: false, responded: false }] } })
      .mockResolvedValueOnce({ data: { data: { average_rating: 4.5, total_feedbacks: 10, rating_distribution: {}, category_averages: {} } } })
    const wrapper = mount(CustomerFeedback, { global: { plugins: [i18n, ToastService], stubs } })
    await vi.waitFor(() => {
      expect(wrapper.text()).toContain('4.5')
      expect(wrapper.text()).toContain('10')
    })
  })
})

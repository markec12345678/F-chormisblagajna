import { describe, it, expect, beforeEach, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import { createPinia, setActivePinia } from 'pinia'
import { createI18n } from 'vue-i18n'
import Home from '@/pages/Home.vue'
import ToastService from 'primevue/toastservice'

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
    getShopMode: 'retail',
    getSettings: {},
    toggleDarkMode: vi.fn(),
    setOrientation: vi.fn(),
  }),
}))

vi.mock('@/services/auth', () => ({
  default: {
    accessToken: { value: 'test-token' },
    currentUser: {
      value: {
        username: 'testuser',
        sub: 'user-1',
        roles: ['admin'],
      },
    },
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
    patch: vi.fn(),
    defaults: { headers: { common: {} } },
  },
}))

vi.mock('chart.js', () => ({
  Chart: { register: vi.fn() },
  Title: {},
  Tooltip: {},
  Legend: {},
  BarElement: {},
  CategoryScale: {},
  LinearScale: {},
  PointElement: {},
  LineElement: {},
  ArcElement: {},
  Filler: {},
}))

vi.mock('vue-chartjs', () => ({
  Line: { template: '<div class="line-stub" />' },
  Pie: { template: '<div class="pie-stub" />' },
}))

const i18n = createI18n({
  legacy: false,
  locale: 'en',
  fallbackLocale: 'en',
  messages: {
    en: {
      search: 'Search',
      product: 'Product',
      products: 'Products',
      subtotal: 'Subtotal',
      discount: 'Discount',
      total: 'Total',
      egp: 'EGP',
      next: 'Next',
      draft: 'Draft',
      add_discount: 'Add Discount',
      signout: 'Sign Out',
      profile: 'Profile',
      notifications: 'Notifications',
      clear_all: 'Clear All',
      current_orders: 'Current Orders',
      stashed_orders: 'Stashed Orders',
      cashier: 'Cashier',
      admin: 'Admin',
      kitchen: 'Kitchen',
      add_comment: 'Add Comment',
      close: 'Close',
      add: 'Add',
      order_details: 'Order Details',
      edit_order_item: 'Edit Order Item',
      payment: 'Payment',
      payment_source: 'Payment Source',
      tips: 'Tips',
      comment: 'Comment',
      customer: 'Customer',
      recap: 'Recap',
      main_details: 'Main Details',
      back: 'Back',
      submit: 'Submit',
      collect: 'Collect',
      pay_later: 'Pay Later',
      paying_later: 'Paying Later',
      collecting: 'Collecting',
      dine_in: 'Dine In',
      takeaway: 'Takeaway',
      delivery: 'Delivery',
      service_type: 'Service Type',
      name: 'Name',
      address: 'Address',
      phone: 'Phone',
      custom_data: 'Custom Data',
      key: 'Key',
      value: 'Value',
      remove: 'Remove',
      confirmation: 'Confirmation',
      start: 'Start',
      now: 'Now',
      location: 'Location',
      take_away: 'Take Away',
      other: 'Other',
      chats: 'Chats',
      write_message: 'Write message',
      no_results: 'No results',
      autostarting: 'Autostarting',
      print_receipt_enabled: 'Print Receipt',
      print_receipt_disabled: 'No Print Receipt',
      auto_finishing: 'Auto Finishing',
      auto_finish: 'Auto Finish',
      client: 'Client',
    },
  },
})

const stubs = {
  Toolbar: { template: '<div class="toolbar-stub"><slot name="start"/><slot name="center"/><slot name="end"/></div>' },
  Button: { template: '<button class="btn-stub">{{ label || "" }}</button>', props: ['label', 'icon', 'severity', 'text', 'outlined', 'size', 'rounded', 'disabled', 'aria-label'] },
  InputText: { template: '<input class="input-text-stub" />', props: ['modelValue', 'placeholder'] },
  Card: { template: '<div class="card-stub"><slot name="content"/></div>' },
  Dialog: { template: '<div class="dialog-stub" v-if="visible"><slot/></div>', props: ['visible', 'modal', 'header'] },
  Badge: { template: '<span class="badge-stub">{{ value }}</span>', props: ['value', 'class'] },
  OverlayPanel: { template: '<div class="overlay-panel-stub"><slot/></div>', props: ['ref'] },
  DataTable: { template: '<div class="datatable-stub"><slot/></div>', props: ['value'] },
  Column: { template: '<div class="column-stub"/>', props: ['field', 'header'] },
  Divider: { template: '<hr class="divider-stub" />' },
  ProgressSpinner: { template: '<div class="progress-spinner-stub" />' },
  Chip: { template: '<span class="chip-stub">{{ label }}</span>', props: ['label'] },
  Avatar: { template: '<div class="avatar-stub" />', props: ['icon', 'size'] },
  ToggleButton: { template: '<input type="checkbox" class="toggle-stub" />', props: ['modelValue', 'onLabel', 'offLabel', 'size', 'class'] },
  ButtonGroup: { template: '<div class="btn-group-stub"><slot/></div>' },
  Drawer: { template: '<div class="drawer-stub"><slot name="container"/></div>', props: ['visible'] },
  IconField: { template: '<div class="iconfield-stub"><slot/></div>' },
  InputIcon: { template: '<span class="inputicon-stub"><slot/></span>' },
  Slider: { template: '<input type="range" class="slider-stub" />', props: ['modelValue'] },
  Popover: { template: '<div class="popover-stub"><slot/></div>', props: ['ref'] },
  Stepper: { template: '<div class="stepper-stub"><slot/></div>' },
  StepList: { template: '<div class="steplist-stub"><slot/></div>' },
  Step: { template: '<div class="step-stub"><slot/></div>', props: ['value'] },
  StepPanels: { template: '<div class="steppanels-stub"><slot/></div>' },
  StepPanel: { template: '<div class="steppanel-stub"><slot :activateCallback="() => {}"/></div>', props: ['value'] },
  Textarea: { template: '<textarea class="textarea-stub" />', props: ['modelValue', 'placeholder'] },
  Select: { template: '<select class="select-stub" />', props: ['modelValue', 'options', 'optionLabel', 'placeholder'] },
  Panel: { template: '<div class="panel-stub"><slot/></div>', props: ['header'] },
  InputGroup: { template: '<div class="input-group-stub"><slot/></div>' },
  InlineMessage: { template: '<div class="inline-msg-stub"><slot/></div>', props: ['severity'] },
  Toast: { template: '<div />' },
  MainSearchResultView: { template: '<div class="main-search-result-stub" />', props: ['order'] },
  OrderView: { template: '<div class="order-view-stub" />', props: ['order'] },
  OrderItemView: { template: '<div class="order-item-view-stub" />', props: ['modelValue'] },
  MealCard: { template: '<div class="meal-card-stub" />', props: ['item'] },
  StashedOrder: { template: '<div class="stashed-order-stub" />', props: ['order'] },
  NotificationView: { template: '<div class="notification-view-stub" />', props: ['notification'] },
  PickCustomer: { template: '<div class="pick-customer-stub" />' },
  RouterLink: { template: '<a class="router-link-stub"><slot/></a>', props: ['to'] },
  Tooltip: { template: '<span class="tooltip-stub"><slot/></span>' },
}

describe('Home', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    mockGet.mockReset()
    mockPost.mockReset()
    mockDelete.mockReset()

    mockGet.mockImplementation((url) => {
      if (url.includes('/api/settings')) {
        return Promise.resolve({
          data: {
            data: {
              payment_sources: [{ name: 'Cash' }],
              language: { code: 'en' },
              shop_mode: '',
              auto_open_cash_drawer: false,
            },
          },
        })
      }
      if (url.includes('/api/languages/en')) {
        return Promise.resolve({
          data: { data: { code: 'en', pack: {}, orientation: 'ltr' } },
        })
      }
      if (url.includes('/api/categories')) {
        return Promise.resolve({ data: { data: [] } })
      }
      if (url.includes('/api/orders')) {
        return Promise.resolve({ data: { data: [] } })
      }
      if (url.includes('/api/products')) {
        return Promise.resolve({ data: { data: [] } })
      }
      return Promise.resolve({ data: { data: [] } })
    })
  })

  it('renders loading spinner initially', () => {
    mockGet.mockReturnValue(new Promise(() => {}))

    const wrapper = mount(Home, {
      global: { plugins: [i18n, ToastService], stubs },
    })

    expect(wrapper.find('.progress-spinner-stub').exists()).toBe(true)
  })

  it('renders toolbar after loading', async () => {
    const wrapper = mount(Home, {
      global: { plugins: [i18n, ToastService], stubs },
    })

    await vi.waitFor(() => {
      expect(wrapper.find('.toolbar-stub').exists()).toBe(true)
    })
  })

  it('shows product search area after loading', async () => {
    const wrapper = mount(Home, {
      global: { plugins: [i18n, ToastService], stubs },
    })

    await vi.waitFor(() => {
      expect(wrapper.findAll('.card-stub').length).toBeGreaterThan(0)
    })
  })

  it('shows total section after loading', async () => {
    const wrapper = mount(Home, {
      global: { plugins: [i18n, ToastService], stubs },
    })

    await vi.waitFor(() => {
      expect(wrapper.text()).toContain('total')
    })
  })

  it('calls settings API on mount', async () => {
    mount(Home, {
      global: { plugins: [i18n, ToastService], stubs },
    })

    await vi.waitFor(() => {
      expect(mockGet).toHaveBeenCalled()
    })
  })
})

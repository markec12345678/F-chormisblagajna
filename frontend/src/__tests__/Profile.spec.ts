import { describe, it, expect, beforeEach, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import { createPinia, setActivePinia } from 'pinia'
import { createI18n } from 'vue-i18n'
import Profile from '@/pages/Profile.vue'
import ToastService from 'primevue/toastservice'

vi.mock('vue-router', () => ({
  useRouter: () => ({
    push: vi.fn(),
  }),
}))

vi.mock('@/services/auth', () => ({
  default: {
    accessToken: { value: 'test-token' },
    currentUser: {
      value: {
        username: 'testuser',
        email: 'test@example.com',
        roles: ['admin', 'cashier'],
      },
    },
    signOut: vi.fn(),
  },
}))

const mockPatch = vi.hoisted(() => vi.fn())

vi.mock('axios', () => ({
  default: {
    get: vi.fn(),
    post: vi.fn(),
    patch: mockPatch,
    delete: vi.fn(),
    isAxiosError: vi.fn(() => false),
  },
}))

const i18n = createI18n({
  legacy: false,
  locale: 'en',
  fallbackLocale: 'en',
  messages: {
    en: {
      username: 'Username',
      email: 'Email',
      role: 'Role',
      change_password: 'Change Password',
      current_password: 'Current Password',
      new_password: 'New Password',
      confirm_password: 'Confirm Password',
      cancel: 'Cancel',
      save: 'Save',
      fill_all_fields: 'Please fill all fields',
      passwords_no_match: 'Passwords do not match',
      password_changed_success: 'Password changed successfully',
      password_change_failed: 'Failed to change password',
    },
  },
})

const stubs = {
  Card: { template: '<div class="card-stub"><slot name="content" /></div>' },
  InputText: {
    template: '<input class="input-text-stub" :disabled="disabled" :value="modelValue" />',
    props: ['modelValue', 'disabled', 'type', 'id'],
  },
  Chip: { template: '<span class="chip-stub">{{ label }}</span>', props: ['label'] },
  Divider: { template: '<hr class="divider-stub" />' },
  Button: {
    template: '<button class="btn-stub" @click="$emit(\'click\')">{{ label }}</button>',
    props: ['label', 'icon', 'severity', 'loading'],
  },
  ButtonGroup: { template: '<div class="btn-group-stub"><slot/></div>' },
  Dialog: {
    template:
      '<div class="dialog-stub" v-if="visible"><slot/><template v-if="$slots.footer"><slot name="footer"/></template></div>',
    props: ['visible', 'modal', 'header', 'style'],
  },
  Toast: { template: '<div />' },
}

describe('Profile', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    mockPatch.mockReset()
  })

  it('renders user info from auth store', async () => {
    const wrapper = mount(Profile, {
      global: { plugins: [i18n, ToastService], stubs },
    })
    await wrapper.vm.$nextTick()

    const inputs = wrapper.findAll('.input-text-stub')
    expect(inputs[0].element.value).toBe('testuser')
    expect(inputs[1].element.value).toBe('test@example.com')
  })

  it('renders user roles as chips', async () => {
    const wrapper = mount(Profile, {
      global: { plugins: [i18n, ToastService], stubs },
    })
    await wrapper.vm.$nextTick()

    const chips = wrapper.findAll('.chip-stub')
    expect(chips.length).toBe(2)
    expect(chips[0].text()).toBe('admin')
    expect(chips[1].text()).toBe('cashier')
  })

  it('renders change password button', () => {
    const wrapper = mount(Profile, {
      global: { plugins: [i18n, ToastService], stubs },
    })

    expect(wrapper.text()).toContain('Change Password')
  })

  it('shows password dialog when change password clicked', async () => {
    const wrapper = mount(Profile, {
      global: { plugins: [i18n, ToastService], stubs },
    })

    const changeBtn = wrapper.findAll('.btn-stub').find((b) => b.text() === 'Change Password')
    expect(changeBtn).toBeDefined()

    await changeBtn!.trigger('click')
    await wrapper.vm.$nextTick()

    expect(wrapper.find('.dialog-stub').exists()).toBe(true)
  })

  it('shows password form fields in dialog', async () => {
    const wrapper = mount(Profile, {
      global: { plugins: [i18n, ToastService], stubs },
    })

    const changeBtn = wrapper.findAll('.btn-stub').find((b) => b.text() === 'Change Password')
    await changeBtn!.trigger('click')
    await wrapper.vm.$nextTick()

    expect(wrapper.text()).toContain('Current Password')
    expect(wrapper.text()).toContain('New Password')
    expect(wrapper.text()).toContain('Confirm Password')
  })

  it('shows username and email labels', () => {
    const wrapper = mount(Profile, {
      global: { plugins: [i18n, ToastService], stubs },
    })

    expect(wrapper.text()).toContain('Username')
    expect(wrapper.text()).toContain('Email')
  })
})

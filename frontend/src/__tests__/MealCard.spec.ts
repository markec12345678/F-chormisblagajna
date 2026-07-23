import { describe, it, expect, beforeEach, vi } from 'vitest'
import { createPinia, setActivePinia } from 'pinia'
import { mount } from '@vue/test-utils'
import MealCard from '@/components/MealCard.vue'
import PrimeVue from 'primevue/config'
import { createI18n } from 'vue-i18n'

vi.mock('@/stores', () => ({
  globalStore: () => ({
    shop_mode: '',
    orientation: 'ltr',
    getColorMode: 'light',
  }),
}))

const i18n = createI18n({
  legacy: false,
  locale: 'en',
  fallbackLocale: 'en',
  messages: {
    en: {
      egp: 'EGP',
      possible: 'Possible',
      add_with_comment: 'Add with comment',
      inventory_consumption_disabled: 'Inventory consumption disabled',
    },
  },
})

const stubs = {
  Button: { template: '<button><slot/></button>', props: ['icon', 'severity', 'size'] },
  Menu: { template: '<div class="menu-stub"></div>' },
}

describe('MealCard', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
  })

  it('renders item name and price', () => {
    const wrapper = mount(MealCard, {
      props: {
        item: {
          name: 'Pizza Margherita',
          price: 12.5,
          image_url: 'pizza.jpg',
          enable_inventory_consumption: false,
          availability: 1,
        },
      },
      global: { plugins: [PrimeVue, i18n], stubs },
    })
    expect(wrapper.text()).toContain('Pizza Margherita')
    expect(wrapper.text()).toContain('12.5')
  })

  it('shows blocking overlay when inventory consumption disabled and availability < 1', () => {
    const wrapper = mount(MealCard, {
      props: {
        item: {
          name: 'Out of Stock',
          price: 5,
          image_url: '',
          enable_inventory_consumption: true,
          availability: 0,
        },
      },
      global: { plugins: [PrimeVue, i18n], stubs },
    })
    const overlay = wrapper.find('.w-full.h-full')
    expect(overlay.exists()).toBe(true)
    expect(overlay.attributes('style')).toContain('not-allowed')
  })

  it('does not show overlay when availability >= 1', () => {
    const wrapper = mount(MealCard, {
      props: {
        item: {
          name: 'Available Item',
          price: 8,
          image_url: '',
          enable_inventory_consumption: true,
          availability: 5,
        },
      },
      global: { plugins: [PrimeVue, i18n], stubs },
    })
    const overlay = wrapper.find('.w-full.h-full')
    expect(overlay.exists()).toBe(false)
  })
})

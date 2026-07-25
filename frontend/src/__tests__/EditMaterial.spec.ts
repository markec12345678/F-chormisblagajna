import { describe, it, expect, vi, beforeEach } from 'vitest'
import { createPinia, setActivePinia } from 'pinia'
import { mount } from '@vue/test-utils'
import EditMaterial from '@/components/EditMaterial.vue'
import PrimeVue from 'primevue/config'
import { createI18n } from 'vue-i18n'
import { Material } from '@/classes/OrderItem'

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
  messages: { en: { name: 'Name', unit: 'Unit', cancel: 'Cancel', done: 'Done' } },
})

const stubs = {
  InputText: { template: '<input />', props: ['modelValue', 'id'] },
  Button: {
    template: '<button @click="$emit(\'click\')">{{ label }}</button>',
    props: ['label', 'severity'],
  },
}

describe('EditMaterial', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
  })

  const makeMaterial = () => {
    const m = new Material()
    m.name = 'Flour'
    m.unit = 'kg'
    return m
  }

  it('renders form with name and unit', () => {
    const wrapper = mount(EditMaterial, {
      props: { material: makeMaterial() },
      global: { plugins: [PrimeVue, i18n], stubs },
    })
    expect(wrapper.text()).toContain('Name')
    expect(wrapper.text()).toContain('Unit')
  })

  it('renders cancel and done buttons', () => {
    const wrapper = mount(EditMaterial, {
      props: { material: makeMaterial() },
      global: { plugins: [PrimeVue, i18n], stubs },
    })
    expect(wrapper.text()).toContain('Cancel')
    expect(wrapper.text()).toContain('Done')
  })

  it('emits returnMaterial with edited material on Done click', async () => {
    const wrapper = mount(EditMaterial, {
      props: { material: makeMaterial() },
      global: { plugins: [PrimeVue, i18n], stubs },
    })
    const buttons = wrapper.findAll('button')
    const doneButton = buttons.find((b) => b.text().includes('Done'))
    await doneButton?.trigger('click')
    expect(wrapper.emitted('returnMaterial')).toBeTruthy()
    expect(wrapper.emitted('returnMaterial')![0][0]).toBeInstanceOf(Material)
  })

  it('initializes edited_material from prop', () => {
    const mat = makeMaterial()
    const wrapper = mount(EditMaterial, {
      props: { material: mat },
      global: { plugins: [PrimeVue, i18n], stubs },
    })
    expect(wrapper.vm.edited_material.name).toBe('Flour')
    expect(wrapper.vm.edited_material.unit).toBe('kg')
  })
})

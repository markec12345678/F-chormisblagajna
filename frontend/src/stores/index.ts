import { defineStore } from 'pinia'

const COLOR_MODE_KEY = 'color_mode'

export interface Settings {
  id: string
  shop_mode: string
  orders: {
    queues: { prefix: string; next: number }[]
    default_cost_calculation_method: string
  }
  language: { code: string; language: string }
  inventory: { stock_alert_treshold: number }
  client_receipt_printer: { host: string }
  kitchen_receipt_printer: { host: string }
  payment_sources: { name: string }[]
  auto_open_cash_drawer: boolean
}

export const globalStore = defineStore('global', {
  state: () => ({
    count: 0,
    orientation: 'ltr',
    settings: null as Settings | null,
    colorMode: localStorage.getItem(COLOR_MODE_KEY) || 'light',
    shopMode: '' as string,
  }),
  getters: {
    double: (state) => state.count * 2,
    getSettings(state) {
      return state.settings
    },
    currentOrientation(state) {
      return state.orientation
    },
    getColorMode(state) {
      return state.colorMode
    },
    getShopMode(state) {
      return state.shopMode
    },
  },
  actions: {
    increment() {
      this.count++
    },
    setOrientation(orientation: string) {
      this.orientation = orientation
    },
    setSettings(settings: Settings) {
      this.settings = settings
    },
    applyDarkModeClass() {
      if (this.colorMode === 'dark') {
        document.documentElement.classList.add('my-app-dark')
      } else {
        document.documentElement.classList.remove('my-app-dark')
      }
    },
    toggleDarkMode() {
      this.colorMode = this.colorMode === 'light' ? 'dark' : 'light'
      localStorage.setItem(COLOR_MODE_KEY, this.colorMode)
      this.applyDarkModeClass()
    },
    setShopMode(mode: string) {
      this.shopMode = mode
    },
  },
})

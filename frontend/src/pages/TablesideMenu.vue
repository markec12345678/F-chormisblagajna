<template>
  <div class="tableside-menu">
    <div class="header text-center py-3">
      <h2>{{ $t('table_menu') }}</h2>
      <p v-if="session" class="text-lg">{{ session.table_label }}</p>
    </div>

    <div class="menu-items px-3" v-if="products.length > 0">
      <div v-for="product in products" :key="product.id" class="menu-card flex align-items-center justify-content-between p-3 border-bottom-1 surface-border">
        <div>
          <div class="font-bold">{{ product.name }}</div>
          <div class="text-sm text-500">{{ product.description }}</div>
          <div class="text-primary font-bold mt-1">{{ formatCurrency(product.price) }}</div>
        </div>
        <div class="flex align-items-center gap-2">
          <InputNumber v-model="quantities[product.id]" :min="0" :max="99" class="w-4rem" />
          <Button icon="pi pi-plus" size="small" @click="addToCart(product)" :disabled="!quantities[product.id] || quantities[product.id] < 1" />
        </div>
      </div>
    </div>

    <div v-else-if="!loading" class="text-center py-5 text-500">
      {{ $t('no_products') }}
    </div>

    <div v-if="cart.length > 0" class="cart-section p-3 surface-section">
      <h3>{{ $t('your_cart') }} ({{ cart.length }})</h3>
      <div v-for="(item, idx) in cart" :key="idx" class="flex justify-content-between align-items-center py-2">
        <div>
          <div class="font-medium">{{ item.name }}</div>
          <div class="text-sm">{{ item.qty }}x {{ formatCurrency(item.price) }}</div>
        </div>
        <div class="flex align-items-center gap-2">
          <span class="font-bold">{{ formatCurrency(item.qty * item.price) }}</span>
          <Button icon="pi pi-trash" severity="danger" size="small" text @click="removeFromCart(idx)" />
        </div>
      </div>
      <div class="flex justify-content-between font-bold text-lg py-2 border-top-1 surface-border mt-2">
        <span>{{ $t('total') }}</span>
        <span>{{ formatCurrency(cartTotal) }}</span>
      </div>
      <Button :label="$t('place_order')" class="w-full mt-2" @click="placeOrder" :loading="ordering" />
    </div>

    <Dialog v-model:visible="orderPlaced" :header="$t('order_placed')" :style="{ width: '350px' }">
      <p>{{ $t('order_placed_message') }}</p>
      <template #footer>
        <Button :label="$t('close')" severity="secondary" @click="orderPlaced = false; cart = []; loadSession()" />
      </template>
    </Dialog>
  </div>
</template>

<script setup lang="ts">
import Button from 'primevue/button'
import InputNumber from 'primevue/inputnumber'
import Dialog from 'primevue/dialog'
import axios from 'axios'
import { ref, computed, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { useToast } from 'primevue/usetoast'
import { useI18n } from 'vue-i18n'

const route = useRoute()
const toast = useToast()
const { t } = useI18n()

interface Product {
  id: string
  name: string
  description: string
  price: number
}

interface Session {
  id: string
  table_label: string
  qr_token: string
  active: boolean
}

const session = ref<Session | null>(null)
const products = ref<Product[]>([])
const quantities = ref<Record<string, number>>({})
const cart = ref<Array<{ id: string; name: string; price: number; qty: number }>>([])
const loading = ref(true)
const ordering = ref(false)
const orderPlaced = ref(false)

const formatCurrency = (amount: number) => {
  return new Intl.NumberFormat('sl-SI', { style: 'currency', currency: 'EUR' }).format(amount)
}

const cartTotal = computed(() => cart.value.reduce((sum, i) => sum + i.price * i.qty, 0))

const addToCart = (product: Product) => {
  const qty = quantities.value[product.id] || 1
  const existing = cart.value.find((i) => i.id === product.id)
  if (existing) {
    existing.qty += qty
  } else {
    cart.value.push({ id: product.id, name: product.name, price: product.price, qty })
  }
  quantities.value[product.id] = 0
}

const removeFromCart = (idx: number) => {
  cart.value.splice(idx, 1)
}

const loadSession = async () => {
  loading.value = true
  try {
    const token = route.params.token as string
    const menuResponse = await axios.get(
      `http://${import.meta.env.VITE_APP_BACKEND_HOST}/tableside/api/menu/${token}/orders`
    )
  } catch {
    // silently fail
  } finally {
    loading.value = false
  }
}

const placeOrder = async () => {
  if (!session.value) return
  ordering.value = true
  try {
    const items = cart.value.map((i) => ({
      product_id: i.id,
      product_name: i.name,
      quantity: i.qty,
      unit_price: i.price,
      notes: '',
    }))

    await axios.post(
      `http://${import.meta.env.VITE_APP_BACKEND_HOST}/tableside/api/menu/place-order`,
      { session_id: session.value.id, items }
    )
    orderPlaced.value = true
  } catch {
    toast.add({ severity: 'error', summary: t('failed'), detail: t('order_failed'), life: 3000 })
  } finally {
    ordering.value = false
  }
}

onMounted(async () => {
  loading.value = true
  try {
    const token = route.params.token as string
    const sessResponse = await axios.get(`http://${import.meta.env.VITE_APP_BACKEND_HOST}/core/api/products`)
    products.value = sessResponse.data.data || []

    const tokenParts = token.split('_')
    if (tokenParts.length > 0) {
      session.value = { id: token, table_label: t('table'), qr_token: token, active: true }
    }
  } catch {
    // silently fail
  } finally {
    loading.value = false
  }
})
</script>

<style scoped>
.tableside-menu {
  max-width: 600px;
  margin: 0 auto;
  min-height: 100vh;
  background: var(--surface-ground);
}
.menu-card:last-child {
  border-bottom: none !important;
}
.cart-section {
  position: sticky;
  bottom: 0;
  border-top: 1px solid var(--surface-border);
}
</style>

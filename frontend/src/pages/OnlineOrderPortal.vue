<template>
  <div class="w-full min-h-screen" style="background: var(--surface-ground)">
    <Toolbar class="mb-4">
      <template #start>
        <div class="flex align-items-center gap-3">
          <img src="@/assets/logo.png" alt="logo" style="height: 30px" v-if="store.getColorMode == 'light'" />
          <h2 class="m-0">{{ $t('order_online') }}</h2>
        </div>
      </template>
      <template #end>
        <Button :label="`${$t('cart')} (${cartItems.length})`" icon="pi pi-shopping-cart" severity="success"
          @click="showCart = true" :badge="cartItems.length.toString()" badgeClass="p-badge-danger" />
      </template>
    </Toolbar>

    <div class="grid mx-2">
      <div class="col-12">
        <div class="flex flex-wrap gap-2 mb-4">
          <Button :label="$t('all')" :severity="selectedCategory === '' ? 'primary' : 'secondary'"
            @click="selectedCategory = ''" />
          <Button v-for="cat in menuCategories" :key="cat.id" :label="cat.name"
            :severity="selectedCategory === cat.id ? 'primary' : 'secondary'"
            @click="selectedCategory = cat.id" />
        </div>
      </div>

      <div class="col-12" v-if="isLoading">
        <div class="text-center py-5">
          <ProgressSpinner style="width: 50px; height: 50px" />
        </div>
      </div>

      <div class="col-12" v-if="!isLoading && filteredProducts.length === 0">
        <Card>
          <template #content>
            <div class="text-center py-5">
              <i class="pi pi-inbox text-5xl text-400 mb-3"></i>
              <p class="text-400">{{ $t('no_products_available') }}</p>
            </div>
          </template>
        </Card>
      </div>

      <div class="col-12 md:col-6 lg:col-3" v-for="product in filteredProducts" :key="product.id">
        <Card class="h-full">
          <template #header>
            <div v-if="product.image_url" class="overflow-hidden">
              <img :src="product.image_url" :alt="product.name" class="w-full" style="height: 180px; object-fit: cover" />
            </div>
            <div v-else class="flex justify-content-center align-items-center" style="height: 180px; background: var(--surface-100)">
              <i class="pi pi-image text-4xl text-300"></i>
            </div>
          </template>
          <template #title>
            <div class="flex justify-content-between align-items-center">
              <span class="text-lg">{{ product.name }}</span>
              <Tag :value="product.category" severity="info" />
            </div>
          </template>
          <template #content>
            <div class="flex justify-content-between align-items-center">
              <span class="text-xl font-bold text-primary">{{ formatCurrency(product.price) }}</span>
              <div class="flex align-items-center gap-2">
                <Button v-if="getCartQuantity(product.id) > 0" icon="pi pi-minus" severity="danger" text rounded
                  @click="removeFromCart(product)" />
                <span v-if="getCartQuantity(product.id) > 0" class="font-bold">{{ getCartQuantity(product.id) }}</span>
                <Button icon="pi pi-plus" severity="success" text rounded @click="addToCart(product)" />
              </div>
            </div>
          </template>
        </Card>
      </div>
    </div>

    <Dialog v-model:visible="showCart" modal :header="$t('your_cart')" :style="{ width: '45rem' }">
      <div v-if="cartItems.length === 0" class="text-center py-5">
        <i class="pi pi-shopping-cart text-5xl text-400 mb-3"></i>
        <p class="text-400">{{ $t('cart_empty') }}</p>
      </div>
      <div v-else>
        <DataTable :value="cartItems">
          <Column field="name" :header="$t('item')"></Column>
          <Column field="quantity" :header="$t('quantity')">
            <template #body="slotProps">
              <div class="flex align-items-center gap-2">
                <Button icon="pi pi-minus" text rounded size="small" @click="decrementCart(slotProps?.data)" />
                <span>{{ slotProps?.data?.quantity }}</span>
                <Button icon="pi pi-plus" text rounded size="small" @click="incrementCart(slotProps?.data)" />
              </div>
            </template>
          </Column>
          <Column header="$t('subtotal')">
            <template #body="slotProps">
              {{ formatCurrency((slotProps?.data?.price || 0) * (slotProps?.data?.quantity || 0)) }}
            </template>
          </Column>
          <Column :header="$t('actions')">
            <template #body="slotProps">
              <Button icon="pi pi-trash" severity="danger" text rounded size="small"
                @click="removeCartItem(slotProps?.data?.product_id)" />
            </template>
          </Column>
        </DataTable>

        <Divider />

        <div class="grid">
          <div class="col-12 md:col-6">
            <div class="flex flex-column gap-3">
              <div class="flex flex-column gap-2">
                <label>{{ $t('customer_name') }} *</label>
                <InputText v-model="checkoutForm.customer_name" />
              </div>
              <div class="flex flex-column gap-2">
                <label>{{ $t('phone') }} *</label>
                <InputText v-model="checkoutForm.customer_phone" />
              </div>
              <div class="flex flex-column gap-2">
                <label>{{ $t('email') }}</label>
                <InputText v-model="checkoutForm.customer_email" />
              </div>
            </div>
          </div>
          <div class="col-12 md:col-6">
            <div class="flex flex-column gap-3">
              <div class="flex flex-column gap-2">
                <label>{{ $t('order_type') }}</label>
                <Dropdown v-model="checkoutForm.order_type" :options="orderTypes" optionLabel="label"
                  optionValue="value" />
              </div>
              <div v-if="checkoutForm.order_type === 'delivery'" class="flex flex-column gap-2">
                <label>{{ $t('delivery_address') }}</label>
                <Textarea v-model="checkoutForm.delivery_addr" rows="2" />
              </div>
              <div class="flex flex-column gap-2">
                <label>{{ $t('notes') }}</label>
                <Textarea v-model="checkoutForm.notes" rows="2" />
              </div>
            </div>
          </div>
        </div>

        <Divider />

        <div class="flex justify-content-between font-bold text-xl">
          <span>{{ $t('total') }}</span>
          <span>{{ formatCurrency(cartTotal) }}</span>
        </div>
      </div>

      <template #footer>
        <Button :label="$t('cancel')" severity="secondary" @click="showCart = false" />
        <Button v-if="cartItems.length > 0" :label="$t('place_order')" icon="pi pi-check" severity="success"
          @click="placeOrder" :loading="isPlacing" />
      </template>
    </Dialog>
  </div>
</template>

<script setup lang="ts">
import Toolbar from 'primevue/toolbar'
import Card from 'primevue/card'
import Button from 'primevue/button'
import Tag from 'primevue/tag'
import Dialog from 'primevue/dialog'
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import InputText from 'primevue/inputtext'
import Textarea from 'primevue/textarea'
import Dropdown from 'primevue/dropdown'
import Divider from 'primevue/divider'
import ProgressSpinner from 'primevue/progressspinner'
import axios from 'axios'
import { ref, computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { useToast } from 'primevue/usetoast'
import { globalStore } from '@/stores'

const { t } = useI18n()
const toast = useToast()
const store = globalStore()

interface Product {
  id: string
  name: string
  price: number
  image_url: string
  category: string
  available: boolean
}

interface MenuCategory {
  id: string
  name: string
  products: Product[]
}

interface CartItem {
  product_id: string
  name: string
  price: number
  quantity: number
}

const menuCategories = ref<MenuCategory[]>([])
const selectedCategory = ref('')
const isLoading = ref(false)
const showCart = ref(false)
const isPlacing = ref(false)
const cartItems = ref<CartItem[]>([])

const checkoutForm = ref({
  customer_name: '',
  customer_phone: '',
  customer_email: '',
  order_type: 'takeaway',
  delivery_addr: '',
  notes: '',
})

const orderTypes = [
  { label: 'Takeaway', value: 'takeaway' },
  { label: 'Delivery', value: 'delivery' },
  { label: 'Dine In', value: 'dine_in' },
]

const filteredProducts = computed(() => {
  const allProducts: Product[] = []
  for (const cat of menuCategories.value) {
    if (selectedCategory.value === '' || cat.id === selectedCategory.value) {
      allProducts.push(...cat.products)
    }
  }
  return allProducts
})

const cartTotal = computed(() => {
  return cartItems.value.reduce((sum, item) => sum + item.price * item.quantity, 0)
})

const formatCurrency = (amount: number) => {
  return new Intl.NumberFormat('sl-SI', { style: 'currency', currency: 'EUR' }).format(amount)
}

const getCartQuantity = (productId: string) => {
  const item = cartItems.value.find((i) => i.product_id === productId)
  return item ? item.quantity : 0
}

const addToCart = (product: Product) => {
  const existing = cartItems.value.find((i) => i.product_id === product.id)
  if (existing) {
    existing.quantity++
  } else {
    cartItems.value.push({
      product_id: product.id,
      name: product.name,
      price: product.price,
      quantity: 1,
    })
  }
}

const removeFromCart = (product: Product) => {
  const existing = cartItems.value.find((i) => i.product_id === product.id)
  if (existing && existing.quantity > 1) {
    existing.quantity--
  } else {
    cartItems.value = cartItems.value.filter((i) => i.product_id !== product.id)
  }
}

const incrementCart = (item: CartItem) => {
  item.quantity++
}

const decrementCart = (item: CartItem) => {
  if (item.quantity > 1) {
    item.quantity--
  } else {
    cartItems.value = cartItems.value.filter((i) => i.product_id !== item.product_id)
  }
}

const removeCartItem = (productId: string) => {
  cartItems.value = cartItems.value.filter((i) => i.product_id !== productId)
}

const placeOrder = async () => {
  if (!checkoutForm.value.customer_name || !checkoutForm.value.customer_phone) {
    toast.add({ severity: 'warn', summary: t('warning'), detail: t('fill_required_fields'), group: 'br', life: 3000 })
    return
  }

  isPlacing.value = true
  try {
    await axios.post(
      `http://${import.meta.env.VITE_APP_BACKEND_HOST}/onlineorder/api/order`,
      {
        customer_name: checkoutForm.value.customer_name,
        customer_phone: checkoutForm.value.customer_phone,
        customer_email: checkoutForm.value.customer_email,
        order_type: checkoutForm.value.order_type,
        delivery_addr: checkoutForm.value.delivery_addr,
        notes: checkoutForm.value.notes,
        items: cartItems.value.map((i) => ({
          product_id: i.product_id,
          product_name: i.name,
          quantity: i.quantity,
          price: i.price,
        })),
      }
    )
    toast.add({ severity: 'success', summary: t('success'), detail: t('order_placed'), group: 'br', life: 5000 })
    cartItems.value = []
    showCart.value = false
    checkoutForm.value = { customer_name: '', customer_phone: '', customer_email: '', order_type: 'takeaway', delivery_addr: '', notes: '' }
  } catch {
    toast.add({ severity: 'error', summary: t('failed'), detail: t('order_failed'), group: 'br', life: 3000 })
  } finally {
    isPlacing.value = false
  }
}

const loadMenu = async () => {
  isLoading.value = true
  try {
    const response = await axios.get(
      `http://${import.meta.env.VITE_APP_BACKEND_HOST}/onlineorder/api/menu`
    )
    menuCategories.value = response.data.data || []
  } catch {
    toast.add({ severity: 'error', summary: t('failed'), detail: t('menu_load_failed'), group: 'br', life: 3000 })
  } finally {
    isLoading.value = false
  }
}

loadMenu()
</script>

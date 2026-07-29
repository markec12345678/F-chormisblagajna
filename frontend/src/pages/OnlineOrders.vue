<template>
  <div class="w-full">
    <div class="grid mx-2">
      <div class="col-12">
        <div class="flex justify-content-between align-items-center">
          <h3>{{ $t('online_orders') }}</h3>
          <div class="flex gap-2">
            <Dropdown
              v-model="statusFilter"
              :options="statusOptions"
              optionLabel="label"
              optionValue="value"
              :placeholder="$t('all_statuses')"
              class="w-12rem"
            />
            <Button
              :label="$t('refresh')"
              icon="pi pi-refresh"
              @click="loadOrders"
              :loading="isLoading"
            />
          </div>
        </div>
      </div>

      <div class="col-12" v-if="!isLoading && orders.length === 0">
        <Card>
          <template #content>
            <div class="text-center py-5">
              <i class="pi pi-shopping-cart text-5xl text-400 mb-3"></i>
              <p class="text-400">{{ $t('no_online_orders') }}</p>
            </div>
          </template>
        </Card>
      </div>

      <div class="col-12 md:col-6 lg:col-4" v-for="order in orders" :key="order.id">
        <Card>
          <template #title>
            <div class="flex justify-content-between align-items-center">
              <span>{{ order.display_id }}</span>
              <Tag :value="order.status" :severity="getStatusSeverity(order.status)" />
            </div>
          </template>
          <template #content>
            <div class="flex flex-column gap-2">
              <div class="flex justify-content-between">
                <span class="text-500">{{ $t('customer') }}:</span>
                <span>{{ order.customer_name }}</span>
              </div>
              <div class="flex justify-content-between">
                <span class="text-500">{{ $t('phone') }}:</span>
                <span>{{ order.customer_phone }}</span>
              </div>
              <div class="flex justify-content-between">
                <span class="text-500">{{ $t('order_type') }}:</span>
                <Tag :value="order.order_type" severity="info" />
              </div>
              <div v-if="order.delivery_addr" class="flex justify-content-between">
                <span class="text-500">{{ $t('address') }}:</span>
                <span class="text-right" style="max-width: 200px">{{ order.delivery_addr }}</span>
              </div>
              <Divider />
              <div
                v-for="item in order.items"
                :key="item.product_id"
                class="flex justify-content-between text-sm"
              >
                <span>{{ item.quantity }}x {{ item.product_name }}</span>
                <span>{{ formatCurrency(item.price * item.quantity) }}</span>
              </div>
              <Divider />
              <div class="flex justify-content-between font-bold">
                <span>{{ $t('total') }}</span>
                <span>{{ formatCurrency(order.total) }}</span>
              </div>
              <div v-if="order.notes" class="text-sm text-500 mt-2">
                <i class="pi pi-comment mr-1"></i> {{ order.notes }}
              </div>
              <div class="text-sm text-400 mt-1">
                {{ formatDateTime(order.created_at) }}
              </div>
            </div>
          </template>
          <template #footer>
            <div class="flex gap-2">
              <Button
                v-if="order.status === 'pending'"
                :label="$t('confirm')"
                icon="pi pi-check"
                severity="success"
                size="small"
                @click="updateStatus(order.id, 'confirmed')"
              />
              <Button
                v-if="order.status === 'confirmed'"
                :label="$t('preparing')"
                icon="pi pi-spin pi-spinner"
                severity="warn"
                size="small"
                @click="updateStatus(order.id, 'preparing')"
              />
              <Button
                v-if="order.status === 'preparing'"
                :label="$t('ready')"
                icon="pi pi-check-circle"
                severity="info"
                size="small"
                @click="updateStatus(order.id, 'ready')"
              />
              <Button
                v-if="order.status === 'ready'"
                :label="$t('delivered')"
                icon="pi pi-verified"
                severity="success"
                size="small"
                @click="updateStatus(order.id, 'delivered')"
              />
              <Button
                v-if="order.status !== 'cancelled' && order.status !== 'delivered'"
                :label="$t('cancel')"
                severity="danger"
                size="small"
                text
                @click="updateStatus(order.id, 'cancelled')"
              />
            </div>
          </template>
        </Card>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import Card from 'primevue/card'
import Button from 'primevue/button'
import Tag from 'primevue/tag'
import Dropdown from 'primevue/dropdown'
import Divider from 'primevue/divider'
import axios from 'axios'
import { ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { useToast } from 'primevue/usetoast'
import auth from '../services/auth'

const { t } = useI18n()
const toast = useToast()

interface OrderItem {
  product_id: string
  product_name: string
  quantity: number
  price: number
  comment: string
}

interface OnlineOrder {
  id: string
  display_id: string
  customer_name: string
  customer_phone: string
  customer_email: string
  items: OrderItem[]
  subtotal: number
  delivery_fee: number
  total: number
  order_type: string
  delivery_addr: string
  status: string
  notes: string
  created_at: string
  updated_at: string
}

const orders = ref<OnlineOrder[]>([])
const isLoading = ref(false)
const statusFilter = ref('')

const statusOptions = [
  { label: 'All', value: '' },
  { label: 'Pending', value: 'pending' },
  { label: 'Confirmed', value: 'confirmed' },
  { label: 'Preparing', value: 'preparing' },
  { label: 'Ready', value: 'ready' },
  { label: 'Delivered', value: 'delivered' },
  { label: 'Cancelled', value: 'cancelled' },
]

const formatCurrency = (amount: number) => {
  return new Intl.NumberFormat('sl-SI', { style: 'currency', currency: 'EUR' }).format(amount)
}

const formatDateTime = (date: string) => {
  if (!date) return '-'
  return new Date(date).toLocaleString()
}

const getStatusSeverity = (status: string) => {
  switch (status) {
    case 'pending':
      return 'warn'
    case 'confirmed':
      return 'info'
    case 'preparing':
      return 'secondary'
    case 'ready':
      return 'success'
    case 'delivered':
      return 'success'
    case 'cancelled':
      return 'danger'
    default:
      return 'secondary'
  }
}

const loadOrders = async () => {
  isLoading.value = true
  try {
    const params: Record<string, string> = {}
    if (statusFilter.value) {
      params.status = statusFilter.value
    }
    const response = await axios.get(
      `http://${import.meta.env.VITE_APP_BACKEND_HOST}/onlineorder/api/orders`,
      {
        headers: { Authorization: `Bearer ${auth.accessToken.value}` },
        params,
      },
    )
    orders.value = response.data.data || []
  } catch {
    toast.add({
      severity: 'error',
      summary: t('failed'),
      detail: t('load_failed'),
      group: 'br',
      life: 3000,
    })
  } finally {
    isLoading.value = false
  }
}

const updateStatus = async (orderId: string, status: string) => {
  try {
    await axios.put(
      `http://${import.meta.env.VITE_APP_BACKEND_HOST}/onlineorder/api/orders/${orderId}/status`,
      { status },
      { headers: { Authorization: `Bearer ${auth.accessToken.value}` } },
    )
    toast.add({
      severity: 'success',
      summary: t('success'),
      detail: t('order_status_updated'),
      group: 'br',
      life: 3000,
    })
    loadOrders()
  } catch {
    toast.add({
      severity: 'error',
      summary: t('failed'),
      detail: t('update_failed'),
      group: 'br',
      life: 3000,
    })
  }
}

loadOrders()
</script>

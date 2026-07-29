<template>
  <div class="p-3">
    <Toast />

    <div class="flex align-items-center justify-content-between mb-3">
      <h2 class="m-0">{{ $t('admin_dashboard') }}</h2>
      <Button
        :label="$t('refresh')"
        icon="pi pi-refresh"
        severity="secondary"
        @click="loadAll"
        :loading="loading"
      />
    </div>

    <div v-if="loading" class="flex justify-content-center py-6">
      <ProgressSpinner style="width: 50px; height: 50px" strokeWidth="4" />
    </div>

    <template v-if="!loading">
      <div class="grid">
        <div class="col-12 md:col-6 lg:col-3">
          <Card>
            <template #title>
              <div class="flex align-items-center gap-2">
                <i class="pi pi-dollar text-green-500" style="font-size: 1.5rem" />
                <span class="text-lg">{{ $t('today_sales') }}</span>
              </div>
            </template>
            <template #content>
              <div class="text-3xl font-bold">{{ formatCurrency(todaySales) }}</div>
              <div
                class="flex align-items-center gap-1 mt-2"
                :class="salesChange >= 0 ? 'text-green-500' : 'text-red-500'"
              >
                <i
                  :class="salesChange >= 0 ? 'pi pi-arrow-up' : 'pi pi-arrow-down'"
                  style="font-size: 0.8rem"
                />
                <span>{{ Math.abs(salesChange).toFixed(1) }}% {{ $t('yesterday') }}</span>
              </div>
            </template>
          </Card>
        </div>

        <div class="col-12 md:col-6 lg:col-3">
          <Card>
            <template #title>
              <div class="flex align-items-center gap-2">
                <i class="pi pi-shopping-cart text-blue-500" style="font-size: 1.5rem" />
                <span class="text-lg">{{ $t('orders_today') }}</span>
              </div>
            </template>
            <template #content>
              <div class="text-3xl font-bold">{{ todayOrders }}</div>
              <div class="text-sm text-gray-500 mt-2">
                {{ $t('total') }} {{ $t('sales') }}: {{ formatCurrency(todaySales) }}
              </div>
            </template>
          </Card>
        </div>

        <div class="col-12 md:col-6 lg:col-3">
          <Card>
            <template #title>
              <div class="flex align-items-center gap-2">
                <i class="pi pi-clock text-orange-500" style="font-size: 1.5rem" />
                <span class="text-lg">{{ $t('active_orders') }}</span>
              </div>
            </template>
            <template #content>
              <div class="text-3xl font-bold">{{ activeOrders }}</div>
              <div class="text-sm text-gray-500 mt-2">
                {{ $t('pending_orders') }}: {{ pendingOrders }}
              </div>
            </template>
          </Card>
        </div>

        <div class="col-12 md:col-6 lg:col-3">
          <Card>
            <template #title>
              <div class="flex align-items-center gap-2">
                <i class="pi pi-exclamation-triangle text-red-500" style="font-size: 1.5rem" />
                <span class="text-lg">{{ $t('low_stock_alerts') }}</span>
              </div>
            </template>
            <template #content>
              <div class="text-3xl font-bold">{{ lowStockCount }}</div>
              <div class="text-sm text-gray-500 mt-2">
                {{ $t('average_rating') }}: {{ avgFeedback.toFixed(1) }}/5
              </div>
            </template>
          </Card>
        </div>
      </div>

      <div class="grid mt-3">
        <div class="col-12 lg:col-8">
          <Card>
            <template #title>
              <i class="pi pi-chart-bar mr-2" />{{ $t('quick_actions') }}
            </template>
            <template #content>
              <div class="grid">
                <div
                  v-for="action in quickActions"
                  :key="action.key"
                  class="col-6 md:col-4 lg:col-3"
                >
                  <Button
                    :label="$t(action.label)"
                    :icon="action.icon"
                    class="w-full p-3"
                    severity="secondary"
                    outlined
                    @click="$router.push(action.link)"
                  />
                </div>
              </div>
            </template>
          </Card>
        </div>

        <div class="col-12 lg:col-4">
          <Card>
            <template #title> <i class="pi pi-star mr-2" />{{ $t('recent_feedback') }} </template>
            <template #content>
              <div v-if="recentFeedback.length === 0" class="text-gray-400 text-center py-3">
                {{ $t('no_feedback') }}
              </div>
              <div
                v-for="fb in recentFeedback.slice(0, 5)"
                :key="fb.id"
                class="flex align-items-center gap-2 mb-2 pb-2"
                style="border-bottom: 1px solid var(--surface-border)"
              >
                <div class="flex-shrink-0">
                  <Rating :modelValue="fb.rating" :cancel="false" readonly />
                </div>
                <div class="flex-1 text-sm">
                  {{
                    fb.comment
                      ? fb.comment.substring(0, 60) + (fb.comment.length > 60 ? '...' : '')
                      : '-'
                  }}
                </div>
              </div>
            </template>
          </Card>
        </div>
      </div>
    </template>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useToast } from 'primevue/usetoast'
import axios from 'axios'

const toast = useToast()
const loading = ref(true)
const todaySales = ref(0)
const yesterdaySales = ref(0)
const salesChange = ref(0)
const todayOrders = ref(0)
const activeOrders = ref(0)
const pendingOrders = ref(0)
const lowStockCount = ref(0)
const avgFeedback = ref(0)
const recentFeedback = ref([])

const API = import.meta.env.VITE_APP_BACKEND_HOST || ''
const CORE = import.meta.env.VITE_APP_MODULE_CORE_API_PREFIX || '/core'

const quickActions = [
  { key: 'orders', label: 'order|plural', icon: 'pi pi-box', link: '/admin/orders' },
  { key: 'inventory', label: 'inventory', icon: 'pi pi-boxes', link: '/admin/inventory' },
  { key: 'products', label: 'product|plural', icon: 'pi pi-barcode', link: '/admin/products' },
  { key: 'customers', label: 'customer|plural', icon: 'pi pi-users', link: '/admin/customers' },
  { key: 'sales', label: 'sales', icon: 'pi pi-chart-line', link: '/admin/sales' },
  {
    key: 'feedback',
    label: 'customer_feedback',
    icon: 'pi pi-comments',
    link: '/admin/customer-feedback',
  },
]

function formatCurrency(v) {
  return (
    (v || 0).toLocaleString('sl-SI', { minimumFractionDigits: 2, maximumFractionDigits: 2 }) +
    ' ' +
    '\u20AC'
  )
}

async function loadAll() {
  loading.value = true
  try {
    const today = new Date()
    const todayStr = today.toISOString().slice(0, 10)
    const yesterdayStr = new Date(today.getTime() - 86400000).toISOString().slice(0, 10)

    const [salesRes, ordersRes, feedbackRes] = await Promise.allSettled([
      axios.get(`${API}${CORE}/api/logs/salesperday`, { params: { page: 1, size: 50 } }),
      axios.get(`${API}${CORE}/api/orders`, { params: { page: 1, size: 100 } }),
      axios.get(`${API}/feedback/api/feedback/summary`),
    ])

    if (salesRes.status === 'fulfilled') {
      const sales = salesRes.value.data?.data || []
      const todayEntry = sales.find((s) => s.date?.startsWith(todayStr))
      const yesterdayEntry = sales.find((s) => s.date?.startsWith(yesterdayStr))
      todaySales.value = todayEntry?.total_sales || 0
      yesterdaySales.value = yesterdayEntry?.total_sales || 0
      salesChange.value =
        yesterdaySales.value > 0
          ? ((todaySales.value - yesterdaySales.value) / yesterdaySales.value) * 100
          : 0
      todayOrders.value = sales
        .filter((s) => s.date?.startsWith(todayStr))
        .reduce((a, b) => a + (b.orders?.length || 0), 0)
    }

    if (ordersRes.status === 'fulfilled') {
      const orders = ordersRes.value.data?.data || []
      activeOrders.value = orders.filter((o) => o.state === 'in_progress').length
      pendingOrders.value = orders.filter((o) => o.state === 'pending').length
    }

    if (feedbackRes.status === 'fulfilled') {
      const summary = feedbackRes.value.data?.data
      avgFeedback.value = summary?.average_rating || 0
      recentFeedback.value = summary?.recent_feedbacks || []
    }
  } catch {
    toast.add({
      severity: 'error',
      summary: 'Error',
      detail: 'Failed to load dashboard data',
      life: 3000,
    })
  } finally {
    loading.value = false
  }
}

onMounted(loadAll)
</script>

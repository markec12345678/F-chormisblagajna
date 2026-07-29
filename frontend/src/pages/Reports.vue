<template>
  <div class="w-full">
    <div class="grid mx-2">
      <div class="col-12">
        <h3>{{ $t('reports') }}</h3>
      </div>

      <div class="col-12 md:col-4">
        <Card>
          <template #title>{{ $t('today_revenue') }}</template>
          <template #content>
            <div class="text-3xl font-bold">{{ dashboard?.today_revenue?.toFixed(2) || '0.00' }} EUR</div>
          </template>
        </Card>
      </div>
      <div class="col-12 md:col-4">
        <Card>
          <template #title>{{ $t('today_orders') }}</template>
          <template #content>
            <div class="text-3xl font-bold">{{ dashboard?.today_orders || 0 }}</div>
          </template>
        </Card>
      </div>
      <div class="col-12 md:col-4">
        <Card>
          <template #title>{{ $t('average_order') }}</template>
          <template #content>
            <div class="text-3xl font-bold">{{ dashboard?.average_order?.toFixed(2) || '0.00' }} EUR</div>
          </template>
        </Card>
      </div>

      <div class="col-12">
        <h3>{{ $t('sales_report') }}</h3>
      </div>
      <div class="col-12">
        <div class="flex gap-2 mb-4">
          <Calendar v-model="startDate" :placeholder="$t('start_date')" showIcon />
          <Calendar v-model="endDate" :placeholder="$t('end_date')" showIcon />
          <Button :label="$t('generate')" icon="pi pi-refresh" @click="loadSalesReport" :loading="salesLoading" />
          <Button v-if="salesReport" label="CSV" icon="pi pi-download" severity="success" @click="exportSalesCSV" />
        </div>
      </div>

      <div v-if="salesReport" class="col-12">
        <Card>
          <template #title>{{ $t('sales_summary') }}</template>
          <template #content>
            <div class="grid">
              <div class="col-12 md:col-3">
                <div class="text-sm text-gray-500">{{ $t('total_revenue') }}</div>
                <div class="text-xl font-bold">{{ salesReport.total_revenue?.toFixed(2) || '0.00' }} EUR</div>
              </div>
              <div class="col-12 md:col-3">
                <div class="text-sm text-gray-500">{{ $t('total_orders') }}</div>
                <div class="text-xl font-bold">{{ salesReport.total_orders || 0 }}</div>
              </div>
              <div class="col-12 md:col-3">
                <div class="text-sm text-gray-500">{{ $t('total_items') }}</div>
                <div class="text-xl font-bold">{{ salesReport.total_items || 0 }}</div>
              </div>
              <div class="col-12 md:col-3">
                <div class="text-sm text-gray-500">{{ $t('net_revenue') }}</div>
                <div class="text-xl font-bold text-green-600">{{ salesReport.net_revenue?.toFixed(2) || '0.00' }} EUR</div>
              </div>
            </div>

            <div v-if="salesReport.top_products?.length" class="mt-4">
              <h4>{{ $t('top_products') }}</h4>
              <DataTable :value="salesReport.top_products" stripedRows>
                <Column field="name" :header="$t('product')"></Column>
                <Column field="quantity" :header="$t('quantity')"></Column>
                <Column field="revenue" :header="$t('revenue')">
                  <template #body="slotProps">
                    {{ slotProps.data.revenue?.toFixed(2) }} EUR
                  </template>
                </Column>
              </DataTable>
            </div>
          </template>
        </Card>
      </div>

      <div class="col-12 mt-4">
        <div class="flex justify-content-between align-items-center">
          <h3>{{ $t('inventory_report') }}</h3>
          <Button :label="$t('refresh')" icon="pi pi-refresh" @click="loadInventoryReport" :loading="inventoryLoading" />
          <Button v-if="inventoryReport" label="CSV" icon="pi pi-download" severity="success" @click="exportInventoryCSV" />
        </div>
      </div>

      <div v-if="inventoryReport" class="col-12">
        <Card>
          <template #title>{{ $t('inventory_status') }}</template>
          <template #content>
            <div class="grid">
              <div class="col-12 md:col-3">
                <div class="text-sm text-gray-500">{{ $t('total_materials') }}</div>
                <div class="text-xl font-bold">{{ inventoryReport.total_materials || 0 }}</div>
              </div>
              <div class="col-12 md:col-3">
                <div class="text-sm text-gray-500">{{ $t('low_stock') }}</div>
                <div class="text-xl font-bold text-yellow-600">{{ inventoryReport.low_stock_count || 0 }}</div>
              </div>
              <div class="col-12 md:col-3">
                <div class="text-sm text-gray-500">{{ $t('out_of_stock') }}</div>
                <div class="text-xl font-bold text-red-600">{{ inventoryReport.out_of_stock_count || 0 }}</div>
              </div>
              <div class="col-12 md:col-3">
                <div class="text-sm text-gray-500">{{ $t('total_value') }}</div>
                <div class="text-xl font-bold">{{ inventoryReport.total_value?.toFixed(2) || '0.00' }} EUR</div>
              </div>
            </div>

            <div v-if="inventoryReport.low_stock_items?.length" class="mt-4">
              <h4>{{ $t('low_stock_items') }}</h4>
              <DataTable :value="inventoryReport.low_stock_items" stripedRows>
                <Column field="name" :header="$t('name')"></Column>
                <Column field="quantity" :header="$t('quantity')"></Column>
                <Column field="unit" :header="$t('unit')"></Column>
                <Column field="value" :header="$t('value')">
                  <template #body="slotProps">
                    {{ slotProps.data.value?.toFixed(2) }} EUR
                  </template>
                </Column>
              </DataTable>
            </div>
          </template>
        </Card>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import Card from 'primevue/card'
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import Button from 'primevue/button'
import Calendar from 'primevue/calendar'
import axios from 'axios'
import { ref, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { useToast } from 'primevue/usetoast'

const { t } = useI18n()
const toast = useToast()

const dashboard = ref<any>(null)
const salesReport = ref<any>(null)
const inventoryReport = ref<any>(null)
const startDate = ref<Date | null>(null)
const endDate = ref<Date | null>(null)
const salesLoading = ref(false)
const inventoryLoading = ref(false)

const formatDate = (d: Date) => {
  if (!d) return ''
  return `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, '0')}-${String(d.getDate()).padStart(2, '0')}`
}

const apiBase = `http://${import.meta.env.VITE_APP_BACKEND_HOST}${import.meta.env.VITE_APP_MODULE_CORE_API_PREFIX}/report/api/reports`

const loadDashboard = () => {
  axios.get(`${apiBase}/dashboard`).then((res) => {
    dashboard.value = res.data.data
  }).catch(() => {})
}

const loadSalesReport = () => {
  salesLoading.value = true
  let url = `${apiBase}/sales?`
  if (startDate.value) url += `start_date=${formatDate(startDate.value)}&`
  if (endDate.value) url += `end_date=${formatDate(endDate.value)}`
  axios.get(url).then((res) => {
    salesReport.value = res.data.data
  }).catch((err) => {
    toast.add({ severity: 'error', summary: t('error'), detail: err.message, life: 3000 })
  }).finally(() => { salesLoading.value = false })
}

const loadInventoryReport = () => {
  inventoryLoading.value = true
  axios.get(`${apiBase}/inventory`).then((res) => {
    inventoryReport.value = res.data.data
  }).catch((err) => {
    toast.add({ severity: 'error', summary: t('error'), detail: err.message, life: 3000 })
  }).finally(() => { inventoryLoading.value = false })
}

const exportSalesCSV = () => {
  let url = `${apiBase}/sales/export?`
  if (startDate.value) url += `start_date=${formatDate(startDate.value)}&`
  if (endDate.value) url += `end_date=${formatDate(endDate.value)}`
  window.open(url, '_blank')
}

const exportInventoryCSV = () => {
  window.open(`${apiBase}/inventory/export`, '_blank')
}

onMounted(() => {
  loadDashboard()
  loadInventoryReport()
})
</script>

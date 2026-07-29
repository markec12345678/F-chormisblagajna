<template>
  <div class="p-4">
    <div class="flex justify-content-between align-items-center mb-4">
      <h2 class="m-0">{{ $t('fiscal_dashboard') }}</h2>
      <div class="flex gap-2">
        <Button label="CSV" icon="pi pi-download" severity="success" @click="exportCSV" />
        <Button
          :label="$t('refresh')"
          icon="pi pi-refresh"
          @click="loadReceipts"
          :loading="loading"
        />
      </div>
    </div>

    <div v-if="loading" class="flex justify-content-center p-8">
      <ProgressSpinner style="width: 50px; height: 50px" strokeWidth="6" />
    </div>

    <div v-else>
      <!-- Summary Cards -->
      <div class="grid mb-4">
        <div class="col-12 md:col-3">
          <Card class="h-full">
            <template #title>{{ $t('total_receipts') }}</template>
            <template #content>
              <div class="text-4xl font-bold text-primary">
                {{ dailySummary?.total_count || 0 }}
              </div>
            </template>
          </Card>
        </div>
        <div class="col-12 md:col-3">
          <Card class="h-full">
            <template #title>{{ $t('total_amount') }}</template>
            <template #content>
              <div class="text-4xl font-bold text-green-500">
                {{ formatCurrency(dailySummary?.total_amount || 0) }}
              </div>
            </template>
          </Card>
        </div>
        <div class="col-12 md:col-3">
          <Card class="h-full">
            <template #title>{{ $t('date') }}</template>
            <template #content>
              <div class="text-2xl font-bold">{{ dailySummary?.date || today }}</div>
            </template>
          </Card>
        </div>
      </div>

      <!-- Date Filter -->
      <div class="flex gap-2 mb-4">
        <Calendar v-model="startDate" :placeholder="$t('start_date')" showIcon />
        <Calendar v-model="endDate" :placeholder="$t('end_date')" showIcon />
        <Button :label="$t('filter')" icon="pi pi-filter" @click="loadReceipts" />
      </div>

      <!-- Receipts Table -->
      <Card>
        <template #title>{{ $t('fiscal_receipts') }}</template>
        <template #content>
          <DataTable :value="receipts" responsiveLayout="scroll" stripedRows>
            <Column field="invoice_number" :header="$t('invoice_number')" sortable />
            <Column field="eor" :header="$t('eor')" sortable />
            <Column field="zoi" :header="$t('zoi')" />
            <Column field="invoice_amount" :header="$t('amount')" sortable>
              <template #body="{ data }">{{ formatCurrency(data.invoice_amount) }}</template>
            </Column>
            <Column field="issued_at" :header="$t('date')" sortable>
              <template #body="{ data }">{{ formatDate(data.issued_at) }}</template>
            </Column>
            <Column field="pending_offline" :header="$t('status')">
              <template #body="{ data }">
                <Tag
                  :value="data.pending_offline ? $t('pending') : $t('fiscalized')"
                  :severity="data.pending_offline ? 'warn' : 'success'"
                />
              </template>
            </Column>
          </DataTable>
        </template>
      </Card>
    </div>

    <Toast />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import Card from 'primevue/card'
import Button from 'primevue/button'
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import Tag from 'primevue/tag'
import Calendar from 'primevue/calendar'
import ProgressSpinner from 'primevue/progressspinner'
import Toast from 'primevue/toast'
import { useToast } from 'primevue/usetoast'
import { useI18n } from 'vue-i18n'
import axios from 'axios'

const { t } = useI18n({ useScope: 'global' })
const toast = useToast()

const loading = ref(true)
const receipts = ref([])
const dailySummary = ref(null)
const startDate = ref<Date | null>(null)
const endDate = ref<Date | null>(null)

const today = new Date().toISOString().split('T')[0]

const apiBase = `http://${import.meta.env.VITE_APP_BACKEND_HOST}${import.meta.env.VITE_APP_MODULE_FISCAL_API_PREFIX}/fiscal/api/fiscal`

const formatCurrency = (value: number) => {
  return new Intl.NumberFormat('sl-SI', { style: 'currency', currency: 'EUR' }).format(value)
}

const formatDate = (dateStr: string) => {
  if (!dateStr) return ''
  return new Date(dateStr).toLocaleString()
}

const formatDateParam = (d: Date) => {
  if (!d) return ''
  return `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, '0')}-${String(d.getDate()).padStart(2, '0')}`
}

const loadReceipts = async () => {
  loading.value = true
  try {
    let url = `${apiBase}/receipts?`
    if (startDate.value) url += `start_date=${formatDateParam(startDate.value)}&`
    if (endDate.value) url += `end_date=${formatDateParam(endDate.value)}`

    const [receiptsRes, summaryRes] = await Promise.all([
      axios.get(url),
      axios.get(`${apiBase}/daily-summary?date=${today}`),
    ])
    receipts.value = receiptsRes.data.data || []
    dailySummary.value = summaryRes.data.data
  } catch {
    toast.add({ severity: 'error', summary: t('error'), detail: t('request_failed'), life: 3000 })
  } finally {
    loading.value = false
  }
}

const exportCSV = () => {
  let url = `${apiBase}/receipts/export?`
  if (startDate.value) url += `start_date=${formatDateParam(startDate.value)}&`
  if (endDate.value) url += `end_date=${formatDateParam(endDate.value)}`
  window.open(url, '_blank')
}

onMounted(() => {
  loadReceipts()
})
</script>

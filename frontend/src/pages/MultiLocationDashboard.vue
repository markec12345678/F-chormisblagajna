<template>
  <div class="p-4">
    <div class="flex justify-content-between align-items-center mb-4">
      <h2 class="m-0">{{ $t('multi_location') }}</h2>
      <Button :label="$t('refresh')" icon="pi pi-refresh" @click="loadDashboard" :loading="loading" />
    </div>

    <div v-if="loading" class="flex justify-content-center p-8">
      <ProgressSpinner style="width: 50px; height: 50px" strokeWidth="6" />
    </div>

    <div v-else>
      <!-- Summary Cards -->
      <div class="grid mb-4">
        <div class="col-12 md:col-3">
          <Card class="h-full">
            <template #title>{{ $t('total_branches') }}</template>
            <template #content>
              <div class="text-4xl font-bold text-primary">{{ dashboard?.total_branches || 0 }}</div>
            </template>
          </Card>
        </div>
        <div class="col-12 md:col-3">
          <Card class="h-full">
            <template #title>{{ $t('consolidated_revenue') }}</template>
            <template #content>
              <div class="text-4xl font-bold text-green-500">
                {{ formatCurrency(dashboard?.total_revenue || 0) }}
              </div>
            </template>
          </Card>
        </div>
        <div class="col-12 md:col-3">
          <Card class="h-full">
            <template #title>{{ $t('consolidated_orders') }}</template>
            <template #content>
              <div class="text-4xl font-bold text-blue-500">
                {{ dashboard?.total_orders || 0 }}
              </div>
            </template>
          </Card>
        </div>
        <div class="col-12 md:col-3">
          <Card class="h-full">
            <template #title>{{ $t('average_order') }}</template>
            <template #content>
              <div class="text-4xl font-bold text-orange-500">
                {{ formatCurrency(dashboard?.avg_order_value || 0) }}
              </div>
            </template>
          </Card>
        </div>
      </div>

      <!-- Branch Details Table -->
      <Card class="mb-4">
        <template #title>{{ $t('location_overview') }}</template>
        <template #content>
          <DataTable :value="dashboard?.branches || []" responsiveLayout="scroll">
            <Column field="branch_name" :header="$t('branch_name')" sortable />
            <Column field="today_revenue" :header="$t('today_revenue')" sortable>
              <template #body="{ data }">{{ formatCurrency(data.today_revenue) }}</template>
            </Column>
            <Column field="today_orders" :header="$t('today_orders')" sortable />
            <Column field="avg_order_value" :header="$t('average_order')" sortable>
              <template #body="{ data }">{{ formatCurrency(data.avg_order_value) }}</template>
            </Column>
            <Column field="week_revenue" :header="$t('week_revenue')" sortable>
              <template #body="{ data }">{{ formatCurrency(data.week_revenue) }}</template>
            </Column>
            <Column field="month_revenue" :header="$t('month_revenue')" sortable>
              <template #body="{ data }">{{ formatCurrency(data.month_revenue) }}</template>
            </Column>
            <Column field="status" :header="$t('status')" sortable>
              <template #body="{ data }">
                <Tag :value="data.status" :severity="data.status === 'active' ? 'success' : 'warn'" />
              </template>
            </Column>
          </DataTable>
        </template>
      </Card>

      <!-- Branch Comparison -->
      <Card v-if="comparison.length > 0">
        <template #title>{{ $t('branch_comparison') }}</template>
        <template #content>
          <DataTable :value="comparison" responsiveLayout="scroll">
            <Column field="metric" :header="$t('metric')" sortable />
            <Column
              v-for="branch in dashboard?.branches || []"
              :key="branch.branch_id"
              :field="'branches'"
              :header="branch.branch_name"
            >
              <template #body="{ data }">
                {{
                  formatCurrency(
                    data.branches?.find((b) => b.branch_id === branch.branch_id)?.value || 0,
                  )
                }}
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
import ProgressSpinner from 'primevue/progressspinner'
import Toast from 'primevue/toast'
import { useToast } from 'primevue/usetoast'
import { useI18n } from 'vue-i18n'
import axios from 'axios'

const { t } = useI18n({ useScope: 'global' })
const toast = useToast()

const loading = ref(true)
const dashboard = ref(null)
const comparison = ref([])

const apiBase = `http://${import.meta.env.VITE_APP_BACKEND_HOST}${import.meta.env.VITE_APP_MODULE_CORE_API_PREFIX}/multilocation/api/multilocation`

const formatCurrency = (value: number) => {
  return new Intl.NumberFormat('sl-SI', { style: 'currency', currency: 'EUR' }).format(value)
}

const loadDashboard = async () => {
  loading.value = true
  try {
    const [dashRes, compRes] = await Promise.all([
      axios.get(`${apiBase}/dashboard`),
      axios.get(`${apiBase}/comparison`),
    ])
    dashboard.value = dashRes.data.data
    comparison.value = compRes.data.data
  } catch {
    toast.add({
      severity: 'error',
      summary: t('error'),
      detail: t('request_failed'),
      life: 3000,
    })
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  loadDashboard()
})
</script>

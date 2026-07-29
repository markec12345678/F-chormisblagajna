<template>
  <div class="w-full">
    <div class="grid mx-2">
      <div class="col-12">
        <div class="flex justify-content-between align-items-center">
          <h3>{{ $t('employee_performance') }}</h3>
          <div class="flex gap-2">
            <Calendar v-model="startDate" :placeholder="$t('start_date')" dateFormat="yy-mm-dd" />
            <Calendar v-model="endDate" :placeholder="$t('end_date')" dateFormat="yy-mm-dd" />
            <Button
              :label="$t('load')"
              icon="pi pi-refresh"
              @click="loadPerformance"
              :loading="isLoading"
            />
          </div>
        </div>
      </div>

      <div class="col-12" v-if="summary">
        <div class="grid">
          <div class="col-12 md:col-3">
            <Card>
              <template #title>{{ $t('total_employees') }}</template>
              <template #content>
                <div class="text-4xl font-bold">{{ summary.totalEmployees }}</div>
              </template>
            </Card>
          </div>
          <div class="col-12 md:col-3">
            <Card>
              <template #title>{{ $t('total_revenue') }}</template>
              <template #content>
                <div class="text-4xl font-bold text-green-500">
                  {{ formatCurrency(summary.total_revenue) }}
                </div>
              </template>
            </Card>
          </div>
          <div class="col-12 md:col-3">
            <Card>
              <template #title>{{ $t('avg_sales_per_hour') }}</template>
              <template #content>
                <div class="text-4xl font-bold text-blue-500">
                  {{ formatCurrency(summary.avg_sales_per_hour) }}
                </div>
              </template>
            </Card>
          </div>
          <div class="col-12 md:col-3">
            <Card>
              <template #title>{{ $t('leaderboard') }}</template>
              <template #content>
                <div class="text-4xl font-bold text-yellow-500">
                  {{ summary.top_performers?.[0]?.employee_name || '-' }}
                </div>
              </template>
            </Card>
          </div>
        </div>
      </div>

      <div class="col-12" v-if="summary">
        <Card>
          <template #title>{{ $t('leaderboard') }}</template>
          <template #content>
            <DataTable :value="summary.top_performers" stripedRows>
              <template #empty>{{ $t('no_data') }}</template>
              <Column field="rank" :header="$t('rank')" style="width: 80px">
                <template #body="slotProps">
                  <Tag
                    :value="`#${slotProps?.data?.rank || 0}`"
                    :severity="getRankSeverity(slotProps?.data?.rank || 0)"
                  />
                </template>
              </Column>
              <Column field="employee_name" :header="$t('employee')"></Column>
              <Column field="total_sales" :header="$t('total_sales')">
                <template #body="slotProps">
                  {{ formatCurrency(slotProps?.data?.total_sales || 0) }}
                </template>
              </Column>
              <Column field="order_count" :header="$t('orders')"></Column>
              <Column field="avg_order_value" :header="$t('avg_order')">
                <template #body="slotProps">
                  {{ formatCurrency(slotProps?.data?.avg_order_value || 0) }}
                </template>
              </Column>
              <Column field="sales_per_hour" :header="$t('sales_per_hour')">
                <template #body="slotProps">
                  {{ formatCurrency(slotProps?.data?.sales_per_hour || 0) }}
                </template>
              </Column>
              <Column field="total_tips" :header="$t('tips')">
                <template #body="slotProps">
                  {{ formatCurrency(slotProps?.data?.totalTips || 0) }}
                </template>
              </Column>
            </DataTable>
          </template>
        </Card>
      </div>

      <div class="col-12" v-if="summary">
        <Card>
          <template #title>{{ $t('top_products_per_employee') }}</template>
          <template #content>
            <div class="grid">
              <div
                class="col-12 md:col-6 lg:col-4"
                v-for="emp in summary.top_performers?.slice(0, 6)"
                :key="emp.employee_id"
              >
                <Card>
                  <template #title>{{ emp.employee_name }}</template>
                  <template #content>
                    <DataTable :value="(emp.top_products || []).slice(0, 5)" size="small">
                      <Column field="product_name" :header="$t('product')"></Column>
                      <Column field="quantity" :header="$t('qty')"></Column>
                      <Column field="revenue" :header="$t('revenue')">
                        <template #body="slotProps">
                          {{ formatCurrency(slotProps?.data?.revenue || 0) }}
                        </template>
                      </Column>
                    </DataTable>
                  </template>
                </Card>
              </div>
            </div>
          </template>
        </Card>
      </div>

      <div class="col-12" v-if="!summary && !isLoading">
        <Card>
          <template #content>
            <div class="flex flex-column align-items-center gap-3 py-6">
              <i class="pi pi-users" style="font-size: 3rem; opacity: 0.3"></i>
              <p class="m-0 text-500">{{ $t('select_date_range') }}</p>
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
import Calendar from 'primevue/calendar'
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import Tag from 'primevue/tag'
import axios from 'axios'
import { ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { useToast } from 'primevue/usetoast'
import auth from '../services/auth'

const { t } = useI18n()
const toast = useToast()

interface ProductStat {
  product_id: string
  product_name: string
  quantity: number
  revenue: number
}

interface EmployeePerf {
  employee_id: string
  employee_name: string
  total_sales: number
  order_count: number
  avg_order_value: number
  totalTips: number
  sales_per_hour: number
  top_products: ProductStat[]
  rank: number
}

interface PerformanceSummary {
  totalEmployees: number
  top_performers: EmployeePerf[]
  total_revenue: number
  avg_sales_per_hour: number
}

const startDate = ref<Date | null>(null)
const endDate = ref<Date | null>(null)
const summary = ref<PerformanceSummary | null>(null)
const isLoading = ref(false)

const formatCurrency = (amount: number) => {
  return new Intl.NumberFormat('sl-SI', { style: 'currency', currency: 'EUR' }).format(amount)
}

const formatDate = (date: Date | null): string => {
  if (!date) return ''
  return date.toISOString().split('T')[0]
}

const getRankSeverity = (rank: number) => {
  if (rank === 1) return 'success'
  if (rank === 2) return 'info'
  if (rank === 3) return 'warn'
  return 'secondary'
}

const loadPerformance = async () => {
  isLoading.value = true
  try {
    const params = new URLSearchParams()
    if (startDate.value) params.append('start_date', formatDate(startDate.value))
    if (endDate.value) params.append('end_date', formatDate(endDate.value))

    const response = await axios.get(
      `http://${import.meta.env.VITE_APP_BACKEND_HOST}/employee/api/performance?${params.toString()}`,
      { headers: { Authorization: `Bearer ${auth.accessToken.value}` } },
    )
    summary.value = response.data.data
  } catch {
    toast.add({
      severity: 'error',
      summary: t('failed'),
      detail: t('performance_load_failed'),
      group: 'br',
      life: 3000,
    })
  } finally {
    isLoading.value = false
  }
}
</script>

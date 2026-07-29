<template>
  <div class="w-full">
    <div class="grid mx-2">
      <div class="col-12">
        <div class="flex justify-content-between align-items-center">
          <h3>{{ $t('menu_engineering') }}</h3>
          <div class="flex gap-2">
            <Calendar v-model="startDate" :placeholder="$t('start_date')" dateFormat="yy-mm-dd" />
            <Calendar v-model="endDate" :placeholder="$t('end_date')" dateFormat="yy-mm-dd" />
            <Button :label="$t('analyze')" icon="pi pi-chart-bar" @click="loadAnalysis" :loading="isLoading" />
          </div>
        </div>
      </div>

      <div class="col-12" v-if="summary">
        <div class="grid">
          <div class="col-12 md:col-3">
            <Card>
              <template #title>{{ $t('total_items') }}</template>
              <template #content>
                <div class="text-4xl font-bold">{{ summary.total_items }}</div>
              </template>
            </Card>
          </div>
          <div class="col-12 md:col-3">
            <Card>
              <template #title>{{ $t('total_revenue') }}</template>
              <template #content>
                <div class="text-4xl font-bold text-green-500">{{ formatCurrency(summary.total_revenue) }}</div>
              </template>
            </Card>
          </div>
          <div class="col-12 md:col-3">
            <Card>
              <template #title>{{ $t('avg_profit') }}</template>
              <template #content>
                <div class="text-4xl font-bold text-blue-500">{{ formatCurrency(summary.avg_profit) }}</div>
              </template>
            </Card>
          </div>
          <div class="col-12 md:col-3">
            <Card>
              <template #title>{{ $t('profitability_matrix') }}</template>
              <template #content>
                <div class="grid text-center">
                  <div class="col-6">
                    <Tag :value="`${summary.top_stars?.length || 0} ${$t('stars')}`" severity="success" class="w-full" />
                  </div>
                  <div class="col-6">
                    <Tag :value="`${summary.top_plowhorses?.length || 0} ${$t('plowhorses')}`" severity="warn" class="w-full" />
                  </div>
                  <div class="col-6">
                    <Tag :value="`${summary.top_puzzles?.length || 0} ${$t('puzzles')}`" severity="info" class="w-full" />
                  </div>
                  <div class="col-6">
                    <Tag :value="`${summary.top_dogs?.length || 0} ${$t('dogs')}`" severity="danger" class="w-full" />
                  </div>
                </div>
              </template>
            </Card>
          </div>
        </div>
      </div>

      <div class="col-12" v-if="summary">
        <div class="grid">
          <div class="col-12 md:col-6">
            <Card>
              <template #title>
                <span class="text-green-500">{{ $t('stars') }}</span> - {{ $t('stars_description') }}
              </template>
              <template #content>
                <DataTable :value="summary.top_stars" stripedRows>
                  <template #empty>{{ $t('no_data') }}</template>
                  <Column field="product_name" :header="$t('product')"></Column>
                  <Column field="total_sold" :header="$t('sold')"></Column>
                  <Column field="profit_per_item" :header="$t('profit_per_item')">
                    <template #body="slotProps">
                      {{ formatCurrency(slotProps?.data?.profit_per_item || 0) }}
                    </template>
                  </Column>
                  <Column field="profit_margin" :header="$t('margin')">
                    <template #body="slotProps">
                      {{ (slotProps?.data?.profit_margin || 0).toFixed(1) }}%
                    </template>
                  </Column>
                </DataTable>
              </template>
            </Card>
          </div>
          <div class="col-12 md:col-6">
            <Card>
              <template #title>
                <span class="text-yellow-500">{{ $t('plowhorses') }}</span> - {{ $t('plowhorses_description') }}
              </template>
              <template #content>
                <DataTable :value="summary.top_plowhorses" stripedRows>
                  <template #empty>{{ $t('no_data') }}</template>
                  <Column field="product_name" :header="$t('product')"></Column>
                  <Column field="total_sold" :header="$t('sold')"></Column>
                  <Column field="profit_per_item" :header="$t('profit_per_item')">
                    <template #body="slotProps">
                      {{ formatCurrency(slotProps?.data?.profit_per_item || 0) }}
                    </template>
                  </Column>
                  <Column field="profit_margin" :header="$t('margin')">
                    <template #body="slotProps">
                      {{ (slotProps?.data?.profit_margin || 0).toFixed(1) }}%
                    </template>
                  </Column>
                </DataTable>
              </template>
            </Card>
          </div>
          <div class="col-12 md:col-6">
            <Card>
              <template #title>
                <span class="text-blue-500">{{ $t('puzzles') }}</span> - {{ $t('puzzles_description') }}
              </template>
              <template #content>
                <DataTable :value="summary.top_puzzles" stripedRows>
                  <template #empty>{{ $t('no_data') }}</template>
                  <Column field="product_name" :header="$t('product')"></Column>
                  <Column field="total_sold" :header="$t('sold')"></Column>
                  <Column field="profit_per_item" :header="$t('profit_per_item')">
                    <template #body="slotProps">
                      {{ formatCurrency(slotProps?.data?.profit_per_item || 0) }}
                    </template>
                  </Column>
                  <Column field="profit_margin" :header="$t('margin')">
                    <template #body="slotProps">
                      {{ (slotProps?.data?.profit_margin || 0).toFixed(1) }}%
                    </template>
                  </Column>
                </DataTable>
              </template>
            </Card>
          </div>
          <div class="col-12 md:col-6">
            <Card>
              <template #title>
                <span class="text-red-500">{{ $t('dogs') }}</span> - {{ $t('dogs_description') }}
              </template>
              <template #content>
                <DataTable :value="summary.top_dogs" stripedRows>
                  <template #empty>{{ $t('no_data') }}</template>
                  <Column field="product_name" :header="$t('product')"></Column>
                  <Column field="total_sold" :header="$t('sold')"></Column>
                  <Column field="profit_per_item" :header="$t('profit_per_item')">
                    <template #body="slotProps">
                      {{ formatCurrency(slotProps?.data?.profit_per_item || 0) }}
                    </template>
                  </Column>
                  <Column field="profit_margin" :header="$t('margin')">
                    <template #body="slotProps">
                      {{ (slotProps?.data?.profit_margin || 0).toFixed(1) }}%
                    </template>
                  </Column>
                </DataTable>
              </template>
            </Card>
          </div>
        </div>
      </div>

      <div class="col-12" v-if="!summary && !isLoading">
        <Card>
          <template #content>
            <div class="flex flex-column align-items-center gap-3 py-6">
              <i class="pi pi-chart-bar" style="font-size: 3rem; opacity: 0.3"></i>
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

interface MenuItem {
  product_id: string
  product_name: string
  total_sold: number
  total_revenue: number
  total_cost: number
  profit_margin: number
  profit_per_item: number
  quadrant: string
}

interface MenuSummary {
  total_items: number
  total_revenue: number
  avg_profit: number
  top_stars: MenuItem[]
  top_plowhorses: MenuItem[]
  top_puzzles: MenuItem[]
  top_dogs: MenuItem[]
}

const startDate = ref<Date | null>(null)
const endDate = ref<Date | null>(null)
const summary = ref<MenuSummary | null>(null)
const isLoading = ref(false)

const formatCurrency = (amount: number) => {
  return new Intl.NumberFormat('sl-SI', { style: 'currency', currency: 'EUR' }).format(amount)
}

const formatDate = (date: Date | null): string => {
  if (!date) return ''
  return date.toISOString().split('T')[0]
}

const loadAnalysis = async () => {
  isLoading.value = true
  try {
    const params = new URLSearchParams()
    if (startDate.value) params.append('start_date', formatDate(startDate.value))
    if (endDate.value) params.append('end_date', formatDate(endDate.value))

    const response = await axios.get(
      `http://${import.meta.env.VITE_APP_BACKEND_HOST}/menuengineering/api/analysis?${params.toString()}`,
      { headers: { Authorization: `Bearer ${auth.accessToken.value}` } }
    )
    summary.value = response.data.data
  } catch {
    toast.add({
      severity: 'error',
      summary: t('failed'),
      detail: t('analysis_failed'),
      group: 'br',
      life: 3000,
    })
  } finally {
    isLoading.value = false
  }
}
</script>

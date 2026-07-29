<template>
  <div class="w-full">
    <div class="grid mx-2">
      <div class="col-12">
        <div class="flex justify-content-between align-items-center">
          <h3>{{ $t('multi_payment') }}</h3>
          <div class="flex gap-2">
            <Calendar v-model="startDate" :placeholder="$t('start_date')" dateFormat="yy-mm-dd" />
            <Calendar v-model="endDate" :placeholder="$t('end_date')" dateFormat="yy-mm-dd" />
            <Button
              icon="pi pi-refresh"
              severity="secondary"
              @click="loadData"
              :loading="isLoading"
            />
          </div>
        </div>
      </div>

      <div class="col-12" v-if="dailyPayments.length > 0">
        <Card>
          <template #title>{{ $t('payment_summary') }}</template>
          <template #content>
            <div class="grid">
              <div class="col-12 md:col-2">
                <div class="text-center">
                  <div class="text-4xl font-bold text-green-500">
                    {{ formatCurrency(grandTotal) }}
                  </div>
                  <div class="text-500">{{ $t('grand_total') }}</div>
                </div>
              </div>
              <div class="col-12 md:col-2">
                <div class="text-center">
                  <div class="text-2xl font-bold">{{ totalCash }}</div>
                  <div class="text-500">{{ $t('cash') }}</div>
                </div>
              </div>
              <div class="col-12 md:col-2">
                <div class="text-center">
                  <div class="text-2xl font-bold">{{ totalCard }}</div>
                  <div class="text-500">{{ $t('card') }}</div>
                </div>
              </div>
              <div class="col-12 md:col-2">
                <div class="text-center">
                  <div class="text-2xl font-bold">{{ totalVoucher }}</div>
                  <div class="text-500">{{ $t('voucher') }}</div>
                </div>
              </div>
              <div class="col-12 md:col-2">
                <div class="text-center">
                  <div class="text-2xl font-bold">{{ totalMobile }}</div>
                  <div class="text-500">{{ $t('mobile_pay') }}</div>
                </div>
              </div>
              <div class="col-12 md:col-2">
                <div class="text-center">
                  <div class="text-2xl font-bold">{{ totalGift }}</div>
                  <div class="text-500">{{ $t('gift_cards') }}</div>
                </div>
              </div>
            </div>
          </template>
        </Card>
      </div>

      <div class="col-12">
        <Card>
          <template #title>{{ $t('daily_breakdown') }}</template>
          <template #content>
            <DataTable :value="dailyPayments" stripedRows :loading="isLoading">
              <template #empty>{{ $t('no_payments') }}</template>
              <Column field="date" :header="$t('date')"></Column>
              <Column field="count" :header="$t('transactions')"></Column>
              <Column field="total_cash" :header="$t('cash')">
                <template #body="slotProps">
                  {{ formatCurrency(slotProps?.data?.total_cash || 0) }}
                </template>
              </Column>
              <Column field="total_card" :header="$t('card')">
                <template #body="slotProps">
                  {{ formatCurrency(slotProps?.data?.total_card || 0) }}
                </template>
              </Column>
              <Column field="total_voucher" :header="$t('voucher')">
                <template #body="slotProps">
                  {{ formatCurrency(slotProps?.data?.total_voucher || 0) }}
                </template>
              </Column>
              <Column field="total_mobile" :header="$t('mobile_pay')">
                <template #body="slotProps">
                  {{ formatCurrency(slotProps?.data?.total_mobile || 0) }}
                </template>
              </Column>
              <Column field="grand_total" :header="$t('total')">
                <template #body="slotProps">
                  <span class="font-bold">{{
                    formatCurrency(slotProps?.data?.grand_total || 0)
                  }}</span>
                </template>
              </Column>
            </DataTable>
          </template>
        </Card>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import Card from 'primevue/card'
import Button from 'primevue/button'
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import Calendar from 'primevue/calendar'
import axios from 'axios'
import { ref, computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { useToast } from 'primevue/usetoast'
import auth from '../services/auth'

const { t } = useI18n()
const toast = useToast()

interface DailyPayments {
  date: string
  total_cash: number
  total_card: number
  total_voucher: number
  total_mobile: number
  total_gift_card: number
  grand_total: number
  count: number
}

const dailyPayments = ref<DailyPayments[]>([])
const startDate = ref<Date | null>(null)
const endDate = ref<Date | null>(null)
const isLoading = ref(false)

const formatCurrency = (amount: number) => {
  return new Intl.NumberFormat('sl-SI', { style: 'currency', currency: 'EUR' }).format(amount)
}

const grandTotal = computed(() =>
  dailyPayments.value.reduce((sum, d) => sum + (d.grand_total || 0), 0),
)
const totalCash = computed(() =>
  formatCurrency(dailyPayments.value.reduce((sum, d) => sum + (d.total_cash || 0), 0)),
)
const totalCard = computed(() =>
  formatCurrency(dailyPayments.value.reduce((sum, d) => sum + (d.total_card || 0), 0)),
)
const totalVoucher = computed(() =>
  formatCurrency(dailyPayments.value.reduce((sum, d) => sum + (d.total_voucher || 0), 0)),
)
const totalMobile = computed(() =>
  formatCurrency(dailyPayments.value.reduce((sum, d) => sum + (d.total_mobile || 0), 0)),
)
const totalGift = computed(() =>
  formatCurrency(dailyPayments.value.reduce((sum, d) => sum + (d.total_gift_card || 0), 0)),
)

const loadData = async () => {
  isLoading.value = true
  try {
    const params: Record<string, string> = {}
    if (startDate.value) params.start_date = startDate.value.toISOString().split('T')[0]
    if (endDate.value) params.end_date = endDate.value.toISOString().split('T')[0]

    const response = await axios.get(
      `http://${import.meta.env.VITE_APP_BACKEND_HOST}/multipayment/api/daily`,
      { headers: { Authorization: `Bearer ${auth.accessToken.value}` }, params },
    )
    dailyPayments.value = response.data.data || []
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

loadData()
</script>

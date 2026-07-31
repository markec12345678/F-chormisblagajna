<template>
  <div class="p-4">
    <div class="flex justify-content-between align-items-center mb-4">
      <h2 class="m-0">{{ $t('employee_tips') }}</h2>
      <div class="flex gap-2">
        <Button
          :label="$t('payout_tips')"
          icon="pi pi-wallet"
          severity="warn"
          @click="showPayoutDialog = true"
        />
        <Button :label="$t('record_tip')" icon="pi pi-plus" @click="showRecordDialog = true" />
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
            <template #title>{{ $t('total_tips') }}</template>
            <template #content>
              <div class="text-3xl font-bold text-green-500">{{ formatCurrency(totalTips) }}</div>
            </template>
          </Card>
        </div>
        <div class="col-12 md:col-3">
          <Card class="h-full">
            <template #title>{{ $t('tip_count') }}</template>
            <template #content>
              <div class="text-3xl font-bold">{{ totalCount }}</div>
            </template>
          </Card>
        </div>
        <div class="col-12 md:col-3">
          <Card class="h-full">
            <template #title>{{ $t('average_tip') }}</template>
            <template #content>
              <div class="text-3xl font-bold text-blue-500">{{ formatCurrency(averageTip) }}</div>
            </template>
          </Card>
        </div>
        <div class="col-12 md:col-3">
          <Card class="h-full">
            <template #title>{{ $t('payout_history') }}</template>
            <template #content>
              <div class="text-3xl font-bold text-yellow-500">
                {{ formatCurrency(totalPaidOut) }}
              </div>
            </template>
          </Card>
        </div>
      </div>

      <!-- Tip Summary by Employee -->
      <Card class="mb-4">
        <template #title>{{ $t('tip_summary') }}</template>
        <template #content>
          <DataTable :value="summaries" responsiveLayout="scroll" stripedRows>
            <Column field="employee_name" :header="$t('employee')" sortable />
            <Column field="total_tips" :header="$t('total_tips')" sortable>
              <template #body="{ data }">{{ formatCurrency(data.total_tips) }}</template>
            </Column>
            <Column field="tip_count" :header="$t('tip_count')" sortable />
            <Column field="average_tip" :header="$t('average_tip')" sortable>
              <template #body="{ data }">{{ formatCurrency(data.average_tip) }}</template>
            </Column>
          </DataTable>
        </template>
      </Card>

      <!-- Payout History -->
      <Card>
        <template #title>{{ $t('payout_history') }}</template>
        <template #content>
          <DataTable :value="payouts" responsiveLayout="scroll" stripedRows>
            <Column field="employee_name" :header="$t('employee')" sortable />
            <Column field="amount" :header="$t('amount')" sortable>
              <template #body="{ data }">{{ formatCurrency(data.amount) }}</template>
            </Column>
            <Column field="payout_method" :header="$t('payout_method')" />
            <Column field="payout_date" :header="$t('date')" sortable>
              <template #body="{ data }">{{
                new Date(data.payout_date).toLocaleString()
              }}</template>
            </Column>
          </DataTable>
        </template>
      </Card>
    </div>

    <!-- Record Tip Dialog -->
    <Dialog
      v-model:visible="showRecordDialog"
      :header="$t('record_tip')"
      :style="{ width: '450px' }"
      modal
    >
      <div class="flex flex-column gap-3">
        <div>
          <label class="block mb-1">{{ $t('employee') }}</label>
          <InputText v-model="newTip.employee_name" class="w-full" />
        </div>
        <div>
          <label class="block mb-1">{{ $t('amount') }}</label>
          <InputNumber
            v-model="newTip.amount"
            mode="currency"
            currency="EUR"
            locale="sl-SI"
            class="w-full"
          />
        </div>
        <div>
          <label class="block mb-1">{{ $t('payment_method') }}</label>
          <Dropdown
            v-model="newTip.payment_method"
            :options="paymentMethods"
            optionLabel="label"
            optionValue="value"
            class="w-full"
          />
        </div>
      </div>
      <template #footer>
        <Button :label="$t('cancel')" severity="secondary" @click="showRecordDialog = false" />
        <Button :label="$t('record_tip')" @click="recordTip" :loading="recording" />
      </template>
    </Dialog>

    <!-- Payout Dialog -->
    <Dialog
      v-model:visible="showPayoutDialog"
      :header="$t('payout_tips')"
      :style="{ width: '450px' }"
      modal
    >
      <div class="flex flex-column gap-3">
        <div>
          <label class="block mb-1">{{ $t('employee') }}</label>
          <InputText v-model="payoutData.employee_name" class="w-full" />
        </div>
        <div>
          <label class="block mb-1">{{ $t('amount') }}</label>
          <InputNumber
            v-model="payoutData.amount"
            mode="currency"
            currency="EUR"
            locale="sl-SI"
            class="w-full"
          />
        </div>
        <div>
          <label class="block mb-1">{{ $t('payout_method') }}</label>
          <Dropdown
            v-model="payoutData.payout_method"
            :options="paymentMethods"
            optionLabel="label"
            optionValue="value"
            class="w-full"
          />
        </div>
      </div>
      <template #footer>
        <Button :label="$t('cancel')" severity="secondary" @click="showPayoutDialog = false" />
        <Button :label="$t('payout')" severity="warn" @click="payoutTips" :loading="payingOut" />
      </template>
    </Dialog>

    <Toast />
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import Card from 'primevue/card'
import Button from 'primevue/button'
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import Dialog from 'primevue/dialog'
import InputText from 'primevue/inputtext'
import InputNumber from 'primevue/inputnumber'
import Dropdown from 'primevue/dropdown'
import ProgressSpinner from 'primevue/progressspinner'
import Toast from 'primevue/toast'
import { useToast } from 'primevue/usetoast'
import { useI18n } from 'vue-i18n'
import axios from 'axios'

const { t } = useI18n({ useScope: 'global' })
const toast = useToast()

const loading = ref(true)
const recording = ref(false)
const payingOut = ref(false)
const summaries = ref([])
const payouts = ref([])
const showRecordDialog = ref(false)
const showPayoutDialog = ref(false)

const newTip = ref({ employee_name: '', amount: 0, payment_method: 'cash' })
const payoutData = ref({ employee_name: '', amount: 0, payout_method: 'cash' })

const paymentMethods = [
  { label: 'Cash', value: 'cash' },
  { label: 'Card', value: 'card' },
]

const apiBase = `http://${import.meta.env.VITE_APP_BACKEND_HOST}/tips/api/tips`

const formatCurrency = (value: number) => {
  return new Intl.NumberFormat('sl-SI', { style: 'currency', currency: 'EUR' }).format(value || 0)
}

const totalTips = computed(() => summaries.value.reduce((sum, s) => sum + (s.total_tips || 0), 0))
const totalCount = computed(() => summaries.value.reduce((sum, s) => sum + (s.tip_count || 0), 0))
const averageTip = computed(() => (totalCount.value > 0 ? totalTips.value / totalCount.value : 0))
const totalPaidOut = computed(() => payouts.value.reduce((sum, p) => sum + (p.amount || 0), 0))

const loadData = async () => {
  loading.value = true
  try {
    const [summaryRes, payoutRes] = await Promise.all([
      axios.get(`${apiBase}/summary`),
      axios.get(`${apiBase}/payouts`),
    ])
    summaries.value = summaryRes.data.data || []
    payouts.value = payoutRes.data.data || []
  } catch {
    toast.add({ severity: 'error', summary: t('error'), detail: t('request_failed'), life: 3000 })
  } finally {
    loading.value = false
  }
}

const recordTip = async () => {
  recording.value = true
  try {
    await axios.post(apiBase, newTip.value)
    toast.add({ severity: 'success', summary: t('success'), detail: t('tip_recorded'), life: 3000 })
    showRecordDialog.value = false
    newTip.value = { employee_name: '', amount: 0, payment_method: 'cash' }
    loadData()
  } catch {
    toast.add({ severity: 'error', summary: t('error'), detail: t('request_failed'), life: 3000 })
  } finally {
    recording.value = false
  }
}

const payoutTips = async () => {
  payingOut.value = true
  try {
    await axios.post(`${apiBase}/payout`, payoutData.value)
    toast.add({ severity: 'success', summary: t('success'), detail: t('tip_paid_out'), life: 3000 })
    showPayoutDialog.value = false
    payoutData.value = { employee_name: '', amount: 0, payout_method: 'cash' }
    loadData()
  } catch {
    toast.add({ severity: 'error', summary: t('error'), detail: t('request_failed'), life: 3000 })
  } finally {
    payingOut.value = false
  }
}

onMounted(() => {
  loadData()
})
</script>

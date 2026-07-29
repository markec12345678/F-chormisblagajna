<template>
  <div class="w-full">
    <div class="grid mx-2">
      <div class="col-12">
        <h3>{{ $t('loyalty_program') }}</h3>
      </div>
      <div class="col-12">
        <DataTable
          @page="onPage"
          :lazy="true"
          :totalRecords="totalRecords"
          :loading="loading"
          paginatorPosition="both"
          paginator
          :rows="rowsPerPage"
          :rowsPerPageOptions="[50, 100, 500]"
          :value="accounts"
          stripedRows
          tableStyle="width: 100%;max-height:60vh;"
        >
          <template #header>
            <div class="flex justify-between align-items-center">
              <Button icon="pi pi-plus" :label="$t('add_account')" rounded raised @click="openAdd" />
              <Dropdown v-model="tierFilter" :options="tierOptions" optionLabel="label" optionValue="value" :placeholder="$t('all_tiers')" showClear />
            </div>
          </template>
          <template #empty>
            <div class="flex flex-column align-items-center gap-2 py-4">
              <i class="pi pi-star" style="font-size: 2rem; opacity: 0.3"></i>
              <p class="m-0 text-slate-400">{{ $t('no_results') }}</p>
            </div>
          </template>
          <Column sortable field="customer_id" :header="$t('customer')"></Column>
          <Column sortable field="points" :header="$t('points')">
            <template #body="slotProps">
              <span class="font-bold">{{ slotProps.data.points }}</span>
            </template>
          </Column>
          <Column sortable field="tier" :header="$t('tier')">
            <template #body="slotProps">
              <Tag :value="slotProps.data.tier" :severity="getTierSeverity(slotProps.data.tier)" />
            </template>
          </Column>
          <Column sortable field="total_spent" :header="$t('total_spent')">
            <template #body="slotProps">
              {{ slotProps.data.total_spent?.toFixed(2) }} EUR
            </template>
          </Column>
          <Column :header="$t('actions')">
            <template #body="slotProps">
              <ButtonGroup>
                <Button icon="pi pi-wallet" severity="success" :aria-label="$t('earn_points')" @click="openEarn(slotProps.data)" />
                <Button icon="pi pi-minus-circle" severity="warn" :aria-label="$t('redeem_points')" @click="openRedeem(slotProps.data)" />
                <Button icon="pi pi-history" severity="secondary" :aria-label="$t('history')" @click="viewHistory(slotProps.data)" />
              </ButtonGroup>
            </template>
          </Column>
        </DataTable>
      </div>
    </div>

    <Dialog v-model:visible="earnDialog" modal :header="$t('earn_points')" :style="{ width: '25rem' }">
      <div class="flex flex-column gap-4">
        <div class="flex flex-column gap-2">
          <label>{{ $t('points') }}</label>
          <InputNumber v-model="earnPoints" :min="1" />
        </div>
        <div class="flex flex-column gap-2">
          <label>{{ $t('description') }}</label>
          <InputText v-model="earnDescription" />
        </div>
      </div>
      <template #footer>
        <ButtonGroup>
          <Button :label="$t('cancel')" severity="secondary" @click="earnDialog = false" />
          <Button class="ml-2" severity="success" @click="submitEarn" :label="$t('add_points')" :loading="submitting" />
        </ButtonGroup>
      </template>
    </Dialog>

    <Dialog v-model:visible="redeemDialog" modal :header="$t('redeem_points')" :style="{ width: '25rem' }">
      <div class="flex flex-column gap-4">
        <div class="flex flex-column gap-2">
          <label>{{ $t('points') }}</label>
          <InputNumber v-model="redeemPoints" :min="1" :max="selectedAccount?.points || 0" />
        </div>
      </div>
      <template #footer>
        <ButtonGroup>
          <Button :label="$t('cancel')" severity="secondary" @click="redeemDialog = false" />
          <Button class="ml-2" severity="warn" @click="submitRedeem" :label="$t('redeem')" :loading="submitting" />
        </ButtonGroup>
      </template>
    </Dialog>

    <Dialog v-model:visible="historyDialog" modal :header="$t('transaction_history')" :style="{ width: '40rem' }" :breakpoints="{ '1199px': '90vw', '575px': '90vw' }">
      <DataTable :value="transactions" stripedRows>
        <template #empty>{{ $t('no_results') }}</template>
        <Column field="type" :header="$t('type')">
          <template #body="slotProps">
            <Tag :value="slotProps.data.type" :severity="slotProps.data.type === 'earn' ? 'success' : 'warn'" />
          </template>
        </Column>
        <Column field="points" :header="$t('points')"></Column>
        <Column field="description" :header="$t('description')"></Column>
        <Column field="created_at" :header="$t('date')"></Column>
      </DataTable>
    </Dialog>
  </div>
</template>

<script setup lang="ts">
import DataTable from 'primevue/datatable'
import Dialog from 'primevue/dialog'
import InputText from 'primevue/inputtext'
import InputNumber from 'primevue/inputnumber'
import Column from 'primevue/column'
import Button from 'primevue/button'
import ButtonGroup from 'primevue/buttongroup'
import Tag from 'primevue/tag'
import Dropdown from 'primevue/dropdown'
import axios from 'axios'
import { ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { useToast } from 'primevue/usetoast'

const { t } = useI18n()
const toast = useToast()

const accounts = ref<any[]>([])
const totalRecords = ref(0)
const loading = ref(false)
const rowsPerPage = ref(50)
const submitting = ref(false)
const tierFilter = ref('')
const earnDialog = ref(false)
const redeemDialog = ref(false)
const historyDialog = ref(false)
const selectedAccount = ref<any>(null)
const earnPoints = ref(10)
const earnDescription = ref('')
const redeemPoints = ref(10)
const transactions = ref<any[]>([])

const tierOptions = [
  { label: 'Bronze', value: 'bronze' },
  { label: 'Silver', value: 'silver' },
  { label: 'Gold', value: 'gold' },
  { label: 'Platinum', value: 'platinum' },
]

const getTierSeverity = (tier: string) => {
  switch (tier) {
    case 'platinum': return 'success'
    case 'gold': return 'warn'
    case 'silver': return 'info'
    default: return 'secondary'
  }
}

const apiBase = `http://${import.meta.env.VITE_APP_BACKEND_HOST}${import.meta.env.VITE_APP_MODULE_CORE_API_PREFIX}/loyalty/api/loyalty`

const getAccounts = (offset = 0, limit = 50) => {
  loading.value = true
  const pageNumber = Math.floor(offset / limit) + 1
  let url = `${apiBase}/accounts?page_number=${pageNumber}&page_size=${limit}`
  if (tierFilter.value) url += `&tier=${tierFilter.value}`
  axios.get(url).then((res) => {
    accounts.value = res.data.data || []
    totalRecords.value = res.data.meta?.total_records || 0
  }).catch((err) => {
    toast.add({ severity: 'error', summary: t('error'), detail: err.message, life: 3000 })
  }).finally(() => { loading.value = false })
}

const onPage = (event: { first: number; rows: number }) => { getAccounts(event.first, event.rows) }
const openAdd = () => { /* TODO: create account dialog */ }
const openEarn = (account: any) => { selectedAccount.value = account; earnPoints.value = 10; earnDescription.value = ''; earnDialog.value = true }
const openRedeem = (account: any) => { selectedAccount.value = account; redeemPoints.value = 10; redeemDialog.value = true }

const submitEarn = () => {
  if (!selectedAccount.value) return
  submitting.value = true
  axios.post(`${apiBase}/earn`, { customer_id: selectedAccount.value.customer_id, points: earnPoints.value, description: earnDescription.value }).then(() => {
    toast.add({ severity: 'success', summary: t('success'), detail: t('points_added'), life: 3000 })
    earnDialog.value = false; getAccounts(0, rowsPerPage.value)
  }).catch((err) => {
    toast.add({ severity: 'error', summary: t('error'), detail: err.message, life: 3000 })
  }).finally(() => { submitting.value = false })
}

const submitRedeem = () => {
  if (!selectedAccount.value) return
  submitting.value = true
  axios.post(`${apiBase}/redeem`, { customer_id: selectedAccount.value.customer_id, points: redeemPoints.value }).then(() => {
    toast.add({ severity: 'success', summary: t('success'), detail: t('points_redeemed'), life: 3000 })
    redeemDialog.value = false; getAccounts(0, rowsPerPage.value)
  }).catch((err) => {
    toast.add({ severity: 'error', summary: t('error'), detail: err.message, life: 3000 })
  }).finally(() => { submitting.value = false })
}

const viewHistory = (account: any) => {
  selectedAccount.value = account
  axios.get(`${apiBase}/transactions/${account.customer_id}`).then((res) => {
    transactions.value = res.data.data || []
    historyDialog.value = true
  }).catch((err) => {
    toast.add({ severity: 'error', summary: t('error'), detail: err.message, life: 3000 })
  })
}

getAccounts()
</script>

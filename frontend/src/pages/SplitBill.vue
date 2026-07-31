<template>
  <div class="w-full">
    <div class="grid mx-2">
      <div class="col-12 flex">
        <div class="grid w-full">
          <div class="col-12">
            <h3>{{ $t('split_bill') }}</h3>
          </div>
          <div class="col-12">
            <DataTable
              @page="updateSplitBillsTableRowsPerPage"
              :lazy="true"
              :totalRecords="splitBillsTableTotalRecords"
              :loading="isSplitBillsTableLoading"
              paginatorPosition="both"
              paginator
              :rows="splitBillsTableRowsPerPage"
              :rowsPerPageOptions="[50, 100, 500]"
              :value="splitBills"
              stripedRows
              tableStyle="width: 100%;max-height:50vh;"
              class="w-full pr-2"
            >
              <template #header>
                <div class="flex justify-between align-items-center">
                  <Button
                    icon="pi pi-plus"
                    :label="$t('create_split_bill')"
                    rounded
                    raised
                    @click="openCreateDialog"
                  />
                </div>
              </template>
              <template #empty>
                <div class="flex flex-column align-items-center gap-2 py-4">
                  <i class="pi pi-wallet" style="font-size: 2rem; opacity: 0.3"></i>
                  <p class="m-0 text-slate-400">{{ $t('no_results') }}</p>
                </div>
              </template>
              <Column sortable field="order_id" :header="$t('order_id')"></Column>
              <Column sortable field="split_type" :header="$t('split_type')">
                <template #body="slotProps">
                  <Tag
                    :value="slotProps.data.split_type"
                    :severity="getSplitTypeSeverity(slotProps.data.split_type)"
                  />
                </template>
              </Column>
              <Column sortable field="status" :header="$t('status')">
                <template #body="slotProps">
                  <Tag
                    :value="slotProps.data.status"
                    :severity="getStatusSeverity(slotProps.data.status)"
                  />
                </template>
              </Column>
              <Column :header="$t('parts')">
                <template #body="slotProps">
                  {{ slotProps.data.parts?.length || 0 }}
                </template>
              </Column>
              <Column :header="$t('paid')">
                <template #body="slotProps">
                  {{ slotProps.data.parts?.filter((p: any) => p.is_paid).length || 0 }} /
                  {{ slotProps.data.parts?.length || 0 }}
                </template>
              </Column>
              <Column :header="$t('actions')">
                <template #body="slotProps">
                  <ButtonGroup>
                    <Button
                      icon="pi pi-eye"
                      severity="secondary"
                      :aria-label="$t('view')"
                      @click="viewSplitBill(slotProps.data)"
                    />
                    <Button
                      v-if="slotProps.data.status !== 'paid'"
                      icon="pi pi-check"
                      severity="success"
                      :aria-label="$t('pay')"
                      @click="openPayDialog(slotProps.data)"
                    />
                  </ButtonGroup>
                </template>
              </Column>
            </DataTable>
          </div>
        </div>
      </div>
    </div>

    <Dialog
      v-model:visible="createDialogVisible"
      modal
      :header="$t('create_split_bill')"
      :style="{ width: '30rem', direction: store.orientation == 'rtl' ? 'rtl' : 'ltr' }"
      :breakpoints="{ '1199px': '90vw', '575px': '90vw' }"
    >
      <div class="flex flex-column gap-4">
        <div class="flex flex-column gap-2">
          <label for="order_id">{{ $t('order_id') }}</label>
          <InputText
            id="order_id"
            v-model="newSplitBill.order_id"
            :class="{ 'p-invalid': createErrors.order_id }"
          />
          <small class="p-error">{{ createErrors.order_id }}</small>
        </div>
        <div class="flex flex-column gap-2">
          <label for="split_type">{{ $t('split_type') }}</label>
          <Dropdown
            id="split_type"
            v-model="newSplitBill.split_type"
            :options="splitTypeOptions"
            optionLabel="label"
            optionValue="value"
          />
        </div>
        <div v-if="newSplitBill.split_type === 'equal'" class="flex flex-column gap-2">
          <label for="split_count">{{ $t('split_count') }}</label>
          <InputNumber id="split_count" v-model="newSplitBill.split_count" :min="2" :max="20" />
        </div>
      </div>
      <template #footer>
        <ButtonGroup>
          <Button :label="$t('cancel')" severity="secondary" @click="createDialogVisible = false" />
          <Button
            class="ml-2"
            severity="primary"
            @click="submitSplitBill"
            :label="$t('create')"
            :loading="isSubmitting"
            :disabled="isSubmitting"
          />
        </ButtonGroup>
      </template>
    </Dialog>

    <Dialog
      v-model:visible="viewDialogVisible"
      modal
      :header="$t('split_bill_details')"
      :style="{ width: '40rem', direction: store.orientation == 'rtl' ? 'rtl' : 'ltr' }"
      :breakpoints="{ '1199px': '90vw', '575px': '90vw' }"
    >
      <div v-if="selectedSplitBill" class="flex flex-column gap-4">
        <div class="grid">
          <div class="col-6">
            <strong>{{ $t('order_id') }}:</strong> {{ selectedSplitBill.order_id }}
          </div>
          <div class="col-6">
            <strong>{{ $t('split_type') }}:</strong>
            <Tag
              :value="selectedSplitBill.split_type"
              :severity="getSplitTypeSeverity(selectedSplitBill.split_type)"
            />
          </div>
        </div>
        <div class="grid">
          <div class="col-6">
            <strong>{{ $t('status') }}:</strong>
            <Tag
              :value="selectedSplitBill.status"
              :severity="getStatusSeverity(selectedSplitBill.status)"
            />
          </div>
        </div>
        <DataTable :value="selectedSplitBill.parts || []" stripedRows>
          <Column field="id" :header="$t('part')">
            <template #body="slotProps"> #{{ slotProps.index + 1 }} </template>
          </Column>
          <Column field="amount" :header="$t('amount')">
            <template #body="slotProps"> {{ slotProps.data.amount?.toFixed(2) }} EUR </template>
          </Column>
          <Column field="payment_method" :header="$t('payment_method')">
            <template #body="slotProps">
              {{ slotProps.data.payment_method || '-' }}
            </template>
          </Column>
          <Column field="is_paid" :header="$t('status')">
            <template #body="slotProps">
              <Tag
                :value="slotProps.data.is_paid ? $t('paid') : $t('pending')"
                :severity="slotProps.data.is_paid ? 'success' : 'warn'"
              />
            </template>
          </Column>
        </DataTable>
      </div>
    </Dialog>

    <Dialog
      v-model:visible="payDialogVisible"
      modal
      :header="$t('pay_split_part')"
      :style="{ width: '25rem', direction: store.orientation == 'rtl' ? 'rtl' : 'ltr' }"
      :breakpoints="{ '1199px': '90vw', '575px': '90vw' }"
    >
      <div v-if="selectedSplitBill" class="flex flex-column gap-4">
        <div class="flex flex-column gap-2">
          <label>{{ $t('select_part') }}</label>
          <Dropdown
            v-model="payPartId"
            :options="unpaidParts"
            optionLabel="label"
            optionValue="value"
            :placeholder="$t('select_part')"
          />
        </div>
        <div class="flex flex-column gap-2">
          <label>{{ $t('payment_method') }}</label>
          <Dropdown
            v-model="payMethod"
            :options="paymentMethodOptions"
            optionLabel="label"
            optionValue="value"
          />
        </div>
      </div>
      <template #footer>
        <ButtonGroup>
          <Button :label="$t('cancel')" severity="secondary" @click="payDialogVisible = false" />
          <Button
            class="ml-2"
            severity="success"
            @click="submitPayPart"
            :label="$t('pay')"
            :loading="isSubmitting"
            :disabled="isSubmitting || !payPartId"
          />
        </ButtonGroup>
      </template>
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
import { ref, computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { useToast } from 'primevue/usetoast'
import { globalStore } from '@/stores'

const { t } = useI18n()
const store = globalStore()
const toast = useToast()

interface SplitPart {
  id: string
  amount: number
  payment_method: string
  is_paid: boolean
  paid_at?: string
}

interface SplitBillItem {
  id: string
  order_id: string
  split_type: string
  parts: SplitPart[]
  status: string
  created_at: string
}

const splitBills = ref<SplitBillItem[]>([])
const splitBillsTableTotalRecords = ref(0)
const isSplitBillsTableLoading = ref(false)
const splitBillsTableRowsPerPage = ref(50)
const isSubmitting = ref(false)
const createDialogVisible = ref(false)
const viewDialogVisible = ref(false)
const payDialogVisible = ref(false)

const newSplitBill = ref({
  order_id: '',
  split_type: 'equal',
  split_count: 2,
})
const createErrors = ref({ order_id: '' })

const selectedSplitBill = ref<SplitBillItem | null>(null)
const payPartId = ref('')
const payMethod = ref('cash')

const splitTypeOptions = computed(() => [
  { label: t('equal'), value: 'equal' },
  { label: t('custom'), value: 'custom' },
  { label: t('by_item'), value: 'by_item' },
])

const paymentMethodOptions = computed(() => [
  { label: t('cash'), value: 'cash' },
  { label: t('card'), value: 'card' },
  { label: t('voucher'), value: 'voucher' },
  { label: t('mobile_pay'), value: 'mobile_pay' },
])

const unpaidParts = computed(() => {
  if (!selectedSplitBill.value) return []
  return selectedSplitBill.value.parts
    .filter((p) => !p.is_paid)
    .map((p, i) => ({
      label: `${t('part')} #${i + 1} - ${p.amount.toFixed(2)} EUR`,
      value: p.id,
    }))
})

const openCreateDialog = () => {
  newSplitBill.value = { order_id: '', split_type: 'equal', split_count: 2 }
  createErrors.value = { order_id: '' }
  createDialogVisible.value = true
}

const viewSplitBill = (bill: SplitBillItem) => {
  selectedSplitBill.value = bill
  viewDialogVisible.value = true
}

const openPayDialog = (bill: SplitBillItem) => {
  selectedSplitBill.value = bill
  payPartId.value = ''
  payMethod.value = 'cash'
  payDialogVisible.value = true
}

const updateSplitBillsTableRowsPerPage = (event: { first: number; rows: number }) => {
  getSplitBills(event.first, event.rows)
}

const getSplitTypeSeverity = (type: string) => {
  switch (type) {
    case 'equal':
      return 'info'
    case 'custom':
      return 'warn'
    case 'by_item':
      return 'success'
    default:
      return 'secondary'
  }
}

const getStatusSeverity = (status: string) => {
  switch (status) {
    case 'paid':
      return 'success'
    case 'partial':
      return 'warn'
    case 'pending':
      return 'secondary'
    default:
      return 'info'
  }
}

const submitSplitBill = () => {
  createErrors.value.order_id = newSplitBill.value.order_id?.trim() ? '' : t('validation_required')
  if (createErrors.value.order_id) return

  isSubmitting.value = true
  axios
    .post(
      `http://${import.meta.env.VITE_APP_BACKEND_HOST}/splitbill/api/split-bills`,
      newSplitBill.value,
    )
    .then(() => {
      toast.add({
        severity: 'success',
        summary: t('success'),
        detail: t('split_bill_created'),
        life: 3000,
      })
      createDialogVisible.value = false
      getSplitBills(0, splitBillsTableRowsPerPage.value)
    })
    .catch((err) => {
      toast.add({ severity: 'error', summary: t('error'), detail: err.message, life: 3000 })
    })
    .finally(() => {
      isSubmitting.value = false
    })
}

const submitPayPart = () => {
  if (!payPartId.value || !selectedSplitBill.value) return

  isSubmitting.value = true
  axios
    .post(
      `http://${import.meta.env.VITE_APP_BACKEND_HOST}/splitbill/api/split-bills/${selectedSplitBill.value.id}/pay`,
      { part_id: payPartId.value, payment_method: payMethod.value, amount: 0 },
    )
    .then(() => {
      toast.add({
        severity: 'success',
        summary: t('success'),
        detail: t('payment_recorded'),
        life: 3000,
      })
      payDialogVisible.value = false
      getSplitBills(0, splitBillsTableRowsPerPage.value)
    })
    .catch((err) => {
      toast.add({ severity: 'error', summary: t('error'), detail: err.message, life: 3000 })
    })
    .finally(() => {
      isSubmitting.value = false
    })
}

const getSplitBills = (offset = 0, limit = 50) => {
  isSplitBillsTableLoading.value = true
  const pageNumber = Math.floor(offset / limit) + 1
  axios
    .get(
      `http://${import.meta.env.VITE_APP_BACKEND_HOST}/splitbill/api/split-bills?page_number=${pageNumber}&page_size=${limit}`,
    )
    .then((res) => {
      splitBills.value = res.data.data || []
      splitBillsTableTotalRecords.value = res.data.meta?.total_records || 0
    })
    .catch((err) => {
      toast.add({ severity: 'error', summary: t('error'), detail: err.message, life: 3000 })
    })
    .finally(() => {
      isSplitBillsTableLoading.value = false
    })
}

getSplitBills()
</script>

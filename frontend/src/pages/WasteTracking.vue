<template>
  <div class="w-full">
    <div class="grid mx-2">
      <div class="col-12">
        <div class="flex justify-content-between align-items-center">
          <h3>{{ $t('waste_tracking') }}</h3>
          <div class="flex gap-2">
            <Calendar
              v-model="summaryStartDate"
              :placeholder="$t('start_date')"
              dateFormat="yy-mm-dd"
            />
            <Calendar
              v-model="summaryEndDate"
              :placeholder="$t('end_date')"
              dateFormat="yy-mm-dd"
            />
            <Button :label="$t('add_waste_entry')" icon="pi pi-plus" @click="showAddDialog" />
          </div>
        </div>
      </div>

      <div class="col-12" v-if="summary">
        <div class="grid">
          <div class="col-12 md:col-4">
            <Card>
              <template #title>{{ $t('total_waste_cost') }}</template>
              <template #content>
                <div class="text-4xl font-bold text-red-500">
                  {{ formatCurrency(summary.total_waste_cost) }}
                </div>
              </template>
            </Card>
          </div>
          <div class="col-12 md:col-4">
            <Card>
              <template #title>{{ $t('total_entries') }}</template>
              <template #content>
                <div class="text-4xl font-bold text-orange-500">{{ summary.total_entries }}</div>
              </template>
            </Card>
          </div>
          <div class="col-12 md:col-4">
            <Card>
              <template #title>{{ $t('by_reason') }}</template>
              <template #content>
                <div class="text-4xl font-bold text-blue-500">
                  {{ summary.by_reason?.length || 0 }}
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
              <template #title>{{ $t('by_reason') }}</template>
              <template #content>
                <DataTable :value="summary.by_reason" stripedRows>
                  <template #empty>{{ $t('no_data') }}</template>
                  <Column field="reason" :header="$t('reason')"></Column>
                  <Column field="total" :header="$t('cost')">
                    <template #body="slotProps">
                      {{ formatCurrency(slotProps?.data?.total || 0) }}
                    </template>
                  </Column>
                  <Column field="count" :header="$t('count')"></Column>
                </DataTable>
              </template>
            </Card>
          </div>
          <div class="col-12 md:col-6">
            <Card>
              <template #title>{{ $t('by_material') }}</template>
              <template #content>
                <DataTable :value="summary.by_material" stripedRows>
                  <template #empty>{{ $t('no_data') }}</template>
                  <Column field="material_name" :header="$t('material')"></Column>
                  <Column field="total_cost" :header="$t('cost')">
                    <template #body="slotProps">
                      {{ formatCurrency(slotProps?.data?.total_cost || 0) }}
                    </template>
                  </Column>
                  <Column field="total_qty" :header="$t('quantity')">
                    <template #body="slotProps">
                      {{ (slotProps?.data?.total_qty || 0).toFixed(1) }}
                    </template>
                  </Column>
                  <Column field="count" :header="$t('count')"></Column>
                </DataTable>
              </template>
            </Card>
          </div>
        </div>
      </div>

      <div class="col-12">
        <Card>
          <template #title>{{ $t('waste_entries') }}</template>
          <template #content>
            <DataTable :value="entries" stripedRows :loading="isLoading">
              <template #empty>{{ $t('no_waste_entries') }}</template>
              <Column field="date" :header="$t('date')">
                <template #body="slotProps">
                  {{ formatDate(slotProps?.data?.date) }}
                </template>
              </Column>
              <Column field="material_name" :header="$t('material')"></Column>
              <Column field="quantity" :header="$t('quantity')">
                <template #body="slotProps">
                  {{ (slotProps?.data?.quantity || 0).toFixed(1) }}
                  {{ slotProps?.data?.unit || '' }}
                </template>
              </Column>
              <Column field="reason" :header="$t('reason')">
                <template #body="slotProps">
                  <Tag
                    :value="slotProps?.data?.reason"
                    :severity="getReasonSeverity(slotProps?.data?.reason)"
                  />
                </template>
              </Column>
              <Column field="cost" :header="$t('cost')">
                <template #body="slotProps">
                  {{ formatCurrency(slotProps?.data?.cost || 0) }}
                </template>
              </Column>
              <Column field="recorded_by" :header="$t('recorded_by')"></Column>
              <Column :header="$t('actions')">
                <template #body="slotProps">
                  <Button
                    icon="pi pi-trash"
                    severity="danger"
                    text
                    rounded
                    @click="deleteEntry(slotProps?.data?.id)"
                  />
                </template>
              </Column>
            </DataTable>
          </template>
        </Card>
      </div>

      <Dialog
        v-model:visible="addDialog"
        modal
        :header="$t('add_waste_entry')"
        :style="{ width: '50rem' }"
      >
        <div class="grid">
          <div class="col-12 md:col-6">
            <div class="flex flex-column gap-2">
              <label>{{ $t('material_name') }}</label>
              <InputText v-model="newEntry.material_name" />
            </div>
          </div>
          <div class="col-12 md:col-6">
            <div class="flex flex-column gap-2">
              <label>{{ $t('material_id') }}</label>
              <InputText v-model="newEntry.material_id" />
            </div>
          </div>
          <div class="col-12 md:col-6">
            <div class="flex flex-column gap-2">
              <label>{{ $t('quantity') }}</label>
              <InputNumber v-model="newEntry.quantity" :minFractionDigits="1" />
            </div>
          </div>
          <div class="col-12 md:col-6">
            <div class="flex flex-column gap-2">
              <label>{{ $t('unit') }}</label>
              <Dropdown v-model="newEntry.unit" :options="units" :placeholder="$t('unit')" />
            </div>
          </div>
          <div class="col-12 md:col-6">
            <div class="flex flex-column gap-2">
              <label>{{ $t('reason') }}</label>
              <Dropdown
                v-model="newEntry.reason"
                :options="reasons"
                :placeholder="$t('select_reason')"
              />
            </div>
          </div>
          <div class="col-12 md:col-6">
            <div class="flex flex-column gap-2">
              <label>{{ $t('cost') }}</label>
              <InputNumber v-model="newEntry.cost" mode="currency" currency="EUR" locale="sl-SI" />
            </div>
          </div>
          <div class="col-12 md:col-6">
            <div class="flex flex-column gap-2">
              <label>{{ $t('recorded_by') }}</label>
              <InputText v-model="newEntry.recorded_by" />
            </div>
          </div>
          <div class="col-12">
            <div class="flex flex-column gap-2">
              <label>{{ $t('notes') }}</label>
              <Textarea v-model="newEntry.notes" rows="3" />
            </div>
          </div>
        </div>
        <template #footer>
          <Button :label="$t('cancel')" severity="secondary" @click="addDialog = false" />
          <Button :label="$t('save')" @click="saveEntry" :loading="isSaving" />
        </template>
      </Dialog>
    </div>
  </div>
</template>

<script setup lang="ts">
import Card from 'primevue/card'
import Button from 'primevue/button'
import Calendar from 'primevue/calendar'
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import Dialog from 'primevue/dialog'
import InputText from 'primevue/inputtext'
import InputNumber from 'primevue/inputnumber'
import Textarea from 'primevue/textarea'
import Dropdown from 'primevue/dropdown'
import Tag from 'primevue/tag'
import axios from 'axios'
import { ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { useToast } from 'primevue/usetoast'
import auth from '../services/auth'

const { t } = useI18n()
const toast = useToast()

interface WasteEntry {
  id: string
  material_id: string
  material_name: string
  quantity: number
  unit: string
  reason: string
  cost: number
  date: string
  recorded_by: string
  notes: string
}

interface ReasonSummary {
  reason: string
  total: number
  count: number
}

interface MaterialSummary {
  material_id: string
  material_name: string
  total_cost: number
  total_qty: number
  count: number
}

interface WasteSummary {
  total_waste_cost: number
  total_entries: number
  by_reason: ReasonSummary[]
  by_material: MaterialSummary[]
}

const entries = ref<WasteEntry[]>([])
const summary = ref<WasteSummary | null>(null)
const summaryStartDate = ref<Date | null>(null)
const summaryEndDate = ref<Date | null>(null)
const isLoading = ref(false)
const isSaving = ref(false)
const addDialog = ref(false)

const newEntry = ref<Partial<WasteEntry>>({})

const reasons = ['expired', 'damaged', 'overcooked', 'other']
const units = ['kg', 'g', 'l', 'ml', 'pcs', 'portion']

const formatCurrency = (amount: number) => {
  return new Intl.NumberFormat('sl-SI', { style: 'currency', currency: 'EUR' }).format(amount)
}

const formatDate = (date: string) => {
  if (!date) return '-'
  return new Date(date).toLocaleDateString()
}

const getReasonSeverity = (reason: string) => {
  switch (reason) {
    case 'expired':
      return 'danger'
    case 'damaged':
      return 'warn'
    case 'overcooked':
      return 'info'
    default:
      return 'secondary'
  }
}

const showAddDialog = () => {
  newEntry.value = {}
  addDialog.value = true
}

const saveEntry = async () => {
  isSaving.value = true
  try {
    await axios.post(
      `http://${import.meta.env.VITE_APP_BACKEND_HOST}/waste/api/waste`,
      newEntry.value,
      { headers: { Authorization: `Bearer ${auth.accessToken.value}` } },
    )
    toast.add({
      severity: 'success',
      summary: t('success'),
      detail: t('waste_entry_saved'),
      group: 'br',
      life: 3000,
    })
    addDialog.value = false
    loadData()
  } catch {
    toast.add({
      severity: 'error',
      summary: t('failed'),
      detail: t('save_failed'),
      group: 'br',
      life: 3000,
    })
  } finally {
    isSaving.value = false
  }
}

const deleteEntry = async (id: string) => {
  try {
    await axios.delete(`http://${import.meta.env.VITE_APP_BACKEND_HOST}/waste/api/waste/${id}`, {
      headers: { Authorization: `Bearer ${auth.accessToken.value}` },
    })
    toast.add({
      severity: 'success',
      summary: t('success'),
      detail: t('waste_entry_deleted'),
      group: 'br',
      life: 3000,
    })
    loadData()
  } catch {
    toast.add({
      severity: 'error',
      summary: t('failed'),
      detail: t('delete_failed'),
      group: 'br',
      life: 3000,
    })
  }
}

const loadData = async () => {
  isLoading.value = true
  try {
    const params: Record<string, string> = {}
    if (summaryStartDate.value) {
      params.start_date = summaryStartDate.value.toISOString().split('T')[0]
    }
    if (summaryEndDate.value) {
      params.end_date = summaryEndDate.value.toISOString().split('T')[0]
    }

    const [entriesRes, summaryRes] = await Promise.all([
      axios.get(`http://${import.meta.env.VITE_APP_BACKEND_HOST}/waste/api/waste`, {
        headers: { Authorization: `Bearer ${auth.accessToken.value}` },
      }),
      axios.get(`http://${import.meta.env.VITE_APP_BACKEND_HOST}/waste/api/waste/summary`, {
        headers: { Authorization: `Bearer ${auth.accessToken.value}` },
        params,
      }),
    ])
    entries.value = entriesRes.data.data || []
    summary.value = summaryRes.data.data
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

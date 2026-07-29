<template>
  <div class="w-full">
    <div class="grid mx-2">
      <div class="col-12 flex">
        <div class="grid w-full">
          <div class="col-12">
            <h3>{{ $t('tables', 3) }}</h3>
          </div>
          <div class="col-12">
            <DataTable
              @page="updateTablesTableRowsPerPage"
              :lazy="true"
              :totalRecords="tablesTableTotalRecords"
              :loading="isTablesTableLoading"
              paginatorPosition="both"
              paginator
              :rows="tablesTableRowsPerPage"
              :rowsPerPageOptions="[50, 100, 500]"
              :value="filteredTables"
              stripedRows
              tableStyle="width: 100%;max-height:50vh;"
              class="w-full pr-2"
            >
              <template #header>
                <div class="flex justify-between align-items-center">
                  <Button
                    icon="pi pi-plus"
                    :label="$t('add_table')"
                    rounded
                    raised
                    @click="openAddDialog"
                  />
                  <div class="flex gap-2">
                    <Dropdown
                      v-model="zoneFilter"
                      :options="zoneOptions"
                      optionLabel="label"
                      optionValue="value"
                      :placeholder="$t('table_zone')"
                      class="w-48"
                      showClear
                    />
                    <Dropdown
                      v-model="statusFilter"
                      :options="statusOptions"
                      optionLabel="label"
                      optionValue="value"
                      :placeholder="$t('table_status')"
                      class="w-48"
                      showClear
                    />
                  </div>
                </div>
              </template>
              <template #empty>
                <div class="flex flex-column align-items-center gap-2 py-4">
                  <i class="pi pi-table" style="font-size: 2rem; opacity: 0.3"></i>
                  <p class="m-0 text-slate-400">{{ $t('no_results') }}</p>
                </div>
              </template>
              <Column sortable field="number" :header="$t('table_number')"></Column>
              <Column sortable field="name" :header="$t('table_name')"></Column>
              <Column sortable field="capacity" :header="$t('table_capacity')"></Column>
              <Column sortable field="zone" :header="$t('table_zone')">
                <template #body="slotProps">
                  <Tag
                    :value="slotProps.data.zone"
                    :severity="getZoneSeverity(slotProps.data.zone)"
                  />
                </template>
              </Column>
              <Column sortable field="status" :header="$t('table_status')">
                <template #body="slotProps">
                  <Tag
                    :value="slotProps.data.status"
                    :severity="getStatusSeverity(slotProps.data.status)"
                  />
                </template>
              </Column>
              <Column :header="$t('qr_code')">
                <template #body="slotProps">
                  <Button
                    icon="pi pi-qrcode"
                    severity="secondary"
                    :aria-label="$t('qr_code')"
                    @click="showQRCode(slotProps.data)"
                  />
                </template>
              </Column>
              <Column :header="$t('actions')">
                <template #body="slotProps">
                  <ConfirmPopup></ConfirmPopup>
                  <ButtonGroup>
                    <Button
                      icon="pi pi-pencil"
                      severity="secondary"
                      :aria-label="$t('edit')"
                      @click="prepareTableToEdit(slotProps.data)"
                    />
                    <Button
                      icon="pi pi-trash"
                      severity="danger"
                      :aria-label="$t('remove')"
                      @click="confirmDeleteTable($event, slotProps.data.id)"
                    />
                  </ButtonGroup>
                </template>
              </Column>
            </DataTable>

            <Dialog
              v-model:visible="tableAddDialog"
              modal
              :header="$t('add_table')"
              :style="{ width: '45rem', direction: store.orientation == 'rtl' ? 'rtl' : 'ltr' }"
              :breakpoints="{ '1199px': '90vw', '575px': '90vw' }"
            >
              <div class="flex flex-column gap-2 w-full">
                <div class="grid">
                  <div class="col-12 md:col-6 flex flex-column gap-2">
                    <label for="table_number">{{ $t('table_number') }}</label>
                    <InputNumber
                      id="table_number"
                      v-model="newTable.number"
                      :class="{ 'p-invalid': newTableErrors.number }"
                      :min="1"
                    />
                    <small class="p-error">{{ newTableErrors.number }}</small>
                  </div>
                  <div class="col-12 md:col-6 flex flex-column gap-2">
                    <label for="table_name">{{ $t('table_name') }}</label>
                    <InputText
                      id="table_name"
                      v-model="newTable.name"
                      :class="{ 'p-invalid': newTableErrors.name }"
                    />
                    <small class="p-error">{{ newTableErrors.name }}</small>
                  </div>
                </div>
                <div class="grid">
                  <div class="col-12 md:col-6 flex flex-column gap-2">
                    <label for="table_capacity">{{ $t('table_capacity') }}</label>
                    <InputNumber
                      id="table_capacity"
                      v-model="newTable.capacity"
                      :min="1"
                      :class="{ 'p-invalid': newTableErrors.capacity }"
                    />
                    <small class="p-error">{{ newTableErrors.capacity }}</small>
                  </div>
                  <div class="col-12 md:col-6 flex flex-column gap-2">
                    <label for="table_zone">{{ $t('table_zone') }}</label>
                    <Dropdown
                      id="table_zone"
                      v-model="newTable.zone"
                      :options="zoneOptions"
                      optionLabel="label"
                      optionValue="value"
                      :placeholder="$t('select_option')"
                    />
                  </div>
                </div>
                <div class="flex flex-column gap-2">
                  <label for="table_status">{{ $t('table_status') }}</label>
                  <Dropdown
                    id="table_status"
                    v-model="newTable.status"
                    :options="statusOptions"
                    optionLabel="label"
                    optionValue="value"
                    :placeholder="$t('select_option')"
                  />
                </div>
              </div>
              <template #footer>
                <ButtonGroup>
                  <Button
                    :label="$t('cancel')"
                    severity="secondary"
                    @click="tableAddDialog = false"
                  />
                  <Button
                    class="ml-2"
                    severity="primary"
                    @click="submitTable"
                    :label="$t('save')"
                    :loading="isSubmitting"
                    :disabled="isSubmitting"
                  />
                </ButtonGroup>
              </template>
            </Dialog>

            <Dialog
              v-model:visible="tableEditDialog"
              modal
              :header="$t('edit') + ' ' + $t('table_number') + ' ' + tableToEdit.number"
              :style="{ width: '45rem', direction: store.orientation == 'rtl' ? 'rtl' : 'ltr' }"
              :breakpoints="{ '1199px': '90vw', '575px': '90vw' }"
            >
              <div class="flex flex-column gap-2 w-full">
                <div class="grid">
                  <div class="col-12 md:col-6 flex flex-column gap-2">
                    <label for="edit_table_number">{{ $t('table_number') }}</label>
                    <InputNumber
                      id="edit_table_number"
                      v-model="tableToEdit.number"
                      :min="1"
                      :class="{ 'p-invalid': editTableErrors.number }"
                    />
                    <small class="p-error">{{ editTableErrors.number }}</small>
                  </div>
                  <div class="col-12 md:col-6 flex flex-column gap-2">
                    <label for="edit_table_name">{{ $t('table_name') }}</label>
                    <InputText
                      id="edit_table_name"
                      v-model="tableToEdit.name"
                      :class="{ 'p-invalid': editTableErrors.name }"
                    />
                    <small class="p-error">{{ editTableErrors.name }}</small>
                  </div>
                </div>
                <div class="grid">
                  <div class="col-12 md:col-6 flex flex-column gap-2">
                    <label for="edit_table_capacity">{{ $t('table_capacity') }}</label>
                    <InputNumber
                      id="edit_table_capacity"
                      v-model="tableToEdit.capacity"
                      :min="1"
                      :class="{ 'p-invalid': editTableErrors.capacity }"
                    />
                    <small class="p-error">{{ editTableErrors.capacity }}</small>
                  </div>
                  <div class="col-12 md:col-6 flex flex-column gap-2">
                    <label for="edit_table_zone">{{ $t('table_zone') }}</label>
                    <Dropdown
                      id="edit_table_zone"
                      v-model="tableToEdit.zone"
                      :options="zoneOptions"
                      optionLabel="label"
                      optionValue="value"
                    />
                  </div>
                </div>
                <div class="flex flex-column gap-2">
                  <label for="edit_table_status">{{ $t('table_status') }}</label>
                  <Dropdown
                    id="edit_table_status"
                    v-model="tableToEdit.status"
                    :options="statusOptions"
                    optionLabel="label"
                    optionValue="value"
                  />
                </div>
              </div>
              <template #footer>
                <ButtonGroup>
                  <Button
                    :label="$t('cancel')"
                    severity="secondary"
                    @click="tableEditDialog = false"
                  />
                  <Button
                    class="ml-2"
                    severity="primary"
                    @click="updateTable"
                    :label="$t('save')"
                    :loading="isSubmitting"
                    :disabled="isSubmitting"
                  />
                </ButtonGroup>
              </template>
            </Dialog>

            <Dialog
              v-model:visible="qrDialogVisible"
              modal
              :header="$t('qr_code') + ' - ' + $t('table_number') + ' ' + qrTableData.number"
              :style="{ width: '30rem', direction: store.orientation == 'rtl' ? 'rtl' : 'ltr' }"
              :breakpoints="{ '1199px': '90vw', '575px': '90vw' }"
            >
              <div class="flex flex-column align-items-center gap-4">
                <div
                  class="flex align-items-center justify-content-center border-2 border-dashed border-round p-4"
                  style="width: 200px; height: 200px"
                >
                  <img
                    v-if="qrTableData.qr_code"
                    :src="qrTableData.qr_code"
                    alt="QR Code"
                    style="max-width: 100%; max-height: 100%"
                  />
                  <div v-else class="text-center text-500">
                    <i class="pi pi-qrcode" style="font-size: 3rem; opacity: 0.3"></i>
                    <p class="mt-2">{{ $t('no_results') }}</p>
                  </div>
                </div>
                <p class="text-center text-sm text-500">
                  {{ $t('table_name') }}: {{ qrTableData.name }}
                </p>
              </div>
              <template #footer>
                <ButtonGroup>
                  <Button
                    :label="$t('generate_qr')"
                    icon="pi pi-refresh"
                    severity="secondary"
                    @click="generateQRCode"
                    :loading="isGeneratingQR"
                  />
                  <Button
                    :label="$t('close')"
                    severity="primary"
                    @click="qrDialogVisible = false"
                  />
                </ButtonGroup>
              </template>
            </Dialog>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import ConfirmPopup from 'primevue/confirmpopup'
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
import { useConfirm } from 'primevue/useconfirm'
import { globalStore } from '@/stores'
import auth from '../services/auth'
import type { TableItem, DataTablePageEvent } from '@/types'

const { t } = useI18n()
const store = globalStore()

const confirm = useConfirm()
const toast = useToast()

const tables = ref<TableItem[]>([])
const tablesTableTotalRecords = ref(0)
const isTablesTableLoading = ref(false)
const tablesTableRowsPerPage = ref(50)
const isSubmitting = ref(false)
const tableAddDialog = ref(false)
const tableEditDialog = ref(false)

const zoneFilter = ref('')
const statusFilter = ref('')

const zoneOptions = computed(() => [
  { label: t('indoor'), value: 'indoor' },
  { label: t('outdoor'), value: 'outdoor' },
  { label: t('bar'), value: 'bar' },
  { label: t('vip'), value: 'vip' },
])

const statusOptions = computed(() => [
  { label: t('available'), value: 'available' },
  { label: t('occupied'), value: 'occupied' },
  { label: t('reserved'), value: 'reserved' },
  { label: t('cleaning'), value: 'cleaning' },
])

const filteredTables = computed(() => {
  let result = tables.value
  if (zoneFilter.value) {
    result = result.filter((tbl) => tbl.zone === zoneFilter.value)
  }
  if (statusFilter.value) {
    result = result.filter((tbl) => tbl.status === statusFilter.value)
  }
  return result
})

const newTable = ref<Partial<TableItem>>({
  number: 1,
  name: '',
  capacity: 4,
  zone: 'indoor',
  status: 'available',
})
const newTableErrors = ref({ number: '', name: '', capacity: '' })

const tableToEdit = ref<Partial<TableItem>>({})
const editTableErrors = ref({ number: '', name: '', capacity: '' })

const qrDialogVisible = ref(false)
const qrTableData = ref<Partial<TableItem>>({})
const isGeneratingQR = ref(false)

const openAddDialog = () => {
  newTable.value = { number: 1, name: '', capacity: 4, zone: 'indoor', status: 'available' }
  newTableErrors.value = { number: '', name: '', capacity: '' }
  tableAddDialog.value = true
}

const prepareTableToEdit = (table: TableItem) => {
  tableToEdit.value = JSON.parse(JSON.stringify(table))
  editTableErrors.value = { number: '', name: '', capacity: '' }
  tableEditDialog.value = true
}

const updateTablesTableRowsPerPage = (event: DataTablePageEvent) => {
  getTables(event.first, event.rows)
}

const getZoneSeverity = (zone: string) => {
  switch (zone) {
    case 'indoor':
      return 'info'
    case 'outdoor':
      return 'success'
    case 'bar':
      return 'warn'
    case 'vip':
      return 'danger'
    default:
      return 'secondary'
  }
}

const getStatusSeverity = (status: string) => {
  switch (status) {
    case 'available':
      return 'success'
    case 'occupied':
      return 'danger'
    case 'reserved':
      return 'warn'
    case 'cleaning':
      return 'secondary'
    default:
      return 'info'
  }
}

const submitTable = () => {
  newTableErrors.value.number = newTable.value.number ? '' : t('validation_required')
  newTableErrors.value.name = newTable.value.name?.trim() ? '' : t('validation_required')
  newTableErrors.value.capacity = newTable.value.capacity ? '' : t('validation_required')

  if (newTableErrors.value.number || newTableErrors.value.name || newTableErrors.value.capacity)
    return

  isSubmitting.value = true

  axios
    .post(
      `http://${import.meta.env.VITE_APP_BACKEND_HOST}${import.meta.env.VITE_APP_MODULE_CORE_API_PREFIX}/table/api/tables`,
      {
        data: newTable.value,
      },
      {
        headers: {
          Authorization: `Bearer ${auth.accessToken.value}`,
        },
      },
    )
    .then(() => {
      tableAddDialog.value = false
      newTable.value = { number: 1, name: '', capacity: 4, zone: 'indoor', status: 'available' }
      getTables()
      setTimeout(() => {
        toast.add({
          severity: 'success',
          summary: t('success'),
          detail: t('table_added_success'),
          group: 'br',
        })
      }, 1000)
    })
    .catch(() => {
      toast.add({
        severity: 'error',
        summary: t('failed'),
        detail: t('table_add_failed'),
        group: 'br',
      })
    })
    .finally(() => {
      isSubmitting.value = false
    })
}

const updateTable = () => {
  editTableErrors.value.number = tableToEdit.value.number ? '' : t('validation_required')
  editTableErrors.value.name = tableToEdit.value.name?.trim() ? '' : t('validation_required')
  editTableErrors.value.capacity = tableToEdit.value.capacity ? '' : t('validation_required')

  if (editTableErrors.value.number || editTableErrors.value.name || editTableErrors.value.capacity)
    return

  isSubmitting.value = true

  axios
    .patch(
      `http://${import.meta.env.VITE_APP_BACKEND_HOST}${import.meta.env.VITE_APP_MODULE_CORE_API_PREFIX}/table/api/tables/${tableToEdit.value.id}`,
      {
        data: tableToEdit.value,
      },
      {
        headers: {
          Authorization: `Bearer ${auth.accessToken.value}`,
        },
      },
    )
    .then(() => {
      toast.add({
        severity: 'success',
        summary: t('success'),
        detail: t('table_updated_success'),
        group: 'br',
      })
      tableEditDialog.value = false
      getTables()
    })
    .catch(() => {
      toast.add({
        severity: 'error',
        summary: t('failed'),
        detail: t('table_update_failed'),
        group: 'br',
      })
    })
    .finally(() => {
      isSubmitting.value = false
    })
}

const deleteTable = (table_id: string) => {
  axios
    .delete(
      `http://${import.meta.env.VITE_APP_BACKEND_HOST}${import.meta.env.VITE_APP_MODULE_CORE_API_PREFIX}/table/api/tables/${table_id}`,
      {
        headers: {
          Authorization: `Bearer ${auth.accessToken.value}`,
        },
      },
    )
    .then(() => {
      setTimeout(() => {
        toast.add({
          severity: 'success',
          summary: t('success'),
          detail: t('table_deleted_success'),
          group: 'br',
        })
      }, 1000)

      getTables()
    })
    .catch(() => {
      toast.add({
        severity: 'error',
        summary: t('failed'),
        detail: t('table_delete_failed'),
        group: 'br',
      })
    })
}

const confirmDeleteTable = (event, table_id: string) => {
  confirm.require({
    target: event.currentTarget,
    message: t('confirm_delete_table'),
    icon: 'pi pi-exclamation-triangle',
    rejectProps: {
      label: t('cancel'),
      severity: 'secondary',
      outlined: true,
    },
    acceptProps: {
      label: t('yes'),
    },
    accept: () => {
      deleteTable(table_id)
    },
    reject: () => {},
  })
}

const showQRCode = (table: TableItem) => {
  qrTableData.value = { ...table }
  qrDialogVisible.value = true
}

const generateQRCode = () => {
  if (!qrTableData.value.id) return

  isGeneratingQR.value = true

  axios
    .post(
      `http://${import.meta.env.VITE_APP_BACKEND_HOST}${import.meta.env.VITE_APP_MODULE_CORE_API_PREFIX}/table/api/tables/${qrTableData.value.id}/qr`,
      {},
      {
        headers: {
          Authorization: `Bearer ${auth.accessToken.value}`,
        },
      },
    )
    .then((response) => {
      qrTableData.value.qr_code = response.data.data.qr_code
      getTables()
      toast.add({
        severity: 'success',
        summary: t('success'),
        detail: t('qr_code_generated'),
        group: 'br',
      })
    })
    .catch(() => {
      toast.add({
        severity: 'error',
        summary: t('failed'),
        detail: t('qr_generate_failed'),
        group: 'br',
      })
    })
    .finally(() => {
      isGeneratingQR.value = false
    })
}

const getTables = (first = 0, rows = tablesTableRowsPerPage.value) => {
  isTablesTableLoading.value = true

  if (first == 0) {
    first = 1
  }

  const page_number = Math.ceil(first / rows)

  axios
    .get(
      `http://${import.meta.env.VITE_APP_BACKEND_HOST}${import.meta.env.VITE_APP_MODULE_CORE_API_PREFIX}/table/api/tables?page[number]=${page_number}&page[size]=${rows}`,
      {
        headers: {
          Authorization: `Bearer ${auth.accessToken.value}`,
        },
      },
    )
    .then((response) => {
      tables.value = response.data.data
      tablesTableTotalRecords.value = response.data.meta.total_records
    })
    .catch(() => {
      toast.add({ severity: 'error', summary: t('failed'), detail: t('tables_load_failed') })
    })
    .finally(() => {
      isTablesTableLoading.value = false
    })
}

getTables()
</script>

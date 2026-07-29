<template>
  <div class="w-full">
    <div class="grid mx-2">
      <div class="col-12">
        <h3>{{ $t('employee_scheduling') }}</h3>
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
          :value="shifts"
          stripedRows
          tableStyle="width: 100%;max-height:60vh;"
        >
          <template #header>
            <div class="flex justify-between align-items-center">
              <Button icon="pi pi-plus" :label="$t('add_shift')" rounded raised @click="openAdd" />
              <div class="flex gap-2">
                <Calendar v-model="dateRange" selectionMode="range" :placeholder="$t('date_range')" showIcon />
              </div>
            </div>
          </template>
          <template #empty>
            <div class="flex flex-column align-items-center gap-2 py-4">
              <i class="pi pi-calendar" style="font-size: 2rem; opacity: 0.3"></i>
              <p class="m-0 text-slate-400">{{ $t('no_results') }}</p>
            </div>
          </template>
          <Column sortable field="employee_id" :header="$t('employee')"></Column>
          <Column sortable field="date" :header="$t('date')"></Column>
          <Column sortable field="start_time" :header="$t('start_time')"></Column>
          <Column sortable field="end_time" :header="$t('end_time')"></Column>
          <Column sortable field="role" :header="$t('role')"></Column>
          <Column sortable field="status" :header="$t('status')">
            <template #body="slotProps">
              <Tag :value="slotProps.data.status" :severity="getStatusSeverity(slotProps.data.status)" />
            </template>
          </Column>
          <Column :header="$t('actions')">
            <template #body="slotProps">
              <ButtonGroup>
                <Button icon="pi pi-pencil" severity="secondary" @click="prepareEdit(slotProps.data)" />
                <Button icon="pi pi-trash" severity="danger" @click="confirmDelete($event, slotProps.data.id)" />
              </ButtonGroup>
            </template>
          </Column>
        </DataTable>
      </div>
    </div>

    <Dialog v-model:visible="addDialog" modal :header="$t('add_shift')" :style="{ width: '30rem' }" :breakpoints="{ '1199px': '90vw', '575px': '90vw' }">
      <div class="flex flex-column gap-4">
        <div class="flex flex-column gap-2">
          <label>{{ $t('employee') }}</label>
          <InputText v-model="form.employee_id" />
        </div>
        <div class="flex flex-column gap-2">
          <label>{{ $t('branch') }}</label>
          <InputText v-model="form.branch_id" />
        </div>
        <div class="flex flex-column gap-2">
          <label>{{ $t('date') }}</label>
          <Calendar v-model="form.date" dateFormat="yy-mm-dd" showIcon />
        </div>
        <div class="grid">
          <div class="col-6 flex flex-column gap-2">
            <label>{{ $t('start_time') }}</label>
            <Calendar v-model="form.start_time" timeOnly hourFormat="24" />
          </div>
          <div class="col-6 flex flex-column gap-2">
            <label>{{ $t('end_time') }}</label>
            <Calendar v-model="form.end_time" timeOnly hourFormat="24" />
          </div>
        </div>
        <div class="flex flex-column gap-2">
          <label>{{ $t('role') }}</label>
          <Dropdown v-model="form.role" :options="roleOptions" optionLabel="label" optionValue="value" />
        </div>
        <div class="flex flex-column gap-2">
          <label>{{ $t('notes') }}</label>
          <InputText v-model="form.notes" />
        </div>
      </div>
      <template #footer>
        <ButtonGroup>
          <Button :label="$t('cancel')" severity="secondary" @click="addDialog = false" />
          <Button class="ml-2" severity="primary" @click="submit" :label="$t('save')" :loading="submitting" :disabled="submitting" />
        </ButtonGroup>
      </template>
    </Dialog>

    <Dialog v-model:visible="editDialog" modal :header="$t('edit_shift')" :style="{ width: '30rem' }" :breakpoints="{ '1199px': '90vw', '575px': '90vw' }">
      <div class="flex flex-column gap-4">
        <div class="flex flex-column gap-2">
          <label>{{ $t('employee') }}</label>
          <InputText v-model="editForm.employee_id" />
        </div>
        <div class="flex flex-column gap-2">
          <label>{{ $t('date') }}</label>
          <Calendar v-model="editForm.date" dateFormat="yy-mm-dd" showIcon />
        </div>
        <div class="grid">
          <div class="col-6 flex flex-column gap-2">
            <label>{{ $t('start_time') }}</label>
            <Calendar v-model="editForm.start_time" timeOnly hourFormat="24" />
          </div>
          <div class="col-6 flex flex-column gap-2">
            <label>{{ $t('end_time') }}</label>
            <Calendar v-model="editForm.end_time" timeOnly hourFormat="24" />
          </div>
        </div>
        <div class="flex flex-column gap-2">
          <label>{{ $t('status') }}</label>
          <Dropdown v-model="editForm.status" :options="statusOptions" optionLabel="label" optionValue="value" />
        </div>
      </div>
      <template #footer>
        <ButtonGroup>
          <Button :label="$t('cancel')" severity="secondary" @click="editDialog = false" />
          <Button class="ml-2" severity="primary" @click="submitEdit" :label="$t('save')" :loading="submitting" :disabled="submitting" />
        </ButtonGroup>
      </template>
    </Dialog>
  </div>
</template>

<script setup lang="ts">
import DataTable from 'primevue/datatable'
import Dialog from 'primevue/dialog'
import InputText from 'primevue/inputtext'
import Column from 'primevue/column'
import Button from 'primevue/button'
import ButtonGroup from 'primevue/buttongroup'
import Tag from 'primevue/tag'
import Dropdown from 'primevue/dropdown'
import Calendar from 'primevue/calendar'
import axios from 'axios'
import { ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { useToast } from 'primevue/usetoast'
import { useConfirm } from 'primevue/useconfirm'

const { t } = useI18n()
const toast = useToast()
const confirm = useConfirm()

const shifts = ref<any[]>([])
const totalRecords = ref(0)
const loading = ref(false)
const rowsPerPage = ref(50)
const submitting = ref(false)
const addDialog = ref(false)
const editDialog = ref(false)
const dateRange = ref<Date[]>([])

const form = ref({ employee_id: '', branch_id: '', date: '', start_time: '', end_time: '', role: 'cashier', notes: '' })
const editForm = ref<any>({})
const editId = ref('')

const roleOptions = [
  { label: 'Cashier', value: 'cashier' },
  { label: 'Chef', value: 'chef' },
  { label: 'Waiter', value: 'waiter' },
  { label: 'Manager', value: 'manager' },
  { label: 'Bartender', value: 'bartender' },
]

const statusOptions = [
  { label: 'Scheduled', value: 'scheduled' },
  { label: 'Confirmed', value: 'confirmed' },
  { label: 'Completed', value: 'completed' },
  { label: 'Cancelled', value: 'cancelled' },
]

const getStatusSeverity = (status: string) => {
  switch (status) {
    case 'confirmed': return 'success'
    case 'completed': return 'info'
    case 'cancelled': return 'danger'
    default: return 'warn'
  }
}

const formatTime = (d: Date) => {
  if (!d) return ''
  return `${String(d.getHours()).padStart(2, '0')}:${String(d.getMinutes()).padStart(2, '0')}`
}

const formatDate = (d: Date) => {
  if (!d) return ''
  return `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, '0')}-${String(d.getDate()).padStart(2, '0')}`
}

const apiBase = `http://${import.meta.env.VITE_APP_BACKEND_HOST}${import.meta.env.VITE_APP_MODULE_CORE_API_PREFIX}/scheduling/api/shifts`

const getShifts = (offset = 0, limit = 50) => {
  loading.value = true
  const pageNumber = Math.floor(offset / limit) + 1
  let url = `${apiBase}?page_number=${pageNumber}&page_size=${limit}`
  if (dateRange.value && dateRange.value.length === 2) {
    url += `&start_date=${formatDate(dateRange.value[0])}&end_date=${formatDate(dateRange.value[1])}`
  }
  axios.get(url).then((res) => {
    shifts.value = res.data.data || []
    totalRecords.value = res.data.meta?.total_records || 0
  }).catch((err) => {
    toast.add({ severity: 'error', summary: t('error'), detail: err.message, life: 3000 })
  }).finally(() => { loading.value = false })
}

const onPage = (event: { first: number; rows: number }) => { getShifts(event.first, event.rows) }
const openAdd = () => { form.value = { employee_id: '', branch_id: '', date: '', start_time: '', end_time: '', role: 'cashier', notes: '' }; addDialog.value = true }
const prepareEdit = (shift: any) => { editForm.value = JSON.parse(JSON.stringify(shift)); editId.value = shift.id; editDialog.value = true }

const submit = () => {
  submitting.value = true
  const payload = {
    ...form.value,
    date: form.value.date ? formatDate(new Date(form.value.date)) : '',
    start_time: form.value.start_time ? formatTime(new Date(form.value.start_time)) : '',
    end_time: form.value.end_time ? formatTime(new Date(form.value.end_time)) : '',
  }
  axios.post(apiBase, payload).then(() => {
    toast.add({ severity: 'success', summary: t('success'), detail: t('shift_added'), life: 3000 })
    addDialog.value = false; getShifts(0, rowsPerPage.value)
  }).catch((err) => {
    toast.add({ severity: 'error', summary: t('error'), detail: err.message, life: 3000 })
  }).finally(() => { submitting.value = false })
}

const submitEdit = () => {
  submitting.value = true
  const payload = {
    ...editForm.value,
    date: editForm.value.date ? formatDate(new Date(editForm.value.date)) : '',
    start_time: editForm.value.start_time ? formatTime(new Date(editForm.value.start_time)) : '',
    end_time: editForm.value.end_time ? formatTime(new Date(editForm.value.end_time)) : '',
  }
  axios.patch(`${apiBase}/${editId.value}`, payload).then(() => {
    toast.add({ severity: 'success', summary: t('success'), detail: t('shift_updated'), life: 3000 })
    editDialog.value = false; getShifts(0, rowsPerPage.value)
  }).catch((err) => {
    toast.add({ severity: 'error', summary: t('error'), detail: err.message, life: 3000 })
  }).finally(() => { submitting.value = false })
}

const confirmDelete = (event: MouseEvent, id: string) => {
  confirm.require({
    target: event.target as HTMLElement,
    message: t('do_you_confirm'),
    icon: 'pi pi-exclamation-triangle',
    rejectClass: 'p-button-secondary p-button-outlined',
    rejectLabel: t('cancel'),
    acceptLabel: t('delete'),
    acceptClass: 'p-button-danger',
    accept: () => {
      axios.delete(`${apiBase}/${id}`).then(() => {
        toast.add({ severity: 'success', summary: t('success'), detail: t('shift_deleted'), life: 3000 })
        getShifts(0, rowsPerPage.value)
      }).catch((err) => {
        toast.add({ severity: 'error', summary: t('error'), detail: err.message, life: 3000 })
      })
    },
  })
}

getShifts()
</script>

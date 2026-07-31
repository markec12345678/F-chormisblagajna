<template>
  <div class="w-full">
    <div class="grid mx-2">
      <div class="col-12">
        <h3>{{ $t('promotions') }}</h3>
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
          :value="promotions"
          stripedRows
          tableStyle="width: 100%;max-height:60vh;"
        >
          <template #header>
            <div class="flex justify-between align-items-center">
              <Button
                icon="pi pi-plus"
                :label="$t('add_promotion')"
                rounded
                raised
                @click="openAdd"
              />
            </div>
          </template>
          <template #empty>
            <div class="flex flex-column align-items-center gap-2 py-4">
              <i class="pi pi-tags" style="font-size: 2rem; opacity: 0.3"></i>
              <p class="m-0 text-slate-400">{{ $t('no_results') }}</p>
            </div>
          </template>
          <Column sortable field="name" :header="$t('name')"></Column>
          <Column sortable field="code" :header="$t('code')"></Column>
          <Column sortable field="type" :header="$t('type')">
            <template #body="slotProps">
              <Tag
                :value="slotProps.data.type"
                :severity="slotProps.data.type === 'percentage' ? 'info' : 'success'"
              />
            </template>
          </Column>
          <Column sortable field="value" :header="$t('value')">
            <template #body="slotProps">
              {{
                slotProps.data.type === 'percentage'
                  ? slotProps.data.value + '%'
                  : slotProps.data.value + ' EUR'
              }}
            </template>
          </Column>
          <Column sortable field="start_date" :header="$t('start_date')"></Column>
          <Column sortable field="end_date" :header="$t('end_date')"></Column>
          <Column sortable field="is_active" :header="$t('status')">
            <template #body="slotProps">
              <Tag
                :value="slotProps.data.is_active ? $t('active') : $t('inactive')"
                :severity="slotProps.data.is_active ? 'success' : 'danger'"
              />
            </template>
          </Column>
          <Column :header="$t('actions')">
            <template #body="slotProps">
              <ButtonGroup>
                <Button
                  icon="pi pi-pencil"
                  severity="secondary"
                  @click="prepareEdit(slotProps.data)"
                />
                <Button
                  icon="pi pi-trash"
                  severity="danger"
                  @click="confirmDelete($event, slotProps.data.id)"
                />
              </ButtonGroup>
            </template>
          </Column>
        </DataTable>
      </div>
    </div>

    <Dialog
      v-model:visible="addDialog"
      modal
      :header="$t('add_promotion')"
      :style="{ width: '35rem' }"
      :breakpoints="{ '1199px': '90vw', '575px': '90vw' }"
    >
      <div class="flex flex-column gap-4">
        <div class="grid">
          <div class="col-6 flex flex-column gap-2">
            <label>{{ $t('name') }}</label>
            <InputText v-model="form.name" />
          </div>
          <div class="col-6 flex flex-column gap-2">
            <label>{{ $t('code') }}</label>
            <InputText v-model="form.code" />
          </div>
        </div>
        <div class="grid">
          <div class="col-4 flex flex-column gap-2">
            <label>{{ $t('type') }}</label>
            <Dropdown
              v-model="form.type"
              :options="typeOptions"
              optionLabel="label"
              optionValue="value"
            />
          </div>
          <div class="col-4 flex flex-column gap-2">
            <label>{{ $t('value') }}</label>
            <InputNumber v-model="form.value" :min="0" />
          </div>
          <div class="col-4 flex flex-column gap-2">
            <label>{{ $t('min_order') }}</label>
            <InputNumber v-model="form.min_order" :min="0" />
          </div>
        </div>
        <div class="grid">
          <div class="col-6 flex flex-column gap-2">
            <label>{{ $t('start_date') }}</label>
            <Calendar v-model="form.start_date" dateFormat="yy-mm-dd" showIcon />
          </div>
          <div class="col-6 flex flex-column gap-2">
            <label>{{ $t('end_date') }}</label>
            <Calendar v-model="form.end_date" dateFormat="yy-mm-dd" showIcon />
          </div>
        </div>
        <div class="flex align-items-center gap-2">
          <Checkbox v-model="form.is_active" :binary="true" inputId="active" />
          <label for="active">{{ $t('active') }}</label>
        </div>
      </div>
      <template #footer>
        <ButtonGroup>
          <Button :label="$t('cancel')" severity="secondary" @click="addDialog = false" />
          <Button
            class="ml-2"
            severity="primary"
            @click="submit"
            :label="$t('save')"
            :loading="submitting"
            :disabled="submitting"
          />
        </ButtonGroup>
      </template>
    </Dialog>

    <Dialog
      v-model:visible="editDialog"
      modal
      :header="$t('edit_promotion')"
      :style="{ width: '35rem' }"
      :breakpoints="{ '1199px': '90vw', '575px': '90vw' }"
    >
      <div class="flex flex-column gap-4">
        <div class="grid">
          <div class="col-6 flex flex-column gap-2">
            <label>{{ $t('name') }}</label>
            <InputText v-model="editForm.name" />
          </div>
          <div class="col-6 flex flex-column gap-2">
            <label>{{ $t('code') }}</label>
            <InputText v-model="editForm.code" />
          </div>
        </div>
        <div class="grid">
          <div class="col-4 flex flex-column gap-2">
            <label>{{ $t('value') }}</label>
            <InputNumber v-model="editForm.value" :min="0" />
          </div>
          <div class="col-4 flex flex-column gap-2">
            <label>{{ $t('min_order') }}</label>
            <InputNumber v-model="editForm.min_order" :min="0" />
          </div>
          <div class="col-4 flex flex-column gap-2">
            <label>{{ $t('usage_limit') }}</label>
            <InputNumber v-model="editForm.usage_limit" :min="0" />
          </div>
        </div>
        <div class="flex align-items-center gap-2">
          <Checkbox v-model="editForm.is_active" :binary="true" inputId="editActive" />
          <label for="editActive">{{ $t('active') }}</label>
        </div>
      </div>
      <template #footer>
        <ButtonGroup>
          <Button :label="$t('cancel')" severity="secondary" @click="editDialog = false" />
          <Button
            class="ml-2"
            severity="primary"
            @click="submitEdit"
            :label="$t('save')"
            :loading="submitting"
            :disabled="submitting"
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
import Calendar from 'primevue/calendar'
import Checkbox from 'primevue/checkbox'
import axios from 'axios'
import { ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { useToast } from 'primevue/usetoast'
import { useConfirm } from 'primevue/useconfirm'

const { t } = useI18n()
const toast = useToast()
const confirm = useConfirm()

interface Promotion { id: string; name: string; code: string; type: string; value: number; min_order: number; start_date: string; end_date: string; active: boolean; usage_count: number; max_uses: number }
interface FormData { name: string; code: string; type: string; value: number; min_order: number; start_date: string; end_date: string; max_uses: number }
const promotions = ref<Promotion[]>([])
const totalRecords = ref(0)
const loading = ref(false)
const rowsPerPage = ref(50)
const submitting = ref(false)
const addDialog = ref(false)
const editDialog = ref(false)

const form = ref<FormData>({
  name: '',
  code: '',
  type: 'percentage',
  value: 0,
  min_order: 0,
  start_date: '',
  end_date: '',
  is_active: true,
})
const editForm = ref<Partial<Promotion>>({})
const editId = ref('')

const typeOptions = [
  { label: 'Percentage', value: 'percentage' },
  { label: 'Fixed Amount', value: 'fixed' },
  { label: 'Free Item', value: 'free_item' },
]

const formatDate = (d: Date) => {
  if (!d) return ''
  return `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, '0')}-${String(d.getDate()).padStart(2, '0')}`
}

const apiBase = `http://${import.meta.env.VITE_APP_BACKEND_HOST}/promotion/api/promotions`

const getPromotions = (offset = 0, limit = 50) => {
  loading.value = true
  const pageNumber = Math.floor(offset / limit) + 1
  axios
    .get(`${apiBase}?page_number=${pageNumber}&page_size=${limit}`)
    .then((res) => {
      promotions.value = res.data.data || []
      totalRecords.value = res.data.meta?.total_records || 0
    })
    .catch((err) => {
      toast.add({ severity: 'error', summary: t('error'), detail: err.message, life: 3000 })
    })
    .finally(() => {
      loading.value = false
    })
}

const onPage = (event: { first: number; rows: number }) => {
  getPromotions(event.first, event.rows)
}
const openAdd = () => {
  form.value = {
    name: '',
    code: '',
    type: 'percentage',
    value: 0,
    min_order: 0,
    start_date: '',
    end_date: '',
    is_active: true,
  }
  addDialog.value = true
}
const prepareEdit = (p: Promotion) => {
  editForm.value = JSON.parse(JSON.stringify(p))
  editId.value = p.id
  editDialog.value = true
}

const submit = () => {
  submitting.value = true
  const payload = {
    ...form.value,
    start_date: form.value.start_date ? formatDate(new Date(form.value.start_date)) : '',
    end_date: form.value.end_date ? formatDate(new Date(form.value.end_date)) : '',
  }
  axios
    .post(apiBase, payload)
    .then(() => {
      toast.add({
        severity: 'success',
        summary: t('success'),
        detail: t('promotion_added'),
        life: 3000,
      })
      addDialog.value = false
      getPromotions(0, rowsPerPage.value)
    })
    .catch((err) => {
      toast.add({ severity: 'error', summary: t('error'), detail: err.message, life: 3000 })
    })
    .finally(() => {
      submitting.value = false
    })
}

const submitEdit = () => {
  submitting.value = true
  axios
    .patch(`${apiBase}/${editId.value}`, editForm.value)
    .then(() => {
      toast.add({
        severity: 'success',
        summary: t('success'),
        detail: t('promotion_updated'),
        life: 3000,
      })
      editDialog.value = false
      getPromotions(0, rowsPerPage.value)
    })
    .catch((err) => {
      toast.add({ severity: 'error', summary: t('error'), detail: err.message, life: 3000 })
    })
    .finally(() => {
      submitting.value = false
    })
}

const confirmDelete = (event: MouseEvent, id: string) => {
  confirm.require({
    target: event.target as HTMLElement,
    message: t('do_you_confirm'),
    icon: 'pi pi-exclamation-triangle',
    rejectLabel: t('cancel'),
    acceptLabel: t('delete'),
    acceptClass: 'p-button-danger',
    accept: () => {
      axios
        .delete(`${apiBase}/${id}`)
        .then(() => {
          toast.add({
            severity: 'success',
            summary: t('success'),
            detail: t('promotion_deleted'),
            life: 3000,
          })
          getPromotions(0, rowsPerPage.value)
        })
        .catch((err) => {
          toast.add({ severity: 'error', summary: t('error'), detail: err.message, life: 3000 })
        })
    },
  })
}

getPromotions()
</script>

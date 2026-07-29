<template>
  <div class="w-full">
    <div class="grid mx-2">
      <div class="col-12 flex justify-content-between align-items-center">
        <h3>{{ $t('reservations') }}</h3>
        <Button :label="$t('add_reservation')" icon="pi pi-plus" @click="openNew" />
      </div>
      <div class="col-12">
        <DataTable :value="list" stripedRows :loading="loading">
          <template #empty>{{ $t('no_reservations') }}</template>
          <Column field="customer_name" :header="$t('customer')"></Column>
          <Column field="customer_phone" :header="$t('phone')"></Column>
          <Column field="guest_count" :header="$t('guests')"></Column>
          <Column field="reservation_date" :header="$t('date')"></Column>
          <Column field="reservation_time" :header="$t('time')"></Column>
          <Column field="status" :header="$t('status')">
            <template #body="s">
              <Tag :severity="s.data.status==='confirmed'?'success':s.data.status==='cancelled'?'danger':'warn'">{{ s.data.status }}</Tag>
            </template>
          </Column>
          <Column field="table_assignment" :header="$t('table')">
            <template #body="s">{{ s.data.table_assignment || '-' }}</template>
          </Column>
          <Column :header="$t('actions')">
            <template #body="s">
              <Button icon="pi pi-check" severity="success" class="mr-1" size="small" v-if="s.data.status==='pending'" @click="confirmRes(s.data)" />
              <Button icon="pi pi-times" severity="danger" size="small" @click="doDelete(s.data)" />
            </template>
          </Column>
        </DataTable>
      </div>
    </div>

    <Dialog v-model:visible="dialog" :header="$t('add_reservation')" :style="{width:'450px'}">
      <div class="grid">
        <div class="col-12"><label>{{ $t('customer_name') }}</label><InputText v-model="form.customer_name" class="w-full" /></div>
        <div class="col-6"><label>{{ $t('phone') }}</label><InputText v-model="form.customer_phone" class="w-full" /></div>
        <div class="col-6"><label>{{ $t('email') }}</label><InputText v-model="form.customer_email" class="w-full" /></div>
        <div class="col-6"><label>{{ $t('date') }}</label><Calendar v-model="formDate" dateFormat="yy-mm-dd" class="w-full" /></div>
        <div class="col-6"><label>{{ $t('time') }}</label><Calendar v-model="formTime" timeOnly class="w-full" /></div>
        <div class="col-6"><label>{{ $t('guests') }}</label><InputNumber v-model="form.guest_count" :min="1" class="w-full" /></div>
        <div class="col-12"><label>{{ $t('notes') }}</label><Textarea v-model="form.notes" class="w-full" rows="2" /></div>
      </div>
      <template #footer>
        <Button :label="$t('cancel')" severity="secondary" @click="dialog=false" />
        <Button :label="$t('save')" @click="doSave" :loading="saving" />
      </template>
    </Dialog>

    <Dialog v-model:visible="assignDialog" :header="$t('assign_table')" :style="{width:'400px'}">
      <div class="grid">
        <div class="col-12"><label>{{ $t('table') }}</label><InputText v-model="assignTable" class="w-full" /></div>
      </div>
      <template #footer>
        <Button :label="$t('cancel')" severity="secondary" @click="assignDialog=false" />
        <Button :label="$t('assign')" @click="doAssign" :loading="assigning" />
      </template>
    </Dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import axios from 'axios'
import { useI18n } from 'vue-i18n'
import { useToast } from 'primevue/usetoast'
import auth from '../services/auth'
import Button from 'primevue/button'
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import Dialog from 'primevue/dialog'
import InputText from 'primevue/inputtext'
import InputNumber from 'primevue/inputnumber'
import Calendar from 'primevue/calendar'
import Textarea from 'primevue/textarea'
import Tag from 'primevue/tag'

const { t } = useI18n()
const toast = useToast()
const list = ref<any[]>([])
const loading = ref(false)
const dialog = ref(false)
const saving = ref(false)
const assignDialog = ref(false)
const assignTable = ref('')
const assignId = ref('')
const assigning = ref(false)
const formDate = ref<Date | null>(null)
const formTime = ref<Date | null>(null)
const form = reactive({ customer_name: '', customer_phone: '', customer_email: '', guest_count: 2, notes: '' })

const openNew = () => {
  form.customer_name = ''; form.customer_phone = ''; form.customer_email = ''; form.guest_count = 2; form.notes = ''
  formDate.value = null; formTime.value = null; dialog.value = true
}

const doSave = async () => {
  saving.value = true
  try {
    await axios.post(`http://${import.meta.env.VITE_APP_BACKEND_HOST}/reservations/api/reservations`, {
      customer_name: form.customer_name,
      customer_phone: form.customer_phone,
      customer_email: form.customer_email,
      guest_count: form.guest_count,
      reservation_date: formDate.value?.toISOString().split('T')[0] || '',
      reservation_time: formTime.value?.toTimeString().slice(0, 5) || '',
      notes: form.notes,
    })
    toast.add({ severity: 'success', summary: t('saved'), group: 'br', life: 2000 })
    dialog.value = false; await load()
  } catch { toast.add({ severity: 'error', summary: t('failed'), group: 'br', life: 3000 })
  } finally { saving.value = false }
}

const confirmRes = (r: any) => { assignId.value = r.id; assignTable.value = ''; assignDialog.value = true }

const doAssign = async () => {
  assigning.value = true
  try {
    await axios.post(`http://${import.meta.env.VITE_APP_BACKEND_HOST}/reservations/api/reservations/${assignId.value}/assign`,
      { table: assignTable.value }, { headers: { Authorization: `Bearer ${auth.accessToken.value}` } })
    toast.add({ severity: 'success', summary: t('saved'), group: 'br', life: 2000 })
    assignDialog.value = false; await load()
  } catch { toast.add({ severity: 'error', summary: t('failed'), group: 'br', life: 3000 })
  } finally { assigning.value = false }
}

const doDelete = async (r: any) => {
  try {
    await axios.delete(`http://${import.meta.env.VITE_APP_BACKEND_HOST}/reservations/api/reservations/${r.id}`,
      { headers: { Authorization: `Bearer ${auth.accessToken.value}` } })
    toast.add({ severity: 'success', summary: t('deleted'), group: 'br', life: 2000 }); await load()
  } catch { toast.add({ severity: 'error', summary: t('failed'), group: 'br', life: 3000 }) }
}

const load = async () => {
  loading.value = true
  try {
    const r = await axios.get(`http://${import.meta.env.VITE_APP_BACKEND_HOST}/reservations/api/reservations`,
      { headers: { Authorization: `Bearer ${auth.accessToken.value}` } })
    list.value = r.data.data || []
  } catch { toast.add({ severity: 'error', summary: t('load_failed'), group: 'br', life: 3000 })
  } finally { loading.value = false }
}
load()
</script>

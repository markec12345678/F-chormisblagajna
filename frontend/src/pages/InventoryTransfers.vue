<template>
  <div class="p-4">
    <div class="flex justify-content-between align-items-center mb-4">
      <h2 class="m-0">{{ $t('inventory_transfers') }}</h2>
      <Button :label="$t('create_transfer')" icon="pi pi-plus" @click="showCreateDialog = true" />
    </div>

    <div v-if="loading" class="flex justify-content-center p-8">
      <ProgressSpinner style="width: 50px; height: 50px" strokeWidth="6" />
    </div>

    <div v-else>
      <DataTable :value="transfers" responsiveLayout="scroll" stripedRows>
        <Column field="material_name" :header="$t('material')" sortable />
        <Column field="quantity" :header="$t('quantity')" sortable>
          <template #body="{ data }">{{ data.quantity }} {{ data.unit }}</template>
        </Column>
        <Column field="from_branch_name" :header="$t('from_branch')" sortable />
        <Column field="to_branch_name" :header="$t('to_branch')" sortable />
        <Column field="status" :header="$t('status')" sortable>
          <template #body="{ data }">
            <Tag :value="data.status" :severity="getStatusSeverity(data.status)" />
          </template>
        </Column>
        <Column field="created_at" :header="$t('date')" sortable>
          <template #body="{ data }">{{ new Date(data.created_at).toLocaleDateString() }}</template>
        </Column>
        <Column :header="$t('actions')">
          <template #body="{ data }">
            <div class="flex gap-1">
              <Button
                v-if="data.status === 'pending'"
                icon="pi pi-check"
                severity="success"
                text
                rounded
                @click="completeTransfer(data)"
              />
              <Button
                v-if="data.status === 'pending'"
                icon="pi pi-times"
                severity="danger"
                text
                rounded
                @click="cancelTransfer(data)"
              />
            </div>
          </template>
        </Column>
      </DataTable>
    </div>

    <!-- Create Dialog -->
    <Dialog
      v-model:visible="showCreateDialog"
      :header="$t('create_transfer')"
      :style="{ width: '500px' }"
      modal
    >
      <div class="flex flex-column gap-3">
        <div>
          <label class="block mb-1">{{ $t('material') }}</label>
          <InputText v-model="newTransfer.material_name" class="w-full" />
        </div>
        <div class="grid">
          <div class="col-6">
            <label class="block mb-1">{{ $t('quantity') }}</label>
            <InputNumber v-model="newTransfer.quantity" class="w-full" />
          </div>
          <div class="col-6">
            <label class="block mb-1">{{ $t('unit') }}</label>
            <InputText v-model="newTransfer.unit" class="w-full" />
          </div>
        </div>
        <div>
          <label class="block mb-1">{{ $t('from_branch') }}</label>
          <InputText v-model="newTransfer.from_branch_name" class="w-full" />
        </div>
        <div>
          <label class="block mb-1">{{ $t('to_branch') }}</label>
          <InputText v-model="newTransfer.to_branch_name" class="w-full" />
        </div>
        <div>
          <label class="block mb-1">{{ $t('notes') }}</label>
          <InputText v-model="newTransfer.notes" class="w-full" />
        </div>
      </div>
      <template #footer>
        <Button :label="$t('cancel')" severity="secondary" @click="showCreateDialog = false" />
        <Button :label="$t('create')" @click="createTransfer" :loading="creating" />
      </template>
    </Dialog>

    <Toast />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import Button from 'primevue/button'
import Dialog from 'primevue/dialog'
import InputText from 'primevue/inputtext'
import InputNumber from 'primevue/inputnumber'
import Tag from 'primevue/tag'
import ProgressSpinner from 'primevue/progressspinner'
import Toast from 'primevue/toast'
import { useToast } from 'primevue/usetoast'
import { useI18n } from 'vue-i18n'
import axios from 'axios'

const { t } = useI18n({ useScope: 'global' })
const toast = useToast()

const loading = ref(true)
const creating = ref(false)
const transfers = ref([])
const showCreateDialog = ref(false)

const newTransfer = ref({
  material_name: '',
  quantity: 0,
  unit: '',
  from_branch_name: '',
  to_branch_name: '',
  notes: '',
})

const apiBase = `http://${import.meta.env.VITE_APP_BACKEND_HOST}${import.meta.env.VITE_APP_MODULE_CORE_API_PREFIX}/inventorytransfer/api/transfers`

const getStatusSeverity = (status: string) => {
  switch (status) {
    case 'pending':
      return 'warn'
    case 'completed':
      return 'success'
    case 'cancelled':
      return 'danger'
    default:
      return 'secondary'
  }
}

const loadTransfers = async () => {
  loading.value = true
  try {
    const res = await axios.get(apiBase)
    transfers.value = res.data.data || []
  } catch {
    toast.add({ severity: 'error', summary: t('error'), detail: t('request_failed'), life: 3000 })
  } finally {
    loading.value = false
  }
}

const createTransfer = async () => {
  creating.value = true
  try {
    await axios.post(apiBase, newTransfer.value)
    toast.add({
      severity: 'success',
      summary: t('success'),
      detail: t('transfer_created'),
      life: 3000,
    })
    showCreateDialog.value = false
    newTransfer.value = {
      material_name: '',
      quantity: 0,
      unit: '',
      from_branch_name: '',
      to_branch_name: '',
      notes: '',
    }
    loadTransfers()
  } catch {
    toast.add({ severity: 'error', summary: t('error'), detail: t('request_failed'), life: 3000 })
  } finally {
    creating.value = false
  }
}

const completeTransfer = async (transfer) => {
  try {
    await axios.put(`${apiBase}/${transfer.id}/status`, { status: 'completed' })
    toast.add({
      severity: 'success',
      summary: t('success'),
      detail: t('transfer_completed'),
      life: 3000,
    })
    loadTransfers()
  } catch {
    toast.add({ severity: 'error', summary: t('error'), detail: t('request_failed'), life: 3000 })
  }
}

const cancelTransfer = async (transfer) => {
  try {
    await axios.put(`${apiBase}/${transfer.id}/status`, { status: 'cancelled' })
    toast.add({
      severity: 'warn',
      summary: t('warning'),
      detail: t('transfer_cancelled'),
      life: 3000,
    })
    loadTransfers()
  } catch {
    toast.add({ severity: 'error', summary: t('error'), detail: t('request_failed'), life: 3000 })
  }
}

onMounted(() => {
  loadTransfers()
})
</script>

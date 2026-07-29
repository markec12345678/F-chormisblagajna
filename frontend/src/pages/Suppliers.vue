<template>
  <div class="w-full">
    <div class="grid mx-2">
      <div class="col-12">
        <div class="flex justify-content-between align-items-center">
          <h3>{{ $t('suppliers') }}</h3>
          <Button :label="$t('add_supplier')" icon="pi pi-plus" @click="showAddDialog" />
        </div>
      </div>

      <div class="col-12">
        <DataTable :value="suppliers" :loading="isLoading" stripedRows>
          <template #empty>{{ $t('no_suppliers') }}</template>
          <Column field="name" :header="$t('name')" sortable></Column>
          <Column field="contact_name" :header="$t('contact_name')"></Column>
          <Column field="email" :header="$t('email')"></Column>
          <Column field="phone" :header="$t('phone')"></Column>
          <Column field="is_active" :header="$t('status')">
            <template #body="slotProps">
              <Tag :value="slotProps?.data?.is_active ? $t('active') : $t('inactive')" :severity="slotProps?.data?.is_active ? 'success' : 'danger'" />
            </template>
          </Column>
          <Column :header="$t('actions')">
            <template #body="slotProps">
              <ButtonGroup>
                <Button icon="pi pi-eye" severity="info" @click="viewSupplier(slotProps?.data)" />
                <Button icon="pi pi-pencil" severity="secondary" @click="editSupplier(slotProps?.data)" />
                <Button icon="pi pi-trash" severity="danger" @click="deleteSupplier(slotProps?.data?.id)" />
              </ButtonGroup>
            </template>
          </Column>
        </DataTable>
      </div>

      <Dialog v-model:visible="supplierDialog" modal :header="editingSupplier ? $t('edit_supplier') : $t('add_supplier')" :style="{ width: '50rem' }">
        <div class="grid">
          <div class="col-12 md:col-6">
            <div class="flex flex-column gap-2">
              <label>{{ $t('name') }}</label>
              <InputText v-model="currentSupplier.name" />
            </div>
          </div>
          <div class="col-12 md:col-6">
            <div class="flex flex-column gap-2">
              <label>{{ $t('contact_name') }}</label>
              <InputText v-model="currentSupplier.contact_name" />
            </div>
          </div>
          <div class="col-12 md:col-6">
            <div class="flex flex-column gap-2">
              <label>{{ $t('email') }}</label>
              <InputText v-model="currentSupplier.email" type="email" />
            </div>
          </div>
          <div class="col-12 md:col-6">
            <div class="flex flex-column gap-2">
              <label>{{ $t('phone') }}</label>
              <InputText v-model="currentSupplier.phone" />
            </div>
          </div>
          <div class="col-12">
            <div class="flex flex-column gap-2">
              <label>{{ $t('address') }}</label>
              <InputText v-model="currentSupplier.address" />
            </div>
          </div>
          <div class="col-12">
            <div class="flex flex-column gap-2">
              <label>{{ $t('website') }}</label>
              <InputText v-model="currentSupplier.website" />
            </div>
          </div>
          <div class="col-12">
            <div class="flex flex-column gap-2">
              <label>{{ $t('notes') }}</label>
              <Textarea v-model="currentSupplier.notes" rows="3" />
            </div>
          </div>
        </div>
        <template #footer>
          <Button :label="$t('cancel')" severity="secondary" @click="supplierDialog = false" />
          <Button :label="$t('save')" @click="saveSupplier" :loading="isSaving" />
        </template>
      </Dialog>

      <Dialog v-model:visible="ordersDialog" modal :header="$t('supplier_orders')" :style="{ width: '70rem' }">
        <DataTable :value="supplierOrders" stripedRows>
          <template #empty>{{ $t('no_orders') }}</template>
          <Column field="order_date" :header="$t('date')">
            <template #body="slotProps">
              {{ formatDate(slotProps?.data?.order_date) }}
            </template>
          </Column>
          <Column field="total_amount" :header="$t('total')">
            <template #body="slotProps">
              {{ formatCurrency(slotProps?.data?.total_amount || 0) }}
            </template>
          </Column>
          <Column field="status" :header="$t('status')">
            <template #body="slotProps">
              <Tag :value="slotProps?.data?.status" :severity="getStatusSeverity(slotProps?.data?.status)" />
            </template>
          </Column>
        </DataTable>
      </Dialog>
    </div>
  </div>
</template>

<script setup lang="ts">
import DataTable from 'primevue/datatable'
import Dialog from 'primevue/dialog'
import InputText from 'primevue/inputtext'
import Textarea from 'primevue/textarea'
import Column from 'primevue/column'
import Button from 'primevue/button'
import ButtonGroup from 'primevue/buttongroup'
import Tag from 'primevue/tag'
import axios from 'axios'
import { ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { useToast } from 'primevue/usetoast'
import auth from '../services/auth'

const { t } = useI18n()
const toast = useToast()

interface Supplier {
  id: string
  name: string
  contact_name: string
  email: string
  phone: string
  address: string
  website: string
  notes: string
  is_active: boolean
}

interface SupplierOrder {
  id: string
  order_date: string
  total_amount: number
  status: string
}

const suppliers = ref<Supplier[]>([])
const supplierOrders = ref<SupplierOrder[]>([])
const isLoading = ref(false)
const isSaving = ref(false)
const supplierDialog = ref(false)
const ordersDialog = ref(false)
const editingSupplier = ref(false)
const currentSupplier = ref<Partial<Supplier>>({})

const formatCurrency = (amount: number) => {
  return new Intl.NumberFormat('sl-SI', { style: 'currency', currency: 'EUR' }).format(amount)
}

const formatDate = (date: string) => {
  if (!date) return '-'
  return new Date(date).toLocaleDateString()
}

const getStatusSeverity = (status: string) => {
  const map: Record<string, string> = {
    delivered: 'success',
    pending: 'warn',
    cancelled: 'danger',
  }
  return map[status] || 'secondary'
}

const showAddDialog = () => {
  editingSupplier.value = false
  currentSupplier.value = {}
  supplierDialog.value = true
}

const editSupplier = (supplier: Supplier) => {
  editingSupplier.value = true
  currentSupplier.value = { ...supplier }
  supplierDialog.value = true
}

const viewSupplier = async (supplier: Supplier) => {
  try {
    const response = await axios.get(
      `http://${import.meta.env.VITE_APP_BACKEND_HOST}/supplier/api/suppliers/${supplier.id}/orders`,
      { headers: { Authorization: `Bearer ${auth.accessToken.value}` } }
    )
    supplierOrders.value = response.data.data || []
    ordersDialog.value = true
  } catch {
    toast.add({ severity: 'error', summary: t('failed'), detail: t('load_failed'), group: 'br', life: 3000 })
  }
}

const saveSupplier = async () => {
  isSaving.value = true
  try {
    if (editingSupplier.value && currentSupplier.value.id) {
      await axios.put(
        `http://${import.meta.env.VITE_APP_BACKEND_HOST}/supplier/api/suppliers/${currentSupplier.value.id}`,
        currentSupplier.value,
        { headers: { Authorization: `Bearer ${auth.accessToken.value}` } }
      )
    } else {
      await axios.post(
        `http://${import.meta.env.VITE_APP_BACKEND_HOST}/supplier/api/suppliers`,
        currentSupplier.value,
        { headers: { Authorization: `Bearer ${auth.accessToken.value}` } }
      )
    }
    toast.add({ severity: 'success', summary: t('success'), detail: t('supplier_saved'), group: 'br', life: 3000 })
    supplierDialog.value = false
    loadSuppliers()
  } catch {
    toast.add({ severity: 'error', summary: t('failed'), detail: t('save_failed'), group: 'br', life: 3000 })
  } finally {
    isSaving.value = false
  }
}

const deleteSupplier = async (id: string) => {
  try {
    await axios.delete(
      `http://${import.meta.env.VITE_APP_BACKEND_HOST}/supplier/api/suppliers/${id}`,
      { headers: { Authorization: `Bearer ${auth.accessToken.value}` } }
    )
    toast.add({ severity: 'success', summary: t('success'), detail: t('supplier_deleted'), group: 'br', life: 3000 })
    loadSuppliers()
  } catch {
    toast.add({ severity: 'error', summary: t('failed'), detail: t('delete_failed'), group: 'br', life: 3000 })
  }
}

const loadSuppliers = async () => {
  isLoading.value = true
  try {
    const response = await axios.get(
      `http://${import.meta.env.VITE_APP_BACKEND_HOST}/supplier/api/suppliers`,
      { headers: { Authorization: `Bearer ${auth.accessToken.value}` } }
    )
    suppliers.value = response.data.data || []
  } catch {
    toast.add({ severity: 'error', summary: t('failed'), detail: t('load_failed'), group: 'br', life: 3000 })
  } finally {
    isLoading.value = false
  }
}

loadSuppliers()
</script>

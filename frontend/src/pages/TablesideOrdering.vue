<template>
  <div class="w-full">
    <div class="grid mx-2">
      <div class="col-12">
        <div class="flex justify-content-between align-items-center">
          <h3>{{ $t('tableside_ordering') }}</h3>
          <Button icon="pi pi-plus" :label="$t('add_table')" @click="openNewDialog" />
        </div>
      </div>

      <div class="col-12">
        <Card>
          <template #title>{{ $t('table_sessions') }}</template>
          <template #content>
            <DataTable :value="sessions" stripedRows :loading="isLoading">
              <template #empty>{{ $t('no_sessions') }}</template>
              <Column field="table_label" :header="$t('table')"></Column>
              <Column field="zone" :header="$t('zone')"></Column>
              <Column field="guest_count" :header="$t('guests')"></Column>
              <Column field="active" :header="$t('status')">
                <template #body="slotProps">
                  <Tag :severity="slotProps?.data?.active ? 'success' : 'secondary'">
                    {{ slotProps?.data?.active ? $t('active') : $t('closed') }}
                  </Tag>
                </template>
              </Column>
              <Column field="opened_at" :header="$t('opened')">
                <template #body="slotProps">
                  {{ slotProps?.data?.opened_at ? new Date(slotProps.data.opened_at).toLocaleString() : '-' }}
                </template>
              </Column>
              <Column :header="$t('actions')">
                <template #body="slotProps">
                  <Button icon="pi pi-qrcode" severity="info" class="mr-1" @click="showQr(slotProps.data)" v-tooltip.left="$t('qr_code')" />
                  <Button icon="pi pi-receipt" class="mr-1" @click="viewOrders(slotProps.data)" v-tooltip.left="$t('view_orders')" />
                  <Button v-if="slotProps.data.active"
                    icon="pi pi-times" severity="danger" @click="closeSession(slotProps.data)" v-tooltip.left="$t('close')" />
                </template>
              </Column>
            </DataTable>
          </template>
        </Card>
      </div>
    </div>

    <Dialog v-model:visible="dialogVisible" :header="$t('add_table')" :style="{ width: '400px' }">
      <div class="grid">
        <div class="col-12">
          <label>{{ $t('table_label') }}</label>
          <InputText v-model="form.table_label" class="w-full" />
        </div>
        <div class="col-12">
          <label>{{ $t('zone') }}</label>
          <Dropdown v-model="form.zone" :options="zoneOptions" class="w-full" />
        </div>
        <div class="col-12">
          <label>{{ $t('guest_count') }}</label>
          <InputNumber v-model="form.guest_count" class="w-full" :min="1" />
        </div>
      </div>
      <template #footer>
        <Button :label="$t('cancel')" severity="secondary" @click="dialogVisible = false" />
        <Button :label="$t('save')" @click="createSession" :loading="saving" />
      </template>
    </Dialog>

    <Dialog v-model:visible="qrDialogVisible" :header="$t('qr_code')" :style="{ width: '400px' }">
      <div class="text-center" v-if="qrInfo">
        <p><strong>{{ qrInfo.table_label }}</strong></p>
        <div class="qr-placeholder my-3" style="width: 200px; height: 200px; background: #f0f0f0; margin: 0 auto; display: flex; align-items: center; justify-content: center; border-radius: 8px;">
          <i class="pi pi-qrcode" style="font-size: 5rem; color: #666;"></i>
        </div>
        <p class="text-sm">{{ qrInfo.url }}</p>
      </div>
      <template #footer>
        <Button :label="$t('close')" severity="secondary" @click="qrDialogVisible = false" />
      </template>
    </Dialog>

    <Dialog v-model:visible="ordersDialogVisible" :header="`${$t('orders')} - ${selectedTable}`" :style="{ width: '600px' }">
      <DataTable :value="tableOrders" stripedRows>
        <template #empty>{{ $t('no_orders') }}</template>
        <Column field="status" :header="$t('status')">
          <template #body="slotProps">
            <Tag :severity="statusSeverity(slotProps?.data?.status)">{{ slotProps?.data?.status }}</Tag>
          </template>
        </Column>
        <Column field="items" :header="$t('items')">
          <template #body="slotProps">
            <div v-for="(item, idx) in slotProps?.data?.items" :key="idx">
              {{ item.quantity }}x {{ item.product_name }}
            </div>
          </template>
        </Column>
        <Column field="subtotal" :header="$t('total')">
          <template #body="slotProps">{{ formatCurrency(slotProps?.data?.subtotal) }}</template>
        </Column>
        <Column field="placed_at" :header="$t('time')">
          <template #body="slotProps">
            {{ slotProps?.data?.placed_at ? new Date(slotProps.data.placed_at).toLocaleTimeString() : '-' }}
          </template>
        </Column>
      </DataTable>
      <template #footer>
        <Button :label="$t('close')" severity="secondary" @click="ordersDialogVisible = false" />
      </template>
    </Dialog>
  </div>
</template>

<script setup lang="ts">
import Card from 'primevue/card'
import Button from 'primevue/button'
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import Dialog from 'primevue/dialog'
import InputText from 'primevue/inputtext'
import InputNumber from 'primevue/inputnumber'
import Dropdown from 'primevue/dropdown'
import Tag from 'primevue/tag'
import axios from 'axios'
import { ref, reactive } from 'vue'
import { useI18n } from 'vue-i18n'
import { useToast } from 'primevue/usetoast'
import auth from '../services/auth'

const { t } = useI18n()
const toast = useToast()

interface TableSession {
  id?: string
  table_label: string
  zone: string
  guest_count: number
  active: boolean
  opened_at?: string
  qr_token?: string
}

interface QrInfo {
  table_label: string
  token: string
  url: string
  host: string
}

const sessions = ref<TableSession[]>([])
const isLoading = ref(false)
const dialogVisible = ref(false)
const qrDialogVisible = ref(false)
const ordersDialogVisible = ref(false)
const saving = ref(false)
const qrInfo = ref<QrInfo | null>(null)
const tableOrders = ref<any[]>([])
const selectedTable = ref('')

const zoneOptions = ['Main Hall', 'Terrace', 'VIP', 'Bar', 'Garden']

const form = reactive({
  table_label: '',
  zone: 'Main Hall',
  guest_count: 2,
})

const formatCurrency = (amount: number) => {
  return new Intl.NumberFormat('sl-SI', { style: 'currency', currency: 'EUR' }).format(amount)
}

const statusSeverity = (status: string) => {
  switch (status) {
    case 'pending': return 'warn'
    case 'preparing': return 'info'
    case 'ready': return 'success'
    case 'delivered': return 'secondary'
    default: return 'info'
  }
}

const openNewDialog = () => {
  form.table_label = ''
  form.zone = 'Main Hall'
  form.guest_count = 2
  dialogVisible.value = true
}

const createSession = async () => {
  saving.value = true
  try {
    await axios.post(
      `http://${import.meta.env.VITE_APP_BACKEND_HOST}/tableside/api/sessions`,
      form,
      { headers: { Authorization: `Bearer ${auth.accessToken.value}` } }
    )
    toast.add({ severity: 'success', summary: t('saved'), group: 'br', life: 2000 })
    dialogVisible.value = false
    await loadSessions()
  } catch {
    toast.add({ severity: 'error', summary: t('failed'), group: 'br', life: 3000 })
  } finally {
    saving.value = false
  }
}

const showQr = async (session: TableSession) => {
  try {
    const host = window.location.origin
    const response = await axios.get(
      `http://${import.meta.env.VITE_APP_BACKEND_HOST}/tableside/api/qr/${session.id}?host=${host}`,
      { headers: { Authorization: `Bearer ${auth.accessToken.value}` } }
    )
    qrInfo.value = response.data.data
    qrDialogVisible.value = true
  } catch {
    toast.add({ severity: 'error', summary: t('load_failed'), group: 'br', life: 3000 })
  }
}

const viewOrders = async (session: TableSession) => {
  try {
    selectedTable.value = session.table_label
    const response = await axios.get(
      `http://${import.meta.env.VITE_APP_BACKEND_HOST}/tableside/api/sessions/${session.id}/orders`,
      { headers: { Authorization: `Bearer ${auth.accessToken.value}` } }
    )
    tableOrders.value = response.data.data || []
    ordersDialogVisible.value = true
  } catch {
    toast.add({ severity: 'error', summary: t('load_failed'), group: 'br', life: 3000 })
  }
}

const closeSession = async (session: TableSession) => {
  try {
    await axios.post(
      `http://${import.meta.env.VITE_APP_BACKEND_HOST}/tableside/api/sessions/${session.id}/close`,
      {},
      { headers: { Authorization: `Bearer ${auth.accessToken.value}` } }
    )
    toast.add({ severity: 'success', summary: t('saved'), group: 'br', life: 2000 })
    await loadSessions()
  } catch {
    toast.add({ severity: 'error', summary: t('failed'), group: 'br', life: 3000 })
  }
}

const loadSessions = async () => {
  isLoading.value = true
  try {
    const response = await axios.get(
      `http://${import.meta.env.VITE_APP_BACKEND_HOST}/tableside/api/sessions`,
      { headers: { Authorization: `Bearer ${auth.accessToken.value}` } }
    )
    sessions.value = response.data.data || []
  } catch {
    toast.add({ severity: 'error', summary: t('failed'), detail: t('load_failed'), group: 'br', life: 3000 })
  } finally {
    isLoading.value = false
  }
}

loadSessions()
</script>

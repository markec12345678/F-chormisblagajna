<template>
  <div class="w-full">
    <div class="grid mx-2">
      <div class="col-12">
        <h3>{{ $t('purchase_orders') }}</h3>
      </div>
      <div class="col-12">
        <DataTable :value="list" stripedRows :loading="loading">
          <template #empty>{{ $t('no_purchase_orders') }}</template>
          <Column field="supplier_name" :header="$t('supplier')"></Column>
          <Column field="items" :header="$t('items')"
            ><template #body="s"
              ><span v-for="(it, i) in s.data.items" :key="i"
                >{{ it.quantity }}x {{ it.material_name }}{{ sep(i, s.data.items) }}</span
              ></template
            ></Column
          >
          <Column field="total_cost" :header="$t('total')"
            ><template #body="s">{{ formatCurrency(s.data.total_cost) }}</template></Column
          >
          <Column field="status" :header="$t('status')"
            ><template #body="s"
              ><Tag
                :severity="
                  s.data.status === 'received'
                    ? 'success'
                    : s.data.status === 'cancelled'
                      ? 'danger'
                      : 'warn'
                "
                >{{ s.data.status }}</Tag
              ></template
            ></Column
          >
          <Column field="ordered_at" :header="$t('date')"
            ><template #body="s">{{
              s.data.ordered_at ? new Date(s.data.ordered_at).toLocaleDateString() : '-'
            }}</template></Column
          >
          <Column :header="$t('actions')"
            ><template #body="s">
              <Button
                v-if="s.data.status === 'pending'"
                icon="pi pi-check"
                severity="success"
                size="small"
                class="mr-1"
                @click="markReceived(s.data.id)"
              />
              <Button
                v-if="s.data.status === 'pending'"
                icon="pi pi-times"
                severity="danger"
                size="small"
                @click="markCancelled(s.data.id)"
              /> </template
          ></Column>
        </DataTable>
      </div>
      <div class="col-12 mt-3">
        <Card
          ><template #title>{{ $t('new_purchase_order') }}</template
          ><template #content>
            <div class="grid">
              <div class="col-4">
                <label>{{ $t('supplier') }}</label
                ><InputText v-model="form.supplier_name" class="w-full" />
              </div>
              <div class="col-4">
                <label>{{ $t('expected_date') }}</label
                ><Calendar v-model="formDate" dateFormat="yy-mm-dd" class="w-full" />
              </div>
              <div class="col-4">
                <label>{{ $t('notes') }}</label
                ><InputText v-model="form.notes" class="w-full" />
              </div>
            </div>
            <div class="mt-2">
              <Button
                :label="$t('add_item')"
                icon="pi pi-plus"
                size="small"
                @click="
                  form.items.push({ material_name: '', quantity: 1, unit_price: 0, total_price: 0 })
                "
              />
            </div>
            <div
              v-for="(it, idx) in form.items"
              :key="idx"
              class="flex gap-2 mt-2 align-items-center"
            >
              <InputText v-model="it.material_name" placeholder="Material" />
              <InputNumber v-model="it.quantity" :min="1" style="width: 80px" />
              <InputNumber
                v-model="it.unit_price"
                :min="0"
                style="width: 100px"
                :placeholder="$t('price')"
              />
              <Button
                icon="pi pi-trash"
                severity="danger"
                size="small"
                @click="form.items.splice(idx, 1)"
              />
            </div>
            <Button
              :label="$t('create_order')"
              class="mt-2"
              @click="createPO"
              :loading="saving"
            /> </template
        ></Card>
      </div>
    </div>
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
import Card from 'primevue/card'
import InputText from 'primevue/inputtext'
import InputNumber from 'primevue/inputnumber'
import Calendar from 'primevue/calendar'
import Tag from 'primevue/tag'
const { t } = useI18n()
const toast = useToast()
interface POItem {
  material_name: string
  quantity: number
  unit_price: number
  total_price: number
}
interface PurchaseOrder {
  id: string
  supplier_name: string
  items: POItem[]
  total_cost: number
  status: string
  ordered_at: string
}
const list = ref<PurchaseOrder[]>([])
const loading = ref(false)
const saving = ref(false)
const formDate = ref<Date | null>(null)
const form = reactive({ supplier_name: '', notes: '', items: [] as POItem[] })
const formatCurrency = (n: number) =>
  new Intl.NumberFormat('sl-SI', { style: 'currency', currency: 'EUR' }).format(n || 0)
const sep = (i: number, items: POItem[]) => (i < items.length - 1 ? ', ' : '')
const createPO = async () => {
  saving.value = true
  try {
    const items = form.items.map((i) => ({ ...i, total_price: i.quantity * i.unit_price }))
    await axios.post(
      `http://${import.meta.env.VITE_APP_BACKEND_HOST}/purchase/api/orders`,
      {
        supplier_name: form.supplier_name,
        items,
        expected_at: formDate.value?.toISOString().split('T')[0],
        notes: form.notes,
      },
      { headers: { Authorization: `Bearer ${auth.accessToken.value}` } },
    )
    toast.add({ severity: 'success', summary: t('saved'), group: 'br', life: 2000 })
    form.supplier_name = ''
    form.notes = ''
    form.items = []
    formDate.value = null
    await load()
  } catch {
    toast.add({ severity: 'error', summary: t('failed'), group: 'br', life: 3000 })
  } finally {
    saving.value = false
  }
}
const markReceived = async (id: string) => {
  try {
    await axios.post(
      `http://${import.meta.env.VITE_APP_BACKEND_HOST}/purchase/api/orders/${id}/receive`,
      {},
      { headers: { Authorization: `Bearer ${auth.accessToken.value}` } },
    )
    await load()
  } catch {
    toast.add({ severity: 'error', summary: t('failed'), group: 'br', life: 3000 })
  }
}
const markCancelled = async (id: string) => {
  try {
    await axios.post(
      `http://${import.meta.env.VITE_APP_BACKEND_HOST}/purchase/api/orders/${id}/cancel`,
      {},
      { headers: { Authorization: `Bearer ${auth.accessToken.value}` } },
    )
    await load()
  } catch {
    toast.add({ severity: 'error', summary: t('failed'), group: 'br', life: 3000 })
  }
}
const load = async () => {
  loading.value = true
  try {
    const r = await axios.get(
      `http://${import.meta.env.VITE_APP_BACKEND_HOST}/purchase/api/orders`,
      { headers: { Authorization: `Bearer ${auth.accessToken.value}` } },
    )
    list.value = r.data.data || []
  } catch {
    toast.add({ severity: 'error', summary: t('load_failed'), group: 'br', life: 3000 })
  } finally {
    loading.value = false
  }
}
load()
</script>

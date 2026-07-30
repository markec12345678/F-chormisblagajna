<template>
  <div class="w-full">
    <div class="grid mx-2">
      <div class="col-12">
        <h3>{{ $t('floor_plan') }}</h3>
      </div>
      <div class="col-12">
        <div class="flex gap-2 mb-2">
          <InputText v-model="newLabel" :placeholder="$t('table_label')" /><Button
            icon="pi pi-plus"
            :label="$t('add_table')"
            @click="addTable"
          />
        </div>
      </div>
      <div class="col-12">
        <div
          class="floor-canvas"
          style="
            position: relative;
            min-height: 400px;
            background: #f8f9fa;
            border-radius: 8px;
            border: 1px solid #ddd;
          "
        >
          <div
            v-for="t in tables"
            :key="t.id"
            class="table-obj"
            :style="{
              left: t.x + 'px',
              top: t.y + 'px',
              width: (t.width || 60) + 'px',
              height: (t.height || 60) + 'px',
              background: t.status === 'available' ? '#e8f5e9' : '#ffebee',
              border: '2px solid ' + (t.status === 'available' ? '#4caf50' : '#f44336'),
              borderRadius: t.shape === 'circle' ? '50%' : '8px',
              position: 'absolute',
              display: 'flex',
              alignItems: 'center',
              justifyContent: 'center',
              cursor: 'pointer',
              fontWeight: 'bold',
            }"
            @click="toggleStatus(t)"
          >
            {{ t.label }}
            <Button
              icon="pi pi-trash"
              severity="danger"
              size="small"
              text
              style="position: absolute; top: -12px; right: -12px"
              @click.stop="deleteTable(t.id)"
            />
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
<script setup lang="ts">
import { ref, onMounted } from 'vue'
import axios from 'axios'
import { useI18n } from 'vue-i18n'
import { useToast } from 'primevue/usetoast'
import auth from '../services/auth'
import Button from 'primevue/button'
import InputText from 'primevue/inputtext'
const { t } = useI18n()
const toast = useToast()
interface TablePlan { id: string; label: string; zone: string; capacity: number; shape: string; x: number; y: number; width: number; height: number; status: string }
const tables = ref<TablePlan[]>([])
const newLabel = ref('')
const addTable = async () => {
  if (!newLabel.value) return
  try {
    await axios.post(
      `http://${import.meta.env.VITE_APP_BACKEND_HOST}/floorplan/api/tables`,
      {
        label: newLabel.value,
        zone: 'Main',
        capacity: 4,
        shape: 'rect',
        x: Math.random() * 300,
        y: Math.random() * 300,
        width: 60,
        height: 60,
        status: 'available',
      },
      { headers: { Authorization: `Bearer ${auth.accessToken.value}` } },
    )
    newLabel.value = ''
    await load()
  } catch {
    toast.add({ severity: 'error', summary: t('failed'), group: 'br', life: 3000 })
  }
}
const toggleStatus = async (t: TablePlan) => {
  try {
    await axios.post(
      `http://${import.meta.env.VITE_APP_BACKEND_HOST}/floorplan/api/tables`,
      { ...t, status: t.status === 'available' ? 'occupied' : 'available' },
      { headers: { Authorization: `Bearer ${auth.accessToken.value}` } },
    )
    await load()
  } catch {
    toast.add({ severity: 'error', summary: t('failed'), group: 'br', life: 3000 })
  }
}
const deleteTable = async (id: string) => {
  try {
    await axios.delete(
      `http://${import.meta.env.VITE_APP_BACKEND_HOST}/floorplan/api/tables/${id}`,
      { headers: { Authorization: `Bearer ${auth.accessToken.value}` } },
    )
    await load()
  } catch {
    toast.add({ severity: 'error', summary: t('failed'), group: 'br', life: 3000 })
  }
}
const load = async () => {
  try {
    const r = await axios.get(
      `http://${import.meta.env.VITE_APP_BACKEND_HOST}/floorplan/api/tables`,
      { headers: { Authorization: `Bearer ${auth.accessToken.value}` } },
    )
    tables.value = r.data.data || []
  } catch {
    toast.add({ severity: 'error', summary: t('load_failed'), group: 'br', life: 3000 })
  }
}
onMounted(load)
</script>

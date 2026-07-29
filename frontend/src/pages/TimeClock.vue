<template>
  <div class="w-full">
    <div class="grid mx-2">
      <div class="col-12">
        <h3>{{ $t('time_clock') }}</h3>
      </div>

      <div class="col-12 md:col-4">
        <Card>
          <template #title>{{ $t('clock_in_out') }}</template>
          <template #content>
            <div class="flex flex-column gap-3">
              <div class="flex flex-column gap-2">
                <label>{{ $t('employee') }}</label>
                <Dropdown v-model="selectedEmployee" :options="employees" optionLabel="name" :placeholder="$t('select_employee')" />
              </div>
              <div class="flex flex-column gap-2">
                <label>{{ $t('notes') }}</label>
                <InputText v-model="clockNotes" />
              </div>
              <Button
                v-if="!isCurrentlyClockedIn"
                :label="$t('clock_in')"
                icon="pi pi-sign-in"
                severity="success"
                class="w-full"
                @click="clockIn"
                :loading="isClocking"
              />
              <Button
                v-else
                :label="$t('clock_out')"
                icon="pi pi-sign-out"
                severity="danger"
                class="w-full"
                @click="clockOut"
                :loading="isClocking"
              />
            </div>
          </template>
        </Card>
      </div>

      <div class="col-12 md:col-8">
        <Card>
          <template #title>{{ $t('currently_clocked_in') }}</template>
          <template #content>
            <DataTable :value="dashboard?.currently_clocked_in || []" stripedRows>
              <template #empty>{{ $t('no_active_shifts') }}</template>
              <Column field="employee_name" :header="$t('employee')"></Column>
              <Column field="clock_in" :header="$t('clock_in_time')">
                <template #body="slotProps">
                  {{ formatTime(slotProps?.data?.clock_in) }}
                </template>
              </Column>
              <Column field="hours_worked" :header="$t('hours')">
                <template #body="slotProps">
                  {{ calculateHours(slotProps?.data?.clock_in) }}
                </template>
              </Column>
              <Column :header="$t('actions')">
                <template #body="slotProps">
                  <Button :label="$t('clock_out')" severity="danger" size="small" @click="clockOutEntry(slotProps?.data?.id)" />
                </template>
              </Column>
            </DataTable>
          </template>
        </Card>
      </div>

      <div class="col-12">
        <Card>
          <template #title>{{ $t('today_summary') }}</template>
          <template #content>
            <DataTable :value="dashboard?.today_summary || []" stripedRows>
              <template #empty>{{ $t('no_data') }}</template>
              <Column field="employee_name" :header="$t('employee')"></Column>
              <Column field="total_hours" :header="$t('total_hours')">
                <template #body="slotProps">
                  {{ (slotProps?.data?.total_hours || 0).toFixed(1) }}h
                </template>
              </Column>
              <Column field="shift_count" :header="$t('shifts')"></Column>
              <Column field="avg_hours_per_shift" :header="$t('avg_per_shift')">
                <template #body="slotProps">
                  {{ (slotProps?.data?.avg_hours_per_shift || 0).toFixed(1) }}h
                </template>
              </Column>
              <Column field="overtime_hours" :header="$t('overtime')">
                <template #body="slotProps">
                  <Tag v-if="(slotProps?.data?.overtime_hours || 0) > 0" :value="`${(slotProps?.data?.overtime_hours || 0).toFixed(1)}h`" severity="warn" />
                  <span v-else>-</span>
                </template>
              </Column>
            </DataTable>
          </template>
        </Card>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import Card from 'primevue/card'
import Button from 'primevue/button'
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import Dialog from 'primevue/dialog'
import InputText from 'primevue/inputtext'
import Dropdown from 'primevue/dropdown'
import Tag from 'primevue/tag'
import axios from 'axios'
import { ref, computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { useToast } from 'primevue/usetoast'
import auth from '../services/auth'

const { t } = useI18n()
const toast = useToast()

interface ClockEntry {
  id: string
  employee_id: string
  employee_name: string
  clock_in: string
  status: string
  hours_worked: number
}

interface TimeClockSummary {
  employee_id: string
  employee_name: string
  total_hours: number
  shift_count: number
  avg_hours_per_shift: number
  overtime_hours: number
}

interface Dashboard {
  currently_clocked_in: ClockEntry[]
  today_summary: TimeClockSummary[]
  week_summary: TimeClockSummary[]
}

const dashboard = ref<Dashboard | null>(null)
const selectedEmployee = ref<{ id: string; name: string } | null>(null)
const clockNotes = ref('')
const isClocking = ref(false)
const isLoading = ref(false)

const employees = [
  { id: 'emp-1', name: 'Janez Novak' },
  { id: 'emp-2', name: 'Maria Skršnik' },
  { id: 'emp-3', name: 'Peter Horvat' },
]

const isCurrentlyClockedIn = computed(() => {
  if (!selectedEmployee.value) return false
  return dashboard.value?.currently_clocked_in?.some(
    (e) => e.employee_id === selectedEmployee.value?.id
  ) || false
})

const formatTime = (dateStr: string) => {
  if (!dateStr) return '-'
  return new Date(dateStr).toLocaleTimeString()
}

const calculateHours = (clockIn: string) => {
  if (!clockIn) return '0h'
  const hours = (Date.now() - new Date(clockIn).getTime()) / (1000 * 60 * 60)
  return `${hours.toFixed(1)}h`
}

const clockIn = async () => {
  if (!selectedEmployee.value) {
    toast.add({ severity: 'warn', summary: t('warning'), detail: t('select_employee'), group: 'br', life: 3000 })
    return
  }

  isClocking.value = true
  try {
    await axios.post(
      `http://${import.meta.env.VITE_APP_BACKEND_HOST}/timeclock/api/clock-in`,
      {
        employee_id: selectedEmployee.value.id,
        employee_name: selectedEmployee.value.name,
        notes: clockNotes.value,
      },
      { headers: { Authorization: `Bearer ${auth.accessToken.value}` } }
    )
    toast.add({ severity: 'success', summary: t('success'), detail: t('clocked_in'), group: 'br', life: 3000 })
    clockNotes.value = ''
    loadDashboard()
  } catch {
    toast.add({ severity: 'error', summary: t('failed'), detail: t('clock_in_failed'), group: 'br', life: 3000 })
  } finally {
    isClocking.value = false
  }
}

const clockOut = async () => {
  if (!selectedEmployee.value) return

  const entry = dashboard.value?.currently_clocked_in?.find(
    (e) => e.employee_id === selectedEmployee.value?.id
  )
  if (!entry) return

  isClocking.value = true
  try {
    await axios.post(
      `http://${import.meta.env.VITE_APP_BACKEND_HOST}/timeclock/api/clock-out/${entry.id}`,
      {},
      { headers: { Authorization: `Bearer ${auth.accessToken.value}` } }
    )
    toast.add({ severity: 'success', summary: t('success'), detail: t('clocked_out'), group: 'br', life: 3000 })
    loadDashboard()
  } catch {
    toast.add({ severity: 'error', summary: t('failed'), detail: t('clock_out_failed'), group: 'br', life: 3000 })
  } finally {
    isClocking.value = false
  }
}

const clockOutEntry = async (entryId: string) => {
  try {
    await axios.post(
      `http://${import.meta.env.VITE_APP_BACKEND_HOST}/timeclock/api/clock-out/${entryId}`,
      {},
      { headers: { Authorization: `Bearer ${auth.accessToken.value}` } }
    )
    toast.add({ severity: 'success', summary: t('success'), detail: t('clocked_out'), group: 'br', life: 3000 })
    loadDashboard()
  } catch {
    toast.add({ severity: 'error', summary: t('failed'), detail: t('clock_out_failed'), group: 'br', life: 3000 })
  }
}

const loadDashboard = async () => {
  isLoading.value = true
  try {
    const response = await axios.get(
      `http://${import.meta.env.VITE_APP_BACKEND_HOST}/timeclock/api/dashboard`,
      { headers: { Authorization: `Bearer ${auth.accessToken.value}` } }
    )
    dashboard.value = response.data.data
  } catch {
    toast.add({ severity: 'error', summary: t('failed'), detail: t('load_failed'), group: 'br', life: 3000 })
  } finally {
    isLoading.value = false
  }
}

loadDashboard()
</script>

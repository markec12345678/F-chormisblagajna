<template>
  <div class="w-full">
    <div class="grid mx-2">
      <div class="col-12">
        <div class="flex justify-content-between align-items-center">
          <h3>{{ $t('inventory_alerts') }}</h3>
          <Button :label="$t('add_rule')" icon="pi pi-plus" @click="showAddRule" />
        </div>
      </div>

      <div class="col-12" v-if="summary">
        <div class="grid">
          <div class="col-12 md:col-3">
            <Card>
              <template #title>{{ $t('active_rules') }}</template>
              <template #content>
                <div class="text-4xl font-bold text-primary">{{ summary.total_active }}</div>
              </template>
            </Card>
          </div>
          <div class="col-12 md:col-3">
            <Card>
              <template #title>{{ $t('unread_alerts') }}</template>
              <template #content>
                <div class="text-4xl font-bold text-orange-500">{{ summary.unread_count }}</div>
              </template>
            </Card>
          </div>
          <div class="col-12 md:col-3">
            <Card>
              <template #title>{{ $t('critical') }}</template>
              <template #content>
                <div class="text-4xl font-bold text-red-500">{{ summary.critical_count }}</div>
              </template>
            </Card>
          </div>
          <div class="col-12 md:col-3">
            <Card>
              <template #title>{{ $t('low_stock') }}</template>
              <template #content>
                <div class="text-4xl font-bold text-yellow-500">{{ summary.low_count }}</div>
              </template>
            </Card>
          </div>
        </div>
      </div>

      <div class="col-12 md:col-6">
        <Card>
          <template #title>{{ $t('alert_rules') }}</template>
          <template #content>
            <DataTable :value="rules" stripedRows :loading="isLoading">
              <template #empty>{{ $t('no_rules') }}</template>
              <Column field="material_name" :header="$t('material')"></Column>
              <Column field="threshold_low" :header="$t('low_threshold')">
                <template #body="slotProps">
                  {{ (slotProps?.data?.threshold_low || 0).toFixed(1) }}
                </template>
              </Column>
              <Column field="threshold_critical" :header="$t('critical_threshold')">
                <template #body="slotProps">
                  {{ (slotProps?.data?.threshold_critical || 0).toFixed(1) }}
                </template>
              </Column>
              <Column field="is_active" :header="$t('status')">
                <template #body="slotProps">
                  <Tag
                    :value="slotProps?.data?.is_active ? $t('active') : $t('inactive')"
                    :severity="slotProps?.data?.is_active ? 'success' : 'secondary'"
                  />
                </template>
              </Column>
              <Column :header="$t('actions')">
                <template #body="slotProps">
                  <Button
                    icon="pi pi-trash"
                    severity="danger"
                    text
                    rounded
                    size="small"
                    @click="deleteRule(slotProps?.data?.id)"
                  />
                </template>
              </Column>
            </DataTable>
          </template>
        </Card>
      </div>

      <div class="col-12 md:col-6">
        <Card>
          <template #title>{{ $t('recent_alerts') }}</template>
          <template #content>
            <DataTable :value="alerts" stripedRows :loading="isLoading">
              <template #empty>{{ $t('no_alerts') }}</template>
              <Column field="created_at" :header="$t('time')">
                <template #body="slotProps">
                  {{ formatTime(slotProps?.data?.created_at) }}
                </template>
              </Column>
              <Column field="material_name" :header="$t('material')"></Column>
              <Column field="severity" :header="$t('severity')">
                <template #body="slotProps">
                  <Tag
                    :value="slotProps?.data?.severity"
                    :severity="getSeverity(slotProps?.data?.severity)"
                  />
                </template>
              </Column>
              <Column field="current_qty" :header="$t('current')">
                <template #body="slotProps">
                  {{ (slotProps?.data?.current_qty || 0).toFixed(1) }}
                </template>
              </Column>
              <Column field="is_read" :header="$t('read')">
                <template #body="slotProps">
                  <Button
                    v-if="!slotProps?.data?.is_read"
                    icon="pi pi-check"
                    text
                    rounded
                    size="small"
                    severity="success"
                    @click="markAsRead(slotProps?.data?.id)"
                  />
                  <span v-else class="text-400"><i class="pi pi-check-circle"></i></span>
                </template>
              </Column>
            </DataTable>
          </template>
        </Card>
      </div>

      <div class="col-12">
        <Card>
          <template #title>{{ $t('add_alert_rule') }}</template>
          <template #content>
            <div class="grid">
              <div class="col-12 md:col-4">
                <div class="flex flex-column gap-2">
                  <label>{{ $t('material_name') }}</label>
                  <InputText v-model="newRule.material_name" />
                </div>
              </div>
              <div class="col-12 md:col-4">
                <div class="flex flex-column gap-2">
                  <label>{{ $t('material_id') }}</label>
                  <InputText v-model="newRule.material_id" />
                </div>
              </div>
              <div class="col-12 md:col-2">
                <div class="flex flex-column gap-2">
                  <label>{{ $t('low_threshold') }}</label>
                  <InputNumber v-model="newRule.threshold_low" :minFractionDigits="1" />
                </div>
              </div>
              <div class="col-12 md:col-2">
                <div class="flex flex-column gap-2">
                  <label>{{ $t('critical_threshold') }}</label>
                  <InputNumber v-model="newRule.threshold_critical" :minFractionDigits="1" />
                </div>
              </div>
            </div>
            <Button
              :label="$t('save_rule')"
              icon="pi pi-save"
              class="mt-3"
              @click="saveRule"
              :loading="isSaving"
            />
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
import Tag from 'primevue/tag'
import InputText from 'primevue/inputtext'
import InputNumber from 'primevue/inputnumber'
import axios from 'axios'
import { ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { useToast } from 'primevue/usetoast'
import auth from '../services/auth'

const { t } = useI18n()
const toast = useToast()

interface AlertRule {
  id: string
  material_id: string
  material_name: string
  threshold_low: number
  threshold_critical: number
  is_active: boolean
}

interface Alert {
  id: string
  material_name: string
  current_qty: number
  threshold: number
  severity: string
  is_read: boolean
  created_at: string
}

interface AlertSummary {
  total_active: number
  unread_count: number
  critical_count: number
  low_count: number
}

const rules = ref<AlertRule[]>([])
const alerts = ref<Alert[]>([])
const summary = ref<AlertSummary | null>(null)
const isLoading = ref(false)
const isSaving = ref(false)

const newRule = ref<Partial<AlertRule>>({
  material_name: '',
  material_id: '',
  threshold_low: 10,
  threshold_critical: 5,
  is_active: true,
})

const formatTime = (date: string) => {
  if (!date) return '-'
  return new Date(date).toLocaleString()
}

const getSeverity = (severity: string) => {
  return severity === 'critical' ? 'danger' : 'warn'
}

const showAddRule = () => {
  newRule.value = {
    material_name: '',
    material_id: '',
    threshold_low: 10,
    threshold_critical: 5,
    is_active: true,
  }
}

const saveRule = async () => {
  if (!newRule.value.material_name) return
  isSaving.value = true
  try {
    await axios.post(
      `http://${import.meta.env.VITE_APP_BACKEND_HOST}/inventoryalerts/api/rules`,
      newRule.value,
      { headers: { Authorization: `Bearer ${auth.accessToken.value}` } },
    )
    toast.add({
      severity: 'success',
      summary: t('success'),
      detail: t('rule_saved'),
      group: 'br',
      life: 3000,
    })
    newRule.value = {
      material_name: '',
      material_id: '',
      threshold_low: 10,
      threshold_critical: 5,
      is_active: true,
    }
    loadData()
  } catch {
    toast.add({
      severity: 'error',
      summary: t('failed'),
      detail: t('save_failed'),
      group: 'br',
      life: 3000,
    })
  } finally {
    isSaving.value = false
  }
}

const deleteRule = async (id: string) => {
  try {
    await axios.delete(
      `http://${import.meta.env.VITE_APP_BACKEND_HOST}/inventoryalerts/api/rules/${id}`,
      { headers: { Authorization: `Bearer ${auth.accessToken.value}` } },
    )
    toast.add({
      severity: 'success',
      summary: t('success'),
      detail: t('rule_deleted'),
      group: 'br',
      life: 3000,
    })
    loadData()
  } catch {
    toast.add({
      severity: 'error',
      summary: t('failed'),
      detail: t('delete_failed'),
      group: 'br',
      life: 3000,
    })
  }
}

const markAsRead = async (id: string) => {
  try {
    await axios.put(
      `http://${import.meta.env.VITE_APP_BACKEND_HOST}/inventoryalerts/api/alerts/${id}/read`,
      {},
      { headers: { Authorization: `Bearer ${auth.accessToken.value}` } },
    )
    loadData()
  } catch {
    // ignore
  }
}

const loadData = async () => {
  isLoading.value = true
  try {
    const [rulesRes, alertsRes, summaryRes] = await Promise.all([
      axios.get(`http://${import.meta.env.VITE_APP_BACKEND_HOST}/inventoryalerts/api/rules`, {
        headers: { Authorization: `Bearer ${auth.accessToken.value}` },
      }),
      axios.get(`http://${import.meta.env.VITE_APP_BACKEND_HOST}/inventoryalerts/api/alerts`, {
        headers: { Authorization: `Bearer ${auth.accessToken.value}` },
      }),
      axios.get(`http://${import.meta.env.VITE_APP_BACKEND_HOST}/inventoryalerts/api/summary`, {
        headers: { Authorization: `Bearer ${auth.accessToken.value}` },
      }),
    ])
    rules.value = rulesRes.data.data || []
    alerts.value = alertsRes.data.data || []
    summary.value = summaryRes.data.data
  } catch {
    toast.add({
      severity: 'error',
      summary: t('failed'),
      detail: t('load_failed'),
      group: 'br',
      life: 3000,
    })
  } finally {
    isLoading.value = false
  }
}

loadData()
</script>

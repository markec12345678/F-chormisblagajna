<template>
  <div class="w-full">
    <div class="grid mx-2">
      <div class="col-12">
        <div class="flex justify-content-between align-items-center">
          <h3>{{ $t('audit_log') }}</h3>
          <div class="flex gap-2">
            <Dropdown
              v-model="actionFilter"
              :options="actionOptions"
              optionLabel="label"
              optionValue="value"
              :placeholder="$t('all_actions')"
              class="w-10rem"
            />
            <Dropdown
              v-model="resourceFilter"
              :options="resourceOptions"
              optionLabel="label"
              optionValue="value"
              :placeholder="$t('all_resources')"
              class="w-10rem"
            />
            <Button
              icon="pi pi-refresh"
              severity="secondary"
              @click="loadData"
              :loading="isLoading"
            />
          </div>
        </div>
      </div>

      <div class="col-12" v-if="summary">
        <div class="grid">
          <div class="col-12 md:col-3">
            <Card>
              <template #title>{{ $t('total_entries') }}</template>
              <template #content>
                <div class="text-4xl font-bold text-primary">{{ summary.total_entries }}</div>
              </template>
            </Card>
          </div>
          <div class="col-12 md:col-3">
            <Card>
              <template #title>{{ $t('by_action') }}</template>
              <template #content>
                <div class="flex flex-column gap-1">
                  <div
                    v-for="a in summary.by_action"
                    :key="a.action"
                    class="flex justify-content-between"
                  >
                    <Tag :value="a.action" :severity="getActionSeverity(a.action)" />
                    <span class="font-bold">{{ a.count }}</span>
                  </div>
                </div>
              </template>
            </Card>
          </div>
          <div class="col-12 md:col-3">
            <Card>
              <template #title>{{ $t('by_resource') }}</template>
              <template #content>
                <div class="flex flex-column gap-1">
                  <div
                    v-for="r in summary.by_resource"
                    :key="r.resource"
                    class="flex justify-content-between"
                  >
                    <span>{{ r.resource }}</span>
                    <span class="font-bold">{{ r.count }}</span>
                  </div>
                </div>
              </template>
            </Card>
          </div>
          <div class="col-12 md:col-3">
            <Card>
              <template #title>{{ $t('by_user') }}</template>
              <template #content>
                <div class="flex flex-column gap-1">
                  <div
                    v-for="u in summary.by_user"
                    :key="u.user_id"
                    class="flex justify-content-between"
                  >
                    <span>{{ u.username || u.user_id }}</span>
                    <span class="font-bold">{{ u.count }}</span>
                  </div>
                </div>
              </template>
            </Card>
          </div>
        </div>
      </div>

      <div class="col-12">
        <Card>
          <template #title>{{ $t('audit_entries') }}</template>
          <template #content>
            <DataTable
              :value="entries"
              stripedRows
              :loading="isLoading"
              :rows="20"
              :paginator="true"
            >
              <template #empty>{{ $t('no_audit_entries') }}</template>
              <Column field="created_at" :header="$t('timestamp')">
                <template #body="slotProps">
                  {{ formatDateTime(slotProps?.data?.created_at) }}
                </template>
              </Column>
              <Column field="action" :header="$t('action')">
                <template #body="slotProps">
                  <Tag
                    :value="slotProps?.data?.action"
                    :severity="getActionSeverity(slotProps?.data?.action)"
                  />
                </template>
              </Column>
              <Column field="resource" :header="$t('resource')"></Column>
              <Column field="resource_id" :header="$t('resource_id')">
                <template #body="slotProps">
                  <span class="text-sm text-500">{{ slotProps?.data?.resource_id || '-' }}</span>
                </template>
              </Column>
              <Column field="username" :header="$t('user')"></Column>
              <Column field="ip_address" :header="$t('ip_address')"></Column>
              <Column field="details" :header="$t('details')">
                <template #body="slotProps">
                  <span
                    v-if="
                      slotProps?.data?.details && Object.keys(slotProps.data.details).length > 0
                    "
                    class="text-sm"
                  >
                    {{ formatDetails(slotProps?.data?.details) }}
                  </span>
                  <span v-else class="text-400">-</span>
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
import Tag from 'primevue/tag'
import Dropdown from 'primevue/dropdown'
import axios from 'axios'
import { ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { useToast } from 'primevue/usetoast'
import auth from '../services/auth'

const { t } = useI18n()
const toast = useToast()

interface AuditLogEntry {
  id: string
  action: string
  resource: string
  resource_id: string
  user_id: string
  username: string
  details: Record<string, string>
  ip_address: string
  created_at: string
}

interface ActionSummary {
  action: string
  count: number
}

interface ResourceSummary {
  resource: string
  count: number
}

interface UserSummary {
  user_id: string
  username: string
  count: number
}

interface AuditSummary {
  total_entries: number
  by_action: ActionSummary[]
  by_resource: ResourceSummary[]
  by_user: UserSummary[]
}

const entries = ref<AuditLogEntry[]>([])
const summary = ref<AuditSummary | null>(null)
const isLoading = ref(false)
const actionFilter = ref('')
const resourceFilter = ref('')

const actionOptions = [
  { label: 'All', value: '' },
  { label: 'Create', value: 'create' },
  { label: 'Update', value: 'update' },
  { label: 'Delete', value: 'delete' },
  { label: 'Login', value: 'login' },
  { label: 'Logout', value: 'logout' },
  { label: 'Fiscalize', value: 'fiscalize' },
]

const resourceOptions = [
  { label: 'All', value: '' },
  { label: 'Order', value: 'order' },
  { label: 'Product', value: 'product' },
  { label: 'Material', value: 'material' },
  { label: 'Customer', value: 'customer' },
  { label: 'Settings', value: 'settings' },
  { label: 'User', value: 'user' },
  { label: 'Category', value: 'category' },
]

const formatDateTime = (date: string) => {
  if (!date) return '-'
  return new Date(date).toLocaleString()
}

const getActionSeverity = (action: string) => {
  switch (action) {
    case 'create':
      return 'success'
    case 'update':
      return 'info'
    case 'delete':
      return 'danger'
    case 'login':
      return 'success'
    case 'logout':
      return 'warn'
    case 'fiscalize':
      return 'secondary'
    default:
      return 'secondary'
  }
}

const formatDetails = (details: Record<string, string>) => {
  if (!details) return ''
  return Object.entries(details)
    .map(([k, v]) => `${k}: ${v}`)
    .join(', ')
}

const loadData = async () => {
  isLoading.value = true
  try {
    const params: Record<string, string> = {}
    if (actionFilter.value) params.action = actionFilter.value
    if (resourceFilter.value) params.resource = resourceFilter.value

    const [logsRes, summaryRes] = await Promise.all([
      axios.get(`http://${import.meta.env.VITE_APP_BACKEND_HOST}/auditlog/api/logs`, {
        headers: { Authorization: `Bearer ${auth.accessToken.value}` },
        params,
      }),
      axios.get(`http://${import.meta.env.VITE_APP_BACKEND_HOST}/auditlog/api/summary`, {
        headers: { Authorization: `Bearer ${auth.accessToken.value}` },
      }),
    ])
    entries.value = logsRes.data.data || []
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

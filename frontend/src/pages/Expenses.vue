<template>
  <div class="w-full">
    <div class="grid mx-2">
      <div class="col-12">
        <div class="flex justify-content-between align-items-center">
          <h3>{{ $t('expense_tracking') }}</h3>
          <div class="flex gap-2">
            <Calendar v-model="startDate" :placeholder="$t('start_date')" dateFormat="yy-mm-dd" />
            <Calendar v-model="endDate" :placeholder="$t('end_date')" dateFormat="yy-mm-dd" />
            <Button :label="$t('add_expense')" icon="pi pi-plus" @click="showAddDialog" />
          </div>
        </div>
      </div>

      <div class="col-12" v-if="summary">
        <div class="grid">
          <div class="col-12 md:col-4">
            <Card>
              <template #title>{{ $t('total_expenses') }}</template>
              <template #content>
                <div class="text-4xl font-bold text-red-500">{{ formatCurrency(summary.total_expenses) }}</div>
              </template>
            </Card>
          </div>
          <div class="col-12 md:col-4">
            <Card>
              <template #title>{{ $t('monthly_total') }}</template>
              <template #content>
                <div class="text-4xl font-bold text-orange-500">{{ formatCurrency(summary.monthly_total) }}</div>
              </template>
            </Card>
          </div>
          <div class="col-12 md:col-4">
            <Card>
              <template #title>{{ $t('categories') }}</template>
              <template #content>
                <div class="text-4xl font-bold text-blue-500">{{ summary.by_category?.length || 0 }}</div>
              </template>
            </Card>
          </div>
        </div>
      </div>

      <div class="col-12" v-if="summary">
        <div class="grid">
          <div class="col-12 md:col-6">
            <Card>
              <template #title>{{ $t('by_category') }}</template>
              <template #content>
                <DataTable :value="summary.by_category" stripedRows>
                  <template #empty>{{ $t('no_data') }}</template>
                  <Column field="category" :header="$t('category')"></Column>
                  <Column field="total" :header="$t('total')">
                    <template #body="slotProps">
                      {{ formatCurrency(slotProps?.data?.total || 0) }}
                    </template>
                  </Column>
                  <Column field="count" :header="$t('count')"></Column>
                </DataTable>
              </template>
            </Card>
          </div>
          <div class="col-12 md:col-6">
            <Card>
              <template #title>{{ $t('recent_expenses') }}</template>
              <template #content>
                <DataTable :value="expenses.slice(0, 10)" stripedRows>
                  <template #empty>{{ $t('no_expenses') }}</template>
                  <Column field="date" :header="$t('date')">
                    <template #body="slotProps">
                      {{ formatDate(slotProps?.data?.date) }}
                    </template>
                  </Column>
                  <Column field="description" :header="$t('description')"></Column>
                  <Column field="amount" :header="$t('amount')">
                    <template #body="slotProps">
                      {{ formatCurrency(slotProps?.data?.amount || 0) }}
                    </template>
                  </Column>
                  <Column field="category" :header="$t('category')"></Column>
                </DataTable>
              </template>
            </Card>
          </div>
        </div>
      </div>

      <Dialog v-model:visible="expenseDialog" modal :header="editingExpense ? $t('edit_expense') : $t('add_expense')" :style="{ width: '50rem' }">
        <div class="grid">
          <div class="col-12 md:col-6">
            <div class="flex flex-column gap-2">
              <label>{{ $t('description') }}</label>
              <InputText v-model="currentExpense.description" />
            </div>
          </div>
          <div class="col-12 md:col-6">
            <div class="flex flex-column gap-2">
              <label>{{ $t('amount') }}</label>
              <InputNumber v-model="currentExpense.amount" mode="currency" currency="EUR" locale="sl-SI" />
            </div>
          </div>
          <div class="col-12 md:col-6">
            <div class="flex flex-column gap-2">
              <label>{{ $t('category') }}</label>
              <Dropdown v-model="currentExpense.category" :options="categories" :placeholder="$t('select_category')" />
            </div>
          </div>
          <div class="col-12 md:col-6">
            <div class="flex flex-column gap-2">
              <label>{{ $t('date') }}</label>
              <Calendar v-model="currentExpense.date" dateFormat="yy-mm-dd" />
            </div>
          </div>
          <div class="col-12">
            <div class="flex flex-column gap-2">
              <label>{{ $t('notes') }}</label>
              <Textarea v-model="currentExpense.notes" rows="3" />
            </div>
          </div>
        </div>
        <template #footer>
          <Button :label="$t('cancel')" severity="secondary" @click="expenseDialog = false" />
          <Button :label="$t('save')" @click="saveExpense" :loading="isSaving" />
        </template>
      </Dialog>
    </div>
  </div>
</template>

<script setup lang="ts">
import Card from 'primevue/card'
import Button from 'primevue/button'
import Calendar from 'primevue/calendar'
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import Dialog from 'primevue/dialog'
import InputText from 'primevue/inputtext'
import InputNumber from 'primevue/inputnumber'
import Textarea from 'primevue/textarea'
import Dropdown from 'primevue/dropdown'
import axios from 'axios'
import { ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { useToast } from 'primevue/usetoast'
import auth from '../services/auth'

const { t } = useI18n()
const toast = useToast()

interface Expense {
  id: string
  description: string
  amount: number
  category: string
  date: string
  notes: string
}

interface CategorySummary {
  category: string
  total: number
  count: number
}

interface ExpenseSummary {
  total_expenses: number
  monthly_total: number
  by_category: CategorySummary[]
}

const expenses = ref<Expense[]>([])
const summary = ref<ExpenseSummary | null>(null)
const startDate = ref<Date | null>(null)
const endDate = ref<Date | null>(null)
const isLoading = ref(false)
const isSaving = ref(false)
const expenseDialog = ref(false)
const editingExpense = ref(false)
const currentExpense = ref<Partial<Expense>>({})

const categories = ['rent', 'utilities', 'supplies', 'marketing', 'maintenance', 'insurance', 'taxes', 'other']

const formatCurrency = (amount: number) => {
  return new Intl.NumberFormat('sl-SI', { style: 'currency', currency: 'EUR' }).format(amount)
}

const formatDate = (date: string) => {
  if (!date) return '-'
  return new Date(date).toLocaleDateString()
}

const showAddDialog = () => {
  editingExpense.value = false
  currentExpense.value = { date: new Date().toISOString().split('T')[0] }
  expenseDialog.value = true
}

const saveExpense = async () => {
  isSaving.value = true
  try {
    if (editingExpense.value && currentExpense.value.id) {
      await axios.put(
        `http://${import.meta.env.VITE_APP_BACKEND_HOST}/expense/api/expenses/${currentExpense.value.id}`,
        currentExpense.value,
        { headers: { Authorization: `Bearer ${auth.accessToken.value}` } }
      )
    } else {
      await axios.post(
        `http://${import.meta.env.VITE_APP_BACKEND_HOST}/expense/api/expenses`,
        currentExpense.value,
        { headers: { Authorization: `Bearer ${auth.accessToken.value}` } }
      )
    }
    toast.add({ severity: 'success', summary: t('success'), detail: t('expense_saved'), group: 'br', life: 3000 })
    expenseDialog.value = false
    loadData()
  } catch {
    toast.add({ severity: 'error', summary: t('failed'), detail: t('save_failed'), group: 'br', life: 3000 })
  } finally {
    isSaving.value = false
  }
}

const loadData = async () => {
  isLoading.value = true
  try {
    const [expensesRes, summaryRes] = await Promise.all([
      axios.get(`http://${import.meta.env.VITE_APP_BACKEND_HOST}/expense/api/expenses`, {
        headers: { Authorization: `Bearer ${auth.accessToken.value}` }
      }),
      axios.get(`http://${import.meta.env.VITE_APP_BACKEND_HOST}/expense/api/expenses/summary`, {
        headers: { Authorization: `Bearer ${auth.accessToken.value}` }
      })
    ])
    expenses.value = expensesRes.data.data || []
    summary.value = summaryRes.data.data
  } catch {
    toast.add({ severity: 'error', summary: t('failed'), detail: t('load_failed'), group: 'br', life: 3000 })
  } finally {
    isLoading.value = false
  }
}

loadData()
</script>

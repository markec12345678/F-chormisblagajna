<template>
  <div class="w-full">
    <div class="grid mx-2">
      <div class="col-12">
        <h3>{{ $t('accounting_export') }}</h3>
      </div>
      <div class="col-12 md:col-6">
        <Card>
          <template #title>QuickBooks</template>
          <template #content>
            <div class="flex flex-column gap-3">
              <p class="m-0">{{ $t('quickbooks_description') }}</p>
              <div class="flex gap-2">
                <Calendar v-model="quickbooksStart" :placeholder="$t('start_date')" dateFormat="yy-mm-dd" />
                <Calendar v-model="quickbooksEnd" :placeholder="$t('end_date')" dateFormat="yy-mm-dd" />
              </div>
              <Button
                :label="$t('export_quickbooks')"
                icon="pi pi-download"
                @click="exportQuickBooks"
                :loading="isExportingQuickBooks"
              />
            </div>
          </template>
        </Card>
      </div>
      <div class="col-12 md:col-6">
        <Card>
          <template #title>Xero</template>
          <template #content>
            <div class="flex flex-column gap-3">
              <p class="m-0">{{ $t('xero_description') }}</p>
              <div class="flex gap-2">
                <Calendar v-model="xeroStart" :placeholder="$t('start_date')" dateFormat="yy-mm-dd" />
                <Calendar v-model="xeroEnd" :placeholder="$t('end_date')" dateFormat="yy-mm-dd" />
              </div>
              <Button
                :label="$t('export_xero')"
                icon="pi pi-download"
                @click="exportXero"
                :loading="isExportingXero"
              />
            </div>
          </template>
        </Card>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import Card from 'primevue/card'
import Button from 'primevue/button'
import Calendar from 'primevue/calendar'
import { ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { useToast } from 'primevue/usetoast'
import auth from '../services/auth'

const { t } = useI18n()
const toast = useToast()

const quickbooksStart = ref<Date | null>(null)
const quickbooksEnd = ref<Date | null>(null)
const xeroStart = ref<Date | null>(null)
const xeroEnd = ref<Date | null>(null)

const isExportingQuickBooks = ref(false)
const isExportingXero = ref(false)

const formatDate = (date: Date | null): string => {
  if (!date) return ''
  return date.toISOString().split('T')[0]
}

const exportQuickBooks = () => {
  isExportingQuickBooks.value = true
  const params = new URLSearchParams()
  if (quickbooksStart.value) params.append('start_date', formatDate(quickbooksStart.value))
  if (quickbooksEnd.value) params.append('end_date', formatDate(quickbooksEnd.value))

  window.open(
    `http://${import.meta.env.VITE_APP_BACKEND_HOST}/accounting/api/export/quickbooks?${params.toString()}`,
    '_blank'
  )

  toast.add({
    severity: 'success',
    summary: t('success'),
    detail: t('export_started'),
    group: 'br',
    life: 3000,
  })
  isExportingQuickBooks.value = false
}

const exportXero = () => {
  isExportingXero.value = true
  const params = new URLSearchParams()
  if (xeroStart.value) params.append('start_date', formatDate(xeroStart.value))
  if (xeroEnd.value) params.append('end_date', formatDate(xeroEnd.value))

  window.open(
    `http://${import.meta.env.VITE_APP_BACKEND_HOST}/accounting/api/export/xero?${params.toString()}`,
    '_blank'
  )

  toast.add({
    severity: 'success',
    summary: t('success'),
    detail: t('export_started'),
    group: 'br',
    life: 3000,
  })
  isExportingXero.value = false
}
</script>

<template>
  <div class="w-full">
    <div class="grid mx-2">
      <div class="col-12">
        <div class="flex justify-content-between align-items-center">
          <h3>{{ $t('receipt_customization') }}</h3>
          <Button :label="$t('add_template')" icon="pi pi-plus" @click="showAddDialog" />
        </div>
      </div>

      <div class="col-12 md:col-6">
        <Card>
          <template #title>{{ $t('receipt_templates') }}</template>
          <template #content>
            <DataTable
              :value="templates"
              stripedRows
              :loading="isLoading"
              selectionMode="single"
              v-model:selection="selectedTemplate"
              @rowSelect="onTemplateSelect"
            >
              <template #empty>{{ $t('no_templates') }}</template>
              <Column field="name" :header="$t('name')">
                <template #body="slotProps">
                  <div class="flex align-items-center gap-2">
                    <span>{{ slotProps?.data?.name }}</span>
                    <Tag
                      v-if="slotProps?.data?.is_default"
                      :value="$t('default')"
                      severity="success"
                    />
                  </div>
                </template>
              </Column>
              <Column field="paper_width" :header="$t('paper_width')">
                <template #body="slotProps"> {{ slotProps?.data?.paper_width || 80 }}mm </template>
              </Column>
              <Column :header="$t('actions')">
                <template #body="slotProps">
                  <div class="flex gap-1">
                    <Button
                      icon="pi pi-pencil"
                      severity="info"
                      text
                      rounded
                      size="small"
                      @click="editTemplate(slotProps?.data)"
                    />
                    <Button
                      icon="pi pi-trash"
                      severity="danger"
                      text
                      rounded
                      size="small"
                      @click="deleteTemplate(slotProps?.data?.id)"
                    />
                  </div>
                </template>
              </Column>
            </DataTable>
          </template>
        </Card>
      </div>

      <div class="col-12 md:col-6">
        <Card>
          <template #title>{{ $t('print_settings') }}</template>
          <template #content>
            <div class="flex flex-column gap-3">
              <div class="flex flex-column gap-2">
                <label>{{ $t('printer_name') }}</label>
                <InputText v-model="printSettings.printer_name" />
              </div>
              <div class="flex flex-column gap-2">
                <label>{{ $t('printer_ip') }}</label>
                <InputText v-model="printSettings.printer_ip" placeholder="192.168.1.100" />
              </div>
              <div class="flex flex-column gap-2">
                <label>{{ $t('print_copies') }}</label>
                <InputNumber v-model="printSettings.print_copies" :min="1" :max="5" />
              </div>
              <div class="flex align-items-center gap-2">
                <Checkbox v-model="printSettings.auto_print" :binary="true" />
                <label>{{ $t('auto_print') }}</label>
              </div>
              <Divider />
              <Button
                :label="$t('save_settings')"
                icon="pi pi-save"
                @click="savePrintSettings"
                :loading="isSaving"
              />
            </div>
          </template>
        </Card>
      </div>

      <div class="col-12" v-if="selectedTemplate">
        <Card>
          <template #title>{{ $t('template_editor') }} — {{ selectedTemplate.name }}</template>
          <template #content>
            <div class="grid">
              <div class="col-12 md:col-6">
                <div class="flex flex-column gap-3">
                  <div class="flex flex-column gap-2">
                    <label>{{ $t('business_name') }}</label>
                    <InputText v-model="selectedTemplate.business_name" />
                  </div>
                  <div class="flex flex-column gap-2">
                    <label>{{ $t('business_address') }}</label>
                    <InputText v-model="selectedTemplate.business_address" />
                  </div>
                  <div class="flex flex-column gap-2">
                    <label>{{ $t('business_phone') }}</label>
                    <InputText v-model="selectedTemplate.business_phone" />
                  </div>
                  <div class="flex flex-column gap-2">
                    <label>{{ $t('business_tax_id') }}</label>
                    <InputText v-model="selectedTemplate.business_tax_id" />
                  </div>
                  <div class="flex flex-column gap-2">
                    <label>{{ $t('paper_width') }}</label>
                    <Dropdown
                      v-model="selectedTemplate.paper_width"
                      :options="paperWidths"
                      optionLabel="label"
                      optionValue="value"
                    />
                  </div>
                </div>
              </div>
              <div class="col-12 md:col-6">
                <div class="flex flex-column gap-3">
                  <div class="flex flex-column gap-2">
                    <label>{{ $t('header') }}</label>
                    <Textarea v-model="selectedTemplate.header" rows="2" />
                  </div>
                  <div class="flex flex-column gap-2">
                    <label>{{ $t('footer') }}</label>
                    <Textarea v-model="selectedTemplate.footer" rows="3" />
                  </div>
                  <div class="flex flex-wrap gap-4">
                    <div class="flex align-items-center gap-2">
                      <Checkbox v-model="selectedTemplate.show_logo" :binary="true" />
                      <label>{{ $t('show_logo') }}</label>
                    </div>
                    <div class="flex align-items-center gap-2">
                      <Checkbox v-model="selectedTemplate.show_tax_id" :binary="true" />
                      <label>{{ $t('show_tax_id') }}</label>
                    </div>
                    <div class="flex align-items-center gap-2">
                      <Checkbox v-model="selectedTemplate.show_qr_code" :binary="true" />
                      <label>{{ $t('show_qr_code') }}</label>
                    </div>
                    <div class="flex align-items-center gap-2">
                      <Checkbox v-model="selectedTemplate.show_server" :binary="true" />
                      <label>{{ $t('show_server') }}</label>
                    </div>
                    <div class="flex align-items-center gap-2">
                      <Checkbox v-model="selectedTemplate.show_table" :binary="true" />
                      <label>{{ $t('show_table') }}</label>
                    </div>
                  </div>
                </div>
              </div>
            </div>
            <Divider />
            <div class="flex gap-2">
              <Button
                :label="$t('save_template')"
                icon="pi pi-save"
                @click="saveTemplate"
                :loading="isSaving"
              />
              <Button
                :label="$t('preview')"
                icon="pi pi-eye"
                severity="secondary"
                @click="previewTemplate"
              />
            </div>
          </template>
        </Card>
      </div>
    </div>

    <Dialog
      v-model:visible="addDialog"
      modal
      :header="$t('add_template')"
      :style="{ width: '30rem' }"
    >
      <div class="flex flex-column gap-2">
        <label>{{ $t('name') }}</label>
        <InputText v-model="newTemplateName" />
      </div>
      <template #footer>
        <Button :label="$t('cancel')" severity="secondary" @click="addDialog = false" />
        <Button :label="$t('create')" @click="createTemplate" :loading="isSaving" />
      </template>
    </Dialog>

    <Dialog
      v-model:visible="previewDialog"
      modal
      :header="$t('receipt_preview')"
      :style="{ width: '350px' }"
    >
      <div
        class="receipt-preview p-3"
        style="
          font-family: monospace;
          font-size: 12px;
          white-space: pre-wrap;
          background: white;
          color: black;
          border: 1px solid #ddd;
          border-radius: 4px;
        "
      >
        {{ previewContent }}
      </div>
      <template #footer>
        <Button :label="$t('close')" @click="previewDialog = false" />
      </template>
    </Dialog>
  </div>
</template>

<script setup lang="ts">
import Card from 'primevue/card'
import Button from 'primevue/button'
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import InputText from 'primevue/inputtext'
import InputNumber from 'primevue/inputnumber'
import Textarea from 'primevue/textarea'
import Dropdown from 'primevue/dropdown'
import Dialog from 'primevue/dialog'
import Divider from 'primevue/divider'
import Checkbox from 'primevue/checkbox'
import Tag from 'primevue/tag'
import axios from 'axios'
import { ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { useToast } from 'primevue/usetoast'
import auth from '../services/auth'

const { t } = useI18n()
const toast = useToast()

interface ReceiptTemplate {
  id: string
  name: string
  business_name: string
  business_address: string
  business_phone: string
  business_tax_id: string
  header: string
  footer: string
  show_logo: boolean
  show_tax_id: boolean
  show_qr_code: boolean
  show_server: boolean
  show_table: boolean
  paper_width: number
  custom_fields: { key: string; value: string }[]
  is_default: boolean
}

interface PrintSettings {
  id: string
  printer_name: string
  printer_ip: string
  auto_print: boolean
  print_copies: number
  template_id: string
  connected: boolean
}

const templates = ref<ReceiptTemplate[]>([])
const selectedTemplate = ref<ReceiptTemplate | null>(null)
const printSettings = ref<PrintSettings>({
  id: '',
  printer_name: '',
  printer_ip: '',
  auto_print: true,
  print_copies: 1,
  template_id: '',
  connected: false,
})
const isLoading = ref(false)
const isSaving = ref(false)
const addDialog = ref(false)
const previewDialog = ref(false)
const previewContent = ref('')
const newTemplateName = ref('')

const paperWidths = [
  { label: '58mm', value: 58 },
  { label: '80mm', value: 80 },
]

const showAddDialog = () => {
  newTemplateName.value = ''
  addDialog.value = true
}

const onTemplateSelect = (event: { data: ReceiptTemplate }) => {
  selectedTemplate.value = { ...event.data }
}

const editTemplate = (tpl: ReceiptTemplate) => {
  selectedTemplate.value = { ...tpl }
}

const createTemplate = async () => {
  if (!newTemplateName.value.trim()) return
  isSaving.value = true
  try {
    await axios.post(
      `http://${import.meta.env.VITE_APP_BACKEND_HOST}/receipt/api/templates`,
      { name: newTemplateName.value, paper_width: 80, show_logo: true, show_tax_id: true },
      { headers: { Authorization: `Bearer ${auth.accessToken.value}` } },
    )
    toast.add({
      severity: 'success',
      summary: t('success'),
      detail: t('template_created'),
      group: 'br',
      life: 3000,
    })
    addDialog.value = false
    loadTemplates()
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

const saveTemplate = async () => {
  if (!selectedTemplate.value) return
  isSaving.value = true
  try {
    await axios.post(
      `http://${import.meta.env.VITE_APP_BACKEND_HOST}/receipt/api/templates`,
      selectedTemplate.value,
      { headers: { Authorization: `Bearer ${auth.accessToken.value}` } },
    )
    toast.add({
      severity: 'success',
      summary: t('success'),
      detail: t('template_saved'),
      group: 'br',
      life: 3000,
    })
    loadTemplates()
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

const deleteTemplate = async (id: string) => {
  try {
    await axios.delete(
      `http://${import.meta.env.VITE_APP_BACKEND_HOST}/receipt/api/templates/${id}`,
      { headers: { Authorization: `Bearer ${auth.accessToken.value}` } },
    )
    toast.add({
      severity: 'success',
      summary: t('success'),
      detail: t('template_deleted'),
      group: 'br',
      life: 3000,
    })
    if (selectedTemplate.value?.id === id) selectedTemplate.value = null
    loadTemplates()
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

const previewTemplate = () => {
  if (!selectedTemplate.value) return
  const tpl = selectedTemplate.value
  let preview = ''
  if (tpl.header) preview += tpl.header + '\n'
  preview += '========================\n'
  preview += (tpl.business_name || 'Business Name') + '\n'
  if (tpl.show_tax_id && tpl.business_tax_id) preview += 'Tax ID: ' + tpl.business_tax_id + '\n'
  if (tpl.business_address) preview += tpl.business_address + '\n'
  if (tpl.business_phone) preview += tpl.business_phone + '\n'
  preview += '========================\n\n'
  preview += '2x Margherita        17.00\n'
  preview += '1x Coca Cola           2.50\n'
  preview += '========================\n'
  preview += 'TOTAL:                19.50\n\n'
  if (tpl.footer) preview += tpl.footer + '\n'
  previewContent.value = preview
  previewDialog.value = true
}

const savePrintSettings = async () => {
  isSaving.value = true
  try {
    await axios.post(
      `http://${import.meta.env.VITE_APP_BACKEND_HOST}/receipt/api/print-settings`,
      printSettings.value,
      { headers: { Authorization: `Bearer ${auth.accessToken.value}` } },
    )
    toast.add({
      severity: 'success',
      summary: t('success'),
      detail: t('settings_saved'),
      group: 'br',
      life: 3000,
    })
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

const loadTemplates = async () => {
  isLoading.value = true
  try {
    const response = await axios.get(
      `http://${import.meta.env.VITE_APP_BACKEND_HOST}/receipt/api/templates`,
      { headers: { Authorization: `Bearer ${auth.accessToken.value}` } },
    )
    templates.value = response.data.data || []
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

const loadPrintSettings = async () => {
  try {
    const response = await axios.get(
      `http://${import.meta.env.VITE_APP_BACKEND_HOST}/receipt/api/print-settings`,
      { headers: { Authorization: `Bearer ${auth.accessToken.value}` } },
    )
    printSettings.value = response.data.data
  } catch {
    // ignore
  }
}

const init = async () => {
  await Promise.all([loadTemplates(), loadPrintSettings()])
}

init()
</script>

<template>
  <div class="w-full">
    <div class="grid mx-2">
      <div class="col-12">
        <div class="flex justify-content-between align-items-center">
          <h3>{{ $t('customer_display') }}</h3>
          <Button icon="pi pi-plus" :label="$t('add_display')" @click="openNewDialog" />
        </div>
      </div>

      <div class="col-12">
        <Card>
          <template #title>{{ $t('display_configs') }}</template>
          <template #content>
            <DataTable :value="configs" stripedRows :loading="isLoading">
              <template #empty>{{ $t('no_displays') }}</template>
              <Column field="display_name" :header="$t('name')"></Column>
              <Column field="active" :header="$t('active')">
                <template #body="slotProps">
                  <Tag :severity="slotProps?.data?.active ? 'success' : 'secondary'">
                    {{ slotProps?.data?.active ? $t('active') : $t('inactive') }}
                  </Tag>
                </template>
              </Column>
              <Column field="auto_slide_interval" :header="$t('slide_interval')">
                <template #body="slotProps">
                  {{ slotProps?.data?.auto_slide_interval || 0 }}s
                </template>
              </Column>
              <Column field="theme" :header="$t('theme')"></Column>
              <Column :header="$t('actions')">
                <template #body="slotProps">
                  <Button
                    icon="pi pi-eye"
                    severity="info"
                    class="mr-1"
                    @click="previewDisplay(slotProps.data)"
                  />
                  <Button
                    icon="pi pi-pencil"
                    class="mr-1"
                    @click="openEditDialog(slotProps.data)"
                  />
                  <Button
                    icon="pi pi-trash"
                    severity="danger"
                    @click="confirmDelete(slotProps.data)"
                  />
                </template>
              </Column>
            </DataTable>
          </template>
        </Card>
      </div>
    </div>

    <Dialog
      v-model:visible="dialogVisible"
      :header="editing ? $t('edit_display') : $t('add_display')"
      :style="{ width: '500px' }"
    >
      <div class="grid">
        <div class="col-12">
          <label>{{ $t('display_name') }}</label>
          <InputText v-model="form.display_name" class="w-full" />
        </div>
        <div class="col-6">
          <label>{{ $t('slide_interval') }} (s)</label>
          <InputNumber v-model="form.auto_slide_interval" class="w-full" :min="3" :max="120" />
        </div>
        <div class="col-6">
          <label>{{ $t('theme') }}</label>
          <Dropdown v-model="form.theme" :options="themeOptions" class="w-full" />
        </div>
        <div class="col-12">
          <label>{{ $t('welcome_message') }}</label>
          <InputText v-model="form.welcome_message" class="w-full" />
        </div>
        <div class="col-6">
          <div class="flex align-items-center">
            <Checkbox v-model="form.show_promotions" :binary="true" inputId="showPromotions" />
            <label for="showPromotions" class="ml-2">{{ $t('show_promotions') }}</label>
          </div>
        </div>
        <div class="col-6">
          <div class="flex align-items-center">
            <Checkbox v-model="form.show_menu" :binary="true" inputId="showMenu" />
            <label for="showMenu" class="ml-2">{{ $t('show_menu') }}</label>
          </div>
        </div>
        <div class="col-6">
          <div class="flex align-items-center">
            <Checkbox v-model="form.show_order_status" :binary="true" inputId="showOrderStatus" />
            <label for="showOrderStatus" class="ml-2">{{ $t('show_order_status') }}</label>
          </div>
        </div>
        <div class="col-6">
          <div class="flex align-items-center">
            <Checkbox v-model="form.show_wait_time" :binary="true" inputId="showWaitTime" />
            <label for="showWaitTime" class="ml-2">{{ $t('show_wait_time') }}</label>
          </div>
        </div>
        <div class="col-6">
          <div class="flex align-items-center">
            <Checkbox v-model="form.active" :binary="true" inputId="activeDisplay" />
            <label for="activeDisplay" class="ml-2">{{ $t('active') }}</label>
          </div>
        </div>
      </div>
      <template #footer>
        <Button :label="$t('cancel')" severity="secondary" @click="dialogVisible = false" />
        <Button :label="$t('save')" @click="saveConfig" :loading="saving" />
      </template>
    </Dialog>

    <Dialog
      v-model:visible="deleteDialogVisible"
      :header="$t('confirm_delete')"
      :style="{ width: '400px' }"
    >
      <p>{{ $t('delete_display_confirm') }}</p>
      <template #footer>
        <Button :label="$t('cancel')" severity="secondary" @click="deleteDialogVisible = false" />
        <Button :label="$t('delete')" severity="danger" @click="deleteConfig" :loading="deleting" />
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
import Checkbox from 'primevue/checkbox'
import Tag from 'primevue/tag'
import axios from 'axios'
import { ref, reactive } from 'vue'
import { useI18n } from 'vue-i18n'
import { useToast } from 'primevue/usetoast'
import auth from '../services/auth'

const { t } = useI18n()
const toast = useToast()

interface DisplayConfig {
  id?: string
  display_name: string
  show_promotions: boolean
  show_menu: boolean
  show_order_status: boolean
  show_wait_time: boolean
  auto_slide_interval: number
  promotion_ids: string[]
  menu_category_ids: string[]
  theme: string
  welcome_message: string
  active: boolean
}

const configs = ref<DisplayConfig[]>([])
const isLoading = ref(false)
const dialogVisible = ref(false)
const deleteDialogVisible = ref(false)
const saving = ref(false)
const deleting = ref(false)
const editing = ref(false)
const selectedConfig = ref<DisplayConfig | null>(null)

const themeOptions = ['light', 'dark', 'colorful']

const form = reactive<DisplayConfig>({
  display_name: '',
  show_promotions: false,
  show_menu: false,
  show_order_status: true,
  show_wait_time: true,
  auto_slide_interval: 10,
  promotion_ids: [],
  menu_category_ids: [],
  theme: 'light',
  welcome_message: '',
  active: true,
})

const openNewDialog = () => {
  editing.value = false
  selectedConfig.value = null
  Object.assign(form, {
    display_name: '',
    show_promotions: false,
    show_menu: false,
    show_order_status: true,
    show_wait_time: true,
    auto_slide_interval: 10,
    promotion_ids: [],
    menu_category_ids: [],
    theme: 'light',
    welcome_message: '',
    active: true,
  })
  dialogVisible.value = true
}

const openEditDialog = (config: DisplayConfig) => {
  editing.value = true
  selectedConfig.value = config
  Object.assign(form, config)
  dialogVisible.value = true
}

const saveConfig = async () => {
  saving.value = true
  try {
    if (editing.value && selectedConfig.value?.id) {
      form.id = selectedConfig.value.id
    }
    await axios.post(
      `http://${import.meta.env.VITE_APP_BACKEND_HOST}/customerdisplay/api/configs`,
      form,
      { headers: { Authorization: `Bearer ${auth.accessToken.value}` } },
    )
    toast.add({ severity: 'success', summary: t('saved'), group: 'br', life: 2000 })
    dialogVisible.value = false
    await loadConfigs()
  } catch {
    toast.add({
      severity: 'error',
      summary: t('failed'),
      detail: t('save_failed'),
      group: 'br',
      life: 3000,
    })
  } finally {
    saving.value = false
  }
}

const confirmDelete = (config: DisplayConfig) => {
  selectedConfig.value = config
  deleteDialogVisible.value = true
}

const deleteConfig = async () => {
  if (!selectedConfig.value?.id) return
  deleting.value = true
  try {
    await axios.delete(
      `http://${import.meta.env.VITE_APP_BACKEND_HOST}/customerdisplay/api/configs/${selectedConfig.value.id}`,
      { headers: { Authorization: `Bearer ${auth.accessToken.value}` } },
    )
    toast.add({ severity: 'success', summary: t('deleted'), group: 'br', life: 2000 })
    deleteDialogVisible.value = false
    await loadConfigs()
  } catch {
    toast.add({
      severity: 'error',
      summary: t('failed'),
      detail: t('delete_failed'),
      group: 'br',
      life: 3000,
    })
  } finally {
    deleting.value = false
  }
}

const previewDisplay = (config: DisplayConfig) => {
  window.open(`/display/${config.id}`, '_blank')
}

const loadConfigs = async () => {
  isLoading.value = true
  try {
    const response = await axios.get(
      `http://${import.meta.env.VITE_APP_BACKEND_HOST}/customerdisplay/api/configs`,
      { headers: { Authorization: `Bearer ${auth.accessToken.value}` } },
    )
    configs.value = response.data.data || []
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

loadConfigs()
</script>

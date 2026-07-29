<template>
  <div class="w-full">
    <div class="grid mx-2">
      <div class="col-12">
        <h3>{{ $t('self_service_kiosk') }}</h3>
      </div>
      <div class="col-12">
        <DataTable :value="configs" stripedRows :loading="loading">
          <Column field="name" :header="$t('name')"></Column
          ><Column field="location" :header="$t('location')"></Column>
          <Column field="theme" :header="$t('theme')"></Column>
          <Column field="active" :header="$t('active')"
            ><template #body="s"
              ><i
                :class="s.data.active ? 'pi pi-check text-green-500' : 'pi pi-times text-300'"
              ></i></template
          ></Column>
          <Column :header="$t('actions')"
            ><template #body="s"
              ><Button icon="pi pi-pencil" size="small" @click="edit(s.data)" /></template
          ></Column>
        </DataTable>
      </div>
      <div class="col-12 mt-3">
        <Card
          ><template #title>{{ editing ? $t('edit') : $t('add_kiosk') }}</template
          ><template #content>
            <div class="grid">
              <div class="col-6">
                <label>{{ $t('name') }}</label
                ><InputText v-model="form.name" class="w-full" />
              </div>
              <div class="col-6">
                <label>{{ $t('location') }}</label
                ><InputText v-model="form.location" class="w-full" />
              </div>
              <div class="col-4">
                <label>{{ $t('theme') }}</label
                ><Dropdown
                  v-model="form.theme"
                  :options="['light', 'dark', 'colorful']"
                  class="w-full"
                />
              </div>
              <div class="col-4 flex align-items-end">
                <div class="flex align-items-center">
                  <Checkbox v-model="form.show_categories" :binary="true" inputId="sc" /><label
                    for="sc"
                    class="ml-2"
                    >{{ $t('show_categories') }}</label
                  >
                </div>
              </div>
              <div class="col-4 flex align-items-end">
                <div class="flex align-items-center">
                  <Checkbox v-model="form.active" :binary="true" inputId="ak" /><label
                    for="ak"
                    class="ml-2"
                    >{{ $t('active') }}</label
                  >
                </div>
              </div>
            </div>
            <Button :label="$t('save')" class="mt-2" @click="doSave" :loading="saving" /> </template
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
import Dropdown from 'primevue/dropdown'
import Checkbox from 'primevue/checkbox'
const { t } = useI18n()
const toast = useToast()
interface KioskConfig {
  id: string
  name: string
  location: string
  theme: string
  show_categories: boolean
  show_images: boolean
  allow_customize: boolean
  auto_send_to_kitchen: boolean
  active: boolean
}
const configs = ref<KioskConfig[]>([])
const loading = ref(false)
const saving = ref(false)
const editing = ref(false)
const form = reactive<Omit<KioskConfig, 'id'> & { id: string }>({
  id: '',
  name: '',
  location: '',
  theme: 'light',
  show_categories: true,
  show_images: true,
  allow_customize: true,
  auto_send_to_kitchen: true,
  active: true,
})
const edit = (d: KioskConfig) => {
  Object.assign(form, d)
  editing.value = true
}
const doSave = async () => {
  saving.value = true
  try {
    await axios.post(`http://${import.meta.env.VITE_APP_BACKEND_HOST}/kiosk/api/configs`, form, {
      headers: { Authorization: `Bearer ${auth.accessToken.value}` },
    })
    toast.add({ severity: 'success', summary: t('saved'), group: 'br', life: 2000 })
    editing.value = false
    form.id = ''
    form.name = ''
    form.location = ''
    await load()
  } catch {
    toast.add({ severity: 'error', summary: t('failed'), group: 'br', life: 3000 })
  } finally {
    saving.value = false
  }
}
const load = async () => {
  loading.value = true
  try {
    const r = await axios.get(`http://${import.meta.env.VITE_APP_BACKEND_HOST}/kiosk/api/configs`, {
      headers: { Authorization: `Bearer ${auth.accessToken.value}` },
    })
    configs.value = r.data.data || []
  } catch {
    toast.add({ severity: 'error', summary: t('load_failed'), group: 'br', life: 3000 })
  } finally {
    loading.value = false
  }
}
load()
</script>

<template>
  <div class="w-full"><div class="grid mx-2"><div class="col-12"><h3>{{ $t('marketing_campaigns') }}</h3></div>
    <div class="col-12"><DataTable :value="campaigns" stripedRows :loading="loading">
      <template #empty>{{ $t('no_results') }}</template>
      <Column field="name" :header="$t('name')"></Column><Column field="type" :header="$t('type')"></Column>
      <Column field="discount_pct" :header="$t('discount')"><template #body="s">{{ s.data.discount_pct }}%</template></Column>
      <Column field="start_date" :header="$t('start_date')"></Column><Column field="end_date" :header="$t('end_date')"></Column>
      <Column field="active" :header="$t('active')"><template #body="s"><i :class="s.data.active?'pi pi-check text-green-500':'pi pi-times text-300'"></i></template></Column>
      <Column :header="$t('actions')"><template #body="s">
        <Button icon="pi pi-power-off" :severity="s.data.active?'warn':'success'" size="small" class="mr-1" @click="toggle(s.data.id)" />
        <Button icon="pi pi-trash" severity="danger" size="small" @click="doDelete(s.data.id)" />
      </template></Column>
    </DataTable></div>
    <div class="col-12 mt-3"><Card><template #title>{{ $t('new_campaign') }}</template><template #content>
      <div class="grid"><div class="col-4"><label>{{ $t('name') }}</label><InputText v-model="form.name" class="w-full" /></div>
      <div class="col-4"><label>{{ $t('type') }}</label><Dropdown v-model="form.type" :options="['discount','loyalty','referral','seasonal']" class="w-full" /></div>
      <div class="col-4"><label>{{ $t('discount') }} (%)</label><InputNumber v-model="form.discount_pct" :min="0" :max="100" class="w-full" /></div>
      <div class="col-6"><label>{{ $t('start_date') }}</label><Calendar v-model="formStart" dateFormat="yy-mm-dd" class="w-full" /></div>
      <div class="col-6"><label>{{ $t('end_date') }}</label><Calendar v-model="formEnd" dateFormat="yy-mm-dd" class="w-full" /></div>
      <div class="col-12"><label>{{ $t('description') }}</label><Textarea v-model="form.description" class="w-full" rows="2" /></div></div>
      <Button :label="$t('create')" class="mt-2" @click="doCreate" :loading="creating" />
    </template></Card></div>
  </div></div>
</template>
<script setup lang="ts">
import { ref, reactive } from 'vue'; import axios from 'axios'; import { useI18n } from 'vue-i18n'; import { useToast } from 'primevue/usetoast'; import auth from '../services/auth'
import Button from 'primevue/button'; import DataTable from 'primevue/datatable'; import Column from 'primevue/column'; import Card from 'primevue/card'; import InputText from 'primevue/inputtext'; import InputNumber from 'primevue/inputnumber'; import Dropdown from 'primevue/dropdown'; import Calendar from 'primevue/calendar'; import Textarea from 'primevue/textarea'
const { t } = useI18n(); const toast = useToast()
const campaigns = ref<any[]>([]); const loading = ref(false); const creating = ref(false); const formStart = ref<Date|null>(null); const formEnd = ref<Date|null>(null)
const form = reactive({ name:'', type:'discount', description:'', discount_pct:10 })
const toggle = async (id:string) => { try { await axios.post(`http://${import.meta.env.VITE_APP_BACKEND_HOST}/marketing/api/campaigns/${id}/toggle`,{},{headers:{Authorization:`Bearer ${auth.accessToken.value}`}}); await load() } catch { toast.add({severity:'error',summary:t('failed'),group:'br',life:3000}) } }
const doDelete = async (id:string) => { try { await axios.delete(`http://${import.meta.env.VITE_APP_BACKEND_HOST}/marketing/api/campaigns/${id}`,{headers:{Authorization:`Bearer ${auth.accessToken.value}`}}); await load() } catch { toast.add({severity:'error',summary:t('failed'),group:'br',life:3000}) } }
const doCreate = async () => { creating.value=true; try { await axios.post(`http://${import.meta.env.VITE_APP_BACKEND_HOST}/marketing/api/campaigns`,{...form,start_date:formStart.value?.toISOString().split('T')[0],end_date:formEnd.value?.toISOString().split('T')[0]},{headers:{Authorization:`Bearer ${auth.accessToken.value}`}}); toast.add({severity:'success',summary:t('saved'),group:'br',life:2000}); form.name='';form.description='';form.discount_pct=10;formStart.value=null;formEnd.value=null; await load() } catch { toast.add({severity:'error',summary:t('failed'),group:'br',life:3000}) } finally { creating.value=false } }
const load = async () => { loading.value=true; try { const r=await axios.get(`http://${import.meta.env.VITE_APP_BACKEND_HOST}/marketing/api/campaigns`,{headers:{Authorization:`Bearer ${auth.accessToken.value}`}}); campaigns.value=r.data.data||[] } catch { toast.add({severity:'error',summary:t('load_failed'),group:'br',life:3000}) } finally { loading.value=false } }
load()
</script>

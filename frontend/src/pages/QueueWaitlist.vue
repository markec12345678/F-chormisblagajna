<template>
  <div class="w-full"><div class="grid mx-2"><div class="col-12"><h3>{{ $t('queue_waitlist') }}</h3></div>
    <div class="col-12 md:col-6"><DataTable :value="queue" stripedRows :loading="loading">
      <template #empty>{{ $t('queue_empty') }}</template>
      <Column field="position" header="#">
        <template #body="s"><Tag :value="s.data.position" /></template>
      </Column>
      <Column field="customer_name" :header="$t('customer')"></Column><Column field="party_size" :header="$t('guests')"></Column>
      <Column field="estimated_min" :header="$t('estimated_wait')"><template #body="s">{{ s.data.estimated_min }} min</template></Column>
      <Column field="status" :header="$t('status')"><template #body="s"><Tag :severity="s.data.status==='waiting'?'warn':'success'">{{ s.data.status }}</Tag></template></Column>
      <Column :header="$t('actions')"><template #body="s">
        <Button v-if="s.data.status==='waiting'" icon="pi pi-check" severity="success" size="small" class="mr-1" @click="seat(s.data.id)" />
        <Button icon="pi pi-times" severity="danger" size="small" @click="remove(s.data.id)" />
      </template></Column>
    </DataTable></div>
    <div class="col-12 md:col-6"><Card><template #title>{{ $t('add_to_queue') }}</template><template #content>
      <div class="grid"><div class="col-6"><label>{{ $t('customer_name') }}</label><InputText v-model="form.customer_name" class="w-full" /></div>
      <div class="col-6"><label>{{ $t('phone') }}</label><InputText v-model="form.phone" class="w-full" /></div>
      <div class="col-6"><label>{{ $t('guests') }}</label><InputNumber v-model="form.party_size" :min="1" class="w-full" /></div>
      <div class="col-12"><label>{{ $t('notes') }}</label><InputText v-model="form.notes" class="w-full" /></div></div>
      <Button :label="$t('add')" class="mt-2" @click="doAdd" :loading="adding" />
    </template></Card></div>
  </div></div>
</template>
<script setup lang="ts">
import { ref, reactive } from 'vue'; import axios from 'axios'; import { useI18n } from 'vue-i18n'; import { useToast } from 'primevue/usetoast'; import auth from '../services/auth'
import Button from 'primevue/button'; import DataTable from 'primevue/datatable'; import Column from 'primevue/column'; import Card from 'primevue/card'; import InputText from 'primevue/inputtext'; import InputNumber from 'primevue/inputnumber'; import Tag from 'primevue/tag'
const { t } = useI18n(); const toast = useToast()
interface QueueEntry { id: string; customer_name: string; phone: string; party_size: number; estimated_min: number; status: string; position: number; notes: string }
const queue = ref<QueueEntry[]>([]); const loading = ref(false); const adding = ref(false); const form = reactive({ customer_name:'', phone:'', party_size:2, notes:'' })
const seat = async (id:string) => { try { await axios.put(`http://${import.meta.env.VITE_APP_BACKEND_HOST}/queue/api/queue/${id}/status`,{status:'seated'},{headers:{Authorization:`Bearer ${auth.accessToken.value}`}}); await load() } catch { toast.add({severity:'error',summary:t('failed'),group:'br',life:3000}) } }
const remove = async (id:string) => { try { await axios.delete(`http://${import.meta.env.VITE_APP_BACKEND_HOST}/queue/api/queue/${id}`,{headers:{Authorization:`Bearer ${auth.accessToken.value}`}}); await load() } catch { toast.add({severity:'error',summary:t('failed'),group:'br',life:3000}) } }
const doAdd = async () => { adding.value=true; try { await axios.post(`http://${import.meta.env.VITE_APP_BACKEND_HOST}/queue/api/queue`,form,{headers:{Authorization:`Bearer ${auth.accessToken.value}`}}); form.customer_name='';form.phone='';form.party_size=2;form.notes=''; await load() } catch { toast.add({severity:'error',summary:t('failed'),group:'br',life:3000}) } finally { adding.value=false } }
const load = async () => { loading.value=true; try { const r=await axios.get(`http://${import.meta.env.VITE_APP_BACKEND_HOST}/queue/api/queue`,{headers:{Authorization:`Bearer ${auth.accessToken.value}`}}); queue.value=r.data.data||[] } catch { toast.add({severity:'error',summary:t('load_failed'),group:'br',life:3000}) } finally { loading.value=false } }
load()
</script>

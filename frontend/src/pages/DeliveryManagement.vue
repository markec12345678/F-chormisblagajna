<template>
  <div class="w-full"><div class="grid mx-2"><div class="col-12"><h3>{{ $t('delivery_management') }}</h3></div>
    <div class="col-12 md:col-6"><TabView><TabPanel :header="$t('zones')">
      <DataTable :value="zones" stripedRows :loading="loading">
        <template #empty>{{ $t('no_results') }}</template>
        <Column field="name" :header="$t('name')"></Column><Column field="fee" :header="$t('fee')"><template #body="s">{{ formatCurrency(s.data.fee) }}</template></Column>
        <Column field="min_order" :header="$t('min_order')"><template #body="s">{{ formatCurrency(s.data.min_order) }}</template></Column>
        <Column field="active" :header="$t('active')"><template #body="s"><i :class="s.data.active?'pi pi-check text-green-500':'pi pi-times text-300'"></i></template></Column>
        <Column :header="$t('actions')"><template #body="s"><Button icon="pi pi-trash" severity="danger" size="small" @click="deleteZone(s.data.id)" /></template></Column>
      </DataTable>
      <div class="flex gap-2 mt-2"><InputText v-model="zoneForm.name" placeholder="Ime" /><InputNumber v-model="zoneForm.fee" placeholder="Cena" /><InputNumber v-model="zoneForm.min_order" placeholder="Min" /><Button icon="pi pi-plus" @click="addZone" /></div>
    </TabPanel><TabPanel :header="$t('orders')">
      <DataTable :value="delOrders" stripedRows :loading="loading">
        <template #empty>{{ $t('no_results') }}</template>
        <Column field="customer_name" :header="$t('customer')"></Column><Column field="address" :header="$t('address')"></Column>
        <Column field="delivery_fee" :header="$t('fee')"><template #body="s">{{ formatCurrency(s.data.delivery_fee) }}</template></Column>
        <Column field="status" :header="$t('status')"><template #body="s"><Tag :severity="s.data.status==='pending'?'warn':s.data.status==='in_transit'?'info':'success'">{{ s.data.status }}</Tag></template></Column>
        <Column :header="$t('actions')"><template #body="s">
          <Button v-if="s.data.status==='pending'" icon="pi pi-truck" size="small" class="mr-1" @click="updateStatus(s.data.id,'in_transit')" />
          <Button v-if="s.data.status==='in_transit'" icon="pi pi-check" severity="success" size="small" @click="updateStatus(s.data.id,'delivered')" />
        </template></Column>
      </DataTable>
    </TabPanel></TabView></div>
  </div></div>
</template>
<script setup lang="ts">
import { ref, reactive } from 'vue'; import axios from 'axios'; import { useI18n } from 'vue-i18n'; import { useToast } from 'primevue/usetoast'; import auth from '../services/auth'
import Button from 'primevue/button'; import DataTable from 'primevue/datatable'; import Column from 'primevue/column'; import InputText from 'primevue/inputtext'; import InputNumber from 'primevue/inputnumber'; import TabView from 'primevue/tabview'; import TabPanel from 'primevue/tabpanel'; import Tag from 'primevue/tag'
const { t } = useI18n(); const toast = useToast()
const zones = ref<any[]>([]); const delOrders = ref<any[]>([]); const loading = ref(false); const zoneForm = reactive({ name: '', fee: 0, min_order: 0 })
const formatCurrency = (n: number) => new Intl.NumberFormat('sl-SI',{style:'currency',currency:'EUR'}).format(n||0)
const addZone = async () => { try { await axios.post(`http://${import.meta.env.VITE_APP_BACKEND_HOST}/delivery/api/zones`,{...zoneForm,active:true},{headers:{Authorization:`Bearer ${auth.accessToken.value}`}}); zoneForm.name='';zoneForm.fee=0;zoneForm.min_order=0; await load() } catch { toast.add({severity:'error',summary:t('failed'),group:'br',life:3000}) } }
const deleteZone = async (id:string) => { try { await axios.delete(`http://${import.meta.env.VITE_APP_BACKEND_HOST}/delivery/api/zones/${id}`,{headers:{Authorization:`Bearer ${auth.accessToken.value}`}}); await load() } catch { toast.add({severity:'error',summary:t('failed'),group:'br',life:3000}) } }
const updateStatus = async (id:string,s:string) => { try { await axios.put(`http://${import.meta.env.VITE_APP_BACKEND_HOST}/delivery/api/orders/${id}/status`,{status:s},{headers:{Authorization:`Bearer ${auth.accessToken.value}`}}); await load() } catch { toast.add({severity:'error',summary:t('failed'),group:'br',life:3000}) } }
const load = async () => { loading.value=true; try { const [z,r]=await Promise.all([axios.get(`http://${import.meta.env.VITE_APP_BACKEND_HOST}/delivery/api/zones`,{headers:{Authorization:`Bearer ${auth.accessToken.value}`}}),axios.get(`http://${import.meta.env.VITE_APP_BACKEND_HOST}/delivery/api/orders`,{headers:{Authorization:`Bearer ${auth.accessToken.value}`}})]); zones.value=z.data.data||[]; delOrders.value=r.data.data||[] } catch { toast.add({severity:'error',summary:t('load_failed'),group:'br',life:3000}) } finally { loading.value=false } }
load()
</script>

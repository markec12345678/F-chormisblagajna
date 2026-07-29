<template>
  <div class="w-full"><div class="grid mx-2"><div class="col-12"><h3>{{ $t('gift_cards') }}</h3></div>
    <div class="col-12"><DataTable :value="cards" stripedRows :loading="loading">
      <template #empty>{{ $t('no_data') }}</template>
      <Column field="code" :header="$t('code')"></Column><Column field="balance" :header="$t('balance')"><template #body="s">{{ formatCurrency(s.data.balance) }}</template></Column>
      <Column field="initial_amt" :header="$t('initial_amount')"><template #body="s">{{ formatCurrency(s.data.initial_amt) }}</template></Column>
      <Column field="issued_to" :header="$t('issued_to')"></Column><Column field="issued_at" :header="$t('issued')"></Column>
      <Column field="active" :header="$t('active')"><template #body="s"><i :class="s.data.active?'pi pi-check text-green-500':'pi pi-times text-300'"></i></template></Column>
      <Column :header="$t('actions')"><template #body="s"><Button v-if="s.data.active" icon="pi pi-ban" severity="danger" size="small" @click="deactivate(s.data.id)" /></template></Column>
    </DataTable></div>
    <div class="col-12 mt-3"><Card><template #title>{{ $t('issue_gift_card') }}</template><template #content>
      <div class="grid"><div class="col-6"><label>{{ $t('issued_to') }}</label><InputText v-model="form.issued_to" class="w-full" /></div>
      <div class="col-6"><label>{{ $t('initial_amount') }}</label><InputNumber v-model="form.initial_amt" :min="1" class="w-full" /></div>
      <div class="col-6"><label>{{ $t('expiry_date') }}</label><Calendar v-model="formExp" dateFormat="yy-mm-dd" class="w-full" /></div></div>
      <Button :label="$t('issue')" class="mt-2" @click="doIssue" :loading="issuing" />
    </template></Card></div>
  </div></div>
</template>
<script setup lang="ts">
import { ref, reactive } from 'vue'; import axios from 'axios'; import { useI18n } from 'vue-i18n'; import { useToast } from 'primevue/usetoast'; import auth from '../services/auth'
import Button from 'primevue/button'; import DataTable from 'primevue/datatable'; import Column from 'primevue/column'; import Card from 'primevue/card'; import InputText from 'primevue/inputtext'; import InputNumber from 'primevue/inputnumber'; import Calendar from 'primevue/calendar'
const { t } = useI18n(); const toast = useToast()
const cards = ref<any[]>([]); const loading = ref(false); const issuing = ref(false); const formExp = ref<Date|null>(null)
const form = reactive({ issued_to:'', initial_amt:10 })
const formatCurrency = (n:number) => new Intl.NumberFormat('sl-SI',{style:'currency',currency:'EUR'}).format(n||0)
const doIssue = async () => { issuing.value=true; try { await axios.post(`http://${import.meta.env.VITE_APP_BACKEND_HOST}/giftcards/api/cards`,{...form,expires_at:formExp.value?.toISOString().split('T')[0]},{headers:{Authorization:`Bearer ${auth.accessToken.value}`}}); toast.add({severity:'success',summary:t('saved'),group:'br',life:2000}); form.issued_to='';form.initial_amt=10;formExp.value=null; await load() } catch { toast.add({severity:'error',summary:t('failed'),group:'br',life:3000}) } finally { issuing.value=false } }
const deactivate = async (id:string) => { try { await axios.post(`http://${import.meta.env.VITE_APP_BACKEND_HOST}/giftcards/api/cards/${id}/deactivate`,{},{headers:{Authorization:`Bearer ${auth.accessToken.value}`}}); await load() } catch { toast.add({severity:'error',summary:t('failed'),group:'br',life:3000}) } }
const load = async () => { loading.value=true; try { const r=await axios.get(`http://${import.meta.env.VITE_APP_BACKEND_HOST}/giftcards/api/cards`,{headers:{Authorization:`Bearer ${auth.accessToken.value}`}}); cards.value=r.data.data||[] } catch { toast.add({severity:'error',summary:t('load_failed'),group:'br',life:3000}) } finally { loading.value=false } }
load()
</script>

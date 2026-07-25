<template>
  <div>
    <DataTable
      filterDisplay="row"
      :loading="loading"
      v-model:filters="filters"
      :globalFilterFields="['name']"
      :value="customers"
      stripedRows
      tableStyle="min-width: 100%"
      class="w-full pr-5"
    >
      <template #header>
        <div class="flex justify-content-start">
          <IconField iconPosition="left">
            <InputIcon>
              <i class="pi pi-search" />
            </InputIcon>
            <InputText v-model="filters['name'].value" :placeholder="$t('search_by_name')" />
          </IconField>
          <Button
            icon="pi pi-plus"
            severity="secondary"
            @click="add_customer_dialog = true"
            class="ml-2"
          />
        </div>
      </template>
      <template #empty>
        <div class="flex flex-column align-items-center py-4">
          <i class="pi pi-inbox text-4xl text-gray-400 mb-2"></i>
          <span class="text-gray-500">{{ $t('no_results') }}</span>
        </div>
      </template>
      <Column field="name" :header="$t('name')"></Column>
      <Column field="phone" :header="$t('phone')"></Column>
      <Column field="address" :header="$t('address')"></Column>
      <Column :header="$t('actions')">
        <template #body="slotProps">
          <ButtonGroup>
            <Button
              :label="$t('choose')"
              severity="secondary"
              aria-label="$t('choose')"
              @click="returnCustomer(slotProps.data)"
            />
          </ButtonGroup>
        </template>
      </Column>
    </DataTable>
    <AddCustomer
      @update:visible="(x) => (add_customer_dialog = x)"
      :visible="add_customer_dialog"
      @customer-added="handleCustomerAdded"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, defineEmits } from 'vue'
import Button from 'primevue/button'
import InputText from 'primevue/inputtext'
import InputIcon from 'primevue/inputicon'
import IconField from 'primevue/iconfield'
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import axios from 'axios'
import AddCustomer from '@/components/AddCustomer.vue'
import auth from '../services/auth'

const add_customer_dialog = ref(false)
const customers = ref([])
const loading = ref(false)

const filters = ref({
  name: { value: null },
})

const emit = defineEmits(['returnCustomer'])

const returnCustomer = (customer) => {
  emit('returnCustomer', customer)
}

const handleCustomerAdded = () => {
  add_customer_dialog.value = false
  GetCustomers()
}

const GetCustomers = (page_number = 1, page_size = 9999999999) => {
  loading.value = true
  axios
    .get(
      `http://${import.meta.env.VITE_APP_BACKEND_HOST}${import.meta.env.VITE_APP_MODULE_CORE_API_PREFIX}/api/customers?page[number]=${page_number}&page[number]=${page_size}`,
      {
        headers: {
          Authorization: `Bearer ${auth.accessToken.value}`,
        },
      },
    )
    .then((response) => {
      customers.value = response.data.data
      loading.value = false
    })
}

GetCustomers()
</script>

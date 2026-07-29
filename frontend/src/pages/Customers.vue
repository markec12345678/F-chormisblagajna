<template>
  <div class="w-full">
    <div class="grid mx-2">
      <div class="col-12 flex">
        <div class="grid w-full">
          <div class="col-12">
            <h3>{{ $t('customer', 3) }}</h3>
          </div>
          <div class="col-12">
            <DataTable
              @page="updatCustomersTableRowsPerPage"
              :lazy="true"
              :totalRecords="customersTableTotalRecords"
              :loading="isCustomersTableLoading"
              paginatorPosition="both"
              paginator
              :rows="customersTableRowsPerPage"
              :rowsPerPageOptions="[50, 100, 500]"
              :value="customers"
              stripedRows
              tableStyle="width: 100%;max-height:50vh;"
              class="w-full pr-2"
            >
              <template #header>
                <div class="flex justify-between align-items-center">
                  <Button
                    icon="pi pi-plus"
                    :label="$t('add_customer')"
                    rounded
                    raised
                    @click="customerAddDialog = true"
                  />
                  <div class="flex gap-2">
                    <Dropdown
                      v-model="selectedTagFilter"
                      :options="availableTags"
                      :placeholder="$t('filter_by_tag')"
                      class="w-48"
                      showClear
                    />
                  </div>
                </div>
              </template>
              <template #empty>
                <div class="flex flex-column align-items-center gap-2 py-4">
                  <i class="pi pi-user" style="font-size: 2rem; opacity: 0.3"></i>
                  <p class="m-0 text-slate-400">{{ $t('no_results') }}</p>
                </div>
              </template>
              <Column sortable field="name" :header="$t('name')"></Column>
              <Column field="email" :header="$t('email')"></Column>
              <Column field="phone" :header="$t('phone', 1)"></Column>
              <Column field="address" :header="$t('address')"></Column>
              <Column :header="$t('tags')">
                <template #body="slotProps">
                  <Tag
                    v-for="tag in slotProps?.data?.tags || []"
                    :key="tag"
                    :value="tag"
                    class="mr-1"
                    severity="info"
                  />
                </template>
              </Column>
              <Column sortable field="loyalty_points" :header="$t('loyalty_points')">
                <template #body="slotProps">
                  <Tag :value="String(slotProps?.data?.loyalty_points || 0)" severity="success" />
                </template>
              </Column>
              <Column sortable field="total_spent" :header="$t('total_spent')">
                <template #body="slotProps">
                  {{ formatCurrency(slotProps?.data?.total_spent || 0) }}
                </template>
              </Column>
              <Column :header="$t('actions')">
                <template #body="slotProps">
                  <ConfirmPopup></ConfirmPopup>
                  <ButtonGroup>
                    <Button
                      icon="pi pi-eye"
                      severity="info"
                      :aria-label="$t('customer_details')"
                      @click="showCustomerDetails(slotProps?.data)"
                    />
                    <Button
                      icon="pi pi-pencil"
                      severity="secondary"
                      :aria-label="$t('edit')"
                      @click="prepareCustomerToEdit(slotProps?.data)"
                    />
                    <Button
                      icon="pi pi-trash"
                      severity="danger"
                      :aria-label="$t('remove')"
                      @click="confirmDeleteCustomer($event, slotProps?.data?.id)"
                    />
                  </ButtonGroup>
                </template>
              </Column>
            </DataTable>

            <Dialog
              v-model:visible="customerAddDialog"
              modal
              :header="$t('add_customer')"
              :style="{ width: '75rem', direction: store.orientation == 'rtl' ? 'rtl' : 'ltr' }"
              :breakpoints="{ '1199px': '50vw', '575px': '90vw' }"
            >
              <div class="grid">
                <div class="col-12 md:col-6">
                  <div class="flex flex-column gap-2">
                    <label for="name">{{ $t('name') }}</label>
                    <InputText
                      id="name"
                      v-model="new_customer_name"
                      aria-describedby="name"
                      :class="{ 'p-invalid': customer_add_errors.name }"
                    />
                    <small class="p-error">{{ customer_add_errors.name }}</small>
                  </div>
                </div>
                <div class="col-12 md:col-6">
                  <div class="flex flex-column gap-2">
                    <label for="email">{{ $t('email') }}</label>
                    <InputText
                      id="email"
                      v-model="new_customer_email"
                      aria-describedby="email"
                      type="email"
                    />
                  </div>
                </div>
                <div class="col-12 md:col-6">
                  <div class="flex flex-column gap-2">
                    <label for="phone">{{ $t('phone') }}</label>
                    <InputText v-model="new_customer_phone" aria-describedby="phone" />
                  </div>
                </div>
                <div class="col-12 md:col-6">
                  <div class="flex flex-column gap-2">
                    <label for="address">{{ $t('address') }}</label>
                    <InputText
                      id="address"
                      v-model="new_customer_address"
                      aria-describedby="address"
                    />
                  </div>
                </div>
                <div class="col-12">
                  <div class="flex flex-column gap-2">
                    <label for="notes">{{ $t('notes') }}</label>
                    <Textarea id="notes" v-model="new_customer_notes" rows="3" />
                  </div>
                </div>
              </div>
              <template #footer>
                <ButtonGroup>
                  <Button :label="$t('cancel')" severity="secondary" :aria-label="$t('cancel')" />
                  <Button
                    class="ml-2"
                    severity="primary"
                    @click="submitCustomer"
                    :label="$t('save')"
                    :aria-label="$t('save')"
                    :loading="isSubmitting"
                    :disabled="isSubmitting"
                  />
                </ButtonGroup>
              </template>
            </Dialog>

            <Dialog
              v-model:visible="customerEditDialog"
              modal
              :header="`${$t('edit_customer')} ${customerToEdit.name}`"
              :style="{ width: '75rem', direction: store.orientation == 'rtl' ? 'rtl' : 'ltr' }"
              :breakpoints="{ '1199px': '50vw', '575px': '90vw' }"
            >
              <div class="grid">
                <div class="col-12 md:col-6">
                  <div class="flex flex-column gap-2">
                    <label for="name">{{ $t('name') }}</label>
                    <InputText
                      id="name"
                      v-model="customerToEdit.name"
                      aria-describedby="name"
                      :class="{ 'p-invalid': customer_edit_errors.name }"
                    />
                    <small class="p-error">{{ customer_edit_errors.name }}</small>
                  </div>
                </div>
                <div class="col-12 md:col-6">
                  <div class="flex flex-column gap-2">
                    <label for="email">{{ $t('email') }}</label>
                    <InputText
                      id="email"
                      v-model="customerToEdit.email"
                      aria-describedby="email"
                      type="email"
                    />
                  </div>
                </div>
                <div class="col-12 md:col-6">
                  <div class="flex flex-column gap-2">
                    <label for="phone">{{ $t('phone') }}</label>
                    <InputText id="phone" v-model="customerToEdit.phone" aria-describedby="phone" />
                  </div>
                </div>
                <div class="col-12 md:col-6">
                  <div class="flex flex-column gap-2">
                    <label for="address">{{ $t('address') }}</label>
                    <InputText
                      id="address"
                      v-model="customerToEdit.address"
                      aria-describedby="address"
                    />
                  </div>
                </div>
                <div class="col-12">
                  <div class="flex flex-column gap-2">
                    <label for="notes">{{ $t('notes') }}</label>
                    <Textarea id="notes" v-model="customerToEdit.notes" rows="3" />
                  </div>
                </div>
              </div>
              <template #footer>
                <ButtonGroup>
                  <Button :label="$t('cancel')" severity="secondary" :aria-label="$t('cancel')" />
                  <Button
                    class="ml-2"
                    severity="primary"
                    @click="updateCustomer"
                    :label="$t('save')"
                    :aria-label="$t('save')"
                    :loading="isSubmitting"
                    :disabled="isSubmitting"
                  />
                </ButtonGroup>
              </template>
            </Dialog>

            <Dialog
              v-model:visible="customerDetailsDialog"
              modal
              :header="$t('customer_details')"
              :style="{ width: '90rem', direction: store.orientation == 'rtl' ? 'rtl' : 'ltr' }"
              :breakpoints="{ '1199px': '70vw', '575px': '95vw' }"
            >
              <div class="grid" v-if="selectedCustomer">
                <div class="col-12 md:col-4">
                  <Card>
                    <template #title>{{ $t('customer_stats') }}</template>
                    <template #content>
                      <div class="flex flex-column gap-3">
                        <div class="flex justify-content-between">
                          <span class="text-500">{{ $t('total_spent') }}:</span>
                          <span class="font-bold">{{
                            formatCurrency(selectedCustomer.total_spent || 0)
                          }}</span>
                        </div>
                        <div class="flex justify-content-between">
                          <span class="text-500">{{ $t('order_count') }}:</span>
                          <span class="font-bold">{{ selectedCustomer.order_count || 0 }}</span>
                        </div>
                        <div class="flex justify-content-between">
                          <span class="text-500">{{ $t('loyalty_points') }}:</span>
                          <Tag :value="selectedCustomer.loyalty_points || 0" severity="success" />
                        </div>
                        <div class="flex justify-content-between">
                          <span class="text-500">{{ $t('last_order_date') }}:</span>
                          <span>{{
                            selectedCustomer.last_order_date
                              ? new Date(selectedCustomer.last_order_date).toLocaleDateString()
                              : '-'
                          }}</span>
                        </div>
                      </div>
                    </template>
                  </Card>
                </div>
                <div class="col-12 md:col-8">
                  <Card>
                    <template #title>{{ $t('order_history') }}</template>
                    <template #content>
                      <DataTable :value="customerOrders" stripedRows tableStyle="width: 100%">
                        <template #empty>
                          <div class="flex flex-column align-items-center gap-2 py-4">
                            <i
                              class="pi pi-shopping-cart"
                              style="font-size: 2rem; opacity: 0.3"
                            ></i>
                            <p class="m-0 text-slate-400">{{ $t('no_orders') }}</p>
                          </div>
                        </template>
                        <Column field="display_id" :header="$t('id')"></Column>
                        <Column field="submitted_at" :header="$t('date')">
                          <template #body="slotProps">
                            {{ new Date(slotProps.data.submitted_at).toLocaleDateString() }}
                          </template>
                        </Column>
                        <Column field="sale_price" :header="$t('total')">
                          <template #body="slotProps">
                            {{
                              formatCurrency(slotProps.data.sale_price - slotProps.data.discount)
                            }}
                          </template>
                        </Column>
                        <Column field="state" :header="$t('status')">
                          <template #body="slotProps">
                            <Tag
                              :value="slotProps.data.state"
                              :severity="getStateSeverity(slotProps.data.state)"
                            />
                          </template>
                        </Column>
                        <Column field="is_paid" :header="$t('payment')">
                          <template #body="slotProps">
                            <Tag
                              :value="slotProps.data.is_paid ? $t('paid') : $t('pending')"
                              :severity="slotProps.data.is_paid ? 'success' : 'warn'"
                            />
                          </template>
                        </Column>
                      </DataTable>
                    </template>
                  </Card>
                </div>
              </div>
            </Dialog>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import ConfirmPopup from 'primevue/confirmpopup'
import DataTable from 'primevue/datatable'
import Dialog from 'primevue/dialog'
import InputText from 'primevue/inputtext'
import Textarea from 'primevue/textarea'
import Column from 'primevue/column'
import Button from 'primevue/button'
import ButtonGroup from 'primevue/buttongroup'
import Card from 'primevue/card'
import Tag from 'primevue/tag'
import Dropdown from 'primevue/dropdown'
import axios from 'axios'
import { ref, computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { useToast } from 'primevue/usetoast'
import { useConfirm } from 'primevue/useconfirm'
import { globalStore } from '../stores'
import auth from '../services/auth'
import type { Customer, DataTablePageEvent } from '@/types'

const { t } = useI18n()
const store = globalStore()
const confirm = useConfirm()

interface Order {
  id: string
  display_id: string
  submitted_at: string
  sale_price: number
  discount: number
  state: string
  is_paid: boolean
}

const customerToEdit = ref<Partial<Customer>>({})
const customerEditDialog = ref(false)
const customer_add_errors = ref({ name: '' })
const customer_edit_errors = ref({ name: '' })

const isSubmitting = ref(false)
const customerAddDialog = ref(false)
const new_customer_name = ref('')
const new_customer_phone = ref('')
const new_customer_id = ref('')
const new_customer_address = ref('')
const new_customer_email = ref('')
const new_customer_notes = ref('')

const customersTableTotalRecords = ref(0)
const customersTableRowsPerPage = ref(50)
const isCustomersTableLoading = ref(true)
const customers = ref<Customer[]>([])
const toast = useToast()

const customerDetailsDialog = ref(false)
const selectedCustomer = ref<Customer | null>(null)
const customerOrders = ref<Order[]>([])
const selectedTagFilter = ref<string | null>(null)

const availableTags = computed(() => {
  const tags = new Set<string>()
  customers.value.forEach((c) => {
    if (c.tags) {
      c.tags.forEach((tag) => tags.add(tag))
    }
  })
  return Array.from(tags).sort()
})

const formatCurrency = (amount: number) => {
  return new Intl.NumberFormat('sl-SI', { style: 'currency', currency: 'EUR' }).format(amount)
}

const getStateSeverity = (state: string) => {
  const map: Record<string, string> = {
    completed: 'success',
    in_progress: 'warn',
    pending: 'info',
    cancelled: 'danger',
  }
  return map[state] || 'secondary'
}

const showCustomerDetails = async (customer: Customer) => {
  selectedCustomer.value = customer
  customerDetailsDialog.value = true

  try {
    const response = await axios.get(
      `http://${import.meta.env.VITE_APP_BACKEND_HOST}${import.meta.env.VITE_APP_MODULE_CORE_API_PREFIX}/api/customers/${customer.id}/orders`,
      { headers: { Authorization: `Bearer ${auth.accessToken.value}` } },
    )
    customerOrders.value = response.data.data || []
  } catch {
    customerOrders.value = []
  }
}

const deleteCustomer = (customer_id: string) => {
  axios
    .delete(
      `http://${import.meta.env.VITE_APP_BACKEND_HOST}${import.meta.env.VITE_APP_MODULE_CORE_API_PREFIX}/api/customers/${customer_id}`,
      {
        headers: {
          Authorization: `Bearer ${auth.accessToken.value}`,
        },
      },
    )
    .then(() => {
      toast.add({
        severity: 'success',
        summary: t('success'),
        detail: t('customer_deleted_success'),
        group: 'br',
        life: 3000,
      })
      loadCustomers()
    })
    .catch(() => {
      toast.add({
        severity: 'error',
        summary: t('failed'),
        detail: t('customer_delete_failed'),
        group: 'br',
        life: 3000,
      })
    })
}

const confirmDeleteCustomer = (event: MouseEvent, customer_id: string) => {
  confirm.require({
    target: event.currentTarget as HTMLElement,
    message: t('confirm_delete_customer'),
    icon: 'pi pi-exclamation-triangle',
    rejectProps: {
      label: t('cancel'),
      severity: 'secondary',
      outlined: true,
    },
    acceptProps: {
      label: t('yes'),
    },
    accept: () => {
      deleteCustomer(customer_id)
    },
    reject: () => {},
  })
}

const prepareCustomerToEdit = (customer: Customer) => {
  customerToEdit.value = JSON.parse(JSON.stringify(customer))
  customerEditDialog.value = true
}

const updateCustomer = () => {
  customer_edit_errors.value.name = customerToEdit.value.name?.trim()
    ? ''
    : t('validation_required')

  if (customer_edit_errors.value.name) return

  isSubmitting.value = true

  axios
    .patch(
      `http://${import.meta.env.VITE_APP_BACKEND_HOST}${import.meta.env.VITE_APP_MODULE_CORE_API_PREFIX}/api/customers/${customerToEdit.value.id}`,
      {
        data: customerToEdit.value,
      },
      {
        headers: {
          Authorization: `Bearer ${auth.accessToken.value}`,
        },
      },
    )
    .then(() => {
      toast.add({
        severity: 'success',
        summary: t('customer_updated'),
        detail: t('done'),
        group: 'br',
        life: 3000,
      })
      customerEditDialog.value = false
      customerToEdit.value = {}
      loadCustomers()
    })
    .catch((error) => {
      toast.add({
        severity: 'error',
        summary: t('failed'),
        detail: error.response?.data?.message || t('error_occurred'),
        group: 'br',
      })
    })
    .finally(() => {
      isSubmitting.value = false
    })
}

const submitCustomer = () => {
  customer_add_errors.value.name = new_customer_name.value?.trim() ? '' : t('validation_required')

  if (customer_add_errors.value.name) return

  isSubmitting.value = true

  const payload = {
    name: new_customer_name.value,
    phone: new_customer_phone.value,
    address: new_customer_address.value,
    email: new_customer_email.value,
    notes: new_customer_notes.value,
  }

  axios
    .post(
      `http://${import.meta.env.VITE_APP_BACKEND_HOST}${import.meta.env.VITE_APP_MODULE_CORE_API_PREFIX}/api/customers`,
      {
        data: payload,
      },
      {
        headers: {
          Authorization: `Bearer ${auth.accessToken.value}`,
        },
      },
    )
    .then((response) => {
      toast.add({
        severity: 'success',
        summary: t('customer_added'),
        detail: t('done'),
        group: 'br',
        life: 3000,
      })
      new_customer_id.value = response.data.data.id
      customerAddDialog.value = false
      resetAddForm()
      loadCustomers()
    })
    .catch((error) => {
      toast.add({
        severity: 'error',
        summary: t('failed'),
        detail: error.response?.data?.data || t('error_occurred'),
        group: 'br',
      })
    })
    .finally(() => {
      isSubmitting.value = false
    })
}

const resetAddForm = () => {
  new_customer_name.value = ''
  new_customer_phone.value = ''
  new_customer_address.value = ''
  new_customer_email.value = ''
  new_customer_notes.value = ''
}

const updatCustomersTableRowsPerPage = (event: DataTablePageEvent) => {
  loadCustomers(event.first, event.rows)
}

const loadCustomers = (first = 0, rows = customersTableRowsPerPage.value) => {
  isCustomersTableLoading.value = true

  if (first == 0) {
    first = 1
  }

  const page_number = Math.ceil(first / rows)

  axios
    .get(
      `http://${import.meta.env.VITE_APP_BACKEND_HOST}${import.meta.env.VITE_APP_MODULE_CORE_API_PREFIX}/api/customers`,
      {
        params: {
          'page[number]': page_number,
          'page[size]': rows,
        },
        headers: {
          Authorization: `Bearer ${auth.accessToken.value}`,
        },
      },
    )
    .then((response) => {
      customers.value = response.data.data
      customersTableTotalRecords.value = response.data.meta.total_records
    })
    .catch(() => {
      toast.add({ severity: 'error', summary: t('failed'), detail: t('customers_load_failed') })
    })
    .finally(() => {
      isCustomersTableLoading.value = false
    })
}

loadCustomers()
</script>

<template>
  <div class="w-full">
    <div class="grid mx-2">
      <div class="col-12 flex">
        <div class="grid w-full">
          <div class="col-12">
            <h3>{{ $t('branches', 3) }}</h3>
          </div>
          <div class="col-12">
            <DataTable
              @page="updateBranchesTableRowsPerPage"
              :lazy="true"
              :totalRecords="branchesTableTotalRecords"
              :loading="isBranchesTableLoading"
              paginatorPosition="both"
              paginator
              :rows="branchesTableRowsPerPage"
              :rowsPerPageOptions="[50, 100, 500]"
              :value="branches"
              stripedRows
              tableStyle="width: 100%;max-height:50vh;"
              class="w-full pr-2"
            >
              <template #header>
                <div class="flex justify-between align-items-center">
                  <Button
                    icon="pi pi-plus"
                    :label="$t('add_branch')"
                    rounded
                    raised
                    @click="openAddDialog"
                  />
                </div>
              </template>
              <template #empty>
                <div class="flex flex-column align-items-center gap-2 py-4">
                  <i class="pi pi-building" style="font-size: 2rem; opacity: 0.3"></i>
                  <p class="m-0 text-slate-400">{{ $t('no_results') }}</p>
                </div>
              </template>
              <Column sortable field="name" :header="$t('branch_name')"></Column>
              <Column sortable field="address" :header="$t('branch_address')"></Column>
              <Column sortable field="phone" :header="$t('branch_phone')"></Column>
              <Column sortable field="email" :header="$t('branch_email')"></Column>
              <Column sortable field="is_active" :header="$t('status')">
                <template #body="slotProps">
                  <Tag
                    :value="slotProps.data.is_active ? $t('active') : $t('inactive')"
                    :severity="slotProps.data.is_active ? 'success' : 'danger'"
                  />
                </template>
              </Column>
              <Column :header="$t('actions')">
                <template #body="slotProps">
                  <ConfirmPopup></ConfirmPopup>
                  <ButtonGroup>
                    <Button
                      icon="pi pi-pencil"
                      severity="secondary"
                      :aria-label="$t('edit')"
                      @click="prepareBranchToEdit(slotProps.data)"
                    />
                    <Button
                      icon="pi pi-trash"
                      severity="danger"
                      :aria-label="$t('remove')"
                      @click="confirmDeleteBranch($event, slotProps.data.id)"
                    />
                  </ButtonGroup>
                </template>
              </Column>
            </DataTable>
          </div>
        </div>
      </div>
    </div>

    <Dialog
      v-model:visible="branchAddDialog"
      modal
      :header="$t('add_branch')"
      :style="{ width: '30rem', direction: store.orientation == 'rtl' ? 'rtl' : 'ltr' }"
      :breakpoints="{ '1199px': '90vw', '575px': '90vw' }"
    >
      <div class="flex flex-column gap-4">
        <div class="grid">
          <div class="col-12 flex flex-column gap-2">
            <label for="new_branch_name">{{ $t('branch_name') }}</label>
            <InputText
              id="new_branch_name"
              v-model="newBranch.name"
              :class="{ 'p-invalid': newBranchErrors.name }"
            />
            <small class="p-error">{{ newBranchErrors.name }}</small>
          </div>
        </div>
        <div class="grid">
          <div class="col-12 flex flex-column gap-2">
            <label for="new_branch_address">{{ $t('branch_address') }}</label>
            <InputText id="new_branch_address" v-model="newBranch.address" />
          </div>
        </div>
        <div class="grid">
          <div class="col-12 md:col-6 flex flex-column gap-2">
            <label for="new_branch_phone">{{ $t('branch_phone') }}</label>
            <InputText id="new_branch_phone" v-model="newBranch.phone" />
          </div>
          <div class="col-12 md:col-6 flex flex-column gap-2">
            <label for="new_branch_email">{{ $t('branch_email') }}</label>
            <InputText id="new_branch_email" v-model="newBranch.email" />
          </div>
        </div>
        <div class="flex align-items-center gap-2">
          <Checkbox id="new_branch_active" v-model="newBranch.is_active" :binary="true" />
          <label for="new_branch_active">{{ $t('active') }}</label>
        </div>
      </div>
      <template #footer>
        <ButtonGroup>
          <Button :label="$t('cancel')" severity="secondary" @click="branchAddDialog = false" />
          <Button
            class="ml-2"
            severity="primary"
            @click="submitBranch"
            :label="$t('save')"
            :loading="isSubmitting"
            :disabled="isSubmitting"
          />
        </ButtonGroup>
      </template>
    </Dialog>

    <Dialog
      v-model:visible="branchEditDialog"
      modal
      :header="$t('edit_branch')"
      :style="{ width: '30rem', direction: store.orientation == 'rtl' ? 'rtl' : 'ltr' }"
      :breakpoints="{ '1199px': '90vw', '575px': '90vw' }"
    >
      <div class="flex flex-column gap-4">
        <div class="grid">
          <div class="col-12 flex flex-column gap-2">
            <label for="edit_branch_name">{{ $t('branch_name') }}</label>
            <InputText
              id="edit_branch_name"
              v-model="branchToEdit.name"
              :class="{ 'p-invalid': editBranchErrors.name }"
            />
            <small class="p-error">{{ editBranchErrors.name }}</small>
          </div>
        </div>
        <div class="grid">
          <div class="col-12 flex flex-column gap-2">
            <label for="edit_branch_address">{{ $t('branch_address') }}</label>
            <InputText id="edit_branch_address" v-model="branchToEdit.address" />
          </div>
        </div>
        <div class="grid">
          <div class="col-12 md:col-6 flex flex-column gap-2">
            <label for="edit_branch_phone">{{ $t('branch_phone') }}</label>
            <InputText id="edit_branch_phone" v-model="branchToEdit.phone" />
          </div>
          <div class="col-12 md:col-6 flex flex-column gap-2">
            <label for="edit_branch_email">{{ $t('branch_email') }}</label>
            <InputText id="edit_branch_email" v-model="branchToEdit.email" />
          </div>
        </div>
        <div class="flex align-items-center gap-2">
          <Checkbox id="edit_branch_active" v-model="branchToEdit.is_active" :binary="true" />
          <label for="edit_branch_active">{{ $t('active') }}</label>
        </div>
      </div>
      <template #footer>
        <ButtonGroup>
          <Button :label="$t('cancel')" severity="secondary" @click="branchEditDialog = false" />
          <Button
            class="ml-2"
            severity="primary"
            @click="updateBranch"
            :label="$t('save')"
            :loading="isSubmitting"
            :disabled="isSubmitting"
          />
        </ButtonGroup>
      </template>
    </Dialog>
  </div>
</template>

<script setup lang="ts">
import ConfirmPopup from 'primevue/confirmpopup'
import DataTable from 'primevue/datatable'
import Dialog from 'primevue/dialog'
import InputText from 'primevue/inputtext'
import Checkbox from 'primevue/checkbox'
import Column from 'primevue/column'
import Button from 'primevue/button'
import ButtonGroup from 'primevue/buttongroup'
import Tag from 'primevue/tag'
import axios from 'axios'
import { ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { useToast } from 'primevue/usetoast'
import { useConfirm } from 'primevue/useconfirm'
import { globalStore } from '@/stores'

const { t } = useI18n()
const store = globalStore()
const confirm = useConfirm()
const toast = useToast()

interface BranchItem {
  id: string
  name: string
  address: string
  phone: string
  email: string
  is_active: boolean
}

const branches = ref<BranchItem[]>([])
const branchesTableTotalRecords = ref(0)
const isBranchesTableLoading = ref(false)
const branchesTableRowsPerPage = ref(50)
const isSubmitting = ref(false)
const branchAddDialog = ref(false)
const branchEditDialog = ref(false)

const newBranch = ref<Partial<BranchItem>>({
  name: '',
  address: '',
  phone: '',
  email: '',
  is_active: true,
})
const newBranchErrors = ref({ name: '' })

const branchToEdit = ref<Partial<BranchItem>>({})
const editBranchErrors = ref({ name: '' })

const openAddDialog = () => {
  newBranch.value = { name: '', address: '', phone: '', email: '', is_active: true }
  newBranchErrors.value = { name: '' }
  branchAddDialog.value = true
}

const prepareBranchToEdit = (branch: BranchItem) => {
  branchToEdit.value = JSON.parse(JSON.stringify(branch))
  editBranchErrors.value = { name: '' }
  branchEditDialog.value = true
}

const updateBranchesTableRowsPerPage = (event: { first: number; rows: number }) => {
  getBranches(event.first, event.rows)
}

const submitBranch = () => {
  newBranchErrors.value.name = newBranch.value.name?.trim() ? '' : t('validation_required')
  if (newBranchErrors.value.name) return

  isSubmitting.value = true
  axios
    .post(
      `http://${import.meta.env.VITE_APP_BACKEND_HOST}/branch/api/branches`,
      newBranch.value,
    )
    .then(() => {
      toast.add({
        severity: 'success',
        summary: t('success'),
        detail: t('branch_added'),
        life: 3000,
      })
      branchAddDialog.value = false
      getBranches(0, branchesTableRowsPerPage.value)
    })
    .catch((err) => {
      toast.add({ severity: 'error', summary: t('error'), detail: err.message, life: 3000 })
    })
    .finally(() => {
      isSubmitting.value = false
    })
}

const updateBranch = () => {
  editBranchErrors.value.name = branchToEdit.value.name?.trim() ? '' : t('validation_required')
  if (editBranchErrors.value.name) return

  isSubmitting.value = true
  axios
    .patch(
      `http://${import.meta.env.VITE_APP_BACKEND_HOST}/branch/api/branches/${branchToEdit.value.id}`,
      branchToEdit.value,
    )
    .then(() => {
      toast.add({
        severity: 'success',
        summary: t('success'),
        detail: t('branch_updated'),
        life: 3000,
      })
      branchEditDialog.value = false
      getBranches(0, branchesTableRowsPerPage.value)
    })
    .catch((err) => {
      toast.add({ severity: 'error', summary: t('error'), detail: err.message, life: 3000 })
    })
    .finally(() => {
      isSubmitting.value = false
    })
}

const confirmDeleteBranch = (event: MouseEvent, id: string) => {
  confirm.require({
    target: event.target,
    message: t('do_you_confirm'),
    icon: 'pi pi-exclamation-triangle',
    rejectClass: 'p-button-secondary p-button-outlined',
    rejectLabel: t('cancel'),
    acceptLabel: t('delete'),
    acceptClass: 'p-button-danger',
    accept: () => {
      axios
        .delete(
          `http://${import.meta.env.VITE_APP_BACKEND_HOST}/branch/api/branches/${id}`,
        )
        .then(() => {
          toast.add({
            severity: 'success',
            summary: t('success'),
            detail: t('branch_deleted'),
            life: 3000,
          })
          getBranches(0, branchesTableRowsPerPage.value)
        })
        .catch((err) => {
          toast.add({ severity: 'error', summary: t('error'), detail: err.message, life: 3000 })
        })
    },
  })
}

const getBranches = (offset = 0, limit = 50) => {
  isBranchesTableLoading.value = true
  const pageNumber = Math.floor(offset / limit) + 1
  axios
    .get(
      `http://${import.meta.env.VITE_APP_BACKEND_HOST}/branch/api/branches?page_number=${pageNumber}&page_size=${limit}`,
    )
    .then((res) => {
      branches.value = res.data.data || []
      branchesTableTotalRecords.value = res.data.meta?.total_records || 0
    })
    .catch((err) => {
      toast.add({ severity: 'error', summary: t('error'), detail: err.message, life: 3000 })
    })
    .finally(() => {
      isBranchesTableLoading.value = false
    })
}

getBranches()
</script>

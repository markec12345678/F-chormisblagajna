<template>
  <div class="w-full">
    <div class="grid mx-2">
      <div class="col-12 flex">
        <div class="grid w-full">
          <div class="col-12">
            <h3>{{ $t('user', 3) }}</h3>
          </div>
          <div class="col-12">
            <DataTable
              :loading="isUsersLoading"
              :value="users"
              stripedRows
              tableStyle="min-width: 50rem;max-height:50vh;"
              class="w-full pr-2"
            >
              <template #header>
                <div class="flex justify-start">
                  <Button
                    icon="pi pi-plus"
                    :label="$t('add_user')"
                    rounded
                    raised
                    @click="userAddDialog = true"
                  />
                </div>
              </template>
              <template #empty>
                <div class="flex flex-column align-items-center gap-2 py-4">
                  <i class="pi pi-users" style="font-size: 2rem; opacity: 0.3"></i>
                  <p class="m-0" style="color: #94a3b8">{{ $t('no_results') }}</p>
                </div>
              </template>
              <Column sortable field="username" :header="$t('username')"></Column>
              <Column sortable field="email" :header="$t('email')"></Column>
              <Column field="roles" :header="$t('role', 3)">
                <template #body="slotProps">
                  <Chip
                    v-for="(role, index) in slotProps.data.roles"
                    :key="index"
                    :label="role"
                    style="height: 1.5rem"
                    class="m-1"
                  />
                </template>
              </Column>
              <Column :header="$t('actions')">
                <template #body="slotProps">
                  <ConfirmPopup></ConfirmPopup>
                  <ButtonGroup>
                    <Button
                      v-if="isSuperuser && !slotProps.data.roles.includes('superuser')"
                      icon="pi pi-pencil"
                      severity="secondary"
                      aria-label="Edit Password"
                      @click="openPasswordEditDialog(slotProps.data)"
                    />
                    <Button
                      v-if="!slotProps.data.roles.includes('superuser')"
                      icon="pi pi-trash"
                      severity="danger"
                      aria-label="Remove"
                      @click="confirmDeleteUser($event, slotProps.data.id)"
                    />
                  </ButtonGroup>
                </template>
              </Column>
            </DataTable>
            <Dialog
              v-model:visible="userAddDialog"
              modal
              :header="$t('add_user')"
              :style="{ width: '75rem' }"
              :breakpoints="{ '1199px': '90vw', '575px': '90vw' }"
            >
              <div class="flex flex-column gap-2 w-5">
                <label for="username">{{ $t('username') }}</label>
                <InputText
                  id="username"
                  v-model="newUser.username"
                  aria-describedby="username"
                  :class="{ 'p-invalid': user_errors.username }"
                />
                <small class="p-error">{{ user_errors.username }}</small>
              </div>
              <div class="flex flex-column gap-2 w-5 mt-2">
                <label for="email">{{ $t('email') }}</label>
                <InputText
                  id="email"
                  v-model="newUser.email"
                  aria-describedby="email"
                  :class="{ 'p-invalid': user_errors.email }"
                />
                <small class="p-error">{{ user_errors.email }}</small>
              </div>
              <div class="flex flex-column gap-2 w-5 mt-2">
                <label for="password">{{ $t('password') }}</label>
                <InputText
                  id="password"
                  v-model="newUser.password"
                  type="password"
                  aria-describedby="password"
                  :class="{ 'p-invalid': user_errors.password }"
                />
                <small class="p-error">{{ user_errors.password }}</small>
              </div>
              <div class="flex flex-column gap-2 w-10 mt-3">
                <label for="roles">{{ $t('role', 3) }}</label>
                <div class="flex flex-wrap gap-2">
                  <div v-for="role in availableRoles" :key="role" class="flex align-items-center">
                    <Checkbox v-model="newUser.roles" :inputId="role" :value="role" />
                    <label :for="role" class="ml-2">{{ role }}</label>
                  </div>
                </div>
              </div>
              <template #footer>
                <ButtonGroup>
                  <Button
                    :label="$t('cancel')"
                    severity="secondary"
                    aria-label="Cancel"
                    @click="userAddDialog = false"
                  />
                  <Button
                    class="ml-2"
                    severity="primary"
                    @click="submitUser"
                    :label="$t('save')"
                    aria-label="Save"
                  />
                </ButtonGroup>
              </template>
            </Dialog>
            <Dialog
              v-model:visible="passwordEditDialog"
              modal
              :header="$t('change_password')"
              :style="{ width: '25rem' }"
            >
              <div class="flex flex-column gap-2">
                <label for="editUsername">{{ $t('username') }}</label>
                <InputText id="editUsername" v-model="editPassword.username" disabled />
              </div>
              <div class="flex flex-column gap-2 mt-2">
                <label for="newPassword">{{ $t('new_password') }}</label>
                <InputText id="newPassword" v-model="editPassword.password" type="password" />
              </div>
              <template #footer>
                <ButtonGroup>
                  <Button
                    :label="$t('cancel')"
                    severity="secondary"
                    aria-label="Cancel"
                    @click="passwordEditDialog = false"
                  />
                  <Button
                    class="ml-2"
                    severity="primary"
                    @click="submitPasswordChange"
                    :label="$t('save')"
                    aria-label="Save"
                  />
                </ButtonGroup>
              </template>
            </Dialog>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useI18n } from 'vue-i18n'
import axios from 'axios'
import { getCurrentInstance } from 'vue'
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import Button from 'primevue/button'
import ButtonGroup from 'primevue/buttongroup'
import Dialog from 'primevue/dialog'
import InputText from 'primevue/inputtext'
import Chip from 'primevue/chip'
import Checkbox from 'primevue/checkbox'
import ConfirmPopup from 'primevue/confirmpopup'
import { useConfirm } from 'primevue/useconfirm'
import { useToast } from 'primevue/usetoast'

const { proxy } = getCurrentInstance()
const confirm = useConfirm()
const toast = useToast()
const { t } = useI18n({ useScope: 'global' })

const backendUrl = `http://${import.meta.env.VITE_APP_BACKEND_HOST}${import.meta.env.VITE_APP_MODULE_CORE_API_PREFIX}`

const users = ref([])
const isUsersLoading = ref(true)
const userAddDialog = ref(false)
const passwordEditDialog = ref(false)
const availableRoles = ['admin', 'cashier', 'chef']
const newUser = ref({
  username: '',
  email: '',
  password: '',
  roles: [],
})
const user_errors = ref({ username: '', email: '', password: '' })
const editPassword = ref({
  userId: '',
  username: '',
  password: '',
})

const auth = proxy.$auth
const isSuperuser = computed(() => auth.currentUser?.value.roles?.includes('superuser') || false)

const loadUsers = () => {
  isUsersLoading.value = true
  axios
    .get(`${backendUrl}/api/auth/users`, {
      headers: {
        Authorization: `Bearer ${auth.accessToken.value}`,
      },
    })
    .then((response) => {
      users.value = response.data
    })
    .catch(() => {
      toast.add({ severity: 'error', summary: 'Error', detail: t('users_load_failed') })
    })
    .finally(() => {
      isUsersLoading.value = false
    })
}

const submitUser = () => {
  user_errors.value.username = newUser.value.username?.trim() ? '' : proxy.$t('validation_required')
  user_errors.value.email = newUser.value.email?.trim() ? '' : proxy.$t('validation_required')
  user_errors.value.password = newUser.value.password?.trim() ? '' : proxy.$t('validation_required')

  if (user_errors.value.username || user_errors.value.email || user_errors.value.password) return

  axios
    .post(`${backendUrl}/api/auth/users`, newUser.value, {
      headers: {
        Authorization: `Bearer ${auth.accessToken.value}`,
      },
    })
    .then(() => {
      toast.add({ severity: 'success', summary: t('user_added'), detail: t('done'), group: 'br' })
      userAddDialog.value = false
      newUser.value = {
        username: '',
        email: '',
        password: '',
        roles: [],
      }
      loadUsers()
    })
    .catch((error) => {
      toast.add({
        severity: 'error',
        summary: 'Error',
        detail: error.response?.data?.message || t('error_occurred'),
        group: 'br',
      })
    })
}

const deleteUser = (user_id: string) => {
  axios
    .delete(`${backendUrl}/api/auth/users?id=${user_id}`, {
      headers: {
        Authorization: `Bearer ${auth.accessToken.value}`,
      },
    })
    .then(() => {
      toast.add({
        severity: 'success',
        summary: 'Success',
        detail: t('user_deleted_success'),
        group: 'br',
        life: 3000,
      })
      loadUsers()
    })
    .catch(() => {
      toast.add({
        severity: 'error',
        summary: 'Error',
        detail: t('user_delete_failed'),
        group: 'br',
        life: 3000,
      })
    })
}

const confirmDeleteUser = (event, user_id) => {
  confirm.require({
    target: event.currentTarget,
    message: t('confirm_delete_user'),
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
      deleteUser(user_id)
    },
    reject: () => {},
  })
}

const openPasswordEditDialog = (user) => {
  editPassword.value = {
    userId: user.id,
    username: user.username,
    password: '',
  }
  passwordEditDialog.value = true
}

const submitPasswordChange = () => {
  if (!editPassword.value.password) {
    toast.add({ severity: 'warn', summary: 'Warning', detail: t('enter_new_password') })
    return
  }

  axios
    .post(
      `${backendUrl}/api/auth/users/password`,
      {
        user_id: editPassword.value.userId,
        password: editPassword.value.password,
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
        summary: 'Success',
        detail: t('password_changed_success'),
        group: 'br',
        life: 3000,
      })
      passwordEditDialog.value = false
      editPassword.value.password = ''
    })
    .catch((error) => {
      toast.add({
        severity: 'error',
        summary: 'Error',
        detail: error.response?.data?.message || t('password_change_failed'),
        group: 'br',
        life: 3000,
      })
    })
}

onMounted(() => {
  loadUsers()
})
</script>

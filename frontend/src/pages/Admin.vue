<template>
  <div v-if="!loading">
    <div class="grid p-0 m-0">
      <div class="col-12 p-0">
        <Toolbar style="border-radius: 0px" class="py-1 lg:py-2">
          <template #start>
            <div @click="version_dialog_visible = true" class="no-underline text-gray-400">
              <img
                src="@/assets/logo.png"
                alt="logo"
                style="height: 25px"
                v-if="store.getColorMode == 'light'"
              />
              <div
                v-else
                class="flex justify-content-center align-items-center"
                style="font-size: 1rem; color: white; font-family: 'FontAwesome'"
              >
                nutrix
              </div>
            </div>
            <router-link v-for="(item, index) in items" :key="index" :to="item.link">
              <Button
                :icon="item.icon"
                :label="$t(`${item.label.title}`, item.label.plural ? 3 : 1)"
                :text="!item.focused"
                severity="secondary"
              />
            </router-link>
          </template>

          <template #end>
            <Button
              outlined
              :icon="`pi pi-${store.getColorMode == 'light' ? 'sun' : 'moon'}`"
              @click="toggleDarkMode()"
            />
            <Button
              severity="secondary"
              size="large"
              text
              rounded
              :aria-label="$t('profile')"
              :label="$t('profile')"
              @click="$router.push('/profile')"
            >
              <span style="font-size: 0.9rem" class="mr-2">{{ user?.username }}</span>
              <span class="p-button-icon pi pi-user"></span>
            </Button>
          </template>
        </Toolbar>
      </div>
      <div class="col-12">
        <div class="grid">
          <div class="col-3 xl:col-2">
            <Tree
              v-model:expandedKeys="expandedKeys"
              v-model:selectionKeys="selectionKeys"
              :value="menu_tree"
              selectionMode="single"
              class="w-full"
              @node-select="(node) => sidemenuNodeSelect(node)"
            >
              <template #default="slotProps">
                <div class="no-underline flex align-items-center w-full">
                  <div>
                    {{ $t(`${slotProps.node.label.title}`, slotProps.node.label.plural ? 3 : 1) }}
                  </div>
                </div>
              </template>
            </Tree>
          </div>
          <div class="col-9 xl:col-10 flex p-0 pt-3 mt-2">
            <RouterView />
          </div>
        </div>
      </div>
    </div>
  </div>
  <div
    style="width: 100vw; height: 100vh; display: flex; justify-content: center; align-items: center"
    v-if="loading"
  >
    <ProgressSpinner
      style="width: 35px; height: 35px"
      strokeWidth="6"
      fill="transparent"
      animationDuration=".5s"
      :aria-label="$t('loading')"
    />
  </div>
  <Dialog v-model:visible="version_dialog_visible" header="Nutrix" :style="{ width: '45rem' }">
    <p class="text-justify">
      {{ $t('about_nutrix') }}
    </p>
    <p>
      {{ $t('for_more_support') }} &nbsp;<a
        style="font-size: large"
        href="https://nutrixpos.com"
        target="_blank"
        ><i class="pi pi-external-link mr-2"></i>https://nutrixpos.com
      </a>
    </p>
    <p>{{ $t('version_info') }} {{ app_version }}</p>
  </Dialog>
</template>

<script setup lang="ts">
import { ref, computed, getCurrentInstance, watch } from 'vue'
import { Toolbar, Dialog } from 'primevue'
import Tree from 'primevue/tree'
import Button from 'primevue/button'
import { useI18n } from 'vue-i18n'
import { useToast } from 'primevue/usetoast'
import { globalStore } from '@/stores'
import axios from 'axios'
import ProgressSpinner from 'primevue/progressspinner'
import { useRoute } from 'vue-router'
const { proxy } = getCurrentInstance()
const route = useRoute()

const store = globalStore()

const app_version = ref('')
const version_dialog_visible = ref(false)

app_version.value = import.meta.env.VITE_APP_APP_VERSION || ''

interface MenuNode {
  key: string
  link?: string
  children?: MenuNode[]
}

const user = computed(() => {
  return proxy.$auth.currentUser.value
})

const toggleDarkMode = () => {
  store.toggleDarkMode()
}

const sidemenuNodeSelect = (node) => {
  if (node.link) {
    proxy.$router.push(node.link)
  }
}

// const selected_list_item = ref ({ name: 'Inventory', icon:'inbox', link:'inventory' })

const expandAll = () => {
  for (const node of menu_tree.value) {
    expandNode(node)
  }

  expandedKeys.value = { ...expandedKeys.value }
}

const expandNode = (node) => {
  if (node.children && node.children.length) {
    expandedKeys.value[node.key] = true

    for (const child of node.children) {
      expandNode(child)
    }
  }
}

const expandedKeys = ref({})
const selectionKeys = ref({})

const findKeyByLink = (nodes: MenuNode[], link: string): string | null => {
  for (const node of nodes) {
    if (node.link === link) return node.key
    if (node.children) {
      const found = findKeyByLink(node.children, link)
      if (found) return found
    }
  }
  return null
}

watch(
  () => route.path,
  (path) => {
    const key = findKeyByLink(menu_tree.value, path)
    if (key) {
      selectionKeys.value = { [key]: true }
    }
  },
  { immediate: true },
)

const menu_tree = ref([
  {
    key: '-1',
    label: {
      title: 'dashboard',
      plural: false,
    },
    data: 'Dashboard',
    icon: 'pi pi-fw pi-home',
    link: '/admin/dashboard',
  },
  {
    key: '0',
    label: {
      title: 'inventory',
      plural: true,
    },
    data: 'Inventory',
    icon: 'fa fa-boxes-stacked',
    link: '/admin/inventory',
  },
  {
    key: '1',
    label: {
      title: 'product',
      plural: true,
    },
    data: 'Products',
    icon: 'fa fa-barcode',
    link: '/admin/products',
  },
  {
    key: '2',
    label: {
      title: 'category',
      plural: true,
    },
    data: 'Categories',
    icon: 'fa fa-sitemap',
    link: '/admin/categories',
  },
  {
    key: '3',
    label: {
      title: 'order',
      plural: true,
    },
    data: 'Orders',
    icon: 'pi pi-fw pi-box',
    link: '/admin/orders',
    children: [
      {
        key: '3-0',
        label: {
          title: 'list',
          plural: false,
        },
        data: 'List orders',
        icon: 'pi pi-fw pi-list',
        link: '/admin/orders',
      },
    ],
  },
  {
    key: '4',
    label: {
      title: 'report',
      plural: true,
    },
    data: 'Reports',
    icon: 'pi pi-fw pi-chart-line',
    link: '/admin/sales',
    children: [
      {
        key: '4-0',
        label: {
          title: 'sales',
          plural: false,
        },
        data: 'Sales',
        icon: 'pi pi-fw pi-percentage',
        link: '/admin/sales',
      },
    ],
  },
  {
    key: '5',
    label: {
      title: 'settings',
      plural: false,
    },
    data: 'Settings',
    icon: 'pi pi-fw pi-cog',
    link: '/admin/settings',
  },
  {
    key: '6',
    label: {
      title: 'customer',
      plural: true,
    },
    data: 'Customers',
    icon: 'pi pi-users',
    link: '/admin/customers',
  },
  {
    key: '7',
    label: {
      title: 'user',
      plural: true,
    },
    data: 'Users',
    icon: 'pi pi-user',
    link: '/admin/users',
  },
  {
    key: '8',
    label: {
      title: 'Hubsync',
      plural: false,
    },
    data: 'Hubsync',
    icon: 'pi pi-sync',
    link: '/admin/hubsync',
  },
  {
    key: '9',
    label: {
      title: 'tables',
      plural: true,
    },
    data: 'Tables',
    icon: 'pi pi-table',
    link: '/admin/tables',
  },
  {
    key: '10',
    label: {
      title: 'branches',
      plural: true,
    },
    data: 'Branches',
    icon: 'pi pi-building',
    link: '/admin/branches',
  },
  {
    key: '11',
    label: {
      title: 'split_bill',
      plural: false,
    },
    data: 'Split Bill',
    icon: 'pi pi-wallet',
    link: '/admin/split-bills',
  },
  {
    key: '12',
    label: {
      title: 'scheduling',
      plural: false,
    },
    data: 'Scheduling',
    icon: 'pi pi-calendar',
    link: '/admin/scheduling',
  },
  {
    key: '13',
    label: {
      title: 'reservations',
      plural: true,
    },
    data: 'Reservations',
    icon: 'pi pi-book',
    link: '/admin/reservations',
  },
  {
    key: '14',
    label: {
      title: 'promotions',
      plural: true,
    },
    data: 'Promotions',
    icon: 'pi pi-tags',
    link: '/admin/promotions',
  },
  {
    key: '15',
    label: {
      title: 'loyalty_program',
      plural: false,
    },
    data: 'Loyalty',
    icon: 'pi pi-star',
    link: '/admin/loyalty',
  },
  {
    key: '16',
    label: {
      title: 'reports',
      plural: true,
    },
    data: 'Reports',
    icon: 'pi pi-chart-line',
    link: '/admin/reports',
  },
  {
    key: '17',
    label: {
      title: 'multi_location',
      plural: false,
    },
    data: 'Locations',
    icon: 'pi pi-globe',
    link: '/admin/locations',
  },
  {
    key: '18',
    label: {
      title: 'ai_search',
      plural: false,
    },
    data: 'AI',
    icon: 'pi pi-sparkles',
    link: '/admin/ai',
  },
  {
    key: '19',
    label: {
      title: 'gift_cards',
      plural: true,
    },
    data: 'Gift Cards',
    icon: 'pi pi-gift',
    link: '/admin/giftcards',
  },
  {
    key: '20',
    label: {
      title: 'inventory_transfers',
      plural: true,
    },
    data: 'Transfers',
    icon: 'pi pi-arrow-right-arrow-left',
    link: '/admin/transfers',
  },
  {
    key: '21',
    label: {
      title: 'fiscal_dashboard',
      plural: false,
    },
    data: 'Fiscal',
    icon: 'pi pi-receipt',
    link: '/admin/fiscal',
  },
  {
    key: '22',
    label: {
      title: 'floor_plan',
      plural: false,
    },
    data: 'Floor Plan',
    icon: 'pi pi-th-large',
    link: '/admin/floorplan',
  },
  {
    key: '23',
    label: {
      title: 'employee_tips',
      plural: true,
    },
    data: 'Tips',
    icon: 'pi pi-dollar',
    link: '/admin/tips',
  },
  {
    key: '24',
    label: {
      title: 'accounting_export',
      plural: false,
    },
    data: 'Accounting Export',
    icon: 'pi pi-file-excel',
    link: '/admin/accounting',
  },
  {
    key: '25',
    label: {
      title: 'menu_engineering',
      plural: false,
    },
    data: 'Menu Engineering',
    icon: 'pi pi-chart-bar',
    link: '/admin/menu-engineering',
  },
  {
    key: '26',
    label: {
      title: 'employee_performance',
      plural: false,
    },
    data: 'Employee Performance',
    icon: 'pi pi-users',
    link: '/admin/employee-performance',
  },
  {
    key: '27',
    label: {
      title: 'suppliers',
      plural: true,
    },
    data: 'Suppliers',
    icon: 'pi pi-truck',
    link: '/admin/suppliers',
  },
  {
    key: '28',
    label: {
      title: 'expense_tracking',
      plural: false,
    },
    data: 'Expenses',
    icon: 'pi pi-wallet',
    link: '/admin/expenses',
  },
  {
    key: '29',
    label: {
      title: 'time_clock',
      plural: false,
    },
    data: 'Time Clock',
    icon: 'pi pi-clock',
    link: '/admin/timeclock',
  },
  {
    key: '30',
    label: {
      title: 'waste_tracking',
      plural: false,
    },
    data: 'Waste Tracking',
    icon: 'pi pi-trash',
    link: '/admin/waste',
  },
  {
    key: '31',
    label: {
      title: 'online_orders',
      plural: false,
    },
    data: 'Online Orders',
    icon: 'pi pi-shopping-cart',
    link: '/admin/online-orders',
  },
  {
    key: '32',
    label: {
      title: 'staff_chat',
      plural: false,
    },
    data: 'Staff Chat',
    icon: 'pi pi-comments',
    link: '/admin/chat',
  },
  {
    key: '33',
    label: {
      title: 'audit_log',
      plural: false,
    },
    data: 'Audit Log',
    icon: 'pi pi-history',
    link: '/admin/audit-log',
  },
  {
    key: '34',
    label: {
      title: 'receipt_customization',
      plural: false,
    },
    data: 'Receipt',
    icon: 'pi pi-file',
    link: '/admin/receipt',
  },
  {
    key: '35',
    label: {
      title: 'inventory_alerts',
      plural: false,
    },
    data: 'Alerts',
    icon: 'pi pi-bell',
    link: '/admin/inventory-alerts',
  },
  {
    key: '36',
    label: {
      title: 'multi_payment',
      plural: false,
    },
    data: 'MultiPay',
    icon: 'pi pi-credit-card',
    link: '/admin/multi-payment',
  },
  {
    key: '37',
    label: {
      title: 'customer_display',
      plural: false,
    },
    data: 'Display',
    icon: 'pi pi-desktop',
    link: '/admin/customer-display',
  },
  {
    key: '38',
    label: {
      title: 'tableside_ordering',
      plural: false,
    },
    data: 'Tableside',
    icon: 'pi pi-table',
    link: '/admin/tableside',
  },
  {
    key: '39',
    label: {
      title: 'staff_training',
      plural: false,
    },
    data: 'Training',
    icon: 'pi pi-graduation-cap',
    link: '/admin/staff-training',
  },
  {
    key: '40',
    label: {
      title: 'customer_feedback',
      plural: false,
    },
    data: 'Feedback',
    icon: 'pi pi-comment',
    link: '/admin/customer-feedback',
  },
  {
    key: '41',
    label: {
      title: 'delivery_management',
      plural: false,
    },
    data: 'Delivery',
    icon: 'pi pi-truck',
    link: '/admin/delivery',
  },
  {
    key: '42',
    label: {
      title: 'self_service_kiosk',
      plural: false,
    },
    data: 'Kiosk',
    icon: 'pi pi-desktop',
    link: '/admin/kiosk',
  },
  {
    key: '43',
    label: {
      title: 'language',
      plural: true,
    },
    data: 'Languages',
    icon: 'pi pi-language',
    link: '/admin/languages',
  },
  {
    key: '44',
    label: {
      title: 'marketing_campaigns',
      plural: false,
    },
    data: 'Marketing',
    icon: 'pi pi-megaphone',
    link: '/admin/marketing',
  },
  {
    key: '45',
    label: {
      title: 'purchase_orders',
      plural: true,
    },
    data: 'Purchase Orders',
    icon: 'pi pi-shopping-bag',
    link: '/admin/purchase-orders',
  },
  {
    key: '46',
    label: {
      title: 'queue_waitlist',
      plural: false,
    },
    data: 'Queue',
    icon: 'pi pi-users',
    link: '/admin/queue',
  },
])

const items = ref([
  {
    label: {
      title: 'cashier',
      plural: false,
    },
    focused: false,
    icon: 'pi pi-desktop',
    link: '/home',
  },
  {
    label: {
      title: 'kitchen',
      plural: false,
    },
    focused: false,
    icon: 'fa fa-kitchen-set',
    link: '/kitchen',
  },
  {
    label: {
      title: 'admin',
      plural: false,
    },
    focused: true,
    icon: 'pi pi-cog',
    link: '/admin',
  },
])

const loading = ref(true)
const { locale, setLocaleMessage } = useI18n({ useScope: 'global' })
const toast = useToast()

const loadLanguage = async () => {
  await axios
    .get(
      `http://${import.meta.env.VITE_APP_BACKEND_HOST}${import.meta.env.VITE_APP_MODULE_CORE_API_PREFIX}/api/settings`,
      {
        headers: {
          Authorization: `Bearer ${proxy.$auth.accessToken.value}`,
        },
      },
    )
    .then(async (response) => {
      await axios
        .get(
          `http://${import.meta.env.VITE_APP_BACKEND_HOST}${import.meta.env.VITE_APP_MODULE_CORE_API_PREFIX}/api/languages/${response.data.data.language.code}`,
          {
            headers: {
              Authorization: `Bearer ${proxy.$auth.accessToken.value}`,
            },
          },
        )
        .then((response2) => {
          setLocaleMessage(response2.data.data.code, response2.data.data.pack)
          locale.value = response2.data.data.code
          store.setOrientation(response2.data.data.orientation)
          loading.value = false

          if (store.getShopMode != 'kitchen') {
            items.value = [
              {
                label: {
                  title: 'cashier',
                  plural: false,
                },
                focused: false,
                icon: 'pi pi-desktop',
                link: '/home',
              },
              {
                label: {
                  title: 'admin',
                  plural: false,
                },
                focused: true,
                icon: 'pi pi-cog',
                link: '/admin',
              },
            ]
          }
        })
        .catch(() => {
          toast.add({
            severity: 'error',
            summary: t('failed'),
            detail: t('request_failed'),
            life: 3000,
            group: 'br',
          })
        })
      loading.value = false
    })
    .catch((err) => {
      if (err.response?.status === 401) {
        proxy.$auth.signOut()
        window.location.href = '/'
      }
    })
}

loadLanguage()
expandAll()
</script>

<style>
html,
body {
  height: 100%;
  margin: 0;
}
</style>

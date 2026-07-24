<template>
  <div v-if="!loading" style="overflow-x: hidden">
    <div class="grid">
      <div class="col-12 pt-5 pl-3">
        <div
          v-if="orders.length <= 0"
          style="height: 100vh; width: 100%"
          class="flex flex-column justify-content-center align-items-center gap-3"
        >
          <i class="pi pi-inbox" style="font-size: 3rem; opacity: 0.3"></i>
          <h2 style="color: cadetblue; margin: 0">{{ $t('order', 3) }}: 0</h2>
          <p style="color: #94a3b8; margin: 0">{{ t('new_orders_appear') }}</p>
        </div>
        <div class="flex flex-wrap">
          <div v-for="(column_orders, index) in dynamic_columns" :key="index">
            <div class="flex flex-column gap-1">
              <QueueOrder
                @finished="orderFinished(order)"
                @openedDialog="openedDialogs++"
                @closedDialog="openedDialogs--"
                v-for="(order, i) in column_orders"
                :key="i"
                :order="order"
                :number="i + 1"
                class="queue-order"
              />
            </div>
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
      aria-label="Custom ProgressSpinner"
    />
  </div>
</template>

<script setup lang="ts">
import QueueOrder from '@/components/QueueOrder.vue'
import axios from 'axios'
import { ref, watch } from 'vue'
import Order from '@/classes/Order'
import { useToast } from 'primevue/usetoast'
import { Notification } from '@/classes/Notification'
import { globalStore } from '@/stores'
import { useI18n } from 'vue-i18n'
import ProgressSpinner from 'primevue/progressspinner'
import auth from '../services/auth'

const store = globalStore()

const toast = useToast()

const notifications = ref<Notification[]>([])

const orders = ref([])
const openedDialogs = ref(0)

interface Chat {
  type: string
  message: string
  topic_name: string
}

const chats = ref<Chat[]>([])

let socket: WebSocket

const loading = ref(true)
const { t, locale, setLocaleMessage } = useI18n({ useScope: 'global' })

const dynamic_columns = ref<Order[][]>([])
const prepareLayout = () => {
  const screenWidth = window.innerWidth
  dynamic_columns.value = []
  for (let i = 0; i < parseInt(`${screenWidth / 16 / 20}rem`); i++) {
    dynamic_columns.value.push([])
  }
}

watch(
  () => orders.value,
  () => {
    displayOrders()
  },
  {
    deep: true,
  },
)

const loadLanguage = async () => {
  await axios
    .get(
      `http://${import.meta.env.VITE_APP_BACKEND_HOST}${import.meta.env.VITE_APP_MODULE_CORE_API_PREFIX}/api/settings`,
      {
        headers: {
          Authorization: `Bearer ${auth.accessToken.value}`,
        },
      },
    )
    .then(async (response) => {
      await axios
        .get(
          `http://${import.meta.env.VITE_APP_BACKEND_HOST}${import.meta.env.VITE_APP_MODULE_CORE_API_PREFIX}/api/languages/${response.data.data.language.code}`,
          {
            headers: {
              Authorization: `Bearer ${auth.accessToken.value}`,
            },
          },
        )
        .then((response2) => {
          setLocaleMessage(response2.data.data.code, response2.data.data.pack)
          locale.value = response2.data.data.code
          store.setOrientation(response2.data.data.orientation)
          loading.value = false
        })
        .catch(() => {})
      loading.value = false
    })
    .catch((err) => {
      if (err.response?.status === 401) {
        auth.signOut()
        window.location.href = '/'
      }
    })
}

const startWebsocket = () => {
  socket = new WebSocket(
    `ws://${import.meta.env.VITE_APP_BACKEND_HOST}${import.meta.env.VITE_APP_MODULE_CORE_API_PREFIX}/ws`,
  )
  socket.onopen = () => {
    socket.send(`{"type":"subscribe","topic_name":"all"}`)
  }

  socket.onmessage = async (event) => {
    const data = JSON.parse(event.data)

    if (data.type == 'topic_message') {
      if (data.topic_name == 'order_finished') {
        toast.removeGroup('br')
        toast.add({
          severity: 'success',
          summary: t('order_finished'),
          detail: `order with id ( ${data.order_id} ) ${t('order_ready_to_serve')}`,
          life: 3000,
          group: 'br',
        })

        const notification = new Notification()
        notification.description = `order with id #${data.order_id} ${t('order_ready_to_serve')}`
        notification.severity = 'success'
        notification.topic_name = t('order_finished')
        notification.type = 'topic_message'
        notifications.value.push(notification)

        orders.value = orders.value.filter((o) => o.id !== data.order_id)
      } else if (data.topic_name == 'order_submitted') {
        orders.value.push(data.order)
        displayOrders()
      } else {
        const notification = new Notification()
        notification.description = data.message
        notification.severity = data.severity
        notification.topic_name = data.topic_name
        notification.type = data.type
        notifications.value.push(notification)

        toast.removeGroup('br')
        toast.add({
          severity: data.severity,
          summary: data.topic_name,
          detail: data.message,
          life: 30000,
          group: 'br',
        })

        chats.value.push({
          type: 'notification',
          message: data.message,
          topic_name: data.topic_name,
        })
      }
    }
  }
  socket.onerror = () => {}
  socket.onclose = () => {
    const retryConnection = async () => {
      if (socket.readyState !== WebSocket.OPEN) {
        await new Promise((r) => setTimeout(r, 5000))
        startWebsocket()
      }
    }
    retryConnection()
  }
}

const orderFinished = (order) => {
  orders.value = orders.value.filter((o) => o.id !== order.id)
}

const loadOrders = () => {
  axios
    .get(
      `http://${import.meta.env.VITE_APP_BACKEND_HOST}${import.meta.env.VITE_APP_MODULE_CORE_API_PREFIX}/api/orders?filter[state]=!finished&filter[state]=!stashed`,
      {
        headers: {
          Authorization: `Bearer ${auth.accessToken.value}`,
        },
      },
    )
    .then((result) => {
      if (result.data.data == null) {
        orders.value = []
      } else {
        orders.value = result.data.data
        displayOrders()
      }
    })
}

const displayOrders = () => {
  for (let i = 0; i < dynamic_columns.value.length; i++) {
    dynamic_columns.value[i] = []
  }

  for (let i = 0; i < orders.value.length; i++) {
    dynamic_columns.value[i % dynamic_columns.value.length].push(orders.value[i])
  }
}

prepareLayout()
loadOrders()
loadLanguage()
startWebsocket()
</script>

<style>
.queue-order {
  margin: 5px;
}
</style>

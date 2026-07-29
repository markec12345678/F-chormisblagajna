<template>
  <div class="w-full">
    <div class="grid mx-2" style="height: calc(100vh - 120px)">
      <div class="col-12 md:col-3 lg:col-2">
        <Card class="h-full">
          <template #title>{{ $t('channels') }}</template>
          <template #content>
            <div class="flex flex-column gap-1">
              <Button v-for="ch in channels" :key="ch.id"
                :label="ch.name"
                :severity="selectedChannel === ch.id ? 'primary' : 'secondary'"
                text class="w-full text-left"
                @click="selectChannel(ch)" />
            </div>
            <Divider />
            <div class="flex flex-column gap-2">
              <label class="text-sm text-500">{{ $t('new_channel') }}</label>
              <InputText v-model="newChannelName" :placeholder="$t('channel_name')" class="w-full" />
              <Button :label="$t('create')" icon="pi pi-plus" size="small" @click="createChannel" :loading="isCreating" />
            </div>
          </template>
        </Card>
      </div>

      <div class="col-12 md:col-9 lg:col-10">
        <Card class="h-full">
          <template #title>
            <div class="flex justify-content-between align-items-center">
              <span># {{ selectedChannelName }}</span>
              <Button icon="pi pi-refresh" severity="secondary" text rounded @click="loadMessages" :loading="isLoadingMessages" />
            </div>
          </template>
          <template #content>
            <div ref="messagesContainer" class="messages-container" style="height: calc(100vh - 300px); overflow-y: auto;">
              <div v-if="messages.length === 0" class="text-center py-5">
                <i class="pi pi-comments text-5xl text-400 mb-3"></i>
                <p class="text-400">{{ $t('no_messages') }}</p>
              </div>
              <div v-for="msg in messages" :key="msg.id" class="mb-3">
                <div class="flex align-items-start gap-2">
                  <Avatar :label="getInitials(msg.sender)" shape="circle" class="mr-2" />
                  <div class="flex-1">
                    <div class="flex align-items-center gap-2">
                      <span class="font-bold">{{ msg.sender }}</span>
                      <span class="text-xs text-400">{{ formatTime(msg.created_at) }}</span>
                    </div>
                    <div class="text-sm mt-1">{{ msg.content }}</div>
                  </div>
                </div>
              </div>
            </div>
          </template>
          <template #footer>
            <div class="flex gap-2">
              <InputText v-model="newMessage" :placeholder="$t('type_message')" class="flex-1"
                @keyup.enter="sendMessage" />
              <Button icon="pi pi-send" severity="success" @click="sendMessage" :loading="isSending"
                :disabled="!newMessage.trim()" />
            </div>
          </template>
        </Card>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import Card from 'primevue/card'
import Button from 'primevue/button'
import InputText from 'primevue/inputtext'
import Avatar from 'primevue/avatar'
import Divider from 'primevue/divider'
import axios from 'axios'
import { ref, nextTick } from 'vue'
import { useI18n } from 'vue-i18n'
import { useToast } from 'primevue/usetoast'
import auth from '../services/auth'
import { getCurrentInstance } from 'vue'

const { t } = useI18n()
const toast = useToast()
const { proxy } = getCurrentInstance()

interface ChatMessage {
  id: string
  channel: string
  sender: string
  sender_id: string
  content: string
  created_at: string
}

interface ChatChannel {
  id: string
  name: string
  description: string
  is_default: boolean
}

const channels = ref<ChatChannel[]>([])
const messages = ref<ChatMessage[]>([])
const selectedChannel = ref('general')
const selectedChannelName = ref('General')
const newMessage = ref('')
const newChannelName = ref('')
const isLoadingMessages = ref(false)
const isSending = ref(false)
const isCreating = ref(false)
const messagesContainer = ref<HTMLElement | null>(null)

let refreshInterval: ReturnType<typeof setInterval> | null = null

const getInitials = (name: string) => {
  return name
    .split(' ')
    .map((n) => n[0])
    .join('')
    .toUpperCase()
    .substring(0, 2)
}

const formatTime = (date: string) => {
  if (!date) return ''
  const d = new Date(date)
  return d.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })
}

const selectChannel = (ch: ChatChannel) => {
  selectedChannel.value = ch.id
  selectedChannelName.value = ch.name
  loadMessages()
}

const scrollToBottom = async () => {
  await nextTick()
  if (messagesContainer.value) {
    messagesContainer.value.scrollTop = messagesContainer.value.scrollHeight
  }
}

const loadMessages = async () => {
  isLoadingMessages.value = true
  try {
    const response = await axios.get(
      `http://${import.meta.env.VITE_APP_BACKEND_HOST}/chat/api/messages`,
      {
        headers: { Authorization: `Bearer ${auth.accessToken.value}` },
        params: { channel: selectedChannel.value }
      }
    )
    messages.value = response.data.data || []
    scrollToBottom()
  } catch {
    toast.add({ severity: 'error', summary: t('failed'), detail: t('load_failed'), group: 'br', life: 3000 })
  } finally {
    isLoadingMessages.value = false
  }
}

const sendMessage = async () => {
  if (!newMessage.value.trim()) return

  isSending.value = true
  const user = proxy.$auth.currentUser?.value
  try {
    await axios.post(
      `http://${import.meta.env.VITE_APP_BACKEND_HOST}/chat/api/messages`,
      {
        channel: selectedChannel.value,
        sender: user?.username || 'Anonymous',
        sender_id: user?.id || '',
        content: newMessage.value,
      },
      { headers: { Authorization: `Bearer ${auth.accessToken.value}` } }
    )
    newMessage.value = ''
    loadMessages()
  } catch {
    toast.add({ severity: 'error', summary: t('failed'), detail: t('send_failed'), group: 'br', life: 3000 })
  } finally {
    isSending.value = false
  }
}

const createChannel = async () => {
  if (!newChannelName.value.trim()) return

  isCreating.value = true
  try {
    await axios.post(
      `http://${import.meta.env.VITE_APP_BACKEND_HOST}/chat/api/channels`,
      {
        name: newChannelName.value,
        description: newChannelName.value,
      },
      { headers: { Authorization: `Bearer ${auth.accessToken.value}` } }
    )
    toast.add({ severity: 'success', summary: t('success'), detail: t('channel_created'), group: 'br', life: 3000 })
    newChannelName.value = ''
    loadChannels()
  } catch {
    toast.add({ severity: 'error', summary: t('failed'), detail: t('create_failed'), group: 'br', life: 3000 })
  } finally {
    isCreating.value = false
  }
}

const loadChannels = async () => {
  try {
    const response = await axios.get(
      `http://${import.meta.env.VITE_APP_BACKEND_HOST}/chat/api/channels`,
      { headers: { Authorization: `Bearer ${auth.accessToken.value}` } }
    )
    channels.value = response.data.data || []
  } catch {
    toast.add({ severity: 'error', summary: t('failed'), detail: t('load_failed'), group: 'br', life: 3000 })
  }
}

const init = async () => {
  await loadChannels()
  await loadMessages()
  refreshInterval = setInterval(loadMessages, 10000)
}

init()

import { onUnmounted } from 'vue'
onUnmounted(() => {
  if (refreshInterval) clearInterval(refreshInterval)
})
</script>

<style scoped>
.messages-container {
  scroll-behavior: smooth;
}
</style>

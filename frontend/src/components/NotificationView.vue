<template>
  <div>
    <Message :severity="notification.severity" @close="emit('closed')" class="my-1">
      <div class="flex flex-column">
        <h4 class="m-0 mb-1">{{ props.notification.topic_name }}</h4>
        {{ props.notification.description }}
        <span class="text-gray-400 mt-1" style="font-size: 0.8rem">{{ timePassed }}</span>
      </div>
    </Message>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import Message from 'primevue/message'
import { formatDistanceToNow } from 'date-fns'
import { enUS } from 'date-fns/locale'
import { Notification } from '@/classes/Notification'
import { useI18n } from 'vue-i18n'
const { locale } = useI18n()

const timePassed = ref()

const emit = defineEmits(['closed'])

const props = defineProps({
  notification: {
    type: Notification,
    required: true,
  },
})

const getLocale = () => {
  if (locale.value === 'sl') return undefined
  if (locale.value === 'ar') return undefined
  return enUS
}

const updateElapsedTime = () => {
  const date = new Date(props.notification.date)
  timePassed.value = formatDistanceToNow(date, { addSuffix: true, locale: getLocale() })
  setInterval(function () {
    timePassed.value = formatDistanceToNow(date, { addSuffix: true, locale: getLocale() })
  }, 10000)
}

const init = () => {
  updateElapsedTime()
}

init()
</script>

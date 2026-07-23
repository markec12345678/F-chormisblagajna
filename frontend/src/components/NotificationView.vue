<template>
  <div>
    <Message :severity="notification.severity" @close="emit('closed')" class="my-1">
      <div class="flex flex-column">
        <h4 class="m-0 mb-1">{{ props.notification.topic_name }}</h4>
        {{ props.notification.description }}
        <span style="color: gray; font-size: 0.8rem" class="mt-1">{{ timePassed }}</span>
      </div>
    </Message>
  </div>
</template>

<script setup lang="ts">
import { defineProps, defineEmits, ref } from 'vue'
import Message from 'primevue/message'
import { formatDistanceToNow } from 'date-fns'
import { Notification } from '@/classes/Notification'

const timePassed = ref()

const emit = defineEmits(['closed'])

const props = defineProps({
  notification: {
    type: Notification,
    required: true,
  },
})

const updateElapsedTime = () => {
  const date = new Date(props.notification.date)
  timePassed.value = formatDistanceToNow(date, { addSuffix: true })
  setInterval(function () {
    timePassed.value = formatDistanceToNow(date, { addSuffix: true })
  }, 10000)
}

const init = () => {
  updateElapsedTime()
}

init()
</script>

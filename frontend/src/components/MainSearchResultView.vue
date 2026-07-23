<template>
  <div>
    <Message severity="secondary">
      <template #container>
        <div class="px-3 py-2 flex justify-content-between align-items-center">
          <span>{{ props.order.display_id }}</span>
          <div class="flex justify-content-center align-items-center">
            <Badge :value="order_status.title" :severity="order_status.severity" />
            <Badge :value="payment_status.title" :severity="payment_status.severity" class="mx-1" />
          </div>
          <span>{{ props.order.items.length }} {{ $t('items') }}</span>
          <Button
            icon="pi pi-book"
            severity="secondary"
            aria-label="Info"
            @click="emit('view-order-pressed')"
          />
        </div>
      </template>
    </Message>
  </div>
</template>

<script setup lang="ts">
import { defineProps, computed, defineEmits } from 'vue'
import { useI18n } from 'vue-i18n'
const { t } = useI18n()
import Message from 'primevue/message'
import Button from 'primevue/button'
import Badge from 'primevue/badge'
import Order from '@/classes/Order'

const emit = defineEmits(['view-order-pressed'])

const props = defineProps({
  order: {
    type: Order,
    required: true,
  },
})

const payment_status = computed(() => {
  if (props.order.is_paid) {
    return {
      title: t('paid'),
      severity: 'success',
    }
  }

  return {
    title: t('unpaid'),
    severity: 'warn',
  }
})

const order_status = computed(() => {
  if (props.order.state == '' || props.order.state == 'pending') {
    return {
      title: t('pending'),
      severity: 'secondary',
    }
  }

  if (props.order.state == 'cancelled') {
    return {
      title: t('cancelled'),
      severity: 'danger',
    }
  }

  if (props.order.state == 'in_progress') {
    return {
      title: t('in_progress'),
      severity: 'info',
    }
  }

  if (props.order.state == 'finished') {
    return {
      title: t('finished'),
      severity: 'success',
    }
  }

  if (props.order.state == 'stashed') {
    return {
      title: t('stashed'),
      severity: 'secondary',
    }
  }

  return {}
})
</script>

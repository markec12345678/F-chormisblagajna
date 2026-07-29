<template>
  <div class="display-fullscreen" :class="themeClass">
    <div class="welcome-section" v-if="content?.welcome_message">
      <h1>{{ content.welcome_message }}</h1>
    </div>

    <div class="slides-section">
      <div
        v-for="(item, idx) in displayItems"
        :key="idx"
        class="slide"
        v-show="currentSlide === idx"
      >
        <div v-if="item.type === 'promotions'" class="promotion-slide">
          <i class="pi pi-star" style="font-size: 3rem"></i>
          <h2>{{ $t('promotions') }}</h2>
        </div>

        <div v-if="item.type === 'menu'" class="menu-slide">
          <h2>{{ $t('our_menu') }}</h2>
          <p>{{ $t('menu_placeholder') }}</p>
        </div>

        <div v-if="item.type === 'order_status'" class="orders-slide">
          <h2>{{ $t('current_orders') }}</h2>
          <div class="orders-grid">
            <div v-for="(order, oi) in item.content" :key="oi" class="order-card">
              <div class="order-id">{{ order?.display_id || order?._id }}</div>
              <Tag :severity="order?.status === 'preparing' ? 'warn' : 'info'">
                {{ order?.status }}
              </Tag>
            </div>
          </div>
        </div>
      </div>
    </div>

    <div class="indicator-section">
      <span
        v-for="(item, idx) in displayItems"
        :key="idx"
        class="dot"
        :class="{ active: currentSlide === idx }"
        @click="currentSlide = idx"
      >
      </span>
    </div>
  </div>
</template>

<script setup lang="ts">
import Tag from 'primevue/tag'
import axios from 'axios'
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useRoute } from 'vue-router'

const route = useRoute()

interface DisplayItem {
  type: string
  content: any
}

interface DisplayContent {
  items: DisplayItem[]
  interval: number
  welcome_message: string
  theme: string
}

const content = ref<DisplayContent | null>(null)
const currentSlide = ref(0)
let intervalTimer: number | null = null

const displayItems = computed(() => content.value?.items || [])
const themeClass = computed(() => {
  const t = content.value?.theme || 'light'
  return `theme-${t}`
})

const loadDisplay = async () => {
  try {
    const id = route.params.id
    const response = await axios.get(
      `http://${import.meta.env.VITE_APP_BACKEND_HOST}/customerdisplay/api/display/${id}`,
    )
    content.value = response.data.data

    if (content.value?.interval && content.value.items.length > 1) {
      startAutoSlide(content.value.interval)
    }
  } catch {
    content.value = { items: [], interval: 0, welcome_message: '', theme: 'light' }
  }
}

const startAutoSlide = (intervalSec: number) => {
  stopAutoSlide()
  intervalTimer = window.setInterval(() => {
    if (displayItems.value.length > 0) {
      currentSlide.value = (currentSlide.value + 1) % displayItems.value.length
    }
  }, intervalSec * 1000)
}

const stopAutoSlide = () => {
  if (intervalTimer !== null) {
    clearInterval(intervalTimer)
    intervalTimer = null
  }
}

onMounted(loadDisplay)
onUnmounted(stopAutoSlide)
</script>

<style scoped>
.display-fullscreen {
  width: 100vw;
  height: 100vh;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  overflow: hidden;
  font-family: sans-serif;
}
.theme-light {
  background: #f8f9fa;
  color: #333;
}
.theme-dark {
  background: #1a1a2e;
  color: #eee;
}
.theme-colorful {
  background: linear-gradient(135deg, #667eea, #764ba2);
  color: #fff;
}
.welcome-section {
  text-align: center;
  margin-bottom: 2rem;
}
.slides-section {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  width: 90%;
}
.slide {
  text-align: center;
  animation: fadeIn 0.5s;
}
@keyframes fadeIn {
  from {
    opacity: 0;
  }
  to {
    opacity: 1;
  }
}
.orders-grid {
  display: flex;
  flex-wrap: wrap;
  gap: 1rem;
  justify-content: center;
}
.order-card {
  background: rgba(255, 255, 255, 0.1);
  padding: 1rem 2rem;
  border-radius: 8px;
}
.indicator-section {
  position: fixed;
  bottom: 2rem;
  display: flex;
  gap: 0.5rem;
}
.dot {
  width: 12px;
  height: 12px;
  border-radius: 50%;
  background: rgba(255, 255, 255, 0.4);
  cursor: pointer;
}
.dot.active {
  background: rgba(255, 255, 255, 0.9);
}
</style>

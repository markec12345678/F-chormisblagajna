<template>
  <div class="w-full">
    <div class="grid mx-2">
      <div class="col-12 flex">
        <div class="grid w-full">
          <div class="col-12">
            <h3>{{ $t('ai_search') }}</h3>
          </div>

          <div class="col-12">
            <Card>
              <template #content>
                <div class="flex gap-2 align-items-center">
                  <div class="flex-1">
                    <IconField iconPosition="left">
                      <InputIcon class="pi pi-search" />
                      <InputText
                        v-model="searchQuery"
                        :placeholder="$t('search') + '...'"
                        class="w-full"
                        @keyup.enter="performSearch"
                      />
                    </IconField>
                  </div>
                  <Button
                    :icon="isRecording ? 'pi pi-stop-circle' : 'pi pi-microphone'"
                    :severity="isRecording ? 'danger' : 'secondary'"
                    @click="toggleVoiceRecording"
                    :aria-label="$t('ai_voice')"
                    rounded
                    raised
                  />
                  <Button
                    icon="pi pi-search"
                    :label="$t('search')"
                    @click="performSearch"
                    :loading="isSearching"
                    :disabled="isSearching"
                    rounded
                    raised
                  />
                </div>
                <div v-if="isRecording" class="mt-2">
                  <div class="flex align-items-center gap-2 text-red-500">
                    <span
                      class="pi pi-circle-fill text-xs"
                      style="animation: blink 1s infinite"
                    ></span>
                    <span class="text-sm font-medium">{{ $t('ai_voice') }}...</span>
                  </div>
                </div>
              </template>
            </Card>
          </div>

          <div class="col-12 lg:col-8">
            <Card class="h-full">
              <template #title>
                <div class="flex justify-between align-items-center">
                  <span>{{ $t('search_results') }}</span>
                  <Tag
                    v-if="searchResults.length > 0"
                    :value="String(searchResults.length)"
                    severity="info"
                  />
                </div>
              </template>
              <template #content>
                <ProgressSpinner
                  v-if="isSearching"
                  style="width: 40px; height: 40px"
                  class="flex mx-auto"
                />
                <div v-else-if="searchResults.length === 0 && hasSearched" class="text-center py-4">
                  <i class="pi pi-search" style="font-size: 2rem; opacity: 0.3"></i>
                  <p class="mt-2 text-slate-400">{{ $t('no_results') }}</p>
                </div>
                <div
                  v-else-if="searchResults.length === 0 && !hasSearched"
                  class="text-center py-4"
                >
                  <i class="pi pi-sparkles" style="font-size: 2rem; opacity: 0.3"></i>
                  <p class="mt-2 text-slate-400">{{ $t('ai_search_description') }}</p>
                </div>
                <DataTable v-else :value="searchResults" stripedRows tableStyle="width: 100%">
                  <Column field="name" :header="$t('name')">
                    <template #body="slotProps">
                      <div>
                        <span class="font-medium">{{ slotProps.data.name }}</span>
                        <div v-if="slotProps.data.description" class="text-sm text-500 mt-1">
                          {{ slotProps.data.description }}
                        </div>
                      </div>
                    </template>
                  </Column>
                  <Column field="category" :header="$t('category', 1)"></Column>
                  <Column field="price" :header="$t('total')">
                    <template #body="slotProps">
                      {{ formatPrice(slotProps.data.price) }}
                    </template>
                  </Column>
                  <Column field="score" :header="$t('score')">
                    <template #body="slotProps">
                      <div class="flex align-items-center gap-2">
                        <div
                          class="h-2 border-round overflow-hidden"
                          style="width: 60px; background: var(--p-surface-200)"
                        >
                          <div
                            class="h-full border-round"
                            :style="{
                              width: slotProps.data.score * 100 + '%',
                              background:
                                slotProps.data.score > 0.7
                                  ? '#22c55e'
                                  : slotProps.data.score > 0.4
                                    ? '#eab308'
                                    : '#ef4444',
                            }"
                          ></div>
                        </div>
                        <span class="text-sm text-500"
                          >{{ Math.round(slotProps.data.score * 100) }}%</span
                        >
                      </div>
                    </template>
                  </Column>
                  <Column :header="$t('actions')">
                    <template #body="slotProps">
                      <Button
                        icon="pi pi-plus"
                        severity="success"
                        :aria-label="$t('add_to_order')"
                        @click="addToOrder(slotProps.data)"
                        rounded
                        raised
                        size="small"
                      />
                    </template>
                  </Column>
                </DataTable>
              </template>
            </Card>
          </div>

          <div class="col-12 lg:col-4">
            <Card class="h-full">
              <template #title>
                <div class="flex justify-between align-items-center">
                  <span>{{ $t('smart_suggestions') }}</span>
                  <Tag
                    v-if="suggestions.length > 0"
                    :value="String(suggestions.length)"
                    severity="warn"
                  />
                </div>
              </template>
              <template #content>
                <ProgressSpinner
                  v-if="isLoadingSuggestions"
                  style="width: 40px; height: 40px"
                  class="flex mx-auto"
                />
                <div v-else-if="suggestions.length === 0" class="text-center py-4">
                  <i class="pi pi-lightbulb" style="font-size: 2rem; opacity: 0.3"></i>
                  <p class="mt-2 text-slate-400">{{ $t('ai_suggestions') }}</p>
                </div>
                <div v-else class="flex flex-column gap-3">
                  <div
                    v-for="suggestion in suggestions"
                    :key="suggestion.product_id"
                    class="border-1 border-round p-3 flex flex-column gap-1"
                    style="border-color: var(--p-surface-200)"
                  >
                    <div class="flex justify-between align-items-center">
                      <span class="font-medium">{{ suggestion.product_name }}</span>
                      <Tag :value="Math.round(suggestion.score * 100) + '%'" severity="info" />
                    </div>
                    <p class="text-sm text-500 m-0">{{ suggestion.reason }}</p>
                    <Button
                      icon="pi pi-plus"
                      :label="$t('add_to_order')"
                      severity="success"
                      size="small"
                      @click="addSuggestionToOrder(suggestion)"
                      class="mt-2"
                    />
                  </div>
                </div>
              </template>
            </Card>

            <Card class="mt-3">
              <template #title>
                <span>{{ $t('voice_ordering') }}</span>
              </template>
              <template #content>
                <p class="text-sm text-500 m-0 mb-3">{{ $t('voice_ordering_description') }}</p>
                <Button
                  :icon="isRecording ? 'pi pi-stop-circle' : 'pi pi-microphone'"
                  :label="isRecording ? $t('stop_recording') : $t('start_recording')"
                  :severity="isRecording ? 'danger' : 'primary'"
                  @click="toggleVoiceRecording"
                  class="w-full"
                  raised
                />
                <div
                  v-if="voiceTranscript"
                  class="mt-3 p-3 border-round"
                  style="background: var(--p-surface-100)"
                >
                  <p class="text-sm font-medium mb-1">{{ $t('transcription') }}:</p>
                  <p class="text-sm m-0">{{ voiceTranscript }}</p>
                </div>
              </template>
            </Card>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import Card from 'primevue/card'
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import Button from 'primevue/button'
import InputText from 'primevue/inputtext'
import IconField from 'primevue/iconfield'
import InputIcon from 'primevue/inputicon'
import Tag from 'primevue/tag'
import ProgressSpinner from 'primevue/progressspinner'
import axios from 'axios'
import { ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { useToast } from 'primevue/usetoast'
import auth from '../services/auth'
import type { AISearchResult, SmartSuggestion } from '@/types'

const { t } = useI18n()
const toast = useToast()

const searchQuery = ref('')
const searchResults = ref<AISearchResult[]>([])
const suggestions = ref<SmartSuggestion[]>([])
const isSearching = ref(false)
const hasSearched = ref(false)
const isLoadingSuggestions = ref(false)

const isRecording = ref(false)
const voiceTranscript = ref('')
let mediaRecorder: MediaRecorder | null = null
let audioChunks: Blob[] = []
interface SpeechRec {
  start: () => void
  stop: () => void
  onresult: (event: Event) => void
  continuous: boolean
  interimResults: boolean
  lang: string
}
let speechRecognition: SpeechRec | null = null

const formatPrice = (price: number) => {
  return new Intl.NumberFormat('en-US', {
    minimumFractionDigits: 2,
    maximumFractionDigits: 2,
  }).format(price)
}

const performSearch = () => {
  if (!searchQuery.value.trim()) return

  isSearching.value = true
  hasSearched.value = true

  axios
    .get(
      `http://${import.meta.env.VITE_APP_BACKEND_HOST}${import.meta.env.VITE_APP_MODULE_CORE_API_PREFIX}/ai/api/ai/search`,
      {
        params: { q: searchQuery.value },
        headers: {
          Authorization: `Bearer ${auth.accessToken.value}`,
        },
      },
    )
    .then((response) => {
      searchResults.value = response.data.data || []
    })
    .catch(() => {
      toast.add({
        severity: 'error',
        summary: t('failed'),
        detail: t('search_failed'),
        group: 'br',
      })
      searchResults.value = []
    })
    .finally(() => {
      isSearching.value = false
    })
}

const toggleVoiceRecording = () => {
  if (isRecording.value) {
    stopRecording()
  } else {
    startRecording()
  }
}

const startRecording = () => {
  if ('webkitSpeechRecognition' in window || 'SpeechRecognition' in window) {
    /* eslint-disable @typescript-eslint/no-explicit-any */
    const SpeechRecognition =
      (window as any).SpeechRecognition || (window as any).webkitSpeechRecognition
    /* eslint-enable @typescript-eslint/no-explicit-any */
    speechRecognition = new SpeechRecognition()
    speechRecognition.continuous = false
    speechRecognition.interimResults = true
    speechRecognition.lang = 'en-US'

    speechRecognition.onresult = (event) => {
      const transcript = Array.from(event.results)
        .map((result) => result[0].transcript)
        .join('')
      voiceTranscript.value = transcript

      if (event.results[0].isFinal) {
        searchQuery.value = transcript
        performVoiceSearch(transcript)
      }
    }

    speechRecognition.onerror = () => {
      isRecording.value = false
      toast.add({
        severity: 'error',
        summary: t('failed'),
        detail: t('voice_recognition_failed'),
        group: 'br',
      })
    }

    speechRecognition.onend = () => {
      isRecording.value = false
    }

    speechRecognition.start()
    isRecording.value = true
    return
  }

  if ('MediaRecorder' in window) {
    navigator.mediaDevices
      .getUserMedia({ audio: true })
      .then((stream) => {
        mediaRecorder = new MediaRecorder(stream)
        audioChunks = []

        mediaRecorder.ondataavailable = (event) => {
          audioChunks.push(event.data)
        }

        mediaRecorder.onstop = () => {
          const audioBlob = new Blob(audioChunks, { type: 'audio/webm' })
          sendVoiceToServer(audioBlob)
          stream.getTracks().forEach((track) => track.stop())
        }

        mediaRecorder.start()
        isRecording.value = true
      })
      .catch(() => {
        toast.add({
          severity: 'error',
          summary: t('failed'),
          detail: t('microphone_access_denied'),
          group: 'br',
        })
      })
  }
}

const stopRecording = () => {
  if (speechRecognition) {
    speechRecognition.stop()
    isRecording.value = false
    return
  }
  if (mediaRecorder && mediaRecorder.state !== 'inactive') {
    mediaRecorder.stop()
    isRecording.value = false
  }
}

const performVoiceSearch = (transcript: string) => {
  isSearching.value = true
  hasSearched.value = true

  axios
    .post(
      `http://${import.meta.env.VITE_APP_BACKEND_HOST}${import.meta.env.VITE_APP_MODULE_CORE_API_PREFIX}/ai/api/ai/voice`,
      { transcript },
      {
        headers: {
          Authorization: `Bearer ${auth.accessToken.value}`,
        },
      },
    )
    .then((response) => {
      searchResults.value = response.data.data || []
    })
    .catch(() => {
      toast.add({
        severity: 'error',
        summary: t('failed'),
        detail: t('voice_search_failed'),
        group: 'br',
      })
    })
    .finally(() => {
      isSearching.value = false
    })
}

const sendVoiceToServer = (audioBlob: Blob) => {
  isSearching.value = true
  hasSearched.value = true

  const formData = new FormData()
  formData.append('audio', audioBlob, 'voice.webm')

  axios
    .post(
      `http://${import.meta.env.VITE_APP_BACKEND_HOST}${import.meta.env.VITE_APP_MODULE_CORE_API_PREFIX}/ai/api/ai/voice`,
      formData,
      {
        headers: {
          Authorization: `Bearer ${auth.accessToken.value}`,
          'Content-Type': 'multipart/form-data',
        },
      },
    )
    .then((response) => {
      searchResults.value = response.data.data || []
      if (response.data.transcript) {
        voiceTranscript.value = response.data.transcript
        searchQuery.value = response.data.transcript
      }
    })
    .catch(() => {
      toast.add({
        severity: 'error',
        summary: t('failed'),
        detail: t('voice_search_failed'),
        group: 'br',
      })
    })
    .finally(() => {
      isSearching.value = false
    })
}

const loadSuggestions = () => {
  isLoadingSuggestions.value = true

  axios
    .get(
      `http://${import.meta.env.VITE_APP_BACKEND_HOST}${import.meta.env.VITE_APP_MODULE_CORE_API_PREFIX}/ai/api/ai/suggestions`,
      {
        headers: {
          Authorization: `Bearer ${auth.accessToken.value}`,
        },
      },
    )
    .then((response) => {
      suggestions.value = response.data.data || []
    })
    .catch(() => {
      suggestions.value = []
    })
    .finally(() => {
      isLoadingSuggestions.value = false
    })
}

const addToOrder = (product: AISearchResult) => {
  toast.add({
    severity: 'success',
    summary: t('success'),
    detail: product.name + ' ' + t('added_to_order'),
    group: 'br',
  })
}

const addSuggestionToOrder = (suggestion: SmartSuggestion) => {
  toast.add({
    severity: 'success',
    summary: t('success'),
    detail: suggestion.product_name + ' ' + t('added_to_order'),
    group: 'br',
  })
}

loadSuggestions()
</script>

<style scoped>
@keyframes blink {
  0%,
  100% {
    opacity: 1;
  }
  50% {
    opacity: 0.3;
  }
}
</style>

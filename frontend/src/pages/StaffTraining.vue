<template>
  <div class="w-full">
    <div class="grid mx-2">
      <div class="col-12">
        <h3>{{ $t('staff_training') }}</h3>
      </div>

      <div class="col-12" v-if="activeSession">
        <Card>
          <template #title>
            <div class="flex justify-content-between align-items-center">
              <span>{{ $t('training_in_progress') }}</span>
              <Tag :value="`${activeSession.steps_done}/${activeSession.total_steps}`" severity="info" />
            </div>
          </template>
          <template #content>
            <div v-if="currentStep" class="text-center py-4">
              <h4>{{ currentStep.title }}</h4>
              <p class="text-500">{{ currentStep.description }}</p>
              <div class="flex justify-content-center gap-2 mt-3">
                <Button :label="$t('complete_step')" @click="doCompleteStep" :loading="completing" />
                <Button :label="$t('finish_training')" severity="secondary" @click="doCompleteSession" :loading="finishing" />
              </div>
            </div>
            <ProgressBar :value="stepProgress" class="mt-3" />
          </template>
        </Card>
      </div>

      <div class="col-12 md:col-6" v-for="mod in modules" :key="mod.key">
        <Card>
          <template #content>
            <div class="flex align-items-center gap-3">
              <i :class="mod.icon" style="font-size: 2rem"></i>
              <div class="flex-1">
                <div class="font-bold text-lg">{{ mod.name }}</div>
                <div class="text-500 text-sm">{{ mod.description }}</div>
                <small class="text-500">{{ mod.steps }} {{ $t('steps') }}</small>
              </div>
              <div class="flex flex-column align-items-center gap-1">
                <div v-if="getProgress(mod.key) as any" class="text-center">
                  <Tag :severity="modProgressSeverity(mod.key)" class="mb-1">
                    {{ getProgress(mod.key)?.completion_pct?.toFixed(0) || 0 }}%
                  </Tag>
                  <div class="text-xs text-500">
                    {{ getProgress(mod.key)?.score }}/{{ getProgress(mod.key)?.max_score }}
                  </div>
                </div>
                <Button :label="getProgress(mod.key)?.completed ? $t('retake') : $t('start')"
                  :severity="getProgress(mod.key)?.completed ? 'secondary' : 'primary'"
                  size="small" @click="startModule(mod)" />
              </div>
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
import Tag from 'primevue/tag'
import ProgressBar from 'primevue/progressbar'
import axios from 'axios'
import { ref, computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { useToast } from 'primevue/usetoast'
import auth from '../services/auth'

const { t } = useI18n()
const toast = useToast()

interface TrainingModule {
  key: string
  name: string
  description: string
  icon: string
  steps: number
}

interface TrainingStep {
  id: string
  title: string
  description: string
}

interface TrainingSession {
  id: string
  module: string
  steps_done: number
  total_steps: number
  score: number
  max_score: number
  completed: boolean
}

interface TrainingProgress {
  module: string
  steps_done: number
  total_steps: number
  score: number
  max_score: number
  completion_pct: number
  completed: boolean
}

const modules = ref<TrainingModule[]>([])
const progressList = ref<TrainingProgress[]>([])
const activeSession = ref<TrainingSession | null>(null)
const currentSteps = ref<TrainingStep[]>([])
const completing = ref(false)
const finishing = ref(false)

const currentStep = computed(() => {
  if (!activeSession.value || currentSteps.value.length === 0) return null
  const idx = activeSession.value.steps_done
  return currentSteps.value[idx] || null
})

const stepProgress = computed(() => {
  if (!activeSession.value) return 0
  return (activeSession.value.steps_done / activeSession.value.total_steps) * 100
})

const getProgress = (moduleKey: string) => {
  return progressList.value.find((p) => p.module === moduleKey)
}

const modProgressSeverity = (key: string) => {
  const p = getProgress(key)
  if (!p) return 'secondary'
  if (p.completed) return 'success'
  if (p.completion_pct > 50) return 'warn'
  return 'secondary'
}

const startModule = async (mod: TrainingModule) => {
  try {
    const response = await axios.post(
      `http://${import.meta.env.VITE_APP_BACKEND_HOST}/training/api/sessions`,
      { user_id: auth.user.value?.id || '', module: mod.key },
      { headers: { Authorization: `Bearer ${auth.accessToken.value}` } }
    )
    activeSession.value = response.data.data

    const stepsResponse = await axios.get(
      `http://${import.meta.env.VITE_APP_BACKEND_HOST}/training/api/modules/${mod.key}/steps`,
      { headers: { Authorization: `Bearer ${auth.accessToken.value}` } }
    )
    currentSteps.value = stepsResponse.data.data || []

    await loadProgress()
  } catch {
    toast.add({ severity: 'error', summary: t('failed'), group: 'br', life: 3000 })
  }
}

const doCompleteStep = async () => {
  if (!activeSession.value) return
  completing.value = true
  try {
    await axios.post(
      `http://${import.meta.env.VITE_APP_BACKEND_HOST}/training/api/sessions/${activeSession.value.id}/step`,
      {},
      { headers: { Authorization: `Bearer ${auth.accessToken.value}` } }
    )
    activeSession.value.steps_done++
    await loadProgress()
  } catch {
    toast.add({ severity: 'error', summary: t('failed'), group: 'br', life: 3000 })
  } finally {
    completing.value = false
  }
}

const doCompleteSession = async () => {
  if (!activeSession.value) return
  finishing.value = true
  try {
    await axios.post(
      `http://${import.meta.env.VITE_APP_BACKEND_HOST}/training/api/sessions/${activeSession.value.id}/complete`,
      {},
      { headers: { Authorization: `Bearer ${auth.accessToken.value}` } }
    )
    activeSession.value.completed = true
    activeSession.value = null
    currentSteps.value = []
    toast.add({ severity: 'success', summary: t('training_completed'), group: 'br', life: 3000 })
    await loadProgress()
  } catch {
    toast.add({ severity: 'error', summary: t('failed'), group: 'br', life: 3000 })
  } finally {
    finishing.value = false
  }
}

const loadProgress = async () => {
  try {
    const userId = auth.user.value?.id || ''
    const response = await axios.get(
      `http://${import.meta.env.VITE_APP_BACKEND_HOST}/training/api/users/${userId}/progress`,
      { headers: { Authorization: `Bearer ${auth.accessToken.value}` } }
    )
    progressList.value = response.data.data || []
  } catch {
    // ignore
  }
}

const loadModules = async () => {
  try {
    const response = await axios.get(
      `http://${import.meta.env.VITE_APP_BACKEND_HOST}/training/api/modules`,
      { headers: { Authorization: `Bearer ${auth.accessToken.value}` } }
    )
    modules.value = response.data.data || []
  } catch {
    toast.add({ severity: 'error', summary: t('load_failed'), group: 'br', life: 3000 })
  }
}

loadModules()
loadProgress()
</script>

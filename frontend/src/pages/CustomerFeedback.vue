<template>
  <div class="w-full">
    <div class="grid mx-2">
      <div class="col-12">
        <h3>{{ $t('customer_feedback') }}</h3>
      </div>

      <div class="col-12" v-if="summary">
        <div class="grid">
          <div class="col-12 md:col-3">
            <Card>
              <template #content>
                <div class="text-center">
                  <div class="text-4xl font-bold text-primary">
                    {{ summary.average_rating?.toFixed(1) }}
                  </div>
                  <div class="text-500">{{ $t('avg_rating') }}</div>
                  <div class="mt-1">{{ renderStars(summary.average_rating) }}</div>
                </div>
              </template>
            </Card>
          </div>
          <div class="col-12 md:col-3">
            <Card>
              <template #content>
                <div class="text-center">
                  <div class="text-4xl font-bold">{{ summary.total_feedbacks }}</div>
                  <div class="text-500">{{ $t('total_reviews') }}</div>
                </div>
              </template>
            </Card>
          </div>
          <div class="col-12 md:col-3">
            <Card>
              <template #content>
                <div class="text-center">
                  <div class="text-2xl font-bold text-green-500">{{ getCategoryAvg('food') }}</div>
                  <div class="text-500">{{ $t('food_rating') }}</div>
                </div>
              </template>
            </Card>
          </div>
          <div class="col-12 md:col-3">
            <Card>
              <template #content>
                <div class="text-center">
                  <div class="text-2xl font-bold text-orange-500">
                    {{ getCategoryAvg('service') }}
                  </div>
                  <div class="text-500">{{ $t('service_rating') }}</div>
                </div>
              </template>
            </Card>
          </div>
        </div>
      </div>

      <div class="col-12">
        <Card>
          <template #title>{{ $t('all_reviews') }}</template>
          <template #content>
            <DataTable :value="feedbacks" stripedRows :loading="isLoading">
              <template #empty>{{ $t('no_feedback') }}</template>
              <Column field="rating" :header="$t('rating')">
                <template #body="slotProps">{{ renderStars(slotProps?.data?.rating) }}</template>
              </Column>
              <Column field="category" :header="$t('category')">
                <template #body="slotProps">
                  <Tag :value="slotProps?.data?.category" severity="info" />
                </template>
              </Column>
              <Column field="comment" :header="$t('comment')" style="min-width: 200px"></Column>
              <Column field="anonymous" :header="$t('anonymous')">
                <template #body="slotProps">
                  <i
                    :class="
                      slotProps?.data?.anonymous
                        ? 'pi pi-check text-green-500'
                        : 'pi pi-times text-300'
                    "
                  ></i>
                </template>
              </Column>
              <Column field="responded" :header="$t('responded')">
                <template #body="slotProps">
                  <Tag :severity="slotProps?.data?.responded ? 'success' : 'secondary'">
                    {{ slotProps?.data?.responded ? $t('yes') : $t('no') }}
                  </Tag>
                </template>
              </Column>
              <Column :header="$t('actions')">
                <template #body="slotProps">
                  <Button
                    v-if="!slotProps.data.responded"
                    icon="pi pi-reply"
                    severity="info"
                    class="mr-1"
                    @click="openRespondDialog(slotProps.data)"
                    v-tooltip.left="$t('respond')"
                  />
                  <Button
                    icon="pi pi-trash"
                    severity="danger"
                    @click="deleteFeedback(slotProps.data)"
                  />
                </template>
              </Column>
            </DataTable>
          </template>
        </Card>
      </div>
    </div>

    <Dialog
      v-model:visible="respondDialogVisible"
      :header="$t('respond_to_feedback')"
      :style="{ width: '500px' }"
    >
      <div class="grid">
        <div class="col-12">
          <label>{{ $t('comment') }}</label>
          <p class="text-500 bg-gray-100 p-2 border-round">{{ selectedFeedback?.comment }}</p>
        </div>
        <div class="col-12">
          <label>{{ $t('your_response') }}</label>
          <Textarea v-model="responseText" class="w-full" rows="4" />
        </div>
      </div>
      <template #footer>
        <Button :label="$t('cancel')" severity="secondary" @click="respondDialogVisible = false" />
        <Button :label="$t('send')" @click="doRespond" :loading="responding" />
      </template>
    </Dialog>
  </div>
</template>

<script setup lang="ts">
import Card from 'primevue/card'
import Button from 'primevue/button'
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import Dialog from 'primevue/dialog'
import Textarea from 'primevue/textarea'
import Tag from 'primevue/tag'
import axios from 'axios'
import { ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { useToast } from 'primevue/usetoast'
import auth from '../services/auth'

const { t } = useI18n()
const toast = useToast()

interface Feedback {
  id: string
  order_id: string
  rating: number
  comment: string
  category: string
  anonymous: boolean
  responded: boolean
  response?: string
}

interface FeedbackSummary {
  total_feedbacks: number
  average_rating: number
  rating_distribution: Record<number, number>
  category_averages: Record<string, number>
}

const feedbacks = ref<Feedback[]>([])
const summary = ref<FeedbackSummary | null>(null)
const isLoading = ref(false)
const respondDialogVisible = ref(false)
const selectedFeedback = ref<Feedback | null>(null)
const responseText = ref('')
const responding = ref(false)

const renderStars = (rating: number) => {
  const full = '★'.repeat(Math.round(rating))
  const empty = '☆'.repeat(5 - Math.round(rating))
  return full + empty
}

const getCategoryAvg = (cat: string) => {
  if (!summary.value?.category_averages) return '-'
  const val = summary.value.category_averages[cat]
  return val ? val.toFixed(1) : '-'
}

const openRespondDialog = (fb: Feedback) => {
  selectedFeedback.value = fb
  responseText.value = fb.response || ''
  respondDialogVisible.value = true
}

const doRespond = async () => {
  if (!selectedFeedback.value) return
  responding.value = true
  try {
    await axios.post(
      `http://${import.meta.env.VITE_APP_BACKEND_HOST}/feedback/api/${selectedFeedback.value.id}/respond`,
      { response: responseText.value },
      { headers: { Authorization: `Bearer ${auth.accessToken.value}` } },
    )
    toast.add({ severity: 'success', summary: t('saved'), group: 'br', life: 2000 })
    respondDialogVisible.value = false
    await loadData()
  } catch {
    toast.add({ severity: 'error', summary: t('failed'), group: 'br', life: 3000 })
  } finally {
    responding.value = false
  }
}

const deleteFeedback = async (fb: Feedback) => {
  try {
    await axios.delete(`http://${import.meta.env.VITE_APP_BACKEND_HOST}/feedback/api/${fb.id}`, {
      headers: { Authorization: `Bearer ${auth.accessToken.value}` },
    })
    toast.add({ severity: 'success', summary: t('deleted'), group: 'br', life: 2000 })
    await loadData()
  } catch {
    toast.add({ severity: 'error', summary: t('failed'), group: 'br', life: 3000 })
  }
}

const loadData = async () => {
  isLoading.value = true
  try {
    const [listRes, summaryRes] = await Promise.all([
      axios.get(`http://${import.meta.env.VITE_APP_BACKEND_HOST}/feedback/api/list`, {
        headers: { Authorization: `Bearer ${auth.accessToken.value}` },
      }),
      axios.get(`http://${import.meta.env.VITE_APP_BACKEND_HOST}/feedback/api/summary`, {
        headers: { Authorization: `Bearer ${auth.accessToken.value}` },
      }),
    ])
    feedbacks.value = listRes.data.data || []
    summary.value = summaryRes.data.data
  } catch {
    toast.add({ severity: 'error', summary: t('load_failed'), group: 'br', life: 3000 })
  } finally {
    isLoading.value = false
  }
}

loadData()
</script>

import { defineStore } from 'pinia'
import { ref } from 'vue'
import { merchantApi, type Merchant, type DashboardStats, type MerchantListParams } from '@/api/merchant'

export const useMerchantStore = defineStore('merchant', () => {
  const list = ref<Merchant[]>([])
  const total = ref(0)
  const loading = ref(false)
  const dashboardStats = ref<DashboardStats | null>(null)

  async function fetchList(params: MerchantListParams) {
    loading.value = true
    try {
      const res = await merchantApi.list(params)
      list.value = res.list
      total.value = res.total
    } finally {
      loading.value = false
    }
  }

  async function fetchDashboard() {
    loading.value = true
    try {
      dashboardStats.value = await merchantApi.dashboard()
    } finally {
      loading.value = false
    }
  }

  return { list, total, loading, dashboardStats, fetchList, fetchDashboard }
})

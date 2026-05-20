import { defineStore } from 'pinia'
import { ref } from 'vue'
import http from '@/api/request'

export interface DashboardData {
  today_sales: number
  today_orders: number
  today_profit: number
  sales_growth: number
  orders_growth: number
  profit_growth: number
}

export interface Alert {
  id: number
  level: 'red' | 'yellow' | 'blue'
  product_name: string
  subtitle: string
  description: string
  action_text: string
  action_type: 'replenish' | 'promo' | 'clearance'
  product_id: number
}

export interface Todo {
  id: number
  content: string
  type: string
  done: boolean
}

export const useMerchantStore = defineStore('merchant', () => {
  const dashboard = ref<DashboardData | null>(null)
  const alerts = ref<Alert[]>([])
  const todos = ref<Todo[]>([])
  const loading = ref(false)

  async function fetchDashboard() {
    loading.value = true
    try {
      dashboard.value = await http.get<DashboardData>('/merchant/dashboard')
    } catch {
      // 错误已在 request 中处理，这里用空占位数据
      dashboard.value = {
        today_sales: 0, today_orders: 0, today_profit: 0,
        sales_growth: 0, orders_growth: 0, profit_growth: 0
      }
    } finally {
      loading.value = false
    }
  }

  async function fetchAlerts() {
    try {
      alerts.value = await http.get<Alert[]>('/merchant/alerts')
    } catch {
      alerts.value = []
    }
  }

  async function fetchTodos() {
    try {
      todos.value = await http.get<Todo[]>('/merchant/todos')
    } catch {
      todos.value = []
    }
  }

  return { dashboard, alerts, todos, loading, fetchDashboard, fetchAlerts, fetchTodos }
})

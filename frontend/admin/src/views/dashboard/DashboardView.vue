<template>
  <div class="dashboard">
    <h1 class="page-title">数据大盘</h1>

    <!-- 指标卡片行 -->
    <div class="stat-grid" v-loading="loading">
      <StatCard
        label="总商家数"
        :value="stats?.total_merchants ?? 0"
        color="primary"
        :growth="stats?.merchant_growth"
      />
      <StatCard
        label="活跃商家"
        :value="stats?.active_merchants ?? 0"
        color="success"
      />
      <StatCard
        label="待审核"
        :value="stats?.pending_audit ?? 0"
        color="warning"
      />
      <StatCard
        label="今日GMV"
        :value="stats?.total_gmv_today ?? 0"
        color="primary"
        prefix="¥"
        :growth="stats?.order_growth"
      />
    </div>

    <!-- 图表区域 -->
    <div class="chart-grid">
      <div class="chart-card">
        <h3 class="chart-title">商家增长趋势</h3>
        <div ref="merchantChartRef" class="chart-placeholder" />
      </div>
      <div class="chart-card">
        <h3 class="chart-title">GMV趋势（近30天）</h3>
        <div ref="gmvChartRef" class="chart-placeholder" />
      </div>
    </div>

    <!-- 排行榜 -->
    <div class="rank-card">
      <h3 class="chart-title">商家销售排行 TOP10</h3>
      <el-table :data="rankList" stripe>
        <el-table-column type="index" label="排名" width="80" />
        <el-table-column prop="shop_name" label="店铺名称" />
        <el-table-column prop="industry" label="行业">
          <template #default="{ row }">
            <el-tag size="small">{{ INDUSTRY_MAP[row.industry] || row.industry }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="gmv" label="GMV" sortable>
          <template #default="{ row }">¥{{ row.gmv?.toLocaleString() }}</template>
        </el-table-column>
        <el-table-column prop="order_count" label="订单数" sortable />
      </el-table>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import * as echarts from 'echarts'
import StatCard from '@/components/common/StatCard.vue'
import { merchantApi, type DashboardStats } from '@/api/merchant'

const INDUSTRY_MAP: Record<string, string> = {
  catering: '餐饮',
  retail: '零售',
  fresh: '生鲜',
  bakery: '烘焙'
}

const loading = ref(false)
const stats = ref<DashboardStats | null>(null)
const rankList = ref<Array<{ shop_name: string; industry: string; gmv: number; order_count: number }>>([])

const merchantChartRef = ref<HTMLElement>()
const gmvChartRef = ref<HTMLElement>()
let merchantChart: echarts.ECharts | null = null
let gmvChart: echarts.ECharts | null = null

function initCharts() {
  if (merchantChartRef.value) {
    merchantChart = echarts.init(merchantChartRef.value)
    merchantChart.setOption({
      tooltip: { trigger: 'axis' },
      xAxis: { type: 'category', data: ['1月', '2月', '3月', '4月', '5月', '6月'] },
      yAxis: { type: 'value' },
      series: [{
        name: '新增商家',
        type: 'line',
        smooth: true,
        data: [12, 28, 35, 42, 58, 75],
        lineStyle: { color: '#4a6cf7' },
        itemStyle: { color: '#4a6cf7' },
        areaStyle: { color: 'rgba(74,108,247,0.1)' }
      }]
    })
  }

  if (gmvChartRef.value) {
    gmvChart = echarts.init(gmvChartRef.value)
    gmvChart.setOption({
      tooltip: { trigger: 'axis' },
      xAxis: { type: 'category', data: ['4/20', '4/25', '4/30', '5/5', '5/10', '5/15', '5/20'] },
      yAxis: { type: 'value', axisLabel: { formatter: '¥{value}' } },
      series: [{
        name: 'GMV',
        type: 'bar',
        data: [12000, 18500, 22000, 15800, 28000, 31000, 35000],
        itemStyle: { color: '#4a6cf7', borderRadius: [4, 4, 0, 0] }
      }]
    })
  }
}

function handleResize() {
  merchantChart?.resize()
  gmvChart?.resize()
}

onMounted(async () => {
  loading.value = true
  try {
    stats.value = await merchantApi.dashboard()
    // 模拟排行榜数据（实际从接口获取）
    rankList.value = [
      { shop_name: '老王烧烤', industry: 'catering', gmv: 58600, order_count: 1240 },
      { shop_name: '鲜果坊', industry: 'fresh', gmv: 42000, order_count: 980 },
      { shop_name: '好利来面包', industry: 'bakery', gmv: 38500, order_count: 750 }
    ]
  } catch {
    // 错误已由拦截器处理，使用空数据占位
  } finally {
    loading.value = false
    initCharts()
  }

  window.addEventListener('resize', handleResize)
})

onUnmounted(() => {
  window.removeEventListener('resize', handleResize)
  merchantChart?.dispose()
  gmvChart?.dispose()
})
</script>

<style scoped>
.dashboard {
  max-width: 1440px;
}

.page-title {
  font-size: 24px;
  font-weight: 700;
  color: var(--color-text-title);
  margin-bottom: var(--space-lg);
}

.stat-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: var(--space-lg);
  margin-bottom: var(--space-lg);
}

.chart-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: var(--space-lg);
  margin-bottom: var(--space-lg);
}

.chart-card,
.rank-card {
  background: var(--color-bg-card);
  border-radius: var(--radius-card);
  box-shadow: var(--shadow-card);
  padding: var(--space-lg);
}

.chart-title {
  font-size: 16px;
  font-weight: 600;
  color: var(--color-text-title);
  margin-bottom: var(--space-md);
}

.chart-placeholder {
  height: 280px;
  width: 100%;
}
</style>

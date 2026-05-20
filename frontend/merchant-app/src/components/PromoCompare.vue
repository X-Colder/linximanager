<template>
  <!-- 促销方案对比卡片 -->
  <view class="promo-compare">
    <view
      v-for="(plan, index) in plans"
      :key="index"
      class="plan-card"
      :class="[
        index === recommendIndex ? 'plan-card--recommend' : '',
        selectedIndex === index ? 'plan-card--selected' : ''
      ]"
    >
      <!-- 推荐标记 -->
      <view v-if="index === recommendIndex" class="recommend-badge">
        <text>⭐ 推荐</text>
      </view>

      <view class="plan-header">
        <text class="plan-name">{{ plan.name }}</text>
        <text v-if="plan.label" class="plan-label">{{ plan.label }}</text>
      </view>

      <view class="plan-detail">
        <view class="detail-row">
          <text class="detail-key">折扣</text>
          <text class="detail-val">{{ discountText(plan.discount) }}</text>
        </view>
        <view class="detail-row">
          <text class="detail-key">预计售罄</text>
          <text class="detail-val">约{{ plan.estimatedDays }}天</text>
        </view>
        <view class="detail-row">
          <text class="detail-key">预计利润</text>
          <text class="detail-val" :class="plan.estimatedProfit < 0 ? 'text--danger' : 'text--success'">
            {{ plan.estimatedProfit >= 0 ? '+' : '' }}¥{{ plan.estimatedProfit }}
          </text>
        </view>
      </view>

      <view v-if="plan.reason" class="plan-reason">
        <text>{{ plan.reason }}</text>
      </view>

      <button
        class="select-btn"
        :class="index === recommendIndex ? 'select-btn--primary' : 'select-btn--outline'"
        @tap="handleSelect(index)"
      >
        {{ index === recommendIndex ? '✅ 开始促销' : '选择此方案' }}
      </button>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref } from 'vue'

export interface PromoPlan {
  name: string
  label?: string
  discount: number
  originalPrice?: number
  promoPrice?: number
  estimatedDays: number
  estimatedProfit: number
  reason?: string
}

const props = defineProps<{
  plans: PromoPlan[]
  recommendIndex?: number
}>()

const emit = defineEmits<{
  select: [plan: PromoPlan, index: number]
}>()

const selectedIndex = ref<number | null>(null)

function discountText(discount: number): string {
  const percent = Math.round(discount * 10)
  return `${percent}折`
}

function handleSelect(index: number) {
  selectedIndex.value = index
  emit('select', props.plans[index], index)
}
</script>

<style scoped>
.promo-compare { display: flex; flex-direction: column; gap: 20rpx; }

.plan-card {
  background: #fff;
  border-radius: 24rpx;
  padding: 28rpx;
  border: 2rpx solid #e5e7eb;
  position: relative;
  transition: all 0.2s;
}

.plan-card--recommend {
  border-color: #f59e0b;
  box-shadow: 0 4rpx 16rpx rgba(245,158,11,0.15);
}

.plan-card--selected {
  border-color: #4a6cf7;
  box-shadow: 0 4rpx 16rpx rgba(74,108,247,0.15);
}

.recommend-badge {
  position: absolute;
  top: -1rpx;
  right: 20rpx;
  background: #f59e0b;
  color: #fff;
  font-size: 22rpx;
  padding: 4rpx 16rpx;
  border-radius: 0 0 12rpx 12rpx;
}

.plan-header {
  display: flex;
  align-items: center;
  gap: 12rpx;
  margin-bottom: 20rpx;
}

.plan-name { font-size: 30rpx; font-weight: 700; color: #1f2937; }

.plan-label {
  font-size: 22rpx;
  color: #4a6cf7;
  background: #e8edff;
  padding: 4rpx 12rpx;
  border-radius: 100rpx;
}

.plan-detail { background: #f9fafb; border-radius: 16rpx; padding: 20rpx; margin-bottom: 16rpx; }

.detail-row {
  display: flex;
  justify-content: space-between;
  margin-bottom: 12rpx;
}

.detail-row:last-child { margin-bottom: 0; }

.detail-key { font-size: 26rpx; color: #9ca3af; }
.detail-val { font-size: 26rpx; color: #1f2937; font-weight: 500; }
.text--success { color: #10b981; }
.text--danger  { color: #ef4444; }

.plan-reason {
  font-size: 24rpx;
  color: #4b5563;
  margin-bottom: 20rpx;
  padding: 12rpx 16rpx;
  background: #fffbeb;
  border-radius: 12rpx;
}

.select-btn {
  width: 100%;
  height: 80rpx;
  border-radius: 16rpx;
  font-size: 28rpx;
  display: flex;
  align-items: center;
  justify-content: center;
}

.select-btn--primary {
  background: #4a6cf7;
  color: #fff;
  border: none;
}

.select-btn--outline {
  background: transparent;
  color: #4a6cf7;
  border: 2rpx solid #4a6cf7;
}
</style>

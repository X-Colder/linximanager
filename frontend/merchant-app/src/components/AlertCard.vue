<template>
  <!-- 预警卡片：支持红/黄/蓝三种级别 -->
  <view class="alert-card" :class="`alert-card--${level}`" @tap="handleAction">
    <view class="alert-card__content">
      <view class="alert-card__icon-wrap">
        <text class="alert-card__level-icon">{{ LEVEL_ICON[level] }}</text>
      </view>
      <view class="alert-card__body">
        <view class="alert-card__title">{{ title }}</view>
        <view class="alert-card__subtitle">{{ subtitle }}</view>
        <view class="alert-card__desc">{{ description }}</view>
      </view>
      <view class="alert-card__action">
        <text class="action-text">{{ actionText }}</text>
        <text class="action-arrow">›</text>
      </view>
    </view>
  </view>
</template>

<script setup lang="ts">
defineProps<{
  level: 'red' | 'yellow' | 'blue'
  title: string
  subtitle: string
  description: string
  actionText: string
}>()

const emit = defineEmits<{
  action: []
}>()

const LEVEL_ICON = { red: '🔴', yellow: '🟡', blue: '🔵' }

function handleAction() {
  emit('action')
}
</script>

<style scoped>
.alert-card {
  border-radius: 24rpx;
  margin-bottom: 20rpx;
  overflow: hidden;
  border-left: 8rpx solid transparent;
}

.alert-card--red {
  background: #fef2f2;
  border-left-color: #ef4444;
}

.alert-card--yellow {
  background: #fffbeb;
  border-left-color: #f59e0b;
}

.alert-card--blue {
  background: #eff6ff;
  border-left-color: #3b82f6;
}

.alert-card__content {
  display: flex;
  align-items: center;
  padding: 24rpx 28rpx;
  gap: 16rpx;
}

.alert-card__icon-wrap {
  flex-shrink: 0;
  font-size: 36rpx;
}

.alert-card__body {
  flex: 1;
  min-width: 0;
}

.alert-card__title {
  font-size: 28rpx;
  font-weight: 600;
  color: #1f2937;
}

.alert-card__subtitle {
  font-size: 24rpx;
  color: #4b5563;
  margin-top: 4rpx;
}

.alert-card__desc {
  font-size: 24rpx;
  color: #9ca3af;
  margin-top: 4rpx;
}

.alert-card__action {
  display: flex;
  align-items: center;
  gap: 4rpx;
  flex-shrink: 0;
}

.action-text {
  font-size: 24rpx;
  color: #4a6cf7;
}

.action-arrow {
  font-size: 32rpx;
  color: #4a6cf7;
  line-height: 1;
}
</style>

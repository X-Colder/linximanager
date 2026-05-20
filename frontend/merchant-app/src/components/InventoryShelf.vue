<template>
  <!-- 库存沙盘可视化组件：货架模式/列表模式切换 -->
  <view class="inventory-shelf">
    <!-- 切换按钮 -->
    <view class="mode-switcher">
      <view
        class="mode-btn"
        :class="{ active: currentMode === 'shelf' }"
        @tap="currentMode = 'shelf'"
      >
        <text>货架</text>
      </view>
      <view
        class="mode-btn"
        :class="{ active: currentMode === 'list' }"
        @tap="currentMode = 'list'"
      >
        <text>列表</text>
      </view>
    </view>

    <!-- 货架模式 -->
    <view v-if="currentMode === 'shelf'" class="shelf-grid">
      <view
        v-for="item in items"
        :key="item.name"
        class="shelf-item"
        :class="item.alertLevel ? `shelf-item--${item.alertLevel}` : ''"
        @tap="emit('itemClick', item)"
      >
        <view class="shelf-bar-wrap">
          <view
            class="shelf-bar-fill"
            :style="{ height: `${Math.min(fillPercent(item), 100)}%`, background: barColor(item) }"
          />
        </view>
        <view class="shelf-item__name">{{ item.name }}</view>
        <view class="shelf-item__percent" :class="item.alertLevel ? `text--${item.alertLevel}` : ''">
          {{ fillPercent(item) }}%
        </view>
        <view class="shelf-item__stock">{{ item.currentStock }}{{ item.unit }}</view>
      </view>
    </view>

    <!-- 列表模式 -->
    <view v-else class="list-mode">
      <view
        v-for="item in items"
        :key="item.name"
        class="list-item"
        @tap="emit('itemClick', item)"
      >
        <view class="list-item__left">
          <text class="list-item__name">{{ item.name }}</text>
          <view class="list-item__progress">
            <view class="progress-track">
              <view
                class="progress-fill"
                :style="{ width: `${Math.min(fillPercent(item), 100)}%`, background: barColor(item) }"
              />
            </view>
            <text class="progress-text">{{ item.currentStock }}/{{ item.maxStock }}{{ item.unit }}</text>
          </view>
        </view>
        <view v-if="item.alertLevel" class="list-item__badge" :class="`badge--${item.alertLevel}`">
          {{ ALERT_LABEL[item.alertLevel] }}
        </view>
      </view>
    </view>

    <!-- 图例 -->
    <view class="legend">
      <view class="legend-item">
        <view class="legend-dot" style="background: #4a6cf7" />
        <text>充足</text>
      </view>
      <view class="legend-item">
        <view class="legend-dot" style="background: #f59e0b" />
        <text>积压</text>
      </view>
      <view class="legend-item">
        <view class="legend-dot" style="background: #ef4444" />
        <text>缺货</text>
      </view>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref } from 'vue'

export interface InventoryItem {
  name: string
  currentStock: number
  maxStock: number
  safetyStock: number
  unit: string
  alertLevel?: 'red' | 'yellow' | 'blue'
}

defineProps<{
  items: InventoryItem[]
}>()

const emit = defineEmits<{
  itemClick: [item: InventoryItem]
}>()

const ALERT_LABEL = { red: '缺货', yellow: '积压', blue: '临期' }

const currentMode = ref<'shelf' | 'list'>('shelf')

function fillPercent(item: InventoryItem): number {
  if (item.maxStock <= 0) return 0
  return Math.round((item.currentStock / item.maxStock) * 100)
}

function barColor(item: InventoryItem): string {
  if (item.alertLevel === 'red') return '#ef4444'
  if (item.alertLevel === 'yellow') return '#f59e0b'
  if (item.alertLevel === 'blue') return '#3b82f6'
  return '#4a6cf7'
}
</script>

<style scoped>
.inventory-shelf { background: #fff; border-radius: 24rpx; padding: 24rpx; }

.mode-switcher {
  display: flex;
  background: #f3f4f6;
  border-radius: 12rpx;
  padding: 6rpx;
  margin-bottom: 24rpx;
  width: 180rpx;
}

.mode-btn {
  flex: 1;
  text-align: center;
  padding: 8rpx 0;
  border-radius: 8rpx;
  font-size: 24rpx;
  color: #9ca3af;
  transition: all 0.2s;
}

.mode-btn.active {
  background: #fff;
  color: #4a6cf7;
  font-weight: 600;
  box-shadow: 0 2rpx 8rpx rgba(0,0,0,0.08);
}

.shelf-grid {
  display: flex;
  flex-wrap: wrap;
  gap: 16rpx;
}

.shelf-item {
  width: calc((100% - 32rpx) / 3);
  background: #f9fafb;
  border-radius: 16rpx;
  padding: 16rpx 12rpx;
  text-align: center;
  border: 2rpx solid transparent;
}

.shelf-item--red    { border-color: #ef4444; }
.shelf-item--yellow { border-color: #f59e0b; }
.shelf-item--blue   { border-color: #3b82f6; }

.shelf-bar-wrap {
  height: 100rpx;
  background: #e5e7eb;
  border-radius: 8rpx;
  margin-bottom: 10rpx;
  position: relative;
  overflow: hidden;
  display: flex;
  align-items: flex-end;
}

.shelf-bar-fill {
  width: 100%;
  border-radius: 8rpx;
  transition: height 0.3s ease;
}

.shelf-item__name {
  font-size: 24rpx;
  color: #1f2937;
  font-weight: 600;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.shelf-item__percent {
  font-size: 22rpx;
  color: #4a6cf7;
  margin-top: 4rpx;
}

.shelf-item__stock {
  font-size: 20rpx;
  color: #9ca3af;
}

.text--red    { color: #ef4444; }
.text--yellow { color: #f59e0b; }
.text--blue   { color: #3b82f6; }

/* 列表模式 */
.list-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 20rpx 0;
  border-bottom: 1rpx solid #f3f4f6;
}

.list-item:last-child { border-bottom: none; }

.list-item__left { flex: 1; }

.list-item__name {
  font-size: 28rpx;
  color: #1f2937;
  font-weight: 500;
}

.list-item__progress {
  display: flex;
  align-items: center;
  gap: 12rpx;
  margin-top: 8rpx;
}

.progress-track {
  flex: 1;
  height: 12rpx;
  background: #e5e7eb;
  border-radius: 6rpx;
  overflow: hidden;
}

.progress-fill {
  height: 100%;
  border-radius: 6rpx;
  transition: width 0.3s;
}

.progress-text { font-size: 22rpx; color: #9ca3af; white-space: nowrap; }

.list-item__badge {
  font-size: 22rpx;
  padding: 4rpx 14rpx;
  border-radius: 100rpx;
  margin-left: 12rpx;
}

.badge--red    { background: #fef2f2; color: #ef4444; }
.badge--yellow { background: #fffbeb; color: #f59e0b; }
.badge--blue   { background: #eff6ff; color: #3b82f6; }

/* 图例 */
.legend {
  display: flex;
  gap: 24rpx;
  margin-top: 20rpx;
  padding-top: 16rpx;
  border-top: 1rpx solid #f3f4f6;
}

.legend-item { display: flex; align-items: center; gap: 8rpx; font-size: 22rpx; color: #9ca3af; }

.legend-dot { width: 16rpx; height: 16rpx; border-radius: 50%; }
</style>

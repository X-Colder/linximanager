<template>
  <view class="page">
    <!-- 分类筛选 -->
    <scroll-view scroll-x class="category-bar">
      <view
        v-for="cat in categories"
        :key="cat.value"
        class="cat-item"
        :class="{ active: currentCategory === cat.value }"
        @tap="currentCategory = cat.value"
      >
        {{ cat.label }}
      </view>
    </scroll-view>

    <!-- 库存沙盘 -->
    <view class="section-card">
      <view class="section-title">库存沙盘</view>
      <InventoryShelf
        :items="filteredItems"
        @item-click="onItemClick"
      />
    </view>

    <!-- 新增按钮 -->
    <view class="fab" @tap="addProduct">
      <text class="fab__icon">+</text>
      <text class="fab__text">新增商品</text>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import InventoryShelf, { type InventoryItem } from '@/components/InventoryShelf.vue'
import http from '@/api/request'

const categories = [
  { label: '全部', value: '' },
  { label: '肉类', value: 'meat' },
  { label: '蔬菜', value: 'vegetable' },
  { label: '酒水', value: 'drink' },
  { label: '调料', value: 'seasoning' }
]

const currentCategory = ref('')
const allItems = ref<Array<InventoryItem & { category: string }>>([])
const loading = ref(false)

const filteredItems = computed(() => {
  if (!currentCategory.value) return allItems.value
  return allItems.value.filter(i => i.category === currentCategory.value)
})

async function loadInventory() {
  loading.value = true
  try {
    const res = await http.get<Array<InventoryItem & { category: string }>>('/merchant/inventory')
    allItems.value = res
  } catch {
    // 使用 mock 数据
    allItems.value = [
      { name: '鸡翅', currentStock: 2, maxStock: 10, safetyStock: 2, unit: 'kg', alertLevel: 'red', category: 'meat' },
      { name: '牛肉', currentStock: 5, maxStock: 10, safetyStock: 3, unit: 'kg', category: 'meat' },
      { name: '羊肉', currentStock: 1, maxStock: 8, safetyStock: 2, unit: 'kg', alertLevel: 'red', category: 'meat' },
      { name: '青岛啤酒', currentStock: 48, maxStock: 50, safetyStock: 24, unit: '瓶', category: 'drink' },
      { name: '可乐', currentStock: 60, maxStock: 60, safetyStock: 24, unit: '瓶', alertLevel: 'yellow', category: 'drink' },
      { name: '果汁', currentStock: 45, maxStock: 50, safetyStock: 20, unit: '瓶', alertLevel: 'yellow', category: 'drink' },
      { name: '面包', currentStock: 12, maxStock: 40, safetyStock: 10, unit: '个', alertLevel: 'blue', category: 'vegetable' }
    ]
  } finally {
    loading.value = false
  }
}

function onItemClick(item: InventoryItem) {
  uni.showActionSheet({
    itemList: ['查看详情', '快速入库', '快速出库'],
    success: ({ tapIndex }) => {
      if (tapIndex === 0) {
        uni.navigateTo({ url: `/pages/product-detail/index?name=${item.name}` })
      }
    }
  })
}

function addProduct() {
  uni.navigateTo({ url: '/pages/product-edit/index' })
}

onMounted(() => loadInventory())
</script>

<style scoped>
.page { min-height: 100vh; background: #f9fafb; padding: 24rpx; }

.category-bar {
  white-space: nowrap;
  margin-bottom: 24rpx;
}

.cat-item {
  display: inline-block;
  padding: 12rpx 28rpx;
  margin-right: 12rpx;
  border-radius: 100rpx;
  font-size: 26rpx;
  color: #4b5563;
  background: #fff;
  border: 2rpx solid #e5e7eb;
  transition: all 0.2s;
}

.cat-item.active {
  background: #4a6cf7;
  color: #fff;
  border-color: #4a6cf7;
}

.section-card {
  background: #fff;
  border-radius: 24rpx;
  padding: 28rpx;
  box-shadow: 0 4rpx 16rpx rgba(0,0,0,0.04);
}

.section-title {
  font-size: 30rpx;
  font-weight: 700;
  color: #1f2937;
  margin-bottom: 20rpx;
}

.fab {
  position: fixed;
  bottom: 120rpx;
  right: 32rpx;
  background: #4a6cf7;
  border-radius: 100rpx;
  padding: 16rpx 28rpx;
  display: flex;
  align-items: center;
  gap: 8rpx;
  box-shadow: 0 8rpx 24rpx rgba(74,108,247,0.4);
}

.fab__icon { font-size: 36rpx; color: #fff; font-weight: 300; }
.fab__text { font-size: 26rpx; color: #fff; font-weight: 600; }
</style>

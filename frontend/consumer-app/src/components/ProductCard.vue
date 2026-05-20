<template>
  <!-- 商品卡片：横向布局，支持促销价/划线价/标签/低库存提示 -->
  <view class="product-card" @tap="emit('add')">
    <view class="product-card__image-wrap">
      <image
        class="product-card__image"
        :src="image || ''"
        mode="aspectFill"
        lazy-load
      />
      <view v-if="tag" class="product-card__tag" :class="`tag--${tagType}`">{{ tag }}</view>
    </view>

    <view class="product-card__info">
      <text class="product-card__name">{{ name }}</text>
      <view class="product-card__price-row">
        <text class="product-card__price">¥{{ price }}</text>
        <text v-if="originalPrice" class="product-card__original">¥{{ originalPrice }}</text>
      </view>
      <text v-if="stock !== undefined && stock <= 10" class="product-card__stock">
        仅剩{{ stock }}份
      </text>
    </view>

    <view class="product-card__add" @tap.stop="emit('add')">
      <text class="add-icon">+</text>
    </view>
  </view>
</template>

<script setup lang="ts">
import { computed } from 'vue'

const props = defineProps<{
  image?: string
  name: string
  price: number
  originalPrice?: number
  tag?: string
  stock?: number
}>()

const emit = defineEmits<{
  add: []
}>()

const tagType = computed(() => {
  if (!props.tag) return ''
  if (props.tag.includes('临期') || props.tag.includes('清仓')) return 'danger'
  if (props.tag.includes('特惠') || props.tag.includes('限时')) return 'promo'
  return 'info'
})
</script>

<style scoped>
.product-card {
  display: flex;
  align-items: center;
  padding: 20rpx 0;
  border-bottom: 1rpx solid #f3f4f6;
  gap: 20rpx;
}

.product-card:last-child { border-bottom: none; }

.product-card__image-wrap {
  position: relative;
  flex-shrink: 0;
}

.product-card__image {
  width: 140rpx;
  height: 140rpx;
  border-radius: 16rpx;
  background: #f3f4f6;
}

.product-card__tag {
  position: absolute;
  top: 0;
  left: 0;
  font-size: 20rpx;
  padding: 4rpx 10rpx;
  border-radius: 8rpx 0 8rpx 0;
  font-weight: 600;
}

.tag--promo  { background: #f97316; color: #fff; }
.tag--danger { background: #ef4444; color: #fff; }
.tag--info   { background: #3b82f6; color: #fff; }

.product-card__info {
  flex: 1;
  min-width: 0;
}

.product-card__name {
  display: block;
  font-size: 28rpx;
  color: #1f2937;
  font-weight: 500;
  line-height: 1.4;
  margin-bottom: 8rpx;
}

.product-card__price-row {
  display: flex;
  align-items: baseline;
  gap: 10rpx;
}

.product-card__price {
  font-size: 32rpx;
  font-weight: 700;
  color: #f97316;
}

.product-card__original {
  font-size: 24rpx;
  color: #9ca3af;
  text-decoration: line-through;
}

.product-card__stock {
  display: block;
  font-size: 22rpx;
  color: #ef4444;
  margin-top: 6rpx;
}

.product-card__add {
  width: 60rpx;
  height: 60rpx;
  background: #4a6cf7;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  box-shadow: 0 4rpx 12rpx rgba(74,108,247,0.3);
}

.add-icon { font-size: 36rpx; color: #fff; font-weight: 300; line-height: 1; }
</style>

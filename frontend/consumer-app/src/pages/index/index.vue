<template>
  <view class="page">
    <!-- 店铺头部 -->
    <view class="shop-banner">
      <view class="shop-header">
        <image class="shop-logo" :src="shop.logo_url || ''" mode="aspectFill" />
        <view class="shop-meta">
          <text class="shop-name">{{ shop.shop_name }}</text>
          <view class="shop-status-row">
            <view class="shop-status-dot" :class="shop.is_open ? 'dot--open' : 'dot--closed'" />
            <text class="shop-status-text">{{ shop.is_open ? '营业中' : '已打烊' }}</text>
            <text v-if="shop.distance" class="shop-distance">📍{{ shop.distance }}</text>
          </view>
        </view>
      </view>
      <view v-if="shop.announcement" class="shop-announcement">
        <text class="announcement-icon">📢</text>
        <text class="announcement-text">{{ shop.announcement }}</text>
      </view>
    </view>

    <!-- 搜索框 -->
    <view class="search-bar" @tap="goSearch">
      <text class="search-icon">🔍</text>
      <text class="search-placeholder">搜索本店商品</text>
    </view>

    <scroll-view scroll-y class="page-content">
      <!-- 优惠专区 -->
      <view v-if="promoList.length" class="section">
        <view class="section-header">
          <text class="section-title">优惠专区</text>
          <text class="section-more" @tap="showAllPromo">更多 ›</text>
        </view>
        <scroll-view scroll-x class="promo-scroll">
          <view
            v-for="product in promoList"
            :key="product.id"
            class="promo-item"
            @tap="goProductDetail(product.id)"
          >
            <image class="promo-img" :src="product.image_url || ''" mode="aspectFill" lazy-load />
            <text class="promo-name">{{ product.name }}</text>
            <text class="promo-price">¥{{ product.promo_price }}</text>
            <text class="promo-original">¥{{ product.sale_price }}</text>
            <view class="promo-tag">{{ promoTag(product) }}</view>
          </view>
        </scroll-view>
      </view>

      <!-- 分类导航 -->
      <view class="category-nav">
        <scroll-view scroll-x>
          <view class="category-list">
            <view
              v-for="cat in categories"
              :key="cat.id"
              class="category-item"
              :class="{ active: currentCategoryId === cat.id }"
              @tap="currentCategoryId = cat.id"
            >
              {{ cat.name }}
            </view>
          </view>
        </scroll-view>
      </view>

      <!-- 商品列表 -->
      <view class="product-list section">
        <ProductCard
          v-for="product in filteredProducts"
          :key="product.id"
          :image="product.image_url"
          :name="product.name"
          :price="product.sale_price"
          :original-price="product.original_price"
          :tag="product.tag"
          :stock="product.stock"
          @add="addToCart(product)"
        />
        <view v-if="!filteredProducts.length && !loading" class="empty-state">
          <text class="empty-text">该分类暂无商品</text>
        </view>
      </view>

      <!-- 附近推荐（联盟） -->
      <view v-if="allianceList.length" class="section">
        <text class="section-title">附近的人也在买</text>
        <scroll-view scroll-x class="alliance-scroll">
          <view
            v-for="item in allianceList"
            :key="item.merchant_id"
            class="alliance-card"
            @tap="goShop(item.merchant_id)"
          >
            <text class="alliance-shop">{{ item.shop_name }}</text>
            <text class="alliance-product">{{ item.product_name }}</text>
            <text class="alliance-price">¥{{ item.price }}</text>
            <view class="alliance-tag">{{ item.tag }}</view>
          </view>
        </scroll-view>
      </view>
    </scroll-view>

    <!-- 购物车浮窗 -->
    <view v-if="cartCount > 0" class="cart-float" @tap="goCart">
      <view class="cart-float-badge">{{ cartCount }}</view>
      <text class="cart-float-icon">🛒</text>
      <text class="cart-float-total">¥{{ cartTotal }}</text>
      <text class="cart-float-btn">去结算</text>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { onLoad } from '@dcloudio/uni-app'
import ProductCard from '@/components/ProductCard.vue'
import http from '@/api/request'

interface Shop {
  shop_name: string
  logo_url: string
  announcement: string
  is_open: boolean
  distance?: string
}

interface Product {
  id: number
  name: string
  image_url: string
  sale_price: number
  promo_price?: number
  original_price?: number
  category_id: number
  tag?: string
  stock?: number
}

interface PromoProduct extends Product {
  promo_type: string
}

interface Category {
  id: number
  name: string
}

interface CartItem extends Product {
  quantity: number
}

interface AllianceItem {
  merchant_id: number
  shop_name: string
  product_name: string
  price: number
  tag: string
}

const merchantId = ref<string>('')

const shop = ref<Shop>({ shop_name: '', logo_url: '', announcement: '', is_open: true })
const categories = ref<Category[]>([{ id: 0, name: '全部' }])
const allProducts = ref<Product[]>([])
const promoList = ref<PromoProduct[]>([])
const allianceList = ref<AllianceItem[]>([])
const currentCategoryId = ref<number>(0)
const loading = ref(false)
const cart = ref<CartItem[]>([])

const filteredProducts = computed(() => {
  if (currentCategoryId.value === 0) return allProducts.value
  return allProducts.value.filter(p => p.category_id === currentCategoryId.value)
})

const cartCount = computed(() => cart.value.reduce((sum, i) => sum + i.quantity, 0))
const cartTotal = computed(() => cart.value.reduce((sum, i) => sum + i.sale_price * i.quantity, 0).toFixed(1))

function promoTag(p: PromoProduct): string {
  if (p.promo_type === 'flash_sale') return '限时'
  if (p.promo_type === 'clearance') return '临期'
  return '特惠'
}

function addToCart(product: Product) {
  const existing = cart.value.find(i => i.id === product.id)
  if (existing) {
    existing.quantity++
  } else {
    cart.value.push({ ...product, quantity: 1 })
  }
  uni.showToast({ title: '已加入购物车', icon: 'success', duration: 800 })
  // 将购物车同步到本地存储
  uni.setStorageSync('cart', JSON.stringify(cart.value))
}

function goCart() {
  uni.switchTab({ url: '/pages/cart/index' })
}

function goSearch() {
  uni.navigateTo({ url: `/pages/search/index?merchantId=${merchantId.value}` })
}

function goProductDetail(id: number) {
  uni.navigateTo({ url: `/pages/product-detail/index?id=${id}&merchantId=${merchantId.value}` })
}

function goShop(id: number) {
  uni.navigateTo({ url: `/pages/index/index?merchantId=${id}` })
}

function showAllPromo() {
  uni.navigateTo({ url: `/pages/promotions/index?merchantId=${merchantId.value}` })
}

async function loadShopData() {
  if (!merchantId.value) return
  loading.value = true
  try {
    const [shopData, products, promotions, alliance] = await Promise.all([
      http.get<Shop>(`/shop/${merchantId.value}`),
      http.get<Product[]>(`/shop/${merchantId.value}/products`),
      http.get<PromoProduct[]>(`/shop/${merchantId.value}/promotions`),
      http.get<AllianceItem[]>(`/shop/${merchantId.value}/nearby`).catch(() => [])
    ])
    shop.value = shopData
    allProducts.value = products
    promoList.value = promotions
    allianceList.value = alliance

    // 提取分类
    const catIds = new Set(products.map(p => p.category_id))
    categories.value = [
      { id: 0, name: '全部' },
      ...Array.from(catIds).map(id => ({ id, name: `分类${id}` }))
    ]
  } catch {
    // mock 数据展示
    shop.value = { shop_name: '老王烧烤', logo_url: '', announcement: '周末特惠！多款酒水5折起', is_open: true, distance: '500m' }
    allProducts.value = [
      { id: 1, name: '烤鸡翅(5串)', image_url: '', sale_price: 15, category_id: 1, stock: 8 },
      { id: 2, name: '烤牛肉(3串)', image_url: '', sale_price: 25, category_id: 1 },
      { id: 3, name: '青岛啤酒', image_url: '', sale_price: 6, category_id: 2 }
    ]
    promoList.value = [
      { id: 4, name: '果汁500ml', image_url: '', sale_price: 4.5, promo_price: 3.15, category_id: 2, promo_type: 'clearance', tag: '7折' },
      { id: 5, name: '全麦面包', image_url: '', sale_price: 8, promo_price: 5, category_id: 3, promo_type: 'clearance', tag: '临期' }
    ]
  } finally {
    loading.value = false
  }
}

onLoad((options) => {
  merchantId.value = options?.merchantId || '1'
})

onMounted(() => {
  // 恢复购物车
  const saved = uni.getStorageSync('cart')
  if (saved) cart.value = JSON.parse(saved)
  loadShopData()
})
</script>

<style scoped>
.page { min-height: 100vh; background: #f9fafb; }

.shop-banner { background: #fff; padding: 24rpx 32rpx; }

.shop-header { display: flex; gap: 20rpx; align-items: flex-start; }

.shop-logo {
  width: 88rpx;
  height: 88rpx;
  border-radius: 16rpx;
  background: #f3f4f6;
  flex-shrink: 0;
}

.shop-name { display: block; font-size: 34rpx; font-weight: 700; color: #1f2937; }

.shop-status-row { display: flex; align-items: center; gap: 10rpx; margin-top: 8rpx; }

.shop-status-dot { width: 16rpx; height: 16rpx; border-radius: 50%; }
.dot--open   { background: #10b981; }
.dot--closed { background: #9ca3af; }

.shop-status-text { font-size: 24rpx; color: #4b5563; }
.shop-distance    { font-size: 24rpx; color: #9ca3af; }

.shop-announcement {
  display: flex;
  align-items: center;
  gap: 8rpx;
  margin-top: 16rpx;
  padding: 12rpx 16rpx;
  background: #fffbeb;
  border-radius: 12rpx;
}

.announcement-icon { font-size: 28rpx; }
.announcement-text { font-size: 24rpx; color: #4b5563; flex: 1; }

.search-bar {
  display: flex;
  align-items: center;
  gap: 12rpx;
  margin: 20rpx 24rpx;
  padding: 18rpx 24rpx;
  background: #fff;
  border-radius: 100rpx;
  box-shadow: 0 2rpx 8rpx rgba(0,0,0,0.06);
}

.search-icon       { font-size: 28rpx; }
.search-placeholder { font-size: 26rpx; color: #9ca3af; }

.page-content { height: calc(100vh - 300rpx); padding: 0 24rpx 160rpx; }

.section { margin-bottom: 24rpx; }

.section-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 16rpx; }
.section-title  { font-size: 30rpx; font-weight: 700; color: #1f2937; margin-bottom: 16rpx; }
.section-more   { font-size: 24rpx; color: #4a6cf7; }

/* 促销横滚 */
.promo-scroll { white-space: nowrap; }

.promo-item {
  display: inline-flex;
  flex-direction: column;
  align-items: center;
  width: 160rpx;
  margin-right: 16rpx;
  background: #fff;
  border-radius: 16rpx;
  padding: 16rpx 12rpx;
  box-shadow: 0 2rpx 8rpx rgba(0,0,0,0.04);
  vertical-align: top;
}

.promo-img {
  width: 120rpx;
  height: 120rpx;
  border-radius: 12rpx;
  background: #f3f4f6;
  margin-bottom: 8rpx;
}

.promo-name     { font-size: 24rpx; color: #1f2937; text-align: center; white-space: nowrap; overflow: hidden; text-overflow: ellipsis; max-width: 100%; }
.promo-price    { font-size: 28rpx; font-weight: 700; color: #f97316; }
.promo-original { font-size: 22rpx; color: #9ca3af; text-decoration: line-through; }
.promo-tag      { font-size: 20rpx; color: #f97316; background: #fff7ed; padding: 2rpx 10rpx; border-radius: 100rpx; margin-top: 4rpx; }

/* 分类 */
.category-nav { background: #fff; padding: 16rpx 0; margin-bottom: 16rpx; border-radius: 16rpx; }

.category-list { display: flex; padding: 0 16rpx; gap: 8rpx; }

.category-item {
  padding: 10rpx 28rpx;
  border-radius: 100rpx;
  font-size: 26rpx;
  color: #4b5563;
  background: #f3f4f6;
  white-space: nowrap;
}

.category-item.active { background: #4a6cf7; color: #fff; }

/* 商品列表 */
.product-list { background: #fff; border-radius: 20rpx; padding: 0 24rpx; }

.empty-state { padding: 40rpx 0; text-align: center; }
.empty-text  { font-size: 26rpx; color: #9ca3af; }

/* 联盟 */
.alliance-scroll { white-space: nowrap; }

.alliance-card {
  display: inline-flex;
  flex-direction: column;
  width: 240rpx;
  background: #fff;
  border-radius: 16rpx;
  padding: 20rpx;
  margin-right: 16rpx;
  box-shadow: 0 2rpx 8rpx rgba(0,0,0,0.04);
  vertical-align: top;
}

.alliance-shop    { font-size: 24rpx; color: #9ca3af; }
.alliance-product { font-size: 26rpx; font-weight: 600; color: #1f2937; margin: 6rpx 0; }
.alliance-price   { font-size: 28rpx; font-weight: 700; color: #f97316; }
.alliance-tag     { font-size: 20rpx; color: #f97316; background: #fff7ed; padding: 2rpx 12rpx; border-radius: 100rpx; margin-top: 8rpx; align-self: flex-start; }

/* 购物车浮窗 */
.cart-float {
  position: fixed;
  bottom: 120rpx;
  left: 32rpx;
  right: 32rpx;
  background: #1f2937;
  border-radius: 100rpx;
  height: 96rpx;
  display: flex;
  align-items: center;
  padding: 0 24rpx;
  box-shadow: 0 8rpx 24rpx rgba(0,0,0,0.2);
}

.cart-float-badge {
  position: absolute;
  top: -8rpx;
  left: 56rpx;
  background: #ef4444;
  color: #fff;
  font-size: 20rpx;
  padding: 2rpx 10rpx;
  border-radius: 100rpx;
  min-width: 32rpx;
  text-align: center;
}

.cart-float-icon  { font-size: 40rpx; }
.cart-float-total { flex: 1; font-size: 28rpx; font-weight: 700; color: #fff; margin-left: 16rpx; }
.cart-float-btn   { background: #4a6cf7; color: #fff; font-size: 26rpx; padding: 14rpx 28rpx; border-radius: 100rpx; }
</style>

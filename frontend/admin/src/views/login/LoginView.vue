<template>
  <div class="login-page">
    <div class="login-card">
      <div class="login-header">
        <div class="logo">
          <span class="logo-icon">&#9775;</span>
          <span class="logo-text">灵犀掌柜</span>
        </div>
        <p class="login-subtitle">智能库存决策平台 · 管理后台</p>
      </div>

      <el-form ref="formRef" :model="form" :rules="rules" class="login-form" @submit.prevent="handleLogin">
        <el-form-item prop="phone">
          <el-input
            v-model="form.phone"
            placeholder="请输入手机号"
            size="large"
            maxlength="11"
            :prefix-icon="Phone"
          />
        </el-form-item>

        <el-form-item prop="code">
          <div class="code-input-row">
            <el-input
              v-model="form.code"
              placeholder="请输入验证码"
              size="large"
              maxlength="6"
              :prefix-icon="Key"
            />
            <el-button
              class="send-code-btn"
              :disabled="sendingCode || countdown > 0"
              :loading="sendingCode"
              size="large"
              @click="handleSendCode"
            >
              {{ countdown > 0 ? `${countdown}s 后重发` : '获取验证码' }}
            </el-button>
          </div>
        </el-form-item>

        <el-form-item>
          <el-button
            type="primary"
            size="large"
            class="login-btn"
            :loading="loading"
            @click="handleLogin"
          >
            登录
          </el-button>
        </el-form-item>
      </el-form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { ElMessage, type FormInstance, type FormRules } from 'element-plus'
import { Phone, Key } from '@element-plus/icons-vue'
import { useAuthStore } from '@/stores/auth'
import { authApi } from '@/api/auth'

const router = useRouter()
const route = useRoute()
const auth = useAuthStore()

const formRef = ref<FormInstance>()
const loading = ref(false)
const sendingCode = ref(false)
const countdown = ref(0)

const form = reactive({ phone: '', code: '' })

const rules: FormRules = {
  phone: [
    { required: true, message: '请输入手机号', trigger: 'blur' },
    { pattern: /^1[3-9]\d{9}$/, message: '手机号格式不正确', trigger: 'blur' }
  ],
  code: [
    { required: true, message: '请输入验证码', trigger: 'blur' },
    { len: 6, message: '验证码为6位数字', trigger: 'blur' }
  ]
}

let timer: ReturnType<typeof setInterval> | null = null

function startCountdown() {
  countdown.value = 60
  timer = setInterval(() => {
    countdown.value--
    if (countdown.value <= 0 && timer) {
      clearInterval(timer)
      timer = null
    }
  }, 1000)
}

async function handleSendCode() {
  if (!form.phone || !/^1[3-9]\d{9}$/.test(form.phone)) {
    ElMessage.warning('请先输入正确的手机号')
    return
  }

  sendingCode.value = true
  try {
    await authApi.sendCode(form.phone)
    ElMessage.success('验证码已发送')
    startCountdown()
  } catch {
    // 错误已由拦截器处理
  } finally {
    sendingCode.value = false
  }
}

async function handleLogin() {
  const valid = await formRef.value?.validate().catch(() => false)
  if (!valid) return

  loading.value = true
  try {
    await auth.login(form)
    ElMessage.success('登录成功')
    const redirect = (route.query.redirect as string) || '/dashboard'
    router.push(redirect)
  } catch {
    // 错误已由拦截器处理
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.login-page {
  display: flex;
  align-items: center;
  justify-content: center;
  min-height: 100vh;
  background: linear-gradient(135deg, var(--color-primary-light) 0%, #f0f4ff 100%);
}

.login-card {
  background: var(--color-bg-card);
  border-radius: var(--radius-modal);
  box-shadow: var(--shadow-modal);
  padding: 48px 40px;
  width: 420px;
}

.login-header {
  text-align: center;
  margin-bottom: 40px;
}

.logo {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 10px;
  margin-bottom: 8px;
}

.logo-icon {
  font-size: 32px;
  color: var(--color-primary);
}

.logo-text {
  font-size: 24px;
  font-weight: 700;
  color: var(--color-text-title);
}

.login-subtitle {
  font-size: 14px;
  color: var(--color-text-hint);
}

.login-form {
  margin-top: 24px;
}

.code-input-row {
  display: flex;
  gap: 12px;
  width: 100%;
}

.code-input-row .el-input {
  flex: 1;
}

.send-code-btn {
  width: 130px;
  flex-shrink: 0;
}

.login-btn {
  width: 100%;
  height: 48px;
  font-size: 16px;
  background-color: var(--color-primary);
  border-color: var(--color-primary);
}

.login-btn:hover {
  background-color: var(--color-primary-dark);
  border-color: var(--color-primary-dark);
}
</style>

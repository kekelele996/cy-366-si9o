<template>
  <div class="recharge-page">
    <van-tabs v-model:active="tab">
      <van-tab title="时长包" name="packages">
        <div class="pkg-grid">
          <div class="pkg-card" v-for="pkg in packages" :key="pkg.id">
            <div class="pkg-name">{{ pkg.name }}</div>
            <div class="pkg-hours">{{ pkg.hours }} 小时</div>
            <div class="pkg-price">¥{{ pkg.price }}</div>
            <van-button size="small" round type="primary" @click="buy(pkg)">立即购买</van-button>
          </div>
        </div>
      </van-tab>
      <van-tab title="会员充值" name="recharge">
        <van-cell-group inset title="为会员充值余额">
          <van-field v-model="rechargeForm.user_id" type="number" label="会员ID" placeholder="输入会员ID" />
          <van-field v-model="rechargeForm.amount" type="number" label="金额" placeholder="输入充值金额" />
          <van-field v-model="rechargeForm.payment_method" label="支付方式" placeholder="balance/cash/wechat/alipay" />
        </van-cell-group>
        <div class="submit-btn"><van-button round block type="primary" @click="doRecharge">确认充值</van-button></div>
        <van-cell-group inset title="我的充值记录">
          <van-cell v-for="r in recharges" :key="r.id" :title="`¥${r.amount}`" :label="PAYMENT_METHOD_TEXT[r.payment_method] || r.payment_method" :value="formatTime(r.created_at)" />
        </van-cell-group>
      </van-tab>
      <van-tab title="我的订单" name="orders">
        <van-cell-group inset>
          <van-cell v-for="o in orders" :key="o.id" :title="o.package_name" :label="o.order_no" :value="`¥${o.amount}`" />
        </van-cell-group>
      </van-tab>
    </van-tabs>
  </div>
</template>

<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue'
import { showSuccessToast, showToast } from 'vant'
import { listActivePackages, buyPackage, recharge as rechargeApi, listMyRecharges, listMyOrders, type TimePackage, type Recharge, type PackageOrder } from '@/api/recharge'
import { PAYMENT_METHOD_TEXT } from '@/constants'
import { formatTime } from '@/utils/format'
import { useAuth } from '@/hooks/useAuth'

const { isStaffOrAdmin } = useAuth()
const tab = ref('packages')
const packages = ref<TimePackage[]>([])
const recharges = ref<Recharge[]>([])
const orders = ref<PackageOrder[]>([])
const rechargeForm = reactive({ user_id: '', amount: '', payment_method: 'cash' })

async function loadPackages() {
  packages.value = await listActivePackages()
}

async function loadMine() {
  const [r, o] = await Promise.all([listMyRecharges({ page: 1, page_size: 20 }), listMyOrders({ page: 1, page_size: 20 })])
  recharges.value = r.list
  orders.value = o.list
}

async function buy(pkg: TimePackage) {
  try {
    await buyPackage({ package_id: pkg.id, payment_method: 'balance' })
    showSuccessToast('购买成功')
    loadMine()
  } catch { /* 余额不足等错误由拦截器提示 */ }
}

async function doRecharge() {
  if (!isStaffOrAdmin.value) {
    showToast('仅店员/管理员可操作充值')
    return
  }
  const uid = Number(rechargeForm.user_id)
  const amount = Number(rechargeForm.amount)
  if (!uid || !amount) {
    showToast('请填写会员ID与金额')
    return
  }
  await rechargeApi({ user_id: uid, amount, payment_method: rechargeForm.payment_method || 'cash' })
  showSuccessToast('充值成功')
  rechargeForm.user_id = ''
  rechargeForm.amount = ''
  loadMine()
}

onMounted(() => {
  loadPackages()
  loadMine()
})
</script>

<style scoped>
.pkg-grid { display: grid; grid-template-columns: repeat(2, 1fr); gap: 12px; padding: 12px; }
.pkg-card { background: #fff; border-radius: 12px; padding: 16px; text-align: center; box-shadow: 0 2px 8px rgba(0,0,0,0.06); }
.pkg-name { font-size: 15px; font-weight: 600; }
.pkg-hours { color: #969799; margin: 6px 0; }
.pkg-price { color: #ee0a24; font-size: 18px; font-weight: 600; margin-bottom: 10px; }
.submit-btn { margin: 12px 16px; }
</style>

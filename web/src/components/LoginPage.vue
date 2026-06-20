<template>
  <div class="login-container">
    <div class="login-box">
      <div class="login-header">
        <h1 class="logo">🧠 知识库管理系统</h1>
        <p class="subtitle">登录您的账户</p>
      </div>
      
      <form @submit.prevent="handleLogin" class="login-form">
        <div class="form-group">
          <label for="username">用户名</label>
          <input 
            id="username"
            v-model="username" 
            type="text" 
            placeholder="请输入用户名"
            class="form-input"
            required
          />
        </div>
        
        <div class="form-group">
          <label for="password">密码</label>
          <input 
            id="password"
            v-model="password" 
            type="password" 
            placeholder="请输入密码"
            class="form-input"
            required
          />
        </div>
        
        <div class="form-group">
          <label for="captcha">验证码</label>
          <div class="captcha-row">
            <input 
              id="captcha"
              v-model="captcha" 
              type="text" 
              placeholder="请输入验证码"
              class="form-input captcha-input"
              maxlength="4"
              required
            />
            <canvas 
              ref="captchaCanvas" 
              class="captcha-img" 
              @click="generateCaptcha"
            ></canvas>
          </div>
        </div>
        
        <button type="submit" class="login-btn" :disabled="isLoading">
          <span v-if="isLoading" class="loading">⏳</span>
          <span v-else>登 录</span>
        </button>
        
        <p v-if="error" class="error-message">{{ error }}</p>
      </form>
      
      <div class="login-footer">
        <p>还没有账户？请联系管理员</p>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import { useRouter, useRoute } from 'vue-router';
import { authService } from '../services/authService';

const router = useRouter();
const route = useRoute();

const username = ref('');
const password = ref('');
const captcha = ref('');
const captchaAnswer = ref('');
const isLoading = ref(false);
const error = ref('');
const captchaCanvas = ref(null);

const generateCaptcha = () => {
  let num1 = Math.floor(Math.random() * 90) + 10;
  let num2 = Math.floor(Math.random() * 90) + 10;
  const operators = ['+', '-'];
  const operator = operators[Math.floor(Math.random() * operators.length)];
  
  let answer;
  let expression;
  
  if (operator === '+') {
    answer = num1 + num2;
    expression = `${num1} + ${num2}`;
  } else {
    if (num1 < num2) {
      [num1, num2] = [num2, num1];
    }
    answer = num1 - num2;
    expression = `${num1} - ${num2}`;
  }
  
  captchaAnswer.value = answer.toString();
  
  const canvas = captchaCanvas.value;
  if (!canvas) return;
  
  const ctx = canvas.getContext('2d');
  const width = 120;
  const height = 40;
  
  canvas.width = width;
  canvas.height = height;
  
  ctx.fillStyle = '#f8fafc';
  ctx.fillRect(0, 0, width, height);
  
  ctx.font = 'bold 20px Arial';
  ctx.fillStyle = '#374151';
  ctx.textAlign = 'center';
  ctx.textBaseline = 'middle';
  ctx.fillText(expression, width / 2, height / 2);
  
  for (let i = 0; i < 4; i++) {
    ctx.strokeStyle = `rgba(${Math.floor(Math.random() * 256)}, ${Math.floor(Math.random() * 256)}, ${Math.floor(Math.random() * 256)}, 0.3)`;
    ctx.lineWidth = 1;
    ctx.beginPath();
    ctx.moveTo(Math.random() * width, Math.random() * height);
    ctx.lineTo(Math.random() * width, Math.random() * height);
    ctx.stroke();
  }
  
  for (let i = 0; i < 20; i++) {
    ctx.fillStyle = `rgba(${Math.floor(Math.random() * 256)}, ${Math.floor(Math.random() * 256)}, ${Math.floor(Math.random() * 256)}, 0.3)`;
    ctx.beginPath();
    ctx.arc(Math.random() * width, Math.random() * height, 1, 0, Math.PI * 2);
    ctx.fill();
  }
};

const handleLogin = async () => {
  error.value = '';
  
  if (!username.value || !password.value) {
    error.value = '请输入用户名和密码';
    return;
  }
  
  if (!captcha.value) {
    error.value = '请输入验证码';
    return;
  }
  
  if (captcha.value !== captchaAnswer.value) {
    error.value = '验证码错误';
    generateCaptcha();
    captcha.value = '';
    return;
  }
  
  isLoading.value = true;
  
  try {
    const loginResult = await authService.login(username.value, password.value);
    const user = loginResult.user || loginResult;
    
    if (user) {
      const redirect = route.query.redirect || '/';
      router.push(redirect);
    }
  } catch (e) {
    error.value = e.message || '登录失败，请重试';
    generateCaptcha();
    captcha.value = '';
  } finally {
    isLoading.value = false;
  }
};

onMounted(() => {
  generateCaptcha();
});
</script>

<style scoped>
.login-container {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  padding: 20px;
}

.login-box {
  background: white;
  border-radius: 16px;
  box-shadow: 0 10px 40px rgba(0, 0, 0, 0.15);
  padding: 40px;
  width: 100%;
  max-width: 400px;
}

.login-header {
  text-align: center;
  margin-bottom: 30px;
}

.logo {
  font-size: 28px;
  font-weight: 700;
  color: #1f2937;
  margin: 0 0 8px 0;
}

.subtitle {
  color: #6b7280;
  margin: 0;
}

.login-form {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.form-group label {
  font-size: 14px;
  font-weight: 500;
  color: #374151;
}

.form-input {
  padding: 12px 16px;
  border: 2px solid #e5e7eb;
  border-radius: 8px;
  font-size: 16px;
  transition: border-color 0.2s, box-shadow 0.2s;
}

.form-input:focus {
  outline: none;
  border-color: #667eea;
  box-shadow: 0 0 0 3px rgba(102, 126, 234, 0.1);
}

.form-input::placeholder {
  color: #9ca3af;
}

.login-btn {
  padding: 14px 24px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  border: none;
  border-radius: 8px;
  font-size: 16px;
  font-weight: 600;
  cursor: pointer;
  transition: transform 0.2s, box-shadow 0.2s;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
}

.login-btn:hover:not(:disabled) {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(102, 126, 234, 0.3);
}

.login-btn:disabled {
  opacity: 0.7;
  cursor: not-allowed;
}

.loading {
  animation: spin 1s linear infinite;
}

@keyframes spin {
  from { transform: rotate(0deg); }
  to { transform: rotate(360deg); }
}

.error-message {
  color: #ef4444;
  font-size: 14px;
  text-align: center;
  margin: 0;
}

.login-footer {
  margin-top: 24px;
  text-align: center;
}

.login-footer p {
  color: #6b7280;
  font-size: 14px;
  margin: 0;
}

.captcha-row {
  display: flex;
  gap: 12px;
  align-items: center;
}

.captcha-input {
  flex: 1;
}

.captcha-img {
  width: 120px;
  height: 40px;
  border: 2px solid #e5e7eb;
  border-radius: 8px;
  cursor: pointer;
  user-select: none;
}

.captcha-img:hover {
  border-color: #667eea;
}
</style>

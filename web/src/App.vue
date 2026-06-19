<script setup>
import { ref } from 'vue';

const notification = ref({
  show: false,
  message: '',
  type: 'success'
});

const showNotification = (message, type = 'success') => {
  notification.value = { show: true, message, type };
  setTimeout(() => {
    notification.value.show = false;
  }, 3000);
};
</script>

<template>
  <div class="app-container">
    <div v-if="notification.show" :class="['notification', notification.type]">
      {{ notification.message }}
    </div>
    
    <router-view />
  </div>
</template>

<style>
* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

body {
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif;
  font-size: 14px;
  line-height: 1.5;
  color: #212529;
  background-color: #fff;
}

.notification {
  position: fixed;
  top: 20px;
  right: 20px;
  padding: 12px 24px;
  border-radius: 4px;
  color: white;
  font-weight: 500;
  z-index: 9999;
  animation: slideIn 0.3s ease;
}

.notification.success {
  background-color: #28a745;
}

.notification.error {
  background-color: #dc3545;
}

@keyframes slideIn {
  from {
    transform: translateX(100%);
    opacity: 0;
  }
  to {
    transform: translateX(0);
    opacity: 1;
  }
}

/* 行内代码样式 */
code {
  font-family: monospace;
  background-color: #f4f4f4;
  padding: 2px 6px;
  border-radius: 4px;
  font-size: 0.9em;
  color: #333;
}

/* 代码块样式 */
pre {
  font-family: monospace;
  background-color: #f8f9fa;
  color: #333;
  padding: 12px 16px;
  border-radius: 8px;
  overflow-x: auto;
  line-height: 1.5;
  font-size: 13px;
  border: 1px solid #e9ecef;
  white-space: pre-wrap;
  word-break: break-all;
}

pre code {
  background-color: transparent;
  padding: 0;
  color: inherit;
}

/* 富文本编辑器中的代码样式 */
.editor-wrapper pre {
  background-color: #2d2d2d !important;
  color: #ccc !important;
  border: none !important;
}

.editor-wrapper code {
  background-color: transparent !important;
  color: inherit !important;
}
</style>

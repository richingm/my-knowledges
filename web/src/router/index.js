import { createRouter, createWebHistory } from 'vue-router';
import LoginPage from '../components/LoginPage.vue';
import KnowledgeBase from '../components/KnowledgeBase.vue';
import { authService } from '../services/authService';

const routes = [
  {
    path: '/login',
    name: 'login',
    component: LoginPage,
    meta: { public: true }
  },
  {
    path: '/',
    redirect: '/knowledges'
  },
  {
    path: '/knowledges',
    name: 'knowledges',
    component: KnowledgeBase,
    meta: { requiresAuth: true }
  },
  {
    path: '/:pathMatch(.*)*',
    redirect: '/knowledges'
  }
];

const router = createRouter({
  history: createWebHistory(),
  routes
});

router.beforeEach(async (to, from, next) => {
  // 检查是否需要登录
  if (to.meta.requiresAuth) {
    const isLoggedIn = authService.isLoggedIn();
    if (!isLoggedIn) {
      // 如果没有登录，跳转到登录页
      return next({ name: 'login', query: { redirect: to.fullPath } });
    }
    
    // 检查token是否有效
    try {
      const user = await authService.getCurrentUser();
      if (!user) {
        // 用户信息获取失败，可能token已过期
        return next({ name: 'login', query: { redirect: to.fullPath } });
      }
      next();
    } catch (error) {
      console.error('路由守卫检查失败:', error);
      return next({ name: 'login', query: { redirect: to.fullPath } });
    }
  } else {
    // 公开路由直接放行
    next();
  }
});

export default router;

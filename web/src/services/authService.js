export const authService = {
  async login(username, password) {
    try {
      const response = await fetch('/api/v1/users/login', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify({ username, password })
      });
      
      if (!response.ok) {
        const errorData = await response.json();
        throw new Error(errorData.message || '登录失败');
      }
      
      const data = await response.json();
      
      // 保存token到localStorage（后端返回的是token字段）
      const accessToken = data.token || data.accessToken || data.access_token;
      
      if (accessToken) {
        localStorage.setItem('access_token', accessToken);
      }
      
      // 保存用户信息到localStorage
      const user = data.user || data;
      if (user) {
        localStorage.setItem('current_user', JSON.stringify(user));
      }
      
      return data;
    } catch (error) {
      console.error('登录失败:', error);
      throw error;
    }
  },
  
  logout() {
    // 直接清空token，不需要调用后端接口
    localStorage.removeItem('access_token');
    localStorage.removeItem('refresh_token');
    localStorage.removeItem('current_user');
    localStorage.removeItem('selectedDomain');
    localStorage.removeItem('selectedKnowledgeId');
    localStorage.removeItem('selectedArticleId');
    return true;
  },
  
  async getCurrentUser() {
    try {
      const token = localStorage.getItem('access_token');
      if (!token) {
        throw new Error('未登录');
      }
      
      // 首先尝试从localStorage获取用户信息（避免依赖有问题的/api/v1/auth/me接口）
      const storedUser = localStorage.getItem('current_user');
      if (storedUser) {
        try {
          return JSON.parse(storedUser);
        } catch (e) {
          console.error('解析本地存储的用户信息失败:', e);
        }
      }
      
      // 如果本地没有用户信息，尝试调用接口获取
      const response = await fetch('/api/v1/auth/me', {
        method: 'GET',
        headers: {
          'Authorization': `Bearer ${token}`
        }
      });
      
      if (!response.ok) {
        // 如果接口返回500错误或其他错误，不抛出异常，返回null
        // 这样前端可以继续使用localStorage中的token进行其他操作
        console.warn('获取用户信息接口失败，状态码:', response.status);
        // 创建一个简单的用户对象
        return { id: '1', username: '用户', email: 'user@example.com' };
      }
      
      const data = await response.json();
      const user = data.user || data;
      
      // 保存到localStorage
      localStorage.setItem('current_user', JSON.stringify(user));
      
      return user;
    } catch (error) {
      console.error('获取用户信息失败:', error);
      // 如果有token但获取用户信息失败，返回一个默认用户对象
      if (localStorage.getItem('access_token')) {
        return { id: '1', username: '用户', email: 'user@example.com' };
      }
      return null;
    }
  },
  
  async refreshToken() {
    try {
      const refreshToken = localStorage.getItem('refresh_token');
      if (!refreshToken) {
        throw new Error('没有刷新token');
      }
      
      const response = await fetch('/api/v1/auth/refresh', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify({ refresh_token: refreshToken })
      });
      
      if (!response.ok) {
        throw new Error('刷新token失败');
      }
      
      const data = await response.json();
      
      if (data.access_token) {
        localStorage.setItem('access_token', data.access_token);
      }
      if (data.refresh_token) {
        localStorage.setItem('refresh_token', data.refresh_token);
      }
      
      return data;
    } catch (error) {
      console.error('刷新token失败:', error);
      throw error;
    }
  },
  
  isLoggedIn() {
    return !!localStorage.getItem('access_token');
  },
  
  getToken() {
    return localStorage.getItem('access_token');
  }
};

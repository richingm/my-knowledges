const getAuthHeaders = () => {
  const token = localStorage.getItem('access_token');
  return token ? { 'Authorization': `Bearer ${token}` } : {};
};

// 文件服务

/**
 * 上传文件
 * @param {File} file - 要上传的文件
 * @returns {Promise<Object>} 上传结果
 */
export const uploadFile = async (file) => {
  try {
    const formData = new FormData();
    formData.append('file', file);
    
    const response = await fetch('/api/v1/files/upload', {
      method: 'POST',
      headers: getAuthHeaders(),
      body: formData,
    });
    
    if (!response.ok) {
      throw new Error(`上传失败: ${response.statusText}`);
    }
    
    const data = await response.json();
    return data;
  } catch (error) {
    console.error('上传文件失败:', error);
    throw error;
  }
};

/**
 * 查看文件
 * @param {string} filename - 文件名
 * @returns {Promise<Response>} 文件响应
 */
export const viewFile = async (filename) => {
  try {
    const response = await fetch(`/api/v1/files/${filename}`, {
      headers: getAuthHeaders()
    });
    
    if (!response.ok) {
      throw new Error(`查看文件失败: ${response.statusText}`);
    }
    
    return response;
  } catch (error) {
    console.error('查看文件失败:', error);
    throw error;
  }
};

/**
 * 下载文件
 * @param {string} filename - 文件名
 * @returns {Promise<Response>} 文件响应
 */
export const downloadFile = async (filename) => {
  try {
    const response = await fetch(`/api/v1/files/${filename}/download`, {
      headers: getAuthHeaders()
    });
    
    if (!response.ok) {
      throw new Error(`下载文件失败: ${response.statusText}`);
    }
    
    return response;
  } catch (error) {
    console.error('下载文件失败:', error);
    throw error;
  }
};

export const fileService = {
  uploadFile,
  viewFile,
  downloadFile
};

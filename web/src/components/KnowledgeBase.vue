<script setup>
import { ref, onMounted, watch, onUnmounted, nextTick } from 'vue';
import { useRouter } from 'vue-router';
import { domainService } from '../services/domainService';
import { knowledgeService } from '../services/knowledgeService';
import { articleService } from '../services/articleService';
import { authService } from '../services/authService';
import KnowledgeTree from './KnowledgeTree.vue';
import ArticleTree from './ArticleTree.vue';
import BlockNoteEditor from './BlockNoteEditor.vue';


const router = useRouter();

const domains = ref([]);
const selectedDomain = ref(null);
const knowledgeTree = ref([]);
const articleTree = ref([]);
const selectedKnowledge = ref(null);
const selectedArticle = ref(null);
const articleDetail = ref(null);
const editableArticle = ref(null);
const mode = ref('default');
const viewMode = ref('tree');
const knowledgeCollapsed = ref(false);
const articleCollapsed = ref(false);

// 目录相关
const tocItems = ref([]);
const viewEditorRef = ref(null);
const editEditorRef = ref(null);

const handleTocUpdate = (items) => {
  tocItems.value = items || [];
};

const handleTocClick = (item) => {
  const ref = mode.value === 'view' ? viewEditorRef.value : editEditorRef.value;
  if (ref && ref.scrollToHeading) {
    ref.scrollToHeading(item.headingIndex);
  }
};

// 文章切换时清空目录
watch(articleDetail, (val) => {
  if (!val) tocItems.value = [];
});

const getAuthHeaders = () => {
  const token = localStorage.getItem('access_token');
  return token ? { 'Authorization': `Bearer ${token}` } : {};
};

const notification = ref({
  show: false,
  message: '',
  type: 'success'
});

const previewImage = ref({
  show: false,
  src: '',
  scale: 1
});

const leftColumnCollapsed = ref(false);
const middleColumnCollapsed = ref(false);

const showNotification = (message, type = 'success') => {
  notification.value = {
    show: true,
    message,
    type
  };
  setTimeout(() => {
    notification.value.show = false;
  }, 3000);
};


const moveKnowledge = async (knowledgeId, newParentId) => {
  const result = await knowledgeService.moveKnowledge(knowledgeId, newParentId, selectedDomain.value?.id);
  if (result && result.success) {
    if (selectedDomain.value?.id) {
      await fetchKnowledgeTree(selectedDomain.value.id);
    }
    showNotification('知识库移动成功!', 'success');
  } else {
    showNotification('知识库移动失败!', 'error');
  }
};

const generateKnowledgeOptions = (nodes, level = 0, excludeId = null) => {
  let options = [];
  const indent = '  '.repeat(level);
  nodes.forEach(node => {
    if (node.id !== excludeId) {
      options.push({
        value: node.id,
        label: `${indent}${node.name}`
      });
      if (node.children && node.children.length > 0) {
        options = options.concat(generateKnowledgeOptions(node.children, level + 1, excludeId));
      }
    }
  });
  return options;
};

const handleMoveKnowledge = (knowledgeId) => {
  const knowledgeOptions = generateKnowledgeOptions(knowledgeTree.value, 0, knowledgeId);
  const optionsHtml = knowledgeOptions.map(opt => 
    `<option value="${opt.value}">${opt.label}</option>`
  ).join('');
  
  const selectHtml = `
    <select id="knowledgeMoveSelect" style="width: 200px; padding: 5px; margin-bottom: 10px;">
      <option value="0">根目录</option>
      ${optionsHtml}
    </select>
  `;
  
  const container = document.createElement('div');
  container.innerHTML = selectHtml;
  const select = container.querySelector('select');
  
  container.style.position = 'fixed';
  container.style.top = '50%';
  container.style.left = '50%';
  container.style.transform = 'translate(-50%, -50%)';
  container.style.padding = '20px';
  container.style.backgroundColor = 'white';
  container.style.border = '1px solid #ccc';
  container.style.borderRadius = '8px';
  container.style.boxShadow = '0 2px 10px rgba(0,0,0,0.2)';
  container.style.zIndex = '1000';
  
  const message = document.createElement('p');
  message.textContent = '请选择目标位置:';
  message.style.marginBottom = '10px';
  container.insertBefore(message, select);
  
  const buttonContainer = document.createElement('div');
  buttonContainer.style.display = 'flex';
  buttonContainer.style.gap = '10px';
  buttonContainer.style.marginTop = '15px';
  
  const confirmBtn = document.createElement('button');
  confirmBtn.textContent = '确定';
  confirmBtn.style.padding = '5px 15px';
  confirmBtn.style.backgroundColor = '#007bff';
  confirmBtn.style.color = 'white';
  confirmBtn.style.border = 'none';
  confirmBtn.style.borderRadius = '4px';
  confirmBtn.style.cursor = 'pointer';
  
  const cancelBtn = document.createElement('button');
  cancelBtn.textContent = '取消';
  cancelBtn.style.padding = '5px 15px';
  cancelBtn.style.backgroundColor = '#6c757d';
  cancelBtn.style.color = 'white';
  cancelBtn.style.border = 'none';
  cancelBtn.style.borderRadius = '4px';
  cancelBtn.style.cursor = 'pointer';
  
  buttonContainer.appendChild(confirmBtn);
  buttonContainer.appendChild(cancelBtn);
  container.appendChild(buttonContainer);
  
  const handleConfirm = () => {
    const newParentId = parseInt(select.value) || 0;
    moveKnowledge(knowledgeId, newParentId);
    document.body.removeChild(container);
  };
  
  const handleCancel = () => {
    document.body.removeChild(container);
  };
  
  confirmBtn.addEventListener('click', handleConfirm);
  cancelBtn.addEventListener('click', handleCancel);
  
  document.body.appendChild(container);
};

const moveArticle = async (articleId, newParentId) => {
  const result = await articleService.moveArticle(articleId, newParentId, selectedKnowledge.value?.id);
  if (result && result.success) {
    if (selectedKnowledge.value?.id) {
      await fetchArticleTree(selectedKnowledge.value.id);
    }
    showNotification('文章移动成功!', 'success');
  } else {
    showNotification('文章移动失败!', 'error');
  }
};

const generateArticleOptions = (nodes, level = 0, excludeId = null) => {
  let options = [];
  const indent = '\u00A0\u00A0'.repeat(level);
  nodes.forEach(node => {
    const nodeId = parseInt(node.id) || node.id;
    const excludeIntId = excludeId !== null ? (parseInt(excludeId) || excludeId) : null;
    
    if (nodeId !== excludeIntId) {
      options.push({
        value: node.id,
        label: `${indent}${node.title}`
      });
      if (node.children && node.children.length > 0) {
        options = options.concat(generateArticleOptions(node.children, level + 1, excludeId));
      }
    }
  });
  return options;
};

const handleMoveArticle = (articleId) => {
  if (!articleTree.value || articleTree.value.length === 0) {
    showNotification('文章列表数据未加载!', 'error');
    return;
  }
  
  const articleOptions = generateArticleOptions(articleTree.value, 0, articleId);
  console.log('Article options generated:', articleOptions);
  
  const optionsHtml = articleOptions.map(opt => 
    `<option value="${opt.value}">${opt.label}</option>`
  ).join('');
  
  const selectHtml = `
    <select id="articleMoveSelect" style="width: 200px; padding: 5px; margin-bottom: 10px;">
      <option value="0">根目录</option>
      ${optionsHtml}
    </select>
  `;
  
  const container = document.createElement('div');
  container.innerHTML = selectHtml;
  const select = container.querySelector('select');
  
  container.style.position = 'fixed';
  container.style.top = '50%';
  container.style.left = '50%';
  container.style.transform = 'translate(-50%, -50%)';
  container.style.padding = '20px';
  container.style.backgroundColor = 'white';
  container.style.border = '1px solid #ccc';
  container.style.borderRadius = '8px';
  container.style.boxShadow = '0 2px 10px rgba(0,0,0,0.2)';
  container.style.zIndex = '1000';
  
  const message = document.createElement('p');
  message.textContent = '请选择目标位置:';
  message.style.marginBottom = '10px';
  container.insertBefore(message, select);
  
  const buttonContainer = document.createElement('div');
  buttonContainer.style.display = 'flex';
  buttonContainer.style.gap = '10px';
  buttonContainer.style.marginTop = '15px';
  
  const confirmBtn = document.createElement('button');
  confirmBtn.textContent = '确定';
  confirmBtn.style.padding = '5px 15px';
  confirmBtn.style.backgroundColor = '#007bff';
  confirmBtn.style.color = 'white';
  confirmBtn.style.border = 'none';
  confirmBtn.style.borderRadius = '4px';
  confirmBtn.style.cursor = 'pointer';
  
  const cancelBtn = document.createElement('button');
  cancelBtn.textContent = '取消';
  cancelBtn.style.padding = '5px 15px';
  cancelBtn.style.backgroundColor = '#6c757d';
  cancelBtn.style.color = 'white';
  cancelBtn.style.border = 'none';
  cancelBtn.style.borderRadius = '4px';
  cancelBtn.style.cursor = 'pointer';
  
  buttonContainer.appendChild(confirmBtn);
  buttonContainer.appendChild(cancelBtn);
  container.appendChild(buttonContainer);
  
  const handleConfirm = () => {
    const newParentId = parseInt(select.value) || 0;
    moveArticle(articleId, newParentId);
    document.body.removeChild(container);
  };
  
  const handleCancel = () => {
    document.body.removeChild(container);
  };
  
  confirmBtn.addEventListener('click', handleConfirm);
  cancelBtn.addEventListener('click', handleCancel);
  
  document.body.appendChild(container);
};

const fetchDomains = async () => {
  try {
    const result = await domainService.listDomains();
    if (result && Array.isArray(result)) {
      domains.value = result;
    }
  } catch (error) {
    console.error('获取领域列表失败:', error);
  }
};

const fetchKnowledgeTree = async (domainId) => {
  try {
    const result = await knowledgeService.getKnowledgeTree(domainId);
    if (result && Array.isArray(result)) {
      knowledgeTree.value = result;
    } else if (result && result.knowledge_tree) {
      knowledgeTree.value = result.knowledge_tree;
    } else if (result && result.items) {
      knowledgeTree.value = result.items;
    } else {
      knowledgeTree.value = [];
    }
  } catch (error) {
    console.error('获取知识库树失败:', error);
  }
};

const fetchArticleTree = async (knowledgeId) => {
  try {
    const result = await articleService.getArticleTree(knowledgeId);
    let articles = [];
    
    if (result && Array.isArray(result)) {
      articles = result;
    } else if (result && result.article_tree) {
      articles = result.article_tree;
    } else if (result && result.items) {
      articles = result.items;
    }
    
    articleTree.value = articles;
  } catch (error) {
    console.error('获取文章树失败:', error);
    articleTree.value = [];
  }
};

const fetchArticle = async (articleId) => {
  try {
    const result = await articleService.getArticle(articleId);
    if (result) {
      const importance = result.importance || result.Importance;
      const level = result.level || result.Level;
      const mappedImportance = importance || (level ? `level-${level}` : 'level-3');
      
      const article = {
        ...result,
        importance: mappedImportance
      };
      
      articleDetail.value = article;
      editableArticle.value = { ...article };
      mode.value = 'view';
    }
  } catch (error) {
    console.error('获取文章失败:', error);
  }
};

const handleDomainChange = async (event) => {
  const domainId = parseInt(event.target.value);
  const domain = domains.value.find(d => d.id && parseInt(d.id) === domainId);
  if (domain) {
    selectedDomain.value = domain;
    localStorage.setItem('selectedDomainId', domainId.toString());
    await fetchKnowledgeTree(domainId);
    selectedKnowledge.value = null;
    selectedArticle.value = null;
    articleDetail.value = null;
    editableArticle.value = null;
    articleTree.value = [];
  }
};

const handleKnowledgeClick = async (knowledge) => {
  selectedKnowledge.value = knowledge;
  localStorage.setItem('selectedKnowledgeId', knowledge.id.toString());
  
  knowledgeTree.value = toggleExpand(knowledgeTree.value, knowledge.id.toString());
  
  await fetchArticleTree(knowledge.id);
  selectedArticle.value = null;
  articleDetail.value = null;
  editableArticle.value = null;
};

const toggleExpand = (nodes, targetId) => {
  return nodes.map(node => {
    const isTarget = node.id && node.id.toString() === targetId;
    const newNode = {
      ...node,
      expanded: isTarget ? !node.expanded : node.expanded
    };
    if (node.children && node.children.length > 0) {
      newNode.children = toggleExpand(node.children, targetId);
    }
    return newNode;
  });
};

const handleArticleClick = async (article) => {
  selectedArticle.value = article;
  localStorage.setItem('selectedArticleId', article.id.toString());
  
  articleTree.value = toggleExpand(articleTree.value, article.id.toString());
  
  await fetchArticle(article.id);
};

const handleCreateKnowledge = async (parentId = null) => {
  const name = prompt('请输入知识库名称');
  if (!name) return;
  
  try {
    const result = await knowledgeService.createKnowledge({
      domain_id: selectedDomain.value?.id,
      parentKnowledgeId: parentId ? parseInt(parentId) : 0,
      name: name,
      description: ''
    });
    
    if (result) {
      await fetchKnowledgeTree(selectedDomain.value?.id);
      showNotification('知识库创建成功!', 'success');
    }
  } catch (error) {
    console.error('创建知识库失败:', error);
    showNotification('创建知识库失败!', 'error');
  }
};

const handleUpdateKnowledge = async (knowledgeId, name, description) => {
  try {
    const result = await knowledgeService.updateKnowledge(
      knowledgeId,
      selectedDomain.value?.id,
      name,
      description
    );
    
    if (result) {
      await fetchKnowledgeTree(selectedDomain.value?.id);
      showNotification('知识库更新成功!', 'success');
    }
  } catch (error) {
    console.error('更新知识库失败:', error);
    showNotification('更新知识库失败!', 'error');
  }
};

const handleDeleteKnowledge = async (knowledgeId) => {
  if (!confirm('确定删除此知识库吗？这将删除所有子知识库和文章！')) {
    return;
  }
  
  try {
    await knowledgeService.deleteKnowledge(knowledgeId);
    await fetchKnowledgeTree(selectedDomain.value?.id);
    
    if (selectedKnowledge.value?.id === knowledgeId) {
      selectedKnowledge.value = null;
      selectedArticle.value = null;
      articleDetail.value = null;
      editableArticle.value = null;
      articleTree.value = [];
    }
    
    showNotification('知识库删除成功!', 'success');
  } catch (error) {
    console.error('删除知识库失败:', error);
    showNotification('删除知识库失败!', 'error');
  }
};

const handleCreateArticle = async (parentArticleId = null) => {
  if (!selectedKnowledge.value) {
    showNotification('请先选择一个知识库!', 'error');
    return;
  }
  
  const title = prompt('请输入文章标题');
  if (!title) return;
  
  try {
    const result = await articleService.createArticle({
      knowledge_id: selectedKnowledge.value.id,
      title: title,
      content: '',
      parent_article_id: parentArticleId ? parseInt(parentArticleId) : 0
    });
    
    if (result && result.id) {
      await fetchArticleTree(selectedKnowledge.value.id);
      
      const newArticleId = result.id.toString();
      selectedArticle.value = findArticleById(articleTree.value, newArticleId);
      
      if (selectedArticle.value) {
        localStorage.setItem('selectedArticleId', newArticleId);
        await fetchArticle(selectedArticle.value.id);
      }
      
      showNotification('文章创建成功!', 'success');
    } else {
      showNotification('创建文章失败，未返回文章ID!', 'error');
    }
  } catch (error) {
    console.error('创建文章失败:', error);
    showNotification('创建文章失败!', 'error');
  }
};

const findArticleById = (nodes, id) => {
  for (const node of nodes) {
    if (node.id && node.id.toString() === id) {
      return node;
    }
    if (node.children && node.children.length > 0) {
      const found = findArticleById(node.children, id);
      if (found) return found;
    }
  }
  return null;
};

const handleUpdateArticle = async () => {
  if (!editableArticle.value) return;
  
  try {
    const content = editableArticle.value.content;
    const importance = editableArticle.value.importance || 'level-3';
    const level = parseInt(importance.replace('level-', ''));
    
    const result = await articleService.updateArticle(
      editableArticle.value.id,
      {
        knowledge_id: selectedKnowledge.value?.id,
        title: editableArticle.value.title,
        content: content,
        level: level
      }
    );
    
    if (result) {
      articleDetail.value = { ...editableArticle.value, content };
      mode.value = 'view';
      await fetchArticleTree(selectedKnowledge.value?.id);
      showNotification('文章更新成功!', 'success');
    }
  } catch (error) {
    console.error('更新文章失败:', error);
    showNotification('更新文章失败!', 'error');
  }
};

const handleDeleteArticle = async (articleId) => {
  if (!confirm('确定删除此文章吗？')) {
    return;
  }
  
  try {
    await articleService.deleteArticle(articleId);
    await fetchArticleTree(selectedKnowledge.value?.id);
    
    if (selectedArticle.value?.id === articleId) {
      selectedArticle.value = null;
      articleDetail.value = null;
      editableArticle.value = null;
    }
    
    showNotification('文章删除成功!', 'success');
  } catch (error) {
    console.error('删除文章失败:', error);
    showNotification('删除文章失败!', 'error');
  }
};

const startEdit = () => {
  mode.value = 'edit';
};

const cancelEdit = () => {
  mode.value = 'view';
  editableArticle.value = { ...articleDetail.value };
};

const getImportanceColor = (importance) => {
  const colorMap = {
    'level-5': '#dc3545',
    'level-4': '#6f42c1',
    'level-3': '#28a745',
    'level-2': '#000',
    'level-1': '#666666'
  };
  return colorMap[importance] || '#6c757d';
};

const getImportanceText = (importance) => {
  const textMap = {
    'level-5': '非常重要',
    'level-4': '比较重要',
    'level-3': '一般',
    'level-2': '不太重要',
    'level-1': '不重要'
  };
  return textMap[importance] || '未设置';
};

const openImagePreview = (src) => {
  previewImage.value = {
    show: true,
    src: src,
    scale: 1
  };
};

const closeImagePreview = () => {
  previewImage.value.show = false;
};

const zoomImage = (delta) => {
  previewImage.value.scale = Math.max(0.1, Math.min(3, previewImage.value.scale + delta));
};

const handleImageDoubleClick = (event) => {
  const target = event.target;
  if (target.tagName === 'IMG') {
    openImagePreview(target.src);
  }
};

const initData = async () => {
  await fetchDomains();
  
  const savedDomainId = localStorage.getItem('selectedDomainId');
  if (savedDomainId) {
    const domain = domains.value.find(d => d.id && d.id.toString() === savedDomainId);
    if (domain) {
      selectedDomain.value = domain;
      await fetchKnowledgeTree(parseInt(savedDomainId));
      
      const savedKnowledgeId = localStorage.getItem('selectedKnowledgeId');
      if (savedKnowledgeId) {
        const knowledge = findKnowledgeById(knowledgeTree.value, savedKnowledgeId);
        if (knowledge) {
          selectedKnowledge.value = knowledge;
          knowledgeTree.value = toggleExpand(knowledgeTree.value, savedKnowledgeId);
          await fetchArticleTree(parseInt(savedKnowledgeId));
          
          const savedArticleId = localStorage.getItem('selectedArticleId');
          if (savedArticleId) {
            const article = findArticleById(articleTree.value, savedArticleId);
            if (article) {
              selectedArticle.value = article;
              articleTree.value = toggleExpand(articleTree.value, savedArticleId);
              await fetchArticle(parseInt(savedArticleId));
            }
          }
        }
      }
    }
  } else if (domains.value.length > 0) {
    selectedDomain.value = domains.value[0];
    localStorage.setItem('selectedDomainId', domains.value[0].id.toString());
    await fetchKnowledgeTree(domains.value[0].id);
  }
};

const findKnowledgeById = (nodes, id) => {
  for (const node of nodes) {
    if (node.id && node.id.toString() === id) {
      return node;
    }
    if (node.children && node.children.length > 0) {
      const found = findKnowledgeById(node.children, id);
      if (found) return found;
    }
  }
  return null;
};

const handleLogout = async () => {
  await authService.logout();
  router.push('/login');
};

const currentUser = ref(null);

const loadCurrentUser = async () => {
  try {
    const user = await authService.getCurrentUser();
    if (user) {
      currentUser.value = user;
    }
  } catch (error) {
    console.error('加载用户信息失败:', error);
  }
};

onMounted(async () => {
  await loadCurrentUser();
  await initData();
  
  document.addEventListener('dblclick', handleImageDoubleClick);
});

onUnmounted(() => {
  document.removeEventListener('dblclick', handleImageDoubleClick);
});
</script>

<template>
  <div class="app-container">
    <div v-if="notification.show" :class="['notification', notification.type]">
      {{ notification.message }}
    </div>
    
    <header class="header">
      <div class="header-left">
        <select v-if="domains.length > 0" :value="selectedDomain?.id" @change="handleDomainChange" class="domain-select">
          <option v-for="domain in domains" :key="domain.id" :value="domain.id">
            {{ domain.name }}
          </option>
        </select>
      </div>
      <div class="header-right">
        
        <div class="user-info">
          <span class="username">{{ currentUser?.name || currentUser?.email || currentUser?.username }}</span>
          <button class="logout-btn" @click="handleLogout">🚪 退出</button>
        </div>
      </div>
    </header>
    
    <main class="main">
      <div v-if="viewMode === 'tree'" class="three-column-layout">
        <button 
          v-if="leftColumnCollapsed" 
          @click="leftColumnCollapsed = false" 
          class="expand-btn column-1-expand" 
          title="展开知识库"
        >
          ▶ 知识库
        </button>
        <div v-if="!leftColumnCollapsed" class="column column-1">
          <div class="column-header">
            <h2>知识库</h2>
            <div class="column-actions">
              <button @click="leftColumnCollapsed = true" class="collapse-btn" title="收缩">
                ◀
              </button>
              <button @click="handleCreateKnowledge()" class="create-btn">+ 新建</button>
            </div>
          </div>
          <KnowledgeTree 
            :tree-data="knowledgeTree" 
            :selected-id="selectedKnowledge?.id"
            @node-click="handleKnowledgeClick"
            @create-child="handleCreateKnowledge"
            @move-node="handleMoveKnowledge"
            @delete-node="handleDeleteKnowledge"
            @update-node="handleUpdateKnowledge"
          />
        </div>
        
        <button 
          v-if="middleColumnCollapsed" 
          @click="middleColumnCollapsed = false" 
          class="expand-btn column-2-expand" 
          title="展开文章"
        >
          ▶ 文章
        </button>
        <div v-if="!middleColumnCollapsed" class="column column-2">
          <div class="column-header">
            <h2>文章</h2>
            <div class="column-actions">
              <button @click="middleColumnCollapsed = true" class="collapse-btn" title="收缩">
                ◀
              </button>
              <button @click="handleCreateArticle()" class="create-btn" :disabled="!selectedKnowledge">+ 新建</button>
            </div>
          </div>
          <ArticleTree 
            :tree-data="articleTree" 
            :selected-id="selectedArticle?.id"
            @node-click="handleArticleClick"
            @create-child="handleCreateArticle"
            @move-node="handleMoveArticle"
            @delete-node="handleDeleteArticle"
          />
        </div>
        
        <div class="column column-3">
          <div v-if="articleDetail" class="article-detail">
            <div class="article-header">
              <div class="article-title-row">
                <h1 v-if="mode === 'view'">{{ articleDetail.title }}</h1>
                <input 
                  v-else 
                  v-model="editableArticle.title" 
                  class="title-input"
                />
              </div>
              <div class="article-meta">
                <div class="importance-selector">
                  <span class="importance-label">重要程度:</span>
                  <div v-if="mode === 'view'" class="importance-badge" :style="{ backgroundColor: getImportanceColor(articleDetail.importance) }">
                    {{ getImportanceText(articleDetail.importance) }}
                  </div>
                  <select v-else v-model="editableArticle.importance" class="importance-select">
                    <option value="level-5">非常重要</option>
                    <option value="level-4">比较重要</option>
                    <option value="level-3">一般</option>
                    <option value="level-2">不太重要</option>
                    <option value="level-1">不重要</option>
                  </select>
                </div>
                <div class="article-actions">
                  <button v-if="mode === 'view'" @click="startEdit" class="edit-btn">编辑</button>
                  <button v-if="mode === 'edit'" @click="handleUpdateArticle" class="save-btn">保存</button>
                  <button v-if="mode === 'edit'" @click="cancelEdit" class="cancel-btn">取消</button>
                  <button v-if="mode === 'view'" @click="handleDeleteArticle(articleDetail.id)" class="delete-btn">删除</button>
                </div>
              </div>
            </div>
            <div class="article-content">
              <div class="toc-sidebar" v-if="tocItems.length > 0">
                <div class="toc-header">目录</div>
                <div class="toc-list">
                  <div
                    v-for="item in tocItems"
                    :key="item.id"
                    class="toc-item"
                    :class="'toc-level-' + item.level"
                    :title="item.text"
                    @click="handleTocClick(item)"
                  >
                    {{ item.text }}
                  </div>
                </div>
              </div>
              <div class="editor-wrapper">
                <BlockNoteEditor
                  v-if="mode === 'view'"
                  ref="viewEditorRef"
                  v-model="articleDetail.content"
                  :readonly="true"
                  @tocUpdate="handleTocUpdate"
                />
                <BlockNoteEditor
                  v-else
                  ref="editEditorRef"
                  v-model="editableArticle.content"
                  @tocUpdate="handleTocUpdate"
                />
              </div>
            </div>
          </div>
          <div v-else class="empty-state">
            <div class="empty-icon">📄</div>
            <p>请选择一篇文章查看详情</p>
          </div>
        </div>
      </div>
    </main>
    
    <div v-if="previewImage.show" class="image-preview-overlay" @click="closeImagePreview">
      <div class="image-preview-container" @click.stop>
        <button class="close-btn" @click="closeImagePreview">&times;</button>
        <div class="image-controls">
          <button @click="zoomImage(-0.1)" class="zoom-btn">-</button>
          <span class="zoom-level">{{ Math.round(previewImage.scale * 100) }}%</span>
          <button @click="zoomImage(0.1)" class="zoom-btn">+</button>
        </div>
        <img 
          :src="previewImage.src" 
          alt="预览图片" 
          class="preview-image"
          :style="{ transform: `scale(${previewImage.scale})` }"
        />
      </div>
    </div>
  </div>
</template>

<style scoped>
.app-container {
  height: 100vh;
  display: flex;
  flex-direction: column;
}

.notification {
  position: fixed;
  top: 20px;
  right: 20px;
  padding: 12px 24px;
  border-radius: 4px;
  color: white;
  font-weight: 500;
  z-index: 1000;
  animation: fadeInOut 3s ease-in-out;
}

.notification.success {
  background-color: #28a745;
}

.notification.error {
  background-color: #dc3545;
}

@keyframes fadeInOut {
  0% { opacity: 0; transform: translateY(-20px); }
  10% { opacity: 1; transform: translateY(0); }
  90% { opacity: 1; transform: translateY(0); }
  100% { opacity: 0; transform: translateY(-20px); }
}

.header {
  padding: 0.5rem;
  background-color: #f8f9fa;
  border-bottom: 1px solid #e9ecef;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.header-left {
  display: flex;
  align-items: center;
}

.domain-select {
  padding: 0.25rem 0.5rem;
  border: 1px solid #ced4da;
  border-radius: 4px;
  font-size: 0.875rem;
}

.header-right {
  display: flex;
  align-items: center;
  gap: 1rem;
}

.user-info {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  padding-left: 1rem;
  border-left: 1px solid #e9ecef;
}

.username {
  font-size: 0.875rem;
  color: #495057;
}

.logout-btn {
  padding: 0.25rem 0.5rem;
  border: 1px solid #dc3545;
  border-radius: 4px;
  font-size: 0.875rem;
  cursor: pointer;
  background-color: white;
  color: #dc3545;
  transition: all 0.2s;
}

.logout-btn:hover {
  background-color: #dc3545;
  color: white;
}

.view-toggle {
  display: flex;
  gap: 0.25rem;
}

.view-btn {
  padding: 0.25rem 0.5rem;
  border: 1px solid #ced4da;
  border-radius: 4px;
  font-size: 0.875rem;
  cursor: pointer;
  background-color: white;
  transition: all 0.2s;
}

.view-btn:hover {
  background-color: #e9ecef;
}

.view-btn.active {
  background-color: #007bff;
  color: white;
  border-color: #007bff;
}

.main {
  flex: 1;
  overflow: hidden;
}

.three-column-layout {
  display: flex;
  height: 100%;
}

.column {
  display: flex;
  flex-direction: column;
  border-right: 1px solid #e9ecef;
  overflow: hidden;
}

.column-1 {
  flex: 2;
  min-width: 200px;
  max-width: 350px;
}

.column-2 {
  flex: 2;
  min-width: 200px;
  max-width: 350px;
}

.column-3 {
  flex: 4;
  min-width: 400px;
}

.column-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0.5rem;
  background-color: #f8f9fa;
  border-bottom: 1px solid #e9ecef;
}

.column-header h2 {
  margin: 0;
  font-size: 1rem;
  font-weight: 600;
}

.create-btn {
  padding: 0.25rem 0.5rem;
  border: 1px solid #28a745;
  border-radius: 4px;
  font-size: 0.875rem;
  cursor: pointer;
  background-color: white;
  color: #28a745;
  transition: all 0.2s;
}

.create-btn:hover:not(:disabled) {
  background-color: #28a745;
  color: white;
}

.create-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.column-actions {
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.collapse-btn {
  padding: 0.25rem 0.5rem;
  border: 1px solid #dc3545;
  border-radius: 4px;
  font-size: 0.75rem;
  cursor: pointer;
  background-color: white;
  color: #dc3545;
  transition: all 0.2s;
}

.collapse-btn:hover {
  background-color: #dc3545;
  color: white;
}

.expand-btn {
  padding: 0.5rem;
  border: 1px solid #007bff;
  border-radius: 4px;
  font-size: 0.875rem;
  cursor: pointer;
  background-color: white;
  color: #007bff;
  transition: all 0.2s;
  writing-mode: vertical-rl;
  text-orientation: mixed;
}

.expand-btn:hover {
  background-color: #007bff;
  color: white;
}

.column-1-expand {
  position: absolute;
  left: 0;
  top: 1rem;
  writing-mode: horizontal-tb;
}

.column-2-expand {
  position: absolute;
  left: 4px;
  top: 1rem;
  writing-mode: horizontal-tb;
}

.article-detail {
  display: flex;
  flex-direction: column;
  height: 100%;
}

.article-header {
  padding: 1rem;
  border-bottom: 1px solid #e9ecef;
  background-color: #f8f9fa;
}

.article-title-row {
  margin-bottom: 0.75rem;
}

.article-title-row h1 {
  margin: 0;
  font-size: 1.25rem;
}

.title-input {
  width: 100%;
  padding: 0.5rem;
  border: 1px solid #ced4da;
  border-radius: 4px;
  font-size: 1.25rem;
}

.article-meta {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.importance-selector {
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.importance-label {
  font-size: 0.875rem;
  color: #6c757d;
}

.importance-badge {
  padding: 0.25rem 0.5rem;
  border-radius: 4px;
  font-size: 0.875rem;
  color: white;
}

.importance-select {
  padding: 0.25rem 0.5rem;
  border: 1px solid #ced4da;
  border-radius: 4px;
  font-size: 0.875rem;
}

.article-actions {
  display: flex;
  gap: 0.5rem;
}

.edit-btn, .save-btn, .cancel-btn, .delete-btn {
  padding: 0.25rem 0.5rem;
  border: 1px solid;
  border-radius: 4px;
  font-size: 0.875rem;
  cursor: pointer;
  transition: all 0.2s;
}

.edit-btn {
  border-color: #007bff;
  color: #007bff;
  background-color: white;
}

.edit-btn:hover {
  background-color: #007bff;
  color: white;
}

.save-btn {
  border-color: #28a745;
  color: #28a745;
  background-color: white;
}

.save-btn:hover {
  background-color: #28a745;
  color: white;
}

.cancel-btn {
  border-color: #6c757d;
  color: #6c757d;
  background-color: white;
}

.cancel-btn:hover {
  background-color: #6c757d;
  color: white;
}

.delete-btn {
  border-color: #dc3545;
  color: #dc3545;
  background-color: white;
}

.delete-btn:hover {
  background-color: #dc3545;
  color: white;
}

.article-content {
  flex: 1;
  overflow: hidden;
  display: flex;
  padding: 0;
}

.toc-sidebar {
  width: 220px;
  min-width: 220px;
  max-height: 100%;
  overflow-y: auto;
  border-right: 1px solid #e9ecef;
  background-color: #f8f9fa;
  padding: 0.75rem 0;
}

.toc-header {
  font-size: 0.8rem;
  font-weight: 600;
  color: #6c757d;
  text-transform: uppercase;
  letter-spacing: 0.05em;
  padding: 0 0.75rem 0.5rem;
  border-bottom: 1px solid #e9ecef;
  margin-bottom: 0.25rem;
}

.toc-list {
  padding: 0.25rem 0;
}

.toc-item {
  padding: 0.25rem 0.75rem;
  font-size: 0.8rem;
  color: #495057;
  cursor: pointer;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  border-left: 2px solid transparent;
  transition: all 0.15s ease;
}

.toc-item:hover {
  background-color: #e9ecef;
  color: #0d6efd;
  border-left-color: #0d6efd;
}

.toc-level-1 { font-weight: 600; font-size: 0.85rem; }
.toc-level-2 { padding-left: 1.5rem; }
.toc-level-3 { padding-left: 2.25rem; font-size: 0.75rem; }
.toc-level-4 { padding-left: 3rem; font-size: 0.75rem; }
.toc-level-5 { padding-left: 3.75rem; font-size: 0.7rem; color: #6c757d; }
.toc-level-6 { padding-left: 4.5rem; font-size: 0.7rem; color: #6c757d; }

.editor-wrapper {
  flex: 1;
  height: 100%;
  overflow: hidden;
}

.empty-state {
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  height: 100%;
  color: #6c757d;
}

.empty-icon {
  font-size: 3rem;
  margin-bottom: 1rem;
}

.image-preview-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: rgba(0, 0, 0, 0.8);
  display: flex;
  justify-content: center;
  align-items: center;
  z-index: 1000;
  cursor: pointer;
}

.image-preview-container {
  position: relative;
  max-width: 90%;
  max-height: 90%;
  cursor: default;
}

.close-btn {
  position: absolute;
  top: -40px;
  right: 0;
  width: 30px;
  height: 30px;
  border: none;
  border-radius: 50%;
  background-color: rgba(255, 255, 255, 0.2);
  color: white;
  font-size: 1.5rem;
  cursor: pointer;
  display: flex;
  justify-content: center;
  align-items: center;
}

.close-btn:hover {
  background-color: rgba(255, 255, 255, 0.4);
}

.image-controls {
  position: absolute;
  bottom: -40px;
  left: 50%;
  transform: translateX(-50%);
  display: flex;
  align-items: center;
  gap: 1rem;
  background-color: rgba(255, 255, 255, 0.2);
  padding: 0.5rem 1rem;
  border-radius: 4px;
}

.zoom-btn {
  width: 30px;
  height: 30px;
  border: none;
  border-radius: 4px;
  background-color: rgba(255, 255, 255, 0.8);
  color: #666666;
  font-size: 1.25rem;
  cursor: pointer;
}

.zoom-level {
  color: white;
  font-size: 0.875rem;
  min-width: 50px;
  text-align: center;
}

.preview-image {
  max-width: 100%;
  max-height: 80vh;
  transition: transform 0.2s ease;
}
</style>

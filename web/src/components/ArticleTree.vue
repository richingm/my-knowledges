<script setup>
import { ref, watch, onMounted } from 'vue';

const props = defineProps({
  treeData: {
    type: Array,
    default: () => []
  },
  selectedId: {
    type: [Number, String],
    default: null
  }
});

const emit = defineEmits(['node-click', 'create-child', 'move-node', 'delete-node']);

const expandedNodes = ref(new Set());

// 递归查找并展开选中节点的所有父节点
const expandSelectedPath = (nodes, targetId) => {
  if (!nodes || nodes.length === 0) return false;
  
  const targetIntId = parseInt(targetId);
  
  for (const node of nodes) {
    const nodeIntId = parseInt(node.id);
    if (nodeIntId === targetIntId || node.id === targetId) {
      return true;
    }
    
    if (node.children && node.children.length > 0) {
      const found = expandSelectedPath(node.children, targetId);
      if (found) {
        expandedNodes.value.add(node.id);
        return true;
      }
    }
  }
  
  return false;
};

// 监听 selectedId 变化，展开选中节点的路径
watch(() => props.selectedId, (newId) => {
  expandedNodes.value.clear();
  if (newId) {
    expandSelectedPath(props.treeData, newId);
  }
}, { immediate: true });

// 监听 treeData 变化，重新展开选中节点的路径
watch(() => props.treeData, (newData) => {
  expandedNodes.value.clear();
  if (props.selectedId) {
    expandSelectedPath(newData, props.selectedId);
  }
}, { deep: true });

// 组件挂载时展开选中节点的路径
onMounted(() => {
  if (props.selectedId) {
    expandSelectedPath(props.treeData, props.selectedId);
  }
});

const toggleNode = (nodeId) => {
  if (expandedNodes.value.has(nodeId)) {
    expandedNodes.value.delete(nodeId);
  } else {
    expandedNodes.value.add(nodeId);
  }
};

const isExpanded = (nodeId) => {
  return expandedNodes.value.has(nodeId);
};

const hasChildren = (node) => {
  return node.children && node.children.length > 0;
};

const isSelected = (nodeId) => {
  if (props.selectedId === null || props.selectedId === undefined) return false;
  const nodeIntId = parseInt(nodeId);
  const selectedIntId = parseInt(props.selectedId);
  return nodeIntId === selectedIntId || nodeId === props.selectedId;
};

const handleNodeClick = (node) => {
  emit('node-click', node);
};
</script>

<template>
  <div class="article-tree">
    <div v-if="treeData.length === 0" class="tree-empty">
      No articles found
    </div>
    <div v-else class="tree-nodes">
      <div 
        v-for="node in treeData" 
        :key="node.id"
        class="tree-node"
      >
        <div class="node-content" @click="toggleNode(node.id)" :class="{ 'node-selected': isSelected(node.id) }">
          <span class="node-toggle" v-if="hasChildren(node)">
            {{ isExpanded(node.id) ? '▼' : '►' }}
          </span>
          <span class="node-name" :class="`level-${node.level || node.importance || '3'}`" @click.stop="handleNodeClick(node)">{{ node.title }}</span>
          <div class="node-actions">
            <button class="action-btn create-btn" @click.stop="emit('create-child', node.id)" title="新建子文章">+</button>
            <button class="action-btn move-btn" @click.stop="emit('move-node', node.id)" title="移动文章">↕</button>
            <button class="action-btn delete-btn" @click.stop="emit('delete-node', node.id)" title="删除文章">×</button>
          </div>
        </div>
        <div 
          v-if="hasChildren(node) && isExpanded(node.id)" 
          class="tree-children"
        >
          <ArticleTree :tree-data="node.children" :selected-id="selectedId" @node-click="handleNodeClick" @create-child="(id) => emit('create-child', id)" @move-node="(id) => emit('move-node', id)" @delete-node="(id) => emit('delete-node', id)" />
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.article-tree {
  width: 100%;
  height: 100%;
  overflow-y: auto;
  overflow-x: hidden;
  flex: 1;
}

.tree-empty {
  padding: 1rem;
  text-align: center;
  color: #6c757d;
}

.tree-nodes {
  padding: 0.5rem 0;
}

.tree-node {
  margin: 0.25rem 0;
}

.node-content {
  display: flex;
  align-items: flex-start;
  padding: 0.5rem;
  cursor: pointer;
  border-radius: 4px;
  transition: background-color 0.2s;
  min-height: 1.5rem;
  position: relative;
}

.node-content:hover {
  background-color: #e9ecef;
}

.node-content.node-selected {
  background-color: #e3f2fd;
  border-left: 3px solid #2196f3;
  font-weight: 600;
}

.node-toggle {
  width: 1rem;
  height: 1rem;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  margin-right: 0.5rem;
  font-size: 0.75rem;
  color: #6c757d;
}

.node-name {
  flex: 1;
  font-size: 0.875rem;
  word-wrap: break-word;
  white-space: normal;
  margin-right: 3rem; /* 为按钮预留空间 */
}

.node-name.level-5 {
  color: #dc3545;
}

.node-name.level-4 {
  color: #6f42c1;
}

.node-name.level-3 {
  color: #333;
}

.node-name.level-2 {
  color: #fd7e14;
}

.node-name.level-1 {
  color: #28a745;
}

.tree-children {
  margin-left: 1.5rem;
  border-left: 1px solid #e9ecef;
  padding-left: 0.5rem;
}

.node-actions {
  display: flex;
  gap: 0.25rem;
  opacity: 0;
  transition: opacity 0.2s;
  flex-shrink: 0;
  position: absolute;
  right: 0.5rem;
  top: 0.5rem;
}

.node-content:hover .node-actions {
  opacity: 1;
}

.action-btn {
  width: 1.25rem;
  height: 1.25rem;
  border: none;
  border-radius: 50%;
  font-size: 0.625rem;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.2s;
}

.move-btn {
  background-color: #ffc107;
  color: white;
}

.move-btn:hover {
  background-color: #e0a800;
  transform: scale(1.1);
}

.delete-btn {
  background-color: #dc3545;
  color: white;
}

.delete-btn:hover {
  background-color: #c82333;
  transform: scale(1.1);
}
</style>
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

const emit = defineEmits(['node-click', 'create-child', 'delete-node', 'move-node']);

const expandedNodes = ref(new Set());

// 递归展开所有节点
const expandAllNodes = (nodes) => {
  if (!nodes || nodes.length === 0) return;
  
  nodes.forEach(node => {
    expandedNodes.value.add(node.id);
    if (node.children && node.children.length > 0) {
      expandAllNodes(node.children);
    }
  });
};

// 监听 treeData 变化，自动展开所有节点
watch(() => props.treeData, (newData) => {
  expandedNodes.value.clear();
  expandAllNodes(newData);
}, { deep: true });

// 组件挂载时展开所有节点
onMounted(() => {
  expandAllNodes(props.treeData);
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

const handleNodeClick = (node) => {
  emit('node-click', node);
};
</script>

<template>
  <div class="knowledge-tree">
    <div v-if="treeData.length === 0" class="tree-empty">
      No knowledge base found
    </div>
    <div v-else class="tree-nodes">
      <div 
        v-for="node in treeData" 
        :key="node.id"
        class="tree-node"
      >
        <div class="node-content" @click="toggleNode(node.id)" :class="{ 'node-selected': selectedId === node.id }">
          <span class="node-toggle" v-if="hasChildren(node)">
            {{ isExpanded(node.id) ? '▼' : '►' }}
          </span>
          <span class="node-name" @click.stop="handleNodeClick(node)">{{ node.name }}</span>
          <div class="node-actions">
            <button class="action-btn create-btn" @click.stop="emit('create-child', node.id)" title="创建子知识库">+</button>
            <button class="action-btn move-btn" @click.stop="emit('move-node', node.id)" title="移动知识库">↕</button>
            <button class="action-btn delete-btn" @click.stop="emit('delete-node', node.id)" title="删除知识库">×</button>
          </div>
        </div>
        <div 
          v-if="hasChildren(node) && isExpanded(node.id)" 
          class="tree-children"
        >
          <KnowledgeTree :tree-data="node.children" :selected-id="selectedId" @node-click="handleNodeClick" @create-child="(id) => emit('create-child', id)" @delete-node="(id) => emit('delete-node', id)" @move-node="(id) => emit('move-node', id)" />
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.knowledge-tree {
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
  background-color: #d1ecf1;
  border-left: 3px solid #17a2b8;
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
  margin-right: 4rem; /* 为按钮预留空间 */
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

.create-btn {
  background-color: #28a745;
  color: white;
}

.create-btn:hover {
  background-color: #218838;
  transform: scale(1.1);
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

.tree-children {
  margin-left: 1.5rem;
  border-left: 1px solid #e9ecef;
  padding-left: 0.5rem;
}
</style>
<script setup>
import { ref, watch, onMounted, nextTick, provide, inject } from 'vue';

// 通过 provide/inject 让递归子组件实例共享同一份拖拽状态
const dragStateKey = 'knowledgeDragState';
let dragState = inject(dragStateKey, null);
if (!dragState) {
  dragState = {
    draggedNodeId: ref(null),
    dragOverId: ref(null),
    dragOverPosition: ref(null)
  };
  provide(dragStateKey, dragState);
}
const { draggedNodeId, dragOverId, dragOverPosition } = dragState;

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

const emit = defineEmits(['node-click', 'create-child', 'delete-node', 'move-node', 'rename-node', 'drag-drop']);

const expandedNodes = ref(new Set());
const editingId = ref(null);
const editingName = ref('');

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

watch(() => props.treeData, (newData) => {
  expandedNodes.value.clear();
  expandAllNodes(newData);
}, { deep: true });

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
  if (editingId.value !== null) return;
  emit('node-click', node);
};

const startEdit = (node, event) => {
  event.stopPropagation();
  editingId.value = node.id;
  editingName.value = node.name;
  nextTick(() => {
    const inputs = document.querySelectorAll('.knowledge-tree .edit-input');
    const inputArray = Array.from(inputs);
    const target = inputArray.find(el => el.dataset.id == node.id);
    if (target) {
      target.focus();
      target.select();
    }
  });
};

const confirmEdit = (node) => {
  const newName = editingName.value.trim();
  if (newName && newName !== node.name) {
    emit('rename-node', node.id, newName);
  }
  editingId.value = null;
  editingName.value = '';
};

const cancelEdit = () => {
  editingId.value = null;
  editingName.value = '';
};

const handleKeydown = (e, node) => {
  if (e.key === 'Enter') {
    confirmEdit(node);
  } else if (e.key === 'Escape') {
    cancelEdit();
  }
};

// === 拖拽逻辑 ===
const onDragStart = (node, event) => {
  if (editingId.value !== null) return;
  draggedNodeId.value = node.id;
  event.dataTransfer.effectAllowed = 'move';
  event.dataTransfer.setData('text/plain', String(node.id));
};

const onDragOver = (node, event) => {
  if (!draggedNodeId.value || draggedNodeId.value === node.id) return;
  event.preventDefault();
  event.dataTransfer.dropEffect = 'move';

  const rect = event.currentTarget.getBoundingClientRect();
  const y = event.clientY - rect.top;
  const h = rect.height;

  if (y < h * 0.25) {
    dragOverPosition.value = 'before';
  } else if (y > h * 0.75) {
    dragOverPosition.value = 'after';
  } else {
    dragOverPosition.value = 'inside';
  }
  dragOverId.value = node.id;
};

const onDragLeave = (node) => {
  if (dragOverId.value === node.id) {
    dragOverId.value = null;
    dragOverPosition.value = null;
  }
};

const onDrop = (node, event) => {
  event.preventDefault();
  event.stopPropagation();
  if (!draggedNodeId.value || draggedNodeId.value === node.id) {
    resetDrag();
    return;
  }

  emit('drag-drop', {
    draggedId: draggedNodeId.value,
    targetId: node.id,
    position: dragOverPosition.value
  });
  resetDrag();
};

const resetDrag = () => {
  draggedNodeId.value = null;
  dragOverId.value = null;
  dragOverPosition.value = null;
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
        <div
          class="node-content"
          :class="{
            'node-selected': selectedId === node.id,
            'drag-over-before': dragOverId === node.id && dragOverPosition === 'before',
            'drag-over-after': dragOverId === node.id && dragOverPosition === 'after',
            'drag-over-inside': dragOverId === node.id && dragOverPosition === 'inside',
            'dragging': draggedNodeId === node.id
          }"
          @click="toggleNode(node.id)"
          draggable="true"
          @dragstart="onDragStart(node, $event)"
          @dragover="onDragOver(node, $event)"
          @dragleave="onDragLeave(node)"
          @drop="onDrop(node, $event)"
          @dragend="resetDrag"
        >
          <span class="node-toggle" v-if="hasChildren(node)">
            {{ isExpanded(node.id) ? '▼' : '►' }}
          </span>
          <span
            v-if="editingId !== node.id"
            class="node-name"
            @click.stop="handleNodeClick(node)"
            @dblclick.stop="startEdit(node, $event)"
          >{{ node.name }}</span>
          <input
            v-else
            v-model="editingName"
            class="edit-input"
            :data-id="node.id"
            @click.stop
            @blur="confirmEdit(node)"
            @keydown="handleKeydown($event, node)"
          />
          <div class="node-actions">
            <button class="action-btn create-btn" @click.stop="emit('create-child', node.id)" title="创建子知识库">+</button>
            <button class="action-btn delete-btn" @click.stop="emit('delete-node', node.id)" title="删除知识库">×</button>
          </div>
        </div>
        <div
          v-if="hasChildren(node) && isExpanded(node.id)"
          class="tree-children"
        >
          <KnowledgeTree
            :tree-data="node.children"
            :selected-id="selectedId"
            @node-click="handleNodeClick"
            @create-child="(id) => emit('create-child', id)"
            @delete-node="(id) => emit('delete-node', id)"
            @move-node="(id) => emit('move-node', id)"
            @rename-node="(id, name) => emit('rename-node', id, name)"
            @drag-drop="(payload) => emit('drag-drop', payload)"
          />
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
  border: 2px solid transparent;
}

.node-content:hover {
  background-color: #e9ecef;
}

.node-content.node-selected {
  background-color: #d1ecf1;
  border-left: 3px solid #17a2b8;
  font-weight: 600;
}

.node-content.dragging {
  opacity: 0.4;
}

.node-content.drag-over-before {
  border-top: 2px solid #0d6efd;
}

.node-content.drag-over-after {
  border-bottom: 2px solid #0d6efd;
}

.node-content.drag-over-inside {
  background-color: #d1ecf1;
  border: 2px solid #0d6efd;
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
  margin-right: 4rem;
}

.edit-input {
  flex: 1;
  font-size: 0.875rem;
  padding: 0.125rem 0.25rem;
  border: 1px solid #0d6efd;
  border-radius: 3px;
  outline: none;
  margin-right: 4rem;
  background-color: white;
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

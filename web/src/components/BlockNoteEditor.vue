<template>
  <div ref="container" class="blocknote-editor-container"></div>
</template>

<script setup>
import { ref, onMounted, onBeforeUnmount, watch } from 'vue';
import { createElement } from 'react';
import { createRoot } from 'react-dom/client';
import { BlockNoteEditor } from '@blocknote/core';
import { BlockNoteView } from '@blocknote/mantine';

import '@blocknote/core/fonts/inter.css';
import '@blocknote/mantine/style.css';

const props = defineProps({
  modelValue: { type: String, default: '' },
  readonly: { type: Boolean, default: false }
});
const emit = defineEmits(['update:modelValue']);

const container = ref(null);
let root = null;
let editor = null;
let isUpdating = false;

const getAuthHeaders = () => {
  const token = localStorage.getItem('access_token');
  return token ? { 'Authorization': `Bearer ${token}` } : {};
};

const uploadFile = async (file) => {
  const formData = new FormData();
  formData.append('file', file);
  const res = await fetch('/api/v1/files/upload', {
    method: 'POST',
    headers: getAuthHeaders(),
    body: formData
  });
  if (!res.ok) throw new Error(`上传失败: ${res.status}`);
  const result = await res.json();
  return result.Url || result.url;
};

onMounted(() => {
  let initialContent;
  if (props.modelValue) {
    try {
      initialContent = undefined; // 先创建空编辑器，再用 tryParseHTMLToBlocks 填充
    } catch {
      initialContent = undefined;
    }
  }

  editor = BlockNoteEditor.create({
    uploadFile,
    ...(initialContent ? { initialContent } : {})
  });

  // 如果有初始 HTML 内容，解析并填入
  if (props.modelValue) {
    try {
      const blocks = editor.tryParseHTMLToBlocks(props.modelValue);
      if (blocks && blocks.length > 0) {
        editor.replaceBlocks(editor.document, blocks);
      }
    } catch (err) {
      console.error('解析HTML内容失败:', err);
    }
  }

  // 用 React createRoot 挂载 BlockNoteView，显式配置 slash 菜单
  root = createRoot(container.value);
  root.render(
    createElement(BlockNoteView, { 
      editor, 
      slashMenu: !props.readonly,
      editable: !props.readonly,
      formattingToolbar: !props.readonly,
      sideMenu: !props.readonly,
      filePanel: !props.readonly,
      tableHandles: !props.readonly
    })
  );

  // 监听编辑器内容变化，同步 HTML 到父组件
  editor.onChange(() => {
    if (isUpdating) return;
    const html = editor.blocksToHTMLLossy(editor.document);
    emit('update:modelValue', html);
  });
});

// 监听外部 modelValue 变化（如切换文章时）
watch(() => props.modelValue, (newVal) => {
  if (!editor || isUpdating) return;
  if (!newVal) {
    editor.replaceBlocks(editor.document, [{ type: 'paragraph', content: [] }]);
    return;
  }
  const currentHtml = editor.blocksToHTMLLossy(editor.document);
  if (currentHtml !== newVal) {
    isUpdating = true;
    try {
      const blocks = editor.tryParseHTMLToBlocks(newVal);
      editor.replaceBlocks(editor.document, blocks);
    } catch (err) {
      console.error('更新编辑器内容失败:', err);
    } finally {
      isUpdating = false;
    }
  }
});

onBeforeUnmount(() => {
  if (root) {
    root.unmount();
  }
  if (editor) {
    editor.destroy?.();
  }
});
</script>

<style scoped>
.blocknote-editor-container {
  height: 100%;
  overflow: auto;
}

.blocknote-editor-container :deep(.bn-container) {
  height: 100%;
}

.blocknote-editor-container :deep(.bn-editor) {
  height: 100%;
  padding: 1rem;
}
</style>

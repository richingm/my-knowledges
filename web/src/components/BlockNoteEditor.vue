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
const emit = defineEmits(['update:modelValue', 'tocUpdate']);

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

// 从编辑器文档中提取 H1-H6 标题
const extractToc = () => {
  if (!editor) return [];
  const toc = [];
  const blocks = editor.document;
  let headingIndex = 0;
  blocks.forEach((block) => {
    if (block.type === 'heading') {
      let text = '';
      if (block.content && Array.isArray(block.content)) {
        text = block.content
          .filter(item => item.type === 'text')
          .map(item => item.text)
          .join('');
      }
      if (text.trim()) {
        toc.push({
          id: `heading-${headingIndex}`,
          headingIndex: headingIndex,
          level: block.props?.level || 1,
          text: text.trim()
        });
      }
      headingIndex++;
    }
  });
  return toc;
};

const emitToc = () => {
  emit('tocUpdate', extractToc());
};

// 滚动到指定标题
const scrollToHeading = (headingIndex) => {
  if (!editor || !container.value) return;

  const editorDom = container.value.querySelector('.bn-editor');
  if (!editorDom) return;

  const headingEls = editorDom.querySelectorAll('[data-content-type="heading"]');
  if (headingIndex < 0 || headingIndex >= headingEls.length) return;

  const target = headingEls[headingIndex];
  if (target) {
    target.scrollIntoView({ behavior: 'smooth', block: 'start' });
  }
};

defineExpose({ scrollToHeading });

onMounted(() => {
  let initialContent;
  if (props.modelValue) {
    try {
      initialContent = undefined;
    } catch {
      initialContent = undefined;
    }
  }

  editor = BlockNoteEditor.create({
    uploadFile,
    ...(initialContent ? { initialContent } : {})
  });

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

  editor.onChange(() => {
    if (isUpdating) return;
    const html = editor.blocksToHTMLLossy(editor.document);
    emit('update:modelValue', html);
    emitToc();
  });

  // 初始加载后延迟发送 toc
  setTimeout(() => emitToc(), 300);
});

watch(() => props.modelValue, (newVal) => {
  if (!editor || isUpdating) return;
  if (!newVal) {
    editor.replaceBlocks(editor.document, [{ type: 'paragraph', content: [] }]);
    emitToc();
    return;
  }
  const currentHtml = editor.blocksToHTMLLossy(editor.document);
  if (currentHtml !== newVal) {
    isUpdating = true;
    try {
      const blocks = editor.tryParseHTMLToBlocks(newVal);
      editor.replaceBlocks(editor.document, blocks);
      setTimeout(() => emitToc(), 300);
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

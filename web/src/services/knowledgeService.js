const getAuthHeaders = () => {
  const token = localStorage.getItem('access_token');
  return token ? { 'Authorization': `Bearer ${token}` } : {};
};

export const knowledgeService = {
  async getKnowledgeTree(domainId) {
    try {
      const response = await fetch(`/api/v1/knowledge-tree?domain_id=${domainId}`, {
        headers: getAuthHeaders()
      });
      if (!response.ok) {
        throw new Error('Failed to fetch knowledge tree');
      }
      const data = await response.json();
      return data.items || [];
    } catch (error) {
      console.error('Error fetching knowledge tree:', error);
      return [];
    }
  },
  
  async createKnowledge(knowledgeData) {
    try {
      const response = await fetch('/api/v1/knowledges', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          ...getAuthHeaders()
        },
        body: JSON.stringify(knowledgeData)
      });
      if (!response.ok) {
        throw new Error('Failed to create knowledge');
      }
      const data = await response.json();
      return data;
    } catch (error) {
      console.error('Error creating knowledge:', error);
      return null;
    }
  },
  
  async deleteKnowledge(knowledgeId) {
    try {
      const response = await fetch(`/api/v1/knowledges/${knowledgeId}`, {
        method: 'DELETE',
        headers: getAuthHeaders()
      });
      if (!response.ok) {
        throw new Error('Failed to delete knowledge');
      }
      return true;
    } catch (error) {
      console.error('Error deleting knowledge:', error);
      return false;
    }
  },

  async updateKnowledge(knowledgeId, domainId, name, description, parentKnowledgeId) {
    try {
      const body = {
        domain_id: domainId,
        name: name,
        description: description
      };
      if (parentKnowledgeId !== undefined) {
        body.parent_knowledge_id = parentKnowledgeId;
      }
      const response = await fetch(`/api/v1/knowledges/${knowledgeId}`, {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
          ...getAuthHeaders()
        },
        body: JSON.stringify(body)
      });
      if (!response.ok) {
        throw new Error('Failed to update knowledge');
      }
      const data = await response.json();
      return data;
    } catch (error) {
      console.error('Error updating knowledge:', error);
      return null;
    }
  },

  async moveKnowledge(knowledgeId, newParentId, domainId) {
    try {
      const response = await fetch('/api/v1/knowledges/move', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          ...getAuthHeaders()
        },
        body: JSON.stringify({
          id: parseInt(knowledgeId),
          newParentId: newParentId ? parseInt(newParentId) : 0,
          domainId: parseInt(domainId)
        })
      });
      if (!response.ok) {
        throw new Error('Failed to move knowledge');
      }
      const data = await response.json();
      return data;
    } catch (error) {
      console.error('Error moving knowledge:', error);
      return null;
    }
  },

  async sortKnowledge(items) {
    try {
      const response = await fetch('/api/v1/knowledges/sort', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          ...getAuthHeaders()
        },
        body: JSON.stringify({
          items: items.map(item => ({
            id: parseInt(item.id),
            bySort: parseInt(item.by_sort)
          }))
        })
      });
      if (!response.ok) {
        throw new Error('Failed to sort knowledge');
      }
      const data = await response.json();
      return data;
    } catch (error) {
      console.error('Error sorting knowledge:', error);
      return null;
    }
  }
};

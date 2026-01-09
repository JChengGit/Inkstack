import { apiClient } from './api';
import { Post, PostListResponse, Comment } from '@/types';

export const postsService = {
  async getPosts(page = 1, pageSize = 10, status?: string): Promise<PostListResponse> {
    const params: any = { page, page_size: pageSize };
    if (status) params.status = status;

    const response = await apiClient.get<PostListResponse>('/posts', { params });
    return response.data;
  },

  async getPostById(id: number): Promise<Post> {
    const response = await apiClient.get<Post>(`/posts/${id}`);
    return response.data;
  },

  async getPostBySlug(slug: string): Promise<Post> {
    const response = await apiClient.get<Post>(`/posts/slug/${slug}`);
    return response.data;
  },

  async createPost(data: {
    title: string;
    content: string;
    excerpt?: string;
    slug?: string;
  }): Promise<Post> {
    const response = await apiClient.post<Post>('/posts', data);
    return response.data;
  },

  async updatePost(
    id: number,
    data: {
      title?: string;
      content?: string;
      excerpt?: string;
      status?: string;
    }
  ): Promise<Post> {
    const response = await apiClient.put<Post>(`/posts/${id}`, data);
    return response.data;
  },

  async deletePost(id: number): Promise<void> {
    await apiClient.delete(`/posts/${id}`);
  },

  async publishPost(id: number): Promise<Post> {
    const response = await apiClient.post<Post>(`/posts/${id}/publish`);
    return response.data;
  },

  async getComments(postId: number): Promise<Comment[]> {
    const response = await apiClient.get<Comment[]>(`/posts/${postId}/comments`);
    return response.data;
  },

  async createComment(
    postId: number,
    content: string,
    parentId?: number
  ): Promise<Comment> {
    const response = await apiClient.post<Comment>(`/posts/${postId}/comments`, {
      content,
      parent_id: parentId,
    });
    return response.data;
  },
};

'use client';

import { useEffect, useState } from 'react';
import { useParams } from 'next/navigation';
import { postsService } from '@/services/posts';
import { Post, Comment } from '@/types';
import { formatDate } from '@/lib/utils';

export default function PostDetailPage() {
  const params = useParams();
  const slug = params.slug as string;

  const [post, setPost] = useState<Post | null>(null);
  const [comments, setComments] = useState<Comment[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');

  useEffect(() => {
    if (slug) {
      loadPost();
    }
  }, [slug]);

  const loadPost = async () => {
    setLoading(true);
    setError('');

    try {
      const postData = await postsService.getPostBySlug(slug);
      setPost(postData);

      // Load comments
      const commentsData = await postsService.getComments(postData.id);
      setComments(commentsData);
    } catch (err: any) {
      setError('Failed to load post. Please try again later.');
      console.error('Error loading post:', err);
    } finally {
      setLoading(false);
    }
  };

  if (loading) {
    return (
      <div className="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
        <div className="text-center py-12">
          <div className="inline-block h-8 w-8 animate-spin rounded-full border-4 border-solid border-blue-600 border-r-transparent"></div>
          <p className="mt-4 text-gray-600">Loading post...</p>
        </div>
      </div>
    );
  }

  if (error || !post) {
    return (
      <div className="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
        <div className="rounded-md bg-red-50 p-4">
          <p className="text-sm text-red-800">
            {error || 'Post not found'}
          </p>
        </div>
      </div>
    );
  }

  return (
    <div className="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
      <article className="bg-white rounded-lg shadow-sm border p-8 mb-8">
        <header className="mb-8">
          <h1 className="text-4xl font-bold text-gray-900 mb-4">
            {post.title}
          </h1>

          <div className="flex items-center gap-4 text-sm text-gray-500">
            <time>{formatDate(post.created_at)}</time>
            <span>•</span>
            <span>{post.view_count} views</span>
            {post.status !== 'published' && (
              <>
                <span>•</span>
                <span className="px-2 py-1 bg-yellow-100 text-yellow-800 rounded text-xs font-medium">
                  {post.status}
                </span>
              </>
            )}
          </div>
        </header>

        <div className="prose max-w-none">
          <div className="whitespace-pre-wrap text-gray-700 leading-relaxed">
            {post.content}
          </div>
        </div>
      </article>

      <section className="bg-white rounded-lg shadow-sm border p-8">
        <h2 className="text-2xl font-bold text-gray-900 mb-6">
          Comments ({comments.length})
        </h2>

        {comments.length === 0 ? (
          <p className="text-gray-600">No comments yet. Be the first to comment!</p>
        ) : (
          <div className="space-y-4">
            {comments.map((comment) => (
              <div
                key={comment.id}
                className="border-l-4 border-blue-500 pl-4 py-2"
              >
                <div className="flex items-center gap-2 text-sm text-gray-500 mb-2">
                  <span>User #{comment.user_id}</span>
                  <span>•</span>
                  <time>{formatDate(comment.created_at)}</time>
                  {comment.status !== 'approved' && (
                    <>
                      <span>•</span>
                      <span className="px-2 py-1 bg-gray-100 text-gray-600 rounded text-xs">
                        {comment.status}
                      </span>
                    </>
                  )}
                </div>
                <p className="text-gray-700">{comment.content}</p>
              </div>
            ))}
          </div>
        )}
      </section>
    </div>
  );
}

import Link from 'next/link';
import { Post } from '@/types';
import { formatDate, truncate } from '@/lib/utils';

interface PostCardProps {
  post: Post;
}

export default function PostCard({ post }: PostCardProps) {
  return (
    <article className="bg-white rounded-lg shadow-sm border p-6 hover:shadow-md transition-shadow">
      <Link href={`/posts/${post.slug}`}>
        <h2 className="text-2xl font-bold text-gray-900 hover:text-blue-600 mb-2">
          {post.title}
        </h2>
      </Link>

      <div className="flex items-center gap-4 text-sm text-gray-500 mb-3">
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

      {post.excerpt && (
        <p className="text-gray-600 mb-4">{truncate(post.excerpt, 200)}</p>
      )}

      <Link
        href={`/posts/${post.slug}`}
        className="text-blue-600 hover:text-blue-800 font-medium text-sm"
      >
        Read more →
      </Link>
    </article>
  );
}

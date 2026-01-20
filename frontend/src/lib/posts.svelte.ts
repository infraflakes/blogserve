// No longer needed: import { onMount } from 'svelte';

export interface PostMetadata {
  title: string;
  date: string;
  tags: string[];
  description: string;
}

export interface Post {
  slug: string;
  content: string;
  metadata: PostMetadata;
}

class PostsState {
  posts = $state<Post[]>([]);
  loading = $state(true);
  error = $state<string | null>(null);

  constructor() {
    if (typeof window !== 'undefined') {
      this.init();
    }
  }

  async fetchPosts() {
    try {
      const res = await fetch('/api/posts');
      if (!res.ok) throw new Error('Failed to fetch posts');
      this.posts = await res.json();
    } catch (e: unknown) {
      this.error = e instanceof Error ? e.message : String(e);
      console.error('Error fetching posts:', e);
    } finally {
      this.loading = false;
    }
  }

  init() {
    this.fetchPosts();

    const eventSource = new EventSource('/api/reload');
    eventSource.onmessage = (event) => {
      if (event.data === 'reload') {
        this.fetchPosts();
      }
    };

    // In Svelte 5, we can use $effect for teardown if this were a component,
    // but since it's a shared state object, we'll keep it simple for now.
    // For a more robust implementation, we could track connections.
  }

  getPost(slug: string) {
    return this.posts.find((p) => p.slug === slug) ?? null;
  }
}

export const postsState = new PostsState();

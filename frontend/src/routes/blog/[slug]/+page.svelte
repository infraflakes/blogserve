<script lang="ts">
import { page } from '$app/state';
import { marked } from 'marked';
import { postsState } from '$lib/posts.svelte';

let post = $derived(postsState.getPost(page.params.slug));

let htmlContent = $derived(post ? marked.parse(post.content || '') : '');

function processHtml(html: string) {
  if (typeof document === 'undefined' || !post) return html;
  const parser = new DOMParser();
  const doc = parser.parseFromString(html, 'text/html');
  const images = doc.querySelectorAll('img');

  images.forEach((img) => {
    const src = img.getAttribute('src');
    if (src && !src.startsWith('http') && !src.startsWith('/')) {
      img.setAttribute('src', `/data/${post.slug}/${src}`);
    }
  });

  return doc.body.innerHTML;
}

let processedHtml = $derived(processHtml(htmlContent));
</script>

<svelte:head>
	{#if post}
		<title>{post.metadata.title} | blogserve</title>
		<meta name="description" content={post.metadata.description} />
		<meta property="og:title" content={post.metadata.title} />
		<meta property="og:description" content={post.metadata.description} />
	{/if}
</svelte:head>

{#if postsState.loading}
	<div class="status-msg">Loading post...</div>
{:else if post}
	<article class="post">
		<header>
			<div class="tags">
				{#each post.metadata.tags || [] as tag}
					<span class="tag">{tag}</span>
				{/each}
			</div>
			<h1>{post.metadata.title || post.slug}</h1>
			<p class="date">{post.metadata.date || ''}</p>
		</header>

		<div class="markdown-body">
			{@html processedHtml}
		</div>
	</article>
{:else}
	<div class="status-msg error">Post not found.</div>
{/if}

<style>
	.post {
		animation: fadeIn 0.5s ease-out;
	}
	header {
		margin-bottom: var(--space-xl);
	}
	h1 {
		font-size: 3.5rem;
		margin: var(--space-xs) 0;
		color: var(--text-primary);
	}
	.date {
		color: var(--text-secondary);
		font-size: 1.1rem;
		font-weight: 500;
	}
	.tags {
		display: flex;
		gap: var(--space-xs);
		margin-bottom: var(--space-md);
	}
	.status-msg {
		text-align: center;
		padding: var(--space-xl);
		color: var(--text-secondary);
	}
	.status-msg.error {
		color: #ef4444;
	}

	@keyframes fadeIn {
		from {
			opacity: 0;
			transform: translateY(10px);
		}
		to {
			opacity: 1;
			transform: translateY(0);
		}
	}
</style>

<script lang="ts">
	import { page } from '$app/state';
	import { onMount, tick } from 'svelte';
	import { marked } from 'marked';

	let post = $state(null);
	let loading = $state(true);
	let slug = $derived(page.params.slug);

	async function fetchPost() {
		try {
			const res = await fetch('http://localhost:8080/api/posts');
			const posts = await res.json();
			post = posts.find((p) => p.slug === slug);
		} catch (e) {
			console.error('Failed to fetch post', e);
		} finally {
			loading = false;
		}
	}

	onMount(() => {
		fetchPost();

		const eventSource = new EventSource('http://localhost:8080/api/reload');
		eventSource.onmessage = (event) => {
			if (event.data === 'reload') {
				fetchPost();
			}
		};

		return () => eventSource.close();
	});

	let htmlContent = $derived(post ? marked.parse(post.content || '') : '');

	function processHtml(html: string) {
		if (typeof document === 'undefined') return html;
		const parser = new DOMParser();
		const doc = parser.parseFromString(html, 'text/html');
		const images = doc.querySelectorAll('img');

		images.forEach((img) => {
			const src = img.getAttribute('src');
			if (src && !src.startsWith('http') && !src.startsWith('/')) {
				img.setAttribute('src', `http://localhost:8080/data/${post.slug}/${src}`);
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

{#if loading}
	<p>Loading post...</p>
{:else if post}
	<article class="post">
		<header>
            <div class="tags">
                {#each (post.metadata.tags || []) as tag}
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
	<p>Post not found.</p>
{/if}

<style>
	.post {
		line-height: 1.7;
	}
	header {
		margin-bottom: 3rem;
	}
	h1 {
		font-size: 3rem;
		margin: 0.5rem 0;
		font-weight: 800;
	}
	.date {
		color: #64748b;
		font-size: 1.1rem;
	}
	.tags {
		display: flex;
		gap: 0.5rem;
		margin-bottom: 1rem;
	}
	.tag {
		background: #1e293b;
		color: #38bdf8;
		padding: 0.2rem 0.6rem;
		border-radius: 9999px;
		font-size: 0.75rem;
		font-weight: 600;
        text-transform: uppercase;
        letter-spacing: 0.05em;
	}
    :global(.markdown-body) {
        font-size: 1.125rem;
        color: #e2e8f0;
    }
	:global(.markdown-body img) {
		max-width: 100%;
		border-radius: 1rem;
        margin: 2rem 0;
        box-shadow: 0 20px 25px -5px rgba(0, 0, 0, 0.1), 0 10px 10px -5px rgba(0, 0, 0, 0.04);
	}
	:global(.markdown-body pre) {
		background: #1e293b;
		padding: 1.5rem;
		border-radius: 0.75rem;
		overflow-x: auto;
        border: 1px solid #334155;
        margin: 2rem 0;
	}
    :global(.markdown-body a) {
        color: #38bdf8;
        text-decoration: none;
        border-bottom: 1px solid transparent;
        transition: border-color 0.2s;
    }
    :global(.markdown-body a:hover) {
        border-bottom-color: #38bdf8;
    }
</style>

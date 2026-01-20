<script lang="ts">
import PostList from '$lib/components/PostList.svelte';
import { postsState } from '$lib/posts.svelte';
</script>

<svelte:head>
	<title>blogserve | Home</title>
	<meta name="description" content="A fast and simple blog engine" />
</svelte:head>

{#if postsState.loading}
	<div class="loading">
		<div class="spinner"></div>
		<p>Fetching your stories...</p>
	</div>
{:else if postsState.error}
	<div class="error-state">
		<p>Failed to load posts: {postsState.error}</p>
	</div>
{:else}
	<PostList posts={postsState.posts} />
{/if}

<style>
	.loading {
		display: flex;
		flex-direction: column;
		align-items: center;
		justify-content: center;
		height: 40vh;
		color: var(--text-muted);
	}
	.spinner {
		width: 40px;
		height: 40px;
		border: 3px solid var(--border);
		border-top-color: var(--accent);
		border-radius: 50%;
		animation: spin 1s linear infinite;
		margin-bottom: var(--space-md);
	}
	@keyframes spin {
		to {
			transform: rotate(360deg);
		}
	}
	.error-state {
		text-align: center;
		color: #ef4444;
		padding: var(--space-xl);
	}
</style>

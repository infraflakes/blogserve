<script lang="ts">
	import { onMount } from 'svelte';
	import PostList from '$lib/components/PostList.svelte';

	let posts = $state([]);
	let loading = $state(true);

	async function fetchPosts() {
		try {
			const res = await fetch('http://localhost:8080/api/posts');
			posts = await res.json();
		} catch (e) {
			console.error('Failed to fetch posts', e);
		} finally {
			loading = false;
		}
	}

	onMount(() => {
		fetchPosts();

		const eventSource = new EventSource('http://localhost:8080/api/reload');
		eventSource.onmessage = (event) => {
			if (event.data === 'reload') {
				fetchPosts();
			}
		};

		return () => eventSource.close();
	});
</script>

<svelte:head>
	<title>blogserve | Home</title>
	<meta name="description" content="A fast and simple blog engine" />
</svelte:head>

{#if loading}
	<div class="loading">
        <div class="spinner"></div>
        <p>Fetching your stories...</p>
    </div>
{:else}
	<PostList {posts} />
{/if}

<style>
	.loading {
		display: flex;
		flex-direction: column;
		align-items: center;
		justify-content: center;
		height: 50vh;
		color: #64748b;
	}
    .spinner {
        width: 40px;
        height: 40px;
        border: 3px solid #1e293b;
        border-top-color: #38bdf8;
        border-radius: 50%;
        animation: spin 1s linear infinite;
        margin-bottom: 1rem;
    }
    @keyframes spin {
        to { transform: rotate(360deg); }
    }
</style>

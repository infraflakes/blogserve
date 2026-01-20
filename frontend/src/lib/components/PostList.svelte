<script lang="ts">
let { posts } = $props();
</script>

<div class="post-list">
	{#each posts as post}
		<a href="/blog/{post.slug}/" class="post-card">
			<div class="card-content">
				<span class="post-date">{post.metadata.date || 'No date'}</span>
				<h2>{post.metadata.title || post.slug}</h2>
				<p>{post.metadata.description || 'No description available.'}</p>
				<div class="tags">
					{#each post.metadata.tags || [] as tag}
						<span class="tag">{tag}</span>
					{/each}
				</div>
			</div>
		</a>
	{/each}
</div>

<style>
	.post-list {
		display: grid;
		grid-template-columns: 20fr;
		gap: var(--space-lg);
	}
	.post-card {
		background: var(--bg-secondary);
		border-radius: 1.25rem;
		border: 1px solid var(--border);
		transition: var(--transition-smooth);
		overflow: hidden;
		display: block;
		text-decoration: none;
		color: inherit;
	}
	.post-card:hover {
		transform: translateY(-4px);
		border-color: var(--border-hover);
		background: var(--surface);
		box-shadow: 0 20px 25px -5px rgba(0, 0, 0, 0.4);
	}
	.card-content {
		padding: var(--space-lg);
	}
	h2 {
		margin: var(--space-xs) 0 var(--space-md) 0;
		color: var(--text-primary);
		font-size: 1.75rem;
		font-weight: 700;
	}
	p {
		color: var(--text-secondary);
		margin: 0 0 var(--space-md) 0;
		line-height: 1.6;
	}
	.post-date {
		display: block;
		font-size: 0.875rem;
		color: var(--accent);
		font-weight: 600;
		text-transform: uppercase;
		letter-spacing: 0.05em;
	}
	.tags {
		display: flex;
		gap: 0.5rem;
		flex-wrap: wrap;
	}
</style>

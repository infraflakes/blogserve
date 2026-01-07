<script lang="ts">
	let { posts } = $props();
</script>

<div class="post-list">
	{#each posts as post}
		<a href="/post/{post.slug}/" class="post-card">
			<div class="card-content">
                <span class="date">{post.metadata.date || 'No date'}</span>
                <h2>{post.metadata.title || post.slug}</h2>
                <p>{post.metadata.description || 'No description available.'}</p>
                <div class="tags">
                    {#each (post.metadata.tags || []) as tag}
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
        grid-template-columns: 1fr;
		gap: 2rem;
	}
	.post-card {
		background:#23262f;
		border-radius: 1.25rem;
		border: 1px solid #334155;
		text-decoration: none;
        color: inherit;
		transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
        overflow: hidden;
        display: block;
	}
	.post-card:hover {
		transform: translateY(-4px);
		border-color: #b3c7ff;
        box-shadow: 0 20px 25px -5px rgba(0, 0, 0, 0.2);
	}
    .card-content {
        padding: 2rem;
    }
	h2 {
		margin: 0.5rem 0 1rem 0;
		color: #f8fafc;
        font-size: 1.75rem;
        font-weight: 700;
	}
	p {
		color: #94a3b8;
		margin: 0 0 1.5rem 0;
        line-height: 1.6;
	}
	.date {
        display: block;
		font-size: 0.875rem;
		color: #b3c7ff;
        font-weight: 600;
        text-transform: uppercase;
        letter-spacing: 0.05em;
	}
	.tags {
		display: flex;
		gap: 0.5rem;
        flex-wrap: wrap;
	}
	.tag {
		background:#17181c
    color: #64748b;
		padding: 0.25rem 0.75rem;
		border-radius: 0.5rem;
        font-size: 0.75rem;
        font-weight: 500;
        border: 1px solid #1e293b;
	}
</style>

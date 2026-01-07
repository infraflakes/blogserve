<script>
	let { children } = $props();
</script>

<div class="app">
	<nav>
		<div class="container">
			<a href="/" class="logo">blogserve</a>
		</div>
	</nav>

	<main class="container">
		{@render children()}
	</main>

	<footer>
		<div class="container">
			<p>&copy; 2026 blogserve. Built with Go + SvelteKit.</p>
		</div>
	</footer>
</div>

<style>
	:global(body) {
		font-family: 'Inter', system-ui, -apple-system, sans-serif;
		background: #0f172a;
		color: #f8fafc;
		margin: 0;
		line-height: 1.5;
		-webkit-font-smoothing: antialiased;
	}

	.container {
		max-width: 900px;
		margin: 0 auto;
		padding: 0 1.5rem;
	}

	nav {
		background: rgba(15, 23, 42, 0.8);
		backdrop-filter: blur(8px);
		border-bottom: 1px solid #1e293b;
		padding: 1rem 0;
		position: sticky;
		top: 0;
		z-index: 10;
	}

	.logo {
		font-size: 1.5rem;
		font-weight: 800;
		text-decoration: none;
		background: linear-gradient(to right, #38bdf8, #818cf8);
		-webkit-background-clip: text;
		-webkit-text-fill-color: transparent;
	}

	main {
		padding-top: 3rem;
		padding-bottom: 5rem;
		min-height: 80vh;
	}

	footer {
		border-top: 1px solid #1e293b;
		padding: 3rem 0;
		color: #64748b;
		text-align: center;
		font-size: 0.875rem;
	}
</style>

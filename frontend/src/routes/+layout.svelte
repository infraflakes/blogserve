<script>
	import '../app.css';
	let { children } = $props();
</script>

<div class="app">
	<header class="navbar">
		<div class="container">
			<a href="/" class="logo">blogserve</a>
			<nav class="nav-links">
				<a href="/" class="nav-link">Blog</a>
				<!-- Future links: Portfolio, Contact, etc. -->
			</nav>
		</div>
	</header>

	<main class="container">
		<div class="page-spacer"></div>
		{@render children()}
	</main>

	<footer>
		<div class="container">
			<p>&copy; 2026 blogserve. Built with Go + SvelteKit.</p>
		</div>
	</footer>
</div>

<style>
	.navbar {
		background: rgba(10, 11, 16, 0.8);
		backdrop-filter: blur(12px);
		border-bottom: 1px solid var(--border);
		padding: var(--space-md) 0;
		position: sticky;
		top: 0;
		z-index: 10;
	}

	.logo {
		font-size: 1.75rem;
		font-weight: 800;
		color: var(--logo-color);
	}

	.nav-links {
		display: inline-flex;
		gap: var(--space-md);
		margin-left: var(--space-xl);
	}

	.nav-link {
		color: var(--text-secondary);
		font-weight: 600;
	}

	.nav-link:hover {
		color: var(--accent);
	}

	.page-spacer {
		height: 6rem;
	}

	footer {
		border-top: 1px solid var(--border);
		padding: var(--space-xl) 0;
		color: var(--text-muted);
		text-align: center;
		font-size: 0.875rem;
	}
</style>

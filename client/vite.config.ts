import { sveltekit } from '@sveltejs/kit/vite';
import { defineConfig } from 'vite';

export default defineConfig({
	plugins: [sveltekit()],
	optimizeDeps: {
		include: ["dayjs/plugin/relativeTime.js"]
	},
	build:  {
		rollupOptions: {
			external: [
				"dayjs/plugin/relativeTime.js"
			]
		}
	}
});

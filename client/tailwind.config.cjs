/** @type {import('tailwindcss').Config}*/
const config = {
	content: ['./src/**/*.{html,js,svelte,ts}'],

	theme: {
		extend: {
			colors: {
				"background": "#0a0a0a"
			}
		}
	},

	plugins: [
		require('tailwindcss-animate')
	]
};

module.exports = config;

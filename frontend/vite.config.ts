import { sveltekit } from '@sveltejs/kit/vite';
import { defineConfig } from 'vite';

export default defineConfig({
	plugins: [sveltekit()],
	server: {
		cors: {
			origin: '*',
			methods: ['GET', 'POST', 'OPTIONS'],
			allowedHeaders: ['Accept', 'Content-Type', 'Content-Length', 'Accept-Encoding', 'Authorization', 'ResponseType'],
		},
		proxy: {
			'/stream.AudioStream/StreamAudio': {
				target: 'http://localhost:50051',
				changeOrigin: true,
				secure: false,
			}
		}
	},
});

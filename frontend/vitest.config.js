import { defineConfig } from 'vitest/config'

// https://vitest.dev/config/
export default defineConfig({
  test: {
    environment: 'happy-dom',
    environmentOptions: {
        happyDom: {
            settings: {
                location: {
                    href: "http://localhost"
                }
            }
        }
    }
  },
});
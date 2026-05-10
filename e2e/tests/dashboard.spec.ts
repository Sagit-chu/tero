// e2e/tests/dashboard.spec.ts
import { test, expect } from '@playwright/test';

test('has title and displays status ok', async ({ page }) => {
  await page.goto('/');

  // Expect a title "to contain" a substring.
  await expect(page).toHaveTitle(/Flvx Monitor Dashboard/);
  
  // Wait for the loading state to disappear and status ok to appear
  await expect(page.locator('text=System Status')).toBeVisible();
  await expect(page.locator('text=Online')).toBeVisible();
});

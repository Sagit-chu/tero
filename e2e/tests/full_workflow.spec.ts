import { test, expect } from '@playwright/test';

test.describe('Flvx Monitor Dashboard Full Workflow', () => {
  test.beforeEach(async ({ page }) => {
    // Navigate to the dashboard
    await page.goto('/');
  });

  test('should display dashboard correctly', async ({ page }) => {
    await expect(page).toHaveTitle(/Flvx Monitor Dashboard/);
    await expect(page.locator('h1')).toHaveText('Flvx Monitor Dashboard');
    
    // Check initial cards
    await expect(page.locator('text=System Status').first()).toBeVisible();
    await expect(page.locator('text=Online').first()).toBeVisible();

    await expect(page.locator('text=Flvx Node Status').first()).toBeVisible();
    await expect(page.locator('text=Alive').first()).toBeVisible();
  });

  test('should open settings, edit config, and save', async ({ page }) => {
    // Open settings dialog
    const settingsButton = page.locator('button:has(.lucide-settings)').first();
    await settingsButton.click();

    // Verify dialog title
    await expect(page.locator('text=Settings & Configuration')).toBeVisible();

    // Fill in the configuration form
    await page.fill('input#flvx_api_url', 'https://api.flvx.example.com');
    await page.fill('input#flvx_account', 'admin');
    await page.fill('input#flvx_password', 'supersecret');
    await page.fill('input#cf_token', 'cf123456');
    await page.fill('input#domain_name', 'node.test.com');
    await page.fill('input#check_api_url', 'https://check.api.test');

    // Save settings
    await page.click('text=Save Settings');

    // Ensure dialog closes
    await expect(page.locator('text=Settings & Configuration')).toBeHidden();
  });

  test('should add a standby node to the pool', async ({ page }) => {
    // Open Add Node dialog
    await page.click('text=Add Node');

    // Verify dialog
    await expect(page.locator('h2:has-text("Add Standby Node")')).toBeVisible();

    // Fill in the form
    await page.fill('input#ip', '192.168.1.100');
    await page.fill('input#port', '2222');
    await page.fill('input#password', 'nodepass123');

    // Save node
    await page.click('button:has-text("Save Node")');

    // Dialog should close
    await expect(page.locator('text=Add Standby Node')).toBeHidden();

    // Verify the node appears in the table
    // It might take a moment to fetch the nodes, wait for it
    await expect(page.locator('table')).toContainText('192.168.1.100');
    await expect(page.locator('table')).toContainText('2222');
    await expect(page.locator('table')).toContainText('standby');
  });

  test('should edit and then delete a standby node', async ({ page }) => {
    // We assume the node from the previous test or a seeded node exists.
    // Let's add one first just to be sure
    await page.click('text=Add Node');
    await page.fill('input#ip', '10.0.0.1');
    await page.fill('input#port', '22');
    await page.fill('input#password', 'temp123');
    await page.click('button:has-text("Save Node")');
    await expect(page.locator('text=Add Standby Node')).toBeHidden();

    // Now edit it
    await page.locator('tr', { hasText: '10.0.0.1' }).first().locator('button:has-text("Edit")').first().click();
    await page.fill('input#ip', '10.0.0.2');
    await page.fill('input#password', 'newpass123');
    await page.click('button:has-text("Save Node")');
    await expect(page.locator('text=Add Standby Node')).toBeHidden();
    await expect(page.locator('table')).toContainText('10.0.0.2');

    // Now delete it
    // Playwright needs to accept the confirm dialog automatically
    page.on('dialog', dialog => dialog.accept());
    await page.locator('tr', { hasText: '10.0.0.2' }).first().locator('button:has-text("Delete")').first().click();
    
    // Verify it's gone
    await expect(page.locator('table')).not.toContainText('10.0.0.2');
  });
});

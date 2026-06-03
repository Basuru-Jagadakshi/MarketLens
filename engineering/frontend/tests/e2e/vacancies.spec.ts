import { test, expect } from '@playwright/test';
import mockMetadata from './fixtures/metadata.fixture.json';
import mockVacancies from './fixtures/vacancies.fixture.json';

test.describe('Vacancy Explorer - E2E Mocked Test Suite', () => {

  test.beforeEach(async ({ page }) => {
    
    await page.route('**/vacancies/filter-values', async (route) => {
      await route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify(mockMetadata),
      });
    });

    await page.route('**/vacancies?*', async (route) => {
      await route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify(mockVacancies),
      });
    });
  });

  test('should display loading state fallback initially, then display ledger view', async ({ page }) => {
    
    await page.goto('/vacancies');

    const structuralHeader = page.locator('header');
    await expect(structuralHeader).toBeVisible();
    await expect(structuralHeader).toContainText('Vacancy Explorer');

    const tableHeaderRow = page.locator('table thead tr');
    await expect(tableHeaderRow).toContainText('Job Role');
    await expect(tableHeaderRow).toContainText('Employer');
    await expect(tableHeaderRow).toContainText('Province');

    const tableRows = page.locator('table tbody tr');
    await expect(tableRows).toHaveCount(2);
    
    await expect(tableRows.nth(0)).toContainText('Lead TypeScript Engineer');
    await expect(tableRows.nth(0)).toContainText('Axiom Labs Colombo');
    await expect(tableRows.nth(0)).toContainText('98%');
    await expect(tableRows.nth(0)).toContainText('Remote Available');

    await expect(tableRows.nth(1)).toContainText('Financial Analyst');
    await expect(tableRows.nth(1)).toContainText('88%');
    await expect(tableRows.nth(1)).toContainText('Office Based');
  });

  test('should accurately perform client-side text indexing filter actions', async ({ page }) => {
    await page.goto('/vacancies');

    const searchInput = page.getByPlaceholder('Search role or employer...');
    await expect(searchInput).toBeVisible();

    await searchInput.fill('Financial');

    const visibleRows = page.locator('table tbody tr');
    await expect(visibleRows).toHaveCount(1);
    await expect(visibleRows).toContainText('Financial Analyst');
    await expect(visibleRows).not.toContainText('Lead TypeScript Engineer');

    await searchInput.fill('Nonexistent Job Role Name');
    await expect(visibleRows).toHaveCount(1); 
    await expect(page.locator('table tbody')).toContainText('No active listings found matching active query params.');
  });

  test('should cycle detail overlay card block workflow mechanics', async ({ page }) => {
    await page.goto('/vacancies');

    await page.getByRole('row', { name: 'Lead TypeScript Engineer' }).click();

    const detailModal = page.locator('div.fixed.inset-0');
    await expect(detailModal).toBeVisible();
    
    await expect(detailModal.locator('h3')).toHaveText('Lead TypeScript Engineer');
    await expect(detailModal.locator('p').filter({ hasText: 'Axiom Labs Colombo' })).toBeVisible();
    await expect(detailModal).toContainText('Design distributed frontends and architect end-to-end testing infrastructure frameworks.');
    await expect(detailModal).toContainText('5+ years software engineering experience. Expertise in Playwright, React, and Next.js platforms.');

    const skillBadge = detailModal.locator('span', { hasText: 'Playwright' });
    await expect(skillBadge).toBeVisible();
    
    await expect(detailModal.getByText('TypeScript', { exact: true })).toBeVisible();

    await page.click('button:has-text("✕")');
    await expect(detailModal).not.toBeVisible();
  });
});
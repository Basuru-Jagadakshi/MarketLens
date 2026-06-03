import { test, expect } from '@playwright/test';
import mockMetadata from './fixtures/metadata.fixture.json';
import mockAnalytics from './fixtures/analytics.fixture.json';

test.describe('Category Analysis - E2E Mocked Test Suite', () => {

  test.beforeEach(async ({ page }) => {
    
    await page.route('**/vacancies/filter-values', async (route) => {
      await route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify(mockMetadata),
      });
    });

    await page.route('**/analytics/category/**', async (route) => {
      await route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify(mockAnalytics),
      });
    });
  });

  test('should load charts, lists, and metadata selectors successfully', async ({ page }) => {
    await page.goto('/categories'); 

    const heading = page.getByRole('heading', { name: 'Category Analysis' });
    await expect(heading).toBeVisible();

    const selectDropdown = page.locator('select');
    await expect(selectDropdown).toBeVisible();
    await expect(selectDropdown).toContainText('All Categories');

    const skillsChartHeading = page.getByRole('heading', { name: 'Top Skills by Volume' });
    await expect(skillsChartHeading).toBeVisible();

    const provinceChartHeading = page.getByRole('heading', { name: 'Province-Wise Job Vacancies' });
    await expect(provinceChartHeading).toBeVisible();

    const tableRows = page.locator('table tbody tr');
    await expect(tableRows).toHaveCount(3);
    
    await expect(tableRows.nth(0)).toContainText('TypeScript');
    await expect(tableRows.nth(0)).toContainText('145');
    await expect(tableRows.nth(2)).toContainText('Financial Modeling');
    await expect(tableRows.nth(2)).toContainText('85');

    const employerCard = page.locator('text=Axiom Labs');
    await expect(employerCard).toBeVisible();
    await expect(page.locator('text=Colombo 03')).toBeVisible();
    await expect(page.locator('text=42')).toBeVisible();
  });

  test('should handle dropdown select modifications seamlessly', async ({ page }) => {
    await page.goto('/categories');

    const selectDropdown = page.locator('select');
    
    await selectDropdown.selectOption('IT');
    
    const specificSkillsHeading = page.getByRole('heading', { name: 'Top Skills by Volume — IT' });
    await expect(specificSkillsHeading).toBeVisible();

    const categoryBadge = page.locator('span', { hasText: 'IT' });
    await expect(categoryBadge).toBeVisible();
  });

  test('should show disconnected layout UI state when analytics streaming fails', async ({ page }) => {
    
    await page.route('**/analytics/category/All*', async (route) => {
      await route.fulfill({
        status: 500,
        contentType: 'application/json',
        body: JSON.stringify({ message: 'Internal Server Error' }),
      });
    });

    await page.goto('/categories');

    const errorHeading = page.getByText('Analytics Disconnected');
    await expect(errorHeading).toBeVisible({ timeout: 15000 });

    const errorDescription = page.getByText('Failed to stream aggregated intelligence values.');
    await expect(errorDescription).toBeVisible({ timeout: 15000 });
  });
});
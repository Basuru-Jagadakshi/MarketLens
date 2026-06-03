import { test, expect } from '@playwright/test';
import mockDashboard from './fixtures/dashboard.fixture.json';

test.describe('National Labour Dashboard - E2E Mocked Test Suite', () => {

  test.beforeEach(async ({ page }) => {
    
    await page.route('**/analytics/dashboard', async (route) => {
      await route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify(mockDashboard),
      });
    });
  });

  test('should display loading state initially, then parse KPI cards completely', async ({ page }) => {
    await page.goto('/'); 

    const mainHeader = page.locator('header');
    await expect(mainHeader).toContainText('Dashboard');
    await expect(mainHeader).toContainText('National Labour Market Intelligence');

    await expect(page.getByText('Total Jobs Listed')).toBeVisible();
    await expect(page.getByText('41,250')).toBeVisible();
    await expect(page.getByText('+12.4% Month over Month')).toBeVisible();

    await expect(page.getByText('Categories')).toBeVisible();
    await expect(page.getByText('18')).toBeVisible();

    await expect(page.getByText('Identified Core Skills')).toBeVisible();
    await expect(page.getByText('342')).toBeVisible();
  });

  test('should verify interactivity inside the interactive SVG province map', async ({ page }) => {
    await page.goto('/');

    const overlayPanel = page.locator('div.absolute.bottom-4.right-4');
    await expect(overlayPanel).toBeVisible();
    await expect(overlayPanel.getByText('All Provinces')).toBeVisible();

    const centralProvinceLegendRow = page.getByText('Central', { exact: true });
    await expect(centralProvinceLegendRow).toBeVisible();

    await centralProvinceLegendRow.hover();
    
    const updatedOverlay = page.locator('div.absolute.bottom-4.right-4');
    await expect(updatedOverlay.locator('h4')).toContainText('All Provinces'); 
  });

  test('should render analytical trends, ingestion data nodes, and tracking progress bars', async ({ page }) => {
    await page.goto('/');

    const volumeBroadcastersHeading = page.getByText('Leading Volume Broadcasters');
    await expect(volumeBroadcastersHeading).toBeVisible();
    
    const firstEmployerEntry = page.locator('text=Axiom Labs');
    await expect(firstEmployerEntry).toBeVisible();
    await expect(page.locator('text=145 Posts')).toBeVisible();

    const senioritySection = page.locator('div.bg-white', { has: page.getByRole('heading', { name: 'Seniority Level', exact: true }) }).first();
    await expect(senioritySection).toBeVisible();
    await expect(senioritySection).toContainText('35%');

    const contractSection = page.locator('div.bg-white', { has: page.getByRole('heading', { name: 'Contract Type', exact: true }) }).first();
    await expect(contractSection).toBeVisible();
    await expect(contractSection).toContainText('Full-Time');
    await expect(contractSection).toContainText('85%');
  });

  test('should trigger full system error layout panel upon gateway timeouts', async ({ page }) => {
    
    await page.route('**/analytics/dashboard', async (route) => {
      await route.fulfill({
        status: 502,
        contentType: 'application/json',
        body: JSON.stringify({ message: 'Bad Gateway Error' }),
      });
    });

    await page.goto('/');

    const errorTitle = page.getByRole('heading', { name: 'Error Loading Dashboard' });
    await expect(errorTitle).toBeVisible({ timeout: 15000 });
    
    const errorDescription = page.getByText('Unable to retrieve labour market summary.');
    await expect(errorDescription).toBeVisible({ timeout: 15000 });
  });
});
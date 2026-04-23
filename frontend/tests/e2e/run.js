const { chromium, devices } = require('playwright');
const http = require('http');
const fs = require('fs');

function waitForServer(url, timeout = 15000) {
  return new Promise((resolve, reject) => {
    const start = Date.now();
    (function ping() {
      const req = http.request(url, { method: 'HEAD' }, (res) => {
        resolve();
      });
      req.on('error', () => {
        if (Date.now() - start > timeout) return reject(new Error('timeout'));
        setTimeout(ping, 200);
      });
      req.end();
    })();
  });
}

(async () => {
  const url = 'http://127.0.0.1:9528/';
  try {
    await waitForServer(url);
  } catch (e) {
    console.error('Server not available:', e);
    process.exit(3);
  }

  const browser = await chromium.launch();
  const context = await browser.newContext({ ...devices['iPhone 12'] });
  const page = await context.newPage();

  try {
    await page.goto(url, { waitUntil: 'domcontentloaded' });
    // open quote editor modal
    const addButton = await page.$('[data-action="quote-add"]');
    if (!addButton) {
      console.error('Add Quote button not found');
      await page.screenshot({ path: 'e2e-no-button.png' });
      await browser.close();
      process.exit(2);
    }
    await addButton.click();
    await page.waitForSelector('.overlay-backdrop .modal', { timeout: 5000 });
    const modalInfo = await page.evaluate(() => {
      const modal = document.querySelector('.overlay-backdrop .modal');
      if (!modal) return null;
      const rect = modal.getBoundingClientRect();
      return {
        top: rect.top,
        bottom: rect.bottom,
        height: rect.height,
        clientHeight: modal.clientHeight,
        scrollHeight: modal.scrollHeight,
        innerHeight: window.innerHeight,
        innerWidth: window.innerWidth
      };
    });
    if (!modalInfo) {
      console.error('Modal not found after click');
      await page.screenshot({ path: 'e2e-no-modal.png' });
      await browser.close();
      process.exit(2);
    }

    const fits = modalInfo.top >= 0 && modalInfo.bottom <= modalInfo.innerHeight;
    const scrollable = modalInfo.scrollHeight > modalInfo.clientHeight;

    await page.screenshot({ path: 'e2e-modal.png', fullPage: false });

    if (!fits && !scrollable) {
      console.error('Modal neither fits within viewport nor is scrollable.', modalInfo);
      process.exit(2);
    }

    console.log('E2E modal test passed', modalInfo);
    await browser.close();
    process.exit(0);
  } catch (err) {
    console.error('E2E test failed:', err);
    await page.screenshot({ path: 'e2e-error.png' });
    await browser.close();
    process.exit(1);
  }
})();

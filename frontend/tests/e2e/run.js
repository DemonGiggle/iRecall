import { chromium, devices } from 'playwright';
import http from 'node:http';

function waitForServer(url, timeout = 15000) {
  return new Promise((resolve, reject) => {
    const start = Date.now();
    (function ping() {
      const req = http.request(url, { method: 'HEAD' }, () => resolve());
      req.on('error', () => {
        if (Date.now() - start > timeout) return reject(new Error('timeout'));
        setTimeout(ping, 200);
      });
      req.end();
    })();
  });
}

function settleWithin(promise, timeout = 5000) {
  return Promise.race([
    promise,
    new Promise((resolve) => {
      const timer = setTimeout(resolve, timeout);
      timer.unref();
    }),
  ]);
}

function startFakeProvider() {
  const observed = { text: 0, images: 0 };
  const server = http.createServer((req, res) => {
    if (req.method !== 'POST' || req.url !== '/v1/chat/completions') {
      res.writeHead(404).end();
      return;
    }
    let raw = '';
    req.setEncoding('utf8');
    req.on('data', (chunk) => { raw += chunk; });
    req.on('end', () => {
      const body = JSON.parse(raw);
      const isImage = body.messages.some((message) =>
        Array.isArray(message.content) && message.content.some((part) => part.type === 'image_url'));
      if (isImage) {
        observed.images += 1;
        // Scenario 1 uses image request 1. Scenario 2 intentionally fails its second image.
        if (observed.images === 3) {
          res.writeHead(400, { 'Content-Type': 'text/plain' });
          res.end('intentional image tag failure');
          return;
        }
      } else {
        observed.text += 1;
      }
      const tags = isImage ? ['diagram', 'blueprint'] : ['distributed systems', 'consensus'];
      res.writeHead(200, { 'Content-Type': 'application/json' });
      res.end(JSON.stringify({ choices: [{ message: { content: JSON.stringify(tags) } }] }));
    });
  });
  return new Promise((resolve, reject) => {
    server.once('error', reject);
    server.listen(0, '127.0.0.1', () => resolve({ server, port: server.address().port, observed }));
  });
}

async function configureProvider(page, baseURL, port) {
  const response = await page.request.post(new URL('/api/app/save-settings', baseURL).toString(), {
    data: {
      Provider: { Host: '127.0.0.1', Port: port, HTTPS: false, APIKey: '', Model: 'e2e-model', KeywordModel: 'e2e-model' },
      Search: { MaxResults: 8, MinRelevance: 0.2 },
      Debug: { MockLLM: false },
      Theme: 'system',
      Web: { Port: 9528 },
    },
  });
  if (!response.ok()) throw new Error(`could not configure E2E provider: ${response.status()} ${await response.text()}`);
}

async function ensureProfile(page) {
  const profileNameInput = page.locator('[data-bind="profile-name"]');
  if (await profileNameInput.count()) {
    await profileNameInput.fill('E2E User');
    await page.locator('[data-action="profile-save"]').click();
    await profileNameInput.waitFor({ state: 'hidden', timeout: 10000 });
  }
}

async function openQuotes(page) {
  const quotesTab = page.locator('[data-action="nav"][data-page="Quotes"]');
  await quotesTab.waitFor({ timeout: 10000 });
  await quotesTab.click();
}

async function assertModalFitsOrScrolls(page) {
  const modalInfo = await page.locator('.overlay-backdrop .modal').evaluate((modal) => {
    const rect = modal.getBoundingClientRect();
    return {
      top: rect.top,
      bottom: rect.bottom,
      clientHeight: modal.clientHeight,
      scrollHeight: modal.scrollHeight,
      innerHeight: window.innerHeight,
    };
  });
  const fits = modalInfo.top >= 0 && modalInfo.bottom <= modalInfo.innerHeight;
  const scrollable = modalInfo.scrollHeight > modalInfo.clientHeight;
  if (!fits && !scrollable) throw new Error(`modal neither fits nor scrolls: ${JSON.stringify(modalInfo)}`);
}

const png = Buffer.from('iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAQAAAC1HAwCAAAAC0lEQVR42mP8/x8AAusB9Wl2nXQAAAAASUVORK5CYII=', 'base64');

async function addQuoteWithImages(page, content, filenames) {
  await page.locator('[data-action="quote-add"]').click();
  await page.locator('.overlay-backdrop .modal').waitFor({ timeout: 5000 });
  await assertModalFitsOrScrolls(page);
  await page.locator('[data-bind="quote-editor-content"]').fill(content);
  await page.locator('[data-bind="quote-editor-images"]').setInputFiles(
    filenames.map((name) => ({ name, mimeType: 'image/png', buffer: png })),
  );
  await page.locator('.attachment-card').filter({ hasText: 'Staged' }).first().waitFor();
  await page.locator('[data-action="quote-editor-save"]').click();
  await page.locator('.overlay-backdrop').waitFor({ state: 'hidden', timeout: 15000 });
}

async function inspectQuote(page, content, expectedImages, expectedTags) {
  const card = page.locator('.quote-card').filter({ hasText: content });
  await card.waitFor({ timeout: 10000 });
  await card.click();
  const detail = page.locator('.detail-stack');
  await detail.getByText(content, { exact: true }).waitFor();
  for (const tag of expectedTags) await detail.locator('.keyword-chip', { hasText: tag }).waitFor();
  if (await detail.locator('.attachment-card img').count() !== expectedImages) {
    throw new Error(`expected ${expectedImages} persisted images for ${content}`);
  }
  for (const image of await detail.locator('.attachment-card img').all()) {
    await image.waitFor({ state: 'visible' });
    if (!(await image.getAttribute('src'))?.includes('/api/app/quote-attachment-content')) {
      throw new Error('attachment preview does not use the persisted backend content endpoint');
    }
  }
  await page.locator('[data-action="overlay-close"]').click();
  await page.locator('.overlay-backdrop').waitFor({ state: 'hidden' });
}

(async () => {
  const url = process.env.E2E_BASE_URL || 'http://127.0.0.1:9528/';
  const fakeProvider = await startFakeProvider();
  let browser;
  try {
    await waitForServer(url);
    browser = await chromium.launch();
    const context = await browser.newContext({ ...devices['iPhone 12'] });
    const page = await context.newPage();
    await page.goto(url, { waitUntil: 'domcontentloaded' });
    await ensureProfile(page);
    await configureProvider(page, url, fakeProvider.port);
    await page.reload({ waitUntil: 'domcontentloaded' });
    await openQuotes(page);

    await addQuoteWithImages(page, 'E2E image happy path', ['happy.png']);
    await inspectQuote(page, 'E2E image happy path', 1, ['distributed systems', 'consensus', 'diagram', 'blueprint']);

    await addQuoteWithImages(page, 'E2E partial image tag failure', ['good.png', 'bad.png']);
    await page.locator('[role="status"]').filter({ hasText: 'bad.png' }).waitFor({ timeout: 5000 });
    await inspectQuote(page, 'E2E partial image tag failure', 2, ['distributed systems', 'consensus', 'diagram', 'blueprint']);

    // Reload to prove the quote, tags, attachment metadata, and attachment bytes survived persistence.
    await page.reload({ waitUntil: 'domcontentloaded' });
    await openQuotes(page);
    await inspectQuote(page, 'E2E image happy path', 1, ['distributed systems', 'consensus', 'diagram', 'blueprint']);
    await inspectQuote(page, 'E2E partial image tag failure', 2, ['distributed systems', 'consensus', 'diagram', 'blueprint']);

    if (fakeProvider.observed.text !== 2 || fakeProvider.observed.images !== 3) {
      throw new Error(`expected 2 text and 3 independent image tag calls, got ${JSON.stringify(fakeProvider.observed)}`);
    }
    await page.screenshot({ path: 'e2e-image-quotes.png', fullPage: true });
    console.log('E2E image quote tests passed', fakeProvider.observed);
  } catch (err) {
    console.error('E2E test failed:', err);
    if (browser) {
      const pages = browser.contexts().flatMap((context) => context.pages());
      if (pages[0]) await pages[0].screenshot({ path: 'e2e-error.png', fullPage: true });
    }
    process.exitCode = 1;
  } finally {
    if (browser) await settleWithin(browser.close());
    await settleWithin(new Promise((resolve) => {
      fakeProvider.server.close(resolve);
      fakeProvider.server.closeAllConnections();
    }));
    process.exit(process.exitCode ?? 0);
  }
})();

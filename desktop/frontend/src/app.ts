const navItems = ["Recall", "Quotes", "Settings"];

export function renderApp(root: HTMLElement): void {
  root.innerHTML = `
    <div class="shell">
      <header class="titlebar">
        <div class="brand">iRecall</div>
        <div class="titlebar-right">
          <div class="greeting">Hi! Desktop User</div>
          <nav class="tabs" aria-label="Primary">
            ${navItems
              .map((item, index) => `<button class="tab${index === 0 ? " active" : ""}">${item}</button>`)
              .join("")}
          </nav>
        </div>
      </header>

      <main class="layout">
        <section class="panel page-shell">
          <div class="section-header">Desktop Shell Scaffold</div>
          <p class="lede">
            This Wails frontend is intentionally lightweight. It mirrors the current iRecall product structure
            and takes its behavior from <code>docs/UI_DESIGN.md</code>.
          </p>
          <div class="cards">
            <article class="card">
              <h2>Recall</h2>
              <p>Question input, keyword line, grounded answer, and reference quotes pane.</p>
            </article>
            <article class="card">
              <h2>Quotes</h2>
              <p>Quote library with add, edit, delete, import, and export/share flows.</p>
            </article>
            <article class="card">
              <h2>Settings</h2>
              <p>Provider and search configuration with explicit save and model fetch actions.</p>
            </article>
          </div>
          <div class="panel inset">
            <div class="section-header">Next Integration Steps</div>
            <ol>
              <li>Bind the frontend shell to <code>desktop/backend.App</code> bootstrap state.</li>
              <li>Implement page modules that follow the contracts in <code>docs/UI_DESIGN.md</code>.</li>
              <li>Replace file-path text inputs with desktop-native file pickers for import/export.</li>
            </ol>
          </div>
        </section>
      </main>
    </div>
  `;
}

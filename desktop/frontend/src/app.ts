type PageName = "Recall" | "Quotes" | "History" | "Settings";
type QuoteContext = "quotes" | "recall" | "history";

interface FocusSnapshot {
  selector: string;
  selectionStart: number | null;
  selectionEnd: number | null;
}

interface Quote {
  ID: number;
  GlobalID: string;
  AuthorUserID: string;
  AuthorName: string;
  SourceUserID: string;
  SourceName: string;
  Content: string;
  Tags: string[];
  Version: number;
  IsOwnedByMe: boolean;
  CreatedAt: string;
  UpdatedAt: string;
}

interface UserProfile {
  UserID: string;
  DisplayName: string;
  CreatedAt: string;
  UpdatedAt: string;
}

interface ProviderConfig {
  Host: string;
  Port: number;
  HTTPS: boolean;
  APIKey: string;
  Model: string;
}

interface SearchConfig {
  MaxResults: number;
  MinRelevance: number;
}

interface SettingsPayload {
  Provider: ProviderConfig;
  Search: SearchConfig;
}

interface BootstrapState {
  productName: string;
  greeting: string;
  profile: UserProfile | null;
  settings: {
    Provider: ProviderConfig;
    Search: SearchConfig;
  };
  paths: {
    rootDir: string;
    dataDir: string;
    configDir: string;
    stateDir: string;
    dbPath: string;
    logPath: string;
  };
  pages: string[];
  docs: Record<string, string>;
}

interface RecallResult {
  question: string;
  keywords: string[];
  quotes: Quote[];
  response: string;
}

interface RecallHistorySummary {
  ID: number;
  Question: string;
  Response: string;
  CreatedAt: string;
}

interface RecallHistoryEntry extends RecallHistorySummary {
  Quotes: Quote[];
}

interface ImportResult {
  Inserted: number;
  Updated: number;
  Duplicates: number;
  Stale: number;
}

interface DesktopBackend {
  BootstrapState(): Promise<BootstrapState>;
  ListQuotes(): Promise<Quote[]>;
  AddQuote(content: string): Promise<Quote>;
  SaveRecallAsQuote(question: string, response: string, keywords: string[]): Promise<Quote>;
  RefineQuoteDraft(content: string): Promise<string>;
  UpdateQuote(id: number, content: string): Promise<Quote>;
  DeleteQuotes(ids: number[]): Promise<void>;
  PreviewQuoteExport(ids: number[]): Promise<string>;
  SelectQuoteExportFile(): Promise<string>;
  ExportQuotesToFile(ids: number[], path: string): Promise<void>;
  SelectQuoteImportFile(): Promise<string>;
  ImportQuotesFromFile(path: string): Promise<ImportResult>;
  SaveUserProfile(name: string): Promise<UserProfile>;
  SaveSettings(settings: SettingsPayload): Promise<SettingsPayload>;
  FetchModels(settings: ProviderConfig): Promise<string[]>;
  RunRecall(question: string): Promise<RecallResult>;
  ListRecallHistory(): Promise<RecallHistorySummary[]>;
  GetRecallHistory(id: number): Promise<RecallHistoryEntry>;
  DeleteRecallHistory(ids: number[]): Promise<void>;
}

declare global {
  interface Window {
    go?: {
      backend?: {
        App?: DesktopBackend;
      };
    };
  }
}

type OverlayState =
  | { type: "namePrompt"; name: string; busy: boolean; status: string; isError: boolean }
  | {
      type: "quoteEditor";
      mode: "add" | "edit";
      quoteId: number | null;
      content: string;
      busy: boolean;
      status: string;
      isError: boolean;
      previewOriginal: string;
      previewRefined: string;
    }
  | { type: "deleteQuotes"; context: QuoteContext; ids: number[]; busy: boolean; status: string; isError: boolean }
  | { type: "deleteHistory"; ids: number[]; busy: boolean; status: string; isError: boolean }
  | { type: "shareQuotes"; context: QuoteContext; ids: number[]; path: string; payload: string; busy: boolean; status: string; isError: boolean }
  | { type: "importQuotes"; path: string; busy: boolean; status: string; isError: boolean; result: ImportResult | null }
  | { type: "notice"; title: string; message: string };

interface SettingsFormState {
  host: string;
  port: string;
  https: boolean;
  apiKey: string;
  modelFilter: string;
  model: string;
  maxResults: string;
  minRelevance: string;
  models: string[];
}

interface AppState {
  bootstrapped: boolean;
  fatalError: string;
  page: PageName;
  bootstrap: BootstrapState | null;
  quotes: Quote[];
  quotesLoading: boolean;
  quotesError: string;
  quotesCursor: number;
  quotesSelected: Set<number>;
  recallQuestion: string;
  recallLastQuestion: string;
  recallKeywords: string[];
  recallQuotes: Quote[];
  recallResponse: string;
  recallBusy: boolean;
  recallError: string;
  recallStatus: string;
  recallStatusIsError: boolean;
  recallCursor: number;
  recallSelected: Set<number>;
  historyEntries: RecallHistorySummary[];
  historyLoading: boolean;
  historyError: string;
  historyCursor: number;
  historySelected: Set<number>;
  historyDetail: RecallHistoryEntry | null;
  historyDetailLoading: boolean;
  historyDetailError: string;
  historyStatus: string;
  historyStatusIsError: boolean;
  historyQuoteCursor: number;
  historyQuoteSelected: Set<number>;
  settings: SettingsFormState;
  settingsBusy: boolean;
  settingsStatus: string;
  settingsIsError: boolean;
  overlay: OverlayState | null;
}

const state: AppState = {
  bootstrapped: false,
  fatalError: "",
  page: "Recall",
  bootstrap: null,
  quotes: [],
  quotesLoading: false,
  quotesError: "",
  quotesCursor: 0,
  quotesSelected: new Set<number>(),
  recallQuestion: "",
  recallLastQuestion: "",
  recallKeywords: [],
  recallQuotes: [],
  recallResponse: "",
  recallBusy: false,
  recallError: "",
  recallStatus: "",
  recallStatusIsError: false,
  recallCursor: 0,
  recallSelected: new Set<number>(),
  historyEntries: [],
  historyLoading: false,
  historyError: "",
  historyCursor: 0,
  historySelected: new Set<number>(),
  historyDetail: null,
  historyDetailLoading: false,
  historyDetailError: "",
  historyStatus: "",
  historyStatusIsError: false,
  historyQuoteCursor: 0,
  historyQuoteSelected: new Set<number>(),
  settings: emptySettingsForm(),
  settingsBusy: false,
  settingsStatus: "",
  settingsIsError: false,
  overlay: null,
};

let rootEl: HTMLElement | null = null;
let handlersInstalled = false;
let bootPromise: Promise<void> | null = null;

const navItems: PageName[] = ["Recall", "History", "Quotes", "Settings"];

export function renderApp(root: HTMLElement): void {
  rootEl = root;
  if (!handlersInstalled) {
    installHandlers(root);
    handlersInstalled = true;
  }
  render();
  if (!bootPromise) {
    bootPromise = initialize();
  }
}

async function initialize(): Promise<void> {
  try {
    const bootstrap = await backend().BootstrapState();
    state.bootstrap = bootstrap;
    state.bootstrapped = true;
    state.page = "Recall";
    state.settings = settingsFormFromBootstrap(bootstrap);
    if (!bootstrap.profile?.DisplayName) {
      state.overlay = {
        type: "namePrompt",
        name: "",
        busy: false,
        status: "",
        isError: false,
      };
    }
    render();
    await loadQuotes();
  } catch (error) {
    state.bootstrapped = true;
    state.fatalError = getErrorMessage(error);
    render();
  }
}

function installHandlers(root: HTMLElement): void {
  root.addEventListener("click", (event) => {
    void handleClick(event);
  });
  root.addEventListener("input", handleInput);
  root.addEventListener("change", handleChange);
  root.addEventListener("submit", (event) => {
    void handleSubmit(event);
  });
  window.addEventListener("keydown", (event) => {
    void handleKeydown(event);
  });
}

async function handleClick(event: MouseEvent): Promise<void> {
  const target = event.target;
  if (!(target instanceof HTMLElement)) {
    return;
  }

  const actionEl = target.closest<HTMLElement>("[data-action]");
  if (!actionEl) {
    return;
  }

  const action = actionEl.dataset.action ?? "";
  switch (action) {
    case "nav":
      await switchPage(actionEl.dataset.page as PageName);
      return;
    case "quotes-refresh":
      await loadQuotes();
      return;
    case "history-refresh":
      await loadHistory();
      return;
    case "history-view-current":
      await openCurrentHistory();
      return;
    case "history-back":
      closeHistoryDetail();
      return;
    case "recall-save-quote":
      await saveRecallAsQuote();
      return;
    case "history-save-quote":
      await saveHistoryAsQuote();
      return;
    case "history-delete-current":
      openDeleteHistoryOverlay();
      return;
    case "history-select-all":
      selectAllHistory();
      return;
    case "history-deselect-all":
      clearHistorySelection();
      return;
    case "quote-add":
      openQuoteEditor("add");
      return;
    case "quote-import":
      openImportOverlay();
      return;
    case "quote-edit-current":
      openCurrentQuoteEditor(actionEl.dataset.context as QuoteContext);
      return;
    case "quote-delete-current":
      openDeleteOverlay(actionEl.dataset.context as QuoteContext);
      return;
    case "quote-share-current":
      await openShareOverlay(actionEl.dataset.context as QuoteContext);
      return;
    case "set-cursor":
      if (target.closest("input, button, label")) {
        return;
      }
      setCursor(actionEl.dataset.context as QuoteContext, Number(actionEl.dataset.index ?? "0"));
      return;
    case "history-set-cursor":
      if (target.closest("input, button, label")) {
        return;
      }
      setHistoryCursor(Number(actionEl.dataset.index ?? "0"));
      return;
    case "history-open":
      await openHistoryDetail(Number(actionEl.dataset.id ?? "0"));
      return;
    case "profile-save":
      await saveProfileName();
      return;
    case "quote-editor-save":
      await saveQuoteEditor();
      return;
    case "quote-editor-refine":
      await refineQuoteEditor();
      return;
    case "quote-editor-apply-refined":
      applyRefinedDraft();
      return;
    case "quote-editor-reject-refined":
      rejectRefinedDraft();
      return;
    case "overlay-close":
      closeOverlay();
      return;
    case "delete-confirm":
      await confirmDelete();
      return;
    case "share-browse":
      await chooseSharePath();
      return;
    case "share-save":
      await saveSharePayload();
      return;
    case "import-browse":
      await chooseImportPath();
      return;
    case "import-run":
      await importQuotes();
      return;
    case "settings-fetch-models":
      await fetchModels();
      return;
    case "settings-save":
      await saveSettings();
      return;
    case "recall-run":
      await runRecall();
      return;
    default:
      return;
  }
}

function handleInput(event: Event): void {
  const target = event.target;
  if (!(target instanceof HTMLInputElement || target instanceof HTMLTextAreaElement)) {
    return;
  }

  const bind = target.dataset.bind ?? "";
  switch (bind) {
    case "recall-question":
      state.recallQuestion = target.value;
      return;
    case "profile-name":
      if (state.overlay?.type === "namePrompt") {
        state.overlay.name = target.value;
      }
      return;
    case "quote-editor-content":
      if (state.overlay?.type === "quoteEditor") {
        state.overlay.content = target.value;
      }
      return;
    case "share-path":
      if (state.overlay?.type === "shareQuotes") {
        state.overlay.path = target.value;
      }
      return;
    case "import-path":
      if (state.overlay?.type === "importQuotes") {
        state.overlay.path = target.value;
      }
      return;
    case "settings-host":
      state.settings.host = target.value;
      return;
    case "settings-port":
      state.settings.port = target.value;
      return;
    case "settings-api-key":
      state.settings.apiKey = target.value;
      return;
    case "settings-model-filter":
      state.settings.modelFilter = target.value;
      syncSelectedModel(state.settings);
      render();
      return;
    case "settings-max-results":
      state.settings.maxResults = target.value;
      return;
    case "settings-min-relevance":
      state.settings.minRelevance = target.value;
      return;
    default:
      return;
  }
}

function handleChange(event: Event): void {
  const target = event.target;
  if (!(target instanceof HTMLInputElement || target instanceof HTMLSelectElement)) {
    return;
  }

  const bind = target.dataset.bind ?? "";
  switch (bind) {
    case "quote-selected":
      toggleSelection(target.dataset.context as QuoteContext, Number(target.dataset.id ?? "0"), (target as HTMLInputElement).checked);
      return;
    case "history-selected":
      toggleHistorySelection(Number(target.dataset.id ?? "0"), (target as HTMLInputElement).checked);
      return;
    case "settings-https":
      if (target instanceof HTMLInputElement) {
        state.settings.https = target.checked;
      }
      return;
    case "settings-model":
      state.settings.model = target.value;
      return;
    default:
      return;
  }
}

async function handleSubmit(event: Event): Promise<void> {
  const target = event.target;
  if (!(target instanceof HTMLFormElement)) {
    return;
  }
  event.preventDefault();
  switch (target.dataset.form) {
    case "recall":
      await runRecall();
      return;
    case "profile":
      await saveProfileName();
      return;
    default:
      return;
  }
}

async function handleKeydown(event: KeyboardEvent): Promise<void> {
  const active = document.activeElement;

  if (event.key === "Escape" && state.overlay && state.overlay.type !== "namePrompt") {
    event.preventDefault();
    closeOverlay();
    return;
  }

  if (event.ctrlKey && event.key.toLowerCase() === "s") {
    if (state.overlay?.type === "quoteEditor") {
      event.preventDefault();
      await saveQuoteEditor();
      return;
    }
    if (!state.overlay && state.page === "Settings") {
      event.preventDefault();
      await saveSettings();
    }
    return;
  }

  if (event.ctrlKey && event.key.toLowerCase() === "r" && state.overlay?.type === "quoteEditor") {
    event.preventDefault();
    await refineQuoteEditor();
    return;
  }

  if (
    event.key === "Enter" &&
    !event.shiftKey &&
    active instanceof HTMLInputElement &&
    active.dataset.bind === "recall-question"
  ) {
    event.preventDefault();
    await runRecall();
  }
}

async function switchPage(page: PageName): Promise<void> {
  state.page = page;
  render();
  if (page === "Quotes") {
    await loadQuotes();
  }
  if (page === "History") {
    await loadHistory();
  }
}

async function loadQuotes(): Promise<void> {
  state.quotesLoading = true;
  state.quotesError = "";
  render();
  try {
    const quotes = await backend().ListQuotes();
    state.quotes = quotes;
    state.quotesCursor = clampCursor(state.quotesCursor, quotes);
    state.quotesSelected = clampSelection(state.quotesSelected, quotes);
    state.quotesError = "";
  } catch (error) {
    state.quotesError = getErrorMessage(error);
  } finally {
    state.quotesLoading = false;
    render();
  }
}

async function loadHistory(): Promise<void> {
  state.historyLoading = true;
  state.historyError = "";
  state.historyStatus = "";
  state.historyStatusIsError = false;
  render();
  try {
    const entries = await backend().ListRecallHistory();
    state.historyEntries = entries;
    state.historyCursor = clampHistoryCursor(state.historyCursor, entries);
    state.historySelected = clampHistorySelection(state.historySelected, entries);
  } catch (error) {
    state.historyError = getErrorMessage(error);
  } finally {
    state.historyLoading = false;
    render();
  }
}

async function openCurrentHistory(): Promise<void> {
  const entry = selectedOrCurrentHistory()[0];
  if (!entry) {
    return;
  }
  await openHistoryDetail(entry.ID);
}

async function openHistoryDetail(id: number): Promise<void> {
  if (!Number.isFinite(id) || id <= 0) {
    return;
  }
  state.historyDetailLoading = true;
  state.historyDetailError = "";
  state.historyStatus = "";
  state.historyStatusIsError = false;
  render();
  try {
    const detail = await backend().GetRecallHistory(id);
    state.historyDetail = detail;
    state.historyQuoteCursor = clampCursor(state.historyQuoteCursor, detail.Quotes);
    state.historyQuoteSelected = clampSelection(state.historyQuoteSelected, detail.Quotes);
  } catch (error) {
    state.historyDetailError = getErrorMessage(error);
  } finally {
    state.historyDetailLoading = false;
    render();
  }
}

function closeHistoryDetail(): void {
  state.historyDetail = null;
  state.historyDetailLoading = false;
  state.historyDetailError = "";
  state.historyQuoteCursor = 0;
  state.historyQuoteSelected = new Set<number>();
  render();
}

function openQuoteEditor(mode: "add" | "edit", quote?: Quote): void {
  state.overlay = {
    type: "quoteEditor",
    mode,
    quoteId: quote?.ID ?? null,
    content: quote?.Content ?? "",
    busy: false,
    status: "",
    isError: false,
    previewOriginal: "",
    previewRefined: "",
  };
  render();
}

function openCurrentQuoteEditor(context: QuoteContext): void {
  const quote = selectedOrCurrentQuotes(context)[0];
  if (!quote) {
    return;
  }
  openQuoteEditor("edit", quote);
}

function openDeleteOverlay(context: QuoteContext): void {
  const ids = selectedOrCurrentQuotes(context).map((quote) => quote.ID);
  if (ids.length === 0) {
    return;
  }
  state.overlay = {
    type: "deleteQuotes",
    context,
    ids,
    busy: false,
    status: "",
    isError: false,
  };
  render();
}

function openDeleteHistoryOverlay(): void {
  const ids = selectedOrCurrentHistory().map((entry) => entry.ID);
  if (ids.length === 0) {
    return;
  }
  state.overlay = {
    type: "deleteHistory",
    ids,
    busy: false,
    status: "",
    isError: false,
  };
  render();
}

async function openShareOverlay(context: QuoteContext): Promise<void> {
  const selected = selectedOrCurrentQuotes(context);
  if (selected.length === 0) {
    return;
  }
  state.overlay = {
    type: "shareQuotes",
    context,
    ids: selected.map((quote) => quote.ID),
    path: "",
    payload: "",
    busy: true,
    status: "",
    isError: false,
  };
  render();
  try {
    const payload = await backend().PreviewQuoteExport(selected.map((quote) => quote.ID));
    if (state.overlay?.type !== "shareQuotes") {
      return;
    }
    state.overlay.payload = payload;
    state.overlay.busy = false;
    state.overlay.status = "Share payload ready. Save it to a file and transfer it manually.";
    state.overlay.isError = false;
  } catch (error) {
    if (state.overlay?.type !== "shareQuotes") {
      return;
    }
    state.overlay.busy = false;
    state.overlay.status = getErrorMessage(error);
    state.overlay.isError = true;
  }
  render();
}

function openImportOverlay(): void {
  state.overlay = {
    type: "importQuotes",
    path: "",
    busy: false,
    status: "",
    isError: false,
    result: null,
  };
  render();
}

async function saveProfileName(): Promise<void> {
  if (state.overlay?.type !== "namePrompt" || state.overlay.busy) {
    return;
  }
  const name = state.overlay.name.trim();
  if (!name) {
    state.overlay.status = "Please enter a name to continue.";
    state.overlay.isError = true;
    render();
    return;
  }

  state.overlay.busy = true;
  state.overlay.status = state.overlay.mode === "add" ? "Saving quote and generating tags..." : "Saving quote and regenerating tags...";
  state.overlay.isError = false;
  render();

  try {
    const profile = await backend().SaveUserProfile(name);
    if (state.bootstrap) {
      state.bootstrap.profile = profile;
      state.bootstrap.greeting = `Hi! ${profile.DisplayName}`;
    }
    state.overlay = null;
  } catch (error) {
    if (state.overlay?.type === "namePrompt") {
      state.overlay.busy = false;
      state.overlay.status = getErrorMessage(error);
      state.overlay.isError = true;
    }
  }
  render();
}

async function saveQuoteEditor(): Promise<void> {
  if (state.overlay?.type !== "quoteEditor" || state.overlay.busy) {
    return;
  }
  const content = state.overlay.content.trim();
  if (!content) {
    state.overlay.status = "Nothing to save.";
    state.overlay.isError = true;
    render();
    return;
  }

  state.overlay.busy = true;
  state.overlay.status = "Refining draft...";
  state.overlay.isError = false;
  render();

  try {
    const quote =
      state.overlay.mode === "add"
        ? await backend().AddQuote(content)
        : await backend().UpdateQuote(state.overlay.quoteId ?? 0, content);
    state.overlay = null;
    applyQuoteUpdate(quote);
    await loadQuotes();
  } catch (error) {
    if (state.overlay?.type === "quoteEditor") {
      state.overlay.busy = false;
      state.overlay.status = getErrorMessage(error);
      state.overlay.isError = true;
    }
    render();
  }
}

async function refineQuoteEditor(): Promise<void> {
  if (state.overlay?.type !== "quoteEditor" || state.overlay.busy) {
    return;
  }
  const content = state.overlay.content.trim();
  if (!content) {
    state.overlay.status = "Nothing to refine.";
    state.overlay.isError = true;
    render();
    return;
  }

  state.overlay.busy = true;
  state.overlay.status = "";
  render();

  try {
    const refined = await backend().RefineQuoteDraft(content);
    if (state.overlay?.type !== "quoteEditor") {
      return;
    }
    state.overlay.busy = false;
    state.overlay.previewOriginal = content;
    state.overlay.previewRefined = refined;
    state.overlay.status = "";
    state.overlay.isError = false;
  } catch (error) {
    if (state.overlay?.type === "quoteEditor") {
      state.overlay.busy = false;
      state.overlay.status = getErrorMessage(error);
      state.overlay.isError = true;
    }
  }
  render();
}

function applyRefinedDraft(): void {
  if (state.overlay?.type !== "quoteEditor") {
    return;
  }
  state.overlay.content = state.overlay.previewRefined;
  state.overlay.previewOriginal = "";
  state.overlay.previewRefined = "";
  state.overlay.status = "Refined draft applied. Review it, then save.";
  state.overlay.isError = false;
  render();
}

function rejectRefinedDraft(): void {
  if (state.overlay?.type !== "quoteEditor") {
    return;
  }
  state.overlay.previewOriginal = "";
  state.overlay.previewRefined = "";
  state.overlay.status = "Refined draft discarded.";
  state.overlay.isError = false;
  render();
}

async function confirmDelete(): Promise<void> {
  if (state.overlay?.type === "deleteHistory") {
    if (state.overlay.busy) {
      return;
    }
    state.overlay.busy = true;
    state.overlay.status = "";
    render();

    try {
      await backend().DeleteRecallHistory(state.overlay.ids);
      removeHistoryEntries(state.overlay.ids);
      state.overlay = null;
      await loadHistory();
    } catch (error) {
      if (state.overlay?.type === "deleteHistory") {
        state.overlay.busy = false;
        state.overlay.status = getErrorMessage(error);
        state.overlay.isError = true;
      }
      render();
    }
    return;
  }
  if (state.overlay?.type !== "deleteQuotes" || state.overlay.busy) {
    return;
  }
  state.overlay.busy = true;
  state.overlay.status = "";
  render();

  try {
    await backend().DeleteQuotes(state.overlay.ids);
    removeQuotes(state.overlay.ids);
    state.overlay = null;
    await loadQuotes();
  } catch (error) {
    if (state.overlay?.type === "deleteQuotes") {
      state.overlay.busy = false;
      state.overlay.status = getErrorMessage(error);
      state.overlay.isError = true;
    }
    render();
  }
}

async function chooseSharePath(): Promise<void> {
  if (state.overlay?.type !== "shareQuotes" || state.overlay.busy) {
    return;
  }
  try {
    const path = await backend().SelectQuoteExportFile();
    if (path && state.overlay?.type === "shareQuotes") {
      state.overlay.path = path;
      render();
    }
  } catch (error) {
    if (state.overlay?.type === "shareQuotes") {
      state.overlay.status = getErrorMessage(error);
      state.overlay.isError = true;
      render();
    }
  }
}

async function saveSharePayload(): Promise<void> {
  if (state.overlay?.type !== "shareQuotes" || state.overlay.busy) {
    return;
  }
  const path = state.overlay.path.trim();
  if (!path) {
    state.overlay.status = "Choose a file path for the export.";
    state.overlay.isError = true;
    render();
    return;
  }
  if (!state.overlay.payload.trim()) {
    state.overlay.status = "Export payload is not ready yet.";
    state.overlay.isError = true;
    render();
    return;
  }

  state.overlay.busy = true;
  state.overlay.status = "";
  render();

  try {
    await backend().ExportQuotesToFile(state.overlay.ids, path);
    if (state.overlay?.type === "shareQuotes") {
      state.overlay.busy = false;
      state.overlay.status = `Saved share payload to ${path}`;
      state.overlay.isError = false;
      render();
    }
  } catch (error) {
    if (state.overlay?.type === "shareQuotes") {
      state.overlay.busy = false;
      state.overlay.status = getErrorMessage(error);
      state.overlay.isError = true;
      render();
    }
  }
}

async function chooseImportPath(): Promise<void> {
  if (state.overlay?.type !== "importQuotes" || state.overlay.busy) {
    return;
  }
  try {
    const path = await backend().SelectQuoteImportFile();
    if (path && state.overlay?.type === "importQuotes") {
      state.overlay.path = path;
      render();
    }
  } catch (error) {
    if (state.overlay?.type === "importQuotes") {
      state.overlay.status = getErrorMessage(error);
      state.overlay.isError = true;
      render();
    }
  }
}

async function importQuotes(): Promise<void> {
  if (state.overlay?.type !== "importQuotes" || state.overlay.busy) {
    return;
  }
  const path = state.overlay.path.trim();
  if (!path) {
    state.overlay.status = "Choose a file to import.";
    state.overlay.isError = true;
    render();
    return;
  }

  state.overlay.busy = true;
  state.overlay.status = "";
  state.overlay.result = null;
  render();

  try {
    const result = await backend().ImportQuotesFromFile(path);
    if (state.overlay?.type !== "importQuotes") {
      return;
    }
    state.overlay.busy = false;
    state.overlay.result = result;
    state.overlay.status = `Imported quotes. inserted=${result.Inserted} updated=${result.Updated} duplicates=${result.Duplicates} stale=${result.Stale}`;
    state.overlay.isError = false;
    await loadQuotes();
  } catch (error) {
    if (state.overlay?.type === "importQuotes") {
      state.overlay.busy = false;
      state.overlay.status = getErrorMessage(error);
      state.overlay.isError = true;
      render();
    }
  }
}

async function runRecall(): Promise<void> {
  if (state.recallBusy) {
    return;
  }
  const question = state.recallQuestion.trim();
  if (!question) {
    state.recallError = "Ask a question first.";
    render();
    return;
  }

  state.recallBusy = true;
  state.recallError = "";
  state.recallStatus = "";
  state.recallStatusIsError = false;
  state.recallLastQuestion = question;
  state.recallKeywords = [];
  state.recallQuotes = [];
  state.recallResponse = "";
  state.recallCursor = 0;
  state.recallSelected = new Set<number>();
  render();

  try {
    const result = await backend().RunRecall(question);
    state.recallKeywords = result.keywords;
    state.recallQuotes = result.quotes;
    state.recallResponse = result.response;
    state.recallLastQuestion = result.question || question;
    state.recallCursor = 0;
    state.recallSelected = new Set<number>();
    state.recallQuestion = "";
  } catch (error) {
    state.recallError = getErrorMessage(error);
  } finally {
    state.recallBusy = false;
    render();
  }
}

async function saveRecallAsQuote(): Promise<void> {
  const question = state.recallLastQuestion.trim();
  const response = state.recallResponse.trim();
  if (!question || !response) {
    state.recallStatus = "Run a recall first before saving it as a quote.";
    state.recallStatusIsError = true;
    render();
    return;
  }
  try {
    const quote = await backend().SaveRecallAsQuote(question, response, state.recallKeywords);
    applyQuoteUpdate(quote);
    await loadQuotes();
    state.recallStatus = "Saved recall as quote.";
    state.recallStatusIsError = false;
    state.overlay = {
      type: "notice",
      title: "Recall Saved as Quote",
      message: "The current question and grounded response were saved as a quote with generated tags.",
    };
  } catch (error) {
    state.recallStatus = getErrorMessage(error);
    state.recallStatusIsError = true;
  }
  render();
}

async function saveHistoryAsQuote(): Promise<void> {
  const entry = state.historyDetail;
  if (!entry) {
    return;
  }
  try {
    const quote = await backend().SaveRecallAsQuote(entry.Question, entry.Response, []);
    applyQuoteUpdate(quote);
    await loadQuotes();
    state.historyStatus = "Saved history entry as quote.";
    state.historyStatusIsError = false;
    state.overlay = {
      type: "notice",
      title: "History Entry Saved as Quote",
      message: "The selected history question and response were saved as a quote with generated tags.",
    };
  } catch (error) {
    state.historyStatus = getErrorMessage(error);
    state.historyStatusIsError = true;
  }
  render();
}

async function fetchModels(): Promise<void> {
  if (state.settingsBusy) {
    return;
  }
  let provider: ProviderConfig;
  try {
    provider = providerConfigFromForm(state.settings);
  } catch (error) {
    state.settingsStatus = getErrorMessage(error);
    state.settingsIsError = true;
    render();
    return;
  }

  state.settingsBusy = true;
  state.settingsStatus = "";
  render();

  try {
    const models = await backend().FetchModels(provider);
    state.settings.models = models;
    syncSelectedModel(state.settings);
    state.settingsStatus = models.length > 0 ? `Fetched ${models.length} models.` : "No models returned.";
    state.settingsIsError = false;
  } catch (error) {
    state.settingsStatus = getErrorMessage(error);
    state.settingsIsError = true;
  } finally {
    state.settingsBusy = false;
    render();
  }
}

async function saveSettings(): Promise<void> {
  if (state.settingsBusy) {
    return;
  }
  let payload: SettingsPayload;
  try {
    payload = settingsPayloadFromForm(state.settings);
  } catch (error) {
    state.settingsStatus = getErrorMessage(error);
    state.settingsIsError = true;
    render();
    return;
  }

  state.settingsBusy = true;
  state.settingsStatus = "";
  render();

  try {
    const saved = await backend().SaveSettings(payload);
    state.settings = settingsFormFromPayload(saved, state.settings.models);
    if (state.bootstrap) {
      state.bootstrap.settings = saved;
    }
    state.settingsStatus = "Saved.";
    state.settingsIsError = false;
  } catch (error) {
    state.settingsStatus = getErrorMessage(error);
    state.settingsIsError = true;
  } finally {
    state.settingsBusy = false;
    render();
  }
}

function closeOverlay(): void {
  if (!state.overlay) {
    return;
  }
  if (state.overlay.type === "namePrompt") {
    return;
  }
  if ("busy" in state.overlay && state.overlay.busy) {
    return;
  }
  state.overlay = null;
  render();
}

function setCursor(context: QuoteContext, index: number): void {
  if (context === "quotes") {
    state.quotesCursor = clampCursor(index, state.quotes);
  } else if (context === "recall") {
    state.recallCursor = clampCursor(index, state.recallQuotes);
  } else {
    const quotes = state.historyDetail?.Quotes ?? [];
    state.historyQuoteCursor = clampCursor(index, quotes);
  }
  render();
}

function toggleSelection(context: QuoteContext, id: number, checked: boolean): void {
  const selected =
    context === "quotes" ? state.quotesSelected : context === "recall" ? state.recallSelected : state.historyQuoteSelected;
  if (checked) {
    selected.add(id);
  } else {
    selected.delete(id);
  }
}

function selectedOrCurrentQuotes(context: QuoteContext): Quote[] {
  const quotes = context === "quotes" ? state.quotes : context === "recall" ? state.recallQuotes : (state.historyDetail?.Quotes ?? []);
  const cursor = context === "quotes" ? state.quotesCursor : context === "recall" ? state.recallCursor : state.historyQuoteCursor;
  const selected =
    context === "quotes" ? state.quotesSelected : context === "recall" ? state.recallSelected : state.historyQuoteSelected;
  const chosen = quotes.filter((quote) => selected.has(quote.ID));
  if (chosen.length > 0) {
    return chosen;
  }
  return quotes[cursor] ? [quotes[cursor]] : [];
}

function applyQuoteUpdate(updated: Quote): void {
  state.quotes = patchQuoteList(state.quotes, updated);
  state.recallQuotes = patchQuoteList(state.recallQuotes, updated);
  if (state.historyDetail) {
    state.historyDetail = { ...state.historyDetail, Quotes: patchQuoteList(state.historyDetail.Quotes, updated) };
  }
  render();
}

function removeQuotes(ids: number[]): void {
  const remove = new Set(ids);
  state.quotes = state.quotes.filter((quote) => !remove.has(quote.ID));
  state.recallQuotes = state.recallQuotes.filter((quote) => !remove.has(quote.ID));
  if (state.historyDetail) {
    state.historyDetail = {
      ...state.historyDetail,
      Quotes: state.historyDetail.Quotes.filter((quote) => !remove.has(quote.ID)),
    };
  }
  state.quotesSelected = new Set([...state.quotesSelected].filter((id) => !remove.has(id)));
  state.recallSelected = new Set([...state.recallSelected].filter((id) => !remove.has(id)));
  state.historyQuoteSelected = new Set([...state.historyQuoteSelected].filter((id) => !remove.has(id)));
  state.quotesCursor = clampCursor(state.quotesCursor, state.quotes);
  state.recallCursor = clampCursor(state.recallCursor, state.recallQuotes);
  state.historyQuoteCursor = clampCursor(state.historyQuoteCursor, state.historyDetail?.Quotes ?? []);
  render();
}

function setHistoryCursor(index: number): void {
  state.historyCursor = clampHistoryCursor(index, state.historyEntries);
  render();
}

function toggleHistorySelection(id: number, checked: boolean): void {
  if (checked) {
    state.historySelected.add(id);
  } else {
    state.historySelected.delete(id);
  }
}

function selectAllHistory(): void {
  state.historySelected = new Set(state.historyEntries.map((entry) => entry.ID));
  render();
}

function clearHistorySelection(): void {
  state.historySelected = new Set<number>();
  render();
}

function selectedOrCurrentHistory(): RecallHistorySummary[] {
  const chosen = state.historyEntries.filter((entry) => state.historySelected.has(entry.ID));
  if (chosen.length > 0) {
    return chosen;
  }
  return state.historyEntries[state.historyCursor] ? [state.historyEntries[state.historyCursor]] : [];
}

function removeHistoryEntries(ids: number[]): void {
  const remove = new Set(ids);
  state.historyEntries = state.historyEntries.filter((entry) => !remove.has(entry.ID));
  state.historySelected = new Set([...state.historySelected].filter((id) => !remove.has(id)));
  state.historyCursor = clampHistoryCursor(state.historyCursor, state.historyEntries);
  if (state.historyDetail && remove.has(state.historyDetail.ID)) {
    closeHistoryDetail();
    return;
  }
  render();
}

function render(): void {
  if (!rootEl) {
    return;
  }
  const focusSnapshot = captureFocusSnapshot();
  rootEl.innerHTML = renderShell();
  restoreFocusSnapshot(focusSnapshot);
}

function captureFocusSnapshot(): FocusSnapshot | null {
  const active = document.activeElement;
  if (!(active instanceof HTMLInputElement || active instanceof HTMLTextAreaElement || active instanceof HTMLSelectElement)) {
    return null;
  }
  const bind = active.dataset.bind;
  if (!bind) {
    return null;
  }
  return {
    selector: `[data-bind="${bind}"]`,
    selectionStart: active instanceof HTMLInputElement || active instanceof HTMLTextAreaElement ? active.selectionStart : null,
    selectionEnd: active instanceof HTMLInputElement || active instanceof HTMLTextAreaElement ? active.selectionEnd : null,
  };
}

function restoreFocusSnapshot(snapshot: FocusSnapshot | null): void {
  if (!rootEl || !snapshot) {
    return;
  }
  const next = rootEl.querySelector<HTMLElement>(snapshot.selector);
  if (!(next instanceof HTMLInputElement || next instanceof HTMLTextAreaElement || next instanceof HTMLSelectElement)) {
    return;
  }
  next.focus({ preventScroll: true });
  if ((next instanceof HTMLInputElement || next instanceof HTMLTextAreaElement) && snapshot.selectionStart !== null && snapshot.selectionEnd !== null) {
    next.setSelectionRange(snapshot.selectionStart, snapshot.selectionEnd);
  }
}

function renderShell(): string {
  if (!state.bootstrapped) {
    return `
      <div class="shell shell-loading">
        <div class="splash">
          <div class="brand">iRecall</div>
          <div class="muted">Loading desktop workspace…</div>
        </div>
      </div>
    `;
  }

  if (state.fatalError) {
    return `
      <div class="shell shell-loading">
        <div class="splash splash-error">
          <div class="brand">iRecall</div>
          <div class="status status-error">${escapeHtml(state.fatalError)}</div>
        </div>
      </div>
    `;
  }

  const greeting = state.bootstrap?.profile?.DisplayName ? `Hi! ${state.bootstrap.profile.DisplayName}` : "";

  return `
    <div class="shell">
      <header class="titlebar">
        <div class="brand-lockup">
          <div class="brand">${escapeHtml(state.bootstrap?.productName ?? "iRecall")}</div>
          <div class="muted subtle">Local-first quote recall desktop</div>
        </div>
        <div class="titlebar-right">
          <div class="greeting">${escapeHtml(greeting)}</div>
          <nav class="tabs" aria-label="Primary">
            ${navItems
              .map(
                (item) => `
                  <button
                    class="tab${state.page === item ? " active" : ""}"
                    data-action="nav"
                    data-page="${item}"
                    type="button"
                  >${item}</button>
                `,
              )
              .join("")}
          </nav>
        </div>
      </header>

      <main class="layout">
        ${renderPage()}
      </main>

      ${state.overlay ? renderOverlay(state.overlay) : ""}
    </div>
  `;
}

function renderPage(): string {
  switch (state.page) {
    case "Recall":
      return renderRecallPage();
    case "Quotes":
      return renderQuotesPage();
    case "History":
      return renderHistoryPage();
    case "Settings":
      return renderSettingsPage();
  }
}

function renderRecallPage(): string {
  const selected = selectedOrCurrentQuotes("recall");
  const response = state.recallResponse.trim()
    ? escapeHtml(state.recallResponse)
    : '<span class="muted">Grounded response will appear here.</span>';
  const keywords =
    state.recallKeywords.length > 0
      ? state.recallKeywords.map((keyword) => `<span class="keyword-chip">${escapeHtml(keyword)}</span>`).join("")
      : '<span class="muted">Keywords: —</span>';

  return `
    <section class="page page-recall">
      <div class="panel page-panel">
        <div class="section-heading">
          <div>
            <div class="section-title">Recall</div>
            <div class="muted">Ask a question, ground the answer in quotes, then manage the reference set.</div>
          </div>
          <div class="toolbar">
            <button class="button" data-action="quote-add" type="button">Add Quote</button>
            <button class="button" data-action="recall-save-quote" type="button" ${!state.recallResponse.trim() ? "disabled" : ""}>Save as Quote</button>
            <button class="button" data-action="quote-edit-current" data-context="recall" type="button" ${selected.length === 0 ? "disabled" : ""}>Edit</button>
            <button class="button button-danger" data-action="quote-delete-current" data-context="recall" type="button" ${selected.length === 0 ? "disabled" : ""}>Delete</button>
            <button class="button" data-action="quote-share-current" data-context="recall" type="button" ${selected.length === 0 ? "disabled" : ""}>Share</button>
          </div>
        </div>

        <form class="question-bar" data-form="recall">
          <input
            class="text-input text-input-lg"
            data-bind="recall-question"
            placeholder="Ask anything..."
            value="${escapeAttribute(state.recallQuestion)}"
          />
          <button class="button button-primary" data-action="recall-run" type="submit" ${state.recallBusy ? "disabled" : ""}>
            ${state.recallBusy ? "Thinking…" : "Ask"}
          </button>
        </form>

        <div class="keyword-row">
          <span class="muted">Keywords:</span>
          <div class="keyword-list">${keywords}</div>
        </div>

        <div class="recall-grid">
          <section class="panel subpanel">
            <div class="subpanel-header">
              <div class="section-title">Response</div>
              <div class="muted">${state.recallBusy ? "Generating grounded answer…" : "Uses the current reference quotes."}</div>
            </div>
            <pre class="response-box">${response}</pre>
          </section>

          <section class="panel subpanel">
            <div class="subpanel-header">
              <div class="section-title">Reference Quotes</div>
              <div class="muted">${selected.length > 0 ? `${selected.length} selected` : `${state.recallQuotes.length} loaded`}</div>
            </div>
            ${renderQuoteList("recall", state.recallQuotes, state.recallCursor, state.recallSelected, false)}
          </section>
        </div>

        ${state.recallError ? `<div class="status status-error">${escapeHtml(state.recallError)}</div>` : ""}
        ${state.recallStatus ? `<div class="status ${state.recallStatusIsError ? "status-error" : "status-ok"}">${escapeHtml(state.recallStatus)}</div>` : ""}
      </div>
    </section>
  `;
}

function renderQuotesPage(): string {
  const selected = selectedOrCurrentQuotes("quotes");
  let content = "";
  if (state.quotesLoading) {
    content = '<div class="empty-state">Loading quotes…</div>';
  } else if (state.quotesError) {
    content = `<div class="status status-error">${escapeHtml(state.quotesError)}</div>`;
  } else {
    content = renderQuoteList("quotes", state.quotes, state.quotesCursor, state.quotesSelected, true);
  }

  return `
    <section class="page page-quotes">
      <div class="panel page-panel">
        <div class="section-heading">
          <div>
            <div class="section-title">Quotes</div>
            <div class="muted">Manage the local quote library, import shared payloads, and export selected notes.</div>
          </div>
          <div class="toolbar">
            <button class="button button-primary" data-action="quote-add" type="button">Add Quote</button>
            <button class="button" data-action="quote-import" type="button">Import</button>
            <button class="button" data-action="quotes-refresh" type="button">Refresh</button>
            <button class="button" data-action="quote-edit-current" data-context="quotes" type="button" ${selected.length === 0 ? "disabled" : ""}>Edit</button>
            <button class="button button-danger" data-action="quote-delete-current" data-context="quotes" type="button" ${selected.length === 0 ? "disabled" : ""}>Delete</button>
            <button class="button" data-action="quote-share-current" data-context="quotes" type="button" ${selected.length === 0 ? "disabled" : ""}>Share</button>
          </div>
        </div>
        <div class="meta-row">
          <span class="muted">Stored Quotes:</span>
          <span>${state.quotes.length}</span>
          <span class="muted">Selection:</span>
          <span>${selected.length > 0 ? selected.length : state.quotes.length > 0 ? 1 : 0}</span>
        </div>
        ${content}
      </div>
    </section>
  `;
}

function renderHistoryPage(): string {
  const selectedEntries = selectedOrCurrentHistory();
  const selectedQuotes = selectedOrCurrentQuotes("history");

  if (state.historyDetailLoading) {
    return `
      <section class="page page-history">
        <div class="panel page-panel">
          <div class="empty-state">Loading history entry…</div>
        </div>
      </section>
    `;
  }

  if (state.historyDetail) {
    return `
      <section class="page page-history">
        <div class="panel page-panel">
          <div class="section-heading">
            <div>
              <div class="section-title">History Detail</div>
              <div class="muted">Full question, response, and the exact quote set used for grounding.</div>
            </div>
            <div class="toolbar">
              <button class="button" data-action="history-back" type="button">Back</button>
              <button class="button" data-action="history-save-quote" type="button">Save as Quote</button>
              <button class="button" data-action="quote-edit-current" data-context="history" type="button" ${selectedQuotes.length === 0 ? "disabled" : ""}>Edit Quote</button>
              <button class="button button-danger" data-action="quote-delete-current" data-context="history" type="button" ${selectedQuotes.length === 0 ? "disabled" : ""}>Delete Quote</button>
              <button class="button" data-action="quote-share-current" data-context="history" type="button" ${selectedQuotes.length === 0 ? "disabled" : ""}>Share Quote</button>
            </div>
          </div>

          ${state.historyDetailError ? `<div class="status status-error">${escapeHtml(state.historyDetailError)}</div>` : ""}

          <div class="recall-grid">
            <section class="panel subpanel">
              <div class="subpanel-header">
                <div class="section-title">History Entry</div>
                <div class="muted">${escapeHtml(formatHistoryCreatedAt(state.historyDetail.CreatedAt))}</div>
              </div>
              <div class="detail-stack">
                <div class="detail-block">
                  <div class="muted">Question</div>
                  <pre class="response-box">${escapeHtml(state.historyDetail.Question)}</pre>
                </div>
                <div class="detail-block">
                  <div class="muted">Response</div>
                  <pre class="response-box">${escapeHtml(state.historyDetail.Response)}</pre>
                </div>
              </div>
            </section>

            <section class="panel subpanel">
              <div class="subpanel-header">
                <div class="section-title">Reference Quotes</div>
                <div class="muted">${selectedQuotes.length > 0 ? `${selectedQuotes.length} selected` : `${state.historyDetail.Quotes.length} loaded`}</div>
              </div>
              ${renderQuoteList("history", state.historyDetail.Quotes, state.historyQuoteCursor, state.historyQuoteSelected, false)}
            </section>
          </div>
          ${state.historyStatus ? `<div class="status ${state.historyStatusIsError ? "status-error" : "status-ok"}">${escapeHtml(state.historyStatus)}</div>` : ""}
        </div>
      </section>
    `;
  }

  let content = "";
  if (state.historyLoading) {
    content = '<div class="empty-state">Loading history…</div>';
  } else if (state.historyError) {
    content = `<div class="status status-error">${escapeHtml(state.historyError)}</div>`;
  } else if (state.historyEntries.length === 0) {
    content = '<div class="empty-state">No recall history yet. Run a recall from the Recall tab to create one.</div>';
  } else {
    content = `
      <div class="history-list">
        ${state.historyEntries
          .map((entry, index) => {
            const isCurrent = index === state.historyCursor;
            const preview = truncateQuotePreview(entry.Response, 140);
            return `
              <article class="quote-card${isCurrent ? " is-current" : ""}" data-action="history-set-cursor" data-index="${index}">
                <div class="quote-topline">
                  <label class="selection-toggle">
                    <input
                      type="checkbox"
                      data-bind="history-selected"
                      data-id="${entry.ID}"
                      ${state.historySelected.has(entry.ID) ? "checked" : ""}
                    />
                    <span>${state.historySelected.has(entry.ID) ? "[x]" : "[ ]"}</span>
                  </label>
                  <div class="quote-topline-meta">
                    <span class="quote-index${isCurrent ? " is-current" : ""}">${isCurrent ? "&gt; " : ""}[${index + 1}]</span>
                    <span class="quote-version">${escapeHtml(formatHistoryCreatedAt(entry.CreatedAt))}</span>
                  </div>
                </div>
                <div class="quote-content">${escapeHtml(truncateQuotePreview(entry.Question, 120))}</div>
                <div class="quote-meta"><span class="muted">Response:</span> <span>${escapeHtml(preview || "(empty response)")}</span></div>
                <div class="toolbar toolbar-inline">
                  <button class="button" data-action="history-open" data-id="${entry.ID}" type="button">View</button>
                </div>
              </article>
            `;
          })
          .join("")}
      </div>
    `;
  }

  return `
    <section class="page page-history">
      <div class="panel page-panel">
        <div class="section-heading">
          <div>
            <div class="section-title">History</div>
            <div class="muted">Review past recall sessions, inspect grounded responses, and manage saved history entries.</div>
          </div>
          <div class="toolbar">
            <button class="button" data-action="history-refresh" type="button">Refresh</button>
            <button class="button" data-action="history-select-all" type="button" ${state.historyEntries.length === 0 ? "disabled" : ""}>Select All</button>
            <button class="button" data-action="history-deselect-all" type="button" ${state.historySelected.size === 0 ? "disabled" : ""}>Deselect All</button>
            <button class="button" data-action="history-view-current" type="button" ${selectedEntries.length === 0 ? "disabled" : ""}>View</button>
            <button class="button button-danger" data-action="history-delete-current" type="button" ${selectedEntries.length === 0 ? "disabled" : ""}>Delete</button>
          </div>
        </div>
        <div class="meta-row">
          <span class="muted">Stored History:</span>
          <span>${state.historyEntries.length}</span>
          <span class="muted">Selection:</span>
          <span>${selectedEntries.length > 0 ? selectedEntries.length : state.historyEntries.length > 0 ? 1 : 0}</span>
        </div>
        ${content}
      </div>
    </section>
  `;
}

function renderSettingsPage(): string {
  const filteredModels = getFilteredModels(state.settings);
  const storagePaths = state.bootstrap?.paths;
  const modelSelect =
    state.settings.models.length > 0 && filteredModels.length > 0
      ? `
        <select class="select-input" data-bind="settings-model">
          ${filteredModels
            .map(
              (model) => `
                <option value="${escapeAttribute(model)}"${model === state.settings.model ? " selected" : ""}>${escapeHtml(model)}</option>
              `,
            )
            .join("")}
        </select>
      `
      : `
        <div class="readonly-model">
          <span>${escapeHtml(state.settings.model || "(none)")}</span>
          <span class="muted">${state.settings.models.length === 0 ? "Fetch models first" : "No matches"}</span>
        </div>
      `;

  return `
    <section class="page page-settings">
      <div class="panel page-panel">
        <div class="section-heading">
          <div>
            <div class="section-title">Settings</div>
            <div class="muted">Configure the OpenAI-compatible endpoint and quote retrieval behavior.</div>
          </div>
          <div class="toolbar">
            <button class="button button-primary" data-action="settings-save" type="button" ${state.settingsBusy ? "disabled" : ""}>Save</button>
          </div>
        </div>

        <div class="settings-grid">
          <section class="panel subpanel">
            <div class="section-title">LLM Provider</div>
            <label class="field">
              <span>Host / IP</span>
              <input class="text-input" data-bind="settings-host" value="${escapeAttribute(state.settings.host)}" />
            </label>
            <label class="field">
              <span>Port</span>
              <input class="text-input" data-bind="settings-port" value="${escapeAttribute(state.settings.port)}" />
            </label>
            <label class="field checkbox-field">
              <input type="checkbox" data-bind="settings-https"${state.settings.https ? " checked" : ""} />
              <span>Use HTTPS</span>
            </label>
            <label class="field">
              <span>API Key</span>
              <input class="text-input" data-bind="settings-api-key" type="password" value="${escapeAttribute(state.settings.apiKey)}" />
            </label>
            <label class="field">
              <span>Filter</span>
              <input class="text-input" data-bind="settings-model-filter" value="${escapeAttribute(state.settings.modelFilter)}" placeholder="Type to filter models" />
            </label>
            <label class="field">
              <span>Model</span>
              <div class="field-inline">
                <div class="field-inline-grow">${modelSelect}</div>
                <button class="button" data-action="settings-fetch-models" type="button" ${state.settingsBusy ? "disabled" : ""}>
                  ${state.settingsBusy ? "Fetching…" : "Fetch Models"}
                </button>
              </div>
            </label>
          </section>

          <section class="panel subpanel">
            <div class="section-title">Search</div>
            <label class="field">
              <span>Max ref quotes</span>
              <input class="text-input" data-bind="settings-max-results" value="${escapeAttribute(state.settings.maxResults)}" />
            </label>
            <label class="field">
              <span>Min relevance</span>
              <input class="text-input" data-bind="settings-min-relevance" value="${escapeAttribute(state.settings.minRelevance)}" placeholder="0.0-1.0" />
            </label>
            <div class="settings-hint muted">
              Saving updates both the persisted settings and the live engine configuration for the desktop session.
              0.0 keeps broad matches. Try 0.3-0.7 for cleaner results; 1.0 is very strict.
            </div>
          </section>

          <section class="panel subpanel">
            <div class="section-title">Local Storage</div>
            <div class="settings-paths">
              <div class="field">
                <span>Data dir</span>
                <div class="readonly-model path-value">${escapeHtml(storagePaths?.dataDir ?? "(unavailable)")}</div>
              </div>
              <div class="field">
                <span>Config dir</span>
                <div class="readonly-model path-value">${escapeHtml(storagePaths?.configDir ?? "(unavailable)")}</div>
              </div>
              <div class="field">
                <span>State dir</span>
                <div class="readonly-model path-value">${escapeHtml(storagePaths?.stateDir ?? "(unavailable)")}</div>
              </div>
              <div class="field">
                <span>Database</span>
                <div class="readonly-model path-value">${escapeHtml(storagePaths?.dbPath ?? "(unavailable)")}</div>
              </div>
            </div>
          </section>
        </div>

        ${state.settingsStatus ? `<div class="status ${state.settingsIsError ? "status-error" : "status-ok"}">${escapeHtml(state.settingsStatus)}</div>` : ""}
      </div>
    </section>
  `;
}

function renderQuoteList(
  context: QuoteContext,
  quotes: Quote[],
  cursor: number,
  selected: Set<number>,
  showTags: boolean,
): string {
  if (quotes.length === 0) {
    return `<div class="empty-state">${context === "quotes" ? "No quotes yet. Add one or import a shared payload." : "No reference quotes for this question yet."}</div>`;
  }

  return `
    <div class="quote-list">
      ${quotes
        .map((quote, index) => {
          const isCurrent = index === cursor;
          const sourceLine =
            !quote.IsOwnedByMe && quote.SourceName
              ? `<div class="quote-meta"><span class="muted">From:</span> <span class="meta-accent">${escapeHtml(quote.SourceName)}</span></div>`
              : "";
          const tagsLine = showTags
            ? `
              <div class="quote-meta">
                <span class="muted">Tags:</span>
                <span>${quote.Tags.length > 0 ? escapeHtml(previewTags(quote.Tags, 3)) : "(none)"}</span>
              </div>
            `
            : "";
          return `
            <article class="quote-card${isCurrent ? " is-current" : ""}" data-action="set-cursor" data-context="${context}" data-index="${index}">
              <div class="quote-topline">
                <label class="selection-toggle">
                  <input
                    type="checkbox"
                    data-bind="quote-selected"
                    data-context="${context}"
                    data-id="${quote.ID}"
                    ${selected.has(quote.ID) ? "checked" : ""}
                  />
                  <span>${selected.has(quote.ID) ? "[x]" : "[ ]"}</span>
                </label>
                <div class="quote-topline-meta">
                  <span class="quote-index${isCurrent ? " is-current" : ""}">${isCurrent ? "&gt; " : ""}[${index + 1}]</span>
                  <span class="quote-version">v${quote.Version}</span>
                </div>
              </div>
              <div class="quote-content">${escapeHtml(truncateQuotePreview(quote.Content, context === "quotes" ? 96 : 120))}</div>
              ${sourceLine}
              ${tagsLine}
            </article>
          `;
        })
        .join("")}
    </div>
  `;
}

function renderOverlay(overlay: OverlayState): string {
  switch (overlay.type) {
    case "namePrompt":
      return `
        <div class="overlay-backdrop">
          <div class="modal">
            <div class="modal-title">Set Your Name</div>
            <p class="modal-copy">
              Your name is attached to quotes you share and shown when other users receive your quotes.
            </p>
            <form class="modal-form" data-form="profile">
              <label class="field">
                <span>Display Name</span>
                <input class="text-input text-input-lg" data-bind="profile-name" value="${escapeAttribute(overlay.name)}" placeholder="Your name" />
              </label>
              ${overlay.status ? `<div class="status ${overlay.isError ? "status-error" : "status-ok"}">${escapeHtml(overlay.status)}</div>` : ""}
              <div class="modal-actions">
                <button class="button button-primary" data-action="profile-save" type="submit" ${overlay.busy ? "disabled" : ""}>
                  ${overlay.busy ? "Saving…" : "Save name and continue"}
                </button>
              </div>
            </form>
          </div>
        </div>
      `;
    case "quoteEditor":
      return `
        <div class="overlay-backdrop">
          <div class="modal modal-wide">
            <div class="modal-title">${overlay.mode === "add" ? "Add Quote" : "Edit Quote"}</div>
            ${
              overlay.previewRefined
                ? `
                  <div class="compare-grid">
                    <section class="panel compare-panel">
                      <div class="section-title">Current Draft</div>
                      <pre class="compare-body">${escapeHtml(overlay.previewOriginal)}</pre>
                    </section>
                    <section class="panel compare-panel">
                      <div class="section-title">Refined Draft</div>
                      <pre class="compare-body">${escapeHtml(overlay.previewRefined)}</pre>
                    </section>
                  </div>
                `
                : `
                  <label class="field">
                    <span>Quote Content</span>
                    <textarea class="text-area" data-bind="quote-editor-content" rows="10" placeholder="Type or paste your note here.">${escapeHtml(overlay.content)}</textarea>
                  </label>
                `
            }
            <div class="muted modal-copy">
              ${
                overlay.previewRefined
                  ? "Compare the current draft with the suggested rewrite before applying it."
                  : "Tags are regenerated automatically by the shared core logic."
              }
            </div>
            ${overlay.status ? `<div class="status ${overlay.isError ? "status-error" : "status-ok"}">${escapeHtml(overlay.status)}</div>` : ""}
            <div class="modal-actions">
              ${
                overlay.previewRefined
                  ? `
                    <button class="button button-primary" data-action="quote-editor-apply-refined" type="button">Apply refined draft</button>
                    <button class="button" data-action="quote-editor-reject-refined" type="button">Keep editing current draft</button>
                  `
                  : `
                    <button class="button button-primary" data-action="quote-editor-save" type="button" ${overlay.busy ? "disabled" : ""}>
                      ${overlay.busy ? "Saving…" : "Save"}
                    </button>
                    <button class="button" data-action="quote-editor-refine" type="button" ${overlay.busy ? "disabled" : ""}>
                      ${overlay.busy ? "Working…" : "Refine"}
                    </button>
                    <button class="button" data-action="overlay-close" type="button" ${overlay.busy ? "disabled" : ""}>Cancel</button>
                  `
              }
            </div>
          </div>
        </div>
      `;
    case "deleteQuotes":
      return `
        <div class="overlay-backdrop">
          <div class="modal">
            <div class="modal-title modal-title-danger">Delete Quotes</div>
            <div class="modal-copy">This permanently removes the selected quote entries from the local library.</div>
            <div class="summary-list">
              ${selectedQuotesByIds(overlay.context, overlay.ids)
                .map((quote, index) => `<div class="summary-item">[${index + 1}] ${escapeHtml(truncate(quote.Content, 140))}</div>`)
                .join("")}
            </div>
            ${overlay.status ? `<div class="status ${overlay.isError ? "status-error" : "status-ok"}">${escapeHtml(overlay.status)}</div>` : ""}
            <div class="modal-actions">
              <button class="button button-danger" data-action="delete-confirm" type="button" ${overlay.busy ? "disabled" : ""}>
                ${overlay.busy ? "Deleting…" : "Delete"}
              </button>
              <button class="button" data-action="overlay-close" type="button" ${overlay.busy ? "disabled" : ""}>Cancel</button>
            </div>
          </div>
        </div>
      `;
    case "deleteHistory":
      return `
        <div class="overlay-backdrop">
          <div class="modal">
            <div class="modal-title modal-title-danger">Delete History</div>
            <div class="modal-copy">This permanently removes the selected recall history entries from the local library.</div>
            <div class="summary-list">
              ${selectedHistoryByIds(overlay.ids)
                .map((entry, index) => `<div class="summary-item">[${index + 1}] ${escapeHtml(truncate(entry.Question, 140))}</div>`)
                .join("")}
            </div>
            ${overlay.status ? `<div class="status ${overlay.isError ? "status-error" : "status-ok"}">${escapeHtml(overlay.status)}</div>` : ""}
            <div class="modal-actions">
              <button class="button button-danger" data-action="delete-confirm" type="button" ${overlay.busy ? "disabled" : ""}>
                ${overlay.busy ? "Deleting…" : "Delete"}
              </button>
              <button class="button" data-action="overlay-close" type="button" ${overlay.busy ? "disabled" : ""}>Cancel</button>
            </div>
          </div>
        </div>
      `;
    case "shareQuotes":
      return `
        <div class="overlay-backdrop">
          <div class="modal modal-wide">
            <div class="modal-title">Share Quotes</div>
            <div class="summary-list">
              ${selectedQuotesByIds(overlay.context, overlay.ids)
                .map((quote, index) => `<div class="summary-item">[${index + 1}] v${quote.Version} ${escapeHtml(truncate(quote.Content, 120))}</div>`)
                .join("")}
            </div>
            <label class="field">
              <span>Save To</span>
              <div class="path-row">
                <input class="text-input" data-bind="share-path" value="${escapeAttribute(overlay.path)}" placeholder="/path/to/irecall-share.json" />
                <button class="button" data-action="share-browse" type="button" ${overlay.busy ? "disabled" : ""}>Browse</button>
              </div>
            </label>
            <div class="muted modal-copy">Export to a JSON file and transfer it manually to the recipient.</div>
            <div class="payload-box"><pre>${escapeHtml(overlay.payload || "Preparing export payload…")}</pre></div>
            ${overlay.status ? `<div class="status ${overlay.isError ? "status-error" : "status-ok"}">${escapeHtml(overlay.status)}</div>` : ""}
            <div class="modal-actions">
              <button class="button button-primary" data-action="share-save" type="button" ${overlay.busy ? "disabled" : ""}>
                ${overlay.busy ? "Working…" : "Save export file"}
              </button>
              <button class="button" data-action="overlay-close" type="button" ${overlay.busy ? "disabled" : ""}>Close</button>
            </div>
          </div>
        </div>
      `;
    case "importQuotes":
      return `
        <div class="overlay-backdrop">
          <div class="modal">
            <div class="modal-title">Import Quotes</div>
            <div class="modal-copy">Import a quote share JSON file exported from another iRecall instance.</div>
            <label class="field">
              <span>Import From</span>
              <div class="path-row">
                <input class="text-input" data-bind="import-path" value="${escapeAttribute(overlay.path)}" placeholder="/path/to/irecall-share.json" />
                <button class="button" data-action="import-browse" type="button" ${overlay.busy ? "disabled" : ""}>Browse</button>
              </div>
            </label>
            ${
              overlay.result
                ? `
                  <div class="result-grid">
                    <div><span class="muted">Inserted:</span> ${overlay.result.Inserted}</div>
                    <div><span class="muted">Updated:</span> ${overlay.result.Updated}</div>
                    <div><span class="muted">Duplicates:</span> ${overlay.result.Duplicates}</div>
                    <div><span class="muted">Stale:</span> ${overlay.result.Stale}</div>
                  </div>
                `
                : ""
            }
            ${overlay.status ? `<div class="status ${overlay.isError ? "status-error" : "status-ok"}">${escapeHtml(overlay.status)}</div>` : ""}
            <div class="modal-actions">
              <button class="button button-primary" data-action="import-run" type="button" ${overlay.busy ? "disabled" : ""}>
                ${overlay.busy ? "Importing…" : "Import file"}
              </button>
              <button class="button" data-action="overlay-close" type="button" ${overlay.busy ? "disabled" : ""}>Close</button>
            </div>
          </div>
        </div>
      `;
    case "notice":
      return `
        <div class="overlay-backdrop">
          <div class="modal">
            <div class="modal-title">${escapeHtml(overlay.title)}</div>
            <div class="modal-copy">${escapeHtml(overlay.message)}</div>
            <div class="modal-actions">
              <button class="button button-primary" data-action="overlay-close" type="button">OK</button>
            </div>
          </div>
        </div>
      `;
  }
}

function selectedQuotesByIds(context: QuoteContext, ids: number[]): Quote[] {
  const source = context === "quotes" ? state.quotes : context === "recall" ? state.recallQuotes : (state.historyDetail?.Quotes ?? []);
  const wanted = new Set(ids);
  return source.filter((quote) => wanted.has(quote.ID));
}

function selectedHistoryByIds(ids: number[]): RecallHistorySummary[] {
  const wanted = new Set(ids);
  return state.historyEntries.filter((entry) => wanted.has(entry.ID));
}

function settingsFormFromBootstrap(bootstrap: BootstrapState): SettingsFormState {
  return settingsFormFromPayload(bootstrap.settings, []);
}

function settingsFormFromPayload(payload: SettingsPayload | BootstrapState["settings"], models: string[]): SettingsFormState {
  const form = {
    host: payload.Provider.Host,
    port: String(payload.Provider.Port),
    https: payload.Provider.HTTPS,
    apiKey: payload.Provider.APIKey,
    modelFilter: "",
    model: payload.Provider.Model,
    maxResults: String(payload.Search.MaxResults),
    minRelevance: String(payload.Search.MinRelevance),
    models,
  };
  syncSelectedModel(form);
  return form;
}

function emptySettingsForm(): SettingsFormState {
  return {
    host: "",
    port: "11434",
    https: false,
    apiKey: "",
    modelFilter: "",
    model: "",
    maxResults: "5",
    minRelevance: "0",
    models: [],
  };
}

function providerConfigFromForm(form: SettingsFormState): ProviderConfig {
  const port = Number.parseInt(form.port.trim(), 10);
  if (!Number.isInteger(port) || port < 1 || port > 65535) {
    throw new Error("Port must be a number between 1 and 65535.");
  }
  return {
    Host: form.host.trim(),
    Port: port,
    HTTPS: form.https,
    APIKey: form.apiKey,
    Model: form.model,
  };
}

function settingsPayloadFromForm(form: SettingsFormState): SettingsPayload {
  const provider = providerConfigFromForm(form);
  const maxResults = Number.parseInt(form.maxResults.trim(), 10);
  if (!Number.isInteger(maxResults) || maxResults < 1 || maxResults > 20) {
    throw new Error("Max ref quotes must be between 1 and 20.");
  }
  const minRelevance = Number.parseFloat(form.minRelevance.trim());
  if (Number.isNaN(minRelevance)) {
    throw new Error("Min relevance must be a decimal number.");
  }
  if (minRelevance < 0 || minRelevance > 1) {
    throw new Error("Min relevance must be between 0.0 and 1.0.");
  }
  return {
    Provider: provider,
    Search: {
      MaxResults: maxResults,
      MinRelevance: minRelevance,
    },
  };
}

function getFilteredModels(form: SettingsFormState): string[] {
  const filter = form.modelFilter.trim().toLowerCase();
  if (!filter) {
    return form.models;
  }
  return form.models.filter((model) => model.toLowerCase().includes(filter));
}

function syncSelectedModel(form: SettingsFormState): void {
  if (form.models.length === 0) {
    return;
  }
  const filteredModels = getFilteredModels(form);
  if (filteredModels.length === 0) {
    return;
  }
  if (!filteredModels.includes(form.model)) {
    form.model = filteredModels[0];
  }
}

function patchQuoteList(list: Quote[], updated: Quote): Quote[] {
  return list.map((quote) => (quote.ID === updated.ID ? updated : quote));
}

function clampCursor(cursor: number, quotes: Quote[]): number {
  if (quotes.length === 0) {
    return 0;
  }
  return Math.min(Math.max(cursor, 0), quotes.length - 1);
}

function clampSelection(selected: Set<number>, quotes: Quote[]): Set<number> {
  const valid = new Set(quotes.map((quote) => quote.ID));
  return new Set([...selected].filter((id) => valid.has(id)));
}

function clampHistoryCursor(cursor: number, entries: RecallHistorySummary[]): number {
  if (entries.length === 0) {
    return 0;
  }
  return Math.min(Math.max(cursor, 0), entries.length - 1);
}

function clampHistorySelection(selected: Set<number>, entries: RecallHistorySummary[]): Set<number> {
  const valid = new Set(entries.map((entry) => entry.ID));
  return new Set([...selected].filter((id) => valid.has(id)));
}

function escapeHtml(value: string): string {
  return value
    .replaceAll("&", "&amp;")
    .replaceAll("<", "&lt;")
    .replaceAll(">", "&gt;")
    .replaceAll('"', "&quot;")
    .replaceAll("'", "&#39;");
}

function escapeAttribute(value: string): string {
  return escapeHtml(value);
}

function truncate(value: string, max: number): string {
  const normalized = value.replace(/\s+/g, " ").trim();
  if (normalized.length <= max) {
    return normalized;
  }
  return `${normalized.slice(0, max - 1).trimEnd()}…`;
}

function formatHistoryCreatedAt(value: string): string {
  const parsed = new Date(value);
  if (Number.isNaN(parsed.getTime())) {
    return value;
  }
  return parsed.toLocaleString();
}

function previewTags(tags: string[], limit: number): string {
  if (tags.length === 0) {
    return "";
  }
  if (limit <= 0 || tags.length <= limit) {
    return tags.join(" · ");
  }
  return `${tags.slice(0, limit).join(" · ")} · +${tags.length - limit} more`;
}

function truncateQuotePreview(content: string, width: number): string {
  return truncate(content, Math.max(8, width));
}

function getErrorMessage(error: unknown): string {
  if (error instanceof Error) {
    return error.message;
  }
  return String(error);
}

function backend(): DesktopBackend {
  const app = window.go?.backend?.App;
  if (!app) {
    throw new Error("Wails backend bridge is unavailable.");
  }
  return app;
}

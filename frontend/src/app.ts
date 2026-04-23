import { applyTheme, themeNames } from "./theme";

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

interface WebConfig {
  Port: number;
}

interface DebugConfig {
  MockLLM: boolean;
}

interface SettingsPayload {
  Provider: ProviderConfig;
  Search: SearchConfig;
  Debug: DebugConfig;
  Theme: string;
  Web: WebConfig;
  RootDir: string;
}

interface BootstrapState {
  productName: string;
  greeting: string;
  profile: UserProfile | null;
  settings: {
    Provider: ProviderConfig;
    Search: SearchConfig;
    Debug: DebugConfig;
    Theme: string;
    Web: WebConfig;
    RootDir: string;
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

interface APITokenStatus {
  hasToken: boolean;
  tokenPrefix: string;
}

interface APITokenCreateResult {
  token: string;
  tokenPrefix: string;
}

interface DesktopBackend {
  AuthStatus(): Promise<AuthStatus>;
  Login(password: string): Promise<void>;
  Logout(): Promise<void>;
  ChangePassword(current: string, next: string, confirm: string): Promise<void>;
  GetAPITokenStatus(): Promise<APITokenStatus>;
  CreateAPIToken(): Promise<APITokenCreateResult>;
  BootstrapState(): Promise<BootstrapState>;
  ListQuotes(): Promise<Quote[]>;
  AddQuote(content: string): Promise<Quote>;
  SaveRecallAsQuote(question: string, response: string, keywords: string[]): Promise<Quote>;
  RefineQuoteDraft(content: string): Promise<string>;
  UpdateQuote(id: number, content: string): Promise<Quote>;
  DeleteQuotes(ids: number[]): Promise<void>;
  PreviewQuoteExport(ids: number[]): Promise<string>;
  ImportQuotesPayload(payload: string): Promise<ImportResult>;
  SelectQuoteExportFile(): Promise<string>;
  ExportQuotesToFile(ids: number[], path: string): Promise<void>;
  SelectQuoteImportFile(): Promise<string>;
  ImportQuotesFromFile(path: string): Promise<ImportResult>;
  SelectRootDir(): Promise<string>;
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
      app?: {
        App?: DesktopBackend;
      };
      main?: {
        App?: DesktopBackend;
      };
    };
  }
}

interface AuthStatus {
  runtime: "desktop" | "web";
  passwordConfigured: boolean;
  authenticated: boolean;
  currentPort: number;
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
  | {
      type: "shareQuotes";
      context: QuoteContext;
      ids: number[];
      path: string;
      payload: string;
      showPayload: boolean;
      busy: boolean;
      status: string;
      isError: boolean;
    }
  | {
      type: "importQuotes";
      path: string;
      payload: string;
      filename: string;
      showPayload: boolean;
      busy: boolean;
      status: string;
      isError: boolean;
      result: ImportResult | null;
    }
  | { type: "quoteInspect"; context: QuoteContext; quote: Quote }
  | { type: "apiTokenReveal"; token: string; tokenPrefix: string }
  | { type: "notice"; title: string; message: string };

interface SettingsFormState {
  host: string;
  port: string;
  https: boolean;
  mockLLM: boolean;
  apiKey: string;
  modelFilter: string;
  model: string;
  maxResults: string;
  minRelevance: string;
  theme: string;
  webPort: string;
  rootDir: string;
  models: string[];
}

interface PasswordFormState {
  current: string;
  next: string;
  confirm: string;
  busy: boolean;
  status: string;
  isError: boolean;
}

interface APITokenState {
  loading: boolean;
  hasToken: boolean;
  tokenPrefix: string;
}

interface ToastState {
  message: string;
  isError: boolean;
}

interface AppState {
  bootstrapped: boolean;
  fatalError: string;
  authChecked: boolean;
  auth: AuthStatus | null;
  authBusy: boolean;
  authPassword: string;
  authConfirmPassword: string;
  authStatus: string;
  authIsError: boolean;
  page: PageName;
  bootstrap: BootstrapState | null;
  quotes: Quote[];
  quotesLoading: boolean;
  quotesError: string;
  quotesCursor: number;
  quotesSelected: Set<number>;
  libraryQuery: string;
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
  passwordForm: PasswordFormState;
  apiToken: APITokenState;
  overlay: OverlayState | null;
  toast: ToastState | null;
}

const state: AppState = {
  bootstrapped: false,
  fatalError: "",
  authChecked: false,
  auth: null,
  authBusy: false,
  authPassword: "",
  authConfirmPassword: "",
  authStatus: "",
  authIsError: false,
  page: "Recall",
  bootstrap: null,
  quotes: [],
  quotesLoading: false,
  quotesError: "",
  quotesCursor: 0,
  quotesSelected: new Set<number>(),
  libraryQuery: "",
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
  passwordForm: {
    current: "",
    next: "",
    confirm: "",
    busy: false,
    status: "",
    isError: false,
  },
  apiToken: {
    loading: false,
    hasToken: false,
    tokenPrefix: "",
  },
  overlay: null,
  toast: null,
};

let rootEl: HTMLElement | null = null;
let handlersInstalled = false;
let bootPromise: Promise<void> | null = null;
let toastTimer: number | null = null;

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
    const app = await waitForBackend();
    state.auth = await app.AuthStatus();
    state.authChecked = true;
    if (state.auth.runtime === "web" && !state.auth.authenticated) {
      render();
      return;
    }
    await finishBootstrap();
  } catch (error) {
    state.authChecked = true;
    state.bootstrapped = true;
    state.fatalError = getErrorMessage(error);
    render();
  }
}

async function finishBootstrap(): Promise<void> {
  await waitForBackend();
  const bootstrap = await backend().BootstrapState();
  state.bootstrap = bootstrap;
  state.bootstrapped = true;
  state.page = "Recall";
  state.settings = settingsFormFromBootstrap(bootstrap);
  applyTheme(state.settings.theme);
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
  await loadAPITokenStatus();
  await loadQuotes();
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

async function submitAuthLogin(): Promise<void> {
  if (!state.auth || state.authBusy) {
    return;
  }
  if (!state.authPassword.trim()) {
    state.authStatus = "Password is required.";
    state.authIsError = true;
    render();
    return;
  }
  state.authBusy = true;
  state.authStatus = "";
  state.authIsError = false;
  render();
  try {
    await backend().Login(state.authPassword);
    state.authPassword = "";
    state.authConfirmPassword = "";
    state.auth = await backend().AuthStatus();
    await finishBootstrap();
  } catch (error) {
    state.authStatus = getErrorMessage(error);
    state.authIsError = true;
    render();
  } finally {
    state.authBusy = false;
  }
}

async function submitAuthLogout(): Promise<void> {
  await backend().Logout();
  state.auth = await backend().AuthStatus();
  state.bootstrapped = false;
  state.bootstrap = null;
  state.overlay = null;
  state.quotes = [];
  state.historyEntries = [];
  state.historyDetail = null;
  state.authPassword = "";
  state.authConfirmPassword = "";
  state.authStatus = "";
  state.authIsError = false;
  state.apiToken = {
    loading: false,
    hasToken: false,
    tokenPrefix: "",
  };
  render();
}

async function loadAPITokenStatus(): Promise<void> {
  state.apiToken.loading = true;
  render();
  try {
    const status = await backend().GetAPITokenStatus();
    state.apiToken = {
      loading: false,
      hasToken: status.hasToken,
      tokenPrefix: status.tokenPrefix,
    };
  } catch (error) {
    state.apiToken.loading = false;
    state.settingsStatus = getErrorMessage(error);
    state.settingsIsError = true;
  }
  render();
}

async function createAPIToken(): Promise<void> {
  if (state.settingsBusy) {
    return;
  }
  const hadToken = state.apiToken.hasToken;
  state.settingsBusy = true;
  state.settingsStatus = "";
  state.settingsIsError = false;
  render();
  try {
    const result = await backend().CreateAPIToken();
    state.apiToken = {
      loading: false,
      hasToken: true,
      tokenPrefix: result.tokenPrefix,
    };
    state.overlay = {
      type: "apiTokenReveal",
      token: result.token,
      tokenPrefix: result.tokenPrefix,
    };
    state.settingsStatus = hadToken ? "API token renewed." : "API token created.";
    state.settingsIsError = false;
  } catch (error) {
    state.settingsStatus = getErrorMessage(error);
    state.settingsIsError = true;
  } finally {
    state.settingsBusy = false;
    render();
  }
}

async function submitPasswordChange(): Promise<void> {
  if (state.passwordForm.busy) {
    return;
  }
  state.passwordForm.busy = true;
  state.passwordForm.status = "";
  render();
  try {
    await backend().ChangePassword(state.passwordForm.current, state.passwordForm.next, state.passwordForm.confirm);
    state.passwordForm = {
      current: "",
      next: "",
      confirm: "",
      busy: false,
      status: "Password updated.",
      isError: false,
    };
  } catch (error) {
    state.passwordForm.busy = false;
    state.passwordForm.status = getErrorMessage(error);
    state.passwordForm.isError = true;
  }
  render();
}

async function loadImportFile(input: HTMLInputElement): Promise<void> {
  const file = input.files?.[0];
  if (!file || state.overlay?.type !== "importQuotes") {
    return;
  }
  const payload = await file.text();
  state.overlay.filename = file.name;
  state.overlay.payload = payload;
  state.overlay.path = file.name;
  state.overlay.status = `Loaded ${file.name}`;
  state.overlay.isError = false;
  render();
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
    case "auth-login":
      await submitAuthLogin();
      return;
    case "auth-logout":
      await submitAuthLogout();
      return;
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
    case "library-clear-filters":
      state.libraryQuery = "";
      state.quotesCursor = 0;
      render();
      return;
    case "quote-select-all":
      selectAllQuotes(actionEl.dataset.context as QuoteContext);
      return;
    case "quote-deselect-all":
      clearQuoteSelection(actionEl.dataset.context as QuoteContext);
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
    case "quote-inspect":
      if (target.closest("input, button, label")) {
        return;
      }
      openQuoteInspectOverlay(actionEl.dataset.context as QuoteContext, Number(actionEl.dataset.index ?? "0"));
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
    case "share-toggle-payload":
      if (state.overlay?.type === "shareQuotes") {
        state.overlay.showPayload = !state.overlay.showPayload;
        render();
      }
      return;
    case "import-toggle-payload":
      if (state.overlay?.type === "importQuotes") {
        state.overlay.showPayload = !state.overlay.showPayload;
        render();
      }
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
    case "settings-browse-root":
      await chooseRootDir();
      return;
    case "settings-clear-root":
      clearRootDir();
      return;
    case "settings-change-password":
      await submitPasswordChange();
      return;
    case "settings-create-api-token":
      await createAPIToken();
      return;
    case "recall-run":
      await runRecall();
      return;
    case "use-last-question":
      state.recallQuestion = state.recallLastQuestion;
      render();
      return;
    case "reuse-history-question":
      if (state.historyDetail) {
        state.recallQuestion = state.historyDetail.Question;
      } else {
        const current = selectedOrCurrentHistory()[0];
        if (current) {
          state.recallQuestion = current.Question;
        }
      }
      state.page = "Recall";
      render();
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
    case "auth-password":
      state.authPassword = target.value;
      return;
    case "auth-confirm-password":
      state.authConfirmPassword = target.value;
      return;
    case "recall-question":
      state.recallQuestion = target.value;
      return;
    case "library-query":
      state.libraryQuery = target.value;
      state.quotesCursor = 0;
      render();
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
    case "settings-theme":
      state.settings.theme = target.value;
      applyTheme(state.settings.theme);
      return;
    case "settings-web-port":
      state.settings.webPort = target.value;
      return;
    case "settings-root-dir":
      state.settings.rootDir = target.value;
      return;
    case "settings-password-current":
      state.passwordForm.current = target.value;
      return;
    case "settings-password-next":
      state.passwordForm.next = target.value;
      return;
    case "settings-password-confirm":
      state.passwordForm.confirm = target.value;
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
    case "settings-mock-llm":
      if (target instanceof HTMLInputElement) {
        state.settings.mockLLM = target.checked;
      }
      return;
    case "settings-model":
      state.settings.model = target.value;
      return;
    case "import-file":
      if (target instanceof HTMLInputElement) {
        void loadImportFile(target);
      }
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
    case "auth-login":
      await submitAuthLogin();
      return;
    case "auth-setup":
      await submitAuthSetup();
      return;
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
  if (page === "Settings") {
    await loadAPITokenStatus();
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
    const activeDetailId = state.historyDetail?.ID ?? null;
    state.historyEntries = entries;
    state.historyCursor = clampHistoryCursor(state.historyCursor, entries);
    state.historySelected = clampHistorySelection(state.historySelected, entries);
    if (activeDetailId === null) {
      closeHistoryDetail();
    } else if (!entries.some((entry) => entry.ID === activeDetailId)) {
      closeHistoryDetail();
    } else {
      void openHistoryDetail(activeDetailId, true);
    }
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

async function openHistoryDetail(id: number, preserveStatus = false): Promise<void> {
  if (!Number.isFinite(id) || id <= 0) {
    return;
  }
  if (state.historyDetailLoading && state.historyDetail?.ID === id) {
    return;
  }
  if (state.historyDetail && state.historyDetail.ID === id && !state.historyDetailError) {
    return;
  }
  state.historyDetailLoading = true;
  state.historyDetailError = "";
  if (!preserveStatus) {
    state.historyStatus = "";
    state.historyStatusIsError = false;
  }
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
  const quote = activeQuoteForContext(context);
  if (!quote) {
    return;
  }
  openQuoteEditor("edit", quote);
}

function openDeleteOverlay(context: QuoteContext): void {
  const ids =
    state.overlay?.type === "quoteInspect" && state.overlay.context === context
      ? [state.overlay.quote.ID]
      : selectedOrCurrentQuotes(context).map((quote) => quote.ID);
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
  const selected =
    state.overlay?.type === "quoteInspect" && state.overlay.context === context
      ? [state.overlay.quote]
      : selectedOrCurrentQuotes(context);
  if (selected.length === 0) {
    return;
  }
  state.overlay = {
    type: "shareQuotes",
    context,
    ids: selected.map((quote) => quote.ID),
    path: "",
    payload: "",
    showPayload: false,
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
    payload: "",
    filename: "",
    showPayload: false,
    busy: false,
    status: "",
    isError: false,
    result: null,
  };
  render();
}

function openQuoteInspectOverlay(context: QuoteContext, index: number): void {
  const quotes = quotesForContext(context);
  const cursor = clampCursor(index, quotes);
  const quote = quotes[cursor];
  if (!quote) {
    return;
  }
  if (context === "quotes") {
    state.quotesCursor = cursor;
  } else if (context === "recall") {
    state.recallCursor = cursor;
  } else {
    state.historyQuoteCursor = cursor;
  }
  state.overlay = { type: "quoteInspect", context, quote };
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
  state.overlay.status = "Saving profile…";
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
  if (isWebRuntime()) {
    state.overlay.path = "irecall-share.json";
    render();
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
  if (isWebRuntime()) {
    const fileName = state.overlay.path.trim() || "irecall-share.json";
    if (!state.overlay.payload.trim()) {
      state.overlay.status = "Export payload is not ready yet.";
      state.overlay.isError = true;
      render();
      return;
    }
    downloadTextFile(fileName, state.overlay.payload);
    state.overlay.status = `Downloaded ${fileName}`;
    state.overlay.isError = false;
    render();
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
  if (isWebRuntime()) {
    const input = document.querySelector<HTMLInputElement>('[data-bind="import-file"]');
    input?.click();
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
  const payload = state.overlay.payload.trim();
  if (isWebRuntime()) {
    if (!payload) {
      state.overlay.status = "Choose a file to import.";
      state.overlay.isError = true;
      render();
      return;
    }
  } else if (!path) {
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
    const result = isWebRuntime() ? await backend().ImportQuotesPayload(payload) : await backend().ImportQuotesFromFile(path);
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
    state.recallError = "Enter a recall question first.";
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
    showToast("Saved the current grounded answer as a quote.");
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
    showToast("Saved the selected activity session as a quote.");
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
    const refreshedBootstrap = await backend().BootstrapState();
    state.bootstrap = refreshedBootstrap;
    state.settings = settingsFormFromPayload(refreshedBootstrap.settings, state.settings.models);
    applyTheme(state.settings.theme);
    const needsRestart =
      state.auth?.runtime === "web" && state.auth.currentPort > 0 && state.auth.currentPort !== saved.Web.Port;
    state.settingsStatus = needsRestart ? "Saved. Restart the web server to apply the new port." : "Saved.";
    state.settingsIsError = false;
  } catch (error) {
    state.settingsStatus = getErrorMessage(error);
    state.settingsIsError = true;
  } finally {
    state.settingsBusy = false;
    render();
  }
}

async function chooseRootDir(): Promise<void> {
  if (state.settingsBusy || isWebRuntime()) {
    return;
  }
  try {
    const path = await backend().SelectRootDir();
    if (!path) {
      return;
    }
    state.settings.rootDir = path;
    state.settingsStatus = "Selected a new storage root. Save changes to apply it.";
    state.settingsIsError = false;
    render();
  } catch (error) {
    state.settingsStatus = getErrorMessage(error);
    state.settingsIsError = true;
    render();
  }
}

function clearRootDir(): void {
  state.settings.rootDir = "";
  state.settingsStatus = "Storage root cleared. Save changes to return to the default app directories.";
  state.settingsIsError = false;
  render();
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

function showToast(message: string, isError = false): void {
  state.toast = { message, isError };
  if (toastTimer !== null) {
    window.clearTimeout(toastTimer);
  }
  toastTimer = window.setTimeout(() => {
    state.toast = null;
    toastTimer = null;
    render();
  }, 2600);
}

function quotesForContext(context: QuoteContext): Quote[] {
  if (context === "quotes") {
    return getFilteredLibraryQuotes();
  }
  if (context === "recall") {
    return state.recallQuotes;
  }
  return state.historyDetail?.Quotes ?? [];
}

function setCursor(context: QuoteContext, index: number): void {
  const quotes = quotesForContext(context);
  if (context === "quotes") {
    state.quotesCursor = clampCursor(index, quotes);
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
  render();
}

function selectAllQuotes(context: QuoteContext): void {
  const quotes = quotesForContext(context);
  const selected = new Set(quotes.map((quote) => quote.ID));
  if (context === "quotes") {
    state.quotesSelected = selected;
  } else if (context === "recall") {
    state.recallSelected = selected;
  } else {
    state.historyQuoteSelected = selected;
  }
  render();
}

function clearQuoteSelection(context: QuoteContext): void {
  if (context === "quotes") {
    state.quotesSelected = new Set<number>();
  } else if (context === "recall") {
    state.recallSelected = new Set<number>();
  } else {
    state.historyQuoteSelected = new Set<number>();
  }
  render();
}

function selectedOrCurrentQuotes(context: QuoteContext): Quote[] {
  const quotes = quotesForContext(context);
  const rawCursor = context === "quotes" ? state.quotesCursor : context === "recall" ? state.recallCursor : state.historyQuoteCursor;
  const cursor = clampCursor(rawCursor, quotes);
  const selected =
    context === "quotes" ? state.quotesSelected : context === "recall" ? state.recallSelected : state.historyQuoteSelected;
  const chosen = quotes.filter((quote) => selected.has(quote.ID));
  if (chosen.length > 0) {
    return chosen;
  }
  return quotes[cursor] ? [quotes[cursor]] : [];
}

function activeQuoteForContext(context: QuoteContext): Quote | null {
  if (state.overlay?.type === "quoteInspect" && state.overlay.context === context) {
    return state.overlay.quote;
  }
  return selectedOrCurrentQuotes(context)[0] ?? null;
}

function getFilteredLibraryQuotes(): Quote[] {
  const query = state.libraryQuery.trim().toLowerCase();
  return state.quotes.filter((quote) => {
    if (!query) {
      return true;
    }
    const haystack = [quote.Content, quote.AuthorName, quote.SourceName, ...quote.Tags].join(" ").toLowerCase();
    return haystack.includes(query);
  });
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
  const current = state.historyEntries[state.historyCursor];
  if (current) {
    void openHistoryDetail(current.ID, true);
  }
  render();
}

function toggleHistorySelection(id: number, checked: boolean): void {
  if (checked) {
    state.historySelected.add(id);
  } else {
    state.historySelected.delete(id);
  }
  render();
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
  if (!state.authChecked) {
    return `
      <div class="shell shell-loading">
        <div class="splash">
          <div class="brand">iRecall</div>
          <div class="muted">Checking workspace access…</div>
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

  if (state.auth?.runtime === "web" && !state.auth.authenticated) {
    return renderAuthShell();
  }

  if (!state.bootstrapped) {
    return `
      <div class="shell shell-loading">
        <div class="splash">
          <div class="brand">iRecall</div>
          <div class="muted">Loading workspace…</div>
        </div>
      </div>
    `;
  }

  const greeting = state.bootstrap?.profile?.DisplayName ? `Hi! ${state.bootstrap.profile.DisplayName}` : "";

  return `
    <div class="shell">
      <header class="titlebar">
        <div class="brand-lockup">
          <div class="brand-row">
            <div class="brand">${escapeHtml(state.bootstrap?.productName ?? "iRecall")}</div>
            ${state.settings.mockLLM ? `<span class="meta-pill meta-pill-accent">Mock LLM on</span>` : ""}
          </div>
          <div class="muted subtle">${isWebRuntime() ? "Local-first knowledge workspace for the web" : "Local-first knowledge workspace for desktop"}</div>
        </div>
        <div class="titlebar-right">
          <div class="greeting">${escapeHtml(greeting)}</div>
          ${
            state.auth?.runtime === "web"
              ? '<button class="button" data-action="auth-logout" type="button">Logout</button>'
              : ""
          }
          <nav class="tabs" aria-label="Primary">
            ${navItems
              .map(
                (item) => `
                  <button
                    class="tab${state.page === item ? " active" : ""}"
                    data-action="nav"
                    data-page="${item}"
                    type="button"
                  >${navLabel(item)}</button>
                `,
              )
              .join("")}
          </nav>
        </div>
      </header>

      <main class="layout">
        ${renderPage()}
      </main>

      ${renderHistoryDetailModal()}
      ${state.overlay ? renderOverlay(state.overlay) : ""}
      ${state.toast ? renderToast(state.toast) : ""}
    </div>
  `;
}

function renderAuthShell(): string {
  const requiresSetup = !state.auth?.passwordConfigured;
  const action = "auth-login";
  const title = requiresSetup ? "Password Required In Terminal" : "Unlock Web UI";
  const copy = requiresSetup
    ? "The web password must be created in the terminal before the server starts listening. Restart the server from a terminal session to finish setup."
    : "Enter the web password to unlock the shared iRecall database.";

  return `
    <div class="shell shell-loading">
      <div class="panel modal">
        <div class="brand">iRecall</div>
        <div class="modal-title">${title}</div>
        <div class="modal-copy">${copy}</div>
        <form class="modal-form" data-form="${action}">
          <label class="field">
            <span>Password</span>
            <input class="text-input" data-bind="auth-password" type="password" value="${escapeAttribute(state.authPassword)}" ${requiresSetup ? "disabled" : ""} />
          </label>
          ${state.authStatus ? `<div class="status ${state.authIsError ? "status-error" : "status-ok"}">${escapeHtml(state.authStatus)}</div>` : ""}
          <div class="modal-actions">
            <button class="button button-primary" data-action="${action}" type="submit" ${state.authBusy || requiresSetup ? "disabled" : ""}>
              ${state.authBusy ? "Working…" : "Login"}
            </button>
          </div>
        </form>
      </div>
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

function navLabel(page: PageName): string {
  switch (page) {
    case "Recall":
      return "Recall";
    case "Quotes":
      return "Quotes";
    case "History":
      return "History";
    case "Settings":
      return "Settings";
  }
}

function renderRecallPage(): string {
  const responseActionsDisabled = !state.recallResponse.trim();
  const mockMode = state.settings.mockLLM;
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
        <div class="page-hero">
          <div>
            <div class="eyebrow">Recall</div>
            <div class="page-title">Question, references, then answer</div>
            <div class="muted page-copy">Run recall once, inspect the retrieved quotes, then read the grounded response.</div>
          </div>
        </div>

        <div class="flow-stack flow-stack-ask">
          <section class="panel subpanel hero-panel">
            <form class="composer-card" data-form="recall">
              <div class="section-title">1. Question</div>
              <div class="muted">Write naturally. iRecall will extract keywords and search your quotes for supporting evidence.</div>
              <textarea
                class="text-area question-input"
                data-bind="recall-question"
                rows="4"
                placeholder="What do I already know about launching the desktop app, debugging the bridge issue, and improving the UX?"
              >${escapeHtml(state.recallQuestion)}</textarea>
              <div class="composer-actions">
                <button class="button button-primary" data-action="recall-run" type="submit" ${state.recallBusy ? "disabled" : ""}>
                  ${state.recallBusy ? "Working…" : "Recall"}
                </button>
                ${
                  state.recallLastQuestion.trim()
                    ? `<button class="button" data-action="use-last-question" type="button">Use previous question</button>`
                    : ""
                }
              </div>
            </form>
          </section>

          <section class="panel subpanel">
            <div class="subpanel-header">
              <div>
                <div class="section-title">2. Reference quotes</div>
                <div class="muted">${state.recallBusy ? "Searching your quotes for relevant evidence…" : mockMode ? `${state.recallQuotes.length} retrieved quotes. Mock LLM uses simple split keywords and deterministic recall behavior.` : `${state.recallQuotes.length} retrieved quotes. Open one to inspect the full note.`}</div>
              </div>
            </div>
            <div class="keyword-row">
              <span class="muted">Keywords</span>
              <div class="keyword-list">${keywords}</div>
            </div>
            ${renderQuoteList("recall", state.recallQuotes, state.recallCursor, state.recallSelected, false)}
          </section>

          <section class="panel subpanel">
            <div class="subpanel-header">
              <div>
                <div class="section-title">3. Response</div>
                <div class="muted">${state.recallBusy ? "Writing a grounded response from the retrieved evidence…" : mockMode ? "Mock LLM combines the retrieved reference quotes into a deterministic placeholder answer." : "The response is generated from the current question and reference set."}</div>
              </div>
              <div class="toolbar toolbar-quiet">
                <button class="button button-primary" data-action="recall-save-quote" type="button" ${responseActionsDisabled ? "disabled" : ""}>Save as Quote</button>
                <button class="button" data-action="nav" data-page="History" type="button" ${responseActionsDisabled ? "disabled" : ""}>Open history</button>
              </div>
            </div>
            <div class="answer-card">
              ${
                state.recallLastQuestion.trim()
                  ? `
                    <div class="answer-anchor">
                      <div class="muted">Current question</div>
                      <div class="answer-question">${escapeHtml(state.recallLastQuestion)}</div>
                    </div>
                  `
                  : ""
              }
              <pre class="response-box">${response}</pre>
            </div>
          </section>
        </div>

        ${state.recallError ? `<div class="status status-error">${escapeHtml(state.recallError)}</div>` : ""}
        ${state.recallStatus ? `<div class="status ${state.recallStatusIsError ? "status-error" : "status-ok"}">${escapeHtml(state.recallStatus)}</div>` : ""}
      </div>
    </section>
  `;
}

function renderQuotesPage(): string {
  const filteredQuotes = getFilteredLibraryQuotes();
  const libraryCursor = clampCursor(state.quotesCursor, filteredQuotes);
  const selectedCount = filteredQuotes.filter((quote) => state.quotesSelected.has(quote.ID)).length;

  return `
    <section class="page page-quotes">
      <div class="panel page-panel">
        <div class="page-hero">
          <div>
            <div class="eyebrow">Quotes</div>
            <div class="page-title">Read, curate, and reuse your stored quotes</div>
            <div class="muted page-copy">This is the persistent source material behind grounded recall. Browse first, then edit, share, or clean up what matters.</div>
          </div>
          <div class="page-hero-actions">
            <button class="button button-primary" data-action="quote-add" type="button">New Quote</button>
            <button class="button" data-action="quote-import" type="button">Import</button>
            <button class="button" data-action="quotes-refresh" type="button">Refresh</button>
          </div>
        </div>

        <div class="workspace workspace-library">
          <section class="panel subpanel">
            <div class="subpanel-header">
              <div>
                <div class="section-title">Quote list</div>
                <div class="muted">${filteredQuotes.length} ${filteredQuotes.length === 1 ? "quote" : "quotes"}</div>
              </div>
              <div class="toolbar toolbar-quiet">
                <button class="button" data-action="quote-select-all" data-context="quotes" type="button" ${filteredQuotes.length === 0 ? "disabled" : ""}>Select all</button>
                <button class="button" data-action="quote-deselect-all" data-context="quotes" type="button" ${selectedCount === 0 ? "disabled" : ""}>Clear</button>
                <button class="button" data-action="quote-share-current" data-context="quotes" type="button" ${selectedCount === 0 ? "disabled" : ""}>Share</button>
              </div>
            </div>
            ${
              state.quotesLoading
                ? '<div class="empty-state">Loading quotes…</div>'
                : state.quotesError
                  ? `<div class="status status-error">${escapeHtml(state.quotesError)}</div>`
                  : renderQuoteList("quotes", filteredQuotes, libraryCursor, state.quotesSelected, true)
            }
          </section>
        </div>
      </div>
    </section>
  `;
}

function renderHistoryPage(): string {
  return `
    <section class="page page-history">
      <div class="panel page-panel">
        <div class="page-hero">
          <div>
            <div class="eyebrow">History</div>
            <div class="page-title">Review what you asked and what the system grounded</div>
            <div class="muted page-copy">History is your recall timeline. Use it to revisit answers, reload earlier prompts, and inspect the exact evidence that supported each response.</div>
          </div>
          <div class="page-hero-actions">
            <button class="button" data-action="history-refresh" type="button">Refresh</button>
          </div>
        </div>

        <div class="flow-stack">
          <section class="panel subpanel">
            <div class="subpanel-header">
              <div>
                <div class="section-title">History list</div>
                <div class="muted">${state.historyEntries.length} saved recall sessions</div>
              </div>
              <div class="toolbar toolbar-quiet">
                <button class="button" data-action="history-select-all" type="button" ${state.historyEntries.length === 0 ? "disabled" : ""}>Select all</button>
                <button class="button" data-action="history-deselect-all" type="button" ${state.historySelected.size === 0 ? "disabled" : ""}>Clear</button>
              </div>
            </div>
            ${
              state.historyLoading
                ? '<div class="empty-state">Loading history…</div>'
                : state.historyError
                  ? `<div class="status status-error">${escapeHtml(state.historyError)}</div>`
                  : state.historyEntries.length === 0
                    ? '<div class="empty-state">No recall history yet. Run a question from the Recall page to create your first grounded session.</div>'
                    : renderHistoryList()
            }
          </section>
        </div>
      </div>
    </section>
  `;
}

function renderHistoryDetailModal(): string {
  if (!state.historyDetailLoading && !state.historyDetail && !state.historyDetailError) {
    return "";
  }

  const detailSummary = state.historyDetail
    ? state.historyEntries.find((entry) => entry.ID === state.historyDetail?.ID) ?? state.historyDetail
    : state.historyEntries[state.historyCursor] ?? null;
  const detail = state.historyDetail && detailSummary && state.historyDetail.ID === detailSummary.ID ? state.historyDetail : null;
  const detailQuotes = detail?.Quotes ?? [];

  return `
    <div class="overlay-backdrop">
      <div class="modal modal-history-detail">
        <div class="subpanel-header">
          <div>
            <div class="modal-title">Session details</div>
            <div class="muted">${
              detailSummary
                ? detail
                  ? `${detailQuotes.length} reference quotes loaded`
                  : "Opening the selected session…"
                : "Loading the selected session…"
            }</div>
          </div>
          <div class="toolbar toolbar-quiet">
            <button class="button" data-action="history-back" type="button">Close</button>
            <button class="button button-primary" data-action="history-save-quote" type="button" ${!detailSummary ? "disabled" : ""}>Save as Quote</button>
            <button class="button" data-action="reuse-history-question" type="button" ${!detailSummary ? "disabled" : ""}>Recall again</button>
          </div>
        </div>

        ${
          detailSummary
            ? `
              <div class="detail-stack">
                <div class="detail-block">
                  <div class="muted">Question</div>
                  <pre class="response-box compact-box">${escapeHtml(detailSummary.Question)}</pre>
                </div>
                <div class="detail-block">
                  <div class="muted">Response</div>
                  <pre class="response-box compact-box">${escapeHtml(detail?.Response ?? detailSummary.Response)}</pre>
                </div>
              </div>
              ${
                state.historyDetailLoading
                  ? '<div class="empty-state">Loading reference quotes…</div>'
                  : detail
                    ? `
                      <div class="subpanel-header nested-header">
                        <div>
                          <div class="section-title">Reference quotes</div>
                          <div class="muted">${detailQuotes.length} retrieved quotes. Open one to inspect the full note.</div>
                        </div>
                      </div>
                      ${renderQuoteList("history", detailQuotes, state.historyQuoteCursor, state.historyQuoteSelected, false)}
                    `
                    : ""
              }
            `
            : ""
        }

        ${state.historyDetailError ? `<div class="status status-error">${escapeHtml(state.historyDetailError)}</div>` : ""}
        ${state.historyStatus ? `<div class="status ${state.historyStatusIsError ? "status-error" : "status-ok"}">${escapeHtml(state.historyStatus)}</div>` : ""}
      </div>
    </div>
  `;
}

function renderHistoryList(): string {
  return `
    <div class="history-list">
      ${state.historyEntries
        .map((entry, index) => {
          const isCurrent = index === state.historyCursor;
          const preview = truncateQuotePreview(entry.Response, 156);
          return `
            <article class="quote-card history-card${isCurrent ? " is-current" : ""}" data-action="history-set-cursor" data-index="${index}">
              <div class="quote-topline">
                <label class="selection-toggle">
                  <input
                    type="checkbox"
                    data-bind="history-selected"
                    data-id="${entry.ID}"
                    ${state.historySelected.has(entry.ID) ? "checked" : ""}
                  />
                </label>
                <div class="quote-topline-meta">
                  <span class="quote-version">${escapeHtml(formatHistoryCreatedAt(entry.CreatedAt))}</span>
                </div>
              </div>
              <div class="quote-content">${escapeHtml(truncateQuotePreview(entry.Question, 132))}</div>
                <div class="quote-meta"><span class="muted">Response preview</span><span>${escapeHtml(preview || "(empty response)")}</span></div>
            </article>
          `;
        })
        .join("")}
    </div>
  `;
}

function renderSettingsPage(): string {
  const filteredModels = getFilteredModels(state.settings);
  const storagePaths = state.bootstrap?.paths;
  const currentPort = state.auth?.currentPort;
  const tokenButtonLabel = state.apiToken.hasToken ? "Renew API Token" : "Create API Token";
  const tokenSummary = state.apiToken.loading
    ? "Loading token status…"
    : state.apiToken.hasToken
      ? `Active token prefix: ${state.apiToken.tokenPrefix || "(unavailable)"}`
      : "No API token has been created yet.";
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
        <div class="page-hero">
          <div>
            <div class="eyebrow">Settings</div>
            <div class="page-title">Configure connection, retrieval, and local preferences</div>
            <div class="muted page-copy">Keep the primary setup visible and move lower-level runtime controls into advanced sections.</div>
          </div>
          <div class="page-hero-actions">
            <button class="button" data-action="settings-fetch-models" type="button" ${state.settingsBusy ? "disabled" : ""}>
              ${state.settingsBusy ? "Fetching…" : "Fetch Models"}
            </button>
            <button class="button button-primary" data-action="settings-save" type="button" ${state.settingsBusy ? "disabled" : ""}>Save Changes</button>
          </div>
        </div>

        <div class="settings-layout">
          <section class="panel subpanel settings-primary">
            <div class="section-title">Connection</div>
            <label class="field">
              <span>Host / IP</span>
              <input class="text-input" data-bind="settings-host" value="${escapeAttribute(state.settings.host)}" />
            </label>
            <div class="settings-row">
              <label class="field">
                <span>Port</span>
                <input class="text-input" data-bind="settings-port" value="${escapeAttribute(state.settings.port)}" />
              </label>
              <label class="field checkbox-field settings-toggle">
                <input type="checkbox" data-bind="settings-https"${state.settings.https ? " checked" : ""} />
                <span>Use HTTPS</span>
              </label>
            </div>
            <label class="field">
              <span>API Key</span>
              <input class="text-input" data-bind="settings-api-key" type="password" value="${escapeAttribute(state.settings.apiKey)}" />
            </label>
            <label class="field">
              <span>Filter models</span>
              <input class="text-input" data-bind="settings-model-filter" value="${escapeAttribute(state.settings.modelFilter)}" placeholder="Type to narrow the model list" />
            </label>
            <label class="field">
              <span>Model</span>
              ${modelSelect}
            </label>
          </section>

          <section class="panel subpanel settings-primary">
            <div class="section-title">Retrieval</div>
            <label class="field">
              <span>Max reference quotes</span>
              <input class="text-input" data-bind="settings-max-results" value="${escapeAttribute(state.settings.maxResults)}" />
            </label>
            <label class="field">
              <span>Minimum relevance</span>
              <input class="text-input" data-bind="settings-min-relevance" value="${escapeAttribute(state.settings.minRelevance)}" placeholder="0.0-1.0" />
            </label>
            <div class="settings-hint muted">
              Lower values keep broader matches. Higher values reduce noise but risk excluding useful evidence. Most real-world sessions land in the 0.3 to 0.7 range.
            </div>
          </section>

          <section class="panel subpanel settings-secondary">
            <div class="section-title">Debug</div>
            <label class="field checkbox-field settings-toggle">
              <input type="checkbox" data-bind="settings-mock-llm"${state.settings.mockLLM ? " checked" : ""} />
              <span>Mock LLM</span>
            </label>
            <div class="settings-hint muted">
              Refine returns the original text, keywords split on spaces, and answers combine reference quotes.
            </div>
          </section>

          <section class="panel subpanel settings-secondary">
            <div class="section-title">Personalization</div>
            <label class="field">
              <span>Theme</span>
              <select class="select-input" data-bind="settings-theme">
                ${themeNames()
                  .map(
                    (theme) => `
                      <option value="${theme}"${theme === state.settings.theme ? " selected" : ""}>${theme}</option>
                    `,
                  )
                  .join("")}
              </select>
            </label>
          </section>

          <section class="panel subpanel settings-secondary">
            <div class="section-title">Security</div>
            <label class="field">
              <span>Current Password</span>
              <input class="text-input" data-bind="settings-password-current" type="password" value="${escapeAttribute(state.passwordForm.current)}" />
            </label>
            <label class="field">
              <span>New Password</span>
              <input class="text-input" data-bind="settings-password-next" type="password" value="${escapeAttribute(state.passwordForm.next)}" />
            </label>
            <label class="field">
              <span>Confirm Password</span>
              <input class="text-input" data-bind="settings-password-confirm" type="password" value="${escapeAttribute(state.passwordForm.confirm)}" />
            </label>
            <div class="muted subtle">Use at least 12 characters and include at least 3 of: uppercase, lowercase, digit, symbol.</div>
            <div class="toolbar">
              <button class="button" data-action="settings-change-password" type="button" ${state.passwordForm.busy ? "disabled" : ""}>
                ${state.passwordForm.busy ? "Updating…" : "Change Password"}
              </button>
            </div>
            ${state.passwordForm.status ? `<div class="status ${state.passwordForm.isError ? "status-error" : "status-ok"}">${escapeHtml(state.passwordForm.status)}</div>` : ""}
          </section>

          <section class="panel subpanel settings-secondary">
            <div class="section-title">REST API Token</div>
            <div class="settings-hint muted">
              Use this token for external REST clients with <code>Authorization: Bearer &lt;token&gt;</code>. The plaintext token is shown only once after creation or renewal.
            </div>
            <div class="readonly-model path-value">${escapeHtml(tokenSummary)}</div>
            <div class="toolbar">
              <button class="button" data-action="settings-create-api-token" type="button" ${state.settingsBusy || state.apiToken.loading ? "disabled" : ""}>
                ${state.settingsBusy ? "Working…" : tokenButtonLabel}
              </button>
            </div>
          </section>

          <section class="panel subpanel settings-secondary">
            <div class="section-title">Advanced</div>
            <label class="field">
              <span>Storage Root</span>
              <input
                class="text-input"
                data-bind="settings-root-dir"
                value="${escapeAttribute(state.settings.rootDir)}"
                placeholder="${escapeAttribute(isWebRuntime() ? "/path/to/irecall-root" : "Choose a folder or enter a full path")}"
              />
            </label>
            <div class="toolbar">
              ${
                isWebRuntime()
                  ? ""
                  : `<button class="button" data-action="settings-browse-root" type="button" ${state.settingsBusy ? "disabled" : ""}>Browse…</button>`
              }
              <button class="button" data-action="settings-clear-root" type="button" ${state.settingsBusy ? "disabled" : ""}>Use Default Paths</button>
            </div>
            <div class="settings-hint muted">
              ${
                isWebRuntime()
                  ? "Enter an absolute path to keep data, config, state, and the database under one custom root. Leave it empty to use the default app directories."
                  : "Choose a folder to keep data, config, state, and the database under one custom root. Leave it empty to use the default app directories."
              }
            </div>
            <label class="field">
              <span>Web Port</span>
              <input class="text-input" data-bind="settings-web-port" value="${escapeAttribute(state.settings.webPort)}" />
            </label>
            <div class="settings-hint muted">
              The web server listens on this port after restart. Current listener: ${escapeHtml(currentPort ? String(currentPort) : "(not running)")}.
            </div>
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
          const isReferenceList = context !== "quotes";
          const sourceLine =
            !quote.IsOwnedByMe && quote.SourceName
              ? `<span class="meta-accent">${escapeHtml(quote.SourceName)}</span>`
              : `<span>${escapeHtml(quote.AuthorName || "You")}</span>`;
          const tagsLine = showTags
            ? `
              <div class="quote-meta">
                <span class="muted">Tags</span>
                <span>${quote.Tags.length > 0 ? escapeHtml(previewTags(quote.Tags, 4)) : "(none)"}</span>
              </div>
            `
            : "";
          const articleAction = "quote-inspect";
          return `
            <article class="quote-card${isCurrent ? " is-current" : ""}${isReferenceList ? " quote-card-minimal" : ""}" data-action="${articleAction}" data-context="${context}" data-index="${index}">
              <div class="quote-topline">
                ${
                  isReferenceList
                    ? `<div class="quote-topline-meta">
                        <span class="quote-badge">${quote.IsOwnedByMe ? "Owned" : "Imported"}</span>
                        <span class="quote-version">${escapeHtml(formatQuoteDate(quote.UpdatedAt))}</span>
                      </div>`
                    : `
                      <label class="selection-toggle">
                        <input
                          type="checkbox"
                          data-bind="quote-selected"
                          data-context="${context}"
                          data-id="${quote.ID}"
                          ${selected.has(quote.ID) ? "checked" : ""}
                        />
                  </label>
                    `
                }
                <div class="quote-topline-meta">
                  ${isReferenceList ? `<span class="quote-source-inline">${sourceLine}</span>` : `<span class="quote-version">${escapeHtml(formatQuoteDate(quote.UpdatedAt))}</span>
                  <span class="quote-badge">${quote.IsOwnedByMe ? "Owned" : "Imported"}</span>`}
                </div>
              </div>
              <div class="quote-content">${escapeHtml(truncateQuotePreview(quote.Content, context === "quotes" ? 160 : 136))}</div>
              ${isReferenceList ? `<div class="quote-actions-inline"><button class="button button-subtle" data-action="quote-inspect" data-context="${context}" data-index="${index}" type="button">Details</button></div>` : `<div class="quote-meta"><span class="muted">${!quote.IsOwnedByMe && quote.SourceName ? "Imported from" : "Author"}</span> ${sourceLine}</div>`}
              ${tagsLine}
            </article>
          `;
        })
        .join("")}
    </div>
  `;
}

function renderQuoteDetail(quote: Quote | null, context: QuoteContext): string {
  if (!quote) {
    return '<div class="empty-state">Select a quote to inspect the full note, provenance, and available actions.</div>';
  }

  return `
    <div class="detail-stack">
      <div class="detail-block">
        <div class="muted">Full quote</div>
        <pre class="response-box compact-box">${escapeHtml(quote.Content)}</pre>
      </div>
      <div class="detail-grid">
        <div class="detail-metric">
          <span class="muted">Author</span>
          <span>${escapeHtml(quote.AuthorName || "You")}</span>
        </div>
        <div class="detail-metric">
          <span class="muted">Version</span>
          <span>v${quote.Version}</span>
        </div>
        <div class="detail-metric">
          <span class="muted">Source</span>
          <span>${escapeHtml(quote.SourceName || "Local library")}</span>
        </div>
        <div class="detail-metric">
          <span class="muted">Updated</span>
          <span>${escapeHtml(formatQuoteDate(quote.UpdatedAt))}</span>
        </div>
      </div>
      <div class="detail-block">
        <div class="muted">Tags</div>
        <div class="keyword-list">
          ${
            quote.Tags.length > 0
              ? quote.Tags.map((tag) => `<span class="keyword-chip">${escapeHtml(tag)}</span>`).join("")
              : '<span class="muted">No tags assigned yet.</span>'
          }
        </div>
      </div>
      <div class="toolbar toolbar-inline">
        <button class="button" data-action="quote-edit-current" data-context="${context}" type="button">Edit</button>
        <button class="button" data-action="quote-share-current" data-context="${context}" type="button">Share</button>
        <button class="button button-danger" data-action="quote-delete-current" data-context="${context}" type="button">Delete</button>
      </div>
    </div>
  `;
}

function formatQuoteDate(value: string): string {
  const date = new Date(value);
  if (Number.isNaN(date.getTime())) {
    return value;
  }
  return date.toLocaleDateString(undefined, { month: "short", day: "numeric", year: "numeric" });
}

function formatBytes(bytes: number): string {
  if (!Number.isFinite(bytes) || bytes <= 0) {
    return "0 B";
  }
  if (bytes < 1024) {
    return `${Math.round(bytes)} B`;
  }
  if (bytes < 1024 * 1024) {
    return `${(bytes / 1024).toFixed(1)} KB`;
  }
  return `${(bytes / (1024 * 1024)).toFixed(1)} MB`;
}

function renderToast(toast: ToastState): string {
  return `
    <div class="toast-stack" role="status" aria-live="polite">
      <div class="toast${toast.isError ? " is-error" : ""}">${escapeHtml(toast.message)}</div>
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
        <div class="overlay-backdrop overlay-backdrop-side">
          <div class="modal modal-side">
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
                  : state.settings.mockLLM
                    ? "Mock LLM is enabled, so Refine returns the original text and skips provider-dependent rewriting."
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
      const shareQuotes = selectedQuotesByIds(overlay.context, overlay.ids);
      return `
        <div class="overlay-backdrop overlay-backdrop-side">
          <div class="modal modal-side">
            <div class="modal-title">Share Quotes</div>
            <div class="modal-copy">Export a portable share file. The file summary comes first; raw JSON is available only if you need to inspect it.</div>
            <div class="summary-list">
              ${shareQuotes
                .map((quote, index) => `<div class="summary-item">[${index + 1}] v${quote.Version} ${escapeHtml(truncate(quote.Content, 120))}</div>`)
                .join("")}
            </div>
            <div class="result-grid">
              <div><span class="muted">Quotes:</span> ${shareQuotes.length}</div>
              <div><span class="muted">Payload size:</span> ${formatBytes(overlay.payload.length)}</div>
            </div>
            <label class="field">
              <span>${isWebRuntime() ? "Download As" : "Save To"}</span>
              ${
                isWebRuntime()
                  ? `<input class="text-input" data-bind="share-path" value="${escapeAttribute(overlay.path || "irecall-share.json")}" placeholder="irecall-share.json" />`
                  : `
                    <div class="path-row">
                      <input class="text-input" data-bind="share-path" value="${escapeAttribute(overlay.path)}" placeholder="/path/to/irecall-share.json" />
                      <button class="button" data-action="share-browse" type="button" ${overlay.busy ? "disabled" : ""}>Browse</button>
                    </div>
                  `
              }
            </label>
            <div class="muted modal-copy">${
              isWebRuntime()
                ? "Download the JSON payload locally, then transfer it manually to the recipient."
                : "Export to a JSON file and transfer it manually to the recipient."
            }</div>
            <div class="toolbar toolbar-inline">
              <button class="button" data-action="share-toggle-payload" type="button" ${!overlay.payload ? "disabled" : ""}>
                ${overlay.showPayload ? "Hide raw JSON" : "Show raw JSON"}
              </button>
            </div>
            ${overlay.showPayload ? `<div class="payload-box"><pre>${escapeHtml(overlay.payload || "Preparing export payload…")}</pre></div>` : ""}
            ${overlay.status ? `<div class="status ${overlay.isError ? "status-error" : "status-ok"}">${escapeHtml(overlay.status)}</div>` : ""}
            <div class="modal-actions">
              <button class="button button-primary" data-action="share-save" type="button" ${overlay.busy ? "disabled" : ""}>
                ${overlay.busy ? "Working…" : isWebRuntime() ? "Download export file" : "Save export file"}
              </button>
              <button class="button" data-action="overlay-close" type="button" ${overlay.busy ? "disabled" : ""}>Close</button>
            </div>
          </div>
        </div>
      `;
    case "importQuotes":
      return `
        <div class="overlay-backdrop overlay-backdrop-side">
          <div class="modal modal-side">
            <div class="modal-title">Import Quotes</div>
            <div class="modal-copy">Import a quote share file exported from another iRecall instance. Start by choosing a file, then review the result summary.</div>
            <label class="field">
              <span>Import From</span>
              ${
                isWebRuntime()
                  ? `
                    <input class="text-input" data-bind="import-path" value="${escapeAttribute(overlay.path)}" placeholder="Choose a local JSON file" readonly />
                    <input data-bind="import-file" type="file" accept="application/json,.json" hidden />
                    <div class="toolbar">
                      <button class="button" data-action="import-browse" type="button" ${overlay.busy ? "disabled" : ""}>Choose File</button>
                    </div>
                  `
                  : `
                    <div class="path-row">
                      <input class="text-input" data-bind="import-path" value="${escapeAttribute(overlay.path)}" placeholder="/path/to/irecall-share.json" />
                      <button class="button" data-action="import-browse" type="button" ${overlay.busy ? "disabled" : ""}>Browse</button>
                    </div>
                `
              }
            </label>
            ${
              overlay.filename || overlay.path
                ? `
                  <div class="result-grid">
                    <div><span class="muted">File:</span> ${escapeHtml(overlay.filename || overlay.path)}</div>
                    <div><span class="muted">Payload size:</span> ${formatBytes(overlay.payload.length)}</div>
                  </div>
                `
                : ""
            }
            ${
              overlay.payload
                ? `
                  <div class="toolbar toolbar-inline">
                    <button class="button" data-action="import-toggle-payload" type="button">
                      ${overlay.showPayload ? "Hide raw JSON" : "Show raw JSON"}
                    </button>
                  </div>
                  ${overlay.showPayload ? `<div class="payload-box"><pre>${escapeHtml(overlay.payload)}</pre></div>` : ""}
                `
                : ""
            }
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
    case "quoteInspect":
      return `
        <div class="overlay-backdrop">
          <div class="modal modal-quote-inspect">
            <div class="modal-title">Quote details</div>
            <p class="modal-copy">
              Review the full note and its provenance without leaving the current flow.
            </p>
            ${renderQuoteDetail(overlay.quote, overlay.context)}
            <div class="modal-actions">
              <button class="button" data-action="overlay-close" type="button">Close</button>
            </div>
          </div>
        </div>
      `;
    case "apiTokenReveal":
      return `
        <div class="overlay-backdrop">
          <div class="modal modal-side">
            <div class="modal-title">Copy API Token Now</div>
            <p class="modal-copy">
              This token is shown only once. Copy it now and store it safely. After you close this dialog, only the prefix
              <strong> ${escapeHtml(overlay.tokenPrefix)} </strong>
              will remain visible in Settings.
            </p>
            <div class="payload-box"><pre>${escapeHtml(overlay.token)}</pre></div>
            <div class="muted modal-copy">Use it with <code>Authorization: Bearer &lt;token&gt;</code> on REST API requests.</div>
            <div class="modal-actions">
              <button class="button" data-action="overlay-close" type="button">Close</button>
            </div>
          </div>
        </div>
      `;
    case "notice":
      return "";
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
    mockLLM: payload.Debug?.MockLLM ?? false,
    apiKey: payload.Provider.APIKey,
    modelFilter: "",
    model: payload.Provider.Model,
    maxResults: String(payload.Search.MaxResults),
    minRelevance: String(payload.Search.MinRelevance),
    theme: payload.Theme || "violet",
    webPort: String(payload.Web?.Port ?? 9527),
    rootDir: payload.RootDir ?? "",
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
    mockLLM: false,
    apiKey: "",
    modelFilter: "",
    model: "",
    maxResults: "5",
    minRelevance: "0",
    theme: "violet",
    webPort: "9527",
    rootDir: "",
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
  const webPort = Number.parseInt(form.webPort.trim(), 10);
  if (!Number.isInteger(maxResults) || maxResults < 1 || maxResults > 20) {
    throw new Error("Max ref quotes must be between 1 and 20.");
  }
  if (!Number.isInteger(webPort) || webPort < 1 || webPort > 65535) {
    throw new Error("Web port must be a number between 1 and 65535.");
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
    Debug: {
      MockLLM: form.mockLLM,
    },
    Theme: form.theme,
    Web: {
      Port: webPort,
    },
    RootDir: form.rootDir.trim(),
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

function isWebRuntime(): boolean {
  return state.auth?.runtime === "web";
}

function downloadTextFile(fileName: string, content: string): void {
  const blob = new Blob([content], { type: "application/json;charset=utf-8" });
  const url = URL.createObjectURL(blob);
  const link = document.createElement("a");
  link.href = url;
  link.download = fileName;
  document.body.appendChild(link);
  link.click();
  link.remove();
  URL.revokeObjectURL(url);
}

function getErrorMessage(error: unknown): string {
  if (error instanceof Error) {
    return error.message;
  }
  return String(error);
}

function resolveBackend(): DesktopBackend | null {
  const namespaces = [window.go?.backend, window.go?.app, window.go?.main];
  for (const namespace of namespaces) {
    if (namespace?.App) {
      return namespace.App;
    }
  }
  return null;
}

async function waitForBackend(timeoutMs = 3000): Promise<DesktopBackend> {
  const start = Date.now();
  for (;;) {
    const app = resolveBackend();
    if (app) {
      return app;
    }
    if (Date.now() - start >= timeoutMs) {
      throw new Error("Wails backend bridge is unavailable.");
    }
    await new Promise((resolve) => window.setTimeout(resolve, 25));
  }
}

function backend(): DesktopBackend {
  const app = resolveBackend();
  if (!app) {
    throw new Error("Wails backend bridge is unavailable.");
  }
  return app;
}

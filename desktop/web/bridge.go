package web

const webBridgeJS = `
(function () {
  async function request(method, url, body) {
    const response = await fetch(url, {
      method,
      credentials: "same-origin",
      headers: body === undefined ? {} : { "Content-Type": "application/json" },
      body: body === undefined ? undefined : JSON.stringify(body),
    });
    const contentType = response.headers.get("content-type") || "";
    const payload = contentType.includes("application/json") ? await response.json() : await response.text();
    if (!response.ok) {
      const message = payload && typeof payload === "object" && payload.error ? payload.error : String(payload || response.statusText);
      throw new Error(message);
    }
    return payload;
  }

	const App = {
	    AuthStatus() {
	      return request("GET", "/api/auth/status");
	    },
	    Login(password) {
	      return request("POST", "/api/auth/login", { password });
	    },
    Logout() {
      return request("POST", "/api/auth/logout");
    },
    ChangePassword(current, next, confirm) {
      return request("POST", "/api/auth/change-password", { current, next, confirm });
    },
    BootstrapState() {
      return request("GET", "/api/app/bootstrap-state");
    },
    ListQuotes() {
      return request("GET", "/api/app/list-quotes");
    },
    AddQuote(content) {
      return request("POST", "/api/app/add-quote", { content });
    },
    SaveRecallAsQuote(question, response, keywords) {
      return request("POST", "/api/app/save-recall-as-quote", { question, response, keywords });
    },
    RefineQuoteDraft(content) {
      return request("POST", "/api/app/refine-quote-draft", { content });
    },
    UpdateQuote(id, content) {
      return request("POST", "/api/app/update-quote", { id, content });
    },
    DeleteQuotes(ids) {
      return request("POST", "/api/app/delete-quotes", { ids });
    },
    PreviewQuoteExport(ids) {
      return request("POST", "/api/app/preview-quote-export", { ids });
    },
    ImportQuotesPayload(payload) {
      return request("POST", "/api/app/import-quotes-payload", { payload });
    },
    SaveUserProfile(name) {
      return request("POST", "/api/app/save-user-profile", { name });
    },
    SaveSettings(settings) {
      return request("POST", "/api/app/save-settings", settings);
    },
    FetchModels(settings) {
      return request("POST", "/api/app/fetch-models", settings);
    },
    RunRecall(question) {
      return request("POST", "/api/app/run-recall", { question });
    },
    ListRecallHistory() {
      return request("GET", "/api/app/list-recall-history");
    },
    GetRecallHistory(id) {
      const url = "/api/app/get-recall-history?id=" + encodeURIComponent(String(id));
      return request("GET", url);
    },
    DeleteRecallHistory(ids) {
      return request("POST", "/api/app/delete-recall-history", { ids });
    },
  };

  window.go = { backend: { App: App } };
})();
`

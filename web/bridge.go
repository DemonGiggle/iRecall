package main

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
    GetAPITokenStatus() {
      return request("GET", "/api/app/get-api-token-status");
    },
    CreateAPIToken() {
      return request("POST", "/api/app/create-api-token");
    },
    BootstrapState() {
      return request("GET", "/api/app/bootstrap-state");
    },
    CountQuotes() {
      return request("GET", "/api/app/count-quotes");
    },
    ListQuotesPage(limit, offset) {
      const query = new URLSearchParams();
      if (Number.isFinite(limit) && limit > 0) {
        query.set("limit", String(limit));
      }
      if (Number.isFinite(offset) && offset > 0) {
        query.set("offset", String(offset));
      }
      const suffix = query.size > 0 ? "?" + query.toString() : "";
      return request("GET", "/api/app/list-quotes" + suffix);
    },
    ListQuotes() {
      return App.ListQuotesPage(0, 0);
    },
    AddQuote(content) {
      return request("POST", "/api/app/add-quote", { content });
    },
    AddQuoteWithImages(content, images) {
      return request("POST", "/api/app/add-quote-with-images", { content, images });
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
    UpdateQuoteWithImages(id, content, retainedIds, images) {
      return request("POST", "/api/app/update-quote-with-images", { id, content, retainedIds, images });
    },
    GetQuoteAttachmentData(id) {
      return request("GET", "/api/app/get-quote-attachment?id=" + encodeURIComponent(id));
    },
    ExportQuoteBundle(ids) {
      return request("POST", "/api/app/export-quote-bundle", { ids });
    },
    ImportQuoteBundle(payloadBase64) {
      return request("POST", "/api/app/import-quote-bundle", { payloadBase64 });
    },
    RegenerateQuoteKeywords(id, globalId) {
      return request("POST", "/api/app/regenerate-quote-keywords", { id, globalId });
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
    CountRecallHistory() {
      return request("GET", "/api/app/count-recall-history");
    },
    ListRecallHistoryPage(limit, offset) {
      const query = new URLSearchParams();
      if (Number.isFinite(limit) && limit > 0) {
        query.set("limit", String(limit));
      }
      if (Number.isFinite(offset) && offset > 0) {
        query.set("offset", String(offset));
      }
      const suffix = query.size > 0 ? "?" + query.toString() : "";
      return request("GET", "/api/app/list-recall-history" + suffix);
    },
    ListRecallHistory() {
      return App.ListRecallHistoryPage(0, 0);
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

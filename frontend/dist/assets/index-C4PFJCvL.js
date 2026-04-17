(function(){const s=document.createElement("link").relList;if(s&&s.supports&&s.supports("modulepreload"))return;for(const i of document.querySelectorAll('link[rel="modulepreload"]'))r(i);new MutationObserver(i=>{for(const n of i)if(n.type==="childList")for(const c of n.addedNodes)c.tagName==="LINK"&&c.rel==="modulepreload"&&r(c)}).observe(document,{childList:!0,subtree:!0});function a(i){const n={};return i.integrity&&(n.integrity=i.integrity),i.referrerPolicy&&(n.referrerPolicy=i.referrerPolicy),i.crossOrigin==="use-credentials"?n.credentials="include":i.crossOrigin==="anonymous"?n.credentials="omit":n.credentials="same-origin",n}function r(i){if(i.ep)return;i.ep=!0;const n=a(i);fetch(i.href,n)}})();const $={violet:{bg:"#0f172a",bgStrong:"#0b1120",panel:"rgba(17, 24, 39, 0.92)",panel2:"rgba(31, 41, 55, 0.82)",border:"#374151",borderStrong:"rgba(167, 139, 250, 0.42)",primary:"#7c3aed",accent:"#a78bfa",muted:"#94a3b8",fg:"#f9fafb",ok:"#10b981",error:"#ef4444",shadow:"0 24px 80px rgba(2, 6, 23, 0.38)",colorScheme:"dark"},forest:{bg:"#071a17",bgStrong:"#041311",panel:"rgba(9, 24, 21, 0.92)",panel2:"rgba(15, 41, 35, 0.82)",border:"#29443f",borderStrong:"rgba(45, 212, 191, 0.42)",primary:"#0f766e",accent:"#2dd4bf",muted:"#9ca3af",fg:"#ecfdf5",ok:"#22c55e",error:"#ef4444",shadow:"0 24px 80px rgba(1, 10, 9, 0.38)",colorScheme:"dark"},sunset:{bg:"#1c0f0a",bgStrong:"#130905",panel:"rgba(33, 17, 12, 0.94)",panel2:"rgba(52, 28, 18, 0.82)",border:"#5c4033",borderStrong:"rgba(251, 146, 60, 0.44)",primary:"#c2410c",accent:"#fb923c",muted:"#d6b8a6",fg:"#fffbeb",ok:"#16a34a",error:"#dc2626",shadow:"0 24px 80px rgba(20, 8, 2, 0.4)",colorScheme:"dark"},ocean:{bg:"#081824",bgStrong:"#06111a",panel:"rgba(11, 25, 38, 0.92)",panel2:"rgba(18, 40, 56, 0.82)",border:"#334155",borderStrong:"rgba(56, 189, 248, 0.42)",primary:"#0369a1",accent:"#38bdf8",muted:"#94a3b8",fg:"#f8fafc",ok:"#10b981",error:"#ef4444",shadow:"0 24px 80px rgba(3, 9, 16, 0.38)",colorScheme:"dark"},paper:{bg:"#f8fafc",bgStrong:"#e2e8f0",panel:"rgba(255, 255, 255, 0.96)",panel2:"rgba(248, 250, 252, 0.94)",border:"#cbd5e1",borderStrong:"rgba(29, 78, 216, 0.28)",primary:"#1d4ed8",accent:"#0f766e",muted:"#64748b",fg:"#111827",ok:"#15803d",error:"#b91c1c",shadow:"0 24px 80px rgba(148, 163, 184, 0.3)",colorScheme:"light"}};function E(e){const s=$[e in $?e:"violet"],a=document.documentElement;a.style.setProperty("--bg",s.bg),a.style.setProperty("--bg-strong",s.bgStrong),a.style.setProperty("--panel",s.panel),a.style.setProperty("--panel-2",s.panel2),a.style.setProperty("--border",s.border),a.style.setProperty("--border-strong",s.borderStrong),a.style.setProperty("--primary",s.primary),a.style.setProperty("--accent",s.accent),a.style.setProperty("--muted",s.muted),a.style.setProperty("--fg",s.fg),a.style.setProperty("--ok",s.ok),a.style.setProperty("--error",s.error),a.style.setProperty("--shadow",s.shadow),a.style.setProperty("color-scheme",s.colorScheme),document.body.dataset.theme=e}function Z(){return Object.keys($)}const t={bootstrapped:!1,fatalError:"",authChecked:!1,auth:null,authBusy:!1,authPassword:"",authConfirmPassword:"",authStatus:"",authIsError:!1,page:"Recall",bootstrap:null,quotes:[],quotesLoading:!1,quotesError:"",quotesCursor:0,quotesSelected:new Set,recallQuestion:"",recallLastQuestion:"",recallKeywords:[],recallQuotes:[],recallResponse:"",recallBusy:!1,recallError:"",recallStatus:"",recallStatusIsError:!1,recallCursor:0,recallSelected:new Set,historyEntries:[],historyLoading:!1,historyError:"",historyCursor:0,historySelected:new Set,historyDetail:null,historyDetailLoading:!1,historyDetailError:"",historyStatus:"",historyStatusIsError:!1,historyQuoteCursor:0,historyQuoteSelected:new Set,settings:zt(),settingsBusy:!1,settingsStatus:"",settingsIsError:!1,passwordForm:{current:"",next:"",confirm:"",busy:!1,status:"",isError:!1},overlay:null};let f=null,R=!1,k=null;const _=["Recall","History","Quotes","Settings"];function tt(e){f=e,R||(st(e),R=!0),o(),k||(k=et())}async function et(){try{if(t.auth=await u().AuthStatus(),t.authChecked=!0,t.auth.runtime==="web"&&!t.auth.authenticated){o();return}await F()}catch(e){t.authChecked=!0,t.bootstrapped=!0,t.fatalError=d(e),o()}}async function F(){var s;const e=await u().BootstrapState();t.bootstrap=e,t.bootstrapped=!0,t.page="Recall",t.settings=Vt(e),E(t.settings.theme),(s=e.profile)!=null&&s.DisplayName||(t.overlay={type:"namePrompt",name:"",busy:!1,status:"",isError:!1}),o(),await b()}function st(e){e.addEventListener("click",s=>{it(s)}),e.addEventListener("input",nt),e.addEventListener("change",lt),e.addEventListener("submit",s=>{ut(s)}),window.addEventListener("keydown",s=>{ct(s)})}async function T(){if(!(!t.auth||t.authBusy)){if(!t.authPassword.trim()){t.authStatus="Password is required.",t.authIsError=!0,o();return}t.authBusy=!0,t.authStatus="",t.authIsError=!1,o();try{await u().Login(t.authPassword),t.authPassword="",t.authConfirmPassword="",t.auth=await u().AuthStatus(),await F()}catch(e){t.authStatus=d(e),t.authIsError=!0,o()}finally{t.authBusy=!1}}}async function at(){await u().Logout(),t.auth=await u().AuthStatus(),t.bootstrapped=!1,t.bootstrap=null,t.overlay=null,t.quotes=[],t.historyEntries=[],t.historyDetail=null,t.authPassword="",t.authConfirmPassword="",t.authStatus="",t.authIsError=!1,o()}async function rt(){if(!t.passwordForm.busy){t.passwordForm.busy=!0,t.passwordForm.status="",o();try{await u().ChangePassword(t.passwordForm.current,t.passwordForm.next,t.passwordForm.confirm),t.passwordForm={current:"",next:"",confirm:"",busy:!1,status:"Password updated.",isError:!1}}catch(e){t.passwordForm.busy=!1,t.passwordForm.status=d(e),t.passwordForm.isError=!0}o()}}async function ot(e){var r,i;const s=(r=e.files)==null?void 0:r[0];if(!s||((i=t.overlay)==null?void 0:i.type)!=="importQuotes")return;const a=await s.text();t.overlay.filename=s.name,t.overlay.payload=a,t.overlay.path=s.name,t.overlay.status=`Loaded ${s.name}`,t.overlay.isError=!1,o()}async function it(e){const s=e.target;if(!(s instanceof HTMLElement))return;const a=s.closest("[data-action]");if(!a)return;switch(a.dataset.action??""){case"auth-login":await T();return;case"auth-logout":await at();return;case"nav":await dt(a.dataset.page);return;case"quotes-refresh":await b();return;case"history-refresh":await q();return;case"history-view-current":await pt();return;case"history-back":A();return;case"recall-save-quote":await Qt();return;case"history-save-quote":await Dt();return;case"history-delete-current":bt();return;case"history-select-all":Lt();return;case"history-deselect-all":Ht();return;case"quote-add":N("add");return;case"quote-import":ft();return;case"quote-edit-current":yt(a.dataset.context);return;case"quote-delete-current":vt(a.dataset.context);return;case"quote-share-current":await ht(a.dataset.context);return;case"set-cursor":if(s.closest("input, button, label"))return;Pt(a.dataset.context,Number(a.dataset.index??"0"));return;case"history-set-cursor":if(s.closest("input, button, label"))return;Rt(Number(a.dataset.index??"0"));return;case"history-open":await M(Number(a.dataset.id??"0"));return;case"profile-save":await B();return;case"quote-editor-save":await O();return;case"quote-editor-refine":await j();return;case"quote-editor-apply-refined":mt();return;case"quote-editor-reject-refined":gt();return;case"overlay-close":U();return;case"delete-confirm":await wt();return;case"share-browse":await $t();return;case"share-save":await St();return;case"import-browse":await Et();return;case"import-run":await qt();return;case"settings-fetch-models":await It();return;case"settings-save":await K();return;case"settings-change-password":await rt();return;case"recall-run":await Q();return;default:return}}function nt(e){var r,i,n,c;const s=e.target;if(!(s instanceof HTMLInputElement||s instanceof HTMLTextAreaElement))return;switch(s.dataset.bind??""){case"auth-password":t.authPassword=s.value;return;case"auth-confirm-password":t.authConfirmPassword=s.value;return;case"recall-question":t.recallQuestion=s.value;return;case"profile-name":((r=t.overlay)==null?void 0:r.type)==="namePrompt"&&(t.overlay.name=s.value);return;case"quote-editor-content":((i=t.overlay)==null?void 0:i.type)==="quoteEditor"&&(t.overlay.content=s.value);return;case"share-path":((n=t.overlay)==null?void 0:n.type)==="shareQuotes"&&(t.overlay.path=s.value);return;case"import-path":((c=t.overlay)==null?void 0:c.type)==="importQuotes"&&(t.overlay.path=s.value);return;case"settings-host":t.settings.host=s.value;return;case"settings-port":t.settings.port=s.value;return;case"settings-api-key":t.settings.apiKey=s.value;return;case"settings-model-filter":t.settings.modelFilter=s.value,C(t.settings),o();return;case"settings-max-results":t.settings.maxResults=s.value;return;case"settings-min-relevance":t.settings.minRelevance=s.value;return;case"settings-theme":t.settings.theme=s.value,E(t.settings.theme);return;case"settings-web-port":t.settings.webPort=s.value;return;case"settings-password-current":t.passwordForm.current=s.value;return;case"settings-password-next":t.passwordForm.next=s.value;return;case"settings-password-confirm":t.passwordForm.confirm=s.value;return;default:return}}function lt(e){const s=e.target;if(!(s instanceof HTMLInputElement||s instanceof HTMLSelectElement))return;switch(s.dataset.bind??""){case"quote-selected":Ct(s.dataset.context,Number(s.dataset.id??"0"),s.checked);return;case"history-selected":kt(Number(s.dataset.id??"0"),s.checked);return;case"settings-https":s instanceof HTMLInputElement&&(t.settings.https=s.checked);return;case"settings-model":t.settings.model=s.value;return;case"import-file":s instanceof HTMLInputElement&&ot(s);return;default:return}}async function ut(e){const s=e.target;if(s instanceof HTMLFormElement)switch(e.preventDefault(),s.dataset.form){case"auth-login":await T();return;case"auth-setup":await submitAuthSetup();return;case"recall":await Q();return;case"profile":await B();return;default:return}}async function ct(e){var a,r;const s=document.activeElement;if(e.key==="Escape"&&t.overlay&&t.overlay.type!=="namePrompt"){e.preventDefault(),U();return}if(e.ctrlKey&&e.key.toLowerCase()==="s"){if(((a=t.overlay)==null?void 0:a.type)==="quoteEditor"){e.preventDefault(),await O();return}!t.overlay&&t.page==="Settings"&&(e.preventDefault(),await K());return}if(e.ctrlKey&&e.key.toLowerCase()==="r"&&((r=t.overlay)==null?void 0:r.type)==="quoteEditor"){e.preventDefault(),await j();return}e.key==="Enter"&&!e.shiftKey&&s instanceof HTMLInputElement&&s.dataset.bind==="recall-question"&&(e.preventDefault(),await Q())}async function dt(e){t.page=e,o(),e==="Quotes"&&await b(),e==="History"&&await q()}async function b(){t.quotesLoading=!0,t.quotesError="",o();try{const e=await u().ListQuotes();t.quotes=e,t.quotesCursor=v(t.quotesCursor,e),t.quotesSelected=z(t.quotesSelected,e),t.quotesError=""}catch(e){t.quotesError=d(e)}finally{t.quotesLoading=!1,o()}}async function q(){t.historyLoading=!0,t.historyError="",t.historyStatus="",t.historyStatusIsError=!1,o();try{const e=await u().ListRecallHistory();t.historyEntries=e,t.historyCursor=x(t.historyCursor,e),t.historySelected=Yt(t.historySelected,e)}catch(e){t.historyError=d(e)}finally{t.historyLoading=!1,o()}}async function pt(){const e=I()[0];e&&await M(e.ID)}async function M(e){if(!(!Number.isFinite(e)||e<=0)){t.historyDetailLoading=!0,t.historyDetailError="",t.historyStatus="",t.historyStatusIsError=!1,o();try{const s=await u().GetRecallHistory(e);t.historyDetail=s,t.historyQuoteCursor=v(t.historyQuoteCursor,s.Quotes),t.historyQuoteSelected=z(t.historyQuoteSelected,s.Quotes)}catch(s){t.historyDetailError=d(s)}finally{t.historyDetailLoading=!1,o()}}}function A(){t.historyDetail=null,t.historyDetailLoading=!1,t.historyDetailError="",t.historyQuoteCursor=0,t.historyQuoteSelected=new Set,o()}function N(e,s){t.overlay={type:"quoteEditor",mode:e,quoteId:(s==null?void 0:s.ID)??null,content:(s==null?void 0:s.Content)??"",busy:!1,status:"",isError:!1,previewOriginal:"",previewRefined:""},o()}function yt(e){const s=h(e)[0];s&&N("edit",s)}function vt(e){const s=h(e).map(a=>a.ID);s.length!==0&&(t.overlay={type:"deleteQuotes",context:e,ids:s,busy:!1,status:"",isError:!1},o())}function bt(){const e=I().map(s=>s.ID);e.length!==0&&(t.overlay={type:"deleteHistory",ids:e,busy:!1,status:"",isError:!1},o())}async function ht(e){var a,r;const s=h(e);if(s.length!==0){t.overlay={type:"shareQuotes",context:e,ids:s.map(i=>i.ID),path:"",payload:"",busy:!0,status:"",isError:!1},o();try{const i=await u().PreviewQuoteExport(s.map(n=>n.ID));if(((a=t.overlay)==null?void 0:a.type)!=="shareQuotes")return;t.overlay.payload=i,t.overlay.busy=!1,t.overlay.status="Share payload ready. Save it to a file and transfer it manually.",t.overlay.isError=!1}catch(i){if(((r=t.overlay)==null?void 0:r.type)!=="shareQuotes")return;t.overlay.busy=!1,t.overlay.status=d(i),t.overlay.isError=!0}o()}}function ft(){t.overlay={type:"importQuotes",path:"",payload:"",filename:"",busy:!1,status:"",isError:!1,result:null},o()}async function B(){var s,a;if(((s=t.overlay)==null?void 0:s.type)!=="namePrompt"||t.overlay.busy)return;const e=t.overlay.name.trim();if(!e){t.overlay.status="Please enter a name to continue.",t.overlay.isError=!0,o();return}t.overlay.busy=!0,t.overlay.status="Saving profile…",t.overlay.isError=!1,o();try{const r=await u().SaveUserProfile(e);t.bootstrap&&(t.bootstrap.profile=r,t.bootstrap.greeting=`Hi! ${r.DisplayName}`),t.overlay=null}catch(r){((a=t.overlay)==null?void 0:a.type)==="namePrompt"&&(t.overlay.busy=!1,t.overlay.status=d(r),t.overlay.isError=!0)}o()}async function O(){var s,a;if(((s=t.overlay)==null?void 0:s.type)!=="quoteEditor"||t.overlay.busy)return;const e=t.overlay.content.trim();if(!e){t.overlay.status="Nothing to save.",t.overlay.isError=!0,o();return}t.overlay.busy=!0,t.overlay.status="Refining draft...",t.overlay.isError=!1,o();try{const r=t.overlay.mode==="add"?await u().AddQuote(e):await u().UpdateQuote(t.overlay.quoteId??0,e);t.overlay=null,D(r),await b()}catch(r){((a=t.overlay)==null?void 0:a.type)==="quoteEditor"&&(t.overlay.busy=!1,t.overlay.status=d(r),t.overlay.isError=!0),o()}}async function j(){var s,a,r;if(((s=t.overlay)==null?void 0:s.type)!=="quoteEditor"||t.overlay.busy)return;const e=t.overlay.content.trim();if(!e){t.overlay.status="Nothing to refine.",t.overlay.isError=!0,o();return}t.overlay.busy=!0,t.overlay.status="",o();try{const i=await u().RefineQuoteDraft(e);if(((a=t.overlay)==null?void 0:a.type)!=="quoteEditor")return;t.overlay.busy=!1,t.overlay.previewOriginal=e,t.overlay.previewRefined=i,t.overlay.status="",t.overlay.isError=!1}catch(i){((r=t.overlay)==null?void 0:r.type)==="quoteEditor"&&(t.overlay.busy=!1,t.overlay.status=d(i),t.overlay.isError=!0)}o()}function mt(){var e;((e=t.overlay)==null?void 0:e.type)==="quoteEditor"&&(t.overlay.content=t.overlay.previewRefined,t.overlay.previewOriginal="",t.overlay.previewRefined="",t.overlay.status="Refined draft applied. Review it, then save.",t.overlay.isError=!1,o())}function gt(){var e;((e=t.overlay)==null?void 0:e.type)==="quoteEditor"&&(t.overlay.previewOriginal="",t.overlay.previewRefined="",t.overlay.status="Refined draft discarded.",t.overlay.isError=!1,o())}async function wt(){var e,s,a,r;if(((e=t.overlay)==null?void 0:e.type)==="deleteHistory"){if(t.overlay.busy)return;t.overlay.busy=!0,t.overlay.status="",o();try{await u().DeleteRecallHistory(t.overlay.ids),Ft(t.overlay.ids),t.overlay=null,await q()}catch(i){((s=t.overlay)==null?void 0:s.type)==="deleteHistory"&&(t.overlay.busy=!1,t.overlay.status=d(i),t.overlay.isError=!0),o()}return}if(!(((a=t.overlay)==null?void 0:a.type)!=="deleteQuotes"||t.overlay.busy)){t.overlay.busy=!0,t.overlay.status="",o();try{await u().DeleteQuotes(t.overlay.ids),xt(t.overlay.ids),t.overlay=null,await b()}catch(i){((r=t.overlay)==null?void 0:r.type)==="deleteQuotes"&&(t.overlay.busy=!1,t.overlay.status=d(i),t.overlay.isError=!0),o()}}}async function $t(){var e,s,a;if(!(((e=t.overlay)==null?void 0:e.type)!=="shareQuotes"||t.overlay.busy)){if(y()){t.overlay.path="irecall-share.json",o();return}try{const r=await u().SelectQuoteExportFile();r&&((s=t.overlay)==null?void 0:s.type)==="shareQuotes"&&(t.overlay.path=r,o())}catch(r){((a=t.overlay)==null?void 0:a.type)==="shareQuotes"&&(t.overlay.status=d(r),t.overlay.isError=!0,o())}}}async function St(){var s,a,r;if(((s=t.overlay)==null?void 0:s.type)!=="shareQuotes"||t.overlay.busy)return;if(y()){const i=t.overlay.path.trim()||"irecall-share.json";if(!t.overlay.payload.trim()){t.overlay.status="Export payload is not ready yet.",t.overlay.isError=!0,o();return}Zt(i,t.overlay.payload),t.overlay.status=`Downloaded ${i}`,t.overlay.isError=!1,o();return}const e=t.overlay.path.trim();if(!e){t.overlay.status="Choose a file path for the export.",t.overlay.isError=!0,o();return}if(!t.overlay.payload.trim()){t.overlay.status="Export payload is not ready yet.",t.overlay.isError=!0,o();return}t.overlay.busy=!0,t.overlay.status="",o();try{await u().ExportQuotesToFile(t.overlay.ids,e),((a=t.overlay)==null?void 0:a.type)==="shareQuotes"&&(t.overlay.busy=!1,t.overlay.status=`Saved share payload to ${e}`,t.overlay.isError=!1,o())}catch(i){((r=t.overlay)==null?void 0:r.type)==="shareQuotes"&&(t.overlay.busy=!1,t.overlay.status=d(i),t.overlay.isError=!0,o())}}async function Et(){var e,s,a;if(!(((e=t.overlay)==null?void 0:e.type)!=="importQuotes"||t.overlay.busy)){if(y()){const r=document.querySelector('[data-bind="import-file"]');r==null||r.click();return}try{const r=await u().SelectQuoteImportFile();r&&((s=t.overlay)==null?void 0:s.type)==="importQuotes"&&(t.overlay.path=r,o())}catch(r){((a=t.overlay)==null?void 0:a.type)==="importQuotes"&&(t.overlay.status=d(r),t.overlay.isError=!0,o())}}}async function qt(){var a,r,i;if(((a=t.overlay)==null?void 0:a.type)!=="importQuotes"||t.overlay.busy)return;const e=t.overlay.path.trim(),s=t.overlay.payload.trim();if(y()){if(!s){t.overlay.status="Choose a file to import.",t.overlay.isError=!0,o();return}}else if(!e){t.overlay.status="Choose a file to import.",t.overlay.isError=!0,o();return}t.overlay.busy=!0,t.overlay.status="",t.overlay.result=null,o();try{const n=y()?await u().ImportQuotesPayload(s):await u().ImportQuotesFromFile(e);if(((r=t.overlay)==null?void 0:r.type)!=="importQuotes")return;t.overlay.busy=!1,t.overlay.result=n,t.overlay.status=`Imported quotes. inserted=${n.Inserted} updated=${n.Updated} duplicates=${n.Duplicates} stale=${n.Stale}`,t.overlay.isError=!1,await b()}catch(n){((i=t.overlay)==null?void 0:i.type)==="importQuotes"&&(t.overlay.busy=!1,t.overlay.status=d(n),t.overlay.isError=!0,o())}}async function Q(){if(t.recallBusy)return;const e=t.recallQuestion.trim();if(!e){t.recallError="Ask a question first.",o();return}t.recallBusy=!0,t.recallError="",t.recallStatus="",t.recallStatusIsError=!1,t.recallLastQuestion=e,t.recallKeywords=[],t.recallQuotes=[],t.recallResponse="",t.recallCursor=0,t.recallSelected=new Set,o();try{const s=await u().RunRecall(e);t.recallKeywords=s.keywords,t.recallQuotes=s.quotes,t.recallResponse=s.response,t.recallLastQuestion=s.question||e,t.recallCursor=0,t.recallSelected=new Set,t.recallQuestion=""}catch(s){t.recallError=d(s)}finally{t.recallBusy=!1,o()}}async function Qt(){const e=t.recallLastQuestion.trim(),s=t.recallResponse.trim();if(!e||!s){t.recallStatus="Run a recall first before saving it as a quote.",t.recallStatusIsError=!0,o();return}try{const a=await u().SaveRecallAsQuote(e,s,t.recallKeywords);D(a),await b(),t.recallStatus="Saved recall as quote.",t.recallStatusIsError=!1,t.overlay={type:"notice",title:"Recall Saved as Quote",message:"The current question and grounded response were saved as a quote with generated tags."}}catch(a){t.recallStatus=d(a),t.recallStatusIsError=!0}o()}async function Dt(){const e=t.historyDetail;if(e){try{const s=await u().SaveRecallAsQuote(e.Question,e.Response,[]);D(s),await b(),t.historyStatus="Saved history entry as quote.",t.historyStatusIsError=!1,t.overlay={type:"notice",title:"History Entry Saved as Quote",message:"The selected history question and response were saved as a quote with generated tags."}}catch(s){t.historyStatus=d(s),t.historyStatusIsError=!0}o()}}async function It(){if(t.settingsBusy)return;let e;try{e=J(t.settings)}catch(s){t.settingsStatus=d(s),t.settingsIsError=!0,o();return}t.settingsBusy=!0,t.settingsStatus="",o();try{const s=await u().FetchModels(e);t.settings.models=s,C(t.settings),t.settingsStatus=s.length>0?`Fetched ${s.length} models.`:"No models returned.",t.settingsIsError=!1}catch(s){t.settingsStatus=d(s),t.settingsIsError=!0}finally{t.settingsBusy=!1,o()}}async function K(){var s;if(t.settingsBusy)return;let e;try{e=Gt(t.settings)}catch(a){t.settingsStatus=d(a),t.settingsIsError=!0,o();return}t.settingsBusy=!0,t.settingsStatus="",o();try{const a=await u().SaveSettings(e);t.settings=W(a,t.settings.models),E(t.settings.theme),t.bootstrap&&(t.bootstrap.settings=a);const r=((s=t.auth)==null?void 0:s.runtime)==="web"&&t.auth.currentPort>0&&t.auth.currentPort!==a.Web.Port;t.settingsStatus=r?"Saved. Restart the web server to apply the new port.":"Saved.",t.settingsIsError=!1}catch(a){t.settingsStatus=d(a),t.settingsIsError=!0}finally{t.settingsBusy=!1,o()}}function U(){t.overlay&&t.overlay.type!=="namePrompt"&&("busy"in t.overlay&&t.overlay.busy||(t.overlay=null,o()))}function Pt(e,s){var a;if(e==="quotes")t.quotesCursor=v(s,t.quotes);else if(e==="recall")t.recallCursor=v(s,t.recallQuotes);else{const r=((a=t.historyDetail)==null?void 0:a.Quotes)??[];t.historyQuoteCursor=v(s,r)}o()}function Ct(e,s,a){const r=e==="quotes"?t.quotesSelected:e==="recall"?t.recallSelected:t.historyQuoteSelected;a?r.add(s):r.delete(s)}function h(e){var n;const s=e==="quotes"?t.quotes:e==="recall"?t.recallQuotes:((n=t.historyDetail)==null?void 0:n.Quotes)??[],a=e==="quotes"?t.quotesCursor:e==="recall"?t.recallCursor:t.historyQuoteCursor,r=e==="quotes"?t.quotesSelected:e==="recall"?t.recallSelected:t.historyQuoteSelected,i=s.filter(c=>r.has(c.ID));return i.length>0?i:s[a]?[s[a]]:[]}function D(e){t.quotes=w(t.quotes,e),t.recallQuotes=w(t.recallQuotes,e),t.historyDetail&&(t.historyDetail={...t.historyDetail,Quotes:w(t.historyDetail.Quotes,e)}),o()}function xt(e){var a;const s=new Set(e);t.quotes=t.quotes.filter(r=>!s.has(r.ID)),t.recallQuotes=t.recallQuotes.filter(r=>!s.has(r.ID)),t.historyDetail&&(t.historyDetail={...t.historyDetail,Quotes:t.historyDetail.Quotes.filter(r=>!s.has(r.ID))}),t.quotesSelected=new Set([...t.quotesSelected].filter(r=>!s.has(r))),t.recallSelected=new Set([...t.recallSelected].filter(r=>!s.has(r))),t.historyQuoteSelected=new Set([...t.historyQuoteSelected].filter(r=>!s.has(r))),t.quotesCursor=v(t.quotesCursor,t.quotes),t.recallCursor=v(t.recallCursor,t.recallQuotes),t.historyQuoteCursor=v(t.historyQuoteCursor,((a=t.historyDetail)==null?void 0:a.Quotes)??[]),o()}function Rt(e){t.historyCursor=x(e,t.historyEntries),o()}function kt(e,s){s?t.historySelected.add(e):t.historySelected.delete(e)}function Lt(){t.historySelected=new Set(t.historyEntries.map(e=>e.ID)),o()}function Ht(){t.historySelected=new Set,o()}function I(){const e=t.historyEntries.filter(s=>t.historySelected.has(s.ID));return e.length>0?e:t.historyEntries[t.historyCursor]?[t.historyEntries[t.historyCursor]]:[]}function Ft(e){const s=new Set(e);if(t.historyEntries=t.historyEntries.filter(a=>!s.has(a.ID)),t.historySelected=new Set([...t.historySelected].filter(a=>!s.has(a))),t.historyCursor=x(t.historyCursor,t.historyEntries),t.historyDetail&&s.has(t.historyDetail.ID)){A();return}o()}function o(){if(!f)return;const e=Tt();f.innerHTML=At(),Mt(e)}function Tt(){const e=document.activeElement;if(!(e instanceof HTMLInputElement||e instanceof HTMLTextAreaElement||e instanceof HTMLSelectElement))return null;const s=e.dataset.bind;return s?{selector:`[data-bind="${s}"]`,selectionStart:e instanceof HTMLInputElement||e instanceof HTMLTextAreaElement?e.selectionStart:null,selectionEnd:e instanceof HTMLInputElement||e instanceof HTMLTextAreaElement?e.selectionEnd:null}:null}function Mt(e){if(!f||!e)return;const s=f.querySelector(e.selector);(s instanceof HTMLInputElement||s instanceof HTMLTextAreaElement||s instanceof HTMLSelectElement)&&(s.focus({preventScroll:!0}),(s instanceof HTMLInputElement||s instanceof HTMLTextAreaElement)&&e.selectionStart!==null&&e.selectionEnd!==null&&s.setSelectionRange(e.selectionStart,e.selectionEnd))}function At(){var s,a,r,i,n;if(!t.authChecked)return`
      <div class="shell shell-loading">
        <div class="splash">
          <div class="brand">iRecall</div>
          <div class="muted">Checking workspace access…</div>
        </div>
      </div>
    `;if(t.fatalError)return`
      <div class="shell shell-loading">
        <div class="splash splash-error">
          <div class="brand">iRecall</div>
          <div class="status status-error">${l(t.fatalError)}</div>
        </div>
      </div>
    `;if(((s=t.auth)==null?void 0:s.runtime)==="web"&&!t.auth.authenticated)return Nt();if(!t.bootstrapped)return`
      <div class="shell shell-loading">
        <div class="splash">
          <div class="brand">iRecall</div>
          <div class="muted">Loading workspace…</div>
        </div>
      </div>
    `;const e=(r=(a=t.bootstrap)==null?void 0:a.profile)!=null&&r.DisplayName?`Hi! ${t.bootstrap.profile.DisplayName}`:"";return`
    <div class="shell">
      <header class="titlebar">
        <div class="brand-lockup">
          <div class="brand">${l(((i=t.bootstrap)==null?void 0:i.productName)??"iRecall")}</div>
          <div class="muted subtle">${y()?"Local-first quote recall web UI":"Local-first quote recall desktop"}</div>
        </div>
        <div class="titlebar-right">
          <div class="greeting">${l(e)}</div>
          ${((n=t.auth)==null?void 0:n.runtime)==="web"?'<button class="button" data-action="auth-logout" type="button">Logout</button>':""}
          <nav class="tabs" aria-label="Primary">
            ${_.map(c=>`
                  <button
                    class="tab${t.page===c?" active":""}"
                    data-action="nav"
                    data-page="${c}"
                    type="button"
                  >${c}</button>
                `).join("")}
          </nav>
        </div>
      </header>

      <main class="layout">
        ${Bt()}
      </main>

      ${t.overlay?Wt(t.overlay):""}
    </div>
  `}function Nt(){var i;const e=!((i=t.auth)!=null&&i.passwordConfigured),s="auth-login";return`
    <div class="shell shell-loading">
      <div class="panel modal">
        <div class="brand">iRecall</div>
        <div class="modal-title">${e?"Password Required In Terminal":"Unlock Web UI"}</div>
        <div class="modal-copy">${e?"The web password must be created in the terminal before the server starts listening. Restart the server from a terminal session to finish setup.":"Enter the web password to unlock the shared iRecall database."}</div>
        <form class="modal-form" data-form="${s}">
          <label class="field">
            <span>Password</span>
            <input class="text-input" data-bind="auth-password" type="password" value="${p(t.authPassword)}" ${e?"disabled":""} />
          </label>
          ${t.authStatus?`<div class="status ${t.authIsError?"status-error":"status-ok"}">${l(t.authStatus)}</div>`:""}
          <div class="modal-actions">
            <button class="button button-primary" data-action="${s}" type="submit" ${t.authBusy||e?"disabled":""}>
              ${t.authBusy?"Working…":"Login"}
            </button>
          </div>
        </form>
      </div>
    </div>
  `}function Bt(){switch(t.page){case"Recall":return Ot();case"Quotes":return jt();case"History":return Kt();case"Settings":return Ut()}}function Ot(){const e=h("recall"),s=t.recallResponse.trim()?l(t.recallResponse):'<span class="muted">Grounded response will appear here.</span>',a=t.recallKeywords.length>0?t.recallKeywords.map(r=>`<span class="keyword-chip">${l(r)}</span>`).join(""):'<span class="muted">Keywords: —</span>';return`
    <section class="page page-recall">
      <div class="panel page-panel">
        <div class="section-heading">
          <div>
            <div class="section-title">Recall</div>
            <div class="muted">Ask a question, ground the answer in quotes, then manage the reference set.</div>
          </div>
          <div class="toolbar">
            <button class="button" data-action="quote-add" type="button">Add Quote</button>
            <button class="button" data-action="recall-save-quote" type="button" ${t.recallResponse.trim()?"":"disabled"}>Save as Quote</button>
            <button class="button" data-action="quote-edit-current" data-context="recall" type="button" ${e.length===0?"disabled":""}>Edit</button>
            <button class="button button-danger" data-action="quote-delete-current" data-context="recall" type="button" ${e.length===0?"disabled":""}>Delete</button>
            <button class="button" data-action="quote-share-current" data-context="recall" type="button" ${e.length===0?"disabled":""}>Share</button>
          </div>
        </div>

        <form class="question-bar" data-form="recall">
          <input
            class="text-input text-input-lg"
            data-bind="recall-question"
            placeholder="Ask anything..."
            value="${p(t.recallQuestion)}"
          />
          <button class="button button-primary" data-action="recall-run" type="submit" ${t.recallBusy?"disabled":""}>
            ${t.recallBusy?"Thinking…":"Ask"}
          </button>
        </form>

        <div class="keyword-row">
          <span class="muted">Keywords:</span>
          <div class="keyword-list">${a}</div>
        </div>

        <div class="recall-grid">
          <section class="panel subpanel">
            <div class="subpanel-header">
              <div class="section-title">Response</div>
              <div class="muted">${t.recallBusy?"Generating grounded answer…":"Uses the current reference quotes."}</div>
            </div>
            <pre class="response-box">${s}</pre>
          </section>

          <section class="panel subpanel">
            <div class="subpanel-header">
              <div class="section-title">Reference Quotes</div>
              <div class="muted">${e.length>0?`${e.length} selected`:`${t.recallQuotes.length} loaded`}</div>
            </div>
            ${P("recall",t.recallQuotes,t.recallCursor,t.recallSelected,!1)}
          </section>
        </div>

        ${t.recallError?`<div class="status status-error">${l(t.recallError)}</div>`:""}
        ${t.recallStatus?`<div class="status ${t.recallStatusIsError?"status-error":"status-ok"}">${l(t.recallStatus)}</div>`:""}
      </div>
    </section>
  `}function jt(){const e=h("quotes");let s="";return t.quotesLoading?s='<div class="empty-state">Loading quotes…</div>':t.quotesError?s=`<div class="status status-error">${l(t.quotesError)}</div>`:s=P("quotes",t.quotes,t.quotesCursor,t.quotesSelected,!0),`
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
            <button class="button" data-action="quote-edit-current" data-context="quotes" type="button" ${e.length===0?"disabled":""}>Edit</button>
            <button class="button button-danger" data-action="quote-delete-current" data-context="quotes" type="button" ${e.length===0?"disabled":""}>Delete</button>
            <button class="button" data-action="quote-share-current" data-context="quotes" type="button" ${e.length===0?"disabled":""}>Share</button>
          </div>
        </div>
        <div class="meta-row">
          <span class="muted">Stored Quotes:</span>
          <span>${t.quotes.length}</span>
          <span class="muted">Selection:</span>
          <span>${e.length>0?e.length:t.quotes.length>0?1:0}</span>
        </div>
        ${s}
      </div>
    </section>
  `}function Kt(){const e=I(),s=h("history");if(t.historyDetailLoading)return`
      <section class="page page-history">
        <div class="panel page-panel">
          <div class="empty-state">Loading history entry…</div>
        </div>
      </section>
    `;if(t.historyDetail)return`
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
              <button class="button" data-action="quote-edit-current" data-context="history" type="button" ${s.length===0?"disabled":""}>Edit Quote</button>
              <button class="button button-danger" data-action="quote-delete-current" data-context="history" type="button" ${s.length===0?"disabled":""}>Delete Quote</button>
              <button class="button" data-action="quote-share-current" data-context="history" type="button" ${s.length===0?"disabled":""}>Share Quote</button>
            </div>
          </div>

          ${t.historyDetailError?`<div class="status status-error">${l(t.historyDetailError)}</div>`:""}

          <div class="recall-grid">
            <section class="panel subpanel">
              <div class="subpanel-header">
                <div class="section-title">History Entry</div>
                <div class="muted">${l(H(t.historyDetail.CreatedAt))}</div>
              </div>
              <div class="detail-stack">
                <div class="detail-block">
                  <div class="muted">Question</div>
                  <pre class="response-box">${l(t.historyDetail.Question)}</pre>
                </div>
                <div class="detail-block">
                  <div class="muted">Response</div>
                  <pre class="response-box">${l(t.historyDetail.Response)}</pre>
                </div>
              </div>
            </section>

            <section class="panel subpanel">
              <div class="subpanel-header">
                <div class="section-title">Reference Quotes</div>
                <div class="muted">${s.length>0?`${s.length} selected`:`${t.historyDetail.Quotes.length} loaded`}</div>
              </div>
              ${P("history",t.historyDetail.Quotes,t.historyQuoteCursor,t.historyQuoteSelected,!1)}
            </section>
          </div>
          ${t.historyStatus?`<div class="status ${t.historyStatusIsError?"status-error":"status-ok"}">${l(t.historyStatus)}</div>`:""}
        </div>
      </section>
    `;let a="";return t.historyLoading?a='<div class="empty-state">Loading history…</div>':t.historyError?a=`<div class="status status-error">${l(t.historyError)}</div>`:t.historyEntries.length===0?a='<div class="empty-state">No recall history yet. Run a recall from the Recall tab to create one.</div>':a=`
      <div class="history-list">
        ${t.historyEntries.map((r,i)=>{const n=i===t.historyCursor,c=S(r.Response,140);return`
              <article class="quote-card${n?" is-current":""}" data-action="history-set-cursor" data-index="${i}">
                <div class="quote-topline">
                  <label class="selection-toggle">
                    <input
                      type="checkbox"
                      data-bind="history-selected"
                      data-id="${r.ID}"
                      ${t.historySelected.has(r.ID)?"checked":""}
                    />
                    <span>${t.historySelected.has(r.ID)?"[x]":"[ ]"}</span>
                  </label>
                  <div class="quote-topline-meta">
                    <span class="quote-index${n?" is-current":""}">${n?"&gt; ":""}[${i+1}]</span>
                    <span class="quote-version">${l(H(r.CreatedAt))}</span>
                  </div>
                </div>
                <div class="quote-content">${l(S(r.Question,120))}</div>
                <div class="quote-meta"><span class="muted">Response:</span> <span>${l(c||"(empty response)")}</span></div>
                <div class="toolbar toolbar-inline">
                  <button class="button" data-action="history-open" data-id="${r.ID}" type="button">View</button>
                </div>
              </article>
            `}).join("")}
      </div>
    `,`
    <section class="page page-history">
      <div class="panel page-panel">
        <div class="section-heading">
          <div>
            <div class="section-title">History</div>
            <div class="muted">Review past recall sessions, inspect grounded responses, and manage saved history entries.</div>
          </div>
          <div class="toolbar">
            <button class="button" data-action="history-refresh" type="button">Refresh</button>
            <button class="button" data-action="history-select-all" type="button" ${t.historyEntries.length===0?"disabled":""}>Select All</button>
            <button class="button" data-action="history-deselect-all" type="button" ${t.historySelected.size===0?"disabled":""}>Deselect All</button>
            <button class="button" data-action="history-view-current" type="button" ${e.length===0?"disabled":""}>View</button>
            <button class="button button-danger" data-action="history-delete-current" type="button" ${e.length===0?"disabled":""}>Delete</button>
          </div>
        </div>
        <div class="meta-row">
          <span class="muted">Stored History:</span>
          <span>${t.historyEntries.length}</span>
          <span class="muted">Selection:</span>
          <span>${e.length>0?e.length:t.historyEntries.length>0?1:0}</span>
        </div>
        ${a}
      </div>
    </section>
  `}function Ut(){var i,n;const e=V(t.settings),s=(i=t.bootstrap)==null?void 0:i.paths,a=(n=t.auth)==null?void 0:n.currentPort,r=t.settings.models.length>0&&e.length>0?`
        <select class="select-input" data-bind="settings-model">
          ${e.map(c=>`
                <option value="${p(c)}"${c===t.settings.model?" selected":""}>${l(c)}</option>
              `).join("")}
        </select>
      `:`
        <div class="readonly-model">
          <span>${l(t.settings.model||"(none)")}</span>
          <span class="muted">${t.settings.models.length===0?"Fetch models first":"No matches"}</span>
        </div>
      `;return`
    <section class="page page-settings">
      <div class="panel page-panel">
        <div class="section-heading">
          <div>
            <div class="section-title">Settings</div>
            <div class="muted">Configure the OpenAI-compatible endpoint and quote retrieval behavior.</div>
          </div>
          <div class="toolbar">
            <button class="button button-primary" data-action="settings-save" type="button" ${t.settingsBusy?"disabled":""}>Save</button>
          </div>
        </div>

        <div class="settings-grid">
          <section class="panel subpanel">
            <div class="section-title">LLM Provider</div>
            <label class="field">
              <span>Host / IP</span>
              <input class="text-input" data-bind="settings-host" value="${p(t.settings.host)}" />
            </label>
            <label class="field">
              <span>Port</span>
              <input class="text-input" data-bind="settings-port" value="${p(t.settings.port)}" />
            </label>
            <label class="field checkbox-field">
              <input type="checkbox" data-bind="settings-https"${t.settings.https?" checked":""} />
              <span>Use HTTPS</span>
            </label>
            <label class="field">
              <span>API Key</span>
              <input class="text-input" data-bind="settings-api-key" type="password" value="${p(t.settings.apiKey)}" />
            </label>
            <label class="field">
              <span>Filter</span>
              <input class="text-input" data-bind="settings-model-filter" value="${p(t.settings.modelFilter)}" placeholder="Type to filter models" />
            </label>
            <label class="field">
              <span>Model</span>
              <div class="field-inline">
                <div class="field-inline-grow">${r}</div>
                <button class="button" data-action="settings-fetch-models" type="button" ${t.settingsBusy?"disabled":""}>
                  ${t.settingsBusy?"Fetching…":"Fetch Models"}
                </button>
              </div>
            </label>
          </section>

          <section class="panel subpanel">
            <div class="section-title">Search</div>
            <label class="field">
              <span>Max ref quotes</span>
              <input class="text-input" data-bind="settings-max-results" value="${p(t.settings.maxResults)}" />
            </label>
            <label class="field">
              <span>Min relevance</span>
              <input class="text-input" data-bind="settings-min-relevance" value="${p(t.settings.minRelevance)}" placeholder="0.0-1.0" />
            </label>
            <div class="settings-hint muted">
              Saving updates both the persisted settings and the live engine configuration for the current session.
              0.0 keeps broad matches. Try 0.3-0.7 for cleaner results; 1.0 is very strict.
            </div>
          </section>

          <section class="panel subpanel">
            <div class="section-title">Appearance + Web UI</div>
            <label class="field">
              <span>Theme</span>
              <select class="select-input" data-bind="settings-theme">
                ${Z().map(c=>`
                      <option value="${c}"${c===t.settings.theme?" selected":""}>${c}</option>
                    `).join("")}
              </select>
            </label>
            <label class="field">
              <span>Web Port</span>
              <input class="text-input" data-bind="settings-web-port" value="${p(t.settings.webPort)}" />
            </label>
            <div class="settings-hint muted">
              The web server listens on this port after restart. Current listener: ${l(a?String(a):"(not running)")}.
            </div>
          </section>

          <section class="panel subpanel">
            <div class="section-title">Change Password</div>
            <label class="field">
              <span>Current Password</span>
              <input class="text-input" data-bind="settings-password-current" type="password" value="${p(t.passwordForm.current)}" />
            </label>
            <label class="field">
              <span>New Password</span>
              <input class="text-input" data-bind="settings-password-next" type="password" value="${p(t.passwordForm.next)}" />
            </label>
            <label class="field">
              <span>Confirm Password</span>
              <input class="text-input" data-bind="settings-password-confirm" type="password" value="${p(t.passwordForm.confirm)}" />
            </label>
            <div class="muted subtle">Use at least 12 characters and include at least 3 of: uppercase, lowercase, digit, symbol.</div>
            <div class="toolbar">
              <button class="button" data-action="settings-change-password" type="button" ${t.passwordForm.busy?"disabled":""}>
                ${t.passwordForm.busy?"Updating…":"Change Password"}
              </button>
            </div>
            ${t.passwordForm.status?`<div class="status ${t.passwordForm.isError?"status-error":"status-ok"}">${l(t.passwordForm.status)}</div>`:""}
          </section>

          <section class="panel subpanel">
            <div class="section-title">Local Storage</div>
            <div class="settings-paths">
              <div class="field">
                <span>Data dir</span>
                <div class="readonly-model path-value">${l((s==null?void 0:s.dataDir)??"(unavailable)")}</div>
              </div>
              <div class="field">
                <span>Config dir</span>
                <div class="readonly-model path-value">${l((s==null?void 0:s.configDir)??"(unavailable)")}</div>
              </div>
              <div class="field">
                <span>State dir</span>
                <div class="readonly-model path-value">${l((s==null?void 0:s.stateDir)??"(unavailable)")}</div>
              </div>
              <div class="field">
                <span>Database</span>
                <div class="readonly-model path-value">${l((s==null?void 0:s.dbPath)??"(unavailable)")}</div>
              </div>
            </div>
          </section>
        </div>

        ${t.settingsStatus?`<div class="status ${t.settingsIsError?"status-error":"status-ok"}">${l(t.settingsStatus)}</div>`:""}
      </div>
    </section>
  `}function P(e,s,a,r,i){return s.length===0?`<div class="empty-state">${e==="quotes"?"No quotes yet. Add one or import a shared payload.":"No reference quotes for this question yet."}</div>`:`
    <div class="quote-list">
      ${s.map((n,c)=>{const g=c===a,Y=!n.IsOwnedByMe&&n.SourceName?`<div class="quote-meta"><span class="muted">From:</span> <span class="meta-accent">${l(n.SourceName)}</span></div>`:"",X=i?`
              <div class="quote-meta">
                <span class="muted">Tags:</span>
                <span>${n.Tags.length>0?l(Xt(n.Tags,3)):"(none)"}</span>
              </div>
            `:"";return`
            <article class="quote-card${g?" is-current":""}" data-action="set-cursor" data-context="${e}" data-index="${c}">
              <div class="quote-topline">
                <label class="selection-toggle">
                  <input
                    type="checkbox"
                    data-bind="quote-selected"
                    data-context="${e}"
                    data-id="${n.ID}"
                    ${r.has(n.ID)?"checked":""}
                  />
                  <span>${r.has(n.ID)?"[x]":"[ ]"}</span>
                </label>
                <div class="quote-topline-meta">
                  <span class="quote-index${g?" is-current":""}">${g?"&gt; ":""}[${c+1}]</span>
                  <span class="quote-version">v${n.Version}</span>
                </div>
              </div>
              <div class="quote-content">${l(S(n.Content,e==="quotes"?96:120))}</div>
              ${Y}
              ${X}
            </article>
          `}).join("")}
    </div>
  `}function Wt(e){switch(e.type){case"namePrompt":return`
        <div class="overlay-backdrop">
          <div class="modal">
            <div class="modal-title">Set Your Name</div>
            <p class="modal-copy">
              Your name is attached to quotes you share and shown when other users receive your quotes.
            </p>
            <form class="modal-form" data-form="profile">
              <label class="field">
                <span>Display Name</span>
                <input class="text-input text-input-lg" data-bind="profile-name" value="${p(e.name)}" placeholder="Your name" />
              </label>
              ${e.status?`<div class="status ${e.isError?"status-error":"status-ok"}">${l(e.status)}</div>`:""}
              <div class="modal-actions">
                <button class="button button-primary" data-action="profile-save" type="submit" ${e.busy?"disabled":""}>
                  ${e.busy?"Saving…":"Save name and continue"}
                </button>
              </div>
            </form>
          </div>
        </div>
      `;case"quoteEditor":return`
        <div class="overlay-backdrop">
          <div class="modal modal-wide">
            <div class="modal-title">${e.mode==="add"?"Add Quote":"Edit Quote"}</div>
            ${e.previewRefined?`
                  <div class="compare-grid">
                    <section class="panel compare-panel">
                      <div class="section-title">Current Draft</div>
                      <pre class="compare-body">${l(e.previewOriginal)}</pre>
                    </section>
                    <section class="panel compare-panel">
                      <div class="section-title">Refined Draft</div>
                      <pre class="compare-body">${l(e.previewRefined)}</pre>
                    </section>
                  </div>
                `:`
                  <label class="field">
                    <span>Quote Content</span>
                    <textarea class="text-area" data-bind="quote-editor-content" rows="10" placeholder="Type or paste your note here.">${l(e.content)}</textarea>
                  </label>
                `}
            <div class="muted modal-copy">
              ${e.previewRefined?"Compare the current draft with the suggested rewrite before applying it.":"Tags are regenerated automatically by the shared core logic."}
            </div>
            ${e.status?`<div class="status ${e.isError?"status-error":"status-ok"}">${l(e.status)}</div>`:""}
            <div class="modal-actions">
              ${e.previewRefined?`
                    <button class="button button-primary" data-action="quote-editor-apply-refined" type="button">Apply refined draft</button>
                    <button class="button" data-action="quote-editor-reject-refined" type="button">Keep editing current draft</button>
                  `:`
                    <button class="button button-primary" data-action="quote-editor-save" type="button" ${e.busy?"disabled":""}>
                      ${e.busy?"Saving…":"Save"}
                    </button>
                    <button class="button" data-action="quote-editor-refine" type="button" ${e.busy?"disabled":""}>
                      ${e.busy?"Working…":"Refine"}
                    </button>
                    <button class="button" data-action="overlay-close" type="button" ${e.busy?"disabled":""}>Cancel</button>
                  `}
            </div>
          </div>
        </div>
      `;case"deleteQuotes":return`
        <div class="overlay-backdrop">
          <div class="modal">
            <div class="modal-title modal-title-danger">Delete Quotes</div>
            <div class="modal-copy">This permanently removes the selected quote entries from the local library.</div>
            <div class="summary-list">
              ${L(e.context,e.ids).map((s,a)=>`<div class="summary-item">[${a+1}] ${l(m(s.Content,140))}</div>`).join("")}
            </div>
            ${e.status?`<div class="status ${e.isError?"status-error":"status-ok"}">${l(e.status)}</div>`:""}
            <div class="modal-actions">
              <button class="button button-danger" data-action="delete-confirm" type="button" ${e.busy?"disabled":""}>
                ${e.busy?"Deleting…":"Delete"}
              </button>
              <button class="button" data-action="overlay-close" type="button" ${e.busy?"disabled":""}>Cancel</button>
            </div>
          </div>
        </div>
      `;case"deleteHistory":return`
        <div class="overlay-backdrop">
          <div class="modal">
            <div class="modal-title modal-title-danger">Delete History</div>
            <div class="modal-copy">This permanently removes the selected recall history entries from the local library.</div>
            <div class="summary-list">
              ${Jt(e.ids).map((s,a)=>`<div class="summary-item">[${a+1}] ${l(m(s.Question,140))}</div>`).join("")}
            </div>
            ${e.status?`<div class="status ${e.isError?"status-error":"status-ok"}">${l(e.status)}</div>`:""}
            <div class="modal-actions">
              <button class="button button-danger" data-action="delete-confirm" type="button" ${e.busy?"disabled":""}>
                ${e.busy?"Deleting…":"Delete"}
              </button>
              <button class="button" data-action="overlay-close" type="button" ${e.busy?"disabled":""}>Cancel</button>
            </div>
          </div>
        </div>
      `;case"shareQuotes":return`
        <div class="overlay-backdrop">
          <div class="modal modal-wide">
            <div class="modal-title">Share Quotes</div>
            <div class="summary-list">
              ${L(e.context,e.ids).map((s,a)=>`<div class="summary-item">[${a+1}] v${s.Version} ${l(m(s.Content,120))}</div>`).join("")}
            </div>
            <label class="field">
              <span>${y()?"Download As":"Save To"}</span>
              ${y()?`<input class="text-input" data-bind="share-path" value="${p(e.path||"irecall-share.json")}" placeholder="irecall-share.json" />`:`
                    <div class="path-row">
                      <input class="text-input" data-bind="share-path" value="${p(e.path)}" placeholder="/path/to/irecall-share.json" />
                      <button class="button" data-action="share-browse" type="button" ${e.busy?"disabled":""}>Browse</button>
                    </div>
                  `}
            </label>
            <div class="muted modal-copy">${y()?"Download the JSON payload locally, then transfer it manually to the recipient.":"Export to a JSON file and transfer it manually to the recipient."}</div>
            <div class="payload-box"><pre>${l(e.payload||"Preparing export payload…")}</pre></div>
            ${e.status?`<div class="status ${e.isError?"status-error":"status-ok"}">${l(e.status)}</div>`:""}
            <div class="modal-actions">
              <button class="button button-primary" data-action="share-save" type="button" ${e.busy?"disabled":""}>
                ${e.busy?"Working…":y()?"Download export file":"Save export file"}
              </button>
              <button class="button" data-action="overlay-close" type="button" ${e.busy?"disabled":""}>Close</button>
            </div>
          </div>
        </div>
      `;case"importQuotes":return`
        <div class="overlay-backdrop">
          <div class="modal">
            <div class="modal-title">Import Quotes</div>
            <div class="modal-copy">Import a quote share JSON file exported from another iRecall instance.</div>
            <label class="field">
              <span>Import From</span>
              ${y()?`
                    <input class="text-input" data-bind="import-path" value="${p(e.path)}" placeholder="Choose a local JSON file" readonly />
                    <input data-bind="import-file" type="file" accept="application/json,.json" hidden />
                    <div class="toolbar">
                      <button class="button" data-action="import-browse" type="button" ${e.busy?"disabled":""}>Choose File</button>
                    </div>
                  `:`
                    <div class="path-row">
                      <input class="text-input" data-bind="import-path" value="${p(e.path)}" placeholder="/path/to/irecall-share.json" />
                      <button class="button" data-action="import-browse" type="button" ${e.busy?"disabled":""}>Browse</button>
                    </div>
                  `}
            </label>
            ${e.result?`
                  <div class="result-grid">
                    <div><span class="muted">Inserted:</span> ${e.result.Inserted}</div>
                    <div><span class="muted">Updated:</span> ${e.result.Updated}</div>
                    <div><span class="muted">Duplicates:</span> ${e.result.Duplicates}</div>
                    <div><span class="muted">Stale:</span> ${e.result.Stale}</div>
                  </div>
                `:""}
            ${e.status?`<div class="status ${e.isError?"status-error":"status-ok"}">${l(e.status)}</div>`:""}
            <div class="modal-actions">
              <button class="button button-primary" data-action="import-run" type="button" ${e.busy?"disabled":""}>
                ${e.busy?"Importing…":"Import file"}
              </button>
              <button class="button" data-action="overlay-close" type="button" ${e.busy?"disabled":""}>Close</button>
            </div>
          </div>
        </div>
      `;case"notice":return`
        <div class="overlay-backdrop">
          <div class="modal">
            <div class="modal-title">${l(e.title)}</div>
            <div class="modal-copy">${l(e.message)}</div>
            <div class="modal-actions">
              <button class="button button-primary" data-action="overlay-close" type="button">OK</button>
            </div>
          </div>
        </div>
      `}}function L(e,s){var i;const a=e==="quotes"?t.quotes:e==="recall"?t.recallQuotes:((i=t.historyDetail)==null?void 0:i.Quotes)??[],r=new Set(s);return a.filter(n=>r.has(n.ID))}function Jt(e){const s=new Set(e);return t.historyEntries.filter(a=>s.has(a.ID))}function Vt(e){return W(e.settings,[])}function W(e,s){var r;const a={host:e.Provider.Host,port:String(e.Provider.Port),https:e.Provider.HTTPS,apiKey:e.Provider.APIKey,modelFilter:"",model:e.Provider.Model,maxResults:String(e.Search.MaxResults),minRelevance:String(e.Search.MinRelevance),theme:e.Theme||"violet",webPort:String(((r=e.Web)==null?void 0:r.Port)??9527),models:s};return C(a),a}function zt(){return{host:"",port:"11434",https:!1,apiKey:"",modelFilter:"",model:"",maxResults:"5",minRelevance:"0",theme:"violet",webPort:"9527",models:[]}}function J(e){const s=Number.parseInt(e.port.trim(),10);if(!Number.isInteger(s)||s<1||s>65535)throw new Error("Port must be a number between 1 and 65535.");return{Host:e.host.trim(),Port:s,HTTPS:e.https,APIKey:e.apiKey,Model:e.model}}function Gt(e){const s=J(e),a=Number.parseInt(e.maxResults.trim(),10),r=Number.parseInt(e.webPort.trim(),10);if(!Number.isInteger(a)||a<1||a>20)throw new Error("Max ref quotes must be between 1 and 20.");if(!Number.isInteger(r)||r<1||r>65535)throw new Error("Web port must be a number between 1 and 65535.");const i=Number.parseFloat(e.minRelevance.trim());if(Number.isNaN(i))throw new Error("Min relevance must be a decimal number.");if(i<0||i>1)throw new Error("Min relevance must be between 0.0 and 1.0.");return{Provider:s,Search:{MaxResults:a,MinRelevance:i},Theme:e.theme,Web:{Port:r}}}function V(e){const s=e.modelFilter.trim().toLowerCase();return s?e.models.filter(a=>a.toLowerCase().includes(s)):e.models}function C(e){if(e.models.length===0)return;const s=V(e);s.length!==0&&(s.includes(e.model)||(e.model=s[0]))}function w(e,s){return e.map(a=>a.ID===s.ID?s:a)}function v(e,s){return s.length===0?0:Math.min(Math.max(e,0),s.length-1)}function z(e,s){const a=new Set(s.map(r=>r.ID));return new Set([...e].filter(r=>a.has(r)))}function x(e,s){return s.length===0?0:Math.min(Math.max(e,0),s.length-1)}function Yt(e,s){const a=new Set(s.map(r=>r.ID));return new Set([...e].filter(r=>a.has(r)))}function l(e){return e.replaceAll("&","&amp;").replaceAll("<","&lt;").replaceAll(">","&gt;").replaceAll('"',"&quot;").replaceAll("'","&#39;")}function p(e){return l(e)}function m(e,s){const a=e.replace(/\s+/g," ").trim();return a.length<=s?a:`${a.slice(0,s-1).trimEnd()}…`}function H(e){const s=new Date(e);return Number.isNaN(s.getTime())?e:s.toLocaleString()}function Xt(e,s){return e.length===0?"":e.length<=s?e.join(" · "):`${e.slice(0,s).join(" · ")} · +${e.length-s} more`}function S(e,s){return m(e,Math.max(8,s))}function y(){var e;return((e=t.auth)==null?void 0:e.runtime)==="web"}function Zt(e,s){const a=new Blob([s],{type:"application/json;charset=utf-8"}),r=URL.createObjectURL(a),i=document.createElement("a");i.href=r,i.download=e,document.body.appendChild(i),i.click(),i.remove(),URL.revokeObjectURL(r)}function d(e){return e instanceof Error?e.message:String(e)}function u(){var s,a;const e=(a=(s=window.go)==null?void 0:s.backend)==null?void 0:a.App;if(!e)throw new Error("Wails backend bridge is unavailable.");return e}const G=document.querySelector("#app");if(!G)throw new Error("Missing #app root");tt(G);

(function(){const s=document.createElement("link").relList;if(s&&s.supports&&s.supports("modulepreload"))return;for(const o of document.querySelectorAll('link[rel="modulepreload"]'))r(o);new MutationObserver(o=>{for(const n of o)if(n.type==="childList")for(const u of n.addedNodes)u.tagName==="LINK"&&u.rel==="modulepreload"&&r(u)}).observe(document,{childList:!0,subtree:!0});function a(o){const n={};return o.integrity&&(n.integrity=o.integrity),o.referrerPolicy&&(n.referrerPolicy=o.referrerPolicy),o.crossOrigin==="use-credentials"?n.credentials="include":o.crossOrigin==="anonymous"?n.credentials="omit":n.credentials="same-origin",n}function r(o){if(o.ep)return;o.ep=!0;const n=a(o);fetch(o.href,n)}})();const Q={violet:{bg:"#0f172a",bgStrong:"#0b1120",panel:"rgba(17, 24, 39, 0.92)",panel2:"rgba(31, 41, 55, 0.82)",border:"#374151",borderStrong:"rgba(167, 139, 250, 0.42)",primary:"#7c3aed",accent:"#a78bfa",muted:"#94a3b8",fg:"#f9fafb",ok:"#10b981",error:"#ef4444",shadow:"0 24px 80px rgba(2, 6, 23, 0.38)",colorScheme:"dark"},forest:{bg:"#071a17",bgStrong:"#041311",panel:"rgba(9, 24, 21, 0.92)",panel2:"rgba(15, 41, 35, 0.82)",border:"#29443f",borderStrong:"rgba(45, 212, 191, 0.42)",primary:"#0f766e",accent:"#2dd4bf",muted:"#9ca3af",fg:"#ecfdf5",ok:"#22c55e",error:"#ef4444",shadow:"0 24px 80px rgba(1, 10, 9, 0.38)",colorScheme:"dark"},sunset:{bg:"#1c0f0a",bgStrong:"#130905",panel:"rgba(33, 17, 12, 0.94)",panel2:"rgba(52, 28, 18, 0.82)",border:"#5c4033",borderStrong:"rgba(251, 146, 60, 0.44)",primary:"#c2410c",accent:"#fb923c",muted:"#d6b8a6",fg:"#fffbeb",ok:"#16a34a",error:"#dc2626",shadow:"0 24px 80px rgba(20, 8, 2, 0.4)",colorScheme:"dark"},ocean:{bg:"#081824",bgStrong:"#06111a",panel:"rgba(11, 25, 38, 0.92)",panel2:"rgba(18, 40, 56, 0.82)",border:"#334155",borderStrong:"rgba(56, 189, 248, 0.42)",primary:"#0369a1",accent:"#38bdf8",muted:"#94a3b8",fg:"#f8fafc",ok:"#10b981",error:"#ef4444",shadow:"0 24px 80px rgba(3, 9, 16, 0.38)",colorScheme:"dark"},paper:{bg:"#f8fafc",bgStrong:"#e2e8f0",panel:"rgba(255, 255, 255, 0.96)",panel2:"rgba(248, 250, 252, 0.94)",border:"#cbd5e1",borderStrong:"rgba(29, 78, 216, 0.28)",primary:"#1d4ed8",accent:"#0f766e",muted:"#64748b",fg:"#111827",ok:"#15803d",error:"#b91c1c",shadow:"0 24px 80px rgba(148, 163, 184, 0.3)",colorScheme:"light"}};function k(e){const s=Q[e in Q?e:"violet"],a=document.documentElement;a.style.setProperty("--bg",s.bg),a.style.setProperty("--bg-strong",s.bgStrong),a.style.setProperty("--panel",s.panel),a.style.setProperty("--panel-2",s.panel2),a.style.setProperty("--border",s.border),a.style.setProperty("--border-strong",s.borderStrong),a.style.setProperty("--primary",s.primary),a.style.setProperty("--accent",s.accent),a.style.setProperty("--muted",s.muted),a.style.setProperty("--fg",s.fg),a.style.setProperty("--ok",s.ok),a.style.setProperty("--error",s.error),a.style.setProperty("--shadow",s.shadow),a.style.setProperty("color-scheme",s.colorScheme),document.body.dataset.theme=e}function rt(){return Object.keys(Q)}const t={bootstrapped:!1,fatalError:"",authChecked:!1,auth:null,authBusy:!1,authPassword:"",authConfirmPassword:"",authStatus:"",authIsError:!1,page:"Recall",bootstrap:null,quotes:[],quotesLoading:!1,quotesError:"",quotesCursor:0,quotesSelected:new Set,recallQuestion:"",recallLastQuestion:"",recallKeywords:[],recallQuotes:[],recallResponse:"",recallBusy:!1,recallError:"",recallStatus:"",recallStatusIsError:!1,recallCursor:0,recallSelected:new Set,historyEntries:[],historyLoading:!1,historyError:"",historyCursor:0,historySelected:new Set,historyDetail:null,historyDetailLoading:!1,historyDetailError:"",historyStatus:"",historyStatusIsError:!1,historyQuoteCursor:0,historyQuoteSelected:new Set,settings:oe(),settingsShowStats:!1,settingsBusy:!1,settingsStatus:"",settingsIsError:!1,passwordForm:{current:"",next:"",confirm:"",busy:!1,status:"",isError:!1},overlay:null,toast:null};let g=null,x=!1,H=null,w=null;const ot=["Recall","History","Quotes","Settings"];function it(e){g=e,x||(lt(e),x=!0),i(),H||(H=nt())}async function nt(){try{const e=await _();if(t.auth=await e.AuthStatus(),t.authChecked=!0,t.auth.runtime==="web"&&!t.auth.authenticated){i();return}await A()}catch(e){t.authChecked=!0,t.bootstrapped=!0,t.fatalError=d(e),i()}}async function A(){var s;await _();const e=await c().BootstrapState();t.bootstrap=e,t.bootstrapped=!0,t.page="Recall",t.settings=G(e),k(t.settings.theme),(s=e.profile)!=null&&s.DisplayName||(t.overlay={type:"namePrompt",name:"",busy:!1,status:"",isError:!1}),i(),await v()}function lt(e){e.addEventListener("click",s=>{pt(s)}),e.addEventListener("input",yt),e.addEventListener("change",vt),e.addEventListener("submit",s=>{ht(s)}),window.addEventListener("keydown",s=>{bt(s)})}async function N(){if(!(!t.auth||t.authBusy)){if(!t.authPassword.trim()){t.authStatus="Password is required.",t.authIsError=!0,i();return}t.authBusy=!0,t.authStatus="",t.authIsError=!1,i();try{await c().Login(t.authPassword),t.authPassword="",t.authConfirmPassword="",t.auth=await c().AuthStatus(),await A()}catch(e){t.authStatus=d(e),t.authIsError=!0,i()}finally{t.authBusy=!1}}}async function ut(){await c().Logout(),t.auth=await c().AuthStatus(),t.bootstrapped=!1,t.bootstrap=null,t.overlay=null,t.quotes=[],t.historyEntries=[],t.historyDetail=null,t.authPassword="",t.authConfirmPassword="",t.authStatus="",t.authIsError=!1,i()}async function ct(){if(!t.passwordForm.busy){t.passwordForm.busy=!0,t.passwordForm.status="",i();try{await c().ChangePassword(t.passwordForm.current,t.passwordForm.next,t.passwordForm.confirm),t.passwordForm={current:"",next:"",confirm:"",busy:!1,status:"Password updated.",isError:!1}}catch(e){t.passwordForm.busy=!1,t.passwordForm.status=d(e),t.passwordForm.isError=!0}i()}}async function dt(e){var r,o;const s=(r=e.files)==null?void 0:r[0];if(!s||((o=t.overlay)==null?void 0:o.type)!=="importQuotes")return;const a=await s.text();t.overlay.filename=s.name,t.overlay.payload=a,t.overlay.path=s.name,t.overlay.status=`Loaded ${s.name}`,t.overlay.isError=!1,i()}async function pt(e){const s=e.target;if(!(s instanceof HTMLElement))return;const a=s.closest("[data-action]");if(!a)return;switch(a.dataset.action??""){case"auth-login":await N();return;case"auth-logout":await ut();return;case"nav":await ft(a.dataset.page);return;case"quotes-refresh":await v();return;case"history-refresh":await $();return;case"history-view-current":await mt();return;case"history-back":O();return;case"recall-save-quote":await Pt();return;case"history-save-quote":await Rt();return;case"history-delete-current":St();return;case"history-select-all":jt();return;case"history-deselect-all":Kt();return;case"quote-add":j("add");return;case"quote-import":Et();return;case"quote-select-all":Mt(a.dataset.context);return;case"quote-deselect-all":At(a.dataset.context);return;case"quote-edit-current":gt(a.dataset.context);return;case"quote-delete-current":wt(a.dataset.context);return;case"quote-share-current":await $t(a.dataset.context);return;case"set-cursor":if(s.closest("input, button, label"))return;Tt(a.dataset.context,Number(a.dataset.index??"0"));return;case"history-set-cursor":if(s.closest("input, button, label"))return;Bt(Number(a.dataset.index??"0"));return;case"history-open":await B(Number(a.dataset.id??"0"));return;case"profile-save":await K();return;case"quote-editor-save":await U();return;case"quote-editor-refine":await W();return;case"quote-editor-apply-refined":qt();return;case"quote-editor-reject-refined":Qt();return;case"overlay-close":Y();return;case"delete-confirm":await Dt();return;case"share-browse":await kt();return;case"share-save":await It();return;case"import-browse":await Ct();return;case"import-run":await Lt();return;case"settings-fetch-models":await xt();return;case"settings-save":await z();return;case"settings-toggle-stats":t.settingsShowStats=!t.settingsShowStats,i();return;case"settings-change-password":await ct();return;case"recall-run":await I();return;default:return}}function yt(e){var r,o,n,u;const s=e.target;if(!(s instanceof HTMLInputElement||s instanceof HTMLTextAreaElement))return;switch(s.dataset.bind??""){case"auth-password":t.authPassword=s.value;return;case"auth-confirm-password":t.authConfirmPassword=s.value;return;case"recall-question":t.recallQuestion=s.value;return;case"profile-name":((r=t.overlay)==null?void 0:r.type)==="namePrompt"&&(t.overlay.name=s.value);return;case"quote-editor-content":((o=t.overlay)==null?void 0:o.type)==="quoteEditor"&&(t.overlay.content=s.value);return;case"share-path":((n=t.overlay)==null?void 0:n.type)==="shareQuotes"&&(t.overlay.path=s.value);return;case"import-path":((u=t.overlay)==null?void 0:u.type)==="importQuotes"&&(t.overlay.path=s.value);return;case"settings-host":t.settings.host=s.value;return;case"settings-port":t.settings.port=s.value;return;case"settings-api-key":t.settings.apiKey=s.value;return;case"settings-model-filter":t.settings.modelFilter=s.value,E(t.settings),i();return;case"settings-max-results":t.settings.maxResults=s.value;return;case"settings-min-relevance":t.settings.minRelevance=s.value;return;case"settings-theme":t.settings.theme=s.value,k(t.settings.theme);return;case"settings-web-port":t.settings.webPort=s.value;return;case"settings-root-dir":t.settings.rootDir=s.value;return;case"settings-password-current":t.passwordForm.current=s.value;return;case"settings-password-next":t.passwordForm.next=s.value;return;case"settings-password-confirm":t.passwordForm.confirm=s.value;return;default:return}}function vt(e){const s=e.target;if(!(s instanceof HTMLInputElement||s instanceof HTMLSelectElement))return;switch(s.dataset.bind??""){case"quote-selected":Ft(s.dataset.context,Number(s.dataset.id??"0"),s.checked);return;case"history-selected":Ot(Number(s.dataset.id??"0"),s.checked);return;case"settings-https":s instanceof HTMLInputElement&&(t.settings.https=s.checked);return;case"settings-model":t.settings.model=s.value;return;case"settings-mock-llm":s instanceof HTMLInputElement&&(t.settings.mockLLM=s.checked);return;case"import-file":s instanceof HTMLInputElement&&dt(s);return;default:return}}async function ht(e){const s=e.target;if(s instanceof HTMLFormElement)switch(e.preventDefault(),s.dataset.form){case"auth-login":await N();return;case"auth-setup":await submitAuthSetup();return;case"recall":await I();return;case"profile":await K();return;default:return}}async function bt(e){var a,r;const s=document.activeElement;if(e.key==="Escape"&&t.overlay&&t.overlay.type!=="namePrompt"){e.preventDefault(),Y();return}if(e.ctrlKey&&e.key.toLowerCase()==="s"){if(((a=t.overlay)==null?void 0:a.type)==="quoteEditor"){e.preventDefault(),await U();return}!t.overlay&&t.page==="Settings"&&(e.preventDefault(),await z());return}if(e.ctrlKey&&e.key.toLowerCase()==="r"&&((r=t.overlay)==null?void 0:r.type)==="quoteEditor"){e.preventDefault(),await W();return}e.key==="Enter"&&!e.shiftKey&&s instanceof HTMLInputElement&&s.dataset.bind==="recall-question"&&(e.preventDefault(),await I())}async function ft(e){t.page=e,i(),e==="Quotes"&&await v(),e==="History"&&await $()}async function v(){t.quotesLoading=!0,t.quotesError="",i();try{const e=await c().ListQuotes();t.quotes=e,t.quotesCursor=h(t.quotesCursor,e),t.quotesSelected=X(t.quotesSelected,e),t.quotesError=""}catch(e){t.quotesError=d(e)}finally{t.quotesLoading=!1,i()}}async function $(){t.historyLoading=!0,t.historyError="",t.historyStatus="",t.historyStatusIsError=!1,i();try{const e=await c().ListRecallHistory();t.historyEntries=e,t.historyCursor=R(t.historyCursor,e),t.historySelected=ne(t.historySelected,e)}catch(e){t.historyError=d(e)}finally{t.historyLoading=!1,i()}}async function mt(){const e=L()[0];e&&await B(e.ID)}async function B(e){if(!(!Number.isFinite(e)||e<=0)){t.historyDetailLoading=!0,t.historyDetailError="",t.historyStatus="",t.historyStatusIsError=!1,i();try{const s=await c().GetRecallHistory(e);t.historyDetail=s,t.historyQuoteCursor=h(t.historyQuoteCursor,s.Quotes),t.historyQuoteSelected=X(t.historyQuoteSelected,s.Quotes)}catch(s){t.historyDetailError=d(s)}finally{t.historyDetailLoading=!1,i()}}}function O(){t.historyDetail=null,t.historyDetailLoading=!1,t.historyDetailError="",t.historyQuoteCursor=0,t.historyQuoteSelected=new Set,i()}function j(e,s){t.overlay={type:"quoteEditor",mode:e,quoteId:(s==null?void 0:s.ID)??null,content:(s==null?void 0:s.Content)??"",busy:!1,status:"",isError:!1,previewOriginal:"",previewRefined:""},i()}function gt(e){const s=f(e)[0];s&&j("edit",s)}function wt(e){const s=f(e).map(a=>a.ID);s.length!==0&&(t.overlay={type:"deleteQuotes",context:e,ids:s,busy:!1,status:"",isError:!1},i())}function St(){const e=L().map(s=>s.ID);e.length!==0&&(t.overlay={type:"deleteHistory",ids:e,busy:!1,status:"",isError:!1},i())}async function $t(e){var a,r;const s=f(e);if(s.length!==0){t.overlay={type:"shareQuotes",context:e,ids:s.map(o=>o.ID),path:"",payload:"",busy:!0,status:"",isError:!1},i();try{const o=await c().PreviewQuoteExport(s.map(n=>n.ID));if(((a=t.overlay)==null?void 0:a.type)!=="shareQuotes")return;t.overlay.payload=o,t.overlay.busy=!1,t.overlay.status="Share payload ready. Save it to a file and transfer it manually.",t.overlay.isError=!1}catch(o){if(((r=t.overlay)==null?void 0:r.type)!=="shareQuotes")return;t.overlay.busy=!1,t.overlay.status=d(o),t.overlay.isError=!0}i()}}function Et(){t.overlay={type:"importQuotes",path:"",payload:"",filename:"",busy:!1,status:"",isError:!1,result:null},i()}async function K(){var s,a;if(((s=t.overlay)==null?void 0:s.type)!=="namePrompt"||t.overlay.busy)return;const e=t.overlay.name.trim();if(!e){t.overlay.status="Please enter a name to continue.",t.overlay.isError=!0,i();return}t.overlay.busy=!0,t.overlay.status="Saving profile…",t.overlay.isError=!1,i();try{const r=await c().SaveUserProfile(e);t.bootstrap&&(t.bootstrap.profile=r,t.bootstrap.greeting=`Hi! ${r.DisplayName}`),t.overlay=null}catch(r){((a=t.overlay)==null?void 0:a.type)==="namePrompt"&&(t.overlay.busy=!1,t.overlay.status=d(r),t.overlay.isError=!0)}i()}async function U(){var s,a;if(((s=t.overlay)==null?void 0:s.type)!=="quoteEditor"||t.overlay.busy)return;const e=t.overlay.content.trim();if(!e){t.overlay.status="Nothing to save.",t.overlay.isError=!0,i();return}t.overlay.busy=!0,t.overlay.status="Refining draft...",t.overlay.isError=!1,i();try{const r=t.overlay.mode==="add"?await c().AddQuote(e):await c().UpdateQuote(t.overlay.quoteId??0,e);t.overlay=null,C(r),await v()}catch(r){((a=t.overlay)==null?void 0:a.type)==="quoteEditor"&&(t.overlay.busy=!1,t.overlay.status=d(r),t.overlay.isError=!0),i()}}async function W(){var s,a,r;if(((s=t.overlay)==null?void 0:s.type)!=="quoteEditor"||t.overlay.busy)return;const e=t.overlay.content.trim();if(!e){t.overlay.status="Nothing to refine.",t.overlay.isError=!0,i();return}t.overlay.busy=!0,t.overlay.status="",i();try{const o=await c().RefineQuoteDraft(e);if(((a=t.overlay)==null?void 0:a.type)!=="quoteEditor")return;t.overlay.busy=!1,t.overlay.previewOriginal=e,t.overlay.previewRefined=o,t.overlay.status="",t.overlay.isError=!1}catch(o){((r=t.overlay)==null?void 0:r.type)==="quoteEditor"&&(t.overlay.busy=!1,t.overlay.status=d(o),t.overlay.isError=!0)}i()}function qt(){var e;((e=t.overlay)==null?void 0:e.type)==="quoteEditor"&&(t.overlay.content=t.overlay.previewRefined,t.overlay.previewOriginal="",t.overlay.previewRefined="",t.overlay.status="Refined draft applied. Review it, then save.",t.overlay.isError=!1,i())}function Qt(){var e;((e=t.overlay)==null?void 0:e.type)==="quoteEditor"&&(t.overlay.previewOriginal="",t.overlay.previewRefined="",t.overlay.status="Refined draft discarded.",t.overlay.isError=!1,i())}async function Dt(){var e,s,a,r;if(((e=t.overlay)==null?void 0:e.type)==="deleteHistory"){if(t.overlay.busy)return;t.overlay.busy=!0,t.overlay.status="",i();try{await c().DeleteRecallHistory(t.overlay.ids),Ut(t.overlay.ids),t.overlay=null,await $()}catch(o){((s=t.overlay)==null?void 0:s.type)==="deleteHistory"&&(t.overlay.busy=!1,t.overlay.status=d(o),t.overlay.isError=!0),i()}return}if(!(((a=t.overlay)==null?void 0:a.type)!=="deleteQuotes"||t.overlay.busy)){t.overlay.busy=!0,t.overlay.status="",i();try{await c().DeleteQuotes(t.overlay.ids),Nt(t.overlay.ids),t.overlay=null,await v()}catch(o){((r=t.overlay)==null?void 0:r.type)==="deleteQuotes"&&(t.overlay.busy=!1,t.overlay.status=d(o),t.overlay.isError=!0),i()}}}async function kt(){var e,s,a;if(!(((e=t.overlay)==null?void 0:e.type)!=="shareQuotes"||t.overlay.busy)){if(y()){t.overlay.path="irecall-share.json",i();return}try{const r=await c().SelectQuoteExportFile();r&&((s=t.overlay)==null?void 0:s.type)==="shareQuotes"&&(t.overlay.path=r,i())}catch(r){((a=t.overlay)==null?void 0:a.type)==="shareQuotes"&&(t.overlay.status=d(r),t.overlay.isError=!0,i())}}}async function It(){var s,a,r;if(((s=t.overlay)==null?void 0:s.type)!=="shareQuotes"||t.overlay.busy)return;if(y()){const o=t.overlay.path.trim()||"irecall-share.json";if(!t.overlay.payload.trim()){t.overlay.status="Export payload is not ready yet.",t.overlay.isError=!0,i();return}ue(o,t.overlay.payload),t.overlay.status=`Downloaded ${o}`,t.overlay.isError=!1,i();return}const e=t.overlay.path.trim();if(!e){t.overlay.status="Choose a file path for the export.",t.overlay.isError=!0,i();return}if(!t.overlay.payload.trim()){t.overlay.status="Export payload is not ready yet.",t.overlay.isError=!0,i();return}t.overlay.busy=!0,t.overlay.status="",i();try{await c().ExportQuotesToFile(t.overlay.ids,e),((a=t.overlay)==null?void 0:a.type)==="shareQuotes"&&(t.overlay.busy=!1,t.overlay.status=`Saved share payload to ${e}`,t.overlay.isError=!1,i())}catch(o){((r=t.overlay)==null?void 0:r.type)==="shareQuotes"&&(t.overlay.busy=!1,t.overlay.status=d(o),t.overlay.isError=!0,i())}}async function Ct(){var e,s,a;if(!(((e=t.overlay)==null?void 0:e.type)!=="importQuotes"||t.overlay.busy)){if(y()){const r=document.querySelector('[data-bind="import-file"]');r==null||r.click();return}try{const r=await c().SelectQuoteImportFile();r&&((s=t.overlay)==null?void 0:s.type)==="importQuotes"&&(t.overlay.path=r,i())}catch(r){((a=t.overlay)==null?void 0:a.type)==="importQuotes"&&(t.overlay.status=d(r),t.overlay.isError=!0,i())}}}async function Lt(){var a,r,o;if(((a=t.overlay)==null?void 0:a.type)!=="importQuotes"||t.overlay.busy)return;const e=t.overlay.path.trim(),s=t.overlay.payload.trim();if(y()){if(!s){t.overlay.status="Choose a file to import.",t.overlay.isError=!0,i();return}}else if(!e){t.overlay.status="Choose a file to import.",t.overlay.isError=!0,i();return}t.overlay.busy=!0,t.overlay.status="",t.overlay.result=null,i();try{const n=y()?await c().ImportQuotesPayload(s):await c().ImportQuotesFromFile(e);if(((r=t.overlay)==null?void 0:r.type)!=="importQuotes")return;t.overlay.busy=!1,t.overlay.result=n,t.overlay.status=`Imported quotes. inserted=${n.Inserted} updated=${n.Updated} duplicates=${n.Duplicates} stale=${n.Stale}`,t.overlay.isError=!1,await v()}catch(n){((o=t.overlay)==null?void 0:o.type)==="importQuotes"&&(t.overlay.busy=!1,t.overlay.status=d(n),t.overlay.isError=!0,i())}}async function I(){if(t.recallBusy)return;const e=t.recallQuestion.trim();if(!e){t.recallError="Ask a question first.",i();return}t.recallBusy=!0,t.recallError="",t.recallStatus="",t.recallStatusIsError=!1,t.recallLastQuestion=e,t.recallKeywords=[],t.recallQuotes=[],t.recallResponse="",t.recallCursor=0,t.recallSelected=new Set,i();try{const s=await c().RunRecall(e);t.recallKeywords=s.keywords,t.recallQuotes=s.quotes,t.recallResponse=s.response,t.recallLastQuestion=s.question||e,t.recallCursor=0,t.recallSelected=new Set,t.recallQuestion=""}catch(s){t.recallError=d(s)}finally{t.recallBusy=!1,i()}}async function Pt(){const e=t.recallLastQuestion.trim(),s=t.recallResponse.trim();if(!e||!s){t.recallStatus="Run a recall first before saving it as a quote.",t.recallStatusIsError=!0,i();return}try{const a=await c().SaveRecallAsQuote(e,s,t.recallKeywords);C(a),await v(),t.recallStatus="Saved recall as quote.",t.recallStatusIsError=!1,t.overlay={type:"notice",title:"Recall Saved as Quote",message:"The current question and grounded response were saved as a quote with generated tags."}}catch(a){t.recallStatus=d(a),t.recallStatusIsError=!0}i()}async function Rt(){const e=t.historyDetail;if(e){try{const s=await c().SaveRecallAsQuote(e.Question,e.Response,[]);C(s),await v(),t.historyStatus="Saved history entry as quote.",t.historyStatusIsError=!1,t.overlay={type:"notice",title:"History Entry Saved as Quote",message:"The selected history question and response were saved as a quote with generated tags."}}catch(s){t.historyStatus=d(s),t.historyStatusIsError=!0}i()}}async function xt(){if(t.settingsBusy)return;let e;try{e=V(t.settings)}catch(s){t.settingsStatus=d(s),t.settingsIsError=!0,i();return}t.settingsBusy=!0,t.settingsStatus="",i();try{const s=await c().FetchModels(e);t.settings.models=s,E(t.settings),t.settingsStatus=s.length>0?`Fetched ${s.length} models.`:"No models returned.",t.settingsIsError=!1}catch(s){t.settingsStatus=d(s),t.settingsIsError=!0}finally{t.settingsBusy=!1,i()}}async function z(){var s,a;if(t.settingsBusy)return;let e;try{e=ie(t.settings)}catch(r){t.settingsStatus=d(r),t.settingsIsError=!0,i();return}t.settingsBusy=!0,t.settingsStatus="",i();try{const r=F(t.settings.rootDir),o=((s=t.auth)==null?void 0:s.currentPort)??0,n=await c().SaveSettings(e);await Ht(n);const u=r!==F(n.RootDir),b=((a=t.auth)==null?void 0:a.runtime)==="web"&&o>0&&o!==n.Web.Port;u&&b?t.settingsStatus="Saved. Switched storage root. Restart the web server to apply the new port.":u?t.settingsStatus="Saved. Switched storage root.":b?t.settingsStatus="Saved. Restart the web server to apply the new port.":t.settingsStatus="Saved.",t.settingsIsError=!1,re(t.settingsStatus)}catch(r){t.settingsStatus=d(r),t.settingsIsError=!0}finally{t.settingsBusy=!1,i()}}async function Ht(e){var r,o;const s=[...t.settings.models],a=await c().BootstrapState();t.bootstrap=a,t.settings=G(a),t.settings.models=s,E(t.settings),k(t.settings.theme),t.recallKeywords=[],t.recallQuotes=[],t.recallResponse="",t.recallLastQuestion="",t.recallSelected=new Set,t.recallCursor=0,t.historyEntries=[],t.historySelected=new Set,t.historyCursor=0,t.historyDetail=null,t.historyDetailError="",t.historyQuoteSelected=new Set,t.historyQuoteCursor=0,(r=a.profile)!=null&&r.DisplayName?((o=t.overlay)==null?void 0:o.type)==="namePrompt"&&(t.overlay=null):t.overlay={type:"namePrompt",name:"",busy:!1,status:"",isError:!1},await v(),await $(),t.bootstrap&&(t.bootstrap.settings=e)}function Y(){t.overlay&&t.overlay.type!=="namePrompt"&&("busy"in t.overlay&&t.overlay.busy||(t.overlay=null,i()))}function Tt(e,s){var a;if(e==="quotes")t.quotesCursor=h(s,t.quotes);else if(e==="recall")t.recallCursor=h(s,t.recallQuotes);else{const r=((a=t.historyDetail)==null?void 0:a.Quotes)??[];t.historyQuoteCursor=h(s,r)}i()}function Ft(e,s,a){const r=e==="quotes"?t.quotesSelected:e==="recall"?t.recallSelected:t.historyQuoteSelected;a?r.add(s):r.delete(s)}function Mt(e){var r;const s=e==="quotes"?t.quotes:e==="recall"?t.recallQuotes:((r=t.historyDetail)==null?void 0:r.Quotes)??[],a=new Set(s.map(o=>o.ID));e==="quotes"?t.quotesSelected=a:e==="recall"?t.recallSelected=a:t.historyQuoteSelected=a,i()}function At(e){e==="quotes"?t.quotesSelected=new Set:e==="recall"?t.recallSelected=new Set:t.historyQuoteSelected=new Set,i()}function f(e){var n;const s=e==="quotes"?t.quotes:e==="recall"?t.recallQuotes:((n=t.historyDetail)==null?void 0:n.Quotes)??[],a=e==="quotes"?t.quotesCursor:e==="recall"?t.recallCursor:t.historyQuoteCursor,r=e==="quotes"?t.quotesSelected:e==="recall"?t.recallSelected:t.historyQuoteSelected,o=s.filter(u=>r.has(u.ID));return o.length>0?o:s[a]?[s[a]]:[]}function C(e){t.quotes=q(t.quotes,e),t.recallQuotes=q(t.recallQuotes,e),t.historyDetail&&(t.historyDetail={...t.historyDetail,Quotes:q(t.historyDetail.Quotes,e)}),i()}function Nt(e){var a;const s=new Set(e);t.quotes=t.quotes.filter(r=>!s.has(r.ID)),t.recallQuotes=t.recallQuotes.filter(r=>!s.has(r.ID)),t.historyDetail&&(t.historyDetail={...t.historyDetail,Quotes:t.historyDetail.Quotes.filter(r=>!s.has(r.ID))}),t.quotesSelected=new Set([...t.quotesSelected].filter(r=>!s.has(r))),t.recallSelected=new Set([...t.recallSelected].filter(r=>!s.has(r))),t.historyQuoteSelected=new Set([...t.historyQuoteSelected].filter(r=>!s.has(r))),t.quotesCursor=h(t.quotesCursor,t.quotes),t.recallCursor=h(t.recallCursor,t.recallQuotes),t.historyQuoteCursor=h(t.historyQuoteCursor,((a=t.historyDetail)==null?void 0:a.Quotes)??[]),i()}function Bt(e){t.historyCursor=R(e,t.historyEntries),i()}function Ot(e,s){s?t.historySelected.add(e):t.historySelected.delete(e)}function jt(){t.historySelected=new Set(t.historyEntries.map(e=>e.ID)),i()}function Kt(){t.historySelected=new Set,i()}function L(){const e=t.historyEntries.filter(s=>t.historySelected.has(s.ID));return e.length>0?e:t.historyEntries[t.historyCursor]?[t.historyEntries[t.historyCursor]]:[]}function Ut(e){const s=new Set(e);if(t.historyEntries=t.historyEntries.filter(a=>!s.has(a.ID)),t.historySelected=new Set([...t.historySelected].filter(a=>!s.has(a))),t.historyCursor=R(t.historyCursor,t.historyEntries),t.historyDetail&&s.has(t.historyDetail.ID)){O();return}i()}function i(){if(!g)return;const e=Wt();g.innerHTML=Yt(),zt(e)}function Wt(){const e=document.activeElement;if(!(e instanceof HTMLInputElement||e instanceof HTMLTextAreaElement||e instanceof HTMLSelectElement))return null;const s=e.dataset.bind;return s?{selector:`[data-bind="${s}"]`,selectionStart:e instanceof HTMLInputElement||e instanceof HTMLTextAreaElement?e.selectionStart:null,selectionEnd:e instanceof HTMLInputElement||e instanceof HTMLTextAreaElement?e.selectionEnd:null}:null}function zt(e){if(!g||!e)return;const s=g.querySelector(e.selector);(s instanceof HTMLInputElement||s instanceof HTMLTextAreaElement||s instanceof HTMLSelectElement)&&(s.focus({preventScroll:!0}),(s instanceof HTMLInputElement||s instanceof HTMLTextAreaElement)&&e.selectionStart!==null&&e.selectionEnd!==null&&s.setSelectionRange(e.selectionStart,e.selectionEnd))}function Yt(){var s,a,r,o,n;if(!t.authChecked)return`
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
    `;if(((s=t.auth)==null?void 0:s.runtime)==="web"&&!t.auth.authenticated)return Gt();if(!t.bootstrapped)return`
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
          <div class="brand">${l(((o=t.bootstrap)==null?void 0:o.productName)??"iRecall")}</div>
          <div class="muted subtle">Ask questions. Read the answer. Keep the notes that help.</div>
        </div>
        <div class="titlebar-right">
          ${e?`<div class="greeting">${l(e)}</div>`:""}
          ${((n=t.auth)==null?void 0:n.runtime)==="web"?'<button class="button" data-action="auth-logout" type="button">Logout</button>':""}
          <nav class="tabs" aria-label="Primary">
            ${ot.map(u=>`
                  <button
                    class="tab${t.page===u?" active":""}"
                    data-action="nav"
                    data-page="${u}"
                    type="button"
                  >${u}</button>
                `).join("")}
          </nav>
        </div>
      </header>

      <main class="layout">
        ${Vt()}
      </main>

      ${t.overlay?te(t.overlay):""}
      ${t.toast?ee(t.toast):""}
    </div>
  `}function Gt(){var o;const e=!((o=t.auth)!=null&&o.passwordConfigured),s="auth-login";return`
    <div class="shell shell-loading">
      <div class="panel modal">
        <div class="brand">iRecall</div>
        <div class="modal-title">${e?"Finish Setup In Terminal":"Unlock iRecall"}</div>
        <div class="modal-copy">${e?"The web password must be created in the terminal before the server starts listening. Restart the server from a terminal session to finish setup.":"Enter the password to open your notes and questions."}</div>
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
  `}function Vt(){switch(t.page){case"Recall":return Jt();case"Quotes":return Xt();case"History":return Zt();case"Settings":return _t()}}function Jt(){const e=f("recall"),s=t.recallResponse.trim()?l(t.recallResponse):'<span class="muted">Your answer will show up here after you ask a question.</span>',a=t.recallKeywords.length>0?t.recallKeywords.map(o=>`<span class="keyword-chip">${l(o)}</span>`).join(""):'<span class="muted">We will pull out helpful search words for you.</span>',r=e.length>0?`${e.length} quotes selected`:t.recallQuotes.length>0?`${t.recallQuotes.length} quotes found`:"No quotes yet";return`
    <section class="page page-recall">
      <div class="panel page-panel">
        <form class="question-bar" data-form="recall">
          <input
            class="text-input text-input-lg question-input"
            data-bind="recall-question"
            placeholder='Try: "What did I learn about SQLite?"'
            value="${p(t.recallQuestion)}"
          />
          <button class="button button-primary" data-action="recall-run" type="submit" ${t.recallBusy?"disabled":""}>
            ${t.recallBusy?"Thinking...":"Ask"}
          </button>
        </form>

        <div class="toolbar">
          <button class="button" data-action="recall-save-quote" type="button" ${t.recallResponse.trim()?"":"disabled"}>Save as Quote</button>
        </div>

        <div class="recall-grid">
          <section class="panel subpanel">
            <div class="subpanel-header">
              <div class="section-title">Answer</div>
            </div>
            <pre class="response-box">${s}</pre>
            <div class="keyword-row">
              <span class="muted">Search words:</span>
              <div class="keyword-list">${a}</div>
            </div>
          </section>

          <section class="panel subpanel">
            <div class="subpanel-header">
              <div>
                <div class="section-title">Reference Quotes</div>
                <div class="muted">These are the quotes iRecall used to make the answer.</div>
              </div>
              <div class="muted">${r}</div>
            </div>
            <div class="toolbar toolbar-soft">
              <button class="button" data-action="quote-edit-current" data-context="recall" type="button" ${e.length===0?"disabled":""}>Edit Quote</button>
              <button class="button button-danger" data-action="quote-delete-current" data-context="recall" type="button" ${e.length===0?"disabled":""}>Delete Quote</button>
              <button class="button" data-action="quote-share-current" data-context="recall" type="button" ${e.length===0?"disabled":""}>Share Quote</button>
            </div>
            ${P("recall",t.recallQuotes,t.recallCursor,t.recallSelected,!1)}
          </section>
        </div>

        ${t.recallError?`<div class="status status-error">${l(t.recallError)}</div>`:""}
        ${t.recallStatus?`<div class="status ${t.recallStatusIsError?"status-error":"status-ok"}">${l(t.recallStatus)}</div>`:""}
      </div>
    </section>
  `}function Xt(){const e=f("quotes");let s="";return t.quotesLoading?s='<div class="empty-state">Loading your notes...</div>':t.quotesError?s=`<div class="status status-error">${l(t.quotesError)}</div>`:s=P("quotes",t.quotes,t.quotesCursor,t.quotesSelected,!0),`
    <section class="page page-quotes">
      <div class="panel page-panel">
        <div class="section-heading">
          <div>
            <div class="section-title">Quotes</div>
            <div class="muted">Keep short notes here, fix them later, and share them when you want.</div>
          </div>
          <div class="toolbar">
            <button class="button button-primary" data-action="quote-add" type="button">Add Quote</button>
            <button class="button" data-action="quote-import" type="button">Import Quotes</button>
            <button class="button" data-action="quote-share-current" data-context="quotes" type="button" ${e.length===0?"disabled":""}>Share Quote</button>
          </div>
        </div>
        <div class="helper-strip">
          <div>
            <div class="helper-title">Simple way to use this page</div>
            <div class="muted">Add a quote, click a quote to focus it, then edit, delete, or share it.</div>
          </div>
          <div class="helper-actions">
            <button class="button" data-action="quotes-refresh" type="button">Refresh</button>
            <button class="button" data-action="quote-select-all" data-context="quotes" type="button" ${t.quotes.length===0?"disabled":""}>Select All</button>
            <button class="button" data-action="quote-deselect-all" data-context="quotes" type="button" ${t.quotesSelected.size===0?"disabled":""}>Clear Picks</button>
            <button class="button" data-action="quote-edit-current" data-context="quotes" type="button" ${e.length===0?"disabled":""}>Edit Quote</button>
            <button class="button button-danger" data-action="quote-delete-current" data-context="quotes" type="button" ${e.length===0?"disabled":""}>Delete Quote</button>
          </div>
        </div>
        ${s}
      </div>
    </section>
  `}function Zt(){const e=L(),s=f("history");if(t.historyDetailLoading)return`
      <section class="page page-history">
        <div class="panel page-panel">
          <div class="empty-state">Loading this past question...</div>
        </div>
      </section>
    `;if(t.historyDetail)return`
      <section class="page page-history">
        <div class="panel page-panel">
          <div class="section-heading">
            <div>
              <div class="section-title">History</div>
              <div class="muted">Open an older answer and see which notes were used.</div>
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
            <div class="section-title">Question and Response</div>
                <div class="muted">${l(M(t.historyDetail.CreatedAt))}</div>
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
                <div class="muted">${s.length>0?`${s.length} notes selected`:`${t.historyDetail.Quotes.length} notes loaded`}</div>
              </div>
              ${P("history",t.historyDetail.Quotes,t.historyQuoteCursor,t.historyQuoteSelected,!1)}
            </section>
          </div>
          ${t.historyStatus?`<div class="status ${t.historyStatusIsError?"status-error":"status-ok"}">${l(t.historyStatus)}</div>`:""}
        </div>
      </section>
    `;let a="";return t.historyLoading?a='<div class="empty-state">Loading past questions...</div>':t.historyError?a=`<div class="status status-error">${l(t.historyError)}</div>`:t.historyEntries.length===0?a='<div class="empty-state">No past questions yet. Ask something on the Ask page and it will show up here.</div>':a=`
      <div class="history-list">
        ${t.historyEntries.map((r,o)=>{const n=o===t.historyCursor,u=D(r.Response,140);return`
              <article class="quote-card${n?" is-current":""}" data-action="history-set-cursor" data-index="${o}">
                <div class="quote-topline">
                  <label class="selection-toggle">
                    <input
                      type="checkbox"
                      data-bind="history-selected"
                      data-id="${r.ID}"
                      ${t.historySelected.has(r.ID)?"checked":""}
                    />
                  </label>
                  <div class="quote-topline-meta">
                    <span class="quote-index${n?" is-current":""}">${n?"&gt; ":""}Question ${o+1}</span>
                    <span class="quote-version">${l(M(r.CreatedAt))}</span>
                  </div>
                </div>
                <div class="quote-content">${l(D(r.Question,120))}</div>
                <div class="quote-meta"><span class="muted">Answer preview:</span> <span>${l(u||"(empty response)")}</span></div>
                <div class="toolbar toolbar-inline">
                  <button class="button" data-action="history-open" data-id="${r.ID}" type="button">Open</button>
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
            <div class="muted">Look back at what you asked before and reopen the answers any time.</div>
          </div>
          <div class="toolbar">
            <button class="button" data-action="history-refresh" type="button">Refresh</button>
            <button class="button" data-action="history-select-all" type="button" ${t.historyEntries.length===0?"disabled":""}>Select All</button>
            <button class="button" data-action="history-deselect-all" type="button" ${t.historySelected.size===0?"disabled":""}>Clear Picks</button>
            <button class="button" data-action="history-view-current" type="button" ${e.length===0?"disabled":""}>Open</button>
            <button class="button button-danger" data-action="history-delete-current" type="button" ${e.length===0?"disabled":""}>Delete</button>
          </div>
        </div>
        <div class="stat-grid stat-grid-wide">
          ${m("History entries",String(t.historyEntries.length))}
          ${m("Picked now",String(e.length>0?e.length:t.historyEntries.length>0?1:0))}
        </div>
        ${a}
      </div>
    </section>
  `}function _t(){var o,n;const e=J(t.settings),s=(o=t.bootstrap)==null?void 0:o.paths,a=(n=t.auth)==null?void 0:n.currentPort,r=t.settings.models.length>0&&e.length>0?`
        <select class="select-input" data-bind="settings-model">
          ${e.map(u=>`
                <option value="${p(u)}"${u===t.settings.model?" selected":""}>${l(u)}</option>
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
            <div class="muted">Keep your name, look and feel, and advanced AI settings in one place.</div>
          </div>
          <div class="toolbar">
            <button class="button" data-action="settings-toggle-stats" type="button">${t.settingsShowStats?"Hide Stats":"Show Stats"}</button>
            <button class="button button-primary" data-action="settings-save" type="button" ${t.settingsBusy?"disabled":""}>Save</button>
          </div>
        </div>

        ${t.settingsShowStats?`
              <div class="stat-grid stat-grid-wide">
                ${m("Stored quotes",String(t.quotes.length))}
                ${m("Stored history",String(t.historyEntries.length))}
                ${m("Reference quotes now",String(t.recallQuotes.length))}
              </div>
            `:""}

        <div class="settings-grid">
          <section class="panel subpanel">
            <div class="section-title">Advanced AI Setup</div>
            <label class="field">
              <span>Host or IP</span>
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
              <span>Find model</span>
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
            <div class="section-title">How Answers Search</div>
            <label class="field">
              <span>How many notes to use</span>
              <input class="text-input" data-bind="settings-max-results" value="${p(t.settings.maxResults)}" />
            </label>
            <label class="field">
              <span>How close the match should be</span>
              <input class="text-input" data-bind="settings-min-relevance" value="${p(t.settings.minRelevance)}" placeholder="0.0-1.0" />
            </label>
            <div class="settings-hint muted">
              0.0 keeps broad matches. Try 0.3 to 0.7 for cleaner results. 1.0 is very strict and may hide useful notes.
            </div>
            <label class="field checkbox-field">
              <input type="checkbox" data-bind="settings-mock-llm"${t.settings.mockLLM?" checked":""} />
              <span>Use Mock LLM</span>
            </label>
            <div class="settings-hint muted">
              Same as TUI mock mode: refine keeps the original quote, keywords split on spaces, and answers are built from matching quotes.
            </div>
          </section>

          <section class="panel subpanel">
            <div class="section-title">Everyday Setup</div>
            <label class="field">
              <span>Theme</span>
              <select class="select-input" data-bind="settings-theme">
                ${rt().map(u=>`
                      <option value="${u}"${u===t.settings.theme?" selected":""}>${u}</option>
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
              <label class="field">
                <span>Config folder root</span>
                <input
                  class="text-input"
                  data-bind="settings-root-dir"
                  value="${p(t.settings.rootDir)}"
                  placeholder="Leave blank to use the default XDG/AppData folders"
                />
              </label>
              <div class="settings-hint muted">
                When you save, iRecall switches to this root immediately. A new root gets <code>data</code>, <code>config</code>, and <code>state</code> subfolders.
              </div>
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

        ${t.settingsStatus&&t.settingsIsError?`<div class="status status-error">${l(t.settingsStatus)}</div>`:""}
      </div>
    </section>
  `}function P(e,s,a,r,o){return s.length===0?`<div class="empty-state">${e==="quotes"?"No quotes yet. Add one or import a shared payload.":"No matching quotes yet for this question."}</div>`:`
    <div class="quote-list">
      ${s.map((n,u)=>{const b=u===a,et=n.IsOwnedByMe?'<span class="pill-badge">Your quote</span>':'<span class="pill-badge pill-badge-soft">Shared quote</span>',st=!n.IsOwnedByMe&&n.SourceName?`<div class="quote-meta"><span class="muted">From:</span> <span class="meta-accent">${l(n.SourceName)}</span></div>`:"",at=o?`
              <div class="quote-meta">
                <span class="muted">Tags:</span>
                <span>${n.Tags.length>0?l(le(n.Tags,3)):"(none)"}</span>
              </div>
            `:"";return`
            <article class="quote-card${b?" is-current":""}" data-action="set-cursor" data-context="${e}" data-index="${u}">
              <div class="quote-topline">
                <label class="selection-toggle">
                  <input
                    type="checkbox"
                    data-bind="quote-selected"
                    data-context="${e}"
                    data-id="${n.ID}"
                    ${r.has(n.ID)?"checked":""}
                  />
                </label>
                <div class="quote-topline-meta">
                  <span class="quote-index${b?" is-current":""}">${b?"&gt; ":""}Quote ${u+1}</span>
                  <span class="quote-version">v${n.Version}</span>
                  ${et}
                </div>
              </div>
              <div class="quote-content">${l(D(n.Content,e==="quotes"?96:120))}</div>
              ${st}
              ${at}
            </article>
          `}).join("")}
    </div>
  `}function te(e){switch(e.type){case"namePrompt":return`
        <div class="overlay-backdrop">
          <div class="modal">
            <div class="modal-title">Tell iRecall Your Name</div>
            <p class="modal-copy">
              Your name is added to quotes you share so other people know where they came from.
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
                    <textarea class="text-area" data-bind="quote-editor-content" rows="10" placeholder="Type or paste your quote here.">${l(e.content)}</textarea>
                  </label>
                `}
            <div class="muted modal-copy">
              ${e.previewRefined?"Compare your draft with the suggested clearer version before you choose one.":"Write a short quote in your own words. Helpful tags are added automatically."}
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
              ${T(e.context,e.ids).map((s,a)=>`<div class="summary-item">[${a+1}] ${l(S(s.Content,140))}</div>`).join("")}
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
              ${se(e.ids).map((s,a)=>`<div class="summary-item">[${a+1}] ${l(S(s.Question,140))}</div>`).join("")}
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
              ${T(e.context,e.ids).map((s,a)=>`<div class="summary-item">[${a+1}] v${s.Version} ${l(S(s.Content,120))}</div>`).join("")}
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
            <div class="muted modal-copy">${y()?"Download the quote file, then send it to someone manually.":"Save the quote file, then send it to someone manually."}</div>
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
            <div class="modal-copy">Open a shared iRecall quote file from another device or person.</div>
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
      `}}function ee(e){return`
    <div class="toast-layer" aria-live="polite" aria-atomic="true">
      <div class="toast ${e.isError?"toast-error":"toast-ok"}">${l(e.message)}</div>
    </div>
  `}function T(e,s){var o;const a=e==="quotes"?t.quotes:e==="recall"?t.recallQuotes:((o=t.historyDetail)==null?void 0:o.Quotes)??[],r=new Set(s);return a.filter(n=>r.has(n.ID))}function se(e){const s=new Set(e);return t.historyEntries.filter(a=>s.has(a.ID))}function G(e){return ae(e.settings,[])}function ae(e,s){var r,o;const a={host:e.Provider.Host,port:String(e.Provider.Port),https:e.Provider.HTTPS,apiKey:e.Provider.APIKey,modelFilter:"",model:e.Provider.Model,maxResults:String(e.Search.MaxResults),minRelevance:String(e.Search.MinRelevance),mockLLM:((r=e.Debug)==null?void 0:r.MockLLM)??!1,theme:e.Theme||"violet",webPort:String(((o=e.Web)==null?void 0:o.Port)??9527),rootDir:e.RootDir??"",models:s};return E(a),a}function re(e,s=!1){t.toast={message:e,isError:s},w!==null&&window.clearTimeout(w),w=window.setTimeout(()=>{t.toast=null,w=null,i()},2200)}function oe(){return{host:"",port:"11434",https:!1,apiKey:"",modelFilter:"",model:"",maxResults:"5",minRelevance:"0",mockLLM:!1,theme:"violet",webPort:"9527",rootDir:"",models:[]}}function V(e){const s=Number.parseInt(e.port.trim(),10);if(!Number.isInteger(s)||s<1||s>65535)throw new Error("Port must be a number between 1 and 65535.");return{Host:e.host.trim(),Port:s,HTTPS:e.https,APIKey:e.apiKey,Model:e.model}}function ie(e){const s=V(e),a=Number.parseInt(e.maxResults.trim(),10),r=Number.parseInt(e.webPort.trim(),10);if(!Number.isInteger(a)||a<1||a>20)throw new Error("Max ref quotes must be between 1 and 20.");if(!Number.isInteger(r)||r<1||r>65535)throw new Error("Web port must be a number between 1 and 65535.");const o=Number.parseFloat(e.minRelevance.trim());if(Number.isNaN(o))throw new Error("Min relevance must be a decimal number.");if(o<0||o>1)throw new Error("Min relevance must be between 0.0 and 1.0.");return{Provider:s,Search:{MaxResults:a,MinRelevance:o},Debug:{MockLLM:e.mockLLM},Theme:e.theme,Web:{Port:r},RootDir:e.rootDir.trim()}}function F(e){return e.trim()}function J(e){const s=e.modelFilter.trim().toLowerCase();return s?e.models.filter(a=>a.toLowerCase().includes(s)):e.models}function E(e){if(e.models.length===0)return;const s=J(e);s.length!==0&&(s.includes(e.model)||(e.model=s[0]))}function q(e,s){return e.map(a=>a.ID===s.ID?s:a)}function h(e,s){return s.length===0?0:Math.min(Math.max(e,0),s.length-1)}function X(e,s){const a=new Set(s.map(r=>r.ID));return new Set([...e].filter(r=>a.has(r)))}function R(e,s){return s.length===0?0:Math.min(Math.max(e,0),s.length-1)}function ne(e,s){const a=new Set(s.map(r=>r.ID));return new Set([...e].filter(r=>a.has(r)))}function l(e){return e.replaceAll("&","&amp;").replaceAll("<","&lt;").replaceAll(">","&gt;").replaceAll('"',"&quot;").replaceAll("'","&#39;")}function p(e){return l(e)}function S(e,s){const a=e.replace(/\s+/g," ").trim();return a.length<=s?a:`${a.slice(0,s-1).trimEnd()}…`}function M(e){const s=new Date(e);return Number.isNaN(s.getTime())?e:s.toLocaleString()}function le(e,s){return e.length===0?"":e.length<=s?e.join(" · "):`${e.slice(0,s).join(" · ")} · +${e.length-s} more`}function m(e,s){return`
    <div class="mini-stat">
      <div class="mini-stat-value">${l(s)}</div>
      <div class="mini-stat-label">${l(e)}</div>
    </div>
  `}function D(e,s){return S(e,Math.max(8,s))}function y(){var e;return((e=t.auth)==null?void 0:e.runtime)==="web"}function ue(e,s){const a=new Blob([s],{type:"application/json;charset=utf-8"}),r=URL.createObjectURL(a),o=document.createElement("a");o.href=r,o.download=e,document.body.appendChild(o),o.click(),o.remove(),URL.revokeObjectURL(r)}function d(e){return e instanceof Error?e.message:String(e)}function Z(){var s,a,r;const e=[(s=window.go)==null?void 0:s.backend,(a=window.go)==null?void 0:a.app,(r=window.go)==null?void 0:r.main];for(const o of e)if(o!=null&&o.App)return o.App;return null}async function _(e=3e3){const s=Date.now();for(;;){const a=Z();if(a)return a;if(Date.now()-s>=e)throw new Error("Wails backend bridge is unavailable.");await new Promise(r=>window.setTimeout(r,25))}}function c(){const e=Z();if(!e)throw new Error("Wails backend bridge is unavailable.");return e}const tt=document.querySelector("#app");if(!tt)throw new Error("Missing #app root");it(tt);

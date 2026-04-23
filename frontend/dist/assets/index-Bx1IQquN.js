(function(){const s=document.createElement("link").relList;if(s&&s.supports&&s.supports("modulepreload"))return;for(const i of document.querySelectorAll('link[rel="modulepreload"]'))r(i);new MutationObserver(i=>{for(const l of i)if(l.type==="childList")for(const d of l.addedNodes)d.tagName==="LINK"&&d.rel==="modulepreload"&&r(d)}).observe(document,{childList:!0,subtree:!0});function a(i){const l={};return i.integrity&&(l.integrity=i.integrity),i.referrerPolicy&&(l.referrerPolicy=i.referrerPolicy),i.crossOrigin==="use-credentials"?l.credentials="include":i.crossOrigin==="anonymous"?l.credentials="omit":l.credentials="same-origin",l}function r(i){if(i.ep)return;i.ep=!0;const l=a(i);fetch(i.href,l)}})();const q={violet:{bg:"#0f172a",bgStrong:"#0b1120",panel:"rgba(17, 24, 39, 0.92)",panel2:"rgba(31, 41, 55, 0.82)",border:"#374151",borderStrong:"rgba(167, 139, 250, 0.42)",primary:"#7c3aed",accent:"#a78bfa",muted:"#94a3b8",fg:"#f9fafb",ok:"#10b981",error:"#ef4444",shadow:"0 24px 80px rgba(2, 6, 23, 0.38)",colorScheme:"dark"},forest:{bg:"#071a17",bgStrong:"#041311",panel:"rgba(9, 24, 21, 0.92)",panel2:"rgba(15, 41, 35, 0.82)",border:"#29443f",borderStrong:"rgba(45, 212, 191, 0.42)",primary:"#0f766e",accent:"#2dd4bf",muted:"#9ca3af",fg:"#ecfdf5",ok:"#22c55e",error:"#ef4444",shadow:"0 24px 80px rgba(1, 10, 9, 0.38)",colorScheme:"dark"},sunset:{bg:"#1c0f0a",bgStrong:"#130905",panel:"rgba(33, 17, 12, 0.94)",panel2:"rgba(52, 28, 18, 0.82)",border:"#5c4033",borderStrong:"rgba(251, 146, 60, 0.44)",primary:"#c2410c",accent:"#fb923c",muted:"#d6b8a6",fg:"#fffbeb",ok:"#16a34a",error:"#dc2626",shadow:"0 24px 80px rgba(20, 8, 2, 0.4)",colorScheme:"dark"},ocean:{bg:"#081824",bgStrong:"#06111a",panel:"rgba(11, 25, 38, 0.92)",panel2:"rgba(18, 40, 56, 0.82)",border:"#334155",borderStrong:"rgba(56, 189, 248, 0.42)",primary:"#0369a1",accent:"#38bdf8",muted:"#94a3b8",fg:"#f8fafc",ok:"#10b981",error:"#ef4444",shadow:"0 24px 80px rgba(3, 9, 16, 0.38)",colorScheme:"dark"},paper:{bg:"#f8fafc",bgStrong:"#e2e8f0",panel:"rgba(255, 255, 255, 0.96)",panel2:"rgba(248, 250, 252, 0.94)",border:"#cbd5e1",borderStrong:"rgba(29, 78, 216, 0.28)",primary:"#1d4ed8",accent:"#0f766e",muted:"#64748b",fg:"#111827",ok:"#15803d",error:"#b91c1c",shadow:"0 24px 80px rgba(148, 163, 184, 0.3)",colorScheme:"light"}};function I(t){const s=q[t in q?t:"violet"],a=document.documentElement;a.style.setProperty("--bg",s.bg),a.style.setProperty("--bg-strong",s.bgStrong),a.style.setProperty("--panel",s.panel),a.style.setProperty("--panel-2",s.panel2),a.style.setProperty("--border",s.border),a.style.setProperty("--border-strong",s.borderStrong),a.style.setProperty("--primary",s.primary),a.style.setProperty("--accent",s.accent),a.style.setProperty("--muted",s.muted),a.style.setProperty("--fg",s.fg),a.style.setProperty("--ok",s.ok),a.style.setProperty("--error",s.error),a.style.setProperty("--shadow",s.shadow),a.style.setProperty("color-scheme",s.colorScheme),document.body.dataset.theme=t}function ne(){return Object.keys(q)}const e={bootstrapped:!1,fatalError:"",authChecked:!1,auth:null,authBusy:!1,authPassword:"",authConfirmPassword:"",authStatus:"",authIsError:!1,page:"Recall",bootstrap:null,quotes:[],quotesLoading:!1,quotesError:"",quotesCursor:0,quotesSelected:new Set,libraryQuery:"",recallQuestion:"",recallLastQuestion:"",recallKeywords:[],recallQuotes:[],recallResponse:"",recallBusy:!1,recallError:"",recallStatus:"",recallStatusIsError:!1,recallCursor:0,recallSelected:new Set,historyEntries:[],historyLoading:!1,historyError:"",historyCursor:0,historySelected:new Set,historyDetail:null,historyDetailLoading:!1,historyDetailError:"",historyStatus:"",historyStatusIsError:!1,historyQuoteCursor:0,historyQuoteSelected:new Set,settings:yt(),settingsBusy:!1,settingsStatus:"",settingsIsError:!1,passwordForm:{current:"",next:"",confirm:"",busy:!1,status:"",isError:!1},apiToken:{loading:!1,hasToken:!1,tokenPrefix:""},overlay:null,toast:null};let b=null,F=!1,N=null,g=null;const le=["Recall","History","Quotes","Settings"];function ue(t){b=t,F||(de(t),F=!0),o(),N||(N=ce())}async function ce(){try{const t=await re();if(e.auth=await t.AuthStatus(),e.authChecked=!0,e.auth.runtime==="web"&&!e.auth.authenticated){o();return}await j()}catch(t){e.authChecked=!0,e.bootstrapped=!0,e.fatalError=c(t),o()}}async function j(){var s;await re();const t=await u().BootstrapState();e.bootstrap=t,e.bootstrapped=!0,e.page="Recall",e.settings=pt(t),I(e.settings.theme),(s=t.profile)!=null&&s.DisplayName||(e.overlay={type:"namePrompt",name:"",busy:!1,status:"",isError:!1}),o(),await U(),await f()}function de(t){t.addEventListener("click",s=>{fe(s)}),t.addEventListener("input",be),t.addEventListener("change",me),t.addEventListener("submit",s=>{ge(s)}),window.addEventListener("keydown",s=>{we(s)})}async function K(){if(!(!e.auth||e.authBusy)){if(!e.authPassword.trim()){e.authStatus="Password is required.",e.authIsError=!0,o();return}e.authBusy=!0,e.authStatus="",e.authIsError=!1,o();try{await u().Login(e.authPassword),e.authPassword="",e.authConfirmPassword="",e.auth=await u().AuthStatus(),await j()}catch(t){e.authStatus=c(t),e.authIsError=!0,o()}finally{e.authBusy=!1}}}async function pe(){await u().Logout(),e.auth=await u().AuthStatus(),e.bootstrapped=!1,e.bootstrap=null,e.overlay=null,e.quotes=[],e.historyEntries=[],e.historyDetail=null,e.authPassword="",e.authConfirmPassword="",e.authStatus="",e.authIsError=!1,e.apiToken={loading:!1,hasToken:!1,tokenPrefix:""},o()}async function U(){e.apiToken.loading=!0,o();try{const t=await u().GetAPITokenStatus();e.apiToken={loading:!1,hasToken:t.hasToken,tokenPrefix:t.tokenPrefix}}catch(t){e.apiToken.loading=!1,e.settingsStatus=c(t),e.settingsIsError=!0}o()}async function ye(){if(e.settingsBusy)return;const t=e.apiToken.hasToken;e.settingsBusy=!0,e.settingsStatus="",e.settingsIsError=!1,o();try{const s=await u().CreateAPIToken();e.apiToken={loading:!1,hasToken:!0,tokenPrefix:s.tokenPrefix},e.overlay={type:"apiTokenReveal",token:s.token,tokenPrefix:s.tokenPrefix},e.settingsStatus=t?"API token renewed.":"API token created.",e.settingsIsError=!1}catch(s){e.settingsStatus=c(s),e.settingsIsError=!0}finally{e.settingsBusy=!1,o()}}async function ve(){if(!e.passwordForm.busy){e.passwordForm.busy=!0,e.passwordForm.status="",o();try{await u().ChangePassword(e.passwordForm.current,e.passwordForm.next,e.passwordForm.confirm),e.passwordForm={current:"",next:"",confirm:"",busy:!1,status:"Password updated.",isError:!1}}catch(t){e.passwordForm.busy=!1,e.passwordForm.status=c(t),e.passwordForm.isError=!0}o()}}async function he(t){var r,i;const s=(r=t.files)==null?void 0:r[0];if(!s||((i=e.overlay)==null?void 0:i.type)!=="importQuotes")return;const a=await s.text();e.overlay.filename=s.name,e.overlay.payload=a,e.overlay.path=s.name,e.overlay.status=`Loaded ${s.name}`,e.overlay.isError=!1,o()}async function fe(t){var i,l;const s=t.target;if(!(s instanceof HTMLElement))return;const a=s.closest("[data-action]");if(!a)return;switch(a.dataset.action??""){case"auth-login":await K();return;case"auth-logout":await pe();return;case"nav":await $e(a.dataset.page);return;case"quotes-refresh":await f();return;case"history-refresh":await P();return;case"history-view-current":await Se();return;case"history-back":$();return;case"recall-save-quote":await He();return;case"history-save-quote":await Ae();return;case"history-delete-current":qe();return;case"history-select-all":Je();return;case"history-deselect-all":Ye();return;case"quote-add":W("add");return;case"quote-import":De();return;case"library-clear-filters":e.libraryQuery="",e.quotesCursor=0,o();return;case"quote-select-all":Oe(a.dataset.context);return;case"quote-deselect-all":je(a.dataset.context);return;case"quote-edit-current":Ee(a.dataset.context);return;case"quote-delete-current":ke(a.dataset.context);return;case"quote-share-current":await Qe(a.dataset.context);return;case"quote-inspect":if(s.closest("input, button, label"))return;Ie(a.dataset.context,Number(a.dataset.index??"0"));return;case"set-cursor":if(s.closest("input, button, label"))return;Ne(a.dataset.context,Number(a.dataset.index??"0"));return;case"history-set-cursor":if(s.closest("input, button, label"))return;We(Number(a.dataset.index??"0"));return;case"history-open":await S(Number(a.dataset.id??"0"));return;case"share-toggle-payload":((i=e.overlay)==null?void 0:i.type)==="shareQuotes"&&(e.overlay.showPayload=!e.overlay.showPayload,o());return;case"import-toggle-payload":((l=e.overlay)==null?void 0:l.type)==="importQuotes"&&(e.overlay.showPayload=!e.overlay.showPayload,o());return;case"profile-save":await z();return;case"quote-editor-save":await J();return;case"quote-editor-refine":await Y();return;case"quote-editor-apply-refined":Pe();return;case"quote-editor-reject-refined":Le();return;case"overlay-close":V();return;case"delete-confirm":await Ce();return;case"share-browse":await xe();return;case"share-save":await Re();return;case"import-browse":await Te();return;case"import-run":await Me();return;case"settings-fetch-models":await Fe();return;case"settings-save":await G();return;case"settings-change-password":await ve();return;case"settings-create-api-token":await ye();return;case"recall-run":await L();return;case"use-last-question":e.recallQuestion=e.recallLastQuestion,o();return;case"reuse-history-question":if(e.historyDetail)e.recallQuestion=e.historyDetail.Question;else{const d=R()[0];d&&(e.recallQuestion=d.Question)}e.page="Recall",o();return;default:return}}function be(t){var r,i,l,d;const s=t.target;if(!(s instanceof HTMLInputElement||s instanceof HTMLTextAreaElement))return;switch(s.dataset.bind??""){case"auth-password":e.authPassword=s.value;return;case"auth-confirm-password":e.authConfirmPassword=s.value;return;case"recall-question":e.recallQuestion=s.value;return;case"library-query":e.libraryQuery=s.value,e.quotesCursor=0,o();return;case"profile-name":((r=e.overlay)==null?void 0:r.type)==="namePrompt"&&(e.overlay.name=s.value);return;case"quote-editor-content":((i=e.overlay)==null?void 0:i.type)==="quoteEditor"&&(e.overlay.content=s.value);return;case"share-path":((l=e.overlay)==null?void 0:l.type)==="shareQuotes"&&(e.overlay.path=s.value);return;case"import-path":((d=e.overlay)==null?void 0:d.type)==="importQuotes"&&(e.overlay.path=s.value);return;case"settings-host":e.settings.host=s.value;return;case"settings-port":e.settings.port=s.value;return;case"settings-api-key":e.settings.apiKey=s.value;return;case"settings-model-filter":e.settings.modelFilter=s.value,M(e.settings),o();return;case"settings-max-results":e.settings.maxResults=s.value;return;case"settings-min-relevance":e.settings.minRelevance=s.value;return;case"settings-theme":e.settings.theme=s.value,I(e.settings.theme);return;case"settings-web-port":e.settings.webPort=s.value;return;case"settings-password-current":e.passwordForm.current=s.value;return;case"settings-password-next":e.passwordForm.next=s.value;return;case"settings-password-confirm":e.passwordForm.confirm=s.value;return;default:return}}function me(t){const s=t.target;if(!(s instanceof HTMLInputElement||s instanceof HTMLSelectElement))return;switch(s.dataset.bind??""){case"quote-selected":Be(s.dataset.context,Number(s.dataset.id??"0"),s.checked);return;case"history-selected":ze(Number(s.dataset.id??"0"),s.checked);return;case"settings-https":s instanceof HTMLInputElement&&(e.settings.https=s.checked);return;case"settings-mock-llm":s instanceof HTMLInputElement&&(e.settings.mockLLM=s.checked);return;case"settings-model":e.settings.model=s.value;return;case"import-file":s instanceof HTMLInputElement&&he(s);return;default:return}}async function ge(t){const s=t.target;if(s instanceof HTMLFormElement)switch(t.preventDefault(),s.dataset.form){case"auth-login":await K();return;case"auth-setup":await submitAuthSetup();return;case"recall":await L();return;case"profile":await z();return;default:return}}async function we(t){var a,r;const s=document.activeElement;if(t.key==="Escape"&&e.overlay&&e.overlay.type!=="namePrompt"){t.preventDefault(),V();return}if(t.ctrlKey&&t.key.toLowerCase()==="s"){if(((a=e.overlay)==null?void 0:a.type)==="quoteEditor"){t.preventDefault(),await J();return}!e.overlay&&e.page==="Settings"&&(t.preventDefault(),await G());return}if(t.ctrlKey&&t.key.toLowerCase()==="r"&&((r=e.overlay)==null?void 0:r.type)==="quoteEditor"){t.preventDefault(),await Y();return}t.key==="Enter"&&!t.shiftKey&&s instanceof HTMLInputElement&&s.dataset.bind==="recall-question"&&(t.preventDefault(),await L())}async function $e(t){e.page=t,o(),t==="Quotes"&&await f(),t==="History"&&await P(),t==="Settings"&&await U()}async function f(){e.quotesLoading=!0,e.quotesError="",o();try{const t=await u().ListQuotes();e.quotes=t,e.quotesCursor=h(e.quotesCursor,t),e.quotesSelected=se(e.quotesSelected,t),e.quotesError=""}catch(t){e.quotesError=c(t)}finally{e.quotesLoading=!1,o()}}async function P(){var t;e.historyLoading=!0,e.historyError="",e.historyStatus="",e.historyStatusIsError=!1,o();try{const s=await u().ListRecallHistory(),a=((t=e.historyDetail)==null?void 0:t.ID)??null;e.historyEntries=s,e.historyCursor=H(e.historyCursor,s),e.historySelected=ht(e.historySelected,s),a===null?$():s.some(r=>r.ID===a)?S(a,!0):$()}catch(s){e.historyError=c(s)}finally{e.historyLoading=!1,o()}}async function Se(){const t=R()[0];t&&await S(t.ID)}async function S(t,s=!1){var a;if(!(!Number.isFinite(t)||t<=0)&&!(e.historyDetailLoading&&((a=e.historyDetail)==null?void 0:a.ID)===t)&&!(e.historyDetail&&e.historyDetail.ID===t&&!e.historyDetailError)){e.historyDetailLoading=!0,e.historyDetailError="",s||(e.historyStatus="",e.historyStatusIsError=!1),o();try{const r=await u().GetRecallHistory(t);e.historyDetail=r,e.historyQuoteCursor=h(e.historyQuoteCursor,r.Quotes),e.historyQuoteSelected=se(e.historyQuoteSelected,r.Quotes)}catch(r){e.historyDetailError=c(r)}finally{e.historyDetailLoading=!1,o()}}}function $(){e.historyDetail=null,e.historyDetailLoading=!1,e.historyDetailError="",e.historyQuoteCursor=0,e.historyQuoteSelected=new Set,o()}function W(t,s){e.overlay={type:"quoteEditor",mode:t,quoteId:(s==null?void 0:s.ID)??null,content:(s==null?void 0:s.Content)??"",busy:!1,status:"",isError:!1,previewOriginal:"",previewRefined:""},o()}function Ee(t){const s=Ke(t);s&&W("edit",s)}function ke(t){var a;const s=((a=e.overlay)==null?void 0:a.type)==="quoteInspect"&&e.overlay.context===t?[e.overlay.quote.ID]:C(t).map(r=>r.ID);s.length!==0&&(e.overlay={type:"deleteQuotes",context:t,ids:s,busy:!1,status:"",isError:!1},o())}function qe(){const t=R().map(s=>s.ID);t.length!==0&&(e.overlay={type:"deleteHistory",ids:t,busy:!1,status:"",isError:!1},o())}async function Qe(t){var a,r,i;const s=((a=e.overlay)==null?void 0:a.type)==="quoteInspect"&&e.overlay.context===t?[e.overlay.quote]:C(t);if(s.length!==0){e.overlay={type:"shareQuotes",context:t,ids:s.map(l=>l.ID),path:"",payload:"",showPayload:!1,busy:!0,status:"",isError:!1},o();try{const l=await u().PreviewQuoteExport(s.map(d=>d.ID));if(((r=e.overlay)==null?void 0:r.type)!=="shareQuotes")return;e.overlay.payload=l,e.overlay.busy=!1,e.overlay.status="Share payload ready. Save it to a file and transfer it manually.",e.overlay.isError=!1}catch(l){if(((i=e.overlay)==null?void 0:i.type)!=="shareQuotes")return;e.overlay.busy=!1,e.overlay.status=c(l),e.overlay.isError=!0}o()}}function De(){e.overlay={type:"importQuotes",path:"",payload:"",filename:"",showPayload:!1,busy:!1,status:"",isError:!1,result:null},o()}function Ie(t,s){const a=E(t),r=h(s,a),i=a[r];i&&(t==="quotes"?e.quotesCursor=r:t==="recall"?e.recallCursor=r:e.historyQuoteCursor=r,e.overlay={type:"quoteInspect",context:t,quote:i},o())}async function z(){var s,a;if(((s=e.overlay)==null?void 0:s.type)!=="namePrompt"||e.overlay.busy)return;const t=e.overlay.name.trim();if(!t){e.overlay.status="Please enter a name to continue.",e.overlay.isError=!0,o();return}e.overlay.busy=!0,e.overlay.status="Saving profile…",e.overlay.isError=!1,o();try{const r=await u().SaveUserProfile(t);e.bootstrap&&(e.bootstrap.profile=r,e.bootstrap.greeting=`Hi! ${r.DisplayName}`),e.overlay=null}catch(r){((a=e.overlay)==null?void 0:a.type)==="namePrompt"&&(e.overlay.busy=!1,e.overlay.status=c(r),e.overlay.isError=!0)}o()}async function J(){var s,a;if(((s=e.overlay)==null?void 0:s.type)!=="quoteEditor"||e.overlay.busy)return;const t=e.overlay.content.trim();if(!t){e.overlay.status="Nothing to save.",e.overlay.isError=!0,o();return}e.overlay.busy=!0,e.overlay.status="Refining draft...",e.overlay.isError=!1,o();try{const r=e.overlay.mode==="add"?await u().AddQuote(t):await u().UpdateQuote(e.overlay.quoteId??0,t);e.overlay=null,x(r),await f()}catch(r){((a=e.overlay)==null?void 0:a.type)==="quoteEditor"&&(e.overlay.busy=!1,e.overlay.status=c(r),e.overlay.isError=!0),o()}}async function Y(){var s,a,r;if(((s=e.overlay)==null?void 0:s.type)!=="quoteEditor"||e.overlay.busy)return;const t=e.overlay.content.trim();if(!t){e.overlay.status="Nothing to refine.",e.overlay.isError=!0,o();return}e.overlay.busy=!0,e.overlay.status="",o();try{const i=await u().RefineQuoteDraft(t);if(((a=e.overlay)==null?void 0:a.type)!=="quoteEditor")return;e.overlay.busy=!1,e.overlay.previewOriginal=t,e.overlay.previewRefined=i,e.overlay.status="",e.overlay.isError=!1}catch(i){((r=e.overlay)==null?void 0:r.type)==="quoteEditor"&&(e.overlay.busy=!1,e.overlay.status=c(i),e.overlay.isError=!0)}o()}function Pe(){var t;((t=e.overlay)==null?void 0:t.type)==="quoteEditor"&&(e.overlay.content=e.overlay.previewRefined,e.overlay.previewOriginal="",e.overlay.previewRefined="",e.overlay.status="Refined draft applied. Review it, then save.",e.overlay.isError=!1,o())}function Le(){var t;((t=e.overlay)==null?void 0:t.type)==="quoteEditor"&&(e.overlay.previewOriginal="",e.overlay.previewRefined="",e.overlay.status="Refined draft discarded.",e.overlay.isError=!1,o())}async function Ce(){var t,s,a,r;if(((t=e.overlay)==null?void 0:t.type)==="deleteHistory"){if(e.overlay.busy)return;e.overlay.busy=!0,e.overlay.status="",o();try{await u().DeleteRecallHistory(e.overlay.ids),Ge(e.overlay.ids),e.overlay=null,await P()}catch(i){((s=e.overlay)==null?void 0:s.type)==="deleteHistory"&&(e.overlay.busy=!1,e.overlay.status=c(i),e.overlay.isError=!0),o()}return}if(!(((a=e.overlay)==null?void 0:a.type)!=="deleteQuotes"||e.overlay.busy)){e.overlay.busy=!0,e.overlay.status="",o();try{await u().DeleteQuotes(e.overlay.ids),Ue(e.overlay.ids),e.overlay=null,await f()}catch(i){((r=e.overlay)==null?void 0:r.type)==="deleteQuotes"&&(e.overlay.busy=!1,e.overlay.status=c(i),e.overlay.isError=!0),o()}}}async function xe(){var t,s,a;if(!(((t=e.overlay)==null?void 0:t.type)!=="shareQuotes"||e.overlay.busy)){if(v()){e.overlay.path="irecall-share.json",o();return}try{const r=await u().SelectQuoteExportFile();r&&((s=e.overlay)==null?void 0:s.type)==="shareQuotes"&&(e.overlay.path=r,o())}catch(r){((a=e.overlay)==null?void 0:a.type)==="shareQuotes"&&(e.overlay.status=c(r),e.overlay.isError=!0,o())}}}async function Re(){var s,a,r;if(((s=e.overlay)==null?void 0:s.type)!=="shareQuotes"||e.overlay.busy)return;if(v()){const i=e.overlay.path.trim()||"irecall-share.json";if(!e.overlay.payload.trim()){e.overlay.status="Export payload is not ready yet.",e.overlay.isError=!0,o();return}mt(i,e.overlay.payload),e.overlay.status=`Downloaded ${i}`,e.overlay.isError=!1,o();return}const t=e.overlay.path.trim();if(!t){e.overlay.status="Choose a file path for the export.",e.overlay.isError=!0,o();return}if(!e.overlay.payload.trim()){e.overlay.status="Export payload is not ready yet.",e.overlay.isError=!0,o();return}e.overlay.busy=!0,e.overlay.status="",o();try{await u().ExportQuotesToFile(e.overlay.ids,t),((a=e.overlay)==null?void 0:a.type)==="shareQuotes"&&(e.overlay.busy=!1,e.overlay.status=`Saved share payload to ${t}`,e.overlay.isError=!1,o())}catch(i){((r=e.overlay)==null?void 0:r.type)==="shareQuotes"&&(e.overlay.busy=!1,e.overlay.status=c(i),e.overlay.isError=!0,o())}}async function Te(){var t,s,a;if(!(((t=e.overlay)==null?void 0:t.type)!=="importQuotes"||e.overlay.busy)){if(v()){const r=document.querySelector('[data-bind="import-file"]');r==null||r.click();return}try{const r=await u().SelectQuoteImportFile();r&&((s=e.overlay)==null?void 0:s.type)==="importQuotes"&&(e.overlay.path=r,o())}catch(r){((a=e.overlay)==null?void 0:a.type)==="importQuotes"&&(e.overlay.status=c(r),e.overlay.isError=!0,o())}}}async function Me(){var a,r,i;if(((a=e.overlay)==null?void 0:a.type)!=="importQuotes"||e.overlay.busy)return;const t=e.overlay.path.trim(),s=e.overlay.payload.trim();if(v()){if(!s){e.overlay.status="Choose a file to import.",e.overlay.isError=!0,o();return}}else if(!t){e.overlay.status="Choose a file to import.",e.overlay.isError=!0,o();return}e.overlay.busy=!0,e.overlay.status="",e.overlay.result=null,o();try{const l=v()?await u().ImportQuotesPayload(s):await u().ImportQuotesFromFile(t);if(((r=e.overlay)==null?void 0:r.type)!=="importQuotes")return;e.overlay.busy=!1,e.overlay.result=l,e.overlay.status=`Imported quotes. inserted=${l.Inserted} updated=${l.Updated} duplicates=${l.Duplicates} stale=${l.Stale}`,e.overlay.isError=!1,await f()}catch(l){((i=e.overlay)==null?void 0:i.type)==="importQuotes"&&(e.overlay.busy=!1,e.overlay.status=c(l),e.overlay.isError=!0,o())}}async function L(){if(e.recallBusy)return;const t=e.recallQuestion.trim();if(!t){e.recallError="Enter a recall question first.",o();return}e.recallBusy=!0,e.recallError="",e.recallStatus="",e.recallStatusIsError=!1,e.recallLastQuestion=t,e.recallKeywords=[],e.recallQuotes=[],e.recallResponse="",e.recallCursor=0,e.recallSelected=new Set,o();try{const s=await u().RunRecall(t);e.recallKeywords=s.keywords,e.recallQuotes=s.quotes,e.recallResponse=s.response,e.recallLastQuestion=s.question||t,e.recallCursor=0,e.recallSelected=new Set,e.recallQuestion=""}catch(s){e.recallError=c(s)}finally{e.recallBusy=!1,o()}}async function He(){const t=e.recallLastQuestion.trim(),s=e.recallResponse.trim();if(!t||!s){e.recallStatus="Run a recall first before saving it as a quote.",e.recallStatusIsError=!0,o();return}try{const a=await u().SaveRecallAsQuote(t,s,e.recallKeywords);x(a),await f(),e.recallStatus="Saved recall as quote.",e.recallStatusIsError=!1,X("Saved the current grounded answer as a quote.")}catch(a){e.recallStatus=c(a),e.recallStatusIsError=!0}o()}async function Ae(){const t=e.historyDetail;if(t){try{const s=await u().SaveRecallAsQuote(t.Question,t.Response,[]);x(s),await f(),e.historyStatus="Saved history entry as quote.",e.historyStatusIsError=!1,X("Saved the selected activity session as a quote.")}catch(s){e.historyStatus=c(s),e.historyStatusIsError=!0}o()}}async function Fe(){if(e.settingsBusy)return;let t;try{t=ee(e.settings)}catch(s){e.settingsStatus=c(s),e.settingsIsError=!0,o();return}e.settingsBusy=!0,e.settingsStatus="",o();try{const s=await u().FetchModels(t);e.settings.models=s,M(e.settings),e.settingsStatus=s.length>0?`Fetched ${s.length} models.`:"No models returned.",e.settingsIsError=!1}catch(s){e.settingsStatus=c(s),e.settingsIsError=!0}finally{e.settingsBusy=!1,o()}}async function G(){var s;if(e.settingsBusy)return;let t;try{t=vt(e.settings)}catch(a){e.settingsStatus=c(a),e.settingsIsError=!0,o();return}e.settingsBusy=!0,e.settingsStatus="",o();try{const a=await u().SaveSettings(t);e.settings=_(a,e.settings.models),I(e.settings.theme),e.bootstrap&&(e.bootstrap.settings=a);const r=((s=e.auth)==null?void 0:s.runtime)==="web"&&e.auth.currentPort>0&&e.auth.currentPort!==a.Web.Port;e.settingsStatus=r?"Saved. Restart the web server to apply the new port.":"Saved.",e.settingsIsError=!1}catch(a){e.settingsStatus=c(a),e.settingsIsError=!0}finally{e.settingsBusy=!1,o()}}function V(){e.overlay&&e.overlay.type!=="namePrompt"&&("busy"in e.overlay&&e.overlay.busy||(e.overlay=null,o()))}function X(t,s=!1){e.toast={message:t,isError:s},g!==null&&window.clearTimeout(g),g=window.setTimeout(()=>{e.toast=null,g=null,o()},2600)}function E(t){var s;return t==="quotes"?Z():t==="recall"?e.recallQuotes:((s=e.historyDetail)==null?void 0:s.Quotes)??[]}function Ne(t,s){var r;const a=E(t);if(t==="quotes")e.quotesCursor=h(s,a);else if(t==="recall")e.recallCursor=h(s,e.recallQuotes);else{const i=((r=e.historyDetail)==null?void 0:r.Quotes)??[];e.historyQuoteCursor=h(s,i)}o()}function Be(t,s,a){const r=t==="quotes"?e.quotesSelected:t==="recall"?e.recallSelected:e.historyQuoteSelected;a?r.add(s):r.delete(s),o()}function Oe(t){const s=E(t),a=new Set(s.map(r=>r.ID));t==="quotes"?e.quotesSelected=a:t==="recall"?e.recallSelected=a:e.historyQuoteSelected=a,o()}function je(t){t==="quotes"?e.quotesSelected=new Set:t==="recall"?e.recallSelected=new Set:e.historyQuoteSelected=new Set,o()}function C(t){const s=E(t),a=t==="quotes"?e.quotesCursor:t==="recall"?e.recallCursor:e.historyQuoteCursor,r=h(a,s),i=t==="quotes"?e.quotesSelected:t==="recall"?e.recallSelected:e.historyQuoteSelected,l=s.filter(d=>i.has(d.ID));return l.length>0?l:s[r]?[s[r]]:[]}function Ke(t){var s;return((s=e.overlay)==null?void 0:s.type)==="quoteInspect"&&e.overlay.context===t?e.overlay.quote:C(t)[0]??null}function Z(){const t=e.libraryQuery.trim().toLowerCase();return e.quotes.filter(s=>t?[s.Content,s.AuthorName,s.SourceName,...s.Tags].join(" ").toLowerCase().includes(t):!0)}function x(t){e.quotes=k(e.quotes,t),e.recallQuotes=k(e.recallQuotes,t),e.historyDetail&&(e.historyDetail={...e.historyDetail,Quotes:k(e.historyDetail.Quotes,t)}),o()}function Ue(t){var a;const s=new Set(t);e.quotes=e.quotes.filter(r=>!s.has(r.ID)),e.recallQuotes=e.recallQuotes.filter(r=>!s.has(r.ID)),e.historyDetail&&(e.historyDetail={...e.historyDetail,Quotes:e.historyDetail.Quotes.filter(r=>!s.has(r.ID))}),e.quotesSelected=new Set([...e.quotesSelected].filter(r=>!s.has(r))),e.recallSelected=new Set([...e.recallSelected].filter(r=>!s.has(r))),e.historyQuoteSelected=new Set([...e.historyQuoteSelected].filter(r=>!s.has(r))),e.quotesCursor=h(e.quotesCursor,e.quotes),e.recallCursor=h(e.recallCursor,e.recallQuotes),e.historyQuoteCursor=h(e.historyQuoteCursor,((a=e.historyDetail)==null?void 0:a.Quotes)??[]),o()}function We(t){e.historyCursor=H(t,e.historyEntries);const s=e.historyEntries[e.historyCursor];s&&S(s.ID,!0),o()}function ze(t,s){s?e.historySelected.add(t):e.historySelected.delete(t),o()}function Je(){e.historySelected=new Set(e.historyEntries.map(t=>t.ID)),o()}function Ye(){e.historySelected=new Set,o()}function R(){const t=e.historyEntries.filter(s=>e.historySelected.has(s.ID));return t.length>0?t:e.historyEntries[e.historyCursor]?[e.historyEntries[e.historyCursor]]:[]}function Ge(t){const s=new Set(t);if(e.historyEntries=e.historyEntries.filter(a=>!s.has(a.ID)),e.historySelected=new Set([...e.historySelected].filter(a=>!s.has(a))),e.historyCursor=H(e.historyCursor,e.historyEntries),e.historyDetail&&s.has(e.historyDetail.ID)){$();return}o()}function o(){if(!b)return;const t=Ve();b.innerHTML=Ze(),Xe(t)}function Ve(){const t=document.activeElement;if(!(t instanceof HTMLInputElement||t instanceof HTMLTextAreaElement||t instanceof HTMLSelectElement))return null;const s=t.dataset.bind;return s?{selector:`[data-bind="${s}"]`,selectionStart:t instanceof HTMLInputElement||t instanceof HTMLTextAreaElement?t.selectionStart:null,selectionEnd:t instanceof HTMLInputElement||t instanceof HTMLTextAreaElement?t.selectionEnd:null}:null}function Xe(t){if(!b||!t)return;const s=b.querySelector(t.selector);(s instanceof HTMLInputElement||s instanceof HTMLTextAreaElement||s instanceof HTMLSelectElement)&&(s.focus({preventScroll:!0}),(s instanceof HTMLInputElement||s instanceof HTMLTextAreaElement)&&t.selectionStart!==null&&t.selectionEnd!==null&&s.setSelectionRange(t.selectionStart,t.selectionEnd))}function Ze(){var s,a,r,i,l;if(!e.authChecked)return`
      <div class="shell shell-loading">
        <div class="splash">
          <div class="brand">iRecall</div>
          <div class="muted">Checking workspace access…</div>
        </div>
      </div>
    `;if(e.fatalError)return`
      <div class="shell shell-loading">
        <div class="splash splash-error">
          <div class="brand">iRecall</div>
          <div class="status status-error">${n(e.fatalError)}</div>
        </div>
      </div>
    `;if(((s=e.auth)==null?void 0:s.runtime)==="web"&&!e.auth.authenticated)return _e();if(!e.bootstrapped)return`
      <div class="shell shell-loading">
        <div class="splash">
          <div class="brand">iRecall</div>
          <div class="muted">Loading workspace…</div>
        </div>
      </div>
    `;const t=(r=(a=e.bootstrap)==null?void 0:a.profile)!=null&&r.DisplayName?`Hi! ${e.bootstrap.profile.DisplayName}`:"";return`
    <div class="shell">
      <header class="titlebar">
        <div class="brand-lockup">
          <div class="brand-row">
            <div class="brand">${n(((i=e.bootstrap)==null?void 0:i.productName)??"iRecall")}</div>
            ${e.settings.mockLLM?'<span class="meta-pill meta-pill-accent">Mock LLM on</span>':""}
          </div>
          <div class="muted subtle">${v()?"Local-first knowledge workspace for the web":"Local-first knowledge workspace for desktop"}</div>
        </div>
        <div class="titlebar-right">
          <div class="greeting">${n(t)}</div>
          ${((l=e.auth)==null?void 0:l.runtime)==="web"?'<button class="button" data-action="auth-logout" type="button">Logout</button>':""}
          <nav class="tabs" aria-label="Primary">
            ${le.map(d=>`
                  <button
                    class="tab${e.page===d?" active":""}"
                    data-action="nav"
                    data-page="${d}"
                    type="button"
                  >${tt(d)}</button>
                `).join("")}
          </nav>
        </div>
      </header>

      <main class="layout">
        ${et()}
      </main>

      ${ot()}
      ${e.overlay?ct(e.overlay):""}
      ${e.toast?ut(e.toast):""}
    </div>
  `}function _e(){var i;const t=!((i=e.auth)!=null&&i.passwordConfigured),s="auth-login";return`
    <div class="shell shell-loading">
      <div class="panel modal">
        <div class="brand">iRecall</div>
        <div class="modal-title">${t?"Password Required In Terminal":"Unlock Web UI"}</div>
        <div class="modal-copy">${t?"The web password must be created in the terminal before the server starts listening. Restart the server from a terminal session to finish setup.":"Enter the web password to unlock the shared iRecall database."}</div>
        <form class="modal-form" data-form="${s}">
          <label class="field">
            <span>Password</span>
            <input class="text-input" data-bind="auth-password" type="password" value="${p(e.authPassword)}" ${t?"disabled":""} />
          </label>
          ${e.authStatus?`<div class="status ${e.authIsError?"status-error":"status-ok"}">${n(e.authStatus)}</div>`:""}
          <div class="modal-actions">
            <button class="button button-primary" data-action="${s}" type="submit" ${e.authBusy||t?"disabled":""}>
              ${e.authBusy?"Working…":"Login"}
            </button>
          </div>
        </form>
      </div>
    </div>
  `}function et(){switch(e.page){case"Recall":return st();case"Quotes":return at();case"History":return rt();case"Settings":return nt()}}function tt(t){switch(t){case"Recall":return"Recall";case"Quotes":return"Quotes";case"History":return"History";case"Settings":return"Settings"}}function st(){const t=!e.recallResponse.trim(),s=e.settings.mockLLM,a=e.recallResponse.trim()?n(e.recallResponse):'<span class="muted">Grounded response will appear here.</span>',r=e.recallKeywords.length>0?e.recallKeywords.map(i=>`<span class="keyword-chip">${n(i)}</span>`).join(""):'<span class="muted">Keywords: —</span>';return`
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
              >${n(e.recallQuestion)}</textarea>
              <div class="composer-actions">
                <button class="button button-primary" data-action="recall-run" type="submit" ${e.recallBusy?"disabled":""}>
                  ${e.recallBusy?"Working…":"Recall"}
                </button>
                ${e.recallLastQuestion.trim()?'<button class="button" data-action="use-last-question" type="button">Use previous question</button>':""}
              </div>
            </form>
          </section>

          <section class="panel subpanel">
            <div class="subpanel-header">
              <div>
                <div class="section-title">2. Reference quotes</div>
                <div class="muted">${e.recallBusy?"Searching your quotes for relevant evidence…":s?`${e.recallQuotes.length} retrieved quotes. Mock LLM uses simple split keywords and deterministic recall behavior.`:`${e.recallQuotes.length} retrieved quotes. Open one to inspect the full note.`}</div>
              </div>
            </div>
            <div class="keyword-row">
              <span class="muted">Keywords</span>
              <div class="keyword-list">${r}</div>
            </div>
            ${T("recall",e.recallQuotes,e.recallCursor,e.recallSelected,!1)}
          </section>

          <section class="panel subpanel">
            <div class="subpanel-header">
              <div>
                <div class="section-title">3. Response</div>
                <div class="muted">${e.recallBusy?"Writing a grounded response from the retrieved evidence…":s?"Mock LLM combines the retrieved reference quotes into a deterministic placeholder answer.":"The response is generated from the current question and reference set."}</div>
              </div>
              <div class="toolbar toolbar-quiet">
                <button class="button button-primary" data-action="recall-save-quote" type="button" ${t?"disabled":""}>Save as Quote</button>
                <button class="button" data-action="nav" data-page="History" type="button" ${t?"disabled":""}>Open history</button>
              </div>
            </div>
            <div class="answer-card">
              ${e.recallLastQuestion.trim()?`
                    <div class="answer-anchor">
                      <div class="muted">Current question</div>
                      <div class="answer-question">${n(e.recallLastQuestion)}</div>
                    </div>
                  `:""}
              <pre class="response-box">${a}</pre>
            </div>
          </section>
        </div>

        ${e.recallError?`<div class="status status-error">${n(e.recallError)}</div>`:""}
        ${e.recallStatus?`<div class="status ${e.recallStatusIsError?"status-error":"status-ok"}">${n(e.recallStatus)}</div>`:""}
      </div>
    </section>
  `}function at(){const t=Z(),s=h(e.quotesCursor,t),a=t.filter(r=>e.quotesSelected.has(r.ID)).length;return`
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
                <div class="muted">${t.length} ${t.length===1?"quote":"quotes"}</div>
              </div>
              <div class="toolbar toolbar-quiet">
                <button class="button" data-action="quote-select-all" data-context="quotes" type="button" ${t.length===0?"disabled":""}>Select all</button>
                <button class="button" data-action="quote-deselect-all" data-context="quotes" type="button" ${a===0?"disabled":""}>Clear</button>
                <button class="button" data-action="quote-share-current" data-context="quotes" type="button" ${a===0?"disabled":""}>Share</button>
              </div>
            </div>
            ${e.quotesLoading?'<div class="empty-state">Loading quotes…</div>':e.quotesError?`<div class="status status-error">${n(e.quotesError)}</div>`:T("quotes",t,s,e.quotesSelected,!0)}
          </section>
        </div>
      </div>
    </section>
  `}function rt(){return`
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
                <div class="muted">${e.historyEntries.length} saved recall sessions</div>
              </div>
              <div class="toolbar toolbar-quiet">
                <button class="button" data-action="history-select-all" type="button" ${e.historyEntries.length===0?"disabled":""}>Select all</button>
                <button class="button" data-action="history-deselect-all" type="button" ${e.historySelected.size===0?"disabled":""}>Clear</button>
              </div>
            </div>
            ${e.historyLoading?'<div class="empty-state">Loading history…</div>':e.historyError?`<div class="status status-error">${n(e.historyError)}</div>`:e.historyEntries.length===0?'<div class="empty-state">No recall history yet. Run a question from the Recall page to create your first grounded session.</div>':it()}
          </section>
        </div>
      </div>
    </section>
  `}function ot(){if(!e.historyDetailLoading&&!e.historyDetail&&!e.historyDetailError)return"";const t=e.historyDetail?e.historyEntries.find(r=>{var i;return r.ID===((i=e.historyDetail)==null?void 0:i.ID)})??e.historyDetail:e.historyEntries[e.historyCursor]??null,s=e.historyDetail&&t&&e.historyDetail.ID===t.ID?e.historyDetail:null,a=(s==null?void 0:s.Quotes)??[];return`
    <div class="overlay-backdrop">
      <div class="modal modal-history-detail">
        <div class="subpanel-header">
          <div>
            <div class="modal-title">Session details</div>
            <div class="muted">${t?s?`${a.length} reference quotes loaded`:"Opening the selected session…":"Loading the selected session…"}</div>
          </div>
          <div class="toolbar toolbar-quiet">
            <button class="button" data-action="history-back" type="button">Close</button>
            <button class="button button-primary" data-action="history-save-quote" type="button" ${t?"":"disabled"}>Save as Quote</button>
            <button class="button" data-action="reuse-history-question" type="button" ${t?"":"disabled"}>Recall again</button>
          </div>
        </div>

        ${t?`
              <div class="detail-stack">
                <div class="detail-block">
                  <div class="muted">Question</div>
                  <pre class="response-box compact-box">${n(t.Question)}</pre>
                </div>
                <div class="detail-block">
                  <div class="muted">Response</div>
                  <pre class="response-box compact-box">${n((s==null?void 0:s.Response)??t.Response)}</pre>
                </div>
              </div>
              ${e.historyDetailLoading?'<div class="empty-state">Loading reference quotes…</div>':s?`
                      <div class="subpanel-header nested-header">
                        <div>
                          <div class="section-title">Reference quotes</div>
                          <div class="muted">${a.length} retrieved quotes. Open one to inspect the full note.</div>
                        </div>
                      </div>
                      ${T("history",a,e.historyQuoteCursor,e.historyQuoteSelected,!1)}
                    `:""}
            `:""}

        ${e.historyDetailError?`<div class="status status-error">${n(e.historyDetailError)}</div>`:""}
        ${e.historyStatus?`<div class="status ${e.historyStatusIsError?"status-error":"status-ok"}">${n(e.historyStatus)}</div>`:""}
      </div>
    </div>
  `}function it(){return`
    <div class="history-list">
      ${e.historyEntries.map((t,s)=>{const a=s===e.historyCursor,r=D(t.Response,156);return`
            <article class="quote-card history-card${a?" is-current":""}" data-action="history-set-cursor" data-index="${s}">
              <div class="quote-topline">
                <label class="selection-toggle">
                  <input
                    type="checkbox"
                    data-bind="history-selected"
                    data-id="${t.ID}"
                    ${e.historySelected.has(t.ID)?"checked":""}
                  />
                </label>
                <div class="quote-topline-meta">
                  <span class="quote-version">${n(ft(t.CreatedAt))}</span>
                </div>
              </div>
              <div class="quote-content">${n(D(t.Question,132))}</div>
                <div class="quote-meta"><span class="muted">Response preview</span><span>${n(r||"(empty response)")}</span></div>
            </article>
          `}).join("")}
    </div>
  `}function nt(){var d,m;const t=te(e.settings),s=(d=e.bootstrap)==null?void 0:d.paths,a=(m=e.auth)==null?void 0:m.currentPort,r=e.apiToken.hasToken?"Renew API Token":"Create API Token",i=e.apiToken.loading?"Loading token status…":e.apiToken.hasToken?`Active token prefix: ${e.apiToken.tokenPrefix||"(unavailable)"}`:"No API token has been created yet.",l=e.settings.models.length>0&&t.length>0?`
        <select class="select-input" data-bind="settings-model">
          ${t.map(y=>`
                <option value="${p(y)}"${y===e.settings.model?" selected":""}>${n(y)}</option>
              `).join("")}
        </select>
      `:`
        <div class="readonly-model">
          <span>${n(e.settings.model||"(none)")}</span>
          <span class="muted">${e.settings.models.length===0?"Fetch models first":"No matches"}</span>
        </div>
      `;return`
    <section class="page page-settings">
      <div class="panel page-panel">
        <div class="page-hero">
          <div>
            <div class="eyebrow">Settings</div>
            <div class="page-title">Configure connection, retrieval, and local preferences</div>
            <div class="muted page-copy">Keep the primary setup visible and move lower-level runtime controls into advanced sections.</div>
          </div>
          <div class="page-hero-actions">
            <button class="button" data-action="settings-fetch-models" type="button" ${e.settingsBusy?"disabled":""}>
              ${e.settingsBusy?"Fetching…":"Fetch Models"}
            </button>
            <button class="button button-primary" data-action="settings-save" type="button" ${e.settingsBusy?"disabled":""}>Save Changes</button>
          </div>
        </div>

        <div class="settings-layout">
          <section class="panel subpanel settings-primary">
            <div class="section-title">Connection</div>
            <label class="field">
              <span>Host / IP</span>
              <input class="text-input" data-bind="settings-host" value="${p(e.settings.host)}" />
            </label>
            <div class="settings-row">
              <label class="field">
                <span>Port</span>
                <input class="text-input" data-bind="settings-port" value="${p(e.settings.port)}" />
              </label>
              <label class="field checkbox-field settings-toggle">
                <input type="checkbox" data-bind="settings-https"${e.settings.https?" checked":""} />
                <span>Use HTTPS</span>
              </label>
            </div>
            <label class="field">
              <span>API Key</span>
              <input class="text-input" data-bind="settings-api-key" type="password" value="${p(e.settings.apiKey)}" />
            </label>
            <label class="field">
              <span>Filter models</span>
              <input class="text-input" data-bind="settings-model-filter" value="${p(e.settings.modelFilter)}" placeholder="Type to narrow the model list" />
            </label>
            <label class="field">
              <span>Model</span>
              ${l}
            </label>
          </section>

          <section class="panel subpanel settings-primary">
            <div class="section-title">Retrieval</div>
            <label class="field">
              <span>Max reference quotes</span>
              <input class="text-input" data-bind="settings-max-results" value="${p(e.settings.maxResults)}" />
            </label>
            <label class="field">
              <span>Minimum relevance</span>
              <input class="text-input" data-bind="settings-min-relevance" value="${p(e.settings.minRelevance)}" placeholder="0.0-1.0" />
            </label>
            <div class="settings-hint muted">
              Lower values keep broader matches. Higher values reduce noise but risk excluding useful evidence. Most real-world sessions land in the 0.3 to 0.7 range.
            </div>
          </section>

          <section class="panel subpanel settings-secondary">
            <div class="section-title">Debug</div>
            <label class="field checkbox-field settings-toggle">
              <input type="checkbox" data-bind="settings-mock-llm"${e.settings.mockLLM?" checked":""} />
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
                ${ne().map(y=>`
                      <option value="${y}"${y===e.settings.theme?" selected":""}>${y}</option>
                    `).join("")}
              </select>
            </label>
          </section>

          <section class="panel subpanel settings-secondary">
            <div class="section-title">Security</div>
            <label class="field">
              <span>Current Password</span>
              <input class="text-input" data-bind="settings-password-current" type="password" value="${p(e.passwordForm.current)}" />
            </label>
            <label class="field">
              <span>New Password</span>
              <input class="text-input" data-bind="settings-password-next" type="password" value="${p(e.passwordForm.next)}" />
            </label>
            <label class="field">
              <span>Confirm Password</span>
              <input class="text-input" data-bind="settings-password-confirm" type="password" value="${p(e.passwordForm.confirm)}" />
            </label>
            <div class="muted subtle">Use at least 12 characters and include at least 3 of: uppercase, lowercase, digit, symbol.</div>
            <div class="toolbar">
              <button class="button" data-action="settings-change-password" type="button" ${e.passwordForm.busy?"disabled":""}>
                ${e.passwordForm.busy?"Updating…":"Change Password"}
              </button>
            </div>
            ${e.passwordForm.status?`<div class="status ${e.passwordForm.isError?"status-error":"status-ok"}">${n(e.passwordForm.status)}</div>`:""}
          </section>

          <section class="panel subpanel settings-secondary">
            <div class="section-title">REST API Token</div>
            <div class="settings-hint muted">
              Use this token for external REST clients with <code>Authorization: Bearer &lt;token&gt;</code>. The plaintext token is shown only once after creation or renewal.
            </div>
            <div class="readonly-model path-value">${n(i)}</div>
            <div class="toolbar">
              <button class="button" data-action="settings-create-api-token" type="button" ${e.settingsBusy||e.apiToken.loading?"disabled":""}>
                ${e.settingsBusy?"Working…":r}
              </button>
            </div>
          </section>

          <section class="panel subpanel settings-secondary">
            <div class="section-title">Advanced</div>
            <label class="field">
              <span>Web Port</span>
              <input class="text-input" data-bind="settings-web-port" value="${p(e.settings.webPort)}" />
            </label>
            <div class="settings-hint muted">
              The web server listens on this port after restart. Current listener: ${n(a?String(a):"(not running)")}.
            </div>
            <div class="settings-paths">
              <div class="field">
                <span>Data dir</span>
                <div class="readonly-model path-value">${n((s==null?void 0:s.dataDir)??"(unavailable)")}</div>
              </div>
              <div class="field">
                <span>Config dir</span>
                <div class="readonly-model path-value">${n((s==null?void 0:s.configDir)??"(unavailable)")}</div>
              </div>
              <div class="field">
                <span>State dir</span>
                <div class="readonly-model path-value">${n((s==null?void 0:s.stateDir)??"(unavailable)")}</div>
              </div>
              <div class="field">
                <span>Database</span>
                <div class="readonly-model path-value">${n((s==null?void 0:s.dbPath)??"(unavailable)")}</div>
              </div>
            </div>
          </section>
        </div>

        ${e.settingsStatus?`<div class="status ${e.settingsIsError?"status-error":"status-ok"}">${n(e.settingsStatus)}</div>`:""}
      </div>
    </section>
  `}function T(t,s,a,r,i){return s.length===0?`<div class="empty-state">${t==="quotes"?"No quotes yet. Add one or import a shared payload.":"No reference quotes for this question yet."}</div>`:`
    <div class="quote-list">
      ${s.map((l,d)=>{const m=d===a,y=t!=="quotes",A=!l.IsOwnedByMe&&l.SourceName?`<span class="meta-accent">${n(l.SourceName)}</span>`:`<span>${n(l.AuthorName||"You")}</span>`,ie=i?`
              <div class="quote-meta">
                <span class="muted">Tags</span>
                <span>${l.Tags.length>0?n(bt(l.Tags,4)):"(none)"}</span>
              </div>
            `:"";return`
            <article class="quote-card${m?" is-current":""}${y?" quote-card-minimal":""}" data-action="quote-inspect" data-context="${t}" data-index="${d}">
              <div class="quote-topline">
                ${y?`<div class="quote-topline-meta">
                        <span class="quote-badge">${l.IsOwnedByMe?"Owned":"Imported"}</span>
                        <span class="quote-version">${n(Q(l.UpdatedAt))}</span>
                      </div>`:`
                      <label class="selection-toggle">
                        <input
                          type="checkbox"
                          data-bind="quote-selected"
                          data-context="${t}"
                          data-id="${l.ID}"
                          ${r.has(l.ID)?"checked":""}
                        />
                  </label>
                    `}
                <div class="quote-topline-meta">
                  ${y?`<span class="quote-source-inline">${A}</span>`:`<span class="quote-version">${n(Q(l.UpdatedAt))}</span>
                  <span class="quote-badge">${l.IsOwnedByMe?"Owned":"Imported"}</span>`}
                </div>
              </div>
              <div class="quote-content">${n(D(l.Content,t==="quotes"?160:136))}</div>
              ${y?`<div class="quote-actions-inline"><button class="button button-subtle" data-action="quote-inspect" data-context="${t}" data-index="${d}" type="button">Details</button></div>`:`<div class="quote-meta"><span class="muted">${!l.IsOwnedByMe&&l.SourceName?"Imported from":"Author"}</span> ${A}</div>`}
              ${ie}
            </article>
          `}).join("")}
    </div>
  `}function lt(t,s){return t?`
    <div class="detail-stack">
      <div class="detail-block">
        <div class="muted">Full quote</div>
        <pre class="response-box compact-box">${n(t.Content)}</pre>
      </div>
      <div class="detail-grid">
        <div class="detail-metric">
          <span class="muted">Author</span>
          <span>${n(t.AuthorName||"You")}</span>
        </div>
        <div class="detail-metric">
          <span class="muted">Version</span>
          <span>v${t.Version}</span>
        </div>
        <div class="detail-metric">
          <span class="muted">Source</span>
          <span>${n(t.SourceName||"Local library")}</span>
        </div>
        <div class="detail-metric">
          <span class="muted">Updated</span>
          <span>${n(Q(t.UpdatedAt))}</span>
        </div>
      </div>
      <div class="detail-block">
        <div class="muted">Tags</div>
        <div class="keyword-list">
          ${t.Tags.length>0?t.Tags.map(a=>`<span class="keyword-chip">${n(a)}</span>`).join(""):'<span class="muted">No tags assigned yet.</span>'}
        </div>
      </div>
      <div class="toolbar toolbar-inline">
        <button class="button" data-action="quote-edit-current" data-context="${s}" type="button">Edit</button>
        <button class="button" data-action="quote-share-current" data-context="${s}" type="button">Share</button>
        <button class="button button-danger" data-action="quote-delete-current" data-context="${s}" type="button">Delete</button>
      </div>
    </div>
  `:'<div class="empty-state">Select a quote to inspect the full note, provenance, and available actions.</div>'}function Q(t){const s=new Date(t);return Number.isNaN(s.getTime())?t:s.toLocaleDateString(void 0,{month:"short",day:"numeric",year:"numeric"})}function B(t){return!Number.isFinite(t)||t<=0?"0 B":t<1024?`${Math.round(t)} B`:t<1024*1024?`${(t/1024).toFixed(1)} KB`:`${(t/(1024*1024)).toFixed(1)} MB`}function ut(t){return`
    <div class="toast-stack" role="status" aria-live="polite">
      <div class="toast${t.isError?" is-error":""}">${n(t.message)}</div>
    </div>
  `}function ct(t){switch(t.type){case"namePrompt":return`
        <div class="overlay-backdrop">
          <div class="modal">
            <div class="modal-title">Set Your Name</div>
            <p class="modal-copy">
              Your name is attached to quotes you share and shown when other users receive your quotes.
            </p>
            <form class="modal-form" data-form="profile">
              <label class="field">
                <span>Display Name</span>
                <input class="text-input text-input-lg" data-bind="profile-name" value="${p(t.name)}" placeholder="Your name" />
              </label>
              ${t.status?`<div class="status ${t.isError?"status-error":"status-ok"}">${n(t.status)}</div>`:""}
              <div class="modal-actions">
                <button class="button button-primary" data-action="profile-save" type="submit" ${t.busy?"disabled":""}>
                  ${t.busy?"Saving…":"Save name and continue"}
                </button>
              </div>
            </form>
          </div>
        </div>
      `;case"quoteEditor":return`
        <div class="overlay-backdrop overlay-backdrop-side">
          <div class="modal modal-side">
            <div class="modal-title">${t.mode==="add"?"Add Quote":"Edit Quote"}</div>
            ${t.previewRefined?`
                  <div class="compare-grid">
                    <section class="panel compare-panel">
                      <div class="section-title">Current Draft</div>
                      <pre class="compare-body">${n(t.previewOriginal)}</pre>
                    </section>
                    <section class="panel compare-panel">
                      <div class="section-title">Refined Draft</div>
                      <pre class="compare-body">${n(t.previewRefined)}</pre>
                    </section>
                  </div>
                `:`
                  <label class="field">
                    <span>Quote Content</span>
                    <textarea class="text-area" data-bind="quote-editor-content" rows="10" placeholder="Type or paste your note here.">${n(t.content)}</textarea>
                  </label>
                `}
            <div class="muted modal-copy">
              ${t.previewRefined?"Compare the current draft with the suggested rewrite before applying it.":e.settings.mockLLM?"Mock LLM is enabled, so Refine returns the original text and skips provider-dependent rewriting.":"Tags are regenerated automatically by the shared core logic."}
            </div>
            ${t.status?`<div class="status ${t.isError?"status-error":"status-ok"}">${n(t.status)}</div>`:""}
            <div class="modal-actions">
              ${t.previewRefined?`
                    <button class="button button-primary" data-action="quote-editor-apply-refined" type="button">Apply refined draft</button>
                    <button class="button" data-action="quote-editor-reject-refined" type="button">Keep editing current draft</button>
                  `:`
                    <button class="button button-primary" data-action="quote-editor-save" type="button" ${t.busy?"disabled":""}>
                      ${t.busy?"Saving…":"Save"}
                    </button>
                    <button class="button" data-action="quote-editor-refine" type="button" ${t.busy?"disabled":""}>
                      ${t.busy?"Working…":"Refine"}
                    </button>
                    <button class="button" data-action="overlay-close" type="button" ${t.busy?"disabled":""}>Cancel</button>
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
              ${O(t.context,t.ids).map((a,r)=>`<div class="summary-item">[${r+1}] ${n(w(a.Content,140))}</div>`).join("")}
            </div>
            ${t.status?`<div class="status ${t.isError?"status-error":"status-ok"}">${n(t.status)}</div>`:""}
            <div class="modal-actions">
              <button class="button button-danger" data-action="delete-confirm" type="button" ${t.busy?"disabled":""}>
                ${t.busy?"Deleting…":"Delete"}
              </button>
              <button class="button" data-action="overlay-close" type="button" ${t.busy?"disabled":""}>Cancel</button>
            </div>
          </div>
        </div>
      `;case"deleteHistory":return`
        <div class="overlay-backdrop">
          <div class="modal">
            <div class="modal-title modal-title-danger">Delete History</div>
            <div class="modal-copy">This permanently removes the selected recall history entries from the local library.</div>
            <div class="summary-list">
              ${dt(t.ids).map((a,r)=>`<div class="summary-item">[${r+1}] ${n(w(a.Question,140))}</div>`).join("")}
            </div>
            ${t.status?`<div class="status ${t.isError?"status-error":"status-ok"}">${n(t.status)}</div>`:""}
            <div class="modal-actions">
              <button class="button button-danger" data-action="delete-confirm" type="button" ${t.busy?"disabled":""}>
                ${t.busy?"Deleting…":"Delete"}
              </button>
              <button class="button" data-action="overlay-close" type="button" ${t.busy?"disabled":""}>Cancel</button>
            </div>
          </div>
        </div>
      `;case"shareQuotes":const s=O(t.context,t.ids);return`
        <div class="overlay-backdrop overlay-backdrop-side">
          <div class="modal modal-side">
            <div class="modal-title">Share Quotes</div>
            <div class="modal-copy">Export a portable share file. The file summary comes first; raw JSON is available only if you need to inspect it.</div>
            <div class="summary-list">
              ${s.map((a,r)=>`<div class="summary-item">[${r+1}] v${a.Version} ${n(w(a.Content,120))}</div>`).join("")}
            </div>
            <div class="result-grid">
              <div><span class="muted">Quotes:</span> ${s.length}</div>
              <div><span class="muted">Payload size:</span> ${B(t.payload.length)}</div>
            </div>
            <label class="field">
              <span>${v()?"Download As":"Save To"}</span>
              ${v()?`<input class="text-input" data-bind="share-path" value="${p(t.path||"irecall-share.json")}" placeholder="irecall-share.json" />`:`
                    <div class="path-row">
                      <input class="text-input" data-bind="share-path" value="${p(t.path)}" placeholder="/path/to/irecall-share.json" />
                      <button class="button" data-action="share-browse" type="button" ${t.busy?"disabled":""}>Browse</button>
                    </div>
                  `}
            </label>
            <div class="muted modal-copy">${v()?"Download the JSON payload locally, then transfer it manually to the recipient.":"Export to a JSON file and transfer it manually to the recipient."}</div>
            <div class="toolbar toolbar-inline">
              <button class="button" data-action="share-toggle-payload" type="button" ${t.payload?"":"disabled"}>
                ${t.showPayload?"Hide raw JSON":"Show raw JSON"}
              </button>
            </div>
            ${t.showPayload?`<div class="payload-box"><pre>${n(t.payload||"Preparing export payload…")}</pre></div>`:""}
            ${t.status?`<div class="status ${t.isError?"status-error":"status-ok"}">${n(t.status)}</div>`:""}
            <div class="modal-actions">
              <button class="button button-primary" data-action="share-save" type="button" ${t.busy?"disabled":""}>
                ${t.busy?"Working…":v()?"Download export file":"Save export file"}
              </button>
              <button class="button" data-action="overlay-close" type="button" ${t.busy?"disabled":""}>Close</button>
            </div>
          </div>
        </div>
      `;case"importQuotes":return`
        <div class="overlay-backdrop overlay-backdrop-side">
          <div class="modal modal-side">
            <div class="modal-title">Import Quotes</div>
            <div class="modal-copy">Import a quote share file exported from another iRecall instance. Start by choosing a file, then review the result summary.</div>
            <label class="field">
              <span>Import From</span>
              ${v()?`
                    <input class="text-input" data-bind="import-path" value="${p(t.path)}" placeholder="Choose a local JSON file" readonly />
                    <input data-bind="import-file" type="file" accept="application/json,.json" hidden />
                    <div class="toolbar">
                      <button class="button" data-action="import-browse" type="button" ${t.busy?"disabled":""}>Choose File</button>
                    </div>
                  `:`
                    <div class="path-row">
                      <input class="text-input" data-bind="import-path" value="${p(t.path)}" placeholder="/path/to/irecall-share.json" />
                      <button class="button" data-action="import-browse" type="button" ${t.busy?"disabled":""}>Browse</button>
                    </div>
                `}
            </label>
            ${t.filename||t.path?`
                  <div class="result-grid">
                    <div><span class="muted">File:</span> ${n(t.filename||t.path)}</div>
                    <div><span class="muted">Payload size:</span> ${B(t.payload.length)}</div>
                  </div>
                `:""}
            ${t.payload?`
                  <div class="toolbar toolbar-inline">
                    <button class="button" data-action="import-toggle-payload" type="button">
                      ${t.showPayload?"Hide raw JSON":"Show raw JSON"}
                    </button>
                  </div>
                  ${t.showPayload?`<div class="payload-box"><pre>${n(t.payload)}</pre></div>`:""}
                `:""}
            ${t.result?`
                  <div class="result-grid">
                    <div><span class="muted">Inserted:</span> ${t.result.Inserted}</div>
                    <div><span class="muted">Updated:</span> ${t.result.Updated}</div>
                    <div><span class="muted">Duplicates:</span> ${t.result.Duplicates}</div>
                    <div><span class="muted">Stale:</span> ${t.result.Stale}</div>
                  </div>
                `:""}
            ${t.status?`<div class="status ${t.isError?"status-error":"status-ok"}">${n(t.status)}</div>`:""}
            <div class="modal-actions">
              <button class="button button-primary" data-action="import-run" type="button" ${t.busy?"disabled":""}>
                ${t.busy?"Importing…":"Import file"}
              </button>
              <button class="button" data-action="overlay-close" type="button" ${t.busy?"disabled":""}>Close</button>
            </div>
          </div>
        </div>
      `;case"quoteInspect":return`
        <div class="overlay-backdrop">
          <div class="modal modal-quote-inspect">
            <div class="modal-title">Quote details</div>
            <p class="modal-copy">
              Review the full note and its provenance without leaving the current flow.
            </p>
            ${lt(t.quote,t.context)}
            <div class="modal-actions">
              <button class="button" data-action="overlay-close" type="button">Close</button>
            </div>
          </div>
        </div>
      `;case"apiTokenReveal":return`
        <div class="overlay-backdrop">
          <div class="modal modal-side">
            <div class="modal-title">Copy API Token Now</div>
            <p class="modal-copy">
              This token is shown only once. Copy it now and store it safely. After you close this dialog, only the prefix
              <strong> ${n(t.tokenPrefix)} </strong>
              will remain visible in Settings.
            </p>
            <div class="payload-box"><pre>${n(t.token)}</pre></div>
            <div class="muted modal-copy">Use it with <code>Authorization: Bearer &lt;token&gt;</code> on REST API requests.</div>
            <div class="modal-actions">
              <button class="button" data-action="overlay-close" type="button">Close</button>
            </div>
          </div>
        </div>
      `;case"notice":return""}}function O(t,s){var i;const a=t==="quotes"?e.quotes:t==="recall"?e.recallQuotes:((i=e.historyDetail)==null?void 0:i.Quotes)??[],r=new Set(s);return a.filter(l=>r.has(l.ID))}function dt(t){const s=new Set(t);return e.historyEntries.filter(a=>s.has(a.ID))}function pt(t){return _(t.settings,[])}function _(t,s){var r,i;const a={host:t.Provider.Host,port:String(t.Provider.Port),https:t.Provider.HTTPS,mockLLM:((r=t.Debug)==null?void 0:r.MockLLM)??!1,apiKey:t.Provider.APIKey,modelFilter:"",model:t.Provider.Model,maxResults:String(t.Search.MaxResults),minRelevance:String(t.Search.MinRelevance),theme:t.Theme||"violet",webPort:String(((i=t.Web)==null?void 0:i.Port)??9527),rootDir:t.RootDir??"",models:s};return M(a),a}function yt(){return{host:"",port:"11434",https:!1,mockLLM:!1,apiKey:"",modelFilter:"",model:"",maxResults:"5",minRelevance:"0",theme:"violet",webPort:"9527",rootDir:"",models:[]}}function ee(t){const s=Number.parseInt(t.port.trim(),10);if(!Number.isInteger(s)||s<1||s>65535)throw new Error("Port must be a number between 1 and 65535.");return{Host:t.host.trim(),Port:s,HTTPS:t.https,APIKey:t.apiKey,Model:t.model}}function vt(t){const s=ee(t),a=Number.parseInt(t.maxResults.trim(),10),r=Number.parseInt(t.webPort.trim(),10);if(!Number.isInteger(a)||a<1||a>20)throw new Error("Max ref quotes must be between 1 and 20.");if(!Number.isInteger(r)||r<1||r>65535)throw new Error("Web port must be a number between 1 and 65535.");const i=Number.parseFloat(t.minRelevance.trim());if(Number.isNaN(i))throw new Error("Min relevance must be a decimal number.");if(i<0||i>1)throw new Error("Min relevance must be between 0.0 and 1.0.");return{Provider:s,Search:{MaxResults:a,MinRelevance:i},Debug:{MockLLM:t.mockLLM},Theme:t.theme,Web:{Port:r},RootDir:t.rootDir}}function te(t){const s=t.modelFilter.trim().toLowerCase();return s?t.models.filter(a=>a.toLowerCase().includes(s)):t.models}function M(t){if(t.models.length===0)return;const s=te(t);s.length!==0&&(s.includes(t.model)||(t.model=s[0]))}function k(t,s){return t.map(a=>a.ID===s.ID?s:a)}function h(t,s){return s.length===0?0:Math.min(Math.max(t,0),s.length-1)}function se(t,s){const a=new Set(s.map(r=>r.ID));return new Set([...t].filter(r=>a.has(r)))}function H(t,s){return s.length===0?0:Math.min(Math.max(t,0),s.length-1)}function ht(t,s){const a=new Set(s.map(r=>r.ID));return new Set([...t].filter(r=>a.has(r)))}function n(t){return t.replaceAll("&","&amp;").replaceAll("<","&lt;").replaceAll(">","&gt;").replaceAll('"',"&quot;").replaceAll("'","&#39;")}function p(t){return n(t)}function w(t,s){const a=t.replace(/\s+/g," ").trim();return a.length<=s?a:`${a.slice(0,s-1).trimEnd()}…`}function ft(t){const s=new Date(t);return Number.isNaN(s.getTime())?t:s.toLocaleString()}function bt(t,s){return t.length===0?"":t.length<=s?t.join(" · "):`${t.slice(0,s).join(" · ")} · +${t.length-s} more`}function D(t,s){return w(t,Math.max(8,s))}function v(){var t;return((t=e.auth)==null?void 0:t.runtime)==="web"}function mt(t,s){const a=new Blob([s],{type:"application/json;charset=utf-8"}),r=URL.createObjectURL(a),i=document.createElement("a");i.href=r,i.download=t,document.body.appendChild(i),i.click(),i.remove(),URL.revokeObjectURL(r)}function c(t){return t instanceof Error?t.message:String(t)}function ae(){var s,a,r;const t=[(s=window.go)==null?void 0:s.backend,(a=window.go)==null?void 0:a.app,(r=window.go)==null?void 0:r.main];for(const i of t)if(i!=null&&i.App)return i.App;return null}async function re(t=3e3){const s=Date.now();for(;;){const a=ae();if(a)return a;if(Date.now()-s>=t)throw new Error("Wails backend bridge is unavailable.");await new Promise(r=>window.setTimeout(r,25))}}function u(){const t=ae();if(!t)throw new Error("Wails backend bridge is unavailable.");return t}const oe=document.querySelector("#app");if(!oe)throw new Error("Missing #app root");ue(oe);

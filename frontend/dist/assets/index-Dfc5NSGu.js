(function(){const s=document.createElement("link").relList;if(s&&s.supports&&s.supports("modulepreload"))return;for(const o of document.querySelectorAll('link[rel="modulepreload"]'))r(o);new MutationObserver(o=>{for(const n of o)if(n.type==="childList")for(const u of n.addedNodes)u.tagName==="LINK"&&u.rel==="modulepreload"&&r(u)}).observe(document,{childList:!0,subtree:!0});function a(o){const n={};return o.integrity&&(n.integrity=o.integrity),o.referrerPolicy&&(n.referrerPolicy=o.referrerPolicy),o.crossOrigin==="use-credentials"?n.credentials="include":o.crossOrigin==="anonymous"?n.credentials="omit":n.credentials="same-origin",n}function r(o){if(o.ep)return;o.ep=!0;const n=a(o);fetch(o.href,n)}})();const D={violet:{bg:"#0f172a",bgStrong:"#0b1120",panel:"rgba(17, 24, 39, 0.92)",panel2:"rgba(31, 41, 55, 0.82)",border:"#374151",borderStrong:"rgba(167, 139, 250, 0.42)",primary:"#7c3aed",accent:"#a78bfa",muted:"#94a3b8",fg:"#f9fafb",ok:"#10b981",error:"#ef4444",shadow:"0 24px 80px rgba(2, 6, 23, 0.38)",colorScheme:"dark"},forest:{bg:"#071a17",bgStrong:"#041311",panel:"rgba(9, 24, 21, 0.92)",panel2:"rgba(15, 41, 35, 0.82)",border:"#29443f",borderStrong:"rgba(45, 212, 191, 0.42)",primary:"#0f766e",accent:"#2dd4bf",muted:"#9ca3af",fg:"#ecfdf5",ok:"#22c55e",error:"#ef4444",shadow:"0 24px 80px rgba(1, 10, 9, 0.38)",colorScheme:"dark"},sunset:{bg:"#1c0f0a",bgStrong:"#130905",panel:"rgba(33, 17, 12, 0.94)",panel2:"rgba(52, 28, 18, 0.82)",border:"#5c4033",borderStrong:"rgba(251, 146, 60, 0.44)",primary:"#c2410c",accent:"#fb923c",muted:"#d6b8a6",fg:"#fffbeb",ok:"#16a34a",error:"#dc2626",shadow:"0 24px 80px rgba(20, 8, 2, 0.4)",colorScheme:"dark"},ocean:{bg:"#081824",bgStrong:"#06111a",panel:"rgba(11, 25, 38, 0.92)",panel2:"rgba(18, 40, 56, 0.82)",border:"#334155",borderStrong:"rgba(56, 189, 248, 0.42)",primary:"#0369a1",accent:"#38bdf8",muted:"#94a3b8",fg:"#f8fafc",ok:"#10b981",error:"#ef4444",shadow:"0 24px 80px rgba(3, 9, 16, 0.38)",colorScheme:"dark"},paper:{bg:"#f8fafc",bgStrong:"#e2e8f0",panel:"rgba(255, 255, 255, 0.96)",panel2:"rgba(248, 250, 252, 0.94)",border:"#cbd5e1",borderStrong:"rgba(29, 78, 216, 0.28)",primary:"#1d4ed8",accent:"#0f766e",muted:"#64748b",fg:"#111827",ok:"#15803d",error:"#b91c1c",shadow:"0 24px 80px rgba(148, 163, 184, 0.3)",colorScheme:"light"}};function L(t){const s=D[t in D?t:"violet"],a=document.documentElement;a.style.setProperty("--bg",s.bg),a.style.setProperty("--bg-strong",s.bgStrong),a.style.setProperty("--panel",s.panel),a.style.setProperty("--panel-2",s.panel2),a.style.setProperty("--border",s.border),a.style.setProperty("--border-strong",s.borderStrong),a.style.setProperty("--primary",s.primary),a.style.setProperty("--accent",s.accent),a.style.setProperty("--muted",s.muted),a.style.setProperty("--fg",s.fg),a.style.setProperty("--ok",s.ok),a.style.setProperty("--error",s.error),a.style.setProperty("--shadow",s.shadow),a.style.setProperty("color-scheme",s.colorScheme),document.body.dataset.theme=t}function ie(){return Object.keys(D)}const e={bootstrapped:!1,fatalError:"",authChecked:!1,auth:null,authBusy:!1,authPassword:"",authConfirmPassword:"",authStatus:"",authIsError:!1,page:"Recall",bootstrap:null,quotes:[],quotesLoading:!1,quotesError:"",quotesCursor:0,quotesSelected:new Set,libraryQuery:"",recallQuestion:"",recallLastQuestion:"",recallKeywords:[],recallQuotes:[],recallResponse:"",recallBusy:!1,recallError:"",recallStatus:"",recallStatusIsError:!1,recallCursor:0,recallSelected:new Set,historyEntries:[],historyLoading:!1,historyError:"",historyCursor:0,historySelected:new Set,historyDetail:null,historyDetailLoading:!1,historyDetailError:"",historyStatus:"",historyStatusIsError:!1,historyQuoteCursor:0,historyQuoteSelected:new Set,settings:dt(),settingsBusy:!1,settingsStatus:"",settingsIsError:!1,passwordForm:{current:"",next:"",confirm:"",busy:!1,status:"",isError:!1},overlay:null,toast:null};let b=null,F=!1,N=null,m=null;const ne=["Recall","History","Quotes","Settings"];function le(t){b=t,F||(ce(t),F=!0),i(),N||(N=ue())}async function ue(){try{const t=await se();if(e.auth=await t.AuthStatus(),e.authChecked=!0,e.auth.runtime==="web"&&!e.auth.authenticated){i();return}await O()}catch(t){e.authChecked=!0,e.bootstrapped=!0,e.fatalError=d(t),i()}}async function O(){var s;await se();const t=await c().BootstrapState();e.bootstrap=t,e.bootstrapped=!0,e.page="Recall",e.settings=ct(t),L(e.settings.theme),(s=t.profile)!=null&&s.DisplayName||(e.overlay={type:"namePrompt",name:"",busy:!1,status:"",isError:!1}),i(),await h()}function ce(t){t.addEventListener("click",s=>{ve(s)}),t.addEventListener("input",he),t.addEventListener("change",be),t.addEventListener("submit",s=>{fe(s)}),window.addEventListener("keydown",s=>{me(s)})}async function j(){if(!(!e.auth||e.authBusy)){if(!e.authPassword.trim()){e.authStatus="Password is required.",e.authIsError=!0,i();return}e.authBusy=!0,e.authStatus="",e.authIsError=!1,i();try{await c().Login(e.authPassword),e.authPassword="",e.authConfirmPassword="",e.auth=await c().AuthStatus(),await O()}catch(t){e.authStatus=d(t),e.authIsError=!0,i()}finally{e.authBusy=!1}}}async function de(){await c().Logout(),e.auth=await c().AuthStatus(),e.bootstrapped=!1,e.bootstrap=null,e.overlay=null,e.quotes=[],e.historyEntries=[],e.historyDetail=null,e.authPassword="",e.authConfirmPassword="",e.authStatus="",e.authIsError=!1,i()}async function pe(){if(!e.passwordForm.busy){e.passwordForm.busy=!0,e.passwordForm.status="",i();try{await c().ChangePassword(e.passwordForm.current,e.passwordForm.next,e.passwordForm.confirm),e.passwordForm={current:"",next:"",confirm:"",busy:!1,status:"Password updated.",isError:!1}}catch(t){e.passwordForm.busy=!1,e.passwordForm.status=d(t),e.passwordForm.isError=!0}i()}}async function ye(t){var r,o;const s=(r=t.files)==null?void 0:r[0];if(!s||((o=e.overlay)==null?void 0:o.type)!=="importQuotes")return;const a=await s.text();e.overlay.filename=s.name,e.overlay.payload=a,e.overlay.path=s.name,e.overlay.status=`Loaded ${s.name}`,e.overlay.isError=!1,i()}async function ve(t){var o,n;const s=t.target;if(!(s instanceof HTMLElement))return;const a=s.closest("[data-action]");if(!a)return;switch(a.dataset.action??""){case"auth-login":await j();return;case"auth-logout":await de();return;case"nav":await ge(a.dataset.page);return;case"quotes-refresh":await h();return;case"history-refresh":await C();return;case"history-view-current":await we();return;case"history-back":w();return;case"recall-save-quote":await Me();return;case"history-save-quote":await He();return;case"history-delete-current":Ee();return;case"history-select-all":We();return;case"history-deselect-all":Je();return;case"quote-add":K("add");return;case"quote-import":Qe();return;case"library-clear-filters":e.libraryQuery="",e.quotesCursor=0,i();return;case"quote-select-all":Ae(a.dataset.context);return;case"quote-deselect-all":Be(a.dataset.context);return;case"quote-edit-current":$e(a.dataset.context);return;case"quote-delete-current":Se(a.dataset.context);return;case"quote-share-current":await qe(a.dataset.context);return;case"quote-inspect":De(a.dataset.context,Number(a.dataset.index??"0"));return;case"set-cursor":if(s.closest("input, button, label"))return;Fe(a.dataset.context,Number(a.dataset.index??"0"));return;case"history-set-cursor":if(s.closest("input, button, label"))return;Ke(Number(a.dataset.index??"0"));return;case"history-open":await $(Number(a.dataset.id??"0"));return;case"share-toggle-payload":((o=e.overlay)==null?void 0:o.type)==="shareQuotes"&&(e.overlay.showPayload=!e.overlay.showPayload,i());return;case"import-toggle-payload":((n=e.overlay)==null?void 0:n.type)==="importQuotes"&&(e.overlay.showPayload=!e.overlay.showPayload,i());return;case"profile-save":await U();return;case"quote-editor-save":await W();return;case"quote-editor-refine":await J();return;case"quote-editor-apply-refined":ke();return;case"quote-editor-reject-refined":Ie();return;case"overlay-close":Y();return;case"delete-confirm":await Le();return;case"share-browse":await Ce();return;case"share-save":await Pe();return;case"import-browse":await Re();return;case"import-run":await xe();return;case"settings-fetch-models":await Te();return;case"settings-save":await z();return;case"settings-change-password":await pe();return;case"recall-run":await P();return;case"use-last-question":e.recallQuestion=e.recallLastQuestion,i();return;case"reuse-history-question":if(e.historyDetail)e.recallQuestion=e.historyDetail.Question;else{const u=q()[0];u&&(e.recallQuestion=u.Question)}e.page="Recall",i();return;default:return}}function he(t){var r,o,n,u;const s=t.target;if(!(s instanceof HTMLInputElement||s instanceof HTMLTextAreaElement))return;switch(s.dataset.bind??""){case"auth-password":e.authPassword=s.value;return;case"auth-confirm-password":e.authConfirmPassword=s.value;return;case"recall-question":e.recallQuestion=s.value;return;case"library-query":e.libraryQuery=s.value,e.quotesCursor=0,i();return;case"profile-name":((r=e.overlay)==null?void 0:r.type)==="namePrompt"&&(e.overlay.name=s.value);return;case"quote-editor-content":((o=e.overlay)==null?void 0:o.type)==="quoteEditor"&&(e.overlay.content=s.value);return;case"share-path":((n=e.overlay)==null?void 0:n.type)==="shareQuotes"&&(e.overlay.path=s.value);return;case"import-path":((u=e.overlay)==null?void 0:u.type)==="importQuotes"&&(e.overlay.path=s.value);return;case"settings-host":e.settings.host=s.value;return;case"settings-port":e.settings.port=s.value;return;case"settings-api-key":e.settings.apiKey=s.value;return;case"settings-model-filter":e.settings.modelFilter=s.value,M(e.settings),i();return;case"settings-max-results":e.settings.maxResults=s.value;return;case"settings-min-relevance":e.settings.minRelevance=s.value;return;case"settings-theme":e.settings.theme=s.value,L(e.settings.theme);return;case"settings-web-port":e.settings.webPort=s.value;return;case"settings-password-current":e.passwordForm.current=s.value;return;case"settings-password-next":e.passwordForm.next=s.value;return;case"settings-password-confirm":e.passwordForm.confirm=s.value;return;default:return}}function be(t){const s=t.target;if(!(s instanceof HTMLInputElement||s instanceof HTMLSelectElement))return;switch(s.dataset.bind??""){case"quote-selected":Ne(s.dataset.context,Number(s.dataset.id??"0"),s.checked);return;case"history-selected":Ue(Number(s.dataset.id??"0"),s.checked);return;case"settings-https":s instanceof HTMLInputElement&&(e.settings.https=s.checked);return;case"settings-mock-llm":s instanceof HTMLInputElement&&(e.settings.mockLLM=s.checked);return;case"settings-model":e.settings.model=s.value;return;case"import-file":s instanceof HTMLInputElement&&ye(s);return;default:return}}async function fe(t){const s=t.target;if(s instanceof HTMLFormElement)switch(t.preventDefault(),s.dataset.form){case"auth-login":await j();return;case"auth-setup":await submitAuthSetup();return;case"recall":await P();return;case"profile":await U();return;default:return}}async function me(t){var a,r;const s=document.activeElement;if(t.key==="Escape"&&e.overlay&&e.overlay.type!=="namePrompt"){t.preventDefault(),Y();return}if(t.ctrlKey&&t.key.toLowerCase()==="s"){if(((a=e.overlay)==null?void 0:a.type)==="quoteEditor"){t.preventDefault(),await W();return}!e.overlay&&e.page==="Settings"&&(t.preventDefault(),await z());return}if(t.ctrlKey&&t.key.toLowerCase()==="r"&&((r=e.overlay)==null?void 0:r.type)==="quoteEditor"){t.preventDefault(),await J();return}t.key==="Enter"&&!t.shiftKey&&s instanceof HTMLInputElement&&s.dataset.bind==="recall-question"&&(t.preventDefault(),await P())}async function ge(t){e.page=t,i(),t==="Quotes"&&await h(),t==="History"&&await C()}async function h(){e.quotesLoading=!0,e.quotesError="",i();try{const t=await c().ListQuotes();e.quotes=t,e.quotesCursor=v(e.quotesCursor,t),e.quotesSelected=ee(e.quotesSelected,t),e.quotesError=""}catch(t){e.quotesError=d(t)}finally{e.quotesLoading=!1,i()}}async function C(){var t;e.historyLoading=!0,e.historyError="",e.historyStatus="",e.historyStatusIsError=!1,i();try{const s=await c().ListRecallHistory(),a=((t=e.historyDetail)==null?void 0:t.ID)??null;e.historyEntries=s,e.historyCursor=H(e.historyCursor,s),e.historySelected=yt(e.historySelected,s),a===null?w():s.some(r=>r.ID===a)?$(a,!0):w()}catch(s){e.historyError=d(s)}finally{e.historyLoading=!1,i()}}async function we(){const t=q()[0];t&&await $(t.ID)}async function $(t,s=!1){var a;if(!(!Number.isFinite(t)||t<=0)&&!(e.historyDetailLoading&&((a=e.historyDetail)==null?void 0:a.ID)===t)&&!(e.historyDetail&&e.historyDetail.ID===t&&!e.historyDetailError)){e.historyDetailLoading=!0,e.historyDetailError="",s||(e.historyStatus="",e.historyStatusIsError=!1),i();try{const r=await c().GetRecallHistory(t);e.historyDetail=r,e.historyQuoteCursor=v(e.historyQuoteCursor,r.Quotes),e.historyQuoteSelected=ee(e.historyQuoteSelected,r.Quotes)}catch(r){e.historyDetailError=d(r)}finally{e.historyDetailLoading=!1,i()}}}function w(){e.historyDetail=null,e.historyDetailLoading=!1,e.historyDetailError="",e.historyQuoteCursor=0,e.historyQuoteSelected=new Set,i()}function K(t,s){e.overlay={type:"quoteEditor",mode:t,quoteId:(s==null?void 0:s.ID)??null,content:(s==null?void 0:s.Content)??"",busy:!1,status:"",isError:!1,previewOriginal:"",previewRefined:""},i()}function $e(t){const s=Oe(t);s&&K("edit",s)}function Se(t){var a;const s=((a=e.overlay)==null?void 0:a.type)==="quoteInspect"&&e.overlay.context===t?[e.overlay.quote.ID]:E(t).map(r=>r.ID);s.length!==0&&(e.overlay={type:"deleteQuotes",context:t,ids:s,busy:!1,status:"",isError:!1},i())}function Ee(){const t=q().map(s=>s.ID);t.length!==0&&(e.overlay={type:"deleteHistory",ids:t,busy:!1,status:"",isError:!1},i())}async function qe(t){var a,r,o;const s=((a=e.overlay)==null?void 0:a.type)==="quoteInspect"&&e.overlay.context===t?[e.overlay.quote]:E(t);if(s.length!==0){e.overlay={type:"shareQuotes",context:t,ids:s.map(n=>n.ID),path:"",payload:"",showPayload:!1,busy:!0,status:"",isError:!1},i();try{const n=await c().PreviewQuoteExport(s.map(u=>u.ID));if(((r=e.overlay)==null?void 0:r.type)!=="shareQuotes")return;e.overlay.payload=n,e.overlay.busy=!1,e.overlay.status="Share payload ready. Save it to a file and transfer it manually.",e.overlay.isError=!1}catch(n){if(((o=e.overlay)==null?void 0:o.type)!=="shareQuotes")return;e.overlay.busy=!1,e.overlay.status=d(n),e.overlay.isError=!0}i()}}function Qe(){e.overlay={type:"importQuotes",path:"",payload:"",filename:"",showPayload:!1,busy:!1,status:"",isError:!1,result:null},i()}function De(t,s){const a=S(t),r=v(s,a),o=a[r];o&&(t==="quotes"?e.quotesCursor=r:t==="recall"?e.recallCursor=r:e.historyQuoteCursor=r,e.overlay={type:"quoteInspect",context:t,quote:o},i())}async function U(){var s,a;if(((s=e.overlay)==null?void 0:s.type)!=="namePrompt"||e.overlay.busy)return;const t=e.overlay.name.trim();if(!t){e.overlay.status="Please enter a name to continue.",e.overlay.isError=!0,i();return}e.overlay.busy=!0,e.overlay.status="Saving profile…",e.overlay.isError=!1,i();try{const r=await c().SaveUserProfile(t);e.bootstrap&&(e.bootstrap.profile=r,e.bootstrap.greeting=`Hi! ${r.DisplayName}`),e.overlay=null}catch(r){((a=e.overlay)==null?void 0:a.type)==="namePrompt"&&(e.overlay.busy=!1,e.overlay.status=d(r),e.overlay.isError=!0)}i()}async function W(){var s,a;if(((s=e.overlay)==null?void 0:s.type)!=="quoteEditor"||e.overlay.busy)return;const t=e.overlay.content.trim();if(!t){e.overlay.status="Nothing to save.",e.overlay.isError=!0,i();return}e.overlay.busy=!0,e.overlay.status="Refining draft...",e.overlay.isError=!1,i();try{const r=e.overlay.mode==="add"?await c().AddQuote(t):await c().UpdateQuote(e.overlay.quoteId??0,t);e.overlay=null,R(r),await h()}catch(r){((a=e.overlay)==null?void 0:a.type)==="quoteEditor"&&(e.overlay.busy=!1,e.overlay.status=d(r),e.overlay.isError=!0),i()}}async function J(){var s,a,r;if(((s=e.overlay)==null?void 0:s.type)!=="quoteEditor"||e.overlay.busy)return;const t=e.overlay.content.trim();if(!t){e.overlay.status="Nothing to refine.",e.overlay.isError=!0,i();return}e.overlay.busy=!0,e.overlay.status="",i();try{const o=await c().RefineQuoteDraft(t);if(((a=e.overlay)==null?void 0:a.type)!=="quoteEditor")return;e.overlay.busy=!1,e.overlay.previewOriginal=t,e.overlay.previewRefined=o,e.overlay.status="",e.overlay.isError=!1}catch(o){((r=e.overlay)==null?void 0:r.type)==="quoteEditor"&&(e.overlay.busy=!1,e.overlay.status=d(o),e.overlay.isError=!0)}i()}function ke(){var t;((t=e.overlay)==null?void 0:t.type)==="quoteEditor"&&(e.overlay.content=e.overlay.previewRefined,e.overlay.previewOriginal="",e.overlay.previewRefined="",e.overlay.status="Refined draft applied. Review it, then save.",e.overlay.isError=!1,i())}function Ie(){var t;((t=e.overlay)==null?void 0:t.type)==="quoteEditor"&&(e.overlay.previewOriginal="",e.overlay.previewRefined="",e.overlay.status="Refined draft discarded.",e.overlay.isError=!1,i())}async function Le(){var t,s,a,r;if(((t=e.overlay)==null?void 0:t.type)==="deleteHistory"){if(e.overlay.busy)return;e.overlay.busy=!0,e.overlay.status="",i();try{await c().DeleteRecallHistory(e.overlay.ids),ze(e.overlay.ids),e.overlay=null,await C()}catch(o){((s=e.overlay)==null?void 0:s.type)==="deleteHistory"&&(e.overlay.busy=!1,e.overlay.status=d(o),e.overlay.isError=!0),i()}return}if(!(((a=e.overlay)==null?void 0:a.type)!=="deleteQuotes"||e.overlay.busy)){e.overlay.busy=!0,e.overlay.status="",i();try{await c().DeleteQuotes(e.overlay.ids),je(e.overlay.ids),e.overlay=null,await h()}catch(o){((r=e.overlay)==null?void 0:r.type)==="deleteQuotes"&&(e.overlay.busy=!1,e.overlay.status=d(o),e.overlay.isError=!0),i()}}}async function Ce(){var t,s,a;if(!(((t=e.overlay)==null?void 0:t.type)!=="shareQuotes"||e.overlay.busy)){if(y()){e.overlay.path="irecall-share.json",i();return}try{const r=await c().SelectQuoteExportFile();r&&((s=e.overlay)==null?void 0:s.type)==="shareQuotes"&&(e.overlay.path=r,i())}catch(r){((a=e.overlay)==null?void 0:a.type)==="shareQuotes"&&(e.overlay.status=d(r),e.overlay.isError=!0,i())}}}async function Pe(){var s,a,r;if(((s=e.overlay)==null?void 0:s.type)!=="shareQuotes"||e.overlay.busy)return;if(y()){const o=e.overlay.path.trim()||"irecall-share.json";if(!e.overlay.payload.trim()){e.overlay.status="Export payload is not ready yet.",e.overlay.isError=!0,i();return}bt(o,e.overlay.payload),e.overlay.status=`Downloaded ${o}`,e.overlay.isError=!1,i();return}const t=e.overlay.path.trim();if(!t){e.overlay.status="Choose a file path for the export.",e.overlay.isError=!0,i();return}if(!e.overlay.payload.trim()){e.overlay.status="Export payload is not ready yet.",e.overlay.isError=!0,i();return}e.overlay.busy=!0,e.overlay.status="",i();try{await c().ExportQuotesToFile(e.overlay.ids,t),((a=e.overlay)==null?void 0:a.type)==="shareQuotes"&&(e.overlay.busy=!1,e.overlay.status=`Saved share payload to ${t}`,e.overlay.isError=!1,i())}catch(o){((r=e.overlay)==null?void 0:r.type)==="shareQuotes"&&(e.overlay.busy=!1,e.overlay.status=d(o),e.overlay.isError=!0,i())}}async function Re(){var t,s,a;if(!(((t=e.overlay)==null?void 0:t.type)!=="importQuotes"||e.overlay.busy)){if(y()){const r=document.querySelector('[data-bind="import-file"]');r==null||r.click();return}try{const r=await c().SelectQuoteImportFile();r&&((s=e.overlay)==null?void 0:s.type)==="importQuotes"&&(e.overlay.path=r,i())}catch(r){((a=e.overlay)==null?void 0:a.type)==="importQuotes"&&(e.overlay.status=d(r),e.overlay.isError=!0,i())}}}async function xe(){var a,r,o;if(((a=e.overlay)==null?void 0:a.type)!=="importQuotes"||e.overlay.busy)return;const t=e.overlay.path.trim(),s=e.overlay.payload.trim();if(y()){if(!s){e.overlay.status="Choose a file to import.",e.overlay.isError=!0,i();return}}else if(!t){e.overlay.status="Choose a file to import.",e.overlay.isError=!0,i();return}e.overlay.busy=!0,e.overlay.status="",e.overlay.result=null,i();try{const n=y()?await c().ImportQuotesPayload(s):await c().ImportQuotesFromFile(t);if(((r=e.overlay)==null?void 0:r.type)!=="importQuotes")return;e.overlay.busy=!1,e.overlay.result=n,e.overlay.status=`Imported quotes. inserted=${n.Inserted} updated=${n.Updated} duplicates=${n.Duplicates} stale=${n.Stale}`,e.overlay.isError=!1,await h()}catch(n){((o=e.overlay)==null?void 0:o.type)==="importQuotes"&&(e.overlay.busy=!1,e.overlay.status=d(n),e.overlay.isError=!0,i())}}async function P(){if(e.recallBusy)return;const t=e.recallQuestion.trim();if(!t){e.recallError="Enter a recall question first.",i();return}e.recallBusy=!0,e.recallError="",e.recallStatus="",e.recallStatusIsError=!1,e.recallLastQuestion=t,e.recallKeywords=[],e.recallQuotes=[],e.recallResponse="",e.recallCursor=0,e.recallSelected=new Set,i();try{const s=await c().RunRecall(t);e.recallKeywords=s.keywords,e.recallQuotes=s.quotes,e.recallResponse=s.response,e.recallLastQuestion=s.question||t,e.recallCursor=0,e.recallSelected=new Set,e.recallQuestion=""}catch(s){e.recallError=d(s)}finally{e.recallBusy=!1,i()}}async function Me(){const t=e.recallLastQuestion.trim(),s=e.recallResponse.trim();if(!t||!s){e.recallStatus="Run a recall first before saving it as a quote.",e.recallStatusIsError=!0,i();return}try{const a=await c().SaveRecallAsQuote(t,s,e.recallKeywords);R(a),await h(),e.recallStatus="Saved recall as quote.",e.recallStatusIsError=!1,V("Saved the current grounded answer as a quote.")}catch(a){e.recallStatus=d(a),e.recallStatusIsError=!0}i()}async function He(){const t=e.historyDetail;if(t){try{const s=await c().SaveRecallAsQuote(t.Question,t.Response,[]);R(s),await h(),e.historyStatus="Saved history entry as quote.",e.historyStatusIsError=!1,V("Saved the selected activity session as a quote.")}catch(s){e.historyStatus=d(s),e.historyStatusIsError=!0}i()}}async function Te(){if(e.settingsBusy)return;let t;try{t=Z(e.settings)}catch(s){e.settingsStatus=d(s),e.settingsIsError=!0,i();return}e.settingsBusy=!0,e.settingsStatus="",i();try{const s=await c().FetchModels(t);e.settings.models=s,M(e.settings),e.settingsStatus=s.length>0?`Fetched ${s.length} models.`:"No models returned.",e.settingsIsError=!1}catch(s){e.settingsStatus=d(s),e.settingsIsError=!0}finally{e.settingsBusy=!1,i()}}async function z(){var s;if(e.settingsBusy)return;let t;try{t=pt(e.settings)}catch(a){e.settingsStatus=d(a),e.settingsIsError=!0,i();return}e.settingsBusy=!0,e.settingsStatus="",i();try{const a=await c().SaveSettings(t);e.settings=X(a,e.settings.models),L(e.settings.theme),e.bootstrap&&(e.bootstrap.settings=a);const r=((s=e.auth)==null?void 0:s.runtime)==="web"&&e.auth.currentPort>0&&e.auth.currentPort!==a.Web.Port;e.settingsStatus=r?"Saved. Restart the web server to apply the new port.":"Saved.",e.settingsIsError=!1}catch(a){e.settingsStatus=d(a),e.settingsIsError=!0}finally{e.settingsBusy=!1,i()}}function Y(){e.overlay&&e.overlay.type!=="namePrompt"&&("busy"in e.overlay&&e.overlay.busy||(e.overlay=null,i()))}function V(t,s=!1){e.toast={message:t,isError:s},m!==null&&window.clearTimeout(m),m=window.setTimeout(()=>{e.toast=null,m=null,i()},2600)}function S(t){var s;return t==="quotes"?G():t==="recall"?e.recallQuotes:((s=e.historyDetail)==null?void 0:s.Quotes)??[]}function Fe(t,s){var r;const a=S(t);if(t==="quotes")e.quotesCursor=v(s,a);else if(t==="recall")e.recallCursor=v(s,e.recallQuotes);else{const o=((r=e.historyDetail)==null?void 0:r.Quotes)??[];e.historyQuoteCursor=v(s,o)}i()}function Ne(t,s,a){const r=t==="quotes"?e.quotesSelected:t==="recall"?e.recallSelected:e.historyQuoteSelected;a?r.add(s):r.delete(s)}function Ae(t){const s=S(t),a=new Set(s.map(r=>r.ID));t==="quotes"?e.quotesSelected=a:t==="recall"?e.recallSelected=a:e.historyQuoteSelected=a,i()}function Be(t){t==="quotes"?e.quotesSelected=new Set:t==="recall"?e.recallSelected=new Set:e.historyQuoteSelected=new Set,i()}function E(t){const s=S(t),a=t==="quotes"?e.quotesCursor:t==="recall"?e.recallCursor:e.historyQuoteCursor,r=v(a,s),o=t==="quotes"?e.quotesSelected:t==="recall"?e.recallSelected:e.historyQuoteSelected,n=s.filter(u=>o.has(u.ID));return n.length>0?n:s[r]?[s[r]]:[]}function Oe(t){var s;return((s=e.overlay)==null?void 0:s.type)==="quoteInspect"&&e.overlay.context===t?e.overlay.quote:E(t)[0]??null}function G(){const t=e.libraryQuery.trim().toLowerCase();return e.quotes.filter(s=>t?[s.Content,s.AuthorName,s.SourceName,...s.Tags].join(" ").toLowerCase().includes(t):!0)}function R(t){e.quotes=Q(e.quotes,t),e.recallQuotes=Q(e.recallQuotes,t),e.historyDetail&&(e.historyDetail={...e.historyDetail,Quotes:Q(e.historyDetail.Quotes,t)}),i()}function je(t){var a;const s=new Set(t);e.quotes=e.quotes.filter(r=>!s.has(r.ID)),e.recallQuotes=e.recallQuotes.filter(r=>!s.has(r.ID)),e.historyDetail&&(e.historyDetail={...e.historyDetail,Quotes:e.historyDetail.Quotes.filter(r=>!s.has(r.ID))}),e.quotesSelected=new Set([...e.quotesSelected].filter(r=>!s.has(r))),e.recallSelected=new Set([...e.recallSelected].filter(r=>!s.has(r))),e.historyQuoteSelected=new Set([...e.historyQuoteSelected].filter(r=>!s.has(r))),e.quotesCursor=v(e.quotesCursor,e.quotes),e.recallCursor=v(e.recallCursor,e.recallQuotes),e.historyQuoteCursor=v(e.historyQuoteCursor,((a=e.historyDetail)==null?void 0:a.Quotes)??[]),i()}function Ke(t){e.historyCursor=H(t,e.historyEntries);const s=e.historyEntries[e.historyCursor];s&&$(s.ID,!0),i()}function Ue(t,s){s?e.historySelected.add(t):e.historySelected.delete(t)}function We(){e.historySelected=new Set(e.historyEntries.map(t=>t.ID)),i()}function Je(){e.historySelected=new Set,i()}function q(){const t=e.historyEntries.filter(s=>e.historySelected.has(s.ID));return t.length>0?t:e.historyEntries[e.historyCursor]?[e.historyEntries[e.historyCursor]]:[]}function ze(t){const s=new Set(t);if(e.historyEntries=e.historyEntries.filter(a=>!s.has(a.ID)),e.historySelected=new Set([...e.historySelected].filter(a=>!s.has(a))),e.historyCursor=H(e.historyCursor,e.historyEntries),e.historyDetail&&s.has(e.historyDetail.ID)){w();return}i()}function i(){if(!b)return;const t=Ye();b.innerHTML=Ge(),Ve(t)}function Ye(){const t=document.activeElement;if(!(t instanceof HTMLInputElement||t instanceof HTMLTextAreaElement||t instanceof HTMLSelectElement))return null;const s=t.dataset.bind;return s?{selector:`[data-bind="${s}"]`,selectionStart:t instanceof HTMLInputElement||t instanceof HTMLTextAreaElement?t.selectionStart:null,selectionEnd:t instanceof HTMLInputElement||t instanceof HTMLTextAreaElement?t.selectionEnd:null}:null}function Ve(t){if(!b||!t)return;const s=b.querySelector(t.selector);(s instanceof HTMLInputElement||s instanceof HTMLTextAreaElement||s instanceof HTMLSelectElement)&&(s.focus({preventScroll:!0}),(s instanceof HTMLInputElement||s instanceof HTMLTextAreaElement)&&t.selectionStart!==null&&t.selectionEnd!==null&&s.setSelectionRange(t.selectionStart,t.selectionEnd))}function Ge(){var s,a,r,o,n;if(!e.authChecked)return`
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
          <div class="status status-error">${l(e.fatalError)}</div>
        </div>
      </div>
    `;if(((s=e.auth)==null?void 0:s.runtime)==="web"&&!e.auth.authenticated)return Xe();if(!e.bootstrapped)return`
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
          <div class="brand">${l(((o=e.bootstrap)==null?void 0:o.productName)??"iRecall")}</div>
          <div class="muted subtle">${y()?"Local-first knowledge workspace for the web":"Local-first knowledge workspace for desktop"}</div>
        </div>
        <div class="titlebar-right">
          <div class="greeting">${l(t)}</div>
          ${((n=e.auth)==null?void 0:n.runtime)==="web"?'<button class="button" data-action="auth-logout" type="button">Logout</button>':""}
          <nav class="tabs" aria-label="Primary">
            ${ne.map(u=>`
                  <button
                    class="tab${e.page===u?" active":""}"
                    data-action="nav"
                    data-page="${u}"
                    type="button"
                  >${_e(u)}</button>
                `).join("")}
          </nav>
        </div>
      </header>

      <main class="layout">
        ${Ze()}
      </main>

      ${at()}
      ${e.overlay?lt(e.overlay):""}
      ${e.toast?nt(e.toast):""}
    </div>
  `}function Xe(){var o;const t=!((o=e.auth)!=null&&o.passwordConfigured),s="auth-login";return`
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
          ${e.authStatus?`<div class="status ${e.authIsError?"status-error":"status-ok"}">${l(e.authStatus)}</div>`:""}
          <div class="modal-actions">
            <button class="button button-primary" data-action="${s}" type="submit" ${e.authBusy||t?"disabled":""}>
              ${e.authBusy?"Working…":"Login"}
            </button>
          </div>
        </form>
      </div>
    </div>
  `}function Ze(){switch(e.page){case"Recall":return et();case"Quotes":return tt();case"History":return st();case"Settings":return ot()}}function _e(t){switch(t){case"Recall":return"Recall";case"Quotes":return"Quotes";case"History":return"History";case"Settings":return"Settings"}}function et(){const t=!e.recallResponse.trim(),s=e.settings.mockLLM,a=e.recallResponse.trim()?l(e.recallResponse):'<span class="muted">Grounded response will appear here.</span>',r=e.recallKeywords.length>0?e.recallKeywords.map(o=>`<span class="keyword-chip">${l(o)}</span>`).join(""):'<span class="muted">Keywords: —</span>';return`
    <section class="page page-recall">
      <div class="panel page-panel">
        <div class="page-hero">
          <div>
            <div class="eyebrow">Recall</div>
            <div class="page-title">Question, references, then answer</div>
            <div class="muted page-copy">Run recall once, inspect the retrieved quotes, then read the grounded response.</div>
          </div>
          ${s?'<div class="meta-row"><span class="meta-pill meta-pill-accent">Mock LLM on</span></div>':""}
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
              >${l(e.recallQuestion)}</textarea>
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
            ${x("recall",e.recallQuotes,e.recallCursor,e.recallSelected,!1)}
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
                      <div class="answer-question">${l(e.recallLastQuestion)}</div>
                    </div>
                  `:""}
              <pre class="response-box">${a}</pre>
            </div>
          </section>
        </div>

        ${e.recallError?`<div class="status status-error">${l(e.recallError)}</div>`:""}
        ${e.recallStatus?`<div class="status ${e.recallStatusIsError?"status-error":"status-ok"}">${l(e.recallStatus)}</div>`:""}
      </div>
    </section>
  `}function tt(){const t=G(),s=v(e.quotesCursor,t),a=E("quotes"),r=t.filter(n=>e.quotesSelected.has(n.ID)).length;return`
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

        <div class="meta-row meta-row-rich">
          ${[`${e.quotes.length} total`,`${e.quotes.filter(n=>n.IsOwnedByMe).length} authored here`,`${e.quotes.filter(n=>!n.IsOwnedByMe).length} imported`].map(n=>`<span class="meta-pill">${l(n)}</span>`).join("")}
          <span class="meta-pill meta-pill-accent">${a.length>0?`${a.length} selected`:"Open a quote to inspect it"}</span>
        </div>

        <div class="workspace workspace-library">
          <section class="panel subpanel">
            <div class="subpanel-header">
              <div>
                <div class="section-title">Quote list</div>
                <div class="muted">${t.length} matching quotes. Choose one to inspect the full note, provenance, and actions.</div>
              </div>
              <div class="toolbar toolbar-quiet">
                <button class="button" data-action="quote-select-all" data-context="quotes" type="button" ${t.length===0?"disabled":""}>Select results</button>
                <button class="button" data-action="quote-deselect-all" data-context="quotes" type="button" ${r===0?"disabled":""}>Clear selection</button>
              </div>
            </div>
            ${e.quotesLoading?'<div class="empty-state">Loading quotes…</div>':e.quotesError?`<div class="status status-error">${l(e.quotesError)}</div>`:x("quotes",t,s,e.quotesSelected,!0)}
          </section>
        </div>
      </div>
    </section>
  `}function st(){const t=q();return`
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
            <button class="button" data-action="reuse-history-question" type="button" ${!!(e.historyDetail??t[0])?"":"disabled"}>Recall again</button>
            <button class="button button-danger" data-action="history-delete-current" type="button" ${t.length===0?"disabled":""}>Delete</button>
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
            ${e.historyLoading?'<div class="empty-state">Loading history…</div>':e.historyError?`<div class="status status-error">${l(e.historyError)}</div>`:e.historyEntries.length===0?'<div class="empty-state">No recall history yet. Run a question from the Recall page to create your first grounded session.</div>':rt()}
          </section>
        </div>
      </div>
    </section>
  `}function at(){if(!e.historyDetailLoading&&!e.historyDetail&&!e.historyDetailError)return"";const t=e.historyDetail?e.historyEntries.find(r=>{var o;return r.ID===((o=e.historyDetail)==null?void 0:o.ID)})??e.historyDetail:e.historyEntries[e.historyCursor]??null,s=e.historyDetail&&t&&e.historyDetail.ID===t.ID?e.historyDetail:null,a=(s==null?void 0:s.Quotes)??[];return`
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
                  <pre class="response-box compact-box">${l(t.Question)}</pre>
                </div>
                <div class="detail-block">
                  <div class="muted">Response</div>
                  <pre class="response-box compact-box">${l((s==null?void 0:s.Response)??t.Response)}</pre>
                </div>
              </div>
              ${e.historyDetailLoading?'<div class="empty-state">Loading reference quotes…</div>':s?`
                      <div class="subpanel-header nested-header">
                        <div>
                          <div class="section-title">Reference quotes</div>
                          <div class="muted">${a.length} retrieved quotes. Open one to inspect the full note.</div>
                        </div>
                      </div>
                      ${x("history",a,e.historyQuoteCursor,e.historyQuoteSelected,!1)}
                    `:""}
            `:""}

        ${e.historyDetailError?`<div class="status status-error">${l(e.historyDetailError)}</div>`:""}
        ${e.historyStatus?`<div class="status ${e.historyStatusIsError?"status-error":"status-ok"}">${l(e.historyStatus)}</div>`:""}
      </div>
    </div>
  `}function rt(){return`
    <div class="history-list">
      ${e.historyEntries.map((t,s)=>{const a=s===e.historyCursor,r=I(t.Response,156);return`
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
                  <span class="quote-version">${l(vt(t.CreatedAt))}</span>
                </div>
              </div>
              <div class="quote-content">${l(I(t.Question,132))}</div>
                <div class="quote-meta"><span class="muted">Response preview</span><span>${l(r||"(empty response)")}</span></div>
            </article>
          `}).join("")}
    </div>
  `}function ot(){var o,n;const t=_(e.settings),s=(o=e.bootstrap)==null?void 0:o.paths,a=(n=e.auth)==null?void 0:n.currentPort,r=e.settings.models.length>0&&t.length>0?`
        <select class="select-input" data-bind="settings-model">
          ${t.map(u=>`
                <option value="${p(u)}"${u===e.settings.model?" selected":""}>${l(u)}</option>
              `).join("")}
        </select>
      `:`
        <div class="readonly-model">
          <span>${l(e.settings.model||"(none)")}</span>
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
              ${r}
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
                ${ie().map(u=>`
                      <option value="${u}"${u===e.settings.theme?" selected":""}>${u}</option>
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
            ${e.passwordForm.status?`<div class="status ${e.passwordForm.isError?"status-error":"status-ok"}">${l(e.passwordForm.status)}</div>`:""}
          </section>

          <section class="panel subpanel settings-secondary">
            <div class="section-title">Advanced</div>
            <label class="field">
              <span>Web Port</span>
              <input class="text-input" data-bind="settings-web-port" value="${p(e.settings.webPort)}" />
            </label>
            <div class="settings-hint muted">
              The web server listens on this port after restart. Current listener: ${l(a?String(a):"(not running)")}.
            </div>
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

        ${e.settingsStatus?`<div class="status ${e.settingsIsError?"status-error":"status-ok"}">${l(e.settingsStatus)}</div>`:""}
      </div>
    </section>
  `}function x(t,s,a,r,o){return s.length===0?`<div class="empty-state">${t==="quotes"?"No quotes yet. Add one or import a shared payload.":"No reference quotes for this question yet."}</div>`:`
    <div class="quote-list">
      ${s.map((n,u)=>{const re=u===a,f=t!=="quotes",T=!n.IsOwnedByMe&&n.SourceName?`<span class="meta-accent">${l(n.SourceName)}</span>`:`<span>${l(n.AuthorName||"You")}</span>`,oe=o?`
              <div class="quote-meta">
                <span class="muted">Tags</span>
                <span>${n.Tags.length>0?l(ht(n.Tags,4)):"(none)"}</span>
              </div>
            `:"";return`
            <article class="quote-card${re?" is-current":""}${f?" quote-card-minimal":""}" data-action="quote-inspect" data-context="${t}" data-index="${u}">
              <div class="quote-topline">
                ${f?`<div class="quote-topline-meta">
                        <span class="quote-badge">${n.IsOwnedByMe?"Owned":"Imported"}</span>
                        <span class="quote-version">${l(k(n.UpdatedAt))}</span>
                      </div>`:`
                      <label class="selection-toggle">
                        <input
                          type="checkbox"
                          data-bind="quote-selected"
                          data-context="${t}"
                          data-id="${n.ID}"
                          ${r.has(n.ID)?"checked":""}
                        />
                  </label>
                    `}
                <div class="quote-topline-meta">
                  ${f?`<span class="quote-source-inline">${T}</span>`:`<span class="quote-version">${l(k(n.UpdatedAt))}</span>
                  <span class="quote-badge">${n.IsOwnedByMe?"Owned":"Imported"}</span>`}
                </div>
              </div>
              <div class="quote-content">${l(I(n.Content,t==="quotes"?160:136))}</div>
              ${f?`<div class="quote-actions-inline"><button class="button button-subtle" data-action="quote-inspect" data-context="${t}" data-index="${u}" type="button">Details</button></div>`:`<div class="quote-meta"><span class="muted">${!n.IsOwnedByMe&&n.SourceName?"Imported from":"Author"}</span> ${T}</div>`}
              ${oe}
            </article>
          `}).join("")}
    </div>
  `}function it(t,s){return t?`
    <div class="detail-stack">
      <div class="detail-block">
        <div class="muted">Full quote</div>
        <pre class="response-box compact-box">${l(t.Content)}</pre>
      </div>
      <div class="detail-grid">
        <div class="detail-metric">
          <span class="muted">Author</span>
          <span>${l(t.AuthorName||"You")}</span>
        </div>
        <div class="detail-metric">
          <span class="muted">Version</span>
          <span>v${t.Version}</span>
        </div>
        <div class="detail-metric">
          <span class="muted">Source</span>
          <span>${l(t.SourceName||"Local library")}</span>
        </div>
        <div class="detail-metric">
          <span class="muted">Updated</span>
          <span>${l(k(t.UpdatedAt))}</span>
        </div>
      </div>
      <div class="detail-block">
        <div class="muted">Tags</div>
        <div class="keyword-list">
          ${t.Tags.length>0?t.Tags.map(a=>`<span class="keyword-chip">${l(a)}</span>`).join(""):'<span class="muted">No tags assigned yet.</span>'}
        </div>
      </div>
      <div class="toolbar toolbar-inline">
        <button class="button" data-action="quote-edit-current" data-context="${s}" type="button">Edit</button>
        <button class="button" data-action="quote-share-current" data-context="${s}" type="button">Share</button>
        <button class="button button-danger" data-action="quote-delete-current" data-context="${s}" type="button">Delete</button>
      </div>
    </div>
  `:'<div class="empty-state">Select a quote to inspect the full note, provenance, and available actions.</div>'}function k(t){const s=new Date(t);return Number.isNaN(s.getTime())?t:s.toLocaleDateString(void 0,{month:"short",day:"numeric",year:"numeric"})}function A(t){return!Number.isFinite(t)||t<=0?"0 B":t<1024?`${Math.round(t)} B`:t<1024*1024?`${(t/1024).toFixed(1)} KB`:`${(t/(1024*1024)).toFixed(1)} MB`}function nt(t){return`
    <div class="toast-stack" role="status" aria-live="polite">
      <div class="toast${t.isError?" is-error":""}">${l(t.message)}</div>
    </div>
  `}function lt(t){switch(t.type){case"namePrompt":return`
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
              ${t.status?`<div class="status ${t.isError?"status-error":"status-ok"}">${l(t.status)}</div>`:""}
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
                      <pre class="compare-body">${l(t.previewOriginal)}</pre>
                    </section>
                    <section class="panel compare-panel">
                      <div class="section-title">Refined Draft</div>
                      <pre class="compare-body">${l(t.previewRefined)}</pre>
                    </section>
                  </div>
                `:`
                  <label class="field">
                    <span>Quote Content</span>
                    <textarea class="text-area" data-bind="quote-editor-content" rows="10" placeholder="Type or paste your note here.">${l(t.content)}</textarea>
                  </label>
                `}
            <div class="muted modal-copy">
              ${t.previewRefined?"Compare the current draft with the suggested rewrite before applying it.":e.settings.mockLLM?"Mock LLM is enabled, so Refine returns the original text and skips provider-dependent rewriting.":"Tags are regenerated automatically by the shared core logic."}
            </div>
            ${t.status?`<div class="status ${t.isError?"status-error":"status-ok"}">${l(t.status)}</div>`:""}
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
              ${B(t.context,t.ids).map((a,r)=>`<div class="summary-item">[${r+1}] ${l(g(a.Content,140))}</div>`).join("")}
            </div>
            ${t.status?`<div class="status ${t.isError?"status-error":"status-ok"}">${l(t.status)}</div>`:""}
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
              ${ut(t.ids).map((a,r)=>`<div class="summary-item">[${r+1}] ${l(g(a.Question,140))}</div>`).join("")}
            </div>
            ${t.status?`<div class="status ${t.isError?"status-error":"status-ok"}">${l(t.status)}</div>`:""}
            <div class="modal-actions">
              <button class="button button-danger" data-action="delete-confirm" type="button" ${t.busy?"disabled":""}>
                ${t.busy?"Deleting…":"Delete"}
              </button>
              <button class="button" data-action="overlay-close" type="button" ${t.busy?"disabled":""}>Cancel</button>
            </div>
          </div>
        </div>
      `;case"shareQuotes":const s=B(t.context,t.ids);return`
        <div class="overlay-backdrop overlay-backdrop-side">
          <div class="modal modal-side">
            <div class="modal-title">Share Quotes</div>
            <div class="modal-copy">Export a portable share file. The file summary comes first; raw JSON is available only if you need to inspect it.</div>
            <div class="summary-list">
              ${s.map((a,r)=>`<div class="summary-item">[${r+1}] v${a.Version} ${l(g(a.Content,120))}</div>`).join("")}
            </div>
            <div class="result-grid">
              <div><span class="muted">Quotes:</span> ${s.length}</div>
              <div><span class="muted">Payload size:</span> ${A(t.payload.length)}</div>
            </div>
            <label class="field">
              <span>${y()?"Download As":"Save To"}</span>
              ${y()?`<input class="text-input" data-bind="share-path" value="${p(t.path||"irecall-share.json")}" placeholder="irecall-share.json" />`:`
                    <div class="path-row">
                      <input class="text-input" data-bind="share-path" value="${p(t.path)}" placeholder="/path/to/irecall-share.json" />
                      <button class="button" data-action="share-browse" type="button" ${t.busy?"disabled":""}>Browse</button>
                    </div>
                  `}
            </label>
            <div class="muted modal-copy">${y()?"Download the JSON payload locally, then transfer it manually to the recipient.":"Export to a JSON file and transfer it manually to the recipient."}</div>
            <div class="toolbar toolbar-inline">
              <button class="button" data-action="share-toggle-payload" type="button" ${t.payload?"":"disabled"}>
                ${t.showPayload?"Hide raw JSON":"Show raw JSON"}
              </button>
            </div>
            ${t.showPayload?`<div class="payload-box"><pre>${l(t.payload||"Preparing export payload…")}</pre></div>`:""}
            ${t.status?`<div class="status ${t.isError?"status-error":"status-ok"}">${l(t.status)}</div>`:""}
            <div class="modal-actions">
              <button class="button button-primary" data-action="share-save" type="button" ${t.busy?"disabled":""}>
                ${t.busy?"Working…":y()?"Download export file":"Save export file"}
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
              ${y()?`
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
                    <div><span class="muted">File:</span> ${l(t.filename||t.path)}</div>
                    <div><span class="muted">Payload size:</span> ${A(t.payload.length)}</div>
                  </div>
                `:""}
            ${t.payload?`
                  <div class="toolbar toolbar-inline">
                    <button class="button" data-action="import-toggle-payload" type="button">
                      ${t.showPayload?"Hide raw JSON":"Show raw JSON"}
                    </button>
                  </div>
                  ${t.showPayload?`<div class="payload-box"><pre>${l(t.payload)}</pre></div>`:""}
                `:""}
            ${t.result?`
                  <div class="result-grid">
                    <div><span class="muted">Inserted:</span> ${t.result.Inserted}</div>
                    <div><span class="muted">Updated:</span> ${t.result.Updated}</div>
                    <div><span class="muted">Duplicates:</span> ${t.result.Duplicates}</div>
                    <div><span class="muted">Stale:</span> ${t.result.Stale}</div>
                  </div>
                `:""}
            ${t.status?`<div class="status ${t.isError?"status-error":"status-ok"}">${l(t.status)}</div>`:""}
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
            ${it(t.quote,t.context)}
            <div class="modal-actions">
              <button class="button" data-action="overlay-close" type="button">Close</button>
            </div>
          </div>
        </div>
      `;case"notice":return""}}function B(t,s){var o;const a=t==="quotes"?e.quotes:t==="recall"?e.recallQuotes:((o=e.historyDetail)==null?void 0:o.Quotes)??[],r=new Set(s);return a.filter(n=>r.has(n.ID))}function ut(t){const s=new Set(t);return e.historyEntries.filter(a=>s.has(a.ID))}function ct(t){return X(t.settings,[])}function X(t,s){var r,o;const a={host:t.Provider.Host,port:String(t.Provider.Port),https:t.Provider.HTTPS,mockLLM:((r=t.Debug)==null?void 0:r.MockLLM)??!1,apiKey:t.Provider.APIKey,modelFilter:"",model:t.Provider.Model,maxResults:String(t.Search.MaxResults),minRelevance:String(t.Search.MinRelevance),theme:t.Theme||"violet",webPort:String(((o=t.Web)==null?void 0:o.Port)??9527),models:s};return M(a),a}function dt(){return{host:"",port:"11434",https:!1,mockLLM:!1,apiKey:"",modelFilter:"",model:"",maxResults:"5",minRelevance:"0",theme:"violet",webPort:"9527",models:[]}}function Z(t){const s=Number.parseInt(t.port.trim(),10);if(!Number.isInteger(s)||s<1||s>65535)throw new Error("Port must be a number between 1 and 65535.");return{Host:t.host.trim(),Port:s,HTTPS:t.https,APIKey:t.apiKey,Model:t.model}}function pt(t){const s=Z(t),a=Number.parseInt(t.maxResults.trim(),10),r=Number.parseInt(t.webPort.trim(),10);if(!Number.isInteger(a)||a<1||a>20)throw new Error("Max ref quotes must be between 1 and 20.");if(!Number.isInteger(r)||r<1||r>65535)throw new Error("Web port must be a number between 1 and 65535.");const o=Number.parseFloat(t.minRelevance.trim());if(Number.isNaN(o))throw new Error("Min relevance must be a decimal number.");if(o<0||o>1)throw new Error("Min relevance must be between 0.0 and 1.0.");return{Provider:s,Search:{MaxResults:a,MinRelevance:o},Debug:{MockLLM:t.mockLLM},Theme:t.theme,Web:{Port:r}}}function _(t){const s=t.modelFilter.trim().toLowerCase();return s?t.models.filter(a=>a.toLowerCase().includes(s)):t.models}function M(t){if(t.models.length===0)return;const s=_(t);s.length!==0&&(s.includes(t.model)||(t.model=s[0]))}function Q(t,s){return t.map(a=>a.ID===s.ID?s:a)}function v(t,s){return s.length===0?0:Math.min(Math.max(t,0),s.length-1)}function ee(t,s){const a=new Set(s.map(r=>r.ID));return new Set([...t].filter(r=>a.has(r)))}function H(t,s){return s.length===0?0:Math.min(Math.max(t,0),s.length-1)}function yt(t,s){const a=new Set(s.map(r=>r.ID));return new Set([...t].filter(r=>a.has(r)))}function l(t){return t.replaceAll("&","&amp;").replaceAll("<","&lt;").replaceAll(">","&gt;").replaceAll('"',"&quot;").replaceAll("'","&#39;")}function p(t){return l(t)}function g(t,s){const a=t.replace(/\s+/g," ").trim();return a.length<=s?a:`${a.slice(0,s-1).trimEnd()}…`}function vt(t){const s=new Date(t);return Number.isNaN(s.getTime())?t:s.toLocaleString()}function ht(t,s){return t.length===0?"":t.length<=s?t.join(" · "):`${t.slice(0,s).join(" · ")} · +${t.length-s} more`}function I(t,s){return g(t,Math.max(8,s))}function y(){var t;return((t=e.auth)==null?void 0:t.runtime)==="web"}function bt(t,s){const a=new Blob([s],{type:"application/json;charset=utf-8"}),r=URL.createObjectURL(a),o=document.createElement("a");o.href=r,o.download=t,document.body.appendChild(o),o.click(),o.remove(),URL.revokeObjectURL(r)}function d(t){return t instanceof Error?t.message:String(t)}function te(){var s,a,r;const t=[(s=window.go)==null?void 0:s.backend,(a=window.go)==null?void 0:a.app,(r=window.go)==null?void 0:r.main];for(const o of t)if(o!=null&&o.App)return o.App;return null}async function se(t=3e3){const s=Date.now();for(;;){const a=te();if(a)return a;if(Date.now()-s>=t)throw new Error("Wails backend bridge is unavailable.");await new Promise(r=>window.setTimeout(r,25))}}function c(){const t=te();if(!t)throw new Error("Wails backend bridge is unavailable.");return t}const ae=document.querySelector("#app");if(!ae)throw new Error("Missing #app root");le(ae);

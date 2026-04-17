(function(){const s=document.createElement("link").relList;if(s&&s.supports&&s.supports("modulepreload"))return;for(const o of document.querySelectorAll('link[rel="modulepreload"]'))r(o);new MutationObserver(o=>{for(const l of o)if(l.type==="childList")for(const u of l.addedNodes)u.tagName==="LINK"&&u.rel==="modulepreload"&&r(u)}).observe(document,{childList:!0,subtree:!0});function a(o){const l={};return o.integrity&&(l.integrity=o.integrity),o.referrerPolicy&&(l.referrerPolicy=o.referrerPolicy),o.crossOrigin==="use-credentials"?l.credentials="include":o.crossOrigin==="anonymous"?l.credentials="omit":l.credentials="same-origin",l}function r(o){if(o.ep)return;o.ep=!0;const l=a(o);fetch(o.href,l)}})();const I={violet:{bg:"#0f172a",bgStrong:"#0b1120",panel:"rgba(17, 24, 39, 0.92)",panel2:"rgba(31, 41, 55, 0.82)",border:"#374151",borderStrong:"rgba(167, 139, 250, 0.42)",primary:"#7c3aed",accent:"#a78bfa",muted:"#94a3b8",fg:"#f9fafb",ok:"#10b981",error:"#ef4444",shadow:"0 24px 80px rgba(2, 6, 23, 0.38)",colorScheme:"dark"},forest:{bg:"#071a17",bgStrong:"#041311",panel:"rgba(9, 24, 21, 0.92)",panel2:"rgba(15, 41, 35, 0.82)",border:"#29443f",borderStrong:"rgba(45, 212, 191, 0.42)",primary:"#0f766e",accent:"#2dd4bf",muted:"#9ca3af",fg:"#ecfdf5",ok:"#22c55e",error:"#ef4444",shadow:"0 24px 80px rgba(1, 10, 9, 0.38)",colorScheme:"dark"},sunset:{bg:"#1c0f0a",bgStrong:"#130905",panel:"rgba(33, 17, 12, 0.94)",panel2:"rgba(52, 28, 18, 0.82)",border:"#5c4033",borderStrong:"rgba(251, 146, 60, 0.44)",primary:"#c2410c",accent:"#fb923c",muted:"#d6b8a6",fg:"#fffbeb",ok:"#16a34a",error:"#dc2626",shadow:"0 24px 80px rgba(20, 8, 2, 0.4)",colorScheme:"dark"},ocean:{bg:"#081824",bgStrong:"#06111a",panel:"rgba(11, 25, 38, 0.92)",panel2:"rgba(18, 40, 56, 0.82)",border:"#334155",borderStrong:"rgba(56, 189, 248, 0.42)",primary:"#0369a1",accent:"#38bdf8",muted:"#94a3b8",fg:"#f8fafc",ok:"#10b981",error:"#ef4444",shadow:"0 24px 80px rgba(3, 9, 16, 0.38)",colorScheme:"dark"},paper:{bg:"#f8fafc",bgStrong:"#e2e8f0",panel:"rgba(255, 255, 255, 0.96)",panel2:"rgba(248, 250, 252, 0.94)",border:"#cbd5e1",borderStrong:"rgba(29, 78, 216, 0.28)",primary:"#1d4ed8",accent:"#0f766e",muted:"#64748b",fg:"#111827",ok:"#15803d",error:"#b91c1c",shadow:"0 24px 80px rgba(148, 163, 184, 0.3)",colorScheme:"light"}};function P(t){const s=I[t in I?t:"violet"],a=document.documentElement;a.style.setProperty("--bg",s.bg),a.style.setProperty("--bg-strong",s.bgStrong),a.style.setProperty("--panel",s.panel),a.style.setProperty("--panel-2",s.panel2),a.style.setProperty("--border",s.border),a.style.setProperty("--border-strong",s.borderStrong),a.style.setProperty("--primary",s.primary),a.style.setProperty("--accent",s.accent),a.style.setProperty("--muted",s.muted),a.style.setProperty("--fg",s.fg),a.style.setProperty("--ok",s.ok),a.style.setProperty("--error",s.error),a.style.setProperty("--shadow",s.shadow),a.style.setProperty("color-scheme",s.colorScheme),document.body.dataset.theme=t}function ne(){return Object.keys(I)}const e={bootstrapped:!1,fatalError:"",authChecked:!1,auth:null,authBusy:!1,authPassword:"",authConfirmPassword:"",authStatus:"",authIsError:!1,page:"Recall",bootstrap:null,quotes:[],quotesLoading:!1,quotesError:"",quotesCursor:0,quotesSelected:new Set,libraryQuery:"",libraryOwnership:"all",libraryTag:null,recallQuestion:"",recallLastQuestion:"",recallKeywords:[],recallQuotes:[],recallResponse:"",recallBusy:!1,recallError:"",recallStatus:"",recallStatusIsError:!1,recallCursor:0,recallSelected:new Set,historyEntries:[],historyLoading:!1,historyError:"",historyCursor:0,historySelected:new Set,historyDetail:null,historyDetailLoading:!1,historyDetailError:"",historyStatus:"",historyStatusIsError:!1,historyQuoteCursor:0,historyQuoteSelected:new Set,settings:dt(),settingsBusy:!1,settingsStatus:"",settingsIsError:!1,passwordForm:{current:"",next:"",confirm:"",busy:!1,status:"",isError:!1},overlay:null,toast:null};let m=null,M=!1,N=null,g=null;const le=["Recall","History","Quotes","Settings"];function ue(t){m=t,M||(de(t),M=!0),i(),N||(N=ce())}async function ce(){try{const t=await re();if(e.auth=await t.AuthStatus(),e.authChecked=!0,e.auth.runtime==="web"&&!e.auth.authenticated){i();return}await j()}catch(t){e.authChecked=!0,e.bootstrapped=!0,e.fatalError=d(t),i()}}async function j(){var s;await re();const t=await c().BootstrapState();e.bootstrap=t,e.bootstrapped=!0,e.page="Recall",e.settings=ct(t),P(e.settings.theme),(s=t.profile)!=null&&s.DisplayName||(e.overlay={type:"namePrompt",name:"",busy:!1,status:"",isError:!1}),i(),await h()}function de(t){t.addEventListener("click",s=>{be(s)}),t.addEventListener("input",he),t.addEventListener("change",fe),t.addEventListener("submit",s=>{me(s)}),window.addEventListener("keydown",s=>{ge(s)})}async function K(){if(!(!e.auth||e.authBusy)){if(!e.authPassword.trim()){e.authStatus="Password is required.",e.authIsError=!0,i();return}e.authBusy=!0,e.authStatus="",e.authIsError=!1,i();try{await c().Login(e.authPassword),e.authPassword="",e.authConfirmPassword="",e.auth=await c().AuthStatus(),await j()}catch(t){e.authStatus=d(t),e.authIsError=!0,i()}finally{e.authBusy=!1}}}async function pe(){await c().Logout(),e.auth=await c().AuthStatus(),e.bootstrapped=!1,e.bootstrap=null,e.overlay=null,e.quotes=[],e.historyEntries=[],e.historyDetail=null,e.authPassword="",e.authConfirmPassword="",e.authStatus="",e.authIsError=!1,i()}async function ye(){if(!e.passwordForm.busy){e.passwordForm.busy=!0,e.passwordForm.status="",i();try{await c().ChangePassword(e.passwordForm.current,e.passwordForm.next,e.passwordForm.confirm),e.passwordForm={current:"",next:"",confirm:"",busy:!1,status:"Password updated.",isError:!1}}catch(t){e.passwordForm.busy=!1,e.passwordForm.status=d(t),e.passwordForm.isError=!0}i()}}async function ve(t){var r,o;const s=(r=t.files)==null?void 0:r[0];if(!s||((o=e.overlay)==null?void 0:o.type)!=="importQuotes")return;const a=await s.text();e.overlay.filename=s.name,e.overlay.payload=a,e.overlay.path=s.name,e.overlay.status=`Loaded ${s.name}`,e.overlay.isError=!1,i()}async function be(t){var o,l;const s=t.target;if(!(s instanceof HTMLElement))return;const a=s.closest("[data-action]");if(!a)return;switch(a.dataset.action??""){case"auth-login":await K();return;case"auth-logout":await pe();return;case"nav":await we(a.dataset.page);return;case"quotes-refresh":await h();return;case"history-refresh":await L();return;case"history-view-current":await $e();return;case"history-back":$();return;case"recall-save-quote":await He();return;case"history-save-quote":await Fe();return;case"history-delete-current":qe();return;case"history-select-all":Je();return;case"history-deselect-all":ze();return;case"quote-add":U("add");return;case"quote-import":De();return;case"library-filter-ownership":e.libraryOwnership=a.dataset.value??"all",e.quotesCursor=0,i();return;case"library-filter-tag":e.libraryTag=a.dataset.value&&a.dataset.value!=="all"?a.dataset.value:null,e.quotesCursor=0,i();return;case"library-clear-filters":e.libraryQuery="",e.libraryOwnership="all",e.libraryTag=null,e.quotesCursor=0,i();return;case"quote-select-all":Oe(a.dataset.context);return;case"quote-deselect-all":Be(a.dataset.context);return;case"quote-edit-current":Se(a.dataset.context);return;case"quote-delete-current":Ee(a.dataset.context);return;case"quote-share-current":await Qe(a.dataset.context);return;case"quote-inspect":Ie(a.dataset.context,Number(a.dataset.index??"0"));return;case"set-cursor":if(s.closest("input, button, label"))return;Me(a.dataset.context,Number(a.dataset.index??"0"));return;case"history-set-cursor":if(s.closest("input, button, label"))return;Ue(Number(a.dataset.index??"0"));return;case"history-open":await S(Number(a.dataset.id??"0"));return;case"share-toggle-payload":((o=e.overlay)==null?void 0:o.type)==="shareQuotes"&&(e.overlay.showPayload=!e.overlay.showPayload,i());return;case"import-toggle-payload":((l=e.overlay)==null?void 0:l.type)==="importQuotes"&&(e.overlay.showPayload=!e.overlay.showPayload,i());return;case"profile-save":await W();return;case"quote-editor-save":await J();return;case"quote-editor-refine":await z();return;case"quote-editor-apply-refined":ke();return;case"quote-editor-reject-refined":Ce();return;case"overlay-close":V();return;case"delete-confirm":await Pe();return;case"share-browse":await Le();return;case"share-save":await xe();return;case"import-browse":await Re();return;case"import-run":await Te();return;case"settings-fetch-models":await Ae();return;case"settings-save":await Y();return;case"settings-change-password":await ye();return;case"recall-run":await x();return;case"use-last-question":e.recallQuestion=e.recallLastQuestion,i();return;case"reuse-history-question":if(e.historyDetail)e.recallQuestion=e.historyDetail.Question;else{const u=Q()[0];u&&(e.recallQuestion=u.Question)}e.page="Recall",i();return;default:return}}function he(t){var r,o,l,u;const s=t.target;if(!(s instanceof HTMLInputElement||s instanceof HTMLTextAreaElement))return;switch(s.dataset.bind??""){case"auth-password":e.authPassword=s.value;return;case"auth-confirm-password":e.authConfirmPassword=s.value;return;case"recall-question":e.recallQuestion=s.value;return;case"library-query":e.libraryQuery=s.value,e.quotesCursor=0,i();return;case"profile-name":((r=e.overlay)==null?void 0:r.type)==="namePrompt"&&(e.overlay.name=s.value);return;case"quote-editor-content":((o=e.overlay)==null?void 0:o.type)==="quoteEditor"&&(e.overlay.content=s.value);return;case"share-path":((l=e.overlay)==null?void 0:l.type)==="shareQuotes"&&(e.overlay.path=s.value);return;case"import-path":((u=e.overlay)==null?void 0:u.type)==="importQuotes"&&(e.overlay.path=s.value);return;case"settings-host":e.settings.host=s.value;return;case"settings-port":e.settings.port=s.value;return;case"settings-api-key":e.settings.apiKey=s.value;return;case"settings-model-filter":e.settings.modelFilter=s.value,H(e.settings),i();return;case"settings-max-results":e.settings.maxResults=s.value;return;case"settings-min-relevance":e.settings.minRelevance=s.value;return;case"settings-theme":e.settings.theme=s.value,P(e.settings.theme);return;case"settings-web-port":e.settings.webPort=s.value;return;case"settings-password-current":e.passwordForm.current=s.value;return;case"settings-password-next":e.passwordForm.next=s.value;return;case"settings-password-confirm":e.passwordForm.confirm=s.value;return;default:return}}function fe(t){const s=t.target;if(!(s instanceof HTMLInputElement||s instanceof HTMLSelectElement))return;switch(s.dataset.bind??""){case"quote-selected":Ne(s.dataset.context,Number(s.dataset.id??"0"),s.checked);return;case"history-selected":We(Number(s.dataset.id??"0"),s.checked);return;case"settings-https":s instanceof HTMLInputElement&&(e.settings.https=s.checked);return;case"settings-model":e.settings.model=s.value;return;case"import-file":s instanceof HTMLInputElement&&ve(s);return;default:return}}async function me(t){const s=t.target;if(s instanceof HTMLFormElement)switch(t.preventDefault(),s.dataset.form){case"auth-login":await K();return;case"auth-setup":await submitAuthSetup();return;case"recall":await x();return;case"profile":await W();return;default:return}}async function ge(t){var a,r;const s=document.activeElement;if(t.key==="Escape"&&e.overlay&&e.overlay.type!=="namePrompt"){t.preventDefault(),V();return}if(t.ctrlKey&&t.key.toLowerCase()==="s"){if(((a=e.overlay)==null?void 0:a.type)==="quoteEditor"){t.preventDefault(),await J();return}!e.overlay&&e.page==="Settings"&&(t.preventDefault(),await Y());return}if(t.ctrlKey&&t.key.toLowerCase()==="r"&&((r=e.overlay)==null?void 0:r.type)==="quoteEditor"){t.preventDefault(),await z();return}t.key==="Enter"&&!t.shiftKey&&s instanceof HTMLInputElement&&s.dataset.bind==="recall-question"&&(t.preventDefault(),await x())}async function we(t){e.page=t,i(),t==="Quotes"&&await h(),t==="History"&&await L()}async function h(){e.quotesLoading=!0,e.quotesError="",i();try{const t=await c().ListQuotes();e.quotes=t,e.quotesCursor=b(e.quotesCursor,t),e.quotesSelected=se(e.quotesSelected,t),e.quotesError=""}catch(t){e.quotesError=d(t)}finally{e.quotesLoading=!1,i()}}async function L(){var t;e.historyLoading=!0,e.historyError="",e.historyStatus="",e.historyStatusIsError=!1,i();try{const s=await c().ListRecallHistory(),a=((t=e.historyDetail)==null?void 0:t.ID)??null;e.historyEntries=s,e.historyCursor=F(e.historyCursor,s),e.historySelected=yt(e.historySelected,s),a===null?$():s.some(r=>r.ID===a)?S(a,!0):$()}catch(s){e.historyError=d(s)}finally{e.historyLoading=!1,i()}}async function $e(){const t=Q()[0];t&&await S(t.ID)}async function S(t,s=!1){var a;if(!(!Number.isFinite(t)||t<=0)&&!(e.historyDetailLoading&&((a=e.historyDetail)==null?void 0:a.ID)===t)&&!(e.historyDetail&&e.historyDetail.ID===t&&!e.historyDetailError)){e.historyDetailLoading=!0,e.historyDetailError="",s||(e.historyStatus="",e.historyStatusIsError=!1),i();try{const r=await c().GetRecallHistory(t);e.historyDetail=r,e.historyQuoteCursor=b(e.historyQuoteCursor,r.Quotes),e.historyQuoteSelected=se(e.historyQuoteSelected,r.Quotes)}catch(r){e.historyDetailError=d(r)}finally{e.historyDetailLoading=!1,i()}}}function $(){e.historyDetail=null,e.historyDetailLoading=!1,e.historyDetailError="",e.historyQuoteCursor=0,e.historyQuoteSelected=new Set,i()}function U(t,s){e.overlay={type:"quoteEditor",mode:t,quoteId:(s==null?void 0:s.ID)??null,content:(s==null?void 0:s.Content)??"",busy:!1,status:"",isError:!1,previewOriginal:"",previewRefined:""},i()}function Se(t){const s=q(t)[0];s&&U("edit",s)}function Ee(t){const s=q(t).map(a=>a.ID);s.length!==0&&(e.overlay={type:"deleteQuotes",context:t,ids:s,busy:!1,status:"",isError:!1},i())}function qe(){const t=Q().map(s=>s.ID);t.length!==0&&(e.overlay={type:"deleteHistory",ids:t,busy:!1,status:"",isError:!1},i())}async function Qe(t){var a,r;const s=q(t);if(s.length!==0){e.overlay={type:"shareQuotes",context:t,ids:s.map(o=>o.ID),path:"",payload:"",showPayload:!1,busy:!0,status:"",isError:!1},i();try{const o=await c().PreviewQuoteExport(s.map(l=>l.ID));if(((a=e.overlay)==null?void 0:a.type)!=="shareQuotes")return;e.overlay.payload=o,e.overlay.busy=!1,e.overlay.status="Share payload ready. Save it to a file and transfer it manually.",e.overlay.isError=!1}catch(o){if(((r=e.overlay)==null?void 0:r.type)!=="shareQuotes")return;e.overlay.busy=!1,e.overlay.status=d(o),e.overlay.isError=!0}i()}}function De(){e.overlay={type:"importQuotes",path:"",payload:"",filename:"",showPayload:!1,busy:!1,status:"",isError:!1,result:null},i()}function Ie(t,s){const a=E(t),r=b(s,a),o=a[r];o&&(t==="quotes"?e.quotesCursor=r:t==="recall"?e.recallCursor=r:e.historyQuoteCursor=r,e.overlay={type:"quoteInspect",context:t,quote:o},i())}async function W(){var s,a;if(((s=e.overlay)==null?void 0:s.type)!=="namePrompt"||e.overlay.busy)return;const t=e.overlay.name.trim();if(!t){e.overlay.status="Please enter a name to continue.",e.overlay.isError=!0,i();return}e.overlay.busy=!0,e.overlay.status="Saving profile…",e.overlay.isError=!1,i();try{const r=await c().SaveUserProfile(t);e.bootstrap&&(e.bootstrap.profile=r,e.bootstrap.greeting=`Hi! ${r.DisplayName}`),e.overlay=null}catch(r){((a=e.overlay)==null?void 0:a.type)==="namePrompt"&&(e.overlay.busy=!1,e.overlay.status=d(r),e.overlay.isError=!0)}i()}async function J(){var s,a;if(((s=e.overlay)==null?void 0:s.type)!=="quoteEditor"||e.overlay.busy)return;const t=e.overlay.content.trim();if(!t){e.overlay.status="Nothing to save.",e.overlay.isError=!0,i();return}e.overlay.busy=!0,e.overlay.status="Refining draft...",e.overlay.isError=!1,i();try{const r=e.overlay.mode==="add"?await c().AddQuote(t):await c().UpdateQuote(e.overlay.quoteId??0,t);e.overlay=null,R(r),await h()}catch(r){((a=e.overlay)==null?void 0:a.type)==="quoteEditor"&&(e.overlay.busy=!1,e.overlay.status=d(r),e.overlay.isError=!0),i()}}async function z(){var s,a,r;if(((s=e.overlay)==null?void 0:s.type)!=="quoteEditor"||e.overlay.busy)return;const t=e.overlay.content.trim();if(!t){e.overlay.status="Nothing to refine.",e.overlay.isError=!0,i();return}e.overlay.busy=!0,e.overlay.status="",i();try{const o=await c().RefineQuoteDraft(t);if(((a=e.overlay)==null?void 0:a.type)!=="quoteEditor")return;e.overlay.busy=!1,e.overlay.previewOriginal=t,e.overlay.previewRefined=o,e.overlay.status="",e.overlay.isError=!1}catch(o){((r=e.overlay)==null?void 0:r.type)==="quoteEditor"&&(e.overlay.busy=!1,e.overlay.status=d(o),e.overlay.isError=!0)}i()}function ke(){var t;((t=e.overlay)==null?void 0:t.type)==="quoteEditor"&&(e.overlay.content=e.overlay.previewRefined,e.overlay.previewOriginal="",e.overlay.previewRefined="",e.overlay.status="Refined draft applied. Review it, then save.",e.overlay.isError=!1,i())}function Ce(){var t;((t=e.overlay)==null?void 0:t.type)==="quoteEditor"&&(e.overlay.previewOriginal="",e.overlay.previewRefined="",e.overlay.status="Refined draft discarded.",e.overlay.isError=!1,i())}async function Pe(){var t,s,a,r;if(((t=e.overlay)==null?void 0:t.type)==="deleteHistory"){if(e.overlay.busy)return;e.overlay.busy=!0,e.overlay.status="",i();try{await c().DeleteRecallHistory(e.overlay.ids),Ye(e.overlay.ids),e.overlay=null,await L()}catch(o){((s=e.overlay)==null?void 0:s.type)==="deleteHistory"&&(e.overlay.busy=!1,e.overlay.status=d(o),e.overlay.isError=!0),i()}return}if(!(((a=e.overlay)==null?void 0:a.type)!=="deleteQuotes"||e.overlay.busy)){e.overlay.busy=!0,e.overlay.status="",i();try{await c().DeleteQuotes(e.overlay.ids),Ke(e.overlay.ids),e.overlay=null,await h()}catch(o){((r=e.overlay)==null?void 0:r.type)==="deleteQuotes"&&(e.overlay.busy=!1,e.overlay.status=d(o),e.overlay.isError=!0),i()}}}async function Le(){var t,s,a;if(!(((t=e.overlay)==null?void 0:t.type)!=="shareQuotes"||e.overlay.busy)){if(v()){e.overlay.path="irecall-share.json",i();return}try{const r=await c().SelectQuoteExportFile();r&&((s=e.overlay)==null?void 0:s.type)==="shareQuotes"&&(e.overlay.path=r,i())}catch(r){((a=e.overlay)==null?void 0:a.type)==="shareQuotes"&&(e.overlay.status=d(r),e.overlay.isError=!0,i())}}}async function xe(){var s,a,r;if(((s=e.overlay)==null?void 0:s.type)!=="shareQuotes"||e.overlay.busy)return;if(v()){const o=e.overlay.path.trim()||"irecall-share.json";if(!e.overlay.payload.trim()){e.overlay.status="Export payload is not ready yet.",e.overlay.isError=!0,i();return}ht(o,e.overlay.payload),e.overlay.status=`Downloaded ${o}`,e.overlay.isError=!1,i();return}const t=e.overlay.path.trim();if(!t){e.overlay.status="Choose a file path for the export.",e.overlay.isError=!0,i();return}if(!e.overlay.payload.trim()){e.overlay.status="Export payload is not ready yet.",e.overlay.isError=!0,i();return}e.overlay.busy=!0,e.overlay.status="",i();try{await c().ExportQuotesToFile(e.overlay.ids,t),((a=e.overlay)==null?void 0:a.type)==="shareQuotes"&&(e.overlay.busy=!1,e.overlay.status=`Saved share payload to ${t}`,e.overlay.isError=!1,i())}catch(o){((r=e.overlay)==null?void 0:r.type)==="shareQuotes"&&(e.overlay.busy=!1,e.overlay.status=d(o),e.overlay.isError=!0,i())}}async function Re(){var t,s,a;if(!(((t=e.overlay)==null?void 0:t.type)!=="importQuotes"||e.overlay.busy)){if(v()){const r=document.querySelector('[data-bind="import-file"]');r==null||r.click();return}try{const r=await c().SelectQuoteImportFile();r&&((s=e.overlay)==null?void 0:s.type)==="importQuotes"&&(e.overlay.path=r,i())}catch(r){((a=e.overlay)==null?void 0:a.type)==="importQuotes"&&(e.overlay.status=d(r),e.overlay.isError=!0,i())}}}async function Te(){var a,r,o;if(((a=e.overlay)==null?void 0:a.type)!=="importQuotes"||e.overlay.busy)return;const t=e.overlay.path.trim(),s=e.overlay.payload.trim();if(v()){if(!s){e.overlay.status="Choose a file to import.",e.overlay.isError=!0,i();return}}else if(!t){e.overlay.status="Choose a file to import.",e.overlay.isError=!0,i();return}e.overlay.busy=!0,e.overlay.status="",e.overlay.result=null,i();try{const l=v()?await c().ImportQuotesPayload(s):await c().ImportQuotesFromFile(t);if(((r=e.overlay)==null?void 0:r.type)!=="importQuotes")return;e.overlay.busy=!1,e.overlay.result=l,e.overlay.status=`Imported quotes. inserted=${l.Inserted} updated=${l.Updated} duplicates=${l.Duplicates} stale=${l.Stale}`,e.overlay.isError=!1,await h()}catch(l){((o=e.overlay)==null?void 0:o.type)==="importQuotes"&&(e.overlay.busy=!1,e.overlay.status=d(l),e.overlay.isError=!0,i())}}async function x(){if(e.recallBusy)return;const t=e.recallQuestion.trim();if(!t){e.recallError="Ask a question first.",i();return}e.recallBusy=!0,e.recallError="",e.recallStatus="",e.recallStatusIsError=!1,e.recallLastQuestion=t,e.recallKeywords=[],e.recallQuotes=[],e.recallResponse="",e.recallCursor=0,e.recallSelected=new Set,i();try{const s=await c().RunRecall(t);e.recallKeywords=s.keywords,e.recallQuotes=s.quotes,e.recallResponse=s.response,e.recallLastQuestion=s.question||t,e.recallCursor=0,e.recallSelected=new Set,e.recallQuestion=""}catch(s){e.recallError=d(s)}finally{e.recallBusy=!1,i()}}async function He(){const t=e.recallLastQuestion.trim(),s=e.recallResponse.trim();if(!t||!s){e.recallStatus="Run a recall first before saving it as a quote.",e.recallStatusIsError=!0,i();return}try{const a=await c().SaveRecallAsQuote(t,s,e.recallKeywords);R(a),await h(),e.recallStatus="Saved recall as quote.",e.recallStatusIsError=!1,G("Saved the current grounded answer as a quote.")}catch(a){e.recallStatus=d(a),e.recallStatusIsError=!0}i()}async function Fe(){const t=e.historyDetail;if(t){try{const s=await c().SaveRecallAsQuote(t.Question,t.Response,[]);R(s),await h(),e.historyStatus="Saved history entry as quote.",e.historyStatusIsError=!1,G("Saved the selected activity session as a quote.")}catch(s){e.historyStatus=d(s),e.historyStatusIsError=!0}i()}}async function Ae(){if(e.settingsBusy)return;let t;try{t=ee(e.settings)}catch(s){e.settingsStatus=d(s),e.settingsIsError=!0,i();return}e.settingsBusy=!0,e.settingsStatus="",i();try{const s=await c().FetchModels(t);e.settings.models=s,H(e.settings),e.settingsStatus=s.length>0?`Fetched ${s.length} models.`:"No models returned.",e.settingsIsError=!1}catch(s){e.settingsStatus=d(s),e.settingsIsError=!0}finally{e.settingsBusy=!1,i()}}async function Y(){var s;if(e.settingsBusy)return;let t;try{t=pt(e.settings)}catch(a){e.settingsStatus=d(a),e.settingsIsError=!0,i();return}e.settingsBusy=!0,e.settingsStatus="",i();try{const a=await c().SaveSettings(t);e.settings=_(a,e.settings.models),P(e.settings.theme),e.bootstrap&&(e.bootstrap.settings=a);const r=((s=e.auth)==null?void 0:s.runtime)==="web"&&e.auth.currentPort>0&&e.auth.currentPort!==a.Web.Port;e.settingsStatus=r?"Saved. Restart the web server to apply the new port.":"Saved.",e.settingsIsError=!1}catch(a){e.settingsStatus=d(a),e.settingsIsError=!0}finally{e.settingsBusy=!1,i()}}function V(){e.overlay&&e.overlay.type!=="namePrompt"&&("busy"in e.overlay&&e.overlay.busy||(e.overlay=null,i()))}function G(t,s=!1){e.toast={message:t,isError:s},g!==null&&window.clearTimeout(g),g=window.setTimeout(()=>{e.toast=null,g=null,i()},2600)}function E(t){var s;return t==="quotes"?X():t==="recall"?e.recallQuotes:((s=e.historyDetail)==null?void 0:s.Quotes)??[]}function Me(t,s){var r;const a=E(t);if(t==="quotes")e.quotesCursor=b(s,a);else if(t==="recall")e.recallCursor=b(s,e.recallQuotes);else{const o=((r=e.historyDetail)==null?void 0:r.Quotes)??[];e.historyQuoteCursor=b(s,o)}i()}function Ne(t,s,a){const r=t==="quotes"?e.quotesSelected:t==="recall"?e.recallSelected:e.historyQuoteSelected;a?r.add(s):r.delete(s)}function Oe(t){const s=E(t),a=new Set(s.map(r=>r.ID));t==="quotes"?e.quotesSelected=a:t==="recall"?e.recallSelected=a:e.historyQuoteSelected=a,i()}function Be(t){t==="quotes"?e.quotesSelected=new Set:t==="recall"?e.recallSelected=new Set:e.historyQuoteSelected=new Set,i()}function q(t){const s=E(t),a=t==="quotes"?e.quotesCursor:t==="recall"?e.recallCursor:e.historyQuoteCursor,r=b(a,s),o=t==="quotes"?e.quotesSelected:t==="recall"?e.recallSelected:e.historyQuoteSelected,l=s.filter(u=>o.has(u.ID));return l.length>0?l:s[r]?[s[r]]:[]}function X(){const t=e.libraryQuery.trim().toLowerCase();return e.quotes.filter(s=>e.libraryOwnership==="owned"&&!s.IsOwnedByMe||e.libraryOwnership==="imported"&&s.IsOwnedByMe||e.libraryTag&&!s.Tags.some(r=>{var o;return r.toLowerCase()===((o=e.libraryTag)==null?void 0:o.toLowerCase())})?!1:t?[s.Content,s.AuthorName,s.SourceName,...s.Tags].join(" ").toLowerCase().includes(t):!0)}function je(t=12){const s=new Map;for(const a of e.quotes)for(const r of a.Tags){const o=r.trim();o&&s.set(o,(s.get(o)??0)+1)}return[...s.entries()].sort((a,r)=>r[1]-a[1]||a[0].localeCompare(r[0])).slice(0,t).map(([a])=>a)}function R(t){e.quotes=D(e.quotes,t),e.recallQuotes=D(e.recallQuotes,t),e.historyDetail&&(e.historyDetail={...e.historyDetail,Quotes:D(e.historyDetail.Quotes,t)}),i()}function Ke(t){var a;const s=new Set(t);e.quotes=e.quotes.filter(r=>!s.has(r.ID)),e.recallQuotes=e.recallQuotes.filter(r=>!s.has(r.ID)),e.historyDetail&&(e.historyDetail={...e.historyDetail,Quotes:e.historyDetail.Quotes.filter(r=>!s.has(r.ID))}),e.quotesSelected=new Set([...e.quotesSelected].filter(r=>!s.has(r))),e.recallSelected=new Set([...e.recallSelected].filter(r=>!s.has(r))),e.historyQuoteSelected=new Set([...e.historyQuoteSelected].filter(r=>!s.has(r))),e.quotesCursor=b(e.quotesCursor,e.quotes),e.recallCursor=b(e.recallCursor,e.recallQuotes),e.historyQuoteCursor=b(e.historyQuoteCursor,((a=e.historyDetail)==null?void 0:a.Quotes)??[]),i()}function Ue(t){e.historyCursor=F(t,e.historyEntries);const s=e.historyEntries[e.historyCursor];s&&S(s.ID,!0),i()}function We(t,s){s?e.historySelected.add(t):e.historySelected.delete(t)}function Je(){e.historySelected=new Set(e.historyEntries.map(t=>t.ID)),i()}function ze(){e.historySelected=new Set,i()}function Q(){const t=e.historyEntries.filter(s=>e.historySelected.has(s.ID));return t.length>0?t:e.historyEntries[e.historyCursor]?[e.historyEntries[e.historyCursor]]:[]}function Ye(t){const s=new Set(t);if(e.historyEntries=e.historyEntries.filter(a=>!s.has(a.ID)),e.historySelected=new Set([...e.historySelected].filter(a=>!s.has(a))),e.historyCursor=F(e.historyCursor,e.historyEntries),e.historyDetail&&s.has(e.historyDetail.ID)){$();return}i()}function i(){if(!m)return;const t=Ve();m.innerHTML=Xe(),Ge(t)}function Ve(){const t=document.activeElement;if(!(t instanceof HTMLInputElement||t instanceof HTMLTextAreaElement||t instanceof HTMLSelectElement))return null;const s=t.dataset.bind;return s?{selector:`[data-bind="${s}"]`,selectionStart:t instanceof HTMLInputElement||t instanceof HTMLTextAreaElement?t.selectionStart:null,selectionEnd:t instanceof HTMLInputElement||t instanceof HTMLTextAreaElement?t.selectionEnd:null}:null}function Ge(t){if(!m||!t)return;const s=m.querySelector(t.selector);(s instanceof HTMLInputElement||s instanceof HTMLTextAreaElement||s instanceof HTMLSelectElement)&&(s.focus({preventScroll:!0}),(s instanceof HTMLInputElement||s instanceof HTMLTextAreaElement)&&t.selectionStart!==null&&t.selectionEnd!==null&&s.setSelectionRange(t.selectionStart,t.selectionEnd))}function Xe(){var s,a,r,o,l;if(!e.authChecked)return`
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
    `;if(((s=e.auth)==null?void 0:s.runtime)==="web"&&!e.auth.authenticated)return Ze();if(!e.bootstrapped)return`
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
          <div class="brand">${n(((o=e.bootstrap)==null?void 0:o.productName)??"iRecall")}</div>
          <div class="muted subtle">${v()?"Local-first knowledge workspace for the web":"Local-first knowledge workspace for desktop"}</div>
        </div>
        <div class="titlebar-right">
          <div class="greeting">${n(t)}</div>
          ${((l=e.auth)==null?void 0:l.runtime)==="web"?'<button class="button" data-action="auth-logout" type="button">Logout</button>':""}
          <nav class="tabs" aria-label="Primary">
            ${le.map(u=>`
                  <button
                    class="tab${e.page===u?" active":""}"
                    data-action="nav"
                    data-page="${u}"
                    type="button"
                  >${et(u)}</button>
                `).join("")}
          </nav>
        </div>
      </header>

      <main class="layout">
        ${_e()}
      </main>

      ${rt()}
      ${e.overlay?lt(e.overlay):""}
      ${e.toast?nt(e.toast):""}
    </div>
  `}function Ze(){var o;const t=!((o=e.auth)!=null&&o.passwordConfigured),s="auth-login";return`
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
  `}function _e(){switch(e.page){case"Recall":return tt();case"Quotes":return st();case"History":return at();case"Settings":return it()}}function et(t){switch(t){case"Recall":return"Ask";case"Quotes":return"Quotes";case"History":return"History";case"Settings":return"Settings"}}function tt(){const t=!e.recallResponse.trim(),s=e.recallResponse.trim()?n(e.recallResponse):'<span class="muted">Grounded response will appear here.</span>',a=e.recallKeywords.length>0?e.recallKeywords.map(r=>`<span class="keyword-chip">${n(r)}</span>`).join(""):'<span class="muted">Keywords: —</span>';return`
    <section class="page page-recall">
      <div class="panel page-panel">
        <div class="page-hero">
          <div>
            <div class="eyebrow">Ask</div>
            <div class="page-title">Question, references, then answer</div>
            <div class="muted page-copy">Ask once, inspect the retrieved quotes, then read the grounded response.</div>
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
                  ${e.recallBusy?"Working…":"Ask"}
                </button>
                ${e.recallLastQuestion.trim()?'<button class="button" data-action="use-last-question" type="button">Use previous question</button>':""}
              </div>
            </form>
          </section>

          <section class="panel subpanel">
            <div class="subpanel-header">
              <div>
                <div class="section-title">2. Reference quotes</div>
                <div class="muted">${e.recallBusy?"Searching your quotes for relevant evidence…":`${e.recallQuotes.length} retrieved quotes. Open one to inspect the full note.`}</div>
              </div>
            </div>
            <div class="keyword-row">
              <span class="muted">Keywords</span>
              <div class="keyword-list">${a}</div>
            </div>
            ${T("recall",e.recallQuotes,e.recallCursor,e.recallSelected,!1)}
          </section>

          <section class="panel subpanel">
            <div class="subpanel-header">
              <div>
                <div class="section-title">3. Response</div>
                <div class="muted">${e.recallBusy?"Writing a grounded response from the retrieved evidence…":"The response is generated from the current question and reference set."}</div>
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
              <pre class="response-box">${s}</pre>
            </div>
          </section>
        </div>

        ${e.recallError?`<div class="status status-error">${n(e.recallError)}</div>`:""}
        ${e.recallStatus?`<div class="status ${e.recallStatusIsError?"status-error":"status-ok"}">${n(e.recallStatus)}</div>`:""}
      </div>
    </section>
  `}function st(){const t=X(),s=b(e.quotesCursor,t),a=q("quotes"),r=a[0]??null,o=t.filter(y=>e.quotesSelected.has(y.ID)).length,l=je();return`
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
          ${[`${e.quotes.length} total`,`${e.quotes.filter(y=>y.IsOwnedByMe).length} authored here`,`${e.quotes.filter(y=>!y.IsOwnedByMe).length} imported`].map(y=>`<span class="meta-pill">${n(y)}</span>`).join("")}
          <span class="meta-pill meta-pill-accent">${a.length>0?`${a.length} selected`:r?"1 in focus":"Nothing selected"}</span>
        </div>

        <div class="workspace workspace-library">
          <section class="panel subpanel filter-panel">
            <div class="section-title">Browse</div>
            <label class="field">
              <span>Search quotes</span>
              <input
                class="text-input"
                data-bind="library-query"
                value="${p(e.libraryQuery)}"
                placeholder="Search content, source, author, or tags"
              />
            </label>
            <div class="field">
              <span>Ownership</span>
              <div class="chip-row">
                ${["all","owned","imported"].map(y=>`
                      <button
                        class="filter-chip${e.libraryOwnership===y?" is-active":""}"
                        data-action="library-filter-ownership"
                        data-value="${y}"
                        type="button"
                      >${y==="all"?"All":y==="owned"?"Authored here":"Imported"}</button>
                    `).join("")}
              </div>
            </div>
            <div class="field">
              <span>Popular tags</span>
              <div class="chip-row">
                <button class="filter-chip${e.libraryTag===null?" is-active":""}" data-action="library-filter-tag" data-value="all" type="button">All tags</button>
                ${l.map(y=>`
                      <button
                        class="filter-chip${e.libraryTag===y?" is-active":""}"
                        data-action="library-filter-tag"
                        data-value="${p(y)}"
                        type="button"
                      >${n(y)}</button>
                    `).join("")}
              </div>
            </div>
            <div class="toolbar toolbar-stack">
              <button class="button" data-action="library-clear-filters" type="button" ${!e.libraryQuery&&e.libraryOwnership==="all"&&e.libraryTag===null?"disabled":""}>Clear filters</button>
              <button class="button" data-action="quote-select-all" data-context="quotes" type="button" ${t.length===0?"disabled":""}>Select results</button>
              <button class="button" data-action="quote-deselect-all" data-context="quotes" type="button" ${o===0?"disabled":""}>Clear selection</button>
            </div>
          </section>

          <section class="panel subpanel">
            <div class="subpanel-header">
              <div>
                <div class="section-title">Quote list</div>
                <div class="muted">${t.length} matching quotes. Choose one to inspect the full note, provenance, and actions.</div>
              </div>
            </div>
            ${e.quotesLoading?'<div class="empty-state">Loading quotes…</div>':e.quotesError?`<div class="status status-error">${n(e.quotesError)}</div>`:T("quotes",t,s,e.quotesSelected,!0)}
          </section>

          <section class="panel subpanel detail-panel">
            <div class="subpanel-header">
              <div>
                <div class="section-title">Quote details</div>
                <div class="muted">${r?"Read the full note and manage it from here.":"Select a quote to inspect its full text and metadata."}</div>
              </div>
              <div class="toolbar toolbar-quiet">
                <button class="button" data-action="quote-edit-current" data-context="quotes" type="button" ${r?"":"disabled"}>Edit</button>
                <button class="button" data-action="quote-share-current" data-context="quotes" type="button" ${r?"":"disabled"}>Share</button>
                <button class="button button-danger" data-action="quote-delete-current" data-context="quotes" type="button" ${r?"":"disabled"}>Delete</button>
              </div>
            </div>
            ${Z(r,"quotes")}
          </section>
        </div>
      </div>
    </section>
  `}function at(){const t=Q();return`
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
            <button class="button" data-action="reuse-history-question" type="button" ${!!(e.historyDetail??t[0])?"":"disabled"}>Ask again</button>
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
            ${e.historyLoading?'<div class="empty-state">Loading history…</div>':e.historyError?`<div class="status status-error">${n(e.historyError)}</div>`:e.historyEntries.length===0?'<div class="empty-state">No recall history yet. Ask a question from the Ask page to create your first grounded session.</div>':ot()}
          </section>
        </div>
      </div>
    </section>
  `}function rt(){if(!e.historyDetailLoading&&!e.historyDetail&&!e.historyDetailError)return"";const t=e.historyDetail?e.historyEntries.find(r=>{var o;return r.ID===((o=e.historyDetail)==null?void 0:o.ID)})??e.historyDetail:e.historyEntries[e.historyCursor]??null,s=e.historyDetail&&t&&e.historyDetail.ID===t.ID?e.historyDetail:null,a=(s==null?void 0:s.Quotes)??[];return`
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
            <button class="button" data-action="reuse-history-question" type="button" ${t?"":"disabled"}>Ask again</button>
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
  `}function ot(){return`
    <div class="history-list">
      ${e.historyEntries.map((t,s)=>{const a=s===e.historyCursor,r=C(t.Response,156);return`
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
                  <span class="quote-version">${n(vt(t.CreatedAt))}</span>
                </div>
              </div>
              <div class="quote-content">${n(C(t.Question,132))}</div>
                <div class="quote-meta"><span class="muted">Response preview</span><span>${n(r||"(empty response)")}</span></div>
            </article>
          `}).join("")}
    </div>
  `}function it(){var o,l;const t=te(e.settings),s=(o=e.bootstrap)==null?void 0:o.paths,a=(l=e.auth)==null?void 0:l.currentPort,r=e.settings.models.length>0&&t.length>0?`
        <select class="select-input" data-bind="settings-model">
          ${t.map(u=>`
                <option value="${p(u)}"${u===e.settings.model?" selected":""}>${n(u)}</option>
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
            <div class="section-title">Personalization</div>
            <label class="field">
              <span>Theme</span>
              <select class="select-input" data-bind="settings-theme">
                ${ne().map(u=>`
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
            ${e.passwordForm.status?`<div class="status ${e.passwordForm.isError?"status-error":"status-ok"}">${n(e.passwordForm.status)}</div>`:""}
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
  `}function T(t,s,a,r,o){return s.length===0?`<div class="empty-state">${t==="quotes"?"No quotes yet. Add one or import a shared payload.":"No reference quotes for this question yet."}</div>`:`
    <div class="quote-list">
      ${s.map((l,u)=>{const y=u===a,f=t!=="quotes",A=!l.IsOwnedByMe&&l.SourceName?`<span class="meta-accent">${n(l.SourceName)}</span>`:`<span>${n(l.AuthorName||"You")}</span>`,ie=o?`
              <div class="quote-meta">
                <span class="muted">Tags</span>
                <span>${l.Tags.length>0?n(bt(l.Tags,4)):"(none)"}</span>
              </div>
            `:"";return`
            <article class="quote-card${y?" is-current":""}${f?" quote-card-minimal":""}" data-action="${f?"quote-inspect":"set-cursor"}" data-context="${t}" data-index="${u}">
              <div class="quote-topline">
                ${f?`<div class="quote-topline-meta">
                        <span class="quote-badge">${l.IsOwnedByMe?"Owned":"Imported"}</span>
                        <span class="quote-version">${n(k(l.UpdatedAt))}</span>
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
                  ${f?`<span class="quote-source-inline">${A}</span>`:`<span class="quote-version">${n(k(l.UpdatedAt))}</span>
                  <span class="quote-badge">${l.IsOwnedByMe?"Owned":"Imported"}</span>`}
                </div>
              </div>
              <div class="quote-content">${n(C(l.Content,t==="quotes"?160:136))}</div>
              ${f?`<div class="quote-actions-inline"><button class="button button-subtle" data-action="quote-inspect" data-context="${t}" data-index="${u}" type="button">Details</button></div>`:`<div class="quote-meta"><span class="muted">${!l.IsOwnedByMe&&l.SourceName?"Imported from":"Author"}</span> ${A}</div>`}
              ${ie}
            </article>
          `}).join("")}
    </div>
  `}function Z(t,s){return t?`
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
          <span>${n(k(t.UpdatedAt))}</span>
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
      </div>
    </div>
  `:'<div class="empty-state">Select a quote to inspect the full note, provenance, and available actions.</div>'}function k(t){const s=new Date(t);return Number.isNaN(s.getTime())?t:s.toLocaleDateString(void 0,{month:"short",day:"numeric",year:"numeric"})}function O(t){return!Number.isFinite(t)||t<=0?"0 B":t<1024?`${Math.round(t)} B`:t<1024*1024?`${(t/1024).toFixed(1)} KB`:`${(t/(1024*1024)).toFixed(1)} MB`}function nt(t){return`
    <div class="toast-stack" role="status" aria-live="polite">
      <div class="toast${t.isError?" is-error":""}">${n(t.message)}</div>
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
              ${t.previewRefined?"Compare the current draft with the suggested rewrite before applying it.":"Tags are regenerated automatically by the shared core logic."}
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
              ${B(t.context,t.ids).map((a,r)=>`<div class="summary-item">[${r+1}] ${n(w(a.Content,140))}</div>`).join("")}
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
              ${ut(t.ids).map((a,r)=>`<div class="summary-item">[${r+1}] ${n(w(a.Question,140))}</div>`).join("")}
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
      `;case"shareQuotes":const s=B(t.context,t.ids);return`
        <div class="overlay-backdrop overlay-backdrop-side">
          <div class="modal modal-side">
            <div class="modal-title">Share Quotes</div>
            <div class="modal-copy">Export a portable share file. The file summary comes first; raw JSON is available only if you need to inspect it.</div>
            <div class="summary-list">
              ${s.map((a,r)=>`<div class="summary-item">[${r+1}] v${a.Version} ${n(w(a.Content,120))}</div>`).join("")}
            </div>
            <div class="result-grid">
              <div><span class="muted">Quotes:</span> ${s.length}</div>
              <div><span class="muted">Payload size:</span> ${O(t.payload.length)}</div>
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
                    <div><span class="muted">Payload size:</span> ${O(t.payload.length)}</div>
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
            ${Z(t.quote,t.context)}
            <div class="modal-actions">
              <button class="button" data-action="overlay-close" type="button">Close</button>
            </div>
          </div>
        </div>
      `;case"notice":return""}}function B(t,s){var o;const a=t==="quotes"?e.quotes:t==="recall"?e.recallQuotes:((o=e.historyDetail)==null?void 0:o.Quotes)??[],r=new Set(s);return a.filter(l=>r.has(l.ID))}function ut(t){const s=new Set(t);return e.historyEntries.filter(a=>s.has(a.ID))}function ct(t){return _(t.settings,[])}function _(t,s){var r;const a={host:t.Provider.Host,port:String(t.Provider.Port),https:t.Provider.HTTPS,apiKey:t.Provider.APIKey,modelFilter:"",model:t.Provider.Model,maxResults:String(t.Search.MaxResults),minRelevance:String(t.Search.MinRelevance),theme:t.Theme||"violet",webPort:String(((r=t.Web)==null?void 0:r.Port)??9527),models:s};return H(a),a}function dt(){return{host:"",port:"11434",https:!1,apiKey:"",modelFilter:"",model:"",maxResults:"5",minRelevance:"0",theme:"violet",webPort:"9527",models:[]}}function ee(t){const s=Number.parseInt(t.port.trim(),10);if(!Number.isInteger(s)||s<1||s>65535)throw new Error("Port must be a number between 1 and 65535.");return{Host:t.host.trim(),Port:s,HTTPS:t.https,APIKey:t.apiKey,Model:t.model}}function pt(t){const s=ee(t),a=Number.parseInt(t.maxResults.trim(),10),r=Number.parseInt(t.webPort.trim(),10);if(!Number.isInteger(a)||a<1||a>20)throw new Error("Max ref quotes must be between 1 and 20.");if(!Number.isInteger(r)||r<1||r>65535)throw new Error("Web port must be a number between 1 and 65535.");const o=Number.parseFloat(t.minRelevance.trim());if(Number.isNaN(o))throw new Error("Min relevance must be a decimal number.");if(o<0||o>1)throw new Error("Min relevance must be between 0.0 and 1.0.");return{Provider:s,Search:{MaxResults:a,MinRelevance:o},Theme:t.theme,Web:{Port:r}}}function te(t){const s=t.modelFilter.trim().toLowerCase();return s?t.models.filter(a=>a.toLowerCase().includes(s)):t.models}function H(t){if(t.models.length===0)return;const s=te(t);s.length!==0&&(s.includes(t.model)||(t.model=s[0]))}function D(t,s){return t.map(a=>a.ID===s.ID?s:a)}function b(t,s){return s.length===0?0:Math.min(Math.max(t,0),s.length-1)}function se(t,s){const a=new Set(s.map(r=>r.ID));return new Set([...t].filter(r=>a.has(r)))}function F(t,s){return s.length===0?0:Math.min(Math.max(t,0),s.length-1)}function yt(t,s){const a=new Set(s.map(r=>r.ID));return new Set([...t].filter(r=>a.has(r)))}function n(t){return t.replaceAll("&","&amp;").replaceAll("<","&lt;").replaceAll(">","&gt;").replaceAll('"',"&quot;").replaceAll("'","&#39;")}function p(t){return n(t)}function w(t,s){const a=t.replace(/\s+/g," ").trim();return a.length<=s?a:`${a.slice(0,s-1).trimEnd()}…`}function vt(t){const s=new Date(t);return Number.isNaN(s.getTime())?t:s.toLocaleString()}function bt(t,s){return t.length===0?"":t.length<=s?t.join(" · "):`${t.slice(0,s).join(" · ")} · +${t.length-s} more`}function C(t,s){return w(t,Math.max(8,s))}function v(){var t;return((t=e.auth)==null?void 0:t.runtime)==="web"}function ht(t,s){const a=new Blob([s],{type:"application/json;charset=utf-8"}),r=URL.createObjectURL(a),o=document.createElement("a");o.href=r,o.download=t,document.body.appendChild(o),o.click(),o.remove(),URL.revokeObjectURL(r)}function d(t){return t instanceof Error?t.message:String(t)}function ae(){var s,a,r;const t=[(s=window.go)==null?void 0:s.backend,(a=window.go)==null?void 0:a.app,(r=window.go)==null?void 0:r.main];for(const o of t)if(o!=null&&o.App)return o.App;return null}async function re(t=3e3){const s=Date.now();for(;;){const a=ae();if(a)return a;if(Date.now()-s>=t)throw new Error("Wails backend bridge is unavailable.");await new Promise(r=>window.setTimeout(r,25))}}function c(){const t=ae();if(!t)throw new Error("Wails backend bridge is unavailable.");return t}const oe=document.querySelector("#app");if(!oe)throw new Error("Missing #app root");ue(oe);

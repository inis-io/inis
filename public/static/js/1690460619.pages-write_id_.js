import{X as e,v as a,r as t,L as l,h as s,K as n,u as o,p as i,m as c,i as u,M as r,a as m,B as d,o as p,c as v,b as g,C as f,e as b,d as h,Y as y,w,k as V,T as x,N as _,F as j,g as k,f as z,j as A,t as L,P as U}from"./1690460619.index.js";import{_ as C,a as E}from"./1690460619.i-tinymce.js";const H={class:"container-fluid container-box"},I={class:"row"},M={class:"col-lg-9"},q={class:"card mb-2"},T={class:"card-body custom",style:{"min-height":"485px"}},P={class:"col-lg-3 custom",id:"page-header-title"},R={class:"card mb-2"},S={class:"card-body px-2 py-0"},$={class:"form-group mb-3"},B=g("span",{class:"ms-1 required"},"标题：",-1),N={class:"form-group"},O={class:"form-label"},D=g("span",{class:"ms-1"},"标签：",-1),F={class:"form-group mb-3"},J=g("span",{class:"ms-1 required"},"唯一键：",-1),K={class:"form-group mb-3"},X=g("span",{class:"ms-1"},"备注：",-1),Y={class:"card"},G={class:"card-body px-2 py-0"},Q={class:"form-group mb-3"},W={class:"form-label"},Z=g("span",{class:"ms-1"},"允许评论：",-1),ee={class:"font-13"},ae={class:"text-muted float-end"},te={class:"form-group mb-3"},le={class:"form-label"},se=g("span",{class:"ms-1"},"显示评论：",-1),ne={class:"font-13"},oe={class:"text-muted float-end"},ie={class:"inis-save"},ce={__name:"pages-write[id]",setup(ce){const{ctx:ue,proxy:re}=k(),me=e(),de=a(),pe=t({item:{id:null,active:["1"],tags:[],menu:{...l,menuList:[{label:"保存",icon:'<svg class="icon" viewBox="0 0 1024 1024" xmlns="http://www.w3.org/2000/svg" width="12" height="12">\n                    <path d="M777.216 0a106.496 106.496 0 0 1 77.824 48.128A107.52 107.52 0 0 1 870.4 102.4v153.6h68.096a102.4 102.4 0 0 1 73.216 51.2 120.832 120.832 0 0 1 12.288 39.936v583.68a113.664 113.664 0 0 1-26.112 56.832 106.496 106.496 0 0 1-65.536 34.304H91.136a113.152 113.152 0 0 1-55.808-25.088A106.496 106.496 0 0 1 0 932.352V349.184a112.64 112.64 0 0 1 15.36-45.056A102.4 102.4 0 0 1 92.16 256H153.6V153.6 88.576a96.256 96.256 0 0 1 20.48-51.2A104.96 104.96 0 0 1 247.296 0z m140.288 540.16H83.456v370.688a35.84 35.84 0 0 0 0 17.92 24.576 24.576 0 0 0 13.824 10.752 17.408 17.408 0 0 0 10.24 0h816.64a21.504 21.504 0 0 0 17.92-24.576v-374.784a51.2 51.2 0 0 0-24.576 0z m-204.8 48.64a40.448 40.448 0 0 1 37.888 39.936 44.544 44.544 0 0 1-17.408 34.816l-256 220.16a41.984 41.984 0 0 1-44.032 6.656A72.192 72.192 0 0 1 409.6 870.4l-115.2-132.608a41.984 41.984 0 0 1-13.312-29.184 38.912 38.912 0 0 1 26.112-38.912 39.936 39.936 0 0 1 47.104 10.24c34.304 37.888 67.584 77.312 102.4 115.2 76.288-65.024 153.6-131.072 227.84-196.608a40.96 40.96 0 0 1 28.16-9.728zM870.4 337.92V460.8a98.304 98.304 0 0 0 34.816 2.56 90.624 90.624 0 0 0 35.328-2.56V370.688a70.144 70.144 0 0 0 0-17.92 21.504 21.504 0 0 0-15.36-14.848c-7.68-5.12-17.92 0-26.624 0s-18.432-4.608-28.16 0zM102.4 338.432a19.968 19.968 0 0 0-20.48 19.968v102.4a32.768 32.768 0 0 0 20.48 2.56h54.784V363.52a51.2 51.2 0 0 0 0-23.04c-20.992-2.56-38.912-3.072-54.784-2.048z m665.6-256H260.096a21.504 21.504 0 0 0-24.064 16.384c-3.072 2.56-2.56 6.656-3.072 10.24v325.12A59.392 59.392 0 0 0 235.52 460.8a28.16 28.16 0 0 0 8.704 2.56h512a69.12 69.12 0 0 0 30.72-2.56 28.16 28.16 0 0 0 0-8.704V108.544a17.92 17.92 0 0 0 0-10.24 22.528 22.528 0 0 0-23.04-16.384zM358.4 204.8h112.64a41.472 41.472 0 0 1 39.424 30.208 39.936 39.936 0 0 1-10.752 39.424 42.496 42.496 0 0 1-30.72 12.288H354.304a58.88 58.88 0 0 1-23.04-3.584 40.448 40.448 0 0 1-24.064-39.424 39.936 39.936 0 0 1 33.28-38.912A82.944 82.944 0 0 1 358.4 204.8z" fill="rgb(var(--assist-color))"></path>\n                </svg>',fn:()=>ve.save()}]}},struct:{content:"",editor:null,json:{comment:{allow:0,show:0}}},select:{tags:[],comment:{allow:[{value:0,label:"继承父级（推荐）"},{value:1,label:"允许"},{value:2,label:"禁止"}],show:[{value:0,label:"继承父级（推荐）"},{value:1,label:"显示"},{value:2,label:"隐藏"}]}}});s((async()=>{document.querySelectorAll(".container-fluid").forEach((e=>{e.classList.remove("container-fluid"),e.classList.add("container-xxl")})),await ve.init(),pe.item.menu.menuList.push({line:!0},...await n())}));const ve={init:async()=>{var e;let a=null==(e=me.params)?void 0:e.id;o.is.empty(a)||(pe.item.id=parseInt(a)),await ve.getTags(),o.is.empty(pe.item.id)?await ve.getConfig():await ve.getPage(pe.item.id)},getTags:async()=>{const{code:e,data:a}=await i.get("/api/tags/column",{field:"id,name"});200===e&&(pe.select.tags=a.map((e=>({value:e.id,label:e.name}))))},getConfig:async()=>{var e;const{code:a,data:t}=await i.get("/api/config/one",{key:"ARTICLE"});200===a&&(o.in.array(null==(e=null==t?void 0:t.json)?void 0:e.editor,["tinymce","vditor"])?pe.struct.editor=t.json.editor:pe.struct.editor="tinymce")},getPage:async(e=null)=>{const{code:a,msg:t,data:l}=await i.get("/api/pages/one",{id:e});if(200!==a)return await de.push({path:"/admin/pages/write"}),c.error(t),c.warn("已为您跳转到页面撰写页！");if(pe.struct={...l,json:Object.assign({},l.json,pe.struct.json)},o.is.empty(l.covers)||(pe.item.fileList=l.covers.split(",").map((e=>({name:e.replace(/.*\//,""),url:e})))),!o.is.empty(l.group)){let e=l.group.split("|").filter((e=>!o.is.empty(e))).map((e=>parseInt(e)));pe.item.group=ve.tree.parse(pe.backup.group,e)}o.is.empty(l.tags)||(pe.item.tags=l.tags.split("|").filter((e=>!o.is.empty(e))).map((e=>parseInt(e))))},save:async()=>{var e,a,t,l,s,n;"vditor"===pe.struct.editor&&(pe.struct.content=re.$refs.vditor.getValue());let u=(null==(l=null==(t=null==(a=null==(e=pe.struct)?void 0:e.content)?void 0:a.replace(/<[^>]+>/g,""))?void 0:t.replace(/\n/g,""))?void 0:l.length)||0;switch(u){case 0:return c.warn("你这页面一个字都没写，糊弄谁呢？");case 1:return c.warn("真就只写一个字呗？");default:if(u<10)return c.warn("你这太水了，10个字都不到。")}if(o.is.empty(null==(s=pe.struct)?void 0:s.title))return c.warn("你可能忘记写标题了");if(o.is.empty(null==(n=pe.struct)?void 0:n.key))return c.warn("唯一键不能为空");pe.struct.tags=o.is.empty(pe.item.tags)?"":`|${pe.item.tags.join("|")}|`;const{code:r,msg:m,data:d}=await i.post("/api/pages/save",{...pe.struct,json:JSON.stringify(pe.struct.json)});if(200!==r)return c.error(m);c.success(m),pe.item.id=d.id,pe.struct.id=d.id,await de.push({path:"/admin/pages/write/"+parseInt(d.id)})},change:{tags:e=>{o.is.empty(e)||e.forEach((async(e,a)=>{if("string"==typeof e){const{code:t,msg:l,data:s}=await i.post("/api/tags/save",{name:e});if(200!==t)return c.error("添加标签失败："+l),pe.item.tags.splice(a,1);pe.item.tags[a]=s.id,pe.select.tags.push({value:s.id,label:e})}}))}}};return u((()=>{var e;return null==(e=me.params)?void 0:e.id}),(e=>{o.is.empty(e)||ve.init()})),document.addEventListener("contextmenu",(e=>{var a,t;e.preventDefault(),(null==(a=null==e?void 0:e.target)?void 0:a.closest("#tinymce"))||null==(t=re.$refs["global-menu"])||t.show(e.x,e.y)})),r((()=>{document.querySelectorAll(".container-xxl").forEach((e=>{e.classList.remove("container-xxl"),e.classList.add("container-fluid")}))})),(e,a)=>{const t=m("i-svg"),l=m("el-tooltip"),s=m("el-input"),n=m("el-option"),i=m("el-select"),c=m("el-collapse-item"),u=m("el-collapse"),r=d("load");return p(),v(j,null,[g("div",H,[g("div",I,[g("div",M,[g("div",q,[f((p(),v("div",T,[f(g("span",null,[h(C,{modelValue:b(pe).struct.content,"onUpdate:modelValue":a[0]||(a[0]=e=>b(pe).struct.content=e),id:"tinymce"},null,8,["modelValue"])],512),[[y,"tinymce"===b(pe).struct.editor]]),f(g("span",null,[h(E,{ref:"vditor",modelValue:b(pe).struct.content,"onUpdate:modelValue":a[1]||(a[1]=e=>b(pe).struct.content=e),opts:{height:600}},null,8,["modelValue"])],512),[[y,"vditor"===b(pe).struct.editor]])])),[[r,b(o).is.empty(b(pe).struct.editor)]])])]),g("div",P,[h(u,{accordion:"",modelValue:b(pe).item.active,"onUpdate:modelValue":a[8]||(a[8]=e=>b(pe).item.active=e)},{default:w((()=>[g("div",R,[g("div",S,[h(c,{name:"1"},{title:w((()=>[z(" 展示信息 ")])),default:w((()=>[g("div",$,[h(l,{content:"（必须）页面的标题",placement:"top"},{default:w((()=>[g("span",null,[h(t,{name:"hint",size:"14px"}),B])])),_:1}),h(s,{modelValue:b(pe).struct.title,"onUpdate:modelValue":a[2]||(a[2]=e=>b(pe).struct.title=e),placeholder:"页面标题"},null,8,["modelValue"])]),g("div",N,[g("label",O,[h(l,{content:"可同时选择多个标签",placement:"top"},{default:w((()=>[g("span",null,[h(t,{name:"hint",size:"14px"}),D])])),_:1})]),h(i,{modelValue:b(pe).item.tags,"onUpdate:modelValue":a[3]||(a[3]=e=>b(pe).item.tags=e),onChange:ve.change.tags,multiple:"","collapse-tags":"",filterable:"","allow-create":"","default-first-option":"",class:"d-block custom",placeholder:"请选择"},{default:w((()=>[(p(!0),v(j,null,A(b(pe).select.tags,(e=>(p(),V(n,{key:e.value,label:e.label,value:e.value},null,8,["label","value"])))),128))])),_:1},8,["modelValue","onChange"])]),g("div",F,[h(l,{content:"（必须）可以用做页面的唯一识别码或页面入口",placement:"top"},{default:w((()=>[g("span",null,[h(t,{name:"hint",size:"14px"}),J])])),_:1}),h(s,{modelValue:b(pe).struct.key,"onUpdate:modelValue":a[4]||(a[4]=e=>b(pe).struct.key=e),autocomplete:"new-password",placeholder:"唯一识别码"},null,8,["modelValue"])]),g("div",K,[h(l,{content:"备注一下",placement:"top"},{default:w((()=>[g("span",null,[h(t,{name:"hint",size:"14px"}),X])])),_:1}),h(s,{modelValue:b(pe).struct.remark,"onUpdate:modelValue":a[5]||(a[5]=e=>b(pe).struct.remark=e),autosize:{minRows:3,maxRows:10},type:"textarea"},null,8,["modelValue"])])])),_:1})])]),g("div",Y,[g("div",G,[h(c,{name:"2"},{title:w((()=>[z(" 高级选项 ")])),default:w((()=>[g("div",Q,[g("label",W,[h(l,{content:"可同时选择多个分类",placement:"top"},{default:w((()=>[g("span",null,[h(t,{name:"hint",size:"14px"}),Z])])),_:1})]),h(i,{modelValue:b(pe).struct.json.comment.allow,"onUpdate:modelValue":a[6]||(a[6]=e=>b(pe).struct.json.comment.allow=e),class:"d-block custom font-13",placeholder:"请选择"},{default:w((()=>[(p(!0),v(j,null,A(b(pe).select.comment.allow,(e=>(p(),V(n,{key:e.value,label:e.label,value:e.value},{default:w((()=>[g("span",ee,L(e.label),1),g("small",ae,L(e.value),1)])),_:2},1032,["label","value"])))),128))])),_:1},8,["modelValue"])]),g("div",te,[g("label",le,[h(l,{content:"可同时选择多个分类",placement:"top"},{default:w((()=>[g("span",null,[h(t,{name:"hint",size:"14px"}),se])])),_:1})]),h(i,{modelValue:b(pe).struct.json.comment.show,"onUpdate:modelValue":a[7]||(a[7]=e=>b(pe).struct.json.comment.show=e),class:"d-block custom font-13",placeholder:"请选择"},{default:w((()=>[(p(!0),v(j,null,A(b(pe).select.comment.show,(e=>(p(),V(n,{key:e.value,label:e.label,value:e.value},{default:w((()=>[g("span",ne,L(e.label),1),g("small",oe,L(e.value),1)])),_:2},1032,["label","value"])))),128))])),_:1},8,["modelValue"])])])),_:1})])])])),_:1},8,["modelValue"])])])]),(p(),V(x,{to:"body"},[g("div",ie,[h(l,{content:"保存页面",placement:"top"},{default:w((()=>[g("button",{onClick:a[9]||(a[9]=e=>ve.save()),type:"button",class:"btn btn-auto mimic"},[h(t,{name:"save",size:"1.6em"})])])),_:1})])])),h(b(U),_({ref:"global-menu"},b(pe).item.menu),null,16)],64)}}};export{ce as default};

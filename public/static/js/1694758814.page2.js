import{r as l,h as a,p as e,m as s,a as t,J as o,o as n,c,K as u,e as i,b as d,d as r,w as m,f as p,F as b,j as v,k as f,t as g}from"./1694758814.index.js";const h={class:"card mb-3"},w={class:"card-body"},y={class:"text-muted text-uppercase mt-0"},x=d("br",null,null,-1),k=d("br",null,null,-1),V={class:"d-inline-flex align-items-center"},j=d("span",{class:"ms-1"},"页面",-1),_={class:"m-b-20"},M=d("span",{class:"badge bg-success font-white"}," 更多 ",-1),C={class:"text-muted"},z=d("strong",{class:"flex-center"},"页面配置",-1),T={class:"row"},U={class:"col-md-6"},E={class:"form-group mb-3"},A={class:"form-label"},G=d("br",null,null,-1),J=d("br",null,null,-1),P=d("span",{class:"ms-1"},"编辑器：",-1),H={class:"font-13"},L={class:"text-muted float-end"},N={class:"col-md-6"},O={class:"form-group mb-3"},S={class:"form-label"},F=d("span",{class:"ms-1"},"审核：",-1),K={class:"font-13"},q={class:"text-muted float-end"},B={class:"col-md-6"},D={class:"form-group mb-3"},I={class:"form-label"},Q=d("span",{class:"ms-1"},"允许评论：",-1),R={class:"font-13"},W={class:"text-muted float-end"},X={class:"col-md-6"},Y={class:"form-group mb-3"},Z={class:"form-label"},$=d("span",{class:"ms-1"},"显示评论：",-1),ll={class:"font-13"},al={class:"text-muted float-end"},el={__name:"page",setup(el,{expose:sl}){const tl=l({struct:{key:"PAGE",json:{editor:"tinymce",comment:{allow:1,show:1},audit:1}},status:{finish:!1,loading:!0,dialog:!1},select:{editor:[{value:"tinymce",label:"富文本"},{value:"vditor",label:"Markdown"}],comment:{allow:[{value:1,label:"允许"},{value:0,label:"禁止"}],show:[{value:1,label:"显示"},{value:0,label:"隐藏"}]},audit:[{value:1,label:"开启",subtitle:"严格一点，防止乱搞"},{value:0,label:"关闭",subtitle:"宽松一点，方便用户"}]}});a((async()=>{await ol.init()}));const ol={init:async()=>{tl.status.finish=!1,tl.status.loading=!0;const{code:l,data:a}=await e.get("/api/config/one",{key:"PAGE"});tl.status.loading=!1,200===l&&(tl.struct=a,tl.status.finish=!0)},change:async l=>{const{code:a,msg:t}=await e.post("/api/config/save",{key:"PAGE",json:JSON.stringify({...tl.struct.json,editor:l})});200!==a&&(tl.struct.json.editor="tinymce"===tl.struct.json.editor?"vditor":"tinymce",s.error(t))},show(){if(!tl.status.finish)return s.warn("配置获取失败，无法进行配置！");tl.status.dialog=!0},save:async()=>{tl.status.wait=!0;const{code:l,msg:a}=await e.post("/api/config/save",{...tl.struct,json:JSON.stringify(tl.struct.json)});if(tl.status.wait=!1,200!==l)return s.error("保存失败："+a);tl.status.dialog=!1}};return sl({init:ol.init}),(l,a)=>{const e=t("i-svg"),s=t("el-tooltip"),el=t("el-switch"),sl=t("el-option"),nl=t("el-select"),cl=t("el-button"),ul=t("el-dialog"),il=o("loading");return n(),c(b,null,[u((n(),c("div",h,[d("div",w,[r(e,{name:"editor",color:"rgb(var(--assist-color))",size:"43px",class:"position-absolute opacity-25",style:{right:"2rem"}}),d("h6",y,[r(s,{placement:"top"},{content:m((()=>[p(" ● Markdown编辑器：Vditor支持所见即所得、即时渲染（类似 Typora）和分屏预览模式。"),x,p(" ● 富文本编辑器：TinyMCE是一个基于浏览器的所见即所得富文本编辑器，用于编辑HTML文档。"),k])),default:m((()=>[d("span",V,[r(e,{name:"hint",color:"rgb(var(--icon-color))",size:"14px"}),j])])),_:1})]),d("h2",_,[r(el,{modelValue:i(tl).struct.json.editor,"onUpdate:modelValue":a[0]||(a[0]=l=>i(tl).struct.json.editor=l),onChange:ol.change,disabled:!i(tl).status.finish,"active-value":"tinymce","inactive-value":"vditor","active-text":"富文本","inactive-text":"Markdown"},null,8,["modelValue","onChange","disabled"])]),M,d("span",C,[p(" 其它配置信息，"),d("span",{onClick:a[1]||(a[1]=l=>ol.show()),class:"text-white pointer"},"点我配置")])])])),[[il,i(tl).status.loading]]),r(ul,{modelValue:i(tl).status.dialog,"onUpdate:modelValue":a[8]||(a[8]=l=>i(tl).status.dialog=l),class:"custom",draggable:"","close-on-click-modal":!1},{header:m((()=>[z])),default:m((()=>[d("div",T,[d("div",U,[d("div",E,[d("label",A,[r(s,{placement:"top"},{content:m((()=>[p(" ● Markdown编辑器：Vditor支持所见即所得、即时渲染（类似 Typora）和分屏预览模式。"),G,p(" ● 富文本编辑器：TinyMCE是一个基于浏览器的所见即所得富文本编辑器，用于编辑HTML文档。"),J])),default:m((()=>[d("span",null,[r(e,{color:"rgb(var(--icon-color))",name:"hint",size:"14px"}),P])])),_:1})]),r(nl,{modelValue:i(tl).struct.json.editor,"onUpdate:modelValue":a[2]||(a[2]=l=>i(tl).struct.json.editor=l),class:"d-block custom font-13",placeholder:"请选择"},{default:m((()=>[(n(!0),c(b,null,v(i(tl).select.editor,(l=>(n(),f(sl,{key:l.value,label:l.label,value:l.value},{default:m((()=>[d("span",H,g(l.label),1),d("small",L,g(l.value),1)])),_:2},1032,["label","value"])))),128))])),_:1},8,["modelValue"])])]),d("div",N,[d("div",O,[d("label",S,[r(s,{content:"用户发布的页面，是否需要审核",placement:"top"},{default:m((()=>[d("span",null,[r(e,{color:"rgb(var(--icon-color))",name:"hint",size:"14px"}),F])])),_:1})]),r(nl,{modelValue:i(tl).struct.json.audit,"onUpdate:modelValue":a[3]||(a[3]=l=>i(tl).struct.json.audit=l),class:"d-block custom font-13",placeholder:"请选择"},{default:m((()=>[(n(!0),c(b,null,v(i(tl).select.audit,(l=>(n(),f(sl,{key:l.value,label:l.label,value:l.value},{default:m((()=>[d("span",K,g(l.label),1),d("small",q,g(l.subtitle),1)])),_:2},1032,["label","value"])))),128))])),_:1},8,["modelValue"])])]),d("div",B,[d("div",D,[d("label",I,[r(s,{content:"是否允许用户评论",placement:"top"},{default:m((()=>[d("span",null,[r(e,{color:"rgb(var(--icon-color))",name:"hint",size:"14px"}),Q])])),_:1})]),r(nl,{modelValue:i(tl).struct.json.comment.allow,"onUpdate:modelValue":a[4]||(a[4]=l=>i(tl).struct.json.comment.allow=l),class:"d-block custom font-13",placeholder:"请选择"},{default:m((()=>[(n(!0),c(b,null,v(i(tl).select.comment.allow,(l=>(n(),f(sl,{key:l.value,label:l.label,value:l.value},{default:m((()=>[d("span",R,g(l.label),1),d("small",W,g(l.value),1)])),_:2},1032,["label","value"])))),128))])),_:1},8,["modelValue"])])]),d("div",X,[d("div",Y,[d("label",Z,[r(s,{content:"是否显示评论",placement:"top"},{default:m((()=>[d("span",null,[r(e,{color:"rgb(var(--icon-color))",name:"hint",size:"14px"}),$])])),_:1})]),r(nl,{modelValue:i(tl).struct.json.comment.show,"onUpdate:modelValue":a[5]||(a[5]=l=>i(tl).struct.json.comment.show=l),class:"d-block custom font-13",placeholder:"请选择"},{default:m((()=>[(n(!0),c(b,null,v(i(tl).select.comment.show,(l=>(n(),f(sl,{key:l.value,label:l.label,value:l.value},{default:m((()=>[d("span",ll,g(l.label),1),d("small",al,g(l.value),1)])),_:2},1032,["label","value"])))),128))])),_:1},8,["modelValue"])])])])])),footer:m((()=>[r(cl,{onClick:a[6]||(a[6]=l=>i(tl).status.dialog=!1)},{default:m((()=>[p("取 消")])),_:1}),r(cl,{onClick:a[7]||(a[7]=l=>ol.save()),loading:i(tl).status.wait},{default:m((()=>[p("保 存")])),_:1},8,["loading"])])),_:1},8,["modelValue"])],64)}}};export{el as _};

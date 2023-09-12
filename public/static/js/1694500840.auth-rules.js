import{X as e,u as l,r as t,R as a,m as s,p as o,h as n,i,a as r,o as c,c as m,d,Y as u,w as p,e as v,b as f,k as b,l as h,t as g,E as w,F as y,j as x,f as k,g as z,P as _,K as q,Z as V,U as C,W as T}from"./1694500840.index.js";import{_ as U,a as E}from"./1694500840.page-limit.js";import{_ as M}from"./1694500840.i-table.js";const H={class:"d-flex justify-content-end"},$={class:"d-flex justify-content-end"},P=["onDblclick"],L={key:0},j={key:1},B={key:2},D=["onDblclick"],R={class:"ms-1"},S=["innerHTML"],I={class:"flex-center"},A={class:"row"},G={class:"col-md-6"},O={class:"form-group mb-3"},W={class:"form-label"},F=f("span",{class:"ms-1"},"名称：",-1),K={class:"col-md-6"},X={class:"form-group mb-3"},Y={class:"form-label"},Z=f("span",{class:"ms-1"},"费用：",-1),J={class:"row"},N={class:"col-md-4"},Q={class:"form-group mb-3"},ee={class:"form-label"},le=f("span",{class:"ms-1"},"请求类型：",-1),te={class:"text-muted float-end"},ae={class:"col-md-4"},se={class:"form-group mb-3"},oe={class:"form-label required"},ne=f("span",{class:"ms-1"},"API：",-1),ie={class:"col-md-4"},re={class:"form-group mb-3"},ce={class:"form-label"},me=f("strong",null,"● 默认：该接口类型需要在权限分组中，赋予指定用户之后，方可拥有",-1),de=f("br",null,null,-1),ue=f("strong",null,"● 公共：即不需要登录就可以使用该接口，可以理解为基础接口",-1),pe=f("br",null,null,-1),ve=f("strong",null,"● 登录：此种接口只有在登录后，才能使用，否则直接拦截掉",-1),fe=f("span",{class:"ms-1"},"接口类型：",-1),be={class:"text-muted float-end"},he={class:"row"},ge={class:"col-lg-12"},we={class:"form-group mb-3"},ye={class:"form-label"},xe=f("span",{class:"ms-1"},"备注：",-1),ke={__name:"auth-rules",props:{type:{type:String,default:"all"},params:{type:Object,default:()=>({order:"hash asc"})},init:{type:Boolean,default:!1}},emits:["refresh","update:init"],setup(_,{expose:q,emit:V}){const C=_,T=e((()=>{let e="left";return l.is.mobile()&&(e=!1),e})),U=e((()=>{let e="right";return l.is.mobile()&&(e=!1),e})),{ctx:E,proxy:ke}=z(),ze=t({item:{table:"auth-rules",dialog:!1,wait:!1},struct:{},opts:{url:"/api/auth-rules/all",params:C.params,columns:[{prop:"name",label:"名称",width:250,slot:!0,fixed:T},{prop:"route",label:"API",width:290,slot:!0},{prop:"remark",label:"备注",width:200,slot:!0},{prop:"update_time",label:"更新时间",width:120,sortable:!0},{prop:"create_time",label:"创建时间",width:120,sortable:!0}],menu:{...a,menuList:[{label:"编辑",icon:'<svg class="icon" viewBox="0 0 1024 1024" xmlns="http://www.w3.org/2000/svg" width="14" height="14"">\n                        <path fill="rgb(var(--menu-icon-color))" d="M943.104 216.064q-8.192 9.216-15.36 16.384l-12.288 12.288q-6.144 6.144-11.264 10.24l-138.24-139.264q8.192-8.192 20.48-19.456t20.48-17.408q20.48-16.384 44.032-14.336t37.888 9.216q15.36 8.192 34.304 28.672t29.184 43.008q5.12 14.336 6.656 33.792t-15.872 36.864zM551.936 329.728l158.72-158.72 138.24 138.24q-87.04 87.04-158.72 157.696-30.72 29.696-59.904 58.88t-53.248 52.224-39.424 38.4l-18.432 18.432q-7.168 7.168-16.384 14.336t-20.48 12.288-31.232 12.288-41.472 13.824-40.96 12.288-29.696 6.656q-19.456 2.048-20.992-3.584t1.536-25.088q1.024-10.24 5.12-30.208t8.192-40.448 8.704-38.4 7.68-25.088q5.12-11.264 10.752-19.456t15.872-18.432zM899.072 478.208q21.504 0 40.96 10.24t19.456 41.984l0 232.448q0 28.672-10.752 52.736t-29.184 41.984-41.984 27.648-48.128 9.728l-571.392 0q-24.576 0-48.128-10.752t-41.472-29.184-29.184-43.52-11.264-53.76l0-570.368q0-20.48 11.264-42.496t29.184-39.936 40.448-29.696 45.056-11.776l238.592 0q28.672 0 40.448 20.992t11.776 42.496-11.776 41.472-40.448 19.968l-187.392 0q-21.504 0-34.816 14.848t-13.312 36.352l0 481.28q0 20.48 13.312 34.304t34.816 13.824l474.112 0q21.504 0 36.864-13.824t15.36-34.304l0-190.464q0-14.336 6.656-24.576t16.384-16.384 21.504-8.704 23.04-2.56z">\n                        </path>\n                    </svg>',fn:e=>_e.edit(e.row),hidden:e=>!l.is.empty(e.select)},{label:"回收站",icon:'<svg class="icon" viewBox="0 0 1024 1024" xmlns="http://www.w3.org/2000/svg" width="20" height="20">\n                        <path fill="rgb(var(--menu-icon-color))" d="M256 298.666667h512v554.666666H256V298.666667z m85.333333 85.333333v384h341.333334V384H341.333333z m42.666667 85.333333h85.333333v213.333334H384v-213.333334z m170.666667 0h85.333333v213.333334h-85.333333v-213.333334zM213.333333 298.666667h597.333334v85.333333H213.333333V298.666667z m170.666667-128h256v85.333333H384V170.666667z">\n                        </path>\n                    </svg>',fn:e=>{if(l.is.empty(e.select))_e.delete(e.row.id);else{const l=e.select.map((e=>e.id));_e.delete(l)}}}]}},select:{method:[{label:"GET",value:"GET"},{label:"PUT",value:"PUT"},{label:"POST",value:"POST"},{label:"DELETE",value:"DELETE"}],type:[{value:"default",label:"默认"},{value:"common",label:"公共"},{value:"login",label:"登录"}]}}),_e={init:async()=>{await ke.$refs["i-table"].init()},save:async(e=ze.struct||{})=>{if(l.is.empty(e))return s.warn("你在想什么？什么都不填！");if(l.is.empty(e.route))return s.warn("API是必填项！");ze.item.wait=!0;const{code:t,msg:a}=await o.post(`/api/${ze.item.table}/save`,e);if(ze.item.wait=!1,200!==t)return s.error(a);ze.item.dialog=!1,await _e.init()},edit:e=>{ze.struct=e,ze.item.dialog=!0},show:()=>ze.item.dialog=!0,async delete(e=[],t=!0){if(l.is.empty(e))return;const a=`/api/${ze.item.table}/${t?"remove":"delete"}`,{code:n,msg:i}=await o.del(a,{ids:e});if(200!==n)return s.error(i);V("refresh","remove"),await _e.init()},async restore(e=[]){if(l.is.empty(e))return;const{code:t,msg:a}=await o.put(`/api/${ze.item.table}/restore`,{ids:e});if(200!==t)return s.error(a);V("refresh","all"),await _e.init()},autoWrap:(e="",t=40,a="<br>")=>l.is.empty(e)?e:e.replace(new RegExp(`(.{${t}})`,"g"),`$1${a}`),copy:(e=null,t="复制成功！")=>{if(!l.is.empty(e))return l.set.copy.text(e),l.is.empty(t)?void 0:s.info(t)},omit:(e=null,t=10,a=" ... ",s="center")=>l.is.empty(e)?"空":l.string.omit(e,t,a,s),color:(e="GET")=>({GET:"success",POST:"warning",PUT:"info",DELETE:"danger"}[e=e.toUpperCase()]||"dark")};return n((async()=>{C.init&&await _e.init()})),i((()=>C.init),(e=>{e&&_e.init()})),i((()=>ze.item.dialog),(e=>{e||(ze.struct={})})),"remove"===C.type&&(ze.opts.menu.menuList=[{label:"编辑",icon:'<svg class="icon" viewBox="0 0 1024 1024" xmlns="http://www.w3.org/2000/svg" width="14" height="14"">\n            <path fill="rgb(var(--menu-icon-color))" d="M943.104 216.064q-8.192 9.216-15.36 16.384l-12.288 12.288q-6.144 6.144-11.264 10.24l-138.24-139.264q8.192-8.192 20.48-19.456t20.48-17.408q20.48-16.384 44.032-14.336t37.888 9.216q15.36 8.192 34.304 28.672t29.184 43.008q5.12 14.336 6.656 33.792t-15.872 36.864zM551.936 329.728l158.72-158.72 138.24 138.24q-87.04 87.04-158.72 157.696-30.72 29.696-59.904 58.88t-53.248 52.224-39.424 38.4l-18.432 18.432q-7.168 7.168-16.384 14.336t-20.48 12.288-31.232 12.288-41.472 13.824-40.96 12.288-29.696 6.656q-19.456 2.048-20.992-3.584t1.536-25.088q1.024-10.24 5.12-30.208t8.192-40.448 8.704-38.4 7.68-25.088q5.12-11.264 10.752-19.456t15.872-18.432zM899.072 478.208q21.504 0 40.96 10.24t19.456 41.984l0 232.448q0 28.672-10.752 52.736t-29.184 41.984-41.984 27.648-48.128 9.728l-571.392 0q-24.576 0-48.128-10.752t-41.472-29.184-29.184-43.52-11.264-53.76l0-570.368q0-20.48 11.264-42.496t29.184-39.936 40.448-29.696 45.056-11.776l238.592 0q28.672 0 40.448 20.992t11.776 42.496-11.776 41.472-40.448 19.968l-187.392 0q-21.504 0-34.816 14.848t-13.312 36.352l0 481.28q0 20.48 13.312 34.304t34.816 13.824l474.112 0q21.504 0 36.864-13.824t15.36-34.304l0-190.464q0-14.336 6.656-24.576t16.384-16.384 21.504-8.704 23.04-2.56z">\n        </path></svg>',fn:e=>_e.edit(e.row),hidden:e=>!l.is.empty(e.select)},{label:"恢复",icon:'<svg class="icon" viewBox="0 0 1024 1024" xmlns="http://www.w3.org/2000/svg" width="14" height="14">\n            <path fill="rgb(var(--menu-icon-color))" d="M716.8 290.133333c-110.933333-102.4-281.6-106.666667-396.8-12.8S170.666667 537.6 247.466667 665.6c59.733333 106.666667 179.2 166.4 302.933333 149.333333s221.866667-102.4 256-221.866666c8.533333-34.133333 42.666667-51.2 76.8-42.666667 34.133333 8.533333 51.2 42.666667 42.666667 76.8-68.266667 226.133333-302.933333 354.133333-524.8 290.133333C174.933333 853.333333 42.666667 618.666667 106.666667 392.533333c42.666667-145.066667 153.6-256 298.666666-298.666666s298.666667 0 405.333334 102.4l81.066666-81.066667c8.533333-8.533333 21.333333-12.8 34.133334-8.533333 4.266667 12.8 12.8 21.333333 12.8 34.133333v264.533333c0 17.066667-12.8 29.866667-29.866667 29.866667h-260.266667c-12.8 0-25.6-8.533333-29.866666-17.066667s0-25.6 8.533333-34.133333l89.6-93.866667z"></path>\n        </svg>',fn:e=>{if(l.is.empty(e.select))_e.restore(e.row.id);else{const l=e.select.map((e=>e.id));_e.restore(l)}}},{label:"删除",icon:'<svg class="icon" viewBox="0 0 1024 1024" xmlns="http://www.w3.org/2000/svg" width="20" height="20">\n            <path fill="rgb(var(--menu-icon-color))" d="M256 298.666667h512v554.666666H256V298.666667z m85.333333 85.333333v384h341.333334V384H341.333333z m42.666667 85.333333h85.333333v213.333334H384v-213.333334z m170.666667 0h85.333333v213.333334h-85.333333v-213.333334zM213.333333 298.666667h597.333334v85.333333H213.333333V298.666667z m170.666667-128h256v85.333333H384V170.666667z">\n        </path></svg>',fn:e=>{if(l.is.empty(e.select))_e.delete(e.row.id,!1);else{const l=e.select.map((e=>e.id));_e.delete(l,!1)}}}]),q({init:_e.init,show:_e.show}),(e,t)=>{const a=r("el-table-column"),s=r("i-svg"),o=r("el-button"),n=r("el-tooltip"),i=r("el-input"),z=r("el-input-number"),_=r("el-option"),q=r("el-select"),V=r("el-dialog");return c(),m(y,null,[d(M,{opts:v(ze).opts,ref:"i-table"},u({start:p((()=>[d(a,{type:"selection",width:"55"})])),"i-name":p((({scope:e={}})=>[f("span",{onDblclick:l=>_e.edit(e),class:"d-flex align-items-center"},[1===parseInt(e.common)?(c(),b(n,{key:0,content:"公共权限，不需要登录即可使用的接口",placement:"top"},{default:p((()=>[d(s,{color:"rgb(var(--icon-color))",name:"!",size:"14px"})])),_:1})):h("",!0),d(n,{content:e.name,disabled:v(l).is.empty(e.name),placement:"top"},{default:p((()=>[f("span",null,g(_e.omit(null==e?void 0:e.name,16," ...","end")),1)])),_:2},1032,["content","disabled"])],40,P)])),"i-route":p((({scope:e={}})=>[d(n,{placement:"top"},{content:p((()=>["login"===e.type?(c(),m("span",L,"登录类型")):"common"===e.type?(c(),m("span",j,"公共类型")):(c(),m("span",B,"默认类型"))])),default:p((()=>[f("span",{class:w("login"===e.type?"me-1":"me-2")},["login"===e.type?(c(),b(s,{key:0,color:"rgb(var(--icon-color))",name:"user",size:"18px"})):"common"===e.type?(c(),b(s,{key:1,color:"rgb(var(--icon-color))",name:"common",size:"18px"})):(c(),b(s,{key:2,color:"rgb(var(--icon-color))",name:"!",size:"16px"}))],2)])),_:2},1024),d(n,{content:"双击复制："+e.route,disabled:v(l).is.empty(e.route),placement:"top"},{default:p((()=>[f("span",{onDblclick:l=>_e.copy(e.route,"复制成功！")},[f("span",{class:w("text-"+_e.color(e.method))},"["+g(null==e?void 0:e.method)+"]",3),f("span",R,g(_e.omit(null==e?void 0:e.route,30," ...","end")),1)],40,D)])),_:2},1032,["content","disabled"])])),"i-remark":p((({scope:e={}})=>[d(n,{disabled:v(l).is.empty(e.remark),placement:"top"},{content:p((()=>[f("span",{innerHTML:_e.autoWrap(e.remark)},null,8,S)])),default:p((()=>[f("span",null,g(_e.omit(null==e?void 0:e.remark)),1)])),_:2},1032,["disabled"])])),_:2},["all"===C.type?{name:"end",fn:p((()=>[d(a,{fixed:v(U),label:"操作",width:"100","class-name":"text-end"},{default:p((e=>[f("span",H,[d(o,{onClick:l=>_e.edit(e.row),size:"small"},{default:p((()=>[d(s,{color:"rgb(var(--icon-color))",name:"edit",size:"16px"})])),_:2},1032,["onClick"]),d(o,{onClick:l=>_e.delete(e.row.id,!0),size:"small",class:"ms-0"},{default:p((()=>[d(s,{color:"rgb(var(--icon-color))",name:"delete",size:"21px"})])),_:2},1032,["onClick"])])])),_:1},8,["fixed"])])),key:"0"}:void 0,"remove"===C.type?{name:"end",fn:p((()=>[d(a,{fixed:v(U),label:"操作",width:"160","class-name":"text-end"},{default:p((e=>[f("span",$,[d(o,{onClick:l=>_e.restore(e.row.id),size:"small"},{default:p((()=>[d(s,{color:"rgb(var(--icon-color))",name:"restore",size:"16px"})])),_:2},1032,["onClick"]),d(o,{onClick:l=>_e.edit(e.row),size:"small",class:"ms-0"},{default:p((()=>[d(s,{color:"rgb(var(--icon-color))",name:"edit",size:"16px"})])),_:2},1032,["onClick"]),d(o,{onClick:l=>_e.delete(e.row.id,!1),size:"small",class:"ms-0"},{default:p((()=>[d(s,{color:"rgb(var(--icon-color))",name:"delete",size:"21px"})])),_:2},1032,["onClick"])])])),_:1},8,["fixed"])])),key:"1"}:void 0]),1032,["opts"]),d(V,{modelValue:v(ze).item.dialog,"onUpdate:modelValue":t[8]||(t[8]=e=>v(ze).item.dialog=e),class:"custom",draggable:"","close-on-click-modal":!1},{header:p((()=>[f("strong",I,g(v(l).is.empty(v(ze).struct.id)?"添 加":"编 辑")+" 权 限 规 则",1)])),default:p((()=>[f("div",A,[f("div",G,[f("div",O,[f("label",W,[d(n,{content:"该接口名称，请遵循以下规则，如：【分组名】API名",placement:"top"},{default:p((()=>[f("span",null,[d(s,{color:"rgb(var(--icon-color))",name:"hint",size:"14px"}),F])])),_:1})]),d(i,{modelValue:v(ze).struct.name,"onUpdate:modelValue":t[0]||(t[0]=e=>v(ze).struct.name=e)},null,8,["modelValue"])])]),f("div",K,[f("div",X,[f("label",Y,[d(n,{content:"可用于接口计费模式",placement:"top"},{default:p((()=>[f("span",null,[d(s,{color:"rgb(var(--icon-color))",name:"hint",size:"14px"}),Z])])),_:1})]),d(z,{modelValue:v(ze).struct.cost,"onUpdate:modelValue":t[1]||(t[1]=e=>v(ze).struct.cost=e),min:0,class:"w-100 d-flex"},null,8,["modelValue"])])])]),f("div",J,[f("div",N,[f("div",Q,[f("label",ee,[d(n,{content:"该接口的请求类型",placement:"top"},{default:p((()=>[f("span",null,[d(s,{color:"rgb(var(--icon-color))",name:"hint",size:"14px"}),le])])),_:1})]),d(q,{modelValue:v(ze).struct.method,"onUpdate:modelValue":t[2]||(t[2]=e=>v(ze).struct.method=e),placeholder:"请选择请求类型",class:"d-block custom font-13"},{default:p((()=>[(c(!0),m(y,null,x(v(ze).select.method,(e=>(c(),b(_,{key:e.value,label:e.value,value:e.label},{default:p((()=>[f("span",{class:w("font-13 text-"+_e.color(e.value))},g(e.value),3),f("small",te,g(e.label)+" 请求",1)])),_:2},1032,["label","value"])))),128))])),_:1},8,["modelValue"])])]),f("div",ae,[f("div",se,[f("label",oe,[d(n,{content:"（必须）接口请求地址",placement:"top"},{default:p((()=>[f("span",null,[d(s,{color:"rgb(var(--icon-color))",name:"hint",size:"14px"}),ne])])),_:1})]),d(i,{modelValue:v(ze).struct.route,"onUpdate:modelValue":t[3]||(t[3]=e=>v(ze).struct.route=e)},null,8,["modelValue"])])]),f("div",ie,[f("div",re,[f("label",ce,[d(n,{placement:"top"},{content:p((()=>[me,de,ue,pe,ve])),default:p((()=>[f("span",null,[d(s,{color:"rgb(var(--icon-color))",name:"hint",size:"14px"}),fe])])),_:1})]),d(q,{modelValue:v(ze).struct.type,"onUpdate:modelValue":t[4]||(t[4]=e=>v(ze).struct.type=e),placeholder:"是否为公共接口",class:"d-block custom font-13"},{default:p((()=>[(c(!0),m(y,null,x(v(ze).select.type,(e=>(c(),b(_,{key:e.value,label:e.label,value:e.value},{default:p((()=>[f("span",null,g(e.label),1),f("small",be,g(e.value),1)])),_:2},1032,["label","value"])))),128))])),_:1},8,["modelValue"])])])]),f("div",he,[f("div",ge,[f("div",we,[f("label",ye,[d(n,{content:"备注而已，页面上不会显示此项",placement:"top"},{default:p((()=>[f("span",null,[d(s,{color:"rgb(var(--icon-color))",name:"hint",size:"14px"}),xe])])),_:1})]),d(i,{modelValue:v(ze).struct.remark,"onUpdate:modelValue":t[5]||(t[5]=e=>v(ze).struct.remark=e),autosize:{minRows:3,maxRows:10},placeholder:"备注一下，避免忘记！",type:"textarea"},null,8,["modelValue"])])])])])),footer:p((()=>[d(o,{onClick:t[6]||(t[6]=e=>v(ze).item.dialog=!1)},{default:p((()=>[k("取 消")])),_:1}),d(o,{onClick:t[7]||(t[7]=e=>_e.save()),loading:v(ze).item.wait},{default:p((()=>[k("保 存")])),_:1},8,["loading"])])),_:1},8,["modelValue"])],64)}}},ze={class:"container-fluid container-box"},_e={class:"row d-none d-lg-flex"},qe={class:"col-lg-6 d-flex"},Ve={class:"el-dropdown-link d-flex align-items-center"},Ce={class:"input-group custom-search me-1"},Te={class:"el-dropdown-link d-flex align-items-center"},Ue={class:"col-lg-6 d-flex justify-content-end",style:{"z-index":"-1"}},Ee={class:"btn btn-auto h-35px mimic",disabled:"",type:"button"},Me={class:"row mt-3"},He={class:"col-12"},$e=f("span",{class:"fw-bolder font-12"},"全部",-1),Pe=f("span",{class:"fw-bolder font-12"},"回收站",-1),Le=f("span",{class:"fw-bolder font-12"},"设置",-1),je={class:"row"},Be={class:"col-md-4"},De={class:"col-md-4"},Re={__name:"auth-rules",setup(e){const{ctx:s,proxy:o}=z(),u=t({item:{timer:null,title:"权限规则",search:null,sort:"排序",type:"类型",tabs:"all",menu:{...a,menuList:[{label:"刷新",icon:'<svg class="icon" viewBox="0 0 1024 1024" xmlns="http://www.w3.org/2000/svg" width="14" height="14">\n                    <path fill="rgb(var(--menu-icon-color))" d="M608 928c-229.76 0-416-186.24-416-416h-0.128c0-0.416 0.128-0.768 0.128-1.184a95.904 95.904 0 1 0-191.872-1.184c0 0.384-0.128 0.768-0.128 1.184l0.032 0.384c0 0.288 0.096 0.544 0.096 0.8H0c0 282.784 229.216 512 512 512 282.016 0 510.592-227.968 511.872-509.632C1022.592 743.072 836.928 928 608 928z"></path>\n                    <path fill="rgb(var(--menu-icon-color))" d="M1023.872 512H1024c0-282.784-229.216-512-512-512C230.016 0 1.408 227.968 0.128 509.632 1.408 280.96 187.072 96 416 96c229.76 0 416 186.24 416 416h0.128c0 0.416-0.128 0.768-0.128 1.184a96 96 0 0 0 96 96 95.872 95.872 0 0 0 95.872-94.816c0-0.416 0.128-0.768 0.128-1.184l-0.032-0.384c0-0.288-0.096-0.544-0.096-0.8z"></path>\n                </svg>',fn:()=>w.refresh()},{label:"添加",icon:'<svg class="icon" viewBox="0 0 1024 1024" xmlns="http://www.w3.org/2000/svg" width="14" height="14">\n                    <path fill="rgb(var(--menu-icon-color))" d="M512 1024C229.229714 1024 0 794.770286 0 512S229.229714 0 512 0s512 229.229714 512 512-229.229714 512-512 512z m0-928C282.258286 96 96 282.258286 96 512S282.258286 928 512 928 928 741.741714 928 512 741.741714 96 512 96z m208.018286 463.981714h-160v160.036572a48.018286 48.018286 0 0 1-96.036572 0v-160.036572H303.981714a47.981714 47.981714 0 0 1 0-95.963428h160V304.018286a48.018286 48.018286 0 0 1 96.036572 0v160h160a47.981714 47.981714 0 0 1 0 95.963428z"></path>\n                </svg>',fn:()=>w.add(),hidden:()=>!u.item.tabs.includes("all")}]}},params:{all:{order:"hash asc"},remove:{order:"hash asc",onlyTrashed:!0}},select:{type:[{value:"default",label:"默认"},{value:"common",label:"公共"},{value:"login",label:"登录"}]},tabs:{all:!1,remove:!1}}),w={order(e="create_time asc",l="排序"){u.item.sort=l;for(let t in u.params)u.params[t].order=e;w.refresh("all","remove")},type(e="default",l="默认"){u.item.type=l;for(let t in u.params)u.params[t].where=[["type","=",e]];w.refresh("all","remove")},add(){o.$refs.all.show()},refresh(...e){let l=["all","remove","qps","page-limit"];e=0===e.length?l:e.filter((e=>l.includes(e)));for(let t of e)o.$refs[t].init()},change:e=>u.tabs[e]=!0};return n((async()=>{u.tabs.all=!0,u.item.menu.menuList.push({line:!0},...await _())})),i((()=>u.item.search),(e=>{var t,a;const s=["all","remove"];for(let o of s)l.is.empty(e)?delete u.params[o].like:u.params[o].like=[["name",`%${e}%`],["route",`%${e}%`],["method",`${null==e?void 0:e.toUpperCase()}`],["remark",`%${e}%`]];clearTimeout(u.item.timer),u.item.timer=setTimeout((()=>w.refresh(...s)),null!=(a=null==(t=globalThis.inis)?void 0:t.lazy_time)?a:500)})),document.addEventListener("contextmenu",(e=>{var l,t;e.preventDefault(),(null==(l=null==e?void 0:e.target)?void 0:l.closest("#tabs-area"))||null==(t=o.$refs.mouse)||t.show(e.x,e.y)})),(e,l)=>{const t=r("i-svg"),a=r("el-dropdown-item"),s=r("el-dropdown"),o=r("el-tab-pane"),n=r("el-tabs");return c(),m(y,null,[f("div",ze,[f("div",_e,[f("div",qe,[v(u).item.tabs.includes("setting")?h("",!0):(c(),b(s,{key:0,class:"custom mimic me-2",trigger:"click"},{dropdown:p((()=>[d(a,{onClick:l[0]||(l[0]=e=>w.order("create_time desc","最新"))},{default:p((()=>[k("最新")])),_:1}),d(a,{onClick:l[1]||(l[1]=e=>w.order("create_time asc","最早"))},{default:p((()=>[k("最早")])),_:1})])),default:p((()=>[f("span",Ve,[k(g(v(u).item.sort)+" ",1),d(t,{name:"down"})])])),_:1})),f("div",Ce,[d(t,{name:"search",color:"rgb(var(--icon-color))",size:"18px"}),q(f("input",{"onUpdate:modelValue":l[2]||(l[2]=e=>v(u).item.search=e),class:"form-control custom search mimic",autocomplete:"new-password",type:"text",placeholder:"名称 | API | 备注 | 请求类型"},null,512),[[V,v(u).item.search]])]),f("button",{onClick:l[3]||(l[3]=e=>w.refresh()),class:"btn btn-auto mx-1 mimic",type:"button"},"刷新"),v(u).item.tabs.includes("setting")?h("",!0):(c(),b(s,{key:1,class:"custom mimic mx-1",trigger:"click"},{dropdown:p((()=>[(c(!0),m(y,null,x(v(u).select.type,((e,l)=>(c(),b(a,{key:l,onClick:l=>w.type(e.value,e.label)},{default:p((()=>[k(g(e.label),1)])),_:2},1032,["onClick"])))),128))])),default:p((()=>[f("span",Te,[k(g(v(u).item.type)+" ",1),d(t,{name:"down"})])])),_:1})),v(u).item.tabs.includes("all")?(c(),m("button",{key:2,onClick:l[4]||(l[4]=e=>w.add()),class:"btn btn-auto ms-1 mimic",type:"button"},"添加")):h("",!0)]),f("div",Ue,[f("button",Ee,g(v(u).item.title),1)])]),f("div",Me,[f("div",He,[d(n,{modelValue:v(u).item.tabs,"onUpdate:modelValue":l[7]||(l[7]=e=>v(u).item.tabs=e),onTabChange:w.change,id:"tabs-area",class:"circle"},{default:p((()=>[d(o,{name:"all"},{label:p((()=>[$e])),default:p((()=>[d(ke,{params:v(u).params.all,init:v(u).tabs.all,"onUpdate:init":l[5]||(l[5]=e=>v(u).tabs.all=e),onRefresh:w.refresh,ref:"all"},null,8,["params","init","onRefresh"])])),_:1}),d(o,{name:"remove"},{label:p((()=>[Pe])),default:p((()=>[d(ke,{params:v(u).params.remove,init:v(u).tabs.remove,"onUpdate:init":l[6]||(l[6]=e=>v(u).tabs.remove=e),onRefresh:w.refresh,ref:"remove",type:"remove"},null,8,["params","init","onRefresh"])])),_:1}),d(o,{name:"setting"},{label:p((()=>[Le])),default:p((()=>[f("div",je,[f("div",Be,[d(U,{ref:"qps"},null,512)]),f("div",De,[d(E,{ref:"page-limit"},null,512)])])])),_:1})])),_:1},8,["modelValue","onTabChange"])])])]),d(v(T),C({ref:"mouse"},v(u).item.menu),null,16)],64)}}};export{Re as default};

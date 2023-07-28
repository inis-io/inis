import{h as e,r as a,p as t,u as o,m as n,a as r,B as l,o as i,c as s,C as c,e as g,k as p,w as d,I as u,F as m,j as y,S as b,t as h,l as f,V as C,b as k,d as w,W as S}from"./1690460619.index.js";const v={class:"mt-2"},z={key:0,class:"float-start"},x={class:"float-end custom"},_={__name:"i-table",props:{opts:{type:Object,default:{url:"",method:"get",params:{},headers:{},columns:[],menu:{}},required:!0},table:{type:Object,default:{defaultSort:{prop:"id",order:"descending"},rowStyle:{backgroundColor:"rgba(var(--theme-color), calc(var(--theme-opacity) * 0.8))"},cellStyle:{backgroundColor:"transparent",border:"unset"},headerRowStyle:{backgroundColor:"transparent !important"},headerCellStyle:{backgroundColor:"rgba(var(--theme-color), var(--theme-opacity))",border:"unset"},style:{background:"rgba(var(--theme-color), calc(var(--theme-opacity) * 0.15))",backdropFilter:"blur(10px)"}}},pagination:{type:Object,default:{count:5,single:!0,class:"custom",background:!0,sizes:[10,50,100,500],layout:"sizes, prev, pager, next"}}},emits:["selection:change"],setup(_,{expose:j,emit:R}){const $=_;e((async()=>{await F.init()}));const A=a({item:{data:[],count:0,page:{code:1,total:1},limit:$.pagination.sizes[0],order:"create_time asc",loading:{data:!1,page:!1},selection:[]},config:{table:{defaultSort:{prop:"id",order:"descending"},rowStyle:{backgroundColor:"rgba(var(--theme-color), calc(var(--theme-opacity) * 0.8))"},cellStyle:{backgroundColor:"transparent",border:"unset"},headerRowStyle:{backgroundColor:"transparent !important"},headerCellStyle:{backgroundColor:"rgba(var(--theme-color), var(--theme-opacity))",border:"unset"},style:{background:"rgba(var(--theme-color), calc(var(--theme-opacity) * 0.65))",backdropFilter:"var(--theme-blur)"},...$.table},pagination:{count:5,single:!0,class:"custom",background:!0,sizes:[10,50,100,500],layout:"sizes, prev, pager, next",...$.pagination},opts:{url:"",method:"get",params:{},headers:{},columns:[],menu:{},...$.opts}}}),F={init:async(e=A.item.page.code,a=A.item.limit)=>{A.item.loading.data=!0;const{data:r,code:l,msg:i}=await t[A.config.opts.method](A.config.opts.url,{page:e,limit:a,order:A.item.order,...A.config.opts.params});if(!o.in.array(l,[200,204]))return n.error(i);A.item.data=r.data,A.item.count=r.count,A.item.page.total=r.page,A.item.page.code=e,A.item.loading.data=!1,A.item.loading.page=!1},empty:e=>o.is.empty(e),inArray:(e,a)=>o.in.array(e,a),format:e=>o.format.number(e),nature:e=>o.time.nature(e)},O={sizeChange:e=>{A.item.limit=e,F.init()},currentChange:e=>F.init(e),mouseMenu(e,a,t){const{x:o,y:n}=t;S({el:t.currentTarget,params:{row:{...e},select:[...A.item.selection.map((e=>({...e})))]},...A.config.opts.menu}).show(o,n),t.preventDefault()},selectionChange(e){A.item.selection=e,R("selection:change",e)}};return j({init:F.init}),(e,a)=>{const t=r("el-table-column"),o=r("el-table"),n=r("el-pagination"),S=l("loading");return i(),s(m,null,[c((i(),p(o,{ref:"selected","element-loading-text":"数据加载中 ...","element-loading-size":"20",onSelectionChange:O.selectionChange,onRowContextmenu:O.mouseMenu,data:g(A).item.data,"row-style":g(A).config.table.rowStyle,"cell-style":g(A).config.table.cellStyle,"header-row-style":g(A).config.table.headerRowStyle,"header-cell-style":g(A).config.table.headerCellStyle,"default-sort":g(A).config.table.defaultSort,style:C(g(A).config.table.style)},{default:d((()=>[u(e.$slots,"start"),u(e.$slots,"default"),(i(!0),s(m,null,y(g(A).config.opts.columns,((a,o)=>(i(),p(t,{key:o,prop:a.prop,label:a.label,width:a.width,"class-name":a.class,fixed:a.fixed,sortable:a.sortable,align:F.inArray(a.prop,["create_time","update_time"])?"center":a.align},b({_:2},[a.slot?{name:"default",fn:d((t=>[u(e.$slots,"i-"+a.prop,{scope:t.row})])),key:"0"}:{name:"default",fn:d((e=>[F.inArray(a.prop,["create_time","update_time"])?(i(),s("span",{key:0},[F.empty(e.row[a.prop])?(i(),s("strong",{key:1},"-")):(i(),s("span",{key:0},h(F.nature(e.row[a.prop])),1))])):f("",!0)])),key:"1"}]),1032,["prop","label","width","class-name","fixed","sortable","align"])))),128)),u(e.$slots,"end")])),_:3},8,["onSelectionChange","onRowContextmenu","data","row-style","cell-style","header-row-style","header-cell-style","default-sort","style"])),[[S,g(A).item.loading.data]]),k("div",v,[g(A).item.count>0?(i(),s("div",z," 总计 "+h(g(A).item.count)+" 条数据 ",1)):f("",!0),k("div",x,[w(n,{background:g(A).config.pagination.background,onSizeChange:O.sizeChange,onCurrentChange:O.currentChange,"page-sizes":g(A).config.pagination.sizes,layout:g(A).config.pagination.layout,"popper-class":g(A).config.pagination.class,"pager-count":g(A).config.pagination.count,"hide-on-single-page":g(A).config.pagination.single,"page-size":g(A).item.limit,"page-count":g(A).item.page.total},null,8,["background","onSizeChange","onCurrentChange","page-sizes","layout","popper-class","pager-count","hide-on-single-page","page-size","page-count"])])])],64)}}};export{_};

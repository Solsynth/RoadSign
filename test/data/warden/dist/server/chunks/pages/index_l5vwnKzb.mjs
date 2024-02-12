/* empty css                           */
import { c as createAstro, d as createComponent, r as renderTemplate, h as renderComponent, m as maybeRenderHead, e as addAttribute } from '../astro_5WdVqH1c.mjs';
import 'kleur/colors';
import 'html-escaper';
import { g as graphQuery, $ as $$PageLayout, a as $$PostList, b as $$RootLayout } from './_slug__TUDhKBhQ.mjs';
import { DocumentRenderer } from '@keystone-6/document-renderer';
import 'clsx';
/* empty css                          */

const $$Astro$2 = createAstro("https://smartsheep.studio");
const prerender$2 = false;
const $$Index$2 = createComponent(async ($$result, $$props, $$slots) => {
  const Astro2 = $$result.createAstro($$Astro$2, $$props, $$slots);
  Astro2.self = $$Index$2;
  const { events } = (await graphQuery(
    `query Query($where: EventWhereInput!) {
  events(where: $where) {
    slug
    title
    description
    content {
      document
    }
    createdAt
  }
}`,
    {
      where: {
        isHistory: {
          equals: true
        }
      }
    }
  )).data;
  return renderTemplate`${renderComponent($$result, "PageLayout", $$PageLayout, { "title": "\u6D3B\u52A8" }, { "default": ($$result2) => renderTemplate` ${maybeRenderHead()}<div class="max-w-[720px] mx-auto"> <div class="card w-full shadow-xl"> <div class="card-body"> <h2 class="card-title">活动</h2> <p>读岁月史书，涨人生阅历</p> <div class="divider"></div> <ul class="timeline timeline-snap-icon max-md:timeline-compact timeline-vertical"> ${events?.map((item, idx) => {
    let align = idx % 2 === 0 ? "timeline-start" : "timeline-end";
    let textAlign = idx % 2 === 0 ? "md:text-right" : "md:text-left";
    return renderTemplate`<li> ${idx > 0 && renderTemplate`<hr>`} <div class="timeline-middle"> <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor" class="h-5 w-5"> <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.857-9.809a.75.75 0 00-1.214-.882l-3.483 4.79-1.88-1.88a.75.75 0 10-1.06 1.061l2.5 2.5a.75.75 0 001.137-.089l4-5.5z" clip-rule="evenodd"></path> </svg> </div> <div${addAttribute(`${align} ${textAlign} mb-10`, "class")}> <time class="font-mono italic"> ${new Date(item.createdAt).toLocaleDateString()} </time> <div class="text-lg font-black">${item.title}</div> ${renderComponent($$result2, "DocumentRenderer", DocumentRenderer, { "document": item.content.document })} </div> <hr> </li>`;
  })} </ul> <div class="text-center max-md:text-left italic">
我们的故事还在继续……
</div> </div> </div> </div> ` })}`;
}, "/Users/littlesheep/Documents/Projects/Capital/src/pages/events/index.astro", void 0);

const $$file$2 = "/Users/littlesheep/Documents/Projects/Capital/src/pages/events/index.astro";
const $$url$2 = "/events";

const index$2 = /*#__PURE__*/Object.freeze(/*#__PURE__*/Object.defineProperty({
  __proto__: null,
  default: $$Index$2,
  file: $$file$2,
  prerender: prerender$2,
  url: $$url$2
}, Symbol.toStringTag, { value: 'Module' }));

const $$Astro$1 = createAstro("https://smartsheep.studio");
const prerender$1 = false;
const $$Index$1 = createComponent(async ($$result, $$props, $$slots) => {
  const Astro2 = $$result.createAstro($$Astro$1, $$props, $$slots);
  Astro2.self = $$Index$1;
  const { posts } = (await graphQuery(
    `query Query($where: PostWhereInput!, $orderBy: [PostOrderByInput!]!) {
  posts(where: $where, orderBy: $orderBy) {
    slug
    type
    title
    description
    cover {
      image {
        url
      }
    }
    content {
      document
    }
    categories {
      name
    }
    tags {
      name
    }
    createdAt
  }
}`,
    {
      orderBy: [
        {
          createdAt: "desc"
        }
      ],
      where: {}
    }
  )).data;
  return renderTemplate`${renderComponent($$result, "PageLayout", $$PageLayout, { "title": "\u8BB0\u5F55" }, { "default": ($$result2) => renderTemplate` ${maybeRenderHead()}<div class="max-w-[720px] mx-auto"> <div class="pt-16 pb-6 px-6"> <h1 class="text-4xl font-bold">记录</h1> <p class="pt-2">记录生活，记录理想，记录记录……</p> </div> ${renderComponent($$result2, "PostList", $$PostList, { "posts": posts })} </div> ` })}`;
}, "/Users/littlesheep/Documents/Projects/Capital/src/pages/posts/index.astro", void 0);

const $$file$1 = "/Users/littlesheep/Documents/Projects/Capital/src/pages/posts/index.astro";
const $$url$1 = "/posts";

const index$1 = /*#__PURE__*/Object.freeze(/*#__PURE__*/Object.defineProperty({
  __proto__: null,
  default: $$Index$1,
  file: $$file$1,
  prerender: prerender$1,
  url: $$url$1
}, Symbol.toStringTag, { value: 'Module' }));

const $$Astro = createAstro("https://smartsheep.studio");
const prerender = false;
const $$Index = createComponent(async ($$result, $$props, $$slots) => {
  const Astro2 = $$result.createAstro($$Astro, $$props, $$slots);
  Astro2.self = $$Index;
  const { events } = (await graphQuery(
    `query Query($where: EventWhereInput!) {
  events(where: $where) {
    slug
    title
    description
    createdAt
  }
}`,
    {
      where: {
        isHistory: {
          equals: true
        }
      }
    }
  )).data;
  return renderTemplate`${renderComponent($$result, "RootLayout", $$RootLayout, { "data-astro-cid-j7pv25f6": true }, { "default": ($$result2) => renderTemplate` ${maybeRenderHead()}<div class="max-h-fullpage mt-header wrapper px-5 snap-y snap-mandatory" data-astro-cid-j7pv25f6> <div id="hello" class="hero h-fullpage snap-start" data-astro-cid-j7pv25f6> <div class="hero-content w-full grid grid-cols-1 md:grid-cols-2 max-md:gap-[60px]" data-astro-cid-j7pv25f6> <div class="max-md:text-center" data-astro-cid-j7pv25f6> <h1 class="text-5xl font-bold" data-astro-cid-j7pv25f6>你好呀 👋</h1> <p class="py-6" data-astro-cid-j7pv25f6>
欢迎来到 SmartSheep Studio
            的官方网站！在这里了解，订阅，跟踪我们的最新消息。
            接触我们最大的官方社区，并且尝试最新产品，参与各种活动，提供反馈，让我们更好的服务您。
</p> <a href="#about" class="btn btn-primary btn-md" data-astro-cid-j7pv25f6>了解更多</a> </div> <div class="flex justify-center md:justify-end max-md:order-first" data-astro-cid-j7pv25f6> <div class="spinning p-3 md:p-5 shadow-2xl aspect-square rounded-[30%] w-[192px] md:w-[256px] lg:w-[384px]" data-astro-cid-j7pv25f6> <img src="/favicon.svg" alt="logo" loading="lazy" data-astro-cid-j7pv25f6> </div> </div> </div> </div> <div id="about" class="hero h-fullpage snap-start" data-astro-cid-j7pv25f6> <div class="hero-content w-full grid grid-cols-1 md:grid-cols-2 max-md:gap-[60px]" data-astro-cid-j7pv25f6> <div class="flex justify-center md:justify-start" data-astro-cid-j7pv25f6> <div class="stats shadow overflow-x-auto" data-astro-cid-j7pv25f6> <div class="stat" data-astro-cid-j7pv25f6> <div class="stat-figure text-secondary" data-astro-cid-j7pv25f6> <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" class="inline-block w-8 h-8 stroke-current" data-astro-cid-j7pv25f6><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" data-astro-cid-j7pv25f6></path></svg> </div> <div class="stat-title" data-astro-cid-j7pv25f6>People</div> <div class="stat-value" data-astro-cid-j7pv25f6>1</div> <div class="stat-desc" data-astro-cid-j7pv25f6>2019 - ${(/* @__PURE__ */ new Date()).getFullYear()}</div> </div> <div class="stat" data-astro-cid-j7pv25f6> <div class="stat-figure text-secondary" data-astro-cid-j7pv25f6> <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" class="inline-block w-8 h-8 stroke-current" data-astro-cid-j7pv25f6><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6V4m0 2a2 2 0 100 4m0-4a2 2 0 110 4m-6 8a2 2 0 100-4m0 4a2 2 0 110-4m0 4v2m0-6V4m6 6v10m6-2a2 2 0 100-4m0 4a2 2 0 110-4m0 4v2m0-6V4" data-astro-cid-j7pv25f6></path></svg> </div> <div class="stat-title" data-astro-cid-j7pv25f6>Clients</div> <div class="stat-value" data-astro-cid-j7pv25f6>180</div> <div class="stat-desc" data-astro-cid-j7pv25f6>↗︎ 80 (44%)</div> </div> <div class="stat" data-astro-cid-j7pv25f6> <div class="stat-figure text-secondary" data-astro-cid-j7pv25f6> <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" class="inline-block w-8 h-8 stroke-current" data-astro-cid-j7pv25f6><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 8h14M5 8a2 2 0 110-4h14a2 2 0 110 4M5 8v10a2 2 0 002 2h10a2 2 0 002-2V8m-9 4h4" data-astro-cid-j7pv25f6></path></svg> </div> <div class="stat-title" data-astro-cid-j7pv25f6>Products</div> <div class="stat-value" data-astro-cid-j7pv25f6>4</div> <div class="stat-desc" data-astro-cid-j7pv25f6>↘︎ 8 (67%)</div> </div> </div> </div> <div class="max-md:text-center" data-astro-cid-j7pv25f6> <h1 class="text-5xl font-bold" data-astro-cid-j7pv25f6>关于我们 🔖</h1> <p class="py-6" data-astro-cid-j7pv25f6>
我们是一群充满活力、对开源充满热情的开发者。成立于 2019
            年。自那年起我们一直在开发让人喜欢的开源软件。在我们这里，“取之于开源，用之于开源”
            不仅是原则，更是我们信仰的座右铭。
</p> <a href="#history" class="btn btn-primary btn-md pl-[24px]" data-astro-cid-j7pv25f6>
查看「岁月史书」
</a> </div> </div> </div> <div id="history" class="flex flex-col justify-center items-center h-fullpage snap-start" data-astro-cid-j7pv25f6> <div class="text-center" data-astro-cid-j7pv25f6> <div data-astro-cid-j7pv25f6> <h1 class="text-4xl font-bold" data-astro-cid-j7pv25f6>岁月史书</h1> <p class="pt-2 pb-4 tracking-[8px]" data-astro-cid-j7pv25f6>但当涉猎，见往事耳</p> <ul class="pb-6 mx-[-20px] max-w-[100vw] px-5 flex justify-center history timeline timeline-horizontal" data-astro-cid-j7pv25f6> ${events?.map((item, idx) => renderTemplate`<li data-astro-cid-j7pv25f6> ${idx > 0 && renderTemplate`<hr data-astro-cid-j7pv25f6>`} <div class="timeline-start" data-astro-cid-j7pv25f6> ${new Date(item.createdAt).toLocaleDateString()} </div> <div class="timeline-middle" data-astro-cid-j7pv25f6> <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor" class="w-5 h-5" data-astro-cid-j7pv25f6> <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.857-9.809a.75.75 0 00-1.214-.882l-3.483 4.79-1.88-1.88a.75.75 0 10-1.06 1.061l2.5 2.5a.75.75 0 001.137-.089l4-5.5z" clip-rule="evenodd" data-astro-cid-j7pv25f6></path> </svg> </div> <div class="timeline-end timeline-box" data-astro-cid-j7pv25f6> <h2 class="font-bold text-lg" data-astro-cid-j7pv25f6>${item.title}</h2> <div class="line-clamp-2" data-astro-cid-j7pv25f6>${item.description}</div> </div> ${idx < events?.length - 1 && renderTemplate`<hr data-astro-cid-j7pv25f6>`} </li>`)} </ul> <a class="btn btn-primary" href="/events" data-astro-cid-j7pv25f6>查看更多</a> </div> </div> </div> </div> ` })}  `;
}, "/Users/littlesheep/Documents/Projects/Capital/src/pages/index.astro", void 0);

const $$file = "/Users/littlesheep/Documents/Projects/Capital/src/pages/index.astro";
const $$url = "";

const index = /*#__PURE__*/Object.freeze(/*#__PURE__*/Object.defineProperty({
  __proto__: null,
  default: $$Index,
  file: $$file,
  prerender,
  url: $$url
}, Symbol.toStringTag, { value: 'Module' }));

export { index$1 as a, index as b, index$2 as i };

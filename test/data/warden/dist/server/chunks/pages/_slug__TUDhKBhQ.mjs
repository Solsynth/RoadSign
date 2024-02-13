/* empty css                           */
import 'html-escaper';
import { c as createAstro, d as createComponent, r as renderTemplate, m as maybeRenderHead, e as addAttribute, f as renderSlot, g as renderTransition, h as renderComponent, i as renderHead } from '../astro_5WdVqH1c.mjs';
import 'kleur/colors';
import 'clsx';
import { DocumentRenderer } from '@keystone-6/document-renderer';
/* empty css                           */
/* empty css                           */

const $$Astro$7 = createAstro("https://smartsheep.studio");
const $$Navbar = createComponent(async ($$result, $$props, $$slots) => {
  const Astro2 = $$result.createAstro($$Astro$7, $$props, $$slots);
  Astro2.self = $$Navbar;
  const items = [
    {
      label: "\u60C5\u62A5",
      children: [
        { href: "/posts", label: "\u8BB0\u5F55" },
        { href: "/events", label: "\u6D3B\u52A8" }
      ]
    }
  ];
  return renderTemplate`${maybeRenderHead()}<div class="fixed top-0 navbar shadow-md bg-base-100 lg:px-5 z-10"> <div class="navbar-start"> <div class="dropdown"> <div tabindex="0" role="button" class="btn btn-ghost lg:hidden"> <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor"> <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h8m-8 6h16"></path> </svg> </div> <ul tabindex="0" class="menu menu-sm dropdown-content mt-3 z-[1] p-2 shadow bg-base-100 rounded-box w-52"> ${items.map((item) => renderTemplate`<li> <a${addAttribute(item.href, "href")}>${item.label}</a> ${item.children && renderTemplate`<ul class="p-2"> ${item.children?.map((child) => renderTemplate`<li> <a${addAttribute(child.href, "href")}>${child.label}</a> </li>`)} </ul>`} </li>`)} </ul> </div> <a class="btn btn-ghost text-xl" href="/">山羊寒舍</a> </div> <div class="navbar-center hidden lg:flex"> <ul class="menu menu-horizontal px-1"> ${items.map((item) => renderTemplate`<li> ${item.children ? renderTemplate`<details> <summary>${item.label}</summary> <ul class="p-2"> ${item.children?.map((child) => renderTemplate`<li> <a${addAttribute(child.href, "href")}>${child.label}</a> </li>`)} </ul> </details>` : renderTemplate`<a${addAttribute(item.href, "href")}>${item.label}</a>`} </li>`)} </ul> </div> <div class="navbar-end"> <label class="swap swap-rotate px-[16px]"> <input type="checkbox" class="theme-controller" value="light" checked> <svg class="swap-on fill-current w-8 h-8" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24"> <path d="M5.64,17l-.71.71a1,1,0,0,0,0,1.41,1,1,0,0,0,1.41,0l.71-.71A1,1,0,0,0,5.64,17ZM5,12a1,1,0,0,0-1-1H3a1,1,0,0,0,0,2H4A1,1,0,0,0,5,12Zm7-7a1,1,0,0,0,1-1V3a1,1,0,0,0-2,0V4A1,1,0,0,0,12,5ZM5.64,7.05a1,1,0,0,0,.7.29,1,1,0,0,0,.71-.29,1,1,0,0,0,0-1.41l-.71-.71A1,1,0,0,0,4.93,6.34Zm12,.29a1,1,0,0,0,.7-.29l.71-.71a1,1,0,1,0-1.41-1.41L17,5.64a1,1,0,0,0,0,1.41A1,1,0,0,0,17.66,7.34ZM21,11H20a1,1,0,0,0,0,2h1a1,1,0,0,0,0-2Zm-9,8a1,1,0,0,0-1,1v1a1,1,0,0,0,2,0V20A1,1,0,0,0,12,19ZM18.36,17A1,1,0,0,0,17,18.36l.71.71a1,1,0,0,0,1.41,0,1,1,0,0,0,0-1.41ZM12,6.5A5.5,5.5,0,1,0,17.5,12,5.51,5.51,0,0,0,12,6.5Zm0,9A3.5,3.5,0,1,1,15.5,12,3.5,3.5,0,0,1,12,15.5Z"></path> </svg> <svg class="swap-off fill-current w-8 h-8" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24"> <path d="M21.64,13a1,1,0,0,0-1.05-.14,8.05,8.05,0,0,1-3.37.73A8.15,8.15,0,0,1,9.08,5.49a8.59,8.59,0,0,1,.25-2A1,1,0,0,0,8,2.36,10.14,10.14,0,1,0,22,14.05,1,1,0,0,0,21.64,13Zm-9.5,6.69A8.14,8.14,0,0,1,7.08,5.22v.27A10.15,10.15,0,0,0,17.22,15.63a9.79,9.79,0,0,0,2.1-.22A8.11,8.11,0,0,1,12.14,19.73Z"></path> </svg> </label> </div> </div>`;
}, "/Users/littlesheep/Documents/Projects/Capital/src/components/Navbar.astro", void 0);

const $$Astro$6 = createAstro("https://smartsheep.studio");
const $$ViewTransitions = createComponent(async ($$result, $$props, $$slots) => {
  const Astro2 = $$result.createAstro($$Astro$6, $$props, $$slots);
  Astro2.self = $$ViewTransitions;
  const { fallback = "animate" } = Astro2.props;
  return renderTemplate`<meta name="astro-view-transitions-enabled" content="true"><meta name="astro-view-transitions-fallback"${addAttribute(fallback, "content")}>`;
}, "/Users/littlesheep/Documents/Projects/Capital/node_modules/astro/components/ViewTransitions.astro", void 0);

var __freeze = Object.freeze;
var __defProp = Object.defineProperty;
var __template = (cooked, raw) => __freeze(__defProp(cooked, "raw", { value: __freeze(raw || cooked.slice()) }));
var _a;
const $$Astro$5 = createAstro("https://smartsheep.studio");
const $$RootLayout = createComponent(async ($$result, $$props, $$slots) => {
  const Astro2 = $$result.createAstro($$Astro$5, $$props, $$slots);
  Astro2.self = $$RootLayout;
  const { title } = Astro2.props;
  return renderTemplate(_a || (_a = __template(['<html lang="en" data-astro-cid-mdysn4oi> <head><meta charset="utf-8"><link rel="icon" type="image/svg+xml" href="/favicon.svg"><meta name="viewport" content="width=device-width"><meta name="generator"', ">", "", "", "", "</head> <body data-astro-cid-mdysn4oi> <!-- Header --> ", " <!-- Content --> <main data-astro-cid-mdysn4oi", "> ", ' </main> <!-- Styles -->   <script async src="https://analytics.smartsheep.studio/script.js" data-website-id="9d676a27-b473-44a3-b444-5a7d851e31e8"><\/script> </body> </html>'])), addAttribute(Astro2.generator, "content"), title && renderTemplate`<title>山羊寒舍 | ${title}</title>`, !title && renderTemplate`<title>山羊寒舍</title>`, renderComponent($$result, "ViewTransitions", $$ViewTransitions, { "data-astro-cid-mdysn4oi": true }), renderHead(), renderComponent($$result, "Navbar", $$Navbar, { "data-astro-cid-mdysn4oi": true }), addAttribute(renderTransition($$result, "53mar5bf", "slide", ""), "data-astro-transition-scope"), renderSlot($$result, $$slots["default"]));
}, "/Users/littlesheep/Documents/Projects/Capital/src/layouts/RootLayout.astro", "self");

const $$Astro$4 = createAstro("https://smartsheep.studio");
const $$PageLayout = createComponent(async ($$result, $$props, $$slots) => {
  const Astro2 = $$result.createAstro($$Astro$4, $$props, $$slots);
  Astro2.self = $$PageLayout;
  const { title } = Astro2.props;
  return renderTemplate`${renderComponent($$result, "RootLayout", $$RootLayout, { "title": title }, { "default": ($$result2) => renderTemplate` ${maybeRenderHead()}<main class="container mx-auto h-fullpage mt-header"> ${renderSlot($$result2, $$slots["default"])} </main> ` })}`;
}, "/Users/littlesheep/Documents/Projects/Capital/src/layouts/PageLayout.astro", void 0);

const defaultCms = "https://smartsheep.studio";
async function graphQuery(query, variables) {
  const response = await fetch(`${process.env.PUBLIC_CMS ?? defaultCms}/api/graphql`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({
      query,
      variables
    })
  });
  return await response.json();
}

const POST_TYPES = {
  article: "文章",
  podcast: "播客",
  announcements: "通告"
};

const $$Astro$3 = createAstro("https://smartsheep.studio");
const $$PostList = createComponent(async ($$result, $$props, $$slots) => {
  const Astro2 = $$result.createAstro($$Astro$3, $$props, $$slots);
  Astro2.self = $$PostList;
  const { posts } = Astro2.props;
  return renderTemplate`${maybeRenderHead()}<div class="grid justify-items-strench shadow-lg"> ${posts?.map((item) => renderTemplate`<a${addAttribute(`/p/${item.slug}`, "href")}> <div class="card sm:card-side hover:bg-base-200 transition-colors sm:max-w-none"> ${item.cover.image.url && renderTemplate`<figure class="mx-auto w-full object-cover p-6 max-sm:pb-0 sm:max-w-[12rem] sm:pe-0"> <img loading="lazy"${addAttribute(item.cover.image.url, "src")} class="border-base-content bg-base-300 rounded-btn border border-opacity-5"${addAttribute(item.title, "alt")}> </figure>`} <div class="card-body"> <h2 class="text-xl">${item.title}</h2> <div class="mx-[-2px] mt-[-4px]"> <span class="badge badge-accent">${POST_TYPES[item.type]}</span> ${item.categories?.map((category) => renderTemplate`<span class="badge badge-primary">${category.name}</span>`)} ${item.tags?.map((tag) => renderTemplate`<span class="badge badge-secondary">${tag.name}</span>`)} </div> <div class="text-xs opacity-60 line-clamp-3"> ${item.description} </div> </div> </div> </a>`)} </div>`;
}, "/Users/littlesheep/Documents/Projects/Capital/src/components/PostList.astro", void 0);

const $$Astro$2 = createAstro("https://smartsheep.studio");
const prerender$2 = false;
const $$slug$2 = createComponent(async ($$result, $$props, $$slots) => {
  const Astro2 = $$result.createAstro($$Astro$2, $$props, $$slots);
  Astro2.self = $$slug$2;
  const { slug } = Astro2.params;
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
      where: { categories: { some: { slug: { equals: slug } } } }
    }
  )).data;
  return renderTemplate`${renderComponent($$result, "PageLayout", $$PageLayout, { "title": "\u5206\u7C7B\u68C0\u7D22" }, { "default": ($$result2) => renderTemplate` ${maybeRenderHead()}<div class="max-w-[720px] mx-auto"> <div class="pt-16 pb-6 px-6"> <h1 class="text-4xl font-bold">分类检索</h1> <p class="pt-3">以下是包含该分类的记录……</p> </div> ${renderComponent($$result2, "PostList", $$PostList, { "posts": posts })} </div> ` })}`;
}, "/Users/littlesheep/Documents/Projects/Capital/src/pages/categories/[slug].astro", void 0);

const $$file$2 = "/Users/littlesheep/Documents/Projects/Capital/src/pages/categories/[slug].astro";
const $$url$2 = "/categories/[slug]";

const _slug_$2 = /*#__PURE__*/Object.freeze(/*#__PURE__*/Object.defineProperty({
  __proto__: null,
  default: $$slug$2,
  file: $$file$2,
  prerender: prerender$2,
  url: $$url$2
}, Symbol.toStringTag, { value: 'Module' }));

const $$Astro$1 = createAstro("https://smartsheep.studio");
const prerender$1 = false;
const $$slug$1 = createComponent(async ($$result, $$props, $$slots) => {
  const Astro2 = $$result.createAstro($$Astro$1, $$props, $$slots);
  Astro2.self = $$slug$1;
  const { slug } = Astro2.params;
  const { post } = (await graphQuery(
    `query Query($where: PostWhereUniqueInput!) {
  post(where: $where) {
    slug
    type
    title
    description
    author {
      name
    }
    assets {
      caption
      url
      type
    }
    cover {
      image {
        url
      }
    }
    content {
      document
    }
    categories {
      slug
      name
    }
    tags {
      slug
      name
    }
    createdAt
  }
}`,
    {
      where: { slug }
    }
  )).data;
  return renderTemplate`${renderComponent($$result, "PageLayout", $$PageLayout, { "title": post.title, "data-astro-cid-gysqo7gh": true }, { "default": ($$result2) => renderTemplate` ${maybeRenderHead()}<div class="wrapper" data-astro-cid-gysqo7gh> <div class="card w-full shadow-xl" data-astro-cid-gysqo7gh> ${post.cover && renderTemplate`<figure data-astro-cid-gysqo7gh> <img${addAttribute(post.cover.image.url, "src")}${addAttribute(post.title, "alt")} data-astro-cid-gysqo7gh> </figure>`} <div class="card-body" data-astro-cid-gysqo7gh> <h2 class="card-title" data-astro-cid-gysqo7gh>${post.title}</h2> <p class="description" data-astro-cid-gysqo7gh>${post.description ?? "No description"}</p> <div class="divider" data-astro-cid-gysqo7gh></div> ${post.assets?.length > 0 && renderTemplate`<div class="mb-5 w-full" data-astro-cid-gysqo7gh> ${renderComponent($$result2, "Media", null, { "client:only": true, "sources": post.assets, "author": post.author, "client:component-hydration": "only", "data-astro-cid-gysqo7gh": true, "client:component-path": "/Users/littlesheep/Documents/Projects/Capital/src/components/posts/Media", "client:component-export": "default" })} </div>`} <div class="prose max-w-none" data-astro-cid-gysqo7gh> ${renderComponent($$result2, "DocumentRenderer", DocumentRenderer, { "document": post.content.document, "data-astro-cid-gysqo7gh": true })} </div> </div> </div> <div class="h-fit sticky top-header" data-astro-cid-gysqo7gh> <div class="card shadow-xl" data-astro-cid-gysqo7gh> <div class="card-body" data-astro-cid-gysqo7gh> <div class="gap-2 text-sm metadata description" data-astro-cid-gysqo7gh> <div data-astro-cid-gysqo7gh> <div data-astro-cid-gysqo7gh>作者</div> <div data-astro-cid-gysqo7gh>${post.author?.name ?? "\u4F5A\u540D"}</div> </div> <div data-astro-cid-gysqo7gh> <div data-astro-cid-gysqo7gh>类型</div> <div class="text-accent" data-astro-cid-gysqo7gh> ${POST_TYPES[post.type]} </div> </div> <div data-astro-cid-gysqo7gh> <div data-astro-cid-gysqo7gh>分类</div> <div class="flex gap-1" data-astro-cid-gysqo7gh> ${post.categories?.map((category) => renderTemplate`<a${addAttribute(`/categories/${category.slug}`, "href")} class="link link-primary" data-astro-cid-gysqo7gh> ${category.name} </a>`)} </div> </div> <div data-astro-cid-gysqo7gh> <div data-astro-cid-gysqo7gh>标签</div> <div class="flex gap-1" data-astro-cid-gysqo7gh> ${post.tags?.map((tag) => renderTemplate`<a${addAttribute(`/tags/${tag.slug}`, "href")} class="link link-secondary" data-astro-cid-gysqo7gh> ${tag.name} </a>`)} </div> </div> <div data-astro-cid-gysqo7gh> <div data-astro-cid-gysqo7gh>发布于</div> <div data-astro-cid-gysqo7gh>${new Date(post.createdAt).toLocaleString()}</div> </div> </div> </div> </div> </div> </div> ` })} `;
}, "/Users/littlesheep/Documents/Projects/Capital/src/pages/posts/[slug].astro", void 0);

const $$file$1 = "/Users/littlesheep/Documents/Projects/Capital/src/pages/posts/[slug].astro";
const $$url$1 = "/posts/[slug]";

const _slug_$1 = /*#__PURE__*/Object.freeze(/*#__PURE__*/Object.defineProperty({
  __proto__: null,
  default: $$slug$1,
  file: $$file$1,
  prerender: prerender$1,
  url: $$url$1
}, Symbol.toStringTag, { value: 'Module' }));

const $$Astro = createAstro("https://smartsheep.studio");
const prerender = false;
const $$slug = createComponent(async ($$result, $$props, $$slots) => {
  const Astro2 = $$result.createAstro($$Astro, $$props, $$slots);
  Astro2.self = $$slug;
  const { slug } = Astro2.params;
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
      where: { tags: { some: { slug: { equals: slug } } } }
    }
  )).data;
  return renderTemplate`${renderComponent($$result, "PageLayout", $$PageLayout, { "title": "\u6807\u7B7E\u68C0\u7D22" }, { "default": ($$result2) => renderTemplate` ${maybeRenderHead()}<div class="max-w-[720px] mx-auto"> <div class="pt-16 pb-6 px-6"> <h1 class="text-4xl font-bold">标签检索</h1> <p class="pt-3">以下是包含该标签的记录……</p> </div> ${renderComponent($$result2, "PostList", $$PostList, { "posts": posts })} </div> ` })}`;
}, "/Users/littlesheep/Documents/Projects/Capital/src/pages/tags/[slug].astro", void 0);

const $$file = "/Users/littlesheep/Documents/Projects/Capital/src/pages/tags/[slug].astro";
const $$url = "/tags/[slug]";

const _slug_ = /*#__PURE__*/Object.freeze(/*#__PURE__*/Object.defineProperty({
  __proto__: null,
  default: $$slug,
  file: $$file,
  prerender,
  url: $$url
}, Symbol.toStringTag, { value: 'Module' }));

export { $$PageLayout as $, _slug_$2 as _, $$PostList as a, $$RootLayout as b, _slug_$1 as c, _slug_ as d, graphQuery as g };

import '@astrojs/internal-helpers/path';
import 'cookie';
import { bold, red, yellow, dim, blue } from 'kleur/colors';
import 'string-width';
import 'html-escaper';
import 'clsx';
import './chunks/astro_5WdVqH1c.mjs';
import { compile } from 'path-to-regexp';

const dateTimeFormat = new Intl.DateTimeFormat([], {
  hour: "2-digit",
  minute: "2-digit",
  second: "2-digit",
  hour12: false
});
const levels = {
  debug: 20,
  info: 30,
  warn: 40,
  error: 50,
  silent: 90
};
function log(opts, level, label, message, newLine = true) {
  const logLevel = opts.level;
  const dest = opts.dest;
  const event = {
    label,
    level,
    message,
    newLine
  };
  if (!isLogLevelEnabled(logLevel, level)) {
    return;
  }
  dest.write(event);
}
function isLogLevelEnabled(configuredLogLevel, level) {
  return levels[configuredLogLevel] <= levels[level];
}
function info(opts, label, message, newLine = true) {
  return log(opts, "info", label, message, newLine);
}
function warn(opts, label, message, newLine = true) {
  return log(opts, "warn", label, message, newLine);
}
function error(opts, label, message, newLine = true) {
  return log(opts, "error", label, message, newLine);
}
function debug(...args) {
  if ("_astroGlobalDebug" in globalThis) {
    globalThis._astroGlobalDebug(...args);
  }
}
function getEventPrefix({ level, label }) {
  const timestamp = `${dateTimeFormat.format(/* @__PURE__ */ new Date())}`;
  const prefix = [];
  if (level === "error" || level === "warn") {
    prefix.push(bold(timestamp));
    prefix.push(`[${level.toUpperCase()}]`);
  } else {
    prefix.push(timestamp);
  }
  if (label) {
    prefix.push(`[${label}]`);
  }
  if (level === "error") {
    return red(prefix.join(" "));
  }
  if (level === "warn") {
    return yellow(prefix.join(" "));
  }
  if (prefix.length === 1) {
    return dim(prefix[0]);
  }
  return dim(prefix[0]) + " " + blue(prefix.splice(1).join(" "));
}
if (typeof process !== "undefined") {
  let proc = process;
  if ("argv" in proc && Array.isArray(proc.argv)) {
    if (proc.argv.includes("--verbose")) ; else if (proc.argv.includes("--silent")) ; else ;
  }
}
class Logger {
  options;
  constructor(options) {
    this.options = options;
  }
  info(label, message, newLine = true) {
    info(this.options, label, message, newLine);
  }
  warn(label, message, newLine = true) {
    warn(this.options, label, message, newLine);
  }
  error(label, message, newLine = true) {
    error(this.options, label, message, newLine);
  }
  debug(label, ...messages) {
    debug(label, ...messages);
  }
  level() {
    return this.options.level;
  }
  forkIntegrationLogger(label) {
    return new AstroIntegrationLogger(this.options, label);
  }
}
class AstroIntegrationLogger {
  options;
  label;
  constructor(logging, label) {
    this.options = logging;
    this.label = label;
  }
  /**
   * Creates a new logger instance with a new label, but the same log options.
   */
  fork(label) {
    return new AstroIntegrationLogger(this.options, label);
  }
  info(message) {
    info(this.options, this.label, message);
  }
  warn(message) {
    warn(this.options, this.label, message);
  }
  error(message) {
    error(this.options, this.label, message);
  }
  debug(message) {
    debug(this.label, message);
  }
}

function getRouteGenerator(segments, addTrailingSlash) {
  const template = segments.map((segment) => {
    return "/" + segment.map((part) => {
      if (part.spread) {
        return `:${part.content.slice(3)}(.*)?`;
      } else if (part.dynamic) {
        return `:${part.content}`;
      } else {
        return part.content.normalize().replace(/\?/g, "%3F").replace(/#/g, "%23").replace(/%5B/g, "[").replace(/%5D/g, "]").replace(/[.*+?^${}()|[\]\\]/g, "\\$&");
      }
    }).join("");
  }).join("");
  let trailing = "";
  if (addTrailingSlash === "always" && segments.length) {
    trailing = "/";
  }
  const toPath = compile(template + trailing);
  return toPath;
}

function deserializeRouteData(rawRouteData) {
  return {
    route: rawRouteData.route,
    type: rawRouteData.type,
    pattern: new RegExp(rawRouteData.pattern),
    params: rawRouteData.params,
    component: rawRouteData.component,
    generate: getRouteGenerator(rawRouteData.segments, rawRouteData._meta.trailingSlash),
    pathname: rawRouteData.pathname || void 0,
    segments: rawRouteData.segments,
    prerender: rawRouteData.prerender,
    redirect: rawRouteData.redirect,
    redirectRoute: rawRouteData.redirectRoute ? deserializeRouteData(rawRouteData.redirectRoute) : void 0,
    fallbackRoutes: rawRouteData.fallbackRoutes.map((fallback) => {
      return deserializeRouteData(fallback);
    })
  };
}

function deserializeManifest(serializedManifest) {
  const routes = [];
  for (const serializedRoute of serializedManifest.routes) {
    routes.push({
      ...serializedRoute,
      routeData: deserializeRouteData(serializedRoute.routeData)
    });
    const route = serializedRoute;
    route.routeData = deserializeRouteData(serializedRoute.routeData);
  }
  const assets = new Set(serializedManifest.assets);
  const componentMetadata = new Map(serializedManifest.componentMetadata);
  const clientDirectives = new Map(serializedManifest.clientDirectives);
  return {
    ...serializedManifest,
    assets,
    componentMetadata,
    clientDirectives,
    routes
  };
}

const manifest = deserializeManifest({"adapterName":"@astrojs/node","routes":[{"file":"","links":[],"scripts":[],"styles":[],"routeData":{"type":"endpoint","isIndex":false,"route":"/_image","pattern":"^\\/_image$","segments":[[{"content":"_image","dynamic":false,"spread":false}]],"params":[],"component":"node_modules/astro/dist/assets/endpoint/node.js","pathname":"/_image","prerender":false,"fallbackRoutes":[],"_meta":{"trailingSlash":"ignore"}}},{"file":"","links":[],"scripts":[{"type":"external","value":"/_astro/hoisted.l-JsOPk0.js"}],"styles":[{"type":"external","src":"/_astro/_slug_.yOjdTrIk.css"},{"type":"external","src":"/_astro/_slug_.bcjV8AoT.css"}],"routeData":{"route":"/events","isIndex":true,"type":"page","pattern":"^\\/events\\/?$","segments":[[{"content":"events","dynamic":false,"spread":false}]],"params":[],"component":"src/pages/events/index.astro","pathname":"/events","prerender":false,"fallbackRoutes":[],"_meta":{"trailingSlash":"ignore"}}},{"file":"","links":[],"scripts":[{"type":"external","value":"/_astro/hoisted.l-JsOPk0.js"}],"styles":[{"type":"external","src":"/_astro/_slug_.yOjdTrIk.css"},{"type":"external","src":"/_astro/_slug_.bcjV8AoT.css"}],"routeData":{"route":"/posts","isIndex":true,"type":"page","pattern":"^\\/posts\\/?$","segments":[[{"content":"posts","dynamic":false,"spread":false}]],"params":[],"component":"src/pages/posts/index.astro","pathname":"/posts","prerender":false,"fallbackRoutes":[],"_meta":{"trailingSlash":"ignore"}}},{"file":"","links":[],"scripts":[{"type":"external","value":"/_astro/hoisted.l-JsOPk0.js"}],"styles":[{"type":"external","src":"/_astro/_slug_.yOjdTrIk.css"},{"type":"external","src":"/_astro/_slug_.bcjV8AoT.css"}],"routeData":{"route":"/categories/[slug]","isIndex":false,"type":"page","pattern":"^\\/categories\\/([^/]+?)\\/?$","segments":[[{"content":"categories","dynamic":false,"spread":false}],[{"content":"slug","dynamic":true,"spread":false}]],"params":["slug"],"component":"src/pages/categories/[slug].astro","prerender":false,"fallbackRoutes":[],"_meta":{"trailingSlash":"ignore"}}},{"file":"","links":[],"scripts":[{"type":"external","value":"/_astro/hoisted.l-JsOPk0.js"}],"styles":[{"type":"external","src":"/_astro/_slug_.yOjdTrIk.css"},{"type":"external","src":"/_astro/_slug_.bcjV8AoT.css"},{"type":"inline","content":".wrapper[data-astro-cid-gysqo7gh]{display:grid;grid-template-columns:1fr;gap:20px}.description[data-astro-cid-gysqo7gh]{color:oklch(var(--bc) / .8)}.metadata[data-astro-cid-gysqo7gh]{display:flex;flex-direction:column;transition:color .3s}@media (min-width: 768px){.wrapper[data-astro-cid-gysqo7gh]{grid-template-columns:2fr 1fr}}\n"},{"type":"external","src":"/_astro/Media.Co8_pG1j.css"}],"routeData":{"route":"/posts/[slug]","isIndex":false,"type":"page","pattern":"^\\/posts\\/([^/]+?)\\/?$","segments":[[{"content":"posts","dynamic":false,"spread":false}],[{"content":"slug","dynamic":true,"spread":false}]],"params":["slug"],"component":"src/pages/posts/[slug].astro","prerender":false,"fallbackRoutes":[],"_meta":{"trailingSlash":"ignore"}}},{"file":"","links":[],"scripts":[{"type":"external","value":"/_astro/hoisted.l-JsOPk0.js"}],"styles":[{"type":"external","src":"/_astro/_slug_.yOjdTrIk.css"},{"type":"external","src":"/_astro/_slug_.bcjV8AoT.css"}],"routeData":{"route":"/tags/[slug]","isIndex":false,"type":"page","pattern":"^\\/tags\\/([^/]+?)\\/?$","segments":[[{"content":"tags","dynamic":false,"spread":false}],[{"content":"slug","dynamic":true,"spread":false}]],"params":["slug"],"component":"src/pages/tags/[slug].astro","prerender":false,"fallbackRoutes":[],"_meta":{"trailingSlash":"ignore"}}},{"file":"","links":[],"scripts":[{"type":"external","value":"/_astro/hoisted.l-JsOPk0.js"}],"styles":[{"type":"external","src":"/_astro/_slug_.yOjdTrIk.css"},{"type":"external","src":"/_astro/_slug_.bcjV8AoT.css"},{"type":"inline","content":".wrapper[data-astro-cid-j7pv25f6]{overflow-y:auto;scrollbar-width:none;scroll-behavior:smooth}.wrapper[data-astro-cid-j7pv25f6]::-webkit-scrollbar{width:0}.history[data-astro-cid-j7pv25f6]{overflow-x:auto;scrollbar-width:none;scroll-behavior:smooth}.history[data-astro-cid-j7pv25f6]::-webkit-scrollbar{width:0}.spinning[data-astro-cid-j7pv25f6]{animation:5s ease-in-out infinite running spinning}@keyframes spinning{0%{rotate:0deg}60%{rotate:360deg}to{rotate:360deg}}\n"}],"routeData":{"route":"/","isIndex":true,"type":"page","pattern":"^\\/$","segments":[],"params":[],"component":"src/pages/index.astro","pathname":"/","prerender":false,"fallbackRoutes":[],"_meta":{"trailingSlash":"ignore"}}},{"file":"","links":[],"scripts":[],"styles":[],"routeData":{"type":"redirect","isIndex":false,"route":"/p/[...slug]","pattern":"^\\/p(?:\\/(.*?))?\\/?$","segments":[[{"content":"p","dynamic":false,"spread":false}],[{"content":"...slug","dynamic":true,"spread":true}]],"params":["...slug"],"component":"/p/[...slug]","prerender":false,"redirect":"/posts/[...slug]","fallbackRoutes":[],"_meta":{"trailingSlash":"ignore"}}}],"site":"https://smartsheep.studio","base":"/","trailingSlash":"ignore","compressHTML":true,"componentMetadata":[["/Users/littlesheep/Documents/Projects/Capital/src/pages/categories/[slug].astro",{"propagation":"in-tree","containsHead":true}],["/Users/littlesheep/Documents/Projects/Capital/src/pages/events/index.astro",{"propagation":"in-tree","containsHead":true}],["/Users/littlesheep/Documents/Projects/Capital/src/pages/posts/[slug].astro",{"propagation":"in-tree","containsHead":true}],["/Users/littlesheep/Documents/Projects/Capital/src/pages/posts/index.astro",{"propagation":"in-tree","containsHead":true}],["/Users/littlesheep/Documents/Projects/Capital/src/pages/tags/[slug].astro",{"propagation":"in-tree","containsHead":true}],["/Users/littlesheep/Documents/Projects/Capital/src/pages/index.astro",{"propagation":"in-tree","containsHead":true}],["/Users/littlesheep/Documents/Projects/Capital/src/layouts/RootLayout.astro",{"propagation":"in-tree","containsHead":false}],["/Users/littlesheep/Documents/Projects/Capital/src/layouts/PageLayout.astro",{"propagation":"in-tree","containsHead":false}],["\u0000@astro-page:src/pages/categories/[slug]@_@astro",{"propagation":"in-tree","containsHead":false}],["\u0000@astrojs-ssr-virtual-entry",{"propagation":"in-tree","containsHead":false}],["\u0000@astro-page:src/pages/events/index@_@astro",{"propagation":"in-tree","containsHead":false}],["\u0000@astro-page:src/pages/posts/[slug]@_@astro",{"propagation":"in-tree","containsHead":false}],["\u0000@astro-page:src/pages/posts/index@_@astro",{"propagation":"in-tree","containsHead":false}],["\u0000@astro-page:src/pages/tags/[slug]@_@astro",{"propagation":"in-tree","containsHead":false}],["\u0000@astro-page:src/pages/index@_@astro",{"propagation":"in-tree","containsHead":false}]],"renderers":[],"clientDirectives":[["idle","(()=>{var i=t=>{let e=async()=>{await(await t())()};\"requestIdleCallback\"in window?window.requestIdleCallback(e):setTimeout(e,200)};(self.Astro||(self.Astro={})).idle=i;window.dispatchEvent(new Event(\"astro:idle\"));})();"],["load","(()=>{var e=async t=>{await(await t())()};(self.Astro||(self.Astro={})).load=e;window.dispatchEvent(new Event(\"astro:load\"));})();"],["media","(()=>{var s=(i,t)=>{let a=async()=>{await(await i())()};if(t.value){let e=matchMedia(t.value);e.matches?a():e.addEventListener(\"change\",a,{once:!0})}};(self.Astro||(self.Astro={})).media=s;window.dispatchEvent(new Event(\"astro:media\"));})();"],["only","(()=>{var e=async t=>{await(await t())()};(self.Astro||(self.Astro={})).only=e;window.dispatchEvent(new Event(\"astro:only\"));})();"],["visible","(()=>{var l=(s,i,o)=>{let r=async()=>{await(await s())()},t=typeof i.value==\"object\"?i.value:void 0,c={rootMargin:t==null?void 0:t.rootMargin},n=new IntersectionObserver(e=>{for(let a of e)if(a.isIntersecting){n.disconnect(),r();break}},c);for(let e of o.children)n.observe(e)};(self.Astro||(self.Astro={})).visible=l;window.dispatchEvent(new Event(\"astro:visible\"));})();"]],"entryModules":{"\u0000@astrojs-ssr-virtual-entry":"entry.mjs","\u0000@astro-renderers":"renderers.mjs","\u0000empty-middleware":"_empty-middleware.mjs","/node_modules/astro/dist/assets/endpoint/node.js":"chunks/pages/node_hIg2I-Kh.mjs","\u0000@astrojs-manifest":"manifest_irk0fM_a.mjs","/Users/littlesheep/Documents/Projects/Capital/node_modules/@astrojs/react/vnode-children.js":"chunks/vnode-children_3wEZly-Z.mjs","\u0000@astro-page:node_modules/astro/dist/assets/endpoint/node@_@js":"chunks/node_0Fr8CwHA.mjs","\u0000@astro-page:src/pages/events/index@_@astro":"chunks/index_6C3b8yBv.mjs","\u0000@astro-page:src/pages/posts/index@_@astro":"chunks/index_Ij1Dwoh1.mjs","\u0000@astro-page:src/pages/categories/[slug]@_@astro":"chunks/_slug__EgGcJ0nJ.mjs","\u0000@astro-page:src/pages/posts/[slug]@_@astro":"chunks/_slug__3BAY271A.mjs","\u0000@astro-page:src/pages/tags/[slug]@_@astro":"chunks/_slug__wxGnmrgA.mjs","\u0000@astro-page:src/pages/index@_@astro":"chunks/index_5_WFUSxR.mjs","/astro/hoisted.js?q=0":"_astro/hoisted.l-JsOPk0.js","@astrojs/react/client.js":"_astro/client.olTvLX7Y.js","/Users/littlesheep/Documents/Projects/Capital/src/components/posts/Media":"_astro/Media.7FWSwaPB.js","astro:scripts/before-hydration.js":""},"assets":["/_astro/ibm-plex-serif-v19-latin-200italic.pJK4yaaG.woff2","/_astro/ibm-plex-serif-v19-latin-200.GFXE_YJc.woff2","/_astro/ibm-plex-serif-v19-latin-100.6qNbweSL.woff2","/_astro/ibm-plex-serif-v19-latin-100italic.E22nrI7z.woff2","/_astro/ibm-plex-serif-v19-latin-300.RVbRgkxX.woff2","/_astro/ibm-plex-serif-v19-latin-300italic.ZdSVgmcR.woff2","/_astro/ibm-plex-serif-v19-latin-regular.HRmMD3sQ.woff2","/_astro/ibm-plex-serif-v19-latin-500.xAA_w-Ac.woff2","/_astro/ibm-plex-serif-v19-latin-italic.MiJiQVsi.woff2","/_astro/ibm-plex-serif-v19-latin-500italic.Unq84pJ7.woff2","/_astro/ibm-plex-serif-v19-latin-600italic.vDhUog1q.woff2","/_astro/ibm-plex-serif-v19-latin-600.cuuqzllG.woff2","/_astro/ibm-plex-serif-v19-latin-700italic.QM1RA0vx.woff2","/_astro/ibm-plex-serif-v19-latin-700.yX9JjmCp.woff2","/_astro/noto-serif-sc-v22-chinese-simplified-200.g4OBZhIi.woff2","/_astro/noto-serif-sc-v22-chinese-simplified-300.yFtdUYoh.woff2","/_astro/noto-serif-sc-v22-chinese-simplified-regular.9muiKgFz.woff2","/_astro/noto-serif-sc-v22-chinese-simplified-600.4n6uFOXj.woff2","/_astro/noto-serif-sc-v22-chinese-simplified-500.exkAspFQ.woff2","/_astro/noto-serif-sc-v22-chinese-simplified-900.ERSRy_0V.woff2","/_astro/noto-serif-sc-v22-chinese-simplified-700.HyiB9Pzv.woff2","/_astro/_slug_.bcjV8AoT.css","/_astro/_slug_.yOjdTrIk.css","/favicon.svg","/_astro/Media.7FWSwaPB.js","/_astro/Media.Co8_pG1j.css","/_astro/client.olTvLX7Y.js","/_astro/hoisted.l-JsOPk0.js","/_astro/index.LFf77hJu.js","/admin/index.html","/media/nicolas-saintot-xkFhOdId7mA-unsplash.jpg"]});

export { AstroIntegrationLogger as A, Logger as L, getEventPrefix as g, levels as l, manifest };

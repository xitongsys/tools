console.log("content.js loaded");

const s = document.createElement("script");
s.src = chrome.runtime.getURL("injected.js");
s.onload = () => s.remove();
document.documentElement.appendChild(s);


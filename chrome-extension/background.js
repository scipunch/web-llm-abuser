chrome.action.onClicked.addListener((tab) => {
  chrome.tabs.create({ url: chrome.runtime.getURL("page.html") });
});

chrome.runtime.onMessage.addListener(function (request, sender, sendResponse) {
    switch (request.action) {
        case "getWebsocketHost":
            chrome.storage.local.get(['websocketHost'], function (result) {
                const host = result.websocketHost || 'ws://localhost:8080';
                sendResponse({ host });
            });
            return true
        case "setWebsocketHost":
            console.debug("setting new websocket host", "request", request)
            const { host } = request
            chrome.storage.local.set({ websocketHost: host }, function () {
                console.log('WebSocket host saved: ' + host);
            });
            return true
        default:
            throw new Error(`Unexpected messate action ${request.action}`)
    }
});


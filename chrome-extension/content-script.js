let ws
let prevMessageLength = 0;
let totalResponses = 0;
let currentIndex = 0;
let inputText = "";
let firstPageLoad = true;

chrome.runtime.sendMessage({ action: "getWebsocketHost" }, function ({ host }) {
  ws = new WebSocket(host);
  ws.onopen = () => {
    console.log("Websocket connection established");
  };

  ws.onmessage = (e) => {
    const msgData = JSON.parse(e.data);
    sendInput(msgData.message);
  };

  ws.onerror = (error) => {
    console.log("Websocket error:", error);
  };
})


function sendInput(message) {
  firstPageLoad = false;
  console.log("Sending input:", message);
  document.querySelector("#prompt-textarea > p").textContent = message
  waitForElement("#composer-submit-button", (el) => el.click())
}

function checkResponses() {
  const responses = document.querySelectorAll(
   ".markdown.prose:not(.result-streaming)"
  );
  console.log("checking responses", "current", responses.length, "total", totalResponses)

  if (responses.length !== totalResponses) {
    const lastMsg = responses[responses.length - 1].textContent
    if (lastMsg.length === 0 || lastMsg.length != prevMessageLength) {
      prevMessageLength = lastMsg.length
      return
    }
    totalResponses = responses.length;
    if (firstPageLoad) {
      return
    }
    ws.send(JSON.stringify({
      role: "model",
      message: lastMsg,
    }));
    console.log("Sent response:", lastMsg);
    prevMessageLength = 0
  }
}

function waitForElement(selector, callback) {
    const interval = setInterval(() => {
        const element = document.querySelector(selector);
        if (element != null) {
            clearInterval(interval);
            callback(element);
        }
    }, 100);
}

setInterval(checkResponses, 100);
console.log("Started interval check");

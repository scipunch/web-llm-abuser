let ws
let prevMessageLength = 0;
let totalResponses = 0;
let currentIndex = 0;
let inputText = "";
let firstPageLoad = true;

chrome.runtime.sendMessage({ action: "getWebsocketHost" }, function ({ host }) {
  ws = new WebSocket(`${host}/responses/ws`);
  ws.onopen = () => {
    console.log("Websocket connection established");
  };

  ws.onmessage = (e) => {
    const msgData = JSON.parse(e.data);
    sendInput(msgData.data);
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
      return null
    }
    totalResponses = responses.length;
    if (firstPageLoad) {
      return null
    }
    return lastMsg
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

setInterval(() => {
  resp = {}
  try {
    const out = checkResponses()
    if (out == null) {
      return
    }
    resp.type = "model-output"
    resp.data = out
  } catch (e) {
    resp.type = "error"
    resp.data = e.toString()
  }

  prevMessageLength = 0

  try {
    ws.send(JSON.stringify(resp));
    console.log("Sent response", resp);
  } catch (e) {
    console.error("could not send response via WebSocket", "cause", e)
  }
}, 100);
console.log("Started interval check");

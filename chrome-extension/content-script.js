chrome.runtime.sendMessage({ action: "getWebsocketHost" }, function ({ host }) {
  const ws = new WebSocket(host);
  ws.onopen = () => {
    console.log("Websocket connection established");
  };

  ws.onmessage = (e) => {
    const msgData = JSON.parse(e.data);
    switch (msgData.type) {
      case "input":
        sendInput(msgData.message);
      default:
        console.error("Unknown message", "data", msgData)
    }
  };

  ws.onerror = (error) => {
    console.log("Websocket error:", error);
  };
})

let totalResponses = 0;
let currentIndex = 0;
let inputText = "";
let firstPageLoad = true;

function sendInput(message) {
  firstPageLoad = false;
  console.log("Sending input:", message);
  document.querySelector("#prompt-textarea > p").textContent = "whats up?"
  document.querySelector("#composer-submit-button").click()
}

function checkResponses() {
  const responses = document.querySelectorAll(
    ".markdown.prose:not(.result-streaming)"
  );
  console.log("checking responses", "current", responses.length, "total", totalResponses)

  if (responses.length !== totalResponses) {
    totalResponses = responses.length;
    if (firstPageLoad) {
      return
    }
    ws.send({
      type: "message",
      message: responses[responses.length - 1].innerHTML,
    });
    console.log("Sent response:", responses[responses.length - 1]);
  }
}

setInterval(checkResponses, 100);
console.log("Started interval check");

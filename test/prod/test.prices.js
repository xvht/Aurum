console.log("Creating WebSocket connection to Aurum Prices Socket...");
const ws = new WebSocket(`wss://aurum.xvh.lol/v1/ws/prices`);

ws.onopen = () => {
  console.log("WebSocket connection opened.");
};

ws.onmessage = (event) => {
  console.log("WebSocket message received. Data length:", event.data.length);
  console.log("-------------------------------");
  try {
    console.log(JSON.parse(event.data));
  } catch (e) {
    console.log(event.data.trim());
  }
  console.log("-------------------------------");
};

ws.onclose = () => {
  console.log("WebSocket connection closed.");
};

ws.onerror = (error) => {
  console.error("WebSocket error occurred:", error);
};

console.log("WebSocket event listeners attached.");

const ws = new WebSocket(`wss://aurum.xvh.lol/v1/ws/prices`);

ws.onopen = () => {
  console.log("connected");
};

ws.onmessage = (event) => {
  console.log(event.data.trim());
};

ws.onclose = () => {
  console.log("disconnected");
};

ws.onerror = (error) => {
  console.log(error);
};

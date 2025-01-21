const ws = new WebSocket(`ws://localhost:8000/v1/ws/prices`);

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

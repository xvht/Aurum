const ws = new WebSocket(`ws://localhost:8000/v1/ws/query`);

ws.onopen = () => {
  console.log("connected");

  const query = {
    queryId: crypto.randomUUID(),
    address: "2GUnfxZavKoPfS9s3VSEjaWDzB3vNf5RojUhprCS1rSx",
    chain: "sol",
  };

  ws.send(JSON.stringify(query));
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

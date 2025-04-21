console.log("Creating WebSocket connection to Aurum Query Socket...");
const ws = new WebSocket(`wss://aurum.xvh.lol/v1/ws/query`);

ws.onopen = () => {
  console.log("WebSocket connection opened.");

  const query = {
    queryId: crypto.randomUUID(),
    address: "2GUnfxZavKoPfS9s3VSEjaWDzB3vNf5RojUhprCS1rSx",
    chain: "sol",
  };
  console.log("Generated query object:", query);

  console.log("Sending query over WebSocket...");
  const serializedQuery = JSON.stringify(query);
  ws.send(serializedQuery);
  console.log("Query sent.");
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

  ws.close();
};

ws.onclose = () => {
  console.log("WebSocket connection closed.");
};

ws.onerror = (error) => {
  console.error("WebSocket error occurred:", error);
};

console.log("WebSocket event listeners attached.");

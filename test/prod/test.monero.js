console.log("Creating WebSocket connection to Bitfinex...");
const ws = new WebSocket("wss://api-pub.bitfinex.com/ws/2");

ws.onopen = () => {
  console.log("WebSocket connection opened.");

  const subscribeMsg = {
    event: "subscribe",
    channel: "ticker",
    symbol: "tXMRUSD",
  };

  ws.send(JSON.stringify(subscribeMsg));
};

ws.onmessage = (event) => {
  /**
   * The parsed ticker data received from the Bitfinex WebSocket stream.
   * @see {@link https://docs.bitfinex.com/reference/ws-public-ticker#stream-data} for the data structure.
   */
  const data = JSON.parse(event.data);

  if (Array.isArray(data) && data[1] !== "hb") {
    console.log("WebSocket message received. Data length:", event.data.length);
    console.log("-------------------------------");

    // LAST_PRICE is at index 6 in the ticker data
    if (Array.isArray(data[1])) {
      const lastPrice = data[1][6];
      console.log("CURRENT XMR PRICE: $" + lastPrice);
    }

    console.log("-------------------------------");
  }
};

ws.onerror = (error) => {
  console.error("WebSocket error occurred:", error);
};

ws.onclose = () => {
  console.log("WebSocket connection closed.");
};

console.log("WebSocket event listeners attached.");

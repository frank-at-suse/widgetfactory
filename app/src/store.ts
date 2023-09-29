import { ref } from "vue";
import type { Ref } from "vue";
import type { HasID, Order, StreamMessage, Widget } from "./types";
import { StreamMessageKind } from "./types";
import axios, { AxiosError } from "axios";

const widgets = ref([] as Widget[]);
const orders = ref([] as Order[]);

var orderSocket: WebSocket;
var widgetSocket: WebSocket;

var host: string;
if (import.meta.env.MODE == "development") {
  host = "localhost:8080";
} else {
  host = window.location.host;
}

var initialized = false;

export const useStore = () => {
  console.log("calling usestore");
  if (!initialized) {
    console.log("initializing");
    initialized = true;

    orderSocket = new WebSocket("ws://" + host + "/ws/orders");
    widgetSocket = new WebSocket("ws://" + host + "/ws/widgets");

    watchSocket(orderSocket, orders);
    watchSocket(widgetSocket, widgets);
  }

  async function createWidget(name: string) {
    const response = await axios.post("http://" + host + "/widget", {
      Name: name,
    });

    if (response.status != 200) {
      console.log("server error: ", response.data);
    }
  }

  async function createOrder(widget: number, quantity: number) {
    const response = await axios.post("http://" + host + "/order", {
      Widget: widget,
      Quantity: quantity,
    });

    if (response.status != 200) {
      console.log("server error: ", response.data);
    }
  }

  async function deleteOrder(ID: number) {
    const response = await axios.delete("http://" + host + "/order", {
      data: { ID: ID },
    });

    if (response.status != 200) {
      console.log("server error: ", response.data);
    }
  }

  async function deleteWidget(ID: number) {
    const response = await axios.delete("http://" + host + "/widget", {
      data: { ID: ID },
    });

    if (response.status != 200) {
      console.log("server error: ", response.data);
    }
  }

  async function sqlQuery(query: string): Promise<string> {
    try {
      const response = await axios.post("http://" + host + "/sql", query);

      return response.data as string;
    } catch (err: any) {
      return "server error: " + err.response.data
    }
  }

  return {
    widgets,
    orders,
    createWidget,
    createOrder,
    deleteOrder,
    deleteWidget,
    sqlQuery,
  };
};

const watchSocket = <T extends HasID>(
  socket: WebSocket,
  collection: Ref<T[]>
) => {
  socket.onmessage = (event) => {
    const sm = JSON.parse(event.data) as StreamMessage;

    switch (sm.Kind) {
      case StreamMessageKind.Create:
        collection.value.push(sm.Object as T);
        break;
      case StreamMessageKind.Delete:
        collection.value.splice(
          collection.value.findIndex((v) => v.ID == (sm.Object as HasID).ID),
          1
        );
        break;
      case StreamMessageKind.Load:
        collection.value = sm.Object as T[];
        break;
      case StreamMessageKind.Error:
        console.log("error from server: ", sm.Object);
        break;
    }
  };
};

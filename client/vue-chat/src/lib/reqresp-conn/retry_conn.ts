import WSSocket from "../socket/wssocket";
import { IRequest, IResponse } from "./reqresp";
import { TSMap } from "typescript-map";
import { HttpStatus, httpStatusTextByCode } from "http-status-ts";
import { chattingEndpoint } from "@/store/endpoints/endpoints";
import { NotifySystem } from "../toasts/notifications";

const sleep = (ms: number) => new Promise((r) => setTimeout(r, ms));

export async function Connect(
  socket: WSSocket,
  timeout: number,

  onConnect: () => void,
  onData: (data: Buffer) => void,
  onClose: () => void,
  onError: (e: Error) => void
): Promise<IResponse> {
  let received = true;
  let response: IResponse = {
    status: "unable to retry connection",
    status_code: HttpStatus.BAD_GATEWAY,

    header: new TSMap<string, string | number>(),
    body: "",
  };

  const request: IRequest = {
    method: "ws",
    proto: "ws",
    header: new TSMap<string, string>([["Content-Type", "application/json"]]),

    url: chattingEndpoint,
    body: "",
  };

  socket.setOnConnect(() => {
    try {
      socket.socketSend(JSON.stringify(request), 0x1, true);
    } catch (e) {
      response.status = httpStatusTextByCode(HttpStatus.BAD_REQUEST);
      response.status_code = HttpStatus.BAD_REQUEST;

      socket.closeSocket();
      received = false;

      console.log("error to send request:", e);
    }
  });

  socket.setOnData((data: Buffer) => {
    try {
      const received_resp = JSON.parse(data.toString());
      response = received_resp;
    } catch (e) {
      socket.closeSocket();
      console.log("error to parse response:", e);
    }

    received = false;
  });

  socket.setOnClose(() => {
    console.log("close");
  });
  socket.setOnError((e: Error) => {
    console.log("error", e);
  });

  setTimeout(() => {
    if (!received) return;

    try {
      socket.closeSocket();
    } catch (e) {
      console.log(e);
    }

    response.status = httpStatusTextByCode(HttpStatus.REQUEST_TIMEOUT);
    response.status_code = HttpStatus.REQUEST_TIMEOUT;
    received = false;
  }, timeout);

  try {
    socket.connectSocket();
  } catch (e) {
    received = false;
  }

  while (received) {
    await sleep(50);
  }

  socket.socket?.removeAllListeners("data");
  socket.socket?.removeAllListeners("connect");
  socket.socket?.removeAllListeners("close");
  socket.socket?.removeAllListeners("error");

  socket.setOnConnect(onConnect);
  socket.setOnData(onData);
  socket.setOnClose(onClose);
  socket.setOnError(onError);

  return response;
}

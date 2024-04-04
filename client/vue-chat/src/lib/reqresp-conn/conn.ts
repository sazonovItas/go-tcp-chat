import { HttpStatus, httpStatusTextByCode } from "http-status-ts";
import { TSMap } from "typescript-map";
import WSSocket from "../socket/wssocket";
import { IRequest, IResponse, successResponse } from "./reqresp";

const sleep = (ms: number) => new Promise((r) => setTimeout(r, ms));
export async function Request(
  host: string,
  port: number,
  timeout: number,
  req: IRequest
): Promise<IResponse> {
  let received = true;
  let resp: IResponse = {
    status: httpStatusTextByCode(HttpStatus.SERVICE_UNAVAILABLE),
    status_code: HttpStatus.SERVICE_UNAVAILABLE,

    header: new TSMap<string, string | number>(),
    body: "",
  };

  const sock = new WSSocket(
    host,
    port,
    () => {
      console.log(
        "connection established",
        sock.getRemoteAddress(),
        sock.getRemotePort()
      );
      try {
        sock.socketSend(JSON.stringify(req), 0x1, true);
      } catch (e) {
        console.log("error to send request to the server:", e);
        received = false;
      }
    },
    (data: Buffer) => {
      try {
        const newResp = JSON.parse(data.toString());
        resp = newResp;
      } catch (e) {
        console.log("error parse response from the server:", e);
      }
      received = false;
      sock.closeSocket();
    }
  );

  setTimeout(() => {
    try {
      sock.closeSocket();
    } catch (e) {
      console.log(e);
    }
    resp.status = httpStatusTextByCode(HttpStatus.REQUEST_TIMEOUT);
    resp.status_code = HttpStatus.REQUEST_TIMEOUT;
    received = false;
  }, timeout);

  try {
    sock.connectSocket();
  } catch (e) {
    console.log("error connect to the server");
    sock.closeSocket();
    return resp;
  }

  while (received) {
    await sleep(50);
  }

  return resp;
}

export async function RequestWS(
  host: string,
  port: number,
  timeout: number,
  req: IRequest
): Promise<WSSocket> {
  let received = true;

  const sock = new WSSocket(
    host,
    port,
    () => {
      console.log(
        "connection established",
        sock.getRemoteAddress(),
        sock.getRemotePort()
      );
      try {
        sock.socketSend(JSON.stringify(req), 0x1, true);
      } catch (e) {
        console.log("error to send request to the server:", e);
        sock.closeSocket();
        received = false;
      }
    },
    (data: Buffer) => {
      try {
        const newResp = JSON.parse(data.toString());
        if (!successResponse(newResp)) {
          sock.closeSocket();
        }
      } catch (e) {
        console.log("error parse response from the server:", e);
        sock.closeSocket();
      }
      received = false;
    }
  );

  setTimeout(() => {
    try {
      sock.closeSocket();
    } catch (e) {
      console.log(e);
    }
    received = false;
  }, timeout);

  try {
    sock.connectSocket();
  } catch (e) {
    console.log("error connect to the server");
    sock.closeSocket();
  }

  while (received) {
    await sleep(50);
  }

  return sock;
}

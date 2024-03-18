import { HttpStatus, httpStatusTextByCode } from "http-status-ts";
import { TSMap } from "typescript-map";
import WSSocket from "../socket/wssocket";
import { IRequest, IResponse } from "./reqresp";

const sleep = (ms: number) => new Promise((r) => setTimeout(r, ms));
export default async function Request(
  host: string,
  port: number,
  timeout: number,
  req: IRequest
): Promise<IResponse> {
  let received = true;
  let resp = {
    status: httpStatusTextByCode(HttpStatus.NOT_FOUND),
    status_code: HttpStatus.NOT_FOUND,

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
    sock.closeSocket();
    resp.status = httpStatusTextByCode(HttpStatus.BAD_GATEWAY);
    resp.status_code = HttpStatus.BAD_GATEWAY;
    received = false;
  }, timeout);

  try {
    sock.connectSocket();
  } catch (e) {
    console.log("error to connect request to the server:", e);
    received = false;
  }

  while (received) {
    await sleep(50);
  }
  return resp;
}

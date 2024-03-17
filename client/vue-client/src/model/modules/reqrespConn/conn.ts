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
    status: "status bad gateaway",
    status_code: 502,

    header: new TSMap<string, string>(),
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

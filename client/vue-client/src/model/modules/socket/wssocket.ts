import net from "net";
import { Buffer } from "buffer";
import { setTimeout, clearTimeout } from "timers";

const preambule: Buffer = Buffer.from([0x5a, 0xa5, 0x5a, 0xa5]);

export default class WSSocket {
  public socket?: net.Socket;
  public host: string;
  public port: number;
  public connected: boolean;

  private timeout?: NodeJS.Timeout;

  constructor(
    host: string,
    port: number,
    onConnect?: () => void,
    onData?: (data: Buffer) => void,
    onError?: (err: Error) => void
  ) {
    this.port = port;
    this.host = host;
    this.connected = false;
    this.timeout = undefined;

    if (onConnect !== undefined) {
      this.onConnect = onConnect;
    }
    if (onData !== undefined) {
      this.onData = (data: Buffer) => {
        const msg = this.unmarshalFrame(data);
        if (msg === undefined) {
          this.closeSocket();
        }

        if (msg !== null) {
          onData(msg!);
        }
      };
    }
    if (onError !== undefined) {
      this.onError = onError;
      this.connected = false;
    }
  }

  public connectSocket() {
    if (!this.socket) {
      this.socket = new net.Socket();

      this.connected = false;
      this.socket.on("connect", this.onConnect);
      this.socket.on("data", this.onData);
      this.socket.on("error", this.onError);
      this.socket.on("drain", () => {
        this.socket?.resume();
      });
      this.socket.on("close", () => {
        this.connected = false;
        console.log("connection is closed");
      });
    }

    if (this.connected || this.socket?.connecting) {
      return;
    }

    try {
      this.socket?.connect(this.port, this.host);
    } catch (e) {
      console.error(e);
    }
  }

  public socketSend(
    data: string | Buffer,
    payloadType: number,
    needMaskingKey: boolean
  ) {
    if (!this.socket || this.connected == false) {
      throw new Error(
        "socket is not initiated or connected is not established"
      );
    }

    if (!Buffer.isBuffer(data)) {
      data = Buffer.from(data, "utf-8");
    }

    const frame: Buffer = this.marshalFrame(data, payloadType, needMaskingKey);
    if (!this.socket.write(frame)) {
      this.socket.pause();
    }
  }

  public setTimeout(timeout: number, onTimeout: () => void) {
    if (this.timeout) {
      clearTimeout(this.timeout);
    }

    this.timeout = setTimeout(() => {
      onTimeout();
      this.setTimeout(timeout, onTimeout);
    }, timeout);
  }

  public closeSocket() {
    if (!this.socket) {
      throw new Error("socket is not initiated");
    }

    if (this.timeout !== undefined) {
      clearTimeout(this.timeout);
    }
    this.socket?.destroy();
  }

  public getSocket(): net.Socket | undefined {
    return this.socket;
  }

  public getLocalAddress(): string | undefined {
    if (this.socket && typeof this.socket.localAddress === "string") {
      return this.socket.localAddress;
    }
  }

  public getLocalPort(): number | undefined {
    if (this.socket && typeof this.socket.localPort === "number") {
      return this.socket.localPort;
    }
  }

  public getRemoteAddress(): string | undefined {
    if (this.socket && typeof this.socket.remoteAddress === "string") {
      return this.socket.remoteAddress;
    }
  }

  public getRemotePort(): number | undefined {
    if (this.socket && typeof this.socket.remotePort === "number") {
      return this.socket.remotePort;
    }
  }

  private onData: (data: Buffer) => void = (data: Buffer) => {
    const msg = this.unmarshalFrame(data);
    if (msg === undefined) {
      this.closeSocket();
    }
    console.log(msg);
  };
  private onConnect: () => void = () => {
    this.connected = true;
    console.log("connection is established");
  };
  private onError: (err: Error) => void = (err: Error) => {
    console.log(err);
  };

  private marshalFrame(
    data: Buffer,
    payloadType: number,
    needMaskingKey: boolean
  ): Buffer {
    let frame: Buffer = Buffer.from(preambule);

    // add fin, rsvs and opcode
    let b: number = 0x80 | (payloadType & 0x0f);
    frame = Buffer.concat([frame, Buffer.from([b])]);

    b = 0x00;
    if (needMaskingKey) {
      b = 0x80;
    }

    let lengthFields = 0;
    if (data.length < 126) {
      b |= data.length & 0xff;
    } else if (data.length < 65536) {
      b |= 126 & 0xff;
      lengthFields = 2;
    } else {
      b |= 127 & 0xff;
      lengthFields = 8;
    }
    frame = Buffer.concat([frame, Buffer.from([b])]);

    for (let i = 0; i < lengthFields; i++) {
      const shift: number = (lengthFields - 1 - i) * 8;
      b = (data.length >> shift) & 0xff;
      frame = Buffer.concat([frame, Buffer.from([b])]);
    }

    if (needMaskingKey) {
      const maskingKey: Buffer = Buffer.from([0x21, 0x54, 0x08, 0x23]);
      frame = Buffer.concat([frame, maskingKey]);

      const msg: Buffer = Buffer.alloc(data.length);
      for (let i = 0; i < msg.length; i++) {
        msg[i] = data[i] ^ maskingKey[i % 4];
      }
      return Buffer.concat([frame, msg]);
    }

    return Buffer.concat([frame, data]);
  }

  private unmarshalFrame(data: Buffer): Buffer | undefined {
    for (let i = 0; i < 4; i++) {
      if (data.length <= i || data[i] !== preambule[i]) {
        return undefined;
      }
    }

    let pos = 4;
    if (data.length < pos) {
      return undefined;
    }
    let b: number = data[pos];

    if ((b & 0xff) === 8) {
      return undefined;
    }

    pos++;
    if (data.length < pos) {
      return undefined;
    }
    b = data[pos];

    const mask: boolean = ((b >> 7) & 1) === 1;
    const checkPayloadLen: number = b & 0x7f;

    let lengthFields = 0;
    let payloadLen = 0;
    switch (checkPayloadLen) {
      case 126:
        lengthFields = 2;
        break;
      case 127:
        lengthFields = 8;
        break;
      default:
        payloadLen = checkPayloadLen;
    }

    pos++;
    for (let i = pos; i < pos + lengthFields; i++) {
      payloadLen = (payloadLen << 8) | (data[i] & 0xff);
    }
    pos += lengthFields;

    if (mask) {
      const maskingKey = data.subarray(pos, pos + 4);
      pos += 4;
      for (let i = pos; i < pos + payloadLen; i++) {
        data[i] ^= maskingKey[(i - pos) % 4];
      }
    }

    return data.subarray(pos, pos + payloadLen);
  }
}

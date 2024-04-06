import net from "net";
import { Buffer } from "buffer";
import { setTimeout, clearTimeout } from "timers";

const preambule: Buffer = Buffer.from([0x5a, 0xa5, 0x5a, 0xa5]);

export default class WSSocket {
  public socket?: net.Socket;
  public host: string;
  public port: number;
  public connected: boolean;
  public try_connecting: boolean;

  private timeout?: NodeJS.Timeout;
  private buf: Buffer;

  public onData: (data: Buffer) => void = (data: Buffer) => {
    console.log(data.toString());
  };
  public onConnect: () => void = () => {
    this.connected = true;
    this.try_connecting = false;
  };
  public onError: (err: Error) => void = (err: Error) => {
    console.log(err);
  };
  public onClose: () => void = () => {
    this.connected = false;
  };

  public setOnConnect(onConnect: () => void) {
    this.onConnect = onConnect;
    const connectHandler = () => {
      this.connected = true;
      this.try_connecting = false;
      console.log("connection is established");
      onConnect();
    };

    this.socket?.removeAllListeners("connect");
    this.socket?.on("connect", connectHandler);
  }

  public setOnData(onData: (data: Buffer) => void) {
    this.onData = onData;
    const dataHandler = (data: Buffer) => {
      this.buf = Buffer.concat([this.buf, data]);

      const msg = this.unmarshalFrame(this.buf);
      if (msg[0] === undefined) {
        if (msg[1] == 0) {
          return;
        }

        this.closeSocket();
        return;
      }

      this.buf = this.buf.subarray(msg[1]);
      if (msg[0] !== undefined) {
        onData(msg[0]!);
      }
    };

    this.socket?.removeAllListeners("data");
    this.socket?.on("data", dataHandler);
  }

  public setOnClose(onClose: () => void) {
    this.onClose = onClose;
    const closeHandler = () => {
      onClose();
      console.log("connection is closed");
      this.connected = false;
      this.try_connecting = false;
      this.socket?.destroy();
      this.socket = undefined;
    };

    this.socket?.removeAllListeners("close");
    this.socket?.on("close", closeHandler);
  }
  public setOnError(onError: (e: Error) => void) {
    this.onError = onError;
    const errorHandler = (e: Error) => {
      onError(e);
    };

    this.socket?.removeAllListeners("error");
    this.socket?.on("error", errorHandler);
  }

  constructor(
    host: string,
    port: number,
    onConnect?: () => void,
    onData?: (data: Buffer) => void,
    onClose?: () => void,
    onError?: (err: Error) => void
  ) {
    this.port = port;
    this.host = host;
    this.connected = false;
    this.try_connecting = false;
    this.timeout = undefined;
    this.buf = Buffer.alloc(0);

    if (onConnect !== undefined) {
      this.onConnect = onConnect;
    }
    if (onData !== undefined) {
      this.onData = onData;
    }
    if (onError !== undefined) {
      this.onError = onError;
    }
    if (onClose !== undefined) {
      this.onClose = onClose;
    }
  }

  public connectSocket() {
    if (!this.socket) {
      this.socket = new net.Socket();
      this.connected = false;
    }

    this.setOnConnect(this.onConnect);
    this.setOnData(this.onData);
    this.setOnClose(this.onClose);
    this.setOnError(this.onError);

    this.socket.removeAllListeners("drain");
    this.socket.on("drain", () => {
      this.socket?.resume();
    });

    if (this.connected || this.socket?.connecting || this.try_connecting) {
      return;
    }

    try {
      this.socket?.connect(this.port, this.host);
      this.try_connecting = true;
    } catch (e) {
      this.try_connecting = false;
      throw e;
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

  private unmarshalFrame(data: Buffer): [Buffer | undefined, number] {
    for (let i = 0; i < 4; i++) {
      if (data.length <= i || data[i] !== preambule[i]) {
        return [undefined, -1];
      }
    }

    let pos = 4;
    if (data.length < pos) {
      return [undefined, -1];
    }
    let b: number = data[pos];

    if ((b & 0xff) === 8) {
      return [undefined, -1];
    }

    pos++;
    if (data.length < pos) {
      return [undefined, -1];
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

    if (data.length < pos + (mask ? 4 : 0) + payloadLen) {
      return [undefined, 0];
    }

    if (mask) {
      const maskingKey = data.subarray(pos, pos + 4);
      pos += 4;
      for (let i = pos; i < pos + payloadLen; i++) {
        data[i] ^= maskingKey[(i - pos) % 4];
      }
    }

    return [data.subarray(pos, pos + payloadLen), pos + payloadLen];
  }
}

import { HttpStatus, httpStatusTextByCode } from "http-status-ts";
import { TSMap } from "typescript-map";

export const ProtoHTTP = "http";
export const ProtoWS = "ws";

export interface IRequest {
  method: string;
  url: string;
  proto: string;

  header: TSMap<string, string | number>;
  body: string;
}

export interface IResponse {
  status: string;
  status_code: number;

  header: TSMap<string, string | number>;
  body: string;
}

export function successResponse(resp: IResponse): boolean {
  return resp.status_code >= 200 && resp.status_code < 300;
}

export function unauthResponse(resp: IResponse): boolean {
  return resp.status_code == HttpStatus.UNAUTHORIZED;
}

import { TSMap } from "typescript-map";

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

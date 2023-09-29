export type HasID = {
  ID: number
}

export type Order = {
  ID: number
  CreatedAt: string
  UpdatedAt: string
  DeletedAt: string
  Widget: number
  Quantity: number
}

export type Widget = {
  ID: number
  CreatedAt: string
  UpdatedAt: string
  DeletedAt: string
  Name: string
}

export enum StreamMessageKind {
  Load = "load",
  Create = "create",
  Delete = "delete",
  Error = "error"
}

export type StreamMessage = {
  Kind: StreamMessageKind
  Object: HasID | HasID[]
}